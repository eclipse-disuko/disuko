// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package notices

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"mercedes-benz.ghe.com/foss/disuko/domain/label"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/labels"

	"mercedes-benz.ghe.com/foss/disuko/domain/license"
	"mercedes-benz.ghe.com/foss/disuko/domain/project"
	"mercedes-benz.ghe.com/foss/disuko/domain/project/components"
	"mercedes-benz.ghe.com/foss/disuko/helper"
	license2 "mercedes-benz.ghe.com/foss/disuko/infra/repository/license"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

func GenerateJSONNotices(requestSession *logy.RequestSession, currentProject project.Project,
	licenseRepository license2.ILicensesRepository, componentInfo components.ComponentInfos, contactMeta project.NoticeContactMeta, labelRepository labels.ILabelRepository) *project.NoticeFileJSON {

	res := project.NoticeFileJSON{
		Components: []project.NoticeComponent{},
		Licenses:   []project.NoticeLicense{},
		Meta: project.NoticeJSONMeta{
			Title:       "",
			Description: "",
		},
	} //	@name	NoticeFile

	groupedLicenses := groupFOSSLicensesByComponent(requestSession, licenseRepository, componentInfo)
	groupedLicenses.Sort(true)

	for _, group := range groupedLicenses {
		if len(group.licenses) == 0 {
			continue
		}
		res.Components = append(res.Components, project.NoticeComponent{
			Name:        group.component.Name,
			Version:     group.component.Version,
			LicenseName: group.licenseNamesString(),
			LicenseID:   group.licenseIdsString(),
			Copyright:   group.component.CopyrightText,
		})

	}

	groupedComponents := groupComponentsByFOSSLicense(requestSession, licenseRepository, componentInfo)
	groupedComponents.Sort(false)

	for _, group := range groupedComponents {
		res.Licenses = append(res.Licenses, project.NoticeLicense{
			Name:      group.license.Name,
			Text:      group.license.Text,
			LicenseID: group.license.LicenseId,
		})
	}

	return &res
}

func GenerateHTMLNotices(requestSession *logy.RequestSession, currentProject project.Project, licenseRepository license2.ILicensesRepository, components components.ComponentInfos, contactMeta project.NoticeContactMeta, labelRepository labels.ILabelRepository) *strings.Builder {
	var sb strings.Builder
	sb.WriteString("<h1>Components:</h1><br>")

	groupedComponents := groupComponentsByFOSSLicense(requestSession, licenseRepository, components)
	groupedComponents.Sort(true)

	for _, group := range groupedComponents {
		sb.WriteString("<h2 style='text-align: center; padding: 20px;'>- " + encodeHtmlTagsAndNextLineForHtmlOutput(group.license.Name) + " -</h2><br>")
		sb.WriteString("<table style='width: 100%'><thead><tr><th style='width: 30%; text-align: left;'>Component Info</th><th style='width: 20%; text-align: left;'>Version</th><th style='width: 50%; text-align: left;'>Copyright</th></tr></thead><tbody>")
		for _, component := range group.components {
			sb.WriteString("<tr><td style='width: 30%;'>")
			sb.WriteString(encodeHtmlTagsAndNextLineForHtmlOutput(component.Name))
			sb.WriteString("</td><td style='width: 20%;'>")
			sb.WriteString(encodeHtmlTagsAndNextLineForHtmlOutput(component.Version))
			sb.WriteString("</td><td style='width: 50%;'>")
			if !helper.IsUnasserted(component.CopyrightText) {
				sb.WriteString(encodeHtmlTagsAndNextLineForHtmlOutput(component.CopyrightText))
			} else {
				sb.WriteString("-")
			}
			sb.WriteString("</td></tr>")
		}
		sb.WriteString("</tbody></table><p style='padding-top: 20px'>" + encodeHtmlTagsAndNextLineForHtmlOutput(group.license.Text) + "</p><br><br><br>")
	}

	return &sb
}

func GenerateTextNotices(requestSession *logy.RequestSession, currentProject project.Project, licenseRepository license2.ILicensesRepository, components components.ComponentInfos, contactMeta project.NoticeContactMeta, labelRepository labels.ILabelRepository) *strings.Builder {
	groupedLicenses := groupFOSSLicensesByComponent(requestSession, licenseRepository, components)
	groupedLicenses.Sort(true)

	var sb strings.Builder
	sb.WriteString("\n\n- Components -\n\n")

	for _, group := range groupedLicenses {
		sb.WriteString(group.component.Name)
		sb.WriteString("/")
		sb.WriteString(group.component.Version)
		sb.WriteString(" : ")
		for i, l := range group.licenses {
			sb.WriteString(fmt.Sprintf("%s (%s)", l.Name, l.LicenseId))
			if i < len(group.licenses)-1 {
				sb.WriteString(string(group.component.GetLicensesEffective().Op))
			}
		}
		sb.WriteString("\n")
	}
	sb.WriteString("\n- Copyright Texts -\n\n")
	alreadyAddedUniqueComponentVersion := make(map[string]bool)
	for _, group := range groupedLicenses {
		if helper.IsUnasserted(group.component.CopyrightText) {
			continue
		}
		uniqueComponentVersion := group.component.Name + "/" + group.component.Version
		if _, ok := alreadyAddedUniqueComponentVersion[uniqueComponentVersion]; ok {
			continue
		}
		alreadyAddedUniqueComponentVersion[uniqueComponentVersion] = true
		sb.WriteString(uniqueComponentVersion)
		sb.WriteString(" : \n\t")
		sb.WriteString(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(group.component.CopyrightText, "\n", "\n\t"), "<", "&lt;"), ">", "&gt;"))
		sb.WriteString("\n\n")
	}

	groupedComponents := groupComponentsByFOSSLicense(requestSession, licenseRepository, components)
	groupedComponents.Sort(false)
	sb.WriteString("\n- Licenses -\n\n")
	for _, group := range groupedComponents {
		sb.WriteString(group.license.LicenseId)
		sb.WriteString(" - ")
		sb.WriteString(group.license.Name)
		sb.WriteString(" : \n\t")
		sb.WriteString(group.license.Text)
		sb.WriteString("\n\n\n\n")
	}

	return &sb
}

