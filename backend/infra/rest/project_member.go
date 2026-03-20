// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"time"

	"github.com/eclipse-disuko/disuko/domain/oauth"
	"github.com/eclipse-disuko/disuko/helper/audit"
	"github.com/eclipse-disuko/disuko/infra/service/locks"
	"go.uber.org/zap/zapcore"

	"github.com/eclipse-disuko/disuko/observermngmt"

	"github.com/eclipse-disuko/disuko/logy"

	audit2 "github.com/eclipse-disuko/disuko/domain/audit"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/user"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/roles"
	"github.com/eclipse-disuko/disuko/helper/validation"
	userRepo "github.com/eclipse-disuko/disuko/infra/repository/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/go-cmp/cmp"
	"github.com/jinzhu/copier"
)

func (projectHandler *ProjectHandler) ProjectChildrenMemberGetAllHandler(w http.ResponseWriter, r *http.Request) {
	currentProject, requestSession := projectHandler.retrieveProject2(r, false)

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, true)
	if !rights.AllowProject.Read {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadProject))
	}

	var result project.ProjectChildrenMembersDto
	result.List = make([]project.ProjectChildMemberCombiDto, 0)

	for _, rootMember := range currentProject.UserManagement.Users {
		memberDto := rootMember.ToDto(requestSession, projectHandler.UserRepository)
		projectChildMemberCombiDto := &project.ProjectChildMemberCombiDto{
			ProjectKey:    currentProject.Key,
			ProjectName:   currentProject.Name,
			ProjectMember: &memberDto,
			UserManagementRights: &oauth.CRUDRights{
				Create: rights.AllowProjectUserManagement.Create || rights.AllowAllProjectUserManagement.Create,
				Read:   rights.AllowProjectUserManagement.Read || rights.AllowAllProjectUserManagement.Read,
				Update: rights.AllowProjectUserManagement.Update || rights.AllowAllProjectUserManagement.Update,
				Delete: rights.AllowProjectUserManagement.Delete || rights.AllowAllProjectUserManagement.Delete,
			},
		}
		result.List = append(result.List, *projectChildMemberCombiDto)
	}

	for _, childKey := range currentProject.Children {

		childProject := projectHandler.ProjectRepository.FindByKey(requestSession, childKey, true)
		if childProject == nil {
			logy.Warnf(requestSession, "Child project not found uuid: %s parent: %s", childKey, currentProject.Key)
			continue
		}

		var userManagementRights *oauth.CRUDRights
		exception.TryCatch(func() {
			_, childProjectRights := roles.GetAndCheckProjectRights(requestSession, r, childProject, false)
			userManagementRights = &oauth.CRUDRights{
				Create: childProjectRights.AllowProjectUserManagement.Create || childProjectRights.AllowAllProjectUserManagement.Create,
				Read:   childProjectRights.AllowProjectUserManagement.Read || childProjectRights.AllowAllProjectUserManagement.Read,
				Update: childProjectRights.AllowProjectUserManagement.Update || childProjectRights.AllowAllProjectUserManagement.Update,
				Delete: childProjectRights.AllowProjectUserManagement.Delete || childProjectRights.AllowAllProjectUserManagement.Delete,
			}
		}, func(exc exception.Exception) {
			if exc.ErrorCode == message.ErrorAAR {
				userManagementRights = &oauth.CRUDRights{
					Create: false,
					Read:   false,
					Update: false,
					Delete: false,
				}
			} else {
				exception.ThrowException(exc)
			}
		})

		if childProject.UserManagement.Users != nil && len(childProject.UserManagement.Users) > 0 {
			for _, childMember := range childProject.UserManagement.Users {
				memberDto := childMember.ToDto(requestSession, projectHandler.UserRepository)
				projectChildMemberCombiDto := &project.ProjectChildMemberCombiDto{
					ProjectKey:           childProject.Key,
					ProjectName:          childProject.Name,
					ProjectMember:        &memberDto,
					UserManagementRights: userManagementRights,
				}
				result.List = append(result.List, *projectChildMemberCombiDto)
			}
		} else {
			projectChildMemberCombiDto := &project.ProjectChildMemberCombiDto{
				ProjectKey:           childProject.Key,
				ProjectName:          childProject.Name,
				UserManagementRights: userManagementRights,
			}
			result.List = append(result.List, *projectChildMemberCombiDto)
		}
	}

	render.JSON(w, r, result)
}

