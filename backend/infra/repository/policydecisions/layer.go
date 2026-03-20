// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package policydecisions

import (
	"github.com/eclipse-disuko/disuko/domain/policydecisions"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
)

const PolicyDecisionsCollectionName = "policydecisions"

type IPolicyDecisionsRepository interface {
	base.IBaseRepositoryWithSoftDelete[*policydecisions.PolicyDecisions]
}
