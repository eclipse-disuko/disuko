// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package obligation

import (
	"github.com/eclipse-disuko/disuko/domain/obligation"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/infra/repository/database"
	"github.com/eclipse-disuko/disuko/logy"
)

type obligationRepositoryStruct struct {
	base.BaseRepositoryWithSoftDelete[*obligation.Obligation]
}

func NewObligationRepository(requestSession *logy.RequestSession) IObligationRepository {
	return &obligationRepositoryStruct{
		BaseRepositoryWithSoftDelete: base.CreateRepositoryWithSoftDelete[*obligation.Obligation](
			requestSession,
			ObligationCollectionName,
			func() *obligation.Obligation {
				return &obligation.Obligation{}
			},
			nil,
			nil,
			[][]string{
				{"name"},
			}),
	}
}

func (pr *obligationRepositoryStruct) FindAllSortedByName(requestSession *logy.RequestSession) []*obligation.Obligation {
	qc := database.New().SetMatcher(
		database.AttributeMatcher(
			"Deleted",
			database.EQ,
			false,
		),
	).SetSort(database.SortConfig{
		database.SortAttribute{
			Name:  "name",
			Order: database.ASC,
		},
	})

	qbObligations := pr.Query(requestSession, qc)
	return qbObligations
}

func (pr *obligationRepositoryStruct) FindByName(requestSession *logy.RequestSession, name string) []*obligation.Obligation {
	qc := database.New().SetMatcher(
		database.AndChain(
			database.AttributeMatcher(
				"name",
				database.EQ,
				name,
			),
			database.AttributeMatcher(
				"Deleted",
				database.EQ,
				false,
			),
		),
	)

	qbObligations := pr.Query(requestSession, qc)
	return qbObligations
}

/*

func (pr *obligationRepositoryStruct) FindByKeys(requestSession *logy.RequestSession, keys []string) []*obligation.Obligation {
	return pr.FindByKeys(requestSession, keys)
		var res []*obligation.Obligation
		for _, key := range keys {
			qc := database.New().SetMatcher(
				database.AttributeMatcher(
					pr.DatabaseConn().GetKeyAttribute(),
					database.EQ,
					key,
				),
			)
			qbLs := pr.Query(requestSession, qc)
			if qbLs != nil {
				res = append(res, qbLs[0])
			}
		}
		return res
}
*/
