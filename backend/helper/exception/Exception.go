// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package exception

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/reflection"
	"github.com/eclipse-disuko/disuko/logy"
	"go.uber.org/zap/zapcore"
)

type ExceptionHandler struct {
	RequestSession *logy.RequestSession
	W              *http.ResponseWriter
	DoNothing      bool `default:"false"`
	IsFatal        bool `default:"false"`
}

type ExceptionType int

const (
	EXCEPTION_TYPE_CLIENT = 0
	EXCEPTION_TYPE_SERVER = 1
)

type Exception struct {
	//musst set
	ErrorCode    string
	ErrorMessage string

	//can be set
	Type          ExceptionType
	HttpErrorCode int
	ErrorRaw      string
	LogLevel      *zapcore.Level
	Error         error
}

const HTTP_CODE_SHOW_NO_REQUEST_ID = 428 //StatusPreconditionRequired

func (exc *Exception) ToString() string {
	str, err := json.Marshal(exc)
	if err != nil {
		return "Error cannot print exception, error: " + err.Error()
	}
	return string(str)
}

func NewExceptionHandler(requestSession *logy.RequestSession, w *http.ResponseWriter) ExceptionHandler {
	return ExceptionHandler{
		RequestSession: requestSession,
		W:              w,
	}
}
func NewExceptionHandler2(requestSession *logy.RequestSession) ExceptionHandler {
	return ExceptionHandler{
		RequestSession: requestSession,
	}
}
func NewExceptionHandlerFatal(requestSession *logy.RequestSession) ExceptionHandler {
	return ExceptionHandler{
		RequestSession: requestSession,
		IsFatal:        true,
	}
}
func CatchException(exceptionHandler ExceptionHandler) {
	if err := recover(); err != nil {
		switch err.(type) {
		case Exception:
			exceptionHandler.HandleException(err.(Exception))
			return
		default:
			logy.Errorf(exceptionHandler.RequestSession, "%s", err)
		}
	}
}

func CatchExceptionWithCustom(exceptionHandler ExceptionHandler, customExceptionHandler func(Exception)) {
	if err := recover(); err != nil {
		switch err.(type) {
		case Exception:
			exceptionHandler.HandleException(err.(Exception))
			customExceptionHandler(err.(Exception))
			return
		default:
			panic(err)
		}
	}
}

func (exceptionHandler ExceptionHandler) HandleException(exception Exception) {

	requestSession := exceptionHandler.RequestSession

	if exceptionHandler.DoNothing {
		//don't take it to hard
		return
	}

	if exceptionHandler.W != nil {
		//we can handle it
		switch exception.Type {
		case EXCEPTION_TYPE_CLIENT:
			if exception.HttpErrorCode != 0 {
				WriteClientErrorWithHttpCode(requestSession, exceptionHandler.W,
					exception.ErrorCode, exception.ErrorMessage, exception.ErrorRaw, exception.HttpErrorCode)
			} else {
				WriteClientError(requestSession, exceptionHandler.W,
					exception.ErrorCode, exception.ErrorMessage, exception.ErrorRaw)
			}
		case EXCEPTION_TYPE_SERVER:
			if exception.HttpErrorCode != 0 {
				WriteServerErrorWithHttpCode(requestSession, exceptionHandler.W,
					exception.ErrorCode, exception.ErrorMessage, exception.ErrorRaw, exception.HttpErrorCode)
			} else {
				WriteServerError(requestSession, exceptionHandler.W,
					exception.ErrorCode, exception.ErrorMessage, exception.ErrorRaw)
			}
		default:
			WriteError(requestSession, exceptionHandler.W,
				exception.ErrorCode, exception.ErrorMessage, exception.ErrorRaw)
		}
	} else {
		LogException(requestSession, exception)
	}
}

func LogException(requestSession *logy.RequestSession, exception Exception) {
	logLevel := zapcore.ErrorLevel
	if exception.LogLevel != nil {
		logLevel = *exception.LogLevel
	}
	logy.Logw(requestSession,
		fmt.Sprintf("%s Error during request.", logy.MsgStageIntermediateCommon),
		logLevel,
		"message", exception.ErrorMessage,
		"code", exception.ErrorCode,
		"raw", exception.ErrorRaw, logy.MsgStage, logy.MsgStageIntermediateCommon)
}

