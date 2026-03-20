// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain/department"
	"github.com/eclipse-disuko/disuko/domain/oauth"
	"github.com/eclipse-disuko/disuko/domain/project/approvable"
	"github.com/eclipse-disuko/disuko/domain/schema"
)

type ProjectDto struct {
	Key                   string                       `json:"_key"`
	Name                  string                       `json:"name"`
	Versions              map[string]VersionSlimDto    `json:"versions"`
	Description           string                       `json:"description"`
	SchemaLabel           string                       `json:"schemaLabel"`
	PolicyLabels          []string                     `json:"policyLabels"`
	ProjectLabels         []string                     `json:"projectLabels"`
	FreeLabels            []string                     `json:"freeLabels"`
	Created               time.Time                    `json:"created,omitempty"`
	Updated               time.Time                    `json:"updated,omitempty"`
	CorrespondingSchema   *schema.SpdxSchemaDto        `json:"correspondingSchema"`
	Token                 []TokenDto                   `json:"token"` // todo #later: deliver with separate request on demand
	Status                ProjectStatusType            `json:"status"`
	DocumentMeta          DisclosureDocumentMetaDto    `json:"documentMeta"`
	NoticeContactMeta     NoticeContactMetaDto         `json:"noticeContactMeta"`
	CustomerMeta          CustomerMetaDto              `json:"customerMeta"`
	AccessRights          oauth.AccessAndRolesRights   `json:"accessRights"`
	Children              []string                     `json:"children"`
	IsGroup               bool                         `json:"isGroup"`
	Parent                string                       `json:"parent"`
	ParentName            string                       `json:"parentName"`
	SupplierExtraData     SupplierExtraDataDto         `json:"supplierExtraData"`
	ApprovableSPDX        approvable.ApprovableSPDXDto `json:"approvablespdx"`
	ParentProjectSettings ProjectSettingsDto           `json:"parentProjectSettings"`
	IsNoFoss              bool                         `json:"isNoFoss"`
	ApplicationMeta       ApplicationMetaDto           `json:"applicationMeta"`
	Subscriptions         *SubscriptionsDto            `json:"subscriptions"`
	Responsible           string                       `json:"responsible"`
	CustomIds             []ProjectCustomIdDto         `json:"customIds"`
	IsDummy               bool                         `json:"isDummy"`
	DummyDeletionDate     time.Time                    `json:"dummyDeletionDate"`
	DeleteDisabledReason  string                       `json:"deleteDisabledReason,omitempty"`
	HasApproval           bool                         `json:"hasApproval"`
	HasChildren           bool                         `json:"hasChildren"`
	HasSBOMToRetain       bool                         `json:"hasSBOMToRetain"`
}

type ProjectChildrenCombiDto struct {
	Project              ProjectSlimDto  `json:"project"`
	Version              *VersionSlimDto `json:"version"`
	ProjectKey           string          `json:"projectKey"`
	HasProjectReadAccess bool            `json:"hasProjectReadAccess"`
}

// ProjectSlimDto used as DTO for Project representation in projects list where only small amount of the entire Project structure is needed.
type ProjectSlimDto struct {
	Key                  string                     `json:"_key"`
	Rev                  string                     `json:"_rev"`
	Name                 string                     `json:"name"`
	Description          string                     `json:"description"`
	ApplicationId        *string                    `json:"applicationId"`
	SchemaLabel          string                     `json:"schemaLabel"`
	PolicyLabels         []string                   `json:"policyLabels"`
	ProjectLabels        []string                   `json:"projectLabels"`
	FreeLabels           []string                   `json:"freeLabels"`
	Children             []string                   `json:"children"`
	IsGroup              bool                       `json:"isGroup"`
	Parent               string                     `json:"parent"`
	ParentName           string                     `json:"parentName"`
	Status               ProjectStatusType          `json:"status"`
	Updated              time.Time                  `json:"updated,omitempty"`
	Created              time.Time                  `json:"created,omitempty"`
	AccessRights         oauth.AccessAndRolesRights `json:"accessRights"`
	IsDeleted            bool                       `json:"isDeleted"`
	Supplier             string                     `json:"supplier"`
	SupplierMissing      bool                       `json:"supplierMissing"`
	Company              string                     `json:"company"`
	Department           string                     `json:"department"`
	DepartmentMissing    bool                       `json:"missing"`
	IsNoFoss             bool                       `json:"isNoFoss"`
	ApplicationMeta      ApplicationMetaDto         `json:"applicationMeta"`
	IsInGroupApproval    bool                       `json:"isInGroupApproval"`
	Responsible          string                     `json:"responsible"`
	CustomIds            []ProjectCustomIdDto       `json:"customIds"`
	IsDummy              bool                       `json:"isDummy"`
	DeleteDisabledReason string                     `json:"deleteDisabledReason,omitempty"`
	HasApproval          bool                       `json:"hasApproval"`
	HasChildren          bool                       `json:"hasChildren"`
	HasSBOMToRetain      bool                       `json:"hasSBOMToRetain"`
}

