// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"html"
	"net/http"
	"strings"
	"unicode"

	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"

	appConnector "github.com/eclipse-disuko/disuko/connector/application"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/go-chi/render"
)

type ApplicationHandler struct {
	Connector *appConnector.Connector
}

func (handler *ApplicationHandler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	rawQuery := r.URL.Query().Get("query")
	if strings.TrimSpace(rawQuery) == "" || len(rawQuery) < 3 {
		exception.ThrowExceptionClient400Message(message.GetI18N(message.RequestApp), "")
	}
	query := strings.TrimLeftFunc(rawQuery, unicode.IsSpace)

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
