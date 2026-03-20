// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"

	"github.com/eclipse-disuko/disuko/domain"
	rt "github.com/eclipse-disuko/disuko/domain/reviewremarks"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/roles"
	"github.com/eclipse-disuko/disuko/helper/validation"
	checklistRepo "github.com/eclipse-disuko/disuko/infra/repository/checklist"
	reviewremarks "github.com/eclipse-disuko/disuko/infra/repository/reviewtemplates"
	"github.com/eclipse-disuko/disuko/infra/service/template"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/eclipse-disuko/disuko/observermngmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type TemplateHandler struct {
	ReviewTemplateRepository reviewremarks.IReviewTemplateRepository
	ChecklistRepository      checklistRepo.IChecklistRepository
}

func (templateHandler *TemplateHandler) CreateReviewTemplate(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	username, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowReviewTemplates.Create {
		exception.ThrowExceptionSendDeniedResponse()
	}
	body := extractReviewTemplateBody(r)
	reviewTemplate := &rt.ReviewTemplate{
		RootEntity:  domain.NewRootEntity(),
		Title:       body.Title,
		Description: body.Description,
		Level:       rt.Level(body.Level),
		Source:      body.Source,
	}
	templateHandler.ReviewTemplateRepository.Save(requestSession, reviewTemplate)

	observermngmt.FireEvent(observermngmt.DatabaseEntryAddedOrDeleted, observermngmt.DatabaseSizeChange{
		RequestSession: requestSession,
		CollectionName: reviewremarks.ReviewTemplateCollectionName,
		Rights:         rights,
		Username:       username,
	})

	w.WriteHeader(http.StatusCreated)
}

func (templateHandler *TemplateHandler) UpdateReviewTemplate(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowReviewTemplates.Update {
		exception.ThrowExceptionSendDeniedResponse()
	}
	updatedDto := extractReviewTemplateBody(r)

	id := chi.URLParam(r, "id")
	foundTemplate := templateHandler.ReviewTemplateRepository.FindByKey(requestSession, id, false)
	updatedEntity := foundTemplate.Update(updatedDto)

	templateHandler.ReviewTemplateRepository.Update(requestSession, updatedEntity)
	w.WriteHeader(http.StatusOK)
}

func (templateHandler *TemplateHandler) GetReviewTemplates(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowReviewTemplates.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}

	templates := templateHandler.ReviewTemplateRepository.FindAll(requestSession, false)
	templateDtos := make([]rt.ReviewTemplateResponseDto, 0)
	for _, reviewTemplate := range templates {
		templateDtos = append(templateDtos, *reviewTemplate.ToDto())
	}
	render.JSON(w, r, templateDtos)
}

func (templateHandler *TemplateHandler) GetReviewTemplate(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowReviewTemplates.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}
	id := chi.URLParam(r, "id")
	reviewTemplate := templateHandler.ReviewTemplateRepository.FindByKey(requestSession, id, false)
	render.JSON(w, r, reviewTemplate.ToDto())
}

func (templateHandler *TemplateHandler) DeleteReviewTemplate(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	username, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowReviewTemplates.Delete {
		exception.ThrowExceptionSendDeniedResponse()
	}
	id := chi.URLParam(r, "id")

	cls := templateHandler.ChecklistRepository.FindAll(requestSession, false)
	for _, cl := range cls {
		for _, i := range cl.Items {
			if i.TargetTemplateKey == id {
				exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorInUse, "checklist "+cl.Key))
			}
		}
	}

	templateHandler.ReviewTemplateRepository.Delete(requestSession, id)

	observermngmt.FireEvent(observermngmt.DatabaseEntryAddedOrDeleted, observermngmt.DatabaseSizeChange{
		RequestSession: requestSession,
		CollectionName: reviewremarks.ReviewTemplateCollectionName,
		Rights:         rights,
		Username:       username,
	})

	w.WriteHeader(http.StatusOK)
}

func (templateHandler *TemplateHandler) CreateCSVHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowReviewTemplates.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}
	template.CreateCSV(&w, requestSession, templateHandler.ReviewTemplateRepository)
}

func extractReviewTemplateBody(r *http.Request) rt.ReviewTemplateRequestDto {
	var requestDto rt.ReviewTemplateRequestDto
	validation.DecodeAndValidate(r, &requestDto, false)
	return requestDto
}
