// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/eclipse-disuko/disuko/infra/repository/licenserules"
	"github.com/eclipse-disuko/disuko/infra/repository/policydecisions"
	"github.com/eclipse-disuko/disuko/infra/repository/user"
	projectLabelService "github.com/eclipse-disuko/disuko/infra/service/project-label"
	"go.uber.org/zap/zapcore"

	"github.com/eclipse-disuko/disuko/domain/project/components"

	"github.com/eclipse-disuko/disuko/infra/repository/auditloglist"
	"github.com/eclipse-disuko/disuko/infra/service/cache"
	sbomLockRetained "github.com/eclipse-disuko/disuko/infra/service/check-sbom-retained"
	"github.com/eclipse-disuko/disuko/infra/service/spdx"
	"github.com/eclipse-disuko/disuko/observermngmt"

	"github.com/eclipse-disuko/disuko/infra/repository/labels"

	sa "github.com/eclipse-disuko/disuko/infra/service/analytics"
	"github.com/eclipse-disuko/disuko/infra/service/locks"

	"github.com/eclipse-disuko/disuko/infra/repository/sbomlist"

	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/notices"
	"github.com/eclipse-disuko/disuko/infra/repository/policyrules"

	"github.com/eclipse-disuko/disuko/helper/validation"
	"github.com/eclipse-disuko/disuko/infra/repository/schema"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain/project"
	sbomlist2 "github.com/eclipse-disuko/disuko/domain/project/sbomlist"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/jwt"
	"github.com/eclipse-disuko/disuko/helper/roles"
	"github.com/eclipse-disuko/disuko/helper/s3Helper"
	license2 "github.com/eclipse-disuko/disuko/infra/repository/license"
	project2 "github.com/eclipse-disuko/disuko/infra/repository/project"
	"github.com/eclipse-disuko/disuko/infra/service"
	sbomlockRetained "github.com/eclipse-disuko/disuko/infra/service/check-sbom-retained"
	"github.com/eclipse-disuko/disuko/infra/service/compare"
	projectService "github.com/eclipse-disuko/disuko/infra/service/project"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type SPDXHandler struct {
	ProjectRepository         project2.IProjectRepository
	SchemaRepository          schema.ISchemaRepository
	LicensesRepository        license2.ILicensesRepository
	PolicyRuleRepository      policyrules.IPolicyRulesRepository
	SbomListRepository        sbomlist.ISbomListRepository
	AnalyticsService          sa.Analytics
	LabelRepository           labels.ILabelRepository
	LockService               *locks.Service
	AuditLogListRepository    auditloglist.IAuditLogListRepository
	LicenseRulesRepository    licenserules.ILicenseRulesRepository
	SpdxService               *spdx.Service
	SbomRetainedService       *sbomLockRetained.Service
	ProjectLabelService       *projectLabelService.ProjectLabelService
	PolicyDecisionsRepository policydecisions.IPolicyDecisionsRepository
	UserRepository            user.IUsersRepository
}

