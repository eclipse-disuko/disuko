// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/google/go-cmp/cmp"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain/audit"
	"github.com/eclipse-disuko/disuko/domain/oauth"
	"github.com/eclipse-disuko/disuko/domain/user"
	auditHelper "github.com/eclipse-disuko/disuko/helper/audit"
	"github.com/eclipse-disuko/disuko/helper/exception"
	jwt2 "github.com/eclipse-disuko/disuko/helper/jwt"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/roles"
	user2 "github.com/eclipse-disuko/disuko/infra/repository/user"
	"github.com/eclipse-disuko/disuko/logy"
	"go.uber.org/zap/zapcore"

	// "encoding/json"

	"github.com/eclipse-disuko/disuko/helper"

	// "gh-bot-platform/oauth"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"golang.org/x/oauth2"
)

type OAuthHandler struct {
	HttpClient     *http.Client
	Ctx            context.Context
	Config         oauth2.Config
	Verifier       *oidc.IDTokenVerifier
	Provider       *oidc.Provider
	UserRepository user2.IUsersRepository
	ProviderURL    string
}

type AuthKey struct {
	Auth jwtauth.JWTAuth
	Kid  string
}

func (handler *OAuthHandler) HandleLoginWithToken(writer http.ResponseWriter, r *http.Request) {
	// todo: does the handler still necessary?
	requestSession := logy.GetRequestSession(r)
	logy.Infof(requestSession, fmt.Sprintf("oauthHandler::HandleLoginWithToken"))
	response := ""

	_, err := writer.Write([]byte(response))
	exception.HandleErrorClientMessage(err, message.GetI18N(message.WritingContent))
}

func (handler *OAuthHandler) HandleRedirectToIAMLogout(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	logy.Infof(requestSession, "oauthHandler::HandleRedirectToIAMLogout")

	cookieRefreshToken := expireRefreshCookie()
	cookieAccessToken := expireAccessCookie()
	http.SetCookie(w, &cookieRefreshToken)
	http.SetCookie(w, &cookieAccessToken)
	http.Redirect(w, r, conf.Config.OAuth2.LogoutEndpoint, http.StatusMovedPermanently)
}

// HandleRedirectToIAM
// This endpoint is called from IAM via request_uri.
func (handler *OAuthHandler) HandleRedirectToIAM(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	logy.Infof(requestSession, "oauthHandler::HandleRedirectToIAM")

	authEndpoint := handler.Config.Endpoint.AuthURL
	if conf.Config.OAuth2.AuthorizationEndpoint != "" {
		authEndpoint = conf.Config.OAuth2.AuthorizationEndpoint
	}

	newUrl := authEndpoint + "?"
	newUrl = newUrl + "response_type=code"
	newUrl = newUrl + "&client_id=" + conf.Config.OAuth2.ClientId
	newUrl = newUrl + "&scope=" + strings.Join(handler.Config.Scopes, "%20")
	newUrl = newUrl + "&redirect_uri=" + url.QueryEscape(conf.Config.Server.DisukoHost+"/api/v1/login")
	logy.Infof(requestSession, "oauthHandler::HandleRedirectToIAM %s", newUrl)
	http.Redirect(w, r, newUrl, http.StatusMovedPermanently)
}

