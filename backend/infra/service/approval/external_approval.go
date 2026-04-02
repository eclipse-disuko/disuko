// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package approval

import (
	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/approval"
	"github.com/eclipse-disuko/disuko/domain/audit"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/project/approvable"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/google/uuid"
)

func (s *ApprovalService) CreateExternalApproval(pr *project.Project, req approval.RequestExternalApprovalDto, creator string, vehicle bool) string {
	if pr.HasParent() {
		exception.ThrowExceptionSendDeniedResponse()
	}
	if s.ProjectLabelService.HasVehiclePlatformLabel(s.RequestSession, pr) != vehicle {
		exception.ThrowExceptionBadRequestResponse()
	}

	info := s.getApprovalInfo(pr, &req.SelectedProjects, false)
	if len(info.Projects) == 0 {
		exception.ThrowExceptionBadRequestResponse()
	}

	if vehicle && (info.CompStats.Denied > 0 || info.CompStats.NoAssertion > 0) {
		exception.ThrowExceptionBadRequestResponse()
	}
	appr := approval.Approval{
		ChildEntity: domain.SetChildEntity(uuid.New().String()),
		Type:        approval.TypeExternal,
		ProjectGuid: pr.Key,
		Creator:     creator,
		Comment:     req.Comment,
		Info:        info,
		External:    approval.ExternalApproval{State: approval.GeneratingDocs, Vehicle: vehicle},

		DocumentFlags: approval.TaskMetaDocument{
			C1: req.MetaDoc.C1,
			C2: req.MetaDoc.C2,
			C3: req.MetaDoc.C3,
			C4: req.MetaDoc.C4,
			C5: req.MetaDoc.C5,
			C6: req.MetaDoc.C6,
		},
	}

	s.addTargetUsers(&appr, approvable.APPROVAL_TYPE_EXTERNAL, pr, creator)
	s.AuditLogListRepo.AddStaticAuditEntryByKey(s.RequestSession, pr.Key, creator, message.ExternalApprovalCreated, appr.ToAudit())
	// auditHelper.CreateAndAddAuditEntry(&res.Container, creator, message.ExternalApprovalCreated, audit.DiffWithReporter, res.ToAudit(), approval.ApprovalAudit{})

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

func (s *ApprovalService) processExternalApprovalUpdate(pr *project.Project, targetApproval *approval.Approval, username string, req approval.UpdateApprovalDto) {
	if m := pr.GetMember(username); m == nil || m.UserType != project.OWNER {
		exception.ThrowExceptionSendDeniedResponse()
	}
	before := targetApproval.ToAudit()

	targetApproval.External.Comment = req.Comment

	var valid bool
	valid, targetApproval.External.State = approval.ParseStateInfo(req.State)
	if !valid {
		exception.ThrowExceptionBadRequestResponse()
	}

	s.AuditLogListRepo.CreateAuditEntryByKey(s.RequestSession, pr.Key, username, message.ExternalApprovalUpdated, audit.DiffWithReporter, targetApproval.ToAudit(), before)
	// auditHelper.CreateAndAddAuditEntry(&targetApproval.Container, username, message.ExternalApprovalUpdated, audit.DiffWithReporter, targetApproval.ToAudit(), before)
}

func (s *ApprovalService) SetExternalGenerated(prKey, key string) {
	approvalList := s.ApprovalListRepo.FindByKey(s.RequestSession, prKey, false)
	if approvalList == nil {
		exception.ThrowExceptionBadRequestResponse()
	}

	appr := approvalList.GetApproval(key)
	if appr == nil {
		exception.ThrowExceptionBadRequestResponse()
	}

	appr.External.State = approval.Pending

	s.ApprovalListRepo.Update(s.RequestSession, approvalList)
}
