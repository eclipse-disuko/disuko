// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"

	"github.com/go-chi/render"
	"mercedes-benz.ghe.com/foss/disuko/conf"
	"mercedes-benz.ghe.com/foss/disuko/connector/application"
	"mercedes-benz.ghe.com/foss/disuko/connector/department"
	"mercedes-benz.ghe.com/foss/disuko/domain/capabilities"
)

type CapabilitiesHandler struct {
	ApplicationConnector *application.Connector
	DepartmentConnector  *department.Connector
}

func (h *CapabilitiesHandler) GetCapabilities(w http.ResponseWriter, r *http.Request) {
	response := capabilities.CapabilitiesDto{
		ApplicationConnector:          h.ApplicationConnector != nil,
		DepartmentConnector:           h.DepartmentConnector != nil,
		EnforceFOSSOfficeConfirmation: conf.Config.Server.EnforceFOSSOfficeConfirmation,
	}
	render.JSON(w, r, response)
}
