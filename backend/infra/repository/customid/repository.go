// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package customid

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/customid"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type customIdRepositoryStruct struct {
	base.BaseRepositoryWithHardDelete[*customid.CustomId]
}

func NewLabelsRepository(requestSession *logy.RequestSession) ICustomIdRepository {
	return &customIdRepositoryStruct{
		BaseRepositoryWithHardDelete: base.CreateRepositoryWithHardDelete(
			requestSession,
			collName,
			func() *customid.CustomId {
				return &customid.CustomId{}
			},
			nil,
			"",
			nil,
			nil),
	}
}
