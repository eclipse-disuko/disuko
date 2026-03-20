// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package middlewareDisco

import (
	"encoding/json"
	"net/http"
	"strings"

	"mercedes-benz.ghe.com/foss/disuko/domain/notification"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

var CurrentNotification *notification.Notification

func CustomHeaders(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if CurrentNotification == nil || strings.HasPrefix(r.URL.RequestURI(), "/api/public/") {
			next.ServeHTTP(w, r)
			return
		}
		notificationHeader, err := json.Marshal(CurrentNotification.ToDto())
		if err != nil {
			requestSession := logy.GetRequestSession(r)
			logy.Errorf(requestSession, "%s", err)
		} else {
			w.Header().Set("X-Notification", string(notificationHeader))
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
