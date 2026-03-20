// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package job

import (
	"fmt"
	"time"

	"mercedes-benz.ghe.com/foss/disuko/conf"
	"mercedes-benz.ghe.com/foss/disuko/domain"
)

type Job struct {
	domain.RootEntity `bson:"inline"`
	Name              string
	JobType           JobType
	Status            Status
	Execution         ExecutionType
	MultiExec         bool
	Log               Log
	Schedule          string
	CustomRes         interface{}
	Config            string
}

type Log []LogEntry

type LogEntry struct {
	domain.ChildEntity `bson:"inline"`
	Msg                string
	Level              Level
	InstanceName       string
}
type Level string

const (
	Error Level = "error"
	Warn  Level = "warn"
	Info  Level = "info"
)

type JobType int

const (
	LicenseRefresh JobType = iota
	TermOfUseUpdate
	DepartmentRefresh
	DepartmentLoadDb
	LicenseAnnouncements
	AnalyticsRebuild
	Notification
	Deprovisioning
	PolicyRuleChangeLogs
	Report
	FOSSDDGen
	DummyProjectDeletion
	DummyMail
	CalculateUserStats
	LabelLoadDb
	TypeLimit
)

type Status int

const (
	Idle Status = iota
	InProgress
	Success
	Failure
)

type ExecutionType int

const (
	Manual ExecutionType = iota
	Periodic
	OneTime
)

const JobLogTTL = time.Hour * 24 * 7

func (l *Log) AddEntry(level Level, format string, args ...interface{}) {
	*l = append(*l, LogEntry{
		ChildEntity:  domain.NewChildEntity(),
		InstanceName: conf.Config.Server.InstanceName,
		Level:        level,
		Msg:          fmt.Sprintf(format, args...),
	})
}

func (j *Job) AddLog(l Log) {
	if l == nil {
		return
	}
	cutOff := 0
	for i, e := range j.Log {
		if time.Since(e.Created) < JobLogTTL {
			cutOff = i
			break
		}
	}
	j.Log = append(j.Log[cutOff:], l...)
}
