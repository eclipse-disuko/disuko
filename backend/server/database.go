// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"github.com/eclipse-disuko/disuko/infra/repository/policydecisions"

	userstatsRepo "github.com/eclipse-disuko/disuko/infra/repository/userstats"

	"github.com/eclipse-disuko/disuko/infra/repository/changeloglist"
	"github.com/eclipse-disuko/disuko/infra/repository/changelogs"
	"github.com/eclipse-disuko/disuko/infra/repository/checklist"
	"github.com/eclipse-disuko/disuko/infra/repository/customid"
	filtersets "github.com/eclipse-disuko/disuko/infra/repository/filterset"
	"github.com/eclipse-disuko/disuko/infra/repository/internalToken"
	"github.com/eclipse-disuko/disuko/infra/repository/licenserules"
	"github.com/eclipse-disuko/disuko/infra/repository/newsbox"
	reviewremarks2 "github.com/eclipse-disuko/disuko/infra/repository/reviewtemplates"
	"github.com/eclipse-disuko/disuko/infra/service/startup"

	"github.com/eclipse-disuko/disuko/infra/repository/analyticsoccurrences"

	"github.com/eclipse-disuko/disuko/infra/repository/analytics"
	"github.com/eclipse-disuko/disuko/infra/repository/analyticscomponents"
	"github.com/eclipse-disuko/disuko/infra/repository/analyticslicenses"
	announcement "github.com/eclipse-disuko/disuko/infra/repository/announcements"
	"github.com/eclipse-disuko/disuko/infra/repository/approvallist"
	"github.com/eclipse-disuko/disuko/infra/repository/auditloglist"
	"github.com/eclipse-disuko/disuko/infra/repository/department"
	"github.com/eclipse-disuko/disuko/infra/repository/dpconfig"
	"github.com/eclipse-disuko/disuko/infra/repository/jobs"
	"github.com/eclipse-disuko/disuko/infra/repository/labels"
	"github.com/eclipse-disuko/disuko/infra/repository/license"
	migration "github.com/eclipse-disuko/disuko/infra/repository/migration"
	"github.com/eclipse-disuko/disuko/infra/repository/obligation"
	"github.com/eclipse-disuko/disuko/infra/repository/policyrules"
	projectRepo "github.com/eclipse-disuko/disuko/infra/repository/project"
	"github.com/eclipse-disuko/disuko/infra/repository/reviewremarks"
	"github.com/eclipse-disuko/disuko/infra/repository/sbomlist"
	schema2 "github.com/eclipse-disuko/disuko/infra/repository/schema"
	"github.com/eclipse-disuko/disuko/infra/repository/spdx_license"
	"github.com/eclipse-disuko/disuko/infra/repository/statistic"
	"github.com/eclipse-disuko/disuko/infra/repository/user"
	"github.com/eclipse-disuko/disuko/logy"
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
