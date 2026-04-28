// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package approval

import (
	"fmt"
	"strconv"
	"time"

	"github.com/eclipse-disuko/disuko/domain/project/approvable"

	"github.com/eclipse-disuko/disuko/observermngmt"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/approval"
	"github.com/eclipse-disuko/disuko/domain/audit"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/user"
	auditHelper "github.com/eclipse-disuko/disuko/helper/audit"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func (s *ApprovalService) CreateInternalApproval(pr *project.Project, req approval.RequestInternalApprovalDto, creator string) string {
	if pr.HasParent() || s.ProjectLabelService.HasVehiclePlatformLabel(s.RequestSession, pr) {
		exception.ThrowExceptionBadRequestResponse()
	}

	info := s.getApprovalInfo(pr, nil, false)
	if len(info.Projects) == 0 {
		exception.ThrowExceptionBadRequestResponse()
	}
	appr := approval.Approval{
		ChildEntity: domain.SetChildEntity(uuid.New().String()),
		Type:        approval.TypeInternal,
		ProjectGuid: pr.Key,
		Creator:     creator,
		Comment:     req.Comment,
		Info:        info,

		Internal: approval.InternalApproval{
			ApproveStates: [4]approval.ApproveState{},
			Approver: [4]string{
				approval.Supplier1: req.SupplierApprover1,
				approval.Supplier2: req.SupplierApprover2,
				approval.Customer1: req.CustomerApprover1,
				approval.Customer2: req.CustomerApprover2,
			},
			ApproveComments: [4]string{},
			Generating:      true,
		},

		DocumentFlags: approval.TaskMetaDocument{
			C1: req.MetaDoc.C1,
			C2: req.MetaDoc.C2,
			C3: req.MetaDoc.C3,
			C4: req.MetaDoc.C4,
			C5: req.MetaDoc.C5,
			C6: pr.IsNoFoss,
		},
	}

	s.setProjectApprovers(&appr, pr)
	s.addTargetUsers(&appr, approvable.APPROVAL_TYPE_INTERNAL, pr, creator)

	s.createApprovalCreatorTask(&appr)
	s.setApproverRole(&appr, approval.Supplier1)
	s.setApproverRole(&appr, approval.Supplier2)

	s.setApprovalSpdxStatus(&appr, approval.Pending)

	s.AuditLogListRepo.AddStaticAuditEntryByKey(s.RequestSession, pr.Key, creator, message.InternalApprovalCreated, appr.ToAudit())
	// auditHelper.CreateAndAddAuditEntry(&res.Container, creator, message.InternalApprovalCreated, audit.DiffWithReporter, res.ToAudit(), approval.ApprovalAudit{})

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

func (s *ApprovalService) processInternalApprovalUpdate(pr *project.Project, targetApproval *approval.Approval, username string, req approval.UpdateApprovalDto) {
	if !targetApproval.Internal.IsActive() {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorDbNotFound), "")
	}

	if targetApproval.Internal.Generating {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorDbNotFound), "")
	}

	valid, stateInfo := approval.ParseStateInfo(req.State)
	if !valid {
		exception.ThrowExceptionBadRequestResponse()
	}

	if stateInfo == approval.Declined && req.Comment == "" {
		exception.ThrowExceptionBadRequestResponse()
	}

	before := targetApproval.ToAudit()
	if stateInfo == approval.Aborted {
		if !targetApproval.Internal.IsActive() {
			exception.ThrowExceptionBadRequestResponse()
		}
		if targetApproval.Creator != username {
			exception.ThrowExceptionSendDeniedResponse()
		}
		s.deletePending(targetApproval)
		s.setTaskDone(targetApproval.Creator, targetApproval, user.ApprovalInfo, user.TaskPending)
		targetApproval.Internal.Aborted = true
		s.setApprovalSpdxStatus(targetApproval, approval.Aborted)
		s.AuditLogListRepo.CreateAuditEntryByKey(s.RequestSession, pr.Key, username, message.InternalApprovalAborted, audit.DiffWithReporter, targetApproval.ToAudit(), before)
		observermngmt.FireEvent(observermngmt.ApprovalFinalized, observermngmt.ApprovalData{
			RequestSession: s.RequestSession,
			Approval:       targetApproval,
		})
		return
	}

	approverRole := targetApproval.Internal.FirstPendingApproverRole(username)
	if approverRole == approval.None {
		exception.ThrowExceptionSendDeniedResponse()
	}
	s.setTaskDone(username, targetApproval, user.Approval, user.TaskActive)
	if stateInfo == approval.Approved {
		now := time.Now()
		targetApproval.Internal.ApproveStates[approverRole] = approval.ApproveState{
			State:   approval.Approved,
			Updated: &now,
		}
		targetApproval.Internal.ApproveComments[approverRole] = req.Comment
		if (approverRole == approval.Supplier1 || approverRole == approval.Supplier2) && targetApproval.Internal.SupplierDone() {
			if targetApproval.Internal.Approver[approval.Customer1] != "" {
				s.setApproverRole(targetApproval, approval.Customer1)
				s.CreateInternalApprovalTask(targetApproval, approval.Customer1)
			}
			if targetApproval.Internal.Approver[approval.Customer2] != "" {
				s.setApproverRole(targetApproval, approval.Customer2)
				s.CreateInternalApprovalTask(targetApproval, approval.Customer2)
			}
		}
		if (approverRole == approval.Customer1 || approverRole == approval.Customer2) && targetApproval.Internal.CustomerDone() {
			s.setApprovalSpdxStatus(targetApproval, approval.Approved)
			s.setTaskDone(targetApproval.Creator, targetApproval, user.ApprovalInfo, user.TaskPending)
			observermngmt.FireEvent(observermngmt.ApprovalFinalized, observermngmt.ApprovalData{
				RequestSession: s.RequestSession,
				Approval:       targetApproval,
			})
		}

	} else if stateInfo == approval.Declined {
		now := time.Now()
		targetApproval.Internal.ApproveStates[approverRole] = approval.ApproveState{
			State:   approval.Declined,
			Updated: &now,
		}
		targetApproval.Internal.ApproveComments[approverRole] = req.Comment
		s.deletePending(targetApproval)
		s.setApprovalSpdxStatus(targetApproval, approval.Declined)
		s.setTaskDone(targetApproval.Creator, targetApproval, user.ApprovalInfo, user.TaskPending)
		observermngmt.FireEvent(observermngmt.ApprovalFinalized, observermngmt.ApprovalData{
			RequestSession: s.RequestSession,
			Approval:       targetApproval,
		})
	}
	s.AuditLogListRepo.CreateAuditEntryByKey(s.RequestSession, pr.Key, username, message.InternalApprovalUpdated, audit.DiffWithReporter, targetApproval.ToAudit(), before)
	// auditHelper.CreateAndAddAuditEntry(&targetApproval.Container, username, message.InternalApprovalUpdated, audit.DiffWithReporter, targetApproval.ToAudit(), before)

	s.stampDisclosureDocument(pr, targetApproval, approverRole, req.PowerOfAttorney, stateInfo == approval.Approved)
}