func (spdxHandler *SPDXHandler) HandleSPDXUploadFile(requestSession *logy.RequestSession, currentProject *project.Project, version *project.ProjectVersion, origin string, uploader string, w http.ResponseWriter, r *http.Request, isPublic bool) {
	validation.CheckExpectedContentType(r, validation.ContentTypeFormData)

	l, acquired := spdxHandler.LockService.Acquire(locks.Options{
		Key:      currentProject.Key,
		Blocking: true,
		Timeout:  time.Minute * 2,
	})
	if !acquired {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ResourceInUse), "")
	}
	defer spdxHandler.LockService.Release(l)

	logy.Infow(requestSession, fmt.Sprintf("%s Start HandleSPDXUploadFile", logy.MsgStageIntermediateCommon), logy.MsgStage, logy.MsgStageIntermediateCommon)
	TryNewFileUpload(requestSession, currentProject.Key, spdxHandler.ProjectRepository)

	file, handler, err := r.FormFile("file")
	if err != nil {
		// max 10mb
		err = r.ParseMultipartForm(10 << 20)
	}
	exception.HandleErrorClientMessage(err, message.GetI18N(message.SpdxFileEmptyOrLarge))
	defer func(file multipart.File) {
		if file == nil {
			exception.HandleErrorClientMessage(err, message.GetI18N(message.SpdxFileEmptyOrLarge))
			return
		}
		err := file.Close()
		if err != nil {
		}
	}(file)

	sbomTag := r.FormValue("sbomTag")

	if handler == nil {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.SpdxFileContentTypeValidation), "FileHeader object is nil")
	}
	validation.CheckExpectedContentType2(handler.Header, []validation.ContentType{
		validation.ContentTypeJson,
		validation.ContentTypeOctets,
	})

	holder := projectService.RepositoryHolder{
		ProjectRepository:  spdxHandler.ProjectRepository,
		LicenseRepository:  spdxHandler.LicensesRepository,
		SBOMListRepository: spdxHandler.SbomListRepository,
		SchemaRepository:   spdxHandler.SchemaRepository,
	}

	spdxFile, validateFailedMsg := service.UploadSbom(requestSession, currentProject,
		version.Key, origin, uploader, file, handler.Filename, sbomTag, holder, spdxHandler.SpdxService)

	if validateFailedMsg != "" {
		// send failed message
		if isPublic {
			render.Status(r, 417)
		}
		render.JSON(w, r, project.SPDXUploadResponse{
			DocIsValid:              false,
			ValidationFailedMessage: validateFailedMsg,
			FileUploaded:            false,
		})
	} else {
		isDummy := hasDummyLabel(currentProject, getDummyLabel(requestSession, spdxHandler.LabelRepository))
		if !isDummy {
			observermngmt.FireEvent(observermngmt.SpdxAdded, observermngmt.SpdxData{
				RequestSession: requestSession,
				Project:        currentProject,
				Version:        version,
				SpdxFile:       spdxFile,
			})
		}
		spdxHandler.AuditLogListRepository.AddStaticAuditEntryByKey(requestSession, version.Key, uploader, message.SpdxFileUploaded, spdxFile)

		if conf.Config.Server.AutoDeleteSbomsAfterUpload && (!conf.IsProdEnv() || conf.Config.Server.ProdAutoDeleteSbomsAfterUpload) {
			sbomList := spdxHandler.SbomListRepository.FindByKey(requestSession, version.Key, false)

			spdxFileHistory := sbomList.SpdxFileHistory
			sort.Slice(spdxFileHistory, func(i, j int) bool {
				return spdxFileHistory[i].Uploaded.UTC().After(spdxFileHistory[j].Uploaded.UTC())
			})

			spdxRemaining := make([]*project.SpdxFileBase, 0)
			unusedSpdxCount := 0
			spdxKeysForDeletion := make([]string, 0)
			spdxDeleted := make([]*project.SpdxFileBase, 0)
			for _, currentSpdx := range spdxFileHistory {
				if IsSpdxInUse(currentSpdx, currentProject, version) {
					spdxRemaining = append(spdxRemaining, currentSpdx)
				} else {
					if unusedSpdxCount < 5 {
						unusedSpdxCount++
						spdxRemaining = append(spdxRemaining, currentSpdx)
					} else {
						spdxKeysForDeletion = append(spdxKeysForDeletion, currentSpdx.Key)
						spdxDeleted = append(spdxDeleted, currentSpdx)
					}
				}
			}

			if len(spdxFileHistory) != len(spdxRemaining) && len(spdxKeysForDeletion) > 0 {
				for _, spdxKey := range spdxKeysForDeletion {
					filename := currentProject.GetFilePathSbom(spdxKey, version.Key)
					s3Helper.DeleteFile(requestSession, filename)
					cacheFilePath := fmt.Sprintf(cache.CachePath, spdxKey)
					s3Helper.DeleteFile(requestSession, cacheFilePath)
				}

				if !isDummy {
					for _, spdx := range spdxDeleted {
						observermngmt.FireEvent(observermngmt.SpdxDeleted, observermngmt.SpdxData{
							RequestSession: requestSession,
							Project:        currentProject,
							Version:        version,
							SpdxFile:       spdx,
						})
					}
				}

				sbomList.SpdxFileHistory = spdxRemaining
				spdxHandler.SbomListRepository.Update(requestSession, sbomList)
				spdxHandler.AuditLogListRepository.AddStaticAuditEntryByKey(requestSession, version.Key, uploader, message.SpdxFileDeleted, spdxDeleted)
			}
		}
		// send success message
		render.JSON(w, r, project.SPDXUploadResponse{
			DocIsValid:              true,
			ValidationFailedMessage: "", Hash: spdxFile.Hash, FileUploaded: true, Id: spdxFile.MetaInfo.SpdxId, SbomGuid: spdxFile.Key,
		})
	}
}

// SPDXUploadFileExternHandler godoc
//
//	@Summary	Upload SBOM as SPDX
//	@Id			uploadSBOMAsSPDX
//	@Produce	json
//	@Accept		mpfd
//	@Param		uuid	path		string						true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Param		version	path		string						true	"Project Version Name (also known as Channel Name) e.g.: main"
//	@Param		file	formData	file						true	"SPDX File"
//	@Param		sbomTag	formData	string						false	"SPDX Tag"
//	@Success	200		{object}	project.SPDXUploadResponse	"SPDX Upload Response"
//	@Failure	401		{object}	exception.HttpError			"Unauthorized Error"
//	@Failure	417		{object}	project.SPDXUploadResponse	"Validation Error"
//	@Router		/v1/projects/{uuid}/versions/{version}/sboms [post]
//	@security	Bearer
func (spdxHandler *SPDXHandler) SPDXUploadFileExternHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	currentProject, version, origin := spdxHandler.retrieveProjectAndVersionFromPublicRequest(requestSession, r)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	ip := jwt.TrimPortFromRemoteAddress(r.RemoteAddr)

	spdxHandler.HandleSPDXUploadFile(requestSession, currentProject, version, origin, ip, w, r, true)
}

func (spdxHandler *SPDXHandler) SPDXUploadFileHandler(w http.ResponseWriter, r *http.Request) {
	err := assertFileSize(w, r)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.MaxFilesize))

	currentProject, version, requestSession := spdxHandler.retrieveProjectAndVersion2(r)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	userName, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowSBOMAction.Upload {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.UploadSbom))
	}

	spdxHandler.HandleSPDXUploadFile(requestSession, currentProject, version, project.OriginUi, userName, w, r, false)
}

