// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package analytics

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/job"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/infra/service/analytics"
	"mercedes-benz.ghe.com/foss/disuko/logy"
	"mercedes-benz.ghe.com/foss/disuko/scheduler"
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
