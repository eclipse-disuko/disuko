// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package confirmtransfer

import (
	domainCompare "github.com/eclipse-disuko/disuko/domain/compare"
	"github.com/eclipse-disuko/disuko/domain/label"
	"github.com/eclipse-disuko/disuko/domain/overallreview"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/project/sbomlist"
	"github.com/eclipse-disuko/disuko/infra/repository/policydecisions"
	"github.com/eclipse-disuko/disuko/infra/repository/policyrules"
	sbomlistRepo "github.com/eclipse-disuko/disuko/infra/repository/sbomlist"
	"github.com/eclipse-disuko/disuko/infra/service/compare"
	projectService "github.com/eclipse-disuko/disuko/infra/service/project"
	projectLabelService "github.com/eclipse-disuko/disuko/infra/service/project-label"
	"github.com/eclipse-disuko/disuko/infra/service/spdx"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/eclipse-disuko/disuko/observermngmt"
)

type ConfirmTransfer struct {
	SbomListRepo         sbomlistRepo.ISbomListRepository
	PolicyRuleRepo       policyrules.IPolicyRulesRepository
	PolicyDecisionsRepo  policydecisions.IPolicyDecisionsRepository
	ProjectLabelService  *projectLabelService.ProjectLabelService
	SpdxService          *spdx.Service
	OverallReviewService *projectService.OverallReviewService
}

func Init(
	sbomListRepo sbomlistRepo.ISbomListRepository,
	policyRuleRepo policyrules.IPolicyRulesRepository,
	policyDecisionsRepo policydecisions.IPolicyDecisionsRepository,
	projectLabelSvc *projectLabelService.ProjectLabelService,
	spdxService *spdx.Service,
	overallReviewService *projectService.OverallReviewService,
) *ConfirmTransfer {
	return &ConfirmTransfer{
		SbomListRepo:         sbomListRepo,
		PolicyRuleRepo:       policyRuleRepo,
		PolicyDecisionsRepo:  policyDecisionsRepo,
		ProjectLabelService:  projectLabelSvc,
		SpdxService:          spdxService,
		OverallReviewService: overallReviewService,
	}
}

func (c *ConfirmTransfer) RegisterHandlers() {
	observermngmt.RegisterHandler(observermngmt.SpdxAdded, c.OnSpdxAdded)
}

func (c *ConfirmTransfer) OnSpdxAdded(_ observermngmt.EventId, arg interface{}) {
	data, ok := arg.(observermngmt.SpdxData)
	if !ok {
		return
	}
	rs := data.RequestSession

	if !c.ProjectLabelService.HasLabelInGroupOrProject(rs, data.Project, label.OFFBOARD, label.POLICY) {
		return
	}

	sbomList := c.SbomListRepo.FindByKey(rs, data.Version.Key, false)
	if sbomList == nil {
		return
	}

	previous := findPreviousUpload(sbomList.SpdxFileHistory, data.SpdxFile.Key)
	if previous == nil || previous.OverallReview == nil || previous.OverallReview.State != overallreview.Audited {
		return
	}

	if c.hasDifferences(rs, data.Project, data.Version, previous, data.SpdxFile) {
		return
	}

	c.OverallReviewService.AddToProject(
		rs,
		data.Project,
		data.Version,
		previous.OverallReview.Creator,
		overallreview.Audited,
		previous.OverallReview.Comment,
		data.SpdxFile.Key,
	)

	sbomList = c.SbomListRepo.FindByKey(rs, data.Version.Key, false)
	n := sbomList.SpdxFileHistory.GetByKey(data.SpdxFile.Key)
	data.SpdxFile.OverallReview = n.OverallReview

	logy.Infof(rs, "confirmtransfer: transferred audited overall review of sbom %s to newly uploaded sbom %s in version %s", previous.Key, data.SpdxFile.Key, data.Version.Key)
}

func findPreviousUpload(history sbomlist.Histories, currentKey string) *project.SpdxFileBase {
	var previous *project.SpdxFileBase
	for _, s := range history {
		if s.Key == currentKey || s.Uploaded == nil {
			continue
		}
		if previous == nil || s.Uploaded.After(*previous.Uploaded) {
			previous = s
		}
	}
	return previous
}

func (c *ConfirmTransfer) hasDifferences(rs *logy.RequestSession, pr *project.Project, version *project.ProjectVersion, oldSpdx, newSpdx *project.SpdxFileBase) bool {
	rules := c.PolicyRuleRepo.FindPolicyRulesForLabel(rs, pr.PolicyLabels)
	policyDecisions := c.PolicyDecisionsRepo.FindByKey(rs, pr.Key, false)
	isVehicle := c.ProjectLabelService.HasVehiclePlatformLabel(rs, pr)

	oldComps := c.SpdxService.GetComponentInfos(rs, pr, version.Key, oldSpdx)
	evalOld := oldComps.EvaluatePolicyRules(rules, policyDecisions, isVehicle, oldSpdx.Uploaded, oldSpdx.Key)

	newComps := c.SpdxService.GetComponentInfos(rs, pr, version.Key, newSpdx)
	evalNew := newComps.EvaluatePolicyRules(rules, policyDecisions, isVehicle, newSpdx.Uploaded, newSpdx.Key)

	diffResult := compare.MultiCompareSpdxFiles(evalOld, evalNew, false)
	for _, d := range diffResult {
		if d.DiffType != domainCompare.UNCHANGED {
			return true
		}
	}
	return false
}
