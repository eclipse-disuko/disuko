// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package policydecisions

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain/decisions"
)

type PolicyDecisionRequestDto struct {
	SBOMId            string     `json:"sbomId" validate:"required"`
	SBOMName          string     `json:"sbomName" validate:"required"`
	SBOMUploaded      *time.Time `json:"sbomUploaded" validate:"required"`
	ComponentSpdxId   string     `json:"componentSpdxId" validate:"required"`
	ComponentName     string     `json:"componentName" validate:"required"`
	ComponentVersion  string     `json:"componentVersion" validate:"required"`
	LicenseExpression string     `json:"licenseExpression" validate:"required"`
	LicenseId         string     `json:"licenseId" validate:"required"`
	PolicyId          string     `json:"policyId" validate:"required"`
	PolicyEvaluated   string     `json:"policyEvaluated" validate:"required"`
	PolicyDecision    string     `json:"policyDecision" validate:"required"`
	Comment           string     `json:"comment" validate:"lte=80"`
	Creator           string     `json:"creator"`
}

type PolicyDecisionSlimDto struct {
	Created           time.Time `json:"created"`
	ComponentName     string    `json:"componentName"`
	ComponentVersion  string    `json:"componentVersion"`
	LicenseExpression string    `json:"licenseExpression"`
	LicenseId         string    `json:"licenseId"`
	PolicyId          string    `json:"policyId"`
	PolicyEvaluated   string    `json:"policyEvaluated"`
	PolicyDecision    string    `json:"policyDecision"`
	Creator           string    `json:"creator"`
	PreviewMode       bool      `json:"previewMode"`
}

func (pd *PolicyDecision) ToSlimDto() *PolicyDecisionSlimDto {
	if pd == nil {
		return nil
	}
	return &PolicyDecisionSlimDto{
		Created:           pd.Created,
		ComponentName:     pd.ComponentName,
		ComponentVersion:  pd.ComponentVersion,
		LicenseExpression: pd.LicenseExpression,
		LicenseId:         pd.LicenseId,
		PolicyId:          pd.PolicyId,
		PolicyEvaluated:   pd.PolicyEvaluated,
		PolicyDecision:    pd.PolicyDecision,
		Creator:           pd.Creator,
		PreviewMode:       pd.PreviewMode,
	}
}

func ToSlimDtos(pds []*PolicyDecision) []*PolicyDecisionSlimDto {
	result := make([]*PolicyDecisionSlimDto, 0)
	for _, entity := range pds {
		result = append(result, entity.ToSlimDto())
	}
	return result
}

func (pd *PolicyDecision) ToDto() *decisions.DecisionDto {
	return &decisions.DecisionDto{
		Key:               pd.Key,
		Created:           pd.Created,
		Updated:           pd.Updated,
		SBOMId:            pd.SBOMId,
		SBOMName:          pd.SBOMName,
		SBOMUploaded:      pd.SBOMUploaded,
		ComponentSpdxId:   pd.ComponentSpdxId,
		ComponentName:     pd.ComponentName,
		ComponentVersion:  pd.ComponentVersion,
		LicenseExpression: pd.LicenseExpression,
		Comment:           pd.Comment,
		Creator:           pd.Creator,
		Active:            pd.Active,

		Type: decisions.PolicyDecision,

		LicenseMatchedId: pd.LicenseId,
		PolicyId:         pd.PolicyId,
		PolicyEvaluated:  pd.PolicyEvaluated,
		PolicyDecision:   pd.PolicyDecision,
	}
}
