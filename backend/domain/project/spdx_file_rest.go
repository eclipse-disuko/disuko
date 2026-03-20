// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package project

import "time"

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
