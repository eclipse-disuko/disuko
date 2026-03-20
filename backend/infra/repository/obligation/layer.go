// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package obligation

import (
	"github.com/eclipse-disuko/disuko/domain/obligation"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
)

const ObligationCollectionName = "obligations"

type IObligationRepository interface {
	base.IBaseRepositoryWithSoftDelete[*obligation.Obligation]
	FindAllSortedByName(requestSession *logy.RequestSession) []*obligation.Obligation
	FindByName(requestSession *logy.RequestSession, key string) []*obligation.Obligation
	// FindByKeys(requestSession *logy.RequestSession, keys []string) []*obligation.Obligation
}
