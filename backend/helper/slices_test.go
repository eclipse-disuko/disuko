// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareSlices_DifferentLenght(t *testing.T) {
	slice1 := []string{"a", "b", "c"}
	slice2 := []string{"a", "b"}
	assert.False(t, EqualsStringSlices(slice1, slice2))
}

func TestCompareSlices_DifferentOrder(t *testing.T) {
	slice1 := []string{"a", "b", "c"}
	slice2 := []string{"b", "a", "c"}
	assert.False(t, EqualsStringSlices(slice1, slice2))
}

func TestCompareSlices_DifferentContent(t *testing.T) {
	slice1 := []string{"a", "b", "c"}
	slice2 := []string{"A", "B", "C"}
	assert.False(t, EqualsStringSlices(slice1, slice2))
}

func TestCompareSlices_Equals(t *testing.T) {
	slice1 := []string{"a", "b", "c"}
	slice2 := []string{"a", "b", "c"}
	assert.True(t, EqualsStringSlices(slice1, slice2))
}

func TestRemoveDuplicates1(t *testing.T) {
	assert.Equal(t, len(RemoveDuplicates(create2dStringSliceWithDuplicates())), 2)
}

func TestRemoveDuplicates2(t *testing.T) {
	assert.Equal(t, len(RemoveDuplicates(create2dStringSliceWithDuplicates2())), 1)
}

func TestRemoveDuplicates3(t *testing.T) {
	assert.Equal(t, len(RemoveDuplicates(create2dStringSliceWithoutDuplicates())), 5)
}

func TestRemoveDuplicates4(t *testing.T) {
	assert.Equal(t, len(RemoveDuplicates(create2dStringSliceEmpty())), 0)
}

func create2dStringSliceWithDuplicates() [][]string {
	slice1 := []string{"c", "b", "a"}
	slice2 := []string{"a", "b", "c"}
	slice3 := []string{"b", "c", "a"}
	slice4 := []string{"c", "d", "b"}
	slice5 := []string{"c", "a", "b"}
	return [][]string{slice1, slice2, slice3, slice4, slice5}
}

func create2dStringSliceWithDuplicates2() [][]string {
	slice1 := []string{"c", "b", "a"}
	slice2 := []string{"a", "b", "c"}
	slice3 := []string{"b", "c", "a"}
	slice4 := []string{"c", "b", "a"}
	slice5 := []string{"a", "c", "b"}
	return [][]string{slice1, slice2, slice3, slice4, slice5}
}

func create2dStringSliceWithoutDuplicates() [][]string {
	slice1 := []string{"a", "b", "c"}
	slice2 := []string{"b", "c", "d"}
	slice3 := []string{"c", "d", "e"}
	slice4 := []string{"d", "e", "f"}
	slice5 := []string{"e", "f", "g"}
	return [][]string{slice1, slice2, slice3, slice4, slice5}
}

func create2dStringSliceEmpty() [][]string {
	return [][]string{}
}
