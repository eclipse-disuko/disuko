// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"testing"

	"mercedes-benz.ghe.com/foss/disuko/domain/project"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/helper/test"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		err      bool
		errorStr string
	}{
		{name: "ProjectRequestDto valid", data: project.ProjectRequestDto{
			Name:         "lala",
			SchemaLabel:  "lala",
			PolicyLabels: []string{"lala", "lala2"},
			FreeLabels:   []string{"lala", "lala2"},
			Description:  "desc",
			Owner:        "owner",
		}, err: false, errorStr: ""},
		{name: "ProjectRequestDto empty policy labels valid", data: project.ProjectRequestDto{
			Name:         "lala",
			SchemaLabel:  "lala",
			PolicyLabels: []string{},
			FreeLabels:   []string{"lala", "lala2"},
			Description:  "desc",
			Owner:        "owner",
		}, err: false, errorStr: ""},
		{name: "ProjectRequestDto empty free labels valid", data: project.ProjectRequestDto{
			Name:         "lala",
			SchemaLabel:  "lala",
			PolicyLabels: []string{"lala", "lala2"},
			FreeLabels:   []string{},
			Description:  "desc",
			Owner:        "owner",
		}, err: false, errorStr: ""},
		{name: "ProjectRequestDto name invalid", data: project.ProjectRequestDto{
			Name:         "la",
			SchemaLabel:  "lala",
			PolicyLabels: []string{"lala", "lala2"},
			FreeLabels:   []string{"lala", "lala2"},
			Description:  "desc",
			Owner:        "owner",
		}, err: true, errorStr: "[Field: ProjectRequestDto.Name / Validator: Greater than or equal / Param: 3] "},
		{name: "ProjectRequestDto SchemaLabel invalid", data: project.ProjectRequestDto{
			Name:         "lala",
			SchemaLabel:  "la",
			PolicyLabels: []string{"lala", "lala2"},
			FreeLabels:   []string{"lala", "lala2"},
			Description:  "desc",
			Owner:        "owner",
		}, err: true, errorStr: "[Field: ProjectRequestDto.SchemaLabel / Validator: Greater than or equal / Param: 3] "},
		{name: "ProjectRequestDto PolicyLabels invalid", data: project.ProjectRequestDto{
			Name:         "lala",
			SchemaLabel:  "lala",
			PolicyLabels: []string{"la", "lala2"},
			FreeLabels:   []string{"lala", "lala2"},
			Description:  "desc",
			Owner:        "owner",
		}, err: true, errorStr: "[Field: ProjectRequestDto.PolicyLabels[0] / Validator: Greater than or equal / Param: 3] "},
		{name: "ProjectRequestDto FreeLabels invalid", data: project.ProjectRequestDto{
			Name:         "lala",
			SchemaLabel:  "lala",
			PolicyLabels: []string{"lala", "lala2"},
			FreeLabels:   []string{"", "lala2"},
			Description:  "desc",
			Owner:        "owner",
		}, err: true, errorStr: "[Field: ProjectRequestDto.FreeLabels[0] / Validator: Greater than or equal / Param: 3] "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err {
				test.ExpectException(t, message.ErrorJsonValidatingForm, func() {
					Validate(tt.data, false, nil)
				})
				test.ExpectException(t, message.ErrorJsonValidatingForm, func() {
					Validate(tt.data, true, nil)
				})
			} else {
				// no exception/error expected
				Validate(tt.data, true, nil)
				Validate(tt.data, false, nil)
			}
		})
	}
}
