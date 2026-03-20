// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package overallreview

import (
	"time"
)

type OverallReviewRequestDto struct {
	Comment      string `json:"comment" validate:"lte=500"`
	State        State  `json:"state" validate:"validateFn"`
	SBOMId       string `json:"sbomId"`
	SBOMName     string `json:"sbomName"`
	SBOMUploaded string `json:"sbomUploaded"`
}
type OverallReviewDto struct {
	Created         time.Time `json:"created"`
	Updated         time.Time `json:"updated"`
	Creator         string    `json:"creator"`
	CreatorFullName string    `json:"creatorFullName"`

	Comment      string `json:"comment" validate:"lte=500"`
	State        State  `json:"state" validate:"validateFn"`
	SBOMId       string `json:"sbomId"`
	SBOMName     string `json:"sbomName"`
	SBOMUploaded string `json:"sbomUploaded"`
}

func (or *OverallReview) ToDto() *OverallReviewDto {
	return &OverallReviewDto{
		Created:         or.Created,
		Updated:         or.Updated,
		Comment:         or.Comment,
		State:           or.State,
		SBOMId:          or.SBOMId,
		SBOMName:        or.SBOMName,
		SBOMUploaded:    or.SBOMUploaded,
		Creator:         or.Creator,
		CreatorFullName: or.CreatorFullName,
	}
}

type OverallReviewPublicResponse struct {
	SBOMId       string     `json:"sbomId,omitempty" example:"dummy-id-----6b9c-44a7-8e01-14e67ef4404a"`
	SBOMName     string     `json:"sbomName,omitempty" example:"SBOM Name"`
	SBOMUploaded string     `json:"sbomUploaded,omitempty" example:"2023-03-12T08:30:17.077559111Z"`
	Comment      string     `json:"comment,omitempty" example:"This is a dummy project."`
	Created      *time.Time `json:"created,omitempty" example:"2023-04-14T09:41:28.077559111Z"`
} //	@name	OverallReviewPublicResponse
