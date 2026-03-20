// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package middlewareDisco

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"mercedes-benz.ghe.com/foss/disuko/domain/internalToken"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	basicauthRepo "mercedes-benz.ghe.com/foss/disuko/infra/repository/internalToken"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type CtxKey string

var TokenKey CtxKey = "auth"

type InternalTokenMW struct {
	repo basicauthRepo.IRepo
}

func InitInternalTokenMW(repo basicauthRepo.IRepo) *InternalTokenMW {
	return &InternalTokenMW{
		repo: repo,
	}
}

func (b *InternalTokenMW) Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestSession := logy.GetRequestSession(r)

		raw := r.Header.Get("Authorization")

		trimmed := strings.TrimPrefix(raw, "Bearer ")
		if trimmed == raw {
			exception.ThrowExceptionBadRequestResponse()
		}

		bb, err := base64.RawStdEncoding.DecodeString(trimmed)
		if err != nil {
			exception.ThrowExceptionBadRequestResponse()
		}

		var auth internalToken.TokenAuth
		if err := json.Unmarshal(bb, &auth); err != nil {
			exception.ThrowExceptionBadRequestResponse()
		}
		it := b.repo.CheckAuth(requestSession, auth.Key, auth.Token)
		if it == nil {
			exception.ThrowExceptionSendDeniedResponse()
		}

		ctx := context.WithValue(r.Context(), TokenKey, it)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
