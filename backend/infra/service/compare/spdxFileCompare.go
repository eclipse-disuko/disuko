// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package compare

import (
	"strings"

	"github.com/eclipse-disuko/disuko/domain/compare"
	"github.com/eclipse-disuko/disuko/domain/project/components"
)

func MultiCompareSpdxFiles(evalOld, evalNew *components.EvaluationResult, isResponsible bool) []*compare.ComponentMultiDiffDto {
	oldComponentsMapList := convertToMapList(evalOld)
	newComponentsMapList := convertToMapList(evalNew)

	var components []*compare.ComponentMultiDiffDto

	// find new
	for componentName, newList := range newComponentsMapList {
		if oldList, exists := oldComponentsMapList[componentName]; exists {
			componentResult := compare.NewComponentMultiDiffDto(compare.UNCHANGED, componentName)
			for newVersion, newComponent := range newList {
				if sameVersionOldComponent, present := oldList[newVersion]; present {
					evaluateMultiDiffs(sameVersionOldComponent, newComponent, componentResult)
					continue
				}
				for oldVersion, oldComponent := range oldList {
					if _, present := newList[oldVersion]; present {
						continue
					}
					evaluateMultiDiffs(oldComponent, newComponent, componentResult)
				}
			}
			for _, componentNew := range newList {
				componentResult.ComponentsNew = append(componentResult.ComponentsNew, componentNew.ToComponentInfoDto(isResponsible, "_", false, nil, nil, nil))
			}
			for _, componentOld := range oldList {
				componentResult.ComponentsOld = append(componentResult.ComponentsOld, componentOld.ToComponentInfoDto(isResponsible, "_", false, nil, nil, nil))
			}
			if len(componentResult.ComponentsNew) != len(componentResult.ComponentsOld) {
				componentResult.DiffType = compare.CHANGED
			}
			components = append(components, componentResult)
		} else {
			componentResult := compare.NewComponentMultiDiffDto(compare.NEW, componentName)
			for _, componentNew := range newList {
				componentResult.ComponentsNew = append(componentResult.ComponentsNew, componentNew.ToComponentInfoDto(isResponsible, "_", false, nil, nil, nil))
			}
			components = append(components, componentResult)
		}
	}

	// find removed entries
	for componentName, oldList := range oldComponentsMapList {
		if _, exists := newComponentsMapList[componentName]; exists {
			continue
		}
		componentResult := compare.NewComponentMultiDiffDto(compare.REMOVED, componentName)
		for _, componentOld := range oldList {
			componentResult.ComponentsOld = append(componentResult.ComponentsOld, componentOld.ToComponentInfoDto(isResponsible, "_", false, nil, nil, nil))
		}
		components = append(components, componentResult)
	}

	return components
}

