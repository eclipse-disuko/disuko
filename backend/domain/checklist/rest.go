// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package checklist

import (
	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/project"
)

type ItemDto struct {
	domain.BaseDto
	Name               string                    `json:"name" validate:"required,gte=3,lte=80"`
	TriggerType        TriggerType               `json:"triggerType"`
	Classifications    []string                  `json:"classifications,omitempty"`
	PolicyStatus       []PolicyStatusType        `json:"policyStatus,omitempty"`
	ScanRemarks        *project.ScanRemarkStatus `json:"scanRemarks,omitempty"`
	LicensesIds        []string                  `json:"licenseIds,omitempty"`
	ComponentNames     []string                  `json:"componentNames,omitempty"`
	TargetTemplateName string                    `json:"targetTemplateName,omitempty"`
	TargetTemplateKey  string                    `json:"targetTemplateKey,omitempty"`
	PolicyLabels       []string                  `json:"policyLabels,omitempty"`
}

type ChecklistDto struct {
	domain.BaseDto
	Name          string    `json:"name" validate:"required,gte=3,lte=80"`
	NameDE        string    `json:"nameDE" validate:"required,gte=3,lte=80"`
	Description   string    `json:"description" validate:"lte=1000"`
	DescriptionDE string    `json:"descriptionDE" validate:"lte=1000"`
	PolicyLabels  []string  `json:"policyLabels" validate:"gte=1"`
	Items         []ItemDto `json:"items"`
	Active        bool      `json:"active"`
}

type ExecuteRequestDto struct {
	Ids []string `json:"ids" validate:"gte=1"`
}

func (i Item) ToDto() ItemDto {
	result := ItemDto{
		BaseDto: domain.BaseDto{
			Key:     i.Key,
			Updated: i.Updated,
			Created: i.Created,
		},
		Name:               i.Name,
		TriggerType:        i.TriggerType,
		Classifications:    i.Classifications,
		PolicyStatus:       i.PolicyStatus,
		ScanRemarks:        i.ScanRemarks,
		LicensesIds:        i.LicenseIds,
		ComponentNames:     i.ComponentNames,
		TargetTemplateName: i.TargetTemplateName,
		TargetTemplateKey:  i.TargetTemplateKey,
		PolicyLabels:       i.PolicyLabels,
	}
	return result
}

func (c *Checklist) ToDto() ChecklistDto {
	result := ChecklistDto{
		Name:          c.Name,
		NameDE:        c.NameDE,
		Description:   c.Description,
		DescriptionDE: c.DescriptionDE,
		PolicyLabels:  c.PolicyLabels,
		Items:         domain.ToDtos(c.Items),
		Active:        c.Active,
	}
	domain.SetBaseValues(c, &result)
	return result
}

func (c *ChecklistDto) ToEntity() Checklist {
	return Checklist{
		Name:          c.Name,
		NameDE:        c.NameDE,
		Description:   c.Description,
		DescriptionDE: c.DescriptionDE,
		PolicyLabels:  c.PolicyLabels,
		Active:        c.Active,
	}
}

func (i ItemDto) ToEntity() Item {
	return Item{
		ChildEntity: domain.ChildEntity{
			Key:     i.Key,
			Created: i.Created,
			Updated: i.Updated,
		},
		Name:               i.Name,
		TriggerType:        i.TriggerType,
		Classifications:    i.Classifications,
		PolicyStatus:       i.PolicyStatus,
		ScanRemarks:        i.ScanRemarks,
		LicenseIds:         i.LicensesIds,
		ComponentNames:     i.ComponentNames,
		TargetTemplateName: i.TargetTemplateName,
		TargetTemplateKey:  i.TargetTemplateKey,
		PolicyLabels:       i.PolicyLabels,
	}
}
