// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package analytics

import (
	"github.com/eclipse-disuko/disuko/domain/job"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/infra/service/analytics"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/eclipse-disuko/disuko/scheduler"
)

type Job struct {
	service analytics.Analytics
}

func Init(service analytics.Analytics) *Job {
	return &Job{
		service: service,
	}
}

func (j *Job) Execute(rs *logy.RequestSession, info job.Job) scheduler.ExecutionResult {
	var (
		log     job.Log
		success = true
	)

	log.AddEntry(job.Info, "started")
	exception.TryCatch(func() {
		j.service.Reinitialise(rs)
	}, func(ex exception.Exception) {
		log.AddEntry(job.Error, "failed with exception %s", ex.ErrorMessage)
		exception.LogException(rs, ex)
		success = false

	})
	if success {
		log.AddEntry(job.Info, "finished")
	}
	return scheduler.ExecutionResult{
		Success: success,
		Log:     log,
	}
}
