// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/eclipse-disuko/disuko/domain"
	approval2 "github.com/eclipse-disuko/disuko/domain/approval"
	"github.com/eclipse-disuko/disuko/domain/newsbox"
	"github.com/eclipse-disuko/disuko/domain/search"
	"go.uber.org/zap/zapcore"

	"github.com/eclipse-disuko/disuko/connector/userrole"
	"github.com/eclipse-disuko/disuko/helper/stopwatch"

	"github.com/eclipse-disuko/disuko/helper/filter"

	"github.com/eclipse-disuko/disuko/helper/sort"

	"github.com/eclipse-disuko/disuko/infra/repository/approvallist"
	labels "github.com/eclipse-disuko/disuko/infra/repository/labels"
	projectRepo "github.com/eclipse-disuko/disuko/infra/repository/project"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain/audit"
	"github.com/eclipse-disuko/disuko/domain/oauth"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/user"
	auditHelper "github.com/eclipse-disuko/disuko/helper/audit"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/jwt"
	jwt2 "github.com/eclipse-disuko/disuko/helper/jwt"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/reflection"
	"github.com/eclipse-disuko/disuko/helper/roles"
	"github.com/eclipse-disuko/disuko/helper/validation"
	"github.com/eclipse-disuko/disuko/infra/repository/database"
	"github.com/eclipse-disuko/disuko/infra/repository/jobs"
	newsboxRepo "github.com/eclipse-disuko/disuko/infra/repository/newsbox"
	user2 "github.com/eclipse-disuko/disuko/infra/repository/user"
	userService "github.com/eclipse-disuko/disuko/infra/service/user"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type UserHandler struct {
	UserRepository         user2.IUsersRepository
	JobRepository          jobs.IJobsRepository
	ApprovalListRepository approvallist.IApprovalListRepository
	ProjectRepository      projectRepo.IProjectRepository
	LabelRepository        labels.ILabelRepository
	UserroleConnector      *userrole.Connector
	NewsBoxRepository      newsboxRepo.IRepo
	DeletionService        *userService.DeletionService
	UserService            *userService.Service
	DeletionAuditService   *userService.DeletionAuditService
}

type DeletePersonalDataEffectedEntities struct {
	UserTasksCount  int `json:"user_tasks_count"`
	UserRolesCount  int `json:"user_roles_count"`
	DataTracesCount int `json:"data_traces_count"`
}
type DeletePersonalDataResponse struct {
	Success          bool                               `json:"success"`
	Message          string                             `json:"message"`
	EntitiesEffected DeletePersonalDataEffectedEntities `json:"entities_effected"`
	DetailedPlan     *userService.DeletionPlan          `json:"detailed_plan,omitempty"`
}

type PersonalDetailEntity struct {
	EntityID            string `json:"entityID"`
	EntityType          string `json:"entityType"`
	EntitySubType       string `json:"entitySubType,omitempty"`
	EntityStatus        string `json:"entityStatus,omitempty"`
	EntityName          string `json:"entityName,omitempty"`
	ProjectID           string `json:"projectID,omitempty"`
	ProjectName         string `json:"projectName,omitempty"`
	DisableDeleteReason string `json:"disableDeleteReason,omitempty"`
}

type PersonalDetailsResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    []PersonalDetailEntity `json:"data"`
}

func (handler *UserHandler) extractRequestBody(r *http.Request) user.UserRequestDto {
	var result user.UserRequestDto
	validation.DecodeAndValidate(r, &result, false)
	return result
}

func (handler *UserHandler) extractLastSeenBody(r *http.Request) user.UserLastSeenDto {
	var result user.UserLastSeenDto
	validation.DecodeAndValidate(r, &result, false)
	return result
}

func (handler *UserHandler) extractRolesRequestBody(r *http.Request) user.UserRolesRequestDto {
	var result user.UserRolesRequestDto
	validation.DecodeAndValidate(r, &result, false)
	return result
}

func (handler *UserHandler) extractDeletePersonalDataBody(r *http.Request) user.DeletePersonalDataDto {
	var result user.DeletePersonalDataDto
	validation.DecodeAndValidate(r, &result, false)
	return result
}

func (handler *UserHandler) SearchHandlerForAdmin(w http.ResponseWriter, r *http.Request) {
	requestVersion := r.Header.Get("X-Client-Version")

	if requestVersion == "2.0" {
		// Or we check the request by
		// var body map[string]interface{}
		// json.NewDecoder(r.Body).Decode(&body)
		//  if _, exists := body["newVue3SpecificField"]; exists
		var searchOptionsVue3 search.RequestSearchOptionsNew
		validation.DecodeAndValidate(r, &searchOptionsVue3, false)
		handler.searchHandlerForAdmin(w, r, &searchOptionsVue3)
	} else {
		var searchOptions search.RequestSearchOptions
		validation.DecodeAndValidate(r, &searchOptions, false)
		handler.searchHandlerForAdmin(w, r, &searchOptions)
	}
}

