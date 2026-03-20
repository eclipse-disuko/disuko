// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package scanremarks

import (
	"slices"
	"strings"
	"unicode"

	"github.com/eclipse-disuko/disuko/domain/label"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/project/components"
	"github.com/eclipse-disuko/disuko/helper"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/infra/repository/labels"
	projectRepo "github.com/eclipse-disuko/disuko/infra/repository/project"
	"github.com/eclipse-disuko/disuko/logy"
)

const OrLinksThreshold = 4

type Service struct {
	ProjectRepo projectRepo.IProjectRepository
	LabelsRepo  labels.ILabelRepository
}

func (s *Service) GetRemarks(rs *logy.RequestSession, p *project.Project, meta *project.MetaInfo, evalRes *components.EvaluationResult) []project.QualityScanRemarks {
	result := make([]project.QualityScanRemarks, 0)

	if meta.HasExternalRefs && len(evalRes.Results) > 0 {
		result = appendQualityScanRemarks(result, components.ComponentResult{Component: &components.ComponentInfo{}}, project.PROBLEM, message.SrContentException, message.SrContainsExternalRefs)
	}

	contactAddress := p.NoticeContactMeta.Address
	if p.HasParent() {
		parentProject := s.ProjectRepo.FindByKey(nil, p.Parent, false)
		if parentProject != nil {
			contactAddress = parentProject.NoticeContactMeta.Address
		}
	}
	if contactAddress == "" && !hasOnboardLabel(rs, p, s.LabelsRepo) {
		result = appendQualityScanRemarks(result, components.ComponentResult{Component: &components.ComponentInfo{}}, project.WARNING,
			message.SrProjectAddressException, message.SrProjectAddressDescription)
	}

	hasNotAllowedUnicodeLetters := func(str string) bool {
		allowedLetters := "<>'@+-©[]"
		for _, letter := range str {
			isAllowed := strings.IndexRune(allowedLetters, letter) > -1
			if isAllowed {
				// is allowed unicode, check next letter
				continue
			}
			if unicode.IsSymbol(letter) {
				return true
			}
		}
		return false
	}

	isMalformedCopyright := func(str string) bool {
		badText := []string{
			"false",
			"copyright:before",
			"null",
			"key",
			"context",
		}
		for _, b := range badText {
			if strings.Contains(str, b) {
				return true
			}
		}
		return false
	}

	hasNonLatinOrNotAllowedLetters := func(str string) bool {
		allowedLetters := "#@:/.-_%=?~&"
		for _, letter := range str {
			isAllowed := strings.IndexRune(allowedLetters, letter) > -1 || unicode.IsDigit(letter)
			if isAllowed {
				continue
			}
			if (letter < 'a' || letter > 'z') && (letter < 'A' || letter > 'Z') {
				return true
			}
		}
		return false
	}

	var projectLevel string
	for _, compRes := range evalRes.Results {
		if compRes.Component.HasAnnotations {
			result = appendQualityScanRemarks(result, compRes, project.INFORMATION,
				message.SrContentException, message.SrContainsAnnotations)
		}
		if compRes.Component.Type == components.SNIPPET {
			result = appendQualityScanRemarks(result, compRes, project.PROBLEM,
				message.SrContentException, message.SrContainsSnippet)
		}
		if len(compRes.Component.Version) == 0 || helper.IsUnasserted(compRes.Component.Version) {
			result = appendQualityScanRemarks(result, compRes, project.PROBLEM,
				message.SrMissingVersion, message.SrMissingVersionDescription)
		}
		if compRes.Component.ComplexExpression {
			result = appendQualityScanRemarks(result, compRes, project.PROBLEM,
				message.SrContentException, message.SrContainsComplex)
		} else if compRes.Component.LicensesDeclared.CountOrLinks() >= OrLinksThreshold ||
			compRes.Component.LicensesConcluded.CountOrLinks() >= OrLinksThreshold {
			result = appendQualityScanRemarks(result, compRes, project.WARNING,
				message.SrTooMuchOrTitle, message.SrContainsTooMuchOr)
		}
		if compRes.Component.ContainsBadChars {
			result = appendQualityScanRemarks(result, compRes, project.PROBLEM,
				message.SrContentException, message.SrContainsBadchars)
		}
		if len(compRes.Component.Name) == 0 || helper.IsUnasserted(compRes.Component.Name) {
			result = appendQualityScanRemarks(result, compRes, project.PROBLEM,
				message.SrMissingName, message.SrMissingNameDescription)
		}

		if len(compRes.Component.GetLicenseEffective()) == 0 || helper.IsUnasserted(compRes.Component.GetLicenseEffective()) {
			result = appendQualityScanRemarks(result, compRes, project.PROBLEM,
				message.SrMissingLicenseId, message.SrMissingLicenseIdDescription)
		}

		// check CopyrightText - exist
		if len(compRes.Component.CopyrightText) == 0 || helper.IsUnasserted(compRes.Component.CopyrightText) {
			if projectLevel == "" {
				projectLevel = copyrightMissingLevel(rs, p, s.LabelsRepo)
			}
			result = appendQualityScanRemarks(result, compRes, projectLevel,
				message.SrMissingCopyrightText, message.SrMissingCopyrightDescription)
		} else {
			// check CopyrightText - content
			if hasNotAllowedUnicodeLetters(compRes.Component.CopyrightText) ||
				// year template
				strings.Contains(strings.ToLower(compRes.Component.CopyrightText), "yyyy") ||
				isMalformedCopyright(compRes.Component.CopyrightText) {
				result = appendQualityScanRemarks(result, compRes, project.INFORMATION,
					message.SrMalformedCopyrightText, message.SrMalformedCopyrightDescription)
			}

			// check CopyrightText - length
			if len(compRes.Component.CopyrightText) > 1000 {
				result = appendQualityScanRemarks(result, compRes, project.INFORMATION,
					message.SrCopyrightLongText, message.SrCopyrightToLongDescription)
			}

		}

		if !helper.IsUnasserted(compRes.Component.License) && !helper.IsUnasserted(compRes.Component.LicenseDeclared) &&
			compRes.Component.License != compRes.Component.LicenseDeclared {
			result = appendQualityScanRemarks(result, compRes, project.WARNING,
				message.SrLicensesDiff, message.SrLicensesDiffDescription)
		}

		for _, l := range compRes.Component.GetLicensesEffective().List {
			if !l.Known && !helper.IsUnasserted(l.OrigName) {
				result = appendQualityScanRemarks(result, compRes, project.PROBLEM,
					message.SrUnknownLicenseUsed, message.SrUnknownLicenseUsedDescription)
			}
			if compRes.ContainsUnmatchedLicense(l.OrigName) {
				result = appendQualityScanRemarks(result, compRes, project.PROBLEM,
					message.SrUnmatchedLicenseUsed, message.SrUnmatchedLicenseUsedDescription)
			}
			if !strings.EqualFold(l.ReferencedLicense, l.OrigName) && l.Known {
				result = appendQualityScanRemarks(result, compRes, project.INFORMATION,
					message.SrAliasingUsed, message.SrAliasingUsedDescription)
			}
		}

		if len(compRes.Component.PURL) > 0 && compRes.Component.PURL != "NOASSERTION" && hasNonLatinOrNotAllowedLetters(compRes.Component.PURL) {
			result = appendQualityScanRemarks(result, compRes, project.WARNING,
				message.SrContainsNonLatinLetters, message.SrContainsNonLatinLettersDescription)
		}

	}

	return result
}