// HandleRequestTokenFromCode
// This endpoint is called from IAM via request_uri.
func (handler *OAuthHandler) HandleRequestTokenFromCode(writer http.ResponseWriter, request *http.Request) {
	requestSession := logy.GetRequestSession(request)
	logy.Infof(requestSession, "oauthHandler::HandleRequestTokenFromCode")

	errorParam := request.URL.Query().Get("error")
	errorDescription := request.URL.Query().Get("error_description")
	if errorParam != "" {
		err := errors.New(fmt.Sprintf("Authentication error: %s - %s", errorParam, errorDescription))
		logErrorAndRedirectToErrorPage(requestSession, writer, request, zapcore.WarnLevel, message.AuthErrorCode, message.GetI18N(message.AuthErrorCode).Text, err)
		return
	}

	code := html.EscapeString(request.URL.Query().Get("code"))
	logy.Infof(requestSession, "oauthHandler::HandleRequestTokenFromCode %v", code)

	codeErr := validateCode(code)
	if codeErr != nil {
		logErrorAndRedirectToErrorPage(requestSession, writer, request, zapcore.ErrorLevel, message.InvalidExchangeCode, message.GetI18N(message.InvalidExchangeCode).Text, codeErr)
		return
	}

	handler.Config.RedirectURL = conf.Config.OAuth2.RedirectURL
	oauth2Token, err := handler.Config.Exchange(handler.Ctx, code)
	logy.Infof(requestSession, "oauthHandler::HandleRequestTokenFromCode %v", conf.Config.OAuth2.RedirectURL)
	if err != nil {
		logErrorAndRedirectToErrorPage(requestSession, writer, request, zapcore.ErrorLevel, message.CreateTokenAuth, message.GetI18N(message.CreateTokenAuth).Text, err)
		return
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		exception.WriteErrorWithCode(requestSession, &writer, message.IdTokenMissing, message.GetI18N(message.IdTokenMissing).Text, "Missing Token", http.StatusUnauthorized, zapcore.WarnLevel)
		return
	}

	// Extract the Access Token from OAuth2 token.
	accessToken, ok := oauth2Token.Extra("access_token").(string)
	if !ok {
		exception.WriteErrorWithCode(requestSession, &writer, message.AccessTokenMissing, message.GetI18N(message.AccessTokenMissing).Text, "Missing access token in oauth answer", http.StatusUnauthorized, zapcore.WarnLevel)
		return
	}

	// Parse and verify ID Token payload.
	idToken, err := handler.Verifier.Verify(handler.Ctx, rawIDToken)
	if err != nil {
		exception.WriteErrorWithCode(requestSession, &writer, message.VerifyError, message.GetI18N(message.VerifyError).Text, err.Error(), http.StatusUnauthorized, zapcore.InfoLevel)
		return
	}
	token := oauth2.Token{
		AccessToken: accessToken,
		Expiry:      oauth2Token.Expiry,
		TokenType:   "Bearer",
	}

	tokenSource := handler.Config.TokenSource(context.Background(), &token)
	userInfo, err := handler.Provider.UserInfo(context.Background(), tokenSource)
	if err != nil {
		exception.WriteErrorWithCode(requestSession, &writer, message.Verify, message.GetI18N(message.Verify).Text, err.Error(), http.StatusUnauthorized, zapcore.ErrorLevel)
		return
	}

	// Extract custom claims
	var claims struct {
		Lastname              string   `json:"last_name"`
		Forename              string   `json:"first_name"`
		Email                 string   `json:"email"`
		Verified              bool     `json:"email_verified"`
		GroupType             string   `json:"group_type"`
		ObjectClass           []string `json:"object_class"`
		CompanyIdentifier     string   `json:"company_identifier"`
		Department            string   `json:"department"`
		DepartmentDescription string   `json:"department_description"`
	}
	if err := idToken.Claims(&claims); err != nil {
		exception.WriteErrorWithCode(requestSession, &writer, message.VerifyClaims, message.GetI18N(message.VerifyClaims).Text, err.Error(), http.StatusForbidden, zapcore.ErrorLevel)
		return
	}

	if claims.GroupType != jwt2.GROUP_TYPE_DAIMLER {
		exception.WriteErrorWithCode(requestSession, &writer, message.Unauthorized, message.GetI18N(message.Unauthorized).Text, "User unauthorized", http.StatusUnauthorized, zapcore.WarnLevel)
		return
	}

	isInternalEmployee := helper.Contains("dcxInternalEmployee", claims.ObjectClass) || helper.Contains(userInfo.Subject, conf.Config.InternalUsersAllowList)

	userInfoClaims := extractUserInfoWithGroups(userInfo)

	existingUser := handler.UserRepository.FindByUserId(requestSession, userInfo.Subject)
	metaData := user.NewMetaData(claims.CompanyIdentifier, claims.Department, claims.DepartmentDescription)
	if existingUser == nil {
		existingUser = user.CreateUser(claims.Forename, claims.Lastname, userInfo.Subject, userInfo.Email, userInfoClaims.EntitlementGroup, metaData, isInternalEmployee)
		newAudit := existingUser.ToUserAudit()
		auditHelper.CreateAndAddAuditEntry(&existingUser.Container, "SYSTEM", message.UserCreated, audit.DiffWithReporter, newAudit, &user.UserAudit{})
		handler.UserRepository.Save(requestSession, existingUser)
	} else {
		before := *existingUser
		existingUser.Lastname = claims.Lastname
		existingUser.Forename = claims.Forename
		existingUser.Roles = userInfoClaims.EntitlementGroup
		existingUser.MetaData = metaData
		existingUser.IsInternal = isInternalEmployee
		if before.Lastname != existingUser.Lastname ||
			before.Forename != existingUser.Forename ||
			!helper.EqualsStringSlicesIgnoreOrder(before.Roles, existingUser.Roles) ||
			!before.MetaData.Equal(existingUser.MetaData) ||
			before.IsInternal != existingUser.IsInternal {
			auditHelper.CreateAndAddAuditEntry(&existingUser.Container, "SYSTEM", message.UserUpdated, cmp.Diff, *existingUser, before)
		}

		handler.UserRepository.Update(requestSession, existingUser)
	}
	userData := jwt2.CreateUserData(requestSession, existingUser, userInfo.Subject, userInfo.Email, userInfoClaims.EntitlementGroup, claims.GroupType, isInternalEmployee, request)
	tokenDetails := jwt2.CreateToken(userData)
	accessData := roles.GetAccessAndRolesRightsFromClaim(userData)
	response := oauth.OAuthTokenResponse{
		Rights:  accessData,
		Profile: existingUser.ToDto(),
	}
	responseJSON, err := json.MarshalIndent(response, "", "  ")
	fmt.Println(responseJSON)
	if err != nil {
		fmt.Println(err)
		return
	}
	cookieRefreshToken := createRefreshCookie(tokenDetails.RefreshToken)
	cookieAccessToken := createAccessCookie(tokenDetails.AccessToken)
	http.SetCookie(writer, &cookieRefreshToken)
	http.SetCookie(writer, &cookieAccessToken)

	redirectUrl := conf.Config.Server.ClientRedirectURL + "/#/oauth/callback"
	logy.Infof(requestSession, "oauthHandler::HandleRequestTokenFromCode redirectUrl %v", redirectUrl)
	http.Redirect(writer, request, redirectUrl, http.StatusMovedPermanently)
}

