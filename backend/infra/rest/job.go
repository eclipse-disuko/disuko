// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"errors"
	"net/http"
	"strconv"

	job2 "github.com/eclipse-disuko/disuko/domain/job"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/roles"
	"github.com/eclipse-disuko/disuko/helper/validation"
	"github.com/eclipse-disuko/disuko/infra/repository/jobs"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/eclipse-disuko/disuko/scheduler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type JobHandler struct {
	JobRepository jobs.IJobsRepository
	Scheduler     *scheduler.Scheduler
}

func (jobHandler *JobHandler) JobGetAllHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.IsApplicationAdmin() {
		exception.ThrowExceptionSendDeniedResponse()
	}

	allJobs := jobHandler.JobRepository.FindAll(requestSession, false)
	render.JSON(w, r, job2.ToDtoList(allJobs))
}

func (jobHandler *JobHandler) JobTriggerRun(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowTools.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}
	jobType := chi.URLParam(r, "jobType")
	typeInt, err := strconv.Atoi(jobType)
	if err != nil {
		exception.ThrowExceptionBadRequestResponse()
	}

	err = jobHandler.Scheduler.ExecuteJobManual(requestSession, job2.JobType(typeInt))
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorStartingJob))
}

func (jobHandler *JobHandler) JobRerunOnetime(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowTools.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}
	key := chi.URLParam(r, "key")
	err := jobHandler.Scheduler.RerunOnetime(requestSession, key)
	msg := message.ErrorStartingJob
	if errors.Is(err, scheduler.ErrNotTimedout) {
		msg = message.ErrorNotTimedout
	}
	exception.HandleErrorServerMessage(err, message.GetI18N(msg))
}

func (jobHandler *JobHandler) JobGetLatestByTypeHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.AllowTools.Read || (rights.AllowUsers.Create && rights.AllowUsers.Read && rights.AllowUsers.Update && rights.AllowUsers.Delete)) {
		exception.ThrowExceptionSendDeniedResponse()
	}
	jobTypeStr := chi.URLParam(r, "jobType")

	jobType, err := strconv.Atoi(jobTypeStr)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.Error, jobTypeStr))

	job := jobHandler.JobRepository.FindLatestJob(requestSession, job2.JobType(jobType))

	render.JSON(w, r, job2.ToDto(job))
}

func extractSetConfigBody(r *http.Request) (res job2.SetConfigDto) {
	validation.DecodeAndValidate(r, &res, false)
	return res
}

func (jobHandler *JobHandler) SetConfig(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowTools.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}
	jobType := chi.URLParam(r, "jobType")
	typeInt, err := strconv.Atoi(jobType)
	if err != nil {
		exception.ThrowExceptionBadRequestResponse()
	}
	reqData := extractSetConfigBody(r)

	j := jobHandler.JobRepository.FindManualJob(requestSession, job2.JobType(typeInt))
	if j == nil {
		exception.ThrowExceptionBadRequestResponse()
	}
	j.Config = reqData.Config
	jobHandler.JobRepository.Update(requestSession, j)
	render.JSON(w, r, SuccessResponse{
		Success: true,
	})
}
