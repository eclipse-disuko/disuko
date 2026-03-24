// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package user

import (
	"slices"

	"github.com/eclipse-disuko/disuko/domain/user"
	projectRepo "github.com/eclipse-disuko/disuko/infra/repository/project"
	userRepo "github.com/eclipse-disuko/disuko/infra/repository/user"
	"github.com/eclipse-disuko/disuko/logy"
)

type DeletionAction struct {
	ActionType  string `json:"action_type"`
	EntityID    string `json:"entity_id"`
	EntityType  string `json:"entity_type"`
	Status      string `json:"status"`
	Reason      string `json:"reason"`
	ProjectID   string `json:"project_id,omitempty"`
	ProjectName string `json:"project_name,omitempty"`
	TaskID      string `json:"task_id,omitempty"`
	TaskType    string `json:"task_type,omitempty"`
	Priority    string `json:"priority,omitempty"`
	RoleName    string `json:"role_name,omitempty"`
	TraceType   string `json:"trace_type,omitempty"`
}

type DeletionPlan struct {
	Username       string           `json:"username"`
	TotalActions   int              `json:"total_actions"`
	TaskActions    []DeletionAction `json:"task_actions"`
	RoleActions    []DeletionAction `json:"role_actions"`
	TraceActions   []DeletionAction `json:"trace_actions"`
	ProfileActions []DeletionAction `json:"profile_actions"`
	service        *DeletionService `json:"-"`
}

type DeletionService struct {
	rs                *logy.RequestSession
	userRepository    userRepo.IUsersRepository
	projectRepository projectRepo.IProjectRepository
}

func NewDeletionService(rs *logy.RequestSession, userRepository userRepo.IUsersRepository, projectRepository projectRepo.IProjectRepository) *DeletionService {
	return &DeletionService{
		rs:                rs,
		userRepository:    userRepository,
		projectRepository: projectRepository,
	}
}

func (s *DeletionService) ExecuteDeletion(username string) (*DeletionPlan, error) {
	plan := &DeletionPlan{
		Username:       username,
		TaskActions:    []DeletionAction{},
		RoleActions:    []DeletionAction{},
		TraceActions:   []DeletionAction{},
		ProfileActions: []DeletionAction{},
		service:        s,
	}

	plan.collectTaskDeletions()
	plan.collectRoleDeletions()
	plan.collectTraceDeletions()
	plan.collectProfileDeletion()

	plan.TotalActions = len(plan.TaskActions) + len(plan.RoleActions) + len(plan.TraceActions) + len(plan.ProfileActions)

	return plan, nil
}

func (plan *DeletionPlan) collectTaskDeletions() {
	mockTasks := []struct {
		taskID      string
		taskType    string
		projectName string
		projectID   string
		priority    string
	}{
		{"task-001", "APPROVAL_REQUEST", "Project Alpha", "proj-alpha-001", "HIGH"},
		{"task-002", "REVIEW_REQUEST", "Project Beta", "proj-beta-002", "MEDIUM"},
		{"task-003", "DOCUMENT_REVIEW", "Project Gamma", "proj-gamma-003", "LOW"},
		{"task-004", "RELEASE_APPROVAL", "Project Alpha", "proj-alpha-001", "CRITICAL"},
	}

	for _, task := range mockTasks {
		action := DeletionAction{
			ActionType:  "cancel_task",
			EntityID:    task.taskID,
			EntityType:  "task",
			Status:      plan.service.getActionStatus(),
			Reason:      plan.service.getActionReason(),
			ProjectID:   task.projectID,
			ProjectName: task.projectName,
			TaskID:      task.taskID,
			TaskType:    task.taskType,
			Priority:    task.priority,
		}
		plan.TaskActions = append(plan.TaskActions, action)

	}
}

