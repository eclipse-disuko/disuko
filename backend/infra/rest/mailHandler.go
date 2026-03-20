// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"

	"github.com/eclipse-disuko/disuko/conf"
	mail2 "github.com/eclipse-disuko/disuko/domain/mail"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/mail"
	"github.com/eclipse-disuko/disuko/helper/roles"
	"github.com/eclipse-disuko/disuko/helper/validation"
	"github.com/eclipse-disuko/disuko/logy"
)

type MailHandler struct {
	Client mail.Client
}

type MailData struct {
	Username string
	Link     string
}

var recipient string = "disclosure-portal@mercedes-benz.com"

func extractRequestMailBody(r *http.Request) mail2.MailRequstDto {
	var projectData mail2.MailRequstDto
	validation.DecodeAndValidate(r, &projectData, false)
	return projectData
}

func (handler *MailHandler) SendMail(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.IsDomainAdmin()) {
		exception.ThrowExceptionSendDeniedResponse()
	}

	request := extractRequestMailBody(r)

	if !handler.Client.IsTeamplateValid(request.MailType) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	mailData := MailData{
		Username: "Superman",
	}

	mailData.Link = conf.Config.Server.DisukoHost + "/#/dashboard/home"

	go func() {
		defer func() {
			if err := recover(); err != nil {
				logy.Errorf(requestSession, "Could not send email %v", err)
			}
		}()

		err := handler.Client.Send(recipient, "taskReview", mailData)
		if err != nil {
			logy.Errorf(requestSession, "Failed to send the email: %v", err)
		} else {
			logy.Infof(requestSession, "Email sent successfully!")
		}
	}()

	w.WriteHeader(http.StatusOK)
}
