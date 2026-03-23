// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package deletionaudit

import (
	"github.com/eclipse-disuko/disuko/domain/deletionaudit"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	db "github.com/eclipse-disuko/disuko/infra/repository/database"
	"github.com/eclipse-disuko/disuko/logy"
)

const DeletionAuditCollectionName = "deletionAuditLog"

type IDeletionAuditRepository interface {
	base.IBaseRepository[*deletionaudit.DeletionAuditEntry]

	FindByOperationID(rs *logy.RequestSession, operationID string) []*deletionaudit.DeletionAuditEntry
	FindByTargetUser(rs *logy.RequestSession, targetUser string) []*deletionaudit.DeletionAuditEntry
	FindByPerformedBy(rs *logy.RequestSession, performedBy string) []*deletionaudit.DeletionAuditEntry
}

type deletionAuditRepositoryStruct struct {
	base.BaseRepository[*deletionaudit.DeletionAuditEntry]
}

func NewDeletionAuditRepository(rs *logy.RequestSession) IDeletionAuditRepository {
	baseRepo := base.BaseRepository[*deletionaudit.DeletionAuditEntry]{
		CollectionName: DeletionAuditCollectionName,
		EntityCreator: func() *deletionaudit.DeletionAuditEntry {
			return &deletionaudit.DeletionAuditEntry{}
		},
		IsSoftDelete: false,
	}
	baseRepo.Database = base.NewDatabase()
	baseRepo.Database.Init(rs, DeletionAuditCollectionName, [][]string{
		{"Created"},
		{"Updated"},
		{"Deleted"},
		{"_id", "Deleted"},
		{"PerformedBy"},
		{"TargetUser"},
		{"OperationID"},
		{"Category"},
		{"Timestamp"},
	})

	return &deletionAuditRepositoryStruct{
		BaseRepository: baseRepo,
	}
}

func (r *deletionAuditRepositoryStruct) FindByOperationID(rs *logy.RequestSession, operationID string) []*deletionaudit.DeletionAuditEntry {
	qc := &db.QueryConfig{}
	qc.SetMatcher(db.AttributeMatcher("OperationID", db.EQ, operationID))
	qc.SetSort(db.SortConfig{
		db.SortAttribute{Name: "Timestamp", Order: db.ASC},
	})
	return r.Query(rs, qc)
}

func (r *deletionAuditRepositoryStruct) FindByTargetUser(rs *logy.RequestSession, targetUser string) []*deletionaudit.DeletionAuditEntry {
	qc := &db.QueryConfig{}
	qc.SetMatcher(db.AttributeMatcher("TargetUser", db.EQ, targetUser))
	qc.SetSort(db.SortConfig{
		db.SortAttribute{Name: "Timestamp", Order: db.DESC},
	})
	return r.Query(rs, qc)
}

func (r *deletionAuditRepositoryStruct) FindByPerformedBy(rs *logy.RequestSession, performedBy string) []*deletionaudit.DeletionAuditEntry {
	qc := &db.QueryConfig{}
	qc.SetMatcher(db.AttributeMatcher("PerformedBy", db.EQ, performedBy))
	qc.SetSort(db.SortConfig{
		db.SortAttribute{Name: "Timestamp", Order: db.DESC},
	})
	return r.Query(rs, qc)
}