func (projectHandler *ProjectHandler) ProjectMemberGetAllHandler(w http.ResponseWriter, r *http.Request) {
	Val.UrlParamUuid(r, "uuid")
	currentProject, requestSession := projectHandler.retrieveProject2(r, false)

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !(rights.AllowProjectUserManagement.Read || rights.AllowAllProjectUserManagement.Read) {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ViewUsers))
	}
	render.JSON(w, r, currentProject.UserManagement.ToDto(requestSession, projectHandler.UserRepository))
}

func (projectHandler *ProjectHandler) ProjectChildrenMemberAddHandler(w http.ResponseWriter, r *http.Request) {
	Val.UrlParamUuid(r, "uuid")
	currentProject, requestSession := projectHandler.retrieveProject2(r, false)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	username, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !(rights.AllowProjectUserManagement.Create || rights.AllowAllProjectUserManagement.Create) {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorAddProjectUser))
	}
	multiProjectUserData := extractMultiProjectMemberRequestBody(r)
	foundUser := projectHandler.UserRepository.FindByUserId(requestSession, multiProjectUserData.TargetUser)

	if !foundUser.IsInternal && multiProjectUserData.UserType == project.OWNER {
		exception.ThrowExceptionClient400Message(message.GetI18N(message.OnlyInternalUsersOwners), "Only internal users can be owners")
	}

	response := make([]project.ProjectChildrenMemberSuccessResponseDto, 0)
	for _, targetProjectKey := range multiProjectUserData.TargetProjects {
		l, acquired := projectHandler.LockService.Acquire(locks.Options{
			Key:      targetProjectKey,
			Blocking: true,
			Timeout:  10 * time.Second,
		})
		if !acquired {
			exception.ThrowExceptionServerMessage(message.GetI18N(message.ResourceInUse), "")
			return
		}

		logy.Infof(requestSession, "Acquired!")
		defer func() {
			projectHandler.LockService.Release(l)
			logy.Infof(requestSession, "Released lock")
		}()
		if !slices.Contains(currentProject.Children, targetProjectKey) && currentProject.Key != targetProjectKey {
			continue
		}

		targetProject := projectHandler.ProjectRepository.FindByKey(requestSession, targetProjectKey, false)
		if targetProject == nil {
			logy.Warnf(requestSession, "Child project not found uuid: %s parent: %s", targetProjectKey, currentProject.Key)
			continue
		}

		if targetProject.IsDeprecated() {
			logy.Warnf(requestSession, "Project '%s' (%s) skipped: Project is deprecated", targetProject.Name, targetProject.Key)
			response = append(response, project.ProjectChildrenMemberSuccessResponseDto{
				Success:     false,
				UserId:      multiProjectUserData.TargetUser,
				ProjectKey:  targetProject.Key,
				ProjectName: targetProject.Name,
				Message:     message.DeprecatedProjectError,
			})
			continue
		}

		shouldContinueToTheNext := false
		exception.TryCatch(func() {
			_, targetProjectRights := roles.GetAndCheckProjectRights(requestSession, r, targetProject, false)
			if !(targetProjectRights.AllowProjectUserManagement.Create || targetProjectRights.AllowAllProjectUserManagement.Create) {
				shouldContinueToTheNext = true
			}
		}, func(exc exception.Exception) {
			if exc.ErrorCode == message.ErrorAAR {
				shouldContinueToTheNext = true
			} else {
				exception.ThrowException(exc)
			}
		})
		if shouldContinueToTheNext {
			logy.Warnf(requestSession, "Project '%s' (%s) skipped: User is not authorized to add user for this project", targetProject.Name, targetProject.Key)
			response = append(response, project.ProjectChildrenMemberSuccessResponseDto{
				Success:     false,
				UserId:      multiProjectUserData.TargetUser,
				ProjectKey:  targetProject.Key,
				ProjectName: targetProject.Name,
				Message:     message.ErrorAddProjectUser,
			})
			continue
		}

		if targetProject.CheckIfUserAlreadyExistsSoft(multiProjectUserData.TargetUser) {
			logy.Warnf(requestSession, "Project '%s' (%s) skipped: Error project member already exists: '%s'", targetProject.Name, targetProject.Key, multiProjectUserData.TargetUser)
			response = append(response, project.ProjectChildrenMemberSuccessResponseDto{
				Success:     false,
				UserId:      multiProjectUserData.TargetUser,
				ProjectKey:  targetProject.Key,
				ProjectName: targetProject.Name,
				Message:     message.ErrorProjectMemberAlreadyExist,
			})
			continue
		}

		var before project.Project
		copier.Copy(&before, targetProject)

		auditEntries := make([]*audit2.Audit, 0)
		if multiProjectUserData.IsResponsible {
			if multiProjectUserData.UserType != project.OWNER {
				logy.Warnf(requestSession, "Project '%s' (%s) skipped: Only owner users can be flagged as project responsibles", targetProject.Name, targetProject.Key)
				response = append(response, project.ProjectChildrenMemberSuccessResponseDto{
					Success:     false,
					UserId:      multiProjectUserData.TargetUser,
					ProjectKey:  targetProject.Key,
					ProjectName: targetProject.Name,
					Message:     message.NonOwnerResponsible,
				})
				continue
			}
			if mc := targetProject.ProjectResponsible(); mc != nil {
				before := *mc
				mc.IsResponsible = false
				beforeAudit := before.ToAudit(requestSession, projectHandler.UserRepository)
				afterAudit := mc.ToAudit(requestSession, projectHandler.UserRepository)
				auditEntries = append(auditEntries, audit.CreateAuditEntry(username, message.ProjectUserUpdated, cmp.Diff, afterAudit, beforeAudit))
			}
		}
		userData := multiProjectUserData.ToUserData()
		targetProject.AddUser(*userData)
		entity := targetProject.GetMember(multiProjectUserData.TargetUser)
		if entity == nil {
			logy.Warnf(requestSession, "Project '%s' (%s) skipped: User '%s' not found in DB", targetProject.Name, targetProject.Key, multiProjectUserData.TargetUser)
			response = append(response, project.ProjectChildrenMemberSuccessResponseDto{
				Success:     false,
				UserId:      multiProjectUserData.TargetUser,
				ProjectKey:  targetProject.Key,
				ProjectName: targetProject.Name,
				Message:     message.ErrorUserNotFound,
			})
			continue
		}
		userAudit := entity.ToAudit(requestSession, projectHandler.UserRepository)
		auditEntries = append(auditEntries, audit.CreateAuditEntry(username, message.ProjectUserCreated, cmp.Diff, userAudit, project.ProjectMemberAudit{}))
		projectHandler.AuditLogListRepository.CreateAuditEntriesByKey(requestSession, targetProject.Key, auditEntries)
		observermngmt.FireEvent(observermngmt.ProjectUpdated, observermngmt.ProjectUpdatedData{
			RequestSession: requestSession,
			Old:            &before,
			New:            targetProject,
		})

		projectHandler.ProjectRepository.Update(requestSession, targetProject)

		response = append(response, project.ProjectChildrenMemberSuccessResponseDto{
			Success:     true,
			UserId:      multiProjectUserData.TargetUser,
			ProjectKey:  targetProject.Key,
			ProjectName: targetProject.Name,
			Message:     message.UserAddedToTheProject,
		})
	}
	render.JSON(w, r, response)
}

