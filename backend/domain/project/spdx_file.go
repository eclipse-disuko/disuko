// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/eclipse-disuko/disuko/domain/licenserules"
	"github.com/eclipse-disuko/disuko/domain/overallreview"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/domain/project/components"
	"github.com/eclipse-disuko/disuko/domain/schema"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/tidwall/gjson"
	"github.com/xeipuuv/gojsonschema"
)

type SpdxFileBase struct {
	domain.ChildEntity `bson:",inline"`
	Hash               string // sha256
	SchemaId           string
	SchemaName         string
	ContentValid       bool
	SchemaValid        bool
	Type               schema.SchemaType
	MetaInfo           MetaInfo
	ApprovalInfo       ApprovalInfo
	ValidationErrors   string
	Uploaded           *time.Time
	Origin             string
	Uploader           string
	Tag                string

	OverallReview *overallreview.OverallReview

	IsInUse  bool // store in DB
	IsLocked bool // store in DB

	IsToDelete bool // do not store in DB, only for Frontend, based on conditions
	IsToRetain bool // do not store in DB, only for Frontend, based on conditions
}

type ApprovalInfo struct {
	IsInApproval bool
	Comment      string
	ApprovalGuid string
	Status       string
}

type MetaInfo struct {
	domain.ChildEntity `bson:",inline"`
	Name               string
	SpdxId             string
	SpdxVersion        string
	DataLicense        string
	Comment            string
	Creators           []string
	CreationData       string
	HasExternalRefs    bool
}

type DetailedLicense[T any] struct {
	License        T
	OrigName       string
	ReferencedName string
}

type Detail struct {
	Key   string
	Value string
}

type ComponentDetails struct {
	UnassertedLicenseText bool

	RawInfo    map[string]json.RawMessage
	Attributes []Detail

	ExtractedLicenses  []ExtractedLicense
	IdentifiedViaAlias []IdentifiedLicense
	UnknownLicenses    []string
	KnownLicenses      []DetailedLicense[*license.LicenseDto]

	Problems []string

	CanChooseLicense   bool
	ChoiceDeniedReason string
	LicenseRuleApplied *licenserules.LicenseRuleSlimDto
	ContainsOr         bool
}

type ComponentLicenses struct {
	UnknownLicenses []string
	KnownLicenses   []DetailedLicense[*license.LicenseNameIdDto]
}

type ExtractedLicense struct {
	LicenseId          string
	ExtractedText      string
	Comment            string
	Name               string
	SeeAlsos           []string
	ExternalDocumentId string
	SpdxDocument       string
}

type IdentifiedLicense struct {
	License       ExtractedLicense
	AliasTargetId string
}

func (spdx *SpdxFileBase) ValidateSpdxContent(requestSession *logy.RequestSession, spdxString string, spdxSchema *schema.SpdxSchema) (bool, error) {
	schemaLoader := gojsonschema.NewStringLoader(spdxSchema.Content)
	documentLoader := gojsonschema.NewStringLoader(spdxString)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return false, err
	}
	if result.Valid() {
		logy.Infof(requestSession, "Sbom with key %s has been validated against schema %s", spdx.Key, spdxSchema.Key)
		return true, nil
	}
	logy.Warnf(requestSession, "The document is not valid. see errors :\n")
	errorMsgs := ""
	for _, desc := range result.Errors() {
		errorMsgs += "- " + desc.String() + "\n"
	}
	return false, errors.New(errorMsgs)
}

type FileContent string

