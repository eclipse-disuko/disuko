// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package scheduler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/valkey-io/valkey-go"

	"github.com/adhocore/gronx"
	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/job"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/infra/repository/jobs"
	"github.com/eclipse-disuko/disuko/infra/service/locks"
	"github.com/eclipse-disuko/disuko/logy"
)

type Scheduler struct {
	rs *logy.RequestSession
	// TODO: fix context hierarchy and just use them
	stopCh       chan struct{}
	repo         jobs.IJobsRepository
	valkeyClient valkey.Client

	wg          sync.WaitGroup
	lockService *locks.Service
	jobCbTable  map[job.JobType]JobExecutor
	gron        *gronx.Gronx
}

const (
	dbLookupInterval = time.Second * 25
	jobLockTTL       = time.Hour
	oneTimeTimeout   = time.Hour * 2
	pubTimeout       = time.Second * 20
)

var (
	errAlreadyRunning = errors.New("job is already being executed")
	ErrNotTimedout    = errors.New("job not timed out yet")
)

type ExecutionResult struct {
	Success   bool
	Log       job.Log
	CustomRes any
	NewConf   *string
}

type JobExecutor interface {
	Execute(*logy.RequestSession, job.Job) ExecutionResult
}

func Init(requestSession *logy.RequestSession, repo jobs.IJobsRepository, lockService *locks.Service) *Scheduler {
	res := Scheduler{
		repo:        repo,
		rs:          requestSession,
		jobCbTable:  make(map[job.JobType]JobExecutor),
		lockService: lockService,
		gron:        gronx.New(),
		stopCh:      make(chan struct{}),
	}

	clientOption := valkey.ClientOption{
		InitAddress: []string{conf.Config.Cache.Host + ":" + strconv.Itoa(conf.Config.Cache.Port)},
		Password:    conf.Config.Cache.Password,
		SelectDB:    0,
	}
	valkeyClient, err := valkey.NewClient(clientOption)
	if err != nil {
		logy.Fatalf(requestSession, "valkey client creation failed: %s", err)
	}
	res.valkeyClient = valkeyClient
	return &res
}

func (s *Scheduler) Wait() {
	s.wg.Wait()
}

func (s *Scheduler) Stop() {
	close(s.stopCh)
}

func (s *Scheduler) AddJobCb(t job.JobType, executor JobExecutor) {
	s.jobCbTable[t] = executor
}

func (s *Scheduler) ExecuteJobManual(rs *logy.RequestSession, t job.JobType) error {
	j := s.repo.FindManualJob(rs, t)
	if j == nil {
		return errors.New("jobtype not found")
	}
	if j.Execution == job.OneTime {
		return errors.New("unexpected job type")
	}
	return s.executeJob(context.Background(), rs, j)
}

func (s *Scheduler) RerunOnetime(rs *logy.RequestSession, key string) error {
	j := s.repo.FindByKey(rs, key, false)
	if j == nil {
		return errors.New("jobtype not found")
	}
	if j.Execution != job.OneTime {
		return errors.New("unexpected job type")
	}
	return s.executeOneTimeJob(rs, j, true)
}

func (s *Scheduler) ExecuteOneTimeJob(rs *logy.RequestSession, name string, t job.JobType, config any) (string, error) {
	jc, err := json.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("marshalling config: %w", err)
	}
	j := job.Job{
		RootEntity: domain.NewRootEntity(),
		Name:       name,
		JobType:    t,
		Status:     job.InProgress,
		Execution:  job.OneTime,
		MultiExec:  false,
		Config:     string(jc),
	}
	s.repo.Save(rs, &j)
	if err := s.executeOneTimeJob(rs, &j, false); err != nil {
		return "", fmt.Errorf("executing job: %w", err)
	}
	return j.Key, nil
}

