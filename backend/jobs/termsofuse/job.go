// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package termsofuse

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain/audit"
	"github.com/eclipse-disuko/disuko/domain/job"
	"github.com/eclipse-disuko/disuko/domain/user"
	auditHelper "github.com/eclipse-disuko/disuko/helper/audit"
	"github.com/eclipse-disuko/disuko/helper/message"
	userRepo "github.com/eclipse-disuko/disuko/infra/repository/user"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/eclipse-disuko/disuko/scheduler"
)

type Job struct {
	repo userRepo.IUsersRepository
}

func Init(repo userRepo.IUsersRepository) *Job {
	return &Job{
		repo: repo,
	}
}

func (j *Job) Execute(rs *logy.RequestSession, info job.Job) scheduler.ExecutionResult {
	var log job.Log
	log.AddEntry(job.Info, "started")

	var customRes struct {
		UserUpdated int `json:"usersUpdated"`
	}
	newTermsOfUseVersion := info.Config
	users := j.repo.FindAll(rs, false)
	usersToUpdate := []*user.User{}
	for _, user := range users {
		if user.TermsOfUseVersion == newTermsOfUseVersion {
			continue
		}
		oldAudit := user.ToUserAudit()
		user.TermsOfUseVersion = newTermsOfUseVersion
		user.TermsOfUse = false
		user.Updated = time.Now()
		newAudit := user.ToUserAudit()
		auditHelper.CreateAndAddAuditEntry(&user.Container, "SYSTEM", message.UserUpdated, audit.DiffWithReporter, newAudit, oldAudit)
		usersToUpdate = append(usersToUpdate, user)
	}
	j.repo.UpdateListSilent(rs, usersToUpdate)

	log.AddEntry(job.Info, "succesfully reset acceptance and updated version (to %s) for %d users.", newTermsOfUseVersion, len(usersToUpdate))
	log.AddEntry(job.Info, "finished")
	customRes.UserUpdated = len(usersToUpdate)
	return scheduler.ExecutionResult{
		Success:   true,
		Log:       log,
		CustomRes: customRes,
	}
}
