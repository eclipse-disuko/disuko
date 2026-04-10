// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"github.com/eclipse-disuko/disuko/helper/validation"
	"github.com/eclipse-disuko/disuko/infra/rest"
	"github.com/eclipse-disuko/disuko/infra/service/rights"
	"github.com/eclipse-disuko/disuko/infra/service/startup"
	"github.com/eclipse-disuko/disuko/infra/service/test"
)

type handlers struct {
	project       rest.ProjectHandler
	mail          rest.MailHandler
	count         rest.CountHandler
	obligation    rest.ObligationsHandler
	checklist     rest.ChecklistHandler
	newsbox       rest.NewsboxHandler
	analyseFiles  rest.AnalyseFilesHandler
	notification  rest.NotificationHandler
	schema        rest.SchemaHandler
	policyRules   rest.PolicyRulesHandler
	licenses      rest.LicensesHandler
	announcements rest.AnnouncementsHandler
	analytics     rest.AnalyticsHandler
	export        rest.ExportHandler
	spdx          rest.SPDXHandler
	job           rest.JobHandler
	application   rest.ApplicationHandler
	auth          rest.OAuthHandler
	label         rest.LabelHandler
	startUp       startup.StartUpHandler
	statistic     rest.StatisticHandler
	sampleData    test.SampleDataHandler
	user          rest.UserHandler
	accessRights  rights.AccessRightsHandler
	department    rest.DepartmentHandler
	filterSet     rest.FilterSetHandler
	template      rest.TemplateHandler
	cap           rest.CapabilitiesHandler
	basicauth     rest.InternalTokenHandler
	customid      rest.CustomidHandler
	publicAuth    rest.PublicAuthHandler
}

