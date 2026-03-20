// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package auditloglist

import (
	"github.com/eclipse-disuko/disuko/domain/audit"
	"github.com/eclipse-disuko/disuko/domain/project/auditloglist"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/google/go-cmp/cmp"
)

const AuditLogListCollectionName = "auditLogList"

type IAuditLogListRepository interface {
	base.IBaseRepositoryWithSoftDelete[*auditloglist.AuditLogList]
	CreateAuditEntryByKey(requestSession *logy.RequestSession, key string, user string, message string, diffFunc audit.DiffFunc, after, before interface{}, ignoreFieldsOption ...cmp.Option)
	AddStaticAuditEntryByKey(requestSession *logy.RequestSession, key string, user string, message string, entryData interface{})
	CreateAuditEntriesByKey(requestSession *logy.RequestSession, key string, auditEntries []*audit.Audit)
}
