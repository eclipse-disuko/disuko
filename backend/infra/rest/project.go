// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain/decisions"
	"github.com/eclipse-disuko/disuko/domain/department"
	policydecisions2 "github.com/eclipse-disuko/disuko/domain/policydecisions"
	"github.com/eclipse-disuko/disuko/infra/repository/policydecisions"
	"go.uber.org/zap/zapcore"

	"github.com/eclipse-disuko/disuko/domain/label"
	"github.com/eclipse-disuko/disuko/domain/project/components"
	sbomlist2 "github.com/eclipse-disuko/disuko/domain/project/sbomlist"
	user2 "github.com/eclipse-disuko/disuko/domain/user"
	"github.com/eclipse-disuko/disuko/helper/s3Helper"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/infra/service/cache"
	sbomLockRetained "github.com/eclipse-disuko/disuko/infra/service/check-sbom-retained"
	checklistService "github.com/eclipse-disuko/disuko/infra/service/checklist"
	"golang.org/x/text/language"

	"github.com/eclipse-disuko/disuko/domain/job"
	licenserules2 "github.com/eclipse-disuko/disuko/domain/licenserules"
	"github.com/eclipse-disuko/disuko/infra/repository/customid"
	"github.com/eclipse-disuko/disuko/infra/repository/jobs"
	"github.com/eclipse-disuko/disuko/infra/repository/licenserules"
	fossddService "github.com/eclipse-disuko/disuko/infra/service/fossdd"
	"github.com/eclipse-disuko/disuko/jobs/fossdd"
	"github.com/eclipse-disuko/disuko/scheduler"

	rt "github.com/eclipse-disuko/disuko/domain/reviewremarks"
	"github.com/eclipse-disuko/disuko/helper"
	audit2 "github.com/eclipse-disuko/disuko/helper/audit"

	"github.com/eclipse-disuko/disuko/connector/application"
	"github.com/eclipse-disuko/disuko/domain/overallreview"

	"github.com/eclipse-disuko/disuko/domain/search"
	"github.com/eclipse-disuko/disuko/helper/filter"
	sort2 "github.com/eclipse-disuko/disuko/helper/sort"
	departmentRepo "github.com/eclipse-disuko/disuko/infra/repository/department"

	"github.com/eclipse-disuko/disuko/infra/repository/dpconfig"
	"github.com/eclipse-disuko/disuko/observermngmt"

	"github.com/eclipse-disuko/disuko/domain"

	"github.com/eclipse-disuko/disuko/infra/repository/analyticscomponents"
	"github.com/eclipse-disuko/disuko/infra/repository/analyticslicenses"
	"github.com/eclipse-disuko/disuko/infra/repository/approvallist"
	"github.com/eclipse-disuko/disuko/infra/repository/auditloglist"

	"github.com/eclipse-disuko/disuko/infra/repository/reviewremarks"

	"github.com/eclipse-disuko/disuko/infra/repository/sbomlist"

	analytics2 "github.com/eclipse-disuko/disuko/domain/analytics"
	approval2 "github.com/eclipse-disuko/disuko/domain/approval"
	"github.com/eclipse-disuko/disuko/domain/audit"
	license2 "github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/domain/oauth"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/project/approvable"
	"github.com/eclipse-disuko/disuko/domain/project/pdocument"
	"github.com/eclipse-disuko/disuko/domain/user/approval"

	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/jwt"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/roles"
	"github.com/eclipse-disuko/disuko/helper/validation"
	"github.com/eclipse-disuko/disuko/infra/repository/analytics"
	"github.com/eclipse-disuko/disuko/infra/repository/labels"
	"github.com/eclipse-disuko/disuko/infra/repository/license"
	"github.com/eclipse-disuko/disuko/infra/repository/obligation"
	"github.com/eclipse-disuko/disuko/infra/repository/policyrules"
	projectRepository "github.com/eclipse-disuko/disuko/infra/repository/project"
	reviewtemplates "github.com/eclipse-disuko/disuko/infra/repository/reviewtemplates"
	"github.com/eclipse-disuko/disuko/infra/repository/schema"
	"github.com/eclipse-disuko/disuko/infra/repository/user"
	sa "github.com/eclipse-disuko/disuko/infra/service/analytics"
	approvalService "github.com/eclipse-disuko/disuko/infra/service/approval"
	"github.com/eclipse-disuko/disuko/infra/service/locks"
	projectService "github.com/eclipse-disuko/disuko/infra/service/project"
	projectLabelService "github.com/eclipse-disuko/disuko/infra/service/project-label"
	"github.com/eclipse-disuko/disuko/infra/service/spdx"
	userService "github.com/eclipse-disuko/disuko/infra/service/user"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

const (
	DiscoBearer        = "DISCO"
	ApprovalLockPrefix = "al-"
)

type ProjectHandler struct {
	// TODO: delete and implement service funcitons for searching analytics
	AnalyticsRepository           analytics.IAnalyticsRepository
	AnalyticsComponentsRepository analyticscomponents.IComponentsRepository
	AnalyticsLicensesRepository   analyticslicenses.ILicensesRepository
	ProjectRepository             projectRepository.IProjectRepository
	SchemaRepository              schema.ISchemaRepository
	LicenseRepository             license.ILicensesRepository
	PolicyRuleRepository          policyrules.IPolicyRulesRepository
	ObligationRepository          obligation.IObligationRepository
	UserRepository                user.IUsersRepository
	LabelRepository               labels.ILabelRepository
	SbomListRepository            sbomlist.ISbomListRepository
	AuditLogListRepository        auditloglist.IAuditLogListRepository
	LockService                   *locks.Service
	ApprovalListRepository        approvallist.IApprovalListRepository
	ReviewRemarksRepository       reviewremarks.IReviewRemarksRepository
	DpConfigRepo                  *dpconfig.DBConfigRepository
	AnalyticsService              *sa.Analytics
	DeparmentRepository           departmentRepo.IDepartmentRepository
	ApplicationConnector          *application.Connector
	ReviewTemplateRepository      reviewtemplates.IReviewTemplateRepository
	LicenseRulesRepository        licenserules.ILicenseRulesRepository
	Scheduler                     *scheduler.Scheduler
	JobRepository                 jobs.IJobsRepository
	SpdxService                   *spdx.Service
	CustomIdRepo                  customid.ICustomIdRepository
	SbomRetainedService           *sbomLockRetained.Service
	ChecklistService              *checklistService.Service
	WizardService                 *projectService.WizardService
	ProjectLabelService           *projectLabelService.ProjectLabelService
	FOSSddService                 *fossddService.Service
	OverallReviewService          *projectService.OverallReviewService
	PolicyDecisionsRepository     policydecisions.IPolicyDecisionsRepository
	UserService                   *userService.Service
}

func (projectHandler *ProjectHandler) ProjectDeprecateHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, true)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	username, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	isOwner := false
	for _, rg := range rights.Groups {
		if rg == string(project.OWNER) {
			isOwner = true
			break
		}
	}
	if !isOwner {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.RequiresOwner))
	}

	if approvalList := projectHandler.ApprovalListRepository.FindByKey(requestSession, currentProject.Key, false); approvalList != nil {
		if hasActiveApprovals(approvalList.Approvals) {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorProjectHasActiveApprovalsOrReviews))
		}
	}

	if currentProject.IsGroup {
		for _, childKey := range currentProject.Children {
			childPrj := projectHandler.ProjectRepository.FindByKey(requestSession, childKey, true)
			if childPrj != nil && !childPrj.IsDeprecated() {
				exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorGroupHasActiveChildren, childPrj.Key))
			}
		}
	}

	oldProject := project.Project{}
	copier.CopyWithOption(&oldProject, currentProject, copier.Option{DeepCopy: true})

	currentProject.DeprecateProject()

	projectHandler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, currentProject.Key, username, message.ProjectUpdated, cmp.Diff, currentProject, &oldProject)

	projectHandler.ProjectRepository.Update(requestSession, currentProject)
	if !hasDummyLabel(currentProject, getDummyLabel(requestSession, projectHandler.LabelRepository)) {
		observermngmt.FireEvent(observermngmt.ProjectUpdated, observermngmt.ProjectUpdatedData{
			RequestSession: requestSession,
			New:            currentProject,
			Old:            &oldProject,
		})
	}

	responseData := SuccessResponse{
		Success: true,
		Message: "Project Status set to DEPRECATED",
	}
	render.JSON(w, r, responseData)
}

func hasActiveApprovals(approvals []approval2.Approval) bool {
	for _, appr := range approvals {
		switch appr.Type {
		case approval2.TypeInternal:
			if appr.Internal.IsActive() {
				return true
			}
		case approval2.TypePlausibility:
			if appr.Plausibility.IsActive() {
				return true
			}
		}
	}
	return false
}

func (projectHandler *ProjectHandler) ProjectRecentHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	username := roles.GetUsernameFromRequest(requestSession, r)
	projects := projectHandler.ProjectRepository.FindRecentByUpdatedForUser(requestSession, username, 5)

	err := roles.FilterProjectsWithoutAccess(requestSession, r, &projects, "")
	if err != nil {
		return
	}

	tokenData := jwt.ExtractTokenMetadata(requestSession, r)
	var result project.ProjectsResponse
	result.Count = len(projects)
	result.Projects = make([]project.ProjectSlimDto, 0)
	dummyLabel := getDummyLabel(requestSession, projectHandler.LabelRepository)
	for _, project := range projects {
		rights, _ := roles.GetProjectAccess(tokenData, project)
		docDep, docMissing, custDep, custMissing := projectHandler.getDeps(requestSession, project)
		slimDto := project.ToSlimDto(docDep, docMissing, custDep, custMissing, hasDummyLabel(project, dummyLabel))
		slimDto.DeleteDisabledReason = projectHandler.CheckProjectDeletionEligibility(requestSession, project, rights)
		// enrich with access rights
		slimDto.AccessRights = *rights
		result.Projects = append(result.Projects, slimDto)
	}

	render.JSON(w, r, result)
}

func (projectHandler *ProjectHandler) ProjectFindApplicableChecklists(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, false)
	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowExecuteChecklist {
		exception.ThrowExceptionSendDeniedResponse()
	}
	lists := projectHandler.ChecklistService.FindApplicableLists(requestSession, currentProject)
	render.JSON(w, r, domain.ToDtos(lists))
}

func (projectHandler *ProjectHandler) ProjectGetPossibleChildrenHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	projectUUIDEscaped := chi.URLParam(r, "uuid")

	projectUUID, _ := url.QueryUnescape(projectUUIDEscaped)
	err := validation.CheckUuid(projectUUID)
	if err != nil {
		return
	}

	username := roles.GetUsernameFromRequest(requestSession, r)
	projects := projectHandler.ProjectRepository.FindAllForUser(requestSession, username)
	err = roles.FilterProjectsWithoutAccess(requestSession, r, &projects, project.OWNER)
	if err != nil {
		return
	}

	tokenData := jwt.ExtractTokenMetadata(requestSession, r)
	var result project.ProjectsResponse
	result.Count = len(projects)

	groupApprovalList := projectHandler.ApprovalListRepository.FindByKey(requestSession, projectUUID, false)
	uniqueProjectKeysUnderApprove := make([]string, 0)
	if groupApprovalList != nil {
		projectsOfGroupUnderApprove := make(map[string]bool)
		for _, approval := range groupApprovalList.Approvals {
			for _, projectApprovable := range approval.Info.Projects {
				projectsOfGroupUnderApprove[projectApprovable.ProjectKey] = true
			}
		}

		if len(projectsOfGroupUnderApprove) > 0 {
			for key := range projectsOfGroupUnderApprove {
				uniqueProjectKeysUnderApprove = append(uniqueProjectKeysUnderApprove, key)
			}
		}
	}

	result.Projects = make([]project.ProjectSlimDto, 0)
	dummyLabel := getDummyLabel(requestSession, projectHandler.LabelRepository)
	for _, project := range projects {
		isPossible := (project.Parent == projectUUID) || (!(len(project.Parent) > 0 || len(project.ParentName) > 0) && !project.IsGroup)
		isDummy := hasDummyLabel(project, dummyLabel)
		if !isPossible || isDummy {
			continue
		}
		rights, _ := roles.GetProjectAccess(tokenData, project)
		docDep, docMissing, custDep, custMissing := projectHandler.getDeps(requestSession, project)
		slimDto := project.ToSlimDto(docDep, docMissing, custDep, custMissing, isDummy)
		if project.Parent == projectUUID && slices.Contains(uniqueProjectKeysUnderApprove, project.Key) {
			slimDto.IsInGroupApproval = true
		}
		// enrich with access rights
		slimDto.AccessRights = *rights
		result.Projects = append(result.Projects, slimDto)
	}

	render.JSON(w, r, result)
}

func (projectHandler *ProjectHandler) ProjectGetAllHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	username := roles.GetUsernameFromRequest(requestSession, r)
	projects := projectHandler.ProjectRepository.FindAllForUser(requestSession, username)

	filtered := make([]*project.Project, len(projects))
	copy(filtered, projects)
	err := roles.FilterProjectsWithoutAccess(requestSession, r, &filtered, "")
	if err != nil {
		return
	}
	tokenData := jwt.ExtractTokenMetadata(requestSession, r)
	var result project.ProjectsResponse
	var docDep *department.Department
	var docMissing bool
	var custDep *department.Department
	var custMissing bool

	result.Count = len(filtered)
	result.Projects = make([]project.ProjectSlimDto, 0)

	dummyLabel := getDummyLabel(requestSession, projectHandler.LabelRepository)
	for _, pr := range filtered {
		rights, _ := roles.GetProjectAccess(tokenData, pr)
		parent := findParent(pr.Parent, projects)

		if parent == nil {
			docDep, docMissing, custDep, custMissing = projectHandler.getDeps(requestSession, pr)
		} else {
			docDep, docMissing, custDep, custMissing = projectHandler.getDeps(requestSession, parent)
		}
		slimDto := pr.ToSlimDto(docDep, docMissing, custDep, custMissing, hasDummyLabel(pr, dummyLabel))

		if !rights.AllowProject.Delete {
			slimDto.DeleteDisabledReason = message.GetI18N(message.DeleteProject).Text
		}

		// enrich with access rights
		slimDto.AccessRights = *rights
		result.Projects = append(result.Projects, slimDto)
	}

	render.JSON(w, r, result)
}

func (projectHandler *ProjectHandler) GetAllDisclosures(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	username := roles.GetUsernameFromRequest(requestSession, r)
	projects := projectHandler.ProjectRepository.FindAllForUser(requestSession, username)

	filtered := make([]*project.Project, len(projects))
	copy(filtered, projects)
	err := roles.FilterProjectsWithoutAccess(requestSession, r, &filtered, "")
	if err != nil {
		return
	}

	tokenData := jwt.ExtractTokenMetadata(requestSession, r)
	var result project.ProjectsResponse
	result.Projects = make([]project.ProjectSlimDto, 0)
	dummyLabel := getDummyLabel(requestSession, projectHandler.LabelRepository)
	for _, project := range filtered {
		if !project.IsGroup {
			continue
		}
		rights, _ := roles.GetProjectAccess(tokenData, project)
		docDep, docMissing, custDep, custMissing := projectHandler.getDeps(requestSession, project)
		slimDto := project.ToSlimDto(docDep, docMissing, custDep, custMissing, hasDummyLabel(project, dummyLabel))
		// enrich with access rights
		slimDto.AccessRights = *rights
		result.Projects = append(result.Projects, slimDto)
	}
	result.Count = len(result.Projects)
	render.JSON(w, r, result)
}

func findParent(key string, prs []*project.Project) *project.Project {
	for _, p := range prs {
		if p.Key == key {
			return p
		}
	}
	return nil
}

