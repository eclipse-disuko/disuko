// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package customid

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/customid"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
)

const collName = "customids"

type ICustomIdRepository interface {
	base.IBaseRepositoryWithHardDelete[*customid.CustomId]
}
