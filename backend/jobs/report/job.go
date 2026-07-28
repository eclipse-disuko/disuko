// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package report

import (
	"github.com/eclipse-disuko/disuko/domain/job"
	"github.com/eclipse-disuko/disuko/infra/service/report"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/eclipse-disuko/disuko/scheduler"
)

type Job struct {
	service *report.Service
}

func Init(reportService *report.Service) *Job {
	return &Job{
		service: reportService,
	}
}

func (j *Job) Execute(rs *logy.RequestSession, info job.Job) scheduler.ExecutionResult {
	var log job.Log
	log.AddEntry(job.Info, "started")

	res := j.service.Generate(rs, &log)

	log.AddEntry(job.Info, "finished")

	return scheduler.ExecutionResult{
		Success:   true,
		Log:       log,
		CustomRes: res,
	}
}