func (projectHandler *ProjectHandler) ProjectGetApprovableInfo(w http.ResponseWriter, r *http.Request) {
	pr, requestSession := projectHandler.retrieveProject2(r, false)

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, pr, true)
	// TODO: which access right is needed here?
	if !rights.AllowProject.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadProject))
	}

	as := approvalService.ApprovalService{
		RequestSession:      requestSession,
		ProjectRepo:         projectHandler.ProjectRepository,
		LicenseRepo:         projectHandler.LicenseRepository,
		SBOMListRepo:        projectHandler.SbomListRepository,
		PolicyRulesRepo:     projectHandler.PolicyRuleRepository,
		SpdxRetriever:       projectHandler,
		LicenseRulesRepo:    projectHandler.LicenseRulesRepository,
		SpdxService:         projectHandler.SpdxService,
		ProjectLabelService: projectHandler.ProjectLabelService,
		PolicyDecisionsRepo: projectHandler.PolicyDecisionsRepository,
	}

	res := as.GetApprovalInfo(pr)
	dto := res.ToDto()

	for _, appPr := range dto.Projects {
		policyDecisions := projectHandler.PolicyDecisionsRepository.FindByKey(requestSession, appPr.ProjectKey, false)
		if hasActiveDeniedDecision(policyDecisions) {
			dto.HasDeniedDecisions = true
			break
		}
	}

	render.JSON(w, r, dto)
}

func (projectHandler *ProjectHandler) ProjectGetChildrenHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, false)

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, true)
	if !rights.AllowProject.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadProject))
	}

	var result project.ProjectsChildren
	result.List = make([]project.ProjectChildrenCombiDto, 0)
	result.Projects = make([]project.ProjectSlimDto, 0)
	dummyLabel := getDummyLabel(requestSession, projectHandler.LabelRepository)
	for _, childKey := range currentProject.Children {
		childProject := projectHandler.ProjectRepository.FindByKeyWithDeleted(requestSession, childKey, false)
		if childProject == nil {
			logy.Warnf(requestSession, "Child project not found uuid: %s parent: %s", childKey, currentProject.Key)
			continue
		}
		nothingAdded := true
		custDep, custMissing, docDep, docMissing := projectHandler.getDeps(requestSession, childProject)

		dummy := hasDummyLabel(childProject, dummyLabel)
		result.Projects = append(result.Projects, childProject.ToSlimDto(custDep, custMissing, docDep, docMissing, dummy))

		// evaluate access rights for child project
		hasChildProjectReadAccess := false
		exception.TryCatch(func() {
			_, childProjectRights := roles.GetAndCheckProjectRights(requestSession, r, childProject, false)
			hasChildProjectReadAccess = childProjectRights.AllowProject.Read
		}, func(exc exception.Exception) {
			if exc.ErrorCode != message.ErrorAAR {
				exception.ThrowException(exc)
			}
		})

		for _, version := range childProject.Versions {
			if version.Deleted {
				continue
			}
			nothingAdded = false
			entryVersion := &project.ProjectChildrenCombiDto{
				ProjectKey:           childProject.Key,
				Version:              version.ToDtoWithParentKey(&childProject.Key),
				Project:              childProject.ToSlimDto(custDep, custMissing, docDep, docMissing, dummy),
				HasProjectReadAccess: hasChildProjectReadAccess,
			}
			result.List = append(result.List, *entryVersion)
		}
		if nothingAdded {
			entry := &project.ProjectChildrenCombiDto{
				ProjectKey:           childProject.Key,
				Project:              childProject.ToSlimDto(custDep, custMissing, docDep, docMissing, dummy),
				HasProjectReadAccess: hasChildProjectReadAccess,
			}
			result.List = append(result.List, *entry)
		}
	}

	render.JSON(w, r, result)
}

func (projectHandler *ProjectHandler) ProjectGetHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, true)

	username, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, true)

	if !rights.AllowProject.Read || roles.CheckProjectTypeAccess(requestSession, rights, currentProject, projectHandler.LabelRepository, oauth.AccessLevelRead) {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadProject))
	}

	projectHandler.HandleProjectGet(currentProject, rights, username, w, r)
}

func (projectHandler *ProjectHandler) ProjectGetSettingsHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, true)

	username, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, true)

	if !rights.AllowProject.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadProject))
	}

	projectHandler.HandleProjectGet(currentProject, rights, username, w, r)
}

func (projectHandler *ProjectHandler) ProjectGetUsageInApprovalOrReviewRequest(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, false)

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProject.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadProject))
	}

	isInUsage := projectHandler.isProjectOrVersionInApprovalOrContainsSbomToRetain(requestSession, currentProject, nil)
	render.JSON(w, r, SuccessResponse{
		Success: isInUsage,
		Message: "Project usage in Approval or Review Request",
	})
}

func (projectHandler *ProjectHandler) ProjectDeleteHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, false)

	username, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)

	// Use the centralized deletion eligibility check
	deleteDisabledReason := projectHandler.CheckProjectDeletionEligibility(requestSession, currentProject, rights)
	if deleteDisabledReason != "" {
		exception.ThrowExceptionClientMessage3(message.I18N{
			Code: "PROJECT_DELETE_DISABLED",
			Text: deleteDisabledReason,
		})
	}

	dummy := hasDummyLabel(currentProject, getDummyLabel(requestSession, projectHandler.LabelRepository))

	if !dummy {
		observermngmt.FireEvent(observermngmt.ProjectDeleted, observermngmt.ProjectDeletedData{
			RequestSession: requestSession,
			Project:        currentProject,
		})
		observermngmt.FireEvent(observermngmt.DatabaseEntryAddedOrDeleted, observermngmt.DatabaseSizeChange{
			RequestSession: requestSession,
			CollectionName: projectRepository.ProjectCollectionName,
			Rights:         rights,
			Username:       username,
		})
	}

	if dummy {
		// Delete each SBOM and Cache for each channel
		sbomLists := make([]*sbomlist2.SbomList, 0)
		reviewRemarkLists := make([]*rt.ReviewRemarks, 0)
		for _, v := range currentProject.GetVersions() {
			sbomList := projectHandler.SbomListRepository.FindByKey(requestSession, v.Key, false)
			reviewRemarkList := projectHandler.ReviewRemarksRepository.FindByKey(requestSession, v.Key, false)

			if sbomList == nil && reviewRemarkList == nil {
				continue
			}

			if sbomList != nil {
				for _, sbom := range sbomList.SpdxFileHistory {
					sbomFile := currentProject.GetFilePathSbom(sbom.Key, v.Key)
					exception.TryCatchAndLog(requestSession, func() {
						s3Helper.DeleteFile(requestSession, sbomFile)
					})
					cacheFilePath := fmt.Sprintf(cache.CachePath, sbom.Key)
					exception.TryCatchAndLog(requestSession, func() {
						s3Helper.DeleteFile(requestSession, cacheFilePath)
					})
				}
				sbomLists = append(sbomLists, sbomList)
			}

			if reviewRemarkList != nil {
				reviewRemarkLists = append(reviewRemarkLists, reviewRemarkList)
			}
		}
		if len(sbomLists) > 0 {
			// Bulk deletion (hard) of all SBOM Lists for each version of the project
			s := projectHandler.SbomListRepository.StartSession(base.DeleteSession, 100)
			for _, sbomList := range sbomLists {
				s.AddEnt(sbomList)
			}
			s.EndSession()
		}
		if len(reviewRemarkLists) > 0 {
			// Bulk deletion (hard) of all Review Remarks Lists for each version of the project
			s := projectHandler.ReviewRemarksRepository.StartSession(base.DeleteSession, 100)
			for _, reviewRemarks := range reviewRemarkLists {
				s.AddEnt(reviewRemarks)
			}
			s.EndSession()
		}

		// Delete Approvals from corresponding Approval and User Tasks
		approvalList := projectHandler.ApprovalListRepository.FindByKey(requestSession, currentProject.Key, false)
		if approvalList != nil {
			userTasksMap := make(map[string][]string)
			for _, appr := range approvalList.Approvals {
				taskUuid := appr.Key
				creator := appr.Creator
				// Only Plausi and Internal Approval produce tasks
				if appr.Type == approval2.TypePlausibility {
					addTaskForUser(userTasksMap, creator, taskUuid)
					addTaskForUser(userTasksMap, appr.Plausibility.Approver, taskUuid)
				} else if appr.Type == approval2.TypeInternal {
					addTaskForUser(userTasksMap, creator, taskUuid)
					for _, approver := range appr.Internal.Approver {
						addTaskForUser(userTasksMap, approver, taskUuid)
					}
				}
			}
			// Deletion (hard) of Approval List of the project
			projectHandler.ApprovalListRepository.DeleteHard(requestSession, approvalList.Key)

			// Delete all tasks for users related to this project
			if len(userTasksMap) > 0 {
				users := make([]*user2.User, 0)
				for userId, taskUuids := range userTasksMap {
					usr := projectHandler.UserRepository.FindByUserId(requestSession, userId)
					newUserTasks := make([]user2.Task, 0)
					for _, userTask := range usr.Tasks {
						if !slices.Contains(taskUuids, userTask.TargetGuid) {
							newUserTasks = append(newUserTasks, userTask)
						}
					}
					usr.Tasks = newUserTasks
					users = append(users, usr)
				}
				projectHandler.UserRepository.UpdateList(requestSession, users)
			}
		}

		// Delete document files of the Project
		// Only Internal and External Approvals produce document files
		// Process currentProject.Documents to resolve each belonging file and delete them
		// Take care about document's VersionIndex, Type, Language to collect the all files
		if currentProject.Documents != nil {
			for _, doc := range currentProject.Documents {
				var langTag *language.Tag
				if doc.Lang != "" {
					if t, err := language.Parse(doc.Lang); err == nil {
						langTag = &t
					}
				}
				versionIndex := doc.VersionIndex
				targetFileName := pdocument.GetFileNameWithIndex(doc.Type, doc.ApprovalId, langTag, int(*versionIndex))
				completeFileNameInS3 := currentProject.GetFilePathDocumentForProject(targetFileName)
				exception.TryCatchAndLog(requestSession, func() {
					s3Helper.DeleteFile(requestSession, completeFileNameInS3)
				})
			}
		}

		// Delete stil remaining project related files, if any
		projectPath := currentProject.GetFilePathBaseProject()
		filesCount := s3Helper.CountFiles(requestSession, projectPath).CntFiles
		projectFiles := s3Helper.ListObjects(requestSession, projectPath)
		if filesCount > 0 {
			// Log as an error for the first time to be notified in Grafana Dashboard
			logy.Errorf(requestSession, "Found %d still remaining files for dummy project %s(%s) after deletion, deleting them now. Enhance deletion process to avoid possible data inconsistency", filesCount, currentProject.Name, currentProject.Key)
			for file := range projectFiles {
				if len(file.Key) < 1 {
					// ignore ghost files, sometime happens on S3 Mock
					logy.Errorf(requestSession, "Found file ghost! ")
					continue
				}

				exception.TryCatchAndLog(requestSession, func() {
					s3Helper.DeleteFile(requestSession, file.Key)
				})
			}
		}

		// Hard delete dummy project
		projectHandler.ProjectRepository.DeleteHard(requestSession, currentProject.Key)

		// Hard delete auditlogs
		projectHandler.AuditLogListRepository.DeleteHard(requestSession, currentProject.Key)

		w.WriteHeader(200)
		return
	}

	projectHandler.ProjectRepository.Delete(requestSession, currentProject.Key)

	if currentProject.HasParent() {
		parentProject := projectHandler.ProjectRepository.FindByKey(requestSession, currentProject.Parent, false)
		existingFlag := parentProject.HasChildren
		parentProject.HasChildren = projectHandler.CountChildren(requestSession, parentProject, parentProject.Children) > 0

		if existingFlag != parentProject.HasChildren {
			projectHandler.ProjectRepository.Update(requestSession, parentProject)
		}

	}

	w.WriteHeader(200)
}

func addTaskForUser(userTasksMap map[string][]string, user string, taskUuid string) {
	if user == "" {
		return
	}
	userTasksMap[user] = append(userTasksMap[user], taskUuid)
}

func checkProjectOwnership(requestSession *logy.RequestSession, r *http.Request, proj *project.Project) string {
	username, groupOwnerRights := roles.GetAndCheckProjectRights(requestSession, r, proj, false)
	isOwner := false
	for _, g := range groupOwnerRights.Groups {
		if g == string(project.OWNER) {
			isOwner = true
			break
		}
	}
	if !isOwner {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.CloneProject))
	}
	return username
}

func (projectHandler *ProjectHandler) ComponentSearchHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowProject.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}

	searchRequest := extractAnalyticsSearchRequestBody(r)
	searchResponse := projectHandler.AnalyticsService.Handler.HandleComponentSearch(requestSession, searchRequest.Component, searchRequest.ExactComponent)

	searchResponse.Components = helper.UniqueNonEmptyElementsOf(searchResponse.Components)

	render.JSON(w, r, searchResponse)
}

func (projectHandler *ProjectHandler) ReinitialiseAnalytics(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowTools.Update {
		exception.ThrowExceptionSendDeniedResponse()
	}

	projectHandler.AnalyticsService.Reinitialise(requestSession)

	w.WriteHeader(200)
}

func (projectHandler *ProjectHandler) LicensesSearchHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowProject.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}

	search := extractAnalyticsSearchRequestBody(r)
	searchResponse := projectHandler.AnalyticsService.Handler.HandleLicenseSearch(requestSession, search.License, search.ExactLicense)
	render.JSON(w, r, searchResponse)
}

func (projectHandler *ProjectHandler) ProjectComponentSearchHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowProject.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}

	all, _ := strconv.ParseBool(r.URL.Query().Get("all"))
	userProjects := make([]*project.Project, 0)
	if !all {
		username := roles.GetUsernameFromRequest(requestSession, r)
		userProjects = projectHandler.ProjectRepository.FindAllForUser(requestSession, username)
	}

	type combinedSearchOptions struct {
		AnalyticsRequestSearchOptions analytics2.RequestSearchOptions `json:"analyticsRequestSearchOptions"`
		RequestSearchOptions          search.RequestSearchOptionsNew  `json:"requestSearchOptions"`
	}

	var search combinedSearchOptions
	validation.DecodeAndValidate(r, &search, false)

	if all && (search.AnalyticsRequestSearchOptions.Component == "" && search.AnalyticsRequestSearchOptions.License == "") {
		searchResponse := analytics2.ResponseAnalyticsSearch{
			Success: true,
			Items:   make([]analytics2.SearchResponseItem, 0),
		}
		render.JSON(w, r, searchResponse)
		return
	}
	var (
		offset int64
		limit  int64
	)
	if search.RequestSearchOptions.HasPaginationActive() {
		offset = (search.RequestSearchOptions.Page - 1) * search.RequestSearchOptions.ItemsPerPage
		limit = search.RequestSearchOptions.ItemsPerPage
	}
	var (
		asc     bool
		sortCol string
	)
	if search.RequestSearchOptions.ShouldOrder() {
		asc = search.RequestSearchOptions.IsSortAsc()
		switch search.RequestSearchOptions.SortBy[0].Key {
		case "name":
			sortCol = "ProjectName"
		case "projectVersionName":
			sortCol = "ProjectVersionName"
		case "sbomName":
			sortCol = "SBomName"
		case "componentVersion":
			sortCol = "ComponentVersion"
		}
	}

	searchResponse := projectHandler.AnalyticsService.Search(sa.SearchOptions{
		Rs:          requestSession,
		Component:   search.AnalyticsRequestSearchOptions.Component,
		License:     search.AnalyticsRequestSearchOptions.License,
		ProjectKeys: keys(userProjects),
		Offset:      int(offset),
		Limit:       int(limit),
		SortCol:     sortCol,
		Asc:         asc,
	})
	render.JSON(w, r, searchResponse)
}

