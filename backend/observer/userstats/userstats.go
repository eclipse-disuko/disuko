// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package userstats

import (
	"github.com/eclipse-disuko/disuko/domain/job"
	"github.com/eclipse-disuko/disuko/domain/userstats"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	job_userstats "github.com/eclipse-disuko/disuko/jobs/userstats"
	"github.com/eclipse-disuko/disuko/observermngmt"
	"github.com/eclipse-disuko/disuko/scheduler"
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
