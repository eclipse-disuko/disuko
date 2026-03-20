// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/eclipse-disuko/disuko/domain/approval"
	"github.com/eclipse-disuko/disuko/domain/checklist"
	"go.uber.org/zap/zapcore"

	"github.com/eclipse-disuko/disuko/domain/department"
	"github.com/eclipse-disuko/disuko/domain/label"
	"github.com/eclipse-disuko/disuko/observermngmt"

	overallreview2 "github.com/eclipse-disuko/disuko/domain/overallreview"

	"golang.org/x/text/language"

	"github.com/eclipse-disuko/disuko/helper/csvutil"
	"github.com/eclipse-disuko/disuko/helper/rest"

	"github.com/eclipse-disuko/disuko/helper"
	"github.com/eclipse-disuko/disuko/helper/jwt"

	"github.com/eclipse-disuko/disuko/helper/reflection"

	"github.com/google/go-cmp/cmp"

	"github.com/eclipse-disuko/disuko/domain/audit"
	"github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/domain/obligation"
	"github.com/eclipse-disuko/disuko/domain/project/components"
	"github.com/eclipse-disuko/disuko/domain/project/pdocument"
	"github.com/eclipse-disuko/disuko/domain/reviewremarks"
	"github.com/jinzhu/copier"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/roles"
	"github.com/eclipse-disuko/disuko/helper/validation"

	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/infra/repository/database"
	"github.com/eclipse-disuko/disuko/infra/repository/labels"
	"github.com/eclipse-disuko/disuko/infra/service"
	detailsService "github.com/eclipse-disuko/disuko/infra/service/componentDetails"
	"github.com/eclipse-disuko/disuko/infra/service/licenseremarks"
	"github.com/eclipse-disuko/disuko/infra/service/locks"
	projectService "github.com/eclipse-disuko/disuko/infra/service/project"
	reviewRemarksService "github.com/eclipse-disuko/disuko/infra/service/reviewremarks"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/schema"
	"github.com/eclipse-disuko/disuko/helper/s3Helper"
	sbomlockRetained "github.com/eclipse-disuko/disuko/infra/service/check-sbom-retained"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const OrLinksThreshold = 4

func (projectHandler *ProjectHandler) ProjectVersionGetAllHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, true)

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProjectVersion.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadVersions))
	}

	var response project.VersionAllResponse
	response.ProjectVersions = currentProject.GetVersions()
	response.Count = len(response.ProjectVersions)
	render.JSON(w, r, response)
}

func (projectHandler *ProjectHandler) HandleProjectVersionCreate(requestSession *logy.RequestSession, projectModel *project.Project,
	username string, w http.ResponseWriter, r *http.Request, handleErrorAsServerException bool,
) {
	if len(projectModel.GetVersions()) >= conf.Config.Server.MaxVersions {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.MaxVersionsReached))
	}

	versionRequestDto := extractVersionRequestBody(r, handleErrorAsServerException)

	var oldVersions map[string]*project.ProjectVersion
	copier.Copy(&oldVersions, projectModel.Versions)

	versionKey := projectModel.CreateNewProjectVersionIfNameNotUsed(versionRequestDto.Name, versionRequestDto.Description)

	projectHandler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, projectModel.Key, username, message.ProjectVersionCreated, cmp.Diff, projectModel.Versions, oldVersions)
	version := projectModel.GetVersion(versionKey)
	projectHandler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, version.Key, username, message.VersionCreated, cmp.Diff, version, project.ProjectVersion{})

	projectHandler.ProjectRepository.Update(requestSession, projectModel)

	responseData := project.VersionCreationResponseMin{
		Success: true,
		Message: "version created",
		Name:    version.Name,
		Uuid:    version.Key,
	}
	render.JSON(w, r, responseData)
}

// ProjectVersionExternCreateHandler godoc
//
//	@Summary	Create project version (also known as channel)
//	@Id			createProjectVersion
//	@Produce	json
//	@Accept		json
//	@Param		uuid	path		string						true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Param		version	body		project.VersionRequestDto	true	"Version (also known as Channel) Details"
//	@Success	200		{object}	rest.SuccessResponse		"Success Response"
//	@Failure	417		{object}	exception.HttpError			"The channel already exist"
//	@Failure	401		{object}	exception.HttpError			"Unauthorized Error"
//	@Router		/projects/{uuid}/versions [post]
//	@security	Bearer
func (projectHandler *ProjectHandler) ProjectVersionExternCreateHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	existingProject, _ := projectHandler.retrieveProjectFromPublicRequest(requestSession, r, true)
	if existingProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}
	projectHandler.HandleProjectVersionCreate(requestSession, existingProject, project.OriginApi, w, r, false)
}

// ProjectSbomSearchHandler godoc
//
//	@Summary	Search SBOM by tag
//	@Id			searchProjectSbom
//	@Produce	json
//	@Accept		json
//	@Param		uuid	path		string						true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Param		tag		body		project.ProjectSearchDto	true	"Tag to be searched"
//	@Success	200		{array}		project.ProjectSearchResDto	"Result list"
//	@Failure	401		{object}	exception.HttpError			"Unauthorized Error"
//	@Router		/projects/{uuid}/search [post]
//	@security	Bearer
func (projectHandler *ProjectHandler) ProjectSbomSearchHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	existingProject, _ := projectHandler.retrieveProjectFromPublicRequest(requestSession, r, true)
	if existingProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	req := extractProjectSearchBody(r, true)

	res := make([]project.ProjectSearchResDto, 0)

	for _, v := range existingProject.Versions {
		sbomList := projectHandler.SbomListRepository.FindByKey(requestSession, v.Key, false)
		if sbomList == nil {
			continue
		}
		for _, sbom := range sbomList.SpdxFileHistory {
			if sbom.Tag == req.Tag {
				res = append(res, project.ProjectSearchResDto{
					ChannelId: v.Key,
					SbomId:    sbom.Key,
				})
			}
		}
	}

	render.JSON(w, r, res)
}

func (projectHandler *ProjectHandler) ProjectVersionCreateHandler(w http.ResponseWriter, r *http.Request) {
	existingProject, requestSession := projectHandler.retrieveProject2(r, true)
	if existingProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}
	username, rights := roles.GetAndCheckProjectRights(requestSession, r, existingProject, false)
	if !rights.AllowProjectVersion.Create {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.CreateVersion))
	}

	projectHandler.HandleProjectVersionCreate(requestSession, existingProject, username, w, r, true)
}

// ProjectVersionGetListExternHandler godoc
//
//	@Summary	Get project versions (also known as channels)
//	@Id			getProjectVersions
//	@Produce	json
//	@Param		uuid	path		string					true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Success	200		{array}		string					"Versions (also known as Channels)"
//	@Failure	404		{object}	exception.HttpError404	"NotFound Error"
//	@Failure	401		{object}	exception.HttpError		"Unauthorized Error"
//	@Router		/projects/{uuid}/versions [get]
//	@security	Bearer
func (projectHandler *ProjectHandler) ProjectVersionGetListExternHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	currentProject, _ := projectHandler.retrieveProjectFromPublicRequest(requestSession, r, true)
	responseData := make([]string, 0)

	for _, version := range currentProject.Versions {
		if version.Deleted {
			continue
		}
		responseData = append(responseData, version.Name)
	}
	sort.Slice(responseData, func(i, j int) bool {
		return strings.Compare(responseData[i], responseData[j]) > 0
	})

	render.JSON(w, r, responseData)
}

// ProjectVersionGetListExternHandlerV2 godoc
//
//	@Summary	Get project versions (also known as channels) with details
//	@Id			getProjectVersionsV2
//	@Produce	json
//	@Param		uuid	path		string								true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Success	200		{array}		project.VersionPublicResponseMin	"Versions (also known as Channels) with UUID and Name"
//	@Failure	404		{object}	exception.HttpError404				"NotFound Error"
//	@Failure	401		{object}	exception.HttpError					"Unauthorized Error"
//	@Router		/projects/{uuid}/versions [get]
//	@security	Bearer
func (projectHandler *ProjectHandler) ProjectVersionGetListExternHandlerV2(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	currentProject, _ := projectHandler.retrieveProjectFromPublicRequest(requestSession, r, true)
	responseData := make([]project.VersionPublicResponseMin, 0)

	for _, version := range currentProject.Versions {
		if version.Deleted {
			continue
		}
		versionDetails := project.VersionPublicResponseMin{
			Name: version.Name,
			Uuid: version.Key,
		}
		responseData = append(responseData, versionDetails)
	}
	sort.Slice(responseData, func(i, j int) bool {
		return strings.Compare(responseData[i].Name, responseData[j].Name) > 0
	})

	render.JSON(w, r, responseData)
}

// ProjectVersionGetExternHandler godoc
//
//	@Summary	Get version (also known as channel) details of project
//	@Id			getProjectVersionDetails
//	@Produce	json
//	@Param		uuid	path		string							true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Param		version	path		string							true	"Project Version Name (also known as Channel Name) e.g.: main"
//	@Success	200		{object}	project.VersionPublicResponse	"Version (also known as Channel) Details"
//	@Failure	404		{object}	exception.HttpError404			"NotFound Error"
//	@Failure	401		{object}	exception.HttpError				"Unauthorized Error"
//	@Router		/projects/{uuid}/versions/{version} [get]
//	@security	Bearer
func (projectHandler *ProjectHandler) ProjectVersionGetExternHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, version, _ := projectHandler.retrieveProjectAndVersionFromPublicRequest(requestSession, r)

	// projectHandler.ExtractPolicyRuleStatus(requestSession, currentProject, version)

	responseData := project.VersionPublicResponse{
		Name:        version.Name,
		Description: version.Description,
		Status:      version.GetStatus(),
		Uuid:        version.Key,
	}

	if version.OverallReviews != nil && len(version.OverallReviews) > 0 {
		overallReviews := version.OverallReviews
		sort.Slice(overallReviews, func(i, j int) bool {
			return overallReviews[i].Created.UTC().After(overallReviews[j].Created.UTC())
		})
		responseData.OverallReview = &overallreview2.OverallReviewPublicResponse{
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
		responseData.LastSbomUploaded = spdxFileHistory[0].Uploaded
	}
	render.JSON(w, r, responseData)
}

// CCSGetListExternHandler godoc
//
//	@Summary	Get external references to source code resources
//	@Id			getProjectVersionExternalSourceCodeReferences
//	@Produce	json
//	@Param		uuid	path		string									true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Param		version	path		string									true	"Project Version Name (also known as Channel Name) e.g.: main"
//	@Success	200		{array}		project.ExternalSourcePublicResponseDto	"External Source"
//	@Failure	404		{object}	exception.HttpError404					"NotFound Error"
//	@Failure	401		{object}	exception.HttpError						"Unauthorized Error"
//	@Router		/projects/{uuid}/versions/{version}/ccs [get]
//	@security	Bearer
func (projectHandler *ProjectHandler) CCSGetListExternHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, version, _ := projectHandler.retrieveProjectAndVersionFromPublicRequest(requestSession, r)

	result := make([]project.ExternalSourcePublicResponseDto, 0)
	for _, source := range version.SourceExternal {
		result = append(result, source.ToPublicDTO())
	}

	render.JSON(w, r, result)
}

// CCSCreateExternHandler godoc
//
//	@Summary	Create external reference to source code resources
//	@Id			createProjectVersionExternalSourceCodeReferences
//	@Produce	json
//	@Param		uuid	path		string						true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Param		version	path		string						true	"Project Version Name (also known as Channel Name) e.g.: main"
//	@Param		source	body		project.SourceExternalDTO	true	"Source"
//	@Success	200		{object}	rest.SuccessResponse		"Success Response"
//	@Failure	401		{object}	exception.HttpError			"Unauthorized Error"
//	@Router		/projects/{uuid}/versions/{version}/ccs [post]
//	@security	Bearer
func (projectHandler *ProjectHandler) CCSCreateExternHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	projectData, version, token := projectHandler.retrieveProjectAndVersionFromPublicRequest(requestSession, r)
	if projectData.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	l, acquired := projectHandler.LockService.Acquire(locks.Options{
		Key:      projectData.Key,
		Blocking: true,
		Timeout:  time.Second * 10,
	})
	if !acquired {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ResourceInUse), "")
	}
	logy.Infof(requestSession, "Acquired!")
	defer func() {
		projectHandler.LockService.Release(l)
		logy.Infof(requestSession, "Released lock")
	}()
	projectData, version, _ = projectHandler.retrieveProjectAndVersionFromPublicRequest(requestSession, r)

	origin := helper.CreateOrigin(project.OriginApi, token.Company, token.Description, token.Key)
	ip := jwt.TrimPortFromRemoteAddress(r.RemoteAddr)

	projectHandler.HandleExternalSourceCreate(requestSession, projectData, version, origin, ip, w, r, false)
}

func (projectHandler *ProjectHandler) ProjectVersionTrailGetAllHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, version, requestSession := projectHandler.retrieveProjectAndVersion2(r)
	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProjectVersion.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadVersions))
	}

	auditEntity := projectHandler.AuditLogListRepository.FindByKey(requestSession, version.Key, false)

	auditTrail := make([]audit.AuditDto, 0)
	if auditEntity != nil && auditEntity.AuditTrail != nil {
		for _, item := range auditEntity.AuditTrail {
			auditTrail = append(auditTrail, item.ToDto())
		}
	}
	render.JSON(w, r, auditTrail)
}

