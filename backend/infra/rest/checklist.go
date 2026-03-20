// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"mercedes-benz.ghe.com/foss/disuko/domain"
	"mercedes-benz.ghe.com/foss/disuko/domain/checklist"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/helper/roles"
	"mercedes-benz.ghe.com/foss/disuko/helper/validation"
	checklistRepo "mercedes-benz.ghe.com/foss/disuko/infra/repository/checklist"
	checklistService "mercedes-benz.ghe.com/foss/disuko/infra/service/checklist"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type ChecklistHandler struct {
	ChecklistRepo    checklistRepo.IChecklistRepository
	ChecklistService checklistService.Service
}

func extractChecklistBody(r *http.Request) checklist.ChecklistDto {
	var dto checklist.ChecklistDto
	validation.DecodeAndValidate(r, &dto, false)
	return dto
}

func extractChecklistItemBody(r *http.Request) checklist.ItemDto {
	var dto checklist.ItemDto
	validation.DecodeAndValidate(r, &dto, false)
	return dto
}

func (h *ChecklistHandler) PostHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowChecklist.Create {
		exception.ThrowExceptionSendDeniedResponse()
	}
	dto := extractChecklistBody(r)
	cl := dto.ToEntity()
	cl.RootEntity = domain.NewRootEntity()
	cl.Items = []checklist.Item{}
	h.ChecklistRepo.Save(requestSession, &cl)
}

func (h *ChecklistHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowChecklist.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}
	all := h.ChecklistRepo.FindAll(requestSession, false)
	render.JSON(w, r, domain.ToDtos(all))
}

func (h *ChecklistHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowChecklist.Update {
		exception.ThrowExceptionSendDeniedResponse()
	}

	id := chi.URLParam(r, "id")
	existing := h.ChecklistRepo.FindByKey(requestSession, id, false)
	if existing == nil {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbNotFound), id+" not found in DB")
	}

	dto := extractChecklistBody(r)
	cl := dto.ToEntity()

	cl.RootEntity = existing.RootEntity
	cl.Items = existing.Items
	h.ChecklistRepo.Update(requestSession, &cl)

	dto = cl.ToDto()
	render.JSON(w, r, dto)
}

func (h *ChecklistHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowChecklist.Delete {
		exception.ThrowExceptionSendDeniedResponse()
	}

	id := chi.URLParam(r, "id")
	existing := h.ChecklistRepo.FindByKey(requestSession, id, false)
	if existing == nil {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbNotFound), id+" not found in DB")
	}
	if existing.Active {
		exception.ThrowExceptionSendDeniedResponse()
	}
	h.ChecklistRepo.Delete(requestSession, id)
	w.WriteHeader(200)
}

func (h *ChecklistHandler) CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowChecklist.Update {
		exception.ThrowExceptionSendDeniedResponse()
	}

	id := chi.URLParam(r, "id")
	cl := h.ChecklistRepo.FindByKey(requestSession, id, false)
	if cl == nil {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbNotFound), id+" not found in DB")
	}

	dto := extractChecklistItemBody(r)
	item := dto.ToEntity()
	item.ChildEntity = domain.NewChildEntity()
	cl.Items = append(cl.Items, item)

	h.ChecklistRepo.Update(requestSession, cl)

	clDto := cl.ToDto()
	render.JSON(w, r, clDto)
}

func (h *ChecklistHandler) UpdateItemHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowChecklist.Update {
		exception.ThrowExceptionSendDeniedResponse()
	}

	id := chi.URLParam(r, "id")
	cl := h.ChecklistRepo.FindByKey(requestSession, id, false)
	if cl == nil {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbNotFound), id+" not found in DB")
	}

	var (
		itemId    = chi.URLParam(r, "itemId")
		itemIndex = -1
	)
	for i, item := range cl.Items {
		if item.Key == itemId {
			itemIndex = i
			break
		}
	}
	if itemIndex == -1 {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbNotFound), itemId+" not found in DB")
	}

	dto := extractChecklistItemBody(r)
	item := dto.ToEntity()
	item.ChildEntity = cl.Items[itemIndex].ChildEntity
	item.ChildEntity.Updated = time.Now()
	cl.Items[itemIndex] = item
	h.ChecklistRepo.Update(requestSession, cl)

	clDto := cl.ToDto()
	render.JSON(w, r, clDto)
}

func (h *ChecklistHandler) DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowChecklist.Update {
		exception.ThrowExceptionSendDeniedResponse()
	}

	id := chi.URLParam(r, "id")
	cl := h.ChecklistRepo.FindByKey(requestSession, id, false)
	if cl == nil {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbNotFound), id+" not found in DB")
	}
	var (
		itemId    = chi.URLParam(r, "itemId")
		itemIndex = -1
	)
	for i, item := range cl.Items {
		if item.Key == itemId {
			itemIndex = i
			break
		}
	}
	if itemIndex == -1 {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbNotFound), itemId+" not found in DB")
	}
	cl.Items = append(cl.Items[:itemIndex], cl.Items[itemIndex+1:]...)
	h.ChecklistRepo.Update(requestSession, cl)

	clDto := cl.ToDto()
	render.JSON(w, r, clDto)
}
