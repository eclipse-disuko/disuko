// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package startup

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/audit"
	"github.com/eclipse-disuko/disuko/domain/job"
	"github.com/eclipse-disuko/disuko/domain/label"
	"github.com/eclipse-disuko/disuko/domain/license"
	auditHelper "github.com/eclipse-disuko/disuko/helper/audit"
	"github.com/eclipse-disuko/disuko/infra/repository/auditloglist"
	"github.com/eclipse-disuko/disuko/infra/repository/jobs"
	"github.com/eclipse-disuko/disuko/infra/repository/labels"
	licenseRepo "github.com/eclipse-disuko/disuko/infra/repository/license"
	"github.com/eclipse-disuko/disuko/infra/repository/licenserules"
	"github.com/eclipse-disuko/disuko/infra/repository/policydecisions"
	"github.com/eclipse-disuko/disuko/infra/rest"
	sbomLockRetained "github.com/eclipse-disuko/disuko/infra/service/check-sbom-retained"

	"github.com/eclipse-disuko/disuko/connector/application"
	"github.com/eclipse-disuko/disuko/domain/approval"
	migrationDomain "github.com/eclipse-disuko/disuko/domain/migration"
	"github.com/eclipse-disuko/disuko/domain/project"
	reviewremarks2 "github.com/eclipse-disuko/disuko/domain/reviewremarks"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/infra/repository/analyticscomponents"
	"github.com/eclipse-disuko/disuko/infra/repository/approvallist"
	"github.com/eclipse-disuko/disuko/infra/repository/database"
	"github.com/eclipse-disuko/disuko/infra/repository/dpconfig"
	"github.com/eclipse-disuko/disuko/infra/repository/migration"
	"github.com/eclipse-disuko/disuko/infra/repository/policyrules"
	projectRepo "github.com/eclipse-disuko/disuko/infra/repository/project"
	"github.com/eclipse-disuko/disuko/infra/repository/reviewremarks"
	sbomlistRepo "github.com/eclipse-disuko/disuko/infra/repository/sbomlist"
	"github.com/eclipse-disuko/disuko/logy"
)

type StartUpHandler struct {
	DpConfigRepo                  *dpconfig.DBConfigRepository
	AuditLogListRepository        auditloglist.IAuditLogListRepository
	PolicyRulesRepository         policyrules.IPolicyRulesRepository
	MigrationRepository           migration.IMigrationRepository
	ReviewRemarkRepository        reviewremarks.IReviewRemarksRepository
	AnalyticsComponentsRepository analyticscomponents.IComponentsRepository
	ApprovalRepository            approvallist.IApprovalListRepository
	ProjectRepository             projectRepo.IProjectRepository
	ApplicationConnector          *application.Connector
	SbomListRepository            sbomlistRepo.ISbomListRepository
	LabelRepository               labels.ILabelRepository
	JobRepository                 jobs.IJobsRepository
	ProjectHandler                *rest.ProjectHandler
	SbomRetainedService           *sbomLockRetained.Service
	LicenseRepository             licenseRepo.ILicensesRepository
	LicenseRulesRepo              licenserules.ILicenseRulesRepository
	PolicyDecisionsRepo           policydecisions.IPolicyDecisionsRepository
}

type Step struct {
	Name string
	Do   func(requestSession *logy.RequestSession)
}

