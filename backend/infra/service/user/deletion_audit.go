// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package user

import (
	"encoding/json"

	"github.com/eclipse-disuko/disuko/domain/deletionaudit"
	deletionAuditRepo "github.com/eclipse-disuko/disuko/infra/repository/deletionaudit"
	"github.com/eclipse-disuko/disuko/logy"
)

// DeletionAuditService wraps every user-deletion action with detailed, immutable audit records.
// It acts as the single integration point between the deletion workflow and the audit trail.
type DeletionAuditService struct {
	rs        *logy.RequestSession
	auditRepo deletionAuditRepo.IDeletionAuditRepository
}

func NewDeletionAuditService(rs *logy.RequestSession, auditRepo deletionAuditRepo.IDeletionAuditRepository) *DeletionAuditService {
	return &DeletionAuditService{
		rs:        rs,
		auditRepo: auditRepo,
	}
}

// ------------------------------------------------------------------
// Core recording methods
// ------------------------------------------------------------------

// RecordDeletionPlan records all actions from a DeletionPlan as audit entries.
// Each action in the plan becomes an individual, queryable audit row.
func (s *DeletionAuditService) RecordDeletionPlan(plan *DeletionPlan, performedBy, ipAddress, requestID string) string {
	operationID := deletionaudit.NewOperationID()
	var entries []*deletionaudit.DeletionAuditEntry

	for _, action := range plan.TaskActions {
		entry := s.buildEntry(performedBy, plan.Username, operationID, deletionaudit.CategoryTask, action, ipAddress, requestID)
		entries = append(entries, entry)
	}

	for _, action := range plan.RoleActions {
		entry := s.buildEntry(performedBy, plan.Username, operationID, deletionaudit.CategoryRole, action, ipAddress, requestID)
		entries = append(entries, entry)
	}

	for _, action := range plan.TraceActions {
		entry := s.buildEntry(performedBy, plan.Username, operationID, deletionaudit.CategoryTrace, action, ipAddress, requestID)
		entries = append(entries, entry)
	}

	for _, action := range plan.ProfileActions {
		entry := s.buildEntry(performedBy, plan.Username, operationID, deletionaudit.CategoryProfile, action, ipAddress, requestID)
		entries = append(entries, entry)
	}

	if len(entries) > 0 {
		s.auditRepo.SaveEntries(s.rs, entries)
	}

	logy.Infof(s.rs, "Recorded deletion audit: operation=%s, performedBy=%s, targetUser=%s, entries=%d",
		operationID, performedBy, plan.Username, len(entries))

	return operationID
}

// RecordSingleEntityDeletion records a single entity deletion (e.g. deleting one task by ID).
func (s *DeletionAuditService) RecordSingleEntityDeletion(
	performedBy, targetUser, entityType, entityID string,
	category deletionaudit.ActionCategory,
	result deletionaudit.ActionResult,
	reason, snapshotJSON, projectID, projectName, ipAddress, requestID string,
) string {
	operationID := deletionaudit.NewOperationID()

	entry := deletionaudit.NewDeletionAuditEntry(performedBy, targetUser, operationID, category, "delete_"+entityType, result)
	entry.EntityID = entityID
	entry.EntityType = entityType
	entry.ProjectID = projectID
	entry.ProjectName = projectName
	entry.Reason = reason
	entry.SnapshotJSON = snapshotJSON
	entry.IPAddress = ipAddress
	entry.RequestID = requestID

	s.auditRepo.SaveEntry(s.rs, entry)

	logy.Infof(s.rs, "Recorded single deletion audit: operation=%s, entity=%s/%s, result=%s, performedBy=%s, target=%s",
		operationID, entityType, entityID, result, performedBy, targetUser)

	return operationID
}

// RecordBulkEntityDeletion records a bulk deletion of all entities of a given type for a user.
func (s *DeletionAuditService) RecordBulkEntityDeletion(
	performedBy, targetUser, entityType string,
	category deletionaudit.ActionCategory,
	result deletionaudit.ActionResult,
	entityIDs []string,
	reason, ipAddress, requestID string,
) string {
	operationID := deletionaudit.NewOperationID()
	var entries []*deletionaudit.DeletionAuditEntry

	for _, entityID := range entityIDs {
		entry := deletionaudit.NewDeletionAuditEntry(performedBy, targetUser, operationID, category, "bulk_delete_"+entityType, result)
		entry.EntityID = entityID
		entry.EntityType = entityType
		entry.Reason = reason
		entry.IPAddress = ipAddress
		entry.RequestID = requestID
		entries = append(entries, entry)
	}

	if len(entries) > 0 {
		s.auditRepo.SaveEntries(s.rs, entries)
	}

	logy.Infof(s.rs, "Recorded bulk deletion audit: operation=%s, type=%s, count=%d, performedBy=%s, target=%s",
		operationID, entityType, len(entries), performedBy, targetUser)

	return operationID
}

