// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package jobs

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/job"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/database"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type jobssRepositoryStruct struct {
	base.BaseRepositoryWithHardDelete[*job.Job]
}

func NewJobsRepository(requestSession *logy.RequestSession) IJobsRepository {
	return &jobssRepositoryStruct{
		BaseRepositoryWithHardDelete: base.CreateRepositoryWithHardDelete[*job.Job](
			requestSession,
			JobCollectionName,
			func() *job.Job {
				return &job.Job{}
			},
			nil,
			"",
			nil,
			nil),
	}
}

func (jr *jobssRepositoryStruct) FindLatestJob(requestSession *logy.RequestSession, jobType job.JobType) *job.Job {
	qc := database.New().SetMatcher(
		database.AttributeMatcher(
			"JobType",
			database.EQ,
			jobType,
		),
	).SetSort(database.SortConfig{
		database.SortAttribute{
			Name:  "Created",
			Order: database.DESC,
		},
	}).SetLimit(0, 1)

	jobs := jr.Query(requestSession, qc)
	if len(jobs) > 0 {
		return jobs[0]
	}
	return nil
}

func (jr *jobssRepositoryStruct) FindByTypeAndExecution(requestSession *logy.RequestSession, jobType job.JobType, executionType job.ExecutionType) *job.Job {
	qc := database.New().SetMatcher(
		database.AndChain(
			database.AttributeMatcher(
				"Execution",
				database.EQ,
				executionType,
			),
			database.AttributeMatcher(
				"JobType",
				database.EQ,
				jobType,
			),
		)).SetLimit(0, 1)

	jobs := jr.Query(requestSession, qc)
	if len(jobs) > 0 {
		return jobs[0]
	}
	return nil
}

func (jr *jobssRepositoryStruct) FindPeriodicJobs(requestSession *logy.RequestSession) []*job.Job {
	qc := database.New().SetMatcher(
		database.AttributeMatcher(
			"Execution",
			database.EQ,
			job.Periodic,
		),
	).SetSort(database.SortConfig{
		database.SortAttribute{
			Name:  "Created",
			Order: database.DESC,
		},
	})

	qbL := jr.Query(requestSession, qc)
	return qbL
}

func (jr *jobssRepositoryStruct) FindManualJob(requestSession *logy.RequestSession, jobType job.JobType) *job.Job {
	qc := database.New().SetMatcher(
		database.AndChain(
			database.AttributeMatcher(
				"JobType",
				database.EQ,
				jobType,
			),
			database.AttributeMatcher(
				"Execution",
				database.EQ,
				job.Manual,
			),
		),
	).SetLimit(0, 1)

	jobs := jr.Query(requestSession, qc)
	if len(jobs) > 0 {
		return jobs[0]
	}
	return nil

}
