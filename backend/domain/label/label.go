// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package label

import (
	"strconv"
	"time"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
)

type Label struct {
	domain.RootEntity `bson:",inline"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	Type              LabelType `json:"type"`
}

func (l *Label) Update(name string, description string) {
	l.Name = name
	l.Description = description
	l.Updated = time.Now()
}

type LabelType int

const (
	SCHEMA = iota
	POLICY
	PROJECT
	INCORRECT
)

const (
	LABELTYPE_NAME_SCHEMA  = "SCHEMA"
	LABELTYPE_NAME_POLICY  = "POLICY"
	LABELTYPE_NAME_PROJECT = "PROJECT"
)

const (
	ENTERPRISE_PLATFORM = "enterprise platform"
	MOBILE_PLATFORM     = "mobile platform"
	OTHER_PLATFORM      = "other platform"
	VEHICLE_PLATFORM    = "vehicle platform"
	FRONTEND_LAYER      = "frontend layer"
	BACKEND_LAYER       = "backend layer"
	COMBINED_LAYER      = "combined layer" // for old projects
	ONBOARD             = "onboard"
	OFFBOARD            = "offboard"
	GROUP_USERS         = "group users"
	COMPANY_USERS       = "company user"
	ENTITY_USERS        = "entity users"
	EXTERNAL_USERS      = "external users" // PROJECT Labels
	CUSTOMER_USERS      = "customer"
	BP_USERS            = "business partner"
	ENTITY_TARGET       = "entity target"
	COMPANY_TARGET      = "company target"
	EXTERNAL_TARGET     = "external target"
	BP_TARGET           = "bp target"
	DUMMY               = "dummy"
	INTERNAL_DEVELOP    = "develop internal"
	EXTERNAL_DEVELOP    = "develop external"
	INHOUSE_DEVELOP     = "develop inhouse"
	COMMON_STANDARD     = "common standard"
)

func ConvertToLabelType(labelTypeStr string) LabelType {
	if labelTypeStr == LABELTYPE_NAME_SCHEMA {
		return SCHEMA
	}

	if labelTypeStr == LABELTYPE_NAME_POLICY {
		return POLICY
	}

	if labelTypeStr == LABELTYPE_NAME_PROJECT {
		return PROJECT
	}

	exception.ThrowExceptionServerMessage(message.GetI18N(message.IncorectLabelType), "")
	return INCORRECT
}

func ConvertLabelTypeToName(labelType LabelType) string {
	if labelType == SCHEMA {
		return LABELTYPE_NAME_SCHEMA
	}

	if labelType == POLICY {
		return LABELTYPE_NAME_POLICY
	}

	if labelType == PROJECT {
		return LABELTYPE_NAME_PROJECT
	}

	exception.ThrowExceptionServerMessage(message.GetI18N(message.IncorectLabelType),
		"labelType="+strconv.Itoa(int(labelType)))
	return "INCORRECT"
}

func (entity *Label) ToDto() *LabelResponseDto {
	dto := &LabelResponseDto{
		Name:        entity.Name,
		Description: entity.Description,
		Type:        ConvertLabelTypeToName(entity.Type),
	}
	domain.SetBaseValues(entity, dto)
	return dto
}
