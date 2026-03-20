// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package userstats

import (
	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/audit"
)

const (
	SystemStatsUser = "SYSTEM_DISCO"
)

type UserStatus struct {
	domain.RootEntity   `bson:",inline"`
	audit.Container     `bson:",inline"`
	User                string
	ProjectCount        int
	LicenseCount        int
	PolicyRuleCount     int
	LabelCount          int
	SchemaCount         int
	ObligationCount     int
	UserCount           int
	DisclosureCount     int
	ReviewTemplateCount int
	ActiveJobCount      int
}
