// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package helper

import "strings"

func IsUnasserted(input string) bool {
	return strings.ToLower(input) == "none" || strings.ToLower(input) == "noassertion" || input == ""
}
