// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package licenserules

import (
	"encoding/json"
	"time"

	"mercedes-benz.ghe.com/foss/disuko/domain"
	"mercedes-benz.ghe.com/foss/disuko/helper/hash"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type LicenseRule struct {
	domain.ChildEntity `bson:"inline"`

	SBOMId              string
	SBOMName            string
	SBOMUploaded        *time.Time
	ComponentSpdxId     string
	ComponentName       string
	ComponentVersion    string
	LicenseExpression   string
	LicenseDecisionId   string
	LicenseDecisionName string
	Comment             string
	Creator             string
	Active              bool
	PreviewMode         bool
}

type LicenseRules struct {
	domain.RootEntity `bson:"inline"`
	domain.SoftDelete `bson:"inline"`

	Rules []*LicenseRule
}

func (r *LicenseRules) GenHash(requestSession *logy.RequestSession) string {
	if r == nil {
		return ""
	}
	ruleStr, err := json.Marshal(r)
	if err != nil {
		logy.Warnf(requestSession, "Error marshalling license rules: %s", r.Key)
		return ""
	}
	return hash.Hash(requestSession, ruleStr)
}
