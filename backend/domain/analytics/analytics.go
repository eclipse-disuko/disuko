// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package analytics

import (
	"time"

	"mercedes-benz.ghe.com/foss/disuko/domain"
	"mercedes-benz.ghe.com/foss/disuko/domain/project/components"
	"mercedes-benz.ghe.com/foss/disuko/domain/user/approval"
)

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
	SBomStatus     approval.ApprovalStatus
	SBomLastUpdate time.Time
}
