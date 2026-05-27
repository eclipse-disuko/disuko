// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package policyruleclassification

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain"
)

type RuleStatus string
type ClassificationKey string

const (
	RuleStatusAllowed   RuleStatus = "allowed"
	RuleStatusWarned    RuleStatus = "warned"
	RuleStatusDenied    RuleStatus = "denied"
	RuleStatusForbidden RuleStatus = "forbidden"
)

type PolicyRuleClassification struct {
	domain.RootEntity `bson:",inline"`
	Name              string                           `json:"name"`
	Rules             map[ClassificationKey]RuleStatus `json:"rules"`
}

func (p *PolicyRuleClassification) Update(name string, rules map[ClassificationKey]RuleStatus) {
	p.Name = name
	p.Rules = rules
	p.Updated = time.Now()
}
