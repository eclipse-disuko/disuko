// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"path/filepath"
	"testing"

	"github.com/eclipse-disuko/disuko/logy"
	"github.com/xeipuuv/gojsonschema"
)

var basePath = "./../.."

// var basePath = "."

/*
*
A list of rules you will find here
https://github.com/json-schema-org/JSON-Schema-Test-Suite/tree/master/tests/draft7
*/
func TestJsonSchema(t *testing.T) {
	logy.Infof(requestSessionTest, "Test 01 JSON Schema with no rules against correct SPDX File\n\n")
	// Test 01 JSON Schema with no rules against correct SPDX File
	if !validateSPDXAgainstSchema("spdxSchemaNoRules.json", "backendSPDX.json") {
		t.Fail()
	}
	logy.Infof(requestSessionTest, "Test 02 JSON Schema with no rules against wrong SPDX File\n\n")
	// Test 02 JSON Schema with no rules against wrong SPDX File
	if !validateSPDXAgainstSchema("spdxSchemaNoRules.json", "backendSPDXWithoutRoot.json") {
		t.Fail()
	}

	logy.Infof(requestSessionTest, "JSON Schema with rules against wrong SPDX File\n\n")
	// Test 03 JSON Schema with rules against wrong SPDX File
	// In schema this rule is set in row 6
	// "additionalProperties": false,
	// In schema this rule is set in row 611
	// "required": ["Document"]
	if validateSPDXAgainstSchema("spdxSchemaWithRuleRoot.json", "backendSPDXWithoutRoot.json") {
		t.Fail()
	}

	logy.Infof(requestSessionTest, "JSON Schema with rules against wrong SPDX File\n\n")
	// Test 04 JSON Schema with rules (check spdx version is correct) against correct SPDX File with wrong version
	// "spdxVersion": "SPDX-2.1", should be "spdxVersion": "SPDX-2.2",
	// In schema this rule is set in row 72
	// "enum": ["SPDX-2.2"]
	if validateSPDXAgainstSchema("spdxSchemaWithRuleVersion.json", "backendSPDXWrongVersion.json") {
		t.Fail()
	}
	// Test 05 JSON Schema with all rules against correct SPDX File with missing one versionInfo in package
	// Here are the rules
	// 6: "additionalProperties": false,
	// 72: "enum": ["SPDX-2.2"]
	// 345: "required": ["versionInfo"]
	// 613: "required": ["Document"]
	if validateSPDXAgainstSchema("spdxSchema.json", "backendSPDXWithoutOnVersionInfo.json") {
		t.Fail()
	}
	// Test 06 JSON Schema with all rules against correct SPDX File
	//if (!validateSPDXAgainstSchema("spdxSchema.json", "backendSPDX.json")) {
	//	t.Fail()
	//}
}

func validateSPDXAgainstSchema(schemaFile string, spdxFile string) bool {
	sPath, _ := filepath.Abs(basePath + "/test_resources/" + schemaFile)
	sPath = "file:///" + filepath.ToSlash(sPath)

	dPath, _ := filepath.Abs(basePath + "/test_resources/" + spdxFile)
	dPath = "file:///" + filepath.ToSlash(dPath)

	// or use gojsonschema.NewStringLoader() when the content is in string
	schemaLoader := gojsonschema.NewReferenceLoader(sPath)
	documentLoader := gojsonschema.NewReferenceLoader(dPath)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		logy.Infof(requestSessionTest, "The document is valid\n")
		return true
	} else {
		logy.Errorf(requestSessionTest, "The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			logy.Errorf(requestSessionTest, "- %s\n", desc)
		}
		return false
	}
}

func TestSimpleJsonSchema(t *testing.T) {
	logy.Infof(requestSessionTest, "Test simpleSchema success\n\n")
	if !validateSPDXAgainstSchema("simpleSchema.json", "simpleSPDXSuccess.json") {
		t.Fail()
	}
	logy.Errorf(requestSessionTest, "Test simpleSchema failed\n\n")
	if validateSPDXAgainstSchema("simpleSchema.json", "simpleSPDXFailed.json") {
		t.Fail()
	}
}