func (handler *UserHandler) searchHandlerForAdmin(w http.ResponseWriter, r *http.Request, searchOptions search.SortableOptions) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.AllowUsers.Create && rights.AllowUsers.Read && rights.AllowUsers.Update && rights.AllowUsers.Delete) {
		exception.ThrowExceptionSendDeniedResponse()
	}

	extractors := map[string]func(dto *user.UserDto) string{
		"isActive":   func(item *user.UserDto) string { return strconv.FormatBool(item.Active) },
		"isInternal": func(item *user.UserDto) string { return strconv.FormatBool(item.IsInternal) },
	}

	users := handler.UserRepository.FindAll(requestSession, false)

	dtos := make([]*user.UserDto, 0)
	for _, entity := range users {
		dto := entity.ToDto()
		if filter.MatchesCriteria(dto, searchOptions, extractors, nil) {
			dtos = append(dtos, dto)
		}
	}
	result := user.AllResponse{
		Users: dtos,
		Count: len(dtos),
	}

	if searchOptions.ShouldOrder() {
		asc := searchOptions.IsSortAsc()
		key := searchOptions.GetSortKey()
		if key == "updated" {
			sort.Sort(result.Users, func(dto *user.UserDto) int64 { return dto.Updated.Unix() }, sort.Int64LessThan, asc)
		} else if key == "isInternal" {
			sort.Sort(result.Users, func(dto *user.UserDto) bool { return dto.IsInternal }, sort.BoolLessThan, asc)
		} else if key == "created" {
			sort.Sort(result.Users, func(dto *user.UserDto) int64 { return dto.Created.Unix() }, sort.Int64LessThan, asc)
		} else if key == "active" {
			sort.Sort(result.Users, func(dto *user.UserDto) bool { return dto.Active }, sort.BoolLessThan, asc)
		} else if key == "user" {
			sort.Sort(result.Users, func(dto *user.UserDto) string { return dto.User }, sort.StringLessThan, asc)
		} else if key == "forename" {
			sort.Sort(result.Users, func(dto *user.UserDto) string { return dto.Forename }, sort.StringLessThan, asc)
		} else if key == "lastname" {
			sort.Sort(result.Users, func(dto *user.UserDto) string { return dto.Lastname }, sort.StringLessThan, asc)
		} else if key == "email" {
			sort.Sort(result.Users, func(dto *user.UserDto) string { return dto.Email }, sort.StringLessThan, asc)
		} else if key == "metaData.department" {
			sort.Sort(result.Users, func(dto *user.UserDto) string { return dto.GetDepartment() }, sort.StringLessThan, asc)
		} else if key == "metaData.departmentDescription" {
			sort.Sort(result.Users, func(dto *user.UserDto) string { return dto.GetDepartmentDescription() }, sort.StringLessThan, asc)
		} else if key == "metaData.companyIdentifier" {
			sort.Sort(result.Users, func(dto *user.UserDto) string { return dto.GetCompanyIdentifier() }, sort.StringLessThan, asc)
		} else if key == "termsOfUseDate" {
			sort.Sort(result.Users, func(dto *user.UserDto) int64 { return dto.GetTermsOfUseDate() }, sort.Int64LessThan, asc)
		} else if key == "termsOfUseVersion" {
			sort.Sort(result.Users, func(dto *user.UserDto) string { return dto.TermsOfUseVersion }, sort.StringLessThan, asc)
		} else if key == "termsOfUse" {
			sort.Sort(result.Users, func(dto *user.UserDto) bool { return dto.TermsOfUse }, sort.BoolLessThan, asc)
		} else if key == "deprovisioned" {
			sort.Sort(result.Users, func(dto *user.UserDto) int64 { return dto.GetDeprovisioned() }, sort.Int64LessThan, asc)
		}
	}

	if searchOptions.HasPaginationActive() && len(result.Users) > 0 {
		lowIndex := (searchOptions.GetPage() - 1) * searchOptions.GetItemsPerPage()
		highIndex := lowIndex + searchOptions.GetItemsPerPage()
		if highIndex > len(result.Users) {
			highIndex = len(result.Users)
		}
		if lowIndex > highIndex {
			lowIndex = 0 // reset page number
		}
		result.Users = result.Users[lowIndex:highIndex]
	}

	render.JSON(w, r, result)
}

func FilterUsers(arr []*user.UserDto, cond func(*user.UserDto) bool) []*user.UserDto {
	newContent := make([]*user.UserDto, 0)
	for i := 0; i < len(arr); i++ {
		if cond(arr[i]) {
			newContent = append(newContent, arr[i])
		}
	}
	return newContent
}

func (handler *UserHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.AllowUsers.Create && rights.AllowUsers.Read && rights.AllowUsers.Update && rights.AllowUsers.Delete) {
		exception.ThrowExceptionSendDeniedResponse()
	}

	users := handler.UserRepository.FindAll(requestSession, false)

	var result user.AllResponse
	result.Users = user.ToDtos(users)
	result.Count = len(users)

	render.JSON(w, r, result)
}

func (handler *UserHandler) GetTermsOfUseCurrentVersionHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.AllowTools.Create && rights.AllowTools.Read && rights.AllowTools.Update && rights.AllowTools.Delete) {
		exception.ThrowExceptionSendDeniedResponse()
	}

	type version struct {
		TermsOfUseCurrentVersion string `json:"termsOfUseCurrentVersion"`
	}

	render.JSON(w, r, version{
		TermsOfUseCurrentVersion: conf.Config.Server.TermsOfUseCurrentVersion,
	})
}

func (handler *UserHandler) GetByUuidHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.AllowUsers.Create && rights.AllowUsers.Read && rights.AllowUsers.Update && rights.AllowUsers.Delete) {
		exception.ThrowExceptionSendDeniedResponse()
	}
	requestedUser := handler.loadRequestedUser(requestSession, r)
	render.JSON(w, r, requestedUser.ToDto())
}

func (handler *UserHandler) GetUserMailByIdHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.AllowUsers.Read || rights.AllowAllProjectUserManagement.Read) {
		exception.ThrowExceptionSendDeniedResponse()
	}
	userIdEscaped := chi.URLParam(r, "userId")
	userId, err := url.QueryUnescape(userIdEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ErrorKeyRequestParamNotValid, "userId"), zapcore.InfoLevel)

	requestedUser := handler.UserRepository.FindByUserId(requestSession, userId)
	render.JSON(w, r, requestedUser.ToUserMailDto())
}

func (handler *UserHandler) GetAuditTrailHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.AllowUsers.Create && rights.AllowUsers.Read && rights.AllowUsers.Update && rights.AllowUsers.Delete) {
		exception.ThrowExceptionSendDeniedResponse()
	}
	requestedUser := handler.loadRequestedUser(requestSession, r)
	auditTrail := make([]audit.AuditDto, 0)
	for _, item := range requestedUser.GetAuditTrail() {
		auditTrail = append(auditTrail, item.ToDto())
	}
	render.JSON(w, r, auditTrail)
}

