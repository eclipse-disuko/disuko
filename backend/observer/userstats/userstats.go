// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package userstats

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/job"
	"mercedes-benz.ghe.com/foss/disuko/domain/userstats"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	job_userstats "mercedes-benz.ghe.com/foss/disuko/jobs/userstats"
	"mercedes-benz.ghe.com/foss/disuko/observermngmt"
	"mercedes-benz.ghe.com/foss/disuko/scheduler"
)

type UserStats struct {
	Scheduler *scheduler.Scheduler
}

func Init(scheduler *scheduler.Scheduler) *UserStats {
	return &UserStats{
		Scheduler: scheduler,
	}
}

func (o *UserStats) RegisterHandlers() {
	observermngmt.RegisterHandler(observermngmt.DatabaseEntryAddedOrDeleted, o.OnDatabaseEntryAddedOrDeleted)
}

func (o *UserStats) OnDatabaseEntryAddedOrDeleted(eventId observermngmt.EventId, arg interface{}) {
	data, ok := arg.(observermngmt.DatabaseSizeChange)
	if !ok {
		return
	}
	go exception.TryCatchAndLog(data.RequestSession, func() {
		_, err := o.Scheduler.ExecuteOneTimeJob(data.RequestSession, "calculate user or system stats", job.CalculateUserStats, job_userstats.Config{Username: data.Username, Rights: data.Rights, CollectionName: data.CollectionName, AdminRequest: data.Username == userstats.SystemStatsUser})
		exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorStartingJob))
	})
}
