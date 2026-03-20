// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package user

import (
	"net/http"
	"time"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/approval"
)

type Task struct {
	domain.ChildEntity `bson:",inline"`

	Type TaskType

	TargetGuid                   string
	ProjectGuid                  string
	Creator                      string
	CreatorDepartment            string
	CreatorDepartmentDescription string
	DelegatedTo                  string

	Status TaskStatus
}

type TaskStatus string

const (
	TaskDone    = "DONE"
	TaskActive  = "ACTIVE"
	TaskPending = "PENDING"
)

type TaskType string

const (
	Approval     TaskType = "APPROVAL"
	ApprovalInfo TaskType = "APPROVALINFO"
)

type TaskDto struct {
	Id                           string                `json:"id" validate:"lte=36"`
	Created                      time.Time             `json:"created" validate:"lte=36"`
	Updated                      time.Time             `json:"updated" validate:"lte=36"`
	TargetGuid                   string                `json:"approvalGuid" validate:"lte=36"`
	TargetType                   approval.ApprovalType `json:"approvalType" validate:"lte=36"`
	Creator                      string                `json:"creator" validate:"lte=200"`
	CreatorFullName              string                `json:"creatorFullName" validate:"lte=200"`
	CreatorDepartment            string                `json:"creatorDepartment" validate:"lte=200"`
	CreatorDepartmentDescription string                `json:"creatorDepartmentDescription" validate:"lte=200"`
	DelegatedTo                  string                `json:"delegatedTo" validate:"lte=200"`
	DelegatedToFullName          string                `json:"delegatedToFullName" validate:"lte=200"`
	ProjectGuid                  string                `json:"projectGuid" validate:"lte=36"`
	ProjectName                  string                `json:"projectName" validate:"required,gte=3,lte=80"`
	IsProjectGroup               bool                  `json:"isProjectGroup"`
	ProjectType                  string                `json:"projectType" validate:"lte=50"`
	Status                       TaskStatus            `json:"status" validate:"lte=36"`
	ResultStatus                 approval.StateInfo    `json:"resultStatus" validate:"lte=36"`
	Type                         TaskType              `json:"type" validate:"lte=200"`
}

type DelegateTaskDto struct {
	DelegateUserId string `json:"delegateUserId" validate:"required"`
}

func (dto *DelegateTaskDto) Bind(r *http.Request) error {
	return nil
}

func (task Task) ToDto(creatorFullName string, targetType approval.ApprovalType, resultStatus approval.StateInfo, ProjectName string, IsProjectGroup bool, delegatedToFullName string, projectType string) TaskDto {
	return TaskDto{
		Id:                           task.GetKey(),
		Creator:                      task.Creator,
		CreatorFullName:              creatorFullName,
		CreatorDepartment:            task.CreatorDepartment,
		CreatorDepartmentDescription: task.CreatorDepartmentDescription,
		DelegatedTo:                  task.DelegatedTo,
		DelegatedToFullName:          delegatedToFullName,
		ProjectGuid:                  task.ProjectGuid,
		ProjectName:                  ProjectName,
		IsProjectGroup:               IsProjectGroup,
		ProjectType:                  projectType,
		Status:                       task.Status,
		Created:                      task.Created,
		Updated:                      task.Updated,
		TargetGuid:                   task.TargetGuid,
		TargetType:                   targetType,
		ResultStatus:                 resultStatus,
		Type:                         task.Type,
	}
}

func (user *User) AddApprovalTask(app approval.Approval, creatorUser *User) Task {
	task := Task{
		ChildEntity:                  domain.NewChildEntity(),
		TargetGuid:                   app.Key,
		Type:                         Approval,
		ProjectGuid:                  app.ProjectGuid,
		Creator:                      app.Creator,
		CreatorDepartment:            creatorUser.MetaData.Department,
		CreatorDepartmentDescription: creatorUser.MetaData.DepartmentDescription,
		Status:                       TaskActive,
	}
	user.Tasks = append(user.Tasks, task)
	return task
}

func (user *User) AddApprovalCreatorTask(app approval.Approval) Task {
	task := Task{
		ChildEntity:                  domain.NewChildEntity(),
		TargetGuid:                   app.Key,
		Type:                         ApprovalInfo,
		ProjectGuid:                  app.ProjectGuid,
		Creator:                      app.Creator,
		CreatorDepartment:            user.MetaData.Department,
		CreatorDepartmentDescription: user.MetaData.DepartmentDescription,
		Status:                       TaskPending,
	}
	user.Tasks = append(user.Tasks, task)
	return task
}
