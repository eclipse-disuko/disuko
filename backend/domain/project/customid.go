// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package project

import "github.com/eclipse-disuko/disuko/domain"

type ProjectCustomIdDto struct {
	Key         string `json:"_key" validate:"lte=36"`
	TechnicalId string `json:"technicalId" validate:"required,gte=3,lte=36"`
	Value       string `json:"value" validate:"required,gte=3,lte=80"`
}

type ProjectCustomId struct {
	domain.ChildEntity `bson:",inline"`
	TechnicalId        string
	Value              string
}
