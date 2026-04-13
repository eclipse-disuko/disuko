// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package fossdd

import (
	"encoding/json"

	"github.com/eclipse-disuko/disuko/domain/approval"
	"github.com/eclipse-disuko/disuko/domain/job"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/infra/repository/approvallist"
	"github.com/eclipse-disuko/disuko/infra/repository/policyrules"
	projectRepo "github.com/eclipse-disuko/disuko/infra/repository/project"
	"github.com/eclipse-disuko/disuko/infra/repository/user"
	approvalService "github.com/eclipse-disuko/disuko/infra/service/approval"
	"github.com/eclipse-disuko/disuko/infra/service/fossdd"
	projectLabelService "github.com/eclipse-disuko/disuko/infra/service/project-label"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/eclipse-disuko/disuko/scheduler"
)

type Config struct {
	ApprovalID string
	ProjectID  string
	WithZIP    bool
	Template   string
}

type Job struct {
	repo                approvallist.IApprovalListRepository
	projectRepo         projectRepo.IProjectRepository
	userRepo            user.IUsersRepository
	projectLabelService *projectLabelService.ProjectLabelService
	externCheckCreator  fossdd.ExternCheckCreator
	fossddService       *fossdd.Service
}

func Init(repo approvallist.IApprovalListRepository, projectRepo projectRepo.IProjectRepository, policyrulesRepo policyrules.IPolicyRulesRepository, userRepo user.IUsersRepository, externCheckCreator fossdd.ExternCheckCreator, fossddService *fossdd.Service, projectLabelService *projectLabelService.ProjectLabelService) *Job {
	return &Job{
		repo:                repo,
		projectRepo:         projectRepo,
		userRepo:            userRepo,
		projectLabelService: projectLabelService,
		externCheckCreator:  externCheckCreator,
		fossddService:       fossddService,
	}
}

func (j *Job) Execute(rs *logy.RequestSession, info job.Job) scheduler.ExecutionResult {
	var log job.Log
	log.AddEntry(job.Info, "started")

	var config Config

	if err := json.Unmarshal([]byte(info.Config), &config); err != nil {
		log.AddEntry(job.Error, "decoding config failed: %s", err)
		return scheduler.ExecutionResult{
			Success: false,
			Log:     log,
		}
	}

	pr := j.projectRepo.FindByKey(rs, config.ProjectID, false)
	if pr == nil {
		log.AddEntry(job.Info, "project not found")
		return scheduler.ExecutionResult{
			Success: false,
			Log:     log,
		}
	}
	apprs := j.repo.FindByKey(rs, config.ProjectID, false)
	if apprs == nil {
		log.AddEntry(job.Info, "approvallist not found")
		return scheduler.ExecutionResult{
			Success: false,
			Log:     log,
		}
	}
	appr := apprs.GetApproval(config.ApprovalID)
	if appr == nil {
		log.AddEntry(job.Info, "approval not found")
		return scheduler.ExecutionResult{
			Success: false,
			Log:     log,
		}

	}

	exc := false
	exception.TryCatch(func() {
		var subs []fossdd.SubProjectRefs
		for _, pr := range appr.Info.Projects {
			sub := fossdd.SubProjectRefs{
				ProjectID: pr.ProjectKey,
			}
			if pr.ApprovableSPDX.SpdxKey != "" && pr.ApprovableSPDX.VersionKey != "" {
				sub.SpdxID = &pr.ApprovableSPDX.SpdxKey
				sub.VersionID = &pr.ApprovableSPDX.VersionKey
			}
			subs = append(subs, sub)
		}
		j.fossddService.Generate(rs, fossdd.GenOpts{
			Approval:        appr,
			MainProjectID:   pr.Key,
			SubProjectsRefs: subs,
			Flags: fossdd.Flags{
				C1: appr.DocumentFlags.C1,
				C2: appr.DocumentFlags.C2,
				C3: appr.DocumentFlags.C3,
				C4: appr.DocumentFlags.C4,
				C5: appr.DocumentFlags.C5,
				C6: appr.DocumentFlags.C6,
			},
			WithZIP:  config.WithZIP,
			Template: config.Template,
		}, &log, j.externCheckCreator)
	}, func(e exception.Exception) {
		exception.LogException(rs, e)
		log.AddEntry(job.Error, "failed to create documents. exception message: %s", e.ErrorMessage)
		exc = true
	})

	as := approvalService.ApprovalService{
		RequestSession:      rs,
		ApprovalListRepo:    j.repo,
		UserRepo:            j.userRepo,
		ProjectLabelService: j.projectLabelService,
	}
	if exc {
		if appr.Type == approval.TypeInternal {
			as.SetGenerationFailed(pr.Key, appr.Key)
		} else {
			as.SetExternalState(pr.Key, appr.Key, approval.GenerationFailed)
		}
		return scheduler.ExecutionResult{
			Success: false,
			Log:     log,
		}
	}

	if appr.Type == approval.TypeInternal {
		log.AddEntry(job.Info, "creating supplier tasks")
		as.CreateInternalApprovalTask(appr, approval.Supplier1)
		as.CreateInternalApprovalTask(appr, approval.Supplier2)
		log.AddEntry(job.Info, "setting generated state")
		as.SetGenerated(pr.Key, appr.Key)
	} else {

		log.AddEntry(job.Info, "setting generated state")
		as.SetExternalState(pr.Key, appr.Key, approval.Pending)
	}
	log.AddEntry(job.Info, "finished with project %s", config.ProjectID)
	return scheduler.ExecutionResult{
		Success: true,
		Log:     log,
	}
}
