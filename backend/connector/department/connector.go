// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

//go:generate oapi-codegen --config=../../resources/connectorapis/department-config.yaml ../../resources/connectorapis/department.yaml

package department

import (
	"context"
	"fmt"

	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/logy"
)

type Connector struct {
	client *ClientWithResponses
}

func Init(rs *logy.RequestSession, host string) *Connector {
	c, err := NewClientWithResponses(host)
	if err != nil {
		logy.Fatalf(rs, "department client init failed: %s", err)
	}
	return &Connector{
		client: c,
	}
}

func (c *Connector) GetDepartments(rs *logy.RequestSession) Departments {
	res, err := c.client.GetDepartmentsWithResponse(context.Background())
	if err != nil {
		exception.HandleErrorServerMessage(err, message.GetI18N(message.ConnectorReqFailed))
	}
	if res.StatusCode() != 200 {
		exception.HandleErrorServerMessage(fmt.Errorf("non-200 status code: %d", res.StatusCode()), message.GetI18N(message.ConnectorReqFailed))
	}
	return *res.JSON200
}
