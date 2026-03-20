// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package analyticsoccurrences

import (
	"github.com/eclipse-disuko/disuko/domain/analytics"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
)

const OccurrencesCollectionName = "analyticsoccurrences"

type IOccurrencesRepository interface {
	base.IBaseRepositoryWithSoftDelete[*analytics.Occurrence]
}
