// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package auditloglist

import (
	"encoding/json"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/audit"
	"github.com/eclipse-disuko/disuko/domain/project/auditloglist"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type auditLogListRepositoryStruct struct {
	base.BaseRepositoryWithSoftDelete[*auditloglist.AuditLogList]
}

func CreateCompareString(diffFunc audit.DiffFunc, after, before interface{}, ignoreFieldsOption ...cmp.Option) string {
	options := []cmp.Option{cmpopts.IgnoreTypes(audit.Container{})}
	if len(ignoreFieldsOption) > 0 {
		options = append(options, ignoreFieldsOption...)
	}
	return diffFunc(before, after, options...)
}

func NewAuditLogListRepository(requestSession *logy.RequestSession) IAuditLogListRepository {
	return &auditLogListRepositoryStruct{
		BaseRepositoryWithSoftDelete: base.CreateRepositoryWithSoftDelete[*auditloglist.AuditLogList](
			requestSession,
			AuditLogListCollectionName,
			func() *auditloglist.AuditLogList {
				return &auditloglist.AuditLogList{}
			},
			nil,
			nil,
			nil)}
}

func (repository *auditLogListRepositoryStruct) CreateAuditEntriesByKey(requestSession *logy.RequestSession, key string, auditEntries []*audit.Audit) {
	auditEntity := repository.FindByKey(requestSession, key, false)
	alreadyExists := auditEntity != nil
	if auditEntity == nil {
		auditEntity = &auditloglist.AuditLogList{
			RootEntity: domain.NewRootEntity(),
		}
		auditEntity.AuditTrail = make([]*audit.Audit, 0)
		auditEntity.Key = key
	}
	auditEntity.AuditTrail = append(auditEntity.AuditTrail, auditEntries...)
	if alreadyExists {
		repository.Update(requestSession, auditEntity)
	} else {
		repository.Save(requestSession, auditEntity)
	}
}

func (repository *auditLogListRepositoryStruct) CreateAuditEntryByKey(requestSession *logy.RequestSession, key string, user string, message string, diffFunc audit.DiffFunc, after, before interface{}, ignoreFieldsOption ...cmp.Option) {
	compareStr := CreateCompareString(diffFunc, after, before, ignoreFieldsOption...)
	auditEntity := repository.FindByKey(requestSession, key, false)
	alreadyExists := auditEntity != nil
	if auditEntity == nil {
		auditEntity = &auditloglist.AuditLogList{
			RootEntity: domain.NewRootEntity(),
		}
		auditEntity.AuditTrail = make([]*audit.Audit, 0)
		auditEntity.Key = key
	}
	auditEntry := audit.NewAudit(user, message, compareStr)
	auditEntity.AuditTrail = append(auditEntity.AuditTrail, auditEntry)
	if alreadyExists {
		repository.Update(requestSession, auditEntity)
	} else {
		repository.Save(requestSession, auditEntity)
	}
}

func (repository *auditLogListRepositoryStruct) AddStaticAuditEntryByKey(requestSession *logy.RequestSession, key string, user string, message string, entryData interface{}) {
	jsonString, err := json.MarshalIndent(entryData, "", "    ")
	if err != nil {
		exception.ThrowExceptionSendDeniedResponse()
	}
	auditEntity := repository.FindByKey(requestSession, key, false)
	alreadyExists := auditEntity != nil
	if auditEntity == nil {
		auditEntity = &auditloglist.AuditLogList{
			RootEntity: domain.NewRootEntity(),
		}
		auditEntity.AuditTrail = make([]*audit.Audit, 0)
		auditEntity.Key = key
	}
	auditEntry := audit.NewAudit(user, message, string(jsonString))
	auditEntity.AuditTrail = append(auditEntity.AuditTrail, auditEntry)
	if alreadyExists {
		repository.Update(requestSession, auditEntity)
	} else {
		repository.Save(requestSession, auditEntity)
	}
}