func (spdxHandler *SPDXHandler) SpdxDeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, version, requestSession := spdxHandler.retrieveProjectAndVersion2(r)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}
	user, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	isOwner := false
	for _, r := range rights.Groups {
		if r == string(project.OWNER) {
			isOwner = true
			break
		}
	}
	if !isOwner && !rights.IsDomainAdmin() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeleteSbom))
	}

	spdxFileKey := chi.URLParam(r, "spdxFileKey")
	if len(spdxFileKey) == 0 {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.SpdxFileKeyNotSet), "")
	}
	l, spdxToDelete := spdxHandler.resolveSbomListAndSpdx(requestSession, version.Key, spdxFileKey)
	isLatest := l.SpdxFileHistory.GetLatest() == spdxToDelete

	if IsSpdxInUse(spdxToDelete, currentProject, version) {
		exception.ThrowExceptionClientWithHttpCode(message.ErrorSpdxInUse, message.GetI18N(message.ErrorSpdxInUse).Text, "", exception.HTTP_CODE_SHOW_NO_REQUEST_ID)
	}

	filename := currentProject.GetFilePathSbom(spdxToDelete.Key, version.Key)
	s3Helper.DeleteFile(requestSession, filename)
	cacheFilePath := fmt.Sprintf(cache.CachePath, spdxToDelete.Key)
	s3Helper.DeleteFile(requestSession, cacheFilePath)

	newSbomHistory := make([]*project.SpdxFileBase, 0)
	for _, spdx := range l.SpdxFileHistory {
		if spdx.Key != spdxToDelete.Key {
			newSbomHistory = append(newSbomHistory, spdx)
		}
	}
	l.SpdxFileHistory = newSbomHistory
	spdxHandler.SbomListRepository.Update(requestSession, l)

	spdxHandler.AuditLogListRepository.AddStaticAuditEntryByKey(requestSession, version.Key, user, message.SpdxFileDeleted, spdxToDelete)
	dummyLabel := getDummyLabel(requestSession, spdxHandler.LabelRepository)
	if !hasDummyLabel(currentProject, dummyLabel) {
		observermngmt.FireEvent(observermngmt.SpdxDeleted, observermngmt.SpdxData{
			RequestSession: requestSession,
			Project:        currentProject,
			Version:        version,
			SpdxFile:       spdxToDelete,
		})
	}

	if len(l.SpdxFileHistory) > 0 && isLatest {
		newLatest := l.SpdxFileHistory.GetLatest()
		if !hasDummyLabel(currentProject, dummyLabel) {
			observermngmt.FireEvent(observermngmt.SpdxUpdatedNewest, observermngmt.SpdxData{
				RequestSession: requestSession,
				Project:        currentProject,
				Version:        version,
				SpdxFile:       newLatest,
			})
		}
	}
	responseData := SuccessResponse{
		Success: true,
		Message: "Spdx file deleted",
	}
	render.JSON(w, r, responseData)
}

func assertFileSize(w http.ResponseWriter, r *http.Request) error {
	maxBytes := mbToBytes(conf.Config.Server.UploadMaxMb)
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
	if err := r.ParseMultipartForm(maxBytes); err != nil {
		return err
	}
	return nil
}

func mbToBytes(mb int64) int64 {
	return mb << 20
}

func (spdxHandler *SPDXHandler) DownloadSPDXHistoryFileHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, version, requestSession := spdxHandler.retrieveProjectAndVersion2(r)

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowSBOMAction.Download {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DownloadSbomHistory))
	}

	spdxFileKey := chi.URLParam(r, "spdxFileKey")
	if len(spdxFileKey) == 0 {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.SpdxFileKeyNotSet), "")
	}
	_, spdxFileFound := spdxHandler.retrieveSbomListAndFile(requestSession, version.Key, spdxFileKey)
	if spdxFileFound == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.SpdxFileNotFound, spdxFileKey), "")
	}
	s3Helper.PerformDownload(requestSession, &w, currentProject.GetFilePathSbom(spdxFileFound.Key, version.Key), spdxFileFound.Hash)
}

func extractUpdateTagReq(r *http.Request, handleErrorAsServerException bool) (body project.SPDXSetTagRequestDto) {
	validation.DecodeAndValidate(r, &body, handleErrorAsServerException)
	return
}

