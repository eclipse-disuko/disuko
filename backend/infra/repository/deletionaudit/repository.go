// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package deletionaudit

import (
	"github.com/eclipse-disuko/disuko/domain/deletionaudit"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
)

type deletionAuditRepositoryStruct struct {
	base.BaseRepositoryWithHardDelete[*deletionaudit.DeletionAuditEntry]
}

func NewDeletionAuditRepository(rs *logy.RequestSession) IDeletionAuditRepository {
	return &deletionAuditRepositoryStruct{
		BaseRepositoryWithHardDelete: base.CreateRepositoryWithHardDelete[*deletionaudit.DeletionAuditEntry](
			rs,
			DeletionAuditCollectionName,
			func() *deletionaudit.DeletionAuditEntry {
				return &deletionaudit.DeletionAuditEntry{}
			},
			nil,
			"",
			nil,
			[][]string{
				{"PerformedBy"},
				{"TargetUser"},
				{"OperationID"},
				{"Category"},
				{"Timestamp"},
			}),
	}
}
