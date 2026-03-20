// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package spdxsubscribe

import (
	"mercedes-benz.ghe.com/foss/disuko/conf"
	"mercedes-benz.ghe.com/foss/disuko/helper/mail"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/user"
	"mercedes-benz.ghe.com/foss/disuko/logy"
	"mercedes-benz.ghe.com/foss/disuko/observermngmt"
)

type SpdxSubscribe struct {
	mailClient mail.Client
	userRepo   user.IUsersRepository
}

func Init(mailClient mail.Client, userRepo user.IUsersRepository) *SpdxSubscribe {
	return &SpdxSubscribe{
		mailClient: mailClient,
		userRepo:   userRepo,
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
	mailData := struct {
		Username    string
		ProjectName string
		ProjectLink string
		VersionName string
		VersionLink string
		SbomLink    string
	}{}

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
			templ := "spdxUploaded"
			err := s.mailClient.Send(targetUser.Email, templ, mailData)
			if err != nil {
				logy.Errorf(data.RequestSession, "Failed to send the email: %v", err)
			} else {
				logy.Infof(data.RequestSession, "Email sent successfully!")
			}
		}()
	}
}
