// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package license

import (
	"github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
)

const LicensesCollectionName = "licenses"

type ILicensesRepository interface {
	base.IBaseRepositoryWithSoftDelete[*license.License]
	FindByName(requestSession *logy.RequestSession, name string) *license.License
	FindById(requestSession *logy.RequestSession, id string) *license.License
	FindByIdCaseInsensitive(requestSession *logy.RequestSession, id string) *license.License
	FindByObligationKey(requestSession *logy.RequestSession, key string) []*license.License
	CountByObligationKey(requestSession *logy.RequestSession, key string) int
	GetLicenseRefs(requestSession *logy.RequestSession) license.LicenseRefs
	FindByIds(requestSession *logy.RequestSession, ids []string) []*license.License
}
