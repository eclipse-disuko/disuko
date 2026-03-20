// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package newsbox

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/newsbox"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type newsboxRepo struct {
	base.BaseRepositoryWithHardDelete[*newsbox.Item]
}

func NewNewsboxRepository(requestSession *logy.RequestSession) IRepo {
	return &newsboxRepo{
		BaseRepositoryWithHardDelete: base.CreateRepositoryWithHardDelete(
			requestSession,
			collName,
			func() *newsbox.Item {
				return &newsbox.Item{}
			},
			nil,
			"",
			nil,
			nil,
		),
	}
}
