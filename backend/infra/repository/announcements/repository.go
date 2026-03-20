// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package announcement

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/announcement"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type announcementsRepositoryStruct struct {
	base.BaseRepositoryWithHardDelete[*announcement.Announcement]
}

func NewAnnouncementsRepository(requestSession *logy.RequestSession) IAnnouncementsRepository {
	return &announcementsRepositoryStruct{
		BaseRepositoryWithHardDelete: base.CreateRepositoryWithHardDelete[*announcement.Announcement](
			requestSession,
			AnnouncementsCollectionName,
			func() *announcement.Announcement {
				return &announcement.Announcement{}
			},
			nil,
			"",
			nil,
			nil),
	}
}
