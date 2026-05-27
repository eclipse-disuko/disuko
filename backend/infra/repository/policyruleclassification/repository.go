// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package policyruleclassification

import (
	"github.com/eclipse-disuko/disuko/domain/policyruleclassification"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
)

type policyRuleClassificationRepositoryStruct struct {
	base.BaseRepositoryWithHardDelete[*policyruleclassification.PolicyRuleClassification]
}

func NewPolicyRuleClassificationRepository(requestSession *logy.RequestSession) IPolicyRuleClassificationRepository {
	return &policyRuleClassificationRepositoryStruct{
		BaseRepositoryWithHardDelete: base.CreateRepositoryWithHardDelete[*policyruleclassification.PolicyRuleClassification](
			requestSession,
			CollectionName,
			func() *policyruleclassification.PolicyRuleClassification {
				return &policyruleclassification.PolicyRuleClassification{}
			},
			nil,
			"",
			nil,
			nil,
		),
	}
}
