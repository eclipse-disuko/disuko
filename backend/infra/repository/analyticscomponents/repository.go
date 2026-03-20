// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package analyticscomponents

import (
	"strings"

	"github.com/eclipse-disuko/disuko/helper"

	"github.com/eclipse-disuko/disuko/domain/analytics"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/infra/repository/database"
	"github.com/eclipse-disuko/disuko/logy"
)

type ComponentsRepository struct{}

type componentsRepositoryStruct struct {
	base.BaseRepositoryWithSoftDelete[*analytics.Component]
	index []string
}

func NewComponentsRepository(requestSession *logy.RequestSession) IComponentsRepository {
	componentsRepositoryStruct := &componentsRepositoryStruct{
		BaseRepositoryWithSoftDelete: base.CreateRepositoryWithSoftDelete[*analytics.Component](
			requestSession,
			ComponentsCollectionName,
			func() *analytics.Component {
				return &analytics.Component{}
			},
			nil,
			nil,
			[][]string{
				{"Name", "Version"},
				{"Name"},
			},
		),
	}
	return componentsRepositoryStruct
}

func (r *componentsRepositoryStruct) search(search string, exact bool) []string {
	return helper.Search(r.index, search, exact)
}

func (r *componentsRepositoryStruct) InitIndex(requestSession *logy.RequestSession) {
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

func (r *componentsRepositoryStruct) SearchByName(requestSession *logy.RequestSession, name string, exact bool) []string {
	return r.search(name, exact)
}

func (r *componentsRepositoryStruct) FindByNameAndVersion(requestSession *logy.RequestSession, name, version string) []*analytics.Component {
	qc := database.New().SetMatcher(database.AndChain(
		database.AttributeMatcher(
			"Name",
			database.EQ,
			name,
		),
		database.AttributeMatcher(
			"Version",
			database.EQ,
			version,
		),
	))
	return r.Query(requestSession, qc)
}

func (r *componentsRepositoryStruct) AddToIndex(requestSession *logy.RequestSession, name string) {
	r.index = append(r.index, name)
}
