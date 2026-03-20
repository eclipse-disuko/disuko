// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"
	"strconv"
	"time"

	"mercedes-benz.ghe.com/foss/disuko/infra/repository/approvallist"

	"github.com/go-chi/render"
	integrity2 "mercedes-benz.ghe.com/foss/disuko/domain/integrity"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/roles"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/dpconfig"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/project"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/sbomlist"
	"mercedes-benz.ghe.com/foss/disuko/infra/service/integrity"
	"mercedes-benz.ghe.com/foss/disuko/logy"
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