func (s *ApprovalService) stampDisclosureDocument(pr *project.Project, app *approval.Approval, role approval.Approver, powerOfAttorney approval.PowerOfAttorneyType, accepted bool) {
	u := s.UserRepo.FindByUserId(s.RequestSession, app.Internal.GetApproverName(role))
	forename := u.Forename
	lastname := u.Lastname
	userId := app.Internal.GetApproverName(role)

	approverRoleText := ""
	if role == approval.Supplier1 || role == approval.Supplier2 {
		approverRoleText = "Supplier"
	} else if role == approval.Customer1 || role == approval.Customer2 {
		approverRoleText = "Customer"
	}

	nowFormatted := time.Now().Format("02.01.2006 15:04 MST")
	text := createDecisionText(accepted, powerOfAttorney, approverRoleText, forename, lastname, userId, nowFormatted)
	s.FOSSddService.SignDocuments(s.RequestSession, pr, app.Key, text, app.Internal.DocVersion)
	app.Internal.DocVersion++
}

func createDecisionText(accepted bool, powerOfAttorney approval.PowerOfAttorneyType, approverRoleText, forename, lastname, userId, date string) string {
	decision := "Accepted"
	if !accepted {
		decision = "Declined"
	}

	powerOfAttorneyText := ""
	switch powerOfAttorney {
	case approval.PowerOfAttorneyPPA:
		powerOfAttorneyText = "ppa. "
	case approval.PowerOfAttorneyIV:
		powerOfAttorneyText = "i.V. "
	}

	text := fmt.Sprintf("%s Decision: %s by %s%s %s (%s)\nDate: %s. Hash Chain: ",
		approverRoleText, decision, powerOfAttorneyText, forename, lastname, userId, date)

	return text
}

