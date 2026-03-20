// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package analyticslicenses

import (
	"strings"

	"mercedes-benz.ghe.com/foss/disuko/helper"

	"mercedes-benz.ghe.com/foss/disuko/domain/analytics"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/database"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type licensesRepositoryStruct struct {
	base.BaseRepositoryWithSoftDelete[*analytics.License]
	index []string
}

func NewLicensesRepository(requestSession *logy.RequestSession) *licensesRepositoryStruct {
	licensesRepositoryStruct := &licensesRepositoryStruct{
		BaseRepositoryWithSoftDelete: base.CreateRepositoryWithSoftDelete[*analytics.License](
			requestSession,
			LicensesCollectionName,
			func() *analytics.License {
				return &analytics.License{}
			},
			nil,
			nil,
			[][]string{
				{"Name"},
			},
		),
	}
	return licensesRepositoryStruct
}
func (r *licensesRepositoryStruct) InitIndex(requestSession *logy.RequestSession) {
	qc := database.New()
	qACs := r.BaseRepository.Query(requestSession, qc)
	var names []string
	for _, c := range qACs {
		if strings.TrimSpace(c.Name) == "" {
			continue
		}
		names = append(names, c.Name)
	}
	r.index = names
}

func (r *licensesRepositoryStruct) search(search string, exact bool) []string {
	return helper.Search(r.index, search, exact)
}

func (r *licensesRepositoryStruct) SearchLicenceByName(requestSession *logy.RequestSession, name string, exact bool) []string {
	return r.search(name, exact)
}

func (r *licensesRepositoryStruct) FindByName(requestSession *logy.RequestSession, name string) []*analytics.License {
	qc := database.New().SetMatcher(database.AndChain(
		database.AttributeMatcher(
			"Name",
			database.EQ,
			name,
		),
	))
	return r.Query(requestSession, qc)
}
func (r *licensesRepositoryStruct) AddToIndex(requestSession *logy.RequestSession, name string) {
	r.index = append(r.index, name)
}