func (handler *UserHandler) Get5BySearchFragmentHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.IsInternal {
		exception.ThrowExceptionSendDeniedResponse()
	}
	searchFragmentEscaped := chi.URLParam(r, "searchFragment")
	searchFragment, err := url.QueryUnescape(searchFragmentEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ErrorKeyRequestParamNotValid, "searchFragment"), zapcore.InfoLevel)

	userFilter := user2.UserFilter{}
	if activeParam := r.URL.Query().Get("active"); activeParam != "" {
		active, parseErr := strconv.ParseBool(activeParam)
		if parseErr == nil {
			userFilter.Active = &active
		}
	}

	users := handler.UserRepository.Find5UsersBySearchFragment(requestSession, searchFragment, userFilter)
	render.JSON(w, r, user.ToDtos(users))
}

func (handler *UserHandler) UpdateHandlerForUser(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	userName, _ := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	currentUser := handler.loadRequestedUser(requestSession, r)

	if currentUser == nil || userName != currentUser.User {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorPermissionDeniedUser, userName))
	}

	userData := handler.extractRequestBody(r)
	if len(userData.User) == 0 && userData.TermsOfUse {
		oldAudit := currentUser.ToUserAudit()
		currentUser.TermsOfUse = userData.TermsOfUse
		currentUser.TermsOfUseDate = reflection.ToPointer(time.Now())
		newAudit := currentUser.ToUserAudit()
		auditHelper.CreateAndAddAuditEntry(&currentUser.Container, userName, message.UserUpdated, audit.DiffWithReporter, newAudit, oldAudit)
	} else {
		currentUser.UpdateData(userData)
	}
	handler.UserRepository.Update(requestSession, currentUser)

	render.JSON(w, r, currentUser.ToDto())
}

func (handler *UserHandler) GetNewsBoxItems(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	user, _ := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)

	currentUserProfile := handler.UserRepository.FindByUserId(requestSession, user)
	if currentUserProfile == nil {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorPermissionDeniedUser, user))
	}

	userLastSeenId := currentUserProfile.NewsboxLastSeenId

	qc := &database.QueryConfig{}
	qc.SetSort(database.SortConfig{
		database.SortAttribute{
			Name:  "Created",
			Order: database.ASC,
		},
	})
	all := handler.NewsBoxRepository.Query(requestSession, qc)

	var (
		nonExpired []*newsbox.Item
		allI       = -1
	)
	for i, item := range all {
		if item.Key == userLastSeenId {
			allI = i
		}
		if item.Expiry.IsZero() || item.Expiry.After(time.Now()) {
			nonExpired = append(nonExpired, item)
		}
	}

	if allI == -1 {
		render.JSON(w, r, newsbox.NewsboxResponse{
			Items:  domain.ToDtos(nonExpired),
			ToShow: 0,
		})
		return
	}

	for i := allI + 1; i < len(all); i++ {
		item := all[i]
		if item.Expiry.IsZero() || item.Expiry.After(time.Now()) {
			render.JSON(w, r, newsbox.NewsboxResponse{
				Items:  domain.ToDtos(nonExpired),
				ToShow: findNewsboxIndex(item.Key, nonExpired),
			})
			return
		}
	}

	currentItem := all[allI]
	if currentItem.Expiry.IsZero() || currentItem.Expiry.After(time.Now()) {
		render.JSON(w, r, newsbox.NewsboxResponse{
			Items:  domain.ToDtos(nonExpired),
			ToShow: findNewsboxIndex(currentItem.Key, nonExpired),
		})
		return
	}

	for i := allI - 1; i >= 0; i-- {
		item := all[i]
		if item.Expiry.IsZero() || item.Expiry.After(time.Now()) {
			render.JSON(w, r, newsbox.NewsboxResponse{
				Items:  domain.ToDtos(nonExpired),
				ToShow: findNewsboxIndex(item.Key, nonExpired),
			})
			return
		}
	}

	render.JSON(w, r, newsbox.NewsboxResponse{
		Items:  domain.ToDtos(nonExpired),
		ToShow: -1,
	})
}

func findNewsboxIndex(key string, items []*newsbox.Item) int {
	for i, item := range items {
		if item.Key == key {
			return i
		}
	}
	return -1
}

func (handler *UserHandler) UpdateLastSeen(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	userName, _ := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	currentUser := handler.loadRequestedUser(requestSession, r)

	if currentUser == nil || userName != currentUser.User {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorPermissionDeniedUser, userName))
	}

	userData := handler.extractLastSeenBody(r)
	currentUser.UpdateNewsboxLastSeenData(userData)
	handler.UserRepository.Update(requestSession, currentUser)

	render.JSON(w, r, currentUser.ToDto())
}

func (handler *UserHandler) UpdateUserRolesHandlerForAdmin(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	userName, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.AllowUsers.Create && rights.AllowUsers.Read && rights.AllowUsers.Update && rights.AllowUsers.Delete) {
		exception.ThrowExceptionSendDeniedResponse()
	}

	currentUser := handler.loadRequestedUser(requestSession, r)
	if currentUser == nil || userName != currentUser.User {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorPermissionDeniedUser, userName))
	}

	userRolesData := handler.extractRolesRequestBody(r)
	currentUser.Roles = userRolesData.Roles

	handler.UserRepository.Update(requestSession, currentUser)
	dto := currentUser.ToDto()
	render.JSON(w, r, dto)
}

func (handler *UserHandler) GetNewTokensHandlerForAdmin(w http.ResponseWriter, r *http.Request) {
	handler.handleGetNewTokens(w, r, false)
}

