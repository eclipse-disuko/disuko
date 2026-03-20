// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package spdx_license

import (
	"github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
)

const SpdxLicensesCollectionName = "spdxLicenses"

type ISpdxLicensesRepository interface {
	base.IBaseRepositoryWithHardDelete[*license.License]
}
