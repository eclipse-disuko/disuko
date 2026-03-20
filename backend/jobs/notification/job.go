// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package notification

import (
	"github.com/eclipse-disuko/disuko/domain/job"
	"github.com/eclipse-disuko/disuko/helper/middlewareDisco"
	"github.com/eclipse-disuko/disuko/infra/repository/dpconfig"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/eclipse-disuko/disuko/scheduler"
)

type Job struct {
	configRepo *dpconfig.DBConfigRepository
}

func Init(configRepo *dpconfig.DBConfigRepository) *Job {
	return &Job{
		configRepo: configRepo,
	}
}

func (j *Job) Execute(rs *logy.RequestSession, info job.Job) scheduler.ExecutionResult {
	middlewareDisco.CurrentNotification = j.configRepo.Notification.Get(rs)
	return scheduler.ExecutionResult{
		Success: true,
	}
}
