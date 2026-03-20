// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package helper

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var UrlRegex = regexp.MustCompile(`(http:\/\/|file:\/\/|https:\/\/|ssh:\/\/)+[^\\]+`)

func Contains(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func StringContainsI(a string, b string) bool {
	return strings.Contains(
		strings.ToLower(a),
		strings.ToLower(b),
	)
}

func RemoveStrFromSlice(s string, slice []string) []string {
	result := make([]string, 0, len(slice))
	for _, v := range slice {
		if v != s {
			result = append(result, v)
		}
	}
	return result
}

func ByteToMB(byteCount int64) float64 {
	return float64(byteCount) / 1024 / 1024
}

func CreateOrigin(origin string, company string, description string, key string) string {
	if len(description) > 0 {
		origin = fmt.Sprintf("%s ('%s', identifier: %s) by '%s'", origin, description, maskUuid(key), company)
	} else {
		origin = fmt.Sprintf("%s (identifier: %s) by '%s'", origin, maskUuid(key), company)
	}
	return origin
}

func maskUuid(uuid string) string {
	asterisks := strings.Repeat("*", len(uuid)-10)
	return uuid[:10] + asterisks
}

func GetPointerOfTimeNow() *time.Time {
	//is ugly, but there is no other way
	now := time.Now()
	return &now
}
