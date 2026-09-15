// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package export

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/eclipse-disuko/disuko/domain/license"
)

func TestExportLicenseKnowledgeBasePreservesCalculatedPolicyRule(t *testing.T) {
	expected := &license.PolicyRules{
		Name:        "calculated-rule",
		Calculated:  true,
		CalculatedConfig: license.CalculatedPolicyConfig{
			BucketDefinition: &license.BucketDefinition{
				DeniedClassifications:  []string{"denied"},
				WarnedClassifications:  []string{"warned"},
				AllowedClassifications: []string{"allowed"},
			},
			LicenseScope: license.CalculatedPolicyScope{
				IsLicenseChart: []bool{true},
				ApprovalState:  []license.ApprovalStatus{license.Approved},
				Family:         []license.FamilyOfLicense{license.Permissive},
				LicenseType:    []license.TypeOfLicenses{license.OpenSource},
				Source:         []license.Source{license.PublicLicenseDb},
			},
		},
	}
	exported := &ExportLicenseKnowledgeBaseDto{
		PolicyRules: []*license.PolicyRules{expected},
	}

	data, err := json.Marshal(exported)
	if err != nil {
		t.Fatalf("marshal export payload: %v", err)
	}

	imported := &ExportLicenseKnowledgeBaseDto{}
	if err := json.Unmarshal(data, imported); err != nil {
		t.Fatalf("unmarshal export payload: %v", err)
	}

	if len(imported.PolicyRules) != 1 {
		t.Fatalf("expected one imported policy rule, got %d", len(imported.PolicyRules))
	}
	if !reflect.DeepEqual(expected.CalculatedConfig, imported.PolicyRules[0].CalculatedConfig) {
		t.Fatalf("calculated policy config changed during export/import: expected %#v, got %#v", expected.CalculatedConfig, imported.PolicyRules[0].CalculatedConfig)
	}
	if expected.Calculated != imported.PolicyRules[0].Calculated {
		t.Fatalf("calculated flag changed during export/import: expected %t, got %t", expected.Calculated, imported.PolicyRules[0].Calculated)
	}
}

func TestExportLicenseKnowledgeBasePreservesMultiplePolicyRules(t *testing.T) {
	exported := &ExportLicenseKnowledgeBaseDto{
		PolicyRules: []*license.PolicyRules{
			{
				Name:        "calculated-rule",
				Description: "calculated description",
				Calculated:  true,
				CalculatedConfig: license.CalculatedPolicyConfig{
					LicenseScope: license.CalculatedPolicyScope{
						ApprovalState: []license.ApprovalStatus{license.Approved},
					},
				},
			},
			{
				Name:        "standard-rule",
				Description: "standard description",
				Calculated:  false,
			},
		},
	}

	data, err := json.Marshal(exported)
	if err != nil {
		t.Fatalf("marshal export payload: %v", err)
	}

	imported := &ExportLicenseKnowledgeBaseDto{}
	if err := json.Unmarshal(data, imported); err != nil {
		t.Fatalf("unmarshal export payload: %v", err)
	}

	if !reflect.DeepEqual(exported.PolicyRules, imported.PolicyRules) {
		t.Fatalf("policy rules changed during export/import: expected %#v, got %#v", exported.PolicyRules, imported.PolicyRules)
	}
}

func TestExportLicenseKnowledgeBasePreservesCalculatedPolicyRuleWithoutBucketDefinition(t *testing.T) {
	expected := &license.PolicyRules{
		Name:       "scope-only-calculated-rule",
		Calculated: true,
		CalculatedConfig: license.CalculatedPolicyConfig{
			LicenseScope: license.CalculatedPolicyScope{
				IsLicenseChart: []bool{false},
				Source:         []license.Source{license.CUSTOM},
			},
		},
	}
	exported := &ExportLicenseKnowledgeBaseDto{
		PolicyRules: []*license.PolicyRules{expected},
	}

	data, err := json.Marshal(exported)
	if err != nil {
		t.Fatalf("marshal export payload: %v", err)
	}

	imported := &ExportLicenseKnowledgeBaseDto{}
	if err := json.Unmarshal(data, imported); err != nil {
		t.Fatalf("unmarshal export payload: %v", err)
	}

	if !reflect.DeepEqual(expected.CalculatedConfig, imported.PolicyRules[0].CalculatedConfig) {
		t.Fatalf("calculated policy config changed during export/import: expected %#v, got %#v", expected.CalculatedConfig, imported.PolicyRules[0].CalculatedConfig)
	}
}

func TestImportLicenseKnowledgeBaseDecodesCalculatedPolicyRule(t *testing.T) {
	data := []byte(`{
		"PolicyRules": [{
			"Name": "imported-calculated-rule",
			"Calculated": true,
			"CalculatedConfig": {
				"bucketDefinition": {
					"deniedClassifications": ["denied"],
					"warnedClassifications": ["warned"],
					"allowedClassifications": ["allowed"]
				},
				"licenseScope": {
					"isLicenseChart": [true],
					"approvalState": ["approved"],
					"family": ["permissive"],
					"licenseType": ["open source"],
					"source": ["spdx"]
				}
			}
		}]
	}`)

	imported := &ExportLicenseKnowledgeBaseDto{}
	if err := json.Unmarshal(data, imported); err != nil {
		t.Fatalf("unmarshal import payload: %v", err)
	}

	if len(imported.PolicyRules) != 1 {
		t.Fatalf("expected one imported policy rule, got %d", len(imported.PolicyRules))
	}
	rule := imported.PolicyRules[0]
	if rule.Name != "imported-calculated-rule" || !rule.Calculated {
		t.Fatalf("calculated policy rule metadata was not imported: %#v", rule)
	}
	if rule.CalculatedConfig.BucketDefinition == nil {
		t.Fatal("calculated policy bucket definition was not imported")
	}
	if !reflect.DeepEqual(rule.CalculatedConfig.BucketDefinition.AllowedClassifications, []string{"allowed"}) {
		t.Fatalf("allowed classifications were not imported: %#v", rule.CalculatedConfig.BucketDefinition.AllowedClassifications)
	}
	if !reflect.DeepEqual(rule.CalculatedConfig.LicenseScope.ApprovalState, []license.ApprovalStatus{license.Approved}) {
		t.Fatalf("approval state was not imported: %#v", rule.CalculatedConfig.LicenseScope.ApprovalState)
	}
	if !reflect.DeepEqual(rule.CalculatedConfig.LicenseScope.Family, []license.FamilyOfLicense{license.Permissive}) {
		t.Fatalf("license family was not imported: %#v", rule.CalculatedConfig.LicenseScope.Family)
	}
	if !reflect.DeepEqual(rule.CalculatedConfig.LicenseScope.LicenseType, []license.TypeOfLicenses{license.OpenSource}) {
		t.Fatalf("license type was not imported: %#v", rule.CalculatedConfig.LicenseScope.LicenseType)
	}
	if !reflect.DeepEqual(rule.CalculatedConfig.LicenseScope.Source, []license.Source{license.PublicLicenseDb}) {
		t.Fatalf("license source was not imported: %#v", rule.CalculatedConfig.LicenseScope.Source)
	}
}
