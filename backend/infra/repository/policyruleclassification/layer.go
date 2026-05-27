// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package policyruleclassification

import (
	"github.com/eclipse-disuko/disuko/domain/policyruleclassification"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
)

const CollectionName = "policy_rule_classifications"

type IPolicyRuleClassificationRepository interface {
	base.IBaseRepositoryWithHardDelete[*policyruleclassification.PolicyRuleClassification]
}
