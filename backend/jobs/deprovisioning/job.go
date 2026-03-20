// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package deprovisioning

import (
	"encoding/json"
	"time"

	"mercedes-benz.ghe.com/foss/disuko/connector/userrole"
	"mercedes-benz.ghe.com/foss/disuko/domain/audit"
	"mercedes-benz.ghe.com/foss/disuko/domain/job"
	"mercedes-benz.ghe.com/foss/disuko/domain/user"
	auditHelper "mercedes-benz.ghe.com/foss/disuko/helper/audit"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/database"
	userRepo "mercedes-benz.ghe.com/foss/disuko/infra/repository/user"
	"mercedes-benz.ghe.com/foss/disuko/logy"
	"mercedes-benz.ghe.com/foss/disuko/scheduler"
)

var (
	desiredChunkSize      = 500
	inactiveDaysThreshold = 60
)

type Config struct {
	InProgress bool `json:"inProgress"`
	CurrChunk  int  `json:"currChunk"`
	ChunkNum   int  `json:"chunkNum"`
	Remaining  int  `json:"remaining"`
}

type Job struct {
	repo      userRepo.IUsersRepository
	connector *userrole.Connector
}

func Init(repo userRepo.IUsersRepository, connector *userrole.Connector) *Job {
	return &Job{
		repo:      repo,
		connector: connector,
	}
}

func (j *Job) Execute(rs *logy.RequestSession, info job.Job) scheduler.ExecutionResult {
	if j.connector == nil {
		return scheduler.ExecutionResult{
			Success: true,
		}
	}
	var log job.Log
	log.AddEntry(job.Info, "started")

	confStr := info.Config
	var config Config

	if info.Config == "" {
		config = Config{}
	} else if err := json.Unmarshal([]byte(confStr), &config); err != nil {
		log.AddEntry(job.Error, "decoding config failed: %s", err)
		return scheduler.ExecutionResult{
			Success: false,
			Log:     log,
		}
	}

	if !config.InProgress {
		config = j.newConfig(rs)
		if !config.InProgress {
			log.AddEntry(job.Info, "nothing to do")
			return scheduler.ExecutionResult{
				Success: true,
				Log:     log,
			}
		}
	}

	j.processNextChunk(rs, &log, &config)

	jsonB, err := json.Marshal(config)
	if err != nil {
		log.AddEntry(job.Error, "encoding config failed: %s", err)
		return scheduler.ExecutionResult{
			Success: false,
			Log:     log,
		}
	}

	jsonS := string(jsonB)

	return scheduler.ExecutionResult{
		Success: true,
		Log:     log,
		NewConf: &jsonS,
	}
}

func (j *Job) newConfig(rs *logy.RequestSession) Config {
	var res Config

	u := j.affectedUsers(rs)
	if len(u) == 0 {
		return res
	}
	res.ChunkNum = int(len(u) / desiredChunkSize)
	res.CurrChunk = 0
	res.Remaining = len(u) % desiredChunkSize
	if res.Remaining > 0 {
		res.ChunkNum++
	}
	res.InProgress = true

	return res
}

func (j *Job) processNextChunk(rs *logy.RequestSession, log *job.Log, conf *Config) {
	lower := conf.CurrChunk * desiredChunkSize
	upper := lower
	if conf.CurrChunk == conf.ChunkNum-1 {
		conf.InProgress = false
		if conf.Remaining > 0 {
			upper += conf.Remaining
		} else {
			upper += desiredChunkSize
		}
	} else {
		upper += desiredChunkSize
	}

	affected := j.affectedUsers(rs)
	if lower >= len(affected) {
		log.AddEntry(job.Info, "lower limit exceeds affected users slice")
		conf.InProgress = false
		return
	}
	if upper >= len(affected) {
		upper = len(affected)
		conf.InProgress = false
	}

	log.AddEntry(job.Info, "processing chunk %d - %d", lower, upper)
	chunk := affected[lower:upper]
	found := 0
	for _, u := range chunk {
		active := j.connector.UserActive(rs, u.User)
		if active {
			continue
		}
		log.AddEntry(job.Info, "setting deprovisioning date and inactive for %s", u.User)
		beforeAudit := u.ToUserAudit()
		u.Deprovisioned = time.Now()
		u.Active = false
		afterAudit := u.ToUserAudit()
		auditHelper.CreateAndAddAuditEntry(&u.Container, "SYSTEM", message.UserUpdated, audit.DiffWithReporter, afterAudit, beforeAudit)
		j.repo.Update(rs, u)
		found++
	}
	conf.CurrChunk++
	log.AddEntry(job.Info, "processed chunk changed %d / %d", found, len(chunk))
}

func (j *Job) affectedUsers(rs *logy.RequestSession) []*user.User {
	all := j.repo.Query(rs, database.New().SetSort(database.SortConfig{
		database.SortAttribute{
			Name:  "Created",
			Order: database.ASC,
		},
	}))

	var filtered []*user.User
	for _, u := range all {
		if !u.Deprovisioned.IsZero() {
			continue
		}
		if time.Since(u.Created) < time.Hour*24*time.Duration(inactiveDaysThreshold) {
			continue
		}
		filtered = append(filtered, u)
	}

	return filtered
}
