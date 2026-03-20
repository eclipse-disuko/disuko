// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package license

import (
	auditHelper "github.com/eclipse-disuko/disuko/helper/audit"
	"github.com/eclipse-disuko/disuko/infra/repository/labels"
	"github.com/eclipse-disuko/disuko/logy"
)

type PolicyRulesAudit struct {
	Key         string
	Name        string     `json:"name"`
	LabelsSets  [][]string `json:"labelsSets"`
	Description string     `json:"description"`

	ComponentsAllow []string `json:"componentsAllow"`
	ComponentsDeny  []string `json:"componentsDeny"`
	ComponentsWarn  []string `json:"componentsWarn"`

	Auxiliary  bool `json:"auxiliary"`
	Deprecated bool `json:"deprecated"`
	Active     bool `json:"active"`
	ApplyToAll bool `json:"applyToAll"`

	Created string `json:"created"`
	Updated string `json:"updated"`
}

func (entity *PolicyRules) ToAudit(requestSession *logy.RequestSession, labelProvider labels.ILabelRepository) *PolicyRulesAudit {
	labelNamesSets := make([][]string, 0)
	if labelProvider != nil {
		for _, labelsSet := range entity.LabelSets {
			labelNamesSet := make([]string, 0)
			for _, labelKey := range labelsSet {
				resolvedLabel := labelProvider.FindByKey(requestSession, labelKey, false)
				labelNamesSet = append(labelNamesSet, resolvedLabel.Name)
			}
			if len(labelNamesSet) > 0 {
				labelNamesSets = append(labelNamesSets, labelNamesSet)
			}
		}
	}

	return &PolicyRulesAudit{
		Key:         entity.Key,
		Name:        entity.Name,
		LabelsSets:  labelNamesSets,
		Description: entity.Description,

		ComponentsAllow: entity.ComponentsAllow,
		ComponentsDeny:  entity.ComponentsDeny,
		ComponentsWarn:  entity.ComponentsWarn,

		Auxiliary:  entity.Auxiliary,
		Deprecated: entity.Deprecated,
		Active:     entity.Active,
		ApplyToAll: entity.ApplyToAll,

		Created: auditHelper.ConvertDateTime(entity.Created),
		Updated: auditHelper.ConvertDateTime(entity.Updated),
	}
}
