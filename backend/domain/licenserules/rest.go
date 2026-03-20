// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package licenserules

import (
	"time"

	"mercedes-benz.ghe.com/foss/disuko/domain/decisions"
)

type LicenseRuleRequestDto struct {
	SBOMId              string     `json:"sbomId" validate:"required"`
	SBOMName            string     `json:"sbomName" validate:"required"`
	SBOMUploaded        *time.Time `json:"sbomUploaded" validate:"required"`
	ComponentSpdxId     string     `json:"componentSpdxId" validate:"required"`
	ComponentName       string     `json:"componentName" validate:"required"`
	ComponentVersion    string     `json:"componentVersion" validate:"required"`
	LicenseExpression   string     `json:"licenseExpression" validate:"required"`
	LicenseDecisionId   string     `json:"licenseDecisionId" validate:"required"`
	LicenseDecisionName string     `json:"licenseDecisionName" validate:"required"`
	Comment             string     `json:"comment" validate:"lte=80"`
	Creator             string     `json:"creator"`
}

type LicenseRuleSlimDto struct {
	Created             time.Time `json:"created"`
	ComponentName       string    `json:"componentName"`
	LicenseExpression   string    `json:"licenseExpression"`
	LicenseDecisionId   string    `json:"licenseDecisionId"`
	LicenseDecisionName string    `json:"licenseDecisionName"`
	Creator             string    `json:"creator"`
	PreviewMode         bool      `json:"previewMode"`
}

func (lr *LicenseRule) ToSlimDto() *LicenseRuleSlimDto {
	if lr == nil {
		return nil
	}
	return &LicenseRuleSlimDto{
		Created:             lr.Created,
		ComponentName:       lr.ComponentName,
		LicenseExpression:   lr.LicenseExpression,
		LicenseDecisionId:   lr.LicenseDecisionId,
		LicenseDecisionName: lr.LicenseDecisionName,
		Creator:             lr.Creator,
		PreviewMode:         lr.PreviewMode,
	}
}

func (lr *LicenseRule) ToDto() *decisions.DecisionDto {
	return &decisions.DecisionDto{
		Key:               lr.Key,
		Created:           lr.Created,
		Updated:           lr.Updated,
		SBOMId:            lr.SBOMId,
		SBOMName:          lr.SBOMName,
		SBOMUploaded:      lr.SBOMUploaded,
		ComponentSpdxId:   lr.ComponentSpdxId,
		ComponentName:     lr.ComponentName,
		ComponentVersion:  lr.ComponentVersion,
		LicenseExpression: lr.LicenseExpression,
		Comment:           lr.Comment,
		Creator:           lr.Creator,
		Active:            lr.Active,

		Type: decisions.LicenseDecision,

		LicenseDecisionId:   lr.LicenseDecisionId,
		LicenseDecisionName: lr.LicenseDecisionName,
	}
}
