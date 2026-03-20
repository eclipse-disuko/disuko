// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package helper

import (
	"strings"
)

func Search(arr []string, search string, exact bool) []string {
	filtered := make([]string, 0)
	for _, str := range arr {
		if (!exact && strings.Contains(str, search)) || (exact && str == search) {
			filtered = append(filtered, str)
		}
	}
	return filtered
}
