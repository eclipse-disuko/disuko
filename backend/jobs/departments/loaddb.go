// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package departments

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/job"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/department"
	"mercedes-benz.ghe.com/foss/disuko/logy"
	"mercedes-benz.ghe.com/foss/disuko/scheduler"
)

type LoadDbJob struct {
	repo department.IDepartmentRepository
}

func InitLoadDb(repo department.IDepartmentRepository) *LoadDbJob {
	return &LoadDbJob{
		repo: repo,
	}
}

func (j *LoadDbJob) Execute(rs *logy.RequestSession, info job.Job) scheduler.ExecutionResult {
	var log job.Log
	log.AddEntry(job.Info, "started")
	amount := j.repo.LoadFromDb(rs)
	log.AddEntry(job.Info, "loaded %d departments from db", amount)
	log.AddEntry(job.Info, "finished")
	return scheduler.ExecutionResult{
		Success: true,
		Log:     log,
	}
}