func (s *ApprovalService) setApproverRole(app *approval.Approval, role approval.Approver) {
	now := time.Now()
	app.Internal.ApproveStates[role] = approval.ApproveState{
		State:   approval.Pending,
		Updated: &now,
	}
}

func (s *ApprovalService) CreateInternalApprovalTask(app *approval.Approval, role approval.Approver) {
	if app.Internal.GetApproverName(role) == "" {
		return
	}
	targetUser := s.UserRepo.FindByUserId(s.RequestSession, app.Internal.GetApproverName(role))
	beforeTask := targetUser.ToUserAudit()

	creatorUser := s.UserRepo.FindByUserId(s.RequestSession, app.Creator)
	task := targetUser.AddApprovalTask(*app, creatorUser)
	afterTask := targetUser.ToUserAudit()
	auditHelper.CreateAndAddAuditEntry(&targetUser.Container, app.Creator, message.ApprovalTaskCreate, audit.DiffWithReporter, afterTask, beforeTask)
	s.UserRepo.Update(s.RequestSession, targetUser)
	if targetUser.Email != "" {
		observermngmt.FireEvent(observermngmt.ApprovalTaskCreated, observermngmt.ApprovalTaskData{
			RequestSession: s.RequestSession,
			Type:           approval.TypeInternal,
			Comment:        app.Comment,
			Creator:        app.Creator,
			TargetUser:     app.Internal.GetApproverName(role),
			ProjectId:      app.ProjectGuid,
			TaskId:         task.Key,
			Approvables:    app.Info.Projects,
		})
	}
}

func (s *ApprovalService) setProjectApprovers(targetApproval *approval.Approval, pr *project.Project) {
	pr.SupplierExtraData.FRI = targetApproval.Internal.GetApproverName(approval.Supplier1)
	pr.SupplierExtraData.SRI = targetApproval.Internal.GetApproverName(approval.Supplier2)

	cFRI := targetApproval.Internal.GetApproverName(approval.Customer1)
	if cFRI != "" {
		pr.CustomerMeta.FRI = cFRI
	}
	cSRI := targetApproval.Internal.GetApproverName(approval.Customer2)
	if cSRI != "" {
		pr.CustomerMeta.SRI = cSRI
	}
	s.ProjectRepo.Update(s.RequestSession, pr)
}

