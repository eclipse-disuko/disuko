// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"github.com/eclipse-disuko/disuko/domain/overallreview"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/audit"
)

type SourcesExternal []*SourceExternal

type ProjectVersion struct {
	domain.ChildEntity `bson:",inline"`
	audit.Container    `bson:",inline"`
	Name               string
	Description        string
	SourceExternal     SourcesExternal
	InternalVersion    int
	Status             ProjectVersionStatusType `json:"status"`
	// Created            time.Time
	// Updated            time.Time
	Deleted        bool
	OverallReviews []overallreview.OverallReview
	Subscriber     []string
}

func (s SourcesExternal) GetUrls() []string {
	externalUrls := make([]string, 0)
	for _, externalUrl := range s {
		externalUrls = append(externalUrls, externalUrl.URL)
	}
	return externalUrls
}

func (version *ProjectVersion) GetSourceExternalAll() []*SourceExternal {
	if version.SourceExternal == nil {
		version.SourceExternal = make([]*SourceExternal, 0)
	}
	return version.SourceExternal
}

func (version *ProjectVersion) GetSourceExternal(index int) *SourceExternal {
	return version.SourceExternal[index]
}

func (version *ProjectVersion) GetStatus() ProjectVersionStatusType {
	if len(version.Status) <= 0 {
		return PV_New
	}
	return version.Status
}

func (version *ProjectVersion) UpdateStatusWhenUploadIsDone() {
	if len(version.Status) == 0 || version.Status == PV_New {
		version.Status = PV_Unreviewed
	}
}
