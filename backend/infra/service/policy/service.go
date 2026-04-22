// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package policy

import (
	"sort"

	"github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/domain/project"
	license2 "github.com/eclipse-disuko/disuko/infra/repository/license"
	"github.com/eclipse-disuko/disuko/infra/repository/policyrules"
	"github.com/eclipse-disuko/disuko/logy"
)

type Service struct {
	PolicyRulesRepository policyrules.IPolicyRulesRepository
	LicenseRepository     license2.ILicensesRepository
}

func (policyRulesHandler *Service) CollectPolicyRulesForProject(requestSession *logy.RequestSession, project *project.Project, licLookup map[string]*license.License) []license.PolicyRulePublicResponseDto {
	lists := policyRulesHandler.PolicyRulesRepository.FindPolicyRulesForLabel(requestSession, project.PolicyLabels)
	responseData := make([]license.PolicyRulePublicResponseDto, 0)
	for _, policyRule := range lists {
		responseData = policyRulesHandler.handlePolicyRulesGetForPublicAddRule(requestSession, policyRule.ComponentsAllow, policyRule, license.ALLOW, responseData, licLookup)
		responseData = policyRulesHandler.handlePolicyRulesGetForPublicAddRule(requestSession, policyRule.ComponentsWarn, policyRule, license.WARN, responseData, licLookup)
		responseData = policyRulesHandler.handlePolicyRulesGetForPublicAddRule(requestSession, policyRule.ComponentsDeny, policyRule, license.DENY, responseData, licLookup)
	}
	return responseData
}

func (policyRulesHandler *Service) handlePolicyRulesGetForPublicAddRule(requestSession *logy.RequestSession,
	components []string, policyRule *license.PolicyRules, listType license.ListType, responseData []license.PolicyRulePublicResponseDto,
	licLookup map[string]*license.License,
) []license.PolicyRulePublicResponseDto {

	licArray := make([]license.PolicyRuleLicensePublicResponse, 0)
	for _, licenseId := range components {
		licenseDto := license.PolicyRuleLicensePublicResponse{
			Identifier: licenseId,
		}
		var licenseEntry *license.License
		if licLookup == nil {
			licenseEntry = policyRulesHandler.LicenseRepository.FindById(requestSession, licenseId)
		} else {
			var ok bool
			licenseEntry, ok = licLookup[licenseId]
			if !ok {
				licenseEntry = policyRulesHandler.LicenseRepository.FindById(requestSession, licenseId)
				licLookup[licenseId] = licenseEntry
			}
		}
		if licenseEntry != nil {
			licenseDto.Key = licenseEntry.Key
			licenseDto.Name = licenseEntry.Name
			licenseDto.Aliases = licenseEntry.Aliases
		}
		licArray = append(licArray, licenseDto)
	}
	newItem := license.PolicyRulePublicResponseDto{
		Key:         policyRule.Key,
		Name:        policyRule.Name,
		Description: policyRule.Description,
		Licenses:    licArray,
		Type:        listType,
		Created:     policyRule.Created,
		Updated:     policyRule.Updated,
	}
	return append(responseData, newItem)
}

func (policyRulesHandler *Service) CalculatePolicyRuleComponents(requestSession *logy.RequestSession, config license.CalculatedPolicyConfig) ([]string, []string, []string) {
	allLicenses := policyRulesHandler.LicenseRepository.FindAll(requestSession, false)

	if config.BucketDefinition == nil {
		return []string{}, []string{}, []string{}
	}

	deniedClassifications := make(map[string]bool)
	for _, key := range config.BucketDefinition.DeniedClassifications {
		deniedClassifications[key] = true
	}
	warnedClassifications := make(map[string]bool)
	for _, key := range config.BucketDefinition.WarnedClassifications {
		warnedClassifications[key] = true
	}
	allowedClassifications := make(map[string]bool)
	for _, key := range config.BucketDefinition.AllowedClassifications {
		allowedClassifications[key] = true
	}

	allowMap := make(map[string]bool)
	warnMap := make(map[string]bool)
	denyMap := make(map[string]bool)

	for _, currentLicense := range allLicenses {
		if currentLicense == nil || currentLicense.LicenseId == "" {
			continue
		}
		if !matchesCalculatedScopeFilters(currentLicense, config) {
			continue
		}

		hasDenied := false
		hasWarned := false
		hasAllowed := false
		for _, classificationKey := range currentLicense.Meta.ObligationsKeyList {
			if deniedClassifications[classificationKey] {
				hasDenied = true
				break
			}
			if warnedClassifications[classificationKey] {
				hasWarned = true
			}
			if allowedClassifications[classificationKey] {
				hasAllowed = true
			}
		}

		if hasDenied {
			denyMap[currentLicense.LicenseId] = true
		} else if hasWarned {
			warnMap[currentLicense.LicenseId] = true
		} else if hasAllowed {
			allowMap[currentLicense.LicenseId] = true
		}
	}

	componentsAllow := make([]string, 0, len(allowMap))
	for licenseID := range allowMap {
		componentsAllow = append(componentsAllow, licenseID)
	}
	componentsWarn := make([]string, 0, len(warnMap))
	for licenseID := range warnMap {
		componentsWarn = append(componentsWarn, licenseID)
	}
	componentsDeny := make([]string, 0, len(denyMap))
	for licenseID := range denyMap {
		componentsDeny = append(componentsDeny, licenseID)
	}

	sort.Strings(componentsAllow)
	sort.Strings(componentsWarn)
	sort.Strings(componentsDeny)

	return componentsAllow, componentsWarn, componentsDeny
}

func matchesCalculatedScopeFilters(currentLicense *license.License, config license.CalculatedPolicyConfig) bool {
	scope := config.LicenseScope
	if len(scope.IsLicenseChart) > 0 && !matchesBoolFilter(currentLicense.Meta.IsLicenseChart, scope.IsLicenseChart) {
		return false
	}
	if len(scope.ApprovalState) > 0 && !matchesApprovalStateFilter(currentLicense.Meta.ApprovalState, scope.ApprovalState) {
		return false
	}
	if len(scope.Family) > 0 && !matchesFamilyFilter(currentLicense.Meta.Family, scope.Family) {
		return false
	}
	if len(scope.LicenseType) > 0 && !matchesLicenseTypeFilter(currentLicense.Meta.LicenseType, scope.LicenseType) {
		return false
	}
	if len(scope.Source) > 0 && !matchesSourceFilter(currentLicense.Source, scope.Source) {
		return false
	}
	return true
}

func matchesBoolFilter(actualValue bool, filterValues []bool) bool {
	for _, filterValue := range filterValues {
		if actualValue == filterValue {
			return true
		}
	}
	return false
}

func matchesApprovalStateFilter(actualValue license.ApprovalStatus, filterValues []license.ApprovalStatus) bool {
	for _, filterValue := range filterValues {
		if actualValue == filterValue {
			return true
		}
	}
	return false
}

func matchesFamilyFilter(actualValue license.FamilyOfLicense, filterValues []license.FamilyOfLicense) bool {
	for _, filterValue := range filterValues {
		if actualValue == filterValue {
			return true
		}
	}
	return false
}

func matchesLicenseTypeFilter(actualValue license.TypeOfLicenses, filterValues []license.TypeOfLicenses) bool {
	for _, filterValue := range filterValues {
		if actualValue == filterValue {
			return true
		}
	}
	return false
}

func matchesSourceFilter(actualValue license.Source, filterValues []license.Source) bool {
	for _, filterValue := range filterValues {
		if actualValue == filterValue {
			return true
		}
	}
	return false
}

