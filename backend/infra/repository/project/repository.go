// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/project"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/database"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type projectRepositoryStruct struct {
	base.BaseRepositoryWithSoftDelete[*project.Project]
}

func NewProjectRepository(requestSession *logy.RequestSession) IProjectRepository {
	repo := &projectRepositoryStruct{
		BaseRepositoryWithSoftDelete: base.CreateRepositoryWithSoftDelete[*project.Project](
			requestSession,
			ProjectCollectionName,
			func() *project.Project {
				return &project.Project{
					Children: []string{},
				}
			},
			nil,
			[]string{
				"Versions",
			},
			[][]string{
				{
					"Deleted",
					"UserManagement.Users.UserId",
				},
			}),
	}

	return repo
}

func (pr *projectRepositoryStruct) CountForUser(requestSession *logy.RequestSession, userId string) int {
	qc := database.New().SetMatcher(
		database.AndChain(
			database.AttributeMatcher(
				"Deleted",
				database.EQ,
				false,
			),
			database.ArrayElemSubfieldMatcher(
				"UserManagement.Users",
				"UserId",
				database.EQ,
				userId,
			),
		),
	).SetKeep([]string{
		"_id",
	})
	qbRes := len(pr.Query(requestSession, qc))
	return qbRes
}

func (pr *projectRepositoryStruct) CountAllGroups(requestSession *logy.RequestSession) int {
	qc := database.New().SetMatcher(
		database.AndChain(
			database.AttributeMatcher(
				"Deleted",
				database.EQ,
				false,
			),
			database.AttributeMatcher(
				"IsGroup",
				database.EQ,
				true,
			),
		),
	).SetKeep([]string{
		"_id",
	})
	qbRes := len(pr.Query(requestSession, qc))
	return qbRes
}

func (pr *projectRepositoryStruct) CountGroupsForUser(requestSession *logy.RequestSession, userId string) int {
	qc := database.New().SetMatcher(
		database.AndChain(
			database.AttributeMatcher(
				"Deleted",
				database.EQ,
				false,
			),
			database.ArrayElemSubfieldMatcher(
				"UserManagement.Users",
				"UserId",
				database.EQ,
				userId,
			),
			database.AttributeMatcher(
				"IsGroup",
				database.EQ,
				true,
			),
		),
	).SetKeep([]string{
		"_id",
	})
	qbRes := len(pr.Query(requestSession, qc))
	return qbRes
}

func (pr *projectRepositoryStruct) FindRecentByUpdatedForUser(requestSession *logy.RequestSession, userId string, limit int) []*project.Project {
	qc := database.New().SetMatcher(database.AndChain(
		database.ArrayElemSubfieldMatcher(
			"UserManagement.Users",
			"UserId",
			database.EQ,
			userId,
		),
		database.AttributeMatcher(
			"Deleted",
			database.EQ,
			false),
	),
	).SetSort(database.SortConfig{
		database.SortAttribute{
			Name:  "Updated",
			Order: database.DESC,
		},
	}).SetUnset([]string{"Versions"}).SetLimit(0, limit)
	qPs := pr.Query(requestSession, qc)
	return qPs
}

func (pr *projectRepositoryStruct) FindAllForUser(requestSession *logy.RequestSession, userId string) []*project.Project {
	qc := database.New().SetMatcher(
		database.AndChain(
			database.AttributeMatcher(
				"Deleted",
				database.EQ,
				false,
			),
			database.ArrayElemSubfieldMatcher(
				"UserManagement.Users",
				"UserId",
				database.EQ,
				userId,
			),
		),
	).SetUnset([]string{"Versions"})
	qbLs := pr.Query(requestSession, qc)
	return qbLs
}

func (pr *projectRepositoryStruct) FindByKeyWithDeleted(requestSession *logy.RequestSession, key string, optimized bool) *project.Project {
	qbL := pr.BaseRepositoryWithSoftDelete.FindByKeyWithDeleted(requestSession, key, optimized)
	return qbL
}

func (pr *projectRepositoryStruct) FindByKey(requestSession *logy.RequestSession, key string, optimized bool) *project.Project {
	qc := database.New().SetMatcher(
		database.AndChain(
			database.AttributeMatcher(
				pr.DatabaseConn().GetKeyAttribute(),
				database.EQ,
				key,
			),
			database.AttributeMatcher(
				"Deleted",
				database.EQ,
				false,
			),
		),
	)
	if optimized {
		qc.SetUnset([]string{"Versions"})
	}
	qbLs := pr.Query(requestSession, qc)
	var qbL *project.Project
	if qbLs != nil {
		qbL = qbLs[0]
	}
	return qbL
}

func (pr *projectRepositoryStruct) ExistsByPolicyLabel(requestSession *logy.RequestSession, label string) bool {
	qc := database.New().SetMatcher(
		database.AndChain(
			database.ArrayElemMatcher(
				"PolicyLabels",
				database.EQ,
				label,
			),
			database.AttributeMatcher(
				"Deleted",
				database.EQ,
				false,
			),
		),
	).SetLimit(0, 1).SetKeep(
		[]string{
			"_id",
		},
	)

	qbRes := len(pr.Query(requestSession, qc)) > 0
	return qbRes
}

func (pr *projectRepositoryStruct) ExistsBySchemaLabel(requestSession *logy.RequestSession, label string) bool {
	qc := database.New().SetMatcher(
		database.AndChain(
			database.AttributeMatcher(
				"SchemaLabel",
				database.EQ,
				label,
			),
			database.AttributeMatcher(
				"Deleted",
				database.EQ,
				false,
			),
		),
	).SetLimit(0, 1).SetKeep(
		[]string{
			"_id",
		},
	)
	qbRes := len(pr.Query(requestSession, qc)) > 0
	return qbRes
}
