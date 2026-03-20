// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package userstats

import (
	"strings"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain/userstats"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/infra/repository/database"
	"github.com/eclipse-disuko/disuko/logy"
)

type usersStatsRepositoryStruct struct {
	base.BaseRepositoryWithHardDelete[*userstats.UserStatus]
}

func NewUsersRepository(requestSession *logy.RequestSession) IUserStatsRepository {
	repo := &usersStatsRepositoryStruct{
		BaseRepositoryWithHardDelete: base.CreateRepositoryWithHardDelete[*userstats.UserStatus](
			requestSession,
			CollectionName,
			func() *userstats.UserStatus {
				return &userstats.UserStatus{}
			},
			nil,
			"",
			nil,
			[][]string{
				{"User"},
			}),
	}
	return repo
}

func (ur *usersStatsRepositoryStruct) FindByUserId(requestSession *logy.RequestSession, name string) *userstats.UserStatus {
	nameStr := name
	if conf.Config.OAuth2.UppercaseUsername {
		nameStr = strings.ToUpper(name)
	}
	qc := database.New().SetMatcher(
		database.AttributeMatcher(
			"User",
			database.EQ,
			nameStr))
	result := ur.Query(requestSession, qc)
	var qU *userstats.UserStatus
	if len(result) > 0 {
		qU = result[0]
	}
	return qU
}
