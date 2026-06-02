// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package policydecisions

import (
	"encoding/json"
	"time"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/helper/hash"
	"github.com/eclipse-disuko/disuko/logy"
)

type PolicyDecision struct {
	domain.ChildEntity `bson:",inline"`

	SBOMId            string
	SBOMName          string
	SBOMUploaded      *time.Time
	ComponentSpdxId   string
	ComponentName     string
	ComponentVersion  string
	LicenseExpression string
	LicenseId         string
	PolicyId          string
	PolicyEvaluated   string
	PolicyDecision    string
	Comment           string
	Creator           string
	Active            bool
	PreviewMode       bool
}

type PolicyDecisions struct {
	domain.RootEntity `bson:",inline"`
	domain.SoftDelete `bson:",inline"`

	Decisions []*PolicyDecision
}

func (pd *PolicyDecisions) GenHash(requestSession *logy.RequestSession) string {
	if pd == nil {
		return ""
	}
	ruleStr, err := json.Marshal(pd)
	if err != nil {
		logy.Warnf(requestSession, "Error marshalling policy decisions: %s", pd.Key)
		return ""
	}
	return hash.Hash(requestSession, ruleStr)
}
