// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"

	export2 "github.com/eclipse-disuko/disuko/domain/export"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/roles"
	"github.com/eclipse-disuko/disuko/helper/validation"
	"github.com/eclipse-disuko/disuko/infra/service/analytics"
	"github.com/eclipse-disuko/disuko/infra/service/export"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/go-chi/render"
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
