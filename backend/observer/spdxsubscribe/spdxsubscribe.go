// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package spdxsubscribe

import (
	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain/mailtemplate"
	"github.com/eclipse-disuko/disuko/infra/repository/user"
	"github.com/eclipse-disuko/disuko/infra/service/mail"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/eclipse-disuko/disuko/observermngmt"
)

type SpdxSubscribe struct {
	mailService *mail.Service
	userRepo    user.IUsersRepository
}

func Init(mailService *mail.Service, userRepo user.IUsersRepository) *SpdxSubscribe {
	return &SpdxSubscribe{
		mailService: mailService,
		userRepo:    userRepo,
	}
}

func (s *SpdxSubscribe) RegisterHandlers() {
	observermngmt.RegisterHandler(observermngmt.SpdxAdded, s.OnSpdxAdded)
}

func (s *SpdxSubscribe) OnSpdxAdded(eventId observermngmt.EventId, arg interface{}) {
	data, ok := arg.(observermngmt.SpdxData)
	if !ok {
		return
	}
	s.sendMail(data)
}

func (s *SpdxSubscribe) sendMail(data observermngmt.SpdxData) {
	mailData := mailtemplate.SpdxUploadedMailData{}

	mailData.ProjectName = data.Project.Name
	projectTypeUrlPart := "projects"
	if data.Project.IsGroup {
		projectTypeUrlPart = "groups"
	}
	mailData.ProjectLink = conf.Config.Server.DisukoHost + "/#/dashboard/" + projectTypeUrlPart + "/" + data.Project.Key

	mailData.VersionName = data.Version.Name
	mailData.VersionLink = mailData.ProjectLink + "/versions/" + data.Version.Key
	mailData.SbomLink = mailData.VersionLink + "/component/" + data.SpdxFile.Key

	for _, u := range data.Project.UserManagement.Users {
		if !u.Subscriptions.Spdx {
			continue
		}
		targetUser := s.userRepo.FindByUserId(data.RequestSession, u.UserId)
		if targetUser.Email == "" || !targetUser.Deprovisioned.IsZero() {
			continue
		}
		mailData.Username = targetUser.Forename + " " + targetUser.Lastname
		go func() {
			defer func() {
				if err := recover(); err != nil {
					logy.Errorf(data.RequestSession, "Could not send email %v", err)
				}
			}()
			err := s.mailService.SendMail(data.RequestSession, targetUser.Email, mailtemplate.MailTemplateKeySpdxUploaded, mailData)
			if err != nil {
				logy.Errorf(data.RequestSession, "Failed to send the email: %v", err)
			} else {
				logy.Infof(data.RequestSession, "Email sent successfully!")
			}
		}()
	}
}
