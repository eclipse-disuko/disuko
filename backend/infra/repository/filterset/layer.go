// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package filtersets

import (
	"github.com/eclipse-disuko/disuko/domain/filterset"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
)

const FilterSetsCollectionName = "filtersets"

type IFilterSetsRepository interface {
	base.IBaseRepositoryWithHardDelete[*filterset.FilterSetEntity]

	FindByTableName(requestSession *logy.RequestSession, tableName string) []*filterset.FilterSetEntity
}
