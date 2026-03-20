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

const SbomListCollectionName = "sbomlist"

type ISbomListRepository interface {
	base.IBaseRepositoryWithSoftDelete[*sbomlist.SbomList]
	FindFile(rs *logy.RequestSession, versionKey, spdxKey string) *project.SpdxFileBase
}
