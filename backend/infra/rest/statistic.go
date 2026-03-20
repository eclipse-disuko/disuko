// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/go-chi/render"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/statistic"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/roles"
	"github.com/eclipse-disuko/disuko/helper/s3Helper"
	"github.com/eclipse-disuko/disuko/infra/repository/database"
	"github.com/eclipse-disuko/disuko/infra/repository/dpconfig"
	"github.com/eclipse-disuko/disuko/infra/repository/labels"
	"github.com/eclipse-disuko/disuko/infra/repository/license"
	"github.com/eclipse-disuko/disuko/infra/repository/obligation"
	"github.com/eclipse-disuko/disuko/infra/repository/policyrules"
	"github.com/eclipse-disuko/disuko/infra/repository/project"
	"github.com/eclipse-disuko/disuko/infra/repository/schema"
	statistic2 "github.com/eclipse-disuko/disuko/infra/repository/statistic"
	"github.com/eclipse-disuko/disuko/infra/repository/user"
	"github.com/eclipse-disuko/disuko/logy"
)

type StatisticHandler struct {
	StatisticRepository   statistic2.IStatisticRepository
	ProjectRepository     project.IProjectRepository
	LicensesRepository    license.ILicensesRepository
	PolicyRulesRepository policyrules.IPolicyRulesRepository
	ObligationRepository  obligation.IObligationRepository
	SchemaRepository      schema.ISchemaRepository
	LabelRepository       labels.ILabelRepository
	UsersRepository       user.IUsersRepository
	DpConfigRepo          *dpconfig.DBConfigRepository
}

func (handler *StatisticHandler) GetSystemStats(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	logy.Infof(requestSession, "GetSystemStats - START")
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)

	if !rights.AllowTools.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}

	dtoResults := handler.getStatisticsForWeb(requestSession)

	if len(dtoResults.DayStats) < 1 {
		handler.UpdateSystemStats(w, r)
	} else {
		render.JSON(w, r, dtoResults)
	}

}

func (handler *StatisticHandler) getStatisticsForWeb(requestSession *logy.RequestSession) *statistic.SystemStatsResponseDto {
	// last 5 days

	qc := database.New().SetSort(database.SortConfig{
		database.SortAttribute{
			Name:  "Created",
			Order: database.DESC,
		},
	}).SetLimit(0, 5)
	qbRes := handler.StatisticRepository.Query(requestSession, qc)

	var statistics []*statistic.SystemStatistic
	statistics = qbRes

	lastDaysStats := statistic.ToSystemStatisticDtoList(statistics)
	handler.checkStatistics(lastDaysStats)

	//last 12 months
	currentMonth := int(time.Now().Month())
	currentYear := time.Now().Year()
	last12MonthsStats := make([]*statistic.SystemStatisticDto, 0)
	for i := 0; i < 12; i++ {
		// calculate year and month
		year := currentYear
		month := currentMonth - i
		if month <= 0 {
			month += 12
			year = currentYear - 1
		}

		//create filter
		filter := strconv.Itoa(year) + "-"
		if month < 10 {
			filter += "0" + strconv.Itoa(month)
		} else {
			filter += strconv.Itoa(month)
		}
		filter += "-%"

		qc := database.New().SetMatcher(
			database.AttributeMatcher(
				"Created",
				database.LIKE,
				filter,
			),
		).SetSort(database.SortConfig{
			database.SortAttribute{
				Name:  "Created",
				Order: database.DESC,
			},
		}).SetLimit(0, 1)

		monthState := handler.StatisticRepository.Query(requestSession, qc)

		//search last statistic for this month
		//no statistic in this month
		if len(monthState) == 0 {
			continue
		}

		//store statistic for this month
		last12MonthsStats = append(last12MonthsStats, monthState[0].ToDto())
	}
	handler.checkStatistics(last12MonthsStats)

	return &statistic.SystemStatsResponseDto{
		DayStats:    lastDaysStats,
		MonthsStats: last12MonthsStats,
	}
}

