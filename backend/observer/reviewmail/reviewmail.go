// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package reviewmail

import (
	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain/overallreview"
	"github.com/eclipse-disuko/disuko/helper/mail"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/infra/repository/user"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/eclipse-disuko/disuko/observermngmt"
)

type ReviewSubscribe struct {
	mailClient mail.Client
	userRepo   user.IUsersRepository
}

func Init(mailClient mail.Client, userRepo user.IUsersRepository) *ReviewSubscribe {
	return &ReviewSubscribe{
		mailClient: mailClient,
		userRepo:   userRepo,
	}
}

func (r *ReviewSubscribe) RegisterHandlers() {
	observermngmt.RegisterHandler(observermngmt.OverallReviewCreated, r.OnReviewAdded)
}

func (r *ReviewSubscribe) OnReviewAdded(eventId observermngmt.EventId, arg interface{}) {
	data, ok := arg.(observermngmt.OverallReviewData)
	if !ok {
		return
	}
	r.sendMail(data)
}

// Make adjustable for MB usage
func translateState(state overallreview.State) (string, string) {
	switch state {
	case overallreview.Unreviewed:
		return message.GetI18N(message.StatusReviewUnreviewed).Text, message.GetI18N(message.StatusReviewUnreviewedDE).Text
	case overallreview.Audited:
		return message.GetI18N(message.StatusReviewAudited).Text, message.GetI18N(message.StatusReviewAuditedDE).Text
	case overallreview.Acceptable:
		return message.GetI18N(message.StatusReviewAcceptable).Text, message.GetI18N(message.StatusReviewAcceptableDE).Text
	case overallreview.NotAcceptable:
		return message.GetI18N(message.StatusReviewNotAcceptable).Text, message.GetI18N(message.StatusReviewNotAcceptableDE).Text
	case overallreview.AcceptableAfterChanges:
		return message.GetI18N(message.StatusReviewAcceptableAfterChanges).Text, message.GetI18N(message.StatusReviewAcceptableAfterChangesDE).Text
	}
	return "", ""
}

func (r *ReviewSubscribe) sendMail(data observermngmt.OverallReviewData) {
	if data.Project.IsGroup {
		return
	}
	mailData := struct {
		Username         string
		ProjectName      string
		ReviewerFullName string
		Status           string
		StatusDE         string
		Comment          string
		ReviewsLink      string
		ProjectLink      string
		VersionName      string
		VersionLink      string
	}{}

	mailData.ProjectName = data.Project.Name
	projectTypeUrlPart := "projects"
	if data.Project.IsGroup {
		projectTypeUrlPart = "groups"
	}
	mailData.ProjectLink = conf.Config.Server.DisukoHost + "/#/dashboard/" + projectTypeUrlPart + "/" + data.Project.Key

	mailData.VersionName = data.Version.Name
	mailData.VersionLink = mailData.ProjectLink + "/versions/" + data.Version.Key
	mailData.ReviewsLink = mailData.VersionLink + "/sbomQuality/reviewRemarks"
	mailData.Status, mailData.StatusDE = translateState(data.Review.State)
	mailData.ReviewerFullName = data.Review.CreatorFullName
	mailData.Comment = data.Review.Comment

	for _, u := range data.Project.UserManagement.Users {
		if !u.Subscriptions.OverallReview {
			continue
		}
		targetUser := r.userRepo.FindByUserId(data.RequestSession, u.UserId)
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
			templ := "reviewCreated"
			err := r.mailClient.Send(targetUser.Email, templ, mailData)
			if err != nil {
				logy.Errorf(data.RequestSession, "Failed to send the email: %v", err)
			} else {
				logy.Infof(data.RequestSession, "Email sent successfully!")
			}
		}()
	}
}
