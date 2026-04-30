// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package approval

import (
	"slices"
	"time"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain/project/approvable"
	"github.com/eclipse-disuko/disuko/helper/roles"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/approval"
	"github.com/eclipse-disuko/disuko/domain/audit"
	"github.com/eclipse-disuko/disuko/domain/overallreview"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/user"
	auditHelper "github.com/eclipse-disuko/disuko/helper/audit"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/observermngmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func (s *ApprovalService) CreatePlausibilityCheck(pr *project.Project, req approval.RequestPlausibilityCheckDto, creator string) string {
	info := s.getApprovalInfo(pr, nil, true)
	appr := approval.Approval{
		ChildEntity: domain.SetChildEntity(uuid.New().String()),
		Type:        approval.TypePlausibility,
		ProjectGuid: pr.Key,
		Creator:     creator,
		Comment:     req.Comment,
		Info:        info,

		Plausibility: approval.PlausibilityCheck{
			Approver: req.Approver,
		},

		DocumentFlags: approval.TaskMetaDocument{
			C1: req.MetaDoc.C1,
			C2: req.MetaDoc.C2,
			C3: req.MetaDoc.C3,
			C4: req.MetaDoc.C4,
			C5: req.MetaDoc.C5,
			C6: req.MetaDoc.C6,
		},
	}

	s.createApprovalCreatorTask(&appr)
	s.addReviewerToProject(req, pr, creator)
	s.createPlausibilityApproverTask(&appr)
	s.addTargetUsers(&appr, approvable.APPROVAL_TYPE_PLAUSI, pr, creator)

	if appr.Plausibility.Approver == conf.Config.Server.FOSSOfficeUserId {
		s.addReviewSubscriber(pr, creator)
	}

	s.AuditLogListRepo.AddStaticAuditEntryByKey(s.RequestSession, pr.Key, creator, message.ReviewCreated, appr.ToAudit())
	// auditHelper.CreateAndAddAuditEntry(&res.Container, creator, message.ReviewCreated, audit.DiffWithReporter, res.ToAudit(), approval.ApprovalAudit{})

	approvalList := s.ApprovalListRepo.FindByKey(s.RequestSession, pr.Key, false)
	if approvalList == nil {
		s.ApprovalListRepo.Save(s.RequestSession, &approval.ApprovalList{
			RootEntity: domain.NewRootEntityWithKey(pr.Key),
			Approvals:  []approval.Approval{appr},
		})
	} else {
		approvalList.Approvals = append(approvalList.Approvals, appr)
		s.ApprovalListRepo.Update(s.RequestSession, approvalList)
	}

	s.markSbomIsInUse(info.Projects)
	return appr.Key
}

func (s *ApprovalService) addReviewSubscriber(pr *project.Project, creator string) {
	m := pr.GetMember(creator)
	if m == nil {
		return
	}
	m.Subscriptions.OverallReview = true
}

func (s *ApprovalService) createPlausibilityApproverTask(app *approval.Approval) {
	// before := app.ToAudit()
	targetUser := s.UserRepo.FindByUserId(s.RequestSession, app.Plausibility.Approver)
	beforeTask := targetUser.ToUserAudit()

	creatorUser := s.UserRepo.FindByUserId(s.RequestSession, app.Creator)
	task := targetUser.AddApprovalTask(*app, creatorUser)
	now := time.Now()
	app.Plausibility.State = approval.ApproveState{
		State:   approval.Pending,
		Updated: &now,
	}
	afterTask := targetUser.ToUserAudit()
	auditHelper.CreateAndAddAuditEntry(&targetUser.Container, app.Creator, message.ApprovalTaskCreate, audit.DiffWithReporter, afterTask, beforeTask)

	s.UserRepo.Update(s.RequestSession, targetUser)
	if targetUser.Email != "" {
		observermngmt.FireEvent(observermngmt.ApprovalTaskCreated, observermngmt.ApprovalTaskData{
			RequestSession: s.RequestSession,
			Type:           approval.TypePlausibility,
			Creator:        app.Creator,
			TargetUser:     app.Plausibility.Approver,
			Comment:        app.Comment,
			ProjectId:      app.ProjectGuid,
			TaskId:         task.Key,
			Approvables:    app.Info.Projects,
		})
	}
	// auditHelper.CreateAndAddAuditEntry(&app.Container, app.Plausibility.Approver, message.ApprovalTaskCreate, audit.DiffWithReporter, app.ToAudit(), before)
}

func (s *ApprovalService) processPlausibilityCheckUpdate(pr *project.Project, targetApproval *approval.Approval, username string, req approval.UpdateApprovalDto) {
	valid, stateInfo := approval.ParseStateInfo(req.State)
	if !valid {
		exception.ThrowExceptionBadRequestResponse()
	}

	if stateInfo == approval.Declined && req.Comment == "" {
		exception.ThrowExceptionBadRequestResponse()
	}
	userDetails := s.UserRepo.FindByUserId(s.RequestSession, username)

	before := targetApproval.ToAudit()
	delegationUser := s.delegationUser(targetApproval.Key, targetApproval.Plausibility)
	if stateInfo == approval.Aborted {
		if !targetApproval.Plausibility.IsActive() {
			exception.ThrowExceptionBadRequestResponse()
		}
		if targetApproval.Creator != username && (!slices.Contains(userDetails.Roles, roles.FossOfficeUser)) {
			exception.ThrowExceptionSendDeniedResponse()
		}
		s.deletePending(targetApproval)
		s.setTaskDone(targetApproval.Creator, targetApproval, user.ApprovalInfo, user.TaskPending)
		targetApproval.Plausibility.Aborted = true
		s.AuditLogListRepo.CreateAuditEntryByKey(s.RequestSession, pr.Key, username, message.ReviewAborted, audit.DiffWithReporter, targetApproval.ToAudit(), before)
		observermngmt.FireEvent(observermngmt.ApprovalFinalized, observermngmt.ApprovalData{
			RequestSession: s.RequestSession,
			Approval:       targetApproval,
			DelegatedTo:    delegationUser,
		})
		return
	}

	approver := targetApproval.Plausibility.GetPendingApprover()
	if approver != username &&
		!(slices.Contains(userDetails.Roles, roles.FossOfficeUser) && approver == conf.Config.Server.FOSSOfficeUserId) {
		exception.ThrowExceptionSendDeniedResponse()
	}

	targetApproval.Plausibility.ApproveComment = req.Comment
	now := time.Now()
	if stateInfo == approval.Approved {
		targetApproval.Plausibility.State = approval.ApproveState{
			State:   approval.Approved,
			Updated: &now,
		}
	} else {
		targetApproval.Plausibility.State = approval.ApproveState{
			State:   approval.Declined,
			Updated: &now,
		}
	}

	s.setTaskDone(approver, targetApproval, user.Approval, user.TaskActive)
	s.setTaskDone(targetApproval.Creator, targetApproval, user.ApprovalInfo, user.TaskPending)
	s.AuditLogListRepo.CreateAuditEntryByKey(s.RequestSession, pr.Key, username, message.ReviewUpdated, audit.DiffWithReporter, targetApproval.ToAudit(), before)

	// Add overall review entry for vehicle projects when approved by FOSSOFFICE user
	if slices.Contains(userDetails.Roles, roles.FossOfficeUser) && s.ProjectLabelService.HasVehiclePlatformLabel(s.RequestSession, pr) {
		s.createOverallReviewForApproval(stateInfo, pr, targetApproval, username)
	}

	observermngmt.FireEvent(observermngmt.ApprovalFinalized, observermngmt.ApprovalData{
		RequestSession: s.RequestSession,
		Approval:       targetApproval,
		DelegatedTo:    delegationUser,
	})
	// auditHelper.CreateAndAddAuditEntry(&targetApproval.Container, username, message.ReviewUpdated, audit.DiffWithReporter, targetApproval.ToAudit(), before)
}

func (s *ApprovalService) adminAbortPlausibility(pr *project.Project, targetApproval *approval.Approval) {
	if !targetApproval.Plausibility.IsActive() {
		return
	}
	before := targetApproval.ToAudit()
	delegationUser := s.delegationUser(targetApproval.Key, targetApproval.Plausibility)
	s.deletePending(targetApproval)
	s.setTaskDone(targetApproval.Creator, targetApproval, user.ApprovalInfo, user.TaskPending)
	targetApproval.Plausibility.Aborted = true
	s.AuditLogListRepo.CreateAuditEntryByKey(s.RequestSession, pr.Key, "SYSTEM", message.ReviewAborted, audit.DiffWithReporter, targetApproval.ToAudit(), before)
	observermngmt.FireEvent(observermngmt.ApprovalFinalized, observermngmt.ApprovalData{
		RequestSession: s.RequestSession,
		Approval:       targetApproval,
		DelegatedTo:    delegationUser,
	})
}

func (s *ApprovalService) delegationUser(appKey string, plausi approval.PlausibilityCheck) string {
	targetUser := s.UserRepo.FindByUserId(s.RequestSession, plausi.Approver)
	task := targetUser.GetTask(appKey, user.Approval, user.TaskActive)
	if task == nil {
		return ""
	}
	return task.DelegatedTo
}

func (s *ApprovalService) addReviewerToProject(req approval.RequestPlausibilityCheckDto, pr *project.Project, creator string) {
	if pr.GetMember(req.Approver) != nil {
		return
	}
	pr.AddUser(project.ProjectMemberRequestDto{
		TargetUser: req.Approver,
		UserType:   project.VIEWER,
	})
	member := pr.GetMember(req.Approver)
	userAudit := member.ToAudit(s.RequestSession, s.UserRepo)
	s.AuditLogListRepo.CreateAuditEntryByKey(s.RequestSession, pr.Key, creator, message.ProjectUserCreated, cmp.Diff, userAudit, project.ProjectMemberAudit{})
	s.ProjectRepo.Update(s.RequestSession, pr)
}

func (s *ApprovalService) createOverallReviewForApproval(state approval.StateInfo, pr *project.Project, targetApproval *approval.Approval, username string) {
	if s.OverallReviewService == nil {
		return
	}

	for _, projectApprovable := range targetApproval.Info.Projects {
		if projectApprovable.ApprovableSPDX.VersionKey == "" || projectApprovable.ApprovableSPDX.SpdxKey == "" {
			continue
		}

		var targetProject *project.Project
		if projectApprovable.ProjectKey == pr.Key {
			targetProject = pr
		} else {
			targetProject = s.ProjectRepo.FindByKey(s.RequestSession, projectApprovable.ProjectKey, false)
			if targetProject == nil {
				continue
			}
		}

		version := targetProject.GetVersion(projectApprovable.ApprovableSPDX.VersionKey)
		if version == nil {
			continue
		}

		if state == approval.Approved {
			s.OverallReviewService.AddToProject(
				s.RequestSession,
				targetProject,
				version,
				username,
				overallreview.Audited,
				message.FossOfficeCommentApproved,
				projectApprovable.ApprovableSPDX.SpdxKey,
			)
		} else {
			s.OverallReviewService.AddToProject(
				s.RequestSession,
				targetProject,
				version,
				username,
				overallreview.AcceptableAfterChanges,
				message.FossOfficeCommentInvest,
				projectApprovable.ApprovableSPDX.SpdxKey,
			)
		}

		if targetProject.Key != pr.Key {
			s.ProjectRepo.Update(s.RequestSession, targetProject)
		}
	}
}
