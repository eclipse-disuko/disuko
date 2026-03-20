// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package analyticscomponents

import (
	"github.com/eclipse-disuko/disuko/domain/analytics"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
)

const ComponentsCollectionName = "analyticscomponents"

type IComponentsRepository interface {
	base.IBaseRepositoryWithSoftDelete[*analytics.Component]
	SearchByName(requestSession *logy.RequestSession, name string, exact bool) []string
	InitIndex(requestSession *logy.RequestSession)
	FindByNameAndVersion(requestSession *logy.RequestSession, name, version string) []*analytics.Component
	AddToIndex(requestSession *logy.RequestSession, name string)
}
