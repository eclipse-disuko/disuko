// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain/overallreview"
)

type MetaInfoDto struct {
	Name        string `json:"name"`
	SpdxId      string `json:"spdxId"`
	SpdxVersion string `json:"spdxVersion"`
	Comment     string `json:"comment"`
}

func (metaInfo MetaInfo) ToDto() MetaInfoDto {
	return MetaInfoDto{
		Name:        metaInfo.Name,
		SpdxId:      metaInfo.SpdxId,
		SpdxVersion: metaInfo.SpdxVersion,
		Comment:     metaInfo.Comment,
	}
}

type ApprovalInfoDto struct {
	IsInApproval bool   `json:"isInApproval"`
	Comment      string `json:"comment"`
	Status       string `json:"status"`
}

func (approvalInfo ApprovalInfo) ToDto() ApprovalInfoDto {
	return ApprovalInfoDto{
		IsInApproval: approvalInfo.IsInApproval,
		Comment:      approvalInfo.Comment,
		Status:       approvalInfo.Status,
	}
}

type SpdxFileDto struct {
	Key          string          `json:"_key"`
	Created      time.Time       `json:"created"`
	Updated      time.Time       `json:"updated"`
	Hash         string          `json:"hash"`
	MetaInfo     MetaInfoDto     `json:"metaInfo"`
	ApprovalInfo ApprovalInfoDto `json:"approvalInfo"`
	Uploaded     *time.Time      `json:"uploaded"`
	Origin       string          `json:"origin"`
	Uploader     string          `json:"uploader"`
	Tag          string          `json:"tag"`

	OverallReview *overallreview.OverallReviewDto `json:"overallReview"`

	IsInUse  bool `json:"isInUse"`
	IsLocked bool `json:"isLocked"`

	IsToDelete bool `json:"isToDelete"`
	IsToRetain bool `json:"isToRetain"`
}

func (entity *SpdxFileBase) ToDto() *SpdxFileDto {
	var overallReviewDto *overallreview.OverallReviewDto
	if entity.OverallReview != nil {
		overallReviewDto = entity.OverallReview.ToDto()
	}
	return &SpdxFileDto{
		Key:           entity.Key,
		Created:       entity.Created,
		Updated:       entity.Updated,
		Hash:          entity.Hash,
		MetaInfo:      entity.MetaInfo.ToDto(),
		ApprovalInfo:  entity.ApprovalInfo.ToDto(),
		Uploaded:      entity.Uploaded,
		Origin:        entity.Origin,
		Uploader:      entity.Uploader,
		Tag:           entity.Tag,
		OverallReview: overallReviewDto,
		IsInUse:       entity.IsInUse,
		IsLocked:      entity.IsLocked,
	}
}

type SpdxFileSlimDto struct {
	Key              string     `json:"_key"`
	ProjectVersionId string     `json:"projectVersionId"`
	Uploaded         *time.Time `json:"uploaded,omitempty"`
	Name             string     `json:"name"`
}

func (entity *SpdxFileBase) ToSlimDto(versionKey string) *SpdxFileSlimDto {
	return &SpdxFileSlimDto{
		Key:              entity.Key,
		ProjectVersionId: versionKey,
		Uploaded:         entity.Uploaded,
		Name:             entity.MetaInfo.Name,
	}
}

type SbomKnownLicenseDto struct {
	Id     string `json:"id"`
	OrigId string `json:"origId"`
	Name   string `json:"name"`
}

type SbomLicensesDto struct {
	Unknown []string              `json:"unknown"`
	Known   []SbomKnownLicenseDto `json:"known"`
}

type NameKeyIdentifier struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type ResponseFlatSbomItem struct {
	VersionKey  string `json:"versionKey"`
	VersionName string `json:"versionName"`
	*SpdxFileDto
}

type ResponseFlatSboms struct {
	Items   []ResponseFlatSbomItem `json:"items"`
	Version []NameKeyIdentifier    `json:"versions"`
}
