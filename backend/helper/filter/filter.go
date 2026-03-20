// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package filter

import (
	"encoding/json"
	"strings"

	"github.com/eclipse-disuko/disuko/domain/search"
	"github.com/eclipse-disuko/disuko/helper"
)

func MatchesCriteria[T any](item T, searchOptions search.SortableOptions, extractors map[string]func(T) string, arrayExtractors map[string]func(T) []string) bool {
	if searchOptions.HasColumnFilters() && (extractors != nil || arrayExtractors != nil) {
		for col, values := range searchOptions.GetFilterBy() {
			if len(values) > 0 {
				if extractor, found := extractors[col]; found {
					extractedValue := extractor(item)
					if !helper.Contains(extractedValue, values) {
						return false
					}
				} else if arrayExtractor, found := arrayExtractors[col]; found {
					extractedValues := arrayExtractor(item)
					if len(extractedValues) == 0 && helper.Contains("", values) {
						continue
					}
					matches := false
					for _, val := range extractedValues {
						if helper.Contains(val, values) {
							matches = true
							break
						}
					}
					if !matches {
						return false
					}
				}
			}
		}
	}
	if searchOptions.HasFilter() {
		out, err := json.Marshal(item)
		if err != nil {
			return false
		}
		jsonString := strings.ToLower(string(out))
		return strings.Contains(jsonString, strings.ToLower(searchOptions.GetFilterString()))
	}
	return true
}
