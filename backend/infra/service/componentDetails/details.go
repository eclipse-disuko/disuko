// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"encoding/json"
	"strings"
	"time"

	"mercedes-benz.ghe.com/foss/disuko/domain/license"

	"github.com/tidwall/gjson"
	"mercedes-benz.ghe.com/foss/disuko/domain/project"
	"mercedes-benz.ghe.com/foss/disuko/domain/project/components"
	"mercedes-benz.ghe.com/foss/disuko/helper"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	license2 "mercedes-benz.ghe.com/foss/disuko/infra/repository/license"
	projectService "mercedes-benz.ghe.com/foss/disuko/infra/service/project"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

func GetComponentDetails(requestSession *logy.RequestSession, holder projectService.RepositoryHolder, projectKey string, spdxId string, spdxContent string, spdxUpload *time.Time, spdxKey string, isResponsible bool) project.ComponentDetails {
	spdxInfo, compType := findCompInFile(spdxContent, spdxId)
	compInfo := buildCompInfoFromJSON(requestSession, holder, projectKey, spdxInfo, compType, spdxUpload, spdxKey)

	var rawDetails map[string]json.RawMessage
	err := json.Unmarshal([]byte(spdxInfo.String()), &rawDetails)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.UnmarshallingSpdxContent))

	containsOr := compInfo.GetLicensesEffective().Op == components.OR

	attributes := []project.Detail{
		{
			Key:   "name",
			Value: compInfo.Name,
		}, {
			Key:   "versionInfo",
			Value: compInfo.Version,
		}, {
			Key:   "purl",
			Value: compInfo.PURL,
		}, {
			Key:   "licenseEffective",
			Value: compInfo.EffectiveLicensesString(),
		},
	}
	attributes = append(attributes, getAttributes(rawDetails)...)

	// In case json does not contain the path "hasExtractedLicensingInfos" the expression value[0].String() is resolved to an empty string ("") letting variable licenses not initialized, if just declared
	// Therefore initialization as an empty slice.
	extractedLicenses := make([]project.ExtractedLicense, 0)
	rawLicenses := gjson.GetMany(spdxContent, "hasExtractedLicensingInfos")
	if rawLicenses[0].String() != "" {
		err = json.Unmarshal([]byte(rawLicenses[0].String()), &extractedLicenses)
		exception.HandleErrorServerMessage(err, message.GetI18N(message.UnmarshallingLicenseContent))
	}

	knownLicenses, unknownLicenses := ProcessComponentLicenses[*license.LicenseDto](requestSession, &compInfo, holder.LicenseRepository, (*license.License).ToDtoWithoutObligations)

	remainingExtractedLicenses := make([]project.ExtractedLicense, 0)
	identifiedViaAliasLicenses := make([]project.IdentifiedLicense, 0)
	licenseRefs := holder.LicenseRepository.GetLicenseRefs(requestSession)
	for _, el := range extractedLicenses {
		if licenseByAlias, ok := licenseRefs[strings.ToLower(el.LicenseId)]; ok {
			identifiedViaAliasLicenses = append(identifiedViaAliasLicenses, project.IdentifiedLicense{License: el, AliasTargetId: licenseByAlias.ID})
		} else {
			remainingExtractedLicenses = append(remainingExtractedLicenses, el)
		}
	}

	problems := make([]string, 0)
	if compInfo.ContainsBadChars {
		problems = append(problems, message.ProblemBadchars)
	}
	if compInfo.ComplexExpression {
		problems = append(problems, message.ProblemComplex)
	}
	var (
		deniedReason string
		canChoose    = compInfo.LicenseRuleApplied == nil && compInfo.GetLicensesEffective().Op == components.OR
	)
	if canChoose {
		if !isResponsible {
			deniedReason = message.ChoiceDeniedResp
		} else if len(compInfo.GetLicensesEffective().List) > 4 {
			deniedReason = message.ChoiceDeniedMassive
		}
	}

	if compInfo.Version == "" && deniedReason == "" {
		deniedReason = message.DecisionDeniedComponentVersionNotSet
	}

	return project.ComponentDetails{
		UnassertedLicenseText: helper.IsUnasserted(compInfo.GetLicenseEffective()),
		RawInfo:               rawDetails,
		Attributes:            attributes,

		ExtractedLicenses:  remainingExtractedLicenses,
		IdentifiedViaAlias: identifiedViaAliasLicenses,
		UnknownLicenses:    unknownLicenses,
		KnownLicenses:      knownLicenses,

		Problems: problems,

		CanChooseLicense:   canChoose,
		ChoiceDeniedReason: deniedReason,
		LicenseRuleApplied: compInfo.LicenseRuleApplied.ToSlimDto(),
		ContainsOr:         containsOr,
	}
}

