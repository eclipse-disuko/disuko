// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package announcement

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/announcement"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
)

const AnnouncementsCollectionName = "announcements"

type IAnnouncementsRepository interface {
	base.IBaseRepositoryWithHardDelete[*announcement.Announcement]
}