func (s *ApprovalService) addTargetUsers(targetApproval *approval.Approval, approvalType string, pr *project.Project, creator string) {
	users := []string{
		targetApproval.Internal.GetApproverName(approval.Supplier1),
		targetApproval.Internal.GetApproverName(approval.Supplier2),
		targetApproval.Internal.GetApproverName(approval.Customer1),
		targetApproval.Internal.GetApproverName(approval.Customer2),
	}

	if approvalType == approvable.APPROVAL_TYPE_PLAUSI {
		users = []string{
			targetApproval.Plausibility.Approver,
		}
	}
	if approvalType == approvable.APPROVAL_TYPE_VEHICLE {
		return
	}
	if approvalType == approvable.APPROVAL_TYPE_EXTERNAL {
		return
	}

	var updated bool
	for _, u := range users {
		if u == "" {
			continue
		}
		missingMember := !(pr.GetMember(u) != nil || u == "")
		if missingMember {
			updated = true
			pr.AddUser(project.ProjectMemberRequestDto{
				TargetUser: u,
				UserType:   project.VIEWER,
			})
			member := pr.GetMember(u)

			userAudit := member.ToAudit(s.RequestSession, s.UserRepo)
			s.AuditLogListRepo.CreateAuditEntryByKey(s.RequestSession, pr.Key, creator, message.ProjectUserCreated, cmp.Diff, userAudit, project.ProjectMemberAudit{})

		}
		if pr.IsGroup {
			for _, child := range pr.Children {
				childProject := s.ProjectRepo.FindByKey(s.RequestSession, child, false)
				if childProject == nil {
					continue
				}
				if childProject.CheckIfUserAlreadyExistsSoft(u) {
					continue
				}
				childProject.AddUser(project.ProjectMemberRequestDto{
					TargetUser: u,
					UserType:   project.VIEWER,
				})
				member := pr.GetMember(u)
				userAudit := member.ToAudit(s.RequestSession, s.UserRepo)
				s.AuditLogListRepo.CreateAuditEntryByKey(s.RequestSession, childProject.Key, creator, message.ProjectUserCreated, cmp.Diff, userAudit, project.ProjectMemberAudit{})
				s.ProjectRepo.Update(s.RequestSession, childProject)
			}
		}
	}
	if updated {
		s.ProjectRepo.Update(s.RequestSession, pr)
	}
}

func (s *ApprovalService) setApprovalSpdxStatus(targetApproval *approval.Approval, state approval.StateInfo) {
	for _, p := range targetApproval.Info.Projects {
		if p.ApprovableSPDX.SpdxKey == "" || p.ApprovableSPDX.VersionKey == "" {
			continue
		}
		sbomList, spdx := s.SpdxRetriever.RetrieveSbomListAndFile(s.RequestSession, p.ApprovableSPDX.VersionKey, p.ApprovableSPDX.SpdxKey)
		if spdx == nil {
			exception.ThrowExceptionServerMessage(message.GetI18N(message.FindingSbomKey), "SPDX not found in history: "+p.ApprovableSPDX.SpdxKey)
		}
		spdx.ApprovalInfo = project.ApprovalInfo{
			IsInApproval: true,
			Comment:      targetApproval.Comment,
			ApprovalGuid: targetApproval.Key,
			Status:       string(state),
		}
		s.SBOMListRepo.Update(s.RequestSession, sbomList)
	}
}

func (s *ApprovalService) adminAbortInternal(pr *project.Project, targetApproval *approval.Approval) {
	if !targetApproval.Internal.IsActive() {
		exception.ThrowExceptionBadRequestResponse()
	}
	before := targetApproval.ToAudit()
	s.deletePending(targetApproval)
	s.setTaskDone(targetApproval.Creator, targetApproval, user.ApprovalInfo, user.TaskPending)
	targetApproval.Internal.Aborted = true
	s.setApprovalSpdxStatus(targetApproval, approval.Aborted)
	s.AuditLogListRepo.CreateAuditEntryByKey(s.RequestSession, pr.Key, "SYSTEM", message.InternalApprovalAborted, audit.DiffWithReporter, targetApproval.ToAudit(), before)
	observermngmt.FireEvent(observermngmt.ApprovalFinalized, observermngmt.ApprovalData{
		RequestSession: s.RequestSession,
		Approval:       targetApproval,
	})
}

