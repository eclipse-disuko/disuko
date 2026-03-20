// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package analyticslicenses

import (
	"github.com/eclipse-disuko/disuko/domain/analytics"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
)

const LicensesCollectionName = "analyticslicenses"

type ILicensesRepository interface {
	base.IBaseRepositoryWithSoftDelete[*analytics.License]
	SearchLicenceByName(requestSession *logy.RequestSession, name string, exact bool) []string
	InitIndex(session *logy.RequestSession)
	FindByName(requestSession *logy.RequestSession, name string) []*analytics.License
	AddToIndex(requestSession *logy.RequestSession, name string)
}
