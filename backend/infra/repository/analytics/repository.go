// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package analytics

import (
	"github.com/eclipse-disuko/disuko/domain/analytics"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/infra/repository/database"
	"github.com/eclipse-disuko/disuko/logy"
)

type AnalyticsRepository struct{}

type analyticsRepositoryStruct struct {
	base.BaseRepositoryWithSoftDelete[*analytics.Analytics]
}

func NewAnalyticsRepository(requestSession *logy.RequestSession) IAnalyticsRepository {
	analyticsRepository := &analyticsRepositoryStruct{
		BaseRepositoryWithSoftDelete: base.CreateRepositoryWithSoftDelete[*analytics.Analytics](
			requestSession,
			AnalyticsCollectionName,
			func() *analytics.Analytics {
				return &analytics.Analytics{}
			},
			nil,
			nil,
			[][]string{
				{"ProjectVersionKey"},
				{"SBomKey"},
				{"SBomName"},
				{"ProjectVersionName"},
				{"ProjectName"},
				{"ProjectKey"},
				{"ComponentVersion"},
				{"Deleted", "ComponentName"},
				{"Deleted", "EntryLicense"},
				{"Deleted", "ProjectKey"},
				{"Deleted", "ComponentName", "ProjectKey"},
				{"Deleted", "EntryLicense", "ProjectKey"},
				{"Deleted", "EntryLicense", "ComponentName"},
				{"Deleted", "ProjectKey", "EntryLicense", "ComponentName"},
				{"Deleted", "OwnerCompanyId"},
				{"OwnerCompanyId"},
			},
		),
	}
	return analyticsRepository
}

func (repository *analyticsRepositoryStruct) FindByNameAndProjectKeysAndLicense(requestSession *logy.RequestSession, component string, keys []string, entryLicense string, offset, limit int, sortCol string, asc bool) []*analytics.Analytics {
	qc := database.New().SetMatcher(
		database.AttributeMatcher(
			"Deleted",
			database.EQ,
			false))

	if component != "" {
		qc.SetMatcher(database.AndChain(
			*qc.Matcher,
			database.AttributeMatcher(
				"ComponentName",
				database.EQ,
				component,
			)))
	}
	if entryLicense != "" {
		qc.SetMatcher(database.AndChain(
			*qc.Matcher,
			database.AttributeMatcher(
				"EntryLicense",
				database.EQ,
				entryLicense,
			)))
	}
	if len(keys) > 0 {
		var projectGroup []database.MatchGroup
		for _, key := range keys {
			projectGroup = append(projectGroup, database.AttributeMatcher(
				"ProjectKey",
				database.EQ,
				key,
			))
		}
		qc.SetMatcher(database.AndChain(
			*qc.Matcher,
			database.OrChain(
				projectGroup...,
			),
		))

	}

	if sortCol != "" {
		order := database.ASC
		if !asc {
			order = database.DESC
		}
		qc.SetSort(database.SortConfig{
			database.SortAttribute{
				Name:  sortCol,
				Order: order,
			},
		})
	}
	if limit > 0 {
		qc.SetLimit(offset, limit)
	}
	qbAs := repository.Query(requestSession, qc)
	return qbAs
}