func (f FileContent) ExtractComponentInfo(requestSession *logy.RequestSession) components.ComponentInfos {
	logy.Infof(requestSession, "project::ExtractComponentInfo")
	componentInfo := make([]components.ComponentInfo, 0)

	jsonDocument := gjson.Parse(string(f))
	jsonPackages := jsonDocument.Get("packages")
	jsonPackages.ForEach(ExtractComponentInfo(&componentInfo, components.PACKAGE))

	jsonFiles := jsonDocument.Get("files")
	jsonFiles.ForEach(ExtractComponentInfo(&componentInfo, components.FILE))
	jsonSnippets := jsonDocument.Get("snippets")
	jsonSnippets.ForEach(ExtractComponentInfo(&componentInfo, components.SNIPPET))

	jsonDescribes := jsonDocument.Get("documentDescribes")
	jsonDescribes.ForEach(func(key, value gjson.Result) bool {
		for i := 0; i < len(componentInfo); i++ {
			if componentInfo[i].SpdxId == value.String() {
				componentInfo[i].Type = components.ROOT
			}
		}
		return true
	})

	jsonRelationships := jsonDocument.Get("relationships")
	jsonRelationships.ForEach(func(key, value gjson.Result) bool {
		if value.Get("relationshipType").String() == "FILE_MODIFIED" {
			for i := 0; i < len(componentInfo); i++ {
				if componentInfo[i].SpdxId == value.Get("spdxElementId").String() {
					componentInfo[i].Modified = true
				}
			}
		}
		if value.Get("relationshipType").String() == "DESCRIBES" && value.Get("spdxElementId").String() == "SPDXRef-DOCUMENT" {
			for i := 0; i < len(componentInfo); i++ {
				if componentInfo[i].SpdxId == value.Get("relatedSpdxElement").String() {
					componentInfo[i].Type = components.ROOT
				}
			}
		}
		return true // keep iterating
	})
	logy.Infof(requestSession, "project::ExtractComponentInfo done got Packages: %d Files: %d, Snippets: %d, Relations: %d", len(jsonPackages.Array()), len(jsonFiles.Array()), len(jsonSnippets.Array()), len(jsonRelationships.Array()))
	return componentInfo
}

func ExtractComponentInfo(componentInfo *[]components.ComponentInfo, componentType components.ComponentType) func(key gjson.Result, value gjson.Result) bool {
	return func(key, value gjson.Result) bool {
		anno := false
		if value.Get("annotations").IsArray() && len(value.Get("annotations").Array()) > 0 {
			anno = true
		}

		name := value.Get("name").String()
		version := value.Get("versionInfo").String()
		if componentType == components.FILE {
			name = value.Get("fileName").String()
			if value.Get("checksums").IsArray() && len(value.Get("checksums").Array()) > 0 {
				version = value.Get("checksums.0.algorithm").String() + "/" + value.Get("checksums.0.checksumValue").String()[:6]
			}
		}

		searchPaths := []string{
			"externalRefs.#(referenceCategory==\"PACKAGE_MANAGER\")",
			"externalRefs.#(referenceCategory==\"PACKAGE-MANAGER\")",
			"externalRefs.#(referenceType==\"purl\")",
		}
		purl := ""
		for _, path := range searchPaths {
			ref := value.Get(path)
			if ref.IsArray() && len(ref.Array()) > 0 {
				purl = ref.Array()[0].Get("referenceLocator").String()
				break
			} else if len(ref.String()) > 0 {
				purl = ref.Get("referenceLocator").String()
				break
			}
		}

		*componentInfo = append(*componentInfo, components.ComponentInfo{
			SpdxId:            value.Get("SPDXID").String(),
			Name:              name,
			License:           value.Get("licenseConcluded").String(),
			LicenseDeclared:   value.Get("licenseDeclared").String(),
			LicenseComments:   value.Get("licenseComments").String(),
			Version:           version,
			CopyrightText:     value.Get("copyrightText").String(),
			Description:       value.Get("description").String(),
			DownloadLocation:  value.Get("downloadLocation").String(),
			PURL:              purl,
			HomepageURL:       value.Get("packageHomepage").String(),
			HasAnnotations:    anno,
			ContainsAuxiliary: false,
			Type:              componentType,
			Modified:          false,
		})
		return true // keep iterating
	}
}

func (spdxFile *SpdxFileBase) ExtractMetaInfo(spdxString string) {
	value := gjson.GetMany(spdxString, "name", "SPDXID", "spdxVersion", "dataLicense",
		"comment", "creationInfo.creators", "creationInfo.created", "externalDocumentRefs")

	var creatorStringArray []string
	creatorResultArray := value[5].Array()

	for _, v := range creatorResultArray {
		creatorStringArray = append(creatorStringArray, v.String())
	}

	refs := false
	if value[7].IsArray() && len(value[7].Array()) > 0 {
		refs = true
	}

	spdxFile.MetaInfo = MetaInfo{
		Name:            value[0].String(),
		SpdxId:          value[1].String(),
		SpdxVersion:     value[2].String(),
		DataLicense:     value[3].String(),
		Comment:         value[4].String(),
		Creators:        creatorStringArray,
		CreationData:    value[6].String(),
		HasExternalRefs: refs,
	}
}
