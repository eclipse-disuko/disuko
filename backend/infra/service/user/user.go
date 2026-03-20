// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package user

import (
	"fmt"
	"slices"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain/approval"
	"github.com/eclipse-disuko/disuko/domain/label"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/user"
	"github.com/eclipse-disuko/disuko/helper/roles"
	approvallistRepo "github.com/eclipse-disuko/disuko/infra/repository/approvallist"
	labelRepo "github.com/eclipse-disuko/disuko/infra/repository/labels"
	projectRepo "github.com/eclipse-disuko/disuko/infra/repository/project"
	userRepo "github.com/eclipse-disuko/disuko/infra/repository/user"
	"github.com/eclipse-disuko/disuko/logy"
)

type Service struct {
	rs                     *logy.RequestSession
	userRepository         userRepo.IUsersRepository
	approvalListRepository approvallistRepo.IApprovalListRepository
	projectRepository      projectRepo.IProjectRepository
	labelRepository        labelRepo.ILabelRepository

	prLookupCache           map[string]*project.Project
	approvalListLookupCache map[string]*approval.ApprovalList
	fullNameLookupCache     map[string]string
}

type TaskEnriched struct {
	Task     user.Task
	Pr       *project.Project
	Approval *approval.Approval
}

type ProjectRole struct {
	Pr     *project.Project
	PrUser *project.ProjectMemberEntity
}

func Init(rs *logy.RequestSession, userRepo userRepo.IUsersRepository, approvallistRepo approvallistRepo.IApprovalListRepository, projectRepo projectRepo.IProjectRepository, labelRepo labelRepo.ILabelRepository) *Service {
	return &Service{
		rs:                     rs,
		userRepository:         userRepo,
		approvalListRepository: approvallistRepo,
		projectRepository:      projectRepo,
		labelRepository:        labelRepo,

		prLookupCache:           make(map[string]*project.Project),
		approvalListLookupCache: make(map[string]*approval.ApprovalList),
		fullNameLookupCache:     make(map[string]string),
	}
}

func (s *Service) GetProjectType(pr *project.Project) string {
	platformLabels := []struct {
		name         string
		platformType string
	}{
		{label.VEHICLE_PLATFORM, string(project.PlatformVehicle)},
		{label.MOBILE_PLATFORM, string(project.PlatformMobile)},
		{label.ENTERPRISE_PLATFORM, string(project.PlatformEnterprise)},
		{label.OTHER_PLATFORM, string(project.PlatformOther)},
	}

	for _, pl := range platformLabels {
		labelEntity := s.labelRepository.FindByNameAndType(s.rs, pl.name, label.POLICY)
		if labelEntity != nil && slices.Contains(pr.PolicyLabels, labelEntity.GetKey()) {
			return pl.platformType
		}
	}

	return ""
}

func (s *Service) EnrichTasks(tasks []user.Task) []TaskEnriched {
	var res []TaskEnriched
	for _, t := range tasks {
		omit, enriched := s.EnrichTask(t)
		if omit {
			continue
		}
		res = append(res, enriched)
	}
	return res
}

func (s *Service) EnrichTask(t user.Task) (bool, TaskEnriched) {
	pr, cached := s.prLookupCache[t.ProjectGuid]
	if !cached {
		pr = s.projectRepository.FindByKey(s.rs, t.ProjectGuid, false)
		if pr == nil {
			logy.Warnf(s.rs, "Omitting task %s because project %s got deleted!",
				t.Key,
				t.ProjectGuid,
			)
			return true, TaskEnriched{}
		}
		s.prLookupCache[pr.Key] = pr
	}

	al, cached := s.approvalListLookupCache[pr.Key]
	if !cached {
		al = s.approvalListRepository.FindByKey(s.rs, pr.Key, false)
		if al == nil {
			logy.Warnf(s.rs, "Omitting task %s because approvallist %s could not be found for project %s!",
				t.Key,
				t.TargetGuid,
				t.ProjectGuid,
			)
			return true, TaskEnriched{}
		}
		s.approvalListLookupCache[pr.Key] = al
	}

	var matchingApproval approval.Approval
	found := false
	for _, a := range al.Approvals {
		if t.TargetGuid == a.Key {
			matchingApproval = a
			found = true
		}
	}

	if !found {
		logy.Warnf(s.rs, "Omitting task %s because approval %s could not be found in project %s!",
			t.Key,
			t.TargetGuid,
			t.ProjectGuid,
		)
		return true, TaskEnriched{}
	}

	return false, TaskEnriched{
		Task:     t,
		Pr:       pr,
		Approval: &matchingApproval,
	}
}