func ThrowException(exception Exception) {
	panic(exception)
}

func ThrowException2(code string, message string, raw string) {
	ThrowException(Exception{
		ErrorCode:    code,
		ErrorMessage: message,
		ErrorRaw:     raw,
	})
}

func ThrowExceptionSendDeniedResponse() {
	ThrowExceptionClientWithHttpCode(message.ErrorAAR, "Access denied", "Access denied",
		http.StatusUnauthorized)
}

func ThrowExceptionSendUserDisabledResponse() {
	ThrowExceptionClientWithHttpCode(message.UserDisabled, "User disabled", "User disabled",
		http.StatusForbidden)
}

func ThrowExceptionConflictResponse() {
	ThrowExceptionClientWithHttpCode(message.Conflict, "Server detected a conflict", "Server detected a conflict",
		http.StatusConflict)
}

func ThrowExceptionBadRequestResponse() {
	ThrowExceptionClientWithHttpCode(message.BadRequest, "Bad request", "Bad request",
		http.StatusBadRequest)
}

func ThrowExceptionSendDeniedResponseRaw(i18n message.I18N, raw string) {
	ThrowException(Exception{
		ErrorCode:     i18n.Code,
		ErrorMessage:  i18n.Text,
		ErrorRaw:      raw,
		HttpErrorCode: http.StatusUnauthorized,
		Type:          EXCEPTION_TYPE_CLIENT,
	})
}

func ThrowExceptionClientMessage3(i18n message.I18N) {
	ThrowExceptionClientMessage(i18n, "")
}
func ThrowExceptionClientMessage(i18n message.I18N, raw string, logLevel ...zapcore.Level) {
	var level *zapcore.Level
	if len(logLevel) > 0 {
		level = &logLevel[0]
	}
	ThrowException(Exception{
		ErrorCode:     i18n.Code,
		ErrorMessage:  i18n.Text,
		ErrorRaw:      raw,
		HttpErrorCode: http.StatusExpectationFailed,
		Type:          EXCEPTION_TYPE_CLIENT,
		LogLevel:      level,
	})
}
func ThrowExceptionClientMessage2(i18n message.I18N, err error, logLevel ...zapcore.Level) {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	var level *zapcore.Level

	if len(logLevel) > 0 {
		level = &logLevel[0]
	}
	ThrowException(Exception{
		ErrorCode:     i18n.Code,
		ErrorMessage:  i18n.Text,
		ErrorRaw:      errStr,
		Error:         err,
		HttpErrorCode: http.StatusExpectationFailed,
		Type:          EXCEPTION_TYPE_CLIENT,
		LogLevel:      level,
	})
}

func ThrowExceptionServerMessage(i18n message.I18N, raw string, logLevel ...zapcore.Level) {
	var level *zapcore.Level
	if len(logLevel) > 0 {
		level = &logLevel[0]
	}

	ThrowException(Exception{
		ErrorCode:     i18n.Code,
		ErrorMessage:  i18n.Text,
		ErrorRaw:      raw,
		Error:         nil,
		HttpErrorCode: http.StatusInternalServerError,
		Type:          EXCEPTION_TYPE_SERVER,
		LogLevel:      level,
	})
}

func ThrowExceptionServerMessageWithError(i18n message.I18N, err error, logLevel ...zapcore.Level) {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	var level *zapcore.Level
	if len(logLevel) > 0 {
		level = &logLevel[0]
	}
	ThrowException(Exception{
		ErrorCode:     i18n.Code,
		ErrorMessage:  i18n.Text,
		ErrorRaw:      errStr,
		Error:         err,
		HttpErrorCode: http.StatusInternalServerError,
		Type:          EXCEPTION_TYPE_SERVER,
		LogLevel:      level,
	})
}

func ThrowExceptionServer404Message(i18n message.I18N, raw string) {
	infoLevel := zapcore.InfoLevel
	ThrowException(Exception{
		ErrorCode:     i18n.Code,
		ErrorMessage:  i18n.Text,
		ErrorRaw:      raw,
		Error:         nil,
		HttpErrorCode: http.StatusNotFound,
		Type:          EXCEPTION_TYPE_SERVER,
		LogLevel:      &infoLevel,
	})
}

func ThrowExceptionServer404(i18n message.I18N) {
	ThrowExceptionServer404Message(i18n, "")
}

