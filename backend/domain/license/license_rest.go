// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package license

import (
	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/obligation"
	obligation2 "github.com/eclipse-disuko/disuko/infra/repository/obligation"
	"github.com/eclipse-disuko/disuko/logy"
)

type AliasDto struct {
	Key         string `json:"_key" validate:"lte=36"`
	LicenseId   string `json:"licenseId"`
	Description string `json:"description"`
}

type MetaDataDto struct {
	Family          FamilyOfLicense             `json:"family" validate:"lte=80"`
	ApprovalState   ApprovalStatus              `json:"approvalState" validate:"lte=80"`
	ReviewState     ReviewStatus                `json:"reviewState" validate:"lte=80"`
	ReviewDateStr   string                      `json:"reviewDate,omitempty"`
	ObligationsList []*obligation.ObligationDto `json:"obligationsList"`
	LicenseUrl      string                      `json:"licenseUrl" validate:"lte=2000"`
	SourceUrl       string                      `json:"sourceUrl" validate:"lte=2000"`
	OSIApproved     bool                        `json:"osiApproved"`
	FSFApproved     bool                        `json:"fsfApproved"`
	Changelog       string                      `json:"changelog" validate:"lte=100000"`
	LicenseType     TypeOfLicenses              `json:"licenseType" validate:"lte=50"`
	Evaluation      string                      `json:"evaluation" validate:"lte=100000"`
	LegalComments   string                      `json:"legalComments" validate:"lte=100000"`
	IsLicenseChart  bool                        `json:"isLicenseChart"`
}

type MetaDataSlimDto struct {
	Family                       FamilyOfLicense                `json:"family" validate:"lte=80"`
	ApprovalState                ApprovalStatus                 `json:"approvalState" validate:"lte=80"`
	LicenseType                  TypeOfLicenses                 `json:"licenseType" validate:"lte=50"`
	IsLicenseChart               bool                           `json:"isLicenseChart"`
	Classifications              []obligation.ObligationSlimDto `json:"classifications"`
	PrevalentClassificationLevel obligation.WarnLevel           `json:"prevalentClassificationLevel" validate:"lte=50"`
}

type LicenseDto struct {
	domain.BaseDto
	IsDeprecatedLicenseId bool        `json:"isDeprecatedLicenseId"`
	LicenseId             string      `json:"licenseId" validate:"gte=2,lte=80"`
	Name                  string      `json:"name" validate:"required,gte=3,lte=100"`
	Text                  string      `json:"text" validate:"lte=100000"`
	Aliases               []AliasDto  `json:"aliases"`
	Meta                  MetaDataDto `json:"meta"`
	Source                Source      `json:"source" validate:"lte=80"`
}

type LicenseNameIdDto struct {
	LicenseId string `json:"licenseId"`
	Name      string `json:"name"`
}

func (entity *License) ToNameIdDto() *LicenseNameIdDto {
	return &LicenseNameIdDto{
		LicenseId: entity.LicenseId,
		Name:      entity.Name,
	}
}

type LicenseSlimDto struct {
	domain.BaseDto
	LicenseId string          `json:"licenseId" validate:"gte=2,lte=80"`
	Name      string          `json:"name" validate:"required,gte=3,lte=80"`
	Meta      MetaDataSlimDto `json:"meta"`
	Source    Source          `json:"source" validate:"lte=80"`
	Aliases   []AliasDto      `json:"aliases"`
}

func (dto *MetaDataDto) ToEntity() *MetaData {
	obligationKeyList := make([]string, 0)
	for _, obligation := range dto.ObligationsList {
		obligationKeyList = append(obligationKeyList, obligation.Key)
	}
	return &MetaData{
		Family:             dto.Family,
		ApprovalState:      dto.ApprovalState,
		ReviewState:        dto.ReviewState,
		ReviewDateStr:      dto.ReviewDateStr,
		ObligationsKeyList: obligationKeyList,
		LicenseUrl:         dto.LicenseUrl,
		SourceUrl:          dto.SourceUrl,
		OSIApproved:        dto.OSIApproved,
		FSFApproved:        dto.FSFApproved,
		Changelog:          dto.Changelog,
		LicenseType:        dto.LicenseType,
		Evaluation:         dto.Evaluation,
		LegalComments:      dto.LegalComments,
		IsLicenseChart:     dto.IsLicenseChart,
	}
}

func (dto *AliasDto) ToEntity() *Alias {
	return &Alias{
		ChildEntity: domain.ChildEntity{
			Key: dto.Key,
		},
		LicenseId:   dto.LicenseId,
		Description: dto.Description,
	}
}

func (entity *Alias) ToDto() *AliasDto {
	return &AliasDto{
		Key:         entity.Key,
		LicenseId:   entity.LicenseId,
		Description: entity.Description,
	}
}

func AliasesToEntity(dtos []AliasDto) []Alias {
	var res []Alias
	for _, alias := range dtos {
		res = append(res, *alias.ToEntity())
	}
	return res
}