func (plan *DeletionPlan) collectRoleDeletions() {
	mockRoles := []struct {
		roleID      string
		projectName string
		projectID   string
		roleName    string
		canDelete   bool
		reason      string
	}{
		{"role-001", "Project Alpha", "proj-alpha-001", "OWNER", false, "Cannot remove: user is Project Responsible Owner"},
		{"role-002", "Project Beta", "proj-beta-002", "CONTRIBUTOR", true, ""},
		{"role-003", "Project Gamma", "proj-gamma-003", "OWNER", false, "Cannot remove: user is the last Owner"},
		{"role-004", "Project Delta", "proj-delta-004", "VIEWER", true, ""},
		{"role-005", "Project Epsilon", "proj-epsilon-005", "REVIEWER", true, ""},
		{"role-006", "Project Zeta", "proj-zeta-006", "APPROVER", false, "Cannot remove: active approvals pending"},
	}

	for _, role := range mockRoles {
		status := "planned"
		reason := ""
		if !role.canDelete {
			status = "skipped"
			reason = role.reason
		}

		action := DeletionAction{
			ActionType:  "remove_role",
			EntityID:    role.roleID,
			EntityType:  "role",
			Status:      status,
			Reason:      reason,
			ProjectID:   role.projectID,
			ProjectName: role.projectName,
			RoleName:    role.roleName,
		}
		plan.RoleActions = append(plan.RoleActions, action)

	}
}

func (plan *DeletionPlan) collectTraceDeletions() {
	mockTraces := []struct {
		traceID      string
		traceType    string
		shouldRetain bool
		reason       string
	}{
		{"trace-001", "PROJECT_AUDIT_LOG", true, "Retained: Data-object audit log kept for compliance"},
		{"trace-002", "APPROVAL_REQUEST", true, "Retained: Approval record kept for traceability"},
		{"trace-003", "REVIEW_REMARK", true, "Retained: Review remark kept for history"},
		{"trace-004", "USER_SESSION_LOG", false, ""},
		{"trace-005", "DISCLOSURE_DOC_REQUEST", true, "Retained: Disclosure request kept for legal reasons"},
		{"trace-006", "EXPORT_HISTORY", false, ""},
		{"trace-007", "LOGIN_ACTIVITY", false, ""},
		{"trace-008", "DOCUMENT_ACCESS_LOG", true, "Retained: Access log required for security audit"},
	}

	for _, trace := range mockTraces {
		actionType := "delete_trace"
		status := "planned"
		reason := trace.reason

		if trace.shouldRetain {
			actionType = "anonymize_trace"
			status = "skipped"
		}

		action := DeletionAction{
			ActionType: actionType,
			EntityID:   trace.traceID,
			EntityType: "trace",
			Status:     status,
			Reason:     reason,
			TraceType:  trace.traceType,
		}
		plan.TraceActions = append(plan.TraceActions, action)

	}
}

func (plan *DeletionPlan) collectProfileDeletion() {
	canDelete := len(plan.TaskActions) == 0 && plan.service.allRolesRemovable(plan.RoleActions)

	status := "planned"
	reason := ""
	if !canDelete {
		status = "skipped"
		reason = "Cannot delete profile: tasks or non-removable roles still exist"
	}

	action := DeletionAction{
		ActionType: "delete_profile",
		EntityID:   plan.Username,
		EntityType: "profile",
		Status:     status,
		Reason:     reason,
	}
	plan.ProfileActions = append(plan.ProfileActions, action)

}

func (s *DeletionService) getActionStatus() string {
	return "completed"
}

func (s *DeletionService) getActionReason() string {
	return ""
}

func (s *DeletionService) allRolesRemovable(roleActions []DeletionAction) bool {
	return !slices.ContainsFunc(roleActions, func(action DeletionAction) bool {
		return action.Status == "skipped"
	})
}

func (s *DeletionService) ExecuteTaskDeletion(username string, taskID string) error {
	return nil
}

func (s *DeletionService) ExecuteRoleDeletion(username string, projectID string, roleID string) error {
	return nil
}

func (s *DeletionService) ExecuteTraceDeletion(username string, traceID string) error {
	return nil
}

func (s *DeletionService) GetUserDeletionStats(u *user.User) (taskCount, roleCount, traceCount int) {
	taskCount = len(u.Tasks)
	roleCount = 5
	traceCount = 8

	return taskCount, roleCount, traceCount
}