func (projectHandler *ProjectHandler) ProjectVersionGetHandler(w http.ResponseWriter, r *http.Request) {
	withProject := "" + html.EscapeString(r.URL.Query().Get("withProject"))
	currentProject, version, requestSession := projectHandler.retrieveProjectAndVersion2(r)

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProjectVersion.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadVersions))
	}

	if withProject == "true" {
		docDep, docMissing, custDep, custMissing := projectHandler.getDeps(requestSession, currentProject)
		projectDto := currentProject.ToDto(docDep, docMissing, custDep, custMissing, hasDummyLabel(currentProject, getDummyLabel(requestSession, projectHandler.LabelRepository)))
		projectDto.AccessRights = *rights

		if currentProject.HasParent() {
			parentProject := projectHandler.ProjectRepository.FindByKey(requestSession, currentProject.Parent, true)
			if parentProject != nil {
				projectDto.ParentProjectSettings = parentProject.ToProjectSettingsDto(projectHandler.getDeps(requestSession, parentProject))
			}
		}

		render.JSON(w, r, projectDto)
	} else {
		render.JSON(w, r, version.ToDto())
	}
}

func (projectHandler *ProjectHandler) getDeps(rs *logy.RequestSession, pr *project.Project) (*department.Department, bool, *department.Department, bool) {
	var (
		custDep     *department.Department
		custMissing bool
		docDep      *department.Department
		docMissing  bool
	)
	if pr.CustomerMeta.DeptId != "" {
		custDep = projectHandler.DeparmentRepository.GetByDeptId(rs, pr.CustomerMeta.DeptId)
		custMissing = (custDep == nil)
	}
	if pr.DocumentMeta.SupplierDeptId != "" {
		docDep = projectHandler.DeparmentRepository.GetByDeptId(rs, pr.DocumentMeta.SupplierDeptId)
		docMissing = (docDep == nil)
	}
	return docDep, docMissing, custDep, custMissing
}

func (projectHandler *ProjectHandler) ProjectVersionUpdateHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, version, requestSession := projectHandler.retrieveProjectAndVersion2(r)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}
	username, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProjectVersion.Update {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.UpdateVersion))
	}

	versionRequestDto := extractVersionRequestBody(r, true)

	if version.Name != versionRequestDto.Name && currentProject.FindVersionByName(versionRequestDto.Name) != nil {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.VersionNameInUse, versionRequestDto.Name))
	}

	oldVersion := project.ProjectVersion{}
	copier.Copy(&oldVersion, version)

	currentProject.UpdateVersion(version, versionRequestDto.Name, versionRequestDto.Description)

	if version.Key == currentProject.ApprovableSPDX.VersionKey {
		currentProject.ApprovableSPDX.VersionName = versionRequestDto.Name
	}

	projectHandler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, version.Key, username, message.VersionUpdated, cmp.Diff, version, &oldVersion)
	projectHandler.ProjectRepository.Update(requestSession, currentProject)

	docDep, docMissing, custDep, custMissing := projectHandler.getDeps(requestSession, currentProject)
	dto := currentProject.ToDto(docDep, docMissing, custDep, custMissing, hasDummyLabel(currentProject, getDummyLabel(requestSession, projectHandler.LabelRepository)))

	// enrich with access rights
	dto.AccessRights = *rights

	if currentProject.HasParent() {
		parentProject := projectHandler.ProjectRepository.FindByKey(requestSession, currentProject.Parent, true)
		if parentProject != nil {
			dto.ParentProjectSettings = parentProject.ToProjectSettingsDto(projectHandler.getDeps(requestSession, currentProject))
		}
	}

	render.JSON(w, r, dto)
}

func (projectHandler *ProjectHandler) ProjectVersionGetUsageInApprovalOrReviewRequest(w http.ResponseWriter, r *http.Request) {
	currentProject, version, requestSession := projectHandler.retrieveProjectAndVersion2(r)

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProjectVersion.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadVersions))
	}

	isInUsage := projectHandler.isProjectOrVersionInApprovalOrContainsSbomToRetain(requestSession, currentProject, version)
	render.JSON(w, r, SuccessResponse{
		Success: isInUsage,
		Message: "Project Version usage in Approval or Review Request",
	})
}

// ProjectVersionDeleteExternHandler godoc
//
//	@Summary	Delete version (also known as channel) of project
//	@Id			deleteProjectVersion
//	@Produce	json
//	@Param		uuid	path		string					true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Param		version	path		string					true	"Project Version Name (also known as Channel Name) e.g.: main"
//	@Success	200		{object}	rest.SuccessResponse	"Success Response"
//	@Failure	417		{object}	exception.HttpError		"Can not be deleted, it is in use"
//	@Failure	404		{object}	exception.HttpError404	"NotFound Error"
//	@Failure	401		{object}	exception.HttpError		"Unauthorized Error"
//	@Router		/projects/{uuid}/versions/{version} [delete]
//	@security	Bearer
func (projectHandler *ProjectHandler) ProjectVersionDeleteExternHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	currentProject, version, _ := projectHandler.retrieveProjectAndVersionFromPublicRequest(requestSession, r)

	ip := jwt.TrimPortFromRemoteAddress(r.RemoteAddr)

	projectHandler.handleProjectVersionDelete(w, requestSession, currentProject, version, ip)
}

func (projectHandler *ProjectHandler) ProjectVersionDeleteHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, version, requestSession := projectHandler.retrieveProjectAndVersion2(r)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}
	username, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProjectVersion.Delete {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeleteVersion))
	}

	projectHandler.handleProjectVersionDelete(w, requestSession, currentProject, version, username)
}

func (projectHandler *ProjectHandler) handleProjectVersionDelete(w http.ResponseWriter, requestSession *logy.RequestSession, currentProject *project.Project, version *project.ProjectVersion, username string) {
	isInUsage := projectHandler.isProjectOrVersionInApprovalOrContainsSbomToRetain(requestSession, currentProject, version)
	if isInUsage {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorInUse, "Version"))
	}

	var oldVersions map[string]*project.ProjectVersion
	copier.Copy(&oldVersions, currentProject.Versions)

	oldVersion := project.ProjectVersion{}
	copier.Copy(&oldVersion, version)

	if !hasDummyLabel(currentProject, getDummyLabel(requestSession, projectHandler.LabelRepository)) {
		observermngmt.FireEvent(observermngmt.ProjectVersionDeleted, observermngmt.VersionData{
			RequestSession: requestSession,
			Project:        currentProject,
			Version:        version,
		})
	}

	currentProject.DeleteVersion(version.Key)
	projectHandler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, currentProject.Key, username, message.ProjectVersionDeleted, cmp.Diff, currentProject.Versions, oldVersions)
	projectHandler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, version.Key, username, message.VersionDeleted, cmp.Diff, version, &oldVersion)
	projectHandler.ProjectRepository.Update(requestSession, currentProject)

	w.WriteHeader(200)
}

func (projectHandler *ProjectHandler) isProjectOrVersionInApprovalOrContainsSbomToRetain(requestSession *logy.RequestSession, currentProject *project.Project, version *project.ProjectVersion) bool {
	// 1. Check if the project (or version) is referenced in any ApprovalList.
	if projectHandler.IsReferencedInApprovalLists(requestSession, currentProject, version) {
		return true
	}

	// 2. If a specific version is provided, check only its retained SBOM status.
	if version != nil {
		return projectHandler.SbomRetainedService.CheckIfRetainedSbom(requestSession, version, currentProject)
	}

	// 3. For a project- (or group-) level deletion (version is nil), check each version (or each child project’s version) for a retained SBOM.
	if projectHandler.SbomRetainedService.HasAnyVersionWithRetainedSbom(requestSession, currentProject) {
		return true
	}

	return false
}

func (projectHandler *ProjectHandler) IsReferencedInApprovalLists(requestSession *logy.RequestSession, currentProject *project.Project, version *project.ProjectVersion) bool {
	if currentProject.HasParent() {
		// Check the group's ApprovalList.
		if groupApprovalList := projectHandler.ApprovalListRepository.FindByKey(requestSession, currentProject.Parent, false); groupApprovalList != nil {
			if isReferencedInApprovalList(groupApprovalList, currentProject, version) {
				return true
			}
		}
		// Also check the project's own ApprovalList.
		// It could be a project, which was later added to the group and therefore, even if in the group, it could have own ApprovalList
		if approvalList := projectHandler.ApprovalListRepository.FindByKey(requestSession, currentProject.Key, false); approvalList != nil {
			// If no version is provided, the existence of an ApprovalList is enough.
			if version == nil || isReferencedInApprovalList(approvalList, currentProject, version) {
				return true
			}
		}
	} else {
		// For loose projects or groups, check only the project's own ApprovalList.
		if approvalList := projectHandler.ApprovalListRepository.FindByKey(requestSession, currentProject.Key, false); approvalList != nil {
			if version == nil || isReferencedInApprovalList(approvalList, currentProject, version) {
				return true
			}
		}
	}
	return false
}

func isReferencedInApprovalList(approvalList *approval.ApprovalList, currentProject *project.Project, version *project.ProjectVersion) bool {
	for _, a := range approvalList.Approvals {
		for _, pa := range a.Info.Projects {
			if pa.ProjectKey == currentProject.Key {
				// If no version is provided, a match on the project is sufficient.
				// Otherwise, the version keys must also match.
				if version == nil || pa.ApprovableSPDX.VersionKey == version.Key {
					return true
				}
			}
		}
	}
	return false
}

func (projectHandler *ProjectHandler) DownloadDocumentByTaskHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, true)

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowDisclosureDocument.Download {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorPermissionDeniedDownload, "Disclosure Document"))
	}

	taskGuid := rest.GetURLParam(r, "taskId")
	fileTypeStr := rest.GetURLParam(r, "fileType")

	documentVersionIndexStr := rest.GetURLParam(r, "docVersion")
	index, err := strconv.Atoi(documentVersionIndexStr)
	if err != nil {
		index = int(pdocument.NONE_VERSION)
	}
	langTag := extractLangTag(r)
	targetFileName := pdocument.GetFileNameWithIndex(pdocument.PDocumentType(fileTypeStr), taskGuid, langTag, int(index))
	var completeFileNameInS3 string

	var document *pdocument.PDocument
	completeFileNameInS3 = currentProject.GetFilePathDocumentForProject(targetFileName)
	document = currentProject.GetDocumentByFileNameWithIndex(targetFileName, index)

	if document == nil {
		exception.ThrowExceptionServer404Message(message.GetI18N(message.ErrorDbNotFound),
			"document not found, project="+currentProject.Key+", document="+completeFileNameInS3)
	}

	s3Helper.PerformDownload(requestSession, &w, completeFileNameInS3, document.Hash)
}

func extractLangTag(r *http.Request) *language.Tag {
	var langTag *language.Tag
	langStr, err := language.Parse(rest.GetURLParam(r, "lang"))
	if err != nil {
		// ignore
		langTag = nil
	} else {
		langTag = &langStr
	}
	return langTag
}

func (projectHandler *ProjectHandler) GetAllExternalSourcesHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, version, requestSession := projectHandler.retrieveProjectAndVersion2(r)

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)

	if !rights.AllowCCSAction.Download {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DownloadExternalSources))
	}

	result := make([]project.SourceExternalDTO, 0)
	for _, source := range version.SourceExternal {
		result = append(result, source.ToDTO())
	}

	render.JSON(w, r, result)
}

func (projectHandler *ProjectHandler) ExternalSourceCreateHandler(w http.ResponseWriter, r *http.Request) {
	projectData, version, requestSession := projectHandler.retrieveProjectAndVersion2(r)
	if projectData.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	userName, rights := roles.GetAndCheckProjectRights(requestSession, r, projectData, false)

	if !rights.AllowCCSAction.Upload {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.UploadExternalSources))
	}

	projectHandler.HandleExternalSourceCreate(requestSession, projectData, version, project.OriginUi, userName, w, r, true)
}

func (projectHandler *ProjectHandler) HandleExternalSourceCreate(requestSession *logy.RequestSession, projectData *project.Project,
	version *project.ProjectVersion, origin string, uploader string, w http.ResponseWriter, r *http.Request, handleErrorAsServerException bool,
) {
	sourceData := extractExternalSourceRequestBody(r, handleErrorAsServerException)

	ent := sourceData.ToEntity()
	ent.Origin = origin
	ent.Uploader = uploader
	version.SourceExternal = append(version.GetSourceExternalAll(), &ent)

	sourceExternalAudit := ent.ToAudit()
	projectHandler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, version.Key, uploader, message.SourceCodeResourceCreated, audit.DiffWithReporter, sourceExternalAudit, "")

	projectHandler.ProjectRepository.Update(requestSession, projectData)

	responseData := SuccessResponse{
		Success: true,
		Message: "External source added",
	}
	render.JSON(w, r, responseData)
}

func (projectHandler *ProjectHandler) ExternalSourceDeleteHandler(w http.ResponseWriter, r *http.Request) {
	projectData, version, requestSession := projectHandler.retrieveProjectAndVersion2(r)
	if projectData.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	userName, rights := roles.GetAndCheckProjectRights(requestSession, r, projectData, false)
	if !rights.AllowProjectVersion.Update {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeleteExternalSources))
	}

	idEscaped := chi.URLParam(r, "sourceId")
	id, err := url.QueryUnescape(idEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamSourceidEmpty))
	if id == "" {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ParamSourceidEmpty))
	}

	target := -1
	for i, source := range version.SourceExternal {
		if source.Key == id {
			target = i
			break
		}
	}

	if target == -1 {
		responseData := SuccessResponse{
			Success: false,
			Message: "ExternalSource not found",
		}
		render.JSON(w, r, responseData)
		return
	}

	sourceExternalAudit := version.GetSourceExternal(target).ToAudit()
	version.SourceExternal = removeExternalSource(version.SourceExternal, target)

	projectHandler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, version.Key, userName, message.SourceCodeResourceDeleted, audit.DiffWithReporter, "", sourceExternalAudit)

	projectHandler.ProjectRepository.Update(requestSession, projectData)

	responseData := SuccessResponse{
		Success: true,
		Message: "External source added",
	}
	render.JSON(w, r, responseData)
}

