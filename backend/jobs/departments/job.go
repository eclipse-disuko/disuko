// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package departments

import (
	"time"

	departmentConn "mercedes-benz.ghe.com/foss/disuko/connector/department"
	"mercedes-benz.ghe.com/foss/disuko/domain"
	"mercedes-benz.ghe.com/foss/disuko/domain/department"
	"mercedes-benz.ghe.com/foss/disuko/domain/job"
	depRepo "mercedes-benz.ghe.com/foss/disuko/infra/repository/department"
	"mercedes-benz.ghe.com/foss/disuko/logy"
	"mercedes-benz.ghe.com/foss/disuko/scheduler"
)

var waitDBWrite = time.Minute * 5

type RefreshJob struct {
	connector *departmentConn.Connector
	repo      depRepo.IDepartmentRepository
	scheduler *scheduler.Scheduler
}

func Init(repo depRepo.IDepartmentRepository, connector *departmentConn.Connector, scheduler *scheduler.Scheduler) *RefreshJob {
	return &RefreshJob{
		connector: connector,
		repo:      repo,
		scheduler: scheduler,
	}
}

func (j *RefreshJob) Execute(rs *logy.RequestSession, info job.Job) scheduler.ExecutionResult {
	var log job.Log
	log.AddEntry(job.Info, "started")
	if j.connector == nil {
		log.AddEntry(job.Info, "no connector set")
		return scheduler.ExecutionResult{
			Success: true,
			Log:     log,
		}
	}
	log.AddEntry(job.Info, "started")

	connDeps := j.connector.GetDepartments(rs)

	var deps []*department.Department
	for _, d := range connDeps {
		deps = append(deps, &department.Department{
			RootEntity:         domain.NewRootEntityWithKey(d.Id),
			ParentDeptId:       d.ParentId,
			DescriptionEnglish: d.Description,
			OrgAbbreviation:    d.Abbrevation,
			Skz:                d.SecondaryId,
			CompanyCode:        d.CompanyCode,
			CompanyName:        d.CompanyName,
			Level:              d.Level,
		})
	}

	j.repo.SaveDepartments(rs, deps)
	log.AddEntry(job.Info, "finished imported %d departments", len(connDeps))

	time.Sleep(waitDBWrite)
	err := j.scheduler.ExecuteJobManual(rs, job.DepartmentLoadDb)
	if err != nil {
		log.AddEntry(job.Error, "could not execute followup job: %s", err.Error())
	}

	return scheduler.ExecutionResult{
		Success: err == nil,
		Log:     log,
	}
}
