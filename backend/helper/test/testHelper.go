// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/stretchr/testify/assert"
)

func HasPrefix(t *testing.T, prefix string, textData string) {
	assert.True(t, strings.HasPrefix(textData, prefix),
		"\ntextData: "+textData+
			"\nprefix: "+prefix)
}

func ExpectException(t *testing.T, expectedErrorCode string, tryFrunc func()) {
	exception.TryCatch(func() {
		tryFrunc()
		assert.FailNow(t, "Missing expected error code: "+expectedErrorCode)
	}, func(exception exception.Exception) {
		assert.Equal(t, expectedErrorCode, exception.ErrorCode)
	})
}

func ConvertToByteBuffer(v interface{}) *bytes.Buffer {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	err := enc.Encode(v)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorJsonEncode))
	return buf
}