func (projectHandler *ProjectHandler) ExternalSourceUpdateHandler(w http.ResponseWriter, r *http.Request) {
	projectData, version, requestSession := projectHandler.retrieveProjectAndVersion2(r)
	if projectData.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	userName, rights := roles.GetAndCheckProjectRights(requestSession, r, projectData, false)
	if !rights.AllowCCSAction.Upload {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.UpdateExternalSource))
	}

	sourceData := extractExternalSourceRequestBody(r, true)

	idEscaped := chi.URLParam(r, "sourceId")
	id, err := url.QueryUnescape(idEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamSourceidEmpty))
	if id == "" {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ParamSourceidEmpty))
	}

	target := -1
	for i, source := range version.SourceExternal {
		if source.Key == id {
			target = i
			break
		}
	}

	if target == -1 {
		responseData := SuccessResponse{
			Success: false,
			Message: "ExternalSource not found",
		}
		render.JSON(w, r, responseData)
		return
	}
	sourceExternal := version.GetSourceExternal(target)
	oldSourceExternalAudit := sourceExternal.ToAudit()
	sourceExternal.Update(sourceData, project.OriginUi, userName)
	newSourceExternalAudit := sourceExternal.ToAudit()

	projectHandler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, version.Key, userName, message.SourceCodeResourceUpdated, audit.DiffWithReporter, newSourceExternalAudit, oldSourceExternalAudit)

	projectHandler.ProjectRepository.Update(requestSession, projectData)

	responseData := SuccessResponse{
		Success: true,
		Message: "External source updated",
	}
	render.JSON(w, r, responseData)
}

func removeExternalSource(slice []*project.SourceExternal, s int) []*project.SourceExternal {
	return append(slice[:s], slice[s+1:]...)
}

func isFilenameValid(filename string) bool {
	return !strings.ContainsAny(filename, conf.Config.Server.DisallowedUploadFilenameChars)
}

func (projectHandler *ProjectHandler) ComponentDetailsForSbomGetHandler(writer http.ResponseWriter, request *http.Request) {
	sbomUuidEscaped := chi.URLParam(request, "sbomUuid")
	sbomUuid, err := url.QueryUnescape(sbomUuidEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamSbomUuidEmpty))
	if sbomUuid == "" {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ParamSbomUuidEmpty))
	}
	err = validation.CheckUuid(sbomUuid)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorKeyRequestParamNotValid, "sbomUuid"), zapcore.InfoLevel)
	spdxIdEscaped := chi.URLParam(request, "spdxId")
	spdxId, err := url.QueryUnescape(spdxIdEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamSpdxidEmpty))
	if spdxId == "" {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ParamSpdxidEmpty))
	}

	currentProject, version, requestSession := projectHandler.retrieveProjectAndVersion2(request)
	userID, rights := roles.GetAndCheckProjectRights(requestSession, request, currentProject, false)
	if !rights.AllowProjectVersion.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadComponentDetails))
	}
	_, spdxFile := projectHandler.RetrieveSbomListAndFile(requestSession, version.Key, sbomUuid)
	if spdxFile == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.FindingSbomKey), "SPDX not found in history: "+sbomUuid)
	}
	sbomContent := s3Helper.ReadTextFile(requestSession, currentProject.GetFilePathSbom(spdxFile.Key, version.Key), spdxFile.Hash)

	holder := projectService.RepositoryHolder{
		LicenseRepository:      projectHandler.LicenseRepository,
		PolicyRulesRepository:  projectHandler.PolicyRuleRepository,
		LicenseRulesRepository: projectHandler.LicenseRulesRepository,
	}

	isProjectResponsible := currentProject.IsResponsible(userID)

	if !isProjectResponsible {
		isProjectResponsible = rights.IsFossOffice() && projectHandler.ProjectLabelService.HasVehiclePlatformLabel(requestSession, currentProject)
	}

	res := detailsService.GetComponentDetails(requestSession, holder, currentProject.Key, spdxId, *sbomContent, spdxFile.Uploaded, spdxFile.Key, isProjectResponsible)
	render.JSON(writer, request, res)
}

func (projectHandler *ProjectHandler) ComponentLicensesGetHandler(writer http.ResponseWriter, request *http.Request) {
	sbomUuidEscaped := chi.URLParam(request, "sbomUuid")
	sbomUuid, err := url.QueryUnescape(sbomUuidEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamSbomUuidEmpty))
	if sbomUuid == "" {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ParamSbomUuidEmpty))
	}
	err = validation.CheckUuid(sbomUuid)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorKeyRequestParamNotValid, "sbomUuid"), zapcore.InfoLevel)

	spdxIdEscaped := chi.URLParam(request, "spdxId")
	spdxId, err := url.QueryUnescape(spdxIdEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamSpdxidEmpty))
	if spdxId == "" {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ParamSpdxidEmpty))
	}

	currentProject, version, requestSession := projectHandler.retrieveProjectAndVersion2(request)
	_, rights := roles.GetAndCheckProjectRights(requestSession, request, currentProject, false)
	if !rights.AllowProjectVersion.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadComponentDetails))
	}
	_, spdxFile := projectHandler.RetrieveSbomListAndFile(requestSession, version.Key, sbomUuid)
	if spdxFile == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.FindingSbomKey), "SPDX not found in history: "+sbomUuid)
	}
	sbomContent := s3Helper.ReadTextFile(requestSession, currentProject.GetFilePathSbom(spdxFile.Key, version.Key), spdxFile.Hash)

	holder := projectService.RepositoryHolder{
		LicenseRepository:      projectHandler.LicenseRepository,
		PolicyRulesRepository:  projectHandler.PolicyRuleRepository,
		LicenseRulesRepository: projectHandler.LicenseRulesRepository,
	}

	response := detailsService.GetComponentLicenses(requestSession, holder, currentProject.Key, spdxId, *sbomContent, spdxFile.Uploaded, spdxFile.Key)

	render.JSON(writer, request, response)
}

func (projectHandler *ProjectHandler) ComponentReviewRemarksGetHandler(w http.ResponseWriter, r *http.Request) {
	spdxIdEscaped := chi.URLParam(r, "spdxId")
	spdxId, err := url.QueryUnescape(spdxIdEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamSpdxidEmpty))
	if spdxId == "" {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ParamSpdxidEmpty))
	}

	currentProject, version, requestSession := projectHandler.retrieveProjectAndVersion2(r)
	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProjectVersion.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadVersionReviewRemarks))
	}

	remarksList := projectHandler.ReviewRemarksRepository.FindByKeyFilteredByComponentId(requestSession, version.Key, spdxId)

	if remarksList == nil {
		render.JSON(w, r, []reviewremarks.RemarkDto{})
		return
	}

	sort.SliceStable(remarksList.Remarks, func(i, j int) bool {
		ri, rj := remarksList.Remarks[i], remarksList.Remarks[j]

		isOpenI := ri.Status == reviewremarks.Open || ri.Status == reviewremarks.InProgress
		isOpenJ := rj.Status == reviewremarks.Open || rj.Status == reviewremarks.InProgress
		if isOpenI != isOpenJ {
			return isOpenI
		}

		levelWeight := map[reviewremarks.Level]int{
			reviewremarks.Red:    3,
			reviewremarks.Yellow: 2,
			reviewremarks.Green:  1,
		}
		return levelWeight[ri.Level] > levelWeight[rj.Level]
	})
	render.JSON(w, r, domain.ToDtos(remarksList.Remarks))
}

func (projectHandler *ProjectHandler) SbomAllLicensesGetHandler(writer http.ResponseWriter, request *http.Request) {
	sbomUuidEscaped := chi.URLParam(request, "sbomUuid")
	sbomUuid, err := url.QueryUnescape(sbomUuidEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamSbomUuidEmpty))
	if sbomUuid == "" {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ParamSbomUuidEmpty))
	}
	err = validation.CheckUuid(sbomUuid)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorKeyRequestParamNotValid, "sbomUuid"), zapcore.InfoLevel)

	currentProject, version, requestSession := projectHandler.retrieveProjectAndVersion2(request)
	_, rights := roles.GetAndCheckProjectRights(requestSession, request, currentProject, false)
	if !rights.AllowProjectVersion.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadComponentDetails))
	}

	_, spdxFile := projectHandler.RetrieveSbomListAndFile(requestSession, version.Key, sbomUuid)
	if spdxFile == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.FindingSbomKey), "SPDX not found in history: "+sbomUuid)
	}
	comps := projectHandler.SpdxService.GetComponentInfos(requestSession, currentProject, version.Key, spdxFile)

	qc := database.New().SetMatcher(
		database.AttributeMatcher(
			"Deleted",
			database.EQ,
			false,
		),
	).SetKeep([]string{"licenseId", "name"})
	lics := projectHandler.LicenseRepository.Query(requestSession, qc)
	licNames := make(map[string]string)
	for _, l := range lics {
		licNames[l.LicenseId] = l.Name
	}

	var result project.SbomLicensesDto
	seenKnown := make(map[string]struct{})
	seenUnknown := make(map[string]struct{})
	for _, comp := range comps {
		unknown, known := detailsService.ProcessCompLicenses(requestSession, &comp, licNames)
		for _, lic := range known {
			id := lic.Id
			if _, exists := seenKnown[id]; !exists {
				seenKnown[id] = struct{}{}
				result.Known = append(result.Known, lic)
			}
		}
		for _, name := range unknown {
			if _, exists := seenUnknown[name]; !exists {
				seenUnknown[name] = struct{}{}
				result.Unknown = append(result.Unknown, name)
			}
		}
	}
	render.JSON(writer, request, result)
}

func (projectHandler *ProjectHandler) ProjectVersionScanRemarksForSbom(w http.ResponseWriter, r *http.Request) {
	result := projectHandler.retrieveProjectAndVersionAndCreateScanRemarksForSbom(r)

	render.JSON(w, r, result)
}

func (projectHandler *ProjectHandler) DownloadScanRemarksForSbomCsvHandler(w http.ResponseWriter, r *http.Request) {
	remarks := projectHandler.retrieveProjectAndVersionAndCreateScanRemarksForSbom(r)
	writeScanRemarksAsCsvIntoResponse(&w, remarks)
	w.WriteHeader(http.StatusOK)
}

func (projectHandler *ProjectHandler) DownloadLicenseRemarksForSbomCsvHandler(w http.ResponseWriter, r *http.Request) {
	evalRes, requestSession := projectHandler.preparePolicyRulesEvaluationResultForLicenseRemarksForSbom(r)
	remarks := projectHandler.createQualityLicenseRemarks(requestSession, evalRes)
	if remarks == nil {
		return
	}
	writeLicenseRemarksAsCsvIntoResponse(&w, remarks)
	w.WriteHeader(http.StatusOK)
}

func (projectHandler *ProjectHandler) retrieveProjectAndVersionAndCreateScanRemarksForSbom(r *http.Request) (res []project.QualityScanRemarks) {
	sbomUuidEscaped := chi.URLParam(r, "sbomUuid")
	sbomUuid, err := url.QueryUnescape(sbomUuidEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamSbomUuidEmpty))
	if sbomUuid == "" {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ParamSbomUuidEmpty))
	}
	err = validation.CheckUuid(sbomUuid)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorKeyRequestParamNotValid, "sbomUuid"), zapcore.InfoLevel)

	res = make([]project.QualityScanRemarks, 0)
	currentProject, version, requestSession := projectHandler.retrieveProjectAndVersion2(r)

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProjectVersion.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadVersionScanRemarks))
	}

	_, selectedSpdx := projectHandler.RetrieveSbomListAndFile(requestSession, version.Key, sbomUuid)
	if selectedSpdx == nil {
		return
	}
	compInfos := projectHandler.SpdxService.GetComponentInfos(requestSession, currentProject, version.Key, selectedSpdx)
	rules := projectHandler.PolicyRuleRepository.FindPolicyRulesForLabel(requestSession, currentProject.PolicyLabels)
	policyDecisions := projectHandler.PolicyDecisionsRepository.FindByKey(requestSession, currentProject.Key, false)
	isVehicle := projectHandler.ProjectLabelService.HasVehiclePlatformLabel(requestSession, currentProject)
	evalRes := compInfos.EvaluatePolicyRules(rules, policyDecisions, isVehicle, selectedSpdx.Uploaded, selectedSpdx.Key)

	// projectHandler.ExtractPolicyRuleStatus(requestSession, currentProject, version)
	return projectHandler.createQualityScanRemarks(requestSession, currentProject, &selectedSpdx.MetaInfo, evalRes, projectHandler.LabelRepository)
}

func writeScanRemarksAsCsvIntoResponse(w *http.ResponseWriter, remarks []project.QualityScanRemarks) {
	csvWriter := csv.NewWriter(*w)
	csvWriter.Comma = ';'
	defer csvWriter.Flush()

	csvHeader := []string{"STATUS", "REMARK", "COMPONENT_NAME", "COMPONENT_VERSION", "DESCRIPTION"}

	if err := csvWriter.Write(csvHeader); err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCsvGeneration, "scan remark", "header"), err)
	}

	for _, remark := range remarks {
		version := remark.Version
		if len(remark.Version) > 0 {
			version = " " + remark.Version
		}
		csvRow := []string{string(remark.Status), message.GetI18N(remark.RemarkKey).Text, remark.Name, version, message.GetI18N(remark.DescriptionKey).Text}
		if err := csvWriter.Write(csvRow); err != nil {
			exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCsvGeneration, "scan remark", "data"), err)
		}
	}
}

