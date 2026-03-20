// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package components

import (
	"github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/domain/licenserules"
	"github.com/eclipse-disuko/disuko/domain/policydecisions"
	"github.com/eclipse-disuko/disuko/helper/message"
)

type PolicyRuleStatusDto struct {
	Key                        string           `json:"key"`
	Name                       string           `json:"name"`
	LicenseMatched             string           `json:"licenseMatched"`
	Type                       license.ListType `json:"type"`
	Used                       bool             `json:"used"`
	Description                string           `json:"description"`
	IsDecisionMade             bool             `json:"isDecisionMade"`
	CanMakeWarnedDecision      bool             `json:"canMakeWarnedDecision"`
	CanMakeDeniedDecision      bool             `json:"canMakeDeniedDecision"`
	DeniedDecisionDeniedReason string           `json:"deniedDecisionDeniedReason"`
}
type UnmatchedLicenseDto struct {
	OrigName       string `json:"orig"`
	ReferencedName string `json:"referenced"`
	Known          bool   `json:"known"`
}

type ComponentInfoDto struct {
	SpdxId             string                 `json:"spdxId"`
	Name               string                 `json:"name"`
	Version            string                 `json:"version"`
	LicenseEffective   string                 `json:"licenseEffective"`
	License            string                 `json:"license"`
	LicenseDeclared    string                 `json:"licenseDeclared"`
	LicenseComments    string                 `json:"licenseComments"`
	WorstFamily        string                 `json:"worstFamily"`
	CopyrightText      string                 `json:"copyrightText"`
	Description        string                 `json:"description"`
	DownloadLocation   string                 `json:"downloadLocation"`
	Type               ComponentType          `json:"type"`
	Modified           bool                   `json:"modified"`
	Questioned         bool                   `json:"questioned"`
	Unasserted         bool                   `json:"unasserted"`
	PolicyRuleStatus   []*PolicyRuleStatusDto `json:"policyRuleStatus"`
	UnmatchedLicenses  []*UnmatchedLicenseDto `json:"unmatchedLicenses"`
	LicenseApplied     LicenseAppliedType     `json:"licenseApplied"`
	PURL               string                 `json:"purl"`
	PrStatus           string                 `json:"prStatus"`
	UsedPolicyRule     string                 `json:"usedPolicyRule"`
	CanChooseLicense   bool                   `json:"canChooseLicense"`
	ChoiceDeniedReason string                 `json:"choiceDeniedReason"`

	LicenseRuleApplied *licenserules.LicenseRuleSlimDto `json:"licenseRuleApplied"`

	PolicyDecisionsApplied     []*policydecisions.PolicyDecisionSlimDto `json:"policyDecisionsApplied"`
	PolicyDecisionDeniedReason string                                   `json:"policyDecisionDeniedReason"`
}

type ComponentInfoSlimDto struct {
	SpdxId            string `json:"spdxId"`
	Name              string `json:"name"`
	Version           string `json:"version"`
	LicenseExpression string `json:"licenseExpression"`
}

func (ci ComponentInfo) ToComponentInfoSlimDto() *ComponentInfoSlimDto {
	return &ComponentInfoSlimDto{
		SpdxId:            ci.SpdxId,
		Name:              ci.Name,
		Version:           ci.Version,
		LicenseExpression: ci.EffectiveLicensesString(),
	}
}

type ComponentsInfoResponse struct {
	ComponentInfo                  []ComponentInfoDto `json:"componentInfo"`
	ComponentStats                 ComponentStats     `json:"componentStats"`
	BulkPolicyDecisionDeniedReason string             `json:"bulkPolicyDecisionDeniedReason"`
}

