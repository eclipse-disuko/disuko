// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package changelogs

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/changeloglist"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
)

const ChangeLogsCollectionName string = "changelogs"

type IChangeLogsRepository interface {
	base.IBaseRepositoryWithHardDelete[*changeloglist.ChangeLog]
}
