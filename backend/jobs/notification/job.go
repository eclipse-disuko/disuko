// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package notification

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/job"
	"mercedes-benz.ghe.com/foss/disuko/helper/middlewareDisco"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/dpconfig"
	"mercedes-benz.ghe.com/foss/disuko/logy"
	"mercedes-benz.ghe.com/foss/disuko/scheduler"
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
