// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"mercedes-benz.ghe.com/foss/disuko/domain"
	"mercedes-benz.ghe.com/foss/disuko/domain/project"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
)

func TestParsing(t *testing.T) {
	expiry, err := time.Parse(time.RFC3339, "2022-08-01"+"T00:00:00.000Z")
	if err != nil {
		t.Errorf("Should not be an error")
	}
	t.Logf("expiry %v", expiry)

	expiry2, err := time.Parse(time.RFC3339, "2022-09-10T13:40:39.504Z")
	if err != nil {
		t.Errorf("Should not be an error")
	}
	t.Logf("expiry %v", expiry2)

}

func TestProjectModel(t *testing.T) {
	var key = "4711"
	var name = key + "Name"
	var description = key + "Description"
	var versionName = "v1.1"
	var p *project.Project
	p = &project.Project{
		RootEntity:  domain.SetRootEntity(key),
		Name:        name,
		Versions:    map[string]*project.ProjectVersion{},
		Description: description,
		FreeLabels:  make([]string, 0),
	}
	if &p == nil {
		t.Errorf("Cannot create project reference (project.Project)")
	}
	if p.Key != key {
		t.Errorf("Missing correct definition of (key)")
	}
	if p.Name != name {
		t.Errorf("Missing correct definition of (name)")
	}
	if p.Description != description {
		t.Errorf("Missing correct definition of (description)")
	}

	if len(p.GetVersions()) != 0 {
		t.Errorf("Versions should be empty")
	}
	p.CreateNewProjectVersionIfNameNotUsed(versionName, "")
	if len(p.GetVersions()) != 1 {
		t.Errorf("Versions should be size 1")
	}
	if p.FindVersionByName(versionName) == nil {
		t.Error("Version with name " + versionName + " should exists")
	}

	var vArray = p.GetVersions()
	var v = vArray[0]

	if len(vArray) != 1 {
		t.Error("Version result size should be 1 but is " + fmt.Sprint(len(vArray)))
	}

	if v.Name != versionName {
		t.Error("First Version name is wrong should be '" + versionName + "' but is '" + v.Name + "'")
	}

	var labelName = "Hello Word"
	p.FreeLabels = append(p.FreeLabels, labelName)
	if len(p.FreeLabels) != 1 {
		t.Error("FreeLabels result size should be 1 but is " + fmt.Sprint(len(p.FreeLabels)))
	}
	if time.Now().Nanosecond() < p.Created.Nanosecond() {
		t.Error("There is something wrong with the created timestamp")
	}

}

func Test_assertToken(t *testing.T) {
	tests := []struct {
		name       string
		discoToken string
		want       string
		error      bool
	}{
		{name: "valid", discoToken: "DISCO 185781c8-5bec-4763-b627-c3a38b00a986", want: "185781c8-5bec-4763-b627-c3a38b00a986", error: false},
		{name: "invalid uuid", discoToken: "DISCO 185781c8-5bec-4763-b627-c3a38b00a986123", want: "185781c8-5bec-4763-b627-c3a38b00a986", error: true},
		{name: "no prefix ", discoToken: "lala", want: "lala", error: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ex := false
			res := ""
			exception.TryCatch(func() {
				res = assertTokenUUID(tt.discoToken, DiscoBearer)
			}, func(exception exception.Exception) {
				ex = true
			})
			assert.Equal(t, ex, tt.error)
			if !tt.error {
				assert.Equal(t, res, tt.want)
			}
		})
	}
}
