// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package obligation

import (
	"mercedes-benz.ghe.com/foss/disuko/domain"
)

type ObligationDto struct {
	domain.BaseDto
	Name          string    `json:"name" validate:"required,gte=3,lte=80"`
	NameDe        string    `json:"nameDe" validate:"required,gte=3,lte=80"`
	Type          Type      `json:"type" validate:"required,gt=1,lte=20"`
	WarnLevel     WarnLevel `json:"warnLevel" validate:"required,gt=1,lte=20"`
	Description   string    `json:"description" validate:"lte=1000"`
	DescriptionDe string    `json:"descriptionDe" validate:"lte=1000"`
	AutoApproved  bool      `json:"autoApproved"`
}

type ObligationSlimDto struct {
	domain.BaseDto
	Name      string    `json:"name" validate:"required,gte=3,lte=80"`
	NameDe    string    `json:"nameDe" validate:"required,gte=3,lte=80"`
	Type      Type      `json:"type" validate:"required,gt=1,lte=20"`
	WarnLevel WarnLevel `json:"warnLevel" validate:"required,gt=1,lte=20"`
}

func (dto *ObligationDto) ToEntity() *Obligation {
	entity := &Obligation{
		Name:          dto.Name,
		NameDe:        dto.NameDe,
		Type:          dto.Type,
		WarnLevel:     dto.WarnLevel,
		Description:   dto.Description,
		DescriptionDe: dto.DescriptionDe,
		AutoApproved:  dto.AutoApproved,
	}
	domain.SetBaseValues(dto, entity)
	return entity
}

func (dto *ObligationDto) ToSlimDto() *ObligationSlimDto {
	return &ObligationSlimDto{
		BaseDto:   domain.BaseDto{},
		Name:      dto.Name,
		NameDe:    dto.NameDe,
		Type:      dto.Type,
		WarnLevel: dto.WarnLevel,
	}
}

func (entity *Obligation) ToDto() *ObligationDto {
	return ToDto(entity)
}

func ToDto(entity *Obligation) *ObligationDto {
	dto := &ObligationDto{
		Name:          entity.Name,
		NameDe:        entity.NameDe,
		Type:          entity.Type,
		WarnLevel:     entity.WarnLevel,
		Description:   entity.Description,
		DescriptionDe: entity.DescriptionDe,
		AutoApproved:  entity.AutoApproved,
	}
	domain.SetBaseValues(entity, dto)
	return dto
}

func ToDtos(source []*Obligation) []*ObligationDto {
	return domain.MapTo(source, ToDto)
}

type AllResponse struct {
	Obligation []*ObligationDto `json:"items"`
	Count      int              `json:"count"`
}

type Request struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}
