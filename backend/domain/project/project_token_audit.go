// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"fmt"
	"time"

	auditHelper "mercedes-benz.ghe.com/foss/disuko/helper/audit"
)

func (token *Token) ToAudit() string {
	if token == nil {
		return ""
	}
	expiryParsedFormatted := token.Expiry
	expiryParsed, err := time.Parse("2006-01-02T15:04:05.000Z", token.Expiry)
	if err == nil {
		expiryParsedFormatted = auditHelper.ConvertDateTime(expiryParsed)
	}

	if len(token.Description) > 0 {
		return fmt.Sprintf("%s (%s), Expiry: %s, Created: %s, Status: %s, Key: %s", token.Company, token.Description, expiryParsedFormatted, auditHelper.ConvertDateTime(token.Created), token.Status, token.Key)
	}
	return fmt.Sprintf("%s, Expiry: %s, Created: %s, Status: %s, Key: %s", token.Company, expiryParsedFormatted, auditHelper.ConvertDateTime(token.Created), token.Status, token.Key)
}
