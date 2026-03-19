// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package checklist

import (
	"slices"
	"strings"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/checklist"
	"github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/project/components"
	"github.com/eclipse-disuko/disuko/domain/reviewremarks"
	"github.com/eclipse-disuko/disuko/helper/exception"
	checklistRepo "github.com/eclipse-disuko/disuko/infra/repository/checklist"
	licRepo "github.com/eclipse-disuko/disuko/infra/repository/license"
	"github.com/eclipse-disuko/disuko/infra/repository/policydecisions"
	policyrulesRepo "github.com/eclipse-disuko/disuko/infra/repository/policyrules"
	reviewRemarksRepo "github.com/eclipse-disuko/disuko/infra/repository/reviewremarks"
	templateRepo "github.com/eclipse-disuko/disuko/infra/repository/reviewtemplates"
	sbomlistRepo "github.com/eclipse-disuko/disuko/infra/repository/sbomlist"
	projectLabelService "github.com/eclipse-disuko/disuko/infra/service/project-label"
	"github.com/eclipse-disuko/disuko/infra/service/scanremarks"
	spdxService "github.com/eclipse-disuko/disuko/infra/service/spdx"
	"github.com/eclipse-disuko/disuko/logy"
)

type Service struct {
	ChecklistRepo       checklistRepo.IChecklistRepository
	TemplateRepo        templateRepo.IReviewTemplateRepository
	SbomListRepo        sbomlistRepo.ISbomListRepository
	PolicyRuleRepo      policyrulesRepo.IPolicyRulesRepository
	LicenseRepo         licRepo.ILicensesRepository
	ReviewRemarkRepo    reviewRemarksRepo.IReviewRemarksRepository
	SpdxService         *spdxService.Service
	ScanRemarksService  *scanremarks.Service
	ProjectLabelService *projectLabelService.ProjectLabelService
	PolicyDecisionsRepo policydecisions.IPolicyDecisionsRepository
}

type execution struct {
	rs          *logy.RequestSession
	service     *Service
	cl          *checklist.Checklist
	pr          *project.Project
	version     *project.ProjectVersion
	spdxBase    *project.SpdxFileBase
	compInfos   components.ComponentInfos
	evalRes     *components.EvaluationResult
	scanRemarks []project.QualityScanRemarks
	executer    string
	rrLookup    map[string]*reviewremarks.ReviewTemplate
	licLookup   map[string]*license.License
	licRepo     licRepo.ILicensesRepository
}

func (s *Service) FindApplicableLists(rs *logy.RequestSession, pr *project.Project) []*checklist.Checklist {
	var (
		all = s.ChecklistRepo.FindAll(rs, false)
		res []*checklist.Checklist
	)
	for _, cl := range all {
		if !cl.Active {
			continue
		}
		if s.checklistApplicable(pr, cl) {
			res = append(res, cl)
		}
	}
	return res
}

