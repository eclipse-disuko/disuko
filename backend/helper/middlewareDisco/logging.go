// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package middlewareDisco

import (
	"fmt"
	"net/http"
	"time"

	"github.com/eclipse-disuko/disuko/logy"
	"github.com/go-chi/chi/v5/middleware"
)

func Logging(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		argsMap := map[string]any{
			logy.MsgStage: logy.MsgStageStartRequest,
			"method":      r.Method,
			"schema":      "http",
			"host":        r.Host,
			"requestURI":  r.RequestURI,
			"protocol":    r.Proto,
			"remoteAddr":  r.RemoteAddr,
		}
		if r.TLS != nil {
			argsMap["schema"] = "https"
		}
		keysAndValues := make([]any, 0)
		for key, value := range argsMap {
			keysAndValues = append(keysAndValues, key)
			keysAndValues = append(keysAndValues, value)
		}
		beforeHandling := fmt.Sprintf("\"%s %s://%s%s %s from %s\"",
			argsMap["method"],
			argsMap["schema"],
			argsMap["host"],
			argsMap["requestURI"],
			argsMap["protocol"],
			argsMap["remoteAddr"],
		)

		rs := logy.GetRequestSession(r)
		rs.Infow(fmt.Sprintf("%s %s", logy.MsgStageStartRequest, beforeHandling), keysAndValues...)

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		startServe := time.Now()
		defer func() {
			duration := time.Since(startServe)
			argsMap[logy.MsgStage] = logy.MsgStageFinishRequest
			argsMap["status"] = ww.Status()
			argsMap["bytesWritten"] = ww.BytesWritten()
			argsMap["duration"] = duration
			argsMap["queries"] = rs.QueryCount
			argsMap["queryTime"] = rs.QueryTime

			keysAndValues = make([]any, 0)
			keysAndValues = append(keysAndValues, "durationNumeric")
			keysAndValues = append(keysAndValues, int64(duration/time.Millisecond))
			for key, value := range argsMap {
				keysAndValues = append(keysAndValues, key)
				keysAndValues = append(keysAndValues, value)
			}

			afterHandling := fmt.Sprintf("%s - %d %dB in %s. ",
				beforeHandling,
				argsMap["status"],
				argsMap["bytesWritten"],
				argsMap["duration"],
			)
			rs.Infow(fmt.Sprintf("%s %s", logy.MsgStageFinishRequest, afterHandling), keysAndValues...)
		}()
		next.ServeHTTP(ww, logy.RequestWithSession(r, rs))
	}
	return http.HandlerFunc(fn)
}
