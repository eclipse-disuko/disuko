// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/connector/application"
	"github.com/eclipse-disuko/disuko/connector/department"
	"github.com/eclipse-disuko/disuko/domain/capabilities"
	"github.com/go-chi/render"
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
