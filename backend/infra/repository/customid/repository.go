// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package customid

import (
	"github.com/eclipse-disuko/disuko/domain/customid"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
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
