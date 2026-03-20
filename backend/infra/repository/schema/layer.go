// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package schema

import (
	"github.com/eclipse-disuko/disuko/domain/schema"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
)

const SpdxSchemaCollectionName = "spdxSchemas"

type ISchemaRepository interface {
	base.IBaseRepositoryWithHardDelete[*schema.SpdxSchema]

	FindSpdxSchemaByNameAndVersion(requestSession *logy.RequestSession, name string, version string) *schema.SpdxSchema
	FindActiveSchemaByLabel(requestSession *logy.RequestSession, label string) *schema.SpdxSchema
	FindActiveSchemas(requestSession *logy.RequestSession) []*schema.SpdxSchema
	ExistsByLabel(requestSession *logy.RequestSession, label string) bool
}