func (projectHandler *ProjectHandler) ProjectMemberAddHandler(w http.ResponseWriter, r *http.Request) {
	Val.UrlParamUuid(r, "uuid")
	currentProject, requestSession := projectHandler.retrieveProject2(r, true)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	username, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !(rights.AllowProjectUserManagement.Create || rights.AllowAllProjectUserManagement.Create) {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorAddProjectUser))
	}

	logRequestBody(requestSession, r)

	userData := extractMemberRequestBody(r)
	foundUser := projectHandler.UserRepository.FindByUserId(requestSession, userData.TargetUser)

	if !foundUser.IsInternal && userData.UserType == project.OWNER {
		exception.ThrowExceptionClient400Message(message.GetI18N(message.OnlyInternalUsersOwners), "Only internal users can be owners")
	}

	currentProject.CheckIfUserAlreadyExists(userData.TargetUser)

	var before project.Project
	copier.Copy(&before, currentProject)

	auditEntries := make([]*audit2.Audit, 0)
	if userData.IsResponsible {
		if userData.UserType != project.OWNER {
			exception.ThrowExceptionClient400Message(message.GetI18N(message.NonOwnerResponsible), "Only owner users can be flagged as project responsibles")
		}
		if mc := currentProject.ProjectResponsible(); mc != nil {
			before := *mc
			mc.IsResponsible = false
			beforeAudit := before.ToAudit(requestSession, projectHandler.UserRepository)
			afterAudit := mc.ToAudit(requestSession, projectHandler.UserRepository)
			auditEntries = append(auditEntries, audit.CreateAuditEntry(username, message.ProjectUserUpdated, cmp.Diff, afterAudit, beforeAudit))
		}
	}
	currentProject.AddUser(userData)
	entity := currentProject.GetMember(userData.TargetUser)
	if entity == nil {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbNotFound), userData.TargetUser+" not found in DB")
	}
	userAudit := entity.ToAudit(requestSession, projectHandler.UserRepository)
	auditEntries = append(auditEntries, audit.CreateAuditEntry(username, message.ProjectUserCreated, cmp.Diff, userAudit, project.ProjectMemberAudit{}))
	projectHandler.AuditLogListRepository.CreateAuditEntriesByKey(requestSession, currentProject.Key, auditEntries)
	if !hasDummyLabel(currentProject, getDummyLabel(requestSession, projectHandler.LabelRepository)) {
		observermngmt.FireEvent(observermngmt.ProjectUpdated, observermngmt.ProjectUpdatedData{
			RequestSession: requestSession,
			Old:            &before,
			New:            currentProject,
		})
	}
	projectHandler.ProjectRepository.Update(requestSession, currentProject)

	responseData := SuccessResponse{
		Success: true,
		Message: "User added",
	}
	render.JSON(w, r, responseData)
}

