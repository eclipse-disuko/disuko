// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package schema

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/schema"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/database"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type schemaRepositoryStruct struct {
	base.BaseRepositoryWithHardDelete[*schema.SpdxSchema]
}

func NewSchemaRepository(requestSession *logy.RequestSession) *schemaRepositoryStruct {
	return &schemaRepositoryStruct{
		BaseRepositoryWithHardDelete: base.CreateRepositoryWithHardDelete[*schema.SpdxSchema](
			requestSession,
			SpdxSchemaCollectionName,
			func() *schema.SpdxSchema {
				return &schema.SpdxSchema{}
			},
			func(requestSession *logy.RequestSession, spdxSchema *schema.SpdxSchema) {
				exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorDbDeleteNotAllowed, SpdxSchemaCollectionName), "It is not allowed to delete a schema!")
			},
			"'content'",
			[]string{
				"content",
			},
			nil),
	}
}

func (sr *schemaRepositoryStruct) FindSpdxSchemaByNameAndVersion(requestSession *logy.RequestSession,
	name string, version string) *schema.SpdxSchema {
	qc := database.New().SetMatcher(
		database.AndChain(
			database.AttributeMatcher(
				"name",
				database.EQ,
				name),
			database.AttributeMatcher(
				"version",
				database.EQ,
				version),
		),
	).SetUnset([]string{"content"}).SetLimit(0, 1)

	qbSchemas := sr.Query(requestSession, qc)
	var qbS *schema.SpdxSchema
	if qbSchemas != nil {
		qbS = qbSchemas[0]
	}
	return qbS
}

func (sr *schemaRepositoryStruct) FindActiveSchemaByLabel(requestSession *logy.RequestSession,
	label string) *schema.SpdxSchema {
	qc := database.New().SetMatcher(
		database.AndChain(
			database.AttributeMatcher(
				"label",
				database.EQ,
				label),
			database.AttributeMatcher(
				"active",
				database.EQ,
				true),
		),
	).SetLimit(0, 1)

	qbSchemas := sr.Query(requestSession, qc)
	var qbS *schema.SpdxSchema
	if qbSchemas != nil {
		qbS = qbSchemas[0]
	}
	return qbS
}

func (sr *schemaRepositoryStruct) FindActiveSchemas(requestSession *logy.RequestSession) []*schema.SpdxSchema {
	qc := database.New().SetMatcher(
		database.AttributeMatcher(
			"active",
			database.EQ,
			true,
		),
	)
	qbS := sr.Query(requestSession, qc)
	return qbS
}

func (sr *schemaRepositoryStruct) ExistsByLabel(requestSession *logy.RequestSession, label string) bool {
	qc := database.New().SetMatcher(
		database.AttributeMatcher(
			"label",
			database.EQ,
			label,
		),
	)
	q := sr.Query(requestSession, qc)
	qbS := false
	if len(q) > 0 {
		qbS = q[0].MatchesProjectLabel(label)
	}
	logy.Infof(requestSession, "found spdx schema for label %s", label)
	return qbS
}
