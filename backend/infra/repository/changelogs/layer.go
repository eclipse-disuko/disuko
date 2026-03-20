// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package changelogs

import (
	"github.com/eclipse-disuko/disuko/domain/changeloglist"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
)

const ChangeLogsCollectionName string = "changelogs"

type IChangeLogsRepository interface {
	base.IBaseRepositoryWithHardDelete[*changeloglist.ChangeLog]
}