func evaluateMultiDiffs(componentOld components.ComponentResult, componentNew components.ComponentResult, multiDiff *compare.ComponentMultiDiffDto) {
	// evaluate and create changes map for every old/new pair with the concatenated versions string "oldVersion_newVersion" as a key
	componentChanges := compare.ComponentChangesDto{}
	versionToVersionLocalDiff := false

	// Modified
	if componentOld.Component.Modified != componentNew.Component.Modified {
		componentChanges.Modified = true
		versionToVersionLocalDiff = true
		multiDiff.DiffType = compare.CHANGED
	}
	// SpdxId
	if notEqualIgnoreCaseSensitive(componentOld.Component.SpdxId, componentNew.Component.SpdxId) {
		componentChanges.SpdxId = true
		versionToVersionLocalDiff = true
		multiDiff.DiffType = compare.CHANGED
	}
	// Name
	if notEqualIgnoreCaseSensitive(componentOld.Component.Name, componentNew.Component.Name) {
		componentChanges.Name = true
		versionToVersionLocalDiff = true
		multiDiff.DiffType = compare.CHANGED
	}
	// Version
	if notEqualIgnoreCaseSensitive(componentOld.Component.Version, componentNew.Component.Version) {
		componentChanges.Version = true
		versionToVersionLocalDiff = true
		multiDiff.DiffType = compare.CHANGED
	}
	// CopyrightText
	if notEqualIgnoreCaseSensitive(componentOld.Component.CopyrightText, componentNew.Component.CopyrightText) {
		componentChanges.CopyrightText = true
		versionToVersionLocalDiff = true
		multiDiff.DiffType = compare.CHANGED
	}
	// LicenseEffective
	if notEqualIgnoreCaseSensitive(componentOld.Component.EffectiveLicensesString(), componentNew.Component.EffectiveLicensesString()) {
		componentChanges.LicenseEffective = true
		versionToVersionLocalDiff = true
		multiDiff.DiffType = compare.CHANGED
	}
	// License
	if notEqualIgnoreCaseSensitive(componentOld.Component.License, componentNew.Component.License) {
		componentChanges.License = true
		versionToVersionLocalDiff = true
		multiDiff.DiffType = compare.CHANGED
	}
	// LicenseDeclared
	if notEqualIgnoreCaseSensitive(componentOld.Component.LicenseDeclared, componentNew.Component.LicenseDeclared) {
		componentChanges.LicenseDeclared = true
		versionToVersionLocalDiff = true
		multiDiff.DiffType = compare.CHANGED
	}
	// LicenseComments
	if notEqualIgnoreCaseSensitive(componentOld.Component.LicenseComments, componentNew.Component.LicenseComments) {
		componentChanges.LicenseComments = true
		versionToVersionLocalDiff = true
		multiDiff.DiffType = compare.CHANGED
	}
	// Type
	if componentOld.Component.Type != componentNew.Component.Type {
		componentChanges.Type = true
		versionToVersionLocalDiff = true
		multiDiff.DiffType = compare.CHANGED
	}
	// Description
	if notEqualIgnoreCaseSensitive(componentOld.Component.Description, componentNew.Component.Description) {
		componentChanges.Description = true
		versionToVersionLocalDiff = true
		multiDiff.DiffType = compare.CHANGED
	}
	// DownloadLocation
	if notEqualIgnoreCaseSensitive(componentOld.Component.DownloadLocation, componentNew.Component.DownloadLocation) {
		componentChanges.DownloadLocation = true
		versionToVersionLocalDiff = true
		multiDiff.DiffType = compare.CHANGED
	}
	// Questioned
	if componentOld.Questioned != componentNew.Questioned {
		componentChanges.Questioned = true
		versionToVersionLocalDiff = true
		multiDiff.DiffType = compare.CHANGED
	}
	// Unasserted
	if componentOld.Unasserted != componentNew.Unasserted {
		componentChanges.Unasserted = true
		versionToVersionLocalDiff = true
		multiDiff.DiffType = compare.CHANGED
	}

	if componentOld.Component.PURL != componentNew.Component.PURL {
		componentChanges.PURL = true
		versionToVersionLocalDiff = true
		multiDiff.DiffType = compare.CHANGED
	}

	// prStatus
	usedPolicyRuleOld, _ := componentOld.GetUsedPolicyRule()
	usedPolicyRuleNew, _ := componentNew.GetUsedPolicyRule()
	if usedPolicyRuleOld != usedPolicyRuleNew {
		componentChanges.PrStatus = true
		versionToVersionLocalDiff = true
		multiDiff.DiffType = compare.CHANGED
	}

	if versionToVersionLocalDiff {
		changesKey := strings.ToLower(componentOld.Component.Version) + "_" + strings.ToLower(componentNew.Component.Version)
		multiDiff.Changes[changesKey] = componentChanges
	}
}

func notEqualIgnoreCaseSensitive(str1 string, str2 string) bool {
	return strings.ToLower(str1) != strings.ToLower(str2)
}

func convertToMapList(evalRes *components.EvaluationResult) map[string]map[string]components.ComponentResult {
	result := make(map[string]map[string]components.ComponentResult)

	for _, compRes := range evalRes.Results {

		key := createID(*compRes.Component)
		componentResultMap, exists := result[key]
		if !exists {
			componentResultMap = make(map[string]components.ComponentResult, 0)
		}
		componentResultMap[compRes.Component.Version] = compRes
		result[key] = componentResultMap
	}
	return result
}

func createID(value components.ComponentInfo) string {
	return strings.ToLower(value.Name)
}
