// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"
	"slices"
	"strconv"

	"github.com/eclipse-disuko/disuko/helper/s3Helper"
	"github.com/eclipse-disuko/disuko/jobs/report"

	"github.com/eclipse-disuko/disuko/domain/analytics"
	da "github.com/eclipse-disuko/disuko/domain/analytics"
	"github.com/eclipse-disuko/disuko/domain/internalToken"
	license2 "github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/roles"
	"github.com/eclipse-disuko/disuko/helper/validation"
	ra "github.com/eclipse-disuko/disuko/infra/repository/analytics"
	"github.com/eclipse-disuko/disuko/infra/repository/database"
	"github.com/eclipse-disuko/disuko/infra/repository/license"
	"github.com/eclipse-disuko/disuko/infra/repository/policyrules"
	projectRepo "github.com/eclipse-disuko/disuko/infra/repository/project"
	"github.com/eclipse-disuko/disuko/infra/repository/sbomlist"
	"github.com/eclipse-disuko/disuko/infra/repository/statistic"
	sa "github.com/eclipse-disuko/disuko/infra/service/analytics"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/go-chi/render"
)

const unreferencedChart = "unreferenced"

type AnalyticsHandler struct {
	ProjectRepository    projectRepo.IProjectRepository
	LicenseRepository    license.ILicensesRepository
	PolicyRuleRepository policyrules.IPolicyRulesRepository
	AnalyticsRepository  ra.IAnalyticsRepository
	SbomListRepository   sbomlist.ISbomListRepository
	AnalyticsService     sa.Analytics
	StatisticRepository  statistic.IStatisticRepository
}

func (handler *AnalyticsHandler) AnalyticsSearchHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowProject.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}

	searchOptions := extractAnalyticsSearchRequestBody(r)

	response := tryBruteForce(handler, requestSession, searchOptions)

	render.JSON(w, r, response)
}

func (handler *AnalyticsHandler) Report(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowProject.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}

	w.Header().Set("Content-Type", "application/octet-stream")

	s3Helper.PerformDownload(requestSession, &w, report.GetReportStorageFileNameOf(report.GetReportAllName()), "")
}

func (handler *AnalyticsHandler) InternalReport(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	auth := extractInternalToken(r.Context())
	if !slices.Contains(auth.Capabilities, internalToken.StatisticsCSV) {
		exception.ThrowExceptionSendDeniedResponse()
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\"disco_dump.csv\"")
	s3Helper.PerformDownload(requestSession, &w, report.GetReportStorageFileNameOf(report.GetReportAllName()), "")
}

func (handler *AnalyticsHandler) Statistic(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowProject.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}

	qc := database.New().SetMatcher(
		database.AttributeMatcher(
			"meta.approvalState",
			database.EQ,
			"forbidden",
		),
	).SetKeep([]string{
		handler.LicenseRepository.DatabaseConn().GetKeyAttribute(),
	})
	lRes := handler.LicenseRepository.Query(requestSession, qc)
	forbidden := len(lRes)

	qc = database.New().SetMatcher(
		database.AttributeMatcher(
			"meta.approvalState",
			database.EQ,
			"",
		),
	).SetKeep([]string{
		handler.LicenseRepository.DatabaseConn().GetKeyAttribute(),
	})
	lRes = handler.LicenseRepository.Query(requestSession, qc)
	unknown := len(lRes)

	qc = database.New().SetSort(database.SortConfig{
		database.SortAttribute{
			Name:  "Created",
			Order: database.DESC,
		},
	}).SetLimit(0, 1)
	qbRes := handler.StatisticRepository.Query(requestSession, qc)
	if len(qbRes) == 0 {
		exception.ThrowExceptionBadRequestResponse()
	}

	res := analytics.StatsDto{
		ProjectCount:        qbRes[0].ProjectCount,
		ProjectActiveCount:  qbRes[0].ProjectActiveCount,
		ProjectDeletedCount: qbRes[0].ProjectDeletedCount,

		LicenseCount:        qbRes[0].LicenseCount,
		LicenseChartCount:   qbRes[0].LicenseChartCount,
		LicenseActiveCount:  qbRes[0].LicenseActiveCount,
		LicenseDeletedCount: qbRes[0].LicenseDeletedCount,
		LicenseForbiden:     forbidden,
		LicenseUnknown:      unknown,

		UploadFileCntSBOM: qbRes[0].UploadedFilesCntSBOM,

		UserCount:                 qbRes[0].UserCount,
		UserActiveCount:           qbRes[0].UserActiveCount,
		UserTermsNotAcceptedCount: qbRes[0].UserTermsNotAcceptedCount,
		UserDeactivateCount:       qbRes[0].UserDeactivateCount,

		CompletedTrainings: 0, // TODO: fetch this
	}

	render.JSON(w, r, res)
}

