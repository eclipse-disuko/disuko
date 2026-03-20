// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package policy

import (
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
