// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"
	"time"

	filterSetDomain "github.com/eclipse-disuko/disuko/domain/filterset"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/roles"
	"github.com/eclipse-disuko/disuko/helper/validation"
	filtersets "github.com/eclipse-disuko/disuko/infra/repository/filterset"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type FilterSetHandler struct {
	FilterSetsRepository filtersets.IFilterSetsRepository
}

func (filterSetHandler *FilterSetHandler) FilterSetsGetByTableHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	tableName := chi.URLParam(r, "tablename")
	filterSets := filterSetHandler.FilterSetsRepository.FindByTableName(requestSession, tableName)

	render.JSON(w, r, filterSetDomain.ToDtos(filterSets))
}

func (filterSetHandler *FilterSetHandler) FilterSetGetHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	id := chi.URLParam(r, "id")
	filterSet := filterSetHandler.FilterSetsRepository.FindByKey(requestSession, id, false)
	if filterSet == nil {
		exception.ThrowExceptionServer404(message.GetI18N(message.FilterSetNotFound))
	}
	render.JSON(w, r, filterSetDomain.ToDto(filterSet))
}

func (filterSetHandler *FilterSetHandler) FilterSetPostHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowLicense.Create {
		exception.ThrowExceptionSendDeniedResponse()
	}

	var filterSetDto filterSetDomain.FilterSetDto
	validation.DecodeAndValidate(r, &filterSetDto, false)
	filterSetEntity := filterSetDto.ToEntity()

	filterSetHandler.FilterSetsRepository.Save(requestSession, filterSetEntity)

	render.JSON(w, r, filterSetDomain.ToDto(filterSetEntity))
}

func (filterSetHandler *FilterSetHandler) FilterSetUpdateHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowLicense.Update {
		exception.ThrowExceptionSendDeniedResponse()
	}
	var filterSetDto filterSetDomain.FilterSetDto
	validation.DecodeAndValidate(r, &filterSetDto, false)

	id := chi.URLParam(r, "id")
	filterSetEntity := filterSetHandler.FilterSetsRepository.FindByKey(requestSession, id, false)

	if filterSetEntity == nil {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorDbNotFound, "filter not found in DB"))
	}
	filterSetEntity.Name = filterSetDto.Name
	filterSetEntity.IncludedFilters = filterSetDto.IncludedFilters
	filterSetEntity.ExcludedFilters = filterSetDto.ExcludedFilters
	filterSetEntity.Updated = time.Now()

	filterSetHandler.FilterSetsRepository.Update(requestSession, filterSetEntity)

	render.JSON(w, r, filterSetDomain.SuccessResponse{
		Success: true,
		Message: "Filterset updated",
	})
}

func (filterSetHandler *FilterSetHandler) FilterSetDeleteHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowLicense.Delete {
		exception.ThrowExceptionSendDeniedResponse()
	}
	id := chi.URLParam(r, "id")
	filterSetHandler.FilterSetsRepository.Delete(requestSession, id)
	w.WriteHeader(http.StatusOK)
}