func hasOnboardLabel(requestSession *logy.RequestSession, currentProject project.Project, labelRepository labels.ILabelRepository) bool {
	onboardLabel := labelRepository.FindByNameAndType(requestSession, label.ONBOARD, label.POLICY)
	if onboardLabel != nil {
		return slices.Contains(currentProject.PolicyLabels, onboardLabel.GetKey())
	} else {
		return false
	}
}

func encodeHtmlTagsAndNextLineForHtmlOutput(source string) string {
	result := strings.ReplaceAll(source, "<", "&lt;")
	result = strings.ReplaceAll(result, ">", "&gt;")
	result = strings.ReplaceAll(result, "\n", "<br>\n")
	return result
}

type componentsGrouped []*componentsGroup

type componentsGroup struct {
	license    *license.License
	components []*components.ComponentInfo
}

func groupComponentsByFOSSLicense(requestSession *logy.RequestSession, licenseRepository license2.ILicensesRepository, componentList components.ComponentInfos) componentsGrouped {
	res := make(componentsGrouped, 0)
	blacklist := make(map[string]bool, 0)
	alreadyAdded := make(map[string]int, 0)

	for compI := range componentList {
		for _, compLicense := range componentList[compI].GetLicensesEffective().List {
			if !compLicense.Known || blacklist[compLicense.ReferencedLicense] {
				continue
			} else if i, ok := alreadyAdded[compLicense.ReferencedLicense]; ok {
				res[i].components = append(res[i].components, &componentList[compI])
			} else {
				lic := licenseRepository.FindById(requestSession, compLicense.ReferencedLicense)
				if lic.Meta.LicenseType == license.OpenSource || lic.Meta.LicenseType == license.PublicDomain {
					res = append(res, &componentsGroup{
						license:    lic,
						components: []*components.ComponentInfo{&componentList[compI]},
					})
					alreadyAdded[compLicense.ReferencedLicense] = len(res) - 1
				} else {
					blacklist[lic.LicenseId] = true
				}
			}
		}
	}

	return res
}

func (c componentsGrouped) Sort(sortInner bool) {
	sort.Slice(c, func(i, j int) bool {
		return strings.ToLower(c[i].license.Name) < strings.ToLower(c[j].license.Name)
	})

	if !sortInner {
		return
	}
	for _, g := range c {
		sort.Slice(g.components, func(i, j int) bool {
			return strings.ToLower(g.components[i].Name) < strings.ToLower(g.components[j].Name)
		})
	}
}

type licensesGrouped []*licenseGroup

type licenseGroup struct {
	component components.ComponentInfo
	licenses  []*license.License
}

func groupFOSSLicensesByComponent(requestSession *logy.RequestSession, licenseRepository license2.ILicensesRepository, components components.ComponentInfos) licensesGrouped {
	res := make(licensesGrouped, 0)
	// save some license fetching
	licenseCache := make(map[string]*license.License, 0)
	alreadyAdded := make(map[string]int, 0)

	for _, component := range components {
		for _, compLicense := range component.GetLicensesEffective().List {
			if !compLicense.Known {
				continue
			}
			lic, cached := licenseCache[compLicense.ReferencedLicense]
			if !cached {
				lic = licenseRepository.FindById(requestSession, compLicense.ReferencedLicense)
				licenseCache[compLicense.ReferencedLicense] = lic
			}
			if lic.Meta.LicenseType != license.OpenSource && lic.Meta.LicenseType != license.PublicDomain {
				continue
			}
			if i, ok := alreadyAdded[component.Name]; ok {
				res[i].licenses = append(res[i].licenses, lic)
			} else {
				res = append(res, &licenseGroup{
					component: component,
					licenses:  []*license.License{lic},
				})
			}
		}
	}

	return res
}

func (l licensesGrouped) Sort(sortInner bool) {
	sort.Slice(l, func(i, j int) bool {
		return strings.ToLower(l[i].component.Name) < strings.ToLower(l[j].component.Name)
	})

	if !sortInner {
		return
	}
	for _, g := range l {
		sort.Slice(g.licenses, func(i, j int) bool {
			return strings.ToLower(g.licenses[i].Name) < strings.ToLower(g.licenses[j].Name)
		})
	}
}

func (g licenseGroup) licenseNamesString() string {
	res := ""

	for _, license := range g.licenses {
		res += license.Name
		res += string(g.component.GetLicensesEffective().Op)
	}

	return strings.TrimSuffix(res, string(g.component.GetLicensesEffective().Op))
}

func (g licenseGroup) licenseIdsString() string {
	res := ""

	for _, license := range g.licenses {
		res += license.LicenseId
		res += string(g.component.GetLicensesEffective().Op)
	}

	return strings.TrimSuffix(res, string(g.component.GetLicensesEffective().Op))
}
