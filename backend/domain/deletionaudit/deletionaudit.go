// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package deletionaudit

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/google/uuid"
)

// ActionCategory classifies which part of user data is being deleted.
type ActionCategory string

const (
	CategoryTask    ActionCategory = "TASK"
	CategoryRole    ActionCategory = "ROLE"
	CategoryTrace   ActionCategory = "TRACE"
	CategoryProfile ActionCategory = "PROFILE"
)

// ActionResult captures the outcome of a single deletion action.
type ActionResult string

const (
	ResultSuccess  ActionResult = "SUCCESS"
	ResultSkipped  ActionResult = "SKIPPED"
	ResultFailed   ActionResult = "FAILED"
	ResultRetained ActionResult = "RETAINED"
	ResultPlanned  ActionResult = "PLANNED"
	ResultDryRun   ActionResult = "DRY_RUN"
)

// DeletionAuditEntry is the core entity stored in the "deletionAuditLog" collection.
// Each entry represents one auditable deletion action (or batch), immutably recorded.
type DeletionAuditEntry struct {
	domain.RootEntity `bson:",inline"`

	// Who performed the deletion (admin user ID / "SYSTEM" for automated jobs).
	PerformedBy string `bson:"performedBy" json:"performedBy"`

	// The target user whose data is being deleted.
	TargetUser string `bson:"targetUser" json:"targetUser"`

	// High-level operation identifier grouping related actions.
	// All actions within one "Delete All" call share the same OperationID.
	OperationID string `bson:"operationId" json:"operationId"`

	// Category of the entity being deleted.
	Category ActionCategory `bson:"category" json:"category"`

	// The specific action taken (e.g. "cancel_task", "remove_role", "delete_trace", "delete_profile").
	ActionType string `bson:"actionType" json:"actionType"`

	// Result of the action.
	Result ActionResult `bson:"result" json:"result"`

	// --- Entity-specific detail fields (populated based on Category) ---

	// EntityID is the database key of the deleted entity.
	EntityID string `bson:"entityId" json:"entityId"`

	// EntityType further qualifies the entity (e.g. "APPROVAL_REQUEST", "OWNER", "USER_SESSION_LOG").
	EntityType string `bson:"entityType" json:"entityType"`

	// ProjectID for project-scoped entities (tasks, roles).
	ProjectID string `bson:"projectId,omitempty" json:"projectId,omitempty"`

	// ProjectName at the time of deletion (denormalized for historical reference).
	ProjectName string `bson:"projectName,omitempty" json:"projectName,omitempty"`

	// Additional context — human-readable reason when an action is skipped/retained/failed.
	Reason string `bson:"reason,omitempty" json:"reason,omitempty"`

	// SnapshotJSON stores a JSON-serialized snapshot of the entity *before* deletion.
	// This ensures we have a complete record of what was removed, even after the original
	// entity no longer exists.
	SnapshotJSON string `bson:"snapshotJson,omitempty" json:"snapshotJson,omitempty"`

	// Timestamp of when the action was executed.
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`

	// IPAddress of the requester (captured at REST boundary for traceability).
	IPAddress string `bson:"ipAddress,omitempty" json:"ipAddress,omitempty"`

	// RequestID ties back to the HTTP request for correlation with application logs.
	RequestID string `bson:"requestId,omitempty" json:"requestId,omitempty"`
}

// NewDeletionAuditEntry creates a new audit entry with a generated UUID key and current timestamp.
func NewDeletionAuditEntry(performedBy, targetUser, operationID string, category ActionCategory, actionType string, result ActionResult) *DeletionAuditEntry {
	return &DeletionAuditEntry{
		RootEntity:  domain.NewRootEntity(),
		PerformedBy: performedBy,
		TargetUser:  targetUser,
		OperationID: operationID,
		Category:    category,
		ActionType:  actionType,
		Result:      result,
		Timestamp:   time.Now(),
	}
}

// NewOperationID generates a unique identifier to group related deletion actions.
func NewOperationID() string {
	return uuid.NewString()
}

// --- DTO for API responses ---

type DeletionAuditEntryDto struct {
	Key         string         `json:"_key"`
	PerformedBy string         `json:"performedBy"`
	TargetUser  string         `json:"targetUser"`
	OperationID string         `json:"operationId"`
	Category    ActionCategory `json:"category"`
	ActionType  string         `json:"actionType"`
	Result      ActionResult   `json:"result"`
	EntityID    string         `json:"entityId"`
	EntityType  string         `json:"entityType"`
	ProjectID   string         `json:"projectId,omitempty"`
	ProjectName string         `json:"projectName,omitempty"`
	Reason      string         `json:"reason,omitempty"`
	Timestamp   time.Time      `json:"timestamp"`
	IPAddress   string         `json:"ipAddress,omitempty"`
	RequestID   string         `json:"requestId,omitempty"`
}

func (e *DeletionAuditEntry) ToDto() DeletionAuditEntryDto {
	return DeletionAuditEntryDto{
		Key:         e.Key,
		PerformedBy: e.PerformedBy,
		TargetUser:  e.TargetUser,
		OperationID: e.OperationID,
		Category:    e.Category,
		ActionType:  e.ActionType,
		Result:      e.Result,
		EntityID:    e.EntityID,
		EntityType:  e.EntityType,
		ProjectID:   e.ProjectID,
		ProjectName: e.ProjectName,
		Reason:      e.Reason,
		Timestamp:   e.Timestamp,
		IPAddress:   e.IPAddress,
		RequestID:   e.RequestID,
	}
}

// ToDtos converts a slice of entries to DTOs.
func ToDtos(entries []*DeletionAuditEntry) []DeletionAuditEntryDto {
	dtos := make([]DeletionAuditEntryDto, 0, len(entries))
	for _, e := range entries {
		dtos = append(dtos, e.ToDto())
	}
	return dtos
}

// DeletionAuditSummaryDto provides a grouped view of a deletion operation.
type DeletionAuditSummaryDto struct {
	OperationID  string                  `json:"operationId"`
	PerformedBy  string                  `json:"performedBy"`
	TargetUser   string                  `json:"targetUser"`
	Timestamp    time.Time               `json:"timestamp"`
	TotalActions int                     `json:"totalActions"`
	Succeeded    int                     `json:"succeeded"`
	Skipped      int                     `json:"skipped"`
	Failed       int                     `json:"failed"`
	Retained     int                     `json:"retained"`
	TaskEntries  []DeletionAuditEntryDto `json:"taskEntries"`
	RoleEntries  []DeletionAuditEntryDto `json:"roleEntries"`
	TraceEntries []DeletionAuditEntryDto `json:"traceEntries"`
	ProfileEntry []DeletionAuditEntryDto `json:"profileEntries"`
}