func (entity *ComponentResult) ToComponentInfoDto(isResponsible bool, policyDecisionDeniedReason string, isAllowDeniedPolicyDecision bool) *ComponentInfoDto {
	status, rule := entity.GetUsedPolicyRule()

	var (
		deniedReason string
		canChoose    = entity.Component.LicenseRuleApplied == nil && entity.Component.GetLicensesEffective().Op == OR
	)
	if canChoose {
		if !isResponsible {
			deniedReason = message.ChoiceDeniedResp
		} else if len(entity.Component.GetLicensesEffective().List) > 4 {
			deniedReason = message.ChoiceDeniedMassive
		}
	}

	if entity.Component.Version == "" {
		if policyDecisionDeniedReason == "" {
			policyDecisionDeniedReason = message.DecisionDeniedComponentVersionNotSet
		}
		if deniedReason == "" {
			deniedReason = message.DecisionDeniedComponentVersionNotSet
		}
	}

	return &ComponentInfoDto{
		SpdxId:                     entity.Component.SpdxId,
		Name:                       entity.Component.Name,
		Version:                    entity.Component.Version,
		LicenseEffective:           entity.Component.EffectiveLicensesString(),
		License:                    entity.Component.License,
		LicenseDeclared:            entity.Component.LicenseDeclared,
		LicenseComments:            entity.Component.LicenseComments,
		WorstFamily:                string(entity.Component.WorstFamily()),
		CopyrightText:              entity.Component.CopyrightText,
		Description:                entity.Component.Description,
		DownloadLocation:           entity.Component.DownloadLocation,
		Type:                       entity.Component.Type,
		LicenseApplied:             entity.Component.GetLicenseAppliedType(),
		Modified:                   entity.Component.Modified,
		Questioned:                 entity.Questioned,
		Unasserted:                 entity.Unasserted,
		PolicyRuleStatus:           ToPolicyStatusDto(entity.Status, isAllowDeniedPolicyDecision),
		UnmatchedLicenses:          ToUnmatchedDto(entity.Unmatched),
		PrStatus:                   status,
		UsedPolicyRule:             rule,
		PURL:                       entity.Component.PURL,
		CanChooseLicense:           canChoose,
		ChoiceDeniedReason:         deniedReason,
		LicenseRuleApplied:         entity.Component.LicenseRuleApplied.ToSlimDto(),
		PolicyDecisionsApplied:     policydecisions.ToSlimDtos(entity.Component.PolicyDecisionsApplied),
		PolicyDecisionDeniedReason: policyDecisionDeniedReason,
	}
}

func (entity *EvaluationResult) ToComponentInfoDtos(isResponsible bool, policyDecisionDeniedReason string, isAllowDeniedPolicyDecision bool) []ComponentInfoDto {
	dtos := make([]ComponentInfoDto, 0)
	for _, compRes := range entity.Results {
		dtos = append(dtos, *compRes.ToComponentInfoDto(isResponsible, policyDecisionDeniedReason, isAllowDeniedPolicyDecision))
	}
	return dtos
}

func ToPolicyStatusDto(status []*PolicyRuleStatus, isAllowDeniedPolicyDecision bool) []*PolicyRuleStatusDto {
	dtos := make([]*PolicyRuleStatusDto, 0)
	for _, s := range status {
		dtos = append(dtos, &PolicyRuleStatusDto{
			Key:                        s.Key,
			Name:                       s.Name,
			LicenseMatched:             s.LicenseMatched,
			Type:                       s.Type,
			Used:                       s.Used,
			Description:                s.Description,
			IsDecisionMade:             s.IsDecisionMade,
			CanMakeWarnedDecision:      s.CanMakeWarnedDecision,
			CanMakeDeniedDecision:      s.CanMakeDeniedDecision && isAllowDeniedPolicyDecision,
			DeniedDecisionDeniedReason: s.DeniedDecisionDeniedReason,
		})
	}
	return dtos
}

func ToUnmatchedDto(unmatched []*UnmatchedLicense) []*UnmatchedLicenseDto {
	dtos := make([]*UnmatchedLicenseDto, 0)
	for _, u := range unmatched {
		dtos = append(dtos, &UnmatchedLicenseDto{
			OrigName:       u.OrigName,
			ReferencedName: u.ReferencedName,
			Known:          u.Known,
		})
	}
	return dtos
}
