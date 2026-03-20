// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package auditloglist

import (
	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/audit"
)

type AuditLogList struct {
	domain.RootEntity `bson:",inline"`
	domain.SoftDelete `bson:",inline"`

	AuditTrail []*audit.Audit
}
