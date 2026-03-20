// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/domain/project/components"
	"github.com/stretchr/testify/assert"
)

func TestSpdxFileBase_GetFileName(t *testing.T) {
	type fields struct {
		Created    time.Time
		SpdxKey    string
		ProjectKey string
		Version    string
	}
	type args struct {
		projectKey string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{name: "test", fields: fields{Created: time.Date(2019, 2, 11, 10, 03, 00, 00, time.UTC), ProjectKey: "p1", SpdxKey: "s1", Version: "v1"}, args: args{projectKey: "123"}, want: "/uploads/2019/02/p1/versions/v1/sbom/s1"},
		{name: "test", fields: fields{Created: time.Date(2019, 10, 11, 10, 03, 00, 00, time.UTC), ProjectKey: "p2", SpdxKey: "s2", Version: "v2"}, args: args{projectKey: "123"}, want: "/uploads/2019/10/p2/versions/v2/sbom/s2"},
		{name: "test", fields: fields{Created: time.Date(2000, 2, 11, 10, 03, 00, 00, time.UTC), ProjectKey: "p3", SpdxKey: "s3", Version: "v3"}, args: args{projectKey: "123"}, want: "/uploads/2000/02/p3/versions/v3/sbom/s3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Project{
				RootEntity: domain.SetRootEntity(tt.fields.ProjectKey),
			}
			p.Created = tt.fields.Created
			s := &SpdxFileBase{
				ChildEntity: domain.SetChildEntity(tt.fields.SpdxKey),
			}
			if got := strings.ReplaceAll(p.GetFilePathSbom(s.Key, tt.fields.Version), "\\", "/"); got != tt.want {
				t.Errorf("GetFilePathSbom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpdxFileBase_GetLicense(t *testing.T) {
	component := components.ComponentInfo{License: "A"}
	assert.Equal(t, component.License, component.GetLicenseEffective())

	component.LicenseDeclared = "B"
	assert.Equal(t, component.License, component.GetLicenseEffective())

	component.License = ""
	assert.Equal(t, component.LicenseDeclared, component.GetLicenseEffective())

	component.License = "NOASSERTION"
	assert.Equal(t, component.LicenseDeclared, component.GetLicenseEffective())
}

func TestSpdxFileBase_GetLicenseAppliedType(t *testing.T) {
	component := components.ComponentInfo{License: "A"}
	assert.Equal(t, components.LicenseConcluded, component.GetLicenseAppliedType())

	component.LicenseDeclared = "B"
	assert.Equal(t, components.LicenseConcluded, component.GetLicenseAppliedType())

	component.License = ""
	assert.Equal(t, components.LicenseDeclared, component.GetLicenseAppliedType())

	component.License = "NOASSERTION"
	assert.Equal(t, components.LicenseDeclared, component.GetLicenseAppliedType())
}

func TestComponentInfos_Stats(t *testing.T) {
	reader, err := os.Open("TestCompareNew.spdx")
	if err != nil {
		t.FailNow()
	}
	spdxBytes, err := io.ReadAll(reader)
	if err != nil {
		t.FailNow()
	}

	type args struct {
		rules []*license.PolicyRules
		refs  license.LicenseRefs
	}
	fileContent := FileContent(spdxBytes)
	tests := []struct {
		name string
		cis  components.ComponentInfos
		args args
		want components.ComponentStats
	}{
		{name: "denied", cis: fileContent.ExtractComponentInfo(nil), args: args{rules: []*license.PolicyRules{{Name: "test", ComponentsDeny: []string{"Apache-2.0"}}}, refs: license.LicenseRefs{"apache-2.0-lala": license.LicenseRef{ID: "Apache-2.0"}}}, want: components.ComponentStats{
			Total:       63,
			Allowed:     0,
			Warned:      0,
			Denied:      1,
			Questioned:  0,
			NoAssertion: 61,
		}},
		{name: "allowed", cis: fileContent.ExtractComponentInfo(nil), args: args{rules: []*license.PolicyRules{{Name: "test", ComponentsAllow: []string{"Apache-2.0"}}}, refs: license.LicenseRefs{"apache-2.0-lala": license.LicenseRef{ID: "Apache-2.0"}}}, want: components.ComponentStats{
			Total:       63,
			Allowed:     1,
			Warned:      0,
			Denied:      0,
			Questioned:  0,
			NoAssertion: 61,
		}},
		{name: "warned", cis: fileContent.ExtractComponentInfo(nil), args: args{rules: []*license.PolicyRules{{Name: "test", ComponentsWarn: []string{"Apache-2.0"}}}, refs: license.LicenseRefs{"apache-2.0-lala": license.LicenseRef{ID: "Apache-2.0"}}}, want: components.ComponentStats{
			Total:       63,
			Allowed:     0,
			Warned:      1,
			Denied:      0,
			Questioned:  0,
			NoAssertion: 61,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cis.EnrichComponentInfos(nil)
			tt.cis.ApplyRefs(tt.args.refs)
			assert.Equalf(t, tt.want, tt.cis.EvaluatePolicyRules(tt.args.rules, nil, false, nil, "").Stats, "Stats(%v, %v)", tt.args.rules, tt.args.refs)
		})
	}
}
