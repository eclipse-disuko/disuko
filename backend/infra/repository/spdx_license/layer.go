// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package spdx_license

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/license"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
)

const SpdxLicensesCollectionName = "spdxLicenses"

type ISpdxLicensesRepository interface {
	base.IBaseRepositoryWithHardDelete[*license.License]
}
