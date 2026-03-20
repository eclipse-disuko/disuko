// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func newValidator() Validator {
	v := validator.New()
	v.RegisterValidation("searchtext", validateSearchText)
	v.RegisterValidation("emea", validateEmea)
	return Validator{validator: v}
}

var searchText = regexp.MustCompile("^[0-9A-Za-z-_,. äöüÄÖÜß\r\n]{3,30}$")
var emea = regexp.MustCompile("^[0-9A-Za-z_]{5,10}$")

func validateSearchText(f validator.FieldLevel) bool {
	return validateRegExp(f, searchText)
}

func validateEmea(f validator.FieldLevel) bool {
	return validateRegExp(f, emea)
}

func validateRegExp(fl validator.FieldLevel, r *regexp.Regexp) bool {
	value := fl.Field().String()
	if !r.MatchString(value) {
		return false
	}
	return true
}
