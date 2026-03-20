// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package pdocument

import (
	"fmt"

	"golang.org/x/text/language"
	"mercedes-benz.ghe.com/foss/disuko/domain"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
)

type PDocument struct {
	domain.ChildEntity `bson:",inline"`
	Description        string
	ApprovalId         string
	Type               PDocumentType
	Lang               string
	Hash               string
	VersionIndex       *DocumentVersion
}

type PDocumentType string

const (
	PD_DISCLOSURE_DOC PDocumentType = "disclosure"
	PD_POLICY_RULES   PDocumentType = "policies"
	PD_POLICY_CHECK   PDocumentType = "policycheck"
	PD_ARCHIVE        PDocumentType = "archive"
)

type DocumentVersion int

const (
	VERSION_1    DocumentVersion = 0
	VERSION_2    DocumentVersion = 1
	VERSION_3    DocumentVersion = 2
	VERSION_4    DocumentVersion = 3
	NONE_VERSION DocumentVersion = 4
)

func GetFileNameWithIndex(fileType PDocumentType, taskGuid string, lang *language.Tag, index int) string {
	targetFileName := ""
	switch fileType {
	case PD_DISCLOSURE_DOC:
		{
			if index == int(NONE_VERSION) {
				targetFileName = fmt.Sprintf("%s_disclosure", taskGuid)
			} else {
				targetFileName = fmt.Sprintf("%s_disclosure_%d", taskGuid, index)
			}
			if lang != nil {
				targetFileName += "-" + lang.String()
			}
			targetFileName += ".pdf"
		}
	case PD_POLICY_RULES:
		{
			targetFileName = taskGuid + "_policyrules.json"
		}
	case PD_POLICY_CHECK:
		{
			targetFileName = taskGuid + "_policycheck.json"
		}
	case PD_ARCHIVE:
		targetFileName = taskGuid + "_archive.zip"
	default:
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorUnknownFileType, fileType), "")
	}
	return targetFileName
}

func GetFileName(fileType PDocumentType, taskGuid string, lang *language.Tag) string {
	return GetFileNameWithIndex(fileType, taskGuid, lang, int(NONE_VERSION))
}

func LangStrToTag(langStr string) *language.Tag {
	var langTag *language.Tag = nil
	switch langStr {
	case "en":
		langTag = &language.English
	case "de":
		langTag = &language.German
	}
	return langTag
}
