// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package user

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/user"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/infra/repository/approvallist"
	"github.com/eclipse-disuko/disuko/infra/repository/auditloglist"
	"github.com/eclipse-disuko/disuko/infra/repository/database"
	projectRepo "github.com/eclipse-disuko/disuko/infra/repository/project"
	sbomListRepo "github.com/eclipse-disuko/disuko/infra/repository/sbomlist"
	userRepo "github.com/eclipse-disuko/disuko/infra/repository/user"
	approvalService "github.com/eclipse-disuko/disuko/infra/service/approval"
	"github.com/eclipse-disuko/disuko/logy"
)

type DeletionService struct {
	UserRepository         userRepo.IUsersRepository
	ProjectRepository      projectRepo.IProjectRepository
	ApprovalListRepository approvallist.IApprovalListRepository
	SbomListRepository     sbomListRepo.ISbomListRepository
	AuditLogListRepository auditloglist.IAuditLogListRepository
}

type deletion struct {
	rs      *logy.RequestSession
	user    *user.User
	service *DeletionService
}

func (s *DeletionService) DeleteUser(rs *logy.RequestSession, user *user.User) {
	d := deletion{
		rs:      rs,
		user:    user,
		service: s,
	}
	d.rs.Infof("aborting approvals for user %s", user.User)
	d.abortApprovals()
	d.rs.Infof("deleting project memberships for user %s", user.User)
	d.deleteProjectMemberships()
	d.rs.Infof("deleting user %s", user.User)
	d.service.UserRepository.Delete(rs, user.Key)
}

func (d *deletion) deleteProjectMemberships() {
	prs := d.service.ProjectRepository.FindAllForUser(d.rs, d.user.User)
	for _, pr := range prs {
		d.rs.Infof("deleting project membership of %s in %s / %s", d.user.User, pr.Name, pr.Key)
		d.transferOwnership(pr.Key)
	}
}

func (d *deletion) transferOwnership(prKey string) {
	pr := d.service.ProjectRepository.FindByKey(d.rs, prKey, false)
	m := pr.GetMember(d.user.User)
	if m.UserType == project.OWNER {
		if !pr.OtherOwnersExists(m.UserId) {
			exception.ThrowExceptionServerMessage(message.GetI18N(message.TransferOwnershipBlocked), "")
		}
		if m.IsResponsible {
			o := pr.OtherOwner(m.UserId)
			d.rs.Infof("transferring responsibility to %s in %s / %s", o.UserId, pr.Name, pr.Key)
			o.IsResponsible = true
		}
	}
	pr.RemoveMember(d.user.User)
	d.service.ProjectRepository.Update(d.rs, pr)
}

func (d *deletion) abortApprovals() {
	approvals := make(map[string]string)
	for _, t := range d.user.Tasks {
		if t.Status == user.TaskDone {
			continue
		}
		approvals[t.TargetGuid] = t.ProjectGuid
	}
	for appUuid, prKey := range approvals {
		d.rs.Infof("aborting approval %s in project %s", appUuid, prKey)
		d.abortApproval(appUuid, prKey)
	}
}

func (d *deletion) abortApproval(appUuid, prKey string) {
	pr := d.service.ProjectRepository.FindByKey(d.rs, prKey, false)
	if pr == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorDbNotFound), "")
	}
	list := d.service.ApprovalListRepository.FindByKey(d.rs, prKey, false)
	if list == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorDbNotFound), "")
	}
	a := list.GetApproval(appUuid)
	if a == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorDbNotFound), "")
	}
	s := approvalService.ApprovalService{
		RequestSession:   d.rs,
		UserRepo:         d.service.UserRepository,
		SBOMListRepo:     d.service.SbomListRepository,
		AuditLogListRepo: d.service.AuditLogListRepository,
	}
	s.AdminAbortRandomApproval(pr, a)
	d.service.ApprovalListRepository.Update(d.rs, list)
}

func (s *DeletionService) IsDeletable(rs *logy.RequestSession, user *user.User) bool {
	prs := s.ProjectRepository.FindAllForUser(rs, user.User)
	for _, pr := range prs {
		m := pr.GetMember(user.User)
		if m.UserType != project.OWNER {
			continue
		}
		if !pr.OtherOwnersExists(m.UserId) {
			return false
		}
	}
	return true
}

func (s *DeletionService) AffectedUsers(rs *logy.RequestSession) []*user.User {
	cutoff := time.Now().UTC().AddDate(0, -3, 0)
	qc := database.New().SetMatcher(
		database.AndChain(
			database.AttributeMatcher(
				"Deprovisioned",
				database.NE,
				time.Time{},
			),
			database.AttributeMatcher(
				"Deprovisioned",
				database.LT,
				cutoff,
			),
		),
	)
	return s.UserRepository.Query(rs, qc)
}
