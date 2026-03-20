// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_checkStringIsProjectNameOrUUID(t *testing.T) {
	projectName := "Disclosure Portal Backend"
	shouldBeFalse := IsUUID(projectName)
	fmt.Println(shouldBeFalse)
	assert.Equal(t, shouldBeFalse, false)

	projectUUID := "c3168451-c768-42da-b5d6-2e236e8df18d"
	shouldBeTrue := IsUUID(projectUUID)
	assert.Equal(t, shouldBeTrue, true)
}
