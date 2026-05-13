// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package license

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain"
)

type PolicyRuleDto struct {
	Key         string    `json:"key"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Auxiliary   bool      `json:"auxiliary"`
	Deprecated  bool      `json:"deprecated"`
	Active      bool      `json:"active"`
}

type PolicyRuleRequestDto struct {
	Name             string                 `validate:"required,gte=3,lte=80"`
	Description      string                 `validate:"lte=1500"`
	LabelSets        [][]string             `validate:"dive,dive,gte=3,lte=80"`
	ComponentsAllow  []string               `validate:"dive,gte=1,lte=80"`
	ComponentsDeny   []string               `validate:"dive,gte=1,lte=80"`
	ComponentsWarn   []string               `validate:"dive,gte=1,lte=80"`
	Auxiliary        bool                   `json:"auxiliary"`
	Active           bool                   `json:"active"`
	ApplyToAll       bool                   `json:"applyToAll"`
	Calculated       bool                   `json:"calculated"`
	CalculatedConfig CalculatedPolicyConfig `json:"calculatedConfig"`
}

type PolicyRuleGetForLabelRequest struct {
	Labels []string
}

type PolicyRuleLicensePublicResponse struct {
	Key        string  `json:"key" example:"dummy-key-----cf422e8a94-03b3-4c6c-9954-8c183d7177cf"`
	Identifier string  `json:"identifier" example:"AGPL-3.0-or-later"`
	Name       string  `json:"name" example:"Name of license"`
	Aliases    []Alias `json:"alias"`
} //	@name	PolicyRuleLicense

type PolicyRulePublicResponseDto struct {
	Key         string                            `json:"key" example:"dummy-key-----4e42-a64a-c0b362976f9c"`
	Name        string                            `json:"name" example:"Policy rule name"`
	Description string                            `json:"description" example:"Example description"`
	Type        ListType                          `json:"type" example:"deny"`
	Created     time.Time                         `json:"created" example:"2023-03-21T04:45:00.806887937Z"`
	Updated     time.Time                         `json:"updated" example:"2023-07-11T11:26:07.440978865Z"`
	Licenses    []PolicyRuleLicensePublicResponse `json:"licenses"`
} //	@name	PolicyRule

type ImportResponse struct {
	JobId string
}

type PolicyRulesStatus string

const (
	StatusActive     PolicyRulesStatus = "active"
	StatusInactive   PolicyRulesStatus = "inactive"
	StatusDeprecated PolicyRulesStatus = "deprecated"
)

type PolicyRulesDto struct {
	domain.BaseDto

	Status      PolicyRulesStatus `json:"status"`
	Name        string            `json:"name"`
	LabelSets   [][]string        `json:"labelSets"`
	Description string            `json:"description"`

	ComponentsAllow []string `json:"componentsAllow"`
	ComponentsDeny  []string `json:"componentsDeny"`
	ComponentsWarn  []string `json:"componentsWarn"`

	Auxiliary        bool                   `json:"auxiliary"`
	Deprecated       bool                   `json:"deprecated"`
	DeprecatedDate   time.Time              `json:"deprecatedDate"`
	Active           bool                   `json:"active"`
	ApplyToAll       bool                   `json:"applyToAll"`
	Calculated       bool                   `json:"calculated"`
	CalculatedConfig CalculatedPolicyConfig `json:"calculatedConfig"`
}

type PolicyRulesAssignmentDto struct {
	Status      PolicyRulesStatus `json:"status"`
	Key         string            `json:"key"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Type        ListType          `json:"type"`
}

type PolicyRulesForLicenseDto struct {
	Id                     string                     `json:"id"`
	PolicyRulesAssignments []PolicyRulesAssignmentDto `json:"policyRulesAssignments"`
}

type CustomLicensesDto struct {
	AllIDs []string          `json:"allIds"`
	Texts  map[string]string `json:"texts"`
}

type LookupRequestDto struct {
	Ids []string `json:"ids" validate:"gte=0"`
}

type LookupResponseDto struct {
	Items []LicenseSlimDto `json:"items"`
}

type Licenses []*License

func (l Licenses) ToPublicResDtos() []PolicyRuleLicensePublicResponse {
	var res []PolicyRuleLicensePublicResponse
	for _, lic := range l {
		res = append(res, PolicyRuleLicensePublicResponse{
			Key:        lic.Key,
			Identifier: lic.LicenseId,
			Name:       lic.Name,
			Aliases:    lic.Aliases,
		})
	}
	return res
}