func (startUpHandler *StartUpHandler) migrateSBOMUploadDates(requestSession *logy.RequestSession) {
	logy.Infof(requestSession, "migrateSBOMUploadDates - START")

	reviewRemarks := startUpHandler.ReviewRemarkRepository.FindAll(requestSession, false)
	for _, rr := range reviewRemarks {
		changed := false
		sbomList := startUpHandler.SbomListRepository.FindByKey(requestSession, rr.Key, false)
		if sbomList == nil {
			continue
		}

		// Create a map of SBOM ID to upload date for quick lookup
		uploadDates := make(map[string]*time.Time)
		for _, sbom := range sbomList.SpdxFileHistory {
			uploadDates[sbom.Key] = sbom.Uploaded
		}

		for _, r := range rr.Remarks {
			if r.SBOMId == "" || r.SBOMUploaded != nil {
				continue
			}

			// Look up the upload date for this SBOM
			if uploadDate, ok := uploadDates[r.SBOMId]; ok {
				r.SBOMUploaded = uploadDate
				changed = true
				logy.Infof(requestSession, "migrateSBOMUploadDates - Updated review remark SBOM upload date: Version: '%s', SBOM: '%s'", rr.Key, r.SBOMId)
			}
		}

		if changed {
			startUpHandler.ReviewRemarkRepository.UpdateWithoutTimestamp(requestSession, rr)
		}
	}
	logy.Infof(requestSession, "migrateSBOMUploadDates - END")
}

func (startUpHandler *StartUpHandler) MigrateDatabase(requestSession *logy.RequestSession, ext ...Step) {
	steps := []Step{
		{Name: "ADD_DELETE_TO_RRS", Do: startUpHandler.addDeleteToRRs},
		{Name: "MIGRATE_CHANGE_GROUP_STATUS", Do: startUpHandler.migrateChangeGroupStatus},
		{Name: "MIGRATE_ACTIVATE_PRS", Do: startUpHandler.migrateActivatePRs},
		{Name: "MIGRAGE_APP_SEC_ID_2", Do: startUpHandler.migrateAppSecId},
		{Name: "MIGRATE_REVIEW_REMARK", Do: startUpHandler.migrateReviewRemarks},
		{Name: "MIGRATE_SBOM_UPLOAD_DATES", Do: startUpHandler.migrateSBOMUploadDates},
		{Name: "MIGRATE_ADD_PROJECT_LABEL_DUMMY", Do: startUpHandler.migrateAddProjectLabelDummy},
		{Name: "MIGRATE_NEW_JOB_DELETE_DUMMY_PROJECTS", Do: startUpHandler.migrateDeleteDummyProjects},
		{Name: "MIGRATE_NEW_JOB_DUMMY_MAIL", Do: startUpHandler.migrateDummyMail},
		{Name: "MIGRATE_PROJECT_FLAGS", Do: startUpHandler.migrateProjectFlags},
		{Name: "MIGRATE_PROJECT_LABELS", Do: startUpHandler.migrateLabels},
		{Name: "MIGRATE_EVAL_LICS", Do: startUpHandler.migrateEvalLics},
		{Name: "MIGRATE_PARENT", Do: startUpHandler.migrateParent},
		{Name: "MIGRATE_OVERALL_REVIEWS_SBOM", Do: startUpHandler.migrateOverallReviews},
		{Name: "MIGRATE_SBOM_UPLOADED_FOR_DECISIONS", Do: startUpHandler.migrateSbomUploadedForDecisions},
		{Name: "MIGRATE_SBOM_RETENTION_FOR_DECISIONS", Do: startUpHandler.migrateSbomRetentionForDecisions},
	}

	steps = append(steps, ext...)

	for _, step := range steps {
		qc := database.New().SetMatcher(
			database.AndChain(
				database.AttributeMatcher(
					"Name",
					database.EQ,
					step.Name,
				),
				/* database.AttributeMatcher( // deleted attribute is not availabe in migrations collections
					"Deleted",
					database.EQ,
					false,
				),*/
			),
		).SetSort(database.SortConfig{
			database.SortAttribute{
				Name:  "Created",
				Order: database.DESC,
			},
		})

		migs := startUpHandler.MigrationRepository.Query(requestSession, qc)

		if len(migs) > 0 {
			logy.Infof(requestSession, "skipping migration: "+step.Name)
			continue
		}
		logy.Infof(requestSession, "starting migration: "+step.Name)
		step.Do(requestSession)
		startUpHandler.MigrationRepository.Save(requestSession, migrationDomain.New(requestSession, step.Name))
		logy.Infof(requestSession, "migration step "+step.Name+" finished")
	}
}