func (s *Service) Execute(rs *logy.RequestSession, pr *project.Project, version *project.ProjectVersion, spdxID string, listIds []string, executer string) {
	lists := s.findApplicableByIds(rs, pr, listIds)

	sbomList := s.SbomListRepo.FindByKey(rs, version.Key, false)
	if sbomList == nil || len(sbomList.SpdxFileHistory) == 0 {
		exception.ThrowExceptionBadRequestResponse()
	}
	var spdxBase *project.SpdxFileBase
	for _, spdx := range sbomList.SpdxFileHistory {
		if spdx.Key == spdxID {
			spdxBase = spdx
			break
		}
	}
	if spdxBase == nil {
		exception.ThrowExceptionBadRequestResponse()
	}

	compInfos := s.SpdxService.GetComponentInfos(rs, pr, version.Key, spdxBase)
	rules := s.PolicyRuleRepo.FindPolicyRulesForLabel(rs, pr.PolicyLabels)
	policyDecisions := s.PolicyDecisionsRepo.FindByKey(rs, pr.Key, false)
	isVehicle := s.ProjectLabelService.HasVehiclePlatformLabel(rs, pr)
	evalRes := compInfos.EvaluatePolicyRules(rules, policyDecisions, isVehicle, spdxBase.Uploaded, spdxBase.Key)

	rrLookup := s.templateLookup(rs)
	licLookup := make(map[string]*license.License)

	scanRemarks := s.ScanRemarksService.GetRemarks(rs, pr, &spdxBase.MetaInfo, evalRes)

	var remarks []*reviewremarks.Remark
	for _, list := range lists {
		e := execution{
			rs:          rs,
			service:     s,
			cl:          list,
			pr:          pr,
			version:     version,
			spdxBase:    spdxBase,
			compInfos:   compInfos,
			evalRes:     evalRes,
			scanRemarks: scanRemarks,
			executer:    executer,
			rrLookup:    rrLookup,
			licLookup:   licLookup,
			licRepo:     s.LicenseRepo,
		}
		remarks = append(remarks, e.do()...)
	}

	if len(remarks) == 0 {
		return
	}

	spdxBase.IsInUse = true
	s.SbomListRepo.Update(rs, sbomList)

	rr := s.ReviewRemarkRepo.FindByKey(rs, version.Key, false)
	if rr == nil {
		rr = &reviewremarks.ReviewRemarks{
			RootEntity: domain.NewRootEntityWithKey(version.Key),
			Remarks:    remarks,
		}
		s.ReviewRemarkRepo.Save(rs, rr)
		return
	}
	rr.Remarks = append(rr.Remarks, remarks...)
	s.ReviewRemarkRepo.Update(rs, rr)
}

func (e *execution) retainSBOM() {
}

func (e *execution) do() []*reviewremarks.Remark {
	var res []*reviewremarks.Remark
	for _, item := range e.cl.Items {
		if !e.service.itemApplicable(e.pr, item) {
			continue
		}
		if rr := e.executeItem(item); rr != nil {
			res = append(res, rr)
		}
	}
	return res
}

func (e *execution) executeItem(item checklist.Item) *reviewremarks.Remark {
	switch item.TriggerType {
	case checklist.Default:
		return e.executeDefaultItem(item)
	case checklist.ClassificationAND:
		return e.executeClass(item, and)
	case checklist.ClassificationOR:
		return e.executeClass(item, or)
	case checklist.License:
		return e.executeLic(item)
	case checklist.PolicyStatus:
		return e.executePol(item)
	case checklist.ScanRemark:
		return e.executeSR(item)
	case checklist.ComponentName:
		return e.executeCompName(item)
	}
	return nil
}

func (e *execution) executeCompName(item checklist.Item) *reviewremarks.Remark {
	triggerComps := make(map[string]components.ComponentInfo)
	for _, evalRes := range e.evalRes.Results {
		matches := slices.ContainsFunc(item.ComponentNames, func(name string) bool {
			return strings.Contains(
				strings.ToLower(evalRes.Component.Name),
				strings.ToLower(name),
			)
		})
		if !matches {
			continue
		}
		if len(item.PolicyStatus) > 0 {
			compMatches, _ := e.matchCompPolStatus(evalRes, item.PolicyStatus)
			if !compMatches {
				continue
			}
		}
		triggerComps[evalRes.Component.SpdxId] = *evalRes.Component
	}
	if len(triggerComps) > 0 {
		return e.createLinkedRemark(item, triggerComps, nil)
	}
	return nil
}

func (e *execution) executeSR(item checklist.Item) *reviewremarks.Remark {
	for _, sr := range e.scanRemarks {
		if sr.Status == *item.ScanRemarks {
			return e.createRemark(item)
		}
	}
	return nil
}

