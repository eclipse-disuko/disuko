// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package userstats

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/userstats"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

const CollectionName = "userstats"

type IUserStatsRepository interface {
	base.IBaseRepositoryWithHardDelete[*userstats.UserStatus]
	FindByUserId(requestSession *logy.RequestSession, name string) *userstats.UserStatus
}