func writeLicenseRemarksAsCsvIntoResponse(w *http.ResponseWriter, remarks []project.QualityLicenseRemarks) {
	csvWriter := csv.NewWriter(*w)
	csvWriter.Comma = ';'
	defer csvWriter.Flush()

	csvHeader := []string{"STATUS", "TYPE", "REMARK", "LICENSE", "COMPONENT_NAME", "COMPONENT_VERSION", "DESCRIPTION"}

	if err := csvWriter.Write(csvHeader); err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCsvGeneration, "license remark", "header"), err)
	}

	for _, remark := range remarks {
		version := remark.Version
		if len(remark.Version) > 0 {
			version = " " + remark.Version
		}

		// compile RegExp parser for one or more new line characters
		re := regexp.MustCompile(`\n+`)
		description := strings.TrimSpace(re.ReplaceAllString(remark.Description, " "))

		// Check if starts with one of the sprecial character and add a leading space
		reSpecial := regexp.MustCompile(`^[=+\-@"]`)
		if reSpecial.MatchString(description) {
			description = " " + description
		}

		csvRow := []string{string(remark.Status), remark.Type, remark.Remark, remark.License, remark.Name, version, description}
		if err := csvWriter.Write(csvRow); err != nil {
			exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCsvGeneration, "license remark", "data"), err)
		}
	}
}

func (projectHandler *ProjectHandler) createQualityScanRemarks(rs *logy.RequestSession, p *project.Project, meta *project.MetaInfo, evalRes *components.EvaluationResult, labelRepository labels.ILabelRepository) []project.QualityScanRemarks {
	result := make([]project.QualityScanRemarks, 0)

	if meta.HasExternalRefs && len(evalRes.Results) > 0 {
		result = appendQualityScanRemarks(result, components.ComponentResult{Component: &components.ComponentInfo{}}, project.PROBLEM,
			message.SrContentException, message.SrContainsExternalRefs)
	}

	contactAddress := p.NoticeContactMeta.Address
	if p.HasParent() {
		parentProject := projectHandler.ProjectRepository.FindByKey(rs, p.Parent, false)
		if parentProject != nil {
			contactAddress = parentProject.NoticeContactMeta.Address
		}
	}
	if contactAddress == "" && !hasOnboardLabel(rs, p, labelRepository) {
		result = appendQualityScanRemarks(result, components.ComponentResult{Component: &components.ComponentInfo{}}, project.WARNING,
			message.SrProjectAddressException, message.SrProjectAddressDescription)
	}

	hasNotAllowedUnicodeLetters := func(str string) bool {
		allowedLetters := "<>'@+-©[]"
		for _, letter := range str {
			isAllowed := strings.IndexRune(allowedLetters, letter) > -1
			if isAllowed {
				// is allowed unicode, check next letter
				continue
			}
			if unicode.IsSymbol(letter) {
				return true
			}
		}
		return false
	}

	isMalformedCopyright := func(str string) bool {
		badText := []string{
			"false",
			"copyright:before",
			"null",
			"key",
			"context",
		}
		for _, b := range badText {
			if strings.Contains(str, b) {
				return true
			}
		}
		return false
	}

	hasNonLatinOrNotAllowedLetters := func(str string) bool {
		allowedLetters := "#@:/.-_%=?~&"
		for _, letter := range str {
			isAllowed := strings.IndexRune(allowedLetters, letter) > -1 || unicode.IsDigit(letter)
			if isAllowed {
				continue
			}
			if (letter < 'a' || letter > 'z') && (letter < 'A' || letter > 'Z') {
				return true
			}
		}
		return false
	}

	var projectLevel string
	for _, compRes := range evalRes.Results {
		if compRes.Component.HasAnnotations {
			result = appendQualityScanRemarks(result, compRes, project.INFORMATION,
				message.SrContentException, message.SrContainsAnnotations)
		}
		if compRes.Component.Type == components.SNIPPET {
			result = appendQualityScanRemarks(result, compRes, project.PROBLEM,
				message.SrContentException, message.SrContainsSnippet)
		}
		if len(compRes.Component.Version) == 0 || helper.IsUnasserted(compRes.Component.Version) {
			result = appendQualityScanRemarks(result, compRes, project.PROBLEM,
				message.SrMissingVersion, message.SrMissingVersionDescription)
		}
		if compRes.Component.ComplexExpression {
			result = appendQualityScanRemarks(result, compRes, project.PROBLEM,
				message.SrContentException, message.SrContainsComplex)
		} else if compRes.Component.LicensesDeclared.CountOrLinks() >= OrLinksThreshold ||
			compRes.Component.LicensesConcluded.CountOrLinks() >= OrLinksThreshold {
			result = appendQualityScanRemarks(result, compRes, project.WARNING,
				message.SrTooMuchOrTitle, message.SrContainsTooMuchOr)
		}
		if compRes.Component.ContainsBadChars {
			result = appendQualityScanRemarks(result, compRes, project.PROBLEM,
				message.SrContentException, message.SrContainsBadchars)
		}
		if len(compRes.Component.Name) == 0 || helper.IsUnasserted(compRes.Component.Name) {
			result = appendQualityScanRemarks(result, compRes, project.PROBLEM,
				message.SrMissingName, message.SrMissingNameDescription)
		}

		if len(compRes.Component.GetLicenseEffective()) == 0 || helper.IsUnasserted(compRes.Component.GetLicenseEffective()) {
			result = appendQualityScanRemarks(result, compRes, project.PROBLEM,
				message.SrMissingLicenseId, message.SrMissingLicenseIdDescription)
		}

		// check CopyrightText - exist
		if len(compRes.Component.CopyrightText) == 0 || helper.IsUnasserted(compRes.Component.CopyrightText) {
			if projectLevel == "" {
				projectLevel = copyrightMissingLevel(rs, p, projectHandler.LabelRepository)
			}
			result = appendQualityScanRemarks(result, compRes, projectLevel,
				message.SrMissingCopyrightText, message.SrMissingCopyrightDescription)
		} else {
			// check CopyrightText - content
			if hasNotAllowedUnicodeLetters(compRes.Component.CopyrightText) ||
				// year template
				strings.Contains(strings.ToLower(compRes.Component.CopyrightText), "yyyy") ||
				isMalformedCopyright(compRes.Component.CopyrightText) {
				result = appendQualityScanRemarks(result, compRes, project.INFORMATION,
					message.SrMalformedCopyrightText, message.SrMalformedCopyrightDescription)
			}

			// check CopyrightText - length
			if len(compRes.Component.CopyrightText) > 1000 {
				result = appendQualityScanRemarks(result, compRes, project.INFORMATION,
					message.SrCopyrightLongText, message.SrCopyrightToLongDescription)
			}

		}

		if !helper.IsUnasserted(compRes.Component.License) && !helper.IsUnasserted(compRes.Component.LicenseDeclared) &&
			compRes.Component.License != compRes.Component.LicenseDeclared {
			result = appendQualityScanRemarks(result, compRes, project.WARNING,
				message.SrLicensesDiff, message.SrLicensesDiffDescription)
		}

		for _, l := range compRes.Component.GetLicensesEffective().List {
			if !l.Known && !helper.IsUnasserted(l.OrigName) {
				result = appendQualityScanRemarks(result, compRes, project.PROBLEM,
					message.SrUnknownLicenseUsed, message.SrUnknownLicenseUsedDescription)
			}
			if compRes.ContainsUnmatchedLicense(l.OrigName) {
				result = appendQualityScanRemarks(result, compRes, project.PROBLEM,
					message.SrUnmatchedLicenseUsed, message.SrUnmatchedLicenseUsedDescription)
			}
			if !strings.EqualFold(l.ReferencedLicense, l.OrigName) && l.Known {
				result = appendQualityScanRemarks(result, compRes, project.INFORMATION,
					message.SrAliasingUsed, message.SrAliasingUsedDescription)
			}
		}

		if len(compRes.Component.PURL) > 0 && compRes.Component.PURL != "NOASSERTION" && hasNonLatinOrNotAllowedLetters(compRes.Component.PURL) {
			result = appendQualityScanRemarks(result, compRes, project.WARNING,
				message.SrContainsNonLatinLetters, message.SrContainsNonLatinLettersDescription)
		}

	}

	return result
}

func hasOnboardLabel(requestSession *logy.RequestSession, currentProject *project.Project, labelRepository labels.ILabelRepository) bool {
	onboardLabel := labelRepository.FindByNameAndType(requestSession, label.ONBOARD, label.POLICY)
	if onboardLabel != nil {
		return slices.Contains(currentProject.PolicyLabels, onboardLabel.GetKey())
	} else {
		return false
	}
}

func copyrightMissingLevel(rs *logy.RequestSession, pr *project.Project, labelRepo labels.ILabelRepository) string {
	ll := make([]string, 0)
	for _, l := range pr.PolicyLabels {
		resolvedLabel := labelRepo.FindByKey(rs, l, false)
		if resolvedLabel == nil {
			continue
		}
		ll = append(ll, resolvedLabel.Name)
	}

	switch {
	case helper.Contains("enterprise platform", ll) && helper.Contains("entity users", ll):
		return project.INFORMATION
	case helper.Contains("enterprise platform", ll) && helper.Contains("group users", ll):
		return project.INFORMATION
	case helper.Contains("vehicle platform", ll) && helper.Contains("onboard", ll):
		return project.PROBLEM
	case helper.Contains("external users", ll):
		return project.PROBLEM
	case helper.Contains("other platform", ll):
		return project.PROBLEM
	}

	return project.WARNING
}

func canBeOmitted(compRes components.ComponentResult, remarkKey string) bool {
	if compRes.Component.Type != components.ROOT || len(compRes.Component.GetLicensesEffective().List) > 0 {
		return false
	}
	omitRemarks := []string{
		message.SrMissingLicenseText,
		message.SrMissingLicenseId,
		message.SrMissingCopyrightText,
		message.SrMissingVersion,
	}
	return helper.Contains(remarkKey, omitRemarks)
}

func appendQualityScanRemarks(result []project.QualityScanRemarks,
	compRes components.ComponentResult, problem string, remarkKey string, descriptionKey string,
) []project.QualityScanRemarks {
	if canBeOmitted(compRes, remarkKey) {
		return result
	}
	componentType := project.ComponentType(compRes.Component.Type)
	if descriptionKey == message.SrProjectAddressDescription {
		componentType = project.PROJECT
	}
	result = append(result, project.QualityScanRemarks{
		Status:            project.ScanRemarkStatus(problem),
		RemarkKey:         remarkKey,
		Name:              "" + compRes.Component.Name,
		SpdxId:            "" + compRes.Component.SpdxId,
		Version:           "" + compRes.Component.Version,
		Type:              componentType,
		DescriptionKey:    descriptionKey,
		PolicyRuleStatus:  components.ToPolicyStatusDto(compRes.Status, false),
		UnmatchedLicenses: components.ToUnmatchedDto(compRes.Unmatched),
	})
	return result
}

func (projectHandler *ProjectHandler) ProjectVersionLicenseRemarksForSbom(w http.ResponseWriter, r *http.Request) {
	evalRes, requestSession := projectHandler.preparePolicyRulesEvaluationResultForLicenseRemarksForSbom(r)
	result := licenseremarks.CreateQualityLicenseRemarks(requestSession, projectHandler.LicenseRepository, projectHandler.ObligationRepository, evalRes)
	if result == nil {
		return
	}
	render.JSON(w, r, result)
}

func (projectHandler *ProjectHandler) preparePolicyRulesEvaluationResultForLicenseRemarksForSbom(r *http.Request) (*components.EvaluationResult, *logy.RequestSession) {
	sbomUuidEscaped := chi.URLParam(r, "sbomUuid")
	sbomUuid, err := url.QueryUnescape(sbomUuidEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamSbomUuidEmpty))
	if sbomUuid == "" {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ParamSbomUuidEmpty))
	}
	err = validation.CheckUuid(sbomUuid)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorKeyRequestParamNotValid, "sbomUuid"), zapcore.InfoLevel)

	currentProject, version, requestSession := projectHandler.retrieveProjectAndVersion2(r)

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProjectVersion.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadVersionScanRemarks))
	}

	_, selectedSpdx := projectHandler.RetrieveSbomListAndFile(requestSession, version.Key, sbomUuid)
	if selectedSpdx == nil {
		return nil, requestSession
	}
	compInfos := projectHandler.SpdxService.GetComponentInfos(requestSession, currentProject, version.Key, selectedSpdx)
	rules := projectHandler.PolicyRuleRepository.FindPolicyRulesForLabel(requestSession, currentProject.PolicyLabels)
	policyDecisions := projectHandler.PolicyDecisionsRepository.FindByKey(requestSession, currentProject.Key, false)
	isVehicle := projectHandler.ProjectLabelService.HasVehiclePlatformLabel(requestSession, currentProject)
	evalRes := compInfos.EvaluatePolicyRules(rules, policyDecisions, isVehicle, selectedSpdx.Uploaded, selectedSpdx.Key)
	return evalRes, requestSession
}

func (projectHandler *ProjectHandler) createQualityLicenseRemarks(requestSession *logy.RequestSession, evalRes *components.EvaluationResult) []project.QualityLicenseRemarks {
	result := make([]project.QualityLicenseRemarks, 0)

	licensesMap := make(map[string]bool, 0)
	for _, compRes := range evalRes.Results {
		for _, license := range compRes.Component.GetLicensesEffective().List {
			licenseName := license.ReferencedLicense
			if helper.IsUnasserted(licenseName) {
				continue
			}
			licensesMap[licenseName] = true
		}
	}
	licenseIds := make([]string, 0)
	for id := range licensesMap {
		licenseIds = append(licenseIds, id)
	}
	licensesSlice := projectHandler.LicenseRepository.FindByIds(requestSession, licenseIds)
	licenses := make(map[string]*license.License)
	for _, license := range licensesSlice {
		licenses[license.LicenseId] = license
	}

	obligationsMap := make(map[string]bool, 0)
	for _, license := range licensesSlice {
		for _, obligationKey := range license.Meta.ObligationsKeyList {
			obligationsMap[obligationKey] = true
		}
	}
	obligationsKeys := make([]string, 0)
	for key := range obligationsMap {
		obligationsKeys = append(obligationsKeys, key)
	}
	obligationsSlice := projectHandler.ObligationRepository.FindByKeys(requestSession, obligationsKeys, false)
	obligations := make(map[string]*obligation.Obligation)
	for _, obligation := range obligationsSlice {
		obligations[obligation.Key] = obligation
	}

	for _, compRes := range evalRes.Results {
		for _, license := range compRes.Component.GetLicensesEffective().List {
			licenseName := license.ReferencedLicense
			if helper.IsUnasserted(licenseName) {
				continue
			}
			license, found := licenses[licenseName]
			if !found {
				continue
			}
			for _, obligationKey := range license.Meta.ObligationsKeyList {
				obligation, found := obligations[obligationKey]
				if !found {
					continue
				}

				result = append(result, project.QualityLicenseRemarks{
					Status:           string(obligation.WarnLevel),
					Remark:           obligation.Name,
					SpdxId:           compRes.Component.SpdxId,
					Type:             string(obligation.Type),
					Name:             "" + compRes.Component.Name,
					Version:          "" + compRes.Component.Version,
					License:          "" + licenseName,
					Description:      obligation.Description,
					PolicyRuleStatus: compRes.Status,
				})
			}
		}
	}
	return result
}

