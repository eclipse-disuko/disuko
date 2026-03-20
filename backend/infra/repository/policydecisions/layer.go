// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package policydecisions

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/policydecisions"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
)

const PolicyDecisionsCollectionName = "policydecisions"

type IPolicyDecisionsRepository interface {
	base.IBaseRepositoryWithSoftDelete[*policydecisions.PolicyDecisions]
}
