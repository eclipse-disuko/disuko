// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/policydecisions"

	userstatsRepo "mercedes-benz.ghe.com/foss/disuko/infra/repository/userstats"

	"mercedes-benz.ghe.com/foss/disuko/infra/repository/changeloglist"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/changelogs"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/checklist"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/customid"
	filtersets "mercedes-benz.ghe.com/foss/disuko/infra/repository/filterset"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/internalToken"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/licenserules"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/newsbox"
	reviewremarks2 "mercedes-benz.ghe.com/foss/disuko/infra/repository/reviewtemplates"
	"mercedes-benz.ghe.com/foss/disuko/infra/service/startup"

	"mercedes-benz.ghe.com/foss/disuko/infra/repository/analyticsoccurrences"

	"mercedes-benz.ghe.com/foss/disuko/infra/repository/analytics"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/analyticscomponents"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/analyticslicenses"
	announcement "mercedes-benz.ghe.com/foss/disuko/infra/repository/announcements"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/approvallist"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/auditloglist"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/department"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/dpconfig"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/jobs"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/labels"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/license"
	migration "mercedes-benz.ghe.com/foss/disuko/infra/repository/migration"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/obligation"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/policyrules"
	projectRepo "mercedes-benz.ghe.com/foss/disuko/infra/repository/project"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/reviewremarks"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/sbomlist"
	schema2 "mercedes-benz.ghe.com/foss/disuko/infra/repository/schema"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/spdx_license"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/statistic"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/user"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type dbRepos struct {
	obligation           obligation.IObligationRepository
	project              projectRepo.IProjectRepository
	migration            migration.IMigrationRepository
	schema               schema2.ISchemaRepository
	licenses             license.ILicensesRepository
	analytics            analytics.IAnalyticsRepository
	analyticsComponents  analyticscomponents.IComponentsRepository
	analyticsLicenses    analyticslicenses.ILicensesRepository
	analyticsOccurrences analyticsoccurrences.IOccurrencesRepository
	policyRules          policyrules.IPolicyRulesRepository
	statistic            statistic.IStatisticRepository
	job                  jobs.IJobsRepository
	label                labels.ILabelRepository
	dpConfig             *dpconfig.DBConfigRepository
	user                 user.IUsersRepository
	sbomList             sbomlist.ISbomListRepository
	auditLogList         auditloglist.IAuditLogListRepository
	department           department.IDepartmentRepository
	spdxLicense          spdx_license.ISpdxLicensesRepository
	approvalList         approvallist.IApprovalListRepository
	announcements        announcement.IAnnouncementsRepository
	reviewRemarks        reviewremarks.IReviewRemarksRepository
	filterSets           filtersets.IFilterSetsRepository
	reviewTemplate       reviewremarks2.IReviewTemplateRepository
	changeLogList        changeloglist.IChangeLogListRepository
	changeLogs           changelogs.IChangeLogsRepository
	licenseRules         licenserules.ILicenseRulesRepository
	basicauth            internalToken.IRepo
	customid             customid.ICustomIdRepository
	projectRepository    projectRepo.IProjectRepository
	checklist            checklist.IChecklistRepository
	newsbox              newsbox.IRepo
	userstats            userstatsRepo.IUserStatsRepository
	policyDecisions      policydecisions.IPolicyDecisionsRepository
}

func (s *Server) setupDatabase(requestSession *logy.RequestSession) {
	s.repos = dbRepos{
		obligation:           obligation.NewObligationRepository(requestSession),
		project:              projectRepo.NewProjectRepository(requestSession),
		migration:            migration.NewMigrationRepository(requestSession),
		schema:               schema2.NewSchemaRepository(requestSession),
		licenses:             license.NewLicenseRepository(requestSession),
		policyRules:          policyrules.NewPolicyRulesRepository(requestSession),
		analytics:            analytics.NewAnalyticsRepository(requestSession),
		analyticsComponents:  analyticscomponents.NewComponentsRepository(requestSession),
		analyticsLicenses:    analyticslicenses.NewLicensesRepository(requestSession),
		analyticsOccurrences: analyticsoccurrences.NewLicensesRepository(requestSession),
		statistic:            statistic.NewSystemStatisticRepository(requestSession),
		job:                  jobs.NewJobsRepository(requestSession),
		label:                labels.NewLabelsRepository(requestSession),
		dpConfig:             dpconfig.NewDbConfigRepository(requestSession),
		user:                 user.NewUsersRepository(requestSession),
		sbomList:             sbomlist.NewSbomListRepository(requestSession),
		auditLogList:         auditloglist.NewAuditLogListRepository(requestSession),
		department:           department.NewDepartmentRepository(requestSession),
		spdxLicense:          spdx_license.NewSpdxLicenseRepository(requestSession),
		approvalList:         approvallist.NewApprovalListRepository(requestSession),
		announcements:        announcement.NewAnnouncementsRepository(requestSession),
		reviewRemarks:        reviewremarks.NewReviewRemarskRepositry(requestSession),
		filterSets:           filtersets.NewFilterSetsRepository(requestSession),
		reviewTemplate:       reviewremarks2.NewReviewTemplateRepositry(requestSession),
		changeLogList:        changeloglist.NewChangeLogListRepository(requestSession),
		changeLogs:           changelogs.NewChangeLogsRepository(requestSession),
		licenseRules:         licenserules.NewLicenseRulesRepository(requestSession),
		basicauth:            internalToken.NewRepo(requestSession),
		customid:             customid.NewLabelsRepository(requestSession),
		projectRepository:    projectRepo.NewProjectRepository(requestSession),
		checklist:            checklist.NewLabelsRepository(requestSession),
		newsbox:              newsbox.NewNewsboxRepository(requestSession),
		userstats:            userstatsRepo.NewUsersRepository(requestSession),
		policyDecisions:      policydecisions.NewPolicyDecisionsRepository(requestSession),
	}
	err := s.repos.seedDb(requestSession)
	if err != nil {
		logy.Fatalf(requestSession, err.Error())
	}
	go s.repos.analyticsComponents.InitIndex(requestSession)
	go s.repos.analyticsLicenses.InitIndex(requestSession)
}

func (s *Server) migrateDatabase(requestSession *logy.RequestSession, ext ...startup.Step) {
	s.handlers.startUp.MigrateDatabase(requestSession, ext...)
}
