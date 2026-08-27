// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package analytics

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/approval"
	"github.com/eclipse-disuko/disuko/domain/project/components"
)

type SbomType string

const (
	SbomTypeLatest            SbomType = "LATEST"
	SbomTypeLatestApproved    SbomType = "LATEST_APPROVED"
	SbomTypeLatestAndApproved SbomType = "LATEST_AND_LATEST_APPROVED"
)

func ParseSbomType(value string) (valid bool, result SbomType) {
	switch value {
	case string(SbomTypeLatest):
		valid, result = true, SbomTypeLatest
	case string(SbomTypeLatestApproved):
		valid, result = true, SbomTypeLatestApproved
	case string(SbomTypeLatestAndApproved):
		valid, result = true, SbomTypeLatestAndApproved
	default:
		valid, result = false, ""
	}
	return
}

type Analytics struct {
	domain.RootEntity `bson:"inline"`
	domain.SoftDelete `bson:"inline"`

	ProjectKey  string
	ProjectName string
	Responsible string

	ProjectVersionKey  string
	ProjectVersionName string

	OwnerDeptId string

	ComponentName    string
	ComponentVersion string

	LicenseConcluded string
	LicenseDeclared  string
	EntryLicense     string
	Licenses         components.LicenseList

	SBomKey        string
	SBomName       string
	SBomStatus     approval.StateInfo
	SBomLastUpdate time.Time
	SBomType       SbomType
}