func (s *DeletionAuditService) RecordRoleDeletionResults(
	performedBy, targetUser string,
	results []RoleDeletionResult,
	ipAddress, requestID string,
) string {
	operationID := deletionaudit.NewOperationID()
	var entries []*deletionaudit.DeletionAuditEntry

	for _, r := range results {
		if r.Skipped {
			continue
		}

		entry := deletionaudit.NewDeletionAuditEntry(
			performedBy, targetUser, operationID,
			deletionaudit.CategoryRole, "remove_role", deletionaudit.ResultSuccess,
		)
		entry.EntityID = r.ProjectID
		entry.EntityType = r.RoleName
		entry.ProjectID = r.ProjectID
		entry.ProjectName = r.ProjectName
		entry.Reason = r.SkipReason
		entry.IPAddress = ipAddress
		entry.RequestID = requestID

		snapshot, err := json.Marshal(r)
		if err == nil {
			entry.SnapshotJSON = string(snapshot)
		}

		entries = append(entries, entry)
	}

	if len(entries) > 0 {
		s.auditRepo.SaveEntries(s.rs, entries)
	}

	logy.Infof(s.rs, "Recorded role deletion audit: operation=%s, performedBy=%s, targetUser=%s, entries=%d",
		operationID, performedBy, targetUser, len(entries))

	return operationID
}

// ------------------------------------------------------------------
// Query / retrieval methods
// ------------------------------------------------------------------

// GetAuditByOperation returns all audit entries for a deletion operation grouped into a summary.
func (s *DeletionAuditService) GetAuditByOperation(operationID string) *deletionaudit.DeletionAuditSummaryDto {
	entries := s.auditRepo.FindByOperationID(s.rs, operationID)
	if len(entries) == 0 {
		return nil
	}

	summary := &deletionaudit.DeletionAuditSummaryDto{
		OperationID:  operationID,
		PerformedBy:  entries[0].PerformedBy,
		TargetUser:   entries[0].TargetUser,
		Timestamp:    entries[0].Timestamp,
		TotalActions: len(entries),
		TaskEntries:  make([]deletionaudit.DeletionAuditEntryDto, 0),
		RoleEntries:  make([]deletionaudit.DeletionAuditEntryDto, 0),
		TraceEntries: make([]deletionaudit.DeletionAuditEntryDto, 0),
		ProfileEntry: make([]deletionaudit.DeletionAuditEntryDto, 0),
	}

	for _, e := range entries {
		dto := e.ToDto()
		switch e.Category {
		case deletionaudit.CategoryTask:
			summary.TaskEntries = append(summary.TaskEntries, dto)
		case deletionaudit.CategoryRole:
			summary.RoleEntries = append(summary.RoleEntries, dto)
		case deletionaudit.CategoryTrace:
			summary.TraceEntries = append(summary.TraceEntries, dto)
		case deletionaudit.CategoryProfile:
			summary.ProfileEntry = append(summary.ProfileEntry, dto)
		}

		switch e.Result {
		case deletionaudit.ResultSuccess:
			summary.Succeeded++
		case deletionaudit.ResultSkipped:
			summary.Skipped++
		case deletionaudit.ResultFailed:
			summary.Failed++
		case deletionaudit.ResultRetained:
			summary.Retained++
		}
	}

	return summary
}

// GetAuditTrailForUser returns all deletion audit entries for a target user.
func (s *DeletionAuditService) GetAuditTrailForUser(targetUser string) []deletionaudit.DeletionAuditEntryDto {
	entries := s.auditRepo.FindByTargetUser(s.rs, targetUser)
	return deletionaudit.ToDtos(entries)
}

// GetAuditTrailByAdmin returns all deletion audit entries performed by a given admin.
func (s *DeletionAuditService) GetAuditTrailByAdmin(adminUser string) []deletionaudit.DeletionAuditEntryDto {
	entries := s.auditRepo.FindByPerformedBy(s.rs, adminUser)
	return deletionaudit.ToDtos(entries)
}

// ------------------------------------------------------------------
// Internal helpers
// ------------------------------------------------------------------

func (s *DeletionAuditService) buildEntry(
	performedBy, targetUser, operationID string,
	category deletionaudit.ActionCategory,
	action DeletionAction,
	ipAddress, requestID string,
) *deletionaudit.DeletionAuditEntry {

	result := mapActionStatusToResult(action.Status)

	entry := deletionaudit.NewDeletionAuditEntry(performedBy, targetUser, operationID, category, action.ActionType, result)
	entry.EntityID = action.EntityID
	entry.EntityType = action.EntityType
	entry.ProjectID = action.ProjectID
	entry.ProjectName = action.ProjectName
	entry.Reason = action.Reason
	entry.IPAddress = ipAddress
	entry.RequestID = requestID

	// Capture a snapshot of the action details for forensic purposes.
	snapshot, err := json.Marshal(action)
	if err == nil {
		entry.SnapshotJSON = string(snapshot)
	}

	return entry
}

func mapActionStatusToResult(status string) deletionaudit.ActionResult {
	switch status {
	case "completed":
		return deletionaudit.ResultSuccess
	case "skipped":
		return deletionaudit.ResultSkipped
	case "planned":
		return deletionaudit.ResultPlanned
	default:
		return deletionaudit.ResultFailed
	}
}