func keys(projects []*project.Project) []string {
	keys := make([]string, 0)
	for _, p := range projects {
		keys = append(keys, p.Key)
	}
	return keys
}

func (projectHandler *ProjectHandler) ProjectSearchHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowProject.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}

	requestVersion := r.Header.Get("X-Client-Version")
	if requestVersion == "2.0" {
		var searchOptionsVue3 search.RequestSearchOptionsNew
		validation.DecodeAndValidate(r, &searchOptionsVue3, false)
		projectHandler.searchProject(w, r, requestSession, &searchOptionsVue3)
	} else {
		var searchOptions search.RequestSearchOptions
		validation.DecodeAndValidate(r, &searchOptions, false)
		projectHandler.searchProject(w, r, requestSession, &searchOptions)
	}
}

func (projectHandler *ProjectHandler) searchProject(w http.ResponseWriter, r *http.Request, requestSession *logy.RequestSession, searchOptions search.SortableOptions) {
	projects := projectHandler.ProjectRepository.FindAll(requestSession, true)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)

	err := roles.FilterProjectsWithoutAccess(requestSession, r, &projects, "")
	if err != nil {
		return
	}

	tokenData := jwt.ExtractTokenMetadata(requestSession, r) // ignore error - it handled already if any in early step "helper.FilterProjectsWithoutAccess"
	dtos := make([]project.ProjectSlimDto, 0)
	extractors := map[string]func(dto project.ProjectSlimDto) string{
		"status":  func(dto project.ProjectSlimDto) string { return string(dto.Status) },
		"isGroup": func(dto project.ProjectSlimDto) string { return strconv.FormatBool(dto.IsGroup) },
	}

	policyLabelFilter, freeLabelFilter, schemaLabelFilter := projectHandler.getLabelFilter(requestSession, searchOptions)
	customIdFilter := projectHandler.getCustomIdFilter(searchOptions)

	projectMap := make(map[string]project.ProjectSlimDto)
	dummyLabel := getDummyLabel(requestSession, projectHandler.LabelRepository)

	platformLabels := map[string]*label.Label{
		label.ENTERPRISE_PLATFORM: getEnterprisePlatformLabel(requestSession, projectHandler.LabelRepository),
		label.MOBILE_PLATFORM:     getMobilePlatformLabel(requestSession, projectHandler.LabelRepository),
		label.VEHICLE_PLATFORM:    getVehiclePlatformLabel(requestSession, projectHandler.LabelRepository),
		label.OTHER_PLATFORM:      getOtherPlatformLabel(requestSession, projectHandler.LabelRepository),
	}

	for _, prj := range projects {
		var projectPlatformLabel *label.Label
		for _, platformLabel := range platformLabels {
			if hasPolicyLabel(prj, platformLabel) {
				projectPlatformLabel = platformLabel
				break
			}
		}

		if projectPlatformLabel != nil && !rights.HasProjectTypeAccess(projectPlatformLabel.Name, oauth.AccessLevelRead) {
			continue
		}

		docDep, docMissing, custDep, custMissing := projectHandler.getDeps(requestSession, prj)
		slimDto := prj.ToSlimDto(docDep, docMissing, custDep, custMissing, hasDummyLabel(prj, dummyLabel))

		if filter.MatchesCriteria(slimDto, searchOptions, extractors, nil) ||
			policyLabelFilter(slimDto) || freeLabelFilter(slimDto) || schemaLabelFilter(slimDto.SchemaLabel) ||
			customIdFilter(slimDto) {
			rights, _ := roles.GetProjectAccess(tokenData, prj)

			slimDto.DeleteDisabledReason = projectHandler.CheckProjectDeletionEligibility(requestSession, prj, rights)

			// enrich with access rights
			slimDto.AccessRights = *rights
			// add to projectMap for later reference
			projectMap[slimDto.Key] = slimDto

			dtos = append(dtos, slimDto)
		}
	}
	// enrich child projects with parent information
	for k, possibleChild := range dtos {
		group, exists := projectMap[possibleChild.Parent]
		if exists {
			dtos[k].Supplier = group.Supplier
			dtos[k].Company = group.Company
			dtos[k].Department = group.Department
		}
	}
	result := project.ProjectsResponse{
		Projects: dtos,
		Count:    len(dtos),
	}

	if searchOptions.ShouldOrder() {
		asc := searchOptions.IsSortAsc()
		key := searchOptions.GetSortKey()
		if key == "updated" {
			sort2.Sort(result.Projects, func(dto project.ProjectSlimDto) int64 { return dto.Updated.Unix() }, sort2.Int64LessThan, asc)
		} else if key == "name" {
			sort2.Sort(result.Projects, func(dto project.ProjectSlimDto) string { return dto.Name }, sort2.StringLessThan, asc)
		} else if key == "description" {
			sort2.Sort(result.Projects, func(dto project.ProjectSlimDto) string { return dto.Description }, sort2.StringLessThan, asc)
		} else if key == "created" {
			sort2.Sort(result.Projects, func(dto project.ProjectSlimDto) int64 { return dto.Created.Unix() }, sort2.Int64LessThan, asc)
		} else if key == "status" {
			sort2.Sort(result.Projects, func(dto project.ProjectSlimDto) string { return string(dto.Status) }, sort2.StringLessThan, asc)
		} else if key == "applicationId" {
			sort2.Sort(result.Projects, func(dto project.ProjectSlimDto) string { return *dto.ApplicationId }, sort2.StringLessThan, asc)
		} else if key == "supplier" {
			sort2.Sort(result.Projects, func(dto project.ProjectSlimDto) string { return dto.Supplier }, sort2.StringLessThan, asc)
		} else if key == "company" {
			sort2.Sort(result.Projects, func(dto project.ProjectSlimDto) string { return dto.Company }, sort2.StringLessThan, asc)
		} else if key == "department" {
			sort2.Sort(result.Projects, func(dto project.ProjectSlimDto) string { return dto.Department }, sort2.StringLessThan, asc)
		} else if key == "isGroup" {
			sort2.Sort(result.Projects, func(dto project.ProjectSlimDto) bool { return dto.IsGroup }, sort2.BoolLessThan, asc)
		}
	}

	if searchOptions.HasPaginationActive() && len(result.Projects) > 0 {
		lowIndex := (searchOptions.GetPage() - 1) * searchOptions.GetItemsPerPage()
		highIndex := lowIndex + searchOptions.GetItemsPerPage()
		if highIndex > len(result.Projects) {
			highIndex = len(result.Projects)
		}
		if lowIndex > highIndex {
			lowIndex = 0 // reset page number
		}
		result.Projects = result.Projects[lowIndex:highIndex]
	}

	render.JSON(w, r, result)
}

func (projectHandler *ProjectHandler) getLabelFilter(requestSession *logy.RequestSession, searchOptions search.SortableOptions) (func(dto project.ProjectSlimDto) bool, func(dto project.ProjectSlimDto) bool, func(key string) bool) {
	possibleLabels := projectHandler.LabelRepository.FindAll(requestSession, true)

	labelNames := make(map[string]string)
	for _, label := range possibleLabels {
		labelNames[label.Key] = label.Name
	}

	policyLabelFilter := func(dto project.ProjectSlimDto) bool {
		if !searchOptions.HasFilter() {
			return false
		}
		for _, key := range dto.PolicyLabels {
			labelName, exists := labelNames[key]
			if exists && strings.Contains(labelName, searchOptions.GetFilterString()) {
				return true
			}
		}
		return false
	}

	freeLabelFilter := func(dto project.ProjectSlimDto) bool {
		if !searchOptions.HasFilter() {
			return false
		}
		for _, key := range dto.FreeLabels {
			labelName, exists := labelNames[key]
			if exists && strings.Contains(labelName, searchOptions.GetFilterString()) {
				return true
			}
		}
		return false
	}

	schemaLabelFilter := func(key string) bool {
		if !searchOptions.HasFilter() {
			return false
		}
		labelName, exisits := labelNames[key]
		if exisits && strings.Contains(labelName, searchOptions.GetFilterString()) {
			return true
		}
		return false
	}
	return policyLabelFilter, freeLabelFilter, schemaLabelFilter
}

func (projectHandler *ProjectHandler) getCustomIdFilter(searchOptions search.SortableOptions) func(dto project.ProjectSlimDto) bool {
	return func(dto project.ProjectSlimDto) bool {
		for _, c := range dto.CustomIds {
			if strings.Contains(strings.ToLower(c.Value), strings.ToLower(searchOptions.GetFilterString())) {
				return true
			}
			if strings.Contains(strings.ToLower(c.TechnicalId), strings.ToLower(searchOptions.GetFilterString())) {
				return true
			}
		}
		return false
	}
}

func (projectHandler *ProjectHandler) ProjectFillCustomer(w http.ResponseWriter, r *http.Request) {
	pr, requestSession := projectHandler.retrieveProject2(r, false)
	if pr.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	username, _ := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)

	appId := extractAppIdFromRequest(r)
	lockKey := ApprovalLockPrefix + pr.Key
	logy.Infof(requestSession, "Acquiring lock %s", lockKey)
	l, acquired := projectHandler.LockService.Acquire(locks.Options{
		Key:      lockKey,
		Blocking: true,
		Timeout:  time.Second * 10,
	})
	if !acquired {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ResourceInUse), "")
	}
	logy.Infof(requestSession, "Acquired!")
	defer func() {
		projectHandler.LockService.Release(l)
		logy.Infof(requestSession, "Released lock %s", lockKey)
	}()
	pr, _ = projectHandler.retrieveProject2(r, true)

	fillCustomerBody := extractFillCustomerBody(r)

	as := approvalService.ApprovalService{
		RequestSession:      requestSession,
		ProjectRepo:         projectHandler.ProjectRepository,
		LicenseRepo:         projectHandler.LicenseRepository,
		SBOMListRepo:        projectHandler.SbomListRepository,
		PolicyRulesRepo:     projectHandler.PolicyRuleRepository,
		SpdxRetriever:       projectHandler,
		ApprovalListRepo:    projectHandler.ApprovalListRepository,
		AuditLogListRepo:    projectHandler.AuditLogListRepository,
		SpdxService:         projectHandler.SpdxService,
		ProjectLabelService: projectHandler.ProjectLabelService,
		UserRepo:            projectHandler.UserRepository,
	}
	as.FillRemainingCustomer(pr, appId, username, &fillCustomerBody)
	response := approval.ResponseApprovalDto{
		Success: true,
	}
	render.JSON(w, r, response)
}

func (projectHandler *ProjectHandler) GetApproverUser(w http.ResponseWriter, r *http.Request) {
	_, requestSession := projectHandler.retrieveProject2(r, false)

	username, _ := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)

	appId := extractAppIdFromRequest(r)
	pr, _ := projectHandler.retrieveProject2(r, true)

	as := approvalService.ApprovalService{
		RequestSession:      requestSession,
		ProjectRepo:         projectHandler.ProjectRepository,
		LicenseRepo:         projectHandler.LicenseRepository,
		SBOMListRepo:        projectHandler.SbomListRepository,
		PolicyRulesRepo:     projectHandler.PolicyRuleRepository,
		SpdxRetriever:       projectHandler,
		ApprovalListRepo:    projectHandler.ApprovalListRepository,
		AuditLogListRepo:    projectHandler.AuditLogListRepository,
		SpdxService:         projectHandler.SpdxService,
		UserRepo:            projectHandler.UserRepository,
		ProjectLabelService: projectHandler.ProjectLabelService,
	}

	approver := chi.URLParam(r, "approver")
	u := as.GetApproverUser(pr, appId, username, approver)
	render.JSON(w, r, u)
}

func (projectHandler *ProjectHandler) ProjectUpdateApproval(w http.ResponseWriter, r *http.Request) {
	pr, requestSession := projectHandler.retrieveProject2(r, false)
	if pr.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	username, _ := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)

	appId := extractAppIdFromRequest(r)
	lockKey := ApprovalLockPrefix + pr.Key
	logy.Infof(requestSession, "Acquiring lock %s", lockKey)
	l, acquired := projectHandler.LockService.Acquire(locks.Options{
		Key:      lockKey,
		Blocking: true,
		Timeout:  time.Second * 10,
	})
	if !acquired {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ResourceInUse), "")
	}
	logy.Infof(requestSession, "Acquired!")
	defer func() {
		projectHandler.LockService.Release(l)
		logy.Infof(requestSession, "Released lock %s", lockKey)
	}()
	pr, _ = projectHandler.retrieveProject2(r, true)

	updateApprovalBody := extractUpdateApprovalBody(r)

	as := approvalService.ApprovalService{
		RequestSession:       requestSession,
		ProjectRepo:          projectHandler.ProjectRepository,
		LicenseRepo:          projectHandler.LicenseRepository,
		SBOMListRepo:         projectHandler.SbomListRepository,
		PolicyRulesRepo:      projectHandler.PolicyRuleRepository,
		ApprovalListRepo:     projectHandler.ApprovalListRepository,
		AuditLogListRepo:     projectHandler.AuditLogListRepository,
		UserRepo:             projectHandler.UserRepository,
		SpdxRetriever:        projectHandler,
		SpdxService:          projectHandler.SpdxService,
		ProjectLabelService:  projectHandler.ProjectLabelService,
		FOSSddService:        projectHandler.FOSSddService,
		OverallReviewService: projectHandler.OverallReviewService,
	}

	updatedApproval := as.ProcessRandomApprovalUpdate(pr, appId, username, updateApprovalBody)

	if updatedApproval.Type != approval2.TypeInternal {
		projectHandler.ProjectRepository.Update(requestSession, pr)
	}

	response := approval.ResponseApprovalDto{
		Success: true,
	}

	render.JSON(w, r, response)
}

func (projectHandler *ProjectHandler) ProjectGetApproval(w http.ResponseWriter, r *http.Request) {
	pr, requestSession := projectHandler.retrieveProject2(r, true)

	username, rights := roles.GetAndCheckProjectRights(requestSession, r, pr, true)
	if !rights.AllowProject.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}

	if !rights.IsFossOffice() && !projectHandler.ProjectLabelService.HasVehiclePlatformLabel(requestSession, pr) && pr.GetMember(username) == nil {
		exception.ThrowExceptionSendDeniedResponse()
	}

	appId := extractAppIdFromRequest(r)
	approvalList := projectHandler.ApprovalListRepository.FindByKey(requestSession, pr.Key, false)
	if approvalList == nil {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorDbNotFound))
	}
	for _, app := range approvalList.Approvals {
		if app.Key != appId {
			continue
		}
		creatorFullName := projectHandler.fullNameForUserSafe(requestSession, app.Creator, nil)
		approverFullNames := projectHandler.getApproverFullNames(requestSession, app, nil)
		dto := app.ToDto(creatorFullName, approverFullNames)

		for _, appPr := range dto.Info.Projects {
			policyDecisions := projectHandler.PolicyDecisionsRepository.FindByKey(requestSession, appPr.ProjectKey, false)
			if hasActiveDeniedDecision(policyDecisions) {
				dto.Info.HasDeniedDecisions = true
				break
			}
		}

		render.JSON(w, r, dto)
		return
	}

	exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorDbNotFound))
}

func (projectHandler *ProjectHandler) ProjectGetApprovalList(w http.ResponseWriter, r *http.Request) {
	pr, requestSession := projectHandler.retrieveProject2(r, true)

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, pr, true)
	if !rights.AllowProject.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}

	approvalList := projectHandler.ApprovalListRepository.FindByKey(requestSession, pr.Key, false)
	if approvalList == nil {
		render.JSON(w, r, []interface{}{})
		return
	}
	var (
		res           []approval2.ApprovalDto
		fullnameCache = make(map[string]string)
	)
	for _, app := range approvalList.Approvals {
		approverFullNames := projectHandler.getApproverFullNames(requestSession, app, fullnameCache)
		dto := app.ToDto(projectHandler.fullNameForUserSafe(requestSession, app.Creator, fullnameCache), approverFullNames)
		for _, doc := range pr.GetDocuments() {
			if doc.ApprovalId == app.Key {
				dto.Documents = append(dto.Documents, doc.ToDto())
			}
		}
		res = append(res, dto)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Updated.After(res[j].Updated)
	})

	render.JSON(w, r, res)
}

