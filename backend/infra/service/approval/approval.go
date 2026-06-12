// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package approval

import (
	"fmt"
	"slices"
	"sort"
	"time"

	"github.com/eclipse-disuko/disuko/domain/approval"
	"github.com/eclipse-disuko/disuko/domain/audit"
	license2 "github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/project/approvable"
	"github.com/eclipse-disuko/disuko/domain/project/components"
	"github.com/eclipse-disuko/disuko/domain/project/sbomlist"
	user2 "github.com/eclipse-disuko/disuko/domain/user"
	auditHelper "github.com/eclipse-disuko/disuko/helper/audit"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/hash"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/infra/repository/approvallist"
	"github.com/eclipse-disuko/disuko/infra/repository/auditloglist"
	"github.com/eclipse-disuko/disuko/infra/repository/labels"
	"github.com/eclipse-disuko/disuko/infra/repository/license"
	"github.com/eclipse-disuko/disuko/infra/repository/licenserules"
	"github.com/eclipse-disuko/disuko/infra/repository/policydecisions"
	"github.com/eclipse-disuko/disuko/infra/repository/policyrules"
	projectRepo "github.com/eclipse-disuko/disuko/infra/repository/project"
	sbomListRepo "github.com/eclipse-disuko/disuko/infra/repository/sbomlist"
	"github.com/eclipse-disuko/disuko/infra/repository/user"
	"github.com/eclipse-disuko/disuko/infra/service/fossdd"
	projectService "github.com/eclipse-disuko/disuko/infra/service/project"
	projectLabelService "github.com/eclipse-disuko/disuko/infra/service/project-label"
	"github.com/eclipse-disuko/disuko/infra/service/spdx"
	"github.com/eclipse-disuko/disuko/logy"
)

type SpdxRetriever interface {
	RetrieveSbomListAndFile(*logy.RequestSession, string, string) (*sbomlist.SbomList, *project.SpdxFileBase)
}

type ApprovalService struct {
	RequestSession *logy.RequestSession

	UserRepo            user.IUsersRepository
	ProjectRepo         projectRepo.IProjectRepository
	LicenseRepo         license.ILicensesRepository
	PolicyRulesRepo     policyrules.IPolicyRulesRepository
	SBOMListRepo        sbomListRepo.ISbomListRepository
	ApprovalListRepo    approvallist.IApprovalListRepository
	AuditLogListRepo    auditloglist.IAuditLogListRepository
	LabelRepo           labels.ILabelRepository
	LicenseRulesRepo    licenserules.ILicenseRulesRepository
	PolicyDecisionsRepo policydecisions.IPolicyDecisionsRepository

	SpdxRetriever SpdxRetriever

	WizardService        *projectService.WizardService
	ProjectLabelService  *projectLabelService.ProjectLabelService
	FOSSddService        *fossdd.Service
	SpdxService          *spdx.Service
	OverallReviewService *projectService.OverallReviewService
}

func (s *ApprovalService) ProcessRandomApprovalUpdate(pr *project.Project, appId, username string, req approval.UpdateApprovalDto) *approval.Approval {
	approvalList := s.ApprovalListRepo.FindByKey(s.RequestSession, pr.Key, false)
	if approvalList == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorDbNotFound), "")
	}
	targetApproval := approvalList.GetApproval(appId)
	if targetApproval == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorDbNotFound), "")
	}

	switch targetApproval.Type {
	case approval.TypeInternal:
		s.processInternalApprovalUpdate(pr, targetApproval, username, req)
	case approval.TypePlausibility:
		s.processPlausibilityCheckUpdate(pr, targetApproval, username, req)
	case approval.TypeExternal:
		s.processExternalApprovalUpdate(pr, targetApproval, username, req)
	default:
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorUnexpectedType), "")
	}
	targetApproval.Updated = time.Now()
	s.ApprovalListRepo.Update(s.RequestSession, approvalList)
	return targetApproval
}

func (s *ApprovalService) GetApprovalInfo(targetProject *project.Project, takeLatestSbom bool) approval.Info {
	return s.getApprovalInfo(targetProject, nil, false, takeLatestSbom)
}

func (s *ApprovalService) AdminAbortRandomApproval(pr *project.Project, app *approval.Approval) {
	switch app.Type {
	case approval.TypeInternal:
		s.adminAbortInternal(pr, app)
	case approval.TypePlausibility:
		s.adminAbortPlausibility(pr, app)
	default:
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorUnexpectedType), "")
	}
}

