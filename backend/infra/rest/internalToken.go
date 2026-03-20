// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"mercedes-benz.ghe.com/foss/disuko/domain/internalToken"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/helper/middlewareDisco"
	"mercedes-benz.ghe.com/foss/disuko/helper/roles"
	"mercedes-benz.ghe.com/foss/disuko/helper/validation"
	internalTokenRepo "mercedes-benz.ghe.com/foss/disuko/infra/repository/internalToken"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

func extractInternalToken(ctx context.Context) *internalToken.InternalToken {
	raw := ctx.Value(middlewareDisco.TokenKey)
	b, ok := raw.(*internalToken.InternalToken)
	if !ok {
		exception.ThrowExceptionSendDeniedResponse()
	}
	return b
}

type InternalTokenHandler struct {
	InternalTokenRepo internalTokenRepo.IRepo
}

func extractInternalTokenKeyFromRequest(r *http.Request) string {
	basicAuthUUID := chi.URLParam(r, "uuid")

	authUUID, err := url.QueryUnescape(basicAuthUUID)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamUuidWrong))

	err = validation.CheckUuid(authUUID)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamUuidWrong))
	return authUUID
}

func extractInternalTokenBody(r *http.Request) internalToken.InternalTokenDto {
	var token internalToken.InternalTokenDto
	validation.DecodeAndValidate(r, &token, false)
	return token
}

func (internalTokenHandler *InternalTokenHandler) List(w http.ResponseWriter, r *http.Request) {
	// Retrieve all records from the database using the repository
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.IsApplicationAdmin()) {
		exception.ThrowExceptionSendDeniedResponse()
	}

	// Fetch all records
	basicauths := internalTokenHandler.InternalTokenRepo.FindAll(requestSession, true)

	// Convert the records to DTOs
	var authDtos []internalToken.InternalTokenDto
	for _, entity := range basicauths {
		authDtos = append(authDtos, entity.ToDto())
	}

	// Use render.JSON to send the response
	render.JSON(w, r, authDtos)
}

func (internalTokenHandler *InternalTokenHandler) Create(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.IsApplicationAdmin()) {
		exception.ThrowExceptionSendDeniedResponse()
	}

	req := extractInternalTokenBody(r)

	if req.Expiry.After(time.Now().AddDate(0, 0, 7)) {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.BasicAuthExpiry, ""), "")
	}

	token, sec := internalTokenHandler.InternalTokenRepo.Add(requestSession, req.Name, req.Description, req.Expiry, req.Capabilities)

	authJ, err := json.Marshal(internalToken.TokenAuth{
		Key:   token.Key,
		Token: sec,
	})
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorJsonMarshalling))
	authB64 := base64.StdEncoding.EncodeToString(authJ)
	res := token.ToDto()
	res.Token = authB64
	render.JSON(w, r, res)
}

func (internalTokenHandler *InternalTokenHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.IsApplicationAdmin()) {
		exception.ThrowExceptionSendDeniedResponse()
	}
	key := extractInternalTokenKeyFromRequest(r)
	internalTokenHandler.InternalTokenRepo.Revoke(requestSession, key)
	responseData := SuccessResponse{
		Success: true,
		Message: "Token revoked successfully",
	}
	render.JSON(w, r, responseData)
}

func (internalTokenHandler *InternalTokenHandler) Renew(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.IsApplicationAdmin()) {
		exception.ThrowExceptionSendDeniedResponse()
	}

	key := extractInternalTokenKeyFromRequest(r)
	token, sec := internalTokenHandler.InternalTokenRepo.Renew(requestSession, key)
	authJ, err := json.Marshal(internalToken.TokenAuth{
		Key:   token.Key,
		Token: sec,
	})
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorJsonMarshalling))
	authB64 := base64.StdEncoding.EncodeToString(authJ)
	res := token.ToDto()
	res.Token = authB64
	render.JSON(w, r, res)
}
