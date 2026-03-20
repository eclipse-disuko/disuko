// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package statistic

import (
	"mercedes-benz.ghe.com/foss/disuko/domain"
	"mercedes-benz.ghe.com/foss/disuko/domain/statistic"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type systemStatisticRepositoryStruct struct {
	base.BaseRepositoryWithHardDelete[*statistic.SystemStatistic]
}

var createEmptyEntityFunc = func() *statistic.SystemStatistic {
	return &statistic.SystemStatistic{
		RootEntity: domain.NewRootEntity(),
	}
}

func NewSystemStatisticRepository(requestSession *logy.RequestSession) IStatisticRepository {
	return &systemStatisticRepositoryStruct{
		BaseRepositoryWithHardDelete: base.CreateRepositoryWithHardDelete[*statistic.SystemStatistic](
			requestSession,
			StatisticCollectionName,
			createEmptyEntityFunc,
			nil,
			"",
			nil,
			nil),
	}
}
