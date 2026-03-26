// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package components

import (
	"sort"
	"strings"
	"time"

	"github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/domain/policydecisions"
	"github.com/eclipse-disuko/disuko/helper/message"
)

type EvaluationResult struct {
	Results []ComponentResult
	Stats   ComponentStats
}

type ComponentResult struct {
	Component         *ComponentInfo
	Status            []*PolicyRuleStatus
	Unmatched         []*UnmatchedLicense
	Questioned        bool
	Unasserted        bool
	ContainsAuxiliary bool
}

type UnmatchedLicense struct {
	OrigName       string
	ReferencedName string
	Known          bool
}
type PolicyRuleStatus struct {
	Key                        string
	Name                       string
	LicenseMatched             string
	Type                       license.ListType
	Used                       bool
	Description                string
	Auxiliary                  bool
	IsDecisionMade             bool
	CanMakeWarnedDecision      bool
	CanMakeDeniedDecision      bool
	DeniedDecisionDeniedReason string
}

type ComponentStats struct {
	Total       int
	Allowed     int
	Warned      int
	Denied      int
	Questioned  int
	NoAssertion int
}

type LicenseFamilyStats struct {
	Total           int
	NetworkCopyLeft int
	StrongCopyLeft  int
	WeakCopyLeft    int
	Permissive      int
	Other           int
}

type ReviewRemarkStats struct {
	Total                  int
	Acceptable             int
	AcceptableAfterChanges int
	NotAcceptable          int
}

type ScanRemarkStats struct {
	Total       int
	Information int
	Warning     int
	Problem     int
}

type ScanRemarkTypeStats struct {
	Total                  int    `json:"total"`
	MissingCopyrights      int    `json:"missingCopyrights"`
	MissingCopyrightsLevel string `json:"missingCopyrightsLevel"`
	MalformedCopyrights    int    `json:"malformedCopyrights"`
}

type NotChartFossLicenseStats struct {
	Total int `json:"total"`
}

type LicenseRemarkStats struct {
	Total       int
	Information int
	Warning     int
	Alarm       int
}

type InApproval struct {
	IsInApproval bool
	ApprovalGuid string
	Status       string
}

type GeneralStats struct {
	SBOMDelivered  bool
	SourceUploaded bool
	ReviewRemark   ReviewRemarkStats
}

type SBOMStats struct {
	PolicyState         ComponentStats
	LicenseFamily       LicenseFamilyStats
	ScanRemark          ScanRemarkStats
	LicenseRemark       LicenseRemarkStats
	ApprovalInfo        InApproval
	ScanRemarkType      ScanRemarkTypeStats      `json:"scanRemarkType"`
	NotChartFossLicense NotChartFossLicenseStats `json:"notChartFossLicense"`
}

var weightedStatus = map[license.ListType]int{
	license.ALLOW: 0,
	license.WARN:  1,
	license.DENY:  2,
}

func (s *ComponentStats) Add(pr license.PolicyResult) {
	s.Total++
	switch pr {
	case license.ALLOWED:
		s.Allowed++
	case license.WARNED:
		s.Warned++
	case license.DENIED:
		s.Denied++
	case license.QUESTIONED:
		s.Questioned++
	case license.UNASSERTED:
		s.NoAssertion++
	}
}

func (s *ComponentStats) AddStats(a ComponentStats) {
	s.Allowed += a.Allowed
	s.Warned += a.Warned
	s.Denied += a.Denied
	s.Questioned += a.Questioned
	s.NoAssertion += a.NoAssertion
	s.Total += a.Total
}

func (entity *ComponentResult) GetUsedPolicyRule() (string, string) {
	if entity.Unasserted {
		return "noassertion", ""
	}
	if entity.Questioned {
		return "questioned", ""
	}
	if len(entity.Status) > 0 {
		for _, status := range entity.Status {
			if status.Used {
				return string(status.Type), status.Name
			}
		}
	}
	return "", ""
}

func (cmpRes *ComponentResult) applyPolicyDecision(policyDecisions *policydecisions.PolicyDecisions, isVehicle bool, sbomUpload *time.Time, sbomKey string) {
	if policyDecisions == nil {
		return
	}
	compLicenseExpression := cmpRes.Component.EffectiveLicensesString()
	for _, prDecision := range policyDecisions.Decisions {
		cmpNameMatches := strings.EqualFold(prDecision.ComponentName, cmpRes.Component.Name)
		licExprMatches := strings.EqualFold(prDecision.LicenseExpression, compLicenseExpression)
		versionMatches := !isVehicle || prDecision.ComponentVersion == cmpRes.Component.Version

		cmpNameAndLicExprAndVersionMatches := cmpNameMatches && licExprMatches && versionMatches

		if !cmpNameAndLicExprAndVersionMatches || sbomUpload == nil {
			continue
		}
		for i := 0; i < len(cmpRes.Status); i++ {
			prStatus := cmpRes.Status[i]

			licMatches := strings.EqualFold(prDecision.LicenseId, prStatus.LicenseMatched)
			prStatusMatches := strings.EqualFold(prDecision.PolicyEvaluated, string(prStatus.Type))
			prKeyMatches := prDecision.PolicyId == prStatus.Key

			if !(licMatches && prStatusMatches && prKeyMatches) {
				continue
			}

			sbomMatches := prDecision.SBOMId == sbomKey
			activeOK := (prDecision.Active && (*sbomUpload).After(prDecision.Created)) || sbomMatches
			previewOK := prDecision.Active && (*sbomUpload).Before(prDecision.Created) && !sbomMatches
			cancelledOK := !prDecision.Active && (*sbomUpload).After(prDecision.Created) && (*sbomUpload).Before(prDecision.Updated)

			if activeOK || previewOK || cancelledOK {
				if previewOK {
					prDecision.PreviewMode = true
				} else {
					prStatus.Type = license.ListType(prDecision.PolicyDecision)
				}
				cmpRes.Component.PolicyDecisionsApplied = append(cmpRes.Component.PolicyDecisionsApplied, prDecision)
				prStatus.IsDecisionMade = true
				prStatus.CanMakeWarnedDecision = false
				prStatus.CanMakeDeniedDecision = false
				break
			}
		}
	}
}

