// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/project"
	"mercedes-benz.ghe.com/foss/disuko/domain/project/sbomlist"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/approvallist"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/auditloglist"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/labels"
	license2 "mercedes-benz.ghe.com/foss/disuko/infra/repository/license"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/licenserules"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/policyrules"
	project2 "mercedes-benz.ghe.com/foss/disuko/infra/repository/project"
	sbomListRepo "mercedes-benz.ghe.com/foss/disuko/infra/repository/sbomlist"
	schema2 "mercedes-benz.ghe.com/foss/disuko/infra/repository/schema"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/user"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type spdxRetriever interface {
	RetrieveSbomListAndFile(*logy.RequestSession, string, string) (*sbomlist.SbomList, *project.SpdxFileBase)
}

type RepositoryHolder struct {
	PolicyRulesRepository  policyrules.IPolicyRulesRepository
	ProjectRepository      project2.IProjectRepository
	LicenseRepository      license2.ILicensesRepository
	LabelRepository        labels.ILabelRepository
	SBOMListRepository     sbomListRepo.ISbomListRepository
	SchemaRepository       schema2.ISchemaRepository
	UserRepository         user.IUsersRepository
	AuditLogListRepository auditloglist.IAuditLogListRepository
	ApprovalListRepository approvallist.IApprovalListRepository
	Retriever              spdxRetriever
	LicenseRulesRepository licenserules.ILicenseRulesRepository
}
