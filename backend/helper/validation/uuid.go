// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func IsUUID(str string) bool {
	if len(str) > 0 {
		var validUUID = regexp.MustCompile(`[a-z]+[0-9]+[-]`)
		return validUUID.MatchString(str)
	}
	return false
}

func CheckUuid(uuid string) error {
	validate := validator.New()
	errs := validate.Var(uuid, "uuid,len=36")

	if errs != nil {
		return errs
	}

	return nil
}
