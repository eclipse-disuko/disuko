// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package validation

import "regexp"

func IsSpdxIdentifier(str string) bool {
	if len(str) > 0 {
		return regexp.MustCompile(`^[a-zA-Z0-9\-._+]*$`).MatchString(str)
	}
	return false
}

func IsSpdxAlias(str string) bool {
	if len(str) > 0 {
		return regexp.MustCompile(`^[a-zA-Z0-9\-._+ ]*$`).MatchString(str)
	}
	return false
}