func (handler *UserHandler) GetNewTokensForNonInternalHandlerForAdmin(w http.ResponseWriter, r *http.Request) {
	handler.handleGetNewTokens(w, r, true)
}

func (handler *UserHandler) handleGetNewTokens(w http.ResponseWriter, r *http.Request, forceNonInternal bool) {
	requestSession := logy.GetRequestSession(r)
	userName, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.AllowUsers.Create && rights.AllowUsers.Read && rights.AllowUsers.Update && rights.AllowUsers.Delete) {
		exception.ThrowExceptionSendDeniedResponse()
	}
	currentUser := handler.loadRequestedUser(requestSession, r)

	if currentUser == nil || userName != currentUser.User {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorPermissionDeniedUser, userName))
	}

	token := jwt.ExtractTokenMetadata(requestSession, r)

	isInternal := token.IsInternalEmployee
	if forceNonInternal {
		isInternal = false
	}

	userData := jwt2.CreateUserData(requestSession, currentUser, token.Username, token.Email, currentUser.Roles, token.GroupType, isInternal, r)
	tokenDetails := jwt2.CreateToken(userData)
	accessData := roles.GetAccessAndRolesRightsFromClaim(userData)
	response := oauth.OAuthTokenResponse{
		Rights:  accessData,
		Profile: currentUser.ToDto(),
	}
	cookieRefreshToken := http.Cookie{
		Name:     "oauth.r",
		Value:    tokenDetails.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	cookieAccessToken := http.Cookie{
		Name:     "oauth.a",
		Value:    tokenDetails.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookieRefreshToken)
	http.SetCookie(w, &cookieAccessToken)

	render.JSON(w, r, response)
}

func (handler *UserHandler) EnableDisableHandlerForAdmin(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	userName, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.AllowUsers.Create && rights.AllowUsers.Read && rights.AllowUsers.Update && rights.AllowUsers.Delete) {
		exception.ThrowExceptionSendDeniedResponse()
	}

	currentUser := handler.loadRequestedUser(requestSession, r)

	oldAudit := currentUser.ToUserAudit()
	userData := handler.extractRequestBody(r)
	currentUser.Active = userData.Active
	newAudit := currentUser.ToUserAudit()
	auditHelper.CreateAndAddAuditEntry(&currentUser.Container, userName, message.UserUpdated, audit.DiffWithReporter, newAudit, oldAudit)

	handler.UserRepository.Update(requestSession, currentUser)

	dto := currentUser.ToDto()
	render.JSON(w, r, dto)
}

func (handler *UserHandler) EnableDisableHandlerForTest(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	userName, _ := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)

	currentUser := handler.loadRequestedUserTest(requestSession, r)

	oldAudit := currentUser.ToUserAudit()
	userData := handler.extractRequestBody(r)
	currentUser.Active = userData.Active
	newAudit := currentUser.ToUserAudit()
	auditHelper.CreateAndAddAuditEntry(&currentUser.Container, userName, message.UserUpdated, audit.DiffWithReporter, newAudit, oldAudit)

	handler.UserRepository.Update(requestSession, currentUser)

	dto := currentUser.ToDto()
	render.JSON(w, r, dto)
}

func (handler *UserHandler) UpdateHandlerForAdmin(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.AllowUsers.Create && rights.AllowUsers.Read && rights.AllowUsers.Update && rights.AllowUsers.Delete) {
		exception.ThrowExceptionSendDeniedResponse()
	}

	currentUser := handler.loadRequestedUser(requestSession, r)
	userData := handler.extractRequestBody(r)
	currentUser.UpdateData(userData)

	handler.UserRepository.Update(requestSession, currentUser)

	dto := currentUser.ToDto()
	render.JSON(w, r, dto)
}

func (handler *UserHandler) loadRequestedUser(requestSession *logy.RequestSession, r *http.Request) *user.User {
	userUUIDEscaped := chi.URLParam(r, "uuid")
	uuid, err := url.QueryUnescape(userUUIDEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ErrorKeyRequestParamNotValid, "uuid"), zapcore.InfoLevel)
	err = validation.CheckUuid(uuid)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ErrorKeyRequestParamNotValid, "uuid"), zapcore.InfoLevel)
	currentUser := handler.UserRepository.FindByKey(requestSession, uuid, false)
	return currentUser
}

func (handler *UserHandler) loadRequestedUserTest(requestSession *logy.RequestSession, r *http.Request) *user.User {
	userId := chi.URLParam(r, "userId")
	currentUser := handler.UserRepository.FindByUserId(requestSession, userId)
	return currentUser
}

func ConvertUsernameWithoutEmea(userName string) string {
	index := strings.LastIndex(userName, "\\")
	if index > -1 {
		return string(userName[index+1:])
	}
	return userName
}

func GetUserByUsername(requestSession *logy.RequestSession, repository user2.IUsersRepository, name string) *user.UserDto {
	entity := repository.FindByUserId(requestSession, ConvertUsernameWithoutEmea(name))
	if entity == nil {
		return nil
	}
	return entity.ToDto()
}

func (handler *UserHandler) GetProfileData(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	userName, _ := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	userData := GetUserByUsername(requestSession, handler.UserRepository, userName)
	if userData == nil {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorPermissionDeniedUser, userName))
	}
	tokenData := jwt.ExtractTokenMetadata(requestSession, r)
	accessData := roles.GetAccessAndRolesRightsFromClaim(*tokenData)

	response := oauth.OAuthTokenResponse{
		Rights:  accessData,
		Profile: userData,
	}
	render.JSON(w, r, response)
}

