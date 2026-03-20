// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"mercedes-benz.ghe.com/foss/disuko/domain"
	"mercedes-benz.ghe.com/foss/disuko/domain/department"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/helper/roles"
	departmentRepo "mercedes-benz.ghe.com/foss/disuko/infra/repository/department"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type DepartmentHandler struct {
	Repo departmentRepo.IDepartmentRepository
}

func (handler *DepartmentHandler) Find(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	//check rights
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.IsInternal || !rights.AllowProject.Create && !rights.AllowProject.Update {
		exception.ThrowExceptionSendDeniedResponse()
	}

	searchStrEscaped := chi.URLParam(r, "searchStr")
	searchStr, err := url.QueryUnescape(searchStrEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ErrorUnexpectError))

	if len(searchStr) < 3 {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorTooFewCharacters), "")
	}

	entities := handler.Repo.FindBySearchStr(requestSession, searchStr)

	dtos := domain.MapToLimit(entities, func(source *department.Department) *department.DepartmentDto {
		return source.ToDto()
	}, 100)

	render.JSON(w, r, dtos)
}
