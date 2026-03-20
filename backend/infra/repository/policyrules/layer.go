// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package policyrules

import (
	"github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
)

const PolicyRulesCollectionName = "rules"

type IPolicyRulesRepository interface {
	base.IBaseRepositoryWithSoftDelete[*license.PolicyRules]
	FindByName(requestSession *logy.RequestSession, name string) *license.PolicyRules
	FindPolicyRulesForLabel(requestSession *logy.RequestSession, label []string) []*license.PolicyRules
	ExistsByLabel(requestSession *logy.RequestSession, label string) bool
}
