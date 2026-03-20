// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package newsbox

import (
	"github.com/eclipse-disuko/disuko/domain/newsbox"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
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
