// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package userstats

import (
	"github.com/eclipse-disuko/disuko/domain/userstats"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
)

const CollectionName = "userstats"

type IUserStatsRepository interface {
	base.IBaseRepositoryWithHardDelete[*userstats.UserStatus]
	FindByUserId(requestSession *logy.RequestSession, name string) *userstats.UserStatus
}
