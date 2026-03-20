// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package analyticsoccurrences

import (
	"github.com/eclipse-disuko/disuko/domain/analytics"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
)

type occurrencesRepositoryStruct struct {
	base.BaseRepositoryWithSoftDelete[*analytics.Occurrence]
}

func NewLicensesRepository(requestSession *logy.RequestSession) *occurrencesRepositoryStruct {
	occurrencesRepositoryStruct := &occurrencesRepositoryStruct{
		BaseRepositoryWithSoftDelete: base.CreateRepositoryWithSoftDelete[*analytics.Occurrence](
			requestSession,
			OccurrencesCollectionName,
			func() *analytics.Occurrence {
				return &analytics.Occurrence{}
			},
			nil,
			nil,
			[][]string{
				{"ReferencedLicense"},
				{"OrigName"},
			},
		),
	}
	return occurrencesRepositoryStruct
}
