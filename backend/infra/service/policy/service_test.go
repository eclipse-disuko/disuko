// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package policy

import (
	"testing"

	"github.com/eclipse-disuko/disuko/domain/license"
	licenseRepository "github.com/eclipse-disuko/disuko/infra/repository/license"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/stretchr/testify/assert"
)

type calculatedRuleLicenseRepository struct {
	licenseRepository.ILicensesRepository
	licenses []*license.License
}

func (repository *calculatedRuleLicenseRepository) FindAll(_ *logy.RequestSession, _ bool) []*license.License {
	return repository.licenses
}

func calculatedRuleLicense(id string, obligations ...string) *license.License {
	return &license.License{
		LicenseId: id,
		Meta: license.MetaData{
			ObligationsKeyList: obligations,
		},
	}
}

func TestCalculatePolicyRuleComponents(t *testing.T) {
	licenses := []*license.License{
		calculatedRuleLicense("z-allowed", "allowed"),
		calculatedRuleLicense("a-warned", "warned", "allowed"),
		calculatedRuleLicense("m-denied", "denied", "warned", "allowed"),
		calculatedRuleLicense("ignored", "unclassified"),
		nil,
		calculatedRuleLicense(""),
	}

	service := &Service{
		LicenseRepository: &calculatedRuleLicenseRepository{licenses: licenses},
	}

	allow, warn, deny := service.CalculatePolicyRuleComponents(nil, license.CalculatedPolicyConfig{
		BucketDefinition: &license.BucketDefinition{
			DeniedClassifications:  []string{"denied"},
			WarnedClassifications:  []string{"warned"},
			AllowedClassifications: []string{"allowed"},
		},
	})

	assert.Equal(t, []string{"z-allowed"}, allow)
	assert.Equal(t, []string{"a-warned"}, warn)
	assert.Equal(t, []string{"m-denied"}, deny)
}

func TestCalculatePolicyRuleComponentsAppliesAllScopeFilters(t *testing.T) {
	matching := calculatedRuleLicense("matching", "allowed")
	matching.Meta.IsLicenseChart = true
	matching.Meta.ApprovalState = license.Approved
	matching.Meta.Family = license.Permissive
	matching.Meta.LicenseType = license.OpenSource
	matching.Source = license.PublicLicenseDb

	config := license.CalculatedPolicyConfig{
		BucketDefinition: &license.BucketDefinition{AllowedClassifications: []string{"allowed"}},
		LicenseScope: license.CalculatedPolicyScope{
			IsLicenseChart: []bool{true},
			ApprovalState:  []license.ApprovalStatus{license.Approved},
			Family:         []license.FamilyOfLicense{license.Permissive},
			LicenseType:    []license.TypeOfLicenses{license.OpenSource},
			Source:         []license.Source{license.PublicLicenseDb},
		},
	}

	filterCases := []struct {
		name      string
		candidate func(*license.License) *license.License
	}{
		{
			name: "is license chart",
			candidate: func(source *license.License) *license.License {
				candidate := *source
				candidate.LicenseId = "wrong-chart"
				candidate.Meta.IsLicenseChart = false
				return &candidate
			},
		},
		{
			name: "approval state",
			candidate: func(source *license.License) *license.License {
				candidate := *source
				candidate.LicenseId = "wrong-approval"
				candidate.Meta.ApprovalState = license.Pending
				return &candidate
			},
		},
		{
			name: "family",
			candidate: func(source *license.License) *license.License {
				candidate := *source
				candidate.LicenseId = "wrong-family"
				candidate.Meta.Family = license.WeakCopyleft
				return &candidate
			},
		},
		{
			name: "license type",
			candidate: func(source *license.License) *license.License {
				candidate := *source
				candidate.LicenseId = "wrong-type"
				candidate.Meta.LicenseType = license.Proprietary
				return &candidate
			},
		},
		{
			name: "source",
			candidate: func(source *license.License) *license.License {
				candidate := *source
				candidate.LicenseId = "wrong-source"
				candidate.Source = license.CUSTOM
				return &candidate
			},
		},
	}

	for _, filterCase := range filterCases {
		t.Run(filterCase.name, func(t *testing.T) {
			service := &Service{
				LicenseRepository: &calculatedRuleLicenseRepository{licenses: []*license.License{
					matching,
					filterCase.candidate(matching),
				}},
			}

			allow, warn, deny := service.CalculatePolicyRuleComponents(nil, config)

			assert.Equal(t, []string{"matching"}, allow)
			assert.Empty(t, warn)
			assert.Empty(t, deny)
		})
	}
}

func TestCalculatePolicyRuleComponentsReturnsEmptyListsWithoutBuckets(t *testing.T) {
	service := &Service{
		LicenseRepository: &calculatedRuleLicenseRepository{licenses: []*license.License{
			calculatedRuleLicense("allowed", "allowed"),
		}},
	}

	allow, warn, deny := service.CalculatePolicyRuleComponents(nil, license.CalculatedPolicyConfig{})

	assert.Empty(t, allow)
	assert.Empty(t, warn)
	assert.Empty(t, deny)
}