func logErrorAndRedirectToErrorPage(requestSession *logy.RequestSession, writer http.ResponseWriter, request *http.Request, level zapcore.Level, errMsgCode string, errMsgText string, err error) {
	exception.LogWithLevel(requestSession, level, errMsgCode, errMsgText, err.Error())
	cookieRefreshToken := expireRefreshCookie()
	cookieAccessToken := expireAccessCookie()
	http.SetCookie(writer, &cookieRefreshToken)
	http.SetCookie(writer, &cookieAccessToken)
	u := conf.Config.Server.ClientRedirectURL + "/#/loginError"
	http.Redirect(writer, request, u, http.StatusMovedPermanently)
	return
}

func validateCode(code string) error {
	// Reject if the input is a URL
	if _, err := url.ParseRequestURI(code); err == nil {
		return errors.New("invalid code: input appears to be a URL")
	}

	codeRegex := regexp.MustCompile(conf.Config.OAuth2.RegexToken)
	if !codeRegex.MatchString(code) {
		return errors.New("invalid code: does not match the expected format")
	}
	return nil
}

func (handler *OAuthHandler) HandleRefreshToken(writer http.ResponseWriter, request *http.Request) {
	requestSession := logy.GetRequestSession(request)
	logy.Infof(requestSession, "oauthHandler::HandleRefreshToken")
	// get refresh token
	var refreshTokenRequest oauth.RefreshTokenRequestDto
	cookie, err := request.Cookie("oauth.r")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			// No cookie was found
			exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.ErrorAAR), "No refresh cookie found")
			return
		}
		exception.ThrowExceptionSendDeniedResponse()
		return
	}
	cookie.HttpOnly = true
	cookie.Secure = true
	refreshTokenRequest.RefreshToken = cookie.Value
	refreshTokenStr := refreshTokenRequest.RefreshToken

	// check if refresh token is valid
	refreshToken := jwt2.ParseAndVerifyToken(requestSession, request, refreshTokenStr, conf.Config.Auth.RefreshSecret)

	// check if token in time
	refreshTokenClaims := jwt2.ExtractClaims(refreshToken)
	refreshTokenExpireUnixSec := jwt2.ExtractClaimInt64(refreshTokenClaims, jwt2.JK_EXPIRES)
	refreshTokenExpireTime := time.Unix(refreshTokenExpireUnixSec, 0)
	if refreshTokenExpireTime.Before(time.Now()) {
		// may not as exception, instead simple send a 401?
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.ErrorAAR, "Invalid token"),
			"refresh token is expired")
	}

	// check if refresh token and the access token belongs to the same user
	accessToken := jwt2.ExtractAndVerifyAccessToken(requestSession, request)
	accessTokenClaims := jwt2.ExtractClaims(accessToken)
	accessTokenUser := jwt2.ExtractClaimStr(accessTokenClaims, jwt2.JK_USERNAME)
	refreshTokenUser := jwt2.ExtractClaimStr(refreshTokenClaims, jwt2.JK_USERNAME)
	if accessTokenUser != refreshTokenUser {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.ErrorAAR, "Invalid token"),
			"refresh token and access toke belongs not to same user")
	}

	// create new access token
	accessTokenStrNew, _ := jwt2.CreateAccessTokenStrFromClaims(accessTokenClaims)
	cookieRefreshToken := createRefreshCookie(refreshTokenStr)
	cookieAccessToken := createAccessCookie(accessTokenStrNew)
	http.SetCookie(writer, &cookieRefreshToken)
	http.SetCookie(writer, &cookieAccessToken)
	result := domain.BaseResponseStatus{
		Success: true,
	}
	render.JSON(writer, request, result)
}

