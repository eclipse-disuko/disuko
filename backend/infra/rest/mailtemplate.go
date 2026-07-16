// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"
	"net/url"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/mailtemplate"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/roles"
	"github.com/eclipse-disuko/disuko/helper/validation"
	"github.com/eclipse-disuko/disuko/infra/repository/mailtemplates"
	user2 "github.com/eclipse-disuko/disuko/infra/repository/user"
	"github.com/eclipse-disuko/disuko/infra/service/mail"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type MailTemplateHandler struct {
	MailTemplatesRepository mailtemplates.IMailTemplatesRepository
	UserRepository          user2.IUsersRepository
	MailService             *mail.Service
}

func (h *MailTemplateHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.IsApplicationAdmin() {
		exception.ThrowExceptionSendDeniedResponse()
	}

	all := h.MailTemplatesRepository.FindAll(requestSession, false)
	render.JSON(w, r, domain.ToDtos(all))
}

func (h *MailTemplateHandler) GetByIdHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.IsApplicationAdmin() {
		exception.ThrowExceptionSendDeniedResponse()
	}

	idEscaped := chi.URLParam(r, "id")
	id, err := url.QueryUnescape(idEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ErrorQueryUnescape))

	tmpl := h.MailTemplatesRepository.FindByKey(requestSession, id, false)
	render.JSON(w, r, tmpl.ToDto())
}

func (h *MailTemplateHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.IsApplicationAdmin() {
		exception.ThrowExceptionSendDeniedResponse()
	}

	idEscaped := chi.URLParam(r, "id")
	id, err := url.QueryUnescape(idEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ErrorQueryUnescape))

	var req mailtemplate.UpdateMailTemplateDto
	validation.DecodeAndValidate(r, &req, false)

	tmpl := h.MailTemplatesRepository.FindByKey(requestSession, id, false)
	tmpl.Subject = req.Subject
	tmpl.Message = req.Message
	tmpl.Bcc = req.Bcc
	tmpl.Cc = req.Cc

	h.MailTemplatesRepository.Update(requestSession, tmpl)
	render.JSON(w, r, tmpl.ToDto())
}

func (h *MailTemplateHandler) TestHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	userName, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.IsApplicationAdmin() {
		exception.ThrowExceptionSendDeniedResponse()
	}

	idEscaped := chi.URLParam(r, "id")
	id, err := url.QueryUnescape(idEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ErrorQueryUnescape))

	var req mailtemplate.TestMailTemplateDto
	validation.DecodeAndValidate(r, &req, false)

	currentUser := GetUserByUsername(requestSession, h.UserRepository, userName)
	if currentUser == nil {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorPermissionDeniedUser, userName))
	}

	err = h.MailService.SendTestMail(requestSession, currentUser.Email, mailtemplate.MailTemplateKey(id), req.Message)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorSendingMail))

	w.WriteHeader(http.StatusOK)
}
