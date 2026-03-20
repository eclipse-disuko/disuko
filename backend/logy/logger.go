// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package logy

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey string

const rsKey ctxKey = "rs"

var logger *zap.SugaredLogger

const (
	MsgStage                   = "msgStage"
	ReqID                      = "reqID"
	MsgStageStartRequest       = ">>> Start"
	MsgStageIntermediateCommon = "---"
	MsgStageFinishRequest      = "<<< Finish"
	InternalReqID              = "SYSTEM"
	Identifier                 = "IDENTIFIER"
	Value                      = "VALUE"
)

type RequestSession struct {
	ReqID      string
	QueryCount int
	QueryTime  time.Duration
}

func NewRequestSession() *RequestSession {
	return &RequestSession{
		ReqID: uuid.NewString(),
	}
}

func RequestWithSession(r *http.Request, rs *RequestSession) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), rsKey, rs))
}

func GetRequestSession(r *http.Request) *RequestSession {
	raw := r.Context().Value(rsKey)
	rs, ok := raw.(*RequestSession)
	if !ok {
		n := &RequestSession{
			ReqID: uuid.NewString(),
		}
		n.Warnf("No request session found in request context")
		return n
	}
	return rs
}

func init() {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel && (lvl > zapcore.DebugLevel || conf.Config.Server.DevLog)
	})

	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	productionEncoderConfig := zap.NewProductionEncoderConfig()
	productionEncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05,000")
	prodEncoder := zapcore.NewJSONEncoder(productionEncoderConfig)

	var cores []zapcore.Core
	if conf.Config.Server.DevLog {
		cores = append(cores, zapcore.NewCore(consoleEncoder, consoleErrors, highPriority))
		cores = append(cores, zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority))
	} else {
		cores = append(cores, zapcore.NewCore(prodEncoder, consoleErrors, highPriority))
		cores = append(cores, zapcore.NewCore(prodEncoder, consoleDebugging, lowPriority))
	}
	core := zapcore.NewTee(
		cores...)

	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel), zap.AddCallerSkip(1)).Sugar()
	defer logger.Sync()
}

func (rs *RequestSession) Infof(format string, keysAndValues ...any) {
	if rs == nil {
		logger.Infow(format, keysAndValues...)
		return
	}
	logger.Infow(fmt.Sprintf("[%s] %s", rs.ReqID, format), append(keysAndValues, ReqID, rs.ReqID)...)
}

func (rs *RequestSession) Warnf(format string, keysAndValues ...any) {
	if rs == nil {
		logger.Warnf(format, keysAndValues...)
		return
	}
	logger.Warnw(fmt.Sprintf("[%s] %s", rs.ReqID, format), append(keysAndValues, ReqID, rs.ReqID)...)
}

func Infow(requestSession *RequestSession, msg string, keysAndValues ...interface{}) {
	if requestSession == nil {
		logger.Infow(msg, keysAndValues...)
	} else {
		keysAndValues = append(keysAndValues, ReqID)
		keysAndValues = append(keysAndValues, requestSession.ReqID)
		logger.Infow(fmt.Sprintf("[%s] %s", requestSession.ReqID, msg), keysAndValues...)
	}
}

func Warnw(requestSession *RequestSession, msg string, keysAndValues ...interface{}) {
	if requestSession == nil {
		logger.Warnw(msg, keysAndValues...)
	} else {
		keysAndValues = append(keysAndValues, ReqID)
		keysAndValues = append(keysAndValues, requestSession.ReqID)
		logger.Warnw(fmt.Sprintf("[%s] %s", requestSession.ReqID, msg), keysAndValues...)
	}
}

func Logw(requestSession *RequestSession, msg string, logLevel zapcore.Level, keysAndValues ...interface{}) {
	switch logLevel {
	case zapcore.WarnLevel:
		Warnw(requestSession, msg, keysAndValues...)
	case zapcore.InfoLevel:
		Infow(requestSession, msg, keysAndValues...)
	default:
		Errorw(requestSession, msg, keysAndValues...)
	}
}

func Infof(requestSession *RequestSession, template string, args ...interface{}) {
	if requestSession == nil {
		logger.Infof(template, args...)
	} else {
		logger.Infof("["+requestSession.ReqID+"] "+template, args...)
	}
}

func Errorw(requestSession *RequestSession, msg string, keysAndValues ...interface{}) {
	if requestSession == nil {
		logger.Errorw(msg, keysAndValues...)
	} else {
		keysAndValues = append(keysAndValues, ReqID)
		keysAndValues = append(keysAndValues, requestSession.ReqID)
		logger.Errorw(fmt.Sprintf("[%s] %s", requestSession.ReqID, msg), keysAndValues...)
	}
}

func Errorf(requestSession *RequestSession, template string, args ...interface{}) {
	if requestSession == nil {
		logger.Errorf(template, args...)
	} else {
		logger.Errorf("["+requestSession.ReqID+"] "+template, args...)
	}
}

func Warnf(requestSession *RequestSession, template string, args ...interface{}) {
	if requestSession == nil {
		logger.Warnf(template, args...)
	} else {
		logger.Warnf("["+requestSession.ReqID+"] "+template, args...)
	}
}

func Fatalw(requestSession *RequestSession, msg string, keysAndValues ...interface{}) {
	if requestSession == nil {
		logger.Fatalw(msg, keysAndValues...)
	} else {
		keysAndValues = append(keysAndValues, ReqID)
		keysAndValues = append(keysAndValues, requestSession.ReqID)
		logger.Fatalw(fmt.Sprintf("[%s] %s", requestSession.ReqID, msg), keysAndValues...)
	}
}

func Fatalf(requestSession *RequestSession, template string, args ...interface{}) {
	if requestSession == nil {
		logger.Fatalf(template, args...)
	} else {
		logger.Fatalf("["+requestSession.ReqID+"] "+template, args...)
	}
}

func Debugw(requestSession *RequestSession, msg string, keysAndValues ...interface{}) {
	if requestSession == nil {
		logger.Debugw(msg, keysAndValues...)
	} else {
		keysAndValues = append(keysAndValues, ReqID)
		keysAndValues = append(keysAndValues, requestSession.ReqID)
		logger.Debugw(fmt.Sprintf("[%s] %s", requestSession.ReqID, msg), keysAndValues...)
	}
}

func Debugf(requestSession *RequestSession, template string, args ...interface{}) {
	if requestSession == nil {
		logger.Debugf(template, args...)
	} else {
		logger.Debugf("["+requestSession.ReqID+"] "+template, args...)
	}
}