func (startUpHandler *StartUpHandler) addDeleteToRRs(requestSession *logy.RequestSession) {
	remarks := startUpHandler.ReviewRemarkRepository.FindAllWithDeleted(requestSession, false)
	for _, rr := range remarks {
		startUpHandler.ReviewRemarkRepository.UpdateWithoutTimestamp(requestSession, rr)
	}
}

func (startUpHandler *StartUpHandler) migrateChangeGroupStatus(requestSession *logy.RequestSession) {
	logy.Infof(requestSession, "Group Status - migrateChangeGroupStatus - START")

	allprojects := startUpHandler.ProjectRepository.FindAllWithDeleted(requestSession, false)
	for _, grp := range allprojects {
		groupChanged := false
		if !grp.IsGroup {
			continue
		}

		approvalList := startUpHandler.ApprovalRepository.FindByKey(requestSession, grp.Key, false)
		if approvalList == nil {
			logy.Infof(requestSession, "Group Status - No approvals available for the Group: %s (%s)", grp.Key, grp.Name)
			grp.Status = project.Ready
			groupChanged = true
		} else {
			for _, appr := range approvalList.Approvals {
				if appr.Type == approval.TypeInternal || appr.Type == approval.TypeExternal {
					logy.Infof(requestSession, "Group Status - Active approvals available for Group: %s (%s)", grp.Key, grp.Name)
					grp.Status = project.Active
					groupChanged = true
					break
				}
			}
		}

		if groupChanged {
			logy.Infof(requestSession, "Group Status - Changing status to '%s' for Group: %s (%s)", grp.Status, grp.Key, grp.Name)
			startUpHandler.ProjectRepository.UpdateWithoutTimestamp(requestSession, grp)
		}
	}
	logy.Infof(requestSession, "Group Status - migrateChangeGroupStatus - END")
}

func (startUpHandler *StartUpHandler) migrateActivatePRs(requestSession *logy.RequestSession) {
	logy.Infof(requestSession, "migrateActivatePRs - START")

	prs := startUpHandler.PolicyRulesRepository.FindAll(requestSession, false)
	for _, pr := range prs {
		pr.Active = true
		pr.Deprecated = false
		logy.Infof(requestSession, "activating '%s' (%s)", pr.Name, pr.Key)
		startUpHandler.PolicyRulesRepository.UpdateWithoutTimestamp(requestSession, pr)
	}
	logy.Infof(requestSession, "migrateActivatePRs - END")
}

func (startUpHandler *StartUpHandler) migrateAppSecId(requestSession *logy.RequestSession) {
	if startUpHandler.ApplicationConnector == nil {
		logy.Infof(requestSession, "migrateAppSecId - no connector")
		return
	}
	logy.Infof(requestSession, "migrateAppSecId - START")

	allprojects := startUpHandler.ProjectRepository.FindAll(requestSession, false)
	for _, p := range allprojects {
		if p.IsGroup || p.ApplicationMeta.Id == "" || p.ApplicationMeta.SecondaryId != "" {
			continue
		}

		var (
			app               application.Application
			exceptionHappened bool
		)
		exception.TryCatch(func() {
			app = startUpHandler.ApplicationConnector.GetApplication(requestSession, p.ApplicationMeta.Id)
		}, func(e exception.Exception) {
			exceptionHappened = true
			exception.LogException(requestSession, e)
			logy.Infof(requestSession, "migrateAppSecId - connector request error for %s", p.Key)
		})
		if exceptionHappened {
			continue
		}
		p.ApplicationMeta = project.ApplicationMeta{
			Id:           app.Id,
			SecondaryId:  app.SecondaryId,
			Name:         app.Name,
			ExternalLink: app.Link,
		}
		startUpHandler.ProjectRepository.UpdateWithoutTimestamp(requestSession, p)
		logy.Infof(requestSession, "migrateAppSecId - updated project %s %s", p.Key, p.Name)
	}
	logy.Infof(requestSession, "migrateAppSecId - END")
}