func (s *Server) setupHandlers() {
	s.handlers.obligation = rest.ObligationsHandler{
		ObligationRepository:   s.repos.obligation,
		LicenseRepository:      s.repos.licenses,
		AuditLogListRepository: s.repos.auditLogList,
	}
	s.handlers.analyseFiles = rest.AnalyseFilesHandler{
		ProjectRepository: s.repos.project,
		DpConfigRepo:      s.repos.dpConfig,
		SbomListRepo:      s.repos.sbomList,
		ApprovalListRepo:  s.repos.approvalList,
	}
	s.handlers.notification = rest.NotificationHandler{
		DpConfigRepo: s.repos.dpConfig,
	}
	s.handlers.count = rest.CountHandler{
		ObligationRepository:     s.repos.obligation,
		ProjectRepository:        s.repos.project,
		LicenseRepository:        s.repos.licenses,
		PolicyRulesRepository:    s.repos.policyRules,
		LabelRepository:          s.repos.label,
		SchemaRepository:         s.repos.schema,
		UserRepository:           s.repos.user,
		ReviewTemplateRepository: s.repos.reviewTemplate,
		UserStatsRepository:      s.repos.userstats,
		Scheduler:                nil, // is set in setupScheduling
		NewsboxRepository:        s.repos.newsbox,
	}
	s.handlers.mail = rest.MailHandler{}
	s.handlers.project = rest.ProjectHandler{
		ObligationRepository:          s.repos.obligation,
		LicenseRepository:             s.repos.licenses,
		ProjectRepository:             s.repos.project,
		AnalyticsRepository:           s.repos.analytics,
		AnalyticsComponentsRepository: s.repos.analyticsComponents,
		AnalyticsLicensesRepository:   s.repos.analyticsLicenses,
		SchemaRepository:              s.repos.schema,
		PolicyRuleRepository:          s.repos.policyRules,
		UserRepository:                s.repos.user,
		LabelRepository:               s.repos.label,
		SbomListRepository:            s.repos.sbomList,
		AuditLogListRepository:        s.repos.auditLogList,
		ApprovalListRepository:        s.repos.approvalList,
		ReviewRemarksRepository:       s.repos.reviewRemarks,
		DpConfigRepo:                  s.repos.dpConfig,
		LockService:                   s.services.lock,
		DeparmentRepository:           s.repos.department,
		ApplicationConnector:          s.connectors.application,
		ReviewTemplateRepository:      s.repos.reviewTemplate,
		LicenseRulesRepository:        s.repos.licenseRules,
		JobRepository:                 s.repos.job,
		SpdxService:                   s.services.spdx,
		AnalyticsService:              &s.services.analytics,
		Scheduler:                     nil, // is set in setupScheduling
		CustomIdRepo:                  s.repos.customid,
		SbomRetainedService:           s.services.sbomRetained,
		ChecklistService:              &s.services.checklist,
		WizardService:                 &s.services.wizard,
		OverallReviewService:          &s.services.overallReview,
		ProjectLabelService:           &s.services.projectLabelService,
		FOSSddService:                 &s.services.fossdd,
		PolicyDecisionsRepository:     s.repos.policyDecisions,
		UserService:                   s.services.userService,
	}

	s.handlers.schema = rest.SchemaHandler{SchemaRepository: s.repos.schema, LabelRepository: s.repos.label}
	s.handlers.policyRules = rest.PolicyRulesHandler{
		LicenseRepository:       s.repos.licenses,
		PolicyRulesRepository:   s.repos.policyRules,
		ProjectRepository:       s.repos.project,
		LabelRepository:         s.repos.label,
		PolicyRulesService:      s.services.policyRules,
		SbomListRepository:      s.repos.sbomList,
		ChangeLogListRepository: s.repos.changeLogList,
	}
	s.handlers.licenses = rest.LicensesHandler{
		PolicyRulesRepository: s.repos.policyRules,
		ObligationRepository:  s.repos.obligation,
		LicenseRepository:     s.repos.licenses,
		JobRepository:         s.repos.job,
		SpdxLicenseRepository: s.repos.spdxLicense,
	}

	s.handlers.announcements = rest.AnnouncementsHandler{
		AnnouncementsRepository: s.repos.announcements,
	}

	s.handlers.analytics = rest.AnalyticsHandler{
		ProjectRepository:    s.repos.project,
		LicenseRepository:    s.repos.licenses,
		PolicyRuleRepository: s.repos.policyRules,
		AnalyticsRepository:  s.repos.analytics,
		SbomListRepository:   s.repos.sbomList,
		AnalyticsService:     s.services.analytics,
		StatisticRepository:  s.repos.statistic,
	}
	s.handlers.export = rest.ExportHandler{
		ExportService:    s.services.export,
		AnalyticsService: &s.services.analytics,
	}
	s.handlers.spdx = rest.SPDXHandler{
		ProjectRepository:         s.repos.project,
		SchemaRepository:          s.repos.schema,
		LicensesRepository:        s.repos.licenses,
		PolicyRuleRepository:      s.repos.policyRules,
		SbomListRepository:        s.repos.sbomList,
		AnalyticsService:          s.services.analytics,
		LabelRepository:           s.repos.label,
		LockService:               s.services.lock,
		AuditLogListRepository:    s.repos.auditLogList,
		LicenseRulesRepository:    s.repos.licenseRules,
		SpdxService:               s.services.spdx,
		SbomRetainedService:       s.services.sbomRetained,
		ProjectLabelService:       &s.services.projectLabelService,
		PolicyDecisionsRepository: s.repos.policyDecisions,
	}
	s.handlers.job = rest.JobHandler{JobRepository: s.repos.job}
	s.handlers.application = rest.ApplicationHandler{
		Connector: s.connectors.application,
	}

	s.handlers.label = rest.LabelHandler{
		LabelRepository:   s.repos.label,
		SchemaRepository:  s.repos.schema,
		ProjectRepository: s.repos.project,
		PolicyRepository:  s.repos.policyRules,
	}

	s.handlers.startUp = startup.StartUpHandler{
		PolicyRulesRepository:         s.repos.policyRules,
		AuditLogListRepository:        s.repos.auditLogList,
		DpConfigRepo:                  s.repos.dpConfig,
		MigrationRepository:           s.repos.migration,
		ReviewRemarkRepository:        s.repos.reviewRemarks,
		AnalyticsComponentsRepository: s.repos.analyticsComponents,
		ApprovalRepository:            s.repos.approvalList,
		ProjectRepository:             s.repos.project,
		ApplicationConnector:          s.connectors.application,
		SbomListRepository:            s.repos.sbomList,
		LabelRepository:               s.repos.label,
		JobRepository:                 s.repos.job,
		SbomRetainedService:           s.services.sbomRetained,
		ProjectHandler:                &s.handlers.project,
		LicenseRepository:             s.repos.licenses,
		LicenseRulesRepo:              s.repos.licenseRules,
		PolicyDecisionsRepo:           s.repos.policyDecisions,
	}
	s.handlers.sampleData = test.SampleDataHandler{
		PolicyRulesRepository: s.repos.policyRules,
		DpConfigRepo:          s.repos.dpConfig,
		ProjectRepository:     s.repos.project,
		LicensesRepository:    s.repos.licenses,
		ObligationRepository:  s.repos.obligation,
		SchemaRepository:      s.repos.schema,
		LabelRepository:       s.repos.label,
		SbomListRepository:    s.repos.sbomList,
		SpdxService:           s.services.spdx,
	}
	s.handlers.statistic = rest.StatisticHandler{
		PolicyRulesRepository: s.repos.policyRules,
		StatisticRepository:   s.repos.statistic,
		ProjectRepository:     s.repos.project,
		LicensesRepository:    s.repos.licenses,
		ObligationRepository:  s.repos.obligation,
		SchemaRepository:      s.repos.schema,
		LabelRepository:       s.repos.label,
		UsersRepository:       s.repos.user,
		DpConfigRepo:          s.repos.dpConfig,
	}
	s.handlers.user = rest.UserHandler{
		UserRepository:         s.repos.user,
		JobRepository:          s.repos.job,
		ApprovalListRepository: s.repos.approvalList,
		ProjectRepository:      s.repos.project,
		LabelRepository:        s.repos.label,
		UserroleConnector:      s.connectors.userrole,
		NewsBoxRepository:      s.repos.newsbox,
		DeletionService:        s.services.deletionService,
		DeletionAuditRepo:      s.repos.deletionAudit,
		UserService:            s.services.userService,
		AuditLogListRepo:       s.repos.auditLogList,
		SBOMListRepo:           s.repos.sbomList,
	}
	s.handlers.auth.UserRepository = s.repos.user
	s.handlers.accessRights = rights.AccessRightsHandler{}
	s.handlers.department = rest.DepartmentHandler{
		Repo: s.repos.department,
	}
	s.handlers.filterSet = rest.FilterSetHandler{
		FilterSetsRepository: s.repos.filterSets,
	}
	validation.UserRepository = s.repos.user
	s.handlers.template = rest.TemplateHandler{
		ReviewTemplateRepository: s.repos.reviewTemplate,
		ChecklistRepository:      s.repos.checklist,
	}
	s.handlers.cap = rest.CapabilitiesHandler{
		ApplicationConnector: s.connectors.application,
	}
	s.handlers.basicauth = rest.InternalTokenHandler{
		InternalTokenRepo: s.repos.basicauth,
	}
	s.handlers.customid = rest.CustomidHandler{
		Repo:        s.repos.customid,
		ProjectRepo: s.repos.project,
	}
	s.handlers.checklist = rest.ChecklistHandler{
		ChecklistRepo:    s.repos.checklist,
		ChecklistService: s.services.checklist,
	}
	s.handlers.newsbox = rest.NewsboxHandler{
		NewsboxRepo: s.repos.newsbox,
	}
	s.handlers.publicAuth = rest.PublicAuthHandler{
		ProjectRepo: s.repos.project,
	}
}
