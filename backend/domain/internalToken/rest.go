// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package internalToken

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain"
)

type InternalTokenDto struct {
	domain.BaseDto
	Name         string       `json:"name" validate:"gte=2,lte=80"`
	Revoked      bool         `json:"revoked"`
	Token        string       `json:"token"`
	Description  string       `json:"description" validate:"lte=1000"`
	Expiry       time.Time    `json:"expiry" validate:"required"`
	Capabilities []Capability `json:"capabilities"  validate:"required,gte=0"`
}

type TokenAuth struct {
	Key   string `json:"key"`
	Token string `json:"token"`
}

func (entity *InternalToken) ToDto() InternalTokenDto {
	res := InternalTokenDto{
		Name:         entity.Name,
		Revoked:      entity.Revoked,
		Description:  entity.Description,
		Expiry:       entity.Expiry,
		Capabilities: entity.Capabilities,
	}
	domain.SetBaseValues(entity, &res)
	return res
}
