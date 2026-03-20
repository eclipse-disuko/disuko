// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain"
)

type SourceExternal struct {
	domain.ChildEntity `bson:",inline"`
	URL                string
	Comment            string
	Origin             string
	Uploader           string
}

type SourceType string //	@name	SourceType

const (
	ExternalLink SourceType = "external"
)

func (dto *SourceExternalDTO) ToEntity() SourceExternal {
	return SourceExternal{
		ChildEntity: domain.NewChildEntity(),
		URL:         dto.URL,
		Comment:     dto.Comment,
	}
}

func (sourceExternal *SourceExternal) Update(dto SourceExternalDTO, origin string, uploader string) {
	sourceExternal.Comment = dto.Comment
	sourceExternal.URL = dto.URL
	sourceExternal.Origin = origin
	sourceExternal.Uploader = uploader
	sourceExternal.Updated = time.Now()
}
