// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package approvals

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain/job"
	"github.com/eclipse-disuko/disuko/infra/repository/approvallist"
	"github.com/eclipse-disuko/disuko/infra/repository/auditloglist"
	projectRepo "github.com/eclipse-disuko/disuko/infra/repository/project"
	sbomlistRepo "github.com/eclipse-disuko/disuko/infra/repository/sbomlist"
	userRepo "github.com/eclipse-disuko/disuko/infra/repository/user"
	approvalService "github.com/eclipse-disuko/disuko/infra/service/approval"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/eclipse-disuko/disuko/scheduler"
)

const abortOnDay = 21

type AbortInactive struct {
	approvalListRepo approvallist.IApprovalListRepository
	projectRepo      projectRepo.IProjectRepository
	userRepo         userRepo.IUsersRepository
	auditLogListRepo auditloglist.IAuditLogListRepository
	sbomListRepo     sbomlistRepo.ISbomListRepository
	spdxRetriever    approvalService.SpdxRetriever
}

func InitAbortInactive(
	approvalListRepo approvallist.IApprovalListRepository,
	projectRepo projectRepo.IProjectRepository,
	userRepo userRepo.IUsersRepository,
	auditLogListRepo auditloglist.IAuditLogListRepository,
	sbomListRepo sbomlistRepo.ISbomListRepository,
	spdxRetriever approvalService.SpdxRetriever,
) *AbortInactive {
	return &AbortInactive{
		approvalListRepo: approvalListRepo,
		projectRepo:      projectRepo,
		userRepo:         userRepo,
		auditLogListRepo: auditLogListRepo,
		sbomListRepo:     sbomListRepo,
		spdxRetriever:    spdxRetriever,
	}
}

func (j *AbortInactive) Execute(rs *logy.RequestSession, info job.Job) scheduler.ExecutionResult {
	var log job.Log
	log.AddEntry(job.Info, "started")

	approvalLists := j.approvalListRepo.FindAll(rs, false)
	for _, list := range approvalLists {
		changed := false
		for i := range list.Approvals {
			appr := &list.Approvals[i]
			if !isOngoing(appr) {
				continue
			}
			if time.Since(appr.Updated) < abortOnDay*24*time.Hour || time.Since(appr.Updated) >= (abortOnDay+1)*24*time.Hour {
				continue
			}
			pr := j.projectRepo.FindByKey(rs, list.Key, false)
			if pr == nil {
				log.AddEntry(job.Error, "project %s not found for approval %s", list.Key, appr.Key)
				continue
			}
			as := approvalService.ApprovalService{
				RequestSession:   rs,
				ApprovalListRepo: j.approvalListRepo,
				UserRepo:         j.userRepo,
				AuditLogListRepo: j.auditLogListRepo,
				SBOMListRepo:     j.sbomListRepo,
				SpdxRetriever:    j.spdxRetriever,
			}
			as.AdminAbortRandomApproval(pr, appr)
			changed = true
			log.AddEntry(job.Info, "aborted approval %s (project %s, type %s)", appr.Key, list.Key, appr.Type)
		}
		if changed {
			j.approvalListRepo.Update(rs, list)
		}
	}

	log.AddEntry(job.Info, "finished")
	return scheduler.ExecutionResult{
		Success: true,
		Log:     log,
	}
}
