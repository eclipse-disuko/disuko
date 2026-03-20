// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package decisions

import "time"

type DecisionType string

const (
	LicenseDecision = "license"
	PolicyDecision  = "policy"
)

type DecisionDto struct {
	Key               string     `json:"key"`
	Created           time.Time  `json:"created"`
	Updated           time.Time  `json:"updated"`
	SBOMId            string     `json:"sbomId"`
	SBOMName          string     `json:"sbomName"`
	SBOMUploaded      *time.Time `json:"sbomUploaded"`
	ComponentSpdxId   string     `json:"componentSpdxId"`
	ComponentName     string     `json:"componentName"`
	ComponentVersion  string     `json:"componentVersion"`
	LicenseExpression string     `json:"licenseExpression"`
	Comment           string     `json:"comment"`
	Creator           string     `json:"creator"`
	Active            bool       `json:"active"`

	Type DecisionType `json:"type"`

	LicenseDecisionId   string `json:"licenseDecisionId"`
	LicenseDecisionName string `json:"licenseDecisionName"`

	LicenseMatchedId string `json:"licenseMatchedId"`
	PolicyId         string `json:"policyId"`
	PolicyEvaluated  string `json:"policyEvaluated"`
	PolicyDecision   string `json:"policyDecision"`
}
