// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package statistic

import (
	"github.com/eclipse-disuko/disuko/domain/statistic"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
)

const StatisticCollectionName = "statistics"

type IStatisticRepository interface {
	base.IBaseRepositoryWithHardDelete[*statistic.SystemStatistic]
}
