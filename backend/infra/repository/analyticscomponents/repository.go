// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package analyticscomponents

import (
	"slices"
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
				{"Name"},
			},
		),
	}
	return componentsRepositoryStruct
}

func (r *componentsRepositoryStruct) InitIndex(requestSession *logy.RequestSession) {
	qc := database.New()
	qACs := r.Query(requestSession, qc)
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
	return helper.Search(r.index, name, exact)
}

func (r *componentsRepositoryStruct) ExistByName(requestSession *logy.RequestSession, name string) bool {
	return slices.Contains(r.index, name)
}

func (r *componentsRepositoryStruct) AddToIndex(requestSession *logy.RequestSession, name string) {
	r.index = append(r.index, name)
}
