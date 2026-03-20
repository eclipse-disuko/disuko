// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"

	"mercedes-benz.ghe.com/foss/disuko/domain/job"
	"mercedes-benz.ghe.com/foss/disuko/domain/statistic"
	userstatsDomain "mercedes-benz.ghe.com/foss/disuko/domain/userstats"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/helper/roles"
	userstatsRepo "mercedes-benz.ghe.com/foss/disuko/infra/repository/userstats"
	"mercedes-benz.ghe.com/foss/disuko/jobs/userstats"
	"mercedes-benz.ghe.com/foss/disuko/scheduler"

	"github.com/go-chi/render"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/labels"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/license"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/newsbox"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/obligation"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/policyrules"
	project2 "mercedes-benz.ghe.com/foss/disuko/infra/repository/project"
	rt "mercedes-benz.ghe.com/foss/disuko/infra/repository/reviewtemplates"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/schema"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/user"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type CountHandler struct {
	ProjectRepository        project2.IProjectRepository
	LicenseRepository        license.ILicensesRepository
	PolicyRulesRepository    policyrules.IPolicyRulesRepository
	LabelRepository          labels.ILabelRepository
	SchemaRepository         schema.ISchemaRepository
	ObligationRepository     obligation.IObligationRepository
	UserRepository           user.IUsersRepository
	ReviewTemplateRepository rt.IReviewTemplateRepository
	UserStatsRepository      userstatsRepo.IUserStatsRepository
	Scheduler                *scheduler.Scheduler
	NewsboxRepository        newsbox.IRepo
}

func (countHandler *CountHandler) GetDashboardCountsHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	userName, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	userStats := countHandler.UserStatsRepository.FindByUserId(requestSession, userName)
	var counts *statistic.DashboardCounts
	if userStats != nil {
		counts = userStats.ToDashboardCounts()
	} else {
		counts = &statistic.DashboardCounts{
			ProjectCount:        -1,
			LicenseCount:        -1,
			PolicyRuleCount:     -1,
			LabelCount:          -1,
			SchemaCount:         -1,
			ObligationCount:     -1,
			UserCount:           -1,
			DisclosureCount:     -1,
			ReviewTemplateCount: -1,
			ActiveJobCount:      -1,
		}
		_, err := countHandler.Scheduler.ExecuteOneTimeJob(requestSession, "calculate user stats", job.CalculateUserStats, userstats.Config{Username: userName, Rights: rights})
		exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorStartingJob))
	}

	render.JSON(w, r, counts)
}

func (countHandler *CountHandler) GetDashboardCountsForAdminHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	userStats := countHandler.UserStatsRepository.FindByUserId(requestSession, userstatsDomain.SystemStatsUser)
	var counts *statistic.DashboardCounts
	if userStats != nil {
		counts = userStats.ToDashboardCounts()
	} else {
		counts = &statistic.DashboardCounts{
			ProjectCount:        -1,
			LicenseCount:        -1,
			PolicyRuleCount:     -1,
			LabelCount:          -1,
			SchemaCount:         -1,
			ObligationCount:     -1,
			UserCount:           -1,
			DisclosureCount:     -1,
			ReviewTemplateCount: -1,
			ActiveJobCount:      -1,
		}
		_, err := countHandler.Scheduler.ExecuteOneTimeJob(requestSession, "calculate user stats", job.CalculateUserStats, userstats.Config{Username: userstatsDomain.SystemStatsUser, Rights: rights, AdminRequest: true})
		exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorStartingJob))
	}

	render.JSON(w, r, counts)
}
