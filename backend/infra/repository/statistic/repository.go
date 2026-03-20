// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package statistic

import (
	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/statistic"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
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