func (s *Scheduler) executeOneTimeJob(rs *logy.RequestSession, j *job.Job, rerun bool) error {
	if j.Status == job.Success {
		return errors.New("job was successfully")
	}
	if rerun && time.Now().Before(j.Updated.Add(oneTimeTimeout)) && j.Status != job.Failure {
		return ErrNotTimedout
	}
	executor, found := s.jobCbTable[j.JobType]
	if !found {
		return fmt.Errorf("no callback for job %s specified", j.Name)
	}
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		logy.Infof(s.rs, "executing job %s", j.Name)
		exceptionHappened := false
		var res ExecutionResult
		exception.TryCatch(func() {
			res = executor.Execute(rs, *j)
		}, func(exception exception.Exception) {
			res.Log.AddEntry(job.Error, "exception: %s", exception.ToString())
			exceptionHappened = true
		})
		if exceptionHappened {
			j.Status = job.Failure
		} else {
			j.Status = job.Success
		}
		j.CustomRes = res.CustomRes
		if !res.Success {
			j.Status = job.Failure
		}
		j.AddLog(res.Log)
		s.repo.Update(s.rs, j)
	}()

	return nil
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(dbLookupInterval)

	completed := s.valkeyClient.B().Subscribe().Channel(conf.Config.Cache.Channel).Build()
	localCtx, cancelLocal := context.WithCancel(ctx)
	subCh := make(chan string)
	subErrCh := make(chan error)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		err := s.valkeyClient.Receive(localCtx, completed, func(msg valkey.PubSubMessage) {
			to := time.NewTicker(pubTimeout)
			select {
			case subCh <- msg.Message:
			case <-to.C:
			}
		})
		if err != nil {
			subErrCh <- err
		}
	}()

	logy.Infof(s.rs, "Entering scheduler loop")
	for {
		select {
		case <-ticker.C:
			s.executePeriodics(ctx)
		case msg := <-subCh:
			s.handleMessage(msg)
		case err := <-subErrCh:
			exception.HandleErrorServerMessage(err, message.GetI18N(message.SchedulerPubSub))
		case <-ctx.Done():
			cancelLocal()
			<-subErrCh
			return
		case <-s.stopCh:
			cancelLocal()
			<-subErrCh
			return

		}
	}
}

func (s *Scheduler) executePeriodics(ctx context.Context) {
	var jobs []*job.Job
	exception.TryCatchAndLog(s.rs, func() {
		jobs = s.repo.FindPeriodicJobs(s.rs)
	})
	for _, job := range jobs {
		if !s.periodicScheduled(job) {
			continue
		}
		err := s.executeJob(ctx, nil, job)
		if err != nil {
			if errors.Is(err, errAlreadyRunning) {
				logy.Infof(s.rs, "%s", err)
			} else {
				logy.Errorf(s.rs, "error executing periodic task: %s", err)
			}
		}
	}
}

// Does have to return an error because it is used in the scheduler goroutine which does not have a recoverer set in
// callstack.
func (s *Scheduler) executeJob(ctx context.Context, rs *logy.RequestSession, j *job.Job) error {
	if rs == nil {
		rs = s.rs
	}
	var (
		l        locks.Lock
		acquired bool
	)
	exception.TryCatchAndLog(rs, func() {
		l, acquired = s.lockService.Acquire(locks.Options{
			// We only use the jobtype here to ensure periodic and manual executions
			// of the same job dont run in parallel.
			Key:      "job" + strconv.Itoa(int(j.JobType)),
			Blocking: false,
		})
	})
	if !acquired {
		return errAlreadyRunning
	}
	fresh := s.repo.FindByKey(rs, j.Key, false)
	if fresh == nil {
		s.lockService.Release(l)
		return errors.New("can't refresh job")
	}
	if fresh.Updated != j.Updated {
		logy.Infof(s.rs, "job %s got executed already", j.Name)
		s.lockService.Release(l)
		return nil
	}
	if j.MultiExec {
		logy.Infof(s.rs, "publishing signal for %s", j.Name)
		cmd := s.valkeyClient.B().Publish().Channel(conf.Config.Cache.Channel).Message(strconv.Itoa(int(j.JobType))).Build()
		vRes := s.valkeyClient.Do(ctx, cmd)
		if err := vRes.Error(); err != nil {
			s.lockService.Release(l)
			return fmt.Errorf("error publishing signal %w", err)
		}
		j.Status = job.InProgress
		exception.TryCatchAndLog(rs, func() {
			s.repo.Update(s.rs, j)
		})
		s.lockService.Release(l)
		return nil
	}
	executor, found := s.jobCbTable[j.JobType]
	if !found {
		s.lockService.Release(l)
		return fmt.Errorf("no callback for job %s specified", j.Name)
	}
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer func() {
			s.lockService.Release(l)
		}()
		logy.Infof(s.rs, "executing job %s", j.Name)
		lastRun := j.Updated
		j.Status = job.InProgress
		s.repo.Update(s.rs, j)
		j.Updated = lastRun
		exceptionHappened := false
		var res ExecutionResult
		exception.TryCatch(func() {
			res = executor.Execute(rs, *j)
		}, func(e exception.Exception) {
			res.Log.AddEntry(job.Error, "exception: %s", e.ToString())
			exception.LogException(s.rs, e)
			exceptionHappened = true
		})
		if exceptionHappened {
			j.Status = job.Failure
		} else {
			j.Status = job.Success
		}
		j.CustomRes = res.CustomRes
		if res.NewConf != nil {
			j.Config = *res.NewConf
		}
		if !res.Success {
			j.Status = job.Failure
		}
		j.AddLog(res.Log)
		s.repo.Update(s.rs, j)
	}()

	return nil
}

