// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"net/textproto"
	"testing"

	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/test"
)

func TestCheckExpectedContentType2_OK_1(t *testing.T) {
	header := textproto.MIMEHeader{}
	contentTypeExpected := "TestHeader"
	header.Add("Content-Type", contentTypeExpected)
	CheckExpectedContentType2(header, []ContentType{
		ContentType(contentTypeExpected),
	})
}

func TestCheckExpectedContentType2_OK_2(t *testing.T) {
	header := textproto.MIMEHeader{}
	contentTypeExpected := "TestHeader"
	header.Add("Content-Type", "adsa;"+contentTypeExpected+";asdhsadg;asdsad")
	CheckExpectedContentType2(header, []ContentType{
		ContentType(contentTypeExpected),
	})
}
func TestCheckExpectedContentType2_OK_3(t *testing.T) {
	header := textproto.MIMEHeader{}
	header.Add("Content-Type", "adsa;"+string(ContentTypeJson)+";asdhsadg;asdsad")
	CheckExpectedContentType2(header, []ContentType{
		ContentTypeJson,
		ContentTypeOctets,
	})
}
func TestCheckExpectedContentType2_FAILED(t *testing.T) {
	header := textproto.MIMEHeader{}
	contentTypeExpected := "TestHeader"
	header.Add("Content-Type", string(ContentTypeFormData))
	test.ExpectException(t, message.ErrorContentTypeWrong, func() {
		CheckExpectedContentType2(header, []ContentType{
			ContentType(contentTypeExpected),
		})
	})
}
