// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldValidateFilename(t *testing.T) {
	//given
	type test struct {
		file  string
		valid bool
	}
	tests := []test{
		{file: "test.zip", valid: true},
		{file: "test.gz", valid: true},
		{file: "test@#$%^&()_.gz", valid: true},
		{file: "a/a.zip", valid: false},
		{file: "a<.zip", valid: false},
		{file: "a>.zip", valid: false},
		{file: "a'.zip", valid: false},
		{file: "a\".zip", valid: false},
		{file: "a\\.zip", valid: false},
		{file: "a*.zip", valid: false},
		{file: "a?.zip", valid: false},
		{file: "a|.zip", valid: false},
	}

	//when
	for _, tc := range tests {
		got := isFilenameValid(tc.file)
		//then
		assert.Equal(t, tc.valid, got, "misfunctioning validation of filename for "+tc.file)
	}

}
