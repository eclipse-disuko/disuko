// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain/overallreview"
)

type VersionSlimDto struct {
	Key             string                            `json:"_key"`
	ParentKey       string                            `json:"parentKey"`
	Name            string                            `json:"name"`
	Description     string                            `json:"description"`
	Created         time.Time                         `json:"created,omitempty"`
	Updated         time.Time                         `json:"updated,omitempty"`
	CurrentSpdxFile *SpdxFileSlimDto                  `json:"currentSpdxFile"`
	SpdxFileHistory []*SpdxFileSlimDto                `json:"spdxFileHistory"`
	Status          ProjectVersionStatusType          `json:"status"`
	IsDeleted       bool                              `json:"isDeleted"`
	OverallReviews  []*overallreview.OverallReviewDto `json:"overallReviews"`
}

func (entity *ProjectVersion) ToDtoWithParentKey(parentKey *string) *VersionSlimDto {
	dto := VersionSlimDto{
		Key:             entity.Key,
		Name:            entity.Name,
		Description:     entity.Description,
		CurrentSpdxFile: &SpdxFileSlimDto{},
		Created:         entity.Created,
		Updated:         entity.Updated,
		Status:          entity.Status,
		IsDeleted:       entity.Deleted,
	}
	if parentKey != nil {
		dto.ParentKey = *parentKey
	}

	if entity.OverallReviews != nil {
		overallReviews := make([]*overallreview.OverallReviewDto, 0)
		for _, overallReview := range entity.OverallReviews {
			overallReviews = append(overallReviews, overallReview.ToDto())
		}
		dto.OverallReviews = overallReviews
	}

	return &dto
}
func (entity *ProjectVersion) ToDto() *VersionSlimDto {
	return entity.ToDtoWithParentKey(nil)
}
