// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package job

import (
	"time"

	"github.com/adhocore/gronx"
	"github.com/eclipse-disuko/disuko/domain"
	"github.com/gorhill/cronexpr"
)

type JobDto struct {
	domain.BaseDto
	Name                   string        `json:"name"`
	JobType                JobType       `json:"jobType"`
	Execution              ExecutionType `json:"execution"`
	Status                 Status        `json:"status"`
	Config                 string        `json:"config"`
	CustomRes              interface{}   `json:"customRes"`
	Log                    []LogEntryDto `json:"log"`
	NextScheduledExecution string        `json:"nextScheduledExecution,omitempty"`
}

type LogEntryDto struct {
	domain.BaseDto
	Msg          string `json:"msg"`
	Level        Level  `json:"level"`
	InstanceName string `json:"instance"`
}

type SetConfigDto struct {
	Config string `json:"config"`
}

type JobStatusDto struct {
	Status Status `json:"status"`
}

func (entity Log) ToDto() (log []LogEntryDto) {
	for _, entry := range entity {
		log = append(log, LogEntryDto{
			BaseDto: domain.BaseDto{
				Created: entry.Created,
				Updated: entry.Updated,
			},
			Msg:          entry.Msg,
			Level:        entry.Level,
			InstanceName: entry.InstanceName,
		})
	}
	return
}

func ToDto(entity *Job) *JobDto {
	dto := &JobDto{
		Name:      entity.Name,
		JobType:   entity.JobType,
		Execution: entity.Execution,
		Status:    entity.Status,
		Log:       entity.Log.ToDto(),
		CustomRes: entity.CustomRes,
		Config:    entity.Config,
	}
	domain.SetBaseValues(entity, dto)

	if entity.Execution == Periodic && entity.Schedule != "" {
		dto.NextScheduledExecution = calculateNextExecution(entity)
	}

	return dto
}

func calculateNextExecution(entity *Job) string {
	if entity.Schedule == "" {
		return ""
	}
	gron := gronx.New()
	if !gron.IsValid(entity.Schedule) {
		return ""
	}
	now := time.Now().Truncate(time.Minute)
	expr, err := cronexpr.Parse(entity.Schedule)
	if err != nil {
		return ""
	}
	nextTime := expr.Next(now)
	if nextTime.IsZero() {
		return ""
	}
	return nextTime.Format("2006-01-02 15:04:05")
}

func ToDtoList(entities []*Job) []*JobDto {
	return domain.MapTo(entities, ToDto)
}