// PublicSpdxStatusHandler godoc
//
//	@Summary	Update SPDX Tag
//	@Id			getSpdxStatus
//	@Produce	json
//	@Accept		json
//	@Param		uuid		path		string								true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Param		version		path		string								true	"Project Version Name (also known as Channel Name) e.g.: main"
//	@Param		sbomUuid	path		string								true	"UUID of the SBOM delivery or 'latest' for the latest SBOM delivery e.g.: dummy-sbom-id---xxx-4413-yyy-24f060311111"
//	@Success	200			{object}	project.SbomStatusPublicResponseDto	"Overall Policy Status Response for SBOM"
//	@Failure	401			{object}	exception.HttpError					"Unauthorized Error"
//	@Failure	417			{object}	exception.HttpError					"Validation Error"
//	@Failure	500			{object}	exception.HttpError					"SPDX not found in history"
//	@Router		/v1/projects/{uuid}/versions/{version}/sboms/{sbomUuid}/status [get]
//	@security	Bearer
func (spdxHandler *SPDXHandler) PublicSpdxStatusHandler(w http.ResponseWriter, r *http.Request) {
	rs := logy.GetRequestSession(r)
	currentProject, version, _ := spdxHandler.retrieveProjectAndVersionFromPublicRequest(rs, r)

	sbomUuidEscaped := chi.URLParam(r, "sbomUuid")
	sbomUuid := ValidateIDOrLatest(sbomUuidEscaped)

	_, spdx := spdxHandler.resolveSbomListAndSpdx(rs, version.Key, sbomUuid)

	compInfos := spdxHandler.SpdxService.GetComponentInfos(rs, currentProject, version.Key, spdx)

	rules := spdxHandler.PolicyRuleRepository.FindPolicyRulesForLabel(rs, currentProject.PolicyLabels)
	policyDecisions := spdxHandler.PolicyDecisionsRepository.FindByKey(rs, currentProject.Key, false)
	isVehicle := spdxHandler.ProjectLabelService.HasVehiclePlatformLabel(rs, currentProject)

	policyEvaluation := compInfos.EvaluatePolicyRules(rules, policyDecisions, isVehicle, spdx.Uploaded, spdx.Key)

	var status string
	switch {
	case policyEvaluation.Stats.Denied > 0 || policyEvaluation.Stats.NoAssertion > 0:
		status = "red"
	case policyEvaluation.Stats.Warned > 0:
		status = "yellow"
	default:
		status = "green"
	}
	response := project.SbomStatusPublicResponseDto{Status: status}

	render.JSON(w, r, response)
}

// PublicSpdxTagUpdateHandler godoc
//
//	@Summary	Update SPDX Tag
//	@Id			updateSpdxTag
//	@Produce	json
//	@Accept		json
//	@Param		uuid		path		string							true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Param		version		path		string							true	"Project Version Name (also known as Channel Name) e.g.: main"
//	@Param		sbomUuid	path		string							true	"UUID of the SBOM delivery or 'latest' for the latest SBOM delivery e.g.: dummy-sbom-id---xxx-4413-yyy-24f060311111"
//	@Param		tag			body		project.SPDXSetTagRequestDto	true	"Tag"
//	@Success	200			{object}	rest.SuccessResponse			"Success Response"
//	@Failure	401			{object}	exception.HttpError				"Unauthorized Error"
//	@Failure	417			{object}	exception.HttpError				"Validation Error"
//	@Failure	500			{object}	exception.HttpError				"SPDX not found in history"
//	@Router		/v1/projects/{uuid}/versions/{version}/sboms/{sbomUuid}/tag [put]
//	@security	Bearer
func (spdxHandler *SPDXHandler) PublicSpdxTagUpdateHandler(w http.ResponseWriter, r *http.Request) {
	rs := logy.GetRequestSession(r)
	currentProject, version, _ := spdxHandler.retrieveProjectAndVersionFromPublicRequest(rs, r)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	body := extractUpdateTagReq(r, true)

	sbomUuidEscaped := chi.URLParam(r, "sbomUuid")
	sbomUuid := ValidateIDOrLatest(sbomUuidEscaped)

	spdxHandler.updateSpdxTag(rs, body.Tag, version.Key, sbomUuid)

	responseData := SuccessResponse{
		Success: true,
		Message: "Spdx tag updated",
	}
	render.JSON(w, r, responseData)
}

func (spdxHandler *SPDXHandler) SpdxTagUpdateHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, version, requestSession := spdxHandler.retrieveProjectAndVersion2(r)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProjectVersion.Update {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.UpdateSbom))
	}

	spdxFileKey := chi.URLParam(r, "spdxFileKey")
	if len(spdxFileKey) == 0 {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.SpdxFileKeyNotSet), "")
	}

	body := extractUpdateTagReq(r, true)

	spdxHandler.updateSpdxTag(requestSession, body.Tag, version.Key, spdxFileKey)

	responseData := SuccessResponse{
		Success: true,
		Message: "Spdx tag updated",
	}
	render.JSON(w, r, responseData)
}

