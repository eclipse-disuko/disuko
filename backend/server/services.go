// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"github.com/eclipse-disuko/disuko/infra/service/analytics"
	sbomRetained "github.com/eclipse-disuko/disuko/infra/service/check-sbom-retained"
	"github.com/eclipse-disuko/disuko/infra/service/checklist"
	"github.com/eclipse-disuko/disuko/infra/service/export"
	"github.com/eclipse-disuko/disuko/infra/service/fossdd"
	"github.com/eclipse-disuko/disuko/infra/service/locks"
	"github.com/eclipse-disuko/disuko/infra/service/policy"
	"github.com/eclipse-disuko/disuko/infra/service/project"
	projectLabelService "github.com/eclipse-disuko/disuko/infra/service/project-label"
	"github.com/eclipse-disuko/disuko/infra/service/scanremarks"
	"github.com/eclipse-disuko/disuko/infra/service/spdx"
	userService "github.com/eclipse-disuko/disuko/infra/service/user"
	"github.com/eclipse-disuko/disuko/logy"
)

type services struct {
	lock                *locks.Service
	spdx                *spdx.Service
	analytics           analytics.Analytics
	policyRules         policy.Service
	export              *export.Service
	checklist           checklist.Service
	sbomRetained        *sbomRetained.Service
	scanRemarks         scanremarks.Service
	wizard              project.WizardService
	projectLabelService projectLabelService.ProjectLabelService
	fossdd              fossdd.Service
	overallReview       project.OverallReviewService
	deletionService     *userService.DeletionService
	userService         *userService.Service
}

func (s *Server) setupServices(rs *logy.RequestSession) {
	lockS := locks.InitService(rs)
	spdxS := &spdx.Service{
		LicenseRepo:      s.repos.licenses,
		LicenseRulesRepo: s.repos.licenseRules,
		LockService:      lockS,
	}
	psS := policy.Service{
		PolicyRulesRepository: s.repos.policyRules,
		LicenseRepository:     s.repos.licenses,
	}
	sbomRetainedS := sbomRetained.NewService(
		s.repos.project,
		s.repos.sbomList,
	)
	srS := scanremarks.Service{
		ProjectRepo: s.repos.project,
		LabelsRepo:  s.repos.label,
	}
	plS := projectLabelService.ProjectLabelService{
		ProjectRepo: s.repos.project,
		LabelRepo:   s.repos.label,
	}
	fossDDS := fossdd.Service{
		ProjectRepo:         s.repos.project,
		LabelRepo:           s.repos.label,
		SbomListRepo:        s.repos.sbomList,
		SpdxService:         spdxS,
		PolicyRuleRepo:      s.repos.policyRules,
		LicenseRepo:         s.repos.licenses,
		ReviewRemarksRepo:   s.repos.reviewRemarks,
		ProjectLabelService: &plS,
		PolicyDecisionsRepo: s.repos.policyDecisions,
	}
	fossDDS.ReadTemplates([]string{"vanilla"})

	userServ := userService.Init(rs, s.repos.user, s.repos.approvalList, s.repos.project, s.repos.label)

	s.services = services{
		lock:         lockS,
		spdx:         spdxS,
		policyRules:  psS,
		sbomRetained: sbomRetainedS,
		scanRemarks:  srS,
		analytics: analytics.Analytics{
			ProjectRepository:    s.repos.project,
			LicenseRepository:    s.repos.licenses,
			PolicyRuleRepository: s.repos.policyRules,
			SbomListrepository:   s.repos.sbomList,
			DepartmentRepo:       s.repos.department,
			Handler: analytics.InitDbHandler(
				s.repos.analytics,
				s.repos.analyticsComponents,
				s.repos.analyticsLicenses,
				s.repos.analyticsOccurrences,
				lockS,
			),
			LicenseRulesRepository: s.repos.licenseRules,
			SpdxService:            spdxS,
			LabelRepository:        s.repos.label,
			ProjectLabelService:    &plS,
			PolicyDecisionsRepo:    s.repos.policyDecisions,
		},
		export: &export.Service{
			LicensesRepository:    s.repos.licenses,
			PolicyRulesRepository: s.repos.policyRules,
			ObligationRepository:  s.repos.obligation,
			LabelRepository:       s.repos.label,
			SchemaRepository:      s.repos.schema,
		},
		checklist: checklist.Service{
			ChecklistRepo:       s.repos.checklist,
			TemplateRepo:        s.repos.reviewTemplate,
			SbomListRepo:        s.repos.sbomList,
			PolicyRuleRepo:      s.repos.policyRules,
			LicenseRepo:         s.repos.licenses,
			ReviewRemarkRepo:    s.repos.reviewRemarks,
			SpdxService:         spdxS,
			ScanRemarksService:  &srS,
			ProjectLabelService: &plS,
			PolicyDecisionsRepo: s.repos.policyDecisions,
			ProjectRepo:         s.repos.project,
		},
		wizard: project.WizardService{
			LabelRepository:        s.repos.label,
			ProjectRepository:      s.repos.project,
			DepartmentRepository:   s.repos.department,
			AuditLogListRepository: s.repos.auditLogList,
			ApplicationConnector:   s.connectors.application,
		},
		fossdd: fossDDS,
		overallReview: project.OverallReviewService{
			AuditlogRepo: s.repos.auditLogList,
			ProjectRepo:  s.repos.project,
			SbomListRepo: s.repos.sbomList,
			UserRepo:     s.repos.user,
		},
		projectLabelService: plS,
		deletionService: &userService.DeletionService{
			UserRepository:         s.repos.user,
			ProjectRepository:      s.repos.project,
			ApprovalListRepository: s.repos.approvalList,
			SbomListRepository:     s.repos.sbomList,
			AuditLogListRepository: s.repos.auditLogList,
		},
		userService: userServ,
	}
}
