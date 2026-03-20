// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"fmt"
	"strconv"

	"mercedes-benz.ghe.com/foss/disuko/conf"
	"mercedes-benz.ghe.com/foss/disuko/connector/application"
	"mercedes-benz.ghe.com/foss/disuko/connector/department"
	"mercedes-benz.ghe.com/foss/disuko/connector/userrole"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type connectors struct {
	application *application.Connector
	userrole    *userrole.Connector
	department  *department.Connector
}

func (s *Server) setupConnector(rs *logy.RequestSession) {
	if conf.Config.Connector.Userrole.Scheme != "" && conf.Config.Connector.Userrole.Host != "" && conf.Config.Connector.Userrole.Port > 0 {
		userroleEndpoint := fmt.Sprintf("%s://%s:%s",
			conf.Config.Connector.Userrole.Scheme,
			conf.Config.Connector.Userrole.Host,
			strconv.Itoa(conf.Config.Connector.Userrole.Port),
		)
		s.connectors.userrole = userrole.Init(rs, userroleEndpoint)
	}
	if conf.Config.Connector.Application.Scheme != "" && conf.Config.Connector.Application.Host != "" && conf.Config.Connector.Application.Port > 0 {
		applicationEndpoint := fmt.Sprintf("%s://%s:%s",
			conf.Config.Connector.Application.Scheme,
			conf.Config.Connector.Application.Host,
			strconv.Itoa(conf.Config.Connector.Application.Port),
		)
		s.connectors.application = application.Init(rs, applicationEndpoint)
	}
	if conf.Config.Connector.Department.Scheme != "" && conf.Config.Connector.Department.Host != "" && conf.Config.Connector.Department.Port > 0 {
		departmentEndpoint := fmt.Sprintf("%s://%s:%s",
			conf.Config.Connector.Department.Scheme,
			conf.Config.Connector.Department.Host,
			strconv.Itoa(conf.Config.Connector.Department.Port),
		)
		s.connectors.department = department.Init(rs, departmentEndpoint)
	}
}
