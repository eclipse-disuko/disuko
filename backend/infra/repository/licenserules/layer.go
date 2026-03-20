// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package licenserules

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/licenserules"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
)

const LicenseRulesCollectionName string = "licenserules"

type ILicenseRulesRepository interface {
	base.IBaseRepositoryWithSoftDelete[*licenserules.LicenseRules]
}