func (projectHandler *ProjectHandler) ProjectMemberDeleteHandler(w http.ResponseWriter, r *http.Request) {
	Val.UrlParamUuid(r, "uuid")
	Val.UrlParamEMEA(r, "userId")
	currentProject, requestSession := projectHandler.retrieveProject2(r, true)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}
	username, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !(rights.AllowProjectUserManagement.Delete || rights.AllowAllProjectUserManagement.Delete) {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.UserDeletion))
	}

	userId := chi.URLParam(r, "userId")

	entity := currentProject.GetMember(userId)
	if entity == nil {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbNotFound), userId+" not found in DB")
	}

	errMsg := projectHandler.isProjectMemberInPendingApprovalOrRequest(requestSession, currentProject.Key, userId)
	if errMsg != "" {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorInUse, "Project Memeber"))
	}

	if entity.IsResponsible {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DelResponsible, "Can't delete Project Responsible"))
	}

	currentProject.DeleteUser(userId)

	userAudit := entity.ToAudit(requestSession, projectHandler.UserRepository)
	projectHandler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, currentProject.Key, username, message.ProjectUserDeleted, cmp.Diff, project.ProjectMemberAudit{}, userAudit)

	projectHandler.ProjectRepository.Update(requestSession, currentProject)

	responseData := SuccessResponse{
		Success: true,
		Message: "User deleted",
	}
	render.JSON(w, r, responseData)
}

func (projectHandler *ProjectHandler) ProjectMemberGetUsageInPendingApprovalOrRequest(w http.ResponseWriter, r *http.Request) {
	Val.UrlParamUuid(r, "uuid")
	Val.UrlParamEMEA(r, "userId")
	currentProject, requestSession := projectHandler.retrieveProject2(r, true)
	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !(rights.AllowProjectUserManagement.Read || rights.AllowAllProjectUserManagement.Read) {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ViewUsers))
	}

	userId := chi.URLParam(r, "userId")

	errMsg := projectHandler.isProjectMemberInPendingApprovalOrRequest(requestSession, currentProject.Key, userId)
	render.JSON(w, r, SuccessResponse{
		Success: errMsg != "",
		Message: "Project Member usage in pending Approval or Review Request",
	})
}

func (projectHandler *ProjectHandler) isProjectMemberInPendingApprovalOrRequest(requestSession *logy.RequestSession, projectKey, userId string) string {
	approvalList := projectHandler.ApprovalListRepository.FindByKey(requestSession, projectKey, true)
	if approvalList == nil {
		return ""
	}
	user := projectHandler.UserRepository.FindByUserId(requestSession, userId)

	if user != nil {
		return projectHandler.UserService.IsProjectMemberInPendingApprovalOrRequestUser(requestSession, user, approvalList)
	}
	return ""
}

