// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package policydecisions

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain"
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
