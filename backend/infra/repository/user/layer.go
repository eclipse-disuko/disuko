// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package user

import (
	"github.com/eclipse-disuko/disuko/domain/user"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
)

const UsersCollectionName = "users"

type IUsersRepository interface {
	base.IBaseRepositoryWithHardDelete[*user.User]
	Find5UsersBySearchFragment(requestSession *logy.RequestSession, searchFragment string, userFilter UserFilter) []*user.User
	FindByUserId(requestSession *logy.RequestSession, name string) *user.User
}
