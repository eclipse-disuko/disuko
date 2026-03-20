// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package checklist

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/checklist"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type checklistRepositoryStruct struct {
	base.BaseRepositoryWithHardDelete[*checklist.Checklist]
}

func NewLabelsRepository(requestSession *logy.RequestSession) IChecklistRepository {
	return &checklistRepositoryStruct{
		BaseRepositoryWithHardDelete: base.CreateRepositoryWithHardDelete(
			requestSession,
			collName,
			func() *checklist.Checklist {
				return &checklist.Checklist{}
			},
			nil,
			"",
			nil,
			nil),
	}
}
