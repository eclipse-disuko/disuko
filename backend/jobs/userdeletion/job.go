// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package userdeletion

import (
	"github.com/eclipse-disuko/disuko/domain/job"
	userService "github.com/eclipse-disuko/disuko/infra/service/user"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/eclipse-disuko/disuko/scheduler"
)

type Job struct {
	service *userService.DeletionService
}

func Init(service *userService.DeletionService) *Job {
	return &Job{
		service: service,
	}
}

func (j *Job) Execute(rs *logy.RequestSession, info job.Job) scheduler.ExecutionResult {
	var log job.Log
	log.AddEntry(job.Info, "started")

	affected := j.service.AffectedUsers(rs)
	log.AddEntry(job.Info, "found %d users eligible for deletion", len(affected))

	for _, u := range affected {
		if !j.service.IsDeletable(rs, u) {
			log.AddEntry(job.Info, "skipping user %s, not deletable", u.User)
			continue
		}
		j.service.DeleteUser(rs, u)
		log.AddEntry(job.Info, "deleted user %s", u.User)
	}

	log.AddEntry(job.Info, "finished")
	return scheduler.ExecutionResult{
		Success: true,
		Log:     log,
	}
}
