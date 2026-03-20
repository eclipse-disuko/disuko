// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"

	"github.com/eclipse-disuko/disuko/domain/job"
	"github.com/eclipse-disuko/disuko/domain/statistic"
	userstatsDomain "github.com/eclipse-disuko/disuko/domain/userstats"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/roles"
	userstatsRepo "github.com/eclipse-disuko/disuko/infra/repository/userstats"
	"github.com/eclipse-disuko/disuko/jobs/userstats"
	"github.com/eclipse-disuko/disuko/scheduler"

	"github.com/eclipse-disuko/disuko/infra/repository/labels"
	"github.com/eclipse-disuko/disuko/infra/repository/license"
	"github.com/eclipse-disuko/disuko/infra/repository/newsbox"
	"github.com/eclipse-disuko/disuko/infra/repository/obligation"
	"github.com/eclipse-disuko/disuko/infra/repository/policyrules"
	project2 "github.com/eclipse-disuko/disuko/infra/repository/project"
	rt "github.com/eclipse-disuko/disuko/infra/repository/reviewtemplates"
	"github.com/eclipse-disuko/disuko/infra/repository/schema"
	"github.com/eclipse-disuko/disuko/infra/repository/user"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/go-chi/render"
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