func GetComponentLicenses(requestSession *logy.RequestSession, holder projectService.RepositoryHolder, projectKey string, spdxId string, spdxContent string, spdxUpload *time.Time, spdxKey string) project.ComponentLicenses {
	spdxInfo, compType := findCompInFile(spdxContent, spdxId)
	compInfo := buildCompInfoFromJSON(requestSession, holder, projectKey, spdxInfo, compType, spdxUpload, spdxKey)

	// In case json does not contain the path "hasExtractedLicensingInfos" the expression value[0].String() is resolved to an empty string ("") letting variable licenses not initialized, if just declared
	// Therefore initialization as an empty slice.
	extractedLicenses := make([]project.ExtractedLicense, 0)
	rawLicenses := gjson.GetMany(spdxContent, "hasExtractedLicensingInfos")
	if rawLicenses[0].String() != "" {
		err := json.Unmarshal([]byte(rawLicenses[0].String()), &extractedLicenses)
		exception.HandleErrorServerMessage(err, message.GetI18N(message.UnmarshallingLicenseContent))
	}

	knownLicenses, unknownLicenses := ProcessComponentLicenses[*license.LicenseNameIdDto](requestSession, &compInfo, holder.LicenseRepository, (*license.License).ToNameIdDto)

	return project.ComponentLicenses{
		UnknownLicenses: unknownLicenses,
		KnownLicenses:   knownLicenses,
	}
}

func findCompInFile(content, spdxId string) (gjson.Result, components.ComponentType) {
	compType := components.PACKAGE
	spdxInfo := gjson.GetMany(content, "packages.#(SPDXID=="+spdxId+")")
	if spdxInfo[0].String() == "" {
		spdxInfo = gjson.GetMany(content, "snippets.#(SPDXID=="+spdxId+")")
		if spdxInfo[0].String() == "" {
			spdxInfo = gjson.GetMany(content, "files.#(SPDXID=="+spdxId+")")
			if spdxInfo[0].String() == "" {
				exception.ThrowExceptionClient404Message3(message.GetI18N(message.FindingSpdxId, spdxId))
			} else {
				compType = components.FILE
			}
		} else {
			compType = components.SNIPPET
		}
	}

	return spdxInfo[0], compType
}

func buildCompInfoFromJSON(requestSession *logy.RequestSession, holder projectService.RepositoryHolder, projectKey string, info gjson.Result, compType components.ComponentType, spdxUpload *time.Time, spdxKey string) components.ComponentInfo {
	var list []components.ComponentInfo
	project.ExtractComponentInfo(&list, compType)(gjson.Result{}, info)
	currentRefs := holder.LicenseRepository.GetLicenseRefs(requestSession)
	components.ComponentInfos(list).EnrichComponentInfos(requestSession)
	components.ComponentInfos(list).ApplyRefs(currentRefs)

	// #6642: apply license rules
	licenseRules := holder.LicenseRulesRepository.FindByKey(requestSession, projectKey, false)
	components.ComponentInfos(list).ApplyLicenseRules(licenseRules, spdxUpload, spdxKey)
	return list[0]
}

var copyAttributes = []string{
	"licenseDeclared",
	"licenseConcluded",
	"copyrightText",
	"homepage",
	"description",
	"comment",
	"SPDXID",
	"snippetFromFile",
	"licenseComments",
	"supplier",
	"summary",
	"sourceInfo",
	"originator",
	"licenseInfoFromFiles",
	"hasFiles",
	"packageFileName",
	"licenseInfoInSnippets",
	"downloadLocation",
	"filesAnalyzed",
	"packageComment",
	"packageHomepage",
}

const HAS_FILES_MAX_LEN = 50

