// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"fmt"
	"net/http"

	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap/zapcore"
)

var Val = newValidator()

type Validator struct {
	validator *validator.Validate
}

func (v Validator) UrlParamUuid(r *http.Request, fieldName string) string {
	return v.UrlParam(r, fieldName, "uuid")
}

func (v Validator) UrlParamEMEA(r *http.Request, fieldName string) string {
	return v.UrlParam(r, fieldName, "emea")
}

func (v Validator) UrlParamSearchText(r *http.Request, fieldName string) string {
	return v.UrlParam(r, fieldName, "searchtext")
}

func (v Validator) UrlParam(r *http.Request, fieldName, tag string) string {
	value := chi.URLParam(r, fieldName)
	v.assertParam(value, fieldName, tag)
	return value
}

func (v Validator) assertParam(field interface{}, fieldName, tag string) {
	if err := v.validator.Var(field, tag); err != nil {
		exception.ThrowExceptionClient400Message(message.GetI18N(message.ErrorKeyRequestParamNotValid), fmt.Sprintf("%s invalid", fieldName), zapcore.InfoLevel)
	}
}
