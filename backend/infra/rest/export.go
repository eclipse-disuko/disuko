// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"

	"github.com/go-chi/render"
	export2 "mercedes-benz.ghe.com/foss/disuko/domain/export"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/helper/roles"
	"mercedes-benz.ghe.com/foss/disuko/helper/validation"
	"mercedes-benz.ghe.com/foss/disuko/infra/service/analytics"
	"mercedes-benz.ghe.com/foss/disuko/infra/service/export"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type ExportHandler struct {
	ExportService    *export.Service
	AnalyticsService *analytics.Analytics
}

func (handler *ExportHandler) ExportLicenseKnowledgeBase(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowTools.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}

	result := handler.ExportService.ExportLicenseKnowledgeBase(requestSession)
	render.JSON(w, r, result)
}

func (handler *ExportHandler) ExportSchemaKnowledgeBase(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowTools.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}

	result := handler.ExportService.ExportSchemaKnowledgeBase(requestSession)
	render.JSON(w, r, result)
}

func (handler *ExportHandler) ImportLicenseKnowledgeBase(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowTools.Update {
		exception.ThrowExceptionSendDeniedResponse()
	}

	validation.CheckExpectedContentType(r, validation.ContentTypeFormData)

	file, fileHandler, err := r.FormFile("file")
	exception.HandleErrorClientMessage(err, message.GetI18N(message.Error))
	defer file.Close()

	validation.CheckExpectedContentType2(fileHandler.Header, []validation.ContentType{
		validation.ContentTypeJson,
		validation.ContentTypeOctets,
	})

	data := &export2.ExportLicenseKnowledgeBaseDto{}
	validation.DecodePartAndValidate(file, data, false)

	handler.ExportService.ImportLicenseKnowledgeBase(requestSession, data)
	go handler.AnalyticsService.Reinitialise(requestSession)

	render.JSON(w, r, export2.ImportResultDto{
		Message: "import successful",
		Success: true,
	})
}

func (handler *ExportHandler) ImportSchemaKnowledgeBase(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowTools.Update {
		exception.ThrowExceptionSendDeniedResponse()
	}

	validation.CheckExpectedContentType(r, validation.ContentTypeFormData)

	file, fileHandler, err := r.FormFile("file")
	exception.HandleErrorClientMessage(err, message.GetI18N(message.Error))
	defer file.Close()

	validation.CheckExpectedContentType2(fileHandler.Header, []validation.ContentType{
		validation.ContentTypeJson,
		validation.ContentTypeOctets,
	})

	data := &export2.ExportSchemaKnowledgeBaseDto{}
	validation.DecodePartAndValidate(file, data, false)

	handler.ExportService.ImportSchemaKnowledgeBase(requestSession, data)

	render.JSON(w, r, export2.ImportResultDto{
		Message: "import successful",
		Success: true,
	})
}
