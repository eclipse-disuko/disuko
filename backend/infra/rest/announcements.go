// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"

	"github.com/go-chi/render"
	announcement2 "mercedes-benz.ghe.com/foss/disuko/domain/announcement"
	announcement "mercedes-benz.ghe.com/foss/disuko/infra/repository/announcements"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type AnnouncementsHandler struct {
	AnnouncementsRepository announcement.IAnnouncementsRepository
}

func (announcementsHandler *AnnouncementsHandler) AnnouncementsGetAllHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	announcements := announcementsHandler.AnnouncementsRepository.FindAll(requestSession, true)

	render.JSON(w, r, announcement2.ToDtos(announcements))
}
