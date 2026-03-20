// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package compare

// todo: write new test for multi diff compare

//
// import (
// 	"testing"
// 	"time"
//
// 	"mercedes-benz.ghe.com/foss/disuko/helper/reflection"
//
// 	"mercedes-benz.ghe.com/foss/disuko/domain"
// 	"mercedes-benz.ghe.com/foss/disuko/domain/compare"
// 	"mercedes-benz.ghe.com/foss/disuko/domain/license"
// 	"mercedes-benz.ghe.com/foss/disuko/domain/project"
// 	"mercedes-benz.ghe.com/foss/disuko/helper/s3Helper"
// 	"mercedes-benz.ghe.com/foss/disuko/logy"
// 	"github.com/stretchr/testify/assert"
// )
//
// var basePath = "./../../.."
//
// var requestSessionTest = &logy.RequestSession{ReqID: "TEST"}
//
// func Test_compareSpdxFiles(t *testing.T) {
// 	//load spdx1
// 	_, ci1 := LoadSpdxFile(t, "TestCompareOld.spdx")
// 	//load spdx2
// 	_, ci2 := LoadSpdxFile(t, "TestCompareNew.spdx")
//
// 	refs := license.LicenseRefs{
// 		"mit":          "MIT",
// 		"mit2":         "MIT2",
// 		"apache-2.0":   "Apache-2.0",
// 		"isc":          "ISC",
// 		"mit_changed":  "MIT_changed",
// 		"bsd-2-clause": "BSD-2-Clause",
// 		"bsd-3-clause": "BSD-3-Clause",
// 		"mit-ref3000":  "MIT",
// 	}
//
// 	ci1.EnrichComponentInfos(requestSessionTest)
// 	ci1.ApplyRefs(requestSessionTest, refs)
// 	ci2.EnrichComponentInfos(requestSessionTest)
// 	ci2.ApplyRefs(requestSessionTest, refs)
//
// 	evalRes1 := ci1.EvaluatePolicyRules([]*license.PolicyRules{})
// 	evalRes2 := ci2.EvaluatePolicyRules([]*license.PolicyRules{})
//
// 	//compare spdx1 und spdx2
// 	// var compareResult = MultiCompareSpdxFiles(evalRes1, evalRes2)
//
// 	assert.NotNil(t, compareResult)
// 	assert.Equal(t, len(compareResult), 40)
//
// 	unchanged := filterByteDiffType(compareResult, compare.UNCHANGED)
// 	newComps := filterByteDiffType(compareResult, compare.NEW)
// 	removed := filterByteDiffType(compareResult, compare.REMOVED)
// 	changed := filterByteDiffType(compareResult, compare.CHANGED)
//
// 	assert.Equal(t, len(unchanged), 37)
// 	assert.Equal(t, len(newComps), 0)
// 	assert.Equal(t, len(removed), 1)
// 	for _, changed := range changed {
// 		if changed.ComponentOld.SpdxId == "SPDXRef-Package-GoTestify-v130" {
// 			assert.Equal(t, changed.ComponentNew.LicenseDeclared, "MIT2")
// 			assert.Equal(t, changed.ComponentOld.LicenseDeclared, "MIT")
// 			assert.False(t, changed.CopyrightText)
// 			assert.False(t, changed.Version)
// 			assert.False(t, changed.Name)
// 			assert.False(t, changed.LicenseEffective)
// 			assert.False(t, changed.Description)
// 			assert.False(t, changed.DownloadLocation)
// 			assert.False(t, changed.Type)
// 			continue
// 		}
// 		if changed.ComponentOld.SpdxId == "SPDXRef-Package-GoTestify-v151" {
// 			assert.False(t, changed.Version)
// 			assert.True(t, changed.CopyrightText)
// 			assert.False(t, changed.Name)
// 			assert.True(t, changed.LicenseEffective)
// 			assert.True(t, changed.Description)
// 			assert.True(t, changed.DownloadLocation)
// 			assert.False(t, changed.Type)
// 			continue
// 		}
// 		if changed.ComponentOld.SpdxId == "SPDXRef-Package-pkgerrors-v091" {
// 			assert.True(t, changed.LicenseEffective)
// 			continue
// 		}
// 		if changed.ComponentOld.SpdxId == "SPDXRef-Package-twinj-uuid-v100" {
// 			assert.True(t, changed.License)
// 			assert.False(t, changed.LicenseEffective)
// 			continue
// 		}
// 		assert.Fail(t, "Unknown change in component SPDXID: "+changed.ComponentOld.SpdxId)
// 	}
// 	assert.Equal(t, len(changed), 2)
// }
//
// func filterByteDiffType(components []compare.ComponentDiffDto, diffType compare.ComponentDiffType) []compare.ComponentDiffDto {
// 	var result []compare.ComponentDiffDto
// 	for _, component := range components {
// 		if component.DiffType == diffType {
// 			result = append(result, component)
// 		}
// 	}
// 	return result
// }
//
// func LoadSpdxFile(t *testing.T, filename string) (*project.SpdxFileBase, *project.ComponentInfos) {
// 	spdxBytes := s3Helper.ReadFileFromLocalFileSystem(basePath + "/test_resources/" + filename)
// 	spdxString := s3Helper.ConvertToStringAndClose(spdxBytes)
//
// 	spdxFile := project.SpdxFileBase{
// 		ChildEntity:  domain.SetChildEntity("XXXX-XXXX-XXXX-XXX1"),
// 		Type:         0,
// 		ContentValid: true,
// 		SchemaValid:  false,
// 		Uploaded:     reflection.ToPointer(time.Now()),
// 	}
//
// 	spdxFile.ExtractMetaInfo(spdxString)
// 	ci := project.FileContent(spdxString).ExtractComponentInfo(requestSessionTest)
// 	s3Helper.MustNotNil(t, ci)
// 	s3Helper.MustTrue(t, len(ci) > 0)
//
// 	return &spdxFile, &ci
// }
//
