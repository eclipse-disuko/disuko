// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package changelogs

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/changeloglist"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type changeLogsRepositoryStruct struct {
	base.BaseRepositoryWithHardDelete[*changeloglist.ChangeLog]
}

func NewChangeLogsRepository(requestSession *logy.RequestSession) IChangeLogsRepository {
	return &changeLogsRepositoryStruct{
		BaseRepositoryWithHardDelete: base.CreateRepositoryWithHardDelete[*changeloglist.ChangeLog](
			requestSession,
			ChangeLogsCollectionName,
			func() *changeloglist.ChangeLog {
				return &changeloglist.ChangeLog{}
			},
			nil,
			"",
			nil,
			nil),
	}
}