func (startUpHandler *StartUpHandler) migrateReviewRemarks(requestSession *logy.RequestSession) {
	logy.Infof(requestSession, "migrateReviewRemarks - START")

	reviewRemarks := startUpHandler.ReviewRemarkRepository.FindAllWithDeleted(requestSession, false)
	for _, rr := range reviewRemarks {
		changed := false
		for _, r := range rr.Remarks {
			if r.ComponentId == "" {
				continue
			}

			r.Components = append(r.Components, reviewremarks2.ComponentMeta{
				ComponentId:      r.ComponentId,
				ComponentName:    r.ComponentName,
				ComponentVersion: r.ComponentVersion,
			})
			r.ComponentId = ""
			r.ComponentName = ""
			r.ComponentVersion = ""
			changed = true
			logy.Infof(requestSession, "migrateReviewRemarks - migrating review remark: UUID: '%s', remark title: '%s', component: %s@%s (%s)", rr.Key, r.Title, r.ComponentName, r.ComponentVersion, r.ComponentId)

			if r.LicenseId == "" {
				continue
			}

			r.Licenses = append(r.Licenses, reviewremarks2.LicenseMeta{
				LicenseId:   r.LicenseId,
				LicenseName: r.LicenseName,
			})
			r.LicenseId = ""
			r.LicenseName = ""
			logy.Infof(requestSession, "migrateReviewRemarks - migrating review remark: UUID: '%s', remark title: '%s', license: %s (%s)", rr.Key, r.Title, r.LicenseName, r.LicenseId)
		}

		if changed {
			startUpHandler.ReviewRemarkRepository.UpdateWithoutTimestamp(requestSession, rr)
		}
	}

	logy.Infof(requestSession, "migrateReviewRemarks - END")
}

func (startUpHandler *StartUpHandler) migrateAddProjectLabelDummy(requestSession *logy.RequestSession) {
	logy.Infof(requestSession, "migrateAddProjectLabelDummy - START")
	existingLabel := startUpHandler.LabelRepository.FindByNameAndType(requestSession, label.DUMMY, label.PROJECT)
	if existingLabel != nil {
		logy.Infof(requestSession, "migrateAddProjectLabelDummy - label exists, skipping migration")
	} else {
		dummyLabel := &label.Label{
			RootEntity:  domain.NewRootEntity(),
			Name:        label.DUMMY,
			Description: "Dummy project label",
			Type:        label.PROJECT,
		}
		startUpHandler.LabelRepository.Save(requestSession, dummyLabel)
		logy.Infof(requestSession, "migrateAddProjectLabelDummy - label '%s' created", dummyLabel.Name)
	}
	logy.Infof(requestSession, "migrateAddProjectLabelDummy - END")
}

func (startUpHandler *StartUpHandler) migrateDeleteDummyProjects(requestSession *logy.RequestSession) {
	logy.Infof(requestSession, "migrateDeleteDummyProjects - START")
	startUpHandler.JobRepository.Save(requestSession, &job.Job{
		RootEntity: domain.NewRootEntity(),
		Name:       "Delete Dummy Projects",
		JobType:    job.DummyProjectDeletion,
		Execution:  job.Periodic,
		Schedule:   "30 4 * * *", // every day 04:30,
	})
	logy.Infof(requestSession, "migrateDeleteDummyProjects - END")
}

func (startUpHandler *StartUpHandler) migrateDummyMail(requestSession *logy.RequestSession) {
	logy.Infof(requestSession, "migrateDummyMail - START")
	startUpHandler.JobRepository.Save(requestSession, &job.Job{
		RootEntity: domain.NewRootEntity(),
		Name:       "Dummy deletion mail",
		JobType:    job.DummyMail,
		Execution:  job.Periodic,
		Schedule:   "00 4 * * *", // every day 04:00,
	})
	logy.Infof(requestSession, "migrateDummyMail - END")
}

