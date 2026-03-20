// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
)

func GetURLParam(r *http.Request, paramName string) string {
	taskGuidEscaped := chi.URLParam(r, paramName)
	result, err := url.QueryUnescape(taskGuidEscaped)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ErrorUnexpectError))
	return result
}
