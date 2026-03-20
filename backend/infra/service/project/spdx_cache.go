// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/project/sbomlist"
	"github.com/eclipse-disuko/disuko/infra/repository/approvallist"
	"github.com/eclipse-disuko/disuko/infra/repository/auditloglist"
	"github.com/eclipse-disuko/disuko/infra/repository/labels"
	license2 "github.com/eclipse-disuko/disuko/infra/repository/license"
	"github.com/eclipse-disuko/disuko/infra/repository/licenserules"
	"github.com/eclipse-disuko/disuko/infra/repository/policyrules"
	project2 "github.com/eclipse-disuko/disuko/infra/repository/project"
	sbomListRepo "github.com/eclipse-disuko/disuko/infra/repository/sbomlist"
	schema2 "github.com/eclipse-disuko/disuko/infra/repository/schema"
	"github.com/eclipse-disuko/disuko/infra/repository/user"
	"github.com/eclipse-disuko/disuko/logy"
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