func (handler *UserHandler) GetTaskListForAdmin(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.AllowUsers.Create && rights.AllowUsers.Read && rights.AllowUsers.Update && rights.AllowUsers.Delete) {
		exception.ThrowExceptionSendDeniedResponse()
	}
	tw1 := stopwatch.StopWatch{}
	tw1.Start()
	logy.Infof(requestSession, "Start GetTaskList")
	requestedUser := handler.loadRequestedUser(requestSession, r)

	service := userService.Init(requestSession, handler.UserRepository, handler.ApprovalListRepository, handler.ProjectRepository, handler.LabelRepository)
	enrichedTasks := service.EnrichTasks(requestedUser.Tasks)

	result := make([]user.TaskDto, 0)
	for _, enriched := range enrichedTasks {
		status := enriched.Approval.ToApprovalDtoStatus()
		delegatedToFullName := ""
		if enriched.Task.DelegatedTo != "" {
			delegatedToFullName = service.FullName(enriched.Task.DelegatedTo)
		}
		projectType := service.GetProjectType(enriched.Pr)
		result = append(result, enriched.Task.ToDto(service.FullName(enriched.Task.Creator), enriched.Approval.Type, status, enriched.Pr.Name, enriched.Pr.IsGroup, delegatedToFullName, projectType))
	}
	tw1.Stop()
	logy.Infof(requestSession, "End GetTaskList: Iterate tasks time: %s cnt:%d", tw1.FormatSeconds(), len(requestedUser.Tasks))

	dtoList := result
	render.JSON(w, r, dtoList)
}

func (handler *UserHandler) GetTaskList(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, handler.prepareTaskDtoList(r))
}

func (handler *UserHandler) GetTaskListCsv(w http.ResponseWriter, r *http.Request) {
	taskList := handler.prepareTaskDtoList(r)

	csvWriter := csv.NewWriter(w)
	csvWriter.Comma = ';'
	defer csvWriter.Flush()

	csvHeader := []string{
		"STATUS",
		"TYPE",
		"CREATOR",
		"DEPARTMENT",
		"DELEGATED TO",
		"REFERENCE",
		"RESULT",
		"UPDATED",
		"CREATED",
	}

	if err := csvWriter.Write(csvHeader); err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCsvGeneration, "user tasks", "header"), err)
	}

	for _, task := range taskList {
		csvRow := []string{
			string(task.Status),
			message.GetI18N(fmt.Sprintf("TASK_TYPE_%s_%s", task.TargetType, task.Type)).Text,
			fmt.Sprintf("%s (%s)", task.CreatorFullName, task.Creator),
			fmt.Sprintf("%s (%s)", task.CreatorDepartmentDescription, task.CreatorDepartment),
			fmt.Sprintf("%s (%s)", task.DelegatedToFullName, task.DelegatedTo),
			fmt.Sprintf("Project: %s", task.ProjectName),
			message.GetI18N(fmt.Sprintf("APPROVAL_STATUS_%s_%s", task.TargetType, task.ResultStatus)).Text,
			task.Updated.Format("02.01.2006"),
			task.Created.Format("02.01.2006"),
		}
		if err := csvWriter.Write(csvRow); err != nil {
			exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCsvGeneration, "user tasks", "data"), err)
		}
	}
}

func (handler *UserHandler) prepareTaskDtoList(r *http.Request) []user.TaskDto {
	requestSession := logy.GetRequestSession(r)
	userName, _ := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)

	tw1 := stopwatch.StopWatch{}
	tw1.Start()
	logy.Infof(requestSession, "Start GetTaskList")

	currentUserProfile := handler.UserRepository.FindByUserId(requestSession, userName)
	if currentUserProfile == nil || userName != currentUserProfile.User {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorPermissionDeniedUser, userName))
	}

	if slices.Contains(currentUserProfile.Roles, roles.FossOfficeUser) && len(strings.TrimSpace(conf.Config.Server.FOSSOfficeUserId)) > 0 {
		poolUserProfile := handler.UserRepository.FindByUserId(requestSession, conf.Config.Server.FOSSOfficeUserId)
		if poolUserProfile != nil {
			currentUserProfile.Tasks = append(currentUserProfile.Tasks, poolUserProfile.Tasks...)
		}

	}
	service := userService.Init(requestSession, handler.UserRepository, handler.ApprovalListRepository, handler.ProjectRepository, handler.LabelRepository)
	enrichedTasks := service.EnrichTasks(currentUserProfile.Tasks)

	result := make([]user.TaskDto, 0)
	for _, enriched := range enrichedTasks {
		status := enriched.Approval.ToApprovalDtoStatus()
		delegatedToFullName := ""
		if enriched.Task.DelegatedTo != "" {
			delegatedToFullName = service.FullName(enriched.Task.DelegatedTo)
		}
		projectType := service.GetProjectType(enriched.Pr)
		result = append(result, enriched.Task.ToDto(service.FullName(enriched.Task.Creator), enriched.Approval.Type, status, enriched.Pr.Name, enriched.Pr.IsGroup, delegatedToFullName, projectType))
	}
	tw1.Stop()
	logy.Infof(requestSession, "End GetTaskList: Iterate tasks time: %s cnt:%d", tw1.FormatSeconds(), len(currentUserProfile.Tasks))
	return result
}

func (handler *UserHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	userName, _ := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)

	currentUserProfile := handler.UserRepository.FindByUserId(requestSession, userName)
	if currentUserProfile == nil || userName != currentUserProfile.User {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorPermissionDeniedUser, userName))
	}
	taskId := chi.URLParam(r, "taskId")
	var task user.Task
	found := false
	for _, taskEntity := range currentUserProfile.Tasks {
		if taskEntity.Key == taskId {
			task = taskEntity
			found = true
			break
		}
	}
	if !found {
		logy.Warnf(requestSession, "Task %s not found for user %s.", taskId, userName)
		exception.ThrowExceptionClientWithHttpCode(message.GetI18N(message.TaskNotFound).Code, message.GetI18N(message.TaskNotFound).Text, "", exception.HTTP_CODE_SHOW_NO_REQUEST_ID)
	}

	service := userService.Init(requestSession, handler.UserRepository, handler.ApprovalListRepository, handler.ProjectRepository, handler.LabelRepository)
	omit, enriched := service.EnrichTask(task)
	if omit {
		exception.ThrowExceptionClientWithHttpCode(message.GetI18N(message.TaskNotFound).Code, message.GetI18N(message.TaskNotFound).Text, "", exception.HTTP_CODE_SHOW_NO_REQUEST_ID)
	}
	delegatedToFullName := ""
	if enriched.Task.DelegatedTo != "" {
		delegatedToFullName = service.FullName(enriched.Task.DelegatedTo)
	}
	projectType := service.GetProjectType(enriched.Pr)
	taskDto := enriched.Task.ToDto(service.FullName(enriched.Task.Creator), enriched.Approval.Type, enriched.Approval.ToApprovalDtoStatus(), enriched.Pr.Name, enriched.Pr.IsGroup, delegatedToFullName, projectType)
	render.JSON(w, r, taskDto)
}

