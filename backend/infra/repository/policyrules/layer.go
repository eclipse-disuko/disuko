// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package policyrules

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/license"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

const PolicyRulesCollectionName = "rules"

type IPolicyRulesRepository interface {
	base.IBaseRepositoryWithSoftDelete[*license.PolicyRules]
	FindByName(requestSession *logy.RequestSession, name string) *license.PolicyRules
	FindPolicyRulesForLabel(requestSession *logy.RequestSession, label []string) []*license.PolicyRules
	ExistsByLabel(requestSession *logy.RequestSession, label string) bool
}
