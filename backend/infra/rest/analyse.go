// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"
	"strconv"
	"time"

	"github.com/eclipse-disuko/disuko/infra/repository/approvallist"

	integrity2 "github.com/eclipse-disuko/disuko/domain/integrity"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/roles"
	"github.com/eclipse-disuko/disuko/infra/repository/dpconfig"
	"github.com/eclipse-disuko/disuko/infra/repository/project"
	"github.com/eclipse-disuko/disuko/infra/repository/sbomlist"
	"github.com/eclipse-disuko/disuko/infra/service/integrity"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/go-chi/render"
)

type AnalyseFilesHandler struct {
	ProjectRepository project.IProjectRepository
	DpConfigRepo      *dpconfig.DBConfigRepository
	SbomListRepo      sbomlist.ISbomListRepository
	ApprovalListRepo  approvallist.IApprovalListRepository
}

func (countHandler *AnalyseFilesHandler) AnalyseFilesHandlerStart(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowTools.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}
	fixIt, err := strconv.ParseBool(r.URL.Query().Get("fixIt"))
	if err != nil {
		fixIt = false
	}

	state := integrity.LoadDbIntegrityResult(requestSession, countHandler.DpConfigRepo)

	if !state.IsRunning {
		state = &integrity2.DbIntegrityResult{}
		state.IsRunning = true
		state.ReqID = requestSession.ReqID
		state.FixIt = fixIt
		state.StartTime = time.Now()

		integrity.SaveDbIntegrityResult(requestSession, state, countHandler.DpConfigRepo)

		exception.RunAsyncAndLogException(requestSession, func() {
			integrity.AnalyseDataIntegrity(requestSession, countHandler.ProjectRepository, countHandler.DpConfigRepo,
				countHandler.SbomListRepo, countHandler.ApprovalListRepo, fixIt)
		})
	}

	render.Status(r, http.StatusOK)
}

func (countHandler *AnalyseFilesHandler) AnalyseFilesHandlerStop(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowTools.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}

	state := integrity.LoadDbIntegrityResult(requestSession, countHandler.DpConfigRepo)

	if state.IsRunning {
		state.IsRunning = false
		integrity.SaveDbIntegrityResult(requestSession, state, countHandler.DpConfigRepo)
	}

	render.Status(r, http.StatusOK)
}

func (countHandler *AnalyseFilesHandler) GetResultHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowTools.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}

	state := integrity.LoadDbIntegrityResult(requestSession, countHandler.DpConfigRepo)

	render.JSON(w, r, state)
}
