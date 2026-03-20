// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"net/http"
	"net/textproto"
	"strings"

	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
)

type ContentType string

const (
	ContentTypeJson     ContentType = "application/json"
	ContentTypeFormData ContentType = "multipart/form-data"
	ContentTypeOctets   ContentType = "application/octet-stream"
)

func CheckExpectedContentType(request *http.Request, expected ContentType) {
	CheckExpectedContentType2(textproto.MIMEHeader(request.Header), []ContentType{
		expected,
	})
}

func CheckExpectedContentType2(header textproto.MIMEHeader, expectedTypes []ContentType) {
	ct := strings.ToLower(strings.TrimSpace(header.Get("Content-Type")))
	if ct == "" {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorContentTypeWrong), "")
		return
	}
	//it could also contains additional metadata like "charset=UTF-8" or "boundary=boundary"
	ctSplit := strings.Split(ct, ";")
	found := false

OUTER:
	for _, ctPart := range ctSplit {
		for _, expected := range expectedTypes {
			if strings.Contains(strings.ToLower(string(expected)), ctPart) {
				found = true
				break OUTER
			}
		}
	}
	if !found {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorContentTypeWrong), "")
	}
}
