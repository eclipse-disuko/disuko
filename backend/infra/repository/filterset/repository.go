// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package filtersets

import (
	"github.com/eclipse-disuko/disuko/domain/filterset"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
)

type filterSetsRepositoryStruct struct {
	base.BaseRepositoryWithHardDelete[*filterset.FilterSetEntity]
}

func (f *filterSetsRepositoryStruct) FindByTableName(requestSession *logy.RequestSession, tableName string) []*filterset.FilterSetEntity {
	var result []*filterset.FilterSetEntity
	for _, filterSet := range f.FindAll(requestSession, false) {
		if filterSet.TableName == tableName {
			result = append(result, filterSet)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func NewFilterSetsRepository(requestSession *logy.RequestSession) IFilterSetsRepository {
	return &filterSetsRepositoryStruct{
		BaseRepositoryWithHardDelete: base.CreateRepositoryWithHardDelete[*filterset.FilterSetEntity](
			requestSession,
			FilterSetsCollectionName,
			func() *filterset.FilterSetEntity {
				return &filterset.FilterSetEntity{}
			},
			nil,
			"",
			nil,
			nil),
	}
}
