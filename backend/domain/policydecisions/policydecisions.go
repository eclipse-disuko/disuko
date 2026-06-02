// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package policydecisions

import (
	"encoding/json"
	"sort"
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

type PolicyDecisionHashEntry struct {
	Key    string
	Active bool
}

type PolicyDecisionsHash struct {
	Key       string
	Decisions []PolicyDecisionHashEntry
}

func (pd *PolicyDecisions) GenHash(requestSession *logy.RequestSession) string {
	if pd == nil {
		return ""
	}

	entries := make([]PolicyDecisionHashEntry, 0, len(pd.Decisions))

	for _, decision := range pd.Decisions {
		if decision == nil {
			continue
		}
		entries = append(entries, PolicyDecisionHashEntry{
			Key:    decision.Key,
			Active: decision.Active,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	hashData := PolicyDecisionsHash{
		Key:       pd.Key,
		Decisions: entries,
	}

	ruleStr, err := json.Marshal(hashData)
	if err != nil {
		logy.Warnf(requestSession, "Error marshalling policy decisions hash data: %s", pd.Key)
		return ""
	}
	return hash.Hash(requestSession, ruleStr)
}