func (entity *Project) ToSlimDto(docDep *department.Department, docMissing bool, custDep *department.Department, custMissing bool, isDummy bool) ProjectSlimDto {
	departmentStr := ""
	companyStr := ""
	supplier := ""
	if entity.SupplierExtraData.External {
		supplier = entity.DocumentMeta.SupplierName
	} else if docDep != nil {
		supplier = docDep.CompanyName + " [" + docDep.CompanyCode + "]"
	}

	if custDep != nil {
		departmentStr = custDep.OrgAbbreviation + " " + custDep.DescriptionEnglish + " [" + custDep.Key + "]"
		companyStr = custDep.CompanyName + " [" + custDep.CompanyCode + "]"
	}
	applicationId := ""
	if entity.ApplicationMeta.Id != "" {
		applicationId = entity.ApplicationMeta.Name
		if entity.ApplicationMeta.SecondaryId != "" {
			applicationId += " (" + entity.ApplicationMeta.SecondaryId + ")"
		}
	}
	if len(applicationId) == 0 && entity.ApplicationId != nil && *entity.ApplicationId != "" {
		applicationId = *entity.ApplicationId
	}
	responsible := ""
	if ru := entity.ProjectResponsible(); ru != nil {
		responsible = ru.UserId
	}
	customIds := make([]ProjectCustomIdDto, 0)
	for _, c := range entity.CustomIds {
		customIds = append(customIds, ProjectCustomIdDto{
			Key:         c.Key,
			TechnicalId: c.TechnicalId,
			Value:       c.Value,
		})
	}

	return ProjectSlimDto{
		Key:               entity.Key,
		Rev:               entity.Rev,
		Name:              entity.Name,
		ApplicationId:     &applicationId,
		Description:       entity.Description,
		SchemaLabel:       entity.SchemaLabel,
		PolicyLabels:      entity.PolicyLabels,
		ProjectLabels:     entity.ToProjectLabelsDto(),
		FreeLabels:        entity.FreeLabels,
		Children:          entity.Children,
		Parent:            entity.Parent,
		ParentName:        entity.ParentName,
		Updated:           entity.Updated,
		Created:           entity.Created,
		Status:            entity.Status,
		IsGroup:           entity.IsGroup,
		IsDeleted:         entity.Deleted,
		Supplier:          supplier,
		SupplierMissing:   !entity.SupplierExtraData.External && docMissing,
		Department:        departmentStr,
		DepartmentMissing: custMissing,
		Company:           companyStr,
		IsNoFoss:          entity.IsNoFoss,
		Responsible:       responsible,
		ApplicationMeta:   entity.ApplicationMeta.ToDto(),
		CustomIds:         customIds,
		IsDummy:           isDummy,
		HasChildren:       entity.HasChildren,
		HasApproval:       entity.HasApproval,
		HasSBOMToRetain:   entity.HasSBOMToRetain,
	}
}

