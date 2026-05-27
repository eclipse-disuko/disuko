// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package policyruleclassification

import "github.com/eclipse-disuko/disuko/domain"

type PolicyRuleClassificationDto struct {
	domain.BaseDto
	Name  string                           `json:"name"`
	Rules map[ClassificationKey]RuleStatus `json:"rules"`
}

type PolicyRuleClassificationRequestDto struct {
	Name  string                           `json:"name" validate:"required,gte=1,lte=200"`
	Rules map[ClassificationKey]RuleStatus `json:"rules" validate:"omitempty,dive,required,oneof=allowed warned denied forbidden"`
}

type ClassificationInfo struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type MatrixResponseDto struct {
	Classifications []ClassificationInfo            `json:"classifications"`
	UseCases        []PolicyRuleClassificationDto   `json:"useCases"`
}

func (p *PolicyRuleClassification) ToDto() PolicyRuleClassificationDto {
	res := PolicyRuleClassificationDto{
		Name:  p.Name,
		Rules: p.Rules,
	}
	domain.SetBaseValues(p, &res)
	return res
}

func (d *PolicyRuleClassificationRequestDto) ToEntity() PolicyRuleClassification {
	return PolicyRuleClassification{
		Name:  d.Name,
		Rules: d.Rules,
	}
}
