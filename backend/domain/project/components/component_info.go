// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package components

import (
	"strings"
	"time"

	"github.com/eclipse-disuko/disuko/domain/licenserules"
	"github.com/eclipse-disuko/disuko/domain/policydecisions"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/helper"
	"github.com/eclipse-disuko/disuko/logy"
)

type ComponentInfoList struct {
	domain.RootEntity

	UsedRefsListHash     string
	UsedLicenseRulesHash string
	ComponentInfos       ComponentInfos
}

type ComponentInfo struct {
	domain.ChildEntity

	Type ComponentType

	SpdxId            string
	Name              string
	Version           string
	License           string
	LicenseDeclared   string
	LicenseComments   string
	CopyrightText     string
	Description       string
	DownloadLocation  string
	PURL              string
	HomepageURL       string
	Modified          bool
	HasAnnotations    bool
	ComplexExpression bool
	ContainsBadChars  bool
	ContainsAuxiliary bool

	LicensesDeclared  LicenseList
	LicensesConcluded LicenseList

	LicenseRuleApplied *licenserules.LicenseRule

	PolicyDecisionsApplied []*policydecisions.PolicyDecision
}

func (cinfo *ComponentInfo) IsAliasUsed() bool {
	for _, lic := range cinfo.GetLicensesEffective().List {
		if lic.Known && !strings.EqualFold(lic.ReferencedLicense, lic.OrigName) {
			return true
		}
	}
	return false
}

type LicenseAppliedType string

const (
	LicenseDeclared  LicenseAppliedType = "declared"
	LicenseConcluded LicenseAppliedType = "concluded"
)

type ComponentType string

const (
	PACKAGE ComponentType = "Package"
	FILE    ComponentType = "File"
	SNIPPET ComponentType = "Snippet"
	ROOT    ComponentType = "Root"
)

type LicenseList struct {
	List []*ComponentLicense
	Op   Operator
}

// That is needed to use the Length within a go lang template
func (lc *LicenseList) Length() int {
	if lc == nil || lc.List == nil {
		return 0
	}
	return len(lc.List)
}

type Operator string

const (
	AND      Operator = " and "
	OR       Operator = " or "
	MULTIPLE Operator = "unsupported"
	NONE     Operator = ""
)

const withExpression string = " with "

var badChars []string = []string{
	";",
	"/",
	"\\",
	"&",
	"|",
}

type ComponentLicense struct {
	OrigName          string
	ReferencedLicense string
	LicenseFamily     license.FamilyOfLicense
	ApprovalState     license.ApprovalStatus
	Known             bool
}

type ComponentInfos []ComponentInfo

func (cis ComponentInfos) ApplyLicenseRules(licenseRules *licenserules.LicenseRules, sbomUpload *time.Time, sbomKey string) {
	if licenseRules == nil {
		return
	}
	for i := 0; i < len(cis); i++ {
		compLicenseExpression := cis[i].EffectiveLicensesString()
		for _, lr := range licenseRules.Rules {
			cmpNameMatches := strings.EqualFold(lr.ComponentName, cis[i].Name)
			licExprMatches := strings.EqualFold(lr.LicenseExpression, compLicenseExpression)

			nameAndExprMatches := cmpNameMatches && licExprMatches

			if !nameAndExprMatches || sbomUpload == nil {
				continue
			}

			sbomMatches := lr.SBOMId == sbomKey
			activeOK := (lr.Active && (*sbomUpload).After(lr.Created)) || sbomMatches
			previewOK := lr.Active && (*sbomUpload).Before(lr.Created) && !sbomMatches
			cancelledOK := !lr.Active && (*sbomUpload).After(lr.Created) && (*sbomUpload).Before(lr.Updated)

			if activeOK || previewOK || cancelledOK {
				if previewOK {
					lr.PreviewMode = true
				}
				cis[i].LicenseRuleApplied = lr
				break
			}
		}
	}
}

