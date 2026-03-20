// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package search

import (
	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/license"
	obligation2 "github.com/eclipse-disuko/disuko/domain/obligation"
	sorthelper "github.com/eclipse-disuko/disuko/helper/sort"
	"testing"
	"time"
)

func Test_Sort(t *testing.T) {
	tests := []struct {
		name     string
		sortBy   []string
		sortDesc []bool
		input    license.LicensesResponse
		expected []string
	}{
		{
			name:     "Sort by name ascending",
			sortBy:   []string{"name"},
			sortDesc: []bool{false},
			input: license.LicensesResponse{
				Licenses: []*license.LicenseSlimDto{
					{Name: "B License"},
					{Name: "A License"},
					{Name: "C License"},
				},
			},
			expected: []string{"A License", "B License", "C License"},
		},
		{
			name:     "Sort by licenseId descending",
			sortBy:   []string{"licenseId"},
			sortDesc: []bool{true},
			input: license.LicensesResponse{
				Licenses: []*license.LicenseSlimDto{
					{LicenseId: "ID2"},
					{LicenseId: "ID1"},
					{LicenseId: "ID3"},
				},
			},
			expected: []string{"ID3", "ID2", "ID1"},
		},
		{
			name:     "Sort by meta.approvalState ascending",
			sortBy:   []string{"meta.approvalState"},
			sortDesc: []bool{false},
			input: license.LicensesResponse{
				Licenses: []*license.LicenseSlimDto{
					{Meta: license.MetaDataSlimDto{ApprovalState: license.Approved}},
					{Meta: license.MetaDataSlimDto{ApprovalState: license.Pending}},
					{Meta: license.MetaDataSlimDto{ApprovalState: license.Deprecated}},
				},
			},
			expected: []string{"approved", "deprecated", "pending"},
		},
		{
			name:     "Sort by updated ascending",
			sortBy:   []string{"updated"},
			sortDesc: []bool{false},
			input: license.LicensesResponse{
				Licenses: []*license.LicenseSlimDto{
					{BaseDto: domain.BaseDto{Updated: time.Date(2024, 8, 3, 10, 0, 3, 0, time.Local)}},
					{BaseDto: domain.BaseDto{Updated: time.Date(2024, 8, 3, 10, 0, 1, 0, time.Local)}},
					{BaseDto: domain.BaseDto{Updated: time.Date(2024, 8, 3, 10, 0, 2, 0, time.Local)}},
				},
			},
			expected: []string{"2024-08-03 10:00:01", "2024-08-03 10:00:02", "2024-08-03 10:00:03"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := &RequestSearchOptions{
				SortBy:   tt.sortBy,
				SortDesc: tt.sortDesc,
			}

			asc := options.IsSortAsc()
			key := options.GetSortKey()
			if key == "name" {
				sorthelper.Sort(tt.input.Licenses, func(dto *license.LicenseSlimDto) string { return dto.Name }, sorthelper.StringLessThan, asc)
			} else if key == "licenseId" {
				sorthelper.Sort(tt.input.Licenses, func(dto *license.LicenseSlimDto) string { return dto.LicenseId }, sorthelper.StringLessThan, asc)
			} else if key == "meta.approvalState" {
				sorthelper.Sort(tt.input.Licenses, func(dto *license.LicenseSlimDto) string { return dto.Meta.ApprovalState.Value() }, sorthelper.StringLessThan, asc)
			} else if key == "meta.family" {
				sorthelper.Sort(tt.input.Licenses, func(dto *license.LicenseSlimDto) string { return dto.Meta.Family.Value() }, sorthelper.StringLessThan, asc)
			} else if key == "meta.licenseType" {
				sorthelper.Sort(tt.input.Licenses, func(dto *license.LicenseSlimDto) string { return dto.Meta.LicenseType.Value() }, sorthelper.StringLessThan, asc)
			} else if key == "source" {
				sorthelper.Sort(tt.input.Licenses, func(dto *license.LicenseSlimDto) string { return string(dto.Source) }, sorthelper.StringLessThan, asc)
			} else if key == "updated" {
				sorthelper.Sort(tt.input.Licenses, func(dto *license.LicenseSlimDto) int64 { return dto.Updated.Unix() }, sorthelper.Int64LessThan, asc)
			} else if key == "meta.classifications" {
				sorthelper.Sort(tt.input.Licenses, func(dto *license.LicenseSlimDto) int64 {
					return obligation2.GetLevelWeight(dto.Meta.PrevalentClassificationLevel)
				}, sorthelper.Int64LessThan, asc)
			}

			for i, expectedName := range tt.expected {
				if resultField := getSortFieldBy(tt.input.Licenses[i], tt.sortBy[0]); resultField != expectedName {
					t.Errorf("expected %s at index %d, got %s", expectedName, i, resultField)
				}
			}
			optionsNew := &RequestSearchOptionsNew{
				SortBy: []SortBy{{
					Key:   tt.sortBy[0],
					Order: sortOrder(tt.sortDesc[0]),
				}},
			}

			asc = optionsNew.IsSortAsc()
			key = optionsNew.GetSortKey()
			if key == "name" {
				sorthelper.Sort(tt.input.Licenses, func(dto *license.LicenseSlimDto) string { return dto.Name }, sorthelper.StringLessThan, asc)
			} else if key == "licenseId" {
				sorthelper.Sort(tt.input.Licenses, func(dto *license.LicenseSlimDto) string { return dto.LicenseId }, sorthelper.StringLessThan, asc)
			} else if key == "meta.approvalState" {
				sorthelper.Sort(tt.input.Licenses, func(dto *license.LicenseSlimDto) string { return dto.Meta.ApprovalState.Value() }, sorthelper.StringLessThan, asc)
			} else if key == "meta.family" {
				sorthelper.Sort(tt.input.Licenses, func(dto *license.LicenseSlimDto) string { return dto.Meta.Family.Value() }, sorthelper.StringLessThan, asc)
			} else if key == "meta.licenseType" {
				sorthelper.Sort(tt.input.Licenses, func(dto *license.LicenseSlimDto) string { return dto.Meta.LicenseType.Value() }, sorthelper.StringLessThan, asc)
			} else if key == "source" {
				sorthelper.Sort(tt.input.Licenses, func(dto *license.LicenseSlimDto) string { return string(dto.Source) }, sorthelper.StringLessThan, asc)
			} else if key == "updated" {
				sorthelper.Sort(tt.input.Licenses, func(dto *license.LicenseSlimDto) int64 { return dto.Updated.Unix() }, sorthelper.Int64LessThan, asc)
			} else if key == "meta.classifications" {
				sorthelper.Sort(tt.input.Licenses, func(dto *license.LicenseSlimDto) int64 {
					return obligation2.GetLevelWeight(dto.Meta.PrevalentClassificationLevel)
				}, sorthelper.Int64LessThan, asc)
			}

			for i, expectedName := range tt.expected {
				if resultField := getSortFieldBy(tt.input.Licenses[i], tt.sortBy[0]); resultField != expectedName {
					t.Errorf("expected %s at index %d, got %s", expectedName, i, resultField)
				}
			}
		})
	}
}

func sortOrder(desc bool) string {
	if desc {
		return "desc"
	} else {
		return "asc"
	}
}

func getSortFieldBy(dto *license.LicenseSlimDto, sortBy string) string {
	switch sortBy {
	case "name":
		return dto.Name
	case "licenseId":
		return dto.LicenseId
	case "meta.approvalState":
		return dto.Meta.ApprovalState.Value()
	case "updated":
		return dto.Updated.Format("2006-01-02 15:04:05")
	default:
		return ""
	}
}
