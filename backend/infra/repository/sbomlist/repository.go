// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package sbomlist

import (
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/project/sbomlist"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
)

type sbomListRepositoryStruct struct {
	base.BaseRepositoryWithSoftDelete[*sbomlist.SbomList]
}

func NewSbomListRepository(requestSession *logy.RequestSession) ISbomListRepository {
	repo := &sbomListRepositoryStruct{
		BaseRepositoryWithSoftDelete: base.CreateRepositoryWithSoftDelete[*sbomlist.SbomList](
			requestSession,
			SbomListCollectionName,
			func() *sbomlist.SbomList {
				return &sbomlist.SbomList{}
			},
			nil,
			nil,
			nil),
	}

	return repo
}

func (r *sbomListRepositoryStruct) FindFile(rs *logy.RequestSession, versionKey, spdxKey string) *project.SpdxFileBase {
	list := r.FindByKey(rs, versionKey, false)
	if list == nil {
		return nil
	}
	for _, spdx := range list.SpdxFileHistory {
		if spdx.Key == spdxKey {
			return spdx
		}
	}
	return nil
}
