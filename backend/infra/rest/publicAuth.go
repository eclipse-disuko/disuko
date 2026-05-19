// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"
	"time"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/publicauth"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/validation"
	project2 "github.com/eclipse-disuko/disuko/infra/repository/project"
	projectRepo "github.com/eclipse-disuko/disuko/infra/repository/project"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v4"
)

var (
	refreshSubject = "public refresh"
	accesSubject   = "public access"
	patSubject     = "pat"
)

type PublicAuthHandler struct {
	ProjectRepo projectRepo.IProjectRepository
}

func extractLoginReq(r *http.Request) *publicauth.LoginReq {
	var item publicauth.LoginReq
	validation.DecodeAndValidate(r, &item, false)
	return &item
}

func (h *PublicAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	rs := logy.GetRequestSession(r)
	req := extractLoginReq(r)
	if err := validation.CheckUuid(req.ProjectUUID); err != nil {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid project UUID"), err.Error())
	}
	if err := validation.CheckUuid(req.Token); err != nil {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid token"), err.Error())
	}
	pr := h.ProjectRepo.FindByKey(rs, req.ProjectUUID, true)
	if pr == nil {
		exception.ThrowExceptionSendDeniedResponse()
	}

	if pr.IsDeprecated() {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DeprecatedProjectError), "")
	}

	expired := pr.ExpireTokens()
	if expired {
		full := h.ProjectRepo.FindByKey(rs, pr.Key, false)
		full.Token = pr.Token
		h.ProjectRepo.Update(rs, full)
	}
	prToken := projectTokenAuth(rs, h.ProjectRepo, pr, req.Token)

	accessClaims := publicauth.AccessClaims{
		ProjectUUID: pr.Key,
		TokenKey:    prToken.Key,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   accesSubject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(conf.Config.PublicAuth.AccessTTLSeconds) * time.Second)),
		},
	}
	refreshClaims := publicauth.RefreshClaims{
		ProjectUUID: pr.Key,
		TokenKey:    prToken.Key,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   refreshSubject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(conf.Config.PublicAuth.RefreshTTLMinutes) * time.Minute)),
		},
	}

	accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessSigned, err := accessJWT.SignedString([]byte(conf.Config.PublicAuth.SigningKey))
	if err != nil {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.Error, "Signing access token"), err.Error())
	}

	refreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshSigned, err := refreshJWT.SignedString([]byte(conf.Config.PublicAuth.SigningKey))
	if err != nil {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.Error, "Signing refresh token"), err.Error())
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access",
		Value:    accessSigned,
		Path:     "/api/public/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh",
		Value:    refreshSigned,
		Path:     "/api/public/auth/refresh",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusOK)
}

func (h *PublicAuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	rs := logy.GetRequestSession(r)

	cookie, err := r.Cookie("refresh")
	if err != nil {
		exception.ThrowExceptionBadRequestResponse()
	}

	token, err := jwt.ParseWithClaims(cookie.Value, &publicauth.RefreshClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(conf.Config.PublicAuth.SigningKey), nil
	})
	if err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.DiscoTokenUnauthorized), err)
	}
	claims, ok := token.Claims.(*publicauth.RefreshClaims)
	if !ok {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid access token"), "Unexpected claims")
	}
	if claims.Subject != refreshSubject {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid access token"), "Unexpected claims")
	}

	pr := h.ProjectRepo.FindByKey(rs, claims.ProjectUUID, true)
	if pr == nil {
		exception.ThrowExceptionClient404Message(message.GetI18N(message.ErrorDbRead, project2.ProjectCollectionName), "project not found: "+claims.ProjectUUID)
	}
	if pr.IsDeprecated() {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DeprecatedProjectError), "")
	}

	expired := pr.ExpireTokens()
	if expired {
		full := h.ProjectRepo.FindByKey(rs, pr.Key, false)
		full.Token = pr.Token
		h.ProjectRepo.Update(rs, full)
	}
	projectTokenKeyAuth(rs, h.ProjectRepo, pr, claims.TokenKey)

	accessClaims := publicauth.AccessClaims{
		ProjectUUID: pr.Key,
		TokenKey:    claims.TokenKey,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   accesSubject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(conf.Config.PublicAuth.AccessTTLSeconds) * time.Second)),
		},
	}
	accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessSigned, err := accessJWT.SignedString([]byte(conf.Config.PublicAuth.SigningKey))
	if err != nil {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.Error, "Signing access token"), err.Error())
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access",
		Value:    accessSigned,
		Path:     "/api/public/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	w.WriteHeader(http.StatusOK)
}

