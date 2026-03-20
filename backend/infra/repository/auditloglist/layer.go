// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package auditloglist

import (
	"github.com/google/go-cmp/cmp"
	"mercedes-benz.ghe.com/foss/disuko/domain/audit"
	"mercedes-benz.ghe.com/foss/disuko/domain/project/auditloglist"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

const AuditLogListCollectionName = "auditLogList"

type IAuditLogListRepository interface {
	base.IBaseRepositoryWithSoftDelete[*auditloglist.AuditLogList]
	CreateAuditEntryByKey(requestSession *logy.RequestSession, key string, user string, message string, diffFunc audit.DiffFunc, after, before interface{}, ignoreFieldsOption ...cmp.Option)
	AddStaticAuditEntryByKey(requestSession *logy.RequestSession, key string, user string, message string, entryData interface{})
	CreateAuditEntriesByKey(requestSession *logy.RequestSession, key string, auditEntries []*audit.Audit)
}