func (projectHandler *ProjectHandler) fullNameForUserSafe(requestSession *logy.RequestSession, userId string, cache map[string]string) string {
	var res string
	if cache != nil {
		var ok bool
		res, ok = cache[userId]
		if ok {
			return res
		}
	}
	user := projectHandler.UserRepository.FindByUserId(requestSession, userId)
	if user != nil {
		res = fmt.Sprintf("%s %s", user.Forename, user.Lastname)
	} else {
		res = userId
	}
	if cache != nil {
		cache[userId] = res
	}
	return res
}

func (projectHandler *ProjectHandler) ProjectCheckVehicleChildren(w http.ResponseWriter, r *http.Request) {
	pr, requestSession := projectHandler.retrieveProject2(r, true)

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, pr, true)
	if !rights.AllowProject.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}

	if !pr.IsGroup {
		exception.ThrowExceptionBadRequestResponse()
	}

	as := approvalService.ApprovalService{
		RequestSession:      requestSession,
		LabelRepo:           projectHandler.LabelRepository,
		ProjectRepo:         projectHandler.ProjectRepository,
		ProjectLabelService: projectHandler.ProjectLabelService,
	}

	render.JSON(w, r, FoundResponse{
		Found: as.ProjectLabelService.HasVehiclePlatformLabel(requestSession, pr),
	})
}

func (projectHandler *ProjectHandler) GroupOnlyVehicleChildren(w http.ResponseWriter, r *http.Request) {
	pr, requestSession := projectHandler.retrieveProject2(r, true)

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, pr, true)
	if !rights.AllowProject.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}

	if !pr.IsGroup {
		exception.ThrowExceptionBadRequestResponse()
	}

	as := approvalService.ApprovalService{
		RequestSession:      requestSession,
		LabelRepo:           projectHandler.LabelRepository,
		ProjectRepo:         projectHandler.ProjectRepository,
		ProjectLabelService: projectHandler.ProjectLabelService,
	}

	render.JSON(w, r, FoundResponse{
		Found: as.ProjectLabelService.OnlyVehicleChildren(requestSession, pr),
	})
}

func (projectHandler *ProjectHandler) ProjectCreateApproval(w http.ResponseWriter, r *http.Request) {
	typeEscaped := chi.URLParam(r, "approvalType")

	approvalType, err := url.QueryUnescape(typeEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ApprovalTypeWrong))

	currentProject, requestSession := projectHandler.retrieveProject2(r, true)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}
	lockKey := ApprovalLockPrefix + currentProject.Key
	logy.Infof(requestSession, "Acquiring lock %s", lockKey)
	l, acquired := projectHandler.LockService.Acquire(locks.Options{
		Key:      lockKey,
		Blocking: true,
		Timeout:  time.Second * 10,
	})
	if !acquired {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ResourceInUse), "")
	}
	logy.Infof(requestSession, "Acquired!")
	defer func() {
		projectHandler.LockService.Release(l)
		logy.Infof(requestSession, "Released lock %s", lockKey)
	}()
	currentProject, _ = projectHandler.retrieveProject2(r, true)
	username, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	as := approvalService.ApprovalService{
		RequestSession:      requestSession,
		ProjectRepo:         projectHandler.ProjectRepository,
		LicenseRepo:         projectHandler.LicenseRepository,
		SBOMListRepo:        projectHandler.SbomListRepository,
		PolicyRulesRepo:     projectHandler.PolicyRuleRepository,
		SpdxRetriever:       projectHandler,
		ApprovalListRepo:    projectHandler.ApprovalListRepository,
		AuditLogListRepo:    projectHandler.AuditLogListRepository,
		UserRepo:            projectHandler.UserRepository,
		LabelRepo:           projectHandler.LabelRepository,
		SpdxService:         projectHandler.SpdxService,
		LicenseRulesRepo:    projectHandler.LicenseRulesRepository,
		WizardService:       projectHandler.WizardService,
		ProjectLabelService: projectHandler.ProjectLabelService,
		PolicyDecisionsRepo: projectHandler.PolicyDecisionsRepository,
	}

	_, docMissing, _, custMissing := projectHandler.getDeps(requestSession, currentProject)

	dummy := hasDummyLabel(currentProject, getDummyLabel(requestSession, projectHandler.LabelRepository))

	var (
		jobKey  string
		apprKey string
	)
	switch approvalType {
	case approvable.APPROVAL_TYPE_INTERNAL:
		if !rights.AllowRequestApproval.Create {
			exception.ThrowExceptionSendDeniedResponse()
		}
		if custMissing || docMissing {
			exception.ThrowExceptionBadRequestResponse()
		}
		if dummy {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorProjectHasDummyLabel))
		}
		requestApproval := extractRequestInternalApprovalBody(r)
		apprKey = as.CreateInternalApproval(currentProject, requestApproval, username)
		projectHandler.activateTargetProjectOrChildren(requestSession, currentProject)
		jobKey, err = projectHandler.Scheduler.ExecuteOneTimeJob(requestSession, "internal approval doc gen", job.FOSSDDGen, fossdd.Config{
			ProjectID:  currentProject.Key,
			ApprovalID: apprKey,
			WithZIP:    requestApproval.WithZip,
			Template:   requestApproval.FOSSVersion,
		})
		exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorStartingJob))
	case approvable.APPROVAL_TYPE_EXTERNAL:
		if !rights.AllowRequestApproval.Create {
			exception.ThrowExceptionSendDeniedResponse()
		}
		if as.ProjectLabelService.HasVehiclePlatformLabel(requestSession, currentProject) {
			exception.ThrowExceptionBadRequestResponse()
		}
		if custMissing || docMissing {
			exception.ThrowExceptionBadRequestResponse()
		}
		if dummy {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorProjectHasDummyLabel))
		}

		requestApproval := extractRequestExternalApprovalBody(r)
		if !projectHandler.FOSSddService.TemplateExist(requestApproval.FOSSVersion) {
			exception.ThrowExceptionBadRequestResponse()
		}
		apprKey = as.CreateExternalApproval(currentProject, requestApproval, username, false)
		projectHandler.activateTargetProjectOrChildren(requestSession, currentProject)
		jobKey, err = projectHandler.Scheduler.ExecuteOneTimeJob(requestSession, "external approval doc gen", job.FOSSDDGen, fossdd.Config{
			ProjectID:  currentProject.Key,
			ApprovalID: apprKey,
			WithZIP:    requestApproval.WithZip,
			Template:   requestApproval.FOSSVersion,
		})
		exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorStartingJob))
	case approvable.APPROVAL_TYPE_VEHICLE:
		if !rights.AllowRequestApproval.Create {
			exception.ThrowExceptionSendDeniedResponse()
		}
		if custMissing || docMissing {
			exception.ThrowExceptionBadRequestResponse()
		}
		if dummy {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorProjectHasDummyLabel))
		}
		requestApproval := extractRequestExternalApprovalBody(r)
		if !projectHandler.FOSSddService.TemplateExist(requestApproval.FOSSVersion) {
			exception.ThrowExceptionBadRequestResponse()
		}
		apprKey = as.CreateExternalApproval(currentProject, requestApproval, username, true)
		projectHandler.activateTargetProjectOrChildren(requestSession, currentProject)
		jobKey, err = projectHandler.Scheduler.ExecuteOneTimeJob(requestSession, "vehicle approval doc gen", job.FOSSDDGen, fossdd.Config{
			ProjectID:  currentProject.Key,
			ApprovalID: apprKey,
			WithZIP:    true,
			Template:   requestApproval.FOSSVersion,
		})
		exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorStartingJob))
	case approvable.APPROVAL_TYPE_PLAUSI:
		if !rights.AllowRequestPlausi.Create {
			exception.ThrowExceptionSendDeniedResponse()
		}
		requestApproval := extractRequestPlausibilityCheckBody(r)
		if dummy && requestApproval.Approver == conf.Config.Server.FOSSOfficeUserId {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorProjectHasDummyLabel))
		}
		apprKey = as.CreatePlausibilityCheck(currentProject, requestApproval, username)
	default:
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ApprovalTypeWrong))
		return
	}

	if len(currentProject.Children) > 0 {
		projects := projectHandler.ProjectRepository.FindByKeys(requestSession, currentProject.Children, false)

		var projectsToUpdate []*project.Project

		for _, pr := range projects {
			if !pr.HasSBOMToRetain {
				pr.HasSBOMToRetain = true
				projectsToUpdate = append(projectsToUpdate, pr)
			}
		}

		if len(projectsToUpdate) > 0 {
			projectHandler.ProjectRepository.UpdateList(requestSession, projectsToUpdate)
		}
	}

	if approvalType != approvable.APPROVAL_TYPE_PLAUSI {
		currentProject.HasApproval = true
	}
	currentProject.HasSBOMToRetain = true

	projectHandler.ProjectRepository.Update(requestSession, currentProject)

	response := approval.ResponseApprovalDto{
		ApprovalGuid: apprKey,
		Success:      true,
		JobKey:       jobKey,
	}
	render.JSON(w, r, response)
}

func (projectHandler *ProjectHandler) activateTargetProjectOrChildren(requestSession *logy.RequestSession, currentProject *project.Project) {
	currentProject.UpdateStatusToActive()
	projectHandler.ProjectRepository.Update(requestSession, currentProject)
	if currentProject.IsGroup {
		for _, childKey := range currentProject.Children {
			childPrj := projectHandler.ProjectRepository.FindByKey(requestSession, childKey, false)
			if childPrj == nil || childPrj.Status == project.Active {
				continue
			}
			childPrj.UpdateStatusToActive()
			projectHandler.ProjectRepository.Update(requestSession, childPrj)
		}
	}
}

// Deprecated: or not?
func (projectHandler *ProjectHandler) ProjectPostHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	username, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowProject.Create {
		exception.ThrowExceptionSendDeniedResponse()
	}

	projectData := extractRequestBody(r)
	if projectData.Owner == "" {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ProjectOwnerMissing))
	}
	newProject := project.CreateNewProject(projectData)
	projectHandler.UpdateProjectData(requestSession, newProject, projectData)

	auditMsg := message.ProjectCreated
	if projectData.IsGroup {
		auditMsg = message.GroupCreated
	}

	projectAuditEntries := make([]*audit.Audit, 0)
	projectAuditEntries = append(projectAuditEntries, audit2.CreateAuditEntry(username, auditMsg, cmp.Diff, newProject, project.Project{}))

	projectHandler.ProjectRepository.Save(requestSession, newProject)
	observermngmt.FireEvent(observermngmt.DatabaseEntryAddedOrDeleted, observermngmt.DatabaseSizeChange{
		RequestSession: requestSession,
		CollectionName: projectRepository.ProjectCollectionName,
		Rights:         rights,
		Username:       username,
	})
	response := project.Response{
		Id:   newProject.Key,
		Name: newProject.Name,
	}
	taskGuid := ""
	if !newProject.IsGroup {
		// done -> create disco doc and upload to s3
		taskGuid = uuid.NewString()

		createVersion := true // !newProject.IsNoFoss create version every ti
		if createVersion {
			var oldVersions map[string]*project.ProjectVersion
			copier.Copy(&oldVersions, newProject.Versions)
			versionRequestDto := project.VersionRequestDto{
				Name:        message.DefaultBranchMainName,
				Description: message.DefaultBranchMainDescription,
			}
			// todo looks like duplicated code here -> start
			versionKey := newProject.CreateNewProjectVersionIfNameNotUsed(versionRequestDto.Name, versionRequestDto.Description)
			projectAuditEntries = append(projectAuditEntries, audit2.CreateAuditEntry(username, message.ProjectVersionCreated, cmp.Diff, newProject.Versions, oldVersions))
			version := newProject.GetVersion(versionKey)
			projectHandler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, version.Key, username, message.VersionCreated, cmp.Diff, version, project.ProjectVersion{})
			projectHandler.ProjectRepository.Update(requestSession, newProject)
			versionRequestDto = project.VersionRequestDto{
				Name:        message.DefaultBranchDevName,
				Description: message.DefaultBranchDevDescription,
			}
			// todo looks like duplicated code here -> start
			versionKey = newProject.CreateNewProjectVersionIfNameNotUsed(versionRequestDto.Name, versionRequestDto.Description)
			projectAuditEntries = append(projectAuditEntries, audit2.CreateAuditEntry(username, message.ProjectVersionCreated, cmp.Diff, newProject.Versions, oldVersions))
			version = newProject.GetVersion(versionKey)
			projectHandler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, version.Key, username, message.VersionCreated, cmp.Diff, version, project.ProjectVersion{})
			projectHandler.ProjectRepository.Update(requestSession, newProject)
		}
	}
	projectHandler.AuditLogListRepository.CreateAuditEntriesByKey(requestSession, newProject.Key, projectAuditEntries)

	response.TaskGuid = taskGuid
	render.JSON(w, r, response)
}

func (projectHandler *ProjectHandler) CloneProjectPostHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, true)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}
	l, acquired := projectHandler.LockService.Acquire(locks.Options{
		Key:      currentProject.Key,
		Blocking: true,
		Timeout:  10 * time.Second,
	})
	if !acquired {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ResourceInUse), "")
	}
	logy.Infof(requestSession, "Acquired!")
	defer func() {
		projectHandler.LockService.Release(l)
		logy.Infof(requestSession, "Released lock")
	}()

	username := checkProjectOwnership(requestSession, r, currentProject)

	var parentProject *project.Project
	if currentProject.Parent != "" {
		parentProject = projectHandler.ProjectRepository.FindByKey(requestSession, currentProject.Parent, false)
		checkProjectOwnership(requestSession, r, parentProject)
	}

	clonedProject := projectHandler.createClonedProject(currentProject, username)
	projectHandler.ProjectRepository.Save(requestSession, clonedProject)
	projectHandler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, clonedProject.Key, username, message.ProjectCreated, cmp.Diff, clonedProject, project.Project{})

	if parentProject != nil {
		projectHandler.prepareParentProjectForCloning(requestSession, parentProject, clonedProject)
	}

	response := project.Response{
		Id:   clonedProject.Key,
		Name: clonedProject.Name,
	}
	render.JSON(w, r, response)
}

