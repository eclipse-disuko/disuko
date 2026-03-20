// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package changeloglist

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/changeloglist"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type changeLogListRepositoryStruct struct {
	base.BaseRepositoryWithHardDelete[*changeloglist.ChangeLogList]
}

func NewChangeLogListRepository(requestSession *logy.RequestSession) IChangeLogListRepository {
	return &changeLogListRepositoryStruct{
		BaseRepositoryWithHardDelete: base.CreateRepositoryWithHardDelete[*changeloglist.ChangeLogList](
			requestSession,
			ChangeLogListCollection,
			func() *changeloglist.ChangeLogList {
				return &changeloglist.ChangeLogList{}
			},
			nil,
			"",
			nil,
			nil),
	}
}