func (s *ApprovalService) getApprovalInfo(targetProject *project.Project, projectFilter *[]string, includeNoFOSS bool, takeLatestSbom bool) approval.Info {
	res := approval.Info{
		CompStats: &components.ComponentStats{},
	}

	var projects []string
	if targetProject.IsGroup {
		if projectFilter != nil {
			for _, s := range *projectFilter {
				if !slices.Contains(targetProject.Children, s) {
					continue
				}
				projects = append(projects, s)
			}
		} else {
			projects = targetProject.Children
		}
	} else {
		projects = []string{targetProject.Key}
	}
	policyRulesAll := s.PolicyRulesRepo.FindAll(s.RequestSession, false)
	licenseRefs := s.LicenseRepo.GetLicenseRefs(s.RequestSession)
	licenseRefsHash := licenseRefs.GenHash(s.RequestSession)
	for _, prKey := range projects {
		pr := s.ProjectRepo.FindByKeyWithDeleted(s.RequestSession, prKey, false)
		if pr == nil {
			logy.Warnf(s.RequestSession, "Child project not found uuid: %s parent: %s", prKey, targetProject.Key)
			continue
		}
		if pr.Deleted {
			logy.Warnf(s.RequestSession, "Child project is marked as deleted, uuid: %s parent: %s", prKey, pr.Key)
			continue
		}
		if pr.IsDeprecated() {
			logy.Warnf(s.RequestSession, "Child project is marked as deprecated, uuid: %s parent: %s", prKey, pr.Key)
			continue
		}

		var supplier *project.ProjectMemberEntity
		for _, u := range pr.UserManagement.Users {
			if u.UserType == project.SUPPLIER {
				supplier = u
				break
			}
		}
		var supplierUserId *string
		if supplier != nil {
			supplierUserId = &supplier.UserId
		}

		approvableSPDX := pr.ApprovableSPDX
		var sbomList *sbomlist.SbomList
		var sbom *project.SpdxFileBase

		hasProjectApprovable := approvableSPDX.SpdxKey != "" && approvableSPDX.VersionKey != ""

		if takeLatestSbom && !hasProjectApprovable {
			approvableSPDX, sbomList, sbom = s.findLatestSpdx(pr)
		}

		if approvableSPDX.SpdxKey == "" || approvableSPDX.VersionKey == "" || (!takeLatestSbom && !includeNoFOSS && pr.IsNoFoss) {
			res.Projects = append(res.Projects, approval.ProjectApprovable{
				ProjectKey:      pr.Key,
				ProjectName:     pr.Name,
				CustomerDiffers: pr.CustomerMeta.Diff(targetProject.CustomerMeta),
				SupplierDiffers: pr.DocumentMeta.Diff(targetProject.DocumentMeta),
				Supplier:        supplierUserId,
			})
			continue
		}
		if sbom == nil || sbomList == nil {
			sbomList, sbom = s.SpdxRetriever.RetrieveSbomListAndFile(s.RequestSession, approvableSPDX.VersionKey, approvableSPDX.SpdxKey)
		}
		if sbom == nil || sbomList == nil {
			res.Projects = append(res.Projects, approval.ProjectApprovable{
				ProjectKey:      pr.Key,
				ProjectName:     pr.Name,
				CustomerDiffers: pr.CustomerMeta.Diff(targetProject.CustomerMeta),
				SupplierDiffers: pr.DocumentMeta.Diff(targetProject.DocumentMeta),
				Supplier:        supplierUserId,
			})
			continue
		}

		spdxFileHistory := sbomList.SpdxFileHistory
		sort.Slice(spdxFileHistory, func(i, j int) bool {
			return spdxFileHistory[i].Uploaded.UTC().After(spdxFileHistory[j].Uploaded.UTC())
		})
		var isSpdxRecent bool
		if sbom.Key == spdxFileHistory[0].Key {
			isSpdxRecent = true
		}

		rules := policyrules.FilterPolicyRulesForLabel(policyRulesAll, pr.PolicyLabels)
		prl := license2.PolicyRulesList(rules)
		licenseRules := s.LicenseRulesRepo.FindByKey(s.RequestSession, pr.Key, false)
		policyDecisions := s.PolicyDecisionsRepo.FindByKey(s.RequestSession, pr.Key, false)
		projectPolicyRulesHash := prl.GenHash(s.RequestSession)

		licenseRulesHash := licenseRules.GenHash(s.RequestSession)
		policyDecisionsHash := policyDecisions.GenHash(s.RequestSession)
		currentTotalStatsHash := new(hash.Hash(s.RequestSession, fmt.Sprintf(
			"%s|%s|%s|%s",
			projectPolicyRulesHash,
			licenseRefsHash,
			licenseRulesHash,
			policyDecisionsHash,
		)))

		var sbomStats components.ComponentStats
		if sbom.TotalStatsHash != nil && *sbom.TotalStatsHash == *currentTotalStatsHash {
			sbomStats = sbom.Stats
		} else {
			compsInfo := s.SpdxService.GetComponentInfos(s.RequestSession, pr, approvableSPDX.VersionKey, sbom)
			isVehicle := s.ProjectLabelService.HasVehiclePlatformLabel(s.RequestSession, pr)
			evalRes := compsInfo.EvaluatePolicyRules(rules, policyDecisions, isVehicle, sbom.Uploaded, sbom.Key)

			sbomStats = evalRes.Stats
			sbom.Stats = sbomStats
			sbom.TotalStatsHash = currentTotalStatsHash
			s.SBOMListRepo.Update(s.RequestSession, sbomList)
		}
		res.CompStats.AddStats(sbomStats)

		res.Projects = append(res.Projects, approval.ProjectApprovable{
			ProjectKey:         pr.Key,
			ProjectName:        pr.Name,
			CustomerDiffers:    pr.CustomerMeta.Diff(targetProject.CustomerMeta),
			SupplierDiffers:    pr.DocumentMeta.Diff(targetProject.DocumentMeta),
			Supplier:           supplierUserId,
			ApprovableSPDX:     approvableSPDX,
			SpdxName:           sbom.MetaInfo.Name,
			SpdxTag:            sbom.Tag,
			ApprovableStats:    sbomStats,
			SpdxUploaded:       sbom.Uploaded,
			IsSpdxRecent:       isSpdxRecent,
			IsSpdxApprovable:   hasProjectApprovable,
			HasProjectApproval: pr.HasApproval,
		})
	}
	return res
}

