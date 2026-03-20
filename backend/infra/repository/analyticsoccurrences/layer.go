// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package analyticsoccurrences

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/analytics"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
)

const OccurrencesCollectionName = "analyticsoccurrences"

type IOccurrencesRepository interface {
	base.IBaseRepositoryWithSoftDelete[*analytics.Occurrence]
}
