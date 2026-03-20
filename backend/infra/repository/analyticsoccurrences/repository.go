// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package analyticsoccurrences

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/analytics"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
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