func (handler *UserHandler) GetProjectRoles(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	userName, _ := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)

	currentUserProfile := handler.UserRepository.FindByUserId(requestSession, userName)
	if currentUserProfile == nil || userName != currentUserProfile.User {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorPermissionDeniedUser, userName))
	}

	service := userService.Init(requestSession, handler.UserRepository, handler.ApprovalListRepository, handler.ProjectRepository, handler.LabelRepository)
	roles := service.Roles(currentUserProfile)
	res := make([]project.ProjectRoleDto, 0)
	for _, r := range roles {
		res = append(res, r.PrUser.ToProjectRoleDto(r.Pr))
	}
	render.JSON(w, r, res)
}

func (handler *UserHandler) GetProjectRolesForAdmin(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.AllowUsers.Create && rights.AllowUsers.Read && rights.AllowUsers.Update && rights.AllowUsers.Delete) {
		exception.ThrowExceptionSendDeniedResponse()
	}
	requestedUser := handler.loadRequestedUser(requestSession, r)

	service := userService.Init(requestSession, handler.UserRepository, handler.ApprovalListRepository, handler.ProjectRepository, handler.LabelRepository)
	roles := service.Roles(requestedUser)
	res := make([]project.ProjectRoleDto, 0)
	for _, r := range roles {
		res = append(res, r.PrUser.ToProjectRoleDto(r.Pr))
	}
	render.JSON(w, r, res)
}

func (handler *UserHandler) DelegateTask(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	userName, _ := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)

	taskId := chi.URLParam(r, "taskId")
	if taskId == "" {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.TaskNotFound))
	}

	var req user.DelegateTaskDto
	if err := render.Bind(r, &req); err != nil {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorUserUpdate))
	}

	service := userService.Init(requestSession, handler.UserRepository, handler.ApprovalListRepository, handler.ProjectRepository, handler.LabelRepository)
	err := service.DelegateTask(taskId, req.DelegateUserId, userName)
	if err != nil {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorUserUpdate))
	}

	render.JSON(w, r, map[string]string{"status": "success"})
}

func (handler *UserHandler) DeletePersonalDataHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)

	if !rights.IsDomainAdmin() {
		exception.ThrowExceptionSendDeniedResponse()
	}

	requestData := handler.extractDeletePersonalDataBody(r)

	var resp DeletePersonalDataResponse
	if requestData.DryRun {
		resp = handler.userPersonalDeletionDryRun(requestSession, requestData.Username)
	} else {
		resp = handler.userPersonalDeletionExecute(requestSession, requestData.Username, requestData.EntityType)
	}

	render.JSON(w, r, resp)
}

func (handler *UserHandler) userPersonalDeletionDryRun(requestSession *logy.RequestSession, userName string) DeletePersonalDataResponse {
	user := handler.UserRepository.FindByUserId(requestSession, userName)
	if user == nil {
		return DeletePersonalDataResponse{
			Success: false,
			Message: message.GetI18N(message.UserManagementUserNotFound).Code,
		}
	}

	plan, err := handler.DeletionService.ExecuteDeletion(userName)
	if err != nil {
		return DeletePersonalDataResponse{
			Success: false,
			Message: message.GetI18N(message.Error).Code,
		}
	}

	taskCount, roleCount, traceCount := handler.DeletionService.GetUserDeletionStats(user)

	return DeletePersonalDataResponse{
		Success: true,
		Message: message.GetI18N(message.UserManagementDryRunSuccess).Code,
		EntitiesEffected: DeletePersonalDataEffectedEntities{
			UserTasksCount:  taskCount,
			UserRolesCount:  roleCount,
			DataTracesCount: traceCount,
		},
		DetailedPlan: plan,
	}
}

func (handler *UserHandler) userPersonalDeletionExecute(requestSession *logy.RequestSession, userName string, entityType string) DeletePersonalDataResponse {
	// Validate entity type
	validEntityTypes := map[string]bool{
		"all":     true,
		"tasks":   true,
		"roles":   true,
		"traces":  true,
		"profile": true,
	}

	if !validEntityTypes[entityType] {
		logy.Warnf(requestSession, "Invalid entity type provided: %s", entityType)
		return DeletePersonalDataResponse{
			Success: false,
			Message: message.GetI18N(message.ErrorKeyRequestParamNotValid, "entity_type").Code,
		}
	}

	if entityType == "all" {

		logy.Infof(requestSession, "Executing full personal data deletion for user: %s", userName)

		_, err := handler.DeletionService.ExecuteDeletion(userName)
		if err != nil {
			logy.Errorf(requestSession, "Failed to execute deletion for user %s: %v", userName, err)
			return DeletePersonalDataResponse{
				Success: false,
				Message: message.GetI18N(message.Error).Code,
			}
		}

		return DeletePersonalDataResponse{
			Success:          true,
			Message:          message.GetI18N(message.UserManagementDeletionSuccess).Code,
			EntitiesEffected: DeletePersonalDataEffectedEntities{},
		}
	}

	logy.Infof(requestSession, "Executing personal data deletion for user: %s, entity_type: %s", userName, entityType)

	switch entityType {
	case "tasks":

		logy.Infof(requestSession, "Deleting all tasks for user: %s", userName)
	case "roles":

		logy.Infof(requestSession, "Deleting all roles for user: %s", userName)
	case "traces":

		logy.Infof(requestSession, "Deleting all traces for user: %s", userName)
	case "profile":

		logy.Infof(requestSession, "Deleting user profile for user: %s", userName)
	}

	return DeletePersonalDataResponse{
		Success:          true,
		Message:          message.GetI18N(message.UserManagementEntityDeleted).Code,
		EntitiesEffected: DeletePersonalDataEffectedEntities{},
	}
}

