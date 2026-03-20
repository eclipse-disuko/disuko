// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package analytics

import (
	"slices"

	"github.com/eclipse-disuko/disuko/domain/analytics"
	"github.com/eclipse-disuko/disuko/domain/label"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/project/components"
	sbomlist2 "github.com/eclipse-disuko/disuko/domain/project/sbomlist"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/infra/repository/labels"
	"github.com/eclipse-disuko/disuko/infra/repository/licenserules"
	"github.com/eclipse-disuko/disuko/infra/repository/policydecisions"
	projectLabelService "github.com/eclipse-disuko/disuko/infra/service/project-label"

	"github.com/eclipse-disuko/disuko/infra/repository/department"
	"github.com/eclipse-disuko/disuko/infra/repository/license"
	"github.com/eclipse-disuko/disuko/infra/repository/policyrules"
	project2 "github.com/eclipse-disuko/disuko/infra/repository/project"
	"github.com/eclipse-disuko/disuko/infra/repository/sbomlist"
	projectService "github.com/eclipse-disuko/disuko/infra/service/project"
	"github.com/eclipse-disuko/disuko/infra/service/spdx"
	"github.com/eclipse-disuko/disuko/logy"
)

type Analytics struct {
	ProjectRepository    project2.IProjectRepository
	LicenseRepository    license.ILicensesRepository
	PolicyRuleRepository policyrules.IPolicyRulesRepository
	DepartmentRepo       department.IDepartmentRepository
	SpdxService          *spdx.Service

	SbomListrepository     sbomlist.ISbomListRepository
	Handler                DataHandler
	LicenseRulesRepository licenserules.ILicenseRulesRepository
	LabelRepository        labels.ILabelRepository
	ProjectLabelService    *projectLabelService.ProjectLabelService
	PolicyDecisionsRepo    policydecisions.IPolicyDecisionsRepository
}

type SpdxAddedOptions struct {
	rs       *logy.RequestSession
	project  *project.Project
	parent   *project.Project
	version  *project.ProjectVersion
	evalRes  *components.EvaluationResult
	spdxFile *project.SpdxFileBase
}

type SearchOptions struct {
	Rs          *logy.RequestSession
	Component   string
	License     string
	ProjectKeys []string
	Offset      int
	Limit       int
	SortCol     string
	Asc         bool
}

type DataHandler interface {
	HandleSpdxAdded(SpdxAddedOptions)
	HandleSpdxDeleted(*logy.RequestSession, string)
	Reset()
	Occurrences(*logy.RequestSession) []*analytics.Occurrence
	HandleSearch(SearchOptions) analytics.ResponseAnalyticsSearch
	HandleComponentSearch(*logy.RequestSession, string, bool) analytics.ResponseComponentsSearch
	HandleLicenseSearch(*logy.RequestSession, string, bool) analytics.ResponseLicensesSearch
	HandleLicenseIdDeleted(*logy.RequestSession, string)
	HandleLicenseIdAdded(*logy.RequestSession, string, string)
	HandleCompanyChanged(*logy.RequestSession, string, string)
	HandleResponsibleChanged(*logy.RequestSession, string, string)
}

func (a *Analytics) ExportSPDX(requestSession *logy.RequestSession, pr *project.Project, version *project.ProjectVersion, spdx *project.SpdxFileBase) {
	compInfo := a.SpdxService.GetComponentInfos(requestSession, pr, version.Key, spdx)
	rules := a.PolicyRuleRepository.FindPolicyRulesForLabel(requestSession, pr.PolicyLabels)
	policyDecisions := a.PolicyDecisionsRepo.FindByKey(requestSession, pr.Key, false)
	isVehicle := a.ProjectLabelService.HasVehiclePlatformLabel(requestSession, pr)
	evalRes := compInfo.EvaluatePolicyRules(rules, policyDecisions, isVehicle, spdx.Uploaded, spdx.Key)
	var parent *project.Project
	if pr.Parent != "" {
		parent = a.ProjectRepository.FindByKey(requestSession, pr.Parent, false)
	}
	a.Handler.HandleSpdxAdded(SpdxAddedOptions{
		rs:       requestSession,
		project:  pr,
		parent:   parent,
		version:  version,
		evalRes:  evalRes,
		spdxFile: spdx,
	})
}

func (a *Analytics) DeleteProject(requestSession *logy.RequestSession, project *project.Project) {
	for _, v := range project.Versions {
		if v.Deleted {
			continue
		}
		a.DeleteVersion(requestSession, project, v)
	}
}

