// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/newsbox"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/roles"
	"github.com/eclipse-disuko/disuko/helper/validation"
	"github.com/eclipse-disuko/disuko/infra/repository/database"
	newsboxRepo "github.com/eclipse-disuko/disuko/infra/repository/newsbox"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type NewsboxHandler struct {
	NewsboxRepo newsboxRepo.IRepo
}

func extractNewsboxItemBody(r *http.Request) newsbox.ItemDto {
	var dto newsbox.ItemDto
	validation.DecodeAndValidate(r, &dto, false)
	return dto
}

func (h *NewsboxHandler) PostHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowNewsbox.Create {
		exception.ThrowExceptionSendDeniedResponse()
	}
	dto := extractNewsboxItemBody(r)
	i := dto.ToEntity()
	i.RootEntity = domain.NewRootEntity()
	h.NewsboxRepo.Save(requestSession, &i)
}

func (h *NewsboxHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowNewsbox.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}
	qc := &database.QueryConfig{}
	qc.SetSort(database.SortConfig{
		database.SortAttribute{
			Name:  "Created",
			Order: database.ASC,
		},
	})
	all := h.NewsboxRepo.Query(requestSession, qc)
	render.JSON(w, r, domain.ToDtos(all))
}

func (h *NewsboxHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowNewsbox.Update {
		exception.ThrowExceptionSendDeniedResponse()
	}

	id := chi.URLParam(r, "id")
	existing := h.NewsboxRepo.FindByKey(requestSession, id, false)
	if existing == nil {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbNotFound), id+" not found in DB")
	}
	dto := extractNewsboxItemBody(r)
	i := dto.ToEntity()
	i.RootEntity = existing.RootEntity
	h.NewsboxRepo.Update(requestSession, &i)
	dto = i.ToDto()
	render.JSON(w, r, dto)
}

func (h *NewsboxHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowNewsbox.Delete {
		exception.ThrowExceptionSendDeniedResponse()
	}
	existing := h.NewsboxRepo.FindByKey(requestSession, id, false)
	if existing == nil {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbNotFound), id+" not found in DB")
	}
	h.NewsboxRepo.Delete(requestSession, existing.Key)
}