func (startUpHandler *StartUpHandler) migrateProjectFlags(requestSession *logy.RequestSession) {
	logy.Infof(requestSession, "migrateProjectFlags - START")

	// Move the actual migration logic here instead of creating a job
	exception.TryCatchAndLog(requestSession, func() {
		projects := startUpHandler.ProjectRepository.FindAllKeys(requestSession)
		logy.Infof(requestSession, "migrateProjectFlags - Found %d projects to process", len(projects))

		processedCount := 0
		failedCount := 0
		successCount := 0

		for _, projectID := range projects {
			exception.TryCatch(func() {
				project := startUpHandler.ProjectRepository.FindByKey(requestSession, projectID, false)
				if project == nil {
					logy.Warnf(requestSession, "migrateProjectFlags - Project not found: %s", projectID)
					failedCount++
					return
				}

				processedCount++

				project.HasChildren = startUpHandler.ProjectHandler.CountChildren(requestSession, project, project.Children) > 0
				project.HasApproval = startUpHandler.ProjectHandler.IsReferencedInApprovalLists(requestSession, project, nil)
				project.HasSBOMToRetain = startUpHandler.SbomRetainedService.HasAnyVersionWithRetainedSbom(requestSession, project)

				startUpHandler.ProjectRepository.UpdateWithoutTimestamp(requestSession, project)
				successCount++

				logy.Infof(requestSession, "migrateProjectFlags - Updated project flags for: %s", project.Key)
			}, func(e exception.Exception) {
				exception.LogException(requestSession, e)
				failedCount++
			})
		}

		logy.Infof(requestSession, "migrateProjectFlags - Project flags migration completed. Processed: %d, Success: %d, Failed: %d",
			processedCount, successCount, failedCount)
	})

	logy.Infof(requestSession, "migrateProjectFlags - END")
}

func (startUpHandler *StartUpHandler) migrateEvalLics(requestSession *logy.RequestSession) {
	logy.Infof(requestSession, "migrateEvalLics - START")

	exception.TryCatchAndLog(requestSession, func() {
		lics := startUpHandler.LicenseRepository.FindAllKeys(requestSession)
		logy.Infof(requestSession, "migrateEvalLics - Found %d licenses to process", len(lics))

		failedCount := 0
		changedCount := 0

		for _, licID := range lics {
			exception.TryCatch(func() {
				lic := startUpHandler.LicenseRepository.FindByKey(requestSession, licID, false)
				if lic == nil {
					logy.Warnf(requestSession, "migrateEvalLics - License not found: %s", lic)
					failedCount++
					return
				}

				if lic.Meta.ApprovalState != license.Pending {
					return
				}

				before := lic.ToAudit(requestSession, nil)
				lic.Meta.ApprovalState = license.NotSet
				after := lic.ToAudit(requestSession, nil)
				auditHelper.CreateAndAddAuditEntry(&lic.Container, "SYSTEM", message.LicenseUpdated, audit.DiffWithReporter, after, before)
				startUpHandler.LicenseRepository.Update(requestSession, lic)
				changedCount++

				logy.Infof(requestSession, "migrateEvalLics - Updated license: %s", lic.Key)
			}, func(e exception.Exception) {
				exception.LogException(requestSession, e)
				failedCount++
			})
		}

		logy.Infof(requestSession, "migrateEvalLics - License eval state migration. Processed: %d, Changed: %d, Failed: %d",
			len(lics), changedCount, failedCount)
	})

	logy.Infof(requestSession, "migrateEvalLics - END")
}