// PublicSpdxLockHandler godoc
//
//	@Summary	Lock SPDX
//	@Id			lockSpdx
//	@Produce	json
//	@Accept		json
//	@Param		uuid		path		string					true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Param		version		path		string					true	"Project Version Name (also known as Channel Name) e.g.: main"
//	@Param		sbomUuid	path		string					true	"UUID of the SBOM delivery or 'latest' for the latest SBOM delivery e.g.: dummy-sbom-id---xxx-4413-yyy-24f060311111"
//	@Success	200			{object}	rest.SuccessResponse	"Success Response"
//	@Failure	401			{object}	exception.HttpError		"Unauthorized Error"
//	@Failure	417			{object}	exception.HttpError		"Validation Error"
//	@Failure	500			{object}	exception.HttpError		"SPDX not found in history"
//	@Router		/v1/projects/{uuid}/versions/{version}/sboms/{sbomUuid}/lock [put]
//	@security	Bearer
func (spdxHandler *SPDXHandler) PublicSpdxLockHandler(w http.ResponseWriter, r *http.Request) {
	rs := logy.GetRequestSession(r)
	currentProject, version, _ := spdxHandler.retrieveProjectAndVersionFromPublicRequest(rs, r)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	sbomUuidEscaped := chi.URLParam(r, "sbomUuid")
	sbomUuid := ValidateIDOrLatest(sbomUuidEscaped)

	l, spdx := spdxHandler.resolveSbomListAndSpdx(rs, version.Key, sbomUuid)

	if spdx.IsLocked {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.SpdxAlreadyLocked))
	} else {
		spdx.IsLocked = true
		spdxHandler.SbomListRepository.Update(rs, l)
		spdxHandler.AuditLogListRepository.AddStaticAuditEntryByKey(rs, version.Key, project.OriginApi, message.SpdxFileLocked, spdx)
		spdxHandler.markProjectSbomRetainFlag(rs, currentProject)
		render.JSON(w, r, SuccessResponse{
			Success: true,
			Message: "Spdx locked",
		})
	}
}

func (spdxHandler *SPDXHandler) markProjectSbomRetainFlag(requestSession *logy.RequestSession, prj *project.Project) {
	if !prj.HasSBOMToRetain {
		prj.HasSBOMToRetain = true
		spdxHandler.ProjectRepository.Update(requestSession, prj)
	}
}

// @Id			unlockSpdx
// @Produce	json
// @Accept		json
// @Param		uuid		path		string					true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
// @Param		version		path		string					true	"Project Version Name (also known as Channel Name) e.g.: main"
// @Param		sbomUuid	path		string					true	"UUID of the SBOM delivery or 'latest' for the latest SBOM delivery e.g.: dummy-sbom-id---xxx-4413-yyy-24f060311111"
// @Success	200			{object}	rest.SuccessResponse	"Success Response"
// @Failure	401			{object}	exception.HttpError		"Unauthorized Error"
// @Failure	417			{object}	exception.HttpError		"Validation Error"
// @Failure	500			{object}	exception.HttpError		"SPDX not found in history"
// @Router		/v1/projects/{uuid}/versions/{version}/sboms/{sbomUuid}/unlock [put]
// @security	Bearer
func (spdxHandler *SPDXHandler) PublicSpdxUnlockHandler(w http.ResponseWriter, r *http.Request) {
	rs := logy.GetRequestSession(r)
	currentProject, version, _ := spdxHandler.retrieveProjectAndVersionFromPublicRequest(rs, r)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	sbomUuidEscaped := chi.URLParam(r, "sbomUuid")
	sbomUuid := ValidateIDOrLatest(sbomUuidEscaped)

	l, spdx := spdxHandler.resolveSbomListAndSpdx(rs, version.Key, sbomUuid)

	if !spdx.IsLocked {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.SpdxNotLocked))
	} else if sbomlockRetained.IsSpdxToRetain(spdx, version) {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.SpdxRetainedForApprovalOrReview))
	} else {
		spdx.IsLocked = false
		spdxHandler.SbomListRepository.Update(rs, l)
		spdxHandler.AuditLogListRepository.AddStaticAuditEntryByKey(rs, version.Key, project.OriginApi, message.SpdxFileUnlocked, spdx)
		// Check if there are still any retained SBOMs(due to review or higher priority ), if not set HasSBOMToRetain to false
		if !sbomlockRetained.HasAnyVersionWithRetainedSbom(rs, spdxHandler.ProjectRepository, spdxHandler.SbomListRepository, currentProject) {
			currentProject.HasSBOMToRetain = false
			spdxHandler.ProjectRepository.Update(rs, currentProject)
		}
		render.JSON(w, r, SuccessResponse{
			Success: true,
			Message: "Spdx unlocked",
		})
	}
}

func (spdxHandler *SPDXHandler) SpdxToggleLockHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := spdxHandler.retrieveProject2(r, true)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}
	user, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	isOwner := false
	for _, r := range rights.Groups {
		if r == string(project.OWNER) {
			isOwner = true
			break
		}
	}
	if !isOwner && !rights.IsDomainAdmin() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.LockSbom))
	}

	versionKey := extractVersionKeyFromRequest(r)
	spdxFileKey := chi.URLParam(r, "spdxFileKey")
	if len(spdxFileKey) == 0 {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.SpdxFileKeyNotSet), "")
	}

	l, spdx := spdxHandler.resolveSbomListAndSpdx(requestSession, versionKey, spdxFileKey)
	spdx.IsLocked = !spdx.IsLocked
	spdxHandler.SbomListRepository.Update(requestSession, l)

	if spdx.IsLocked {
		currentProject.HasSBOMToRetain = true
		spdxHandler.ProjectRepository.Update(requestSession, currentProject)
	}

	// Check if there are still any retained SBOMs(due to review or higher priority ), if not set HasSBOMToRetain to false
	if !spdx.IsLocked && !sbomlockRetained.HasAnyVersionWithRetainedSbom(requestSession, spdxHandler.ProjectRepository, spdxHandler.SbomListRepository, currentProject) {
		currentProject.HasSBOMToRetain = false
		spdxHandler.ProjectRepository.Update(requestSession, currentProject)
	}
	if spdx.IsLocked {
		spdxHandler.AuditLogListRepository.AddStaticAuditEntryByKey(requestSession, versionKey, user, message.SpdxFileLocked, spdx)
	} else {
		spdxHandler.AuditLogListRepository.AddStaticAuditEntryByKey(requestSession, versionKey, user, message.SpdxFileUnlocked, spdx)
	}
	responseData := SuccessResponse{
		Success: true,
		Message: "Spdx lock toggled",
	}
	render.JSON(w, r, responseData)
}