func (cmpRes *ComponentResult) evaluateComponentPolicyRules(rules []*license.PolicyRules) bool {
	containsUnmatchedLicense := false
	for _, l := range cmpRes.Component.GetLicensesEffective().List {
		matchingRule := 0
		for _, policy := range rules {
			if l.Known {
				if cmpRes.processLicense(l.ReferencedLicense, l.ApprovalState, policy) {
					matchingRule++
				}
			}
		}
		if matchingRule == 0 {
			cmpRes.Unmatched = append(cmpRes.Unmatched, &UnmatchedLicense{
				OrigName:       l.OrigName,
				ReferencedName: l.ReferencedLicense,
				Known:          l.Known,
			})
			containsUnmatchedLicense = true
		}
	}
	return containsUnmatchedLicense
}

func (cis ComponentInfos) EvaluatePolicyRules(rules []*license.PolicyRules, policyDecisions *policydecisions.PolicyDecisions, isVehicle bool, sbomUpload *time.Time, sbomKey string) *EvaluationResult {
	var res EvaluationResult

	for i := 0; i < len(cis); i++ {
		compRes := ComponentResult{
			Component: &cis[i],
		}

		if compRes.Component.Type == ROOT && len(compRes.Component.GetLicensesEffective().List) == 0 {
			res.Results = append(res.Results, compRes)
			res.Stats.Total++
			continue
		}

		containsUnmatchedLicense := compRes.evaluateComponentPolicyRules(rules)

		// first apply policy rule decision
		compRes.applyPolicyDecision(policyDecisions, isVehicle, sbomUpload, sbomKey)
		// then calculate worse used
		worseUsed := compRes.useWorse() // allow, warn or deny

		if compRes.containsDenied() {
			res.Stats.Add(worseUsed)
			sort.Slice(compRes.Status, func(i, j int) bool {
				return weightedStatus[compRes.Status[i].Type] > weightedStatus[compRes.Status[j].Type]
			})
			res.Results = append(res.Results, compRes)
			continue
		}
		if compRes.Component.HasUnknownLicenses() || len(compRes.Component.GetLicensesEffective().List) == 0 || containsUnmatchedLicense {
			compRes.Unasserted = true
			res.Stats.Add(license.UNASSERTED)
		} else {
			if compRes.Component.GetLicensesEffective().Op == OR && compRes.containsAllowedOnly() {
				compRes.Questioned = true
				compRes.Unasserted = false
				res.Stats.Add(license.QUESTIONED)
			} else {
				compRes.Unasserted = false
				res.Stats.Add(worseUsed)
			}
		}
		sort.Slice(compRes.Status, func(i, j int) bool {
			return weightedStatus[compRes.Status[i].Type] > weightedStatus[compRes.Status[j].Type]
		})
		res.Results = append(res.Results, compRes)
	}
	return &res
}

func (c *ComponentResult) containsAllowedOnly() bool {
	for _, res := range c.Status {
		if res.Type != license.ALLOW {
			return false
		}
	}
	return len(c.Status) != 0
}

func (cmpRes *ComponentResult) useWorse() license.PolicyResult {
	worseCase := license.ALLOW

	for _, status := range cmpRes.Status {
		if weightedStatus[status.Type] > weightedStatus[worseCase] {
			worseCase = status.Type
		}
	}

	for i := 0; i < len(cmpRes.Status); i++ {
		if cmpRes.Status[i].Type == worseCase {
			cmpRes.Status[i].Used = true
		}
	}

	return license.PolicyResult(weightedStatus[worseCase])
}

func (cmpRes *ComponentResult) processLicense(licenseName string, approvalState license.ApprovalStatus, rules *license.PolicyRules) bool {
	status := &PolicyRuleStatus{
		Key:                        rules.Key,
		Name:                       rules.Name,
		Description:                rules.Description,
		Used:                       false,
		LicenseMatched:             licenseName,
		Auxiliary:                  rules.Auxiliary,
		IsDecisionMade:             false,
		CanMakeWarnedDecision:      false,
		CanMakeDeniedDecision:      false,
		DeniedDecisionDeniedReason: "",
	}

	for _, allow := range rules.ComponentsAllow {
		if licenseName == allow {
			status.Type = license.ALLOW
			cmpRes.Status = append(cmpRes.Status, status)
			return true
		}
	}
	for _, deny := range rules.ComponentsDeny {
		if licenseName == deny {
			status.Type = license.DENY
			status.CanMakeDeniedDecision = true
			if approvalState == license.Forbidden {
				status.DeniedDecisionDeniedReason = message.PolicyDecisionDeniedForbiddenLicense
			}
			cmpRes.Status = append(cmpRes.Status, status)
			return true
		}
	}
	for _, warn := range rules.ComponentsWarn {
		if licenseName == warn {
			status.Type = license.WARN
			status.CanMakeWarnedDecision = true
			cmpRes.Status = append(cmpRes.Status, status)
			return true
		}
	}

	return false
}

func (r *ComponentResult) ContainsUnmatchedLicense(id string) bool {
	for _, l := range r.Unmatched {
		if l.OrigName == id {
			return true
		}
	}
	return false
}

func (r *ComponentResult) containsDenied() bool {
	for _, status := range r.Status {
		if status.Type == license.DENY {
			return true
		}
	}
	return false
}
