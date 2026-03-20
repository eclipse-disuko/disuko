// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package filtersets

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/filterset"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

const FilterSetsCollectionName = "filtersets"

type IFilterSetsRepository interface {
	base.IBaseRepositoryWithHardDelete[*filterset.FilterSetEntity]

	FindByTableName(requestSession *logy.RequestSession, tableName string) []*filterset.FilterSetEntity
}
