// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package middlewareDisco

import (
	"context"
	"net/http"
	"strings"

	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/infra/service/patauth"
	"github.com/eclipse-disuko/disuko/logy"
)

type CtxKey string

var PATUserKey CtxKey = "pat"

type PATAuthMW struct {
	service *patauth.Service
}

func InitPATAuthMW(service *patauth.Service) *PATAuthMW {
	return &PATAuthMW{
		service: service,
	}
}

func (mw *PATAuthMW) Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestSession := logy.GetRequestSession(r)

		raw := r.Header.Get("Authorization")
		trimmed := strings.TrimPrefix(raw, "Bearer ")
		if trimmed == raw {
			exception.ThrowExceptionBadRequestResponse()
		}

		user, _ := mw.service.Validate(requestSession, trimmed)

		ctx := context.WithValue(r.Context(), PATUserKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