func (spdxHandler *SPDXHandler) resolveSbomListAndSpdx(requestSession *logy.RequestSession, versionKey, sbomUuid string) (l *sbomlist2.SbomList, spdx *project.SpdxFileBase) {
	if sbomUuid == "latest" {
		l, spdx = spdxHandler.retrieveSbomListAndLatestFile(requestSession, versionKey)
	} else {
		l, spdx = spdxHandler.retrieveSbomListAndFile(requestSession, versionKey, sbomUuid)
	}
	if spdx == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.FindingSbomKey), "SPDX not found in history: "+sbomUuid)
	}
	return l, spdx
}

func (spdxHandler *SPDXHandler) updateSpdxTag(requestSession *logy.RequestSession, tag, versionKey, sbomUuid string) {
	l, spdx := spdxHandler.resolveSbomListAndSpdx(requestSession, versionKey, sbomUuid)
	spdx.Tag = tag
	spdxHandler.SbomListRepository.Update(requestSession, l)
}

func (spdxHandler *SPDXHandler) ExportNoticeFileForSbomAsTextHandler(w http.ResponseWriter, r *http.Request) {
	requestSession, currentProject, compInfos, contactMeta := spdxHandler.prepareExportNoticeFileForSbom(r)
	if len(compInfos) == 0 {
		w.WriteHeader(200)
		return
	}
	sb := notices.GenerateTextNotices(requestSession, *currentProject, spdxHandler.LicensesRepository, compInfos, contactMeta, spdxHandler.LabelRepository)

	_, err := w.Write([]byte(sb.String()))
	exception.HandleErrorClientMessage(err, message.GetI18N(message.WritingContent))

	w.WriteHeader(200)
}

// ExportTextNoticeExtern godoc
//
//	@Summary	Get notice file for specified SBOM formatted as text
//	@Id			getSBOMNoticeFileText
//	@Produce	plain
//	@Param		uuid		path		string					true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Param		version		path		string					true	"Project Version Name (also known as Channel Name) e.g.: main"
//	@Param		sbomUuid	path		string					true	"UUID of the SBOM delivery or 'latest' for the latest SBOM delivery e.g.: dummy-sbom-id---xxx-4413-yyy-24f060311111"
//	@Success	200			{string}	string					"Notice File"
//	@Failure	404			{object}	exception.HttpError404	"NotFound Error"
//	@Failure	401			{object}	exception.HttpError		"Unauthorized Error"
//	@Router		/v1/projects/{uuid}/versions/{version}/sboms/{sbomUuid}/notice/text [get]
//	@security	Bearer
func (spdxHandler *SPDXHandler) ExportTextNoticeExtern(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	currentProject, version, _ := spdxHandler.retrieveProjectAndVersionFromPublicRequest(requestSession, r)
	sbomUuidEscaped := chi.URLParam(r, "sbomUuid")
	sbomUuid := ValidateIDOrLatest(sbomUuidEscaped)

	var spdx *project.SpdxFileBase
	if sbomUuid == "latest" {
		_, spdx = spdxHandler.retrieveSbomListAndLatestFile(requestSession, version.Key)
	} else {
		_, spdx = spdxHandler.retrieveSbomListAndFile(requestSession, version.Key, sbomUuid)
	}
	if spdx == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.SpdxFileNotFound, sbomUuid), "")
	}
	compInfos := spdxHandler.SpdxService.GetComponentInfos(requestSession, currentProject, version.Key, spdx)
	if len(compInfos) == 0 {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.WarnNotExists), "Currently notice file not exists")
	}

	contactMeta := getContactMetaOfGroupOrProject(requestSession, spdxHandler.ProjectRepository, currentProject)

	sb := notices.GenerateTextNotices(requestSession, *currentProject, spdxHandler.LicensesRepository, compInfos, contactMeta, spdxHandler.LabelRepository)
	_, err := w.Write([]byte(sb.String()))
	exception.HandleErrorClientMessage(err, message.GetI18N(message.WritingContent))
	w.WriteHeader(200)
}