func (s *Scheduler) periodicScheduled(job *job.Job) bool {
	if job.Schedule == "" {
		return false
	}

	if !s.gron.IsValid(job.Schedule) {
		logy.Errorf(s.rs, "invalid schedule %s for job %s", job.Schedule, job.Name)
		return false
	}

	nextExecution, err := gronx.NextTickAfter(job.Schedule, job.Updated, false)
	if err != nil {
		logy.Errorf(s.rs, "cant get next execution time for job: %s for %s", job.Schedule, job.Name)
		return false
	}
	return time.Now().After(nextExecution)
}

func (s *Scheduler) handleMessage(msg string) {
	logy.Infof(s.rs, "handling message: %s", msg)
	typeI, err := strconv.Atoi(msg)
	if err != nil {
		logy.Errorf(s.rs, "malformed message %s", err)
		return
	}
	if typeI < 0 || typeI > int(job.TypeLimit) {
		logy.Errorf(s.rs, "invalid job type %d", typeI)
		return
	}
	var j *job.Job
	exception.TryCatchAndLog(s.rs, func() {
		j = s.repo.FindLatestJob(s.rs, job.JobType(typeI))
	})
	if j == nil {
		logy.Errorf(s.rs, "job type not found in db %d", typeI)
		return
	}
	if !j.MultiExec {
		logy.Errorf(s.rs, "job multiexec is not set %s", j.Name)
		return
	}
	executor, found := s.jobCbTable[j.JobType]
	if !found {
		logy.Errorf(s.rs, "no callback for job %s specified", j.Name)
		return
	}
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		var res ExecutionResult
		exception.TryCatch(func() {
			res = executor.Execute(s.rs, *j)
		}, func(exception exception.Exception) {
			res.Log.AddEntry(job.Error, "exception: %s", exception.ToString())
		})
		var (
			l        locks.Lock
			acquired bool
		)
		exception.TryCatchAndLog(s.rs, func() {
			l, acquired = s.lockService.Acquire(locks.Options{
				// We only use the jobtype here to ensure periodic and manual executions
				// of the same job dont run in parallel.
				Key:      "job" + strconv.Itoa(int(j.JobType)),
				Blocking: true,
				Timeout:  jobLockTTL,
			})
		})
		if !acquired {
			logy.Errorf(s.rs, "could not acquire job lock for writing results %s", j.Name)
			return
		}
		var j *job.Job
		exception.TryCatchAndLog(s.rs, func() {
			j = s.repo.FindLatestJob(s.rs, job.JobType(typeI))
		})
		if j.Status == job.InProgress {
			j.Status = job.Success
			if !res.Success {
				j.Status = job.Failure
			}
		} else if j.Status == job.Success && !res.Success {
			j.Status = job.Failure
		}
		j.AddLog(res.Log)
		exception.TryCatchAndLog(s.rs, func() {
			s.repo.Update(s.rs, j)
		})
		s.lockService.Release(l)
	}()
}
