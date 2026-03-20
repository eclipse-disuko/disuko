// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package license

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/audit"
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

	Auxiliary      bool
	Active         bool
	ApplyToAll     bool
	Deprecated     bool
	DeprecatedDate time.Time
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
		Status:          status,
		Name:            entity.Name,
		LabelSets:       entity.LabelSets,
		Description:     entity.Description,
		ComponentsAllow: entity.ComponentsAllow,
		ComponentsDeny:  entity.ComponentsDeny,
		ComponentsWarn:  entity.ComponentsWarn,
		Auxiliary:       entity.Auxiliary,
		Deprecated:      entity.Deprecated,
		DeprecatedDate:  entity.DeprecatedDate,
		Active:          entity.Active,
		ApplyToAll:      entity.ApplyToAll,
	}
	domain.SetBaseValues(entity, policyRule)
	return policyRule
}

func ToPolicyRuleDtoList(source []*PolicyRules) []*PolicyRulesDto {
	return domain.MapTo(source, ToDto)
}
