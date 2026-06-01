// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/policyruleclassification"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/roles"
	"github.com/eclipse-disuko/disuko/helper/validation"
	"github.com/eclipse-disuko/disuko/infra/repository/obligation"
	prcRepo "github.com/eclipse-disuko/disuko/infra/repository/policyruleclassification"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type PolicyRuleClassificationHandler struct {
	PolicyRuleClassificationRepo prcRepo.IPolicyRuleClassificationRepository
	ObligationRepo               obligation.IObligationRepository
}

func (h *PolicyRuleClassificationHandler) findByKeyOrThrow(requestSession *logy.RequestSession, id string) *policyruleclassification.PolicyRuleClassification {
	existing := h.PolicyRuleClassificationRepo.FindByKey(requestSession, id, false)
	if existing == nil {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbNotFound), id+" not found in DB")
	}
	return existing
}

func (h *PolicyRuleClassificationHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowTools.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}
	all := h.PolicyRuleClassificationRepo.FindAll(requestSession, false)
	render.JSON(w, r, domain.ToDtos(all))
}

func (h *PolicyRuleClassificationHandler) GetMatrixHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowTools.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}

	all := h.PolicyRuleClassificationRepo.FindAll(requestSession, false)
	obligations := h.ObligationRepo.FindAll(requestSession, false)

	classifications := make([]policyruleclassification.ClassificationInfo, 0, len(obligations))
	for _, o := range obligations {
		classifications = append(classifications, policyruleclassification.ClassificationInfo{
			Key:  o.Key,
			Name: o.Name,
		})
	}

	useCaseDtos := domain.ToDtos(all)

	response := policyruleclassification.MatrixResponseDto{
		Classifications: classifications,
		UseCases:        useCaseDtos,
	}
	render.JSON(w, r, response)
}

func (h *PolicyRuleClassificationHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowLicense.Create {
		exception.ThrowExceptionSendDeniedResponse()
	}

	var dto policyruleclassification.PolicyRuleClassificationRequestDto
	validation.DecodeAndValidate(r, &dto, false)

	entity := dto.ToEntity()
	entity.RootEntity = domain.NewRootEntity()
	h.PolicyRuleClassificationRepo.Save(requestSession, &entity)
	render.JSON(w, r, entity.ToDto())
}

func (h *PolicyRuleClassificationHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowLicense.Update {
		exception.ThrowExceptionSendDeniedResponse()
	}

	id := chi.URLParam(r, "id")
	existing := h.findByKeyOrThrow(requestSession, id)

	var dto policyruleclassification.PolicyRuleClassificationRequestDto
	validation.DecodeAndValidate(r, &dto, false)

	existing.Update(dto.Name, dto.Rules)
	h.PolicyRuleClassificationRepo.Update(requestSession, existing)
	render.JSON(w, r, existing.ToDto())
}

func (h *PolicyRuleClassificationHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowLicense.Delete {
		exception.ThrowExceptionSendDeniedResponse()
	}

	id := chi.URLParam(r, "id")
	existing := h.findByKeyOrThrow(requestSession, id)
	h.PolicyRuleClassificationRepo.Delete(requestSession, existing.Key)
	w.WriteHeader(200)
}
