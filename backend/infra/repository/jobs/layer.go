// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package jobs

import (
	"github.com/eclipse-disuko/disuko/domain/job"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
)

const JobCollectionName = "jobs"

type IJobsRepository interface {
	base.IBaseRepositoryWithHardDelete[*job.Job]
	FindLatestJob(requestSession *logy.RequestSession, jobType job.JobType) *job.Job
	FindPeriodicJobs(requestSession *logy.RequestSession) []*job.Job
	FindByTypeAndExecution(requestSession *logy.RequestSession, jobType job.JobType, executionType job.ExecutionType) *job.Job
	FindManualJob(requestSession *logy.RequestSession, jobType job.JobType) *job.Job
}