func (s *ApprovalService) findLatestSpdx(pr *project.Project) (approvable.ApprovableSPDX, *sbomlist.SbomList, *project.SpdxFileBase) {
	var latest *project.SpdxFileBase
	var latestSBOMList *sbomlist.SbomList
	var latestVersionKey string
	var latestVersionName string

	for _, version := range pr.Versions {
		if version.Deleted {
			continue
		}

		sbomList := s.SBOMListRepo.FindByKey(s.RequestSession, version.Key, false)
		if sbomList == nil {
			continue
		}

		for _, spdx := range sbomList.SpdxFileHistory {
			if latest == nil || spdx.Uploaded.After(*latest.Uploaded) {
				latest = spdx
				latestSBOMList = sbomList
				latestVersionKey = version.Key
				latestVersionName = version.Name
			}
		}
	}

	if latest == nil {
		return approvable.ApprovableSPDX{}, nil, nil
	}

	approvableSpdx := approvable.ApprovableSPDX{
		SpdxKey:     latest.Key,
		VersionKey:  latestVersionKey,
		VersionName: latestVersionName,
	}
	return approvableSpdx, latestSBOMList, latest
}

func (s *ApprovalService) setTaskDone(username string, app *approval.Approval, taskType user2.TaskType, taskStatus user2.TaskStatus) {
	targetUser := s.UserRepo.FindByUserId(s.RequestSession, username)
	targetBefore := targetUser.ToUserAudit()
	task := targetUser.GetTask(app.Key, taskType, taskStatus)
	if task != nil {
		task.Status = user2.TaskDone
		targetAfter := targetUser.ToUserAudit()
		auditHelper.CreateAndAddAuditEntry(&targetUser.Container, app.Creator, message.ApprovalTaskUpdate, audit.DiffWithReporter, targetAfter, targetBefore)
		s.UserRepo.Update(s.RequestSession, targetUser)
	} else {
		logy.Warnf(nil, "setTaskDone but user does not have any task for approval %s %v %v", app.Key, taskType, taskStatus)
	}
}

func (s *ApprovalService) createApprovalCreatorTask(app *approval.Approval) {
	targetUser := s.UserRepo.FindByUserId(s.RequestSession, app.Creator)
	targetUser.AddApprovalCreatorTask(*app)
	s.UserRepo.Update(s.RequestSession, targetUser)
}

func (s *ApprovalService) deletePending(app *approval.Approval) {
	switch app.Type {
	case approval.TypeInternal:
		for i := 0; i < 4; i++ {
			if app.Internal.ApproveStates[i].State != approval.Pending {
				continue
			}
			if app.Internal.Approver[i] == "" {
				continue
			}
			targetUser := app.Internal.GetApproverName(approval.Approver(i))
			s.setTaskDone(targetUser, app, user2.Approval, user2.TaskActive)
			app.Internal.ApproveStates[i].State = approval.Unset
			app.Internal.ApproveStates[i].Updated = nil
		}
	case approval.TypePlausibility:
		s.setTaskDone(app.Plausibility.Approver, app, user2.Approval, user2.TaskActive)
		app.Plausibility.State.State = approval.Unset
		app.Plausibility.State.Updated = nil
	}
}

func (s *ApprovalService) markSbomIsInUse(projects []approval.ProjectApprovable) {
	for _, projectApprovable := range projects {
		spdxKey := projectApprovable.ApprovableSPDX.SpdxKey
		versionKey := projectApprovable.ApprovableSPDX.VersionKey
		if spdxKey == "" || versionKey == "" {
			continue
		}

		sbomList := s.SBOMListRepo.FindByKey(s.RequestSession, versionKey, false)
		if sbomList == nil {
			continue
		}

		for _, sbom := range sbomList.SpdxFileHistory {
			if sbom.Key != spdxKey {
				continue
			}
			if sbom.IsInUse {
				break
			}

			sbom.IsInUse = true
			s.SBOMListRepo.Update(s.RequestSession, sbomList)
			break
		}
	}
}
