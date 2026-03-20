// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"
	"net/url"
	"regexp"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"mercedes-benz.ghe.com/foss/disuko/domain"
	"mercedes-benz.ghe.com/foss/disuko/domain/customid"
	"mercedes-benz.ghe.com/foss/disuko/domain/project"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/helper/roles"
	"mercedes-benz.ghe.com/foss/disuko/helper/validation"
	customidRepo "mercedes-benz.ghe.com/foss/disuko/infra/repository/customid"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/database"
	projectRepo "mercedes-benz.ghe.com/foss/disuko/infra/repository/project"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

var keyRX = regexp.MustCompile("^[a-zA-Z0-9-]{4,36}$")

type CustomidHandler struct {
	Repo        customidRepo.ICustomIdRepository
	ProjectRepo projectRepo.IProjectRepository
}

func extractCustomIdBody(r *http.Request) customid.CustomIdDto {
	var res customid.CustomIdDto
	validation.DecodeAndValidate(r, &res, false)
	return res
}

func extractCustomIdKey(r *http.Request) string {
	basicAuthUUID := chi.URLParam(r, "uuid")

	authUUID, err := url.QueryUnescape(basicAuthUUID)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ParamUuidWrong))

	return authUUID
}

func (h *CustomidHandler) List(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	ids := h.Repo.FindAll(requestSession, true)
	dtos := make([]customid.CustomIdDto, 0)
	for _, e := range ids {
		dtos = append(dtos, e.ToDto())
	}
	render.JSON(w, r, dtos)
}

func (h *CustomidHandler) Create(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.IsDomainAdmin()) {
		exception.ThrowExceptionSendDeniedResponse()
	}

	data := extractCustomIdBody(r)
	if !keyRX.MatchString(data.Key) {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.CustomIdKeyMalformed), "")
	}
	e := data.ToEntity()
	e.RootEntity = domain.NewRootEntityWithKey(data.Key)
	h.Repo.Save(requestSession, &e)
	res := SuccessResponse{
		Success: true,
	}
	render.JSON(w, r, res)
}

func (h *CustomidHandler) Update(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.IsApplicationAdmin()) {
		exception.ThrowExceptionSendDeniedResponse()
	}

	data := extractCustomIdBody(r)
	key := extractCustomIdKey(r)

	existing := h.Repo.FindByKey(requestSession, key, false)
	if existing == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorDbNotFound), "")
	}

	existing.Name = data.Name
	existing.NameDE = data.NameDE
	existing.Description = data.Description
	existing.DescriptionDE = data.DescriptionDE
	existing.LinkTemplate = data.LinkTemplate

	h.Repo.Update(requestSession, existing)

	res := SuccessResponse{
		Success: true,
	}
	render.JSON(w, r, res)
}

func (h *CustomidHandler) Delete(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.IsApplicationAdmin()) {
		exception.ThrowExceptionSendDeniedResponse()
	}

	uuid := extractCustomIdKey(r)

	qc := database.New().SetMatcher(
		database.ArrayElemSubfieldMatcher(
			"CustomIds",
			"TechnicalId",
			database.EQ,
			uuid,
		),
	)

	prs := h.ProjectRepo.Query(requestSession, qc)
	for _, p := range prs {
		cleaned := make([]project.ProjectCustomId, 0)
		for _, c := range p.CustomIds {
			if c.TechnicalId != uuid {
				cleaned = append(cleaned, c)
			}
		}
		p.CustomIds = cleaned
		h.ProjectRepo.Update(requestSession, p)
	}

	h.Repo.Delete(requestSession, uuid)

	res := SuccessResponse{
		Success: true,
	}
	render.JSON(w, r, res)
}

func (h *CustomidHandler) Usage(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.IsApplicationAdmin()) {
		exception.ThrowExceptionSendDeniedResponse()
	}

	uuid := extractCustomIdKey(r)
	qc := database.New().SetMatcher(
		database.ArrayElemSubfieldMatcher(
			"CustomIds",
			"TechnicalId",
			database.EQ,
			uuid,
		),
	).SetKeep(
		[]string{h.ProjectRepo.DatabaseConn().GetKeyAttribute()},
	)
	prs := h.ProjectRepo.Query(requestSession, qc)

	res := customid.CustomIdUsage{
		Count: len(prs),
	}

	render.JSON(w, r, res)
}
