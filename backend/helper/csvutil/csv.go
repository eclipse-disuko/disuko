// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package csvutil

import "strings"

func PrepareFieldForCsv(value string) string {
	value = strings.TrimSpace(strings.ReplaceAll(value, "\n", " "))

	if strings.Contains(value, ",") {
		value = strings.ReplaceAll(value, "\"", "\"\"")
		value = "\"" + value + "\""
		return value
	}

	// For non-comma fields, just protect Excel formulas
	if len(value) > 0 && strings.ContainsAny(string(value[0]), "=+-@") {
		value = "'" + value
	}

	return value
}