func (cis ComponentInfos) CleanLicenseRules() {
	for i := 0; i < len(cis); i++ {
		cis[i].LicenseRuleApplied = nil
	}
}

func (cis ComponentInfos) ApplyRefs(refs license.LicenseRefs) {
	for i := 0; i < len(cis); i++ {
		applyRefs(refs, cis[i].LicensesConcluded)
		applyRefs(refs, cis[i].LicensesDeclared)
	}
}

func (cis ComponentInfos) EnrichComponentInfos(requestSession *logy.RequestSession) {
	processedComps := 0
	processedLicenses := 0

	for i := 0; i < len(cis); i++ {
		text := strings.ToLower(cis[i].GetLicenseEffective())
		op := GetOperator(text)
		cis[i].ComplexExpression = false
		// todo #6801: adjust after implementation
		if op == MULTIPLE {
			// logy.Infof(requestSession, "Ignoring %s because of multiple operators", cis[i].Name)
			cis[i].ComplexExpression = true
			continue
		}

		// todo #6801: adjust after implementation
		cis[i].ContainsBadChars = false
		if checkBadChars(text) {
			cis[i].ContainsBadChars = true
			// logy.Infof(requestSession, "Ignoring %s because of badchars", cis[i].Name)
			continue
		}

		declaredNum := 0
		declaredNum, cis[i].LicensesDeclared = parseLicenseText(requestSession, cis[i].LicenseDeclared)
		processedLicenses += declaredNum
		concludedNum := 0
		concludedNum, cis[i].LicensesConcluded = parseLicenseText(requestSession, cis[i].License)
		processedLicenses += concludedNum
		processedComps++
	}
	logy.Infof(requestSession, "EnrichComponentInfos done, processed comps: %d processed licenses %d", processedComps, processedLicenses)
}

// todo #6801: adjust after implementation
func parseLicenseText(requestSession *logy.RequestSession, text string) (processedLicenses int, res LicenseList) {
	text = strings.ToLower(text)
	// Currently we do not support brackets in licenses
	text = helper.RemoveAllBracketsInLicenseText(text)
	op := GetOperator(text)
	if op == MULTIPLE {
		return
	}
	names := ExtractNames(text, op)
	if len(names) == 0 {
		logy.Infof(requestSession, "No names could be extracted from %s op %s", text, op)
	}
	res.Op = op
	for _, name := range names {
		if strings.ToLower(name) == "noassertion" {
			continue
		}
		if len(name) == 0 {
			// logy.Infof(requestSession, "Empty name from %s op %s", text, op)
			continue
		}
		res.List = append(res.List, &ComponentLicense{
			OrigName: name,
		})
		processedLicenses++
	}
	return
}

func checkBadChars(text string) bool {
	for _, c := range badChars {
		if strings.Contains(text, c) {
			return true
		}
	}
	return false
}

func ExtractNames(text string, op Operator) []string {
	if op == NONE {
		return []string{strings.Trim(strings.TrimSpace(text), "()")}
	}

	if strings.HasPrefix(text, "(") {
		text = strings.Trim(text, "()")
	}
	res := make([]string, 0)
	spl := strings.Split(text, string(op))
	for _, s := range spl {
		res = append(res, strings.TrimSpace(s))
	}
	return res
}

// todo #6801: adjust after implementation
func GetOperator(text string) Operator {
	if strings.Contains(text, string(OR)) && strings.Contains(text, string(AND)) {
		return MULTIPLE
	} else if strings.Contains(text, string(OR)) {
		return OR
	} else if strings.Contains(text, string(AND)) {
		return AND
	}
	return NONE
}

func applyRefs(refs license.LicenseRefs, licenses LicenseList) {
	totalNum := 0
	knownNum := 0
	for _, l := range licenses.List {
		l.Known = false
		l.ReferencedLicense = ""
		totalNum++
		if referencedLicense, ok := refs[strings.ToLower(l.OrigName)]; ok {
			knownNum++
			l.Known = true
			l.ReferencedLicense = referencedLicense.ID
			l.LicenseFamily = referencedLicense.Family
			l.ApprovalState = referencedLicense.ApprovalState
		}
	}
	// logy.Infof(requestSession, "CheckKnownLicenses known licenses %d/%d.", knownNum, totalNum)
}

