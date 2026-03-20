// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

//go:generate oapi-codegen --config=../../resources/connectorapis/application-config.yaml ../../resources/connectorapis/application.yaml

package application

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
		logy.Fatalf(rs, "userrole client init failed: %s", err)
	}
	return &Connector{
		client: c,
	}
}

func (c *Connector) Search(rs *logy.RequestSession, query string) SearchRes {
	res, err := c.client.PostSearchWithResponse(context.Background(), SearchQuery{
		Query: query,
	})
	if err != nil {
		exception.HandleErrorServerMessage(err, message.GetI18N(message.ConnectorReqFailed))
	}
	if res.StatusCode() != 200 {
		exception.HandleErrorServerMessage(fmt.Errorf("non-200 status code: %d", res.StatusCode()), message.GetI18N(message.ConnectorReqFailed))
	}
	return *res.JSON200
}

func (c *Connector) GetApplication(rs *logy.RequestSession, id string) Application {
	res, err := c.client.GetApplicationIdWithResponse(context.Background(), id)
	if err != nil {
		exception.HandleErrorServerMessage(err, message.GetI18N(message.ConnectorReqFailed))
	}
	if res.StatusCode() != 200 {
		exception.HandleErrorServerMessage(fmt.Errorf("non-200 status code: %d", res.StatusCode()), message.GetI18N(message.ConnectorReqFailed))
	}
	return *res.JSON200
}
