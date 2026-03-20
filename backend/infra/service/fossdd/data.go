// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package fossdd

import (
	"github.com/eclipse-disuko/disuko/domain/job"
	"github.com/eclipse-disuko/disuko/domain/label"
	"github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/project/components"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
)

type subData struct {
	refs        SubProjectRefs
	pr          *project.Project
	schemaLabel *label.Label
	rules       []*ruleWithLics
	lics        []*license.License
	spdx        *project.SpdxFileBase
	compInfos   components.ComponentInfos
	evalRes     *components.EvaluationResult
}

type data struct {
	pr *project.Project

	allPolicyLabels []*label.Label
	allRules        map[string]*ruleWithLics
	allLicenses     map[string]*license.License

	subDatas []subData

	licCache map[string]*license.License
}

type ruleWithLics struct {
	*license.PolicyRules
	allowed []*license.License
	warned  []*license.License
	denied  []*license.License
}

func (g *gen) collectData() {
	g.data.allPolicyLabels = g.service.LabelRepo.FindAllByType(g.rs, label.POLICY)
	g.data.licCache = make(map[string]*license.License)
	g.data.allRules = make(map[string]*ruleWithLics)
	g.data.allLicenses = make(map[string]*license.License)

	g.data.pr = g.service.ProjectRepo.FindByKey(g.rs, g.opts.MainProjectID, false)
	if g.data.pr == nil {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbNotFound), g.opts.MainProjectID+" not found in DB")
	}

	for _, refs := range g.opts.SubProjectsRefs {
		g.collectSubData(refs)
	}
}

func (g *gen) collectSubData(refs SubProjectRefs) {
	g.jobLog.AddEntry(job.Info, "collecting data for sub project %s", refs.ProjectID)
	res := subData{
		refs: refs,
	}

	res.pr = g.service.ProjectRepo.FindByKey(g.rs, refs.ProjectID, false)
	if res.pr == nil {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbNotFound), refs.ProjectID+" not found in DB")
	}

	res.schemaLabel = g.service.LabelRepo.FindByKey(g.rs, res.pr.SchemaLabel, false)
	if res.schemaLabel == nil {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbNotFound), res.pr.SchemaLabel+" not found in DB")
	}

	g.collectSubRules(&res)
	g.collectSpdxData(&res)
	g.collectSubLics(&res)

	g.data.subDatas = append(g.data.subDatas, res)
}

func (g *gen) collectSubRules(subData *subData) {
	g.jobLog.AddEntry(job.Info, "collecting rules for sub project %s", subData.refs.ProjectID)
	rules := g.service.PolicyRuleRepo.FindPolicyRulesForLabel(g.rs, subData.pr.PolicyLabels)
	for _, r := range rules {
		rwl := &ruleWithLics{
			PolicyRules: r,
			allowed:     g.lookupLics(r.ComponentsAllow),
			warned:      g.lookupLics(r.ComponentsWarn),
			denied:      g.lookupLics(r.ComponentsDeny),
		}

		if !r.Auxiliary {
			g.data.allRules[r.Key] = rwl
		}

		subData.rules = append(subData.rules, rwl)
	}
}

func (g *gen) collectSpdxData(subData *subData) {
	if subData.refs.SpdxID == nil || subData.refs.VersionID == nil {
		return
	}
	g.jobLog.AddEntry(job.Info, "collecting spdx data for sub project %s", subData.refs.ProjectID)

	subData.spdx = g.service.SbomListRepo.FindFile(g.rs, *subData.refs.VersionID, *subData.refs.SpdxID)
	if subData.spdx == nil {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbNotFound), *subData.refs.SpdxID+" not found in DB")
	}
	subData.compInfos = g.service.SpdxService.GetComponentInfos(g.rs, subData.pr, *subData.refs.VersionID, subData.spdx)
	var evalRules []*license.PolicyRules
	for _, r := range subData.rules {
		evalRules = append(evalRules, r.PolicyRules)
	}
	policyDecisions := g.service.PolicyDecisionsRepo.FindByKey(g.rs, subData.pr.Key, false)
	isVehicle := g.service.ProjectLabelService.HasVehiclePlatformLabel(g.rs, subData.pr)
	subData.evalRes = subData.compInfos.EvaluatePolicyRules(evalRules, policyDecisions, isVehicle, subData.spdx.Uploaded, subData.spdx.Key)
}

func (g *gen) collectSubLics(subData *subData) {
	if subData.evalRes == nil {
		return
	}
	g.jobLog.AddEntry(job.Info, "collecting licenses for sub project %s", subData.refs.ProjectID)

	licensesDupCheck := make(map[string]struct{})
resLoop:
	for _, compRes := range subData.evalRes.Results {
		if compRes.Unasserted {
			continue
		}

		for _, status := range compRes.Status {
			if status.Auxiliary {
				continue resLoop
			}
		}

		for _, compLic := range compRes.Component.GetLicensesEffective().List {
			lic, ok := g.data.licCache[compLic.ReferencedLicense]
			if !ok {
				lic = g.service.LicenseRepo.FindById(g.rs, compLic.ReferencedLicense)
				if lic == nil {
					exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbNotFound), compLic.ReferencedLicense+" not found in DB")
				}
				g.data.licCache[lic.LicenseId] = lic
			}
			g.data.allLicenses[lic.LicenseId] = lic
			if _, ok = licensesDupCheck[lic.LicenseId]; !ok {
				licensesDupCheck[lic.LicenseId] = struct{}{}
				subData.lics = append(subData.lics, lic)
			}
		}
	}
}

func (g *gen) lookupLics(ids []string) []*license.License {
	var res []*license.License
	for _, id := range ids {
		lic, ok := g.data.licCache[id]
		if !ok {
			lic = g.service.LicenseRepo.FindById(g.rs, id)
			if lic == nil {
				exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbNotFound), id+" not found in DB")
			}
			g.data.licCache[lic.LicenseId] = lic
		}
		res = append(res, lic)

	}
	return res
}
