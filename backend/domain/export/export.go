// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package export

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/label"
	"mercedes-benz.ghe.com/foss/disuko/domain/license"
	"mercedes-benz.ghe.com/foss/disuko/domain/obligation"
	"mercedes-benz.ghe.com/foss/disuko/domain/schema"
)

type ExportSchemaKnowledgeBaseDto struct {
	SchemaLabels []*label.Label       `validate:"required"`
	Schemas      []*schema.SpdxSchema `validate:"required"`
}

type ExportLicenseKnowledgeBaseDto struct {
	PolicyLabels []*label.Label           `validate:"required"`
	PolicyRules  []*license.PolicyRules   `validate:"required"`
	Obligations  []*obligation.Obligation `validate:"required"`
	Licenses     []*license.License       `validate:"required"`
}

type ImportResultDto struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}
