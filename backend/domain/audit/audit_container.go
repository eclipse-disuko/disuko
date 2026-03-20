// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package audit

import (
	"github.com/google/go-cmp/cmp"
)

type IContainer interface {
	GetAuditTrail() []*Audit
	AddAuditEntry(audit *Audit)
}

type Container struct {
	AuditTrail []*Audit
}

func (container *Container) GetAuditTrail() []*Audit {
	if container.AuditTrail == nil {
		container.AuditTrail = make([]*Audit, 0)
	}
	return container.AuditTrail
}

func (container *Container) AddAuditEntry(audit *Audit) {
	auditLog := container.GetAuditTrail()
	container.AuditTrail = append(auditLog, audit)
}

type DiffFunc func(before, after interface{}, options ...cmp.Option) string