func (a *Analytics) DeleteVersion(requestSession *logy.RequestSession, project *project.Project, version *project.ProjectVersion) {
	versionSBoms := a.getSbomList(requestSession, version.Key)
	if versionSBoms == nil || len(versionSBoms.SpdxFileHistory) == 0 || version.Deleted {
		return
	}
	spdx := versionSBoms.SpdxFileHistory.GetLatest()
	go a.Handler.HandleSpdxDeleted(requestSession, spdx.Key)
}

func (a *Analytics) Reinitialise(requestSession *logy.RequestSession) {
	logy.Infof(requestSession, "reinitialising analytics")

	a.Handler.Reset()
	projects := a.ProjectRepository.FindAllKeys(requestSession)
	for _, pK := range projects {
		p := a.ProjectRepository.FindByKey(requestSession, pK, false)
		if hasDummyLabel(requestSession, p, a.LabelRepository) {
			continue
		}
		exception.TryCatch(func() {
			a.exportProject(requestSession, p, repos(a))
		}, func(ex exception.Exception) {
			logy.Errorw(requestSession, "Error %s", ex.ErrorMessage)
			exception.LogException(requestSession, ex)
		})

	}

	logy.Infof(requestSession, "reinitialised analytics")
}

func hasDummyLabel(requestSession *logy.RequestSession, currentProject *project.Project, labelRepository labels.ILabelRepository) bool {
	if currentProject == nil {
		return false
	}
	if currentProject.ProjectLabels == nil {
		return false
	}
	dummyLabel := labelRepository.FindByNameAndType(requestSession, label.DUMMY, label.PROJECT)
	if dummyLabel != nil {
		return slices.Contains(currentProject.ProjectLabels, dummyLabel.GetKey())
	} else {
		return false
	}
}

func (a *Analytics) Search(options SearchOptions) analytics.ResponseAnalyticsSearch {
	res := a.Handler.HandleSearch(options)
	for i, e := range res.Items {
		if e.OwnerDeptId == "" {
			continue
		}
		dept := a.DepartmentRepo.GetByDeptId(options.Rs, e.OwnerDeptId)
		if dept == nil {
			res.Items[i].OwnerDeptMissing = true
			continue
		}
		res.Items[i].OwnerDept = dept.OrgAbbreviation + " " + dept.DescriptionEnglish + " [" + dept.Key + "]"
		res.Items[i].OwnerCompany = dept.CompanyName + " [" + dept.CompanyCode + "]"
	}
	return res
}

func repos(a *Analytics) projectService.RepositoryHolder {
	return projectService.RepositoryHolder{
		ProjectRepository:      a.ProjectRepository,
		LicenseRepository:      a.LicenseRepository,
		SBOMListRepository:     a.SbomListrepository,
		LicenseRulesRepository: a.LicenseRulesRepository,
	}
}

func (a *Analytics) exportProject(requestSession *logy.RequestSession, pr *project.Project, repos projectService.RepositoryHolder) {
	if pr.Deleted {
		return
	}
	for _, version := range pr.Versions {
		if version.Deleted {
			continue
		}
		versionSBoms := a.getSbomList(requestSession, version.Key)
		if versionSBoms == nil || len(versionSBoms.SpdxFileHistory) == 0 || version.Deleted {
			continue
		}
		var parent *project.Project
		if pr.Parent != "" {
			parent = a.ProjectRepository.FindByKey(requestSession, pr.Parent, false)
		}
		spdx := versionSBoms.SpdxFileHistory.GetLatest()
		compInfo := a.SpdxService.GetComponentInfos(requestSession, pr, version.Key, spdx)
		rules := a.PolicyRuleRepository.FindPolicyRulesForLabel(requestSession, pr.PolicyLabels)
		policyDecisions := a.PolicyDecisionsRepo.FindByKey(requestSession, pr.Key, false)
		isVehicle := a.ProjectLabelService.HasVehiclePlatformLabel(requestSession, pr)
		evalRes := compInfo.EvaluatePolicyRules(rules, policyDecisions, isVehicle, spdx.Uploaded, spdx.Key)
		a.Handler.HandleSpdxAdded(SpdxAddedOptions{
			rs:       requestSession,
			project:  pr,
			parent:   parent,
			version:  version,
			evalRes:  evalRes,
			spdxFile: spdx,
		})
	}
}

func (a *Analytics) getSbomList(requestSession *logy.RequestSession, key string) *sbomlist2.SbomList {
	sbomList := a.SbomListrepository.FindByKey(requestSession, key, false)
	if sbomList == nil {
		return nil
	}
	return sbomList
}
