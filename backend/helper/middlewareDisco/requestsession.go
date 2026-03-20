// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package middlewareDisco

import (
	"net/http"

	"mercedes-benz.ghe.com/foss/disuko/logy"
)

func RequestSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rs := logy.NewRequestSession()
		next.ServeHTTP(w, logy.RequestWithSession(r, rs))
	})
}
