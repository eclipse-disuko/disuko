// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package license

import (
	"encoding/json"
	"sort"
	"strings"
	"time"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/audit"
	"github.com/eclipse-disuko/disuko/helper/hash"
	"github.com/eclipse-disuko/disuko/logy"
)

type PolicyRules struct {
	domain.RootEntity `bson:"inline"`
	domain.SoftDelete `bson:"inline"`
	audit.Container   `bson:"inline"`
	Name              string
	LabelSets         [][]string
	Description       string

	ComponentsAllow []string
	ComponentsDeny  []string
	ComponentsWarn  []string

	Auxiliary        bool
	Active           bool
	ApplyToAll       bool
	Calculated       bool
	CalculatedConfig CalculatedPolicyConfig
	Deprecated       bool
	DeprecatedDate   time.Time
}

type PolicyRuleHashEntry struct {
	Key             string
	ComponentsAllow []string
	ComponentsDeny  []string
	ComponentsWarn  []string
}

type PolicyRulesList []*PolicyRules

func (r *PolicyRulesList) GenHash(requestSession *logy.RequestSession) string {
	if r == nil {
		return ""
	}

	entries := make([]PolicyRuleHashEntry, 0, len(*r))
	for _, rule := range *r {
		if rule == nil {
			continue
		}
		entries = append(entries, PolicyRuleHashEntry{
			Key:             rule.Key,
			ComponentsAllow: normalizeAndSort(rule.ComponentsAllow),
			ComponentsDeny:  normalizeAndSort(rule.ComponentsDeny),
			ComponentsWarn:  normalizeAndSort(rule.ComponentsWarn),
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	prStr, err := json.Marshal(entries)
	if err != nil {
		logy.Warnf(requestSession, "Error marshalling slice of policy rule hash entries")
		return ""
	}
	return hash.Hash(requestSession, prStr)
}

func normalizeAndSort(values []string) []string {
	res := make([]string, 0, len(values))
	for _, v := range values {
		res = append(res, strings.ToLower(strings.TrimSpace(v)))
	}
	sort.Strings(res)
	return res
}

type BucketDefinition struct {
	DeniedClassifications  []string `json:"deniedClassifications" validate:"dive,gte=1,lte=80"`
	WarnedClassifications  []string `json:"warnedClassifications" validate:"dive,gte=1,lte=80"`
	AllowedClassifications []string `json:"allowedClassifications" validate:"dive,gte=1,lte=80"`
}

type CalculatedPolicyConfig struct {
	BucketDefinition *BucketDefinition     `json:"bucketDefinition"`
	LicenseScope     CalculatedPolicyScope `json:"licenseScope"`
}

type CalculatedPolicyScope struct {
	IsLicenseChart []bool            `json:"isLicenseChart" validate:"omitempty,unique"`
	ApprovalState  []ApprovalStatus  `json:"approvalState" validate:"omitempty,unique,dive,lte=80"`
	Family         []FamilyOfLicense `json:"family" validate:"omitempty,unique,dive,lte=80"`
	LicenseType    []TypeOfLicenses  `json:"licenseType" validate:"omitempty,unique,dive,lte=50"`
	Source         []Source          `json:"source" validate:"omitempty,unique,dive,lte=80"`
}

type PolicyResult int

const (
	ALLOWED PolicyResult = iota
	WARNED
	DENIED
	QUESTIONED
	UNASSERTED
)

type ListType string //	@name	ListType

const (
	ALLOW   ListType = "allow"
	DENY    ListType = "deny"
	WARN    ListType = "warn"
	NOT_SET ListType = "NOT_SET"
)

func (entity *PolicyRules) ToDto() *PolicyRulesDto {
	return ToDto(entity)
}

func ToDto(entity *PolicyRules) *PolicyRulesDto {
	status := StatusActive
	if !entity.Active {
		status = StatusInactive
	}
	if entity.Deprecated {
		status = StatusDeprecated
	}
	policyRule := &PolicyRulesDto{
		Status:           status,
		Name:             entity.Name,
		LabelSets:        entity.LabelSets,
		Description:      entity.Description,
		ComponentsAllow:  entity.ComponentsAllow,
		ComponentsDeny:   entity.ComponentsDeny,
		ComponentsWarn:   entity.ComponentsWarn,
		Auxiliary:        entity.Auxiliary,
		Deprecated:       entity.Deprecated,
		DeprecatedDate:   entity.DeprecatedDate,
		Active:           entity.Active,
		ApplyToAll:       entity.ApplyToAll,
		Calculated:       entity.Calculated,
		CalculatedConfig: entity.CalculatedConfig,
	}
	domain.SetBaseValues(entity, policyRule)
	return policyRule
}

func ToPolicyRuleDtoList(source []*PolicyRules) []*PolicyRulesDto {
	return domain.MapTo(source, ToDto)
}
