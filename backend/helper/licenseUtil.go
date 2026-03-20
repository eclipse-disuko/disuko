// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package helper

import "strings"

func RemoveAllBracketsInLicenseText(license string) string {
	license = strings.ReplaceAll(strings.ReplaceAll(license, ")", ""), "(", "")
	return license
}

func SplitLicenses(license string) []string {
	if IsOrDualLicensed(license) {
		return strings.Split(strings.ReplaceAll(strings.ReplaceAll(license, ")", ""), "(", ""), " OR ")
	} else if isAndDualLicensed(license) {
		return strings.Split(strings.ReplaceAll(strings.ReplaceAll(license, ")", ""), "(", ""), " AND ")
	} else {
		return []string{license}
	}
}

func IsDualLicensed(license string) bool {
	return isAndDualLicensed(license) || IsOrDualLicensed(license)
}

func isAndDualLicensed(license string) bool {
	return strings.Contains(license, " AND ")
}

func IsOrDualLicensed(license string) bool {
	return strings.Contains(license, " OR ")
}