// ExportHTMLNoticeExtern godoc
//
//	@Summary	Get notice file for specified SBOM formatted as HTML
//	@Id			getSBOMNoticeFileHTML
//	@Produce	html
//	@Param		uuid		path		string					true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Param		version		path		string					true	"Project Version Name (also known as Channel Name) e.g.: main"
//	@Param		sbomUuid	path		string					true	"UUID of the SBOM delivery or 'latest' for the latest SBOM delivery e.g.: dummy-sbom-id---xxx-4413-yyy-24f060311111"
//	@Success	200			{string}	string					"Notice File"
//	@Failure	404			{object}	exception.HttpError404	"NotFound Error"
//	@Failure	401			{object}	exception.HttpError		"Unauthorized Error"
//	@Router		/v1/projects/{uuid}/versions/{version}/sboms/{sbomUuid}/notice/html [get]
//	@security	Bearer
func (spdxHandler *SPDXHandler) ExportHTMLNoticeExtern(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	currentProject, version, _ := spdxHandler.retrieveProjectAndVersionFromPublicRequest(requestSession, r)

	sbomUuidEscaped := chi.URLParam(r, "sbomUuid")
	sbomUuid := ValidateIDOrLatest(sbomUuidEscaped)
	var spdx *project.SpdxFileBase
	if sbomUuid == "latest" {
		_, spdx = spdxHandler.retrieveSbomListAndLatestFile(requestSession, version.Key)
	} else {
		_, spdx = spdxHandler.retrieveSbomListAndFile(requestSession, version.Key, sbomUuid)
	}
	if spdx == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.SpdxFileNotFound, sbomUuid), "")
	}
	compInfos := spdxHandler.SpdxService.GetComponentInfos(requestSession, currentProject, version.Key, spdx)
	if len(compInfos) == 0 {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.WarnNotExists), "Currently notice file not exists")
	}

	contactMeta := getContactMetaOfGroupOrProject(requestSession, spdxHandler.ProjectRepository, currentProject)

	sb := notices.GenerateHTMLNotices(requestSession, *currentProject, spdxHandler.LicensesRepository, compInfos, contactMeta, spdxHandler.LabelRepository)

	_, err := w.Write([]byte(sb.String()))
	exception.HandleErrorClientMessage(err, message.GetI18N(message.WritingContent))

	w.WriteHeader(200)
}

// ExportJSONNoticeExtern godoc
//
//	@Summary	Get notice file for specified SBOM formatted as JSON
//	@Id			getSBOMNoticeFileJSON
//	@Produce	json
//	@Param		uuid		path		string					true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Param		version		path		string					true	"Project Version Name (also known as Channel Name) e.g.: main"
//	@Param		sbomUuid	path		string					true	"UUID of the SBOM delivery or 'latest' for the latest SBOM delivery e.g.: dummy-sbom-id---xxx-4413-yyy-24f060311111"
//	@Success	200			{object}	project.NoticeFileJSON	"Notice File"
//	@Failure	404			{object}	exception.HttpError404	"NotFound Error"
//	@Failure	401			{object}	exception.HttpError		"Unauthorized Error"
//	@Router		/v1/projects/{uuid}/versions/{version}/sboms/{sbomUuid}/notice/json [get]
//	@security	Bearer
func (spdxHandler *SPDXHandler) ExportJSONNoticeExtern(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	currentProject, version, _ := spdxHandler.retrieveProjectAndVersionFromPublicRequest(requestSession, r)
	sbomUuidEscaped := chi.URLParam(r, "sbomUuid")
	sbomUuid := ValidateIDOrLatest(sbomUuidEscaped)

	var spdx *project.SpdxFileBase
	if sbomUuid == "latest" {
		_, spdx = spdxHandler.retrieveSbomListAndLatestFile(requestSession, version.Key)
	} else {
		_, spdx = spdxHandler.retrieveSbomListAndFile(requestSession, version.Key, sbomUuid)
	}
	if spdx == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.SpdxFileNotFound, sbomUuid), "")
	}
	compInfos := spdxHandler.SpdxService.GetComponentInfos(requestSession, currentProject, version.Key, spdx)
	if len(compInfos) == 0 {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.WarnNotExists), "Currently notice file not exists")
	}

	contactMeta := getContactMetaOfGroupOrProject(requestSession, spdxHandler.ProjectRepository, currentProject)

	notice := notices.GenerateJSONNotices(requestSession, *currentProject, spdxHandler.LicensesRepository, compInfos, contactMeta, spdxHandler.LabelRepository)

	render.JSON(w, r, notice)
}

func (spdxHandler SPDXHandler) prepareExportNoticeFileForSbom(r *http.Request) (*logy.RequestSession, *project.Project, components.ComponentInfos, project.NoticeContactMeta) {
	sbomUuidEscaped := chi.URLParam(r, "sbomUuid")
	sbomUuid, err := url.QueryUnescape(sbomUuidEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamSbomUuidEmpty))
	if sbomUuid == "" {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ParamSbomUuidEmpty))
	}
	err = validation.CheckUuid(sbomUuid)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorKeyRequestParamNotValid, "sbomUuid"), zapcore.InfoLevel)

	currentProject, version, requestSession := spdxHandler.retrieveProjectAndVersion2(r)

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProjectVersion.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ViewNoticeFile))
	}

	_, selectedSpdx := spdxHandler.retrieveSbomListAndFile(requestSession, version.Key, sbomUuid)
	if selectedSpdx == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.NoSbom), "")
	}
	compInfos := spdxHandler.SpdxService.GetComponentInfos(requestSession, currentProject, version.Key, selectedSpdx)

	contactMeta := getContactMetaOfGroupOrProject(requestSession, spdxHandler.ProjectRepository, currentProject)
	return requestSession, currentProject, compInfos, contactMeta
}