func (p *ProjectHandler) ProjectVersionSPDXHistory(w http.ResponseWriter, r *http.Request) {
	currentProject, version, requestSession := p.retrieveProjectAndVersion2(r)

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProjectVersion.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ViewSbom))
	}

	sbomList := p.SbomListRepository.FindByKey(requestSession, version.Key, false)
	if sbomList == nil {
		render.JSON(w, r, []interface{}{})
		return
	}
	spdxFileHistory := sbomList.SpdxFileHistory

	sort.Slice(spdxFileHistory, func(i, j int) bool {
		return spdxFileHistory[i].Uploaded.UTC().After(spdxFileHistory[j].Uploaded.UTC())
	})
	count := len(spdxFileHistory)
	if count > 0 {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err == nil && limit > 0 {
			if limit > count {
				limit = count
			}
			spdxFileHistory = spdxFileHistory[:limit]
		}

		approvableIncluded := false
		unusedSpdxCount := 0
		for i, spdxFile := range spdxFileHistory {
			if limit > 0 && len(currentProject.ApprovableSPDX.SpdxKey) > 0 && spdxFile.Key == currentProject.ApprovableSPDX.SpdxKey {
				approvableIncluded = true
			}

			if sbomlockRetained.IsSpdxToRetain(spdxFile, version) {
				spdxFileHistory[i].IsToRetain = true
			}
			if !IsSpdxInUse(spdxFile, currentProject, version) {
				if unusedSpdxCount < 5 {
					unusedSpdxCount++
				} else {
					spdxFileHistory[i].IsToDelete = true
				}
			} else {
				spdxFileHistory[i].IsInUse = true
			}
		}
		if !approvableIncluded && limit > 0 && len(currentProject.ApprovableSPDX.SpdxKey) > 0 {
			for _, spdxFile := range sbomList.SpdxFileHistory {
				if spdxFile.Key == currentProject.ApprovableSPDX.SpdxKey {
					spdxFileHistory = append(spdxFileHistory, spdxFile)
				}
			}
		}
	}

	render.JSON(w, r, spdxFileHistory)
}

func ValidateIDOrLatest(escaped string) string {
	sbomUuid, err := url.QueryUnescape(escaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamSbomUuidEmpty))
	if sbomUuid == "" {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ParamSbomUuidEmpty))
	}
	if sbomUuid == "latest" {
		return sbomUuid
	}
	err = validation.CheckUuid(sbomUuid)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorKeyRequestParamNotValid, "sbomUuid"), zapcore.InfoLevel)
	return sbomUuid
}

// ProjectVersionSPDXMetaByIDExtern godoc
//
//	@Summary	Get SBOM meta data for a specific delivery
//	@Id			getSBOMMetaForDelivery
//	@Produce	json
//	@Param		uuid		path		string							true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Param		version		path		string							true	"Project Version Name (also known as Channel Name) e.g.: main"
//	@Param		sbomUuid	path		string							true	"UUID of the SBOM delivery or 'latest' for the latest SBOM delivery e.g.: dummy-sbom-id---xxx-4413-yyy-24f060311111"
//	@Success	200			{object}	project.SPDXMetaPublicResponse	"SPDX Meta Data"
//	@Failure	404			{object}	exception.HttpError404			"NotFound Error"
//	@Failure	401			{object}	exception.HttpError				"Unauthorized Error"
//	@Router		/projects/{uuid}/versions/{version}/sboms/{sbomUuid} [get]
//	@security	Bearer
func (projectHandler *ProjectHandler) ProjectVersionSPDXMetaByIDExtern(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	currentProject, version, _ := projectHandler.retrieveProjectAndVersionFromPublicRequest(requestSession, r)

	sbomUuidEscaped := chi.URLParam(r, "sbomUuid")
	sbomUuid := ValidateIDOrLatest(sbomUuidEscaped)

	var spdxFile *project.SpdxFileBase
	if sbomUuid == "latest" {
		_, spdxFile = projectHandler.retrieveSbomListAndLatestFile(requestSessionTest, version.Key)
	} else {
		_, spdxFile = projectHandler.RetrieveSbomListAndFile(requestSessionTest, version.Key, sbomUuid)
	}
	if spdxFile == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.FindingSbomKey), "SPDX not found in history: "+sbomUuid)
	}

	responseData := project.SPDXMetaPublicResponse{
		Name:     spdxFile.MetaInfo.Name,
		Id:       spdxFile.MetaInfo.SpdxId,
		Version:  spdxFile.MetaInfo.SpdxVersion,
		Creators: strings.Join(spdxFile.MetaInfo.Creators, ","),
		Created:  spdxFile.Updated,
		Uploaded: spdxFile.Updated,
		Status:   true,
		IsRetain: spdxFile.IsToRetain || IsSpdxRetained(spdxFile, currentProject, version),
	}
	render.JSON(w, r, responseData)
}

// ProjectVersionSPDXHistoryExtern godoc
//
//	@Summary	Get SPDX list of project version (also known as channel)
//	@Id			getProjectVersionSPDXList
//	@Produce	json
//	@Param		uuid	path		string									true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Param		version	path		string									true	"Project Version Name (also known as Channel Name) e.g.: main"
//	@Success	200		{array}		project.VersionHistoryPublicResponse	"Version (also known as Channel) History"
//	@Failure	404		{object}	exception.HttpError404					"NotFound Error"
//	@Failure	401		{object}	exception.HttpError						"Unauthorized Error"
//	@Router		/projects/{uuid}/versions/{version}/sboms [get]
//	@security	Bearer
func (projectHandler *ProjectHandler) ProjectVersionSPDXHistoryExtern(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, version, _ := projectHandler.retrieveProjectAndVersionFromPublicRequest(requestSession, r)

	sbomList := projectHandler.SbomListRepository.FindByKey(requestSession, version.Key, false)
	if sbomList == nil {
		render.JSON(w, r, []interface{}{})
		return
	}

	responseData := make([]project.VersionHistoryPublicResponse, 0)
	for _, item := range sbomList.SpdxFileHistory {
		responseData = append(responseData, project.VersionHistoryPublicResponse{
			Name:    item.MetaInfo.Name,
			Updated: item.Updated,
			Id:      item.ChildEntity.Key,
		})
	}

	render.JSON(w, r, responseData)
}

func (projectHandler *ProjectHandler) ProjectVersionComponentsForSbom(w http.ResponseWriter, r *http.Request) {
	sbomUuidEscaped := chi.URLParam(r, "sbomUuid")
	sbomUuid, err := url.QueryUnescape(sbomUuidEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamSbomUuidEmpty))
	if sbomUuid == "" {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ParamSbomUuidEmpty))
	}
	err = validation.CheckUuid(sbomUuid)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorKeyRequestParamNotValid, "sbomUuid"), zapcore.InfoLevel)

	currentProject, version, requestSession := projectHandler.retrieveProjectAndVersion2(r)

	userID, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProjectVersion.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ViewComponents))
	}

	_, selectedSpdx := projectHandler.RetrieveSbomListAndFile(requestSession, version.Key, sbomUuid)
	if selectedSpdx == nil {
		render.JSON(w, r, components.ComponentsInfoResponse{})
		return
	}

	compInfos := projectHandler.SpdxService.GetComponentInfos(requestSession, currentProject, version.Key, selectedSpdx)
	rules := projectHandler.PolicyRuleRepository.FindPolicyRulesForLabel(requestSession, currentProject.PolicyLabels)

	policyDecisions := projectHandler.PolicyDecisionsRepository.FindByKey(requestSession, currentProject.Key, false)
	isVehicle := projectHandler.ProjectLabelService.HasVehiclePlatformLabel(requestSession, currentProject)

	evalRes := compInfos.EvaluatePolicyRules(rules, policyDecisions, isVehicle, selectedSpdx.Uploaded, selectedSpdx.Key)

	isResponsible := currentProject.IsResponsible(userID)
	policyDecisionDeniedReason := evaluatePolicyDecisionDeniedReason(isResponsible, rights.IsFossOffice(), isVehicle)
	isAllowDeniedPolicyDecision := evaluateIsAllowDeniedPolicyDecision(rights.IsDomainAdmin(), rights.IsFossOffice(), isVehicle)

	response := components.ComponentsInfoResponse{
		ComponentInfo:                  evalRes.ToComponentInfoDtos(isResponsible, policyDecisionDeniedReason, isAllowDeniedPolicyDecision),
		ComponentStats:                 evalRes.Stats,
		BulkPolicyDecisionDeniedReason: policyDecisionDeniedReason,
	}
	render.JSON(w, r, response)
}

func evaluatePolicyDecisionDeniedReason(isResponsible, isFossOfficeUser, isVehicle bool) string {
	if isVehicle {
		if !isFossOfficeUser {
			return message.PolicyDecisionDeniedNotFossOfficeUser
		}
	} else {
		if !isResponsible {
			return message.PolicyDecisionDeniedNotResponsible
		}
	}
	return ""
}

func evaluateIsAllowDeniedPolicyDecision(isDomainAdminUser, isFossOfficeUser, isVehicle bool) bool {
	if isVehicle {
		return isFossOfficeUser
	}
	return isDomainAdminUser
}

func (projectHandler *ProjectHandler) ExecuteChecklistsHandler(w http.ResponseWriter, r *http.Request) {
	sbomUuidEscaped := chi.URLParam(r, "sbomUuid")
	sbomUuid, err := url.QueryUnescape(sbomUuidEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamSbomUuidEmpty))
	if sbomUuid == "" {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ParamSbomUuidEmpty))
	}
	err = validation.CheckUuid(sbomUuid)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorKeyRequestParamNotValid, "sbomUuid"), zapcore.InfoLevel)

	body := extractExecuteChecklistsBody(r)

	currentProject, version, requestSession := projectHandler.retrieveProjectAndVersion2(r)
	userID, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowExecuteChecklist {
		exception.ThrowExceptionSendDeniedResponse()
	}
	projectHandler.ChecklistService.Execute(requestSession, currentProject, version, sbomUuid, body.Ids, userID)
	render.JSON(w, r, SuccessResponse{
		Success: true,
	})
}

func (projectHandler *ProjectHandler) GetGeneralVersionStatsHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, version, requestSession := projectHandler.retrieveProjectAndVersion2(r)

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProjectVersion.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadVersions))
	}
	generalStats := components.GeneralStats{
		SBOMDelivered:  false,
		SourceUploaded: false,
	}
	if projectHandler.SbomListRepository.ExistByKey(requestSession, version.Key) {
		generalStats.SBOMDelivered = true
	}
	if len(version.GetSourceExternalAll()) > 0 {
		generalStats.SourceUploaded = true
	}
	generalStats.ReviewRemark = projectHandler.calculateReviewRemarkStats(requestSession, version.Key)

	render.JSON(w, r, generalStats)
}

func (projectHandler *ProjectHandler) calculateReviewRemarkStats(requestSession *logy.RequestSession, versionKey string) components.ReviewRemarkStats {
	remarksList := projectHandler.ReviewRemarksRepository.FindByKey(requestSession, versionKey, false)
	if remarksList == nil {
		return components.ReviewRemarkStats{}
	}
	reviewRemarkStats := components.ReviewRemarkStats{
		Total:                  0,
		Acceptable:             0,
		AcceptableAfterChanges: 0,
		NotAcceptable:          0,
	}

	for _, remark := range remarksList.Remarks {
		if remark.Status == reviewremarks.Open || remark.Status == reviewremarks.Closed {
			switch remark.Level {
			case reviewremarks.Green:
				reviewRemarkStats.Acceptable++
			case reviewremarks.Yellow:
				reviewRemarkStats.AcceptableAfterChanges++
			case reviewremarks.Red:
				reviewRemarkStats.NotAcceptable++
			}
			reviewRemarkStats.Total++
		}
	}
	return reviewRemarkStats
}