func extractUserInfoWithGroups(userInfo *oidc.UserInfo) *oauth.UserInfoWithGroups {
	userInfoClaimsDto := oauth.UserInfoWithGroupsDto{}
	err := userInfo.Claims(&userInfoClaimsDto)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorJsonUnmarshall))
	return convertUserInfoWithGroupsDto(&userInfoClaimsDto)
}

func convertUserInfoWithGroupsDto(dto *oauth.UserInfoWithGroupsDto) *oauth.UserInfoWithGroups {
	userInfoClaims := oauth.UserInfoWithGroups{
		UserInfo: oauth.UserInfo{
			Sub:         dto.Sub,
			AppId:       dto.AppId,
			AccessGroup: dto.AccessGroup,
		},
		EntitlementGroup: []string{},
		GroupType:        dto.GroupType,
	}

	switch groups := dto.EntitlementGroup.(type) {
	case string:
		// 1 role
		userInfoClaims.EntitlementGroup = []string{groups}
	case []interface{}:
		// multiple roles
		groupsStrings := make([]string, len(groups))
		for i, v := range groups {
			groupsStrings[i] = v.(string)
		}
		userInfoClaims.EntitlementGroup = groupsStrings
	case nil:
		// 0 roles
	default:
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorUnexpectedType), "")
	}
	return &userInfoClaims
}

func createRefreshCookie(refreshToken string) http.Cookie {
	cookieRefreshToken := createCookie("oauth.r", refreshToken, time.Time{})
	return cookieRefreshToken
}

func createAccessCookie(accessToken string) http.Cookie {
	cookieAccessToken := createCookie("oauth.a", accessToken, time.Time{})
	return cookieAccessToken
}

func expireRefreshCookie() http.Cookie {
	expiredTime := time.Unix(0, 0)
	cookieRefreshToken := createCookie("oauth.r", "", expiredTime)
	return cookieRefreshToken
}

func expireAccessCookie() http.Cookie {
	expiredTime := time.Unix(0, 0)
	cookieAccessToken := createCookie("oauth.a", "", expiredTime)
	return cookieAccessToken
}

func createCookie(name, value string, expires time.Time) http.Cookie {
	return http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  expires,
	}
}