func (entity *Project) ToDto(docDep *department.Department, docMissing bool, custDep *department.Department, custMissing bool, isDummy bool) ProjectDto {
	versions := make(map[string]VersionSlimDto)
	for key, version := range entity.Versions {
		if version.Deleted {
			continue
		}
		versions[key] = *version.ToDto()
	}
	token := make([]TokenDto, 0)
	for _, t := range entity.Token {
		token = append(token, t.ToDto())
	}
	responsible := ""
	if ru := entity.ProjectResponsible(); ru != nil {
		responsible = ru.UserId
	}
	cids := make([]ProjectCustomIdDto, 0)
	for _, cid := range entity.CustomIds {
		cids = append(cids, ProjectCustomIdDto{
			Key:         cid.Key,
			TechnicalId: cid.TechnicalId,
			Value:       cid.Value,
		})
	}

	var del time.Time
	if isDummy {
		del = entity.Created.UTC().AddDate(0, 3, 0)
	}
	return ProjectDto{
		Key:                 entity.Key,
		Name:                entity.Name,
		Versions:            versions,
		Description:         entity.Description,
		SchemaLabel:         entity.SchemaLabel,
		PolicyLabels:        entity.PolicyLabels,
		ProjectLabels:       entity.ToProjectLabelsDto(),
		FreeLabels:          entity.FreeLabels,
		Created:             entity.Created,
		Updated:             entity.Updated,
		CorrespondingSchema: entity.CorrespondingSchema.ToDto(),
		Status:              entity.Status,
		DocumentMeta:        entity.DocumentMeta.ToDto(docDep, docMissing),
		CustomerMeta:        entity.CustomerMeta.ToDto(custDep, custMissing),
		NoticeContactMeta:   entity.NoticeContactMeta.ToDto(),
		Token:               token,
		Parent:              entity.Parent,
		ParentName:          entity.ParentName,
		Children:            entity.Children,
		IsGroup:             entity.IsGroup,
		SupplierExtraData:   entity.SupplierExtraData.ToDTO(),
		ApprovableSPDX:      entity.ApprovableSPDX.ToDto(),
		IsNoFoss:            entity.IsNoFoss,
		ApplicationMeta:     ApplicationMetaDto(entity.ApplicationMeta),
		Responsible:         responsible,
		CustomIds:           cids,
		IsDummy:             isDummy,
		DummyDeletionDate:   del,
		HasApproval:         entity.HasApproval,
		HasChildren:         entity.HasChildren,
		HasSBOMToRetain:     entity.HasSBOMToRetain,
	}
}

func (entity *Project) ToProjectLabelsDto() []string {
	if entity.ProjectLabels == nil {
		return []string{}
	}
	return entity.ProjectLabels
}

func (entity *Project) ToProjectSettingsDto(docDep *department.Department, docMissing bool, custDep *department.Department, custMissing bool) ProjectSettingsDto {
	return ProjectSettingsDto{
		DocumentMeta:      entity.DocumentMeta.ToDto(docDep, docMissing),
		CustomerMeta:      entity.CustomerMeta.ToDto(custDep, custMissing),
		SupplierExtraData: entity.SupplierExtraData.ToDTO(),
		NoticeContactMeta: entity.NoticeContactMeta.ToDto(),
	}
}

type ApplicationMetaDto struct {
	Id           string `json:"id"`
	SecondaryId  string `json:"secondaryId"`
	Name         string `json:"name"`
	ExternalLink string `json:"externalLink"`
}

func (dto ApplicationMetaDto) ToEntity() ApplicationMeta {
	return ApplicationMeta(dto)
}

func (dto ApplicationMetaDto) Diff(cmp ApplicationMeta) bool {
	return cmp.Name != dto.Name || cmp.ExternalLink != dto.ExternalLink || cmp.Id != dto.Id || cmp.SecondaryId != dto.SecondaryId
}

type ProjectChildMemberCombiDto struct {
	ProjectKey           string            `json:"projectKey"`
	ProjectName          string            `json:"projectName"`
	ProjectMember        *ProjectMemberDto `json:"projectMember"`
	UserManagementRights *oauth.CRUDRights `json:"userManagementRights"`
}
type ProjectChildrenMembersDto struct {
	List []ProjectChildMemberCombiDto `json:"list"`
}

type ProjectChildrenMemberSuccessResponseDto struct {
	Success     bool   `json:"success"`
	UserId      string `json:"userId"`
	ProjectKey  string `json:"projectKey"`
	ProjectName string `json:"projectName"`
	Message     string `json:"message"`
}
