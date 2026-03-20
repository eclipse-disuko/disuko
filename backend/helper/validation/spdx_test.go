// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSpdxIdentifier(t *testing.T) {
	assert.True(t, IsSpdxIdentifier("abasda"))
	assert.True(t, IsSpdxIdentifier("ab-asda"))
	assert.True(t, IsSpdxIdentifier("ab_asda"))
	assert.True(t, IsSpdxIdentifier("ab_asda_"))
	assert.True(t, IsSpdxIdentifier("ab-asdaAsd"))
	assert.True(t, IsSpdxIdentifier("ab-asdaAsd+"))
	assert.True(t, IsSpdxIdentifier("ab-asdaAsd-999.098"))
	assert.True(t, IsSpdxIdentifier("LicenseRef-MB-Flavor-Artistic"))

	assert.False(t, IsSpdxIdentifier("ab asda"))
	assert.False(t, IsSpdxIdentifier("ab "))
	assert.False(t, IsSpdxIdentifier(""))
	assert.False(t, IsSpdxIdentifier("ä"))
	assert.False(t, IsSpdxIdentifier("ß"))
	assert.False(t, IsSpdxIdentifier("LicenseRef-MB-Flavor Artistic"))
}