func (projectHandler *ProjectHandler) GetSBOMStatsHandler(w http.ResponseWriter, r *http.Request) {
	sbomUuidEscaped := chi.URLParam(r, "sbomUuid")
	sbomUuid, err := url.QueryUnescape(sbomUuidEscaped)
	if err != nil || sbomUuid == "" {
		exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamSbomUuidEmpty))
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ParamSbomUuidEmpty))
	}
	err = validation.CheckUuid(sbomUuid)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorKeyRequestParamNotValid, "sbomUuid"), zapcore.InfoLevel)

	currentProject, version, requestSession := projectHandler.retrieveProjectAndVersion2(r)
	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProjectVersion.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ViewComponents))
		return
	}

	_, selectedSpdx := projectHandler.RetrieveSbomListAndFile(requestSession, version.Key, sbomUuid)
	if selectedSpdx == nil {
		render.JSON(w, r, components.SBOMStats{})
		return
	}

	compInfos := projectHandler.SpdxService.GetComponentInfos(requestSession, currentProject, version.Key, selectedSpdx)

	notChartFossLicense := components.NotChartFossLicenseStats{
		Total: 0,
	}

	fossOnly := r.URL.Query().Get("fossOnly")
	// narrow components with FOSS licenses only
	if len(fossOnly) > 0 {
		compInfosWithFossOnly := make([]components.ComponentInfo, 0)
		for compInfoIndex := range compInfos {
			compAdded := false
			compCounted := false
			for _, compLicense := range compInfos[compInfoIndex].GetLicensesEffective().List {
				if !compLicense.Known {
					continue
				}
				lic := projectHandler.LicenseRepository.FindById(requestSession, compLicense.ReferencedLicense)
				if lic.Meta.LicenseType == license.OpenSource || lic.Meta.LicenseType == license.PublicDomain {
					if !compAdded {
						compInfosWithFossOnly = append(compInfosWithFossOnly, compInfos[compInfoIndex])
						compAdded = true
					}
					if !lic.Meta.IsLicenseChart && !compCounted {
						notChartFossLicense.Total++
						compCounted = true
					}
					if compAdded && compCounted {
						break
					}
				}
			}
		}
		compInfos = compInfosWithFossOnly
	}

	rules := projectHandler.PolicyRuleRepository.FindPolicyRulesForLabel(requestSession, currentProject.PolicyLabels)
	policyDecisions := projectHandler.PolicyDecisionsRepository.FindByKey(requestSession, currentProject.Key, false)
	isVehicle := projectHandler.ProjectLabelService.HasVehiclePlatformLabel(requestSession, currentProject)

	policyEvaluation := compInfos.EvaluatePolicyRules(rules, policyDecisions, isVehicle, selectedSpdx.Uploaded, selectedSpdx.Key)
	licenseStats := calculateLicenseStats(compInfos)
	scanRemarkStats, scanRemarkTypeStats := projectHandler.calculateScanRemarkStats(requestSession, currentProject, selectedSpdx, policyEvaluation)
	licenseRemarkStats := projectHandler.calculateLicenseRemarkStats(requestSession, policyEvaluation)
	approvalInfo := components.InApproval{
		IsInApproval: selectedSpdx.ApprovalInfo.IsInApproval,
		ApprovalGuid: selectedSpdx.ApprovalInfo.ApprovalGuid,
		Status:       selectedSpdx.ApprovalInfo.Status,
	}

	render.JSON(w, r, components.SBOMStats{
		PolicyState:         policyEvaluation.Stats,
		LicenseFamily:       licenseStats,
		ScanRemark:          scanRemarkStats,
		LicenseRemark:       licenseRemarkStats,
		ApprovalInfo:        approvalInfo,
		ScanRemarkType:      scanRemarkTypeStats,
		NotChartFossLicense: notChartFossLicense,
	})
}

func (projectHandler *ProjectHandler) calculateLicenseRemarkStats(requestSession *logy.RequestSession,
	policyEvaluation *components.EvaluationResult,
) components.LicenseRemarkStats {
	licenseRemarkStats := components.LicenseRemarkStats{
		Total:       0,
		Information: 0,
		Warning:     0,
		Alarm:       0,
	}
	licenseRemarks := licenseremarks.CreateQualityLicenseRemarks(requestSession, projectHandler.LicenseRepository, projectHandler.ObligationRepository, policyEvaluation)
	if licenseRemarks == nil {
		return components.LicenseRemarkStats{}
	}
	for _, licenseRemark := range licenseRemarks {
		for _, obl := range licenseRemark.Obligations {
			switch strings.ToLower(string(obl.WarnLevel)) {
			case obligation.Information:
				licenseRemarkStats.Information++
			case obligation.Warning:
				licenseRemarkStats.Warning++
			case obligation.Alarm:
				licenseRemarkStats.Alarm++
			}
			licenseRemarkStats.Total++
		}
	}
	return licenseRemarkStats
}

func (projectHandler *ProjectHandler) calculateScanRemarkStats(requestSession *logy.RequestSession, currentProject *project.Project,
	selectedSpdx *project.SpdxFileBase, policyEvaluation *components.EvaluationResult,
) (components.ScanRemarkStats, components.ScanRemarkTypeStats) {
	scanRemarkStats := components.ScanRemarkStats{
		Total:       0,
		Information: 0,
		Warning:     0,
		Problem:     0,
	}
	scanRemarkTypeStats := components.ScanRemarkTypeStats{
		Total:               0,
		MissingCopyrights:   0,
		MalformedCopyrights: 0,
	}
	scanRemarks := projectHandler.createQualityScanRemarks(requestSession, currentProject, &selectedSpdx.MetaInfo, policyEvaluation, projectHandler.LabelRepository)
	if scanRemarks == nil {
		return components.ScanRemarkStats{}, components.ScanRemarkTypeStats{}
	}
	for _, scanremark := range scanRemarks {
		switch scanremark.Status {
		case project.ScanRemarkStatus(project.INFORMATION):
			scanRemarkStats.Information++
		case project.ScanRemarkStatus(project.WARNING):
			scanRemarkStats.Warning++
		case project.ScanRemarkStatus(project.PROBLEM):
			scanRemarkStats.Problem++
		}
		scanRemarkStats.Total++
		switch scanremark.RemarkKey {
		case message.SrMissingCopyrightText:
			scanRemarkTypeStats.MissingCopyrights++
			if len(scanRemarkTypeStats.MissingCopyrightsLevel) == 0 {
				scanRemarkTypeStats.MissingCopyrightsLevel = string(scanremark.Status)
			}
		case message.SrMalformedCopyrightText:
			scanRemarkTypeStats.MalformedCopyrights++
		}
		scanRemarkTypeStats.Total++
	}
	return scanRemarkStats, scanRemarkTypeStats
}

func calculateLicenseStats(compInfos []components.ComponentInfo) components.LicenseFamilyStats {
	var licenseStats components.LicenseFamilyStats

	for _, ci := range compInfos {
		if ci.Type == "Root" && ci.License == "NOASSERTION" {
			continue
		}
		licenseStats.Total++
		if helper.IsUnasserted(ci.GetLicenseEffective()) {
			licenseStats.Other++
			continue
		}

		worst := ci.WorstFamily()
		switch worst {
		case license.NetworkCopyleft:
			licenseStats.NetworkCopyLeft++
		case license.StrongCopyleft:
			licenseStats.StrongCopyLeft++
		case license.WeakCopyleft:
			licenseStats.WeakCopyLeft++
		case license.Permissive:
			licenseStats.Permissive++
		default:
			licenseStats.Other++
		}
	}
	return licenseStats
}

func (projectHandler *ProjectHandler) ProjectVersionComponentsBySearch(w http.ResponseWriter, r *http.Request) {
	sbomUuidEscaped := chi.URLParam(r, "sbomUuid")
	sbomUuid, err := url.QueryUnescape(sbomUuidEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamSbomUuidEmpty))
	if sbomUuid == "" {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ParamSbomUuidEmpty))
	}
	err = validation.CheckUuid(sbomUuid)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorKeyRequestParamNotValid, "sbomUuid"), zapcore.InfoLevel)

	currentProject, version, requestSession := projectHandler.retrieveProjectAndVersion2(r)

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProjectVersion.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ViewComponents))
	}

	_, latestSpdx := projectHandler.RetrieveSbomListAndFile(requestSession, version.Key, sbomUuid)
	if latestSpdx == nil {
		render.JSON(w, r, components.ComponentsInfoResponse{})
		return
	}

	searchFragmentEscaped := chi.URLParam(r, "searchFragment")
	searchFragment, err := url.QueryUnescape(searchFragmentEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamSearchFragmentEmpty))

	compInfos := projectHandler.SpdxService.GetComponentInfos(requestSession, currentProject, version.Key, latestSpdx)
	search := compInfos.FindComponentsByNameFragment(searchFragment)

	response := make([]components.ComponentInfoSlimDto, 0)
	for _, info := range search {
		response = append(response, *info.ToComponentInfoSlimDto())
	}

	render.JSON(w, r, response)
}

func extractExecuteChecklistsBody(r *http.Request) checklist.ExecuteRequestDto {
	var res checklist.ExecuteRequestDto
	validation.DecodeAndValidate(r, &res, false)
	return res
}

func extractVersionRequestBody(r *http.Request, handleErrorAsServerException bool) project.VersionRequestDto {
	var versionData project.VersionRequestDto
	validation.DecodeAndValidate(r, &versionData, handleErrorAsServerException)
	return versionData
}

func extractProjectSearchBody(r *http.Request, handleErrorAsServerException bool) project.ProjectSearchDto {
	var req project.ProjectSearchDto
	validation.DecodeAndValidate(r, &req, handleErrorAsServerException)
	return req
}

func isSchemeAllowed(scheme string) bool {
	return scheme == "https" || scheme == "file" || scheme == "http" || scheme == "ssh"
}

func extractExternalSourceRequestBody(r *http.Request, handleErrorAsServerException bool) project.SourceExternalDTO {
	var sourceData project.SourceExternalDTO
	validation.DecodeAndValidate(r, &sourceData, handleErrorAsServerException)

	url, err := url.ParseRequestURI(sourceData.URL)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.UrlInvalid))
	if !isSchemeAllowed(url.Scheme) {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.UrlInvalid), "")
	}

	if !helper.UrlRegex.MatchString(sourceData.URL) {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.UrlInvalid), "")
	}

	return sourceData
}

func extractVersionKeyFromRequest(r *http.Request) string {
	versionEscaped := chi.URLParam(r, "version")
	versionKey, err := url.QueryUnescape(versionEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamVersionWrong))

	err = validation.CheckUuid(versionKey)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamVersionWrong))
	return versionKey
}

// ProjectSPDXExternCheckOnDemand godoc
//
//	@Summary	Get status information of uploaded SBOM file
//	@Id			getSBOMStatusInformation
//	@Produce	json
//	@Accept		mpfd
//	@Param		uuid	path		string							true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Param		file	formData	file							true	"SBOM file"
//	@Success	200		{object}	project.SpdxStatusInformation	"SPDX Status Information"
//	@Failure	401		{object}	exception.HttpError				"Unauthorized Error"
//	@Failure	417		{object}	project.SPDXUploadResponse		"Validation Error"
//	@Router		/projects/{uuid}/sbomcheck [post]
//	@security	Bearer
func (projectHandler *ProjectHandler) ProjectSPDXExternCheckOnDemand(w http.ResponseWriter, r *http.Request) {
	validation.CheckExpectedContentType(r, validation.ContentTypeFormData)

	requestSession := logy.GetRequestSession(r)

	currentProject, _ := projectHandler.retrieveProjectFromPublicRequest(requestSession, r, false)
	l, acquired := projectHandler.LockService.Acquire(locks.Options{
		Key:      currentProject.Key,
		Blocking: true,
		Timeout:  time.Second * 10,
	})
	if !acquired {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ResourceInUse), "")
	}
	logy.Infof(requestSession, "Acquired!")
	defer func() {
		projectHandler.LockService.Release(l)
		logy.Infof(requestSession, "Released lock")
	}()
	currentProject, _ = projectHandler.retrieveProjectFromPublicRequest(requestSession, r, false)
	TryNewFileUpload(requestSession, currentProject.Key, projectHandler.ProjectRepository)

	file, handler, err := r.FormFile("file")
	if err != nil {
		// max 10mb
		err = r.ParseMultipartForm(10 << 20)
	}
	exception.HandleErrorClientMessage(err, message.GetI18N(message.SpdxFileEmptyOrLarge))
	defer file.Close()

	spdxBytes, err := io.ReadAll(file)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.SpdxFileRead))
	spdxString := string(spdxBytes)

	validation.CheckExpectedContentType2(handler.Header, []validation.ContentType{
		validation.ContentTypeJson,
		validation.ContentTypeOctets,
	})

	spdxFile := &project.SpdxFileBase{
		ChildEntity:  domain.NewChildEntity(),
		Type:         schema.JSON,
		ContentValid: true,
		Uploaded:     reflection.ToPointer(time.Now()),
	}

	err = service.ValidateSbom(requestSession, spdxString, spdxFile, currentProject, projectHandler.SchemaRepository)
	if err != nil {
		render.Status(r, 417)
		render.JSON(w, r, project.SPDXUploadResponse{
			DocIsValid:              false,
			ValidationFailedMessage: err.Error(), FileUploaded: false,
		})
		return
	}

	currentRefs := projectHandler.LicenseRepository.GetLicenseRefs(requestSession)
	ci := project.FileContent(spdxString).ExtractComponentInfo(requestSession)
	ci.EnrichComponentInfos(requestSession)
	ci.ApplyRefs(currentRefs)

	// #6642: apply license rules
	licenseRules := projectHandler.LicenseRulesRepository.FindByKey(requestSession, currentProject.Key, false)
	ci.ApplyLicenseRules(licenseRules, spdxFile.Uploaded, spdxFile.Key)

	projectHandler.ProjectSPDXExternCheck(w, r, requestSession, currentProject, ci, spdxFile.Uploaded, spdxFile.Key)
}

