// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"testing"
)

func TestConvertUserName(t *testing.T) {
	username := "EMEA\\MEBELKO"
	should := "MEBELKO"
	userNameConverted := ConvertUsernameWithoutEmea(username)
	if userNameConverted != should {
		t.Fail()
	}

}
