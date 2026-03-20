// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package helper

import (
	"sort"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func EqualsStringSlices(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, v := range s1 {
		if v != s2[i] {
			return false
		}
	}
	return true
}

func EqualsStringSlicesIgnoreOrder(s1, s2 []string) bool {
	diff := cmp.Diff(s1, s2, cmpopts.SortSlices(
		func(a, b string) bool {
			return a < b
		},
	)) == ""
	return diff
}

func RemoveDuplicates(original [][]string) [][]string {
	// Sort content of each inner slice before remove duplicates
	for _, sliceCurrent := range original {
		sort.Sort(sort.StringSlice(sliceCurrent))
	}
	reduced := make([][]string, 0)
	for i, sliceIth := range original {
		duplicatePresent := false
		for j, sliceJth := range original {
			if i < j && EqualsStringSlices(sliceIth, sliceJth) {
				duplicatePresent = true
				break
			}
		}
		if !duplicatePresent {
			reduced = append(reduced, sliceIth)
		}
	}
	return reduced
}
