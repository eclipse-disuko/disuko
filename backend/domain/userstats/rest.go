// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package userstats

import (
	"github.com/eclipse-disuko/disuko/domain/statistic"
)

func (entity *UserStatus) ToDashboardCounts() *statistic.DashboardCounts {
	result := &statistic.DashboardCounts{
		ProjectCount:        entity.ProjectCount,
		LicenseCount:        entity.LicenseCount,
		PolicyRuleCount:     entity.PolicyRuleCount,
		LabelCount:          entity.LabelCount,
		SchemaCount:         entity.SchemaCount,
		ObligationCount:     entity.ObligationCount,
		UserCount:           entity.UserCount,
		DisclosureCount:     entity.DisclosureCount,
		ReviewTemplateCount: entity.ReviewTemplateCount,
		ActiveJobCount:      entity.ActiveJobCount,
	}
	return result
}