func (projectHandler *ProjectHandler) createClonedProject(originalProject *project.Project, owner string) *project.Project {
	clonedProject := &project.Project{
		RootEntity: domain.NewRootEntity(),
		Name:       originalProject.Name + "-Clone",
		Status:     project.Ready,
	}

	// Set up ownership
	var usersArray []*project.ProjectMemberEntity
	for _, existingUser := range originalProject.UserManagement.Users {
		if existingUser.UserId == owner {
			existingUser.IsResponsible = true
			existingUser.UserType = project.OWNER
		} else {
			existingUser.IsResponsible = false
		}
		usersArray = append(usersArray, existingUser)
	}

	clonedProject.UserManagement = project.UserManagementEntity{
		ChildEntity: domain.NewChildEntity(),
		Users:       usersArray,
	}

	// Copy essential project metadata
	clonedProject.Description = originalProject.Description
	clonedProject.IsGroup = originalProject.IsGroup
	clonedProject.ApplicationMeta = originalProject.ApplicationMeta
	clonedProject.IsNoFoss = originalProject.IsNoFoss

	// Copy labels and configuration
	clonedProject.SchemaLabel = originalProject.SchemaLabel
	if originalProject.PolicyLabels != nil {
		clonedProject.PolicyLabels = make([]string, len(originalProject.PolicyLabels))
		copy(clonedProject.PolicyLabels, originalProject.PolicyLabels)
	}
	if originalProject.ProjectLabels != nil {
		clonedProject.ProjectLabels = make([]string, len(originalProject.ProjectLabels))
		copy(clonedProject.ProjectLabels, originalProject.ProjectLabels)
	}
	if originalProject.FreeLabels != nil {
		clonedProject.FreeLabels = make([]string, len(originalProject.FreeLabels))
		copy(clonedProject.FreeLabels, originalProject.FreeLabels)
	}

	// Initialize empty collections for new project state
	clonedProject.Updated = time.Now()
	clonedProject.Status = project.Ready
	clonedProject.DocumentMeta = originalProject.DocumentMeta
	clonedProject.CustomerMeta = originalProject.CustomerMeta
	clonedProject.SupplierExtraData = originalProject.SupplierExtraData
	clonedProject.Children = make([]string, 0)

	// Clone versions with reset state
	clonedProject.Versions = make(map[string]*project.ProjectVersion)
	for _, originalVersion := range originalProject.Versions {
		if originalVersion.Deleted {
			continue // Skip deleted versions
		}

		clonedVersion := &project.ProjectVersion{
			ChildEntity: domain.NewChildEntity(),
			Name:        originalVersion.Name,
			Description: originalVersion.Description,
			Status:      project.PV_New,
		}

		clonedProject.Versions[clonedVersion.Key] = clonedVersion
	}

	return clonedProject
}

func (projectHandler *ProjectHandler) prepareParentProjectForCloning(rs *logy.RequestSession, parentProject *project.Project, clonedProject *project.Project) {
	parentProject.Children = append(parentProject.Children, clonedProject.Key)
	projectHandler.ProjectRepository.Update(rs, parentProject)

	clonedProject.Parent = parentProject.Key
	clonedProject.ParentName = parentProject.Name
	projectHandler.ProjectRepository.Update(rs, clonedProject)
}

func (projectHandler *ProjectHandler) UpdateProjectData(requestSession *logy.RequestSession, target *project.Project, source project.ProjectRequestDto) {
	CheckIfLabelExistOrThrowException(requestSession, projectHandler.LabelRepository, source.SchemaLabel)
	CheckIfLabelsExistOrThrowException(requestSession, projectHandler.LabelRepository, source.PolicyLabels)
	CheckIfLabelsExistOrThrowException(requestSession, projectHandler.LabelRepository, source.ProjectLabels)

	if projectHandler.ApplicationConnector == nil && source.ApplicationMeta.Diff(target.ApplicationMeta) {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ChangeWithoutConnector))
	}
	if projectHandler.ApplicationConnector != nil && source.ApplicationMeta.Id != "" {
		app := projectHandler.ApplicationConnector.GetApplication(requestSession, source.ApplicationMeta.Id)
		source.ApplicationMeta = project.ApplicationMetaDto{
			Id:           app.Id,
			SecondaryId:  app.SecondaryId,
			Name:         app.Name,
			ExternalLink: app.Link,
		}
	}

	target.UpdateProjectData(source, projectHandler.ApplicationConnector != nil)
	projectHandler.setProjectSettings(target, source.ProjectSettings)
}

func (projectHandler *ProjectHandler) CountChildren(requestSession *logy.RequestSession, currentProject *project.Project, existing []string) int {
	if len(currentProject.Children) == 0 {
		return 0
	}
	return len(currentProject.Children) - projectHandler.getDeletedProjectCount(requestSession, existing)
}

func (projectHandler *ProjectHandler) getDeletedProjectCount(requestSession *logy.RequestSession, existingProjects []string) int {
	projectsWhichAreNotDeleted := projectHandler.ProjectRepository.FindByKeys(requestSession, existingProjects, false)

	return len(existingProjects) - len(projectsWhichAreNotDeleted)
}

func (projectHandler *ProjectHandler) ProjectUpdateHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, true)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}
	l, acquired := projectHandler.LockService.Acquire(locks.Options{
		Key:      currentProject.Key,
		Blocking: true,
		Timeout:  10 * time.Second,
	})
	if !acquired {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ResourceInUse), "")
	}
	logy.Infof(requestSession, "Acquired!")
	defer func() {
		projectHandler.LockService.Release(l)
		logy.Infof(requestSession, "Released lock")
	}()

	oldProject := project.Project{}
	copier.Copy(&oldProject, currentProject)

	username, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProject.Update {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.UpdateProject))
	}

	projectData := extractRequestBody(r)

	projectHandler.UpdateProjectData(requestSession, currentProject, projectData)

	if projectData.IsGroup {
		if !rights.AllowProjectGroup.Update {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.UpdateProjectGroup))
		}
		dummyLabel := getDummyLabel(requestSession, projectHandler.LabelRepository)
		added, removed, existing := currentProject.PrepareUpdateChild(projectData, func(uuid string) bool {
			p := projectHandler.ProjectRepository.FindByKey(requestSession, uuid, true)
			if p == nil {
				return false
			}
			if hasDummyLabel(p, dummyLabel) {
				logy.Warnf(requestSession, "Project '%s'(%s) is dummy - skipping as grouping as not allowed for dummy.", p.Name, p.Key)
				return false
			}
			_, rights := roles.GetAndCheckProjectRights(requestSession, r, p, false)
			for _, g := range rights.Groups {
				if g == string(project.OWNER) {
					return true
				}
			}
			return false
		})

		groupApprovalList := projectHandler.ApprovalListRepository.FindByKey(requestSession, currentProject.Key, false)
		if groupApprovalList != nil {
			projectsOfGroupUnderApprove := make(map[string]bool)
			for _, approval := range groupApprovalList.Approvals {
				for _, projectApprovable := range approval.Info.Projects {
					if slices.Contains(removed, projectApprovable.ProjectKey) {
						projectsOfGroupUnderApprove[projectApprovable.ProjectName] = true
					}
				}
			}

			if len(projectsOfGroupUnderApprove) > 0 {
				projectNames := make([]string, 0)
				for name := range projectsOfGroupUnderApprove {
					projectNames = append(projectNames, name)
				}
				exception.ThrowExceptionClientWithHttpCode(message.ErrorProjectDecoupling, message.GetI18N(message.ErrorProjectDecoupling, strings.Join(projectNames, ", ")).Text, "", exception.HTTP_CODE_SHOW_NO_REQUEST_ID)
			}
		}

		setProjectParentsOfArray(requestSession, projectHandler.ProjectRepository, added, currentProject, true)
		setProjectParentsOfArray(requestSession, projectHandler.ProjectRepository, removed, nil, true)
		if oldProject.Name != projectData.Name {
			setProjectParentsOfArray(requestSession, projectHandler.ProjectRepository, existing, currentProject, true)
		}
		currentProject.HasChildren = projectHandler.CountChildren(requestSession, currentProject, existing) > 0
	}

	projectHandler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, currentProject.Key, username, message.ProjectUpdated, cmp.Diff, currentProject, &oldProject)

	projectHandler.ProjectRepository.Update(requestSession, currentProject)
	dummy := hasDummyLabel(currentProject, getDummyLabel(requestSession, projectHandler.LabelRepository))
	if !dummy {
		observermngmt.FireEvent(observermngmt.ProjectUpdated, observermngmt.ProjectUpdatedData{
			RequestSession: requestSession,
			New:            currentProject,
			Old:            &oldProject,
		})
	}

	docDep, docMissing, custDep, custMissing := projectHandler.getDeps(requestSession, currentProject)
	dto := currentProject.ToDto(docDep, docMissing, custDep, custMissing, dummy)

	// TODO: delete after complete deletion of SpdxValid field as the frontend only uses this for SpdxValid
	for _, v := range dto.Versions {
		sbomList := projectHandler.SbomListRepository.FindByKeyWithDeleted(requestSession, v.Key, false)
		if sbomList == nil {
			continue
		}
		spdxFileHistory := make([]*project.SpdxFileSlimDto, 0)
		for _, spdxFile := range sbomList.SpdxFileHistory {
			spdxFileHistory = append(spdxFileHistory, spdxFile.ToSlimDto(sbomList.Key))
		}
		v.SpdxFileHistory = spdxFileHistory
		if len(v.SpdxFileHistory) > 0 {
			v.CurrentSpdxFile = v.SpdxFileHistory[len(v.SpdxFileHistory)-1]
		}
	}

	// enrich with access rights
	dto.AccessRights = *rights

	if currentProject.HasParent() {
		parentProject := projectHandler.ProjectRepository.FindByKey(requestSession, currentProject.Parent, true)
		if parentProject != nil {
			dto.ParentProjectSettings = parentProject.ToProjectSettingsDto(projectHandler.getDeps(requestSession, parentProject))
		}
	}
	render.JSON(w, r, dto)
}

func (projectHandler *ProjectHandler) ProjectTokenGetAllHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, false)

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !(rights.AllowProjectTokenManagement.Create || rights.AllowAllProjectTokenManagement.Create) {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.CreateToken))
	}

	statusChanged := currentProject.ExpireTokens()
	if statusChanged {
		tokens := currentProject.Token
		currentProjectFull := projectHandler.ProjectRepository.FindByKey(requestSession, currentProject.Key, false)
		currentProjectFull.Token = tokens
		projectHandler.ProjectRepository.UpdateWithoutTimestamp(requestSession, currentProjectFull)
	}

	tokens := make([]project.TokenDto, 0)
	for _, token := range currentProject.Token {
		tokens = append(tokens, token.ToDto())
	}
	render.JSON(w, r, tokens)
}

func (projectHandler *ProjectHandler) ProjectTokenAddHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, true)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	username, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !(rights.AllowProjectTokenManagement.Create || rights.AllowAllProjectTokenManagement.Create) {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.CreateToken))
	}

	tokenData := extractTokenRequestBody(r)
	if tokenData.Expiry != "" && tokenData.IsExpired() {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.TokenExpired), "")
	}

	expiry, _ := tokenData.GetExpired()
	if expiry.After(time.Now().AddDate(2, 0, 0)) {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.TokenExpiryExeedsMax), "")
	}

	newToken := currentProject.GenerateAndAddToken(tokenData)
	tokenAudit := newToken.ToAudit()

	projectHandler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, currentProject.Key, username, message.ProjectTokenCreated, audit.DiffWithReporter, tokenAudit, "")
	projectHandler.ProjectRepository.Update(requestSession, currentProject)
	render.JSON(w, r, newToken.ToDtoWithSecret())
}

func (projectHandler *ProjectHandler) ProjectTokenRenewHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, true)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	username, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !(rights.AllowProjectTokenManagement.Update || rights.AllowAllProjectTokenManagement.Update) {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.RenewToken))
	}

	tokenKey := chi.URLParam(r, "token")

	oldToken, _ := currentProject.GetToken(tokenKey)
	oldTokenAudit := oldToken.ToAudit()

	renewedToken := currentProject.RenewToken(tokenKey)

	renewedTokenAudit := renewedToken.ToAudit()
	projectHandler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, currentProject.Key, username, message.ProjectTokenUpdated, audit.DiffWithReporter, renewedTokenAudit, oldTokenAudit)
	projectHandler.ProjectRepository.Update(requestSession, currentProject)
	render.JSON(w, r, renewedToken.ToDtoWithSecret())
}

func (projectHandler *ProjectHandler) ProjectTokenRevokeHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, true)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	username, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !(rights.AllowProjectTokenManagement.Delete || rights.AllowAllProjectTokenManagement.Delete) {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.RevokingToken))
	}

	tokenKey := chi.URLParam(r, "token")
	oldToken, _ := currentProject.GetToken(tokenKey)
	oldTokenAudit := oldToken.ToAudit()
	currentProject.RevokeToken(tokenKey)

	revokedToken, _ := currentProject.GetToken(tokenKey)
	revokedTokenAudit := revokedToken.ToAudit()
	projectHandler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, currentProject.Key, username, message.ProjectTokenDeleted, audit.DiffWithReporter, revokedTokenAudit, oldTokenAudit)
	projectHandler.ProjectRepository.Update(requestSession, currentProject)

	responseData := SuccessResponse{
		Success: true,
		Message: "Token revoked",
	}
	render.JSON(w, r, responseData)
}

func (projectHandler *ProjectHandler) ProjectDocumentsGetAllHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, false)

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProject.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ProjectRead))
	}

	documents := make([]pdocument.PDocumentDto, 0)
	for _, item := range currentProject.GetDocuments() {
		documents = append(documents, item.ToDto())
	}
	render.JSON(w, r, documents)
}

func (projectHandler *ProjectHandler) ProjectTrailGetAllHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, false)

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProject.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ProjectRead))
	}

	auditEntity := projectHandler.AuditLogListRepository.FindByKey(requestSession, currentProject.Key, false)

	auditTrail := make([]audit.AuditDto, 0)
	if auditEntity != nil && auditEntity.AuditTrail != nil {
		for _, item := range auditEntity.AuditTrail {
			auditTrail = append(auditTrail, item.ToDto())
		}
	}
	render.JSON(w, r, auditTrail)
}

func (projectHandler *ProjectHandler) ProjectSchemaGetHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, false)
	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProject.Read { // todo #751: check if the check is correct - may be AllowProjectVersion
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.FindActiveSchemas))
	}

	activeSchemas := projectHandler.SchemaRepository.FindActiveSchemas(requestSession)

	for _, activeSchema := range activeSchemas {
		if activeSchema.MatchesProjectLabel(currentProject.SchemaLabel) {
			render.JSON(w, r, project.ProjectSchemaResponse{
				Content: activeSchema.Content,
				Name:    activeSchema.Name,
				Version: activeSchema.Version,
			})
			return
		}
	}
}

func extractRequestInternalApprovalBody(r *http.Request) approval2.RequestInternalApprovalDto {
	var dto approval2.RequestInternalApprovalDto
	validation.DecodeAndValidate(r, &dto, false)
	return dto
}

func extractRequestExternalApprovalBody(r *http.Request) approval2.RequestExternalApprovalDto {
	var dto approval2.RequestExternalApprovalDto
	validation.DecodeAndValidate(r, &dto, false)
	return dto
}

func extractRequestPlausibilityCheckBody(r *http.Request) approval2.RequestPlausibilityCheckDto {
	var dto approval2.RequestPlausibilityCheckDto
	validation.DecodeAndValidate(r, &dto, false)
	return dto
}

func extractUpdateApprovalBody(r *http.Request) approval2.UpdateApprovalDto {
	var dto approval2.UpdateApprovalDto
	validation.DecodeAndValidate(r, &dto, false)
	return dto
}

func extractFillCustomerBody(r *http.Request) approval2.FillCustomerDto {
	var dto approval2.FillCustomerDto
	validation.DecodeAndValidate(r, &dto, false)
	return dto
}

func extractRequestBody(r *http.Request) project.ProjectRequestDto {
	var projectData project.ProjectRequestDto
	validation.DecodeAndValidate(r, &projectData, false)
	return projectData
}

func extractTokenRequestBody(r *http.Request) project.Token {
	var tokenData project.Token
	validation.DecodeAndValidate(r, &tokenData, false)
	return tokenData
}

func extractApprovableSPDXBody(r *http.Request) approvable.ApprovableSPDXDto {
	var data approvable.ApprovableSPDXDto
	validation.DecodeAndValidate(r, &data, false)
	return data
}

func assertToken(discoToken, bearer string) string {
	tokenSplit := strings.Split(discoToken, " ")
	if len(tokenSplit) != 2 {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid disco token"), "Malformed token provided")
	}
	if tokenSplit[0] != bearer {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid disco token"), "Malformed token provided")
	}
	return tokenSplit[1]
}

func assertTokenUUID(discoToken, bearer string) string {
	discoToken = assertToken(discoToken, bearer)
	err := validation.CheckUuid(discoToken)
	if err != nil {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid disco token"), err.Error())
	}
	return discoToken
}

