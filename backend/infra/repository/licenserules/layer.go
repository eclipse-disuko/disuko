// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package licenserules

import (
	"github.com/eclipse-disuko/disuko/domain/licenserules"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
)

const LicenseRulesCollectionName string = "licenserules"

type ILicenseRulesRepository interface {
	base.IBaseRepositoryWithSoftDelete[*licenserules.LicenseRules]
}
