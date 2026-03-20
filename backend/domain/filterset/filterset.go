// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package filterset

import (
	"mercedes-benz.ghe.com/foss/disuko/domain"
)

type Filter struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}
type FilterSetEntity struct {
	domain.RootEntity `bson:"inline"`
	Name              string
	IncludedFilters   []Filter
	ExcludedFilters   []Filter
	TableName         string
}

type FilterSetDto struct {
	Key             string   `json:"_key"`
	Name            string   `json:"name"`
	IncludedFilters []Filter `json:"includedFilters"`
	ExcludedFilters []Filter `json:"excludedFilters"`
	TableName       string   `json:"tableName"`
}

func (dto *FilterSetDto) ToEntity() *FilterSetEntity {
	return &FilterSetEntity{
		RootEntity:      domain.NewRootEntity(),
		Name:            dto.Name,
		IncludedFilters: dto.IncludedFilters,
		ExcludedFilters: dto.ExcludedFilters,
		TableName:       dto.TableName,
	}
}

func ToDto(entity *FilterSetEntity) *FilterSetDto {
	return &FilterSetDto{
		Key:             entity.Key,
		Name:            entity.Name,
		IncludedFilters: entity.IncludedFilters,
		ExcludedFilters: entity.ExcludedFilters,
		TableName:       entity.TableName,
	}
}

func ToDtos(source []*FilterSetEntity) []*FilterSetDto {
	return domain.MapTo(source, ToDto)
}

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message" example:"Resource created"`
}