func (handler *UserHandler) DeletePersonalDataDryRunHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)

	if !rights.IsDomainAdmin() {
		exception.ThrowExceptionSendDeniedResponse()
	}

	userName := r.URL.Query().Get("username")
	if userName == "" {
		render.JSON(w, r, DeletePersonalDataResponse{
			Success: false,
			Message: message.GetI18N(message.ErrorKeyRequestParamNotValid, "username").Code,
		})
		return
	}

	user := handler.UserRepository.FindByUserId(requestSession, userName)
	if user == nil {
		render.JSON(w, r, DeletePersonalDataResponse{
			Success: false,
			Message: message.GetI18N(message.UserManagementUserNotFound).Code,
		})
		return
	}

	taskCount, roleCount, traceCount := handler.DeletionService.GetUserDeletionStats(user)

	render.JSON(w, r, DeletePersonalDataResponse{
		Success: true,
		Message: message.GetI18N(message.UserManagementDryRunSuccess).Code,
		EntitiesEffected: DeletePersonalDataEffectedEntities{
			UserTasksCount:  taskCount,
			UserRolesCount:  roleCount,
			DataTracesCount: traceCount,
		},
	})
}

func (handler *UserHandler) GetPersonalDetailsHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)

	if !rights.IsDomainAdmin() {
		exception.ThrowExceptionSendDeniedResponse()
	}

	userName := chi.URLParam(r, "username")
	entity := r.URL.Query().Get("entity")

	if userName == "" {
		render.JSON(w, r, PersonalDetailsResponse{
			Success: false,
			Message: message.GetI18N(message.ErrorKeyRequestParamNotValid, "username").Code,
			Data:    []PersonalDetailEntity{},
		})
		return
	}

	if entity == "" {
		render.JSON(w, r, PersonalDetailsResponse{
			Success: false,
			Message: message.GetI18N(message.ErrorKeyRequestParamNotValid, "entity").Code,
			Data:    []PersonalDetailEntity{},
		})
		return
	}

	user := handler.UserRepository.FindByUserId(requestSession, userName)
	if user == nil {
		render.JSON(w, r, PersonalDetailsResponse{
			Success: false,
			Message: message.GetI18N(message.UserManagementUserNotFound).Code,
			Data:    []PersonalDetailEntity{},
		})
		return
	}

	var entities []PersonalDetailEntity

	switch entity {
	case "tasks":
		projectIDNameMap := map[string]string{}

		projectIds := []string{}
		for _, val := range user.Tasks {
			projectIds = append(projectIds, val.ProjectGuid)
		}
		projects := handler.ProjectRepository.FindByKeys(requestSession, projectIds, false)

		for _, project := range projects {
			projectIDNameMap[project.Key] = project.Name
		}

		for _, task := range user.Tasks {
			entities = append(entities, PersonalDetailEntity{
				EntityID:      task.Key,
				EntityType:    "tasks",
				EntitySubType: string(task.Type),
				EntityStatus:  string(task.Status),
				ProjectID:     task.ProjectGuid,
				ProjectName:   projectIDNameMap[task.ProjectGuid],
			})
		}
	case "roles":
		projects := handler.ProjectRepository.FindAllForUser(requestSession, userName)

		projectIds := []string{}
		for _, projectLocal := range projects {
			projectIds = append(projectIds, projectLocal.Key)
		}
		approvalLists := handler.ApprovalListRepository.FindByKeys(requestSession, projectIds, false)

		projectIDApprovalListMaps := map[string]*approval2.ApprovalList{}

		for _, approvalList := range approvalLists {
			if approvalList != nil && len(approvalList.Approvals) > 0 {
				projectIDApprovalListMaps[approvalList.Approvals[0].ProjectGuid] = approvalList
			}
		}

		for _, projectLocal := range projects {
			var userLocal *project.ProjectMemberEntity
			for _, member := range projectLocal.UserManagement.Users {
				if member.UserId == userName {
					userLocal = member
				}
			}
			if userLocal == nil {
				continue
			}

			disableDeleteReason := ""
			if userLocal.IsResponsible {
				disableDeleteReason = "DLG_CAN_NOT_DELETE_RESPONSIBLE"
			}
			if disableDeleteReason == "" {
				if approvalList, ok := projectIDApprovalListMaps[projectLocal.Key]; ok {
					disableDeleteReason = handler.UserService.IsProjectMemberInPendingApprovalOrRequestUser(requestSession, user, approvalList)
				}
			}
			entities = append(entities, PersonalDetailEntity{
				EntityID:            projectLocal.Key,
				EntityType:          "roles",
				EntityStatus:        string(userLocal.UserType),
				EntityName:          string(userLocal.UserType),
				ProjectName:         projectLocal.Name,
				ProjectID:           projectLocal.Key,
				DisableDeleteReason: disableDeleteReason,
			})
		}
	case "logs":
		logy.Infof(requestSession, "Fetching logs for user: %s", userName)

		mockTraces := []struct {
			traceID   string
			traceType string
		}{
			{"trace-004", "USER_SESSION_LOG"},
			{"trace-006", "EXPORT_HISTORY"},
			{"trace-007", "LOGIN_ACTIVITY"},
		}
		for _, trace := range mockTraces {
			entities = append(entities, PersonalDetailEntity{
				EntityID:      trace.traceID,
				EntityType:    "logs",
				EntitySubType: trace.traceType,
				EntityStatus:  "",
				EntityName:    trace.traceType,
			})
		}
	default:
		render.JSON(w, r, PersonalDetailsResponse{
			Success: false,
			Message: message.GetI18N(message.ErrorKeyRequestParamNotValid, "entity").Code,
			Data:    []PersonalDetailEntity{},
		})
		return
	}

	render.JSON(w, r, PersonalDetailsResponse{
		Success: true,
		Message: "SUCCESS",
		Data:    entities,
	})
}