func (startUpHandler *StartUpHandler) migrateParent(requestSession *logy.RequestSession) {
	logy.Infof(requestSession, "migrateParent - START")

	exception.TryCatchAndLog(requestSession, func() {
		projects := startUpHandler.ProjectRepository.FindAllKeys(requestSession)
		logy.Infof(requestSession, "migrateParent - Found %d projects to process", len(projects))

		for _, projectID := range projects {
			exception.TryCatch(func() {
				project := startUpHandler.ProjectRepository.FindByKey(requestSession, projectID, false)
				if project == nil {
					logy.Warnf(requestSession, "migrateParent - Project not found: %s", projectID)
					return
				}
				if project.Parent == "" {
					return
				}
				parentProject := startUpHandler.ProjectRepository.FindByKey(requestSession, project.Parent, true)
				if parentProject != nil {
					project.ParentName = parentProject.Name
					startUpHandler.ProjectRepository.UpdateWithoutTimestamp(requestSession, project)
					logy.Infof(requestSession, "migrateParent - Updated parent name for project: %s (parent: %s)", project.Key, parentProject.Name)
				} else {
					logy.Warnf(requestSession, "migrateParent - Parent project not found for: %s (parent key: %s)", project.Key, project.Parent)
				}
			}, func(e exception.Exception) {
				exception.LogException(requestSession, e)
				logy.Errorf(requestSession, "migrateParent - Failed to update project: %s", projectID)
			})
		}

		logy.Infof(requestSession, "migrateParent - Parent name migration completed")
	})

	logy.Infof(requestSession, "migrateParent - END")
}

func (startUpHandler *StartUpHandler) migrateSbomUploadedForDecisions(requestSession *logy.RequestSession) {
	logy.Infof(requestSession, "migrateSbomUploadedForDecisions - START")

	sbomIdToUploadedMap := make(map[string]*time.Time)
	processedProjects := make(map[string]struct{})

	licenseRulesList := startUpHandler.LicenseRulesRepo.FindAll(requestSession, false)
	for _, licenseRules := range licenseRulesList {
		changed := false

		startUpHandler.processProject(requestSession, licenseRules.Key, processedProjects, sbomIdToUploadedMap)

		for i := range licenseRules.Rules {
			if uploaded, ok := sbomIdToUploadedMap[licenseRules.Rules[i].SBOMId]; ok {
				licenseRules.Rules[i].SBOMUploaded = uploaded
				changed = true
				logy.Infof(requestSession, "migrateSbomUploadedForDecisions - Updated LicenseRule SBOM uploaded for Project-UUID: '%s', SBOM-UUID: '%s'", licenseRules.Key, licenseRules.Rules[i].SBOMId)
			}
		}

		if changed {
			startUpHandler.LicenseRulesRepo.UpdateWithoutTimestamp(requestSession, licenseRules)
		}
	}

	policyDecisionsList := startUpHandler.PolicyDecisionsRepo.FindAll(requestSession, false)
	for _, policyDecisions := range policyDecisionsList {
		changed := false

		startUpHandler.processProject(requestSession, policyDecisions.Key, processedProjects, sbomIdToUploadedMap)

		for i := range policyDecisions.Decisions {
			if uploaded, ok := sbomIdToUploadedMap[policyDecisions.Decisions[i].SBOMId]; ok {
				policyDecisions.Decisions[i].SBOMUploaded = uploaded
				changed = true
				logy.Infof(requestSession, "migrateSbomUploadedForDecisions - Updated PolicyDecision SBOM uploaded for Project-UUID: '%s', SBOM-UUID: '%s'", policyDecisions.Key, policyDecisions.Decisions[i].SBOMId)
			}
		}

		if changed {
			startUpHandler.PolicyDecisionsRepo.UpdateWithoutTimestamp(requestSession, policyDecisions)
		}
	}

	logy.Infof(requestSession, "migrateSbomUploadedForDecisions - END")
}

