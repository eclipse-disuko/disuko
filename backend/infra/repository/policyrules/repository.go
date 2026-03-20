// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package policyrules

import (
	"github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/helper"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/infra/repository/database"
	"github.com/eclipse-disuko/disuko/logy"
)

type PolicyRulesRepository struct{}

type policyRulesRepositoryStruct struct {
	base.BaseRepositoryWithSoftDelete[*license.PolicyRules]
}

func NewPolicyRulesRepository(requestSession *logy.RequestSession) IPolicyRulesRepository {
	return &policyRulesRepositoryStruct{
		BaseRepositoryWithSoftDelete: base.CreateRepositoryWithSoftDelete[*license.PolicyRules](
			requestSession,
			PolicyRulesCollectionName,
			func() *license.PolicyRules {
				return &license.PolicyRules{}
			},
			nil,
			nil,
			nil),
	}
}

func (repository *policyRulesRepositoryStruct) FindByName(requestSession *logy.RequestSession, name string) *license.PolicyRules {
	qc := database.New().SetMatcher(
		database.AndChain(
			database.AttributeMatcher(
				"Name",
				database.EQ,
				name),
			database.AttributeMatcher(
				"Deleted",
				database.EQ,
				false)))
	var qPR *license.PolicyRules
	qPolicyRules := repository.Query(requestSession, qc)
	if len(qPolicyRules) > 0 {
		qPR = qPolicyRules[0]
	}
	return qPR
}

func (repository *policyRulesRepositoryStruct) FindPolicyRulesForLabel(requestSession *logy.RequestSession, labels []string) []*license.PolicyRules {
	all := repository.FindAll(requestSession, false)
	var res []*license.PolicyRules
	for _, r := range all {
		if !r.Active {
			continue
		}
		if r.ApplyToAll {
			res = append(res, r)
			continue
		}
		for _, s := range r.LabelSets {
			if !helper.EqualsStringSlicesIgnoreOrder(s, labels) {
				continue
			}
			res = append(res, r)
			break
		}
	}
	return res
}

func (repository *policyRulesRepositoryStruct) ExistsByLabel(requestSession *logy.RequestSession, label string) bool {
	all := repository.FindAll(requestSession, false)
	for _, r := range all {
		for _, s := range r.LabelSets {
			if !helper.Contains(label, s) {
				continue
			}
			return true
		}
	}
	return false
}
