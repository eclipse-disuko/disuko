// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/go-cmp/cmp"
	"mercedes-benz.ghe.com/foss/disuko/domain/audit"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/helper/validation"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/auditloglist"
	obligation2 "mercedes-benz.ghe.com/foss/disuko/infra/repository/obligation"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"mercedes-benz.ghe.com/foss/disuko/domain/obligation"
	"mercedes-benz.ghe.com/foss/disuko/helper/roles"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/license"
	"mercedes-benz.ghe.com/foss/disuko/infra/service/obligationcsv"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type ObligationsHandler struct {
	ObligationRepository   obligation2.IObligationRepository
	LicenseRepository      license.ILicensesRepository
	AuditLogListRepository auditloglist.IAuditLogListRepository
}

func (handler *ObligationsHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	username, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowObligation.Update {
		exception.ThrowExceptionSendDeniedResponse()
	}
	request := extractObligationRequestBody(r)

	alreadyExistingWithSameName := handler.ObligationRepository.FindByName(requestSession, request.Name)
	if len(alreadyExistingWithSameName) > 1 || (len(alreadyExistingWithSameName) == 1 && alreadyExistingWithSameName[0].Key != request.Key) {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ObligationAlreadyExists))
	}

	current := handler.ObligationRepository.FindByKey(requestSession, request.Key, false)
	oldObligationAudit := current.ToAudit()

	current.Update(request)

	newObligationAudit := current.ToAudit()
	handler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, current.Key, username, message.ClassificationUpdated, audit.DiffWithReporter, newObligationAudit, oldObligationAudit)

	handler.ObligationRepository.Update(requestSession, current)

	w.WriteHeader(200)
	result := SuccessResponse{
		Success: true,
		Message: "done",
	}

	render.JSON(w, r, result)
}

func (handler *ObligationsHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	obligations := handler.ObligationRepository.FindAll(requestSession, false)

	var result obligation.AllResponse
	result.Obligation = obligation.ToDtos(obligations)
	result.Count = len(obligations)

	render.JSON(w, r, result)
}

func (handler *ObligationsHandler) PostHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	username, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowObligation.Create {
		exception.ThrowExceptionSendDeniedResponse()
	}
	requestItem := extractObligationRequestBody(r)

	if len(requestItem.Name) <= 0 {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DataMissingObligation))
	}

	alreadyExistingWithSameName := handler.ObligationRepository.FindByName(requestSession, requestItem.Name)
	if len(alreadyExistingWithSameName) > 0 {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ObligationAlreadyExists))
	}

	newItem := obligation.CreateFrom(requestItem)

	obligationAudit := newItem.ToAudit()
	handler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, newItem.Key, username, message.ClassificationCreated, audit.DiffWithReporter, obligationAudit, &obligation.ObligationAudit{})

	handler.ObligationRepository.Save(requestSession, newItem)

	w.WriteHeader(200)
	render.JSON(w, r, newItem.ToDto())
}

func (handler *ObligationsHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	username, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowObligation.Delete {
		exception.ThrowExceptionSendDeniedResponse()
	}
	currentItem := retrieveObligationForKey(requestSession, handler.ObligationRepository, r)

	// Find usages of the classification over licenses by the key.
	licenses := handler.LicenseRepository.FindByObligationKey(requestSession, currentItem.Key)

	// perform update every license by removing obligation key
	if len(licenses) > 0 {
		// TODO let this for future use!
		const userHasRightsToDeleteEvenIsUsed = false
		if userHasRightsToDeleteEvenIsUsed {
			for _, license := range licenses {
				obligationsKeyList := make([]string, 0)
				for _, obligationKey := range license.Meta.ObligationsKeyList {
					if obligationKey != currentItem.Key {
						obligationsKeyList = append(obligationsKeyList, obligationKey)
					}
				}
				license.Meta.ObligationsKeyList = obligationsKeyList
				// TODO add License Audit Entry as License is changing
				handler.LicenseRepository.Update(requestSession, license)
			}
		} else {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorInUse, "Classification"))
		}
	}

	obligationAudit := currentItem.ToAudit()
	handler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, currentItem.Key, username, message.ClassificationDeleted, cmp.Diff, obligation.ObligationAudit{}, obligationAudit)

	handler.ObligationRepository.Delete(requestSession, currentItem.Key)

	w.WriteHeader(200)

	fmt.Fprintf(w, "Successfully deleted obligation with id %s \n", currentItem.Key)
}

func (handler *ObligationsHandler) GetByIdHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	uuidEscaped := chi.URLParam(r, "uuid")
	uuid, err := url.QueryUnescape(uuidEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ErrorQueryUnescape))
	obligation := handler.ObligationRepository.FindByKey(requestSession, uuid, false)

	render.JSON(w, r, obligation.ToDto())
}

func retrieveObligationForKey(requestSession *logy.RequestSession, repo obligation2.IObligationRepository, r *http.Request) *obligation.Obligation {
	labelKey := chi.URLParam(r, "id")
	itemFromDatabase := repo.FindByKey(requestSession, labelKey, false)
	if len(itemFromDatabase.Key) == 0 {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ObligationKeyMissing, labelKey), "")
	}
	return itemFromDatabase
}

func extractObligationRequestBody(r *http.Request) *obligation.ObligationDto {
	var item obligation.ObligationDto
	validation.DecodeAndValidate(r, &item, false)
	return &item
}

func (handler *ObligationsHandler) CreateCSVHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowObligation.Read || !rights.AllowLicense.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}
	obligationcsv.CreateCSV(&w, requestSession, handler.ObligationRepository, handler.LicenseRepository)
}

func (handler *ObligationsHandler) GetAuditTrailHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.IsLicenseManager() || rights.IsDomainAdmin()) {
		exception.ThrowExceptionSendDeniedResponse()
	}

	idEscaped := chi.URLParam(r, "id")
	id, err := url.QueryUnescape(idEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ErrorQueryUnescape))

	auditList := handler.AuditLogListRepository.FindByKey(requestSession, id, false)

	auditTrail := make([]audit.AuditDto, 0)
	if auditList != nil && auditList.AuditTrail != nil {
		for _, item := range auditList.AuditTrail {
			auditTrail = append(auditTrail, item.ToDto())
		}
	}
	render.JSON(w, r, auditTrail)
}
