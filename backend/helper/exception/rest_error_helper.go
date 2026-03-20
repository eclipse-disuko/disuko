// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package exception

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eclipse-disuko/disuko/logy"
	"go.uber.org/zap/zapcore"
)

type HttpError struct {
	Code    string `json:"code" example:"DISCOTOKEN_UNAUTHORIZED"`
	Message string `json:"message" example:"DISCOTOKEN_UNAUTHORIZED"`
	ReqID   string `json:"reqID" example:"dummy-id----7c1ca33cb0d9c78dd5e9d0"`
	Raw     string `json:"raw" example:"Key: '' Error:Field validation for '' failed on the 'uuid' tag"`
} //	@name	HttpError

type HttpError404 struct {
	Code    string `json:"code" example:"PROJECT_TOKEN"`
	Message string `json:"message" example:"Project uuid wrong"`
	ReqID   string `json:"reqID" example:"dummy-id----4c997f0cc6baaf11f936"`
	Raw     string `json:"raw" example:"Key: '' Error:Field validation for '' failed on the 'uuid' tag"`
} //	@name	HttpError404

// WriteError TODO many of the usage should be use WriteServerError; StatusExpectationFailed is a 4xx error and that means the client has made a mistake/error
func WriteError(requestSession *logy.RequestSession, w *http.ResponseWriter, code string, message string, raw string) {
	WriteErrorWithCode(requestSession, w, code, message, raw, http.StatusExpectationFailed, zapcore.ErrorLevel)
}

func WriteClientError(requestSession *logy.RequestSession, w *http.ResponseWriter, code string, message string, raw string) {
	WriteClientErrorWithHttpCode(requestSession, w, code, message, raw, http.StatusExpectationFailed)
}
func WriteClientErrorWithHttpCode(requestSession *logy.RequestSession, w *http.ResponseWriter, code string, message string, raw string, httpCode int) {
	WriteErrorWithCode(requestSession, w, code, message, raw, httpCode, zapcore.WarnLevel)
}

func WriteServerError(requestSession *logy.RequestSession, w *http.ResponseWriter, code string, message string, raw string) {
	WriteServerErrorWithHttpCode(requestSession, w, code, message, raw, http.StatusInternalServerError)
}
func WriteServerErrorWithHttpCode(requestSession *logy.RequestSession, w *http.ResponseWriter, code string, message string, raw string, httpCode int) {
	WriteErrorWithCode(requestSession, w, code, message, raw, httpCode, zapcore.ErrorLevel)
}
func WriteErrorWithCode(requestSession *logy.RequestSession, w *http.ResponseWriter, code string, message string, raw string, httpCode int, logLevel zapcore.Level) {
	reqID := requestSession.ReqID
	logy.Logw(requestSession, fmt.Sprintf("%s Error during request.", logy.MsgStageIntermediateCommon), logLevel, "message", message, "code", code, "raw", raw, logy.MsgStage, logy.MsgStageIntermediateCommon)

	errorString, err := json.Marshal(HttpError{Code: code, Message: message, ReqID: reqID, Raw: raw})
	if code == "404" {
		errorString, err = json.Marshal(HttpError404{Code: code, Message: message, ReqID: reqID, Raw: raw})
	}
	if err != nil {
		logy.Logw(requestSession, fmt.Sprintf("%s Error marshalling error message for error code: %s raw: %s", logy.MsgStageIntermediateCommon, code, raw), logLevel, "code", code, logy.MsgStage, logy.MsgStageIntermediateCommon)
		http.Error(*w, "Error marshalling error message for error code: "+code, http.StatusExpectationFailed)
		return
	}
	(*w).Header().Set("Content-Type", "application/json; charset=utf-8")
	http.Error(*w, string(errorString), httpCode)
}

func LogWithLevel(requestSession *logy.RequestSession, level zapcore.Level, code string, message string, raw string) {
	logy.Logw(requestSession, fmt.Sprintf("%s Error during request.", logy.MsgStageIntermediateCommon), level, "message", message, "code", code, "raw", raw, logy.MsgStage, logy.MsgStageIntermediateCommon)
}