func (spdxHandler *SPDXHandler) ExportNoticeFileForSbomAsJSONHandler(w http.ResponseWriter, r *http.Request) {
	requestSession, currentProject, compInfos, contactMeta := spdxHandler.prepareExportNoticeFileForSbom(r)
	if len(compInfos) == 0 {
		w.WriteHeader(200)
		return
	}
	notice := notices.GenerateJSONNotices(requestSession, *currentProject, spdxHandler.LicensesRepository, compInfos, contactMeta, spdxHandler.LabelRepository)
	render.JSON(w, r, notice)
}

func (spdxHandler *SPDXHandler) ExportNoticeFileForSbomAsHTMLHandler(w http.ResponseWriter, r *http.Request) {
	requestSession, currentProject, compInfos, contactMeta := spdxHandler.prepareExportNoticeFileForSbom(r)
	if len(compInfos) == 0 {
		w.WriteHeader(200)
		return
	}
	sb := notices.GenerateHTMLNotices(requestSession, *currentProject, spdxHandler.LicensesRepository, compInfos, contactMeta, spdxHandler.LabelRepository)

	_, err := w.Write([]byte(sb.String()))
	exception.HandleErrorClientMessage(err, message.GetI18N(message.WritingContent))

	w.WriteHeader(200)
}

func checkInputAndResolveParameterValue(r *http.Request, paramName string, i18n message.I18N) string {
	paramValueEscaped := chi.URLParam(r, paramName)
	paramValue, err := url.QueryUnescape(paramValueEscaped)
	exception.HandleErrorServerMessage(err, i18n)
	return paramValue
}

func (spdxHandler *SPDXHandler) SPDXCompareHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	checkInputAndResolveParameterValue(r, "uuid", message.GetI18N(message.ParamProjectUuidEmpty, "uuid"))
	versionUuidOld := checkInputAndResolveParameterValue(r, "versionOld", message.GetI18N(message.ParamVersionOldEmpty, "versionOld"))
	spdxOldID := checkInputAndResolveParameterValue(r, "spdxOld", message.GetI18N(message.ParamSpdxOldWrong, "spdxOld"))
	versionUuidNew := checkInputAndResolveParameterValue(r, "versionNew", message.GetI18N(message.ParamVersionNewEmpty, "versionNew"))
	spdxNewID := checkInputAndResolveParameterValue(r, "spdxNew", message.GetI18N(message.ParamSpdxNewWrong, "spdxNew"))

	currentProject, _ := spdxHandler.retrieveProject2(r, true)

	userID, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProjectVersion.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadProject))
	}

	rules := spdxHandler.PolicyRuleRepository.FindPolicyRulesForLabel(requestSession, currentProject.PolicyLabels)
	policyDecisions := spdxHandler.PolicyDecisionsRepository.FindByKey(requestSession, currentProject.Key, false)
	isVehicle := spdxHandler.ProjectLabelService.HasVehiclePlatformLabel(requestSession, currentProject)

	versionOld := currentProject.GetVersion(versionUuidOld)
	_, spdxOld := spdxHandler.retrieveSbomListAndFile(requestSession, versionOld.Key, spdxOldID)
	if spdxOld == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.SpdxFileNotFound, spdxOld), "")
	}
	compsOld := spdxHandler.SpdxService.GetComponentInfos(requestSession, currentProject, versionOld.Key, spdxOld)
	evalOld := compsOld.EvaluatePolicyRules(rules, policyDecisions, isVehicle, spdxOld.Uploaded, spdxOld.Key)

	versionNew := currentProject.GetVersion(versionUuidNew)
	_, spdxNew := spdxHandler.retrieveSbomListAndFile(requestSession, versionNew.Key, spdxNewID)
	if spdxNew == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.SpdxFileNotFound, spdxOld), "")
	}
	compsNew := spdxHandler.SpdxService.GetComponentInfos(requestSession, currentProject, versionNew.Key, spdxNew)
	evalNew := compsNew.EvaluatePolicyRules(rules, policyDecisions, isVehicle, spdxNew.Uploaded, spdxNew.Key)

	diffResult := compare.MultiCompareSpdxFiles(evalOld, evalNew, currentProject.IsResponsible(userID))

	render.JSON(w, r, diffResult)
}

func getContactMetaOfGroupOrProject(requestSession *logy.RequestSession, projectRepo project2.IProjectRepository, p *project.Project) project.NoticeContactMeta {
	contactMeta := p.NoticeContactMeta
	if len(p.Parent) > 0 {
		parentProject := projectRepo.FindByKey(requestSession, p.Parent, false)
		contactMeta = parentProject.NoticeContactMeta
	}

	return contactMeta
}

func IsSpdxInUse(spdx *project.SpdxFileBase, prj *project.Project, version *project.ProjectVersion) bool {
	spdxIsInUse := spdx.Key == prj.ApprovableSPDX.SpdxKey ||
		sbomlockRetained.AnyOverallReviewMatches(spdx.Key, version.OverallReviews) ||
		spdx.ApprovalInfo.IsInApproval ||
		spdx.IsLocked ||
		spdx.IsInUse
	return spdxIsInUse
}