// ProjectVersionSPDXExternCheck godoc
//
//	@Summary	Get SBOM status information of project version (also known as channel)
//	@Id			getProjectVersionSBOMStatus
//	@Produce	json
//	@Param		uuid		path		string							true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Param		version		path		string							true	"Project Version Name (also known as Channel Name) e.g.: main"
//	@Param		sbomUuid	path		string							true	"UUID of the SBOM delivery or 'latest' for the latest SBOM delivery e.g.: dummy-sbom-id---xxx-4413-yyy-24f060311111"
//	@Success	200			{object}	project.SpdxStatusInformation	"SPDX Status Information"
//	@Failure	404			{object}	exception.HttpError404			"NotFound Error"
//	@Failure	401			{object}	exception.HttpError				"Unauthorized Error"
//	@Router		/projects/{uuid}/versions/{version}/sboms/{sbomUuid}/check [get]
//	@security	Bearer
func (projectHandler *ProjectHandler) ProjectVersionSPDXExternCheck(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	currentProject, version, _ := projectHandler.retrieveProjectAndVersionFromPublicRequest(requestSession, r)

	sbomUuidEscaped := chi.URLParam(r, "sbomUuid")
	sbomUuid := ValidateIDOrLatest(sbomUuidEscaped)

	var spdx *project.SpdxFileBase
	if sbomUuid == "latest" {
		_, spdx = projectHandler.retrieveSbomListAndLatestFile(requestSession, version.Key)
	} else {
		_, spdx = projectHandler.RetrieveSbomListAndFile(requestSession, version.Key, sbomUuid)
	}
	if spdx == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.FindingSbomKey), "SPDX not found in history: "+sbomUuid)
		return
	}
	compInfos := projectHandler.SpdxService.GetComponentInfos(requestSession, currentProject, version.Key, spdx)
	projectHandler.ProjectSPDXExternCheck(w, r, requestSession, currentProject, compInfos, spdx.Uploaded, spdx.Key)
}

func (projectHandler *ProjectHandler) ProjectSPDXExternCheck(w http.ResponseWriter, r *http.Request,
	requestSession *logy.RequestSession, currentProject *project.Project, componentInfo components.ComponentInfos, sbomUpload *time.Time, sbomKey string,
) {
	spdxStatusInformation := projectHandler.CreateProjectSPDXExternCheck(requestSession, currentProject, componentInfo, sbomUpload, sbomKey)
	render.JSON(w, r, spdxStatusInformation)
}

func (projectHandler *ProjectHandler) CreateProjectSPDXExternCheck(requestSession *logy.RequestSession,
	currentProject *project.Project, componentInfo components.ComponentInfos, sbomUpload *time.Time, sbomKey string,
) project.SpdxStatusInformation {
	rules := projectHandler.PolicyRuleRepository.FindPolicyRulesForLabel(requestSession, currentProject.PolicyLabels)
	policyDecisions := projectHandler.PolicyDecisionsRepository.FindByKey(requestSession, currentProject.Key, false)
	isVehicle := projectHandler.ProjectLabelService.HasVehiclePlatformLabel(requestSession, currentProject)
	evalRes := componentInfo.EvaluatePolicyRules(rules, policyDecisions, isVehicle, sbomUpload, sbomKey)

	qualityLicenseRemarks := projectHandler.createQualityLicenseRemarks(requestSession, evalRes)
	qualityScanRemarks := projectHandler.createQualityScanRemarks(requestSession, currentProject, &project.MetaInfo{}, evalRes, projectHandler.LabelRepository)

	disclaimerContent := readDisclaimerTextFromFile(requestSession, "./resources/disclaimerContent.md")
	scanRemarksContent := readDisclaimerTextFromFile(requestSession, "./resources/scanRemarksContent.md")
	licenseRemarksContent := readDisclaimerTextFromFile(requestSession, "./resources/licenseRemarksContent.md")
	generalRemarksContent := readDisclaimerTextFromFile(requestSession, "./resources/generalRemarksContent.md")

	if disclaimerContent == nil || scanRemarksContent == nil || generalRemarksContent == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.StaticReadFile), "")
	}

	components := make([]project.SpdxStatusComponent, 0)
	for _, compRes := range evalRes.Results {
		componentInfoElement := compRes.Component
		statusComponent := project.SpdxStatusComponent{
			SpdxId:  componentInfoElement.SpdxId,
			License: componentInfoElement.GetLicenseEffective(),
			Name:    componentInfoElement.Name,
			Version: componentInfoElement.Version,
		}

		if compRes.Status != nil {
			policyRuleStatus := make([]project.SpdxStatusPolicy, 0)
			for _, policyRule := range compRes.Status {
				policyRuleStatus = append(policyRuleStatus, project.SpdxStatusPolicy{
					Name:           policyRule.Name,
					LicenseMatched: policyRule.LicenseMatched,
					Type:           policyRule.Type,
					Used:           policyRule.Used,
					Description:    policyRule.Description,
				})
			}
			statusComponent.PolicyRuleStatus = policyRuleStatus
		}

		statusScanRemarks := make([]project.SpdxStatusScanRemarks, 0)
		for _, scanRemarks := range qualityScanRemarks {
			if componentInfoElement.SpdxId == scanRemarks.SpdxId {
				statusScanRemarks = append(statusScanRemarks, project.SpdxStatusScanRemarks{
					Status:      scanRemarks.Status,
					Remark:      message.GetI18N(scanRemarks.RemarkKey).Text,
					Description: message.GetI18N(scanRemarks.DescriptionKey).Text,
				})
			}
		}
		if len(statusScanRemarks) > 0 {
			statusComponent.ScanRemarks = statusScanRemarks
		}

		statusLicenseRemarks := make([]project.SpdxStatusLicenseRemarks, 0)
		for _, licenseRemarks := range qualityLicenseRemarks {
			if componentInfoElement.SpdxId == licenseRemarks.SpdxId {
				statusLicenseRemarks = append(statusLicenseRemarks, project.SpdxStatusLicenseRemarks{
					Status:         licenseRemarks.Status,
					Remark:         licenseRemarks.Remark,
					Type:           licenseRemarks.Type,
					LicenseMatched: licenseRemarks.License,
					Description:    licenseRemarks.Description,
				})
			}
		}
		if len(statusLicenseRemarks) > 0 {
			statusComponent.LicenseRemarks = statusLicenseRemarks
		}

		usedAliases := make([]project.UsedAlias, 0)
		for _, l := range compRes.Component.GetLicensesEffective().List {
			if !l.Known {
				continue
			}
			usedAliases = append(usedAliases, project.UsedAlias{
				Name:           l.OrigName,
				ReferencedName: l.ReferencedLicense,
			})
		}
		statusComponent.UsedAliases = usedAliases

		if compRes.Component.LicenseRuleApplied != nil {
			statusComponent.UsedDecision = &project.UsedDecision{
				Expression:  compRes.Component.LicenseRuleApplied.LicenseExpression,
				LicenseID:   compRes.Component.LicenseRuleApplied.LicenseDecisionId,
				LicenseName: compRes.Component.LicenseRuleApplied.LicenseDecisionName,
			}
		}

		components = append(components, statusComponent)
	}

	spdxStatusInformation := project.SpdxStatusInformation{
		Disclaimer:     string(disclaimerContent),
		ScanRemarks:    string(scanRemarksContent),
		LicenseRemarks: string(licenseRemarksContent),
		GeneralRemarks: string(generalRemarksContent),
		Components:     components,
	}
	return spdxStatusInformation
}

func readDisclaimerTextFromFile(requestSession *logy.RequestSession, path string) []byte {
	content, err := os.ReadFile(path)
	if err != nil {
		logy.Errorf(requestSession, "Could not read file: "+path, err)
		return nil
	}
	return content
}

func (projectHandler *ProjectHandler) CreateOverallReview(w http.ResponseWriter, r *http.Request) {
	currentProject, version, requestSession := projectHandler.retrieveProjectAndVersion2(r)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	username, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProjectVersion.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadVersionReviewRemarks))
	}

	createData := extractOverallReviewBody(r, true)

	if createData.State == overallreview2.Audited && !rights.IsDomainAdmin() && !rights.IsFossOffice() && !projectHandler.ProjectLabelService.HasVehiclePlatformLabel(requestSession, currentProject) {
		exception.ThrowExceptionSendDeniedResponse()
	}

	projectHandler.OverallReviewService.AddToProjectFromDTO(
		requestSession,
		currentProject,
		version,
		username,
		createData,
	)

	responseData := SuccessResponse{
		Success: true,
		Message: "overall review created",
	}
	render.JSON(w, r, responseData)
}

func logRequestBody(requestSession *logy.RequestSession, r *http.Request) {
	if r.Body != nil {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			logy.Warnf(requestSession, "Error reading request body: %v", err)
		} else {
			logy.Infof(requestSession, "Request Body: %s", string(bodyBytes))
		}
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
}

func (projectHandler *ProjectHandler) GetReviewRemarks(w http.ResponseWriter, r *http.Request) {
	remarksList, _ := projectHandler.retrieveReviewRemarksWithAuth(r)

	res := make([]reviewremarks.RemarkDto, 0)
	if remarksList == nil {
		render.JSON(w, r, res)
		return
	}
	for _, r := range remarksList.Remarks {
		res = append(res, r.ToDto())
	}
	render.JSON(w, r, res)
}

func (projectHandler *ProjectHandler) DownloadReviewRemarksHandler(w http.ResponseWriter, r *http.Request) {
	remarksList, requestSession := projectHandler.retrieveReviewRemarksWithAuth(r)

	// Convert to the format expected by CSV writer
	var remarks []*reviewremarks.ReviewRemarks
	if remarksList != nil {
		remarks = []*reviewremarks.ReviewRemarks{remarksList}
	} else {
		remarks = make([]*reviewremarks.ReviewRemarks, 0)
	}

	writeReviewRemarksAsCsvIntoResponse(requestSession, &w, remarks)
	w.WriteHeader(http.StatusOK)
}

// Shared method that handles validation, authorization, and data retrieval
func (projectHandler *ProjectHandler) retrieveReviewRemarksWithAuth(r *http.Request) (*reviewremarks.ReviewRemarks, *logy.RequestSession) {
	currentProject, version, requestSession := projectHandler.retrieveProjectAndVersion2(r)

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProjectVersion.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadVersionReviewRemarks))
	}

	return projectHandler.ReviewRemarksRepository.FindByKey(requestSession, version.Key, false), requestSession
}

func writeReviewRemarksAsCsvIntoResponse(requestSession *logy.RequestSession, w *http.ResponseWriter, reviewRemarks []*reviewremarks.ReviewRemarks) {
	csvWriter := csv.NewWriter(*w)
	// Use semicolon as delimiter for better Excel compatibility with comma-containing content
	csvWriter.Comma = ';'
	csvWriter.UseCRLF = true
	defer csvWriter.Flush()

	// Add BOM for Excel to correctly detect UTF-8
	_, err := (*w).Write([]byte{0xEF, 0xBB, 0xBF})
	if err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCsvGeneration, "BOM header", "write"), err)
	}

	csvHeader := []string{"ENTRY_TYPE", "UUID", "STATUS", "LEVEL", "REVIEW REMARK", "DESCRIPTION", "CREATOR", "CREATED", "COMPONENTS", "LICENSES", "COMMENT UPDATE"}

	if err := csvWriter.Write(csvHeader); err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCsvGeneration, "review remark", "header"), err)
	}

	for _, reviewRemark := range reviewRemarks {
		if reviewRemark.Remarks == nil {
			continue
		}

		for _, remark := range reviewRemark.Remarks {

			// Other fields can be prepared normally
			csvRow := []string{
				"REMARK",
				remark.Key, // Include the UUID
				string(remark.Status),
				string(remark.Level),
				csvutil.PrepareFieldForCsv(remark.Title),
				csvutil.PrepareFieldForCsv(remark.Description),
				csvutil.PrepareFieldForCsv(remark.Author),
				remark.Created.Format("2006-01-02 15:04:05"),
				prepareComponents(remark),
				prepareLicenses(remark),
				"", // Empty comment content for main remark
			}
			if err := csvWriter.Write(csvRow); err != nil {
				exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCsvGeneration, "review remark", "data"), err)
			}

			// Write all comments for this remark in chronological order
			if len(remark.Events) > 0 {
				for _, event := range remark.Events {
					if event.Type != "COMMENT" {
						continue
					}

					var commentData reviewremarks.Comment
					if err := json.Unmarshal(event.Content, &commentData); err != nil {
						logy.Warnf(requestSession, "Error during unmarshal: %v", err)
					}

					commentRow := []string{
						"COMMENT",
						remark.Key,                               // Include parent remark UUID
						string(remark.Status),                    // Keep parent remark status
						string(remark.Level),                     // Keep parent remark level
						csvutil.PrepareFieldForCsv(remark.Title), // Keep parent remark title
						"",                                       // Empty description for comments
						csvutil.PrepareFieldForCsv(event.Author),
						event.Created.Format("2006-01-02 15:04:05"),
						prepareComponents(remark),
						prepareLicenses(remark),
						csvutil.PrepareFieldForCsv(string(commentData)),
					}
					if err := csvWriter.Write(commentRow); err != nil {
						exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCsvGeneration, "review remark comment", "data"), err)
					}
				}
			}
		}
	}
}

func prepareComponents(remark *reviewremarks.Remark) string {
	if remark.Components != nil && len(remark.Components) > 0 {
		components := make([]string, len(remark.Components))
		for i, comp := range remark.Components {
			components[i] = fmt.Sprintf("%s (%s)", csvutil.PrepareFieldForCsv(comp.ComponentName), csvutil.PrepareFieldForCsv(comp.ComponentVersion))
		}
		return strings.Join(components, "; ")
	}
	return ""
}

func prepareLicenses(remark *reviewremarks.Remark) string {
	if remark.Licenses != nil && len(remark.Licenses) > 0 {
		licenses := make([]string, len(remark.Licenses))
		for i, lic := range remark.Licenses {
			if lic.LicenseName == "" {
				licenses[i] = fmt.Sprintf("%s (unknown)", csvutil.PrepareFieldForCsv(lic.LicenseId))
			} else {
				licenses[i] = csvutil.PrepareFieldForCsv(lic.LicenseName)
			}
		}
		return strings.Join(licenses, "; ")
	}
	return ""
}

func extractOverallReviewBody(r *http.Request, handleErrorAsServerException bool) (createData overallreview2.OverallReviewRequestDto) {
	validation.DecodeAndValidate(r, &createData, handleErrorAsServerException)
	return
}

func extractReviewRemarkBody(r *http.Request, handleErrorAsServerException bool) (createData reviewremarks.ReviewRemarkRequestDto) {
	validation.DecodeAndValidate(r, &createData, handleErrorAsServerException)
	return
}

