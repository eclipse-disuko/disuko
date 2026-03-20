// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"

	"github.com/eclipse-disuko/disuko/domain/notification"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/roles"
	"github.com/eclipse-disuko/disuko/helper/validation"
	"github.com/eclipse-disuko/disuko/infra/repository/dpconfig"

	"github.com/eclipse-disuko/disuko/logy"
	"github.com/go-chi/render"
)

type NotificationHandler struct {
	DpConfigRepo *dpconfig.DBConfigRepository
}

func (handler *NotificationHandler) NotificationGetHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	response := handler.DpConfigRepo.Notification.Get(requestSession)

	render.JSON(w, r, response)
}

func (handler *NotificationHandler) NotificationSetHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowTools.Create || !rights.AllowTools.Update || !rights.AllowTools.Delete {
		exception.ThrowExceptionSendDeniedResponse()
	}

	var notificationData notification.NotificationDto
	validation.DecodeAndValidate(r, &notificationData, false)

	handler.DpConfigRepo.Notification.Save(requestSession, &notification.Notification{
		Text:    notificationData.Text,
		Enabled: notificationData.Enabled,
	})

	w.WriteHeader(200)
}