func (dto *LicenseDto) ToEntity() *License {
	result := License{
		LicenseId: dto.LicenseId,
		Name:      dto.Name,
		Text:      dto.Text,
		Meta:      *dto.Meta.ToEntity(),
		Source:    dto.Source,
	}
	result.Rev = dto.Rev
	result.Key = dto.Key
	result.Created = dto.Created
	result.Updated = dto.Updated
	result.Aliases = AliasesToEntity(dto.Aliases)
	return &result
}

func (entity *MetaData) ToSlimDto() *MetaDataSlimDto {
	return &MetaDataSlimDto{
		Family:                       entity.Family,
		ApprovalState:                entity.ApprovalState,
		LicenseType:                  entity.LicenseType,
		IsLicenseChart:               entity.IsLicenseChart,
		Classifications:              []obligation.ObligationSlimDto{},
		PrevalentClassificationLevel: "",
	}
}

func (entity *MetaData) ToDto(requestSession *logy.RequestSession,
	obligationProvider obligation2.IObligationRepository) *MetaDataDto {

	obligations := make([]*obligation.ObligationDto, 0)
	if obligationProvider != nil {
		//with loaded obligations?
		for _, obligationKey := range entity.ObligationsKeyList {
			obligation := obligationProvider.FindByKey(requestSession, obligationKey, false)
			if obligation != nil {
				obligations = append(obligations, obligation.ToDto())
			}
		}
	}
	return &MetaDataDto{
		Family:          entity.Family,
		ApprovalState:   entity.ApprovalState,
		ReviewState:     entity.ReviewState,
		ReviewDateStr:   entity.ReviewDateStr,
		ObligationsList: obligations,
		LicenseUrl:      entity.LicenseUrl,
		SourceUrl:       entity.SourceUrl,
		OSIApproved:     entity.OSIApproved,
		FSFApproved:     entity.FSFApproved,
		Changelog:       entity.Changelog,
		LicenseType:     entity.LicenseType,
		Evaluation:      entity.Evaluation,
		LegalComments:   entity.LegalComments,
		IsLicenseChart:  entity.IsLicenseChart,
	}
}
func (entity *License) ToSlimDto() *LicenseSlimDto {
	res := &LicenseSlimDto{
		LicenseId: entity.LicenseId,
		Name:      entity.Name,
		Meta:      *entity.Meta.ToSlimDto(),
		Source:    entity.Source,
		Aliases:   []AliasDto{},
	}
	for _, alias := range entity.Aliases {
		res.Aliases = append(res.Aliases, *alias.ToDto())
	}
	domain.SetBaseValues(entity, res)
	return res
}

func (entity *License) ToDtoWithoutObligations() *LicenseDto {
	return entity.ToDto(nil, nil)
}

func (entity *License) ToDto(requestSession *logy.RequestSession,
	obligationProvider obligation2.IObligationRepository) *LicenseDto {

	meta := entity.Meta.ToDto(requestSession, obligationProvider)

	res := &LicenseDto{
		IsDeprecatedLicenseId: entity.IsDeprecatedLicenseId,
		LicenseId:             entity.LicenseId,
		Name:                  entity.Name,
		Text:                  entity.Text,
		Meta:                  *meta,
		Source:                entity.Source,
	}

	domain.SetBaseValues(entity, res)

	for _, alias := range entity.Aliases {
		res.Aliases = append(res.Aliases, *alias.ToDto())
	}

	return res
}

func ToSlimDtos(licenses []*License) []LicenseSlimDto {
	dtos := make([]LicenseSlimDto, 0)
	for _, license := range licenses {
		dto := license.ToSlimDto()
		dtos = append(dtos, *dto)
	}
	return dtos
}

func ToDtos(requestSession *logy.RequestSession, licenses []*License,
	obligationProvider obligation2.IObligationRepository) []LicenseDto {
	dtos := make([]LicenseDto, 0)
	for _, license := range licenses {
		dto := license.ToDto(requestSession, obligationProvider)
		dtos = append(dtos, *dto)
	}
	return dtos
}

type LicenseDiffDto struct {
	LicenseId  string      `json:"licenseId"`
	OldLicense *LicenseDto `json:"oldLicense"`
	NewLicense *LicenseDto `json:"newLicense"`
}

type LicenseWithSimilarityDto struct {
	License    LicenseDto `json:"license"`
	Similarity float64    `json:"similarity"`
}

type ClassificationWithCount struct {
	Classification *obligation.ObligationDto `json:"classification"`
	Count          int                       `json:"count"`
}

type PossibleFilterValues struct {
	PossibleCharts          map[string]int            `json:"possibleCharts"`
	PossibleSources         map[string]int            `json:"possibleSources"`
	PossibleFamilies        map[string]int            `json:"possibleFamilies"`
	PossibleApproval        map[string]int            `json:"possibleApproval"`
	PossibleType            map[string]int            `json:"possibleType"`
	PossibleClassifications []ClassificationWithCount `json:"possibleClassifications"`
}

type LicensesResponse struct {
	Licenses []*LicenseSlimDto    `json:"licenses"`
	Count    int                  `json:"count"`
	Meta     PossibleFilterValues `json:"meta"`
}
