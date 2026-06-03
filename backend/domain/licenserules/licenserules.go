// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package licenserules

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/helper/hash"
	"github.com/eclipse-disuko/disuko/logy"
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

type LicenseRulesHashEntry struct {
	Key    string
	Active bool
}

type LicenseRulesHash struct {
	Key   string
	Rules []LicenseRulesHashEntry
}

func (r *LicenseRules) GenHash(requestSession *logy.RequestSession) string {
	if r == nil {
		return ""
	}

	entries := make([]LicenseRulesHashEntry, 0, len(r.Rules))

	for _, rule := range r.Rules {
		if rule == nil {
			continue
		}
		entries = append(entries, LicenseRulesHashEntry{
			Key:    rule.Key,
			Active: rule.Active,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	hashData := LicenseRulesHash{
		Key:   r.Key,
		Rules: entries,
	}

	ruleStr, err := json.Marshal(hashData)
	if err != nil {
		logy.Warnf(requestSession, "Error marshalling license rules hash data: %s", r.Key)
		return ""
	}
	return hash.Hash(requestSession, ruleStr)
}
