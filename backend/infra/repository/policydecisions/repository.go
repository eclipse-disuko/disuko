// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package policydecisions

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/policydecisions"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type policyDecisionsRepositoryStruct struct {
	base.BaseRepositoryWithSoftDelete[*policydecisions.PolicyDecisions]
}

func NewPolicyDecisionsRepository(requestSession *logy.RequestSession) IPolicyDecisionsRepository {
	return &policyDecisionsRepositoryStruct{
		BaseRepositoryWithSoftDelete: base.CreateRepositoryWithSoftDelete[*policydecisions.PolicyDecisions](
			requestSession,
			PolicyDecisionsCollectionName,
			func() *policydecisions.PolicyDecisions {
				return &policydecisions.PolicyDecisions{}
			},
			nil,
			nil,
			nil),
	}
}
