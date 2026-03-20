// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package integrity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_append_inside_function(t *testing.T) {
	sliceTest := make([]string, 0)
	appendInFunc("test", &sliceTest)

	assert.Equal(t, 1, len(sliceTest))
	assert.Equal(t, "test", sliceTest[0])
}

func appendInFunc(value string, slices *[]string) {
	*slices = append(*slices, value)
}

func Test_modal(t *testing.T) {
	t.Logf("100 mod 10 = %d", 100%10)
	t.Logf("101 mod 10 = %d", 101%10)
}
