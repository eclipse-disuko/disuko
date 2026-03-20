// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package sbomlist

import (
	"sort"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/project"
)

type Histories []*project.SpdxFileBase

type SbomList struct {
	domain.RootEntity `bson:",inline"`
	domain.SoftDelete `bson:",inline"`

	SpdxFileHistory Histories
}

func (h Histories) GetLatest() *project.SpdxFileBase {
	sort.Slice(h, func(i, j int) bool {
		return h[i].Updated.After(h[j].Updated)
	})

	return h[0]
}

func (h Histories) GetByKey(key string) *project.SpdxFileBase {
	for _, s := range h {
		if s.Key == key {
			return s
		}
	}
	return nil
}