func setProjectParentsOfArray(requestSession *logy.RequestSession, projectRepo projectRepository.IProjectRepository, projectUuids []string, parent *project.Project, update bool) {
	var (
		parentName string
		parentKey  string
	)
	if parent != nil {
		parentName = parent.Name
		parentKey = parent.Key
	}

	for _, uuidProject := range projectUuids {
		pr := projectRepo.FindByKeyWithDeleted(requestSession, uuidProject, false)
		if pr == nil {
			continue
		}

		var before project.Project
		copier.Copy(&before, pr)
		if len(pr.Parent) > 0 {
			if pr.Parent != parentKey && len(parentKey) != 0 {
				exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorGroupsParentAlreadyExists))
			}
		}
		pr.Parent = parentKey
		pr.ParentName = parentName
		projectRepo.Update(requestSession, pr)
		if update {
			observermngmt.FireEvent(observermngmt.ProjectUpdated, observermngmt.ProjectUpdatedData{
				RequestSession: requestSession,
				Old:            &before,
				New:            pr,
				NewParent:      parent,
			})
		}
	}
}

func extractProjectKeyFromRequest(r *http.Request) string {
	projectUUIDEscaped := chi.URLParam(r, "uuid")

	projectUUID, err := url.QueryUnescape(projectUUIDEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamUuidWrong))

	err = validation.CheckUuid(projectUUID)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamUuidWrong))
	return projectUUID
}

func extractAppIdFromRequest(r *http.Request) string {
	appUUIDEscaped := chi.URLParam(r, "appId")

	appUUID, err := url.QueryUnescape(appUUIDEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamUuidWrong))

	err = validation.CheckUuid(appUUID)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamUuidWrong))
	return appUUID
}

// ProjectGetExternHandler godoc
//
//	@Summary	Get project details
//	@Id			getProjectDetails
//	@Produce	json
//	@Param		uuid	path		string							true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Success	200		{object}	project.ProjectPublicResponse	"Project"
//	@Failure	404		{object}	exception.HttpError404			"NotFound Error"
//	@Failure	401		{object}	exception.HttpError				"Unauthorized Error"
//	@Router		/v1/projects/{uuid} [get]
//	@security	Bearer
func (projectHandler *ProjectHandler) ProjectGetExternHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	currentProject, _ := projectHandler.retrieveProjectFromPublicRequest(requestSession, r, false)
	projectHandler.HandleProjectGetForPublicResponse(requestSession, currentProject, w, r)
}

// ProjectGetChildrenExternHandler godoc
//
//	@Summary	Get children of a group project
//	@Id			getChildrenProjects
//	@Produce	json
//	@Param		uuid	path		string							true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Success	200		{object}	[]project.ProjectPublicResponse	"Project"
//	@Failure	404		{object}	exception.HttpError404			"NotFound Error"
//	@Failure	401		{object}	exception.HttpError				"Unauthorized Error"
//	@Router		/v1/groups/{uuid}/children [get]
//	@security	Bearer
func (projectHandler *ProjectHandler) ProjectGetChildrenExternHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	currentProject, _ := projectHandler.retrieveProjectFromPublicRequest(requestSession, r, false)
	if !currentProject.IsGroup {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ProjectGroupRequired))
	}
	activeSchemas := projectHandler.SchemaRepository.FindActiveSchemas(requestSession)

	response := make([]project.ProjectPublicResponse, 0)
	for _, chPrj := range currentProject.Children {
		childProject := projectHandler.ProjectRepository.FindByKey(requestSession, chPrj, false)
		if childProject == nil {
			continue
		}

		childProjectPublic := project.ProjectPublicResponse{
			Name:        childProject.Name,
			Uuid:        childProject.UUID(),
			Created:     childProject.Created,
			Updated:     childProject.Updated,
			Description: childProject.Description,
		}
		currentActiveSchema := childProject.FindCorrespondingSchema(activeSchemas)

		if currentActiveSchema != nil {
			childProjectPublic.Schema = currentActiveSchema.Name
		}
		response = append(response, childProjectPublic)
	}
	render.JSON(w, r, response)
}

// ProjectStatusExternHandler godoc
//
//	@Summary	Get project status
//	@Id			getProjectStatus
//	@Produce	json
//	@Param		uuid	path		string								true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Success	200		{object}	project.ProjectStatusPublicResponse	"Project Status"
//	@Failure	404		{object}	exception.HttpError404				"NotFound Error"
//	@Failure	401		{object}	exception.HttpError					"Unauthorized Error"
//	@Router		/v1/projects/{uuid}/status [get]
//	@security	Bearer
func (projectHandler *ProjectHandler) ProjectStatusExternHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	currentProject, _ := retrieveProjectFromPublicRequest(requestSession, projectHandler.ProjectRepository, r, true, false)
	projectHandler.HandleProjectStatus(requestSession, currentProject, w, r)
}

func (projectHandler *ProjectHandler) HandleProjectGet(currentProject *project.Project, rights *oauth.AccessAndRolesRights, username string, w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	activeSchemas := projectHandler.SchemaRepository.FindActiveSchemas(requestSession)

	currentActiveSchema := currentProject.FindCorrespondingSchema(activeSchemas)
	if currentActiveSchema != nil {
		currentProject.CorrespondingSchema = currentActiveSchema
		currentProject.CorrespondingSchema.Content = ""
	}

	docDep, docMissing, custDep, custMissing := projectHandler.getDeps(requestSession, currentProject)
	dto := currentProject.ToDto(docDep, docMissing, custDep, custMissing, hasDummyLabel(currentProject, getDummyLabel(requestSession, projectHandler.LabelRepository)))
	dto.SupplierExtraData.UserFRI = GetUserByUsername(requestSession, projectHandler.UserRepository, dto.SupplierExtraData.FRI)
	dto.SupplierExtraData.UserSRI = GetUserByUsername(requestSession, projectHandler.UserRepository, dto.SupplierExtraData.SRI)

	dto.CustomerMeta.UserFRI = GetUserByUsername(requestSession, projectHandler.UserRepository, dto.CustomerMeta.FRI)
	dto.CustomerMeta.UserSRI = GetUserByUsername(requestSession, projectHandler.UserRepository, dto.CustomerMeta.SRI)

	dto.AccessRights = *rights

	// Check if project deletion is disabled and why
	dto.DeleteDisabledReason = projectHandler.CheckProjectDeletionEligibility(requestSession, currentProject, rights)

	if currentProject.HasParent() {
		parentProject := projectHandler.ProjectRepository.FindByKey(requestSession, currentProject.Parent, true)
		if parentProject != nil {
			parentProjectSettingsDto := parentProject.ToProjectSettingsDto(projectHandler.getDeps(requestSession, parentProject))

			parentProjectSettingsDto.SupplierExtraData.UserFRI = GetUserByUsername(requestSession, projectHandler.UserRepository, parentProject.SupplierExtraData.FRI)
			parentProjectSettingsDto.SupplierExtraData.UserSRI = GetUserByUsername(requestSession, projectHandler.UserRepository, parentProject.SupplierExtraData.SRI)

			parentProjectSettingsDto.CustomerMeta.UserFRI = GetUserByUsername(requestSession, projectHandler.UserRepository, parentProject.CustomerMeta.FRI)
			parentProjectSettingsDto.CustomerMeta.UserSRI = GetUserByUsername(requestSession, projectHandler.UserRepository, parentProject.CustomerMeta.SRI)

			dto.ParentProjectSettings = parentProjectSettingsDto
		}
	}

	m := currentProject.GetMember(username)
	if m != nil {
		subs := m.Subscriptions.ToDto()
		dto.Subscriptions = &subs
	}

	render.JSON(w, r, dto)
}

func (projectHandler *ProjectHandler) HandleProjectGetForPublicResponse(requestSession *logy.RequestSession, currentProject *project.Project, w http.ResponseWriter, r *http.Request) {
	activeSchemas := projectHandler.SchemaRepository.FindActiveSchemas(requestSession)

	responseData := project.ProjectPublicResponse{
		Name:        currentProject.Name,
		Uuid:        currentProject.UUID(),
		Created:     currentProject.Created,
		Updated:     currentProject.Updated,
		Description: currentProject.Description,
		IsGroup:     currentProject.IsGroup,
	}

	currentActiveSchema := currentProject.FindCorrespondingSchema(activeSchemas)
	if currentActiveSchema != nil {
		currentProject.CorrespondingSchema = currentActiveSchema
		currentProject.CorrespondingSchema.Content = ""
		responseData.Schema = currentProject.CorrespondingSchema.Name
	}

	render.JSON(w, r, responseData)
}

func (projectHandler *ProjectHandler) HandleProjectStatus(requestSession *logy.RequestSession, currentProject *project.Project, w http.ResponseWriter, r *http.Request) {
	responseData := project.ProjectStatusPublicResponse{
		Status: currentProject.GetStatus(),
	}

	if len(currentProject.GetVersions()) > 0 {
		versionStatus := make([]project.VersionStatusPublicResponse, 0)
		for _, version := range currentProject.GetVersions() {
			currentVersionStatus := project.VersionStatusPublicResponse{
				Name:   version.Name,
				Status: version.GetStatus(),
			}
			if version.OverallReviews != nil && len(version.OverallReviews) > 0 {
				overallReviews := version.OverallReviews
				sort.Slice(overallReviews, func(i, j int) bool {
					return overallReviews[i].Created.UTC().After(overallReviews[j].Created.UTC())
				})
				currentVersionStatus.OverallReview = &overallreview.OverallReviewPublicResponse{
					SBOMId:       overallReviews[0].SBOMId,
					SBOMName:     overallReviews[0].SBOMName,
					SBOMUploaded: overallReviews[0].SBOMUploaded,
					Comment:      overallReviews[0].Comment,
					Created:      &overallReviews[0].Created,
				}
			}
			sbomList := projectHandler.SbomListRepository.FindByKey(requestSession, version.Key, false)
			if sbomList != nil && len(sbomList.SpdxFileHistory) > 0 {
				spdxFileHistory := sbomList.SpdxFileHistory
				sort.Slice(spdxFileHistory, func(i, j int) bool {
					return spdxFileHistory[i].Uploaded.UTC().After(spdxFileHistory[j].Uploaded.UTC())
				})
				currentVersionStatus.LastSbomUploaded = spdxFileHistory[0].Uploaded
			}
			versionStatus = append(versionStatus, currentVersionStatus)
		}
		responseData.VersionStatus = versionStatus
	}

	render.JSON(w, r, responseData)
}

// ProjectSchemaExternHandler godoc
//
//	@Summary		Get project schema
//	@Id				getProjectSchema
//	@Description	some description
//	@Produce		json
//	@Param			uuid	path		string					true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Success		200		{string}	string					"Schema Details"
//	@Failure		404		{object}	exception.HttpError404	"NotFound Error"
//	@Failure		401		{object}	exception.HttpError		"Unauthorized Error"
//	@Router			/v1/projects/{uuid}/schema [get]
//	@security		Bearer
func (projectHandler *ProjectHandler) ProjectSchemaExternHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	currentProject, _ := projectHandler.retrieveProjectFromPublicRequest(requestSession, r, false)

	activeSchemas := projectHandler.SchemaRepository.FindActiveSchemas(requestSession)

	currentActiveSchema := currentProject.FindCorrespondingSchema(activeSchemas)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	_, err := w.Write([]byte(currentActiveSchema.Content))
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParseSchema))
}

func extractProjectSettings(r *http.Request) *project.ProjectSettingsRequest {
	var projectSettingsRequest project.ProjectSettingsRequest
	validation.DecodeAndValidate(r, &projectSettingsRequest, false)
	return &projectSettingsRequest
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func (projectHandler *ProjectHandler) ProjectUpdateSettingsHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, true)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}
	oldProject := project.Project{}
	copier.Copy(&oldProject, currentProject)

	username, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProject.Update {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.UpdateProject))
	}

	projectSettingsDto := extractProjectSettings(r)

	// uniqId := make(map[string]bool)
	for _, c := range projectSettingsDto.CustomIds {
		// if _, ok := uniqId[c.TechnicalId]; ok {
		// 	exception.ThrowExceptionBadRequestResponse()
		// }
		// uniqId[c.TechnicalId] = true
		dbc := projectHandler.CustomIdRepo.FindByKey(requestSession, c.TechnicalId, false)
		if dbc == nil {
			exception.ThrowExceptionBadRequestResponse()
		}

	}

	projectHandler.setProjectSettings(currentProject, projectSettingsDto)
	currentProject.Updated = time.Now()
	projectHandler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, currentProject.Key, username, message.ProjectUpdated, cmp.Diff, currentProject, &oldProject)
	projectHandler.ProjectRepository.Update(requestSession, currentProject)
	dummyLabel := getDummyLabel(requestSession, projectHandler.LabelRepository)
	if !hasDummyLabel(currentProject, dummyLabel) {
		observermngmt.FireEvent(observermngmt.ProjectUpdated, observermngmt.ProjectUpdatedData{
			RequestSession: requestSession,
			Old:            &oldProject,
			New:            currentProject,
		})
	}
	for _, c := range currentProject.Children {
		child := projectHandler.ProjectRepository.FindByKey(requestSession, c, false)
		if child == nil {
			continue
		}
		if !hasDummyLabel(child, dummyLabel) {
			observermngmt.FireEvent(observermngmt.ProjectUpdated, observermngmt.ProjectUpdatedData{
				RequestSession: requestSession,
				Old:            child,
				New:            child,
				NewParent:      currentProject,
			})
		}
	}

	w.WriteHeader(200)
}

func (projectHandler *ProjectHandler) setProjectSettings(currentProject *project.Project, projectSettingsDto *project.ProjectSettingsRequest) {
	if projectSettingsDto == nil {
		// should not be changed
		return
	}
	currentProject.SetProjectSettings(projectSettingsDto)
}

func (projectHandler *ProjectHandler) ProjectGetPolicyRulesHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, false)
	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, true)

	if !rights.AllowProjectPolicy.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadProject))
	}

	rules := projectHandler.PolicyRuleRepository.FindPolicyRulesForLabel(requestSession, currentProject.PolicyLabels)
	render.JSON(w, r, policyListToDto(rules))
}

func (projectHandler *ProjectHandler) ProjectGetPolicyRulesByIdHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, false)

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, true)

	if !rights.AllowProjectPolicy.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadProject))
	}

	id := chi.URLParam(r, "id")
	policyRule := projectHandler.PolicyRuleRepository.FindByKey(requestSession, id, false)
	if policyRule == nil {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorDbNotFound))
	}

	render.JSON(w, r, policyRule.ToDto())
}

func (p *ProjectHandler) ProjectGetAllSbom(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := p.retrieveProject2(r, true)

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProjectVersion.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ViewSbom))
	}

	newResult := project.ResponseFlatSboms{}
	newResult.Items = make([]project.ResponseFlatSbomItem, 0)
	newResult.Version = make([]project.NameKeyIdentifier, 0)

	var vs []*project.ProjectVersion
	for _, v := range currentProject.Versions {
		vs = append(vs, v)
	}

	sort.Slice(vs, func(i, j int) bool {
		return vs[i].Created.After(vs[j].Created)
	})

	for _, version := range vs {
		if version.Deleted {
			continue
		}
		nameKey := project.NameKeyIdentifier{
			Key:  version.Key,
			Name: version.Name,
		}
		newResult.Version = append(newResult.Version, nameKey)
		sbomList := p.SbomListRepository.FindByKey(requestSession, version.Key, false)
		if sbomList == nil {
			continue
		}

		spdxFileHistory := sbomList.SpdxFileHistory
		sort.Slice(spdxFileHistory, func(i, j int) bool {
			return spdxFileHistory[i].Uploaded.UTC().After(spdxFileHistory[j].Uploaded.UTC())
		})

		unusedSpdxCount := 0
		for _, sbomEntity := range sbomList.SpdxFileHistory {
			spdxFileDto := sbomEntity.ToDto()

			if sbomLockRetained.IsSpdxToRetain(sbomEntity, version) {
				spdxFileDto.IsToRetain = true
			}
			if !IsSpdxInUse(sbomEntity, currentProject, version) {
				if unusedSpdxCount < 5 {
					unusedSpdxCount++
				} else {
					spdxFileDto.IsToDelete = true
				}
			} else {
				spdxFileDto.IsInUse = true
			}

			newResult.Items = append(newResult.Items, project.ResponseFlatSbomItem{
				VersionKey:  version.Key,
				VersionName: version.Name,
				SpdxFileDto: spdxFileDto,
			})
		}
	}
	sort.Slice(newResult.Items, func(i, j int) bool {
		return newResult.Items[i].Uploaded.After(*newResult.Items[j].Uploaded)
	})
	render.JSON(w, r, newResult)
}

