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

// IDeletionAuditRepository defines CRUD + query operations for deletion audit entries.
type IDeletionAuditRepository interface {
	base.IBaseRepository[*deletionaudit.DeletionAuditEntry]

	// SaveEntry persists a single audit entry.
	SaveEntry(rs *logy.RequestSession, entry *deletionaudit.DeletionAuditEntry)

	// SaveEntries persists multiple audit entries in bulk.
	SaveEntries(rs *logy.RequestSession, entries []*deletionaudit.DeletionAuditEntry)

	// FindByOperationID returns all entries for a given operation.
	FindByOperationID(rs *logy.RequestSession, operationID string) []*deletionaudit.DeletionAuditEntry

	// FindByTargetUser returns all deletion audit entries for a given target user.
	FindByTargetUser(rs *logy.RequestSession, targetUser string) []*deletionaudit.DeletionAuditEntry

	// FindByPerformedBy returns all entries initiated by a given admin.
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
		{"performedBy"},
		{"targetUser"},
		{"operationId"},
		{"category"},
		{"timestamp"},
	})

	return &deletionAuditRepositoryStruct{
		BaseRepository: baseRepo,
	}
}

func (r *deletionAuditRepositoryStruct) SaveEntry(rs *logy.RequestSession, entry *deletionaudit.DeletionAuditEntry) {
	r.Save(rs, entry)
}

func (r *deletionAuditRepositoryStruct) SaveEntries(rs *logy.RequestSession, entries []*deletionaudit.DeletionAuditEntry) {
	r.SaveList(rs, entries, true)
}

func (r *deletionAuditRepositoryStruct) FindByOperationID(rs *logy.RequestSession, operationID string) []*deletionaudit.DeletionAuditEntry {
	qc := &db.QueryConfig{}
	qc.SetMatcher(db.AttributeMatcher("operationId", db.EQ, operationID))
	qc.SetSort(db.SortConfig{
		db.SortAttribute{Name: "timestamp", Order: db.ASC},
	})
	return r.Query(rs, qc)
}

func (r *deletionAuditRepositoryStruct) FindByTargetUser(rs *logy.RequestSession, targetUser string) []*deletionaudit.DeletionAuditEntry {
	qc := &db.QueryConfig{}
	qc.SetMatcher(db.AttributeMatcher("targetUser", db.EQ, targetUser))
	qc.SetSort(db.SortConfig{
		db.SortAttribute{Name: "timestamp", Order: db.DESC},
	})
	return r.Query(rs, qc)
}

func (r *deletionAuditRepositoryStruct) FindByPerformedBy(rs *logy.RequestSession, performedBy string) []*deletionaudit.DeletionAuditEntry {
	qc := &db.QueryConfig{}
	qc.SetMatcher(db.AttributeMatcher("performedBy", db.EQ, performedBy))
	qc.SetSort(db.SortConfig{
		db.SortAttribute{Name: "timestamp", Order: db.DESC},
	})
	return r.Query(rs, qc)
}
