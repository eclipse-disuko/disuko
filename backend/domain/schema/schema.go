// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package schema

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain"
)

type SpdxSchema struct {
	domain.RootEntity `bson:",inline"`
	Name              string     `json:"name"`
	Version           string     `json:"version"`
	Description       string     `json:"description"`
	Type              SchemaType `json:"type"`
	Active            bool       `json:"active"`
	Content           string     `json:"content,omitempty"`
	Label             string     `json:"label,omitempty"` // key of the label
}

type SchemaType int

const (
	JSON = iota
	OCTET
)

func (schema *SpdxSchema) ActivateSchema() {
	if !schema.Active {
		schema.Updated = time.Now()
		schema.Active = true
	}
}

func (schema *SpdxSchema) DeactivateSchema() {
	if schema.Active {
		schema.Updated = time.Now()
		schema.Active = false
	}
}

func (schema *SpdxSchema) MatchesProjectLabel(label string) bool {
	return schema.Label == label
}

func (schema *SpdxSchema) ContentTypeAsString() string {
	if schema.Type == JSON {
		return "application/json"
	} else {
		return "application/octet-stream"
	}
}

func (schema *SpdxSchema) ToDto() *SpdxSchemaDto {
	return ToDto(schema)
}

func ToDto(schema *SpdxSchema) *SpdxSchemaDto {
	if schema == nil {
		return nil
	}
	dto := &SpdxSchemaDto{
		Name:        schema.Name,
		Version:     schema.Version,
		Description: schema.Description,
		Type:        schema.Type,
		Active:      schema.Active,
		Content:     schema.Content,
		Label:       schema.Label,
	}
	domain.SetBaseValues(schema, dto)
	return dto
}

func ToSpdxSchemaDtoList(source []*SpdxSchema) []*SpdxSchemaDto {
	return domain.MapTo(source, ToDto)
}