func policyListToDto(rules []*license2.PolicyRules) []license2.PolicyRuleDto {
	res := make([]license2.PolicyRuleDto, 0)
	for _, rule := range rules {
		res = append(res, license2.PolicyRuleDto{
			Key:         rule.Key,
			Name:        rule.Name,
			Description: rule.Description,
			Created:     rule.Created,
			Updated:     rule.Updated,
			Auxiliary:   rule.Auxiliary,
			Deprecated:  rule.Deprecated,
			Active:      rule.Active,
		})
	}
	return res
}

func (p *ProjectHandler) ProjectUpdateTaskApprovableSPDX(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := p.retrieveProject2(r, true)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	isOwner := false
	for _, r := range rights.Groups {
		if r == string(project.OWNER) {
			isOwner = true
			break
		}
	}
	if !isOwner {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.RequiresOwner))
	}
	reqData := extractApprovableSPDXBody(r)
	if (reqData.VersionKey != "" && reqData.SpdxKey == "") || (reqData.VersionKey == "" && reqData.SpdxKey != "") {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ApprovableSPDXParamEmpty))
	}
	version, versionOk := currentProject.Versions[reqData.VersionKey]
	if reqData.VersionKey != "" && !versionOk {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ParamVersionWrong))
	}
	if _, spdxFile := p.RetrieveSbomListAndFile(requestSession, reqData.VersionKey, reqData.SpdxKey); reqData.SpdxKey != "" && spdxFile == nil {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.FindingSpdxId, reqData.SpdxKey))
	}
	currentProject.ApprovableSPDX.SpdxKey = reqData.SpdxKey
	currentProject.ApprovableSPDX.VersionKey = reqData.VersionKey
	if version == nil {
		currentProject.ApprovableSPDX.VersionName = ""
	} else {
		currentProject.ApprovableSPDX.VersionName = version.Name
	}

	currentProject.HasSBOMToRetain = true
	p.ProjectRepository.Update(requestSession, currentProject)

	w.WriteHeader(200)
}

func (p *ProjectHandler) getApproverFullNames(requestSession *logy.RequestSession, app approval2.Approval, cache map[string]string) [4]string {
	var approverFullNames [4]string
	switch app.Type {
	case approval2.TypeInternal:
		for i, approver := range app.Internal.Approver {
			approverFullNames[i] = p.fullNameForUserSafe(requestSession, approver, cache)
		}
	case approval2.TypePlausibility:
		approverFullNames[0] = p.fullNameForUserSafe(requestSession, app.Plausibility.Approver, cache)
	}
	return approverFullNames
}

func (p *ProjectHandler) JwtTest(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	rawToken := assertToken(r.Header.Get("Authorization"), "Bearer")
	token := jwt.ExtractTokenMetadataExternal(requestSession, rawToken)
	render.JSON(w, r, token)
}

func extractSubscriptionsReq(r *http.Request) (body project.SubscriptionsDto) {
	validation.DecodeAndValidate(r, &body, false)
	return
}

func (projectHandler *ProjectHandler) SetSubscribedHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, true)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}
	if currentProject.IsGroup {
		exception.ThrowExceptionBadRequestResponse()
	}

	username, _ := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	m := currentProject.GetMember(username)
	if m == nil {
		exception.ThrowExceptionSendDeniedResponse()
	}
	req := extractSubscriptionsReq(r)
	m.Subscriptions = req.ToEntity()

	projectHandler.ProjectRepository.Update(requestSession, currentProject)
	render.JSON(w, r, req)
}

func (projectHandler *ProjectHandler) GetReviewTemplates(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, false)

	username, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if currentProject.GetMember(username) == nil && !rights.IsProjectAnalyst() && !rights.IsDomainAdmin() {
		if !roles.CanAccessVehicleProjectOperations(rights, projectHandler.ProjectLabelService.HasVehiclePlatformLabel(requestSession, currentProject)) {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadReviewTemplates))
		}
	}

	templates := projectHandler.ReviewTemplateRepository.FindAll(requestSession, false)
	templateDtos := make([]rt.ReviewTemplateResponseDto, 0)
	for _, reviewTemplate := range templates {
		templateDtos = append(templateDtos, *reviewTemplate.ToDto())
	}
	render.JSON(w, r, templateDtos)
}

func (projectHandler *ProjectHandler) GetReviewTemplate(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, false)

	username, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if currentProject.GetMember(username) == nil && !rights.IsProjectAnalyst() && !rights.IsDomainAdmin() {
		if !roles.CanAccessVehicleProjectOperations(rights, projectHandler.ProjectLabelService.HasVehiclePlatformLabel(requestSession, currentProject)) {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadReviewTemplates))
		}
	}

	id := chi.URLParam(r, "id")
	reviewTemplate := projectHandler.ReviewTemplateRepository.FindByKey(requestSession, id, false)
	render.JSON(w, r, reviewTemplate.ToDto())
}

func (projectHandler *ProjectHandler) GetDecisions(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, false)

	// todo: Clarify which roles are allowed to see which decisions

	licenseRules := projectHandler.LicenseRulesRepository.FindByKey(requestSession, currentProject.Key, false)
	policyDecisions := projectHandler.PolicyDecisionsRepository.FindByKey(requestSession, currentProject.Key, false)

	result := make([]*decisions.DecisionDto, 0)
	if licenseRules == nil && policyDecisions == nil {
		render.JSON(w, r, result)
		return
	}

	if licenseRules != nil {
		result = append(result, domain.ToDtos(licenseRules.Rules)...)
	}
	if policyDecisions != nil {
		result = append(result, domain.ToDtos(policyDecisions.Decisions)...)
	}

	render.JSON(w, r, result)
}

func (projectHandler *ProjectHandler) CreateBulkPolicyDecisions(w http.ResponseWriter, r *http.Request) {
	currentProject, currentVersion, requestSession := projectHandler.retrieveProjectAndVersion2(r)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	username, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)

	isVehicle := projectHandler.ProjectLabelService.HasVehiclePlatformLabel(requestSession, currentProject)
	isResponsible := currentProject.IsResponsible(username)

	policyDecisionDeniedReason := evaluatePolicyDecisionDeniedReason(isResponsible, rights.IsFossOffice(), isVehicle)
	if len(policyDecisionDeniedReason) > 0 {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.PolicyDecisionOperationNotAuthorized))
	}

	policyDecisionsData := extractBulkPolicyDecisionRequestBody(r, true)
	if len(policyDecisionsData) == 0 {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.InvalidBulkPolicyDecisionData))
	}

	bulkDecision := policyDecisionsData[0].PolicyDecision
	if !strings.EqualFold(bulkDecision, string(license2.ALLOW)) &&
		!strings.EqualFold(bulkDecision, string(license2.DENY)) {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.InvalidPolicyDecisionData))
	}

	sbomId := policyDecisionsData[0].SBOMId

	seenIncoming := make(map[pdKey]struct{}, len(policyDecisionsData))
	newPolicyDecisions := make([]*policydecisions2.PolicyDecision, 0, len(policyDecisionsData))
	for _, pdData := range policyDecisionsData {
		if strings.EqualFold(pdData.PolicyEvaluated, string(license2.DENY)) {
			logy.Errorf(requestSession, "Ignoring 'DENY' Policy Decision while creating Bulk 'WARN' Decisions. Project(UUID)/Version/SBOMId: %s(%s)/%s/%s", currentProject.Name, currentProject.Key, currentVersion.Name, sbomId)
			continue
		}

		if !strings.EqualFold(pdData.PolicyDecision, bulkDecision) {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.InvalidPolicyDecisionData))
		}

		pd := &policydecisions2.PolicyDecision{
			ChildEntity:       domain.NewChildEntity(),
			SBOMId:            pdData.SBOMId,
			SBOMName:          pdData.SBOMName,
			SBOMUploaded:      pdData.SBOMUploaded,
			ComponentSpdxId:   pdData.ComponentSpdxId,
			ComponentName:     pdData.ComponentName,
			ComponentVersion:  pdData.ComponentVersion,
			LicenseExpression: pdData.LicenseExpression,
			LicenseId:         pdData.LicenseId,
			PolicyId:          pdData.PolicyId,
			PolicyEvaluated:   pdData.PolicyEvaluated,
			PolicyDecision:    pdData.PolicyDecision,
			Comment:           pdData.Comment,
			Creator:           pdData.Creator,
			Active:            true,
		}

		k := makePdKey(pd, isVehicle)
		if _, exists := seenIncoming[k]; exists {
			continue
		}
		seenIncoming[k] = struct{}{}

		newPolicyDecisions = append(newPolicyDecisions, pd)
	}

	auditEntries := make([]*audit.Audit, 0)

	policyDecisions := projectHandler.PolicyDecisionsRepository.FindByKey(requestSession, currentProject.Key, false)
	if policyDecisions == nil {
		policyDecisions = &policydecisions2.PolicyDecisions{
			RootEntity: domain.NewRootEntityWithKey(currentProject.Key),
			Decisions:  newPolicyDecisions,
		}
		projectHandler.PolicyDecisionsRepository.Save(requestSession, policyDecisions)

		for _, newPd := range newPolicyDecisions {
			auditEntries = append(auditEntries, audit2.CreateAuditEntry(username, message.PolicyDecisionCreated, cmp.Diff, newPd, policydecisions2.PolicyDecision{}))
		}
	} else {
		existingKeys := make(map[pdKey]struct{}, len(policyDecisions.Decisions))
		for _, pd := range policyDecisions.Decisions {
			if !pd.Active {
				continue
			}
			k := makePdKey(pd, isVehicle)
			existingKeys[k] = struct{}{}
		}

		newPolicyDecisionsToSave := make([]*policydecisions2.PolicyDecision, 0, len(newPolicyDecisions))
		for _, newPd := range newPolicyDecisions {
			k := makePdKey(newPd, isVehicle)
			if _, exists := existingKeys[k]; exists {
				continue
			}
			existingKeys[k] = struct{}{}
			newPolicyDecisionsToSave = append(newPolicyDecisionsToSave, newPd)
			auditEntries = append(auditEntries, audit2.CreateAuditEntry(username, message.PolicyDecisionCreated, cmp.Diff, newPd, policydecisions2.PolicyDecision{}))
		}

		policyDecisions.Decisions = append(policyDecisions.Decisions, newPolicyDecisionsToSave...)
		projectHandler.PolicyDecisionsRepository.Update(requestSession, policyDecisions)
	}

	projectHandler.AuditLogListRepository.CreateAuditEntriesByKey(requestSession, currentProject.Key, auditEntries)
	projectHandler.markSbomUsageFlags(requestSession, currentProject, currentVersion, sbomId)

	render.JSON(w, r, SuccessResponse{
		Success: true,
		Message: "bulk policy decisions created",
	})
}

type pdKey struct {
	cmp   string
	le    string
	ver   string
	lid   string
	ps    string
	psKey string
}

func makePdKey(pd *policydecisions2.PolicyDecision, isVehicle bool) pdKey {
	norm := func(s string) string { return strings.ToLower(strings.TrimSpace(s)) }
	k := pdKey{
		cmp:   norm(pd.ComponentName),
		le:    norm(pd.LicenseExpression),
		lid:   norm(pd.LicenseId),
		ps:    norm(pd.PolicyEvaluated),
		psKey: pd.PolicyId,
	}
	if isVehicle {
		k.ver = norm(pd.ComponentVersion)
	}
	return k
}

func (projectHandler *ProjectHandler) CreatePolicyDecision(w http.ResponseWriter, r *http.Request) {
	currentProject, currentVersion, requestSession := projectHandler.retrieveProjectAndVersion2(r)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	username, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)

	isVehicle := projectHandler.ProjectLabelService.HasVehiclePlatformLabel(requestSession, currentProject)
	isResponsible := currentProject.IsResponsible(username)

	policyDecisionData := extractPolicyDecisionRequestBody(r, true)

	lic := projectHandler.LicenseRepository.FindByIdCaseInsensitive(requestSession, policyDecisionData.LicenseId)
	if lic == nil {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.InvalidPolicyDecisionLicenseData))
	}

	evaluated := policyDecisionData.PolicyEvaluated
	decision := policyDecisionData.PolicyDecision

	if strings.EqualFold(evaluated, string(license2.WARN)) {
		policyDecisionDeniedReason := evaluatePolicyDecisionDeniedReason(isResponsible, rights.IsFossOffice(), isVehicle)
		if len(policyDecisionDeniedReason) > 0 {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.PolicyDecisionOperationNotAuthorized))
		}

		if !strings.EqualFold(decision, string(license2.ALLOW)) &&
			!strings.EqualFold(decision, string(license2.DENY)) {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.InvalidPolicyDecisionData))
		}
	}

	if strings.EqualFold(evaluated, string(license2.DENY)) {
		isAllowDeniedPolicyDecision := evaluateIsAllowDeniedPolicyDecision(rights.IsDomainAdmin(), rights.IsFossOffice(), isVehicle)
		if !isAllowDeniedPolicyDecision {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.PolicyDecisionOperationNotAuthorized))
		}

		if !strings.EqualFold(decision, string(license2.ALLOW)) {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.InvalidPolicyDecisionData))
		}

		if lic.Meta.ApprovalState == license2.Forbidden {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.InvalidPolicyDecisionLicenseApprovalState))
		}
	}

	newPolicyDecision := policydecisions2.PolicyDecision{
		ChildEntity:       domain.NewChildEntity(),
		SBOMId:            policyDecisionData.SBOMId,
		SBOMName:          policyDecisionData.SBOMName,
		SBOMUploaded:      policyDecisionData.SBOMUploaded,
		ComponentSpdxId:   policyDecisionData.ComponentSpdxId,
		ComponentName:     policyDecisionData.ComponentName,
		ComponentVersion:  policyDecisionData.ComponentVersion,
		LicenseExpression: policyDecisionData.LicenseExpression,
		LicenseId:         policyDecisionData.LicenseId,
		PolicyId:          policyDecisionData.PolicyId,
		PolicyEvaluated:   evaluated,
		PolicyDecision:    decision,
		Comment:           policyDecisionData.Comment,
		Creator:           policyDecisionData.Creator,
		Active:            true,
	}

	policyDecisions := projectHandler.PolicyDecisionsRepository.FindByKey(requestSession, currentProject.Key, false)
	if policyDecisions == nil {
		policyDecisions = &policydecisions2.PolicyDecisions{
			RootEntity: domain.NewRootEntityWithKey(currentProject.Key),
			Decisions: []*policydecisions2.PolicyDecision{
				&newPolicyDecision,
			},
		}
		projectHandler.PolicyDecisionsRepository.Save(requestSession, policyDecisions)
	} else {
		for _, prDecision := range policyDecisions.Decisions {

			cmpNameMatches := strings.EqualFold(prDecision.ComponentName, newPolicyDecision.ComponentName)
			licExprMatches := strings.EqualFold(prDecision.LicenseExpression, newPolicyDecision.LicenseExpression)
			versionMatches := !isVehicle || prDecision.ComponentVersion == newPolicyDecision.ComponentVersion
			licMatches := strings.EqualFold(prDecision.LicenseId, newPolicyDecision.LicenseId)
			prStatusMatches := strings.EqualFold(prDecision.PolicyEvaluated, newPolicyDecision.PolicyEvaluated)
			prKeyMatches := prDecision.PolicyId == newPolicyDecision.Key

			allMatches := cmpNameMatches && licExprMatches && versionMatches && licMatches && prStatusMatches && prKeyMatches

			if prDecision.Active && allMatches {
				render.JSON(w, r, SuccessResponse{
					Success: false,
					Message: message.ActivePolicyDecisionExists,
				})
				return
			}
		}

		policyDecisions.Decisions = append(policyDecisions.Decisions, &newPolicyDecision)
		projectHandler.PolicyDecisionsRepository.Update(requestSession, policyDecisions)
	}

	projectHandler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, currentProject.Key, username, message.PolicyDecisionCreated, cmp.Diff, newPolicyDecision, policydecisions2.PolicyDecision{})
	projectHandler.markSbomUsageFlags(requestSession, currentProject, currentVersion, policyDecisionData.SBOMId)
	render.JSON(w, r, SuccessResponse{
		Success: true,
		Message: "policy decision created",
	})
}

