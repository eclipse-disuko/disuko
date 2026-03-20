// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package changeloglist

import (
	"encoding/json"
	"time"

	"mercedes-benz.ghe.com/foss/disuko/domain"
)

const (
	TYPE_POLICY_RULE = "policy_rule_change"
)

type ChangeLog struct {
	domain.RootEntity `bson:"inline"`
	When              time.Time
	Type              string
	RefKey            string
	RefName           string
	Content           string
}

type PolicyRuleChangeLogContent struct {
	LicenseName  string `json:"licenseName"`
	LicenseId    string `json:"licenseId"`
	PolicyStatus string `json:"policyStatus"`
	Change       string `json:"change"`
}
type PolicyRuleChangeLog struct {
	ChangeLog
	PolicyRuleChangeLogContent
}

type ChangeLogDto struct {
	Key     string    `json:"_key"`
	When    time.Time `json:"when"`
	Type    string    `json:"type"`
	RefKey  string    `json:"refKey"`
	RefName string    `json:"refName"`
	Content string    `json:"content"`
}

type ChangeLogGenerationResultDto struct {
	ReqID           string
	ChangeLogsCount int
	Errors          []string
}

func NewPolicyRuleChangeLog() *PolicyRuleChangeLog {
	return &PolicyRuleChangeLog{
		ChangeLog: ChangeLog{
			RootEntity: domain.NewRootEntity(),
			When:       time.Now(),
			Type:       TYPE_POLICY_RULE,
		},
		PolicyRuleChangeLogContent: PolicyRuleChangeLogContent{
			LicenseName:  "",
			LicenseId:    "",
			PolicyStatus: "",
			Change:       "",
		},
	}
}

func (entity *PolicyRuleChangeLog) ToChangeLog() (*ChangeLog, error) {
	content := PolicyRuleChangeLogContent{
		LicenseName:  entity.LicenseName,
		LicenseId:    entity.LicenseId,
		PolicyStatus: entity.PolicyStatus,
		Change:       entity.Change,
	}

	contentJson, err := json.Marshal(content)

	return &ChangeLog{
		RootEntity: domain.NewRootEntity(),
		When:       entity.When,
		Type:       entity.Type,
		RefKey:     entity.RefKey,
		RefName:    entity.RefName,
		Content:    string(contentJson),
	}, err
}

func (entity *ChangeLog) ToDto() *ChangeLogDto {
	return ToDto(entity)
}

func ToDto(entity *ChangeLog) *ChangeLogDto {
	return &ChangeLogDto{
		Key:     entity.Key,
		When:    entity.When,
		Type:    entity.Type,
		RefKey:  entity.RefKey,
		RefName: entity.RefName,
		Content: entity.Content,
	}
}

func ToDtos(source []*ChangeLog) []*ChangeLogDto {
	return domain.MapTo(source, ToDto)
}
