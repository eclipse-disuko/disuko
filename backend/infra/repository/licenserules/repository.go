// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package licenserules

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/licenserules"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type licenseRulesRepositoryStruct struct {
	base.BaseRepositoryWithSoftDelete[*licenserules.LicenseRules]
}

func NewLicenseRulesRepository(requestSession *logy.RequestSession) ILicenseRulesRepository {
	return &licenseRulesRepositoryStruct{
		BaseRepositoryWithSoftDelete: base.CreateRepositoryWithSoftDelete[*licenserules.LicenseRules](
			requestSession,
			LicenseRulesCollectionName,
			func() *licenserules.LicenseRules {
				return &licenserules.LicenseRules{}
			},
			nil,
			nil,
			nil),
	}
}
