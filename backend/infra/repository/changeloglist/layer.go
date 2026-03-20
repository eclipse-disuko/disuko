// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package changeloglist

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/changeloglist"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
)

const ChangeLogListCollection = "changeloglist"

type IChangeLogListRepository interface {
	base.IBaseRepositoryWithHardDelete[*changeloglist.ChangeLogList]
}