func appendQualityScanRemarks(result []project.QualityScanRemarks,
	compRes components.ComponentResult, problem string, remarkKey string, descriptionKey string,
) []project.QualityScanRemarks {
	if canBeOmitted(compRes, remarkKey) {
		return result
	}
	componentType := project.ComponentType(compRes.Component.Type)
	if descriptionKey == message.SrProjectAddressDescription {
		componentType = project.PROJECT
	}
	result = append(result, project.QualityScanRemarks{
		Status:            project.ScanRemarkStatus(problem),
		RemarkKey:         remarkKey,
		Name:              "" + compRes.Component.Name,
		SpdxId:            "" + compRes.Component.SpdxId,
		Version:           "" + compRes.Component.Version,
		Type:              componentType,
		DescriptionKey:    descriptionKey,
		PolicyRuleStatus:  components.ToPolicyStatusDto(compRes.Status, false),
		UnmatchedLicenses: components.ToUnmatchedDto(compRes.Unmatched),
	})
	return result
}

func canBeOmitted(compRes components.ComponentResult, remarkKey string) bool {
	if compRes.Component.Type != components.ROOT || len(compRes.Component.GetLicensesEffective().List) > 0 {
		return false
	}
	omitRemarks := []string{
		message.SrMissingLicenseText,
		message.SrMissingLicenseId,
		message.SrMissingCopyrightText,
		message.SrMissingVersion,
	}
	return helper.Contains(remarkKey, omitRemarks)
}

func hasOnboardLabel(requestSession *logy.RequestSession, currentProject *project.Project, labelRepository labels.ILabelRepository) bool {
	onboardLabel := labelRepository.FindByNameAndType(requestSession, label.ONBOARD, label.POLICY)
	if onboardLabel != nil {
		return slices.Contains(currentProject.PolicyLabels, onboardLabel.GetKey())
	} else {
		return false
	}
}

func copyrightMissingLevel(rs *logy.RequestSession, pr *project.Project, labelRepo labels.ILabelRepository) string {
	ll := make([]string, 0)
	for _, l := range pr.PolicyLabels {
		resolvedLabel := labelRepo.FindByKey(rs, l, false)
		if resolvedLabel == nil {
			continue
		}
		ll = append(ll, resolvedLabel.Name)
	}

	switch {
	case helper.Contains("enterprise platform", ll) && helper.Contains("entity users", ll):
		return project.INFORMATION
	case helper.Contains("enterprise platform", ll) && helper.Contains("group users", ll):
		return project.INFORMATION
	case helper.Contains("vehicle platform", ll) && helper.Contains("onboard", ll):
		return project.PROBLEM
	case helper.Contains("external users", ll):
		return project.PROBLEM
	case helper.Contains("other platform", ll):
		return project.PROBLEM
	}

	return project.WARNING
}