func (s *Service) FullName(id string) string {
	full, cached := s.fullNameLookupCache[id]
	if !cached {
		u := s.userRepository.FindByUserId(s.rs, id)
		if u == nil {
			return "-"
		}
		full = fmt.Sprintf("%s %s", u.Forename, u.Lastname)
		s.fullNameLookupCache[id] = full
	}
	return full
}

func (s *Service) Roles(user *user.User) []ProjectRole {
	prs := s.projectRepository.FindAllForUser(s.rs, user.User)
	var res []ProjectRole
	for _, p := range prs {
		for _, m := range p.UserManagement.Users {
			if m.UserId == user.User {
				res = append(res, ProjectRole{
					Pr:     p,
					PrUser: m,
				})
				break
			}
		}
	}
	return res
}

func (s *Service) DelegateTask(taskId string, delegateUserId string, currentUserId string) error {
	currentUser := s.userRepository.FindByUserId(s.rs, currentUserId)
	if currentUser == nil {
		return fmt.Errorf("current user not found")
	}

	// Check if user has tasks (could be their own or pool user tasks for FOSS Office User)
	var targetTask *user.Task
	var taskOwner *user.User // Track which user owns the task

	// First check in current user's tasks
	for i := range currentUser.Tasks {
		if currentUser.Tasks[i].Key == taskId {
			targetTask = &currentUser.Tasks[i]
			taskOwner = currentUser
			break
		}
	}

	// If not found and user is from FOSS Office, check pool user's tasks
	if targetTask == nil && slices.Contains(currentUser.Roles, roles.FossOfficeUser) && len(conf.Config.Server.FOSSOfficeUserId) > 0 {
		// The task might belong to the pool user that was merged into the list
		// We need to find it in the pool user's actual task list
		poolUser := s.userRepository.FindByUserId(s.rs, conf.Config.Server.FOSSOfficeUserId)
		if poolUser != nil {
			for i := range poolUser.Tasks {
				if poolUser.Tasks[i].Key == taskId {
					targetTask = &poolUser.Tasks[i]
					taskOwner = poolUser
					break
				}
			}
		}
	}

	if targetTask == nil {
		return fmt.Errorf("task not found")
	}

	if targetTask.Status != user.TaskActive {
		return fmt.Errorf("only active tasks can be delegated")
	}

	// Verify this is a PLAUSIBILITY approval task
	approvalList := s.approvalListRepository.FindByKey(s.rs, targetTask.ProjectGuid, false)
	if approvalList == nil {
		return fmt.Errorf("approval list not found")
	}

	var targetApproval *approval.Approval
	for i := range approvalList.Approvals {
		if approvalList.Approvals[i].Key == targetTask.TargetGuid {
			targetApproval = &approvalList.Approvals[i]
			break
		}
	}

	if targetApproval == nil {
		return fmt.Errorf("approval not found")
	}

	if targetApproval.Type != approval.TypePlausibility {
		return fmt.Errorf("only PLAUSIBILITY type tasks can be delegated")
	}

	// Verify delegate user exists
	delegateUser := s.userRepository.FindByUserId(s.rs, delegateUserId)
	if delegateUser == nil {
		return fmt.Errorf("delegate user not found")
	}

	// Update the task - use the task owner (either current user or pool user)
	targetTask.DelegatedTo = delegateUserId
	s.userRepository.Update(s.rs, taskOwner)

	creatorUser := s.userRepository.FindByUserId(s.rs, targetTask.Creator)
	if creatorUser != nil {
		for i := range creatorUser.Tasks {
			if creatorUser.Tasks[i].TargetGuid == targetTask.TargetGuid &&
				creatorUser.Tasks[i].Type == user.ApprovalInfo &&
				creatorUser.Tasks[i].Status == user.TaskPending {
				creatorUser.Tasks[i].DelegatedTo = delegateUserId
				s.userRepository.Update(s.rs, creatorUser)
				break
			}
		}
	}

	return nil
}

func (s *Service) IsProjectMemberInPendingApprovalOrRequestUser(requestSession *logy.RequestSession, user *user.User, approvalList *approval.ApprovalList) string {
	if user != nil {
		for _, task := range user.Tasks {
			for _, appr := range approvalList.Approvals {
				if appr.Key == task.TargetGuid {
					approvalStatus := appr.ToApprovalDtoStatus()
					if approval.Pending == approvalStatus {
						return "PROJECT_TASK_APPROVAL_PENDING"
					}
				}
			}
		}
	}
	return ""
}