func (handler *UserHandler) DeletePersonalDataByEntityIdHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)

	if !rights.IsDomainAdmin() {
		exception.ThrowExceptionSendDeniedResponse()
	}

	entity := chi.URLParam(r, "entity")
	id := chi.URLParam(r, "id")

	if entity == "" || id == "" {
		render.JSON(w, r, SuccessResponse{
			Success: false,
			Message: message.GetI18N(message.ErrorKeyRequestParamNotValid).Code,
		})
		return
	}

	validEntityTypes := map[string]bool{
		"tasks": true,
		"roles": true,
		"logs":  true,
	}

	if !validEntityTypes[entity] {
		logy.Warnf(requestSession, "Invalid entity type provided: %s", entity)
		render.JSON(w, r, SuccessResponse{
			Success: false,
			Message: message.GetI18N(message.ErrorKeyRequestParamNotValid, "entity").Code,
		})
		return
	}

	logy.Infof(requestSession, "Deleting entity: %s with ID: %s", entity, id)

	switch entity {
	case "tasks":
		logy.Infof(requestSession, "Deleting task with ID: %s", id)
	case "roles":
		logy.Infof(requestSession, "Deleting role with project ID: %s", id)
	case "logs":
		logy.Infof(requestSession, "Deleting log with ID: %s", id)
	}

	render.JSON(w, r, SuccessResponse{
		Success: true,
		Message: message.GetI18N(message.UserManagementEntityDeleted).Code,
	})
}

func (handler *UserHandler) DeletePersonalDataByEntityHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)

	if !rights.IsDomainAdmin() {
		exception.ThrowExceptionSendDeniedResponse()
	}

	entity := chi.URLParam(r, "entity")

	if entity == "" {
		render.JSON(w, r, SuccessResponse{
			Success: false,
			Message: message.GetI18N(message.ErrorKeyRequestParamNotValid, "entity").Code,
		})
		return
	}

	userName := r.URL.Query().Get("username")
	if userName == "" {
		render.JSON(w, r, SuccessResponse{
			Success: false,
			Message: message.GetI18N(message.ErrorKeyRequestParamNotValid, "username").Code,
		})
		return
	}

	validEntityTypes := map[string]bool{
		"tasks": true,
		"roles": true,
		"logs":  true,
	}

	if !validEntityTypes[entity] {
		logy.Warnf(requestSession, "Invalid entity type provided: %s", entity)
		render.JSON(w, r, SuccessResponse{
			Success: false,
			Message: message.GetI18N(message.ErrorKeyRequestParamNotValid, "entity").Code,
		})
		return
	}

	logy.Infof(requestSession, "Deleting all entities of type: %s for user: %s", entity, userName)

	switch entity {
	case "tasks":
		logy.Infof(requestSession, "Deleting all tasks for user: %s", userName)
	case "roles":
		logy.Infof(requestSession, "Deleting all roles for user: %s", userName)
	case "logs":
		logy.Infof(requestSession, "Deleting all logs for user: %s", userName)
	}

	render.JSON(w, r, SuccessResponse{
		Success: true,
		Message: message.GetI18N(message.UserManagementEntityDeleted).Code,
	})
}

func (handler *UserHandler) GetDeletionAuditTrailHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)

	if !rights.IsDomainAdmin() {
		exception.ThrowExceptionSendDeniedResponse()
	}

	username := chi.URLParam(r, "username")
	if username == "" {
		render.JSON(w, r, SuccessResponse{
			Success: false,
			Message: message.GetI18N(message.ErrorKeyRequestParamNotValid, "username").Code,
		})
		return
	}

	entries := handler.DeletionAuditService.GetAuditTrailForUser(username)
	render.JSON(w, r, entries)
}

// GetDeletionAuditByOperationHandler returns a grouped summary of a deletion operation.
func (handler *UserHandler) GetDeletionAuditByOperationHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)

	if !rights.IsDomainAdmin() {
		exception.ThrowExceptionSendDeniedResponse()
	}

	operationId := chi.URLParam(r, "operationId")
	if operationId == "" {
		render.JSON(w, r, SuccessResponse{
			Success: false,
			Message: message.GetI18N(message.ErrorKeyRequestParamNotValid, "operationId").Code,
		})
		return
	}

	summary := handler.DeletionAuditService.GetAuditByOperation(operationId)
	if summary == nil {
		render.JSON(w, r, SuccessResponse{
			Success: false,
			Message: "No audit entries found for this operation",
		})
		return
	}

	render.JSON(w, r, summary)
}

// GetDeletionAuditByAdminHandler returns all deletion audit entries performed by a given admin.
func (handler *UserHandler) GetDeletionAuditByAdminHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)

	if !rights.IsDomainAdmin() {
		exception.ThrowExceptionSendDeniedResponse()
	}

	adminUser := chi.URLParam(r, "adminUser")
	if adminUser == "" {
		render.JSON(w, r, SuccessResponse{
			Success: false,
			Message: message.GetI18N(message.ErrorKeyRequestParamNotValid, "adminUser").Code,
		})
		return
	}

	entries := handler.DeletionAuditService.GetAuditTrailByAdmin(adminUser)
	render.JSON(w, r, entries)
}
