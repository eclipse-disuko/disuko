// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package user

import (
	"strings"

	"github.com/eclipse-disuko/disuko/conf"

	"github.com/eclipse-disuko/disuko/domain/user"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/infra/repository/database"
	"github.com/eclipse-disuko/disuko/logy"
)

type UserFilter struct {
	Active *bool
}
type usersRepositoryStruct struct {
	base.BaseRepositoryWithHardDelete[*user.User]
}

func NewUsersRepository(requestSession *logy.RequestSession) IUsersRepository {
	repo := &usersRepositoryStruct{
		BaseRepositoryWithHardDelete: base.CreateRepositoryWithHardDelete[*user.User](
			requestSession,
			UsersCollectionName,
			func() *user.User {
				return &user.User{}
			},
			nil,
			"",
			nil,
			[][]string{
				{"Forename"},
				{"User"},
			}),
	}
	repo.PreUpdate = func(requestSession *logy.RequestSession, u *user.User) {
		loadedUser := repo.FindByUserId(requestSession, u.User)
		if loadedUser == nil {
			return
		}
		if loadedUser.Key == u.Key {
			return
		}
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorUserIdInUse, u.User), "")
	}
	repo.PreSave = func(requestSession *logy.RequestSession, u *user.User) {
		loadedUser := repo.FindByUserId(requestSession, u.User)
		if loadedUser == nil {
			return
		}
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorUserIdInUse, u.User), "")
	}
	return repo
}

func (ur *usersRepositoryStruct) Find5UsersBySearchFragment(requestSession *logy.RequestSession, searchFragment string, userFilter UserFilter) []*user.User {
	searchFragment = strings.TrimSpace(searchFragment)
	searchFragments := strings.Fields(searchFragment)

	filter := ""
	qc := &database.QueryConfig{}
	if len(searchFragments) == 1 {
		qc = database.New().SetMatcher(
			database.OrChain(
				database.AttributeMatcher(
					"User",
					database.LIKEI,
					"%"+searchFragment+"%",
				),
				database.AttributeMatcher(
					"Forename",
					database.LIKEI,
					"%"+searchFragment+"%",
				),
				database.AttributeMatcher(
					"Lastname",
					database.LIKEI,
					"%"+searchFragment+"%",
				),
			),
		)

		// searching for user || forename || lastname, as only 1 input is given
		filter += "(LIKE(entity.User, @searchFragment, true) OR LIKE(entity.Forename, @searchFragment, true) OR LIKE(entity.Lastname, @searchFragment, true))"
	} else {
		qc = database.New().SetMatcher(
			database.OrChain(
				database.AndChain(
					database.AttributeMatcher(
						"Forename",
						database.LIKEI,
						"%"+strings.TrimSpace(searchFragments[0])+"%",
					),
					database.AttributeMatcher(
						"Lastname",
						database.LIKEI,
						"%"+strings.TrimSpace(searchFragments[1])+"%",
					),
				),
				database.AndChain(
					database.AttributeMatcher(
						"Forename",
						database.LIKEI,
						"%"+strings.TrimSpace(searchFragments[1])+"%",
					),
					database.AttributeMatcher(
						"Lastname",
						database.LIKEI,
						"%"+strings.TrimSpace(searchFragments[0])+"%",
					),
				),
			),
		)
		// searching for forename && lastname, as 2 inputs are given
		filter += "((LIKE(entity.Forename, @searchFragment0, true) AND LIKE(entity.Lastname, @searchFragment1, true)) OR (LIKE(entity.Forename, @searchFragment1, true) AND LIKE(entity.Lastname, @searchFragment0, true)))"
	}

	if userFilter.Active != nil {
		qc.SetMatcher(
			database.AndChain(*qc.Matcher,
				database.AttributeMatcher(
					"Active",
					database.EQ,
					*userFilter.Active,
				),
			),
		)
	}
	qc.SetLimit(0, 5).SetSort(
		database.SortConfig{
			database.SortAttribute{
				Name:  "Forename",
				Order: database.ASC,
			},
		},
	)

	qbUsers := ur.Query(requestSession, qc)

	return qbUsers
}

func (ur *usersRepositoryStruct) FindByUserId(requestSession *logy.RequestSession, name string) *user.User {
	nameStr := name
	if conf.Config.OAuth2.UppercaseUsername {
		nameStr = strings.ToUpper(name)
	}
	qc := database.New().SetMatcher(
		database.AttributeMatcher(
			"User",
			database.EQ,
			nameStr))
	qUser := ur.Query(requestSession, qc)
	var qU *user.User
	if len(qUser) > 0 {
		qU = qUser[0]
	}
	return qU
}