func (startUpHandler *StartUpHandler) processProject(
	requestSession *logy.RequestSession,
	key string,
	processedProjects map[string]struct{},
	sbomIdToUploadedMap map[string]*time.Time,
) {
	if _, found := processedProjects[key]; !found {
		prj := startUpHandler.ProjectRepository.FindByKey(requestSession, key, false)
		if prj == nil {
			return
		}
		for _, v := range prj.Versions {
			sbomList := startUpHandler.SbomListRepository.FindByKey(requestSession, v.Key, false)
			if sbomList == nil {
				continue
			}

			for _, sbom := range sbomList.SpdxFileHistory {
				if sbom.Uploaded != nil {
					sbomIdToUploadedMap[sbom.Key] = sbom.Uploaded
				}
			}
		}
		processedProjects[prj.Key] = struct{}{}
	}
}

func (startUpHandler *StartUpHandler) migrateSbomRetentionForDecisions(rs *logy.RequestSession) {
	logy.Infof(rs, "migrateSbomRetentionForDecisions - START")

	refs := map[string]map[string]struct{}{}

	addRef := func(projectId, sbomId string) {
		if projectId == "" || sbomId == "" {
			return
		}
		m, ok := refs[projectId]
		if !ok {
			m = map[string]struct{}{}
			refs[projectId] = m
		}
		m[sbomId] = struct{}{}
	}

	licenseRulesList := startUpHandler.LicenseRulesRepo.FindAll(rs, false)
	for _, licenseRules := range licenseRulesList {
		for _, lr := range licenseRules.Rules {
			addRef(licenseRules.Key, lr.SBOMId)
		}
	}

	policyDecisionsList := startUpHandler.PolicyDecisionsRepo.FindAll(rs, false)
	for _, policyDecisions := range policyDecisionsList {
		for _, pd := range policyDecisions.Decisions {
			addRef(policyDecisions.Key, pd.SBOMId)
		}
	}

	logy.Infof(rs, "migrateSbomRetentionForDecisions - Collected references from %d projects", len(refs))

	totalMarked := 0
	totalNotFound := 0

	for projectId, wanted := range refs {
		prj := startUpHandler.ProjectRepository.FindByKey(rs, projectId, false)
		if prj == nil {
			logy.Warnf(rs, "migrateSbomRetentionForDecisions - Project not found: %s", projectId)
			continue
		}

		found := make(map[string]struct{}, len(wanted))

		for _, v := range prj.Versions {
			sbomList := startUpHandler.SbomListRepository.FindByKey(rs, v.Key, false)
			if sbomList == nil || len(sbomList.SpdxFileHistory) == 0 {
				continue
			}

			changed := false

			for i := range sbomList.SpdxFileHistory {
				sbomKey := sbomList.SpdxFileHistory[i].Key

				if _, ok := wanted[sbomKey]; !ok {
					continue
				}

				found[sbomKey] = struct{}{}

				if !sbomList.SpdxFileHistory[i].IsToRetain {
					sbomList.SpdxFileHistory[i].IsToRetain = true
					changed = true
					totalMarked++
					logy.Infof(rs, "migrateSbomRetentionForDecisions - Marked SBOM for retention - Project: '%s', Version: '%s', SBOM: '%s'", projectId, v.Key, sbomKey)
				}
			}

			if changed {
				startUpHandler.SbomListRepository.UpdateWithoutTimestamp(rs, sbomList)
			}
		}

		for sbomId := range wanted {
			if _, ok := found[sbomId]; !ok {
				totalNotFound++
				logy.Warnf(rs, "migrateSbomRetentionForDecisions - Referenced SBOM not found - Project: %s, SBOM: %s", projectId, sbomId)
			}
		}
	}

	logy.Infof(rs, "migrateSbomRetentionForDecisions - Marked %d SBOMs for retention, %d SBOMs not found", totalMarked, totalNotFound)
	logy.Infof(rs, "migrateSbomRetentionForDecisions - END")
}