func ThrowExceptionClient404Message(i18n message.I18N, raw string) {
	ThrowException(Exception{
		ErrorCode:     i18n.Code,
		ErrorMessage:  i18n.Text,
		ErrorRaw:      raw,
		Error:         nil,
		HttpErrorCode: http.StatusNotFound,
		Type:          EXCEPTION_TYPE_CLIENT,
	})
}

func ThrowExceptionClient400Message(i18n message.I18N, raw string, logLevel ...zapcore.Level) {
	var level *zapcore.Level
	if len(logLevel) > 0 {
		level = &logLevel[0]
	}
	ThrowException(Exception{
		ErrorCode:     i18n.Code,
		ErrorMessage:  i18n.Text,
		ErrorRaw:      raw,
		Error:         nil,
		HttpErrorCode: http.StatusBadRequest,
		Type:          EXCEPTION_TYPE_CLIENT,
		LogLevel:      level,
	})
}

func ThrowExceptionClient404Message3(i18n message.I18N) {
	ThrowExceptionClient404Message(i18n, "")
}

func ThrowExceptionClientWithHttpCode(code string, message string,
	raw string, httpCode int) {
	ThrowException(Exception{
		ErrorCode:     code,
		ErrorMessage:  message,
		ErrorRaw:      raw,
		Type:          EXCEPTION_TYPE_CLIENT,
		HttpErrorCode: httpCode,
	})
}

func CatchPanic(executeOnPanic func(any)) {
	if err := recover(); err != nil {
		executeOnPanic(err)
	}
}
func HandleErrorServerMessage(err error, i18n message.I18N, logLevel ...zapcore.Level) {
	if err == nil {
		return
	}
	ThrowExceptionServerMessageWithError(i18n, err, logLevel...)
}

func HandleErrorClientMessage(err error, i18n message.I18N, logLevel ...zapcore.Level) {
	if err == nil {
		return
	}
	ThrowExceptionClientMessage2(i18n, err, logLevel...)
}

func TryCatchAndLog(requestSession *logy.RequestSession, tryFunc func()) {
	TryCatch(tryFunc, func(exception Exception) {
		LogException(requestSession, exception)
	})
}

func TryCatch(tryFunc func(), catchFunc func(exception Exception)) {
	defer CatchPanic(func(err any) {
		if err != nil {
			var exp Exception
			switch err.(type) {
			case Exception:
				exp = err.(Exception)
			case error:
				exp = Exception{
					ErrorCode: message.ErrorUnexpectError,
					Type:      EXCEPTION_TYPE_SERVER,
					LogLevel:  reflection.ToPointer(zapcore.ErrorLevel),
					ErrorRaw:  err.(error).Error(),
					Error:     err.(error),
				}
			default:
				exp = Exception{
					ErrorCode: message.ErrorUnexpectPanic,
					ErrorRaw:  fmt.Sprintf("%#v", exp),
					Type:      EXCEPTION_TYPE_SERVER,
					LogLevel:  reflection.ToPointer(zapcore.ErrorLevel),
				}
			}

			catchFunc(exp)
		}
	})
	tryFunc()
}
func TryCatchAndThrow(tryFunc func(), catchFunc func(exception Exception) Exception) {
	defer CatchPanic(func(err any) {
		if err != nil {
			switch err.(type) {
			case Exception:
				newException := catchFunc(err.(Exception))
				ThrowException(newException)
				return
			default:
				panic(err)
			}
		}
	})
	tryFunc()
}

func RunAsyncAndLogException(requestSession *logy.RequestSession, tryFunc func()) {
	go func() {
		TryCatch(tryFunc, func(exception Exception) {
			LogException(requestSession, exception)
		})
	}()
}
func RunAsyncAndLogExceptionAndInform(requestSession *logy.RequestSession, tryFunc func(), catchFunc func(exception Exception) Exception) {
	go func() {
		TryCatch(tryFunc, func(exception Exception) {
			LogException(requestSession, exception)
			catchFunc(exception)
		})
	}()
}

func LogExceptionAndThrow(requestSession *logy.RequestSession, tryFrunc func()) {
	TryCatch(func() {
		tryFrunc()
	}, func(exception2 Exception) {
		LogException(requestSession, exception2)
		ThrowException(exception2)
	})
}
