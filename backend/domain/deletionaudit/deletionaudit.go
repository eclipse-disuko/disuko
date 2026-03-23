// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package deletionaudit

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/google/uuid"
)

type ActionCategory string

const (
	CategoryTask    ActionCategory = "TASK"
	CategoryRole    ActionCategory = "ROLE"
	CategoryTrace   ActionCategory = "TRACE"
	CategoryProfile ActionCategory = "PROFILE"
)

type ActionResult string

const (
	ResultSuccess  ActionResult = "SUCCESS"
	ResultSkipped  ActionResult = "SKIPPED"
	ResultFailed   ActionResult = "FAILED"
	ResultRetained ActionResult = "RETAINED"
	ResultPlanned  ActionResult = "PLANNED"
	ResultDryRun   ActionResult = "DRY_RUN"
)

type DeletionAuditEntry struct {
	domain.RootEntity `bson:",inline"`
	PerformedBy       string
	TargetUser        string
	OperationID       string
	Category          ActionCategory
	ActionType        string
	Result            ActionResult
	EntityID          string
	EntityType        string
	ProjectID         string
	ProjectName       string
	Reason            string
	SnapshotJSON      string
	Timestamp         time.Time
	IPAddress         string
	RequestID         string
}

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

func NewOperationID() string {
	return uuid.NewString()
}

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

func ToDtos(entries []*DeletionAuditEntry) []DeletionAuditEntryDto {
	dtos := make([]DeletionAuditEntryDto, 0, len(entries))
	for _, e := range entries {
		dtos = append(dtos, e.ToDto())
	}
	return dtos
}

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
