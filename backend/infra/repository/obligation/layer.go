// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package obligation

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/obligation"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

const ObligationCollectionName = "obligations"

type IObligationRepository interface {
	base.IBaseRepositoryWithSoftDelete[*obligation.Obligation]
	FindAllSortedByName(requestSession *logy.RequestSession) []*obligation.Obligation
	FindByName(requestSession *logy.RequestSession, key string) []*obligation.Obligation
	// FindByKeys(requestSession *logy.RequestSession, keys []string) []*obligation.Obligation
}
