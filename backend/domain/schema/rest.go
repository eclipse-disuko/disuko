// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package schema

import "mercedes-benz.ghe.com/foss/disuko/domain"

type SpdxSchemaDto struct {
	domain.BaseDto
	Name        string     `json:"name"`
	Version     string     `json:"version"`
	Description string     `json:"description"`
	Type        SchemaType `json:"type"`
	Active      bool       `json:"active"`
	Content     string     `json:"content,omitempty"`
	Label       string     `json:"label,omitempty"` //key of the label
}

type SchemaRequestDto struct {
	Name        string `json:"name" validate:"required,gte=3,lte=80"`
	Version     string `json:"version" validate:"required,gte=3,lte=80"`
	Description string `json:"description" validate:"lte=1000"`
	Label       string `json:"label" validate:"required,gte=3,lte=80"` //key of the label
}

type StatusResponseDto struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
