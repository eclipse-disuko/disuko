// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"fmt"

	auditHelper "mercedes-benz.ghe.com/foss/disuko/helper/audit"
)

func (entity SourceExternal) ToAudit() string {
	return fmt.Sprintf("URL: %s (%s), Origin: %s, Uploader: %s, Updated: %s, Created: %s, Key %s", entity.URL, entity.Comment, entity.Origin, entity.Uploader, auditHelper.ConvertDateTime(entity.Updated), auditHelper.ConvertDateTime(entity.Created), entity.Key)
}
