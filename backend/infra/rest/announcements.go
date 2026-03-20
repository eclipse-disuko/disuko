// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"

	announcement2 "github.com/eclipse-disuko/disuko/domain/announcement"
	announcement "github.com/eclipse-disuko/disuko/infra/repository/announcements"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/go-chi/render"
)

type AnnouncementsHandler struct {
	AnnouncementsRepository announcement.IAnnouncementsRepository
}

func (announcementsHandler *AnnouncementsHandler) AnnouncementsGetAllHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	announcements := announcementsHandler.AnnouncementsRepository.FindAll(requestSession, true)

	render.JSON(w, r, announcement2.ToDtos(announcements))
}
