// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package sbomlist

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/project"
	"mercedes-benz.ghe.com/foss/disuko/domain/project/sbomlist"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

const SbomListCollectionName = "sbomlist"

type ISbomListRepository interface {
	base.IBaseRepositoryWithSoftDelete[*sbomlist.SbomList]
	FindFile(rs *logy.RequestSession, versionKey, spdxKey string) *project.SpdxFileBase
}