func (e *execution) executePol(item checklist.Item) *reviewremarks.Remark {
	var (
		triggerComps = make(map[string]components.ComponentInfo)
		triggerLics  = make(map[string]*license.License)
	)
	for _, evalRes := range e.evalRes.Results {
		compMatches, matchingLics := e.matchCompPolStatus(evalRes, item.PolicyStatus)
		if !compMatches {
			continue
		}
		triggerComps[evalRes.Component.SpdxId] = *evalRes.Component
		for _, lic := range matchingLics {
			triggerLics[lic.LicenseId] = lic
		}

	}
	if len(triggerComps) > 0 || len(triggerLics) > 0 {
		return e.createLinkedRemark(item, triggerComps, triggerLics)
	}
	return nil
}

func (e *execution) matchCompPolStatus(evalRes components.ComponentResult, filter []checklist.PolicyStatusType) (bool, []*license.License) {
	var (
		compMatches  bool
		matchingLics []*license.License
	)
	if evalRes.Unasserted && slices.Contains(filter, checklist.Unasserted) {
		compMatches = true
	}
	if evalRes.Questioned && slices.Contains(filter, checklist.Questioned) {
		compMatches = true
	}
	for _, status := range evalRes.Status {
		var toSearch checklist.PolicyStatusType
		switch status.Type {
		case license.ALLOW:
			toSearch = checklist.Allowed
		case license.DENY:
			toSearch = checklist.Denied
		case license.WARN:
			toSearch = checklist.Warned
		}
		if slices.Contains(filter, toSearch) {
			lic, ok := e.licLookup[status.LicenseMatched]
			if !ok {
				lic = e.licRepo.FindById(e.rs, status.LicenseMatched)
				e.licLookup[lic.LicenseId] = lic
			}
			matchingLics = append(matchingLics, lic)
			compMatches = true

		}
	}
	return compMatches, matchingLics
}

func (e *execution) executeLic(item checklist.Item) *reviewremarks.Remark {
	var (
		triggerComps = make(map[string]components.ComponentInfo)
		triggerLics  = make(map[string]*license.License)
	)
	for _, evalRes := range e.evalRes.Results {
		if len(item.PolicyStatus) > 0 {
			compMatches, _ := e.matchCompPolStatus(evalRes, item.PolicyStatus)
			if !compMatches {
				continue
			}
		}
		for _, l := range evalRes.Component.GetLicensesEffective().List {
			if !l.Known {
				continue
			}
			if slices.Contains(item.LicenseIds, l.ReferencedLicense) {
				lic, ok := e.licLookup[l.ReferencedLicense]
				if !ok {
					lic = e.licRepo.FindById(e.rs, l.ReferencedLicense)
					e.licLookup[lic.LicenseId] = lic
				}
				triggerLics[lic.LicenseId] = lic
				triggerComps[evalRes.Component.SpdxId] = *evalRes.Component
			}
		}
	}
	if len(triggerComps) > 0 || len(triggerLics) > 0 {
		return e.createLinkedRemark(item, triggerComps, triggerLics)
	}
	return nil
}

func (e *execution) executeClass(item checklist.Item, logic logic) *reviewremarks.Remark {
	var (
		triggerComps = make(map[string]components.ComponentInfo)
		triggerLics  = make(map[string]*license.License)
	)
	for _, evalRes := range e.evalRes.Results {
		if len(item.PolicyStatus) > 0 {
			compMatches, _ := e.matchCompPolStatus(evalRes, item.PolicyStatus)
			if !compMatches {
				continue
			}
		}
		trigger := false
		for _, l := range evalRes.Component.GetLicensesEffective().List {
			if !l.Known {
				continue
			}
			lic, ok := e.licLookup[l.ReferencedLicense]
			if !ok {
				lic = e.licRepo.FindById(e.rs, l.ReferencedLicense)
				e.licLookup[lic.LicenseId] = lic
			}
			match := sliceMatch(lic.Meta.ObligationsKeyList, item.Classifications, logic)
			if match {
				trigger = true
				triggerLics[lic.Key] = lic
			}
		}
		if trigger {
			triggerComps[evalRes.Component.SpdxId] = *evalRes.Component
		}
	}
	if len(triggerComps) > 0 || len(triggerLics) > 0 {
		return e.createLinkedRemark(item, triggerComps, triggerLics)
	}
	return nil
}

