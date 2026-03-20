// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package analytics

import (
	da "mercedes-benz.ghe.com/foss/disuko/domain/analytics"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

const AnalyticsCollectionName = "analytics"

type IAnalyticsRepository interface {
	base.IBaseRepositoryWithSoftDelete[*da.Analytics]
	FindByNameAndProjectKeysAndLicense(requestSession *logy.RequestSession, componentName string, keys []string, licenceEffective string, offset, limit int, sortCol string, asc bool) []*da.Analytics
}