func getAttributes(details map[string]json.RawMessage) (res []project.Detail) {
	for _, copy := range copyAttributes {
		rawValue, ok := details[copy]
		value := ""
		if ok {
			if copy == "licenseInfoFromFiles" {
				var values []string
				err := json.Unmarshal(rawValue, &values)
				exception.HandleErrorServerMessage(err, message.GetI18N(message.UnmarshallingSpdxContent))
				value = strings.Join(values, ", ")
			} else if copy == "hasFiles" {
				err := json.Unmarshal(rawValue, &value)
				if err != nil {
					var values []string
					err := json.Unmarshal(rawValue, &values)
					exception.HandleErrorServerMessage(err, message.GetI18N(message.UnmarshallingSpdxContent))
					value = strings.Join(values, ", ")
					if len(value) > HAS_FILES_MAX_LEN {
						value = value[:HAS_FILES_MAX_LEN] + "..."
					}
				}
			} else if copy == "filesAnalyzed" {
				var boolValue bool
				err := json.Unmarshal(rawValue, &boolValue)
				exception.HandleErrorServerMessage(err, message.GetI18N(message.UnmarshallingSpdxContent))
				if boolValue {
					value = "true"
				} else {
					value = "false"
				}
			} else {
				err := json.Unmarshal(rawValue, &value)
				exception.HandleErrorServerMessage(err, message.GetI18N(message.UnmarshallingSpdxContent))
			}
		}
		res = append(res, project.Detail{
			Key:   copy,
			Value: value,
		})
	}

	return res
}

func ProcessCompLicenses(rs *logy.RequestSession, c *components.ComponentInfo, licNames map[string]string) (unknown []string, known []project.SbomKnownLicenseDto) {
	decUnknown, decKnown := processList(rs, c.LicensesDeclared, licNames)
	conUnknown, conKnown := processList(rs, c.LicensesConcluded, licNames)
	unknown = append(decUnknown, conUnknown...)
	known = append(decKnown, conKnown...)
	return
}

func processList(rs *logy.RequestSession, list components.LicenseList, licNames map[string]string) (unknown []string, known []project.SbomKnownLicenseDto) {
	for _, l := range list.List {
		if !l.Known {
			unknown = append(unknown, l.OrigName)
			continue
		}
		name, ok := licNames[l.ReferencedLicense]
		if !ok {
			logy.Warnf(rs, "no name found for license %s", l.ReferencedLicense)
			continue
		}
		known = append(known, project.SbomKnownLicenseDto{
			Id:     l.ReferencedLicense,
			OrigId: l.OrigName,
			Name:   name,
		})
	}
	return
}

func ProcessComponentLicenses[T any](requestSession *logy.RequestSession, c *components.ComponentInfo, licenseRepo license2.ILicensesRepository, toDtoFunc func(*license.License) T) ([]project.DetailedLicense[T], []string) {
	knownLicenses := []project.DetailedLicense[T]{}
	unknown := []string{}
	fetched := make(map[string]struct{})

	knownDeclared := processLicenseList[T](requestSession, licenseRepo, c.LicensesDeclared.List, fetched, &unknown, toDtoFunc)
	knownConcluded := processLicenseList[T](requestSession, licenseRepo, c.LicensesConcluded.List, fetched, &unknown, toDtoFunc)

	knownLicenses = append(knownLicenses, knownConcluded...)
	knownLicenses = append(knownLicenses, knownDeclared...)

	return knownLicenses, unknown
}

func processLicenseList[T any](
	requestSession *logy.RequestSession,
	licenseRepo license2.ILicensesRepository,
	in []*components.ComponentLicense,
	fetched map[string]struct{},
	unknown *[]string,
	toDtoFunc func(*license.License) T,
) []project.DetailedLicense[T] {
	var res []project.DetailedLicense[T]
	for _, l := range in {
		if !l.Known {
			alreadyIn := false
			for _, u := range *unknown {
				if u == l.OrigName {
					alreadyIn = true
					break
				}
			}
			if !alreadyIn {
				*unknown = append(*unknown, l.OrigName)
			}
			continue
		}
		_, ok := fetched[l.ReferencedLicense]
		if ok {
			continue
		}
		licenseData := licenseRepo.FindById(requestSession, l.ReferencedLicense)
		fetched[licenseData.LicenseId] = struct{}{}
		res = append(res, project.DetailedLicense[T]{
			License:        toDtoFunc(licenseData),
			OrigName:       l.OrigName,
			ReferencedName: l.ReferencedLicense,
		})
	}
	return res
}