func (s *ApprovalService) FillRemainingCustomer(pr *project.Project, appId, username string, req *approval.FillCustomerDto) {
	approvalList := s.ApprovalListRepo.FindByKey(s.RequestSession, pr.Key, false)
	if approvalList == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorDbNotFound), "")
	}
	targetApproval := approvalList.GetApproval(appId)
	if targetApproval == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorDbNotFound), "")
	}

	if targetApproval.Type != approval.TypeInternal || !targetApproval.Internal.IsActive() {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorDbNotFound), "")
	}
	if targetApproval.Internal.Generating {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorDbNotFound), "")
	}
	if targetApproval.Creator != username {
		exception.ThrowExceptionSendDeniedResponse()
	}
	cus := targetApproval.Internal.GetApproverName(approval.Customer1)
	cusChanged := map[approval.Approver]string{}

	if cus != "" {
		if cus != req.CustomerApprover1 {
			exception.ThrowExceptionBadRequestResponse()
		}
	} else {
		cusChanged[approval.Customer1] = req.CustomerApprover1
	}
	cus = targetApproval.Internal.GetApproverName(approval.Customer2)
	if cus != "" {
		if cus != req.CustomerApprover2 {
			exception.ThrowExceptionBadRequestResponse()
		}
	} else {
		cusChanged[approval.Customer2] = req.CustomerApprover2
	}
	before := targetApproval.ToAudit()
	for approver, name := range cusChanged {
		targetApproval.Internal.Approver[approver] = name
		if targetApproval.Internal.SupplierDone() {
			s.setApproverRole(targetApproval, approver)
			s.CreateInternalApprovalTask(targetApproval, approver)
		}
	}
	if len(cusChanged) > 0 {
		s.setProjectApprovers(targetApproval, pr)
	}
	s.addTargetUsers(targetApproval, approvable.APPROVAL_TYPE_INTERNAL, pr, username)
	s.ApprovalListRepo.Update(s.RequestSession, approvalList)
	s.AuditLogListRepo.CreateAuditEntryByKey(s.RequestSession, pr.Key, username, message.InternalApprovalUpdated, audit.DiffWithReporter, targetApproval.ToAudit(), before)
}

func (s *ApprovalService) GetApproverUser(pr *project.Project, appId, username string, approverRole string) *user.UserDto {
	approvalList := s.ApprovalListRepo.FindByKey(s.RequestSession, pr.Key, false)
	if approvalList == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorDbNotFound), "")
	}
	targetApproval := approvalList.GetApproval(appId)
	if targetApproval == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorDbNotFound), "")
	}

	if targetApproval.Type != approval.TypeInternal {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorDbNotFound), "")
	}
	if targetApproval.Creator != username {
		exception.ThrowExceptionSendDeniedResponse()
	}

	role, err := strconv.Atoi(approverRole)
	if err != nil {
		exception.ThrowExceptionBadRequestResponse()
	}
	if role < int(approval.Supplier1) || role > int(approval.Customer2) {
		exception.ThrowExceptionBadRequestResponse()
	}
	cus := targetApproval.Internal.GetApproverName(approval.Approver(role))
	if cus == "" {
		return nil
	}
	user := s.UserRepo.FindByUserId(s.RequestSession, cus)
	return user.ToDto()
}

func (s *ApprovalService) SetGenerated(prKey, key string) {
	approvalList := s.ApprovalListRepo.FindByKey(s.RequestSession, prKey, false)
	if approvalList == nil {
		exception.ThrowExceptionBadRequestResponse()
	}

	appr := approvalList.GetApproval(key)
	if appr == nil {
		exception.ThrowExceptionBadRequestResponse()
	}

	appr.Internal.Generating = false

	s.ApprovalListRepo.Update(s.RequestSession, approvalList)
}

func (s *ApprovalService) SetGenerationFailed(prKey, key string) {
	approvalList := s.ApprovalListRepo.FindByKey(s.RequestSession, prKey, false)
	if approvalList == nil {
		exception.ThrowExceptionBadRequestResponse()
	}

	appr := approvalList.GetApproval(key)
	if appr == nil {
		exception.ThrowExceptionBadRequestResponse()
	}

	appr.Internal.Generating = false
	appr.Internal.GenerationFailed = true

	s.ApprovalListRepo.Update(s.RequestSession, approvalList)
}
