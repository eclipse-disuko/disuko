// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"

	"mercedes-benz.ghe.com/foss/disuko/domain/notification"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/roles"
	"mercedes-benz.ghe.com/foss/disuko/helper/validation"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/dpconfig"

	"github.com/go-chi/render"
	"mercedes-benz.ghe.com/foss/disuko/logy"
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