func findDecisionByKey(decisions []*policydecisions2.PolicyDecision, key string) *policydecisions2.PolicyDecision {
	for i := range decisions {
		if decisions[i].Key == key {
			return decisions[i]
		}
	}
	return nil
}

func (projectHandler *ProjectHandler) CancelPolicyDecision(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, false)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	username, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)

	isVehicle := projectHandler.ProjectLabelService.HasVehiclePlatformLabel(requestSession, currentProject)
	isResponsible := currentProject.IsResponsible(username)

	policyDecisionUuidEscaped := chi.URLParam(r, "policyDecisionId")
	policyDecisionUuid, err := url.QueryUnescape(policyDecisionUuidEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamPolicyDecisionUuidEmpty))
	if policyDecisionUuid == "" {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ParamPolicyDecisionUuidEmpty))
	}
	err = validation.CheckUuid(policyDecisionUuid)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorKeyRequestParamNotValid, "policyDecisionId"), zapcore.InfoLevel)

	policyDecisions := projectHandler.PolicyDecisionsRepository.FindByKey(requestSession, currentProject.Key, false)
	policyDecision := findDecisionByKey(policyDecisions.Decisions, policyDecisionUuid)

	if policyDecision == nil {
		responseData := SuccessResponse{
			Success: false,
			Message: fmt.Sprintf("policy decision for uuid %s not found", policyDecisionUuid),
		}
		render.JSON(w, r, responseData)
		return
	}

	evaluated := policyDecision.PolicyEvaluated
	if strings.EqualFold(evaluated, string(license2.WARN)) {
		policyDecisionDeniedReason := evaluatePolicyDecisionDeniedReason(isResponsible, rights.IsFossOffice(), isVehicle)
		if len(policyDecisionDeniedReason) > 0 {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.PolicyDecisionOperationNotAuthorized))
		}
	}

	if strings.EqualFold(evaluated, string(license2.DENY)) {
		isAllowDeniedPolicyDecision := evaluateIsAllowDeniedPolicyDecision(rights.IsDomainAdmin(), rights.IsFossOffice(), isVehicle)
		if !isAllowDeniedPolicyDecision {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.PolicyDecisionOperationNotAuthorized))
		}
	}

	oldPolicyDecision := policydecisions2.PolicyDecision{}
	copier.Copy(&oldPolicyDecision, policyDecision)

	policyDecision.Updated = time.Now()
	policyDecision.Active = false

	projectHandler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, currentProject.Key, username, message.PolicyDecisionUpdated, cmp.Diff, policyDecision, &oldPolicyDecision)
	projectHandler.PolicyDecisionsRepository.Update(requestSession, policyDecisions)

	responseData := SuccessResponse{
		Success: true,
		Message: "policy decision canceled",
	}
	render.JSON(w, r, responseData)
}

func (projectHandler *ProjectHandler) CreateLicenseRule(w http.ResponseWriter, r *http.Request) {
	currentProject, currentVersion, requestSession := projectHandler.retrieveProjectAndVersion2(r)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	username, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	member := currentProject.GetMember(username)
	if member == nil {
		if !roles.CanAccessVehicleProjectOperations(rights, projectHandler.ProjectLabelService.HasVehiclePlatformLabel(requestSession, currentProject)) {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.CreateLicenseRule))
		}
	}
	if member != nil && !member.IsResponsible {
		if roles.CanAccessVehicleProjectOperations(rights, projectHandler.ProjectLabelService.HasVehiclePlatformLabel(requestSession, currentProject)) {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.CreateLicenseRule))
		}
	}

	licenseRuleData := extractLicenseRuleRequestBody(r, true)
	licenseRule := licenserules2.LicenseRule{
		ChildEntity:         domain.NewChildEntity(),
		SBOMId:              licenseRuleData.SBOMId,
		SBOMName:            licenseRuleData.SBOMName,
		SBOMUploaded:        licenseRuleData.SBOMUploaded,
		ComponentSpdxId:     licenseRuleData.ComponentSpdxId,
		ComponentName:       licenseRuleData.ComponentName,
		ComponentVersion:    licenseRuleData.ComponentVersion,
		LicenseExpression:   licenseRuleData.LicenseExpression,
		LicenseDecisionId:   licenseRuleData.LicenseDecisionId,
		LicenseDecisionName: licenseRuleData.LicenseDecisionName,
		Comment:             licenseRuleData.Comment,
		Creator:             licenseRuleData.Creator,
		Active:              true,
	}

	if components.GetOperator(strings.ToLower(licenseRule.LicenseExpression)) != components.OR {
		exception.ThrowExceptionBadRequestResponse()
	}

	if len(components.ExtractNames(strings.ToLower(licenseRule.LicenseExpression), components.OR)) > 4 {
		exception.ThrowExceptionBadRequestResponse()
	}

	licenseRules := projectHandler.LicenseRulesRepository.FindByKey(requestSession, currentProject.Key, false)
	if licenseRules == nil {
		licenseRules = &licenserules2.LicenseRules{
			RootEntity: domain.NewRootEntityWithKey(currentProject.Key),
			Rules: []*licenserules2.LicenseRule{
				&licenseRule,
			},
		}
		projectHandler.LicenseRulesRepository.Save(requestSession, licenseRules)
	} else {
		for _, lr := range licenseRules.Rules {

			cmpNameMatches := strings.EqualFold(lr.ComponentName, licenseRule.ComponentName)
			licExprMatches := strings.EqualFold(lr.LicenseExpression, licenseRule.LicenseExpression)

			nameAndExprMatches := cmpNameMatches && licExprMatches

			if lr.Active && nameAndExprMatches {
				render.JSON(w, r, SuccessResponse{
					Success: false,
					Message: message.ActiveLicenseRuleExists,
				})
				return
			}
		}

		licenseRules.Rules = append(licenseRules.Rules, &licenseRule)
		projectHandler.LicenseRulesRepository.Update(requestSession, licenseRules)
	}

	projectHandler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, currentProject.Key, username, message.LicenseRuleCreated, cmp.Diff, licenseRule, licenserules2.LicenseRule{})
	projectHandler.markSbomUsageFlags(requestSession, currentProject, currentVersion, licenseRule.SBOMId)
	render.JSON(w, r, SuccessResponse{
		Success: true,
		Message: "license rule created",
	})
}

func findRuleByKey(rules []*licenserules2.LicenseRule, key string) *licenserules2.LicenseRule {
	for i := range rules {
		if rules[i].Key == key {
			return rules[i]
		}
	}
	return nil
}

func (projectHandler *ProjectHandler) CancelLicenseRule(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, false)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	username, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	member := currentProject.GetMember(username)
	if member == nil || !member.IsResponsible {
		if !roles.CanAccessVehicleProjectOperations(rights, projectHandler.ProjectLabelService.HasVehiclePlatformLabel(requestSession, currentProject)) {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.EditLicenseRule))
		}
	}

	licenseRuleUuidEscaped := chi.URLParam(r, "licenseRuleId")
	licenseRuleUuid, err := url.QueryUnescape(licenseRuleUuidEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamLicenseRuleUuidEmpty))
	if licenseRuleUuid == "" {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ParamLicenseRuleUuidEmpty))
	}
	err = validation.CheckUuid(licenseRuleUuid)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorKeyRequestParamNotValid, "licenseRuleId"), zapcore.InfoLevel)

	licenseRules := projectHandler.LicenseRulesRepository.FindByKey(requestSession, currentProject.Key, false)
	licenseRule := findRuleByKey(licenseRules.Rules, licenseRuleUuid)

	if licenseRule == nil {
		responseData := SuccessResponse{
			Success: false,
			Message: fmt.Sprintf("license rule for uuid %s not found", licenseRuleUuid),
		}
		render.JSON(w, r, responseData)
		return
	}

	oldLicenseRule := licenserules2.LicenseRule{}
	copier.Copy(&oldLicenseRule, licenseRule)

	licenseRule.Updated = time.Now()
	licenseRule.Active = false

	projectHandler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, currentProject.Key, username, message.LicenseRuleUpdated, cmp.Diff, licenseRule, &oldLicenseRule)
	projectHandler.LicenseRulesRepository.Update(requestSession, licenseRules)

	responseData := SuccessResponse{
		Success: true,
		Message: "license rule canceled",
	}
	render.JSON(w, r, responseData)
}

func (projectHandler *ProjectHandler) JobGetOnetimeStatus(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, true)
	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, true)
	if !rights.AllowProject.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadProject))
	}
	key := chi.URLParam(r, "key")
	j := projectHandler.JobRepository.FindByKey(requestSession, key, false)
	if j == nil {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbNotFound), "")
	}
	if j.Execution != job.OneTime {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbNotFound), "")
	}
	if j.JobType != job.FOSSDDGen {
		exception.ThrowExceptionBadRequestResponse()
	}
	var config fossdd.Config
	err := json.Unmarshal([]byte(j.Config), &config)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorJsonDecodingInput))
	if config.ProjectID != currentProject.Key {
		exception.ThrowExceptionBadRequestResponse()
	}
	render.JSON(w, r, job.JobStatusDto{
		Status: j.Status,
	})
}

func extractLicenseRuleRequestBody(r *http.Request, handleErrorAsServerException bool) (data licenserules2.LicenseRuleRequestDto) {
	validation.DecodeAndValidate(r, &data, handleErrorAsServerException)
	return
}

func extractPolicyDecisionRequestBody(r *http.Request, handleErrorAsServerException bool) (data policydecisions2.PolicyDecisionRequestDto) {
	validation.DecodeAndValidate(r, &data, handleErrorAsServerException)
	return
}

func extractBulkPolicyDecisionRequestBody(r *http.Request, handleErrorAsServerException bool) (data []policydecisions2.PolicyDecisionRequestDto) {
	validation.DecodePartAndValidateList(r.Body, &data, handleErrorAsServerException)
	return
}

func hasDummyLabel(currentProject *project.Project, dummyLabel *label.Label) bool {
	if dummyLabel == nil {
		return false
	}
	return slices.Contains(currentProject.ProjectLabels, dummyLabel.GetKey())
}

func hasPolicyLabel(currentProject *project.Project, label *label.Label) bool {
	if label == nil {
		return false
	}
	return slices.Contains(currentProject.PolicyLabels, label.GetKey())
}

func getDummyLabel(requestSession *logy.RequestSession, labelRepository labels.ILabelRepository) *label.Label {
	dummyLabel := labelRepository.FindByNameAndType(requestSession, label.DUMMY, label.PROJECT)
	return dummyLabel
}

func getVehiclePlatformLabel(requestSession *logy.RequestSession, labelRepository labels.ILabelRepository) *label.Label {
	return labelRepository.FindByNameAndType(requestSession, label.VEHICLE_PLATFORM, label.POLICY)
}

func getEnterprisePlatformLabel(requestSession *logy.RequestSession, labelRepository labels.ILabelRepository) *label.Label {
	return labelRepository.FindByNameAndType(requestSession, label.ENTERPRISE_PLATFORM, label.POLICY)
}

func getMobilePlatformLabel(requestSession *logy.RequestSession, labelRepository labels.ILabelRepository) *label.Label {
	return labelRepository.FindByNameAndType(requestSession, label.MOBILE_PLATFORM, label.POLICY)
}

func getOtherPlatformLabel(requestSession *logy.RequestSession, labelRepository labels.ILabelRepository) *label.Label {
	return labelRepository.FindByNameAndType(requestSession, label.OTHER_PLATFORM, label.POLICY)
}

// CheckProjectDeletionEligibility checks if a project can be deleted and returns the reason if it cannot
func (projectHandler *ProjectHandler) CheckProjectDeletionEligibility(
	requestSession *logy.RequestSession,
	currentProject *project.Project,
	rights *oauth.AccessAndRolesRights,
) string {
	// Check delete permission
	if !rights.AllowProject.Delete {
		return message.GetI18N(message.DeleteProject).Text
	}

	dummy := hasDummyLabel(currentProject, getDummyLabel(requestSession, projectHandler.LabelRepository))

	// Check if project is in use (only for non-dummy projects)
	if !dummy {
		isInUsage := projectHandler.isProjectOrVersionInApprovalOrContainsSbomToRetain(requestSession, currentProject, nil)
		if isInUsage {
			return message.GetI18N(message.ErrorInUse, "Project").Text
		}
	}

	// Check if group has children
	if currentProject.IsGroup && len(currentProject.Children) > 0 {
		for _, childKey := range currentProject.Children {
			childProj := projectHandler.ProjectRepository.FindByKey(requestSession, childKey, true)
			if childProj != nil {
				return message.GetI18N(message.ErrorHasChildren).Text
			}
		}
	}

	// Return empty string if deletion is allowed
	return ""
}

func (projectHandler *ProjectHandler) markSbomUsageFlags(requestSession *logy.RequestSession, prj *project.Project, version *project.ProjectVersion, sbomUuid string) {
	projectHandler.markSbomIsInUse(requestSession, version, sbomUuid)
	projectHandler.markProjectSbomRetainFlag(requestSession, prj)
}

func (projectHandler *ProjectHandler) markSbomIsInUse(requestSession *logy.RequestSession, version *project.ProjectVersion, sbomUuid string) {
	sbomList := projectHandler.SbomListRepository.FindByKey(requestSession, version.Key, false)
	if sbomList == nil || len(sbomList.SpdxFileHistory) == 0 {
		exception.ThrowExceptionBadRequestResponse()
	}

	for _, spdx := range sbomList.SpdxFileHistory {
		if spdx.Key != sbomUuid {
			continue
		}
		if spdx.IsInUse {
			return
		}

		spdx.IsInUse = true
		projectHandler.SbomListRepository.Update(requestSession, sbomList)
		return
	}

	exception.ThrowExceptionBadRequestResponse()
}

func (projectHandler *ProjectHandler) markProjectSbomRetainFlag(requestSession *logy.RequestSession, prj *project.Project) {
	if !prj.HasSBOMToRetain {
		prj.HasSBOMToRetain = true
		projectHandler.ProjectRepository.Update(requestSession, prj)
	}
}

func hasActiveDeniedDecision(policyDecisions *policydecisions2.PolicyDecisions) bool {
	if policyDecisions == nil {
		return false
	}
	for _, pd := range policyDecisions.Decisions {
		if pd.Active && pd.PolicyEvaluated == string(license2.DENY) {
			return true
		}
	}
	return false
}
