// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package compare

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/project/components"
)

type ComponentDiffType string

const (
	UNCHANGED = "UNCHANGED"
	NEW       = "NEW"
	REMOVED   = "REMOVED"
	CHANGED   = "CHANGED"
)

type ComponentChangesDto struct {
	SpdxId           bool
	Name             bool
	Version          bool
	License          bool
	LicenseDeclared  bool
	LicenseEffective bool
	LicenseComments  bool
	CopyrightText    bool
	Description      bool
	DownloadLocation bool
	Type             bool
	Modified         bool
	Questioned       bool
	Unasserted       bool
	PURL             bool
	PrStatus         bool `json:"prStatus"`
}

type ComponentMultiDiffDto struct {
	DiffType      ComponentDiffType
	Name          string
	Changes       map[string]ComponentChangesDto
	ComponentsOld []*components.ComponentInfoDto
	ComponentsNew []*components.ComponentInfoDto
}

func NewComponentMultiDiffDto(diffType ComponentDiffType, componentName string) *ComponentMultiDiffDto {
	return &ComponentMultiDiffDto{
		DiffType:      diffType,
		Name:          componentName,
		Changes:       make(map[string]ComponentChangesDto, 0),
		ComponentsOld: make([]*components.ComponentInfoDto, 0),
		ComponentsNew: make([]*components.ComponentInfoDto, 0),
	}
}
