// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package analyticscomponents

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/analytics"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

const ComponentsCollectionName = "analyticscomponents"

type IComponentsRepository interface {
	base.IBaseRepositoryWithSoftDelete[*analytics.Component]
	SearchByName(requestSession *logy.RequestSession, name string, exact bool) []string
	InitIndex(requestSession *logy.RequestSession)
	FindByNameAndVersion(requestSession *logy.RequestSession, name, version string) []*analytics.Component
	AddToIndex(requestSession *logy.RequestSession, name string)
}