func (handler *StatisticHandler) checkStatistics(stats []*statistic.SystemStatisticDto) {
	for i, statistic := range stats {
		if i >= (len(stats) - 1) {
			break
		}
		before := stats[i+1]
		statistic.MissingProjects = statistic.ProjectCount < before.ProjectCount
		statistic.MissingLicenses = statistic.LicenseCount < before.LicenseCount
		statistic.MissingPolicyRules = statistic.PolicyRuleCount < before.PolicyRuleCount
		statistic.MissingObligations = statistic.ObligationCount < before.ObligationCount
		statistic.MissingUploadFiles = statistic.UploadFileCnt < before.UploadFileCnt
		statistic.MissingUsers = statistic.UserCount < before.UserCount
	}
}

func (handler *StatisticHandler) UpdateSystemStats(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowTools.Update {
		exception.ThrowExceptionSendDeniedResponse()
	}
	handler.triggerUpdateSystemStats(requestSession)

	dtoResults := handler.getStatisticsForWeb(requestSession)

	render.JSON(w, r, dtoResults)
}

func (handler *StatisticHandler) TriggerUpdateSystemStatsHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	handler.triggerUpdateSystemStats(requestSession)
	render.Status(r, 200)
}

func (handler *StatisticHandler) triggerUpdateSystemStats(requestSession *logy.RequestSession) {
	logy.Infof(requestSession, "UpdateSystemStats - START")
	projectCount := handler.ProjectRepository.CountAllWithDeleted(requestSession)

	projectActiveCount := handler.ProjectRepository.CountAll(requestSession)

	licensesCount := handler.LicensesRepository.CountAllWithDeleted(requestSession)

	qc := database.New().SetMatcher(
		database.AttributeMatcher(
			"meta.IsLicenseChart",
			database.EQ,
			true,
		),
	)
	licensesChartCount := len(handler.LicensesRepository.Query(requestSession, qc))

	licensesActiveCount := handler.LicensesRepository.CountAll(requestSession)

	policiesCount := handler.PolicyRulesRepository.CountAllWithDeleted(requestSession)

	policiesActiveCount := handler.PolicyRulesRepository.CountAll(requestSession)

	labelCount := handler.LabelRepository.CountAll(requestSession)

	allSchemas := handler.SchemaRepository.CountAll(requestSession)

	obligationsCount := handler.ObligationRepository.CountAllWithDeleted(requestSession)

	// TODO: delete this, obligations are saved in harddelete repo
	obligationsActiveCount := handler.ObligationRepository.CountAll(requestSession)

	usersCount := handler.UsersRepository.CountAll(requestSession)

	qc.SetMatcher(
		database.AttributeMatcher(
			"Active",
			database.EQ,
			true,
		),
	)
	usersActiveCount := len(handler.UsersRepository.Query(requestSession, qc))

	qc.SetMatcher(
		database.AttributeMatcher(
			"Deprovisioned",
			database.NE,
			"0001-01-01T00:00:00Z",
		),
	)
	usersDeprovisionedCount := len(handler.UsersRepository.Query(requestSession, qc))

	qc.SetMatcher(
		database.AttributeMatcher(
			"Active",
			database.EQ,
			false,
		),
	)
	userDeactivateCount := len(handler.UsersRepository.Query(requestSession, qc))

	qc.SetMatcher(
		database.AttributeMatcher(
			"TermsOfUse",
			database.EQ,
			false,
		),
	)
	userTermsNotAcceptedCount := len(handler.UsersRepository.Query(requestSession, qc))

	//count over folders
	uploadedFilesCountByType := s3Helper.CountFiles(requestSession, conf.Config.Server.GetUploadPath())
	dbBackupFilesCount := s3Helper.CountFiles(requestSession, conf.Config.Server.BackupPath).CntFiles

	//count version over limit and max version count in one project
	allProjectKeysNotDeleted := handler.ProjectRepository.FindAllKeys(requestSession)
	maxVersionsInOneProject := 0
	projectsOverOrAtVersionLimit := 0
	for _, projectKey := range allProjectKeysNotDeleted {
		projectEntity := handler.ProjectRepository.FindByKey(requestSession, projectKey, false)

		if projectEntity == nil {
			continue
		}

		versionCount := 0
		for _, version := range projectEntity.Versions {
			if !version.Deleted {
				versionCount++
			}
		}
		if versionCount > maxVersionsInOneProject {
			maxVersionsInOneProject = versionCount
		}
		if versionCount >= conf.Config.Server.MaxVersions {
			projectsOverOrAtVersionLimit++
		}
	}

	updateStatistic := &statistic.SystemStatistic{
		RootEntity:                   domain.NewRootEntity(),
		ProjectCount:                 projectCount,
		ProjectActiveCount:           projectActiveCount,
		ProjectDeletedCount:          projectCount - projectActiveCount,
		LicenseCount:                 licensesCount,
		LicenseChartCount:            licensesChartCount,
		LicenseActiveCount:           licensesActiveCount,
		LicenseDeletedCount:          licensesCount - licensesActiveCount,
		PolicyRuleCount:              policiesCount,
		PolicyRuleActiveCount:        policiesActiveCount,
		PolicyRuleDeletedCount:       policiesCount - policiesActiveCount,
		LabelCount:                   labelCount,
		SchemaCount:                  allSchemas,
		ObligationCount:              obligationsCount,
		ObligationActiveCount:        obligationsActiveCount,
		ObligationDeletedCount:       obligationsCount - obligationsActiveCount,
		UploadedFilesCnt:             uploadedFilesCountByType.CntFiles,
		UploadedFilesCntPDF:          uploadedFilesCountByType.CntPDF,
		UploadedFilesCntJSON:         uploadedFilesCountByType.CntJson,
		UploadedFilesCntSBOM:         uploadedFilesCountByType.CntSBOM,
		DbBackupFilesCnt:             dbBackupFilesCount,
		UserCount:                    usersCount,
		UserActiveCount:              usersActiveCount,
		UserDeactivateCount:          userDeactivateCount,
		UserDeprovisionedCount:       usersDeprovisionedCount,
		UserTermsNotAcceptedCount:    userTermsNotAcceptedCount,
		ProjectsOverOrAtVersionLimit: projectsOverOrAtVersionLimit,
		MaxVersionsInOneProject:      maxVersionsInOneProject,
		VersionLimit:                 conf.Config.Server.MaxVersions,
	}

	//check changes
	lastStatistics := handler.getStatisticsForWeb(requestSession)
	if len(lastStatistics.DayStats) > 0 {
		lastStatistic := lastStatistics.DayStats[0]
		if lastStatistic.ProjectCount > updateStatistic.ProjectCount {
			logy.Errorf(requestSession, "Project count was reduced! It was: %d and is now: %d",
				lastStatistic.ProjectCount, updateStatistic.ProjectCount)
		}
		if lastStatistic.LicenseCount > updateStatistic.LicenseCount {
			logy.Errorf(requestSession, "License count was reduced! It was: %d and is now: %d",
				lastStatistic.LicenseCount, updateStatistic.LicenseCount)
		}
		if lastStatistic.PolicyRuleCount > updateStatistic.PolicyRuleCount {
			logy.Errorf(requestSession, "PolicyRule count was reduced! It was: %d and is now: %d",
				lastStatistic.PolicyRuleCount, updateStatistic.PolicyRuleCount)
		}
		if lastStatistic.ObligationCount > updateStatistic.ObligationCount {
			logy.Errorf(requestSession, "Obligation count was reduced! It was: %d and is now: %d",
				lastStatistic.ObligationCount, updateStatistic.ObligationCount)
		}
		if lastStatistic.UploadFileCnt > updateStatistic.UploadedFilesCnt {
			logy.Errorf(requestSession, "File count was reduced! It was: %d and is now: %d",
				lastStatistic.UploadFileCnt, updateStatistic.UploadedFilesCnt)
		}
		if lastStatistic.UserCount > updateStatistic.UserCount {
			logy.Errorf(requestSession, "User count was reduced! It was: %d and is now: %d",
				lastStatistic.UserCount, updateStatistic.UserCount)
		}
	}

	handler.StatisticRepository.Save(requestSession, updateStatistic)
}

func (handler *StatisticHandler) GetSystemProfileStats(w http.ResponseWriter, r *http.Request) {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)
	render.JSON(w, r, stats)
}