func (h *PublicAuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access",
		Value:    "",
		Path:     "/api/public/",
		Expires:  time.Now().Add(-24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh",
		Value:    "",
		Path:     "/api/public/auth/refresh",
		Expires:  time.Now().Add(-24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	w.WriteHeader(http.StatusOK)
}

func (h *PublicAuthHandler) Info(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("access")
	if err != nil {
		exception.ThrowExceptionBadRequestResponse()
	}
	token, err := jwt.ParseWithClaims(cookie.Value, &publicauth.AccessClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(conf.Config.PublicAuth.SigningKey), nil
	})
	if err != nil {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid access token"), err.Error())
	}
	claims, ok := token.Claims.(*publicauth.AccessClaims)
	if !ok {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid access token"), "Unexpected claims")
	}
	if claims.Subject != accesSubject {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid access token"), "Unexpected claims")
	}
	requestSession := logy.GetRequestSession(r)
	pr := h.ProjectRepo.FindByKey(requestSession, claims.ProjectUUID, true)
	if pr == nil {
		exception.ThrowExceptionClient404Message(message.GetI18N(message.ErrorDbRead, project2.ProjectCollectionName), "project not found: "+claims.ProjectUUID)
	}
	res := publicauth.InfoDto{
		ProjectUUID: pr.Key,
		IsGroup:     pr.IsGroup,
	}
	render.JSON(w, r, res)
}

func projectAccessAuth(rs *logy.RequestSession, repo projectRepo.IProjectRepository, pr *project.Project, cookie *http.Cookie) string {
	token, err := jwt.ParseWithClaims(cookie.Value, &publicauth.AccessClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(conf.Config.PublicAuth.SigningKey), nil
	})
	if err != nil {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid access token"), err.Error())
	}
	claims, ok := token.Claims.(*publicauth.AccessClaims)
	if !ok {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid access token"), "Unexpected claims")
	}
	if claims.Subject != accesSubject {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid access token"), "Unexpected claims")
	}
	if claims.ProjectUUID != pr.Key && claims.ProjectUUID != pr.Parent {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid access token"), "Project UUID mismatch")
	}
	return projectTokenKeyAuth(rs, repo, pr, claims.TokenKey)
}

func projectTokenKeyAuth(rs *logy.RequestSession, repo projectRepo.IProjectRepository, pr *project.Project, key string) string {
	prToken := pr.GetActiveTokenByKey(key)
	if prToken != nil {
		return prToken.Origin()
	}
	if !pr.HasParent() {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid access token"), "Project token not found or expired")
	}
	parent := repo.FindByKey(rs, pr.Parent, false)
	if parent == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorDbNotFound), "Parent project not found")
	}
	expired := parent.ExpireTokens()
	if expired {
		repo.Update(rs, parent)
	}
	parentToken := parent.GetActiveTokenByKey(key)
	if parentToken == nil {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid access token"), "Project token not found or expired")
	}
	return parentToken.Origin()
}

func projectTokenAuth(rs *logy.RequestSession, repo projectRepo.IProjectRepository, pr *project.Project, token string) *project.Token {
	prToken := pr.GetActiveToken(token)
	if prToken != nil {
		return prToken
	}
	if !pr.HasParent() {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid disco token"), "Project token not found or expired")
	}
	parent := repo.FindByKey(rs, pr.Parent, false)
	if parent == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorDbNotFound), "Parent project not found")
	}
	expired := parent.ExpireTokens()
	if expired {
		repo.Update(rs, parent)
	}
	parentToken := parent.GetActiveToken(token)
	if parentToken == nil {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid disco token"), "Project token not found or expired")
	}
	return parentToken
}