func (e *execution) executeDefaultItem(item checklist.Item) *reviewremarks.Remark {
	return e.createRemark(item)
}

func (e *execution) createRemark(item checklist.Item) *reviewremarks.Remark {
	temp := e.rrLookup[item.TargetTemplateKey]
	return &reviewremarks.Remark{
		ChildEntity:  domain.NewChildEntity(),
		Author:       e.executer,
		Origin:       "Checklist " + e.cl.Name,
		Title:        temp.Title,
		Level:        temp.Level,
		Description:  temp.Description,
		Status:       reviewremarks.Open,
		SBOMId:       e.spdxBase.Key,
		SBOMName:     e.spdxBase.MetaInfo.Name,
		SBOMUploaded: &e.spdxBase.Created,
	}
}

func (e *execution) createLinkedRemark(item checklist.Item, triggerComps map[string]components.ComponentInfo, triggerLics map[string]*license.License) *reviewremarks.Remark {
	temp := e.rrLookup[item.TargetTemplateKey]
	res := reviewremarks.Remark{
		ChildEntity:  domain.NewChildEntity(),
		Author:       e.executer,
		Origin:       "Checklist " + e.cl.Name,
		Title:        temp.Title,
		Level:        temp.Level,
		Description:  temp.Description,
		Status:       reviewremarks.Open,
		SBOMId:       e.spdxBase.Key,
		SBOMName:     e.spdxBase.MetaInfo.Name,
		SBOMUploaded: &e.spdxBase.Created,
	}
	for _, comp := range triggerComps {
		res.Components = append(res.Components, reviewremarks.ComponentMeta{
			ComponentId:      comp.SpdxId,
			ComponentName:    comp.Name,
			ComponentVersion: comp.Version,
		})
	}
	for _, lic := range triggerLics {
		res.Licenses = append(res.Licenses, reviewremarks.LicenseMeta{
			LicenseId:   lic.LicenseId,
			LicenseName: lic.Name,
		})
	}
	return &res
}

func (s *Service) findApplicableByIds(rs *logy.RequestSession, pr *project.Project, ids []string) []*checklist.Checklist {
	var res []*checklist.Checklist
	for _, id := range ids {
		cl := s.ChecklistRepo.FindByKey(rs, id, false)
		if cl == nil {
			exception.ThrowExceptionBadRequestResponse()
		}
		if !cl.Active {
			exception.ThrowExceptionBadRequestResponse()
		}
		if !s.checklistApplicable(pr, cl) {
			exception.ThrowExceptionBadRequestResponse()
		}
		res = append(res, cl)
	}
	return res
}

func (s *Service) checklistApplicable(pr *project.Project, cl *checklist.Checklist) bool {
	return sliceMatch(pr.PolicyLabels, cl.PolicyLabels, or)
}

func (s *Service) itemApplicable(pr *project.Project, item checklist.Item) bool {
	if len(item.PolicyLabels) == 0 {
		return true
	}
	return sliceMatch(pr.PolicyLabels, item.PolicyLabels, or)
}

type logic int

const (
	and logic = iota
	or
)

func sliceMatch[S ~[]E, E comparable](all S, find S, l logic) bool {
	for _, f := range find {
		contains := slices.Contains(all, f)
		if l == or && contains {
			return true
		}
		if l == and && !contains {
			return false
		}
	}
	return l == and
}

func (s *Service) templateLookup(rs *logy.RequestSession) map[string]*reviewremarks.ReviewTemplate {
	templates := s.TemplateRepo.FindAll(rs, false)
	res := make(map[string]*reviewremarks.ReviewTemplate)
	for _, t := range templates {
		res[t.Key] = t
	}
	return res
}
