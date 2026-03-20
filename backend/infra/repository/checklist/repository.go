// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package checklist

import (
	"github.com/eclipse-disuko/disuko/domain/checklist"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
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
