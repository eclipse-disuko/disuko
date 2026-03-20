// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"html"
	"net/http"
	"strings"

	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"

	"github.com/go-chi/render"
	appConnector "mercedes-benz.ghe.com/foss/disuko/connector/application"
	"mercedes-benz.ghe.com/foss/disuko/domain/project"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type ApplicationHandler struct {
	Connector *appConnector.Connector
}

func (handler *ApplicationHandler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	query := strings.TrimSpace(r.URL.Query().Get("query"))
	if len(query) < 4 {
		exception.ThrowExceptionClient400Message(message.GetI18N(message.RequestApp), "")
	}

	response := make([]*project.ApplicationMetaDto, 0)
	if handler.Connector == nil {
		render.JSON(w, r, response)
		return
	}

	appRes := handler.Connector.Search(requestSession, html.EscapeString(query))
	for _, a := range appRes {
		response = append(response, &project.ApplicationMetaDto{
			Name:         a.Name,
			Id:           a.Id,
			ExternalLink: a.Link,
		})
	}
	render.JSON(w, r, response)
}