func (projectHandler *ProjectHandler) ProjectMemberUpdateHandler(w http.ResponseWriter, r *http.Request) {
	Val.UrlParamUuid(r, "uuid")
	Val.UrlParamEMEA(r, "userId")
	currentProject, requestSession := projectHandler.retrieveProject2(r, true)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}
	username, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !(rights.AllowProjectUserManagement.Update || rights.AllowAllProjectUserManagement.Update) {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.UserUpdateNotAuthorized))
	}

	userData := extractMemberRequestBody(r)
	userId := chi.URLParam(r, "userId")
	foundUser := projectHandler.UserRepository.FindByUserId(requestSession, userData.TargetUser)

	if !foundUser.IsInternal && userData.UserType == project.OWNER {
		exception.ThrowExceptionClient400Message(message.GetI18N(message.OnlyInternalUsersOwners), "Only internal users can be owners")
	}

	userBeforeUpdate := currentProject.GetMember(userId)
	if userBeforeUpdate == nil {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbNotFound), userId+" not found in DB")
	}

	errMsg := projectHandler.isProjectMemberInPendingApprovalOrRequest(requestSession, currentProject.Key, userId)
	if errMsg != "" && userData.TargetUser != userId {
		exception.ThrowExceptionClientWithHttpCode(message.ErrorInUseOnUpdate, message.GetI18N(message.ErrorInUseOnUpdate).Text, "", exception.HTTP_CODE_SHOW_NO_REQUEST_ID)
	}

	if userData.IsResponsible && userData.UserType != project.OWNER {
		exception.ThrowExceptionClient400Message(message.GetI18N(message.NonOwnerResponsible), "Only owner users can be flagged as project responsible")
	}

	if userBeforeUpdate.UserType == project.OWNER && userData.UserType != project.OWNER && userData.IsResponsible {
		exception.ThrowExceptionClient400Message(message.GetI18N(message.NonOwnerResponsible), "Only owner users can be flagged as project responsible")
	}

	var before project.Project
	copier.Copy(&before, currentProject)
	auditEntries := make([]*audit2.Audit, 0)
	if mc := currentProject.ProjectResponsible(); mc != nil {
		if mc.UserId == userId {
			if !userData.IsResponsible {
				exception.ThrowExceptionClient400Message(message.GetI18N(message.OneResponsibleOnly), "There has to be exactly one project responsible per project.")
			}
		} else if userData.IsResponsible {
			before := *mc
			mc.IsResponsible = false
			beforeAudit := before.ToAudit(requestSession, projectHandler.UserRepository)
			afterAudit := mc.ToAudit(requestSession, projectHandler.UserRepository)
			auditEntries = append(auditEntries, audit.CreateAuditEntry(username, message.ProjectUserUpdated, cmp.Diff, afterAudit, beforeAudit))
		}
	}

	userAuditBefore := userBeforeUpdate.ToAudit(requestSession, projectHandler.UserRepository)

	currentProject.UpdateProjectMember(userId, userData)

	userAfterUpdate := currentProject.GetMember(userData.TargetUser)
	if userAfterUpdate == nil {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbNotFound), userData.TargetUser+" not found in DB")
	}

	userAuditAfter := userAfterUpdate.ToAudit(requestSession, projectHandler.UserRepository)
	auditEntries = append(auditEntries, audit.CreateAuditEntry(username, message.ProjectUserUpdated, audit2.DiffWithReporter, userAuditAfter, userAuditBefore))
	projectHandler.AuditLogListRepository.CreateAuditEntriesByKey(requestSession, currentProject.Key, auditEntries)

	if !hasDummyLabel(currentProject, getDummyLabel(requestSession, projectHandler.LabelRepository)) {
		observermngmt.FireEvent(observermngmt.ProjectUpdated, observermngmt.ProjectUpdatedData{
			RequestSession: requestSession,
			Old:            &before,
			New:            currentProject,
		})
	}
	projectHandler.ProjectRepository.Update(requestSession, currentProject)

	responseData := SuccessResponse{
		Success: true,
		Message: "User updated",
	}
	render.JSON(w, r, responseData)
}

func (projectHandler *ProjectHandler) ProjectUserGetAllBySearchFragmentHandler(w http.ResponseWriter, r *http.Request) {
	Val.UrlParamUuid(r, "uuid")
	Val.UrlParamSearchText(r, "searchFragment")

	currentProject, requestSession := projectHandler.retrieveProject2(r, false)
	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !(rights.AllowProjectUserManagement.Create || rights.AllowAllProjectUserManagement.Create) {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorProjectMemberSearch))
	}

	searchFragmentEscaped := chi.URLParam(r, "searchFragment")
	searchFragment, err := url.QueryUnescape(searchFragmentEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ErrorKeyRequestParamNotValid, "searchFragment"), zapcore.InfoLevel)

	userFilter := userRepo.UserFilter{}
	if activeParam := r.URL.Query().Get("active"); activeParam != "" {
		active, parseErr := strconv.ParseBool(activeParam)
		if parseErr == nil {
			userFilter.Active = &active
		}
	}

	users := projectHandler.UserRepository.Find5UsersBySearchFragment(requestSession, searchFragment, userFilter)
	render.JSON(w, r, user.ToDtos(users))
}

func extractMemberRequestBody(r *http.Request) project.ProjectMemberRequestDto {
	var userData project.ProjectMemberRequestDto
	validation.DecodeAndValidate(r, &userData, false)
	return userData
}

func extractMultiProjectMemberRequestBody(r *http.Request) project.MultiProjectMemberRequestDto {
	var userData project.MultiProjectMemberRequestDto
	validation.DecodeAndValidate(r, &userData, false)
	return userData
}