func (ci ComponentInfo) GetLicensesEffective() LicenseList {
	res := ci.LicensesDeclared
	if !helper.IsUnasserted(ci.License) {
		res = ci.LicensesConcluded
	}
	// #6642: check here also if a decision was made and use!
	if ci.LicenseRuleApplied == nil || ci.LicenseRuleApplied.PreviewMode {
		return res
	}
	for _, l := range res.List {
		cmpLicenseId := l.ReferencedLicense
		if cmpLicenseId == "" {
			cmpLicenseId = l.OrigName
		}
		if strings.EqualFold(cmpLicenseId, ci.LicenseRuleApplied.LicenseDecisionId) {
			res = LicenseList{
				List: []*ComponentLicense{l},
				Op:   NONE,
			}
			break
		}
	}
	return res
}

// GetLicenseEffective will return the right license value
// When license (licenseConcluded) is defined we should use this attribute.
// In all other circumstances we use the default licenseDeclared attribute
// https://github.com/eclipse-disuko/disuko/issues/862
func (ci ComponentInfo) GetLicenseEffective() string {
	res := ci.LicenseDeclared
	if !helper.IsUnasserted(ci.License) {
		res = ci.License
	}
	// #6642: check here also if a decision was made and use!
	if ci.LicenseRuleApplied != nil && !ci.LicenseRuleApplied.PreviewMode {
		res = ci.LicenseRuleApplied.LicenseDecisionId
	}
	return res
}

// GetLicenseAppliedType will return the type of the right license value
// When license (licenseConcluded) is defined we should use this attribute.
// In all other circumstances we use the default licenseDeclared attribute
// https://github.com/eclipse-disuko/disuko/issues/862
func (ci ComponentInfo) GetLicenseAppliedType() LicenseAppliedType {
	if !helper.IsUnasserted(ci.License) {
		return LicenseConcluded
	}
	return LicenseDeclared
}

func (ci ComponentInfo) HasUnknownLicenses() bool {
	for _, license := range ci.GetLicensesEffective().List {
		if !license.Known {
			return true
		}
	}
	return false
}

func (ci ComponentInfo) EffectiveLicensesString() string {
	res := ""

	if ci.ComplexExpression || ci.ContainsBadChars {
		return ci.GetLicenseEffective()
	}

	for _, license := range ci.GetLicensesEffective().List {
		if !license.Known {
			res += license.OrigName
		} else {
			res += license.ReferencedLicense
		}
		res += strings.ToUpper(string(ci.GetLicensesEffective().Op))
	}
	if res == "" {
		return "NOASSERTION"
	}

	return strings.TrimSuffix(res, strings.ToUpper(string(ci.GetLicensesEffective().Op)))
}

func (ci ComponentInfo) WorstFamily() license.FamilyOfLicense {
	effective := ci.GetLicensesEffective()
	if len(effective.List) == 0 {
		return license.NotDeclared
	}
	ranking := map[license.FamilyOfLicense]int{
		license.Permissive:      0,
		license.WeakCopyleft:    1,
		license.StrongCopyleft:  2,
		license.NetworkCopyleft: 3,
		license.NotDeclared:     4,
	}
	res := license.Permissive
	for _, l := range effective.List {
		if ranking[l.LicenseFamily] > ranking[res] {
			res = l.LicenseFamily
		}
	}
	return res
}

func (cis ComponentInfos) FindComponentsByNameFragment(searchFragment string) []ComponentInfo {
	var res []ComponentInfo
	for _, c := range cis {
		if !strings.Contains(c.Name, searchFragment) {
			continue
		}
		res = append(res, c)
	}
	return res
}

func (l *LicenseList) CountOrLinks() int {
	if l.Op != OR {
		return 0
	}
	return len(l.List)
}
