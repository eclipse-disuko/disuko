// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package jobs

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/job"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

const JobCollectionName = "jobs"

type IJobsRepository interface {
	base.IBaseRepositoryWithHardDelete[*job.Job]
	FindLatestJob(requestSession *logy.RequestSession, jobType job.JobType) *job.Job
	FindPeriodicJobs(requestSession *logy.RequestSession) []*job.Job
	FindByTypeAndExecution(requestSession *logy.RequestSession, jobType job.JobType, executionType job.ExecutionType) *job.Job
	FindManualJob(requestSession *logy.RequestSession, jobType job.JobType) *job.Job
}
