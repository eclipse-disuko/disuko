// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package approval

import (
	"slices"
	"sort"
	"time"

	"mercedes-benz.ghe.com/foss/disuko/infra/repository/licenserules"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/policydecisions"

	"mercedes-benz.ghe.com/foss/disuko/domain/approval"
	"mercedes-benz.ghe.com/foss/disuko/domain/audit"
	"mercedes-benz.ghe.com/foss/disuko/domain/project"
	"mercedes-benz.ghe.com/foss/disuko/domain/project/components"
	"mercedes-benz.ghe.com/foss/disuko/domain/project/sbomlist"
	user2 "mercedes-benz.ghe.com/foss/disuko/domain/user"
	auditHelper "mercedes-benz.ghe.com/foss/disuko/helper/audit"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/approvallist"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/auditloglist"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/labels"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/license"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/policyrules"
	projectRepo "mercedes-benz.ghe.com/foss/disuko/infra/repository/project"
	sbomListRepo "mercedes-benz.ghe.com/foss/disuko/infra/repository/sbomlist"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/user"
	"mercedes-benz.ghe.com/foss/disuko/infra/service/fossdd"
	projectService "mercedes-benz.ghe.com/foss/disuko/infra/service/project"
	projectLabelService "mercedes-benz.ghe.com/foss/disuko/infra/service/project-label"
	"mercedes-benz.ghe.com/foss/disuko/infra/service/spdx"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type spdxRetriever interface {
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

	SpdxRetriever spdxRetriever

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

func (s *ApprovalService) GetApprovalInfo(targetProject *project.Project) approval.Info {
	return s.getApprovalInfo(targetProject, nil)
}

func (s *ApprovalService) getApprovalInfo(targetProject *project.Project, projectFilter *[]string) approval.Info {
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
		if pr.ApprovableSPDX.SpdxKey == "" || pr.ApprovableSPDX.VersionKey == "" || pr.IsNoFoss {
			res.Projects = append(res.Projects, approval.ProjectApprovable{
				ProjectKey:      pr.Key,
				ProjectName:     pr.Name,
				CustomerDiffers: pr.CustomerMeta.Diff(targetProject.CustomerMeta),
				SupplierDiffers: pr.DocumentMeta.Diff(targetProject.DocumentMeta),
			})
			continue
		}
		sbomList, sbom := s.SpdxRetriever.RetrieveSbomListAndFile(s.RequestSession, pr.ApprovableSPDX.VersionKey, pr.ApprovableSPDX.SpdxKey)
		if sbom == nil {
			res.Projects = append(res.Projects, approval.ProjectApprovable{
				ProjectKey:      pr.Key,
				ProjectName:     pr.Name,
				CustomerDiffers: pr.CustomerMeta.Diff(targetProject.CustomerMeta),
				SupplierDiffers: pr.DocumentMeta.Diff(targetProject.DocumentMeta),
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

		compsInfo := s.SpdxService.GetComponentInfos(s.RequestSession, pr, pr.ApprovableSPDX.VersionKey, sbom)
		rules := s.PolicyRulesRepo.FindPolicyRulesForLabel(s.RequestSession, pr.PolicyLabels)
		policyDecisions := s.PolicyDecisionsRepo.FindByKey(s.RequestSession, pr.Key, false)
		isVehicle := s.ProjectLabelService.HasVehiclePlatformLabel(s.RequestSession, pr)
		evalRes := compsInfo.EvaluatePolicyRules(rules, policyDecisions, isVehicle, sbom.Uploaded, sbom.Key)
		res.CompStats.AddStats(evalRes.Stats)
		res.Projects = append(res.Projects, approval.ProjectApprovable{
			ProjectKey:      pr.Key,
			ProjectName:     pr.Name,
			CustomerDiffers: pr.CustomerMeta.Diff(targetProject.CustomerMeta),
			SupplierDiffers: pr.DocumentMeta.Diff(targetProject.DocumentMeta),
			ApprovableSPDX:  pr.ApprovableSPDX,
			SpdxName:        sbom.MetaInfo.Name,
			SpdxTag:         sbom.Tag,
			ApprovableStats: evalRes.Stats,
			SpdxUploaded:    sbom.Uploaded,
			IsSpdxRecent:    isSpdxRecent,
		})
	}
	return res
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
		if projectApprovable.ApprovableSPDX.SpdxKey == "" || projectApprovable.ApprovableSPDX.VersionKey == "" {
			continue
		}

		changed := false
		sbomList := s.SBOMListRepo.FindByKey(s.RequestSession, projectApprovable.ApprovableSPDX.VersionKey, false)
		for _, sbom := range sbomList.SpdxFileHistory {
			if sbom.Key == projectApprovable.ApprovableSPDX.SpdxKey && !sbom.IsInUse {
				sbom.IsInUse = true
				changed = true
			}
		}
		if changed {
			s.SBOMListRepo.Update(s.RequestSession, sbomList)
		}
	}
}