func extractCommentReviewRemarkBody(r *http.Request, handleErrorAsServerException bool) (commentData reviewremarks.CommentRequestDto) {
	validation.DecodeAndValidate(r, &commentData, handleErrorAsServerException)
	return
}

func extractCommentReviewRemarkExternBody(r *http.Request, handleErrorAsServerException bool) (commentData reviewremarks.RRCommentExternDTO) {
	validation.DecodeAndValidate(r, &commentData, handleErrorAsServerException)
	return
}

func extractSetReviewRemarkStatusBody(r *http.Request, handleErrorAsServerException bool) (statusData reviewremarks.SetRemarkStatusRequestDto) {
	validation.DecodeAndValidate(r, &statusData, handleErrorAsServerException)
	return
}

func (projectHandler *ProjectHandler) CreateReviewRemark(w http.ResponseWriter, r *http.Request) {
	currentProject, version, requestSession := projectHandler.retrieveProjectAndVersion2(r)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	username, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if currentProject.GetMember(username) == nil && !rights.IsProjectAnalyst() && !rights.IsDomainAdmin() {
		if !roles.CanAccessVehicleProjectOperations(rights, projectHandler.ProjectLabelService.HasVehiclePlatformLabel(requestSession, currentProject)) {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadVersionReviewRemarks))
		}
	}

	createData := extractReviewRemarkBody(r, true)

	rrs := reviewRemarksService.ReviewRemarksService{
		RequestSession:          requestSession,
		LicenseRepo:             projectHandler.LicenseRepository,
		AuditLogListRepository:  projectHandler.AuditLogListRepository,
		ReviewRemarksRepository: projectHandler.ReviewRemarksRepository,
		Retriever:               projectHandler,
		LicenseRulesRepo:        projectHandler.LicenseRulesRepository,
		SpdxService:             projectHandler.SpdxService,
	}
	if !rrs.CreateReviewRemark(currentProject, version.Key, createData, username) {
		exception.ThrowExceptionBadRequestResponse()
		return
	}

	if len(createData.SBOMId) > 0 {
		sbomList := projectHandler.SbomListRepository.FindByKey(requestSession, version.Key, false)
		changed := false
		for _, sbom := range sbomList.SpdxFileHistory {
			if sbom.Key == createData.SBOMId && !sbom.IsInUse {
				sbom.IsInUse = true
				changed = true
			}
		}
		if changed {
			projectHandler.SbomListRepository.Update(requestSession, sbomList)
		}
	}

	responseData := SuccessResponse{
		Success: true,
		Message: "remark created",
	}
	render.JSON(w, r, responseData)
}

func (projectHandler *ProjectHandler) EditReviewRemark(w http.ResponseWriter, r *http.Request) {
	currentProject, version, requestSession := projectHandler.retrieveProjectAndVersion2(r)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	username, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if currentProject.GetMember(username) == nil && !rights.IsProjectAnalyst() && !rights.IsDomainAdmin() {
		if !roles.CanAccessVehicleProjectOperations(rights, projectHandler.ProjectLabelService.HasVehiclePlatformLabel(requestSession, currentProject)) {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadVersionReviewRemarks))
		}
	}

	remarkUuidEscaped := chi.URLParam(r, "remarkId")
	remarkUuid, err := url.QueryUnescape(remarkUuidEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamRemarkUuidEmpty))
	if remarkUuid == "" {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ParamRemarkUuidEmpty))
	}
	err = validation.CheckUuid(remarkUuid)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorKeyRequestParamNotValid, "remarkUuid"), zapcore.InfoLevel)

	editData := extractReviewRemarkBody(r, true)

	rrs := reviewRemarksService.ReviewRemarksService{
		RequestSession:          requestSession,
		LicenseRepo:             projectHandler.LicenseRepository,
		AuditLogListRepository:  projectHandler.AuditLogListRepository,
		ReviewRemarksRepository: projectHandler.ReviewRemarksRepository,
		Retriever:               projectHandler,
		LicenseRulesRepo:        projectHandler.LicenseRulesRepository,
		SpdxService:             projectHandler.SpdxService,
	}

	if !rrs.EditReviewRemark(currentProject, version.Key, remarkUuid, username, projectHandler.fullNameForUserSafe(requestSession, username), editData) {
		exception.ThrowExceptionBadRequestResponse()
		return
	}
	responseData := SuccessResponse{
		Success: true,
		Message: "remark edited",
	}
	render.JSON(w, r, responseData)
}

func (projectHandler *ProjectHandler) CommentReviewRemark(w http.ResponseWriter, r *http.Request) {
	currentProject, version, requestSession := projectHandler.retrieveProjectAndVersion2(r)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	username, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if currentProject.GetMember(username) == nil && !rights.IsProjectAnalyst() && !rights.IsDomainAdmin() {
		if !roles.CanAccessVehicleProjectOperations(rights, projectHandler.ProjectLabelService.HasVehiclePlatformLabel(requestSession, currentProject)) {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadVersionReviewRemarks))
		}
	}

	remarkUuidEscaped := chi.URLParam(r, "remarkId")
	remarkUuid, err := url.QueryUnescape(remarkUuidEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamRemarkUuidEmpty))
	if remarkUuid == "" {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ParamRemarkUuidEmpty))
	}
	err = validation.CheckUuid(remarkUuid)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorKeyRequestParamNotValid, "remarkUuid"), zapcore.InfoLevel)
	commentData := extractCommentReviewRemarkBody(r, true)

	remarks := projectHandler.ReviewRemarksRepository.FindByKey(requestSession, version.Key, false)
	var remark *reviewremarks.Remark
	for _, r := range remarks.Remarks {
		if r.Key == remarkUuid {
			remark = r
			break
		}
	}
	if remark == nil {
		exception.ThrowExceptionBadRequestResponse()
	}
	var before reviewremarks.Remark
	copier.Copy(&before, remark)
	remark.Comment(username, projectHandler.fullNameForUserSafe(requestSession, username), commentData.Content)
	projectHandler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, version.Key, username, message.ReviewRemarkCommented, cmp.Diff, *remark, before)
	projectHandler.ReviewRemarksRepository.Update(requestSession, remarks)
	responseData := SuccessResponse{
		Success: true,
		Message: "comment created",
	}
	render.JSON(w, r, responseData)
}

func (projectHandler *ProjectHandler) SetReviewRemarkStatus(w http.ResponseWriter, r *http.Request) {
	currentProject, version, requestSession := projectHandler.retrieveProjectAndVersion2(r)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	username, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if currentProject.GetMember(username) == nil && !rights.IsProjectAnalyst() && !rights.IsDomainAdmin() {
		if !roles.CanAccessVehicleProjectOperations(rights, projectHandler.ProjectLabelService.HasVehiclePlatformLabel(requestSession, currentProject)) {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadVersionReviewRemarks))
		}
	}

	remarkUuidEscaped := chi.URLParam(r, "remarkId")
	remarkUuid, err := url.QueryUnescape(remarkUuidEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamRemarkUuidEmpty))
	if remarkUuid == "" {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ParamRemarkUuidEmpty))
	}
	err = validation.CheckUuid(remarkUuid)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorKeyRequestParamNotValid, "remarkUuid"), zapcore.InfoLevel)
	statusData := extractSetReviewRemarkStatusBody(r, true)

	valid, status := reviewremarks.ParseStatus(statusData.Status)
	if !valid {
		exception.ThrowExceptionBadRequestResponse()
	}

	remarks := projectHandler.ReviewRemarksRepository.FindByKey(requestSession, version.Key, false)
	var remark *reviewremarks.Remark
	for _, r := range remarks.Remarks {
		if r.Key == remarkUuid {
			remark = r
			break
		}
	}
	if remark == nil {
		exception.ThrowExceptionBadRequestResponse()
	}

	var before reviewremarks.Remark
	copier.Copy(&before, remark)
	if status == reviewremarks.Closed {
		remark.Close(username, projectHandler.fullNameForUserSafe(requestSession, username))
	} else if status == reviewremarks.Cancelled {
		remark.Cancel(username, projectHandler.fullNameForUserSafe(requestSession, username))
	} else if status == reviewremarks.InProgress {
		remark.InProgress(username, projectHandler.fullNameForUserSafe(requestSession, username))
	} else if status == reviewremarks.Open {
		remark.Reopen(username, projectHandler.fullNameForUserSafe(requestSession, username))
	} else {
		exception.ThrowExceptionBadRequestResponse()
	}
	projectHandler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, version.Key, username, message.ReviewRemarkChanged, cmp.Diff, *remark, before)
	projectHandler.ReviewRemarksRepository.Update(requestSession, remarks)
	responseData := SuccessResponse{
		Success: true,
		Message: "status updated",
	}
	render.JSON(w, r, responseData)
}

// ProjectVersionReviewRemarksExtern godoc
//
//	@Summary	Get review remarks for version (also known as channel)
//	@Id			getProjectVersionReviewRemarks
//	@Produce	json
//	@Param		uuid	path		string							true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Param		version	path		string							true	"Project Version Name (also known as Channel Name) e.g.: main"
//	@Success	200		{array}		reviewremarks.RemarkDtoExternV1	"Review remarks"
//	@Failure	404		{object}	exception.HttpError404			"NotFound Error"
//	@Failure	401		{object}	exception.HttpError				"Unauthorized Error"
//	@Router		/projects/{uuid}/versions/{version}/reviewremarks [get]
//	@security	Bearer
func (projectHandler *ProjectHandler) ProjectVersionReviewRemarksExtern(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, version, _ := projectHandler.retrieveProjectAndVersionFromPublicRequest(requestSession, r)

	res := make([]reviewremarks.RemarkDtoExternV1, 0)
	remarksList := projectHandler.ReviewRemarksRepository.FindByKey(requestSession, version.Key, false)
	if remarksList == nil {
		render.JSON(w, r, res)
		return
	}
	for _, r := range remarksList.Remarks {
		res = append(res, r.ToExternV1Dto())
	}
	render.JSON(w, r, res)
}

// ProjectVersionReviewRemarksExternV2 godoc
//
//	@Summary	Get review remarks for version (also known as channel)
//	@Id			getProjectVersionReviewRemarks
//	@Produce	json
//	@Param		uuid	path		string							true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Param		version	path		string							true	"Project Version Name (also known as Channel Name) e.g.: main"
//	@Success	200		{array}		reviewremarks.RemarkDtoExternV2	"Review remarks"
//	@Failure	404		{object}	exception.HttpError404			"NotFound Error"
//	@Failure	401		{object}	exception.HttpError				"Unauthorized Error"
//	@Router		/projects/{uuid}/versions/{version}/reviewremarks [get]
//	@security	Bearer
func (projectHandler *ProjectHandler) ProjectVersionReviewRemarksExternV2(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, version, _ := projectHandler.retrieveProjectAndVersionFromPublicRequest(requestSession, r)
	res := make([]reviewremarks.RemarkDtoExternV2, 0)
	remarksList := projectHandler.ReviewRemarksRepository.FindByKey(requestSession, version.Key, false)
	if remarksList == nil {
		render.JSON(w, r, res)
		return
	}
	for _, r := range remarksList.Remarks {
		res = append(res, r.ToExternV2Dto())
	}
	render.JSON(w, r, res)
}

// ProjectVersionReviewRemarksCommentExtern godoc
//
//	@Summary	Comment on a review remark
//	@Id			projectVersionReviewRemarkcomment
//	@Produce	json
//	@Param		uuid				path		string								true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Param		version				path		string								true	"Project Version Name (also known as Channel Name) e.g.: main"
//	@Param		reviewRemarkUuid	path		string								true	"Review remark UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Param		source				body		reviewremarks.RRCommentExternDTO	true	"Comment"
//	@Success	200					{object}	rest.SuccessResponse				"Success Response"
//	@Failure	401					{object}	exception.HttpError					"Unauthorized Error"
//	@Failure	417					{object}	exception.HttpError					"Validation error"
//	@Failure	500					{object}	exception.HttpError					"Reivew remark not found"
//	@Router		/projects/{uuid}/versions/{version}/reviewremarks/{reviewRemarkUuid} [post]
//	@security	Bearer
func (projectHandler *ProjectHandler) ProjectVersionReviewRemarksCommentExtern(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	currentProject, version, token := projectHandler.retrieveProjectAndVersionFromPublicRequest(requestSession, r)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	remarkUuidEscaped := chi.URLParam(r, "reviewRemarkUuid")
	remarkUuid, err := url.QueryUnescape(remarkUuidEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamRemarkUuidEmpty))
	if remarkUuid == "" {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ParamRemarkUuidEmpty))
	}
	err = validation.CheckUuid(remarkUuid)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ErrorKeyRequestParamNotValid, "reviewRemarkUuid"), zapcore.InfoLevel)
	commentData := extractCommentReviewRemarkExternBody(r, true)

	remarks := projectHandler.ReviewRemarksRepository.FindByKey(requestSession, version.Key, false)
	var remark *reviewremarks.Remark
	for _, r := range remarks.Remarks {
		if r.Key == remarkUuid {
			remark = r
			break
		}
	}
	if remark == nil {
		exception.ThrowExceptionClient404Message3(message.GetI18N(message.ErrorDbNotFound, remarkUuid))
	}

	origin := helper.CreateOrigin(project.OriginApi, token.Company, token.Description, token.Key)
	ip := jwt.TrimPortFromRemoteAddress(r.RemoteAddr)

	var before reviewremarks.Remark
	copier.Copy(&before, remark)
	remark.Comment(ip, origin, commentData.Content)
	projectHandler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, version.Key, ip, message.ReviewRemarkCommented, cmp.Diff, *remark, before)
	projectHandler.ReviewRemarksRepository.Update(requestSession, remarks)
	responseData := SuccessResponse{
		Success: true,
		Message: "comment created",
	}
	render.JSON(w, r, responseData)
}
