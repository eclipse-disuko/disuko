// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package analytics

import (
	"mercedes-benz.ghe.com/foss/disuko/domain"
	"mercedes-benz.ghe.com/foss/disuko/domain/license"
)

type Occurrence struct {
	domain.RootEntity `bson:"inline"`
	domain.SoftDelete `bson:"inline"`
	OrigName          string
	ReferencedLicense string
	Count             int
}

func (entity *Occurrence) ToDto(lic *license.License) *OccurrenceDto {
	var res OccurrenceDto
	res.Count = entity.Count
	res.OrigName = entity.OrigName
	res.ReferencedLicense = entity.ReferencedLicense
	if lic != nil {
		res.License = lic.ToSlimDto()
	}
	return &res
}