func (handler *AnalyticsHandler) LicenseOccurrences(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowProject.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}

	os := handler.AnalyticsService.Handler.Occurrences(requestSession)
	licCache := make(map[string]*license2.License)
	var list []analytics.OccurrenceDto

	possibleCharts := make(map[string]int)
	possibleSources := make(map[string]int)
	possibleFamilies := make(map[string]int)
	possibleApproval := map[string]int{
		license2.NotSet:     0,
		license2.Pending:    0,
		license2.Check:      0,
		license2.Assigning:  0,
		license2.Approved:   0,
		license2.Forbidden:  0,
		license2.Deprecated: 0,
	}
	possibleType := make(map[string]int)

	for _, o := range os {
		if o.ReferencedLicense == "" {
			possibleCharts[unreferencedChart]++
			list = append(list, *o.ToDto(nil))
			continue
		}
		lic, found := licCache[o.ReferencedLicense]
		if !found {
			lic = handler.LicenseRepository.FindById(requestSession, o.ReferencedLicense)
			if lic == nil {
				logy.Infof(requestSession, "lic not found %#v", o)
				continue
			}
			licCache[o.ReferencedLicense] = lic
		}
		possibleCharts[strconv.FormatBool(lic.Meta.IsLicenseChart)]++
		possibleSources[string(lic.Source)]++
		possibleFamilies[string(lic.Meta.Family)]++
		possibleApproval[string(lic.Meta.ApprovalState)]++
		possibleType[string(lic.Meta.LicenseType)]++

		list = append(list, *o.ToDto(lic))
	}
	render.JSON(w, r, analytics.OccurrencesResDto{
		List: list,
		PossibleValues: license2.PossibleFilterValues{
			PossibleCharts:   possibleCharts,
			PossibleSources:  possibleSources,
			PossibleFamilies: possibleFamilies,
			PossibleApproval: possibleApproval,
			PossibleType:     possibleType,
		},
	})
}

func tryBruteForce(handler *AnalyticsHandler, requestSession *logy.RequestSession, searchOptions da.RequestSearchOptions) interface{} {
	response := da.ResponseAnalyticsSearch{
		Success: false,
		Items:   make([]da.SearchResponseItem, 0),
		Stats: da.Statistic{
			CountProjects:   0,
			CountComponents: 0,
			CountVersions:   0,
		},
	}

	qc := database.New().SetMatcher(
		database.AndChain(
			database.AttributeMatcher(
				"Deleted",
				database.EQ,
				false,
			),
			database.AttributeMatcher(
				"Name",
				database.LIKE,
				"%"+searchOptions.Component+"%",
			),
		),
	)
	projects := handler.ProjectRepository.Query(requestSession, qc)
	response.Stats.CountProjects = len(projects)
	response.Stats.TotalProjects = handler.ProjectRepository.CountAll(requestSession)

	for _, prjKey := range projects {
		prj := handler.ProjectRepository.FindByKey(requestSession, prjKey.Key, true)
		response.Items = append(response.Items, da.SearchResponseItem{
			Name: prj.Name,
			Type: da.PROJECT,
		})
	}
	response.Success = true
	return response
}

func extractAnalyticsSearchRequestBody(r *http.Request) da.RequestSearchOptions {
	var searchOptions da.RequestSearchOptions
	validation.DecodeAndValidate(r, &searchOptions, false)
	return searchOptions
}
