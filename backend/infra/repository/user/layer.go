// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package user

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/user"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

const UsersCollectionName = "users"

type IUsersRepository interface {
	base.IBaseRepositoryWithHardDelete[*user.User]
	Find5UsersBySearchFragment(requestSession *logy.RequestSession, searchFragment string, userFilter UserFilter) []*user.User
	FindByUserId(requestSession *logy.RequestSession, name string) *user.User
}
