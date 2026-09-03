// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package report

import (
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain/approval"
	"github.com/eclipse-disuko/disuko/domain/label"
	license2 "github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/domain/overallreview"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/project/components"
	"github.com/eclipse-disuko/disuko/domain/report"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/logy"
)

type generation struct {
	rs      *logy.RequestSession
	service *Service

	customIdNames []string
	dummyLabel    *label.Label

	prCache  prCache
	licCache licCache
	oblCache oblCache
}

func newGeneration(rs *logy.RequestSession, s *Service) *generation {
	g := &generation{
		rs:       rs,
		service:  s,
		prCache:  make(prCache),
		licCache: make(licCache),
		oblCache: make(oblCache),
	}
	g.customIdNames = g.loadCustomIdNames()
	g.dummyLabel = s.RepoLabel.FindByNameAndType(rs, label.DUMMY, label.PROJECT)
	return g
}

func (g *generation) loadCustomIdNames() []string {
	var customIdNames []string
	cids := g.service.RepoCustomId.FindAll(g.rs, false)
	for _, cid := range cids {
		customIdNames = append(customIdNames, cid.Key)
	}
	return customIdNames
}

func (g *generation) findProject(key string) *project.Project {
	pr, ok := g.prCache[key]
	if !ok {
		pr = g.service.Repo.FindByKey(g.rs, key, false)
		g.prCache[key] = pr
	}
	return pr
}

func (g *generation) buildProjects() []report.Project {
	var projects []report.Project
	projectKeys := g.service.Repo.FindAllKeys(g.rs)
	for _, projectKey := range projectKeys {
		pr := g.findProject(projectKey)
		if pr == nil {
			continue
		}
		projects = append(projects, g.project(pr, hasDummyLabel(pr, g.dummyLabel)))
	}
	return projects
}

func (g *generation) project(pr *project.Project, isDummy bool) report.Project {
	res := report.Project{}
	g.fillBasicProjectInfo(pr, &res, isDummy)
	g.fillParentAndSupplierInfo(pr, &res)
	g.fillLabelsAndTags(pr, &res)
	g.fillLicenseDecisionRuleStats(pr, &res)
	g.fillPolicyDecisionRuleStats(pr, &res)
	g.fillDeniedPolicyDecisionStats(pr, &res)
	g.fillReviewStats(pr, &res)
	g.fillSourceStats(pr, &res)
	g.fillTokenStats(pr, &res)
	g.fillCustomIds(pr, &res)
	g.fillSbomStats(pr, &res)
	g.fillApprovalStats(pr, &res)
	return res
}

func (g *generation) fillBasicProjectInfo(pr *project.Project, res *report.Project, isDummy bool) {
	res.Name = pr.Name
	res.Group = renderBool(pr.IsGroup)
	res.NonFoss = renderBool(pr.IsNoFoss)
	res.IsDummy = renderBool(isDummy)
	res.Status = string(pr.Status)
	res.Guid = pr.Key
	res.Created = pr.Created.String()
	res.Updated = pr.Updated.String()
	res.Link = conf.Config.Server.DisukoHost + "/#/dashboard/"
	res.Subscribers = strconv.Itoa(getSbomSubscribersCount(pr))
	if pr.IsGroup {
		res.Link += "groups/" + pr.Key
	} else {
		res.Link += "projects/" + pr.Key
	}
}

func (g *generation) fillParentAndSupplierInfo(pr *project.Project, res *report.Project) {
	depPrj := pr
	if pr.Parent != "" {
		depPrj = g.findProject(pr.Parent)
		if depPrj == nil {
			logy.Warnf(g.rs, "parent project %s for %s not found anymore", pr.Parent, pr.Key)
			return
		}
	}
	if depPrj.CustomerMeta.DeptId != "" {
		dep := g.service.RepoDept.GetByDeptId(g.rs, depPrj.CustomerMeta.DeptId)
		if dep == nil {
			res.OwnerCompanyName = "deleted in the meantime"
			res.OwnerCompanyId = "deleted in the meantime"
			res.OwnerDepartmentId = "deleted in the meantime"
			res.OwnerDepartmentTitle = "deleted in the meantime"
			res.OwnerDepartmentAbbreviation = "deleted in the meantime"
		} else {
			res.OwnerCompanyName = dep.CompanyName
			res.OwnerCompanyId = dep.CompanyCode
			res.OwnerDepartmentId = dep.Key
			res.OwnerDepartmentTitle = dep.DescriptionEnglish
			res.OwnerDepartmentAbbreviation = dep.OrgAbbreviation
		}
	}
	if depPrj.DocumentMeta.SupplierDeptId != "" && !depPrj.SupplierExtraData.External {
		dep := g.service.RepoDept.GetByDeptId(g.rs, depPrj.DocumentMeta.SupplierDeptId)
		if dep == nil {
			res.SupplierCompanyName = "deleted in the meantime"
			res.SupplierCompanyId = "deleted in the meantime"
			res.SupplierDepartmentId = "deleted in the meantime"
			res.SupplierDepartmentTitle = "deleted in the meantime"
			res.SupplierDepartmentAbbreviation = "deleted in the meantime"
		} else {
			res.SupplierCompanyName = dep.CompanyName
			res.SupplierCompanyId = dep.CompanyCode
			res.SupplierDepartmentId = dep.Key
			res.SupplierDepartmentTitle = dep.DescriptionEnglish
			res.SupplierDepartmentAbbreviation = dep.OrgAbbreviation
		}
	} else if depPrj.SupplierExtraData.External && depPrj.DocumentMeta.SupplierName != "" {
		res.SupplierCompanyName = depPrj.DocumentMeta.SupplierName
	}
	res.SupplierDepartmentExternal = renderBool(depPrj.SupplierExtraData.External)
	responsible := pr.ProjectResponsible()
	if responsible != nil {
		res.ProjectResponsibleUserid = responsible.UserId
		// TODO: maybe cache that too?
		user := g.service.RepoUser.FindByUserId(g.rs, responsible.UserId)
		if user != nil {
			res.ProjectResponsibleEmail = user.Email
			res.ProjectResponsibleFullName = user.Lastname + "," + user.Forename
		}
	}
	res.ApplicationId = pr.ApplicationMeta.Id
	res.ApplicationSecondaryId = pr.ApplicationMeta.SecondaryId
	res.ApplicationName = pr.ApplicationMeta.Name
	res.GroupId = pr.Parent
	res.GroupName = pr.ParentName
}

func (g *generation) fillLicenseDecisionRuleStats(pr *project.Project, res *report.Project) {
	prs := g.childProjectsOrSelf(pr)
	var (
		active   int
		inactive int
	)
	for _, iP := range prs {
		licenseRules := g.service.RepoLic.FindByKey(g.rs, iP.Key, false)
		if licenseRules == nil {
			continue
		}
		for _, lr := range licenseRules.Rules {
			if lr.Active {
				active++
			} else {
				inactive++
			}
		}

	}
	res.ActiveLicenseDecisionRules = strconv.Itoa(active)
	res.InactiveLicenseDecisionRules = strconv.Itoa(inactive)
}

func (g *generation) fillPolicyDecisionRuleStats(pr *project.Project, res *report.Project) {
	prs := g.childProjectsOrSelf(pr)
	var (
		active   int
		inactive int
	)
	for _, iP := range prs {
		policyRules := g.service.PolicyDecisionsRepo.FindByKey(g.rs, iP.Key, false)
		if policyRules == nil {
			continue
		}
		for _, lr := range policyRules.Decisions {
			if lr.Active {
				active++
			} else {
				inactive++
			}
		}

	}
	res.ActivePolicyDecisionRules = strconv.Itoa(active)
	res.InactivePolicyDecisionRules = strconv.Itoa(inactive)
}

func (g *generation) fillDeniedPolicyDecisionStats(pr *project.Project, res *report.Project) {
	prs := g.childProjectsOrSelf(pr)
	var (
		active   int
		inactive int
	)
	for _, iP := range prs {
		policyRules := g.service.PolicyDecisionsRepo.FindByKey(g.rs, iP.Key, false)
		if policyRules == nil {
			continue
		}
		for _, lr := range policyRules.Decisions {
			evaluated := lr.PolicyEvaluated
			if strings.EqualFold(evaluated, string(license2.DENY)) {
				if lr.Active {
					active++
				} else {
					inactive++
				}
			}
		}

	}
	res.ActiveDeniedPolicyDecision = strconv.Itoa(active)
	res.InactiveDeniedPolicyDecision = strconv.Itoa(inactive)
}

func (g *generation) childProjectsOrSelf(pr *project.Project) []*project.Project {
	if !pr.IsGroup {
		return []*project.Project{pr}
	}
	prs := make([]*project.Project, 0)
	for _, ck := range pr.Children {
		child := g.findProject(ck)
		if child == nil || child.Deleted {
			continue
		}
		prs = append(prs, child)
	}
	return prs
}

func (g *generation) fillLabelsAndTags(pr *project.Project, res *report.Project) {
	if pr.SchemaLabel != "" {
		l := g.service.RepoLabel.FindByKey(g.rs, pr.SchemaLabel, false)
		if l != nil {
			res.SchemaLabel = l.Name
		}
	}
	for _, k := range pr.PolicyLabels {
		l := g.service.RepoLabel.FindByKey(g.rs, k, false)
		if l != nil {
			res.PolicyLabels += l.Name + ","
		}
	}
	res.PolicyLabels = strings.TrimSuffix(res.PolicyLabels, ",")
	if pr.FreeLabels != nil {
		res.Tags = strings.Join(pr.FreeLabels, ",")
	}
	for _, k := range pr.ProjectLabels {
		l := g.service.RepoLabel.FindByKey(g.rs, k, false)
		if l != nil {
			res.ProjectLabels += l.Name + ","
		}
	}
	res.ProjectLabels = strings.TrimSuffix(res.ProjectLabels, ",")
}

func (g *generation) fillReviewStats(pr *project.Project, res *report.Project) {
	var (
		latestReviewState   overallreview.State
		latestReviewDate    time.Time
		latestReviewComment string
		hasAuditeReview     bool
	)
	for k := range pr.Versions {
		if pr.Versions[k].Deleted {
			continue
		}
		for _, review := range pr.Versions[k].OverallReviews {
			if review.State == overallreview.Audited {
				hasAuditeReview = true
			}
			if review.Created.After(latestReviewDate) {
				latestReviewDate = review.Created
				latestReviewState = review.State
				latestReviewComment = strings.TrimSpace(review.Comment)
			}
		}
	}
	if latestReviewDate.IsZero() {
		return
	}
	res.LatestStatusReviewDate = latestReviewDate.Format(time.RFC3339)
	res.LatestStatusReviewStatus = string(latestReviewState)
	if latestReviewState == overallreview.Audited {
		res.LatestE2ReviewDate = latestReviewDate.Format(time.RFC3339)
		res.LatestE2ReviewComment = strings.ReplaceAll(latestReviewComment, "\n", " ")
	}
	if hasAuditeReview {
		res.HasAuditedE2Review = string(overallreview.Audited)
	}
}

func (g *generation) fillSourceStats(pr *project.Project, res *report.Project) {
	codeReferenceCount := 0
	for k := range pr.Versions {
		if pr.Versions[k].Deleted {
			continue
		}
		codeReferenceCount += countSourceRefs(pr.Versions[k])
	}
	res.NumberOfCodeReference = strconv.Itoa(codeReferenceCount)
}

func (g *generation) fillSbomStats(pr *project.Project, res *report.Project) {
	prs := g.childProjectsOrSelf(pr)

	var (
		sbomFound                     bool
		totalLockedSboms              int
		manuallyLockedCount           int
		uploaded                      int
		latestSbomSourceCodeReference int
		compStats                     components.ComponentStats
		stats                         sbomStats
	)
	for _, iP := range prs {
		var (
			latestSbom        *project.SpdxFileBase
			latestSbomVersion *project.ProjectVersion
		)
		for k := range iP.Versions {
			if iP.Versions[k].Deleted {
				continue
			}
			sboms := g.service.RepoSboms.FindByKey(g.rs, k, false)
			if sboms == nil {
				continue
			}
			uploaded += len(sboms.SpdxFileHistory)
			if len(sboms.SpdxFileHistory) == 0 {
				continue
			}

			latest := sboms.SpdxFileHistory.GetLatest()
			if latestSbom == nil || latest.Created.After(latestSbom.Created) {
				latestSbom = latest
				latestSbomVersion = iP.Versions[k]
			}
			t, m := countLockedSboms(sboms)
			totalLockedSboms += t
			manuallyLockedCount += m
		}
		if latestSbom == nil || latestSbomVersion == nil {
			continue
		}

		sbomFound = true
		if !pr.IsGroup {
			res.LastUpload = latestSbom.Created.String()
		}

		latestSbomSourceCodeReference += countSourceRefs(latestSbomVersion)
		var (
			excHappened bool
			compInfo    components.ComponentInfos
		)
		exception.TryCatch(func() {
			compInfo = g.service.SpdxService.GetComponentInfos(g.rs, pr, latestSbomVersion.Key, latestSbom)
		}, func(exception exception.Exception) {
			excHappened = true
		})
		if excHappened {
			continue
		}

		sbomStats, evalRes := g.processSbom(iP, compInfo, nil, "", true)
		compStats.Allowed += evalRes.Stats.Allowed
		compStats.Denied += evalRes.Stats.Denied
		compStats.NoAssertion += evalRes.Stats.NoAssertion
		compStats.Questioned += evalRes.Stats.Questioned
		compStats.Total += evalRes.Stats.Total
		compStats.Warned += evalRes.Stats.Warned
		stats.andCount += sbomStats.andCount
		stats.keepOfSourceCodeCount += sbomStats.keepOfSourceCodeCount
		stats.massiveAnd += sbomStats.massiveAnd
		stats.massiveOr += sbomStats.massiveOr
		stats.mixedCount += sbomStats.mixedCount
		stats.networkCopyLeftCount += sbomStats.networkCopyLeftCount
		stats.noFossCount += sbomStats.noFossCount
		stats.notDeclaredCount += sbomStats.notDeclaredCount
		stats.orCount += sbomStats.orCount
		stats.permissiveCount += sbomStats.permissiveCount
		stats.strongCopyLeftCount += sbomStats.strongCopyLeftCount
		stats.totalComponentCount += sbomStats.totalComponentCount
		stats.weakCopyLeftCount += sbomStats.weakCopyLeftCount
		stats.withCount += sbomStats.withCount
		stats.GNU_CCSObligationCount += sbomStats.GNU_CCSObligationCount

	}

	res.ManuallyLockedSBOM = strconv.Itoa(manuallyLockedCount)
	res.TotalLockedSBOM = strconv.Itoa(totalLockedSboms)
	res.ProjectSboms = strconv.Itoa(uploaded)

	if !sbomFound {
		return
	}

	res.LatestSbomSourceCodeReference = strconv.Itoa(latestSbomSourceCodeReference)
	res.LatestSbomAllowed = strconv.Itoa(compStats.Allowed)
	res.LatestSbomDenied = strconv.Itoa(compStats.Denied)
	res.LatestSbomUnasserted = strconv.Itoa(compStats.NoAssertion)
	res.LatestSbomWarned = strconv.Itoa(compStats.Warned)
	res.LatestSbomQuestioned = strconv.Itoa(compStats.Questioned)
	res.LatestSbomWeakCopyLeft = strconv.Itoa(stats.weakCopyLeftCount)
	res.LatestSbomStrongCopyLeft = strconv.Itoa(stats.strongCopyLeftCount)
	res.LatestSbomNetworkCopyLeft = strconv.Itoa(stats.networkCopyLeftCount)
	res.LatestSbomAndLicenseExp = strconv.Itoa(stats.andCount)
	res.LatestSbomOrLicenseExp = strconv.Itoa(stats.orCount)
	res.LatestSbomWithLicenseExp = strconv.Itoa(stats.withCount)
	res.LatestSbomMixedLicenseExp = strconv.Itoa(stats.mixedCount)
	res.LatestSbomMassiveAndExp = strconv.Itoa(stats.massiveAnd)
	res.LatestSbomMassiveOrExp = strconv.Itoa(stats.massiveOr)
	res.LatestSbomKeepSourceCode = strconv.Itoa(stats.keepOfSourceCodeCount)
	res.LatestSbomGNU_CCSObligation = strconv.Itoa(stats.GNU_CCSObligationCount)
	res.LatestSbomNoFoss = strconv.Itoa(stats.noFossCount)
	res.LatestSbomTotal = strconv.Itoa(stats.totalComponentCount)
}

func (g *generation) fillTokenStats(pr *project.Project, res *report.Project) {
	res.ProjectTokens = strconv.Itoa(len(pr.Token))
	active := 0
	for _, t := range pr.Token {
		if t.IsExpired() || t.Status == project.REVOKED {
			continue
		}
		active++
	}
	res.ActiveTokens = strconv.Itoa(active)
}

func (g *generation) fillApprovalStats(pr *project.Project, res *report.Project) {
	approvals := g.service.RepoApprovals.FindByKey(g.rs, pr.Key, false)
	if approvals == nil {
		return
	}
	var (
		internals                     []approval.Approval
		latestApproved                approval.Approval
		approvedFound                 bool
		externalFound                 bool
		latestExternal                approval.Approval
		latestApprovalUpdated         time.Time
		approvedApprovalUpdated       time.Time
		latestExternalApprovalUpdated time.Time
	)
	for _, a := range approvals.Approvals {
		switch a.Type {
		case approval.TypeInternal:
			if a.Internal.CustomerDone() {
				approvedFound = true
				latestApproved = a
			}
			internals = append(internals, a)
		case approval.TypeExternal:
			externalFound = true
			latestExternal = a
		}
	}
	if len(internals) > 0 {
		latest := internals[len(internals)-1]
		stats := g.sumSbomStats(latest)
		res.LatestApprovalWeakCopyLeft = strconv.Itoa(stats.weakCopyLeftCount)
		res.LatestApprovalStrongCopyLeft = strconv.Itoa(stats.strongCopyLeftCount)
		res.LatestApprovalNetworkCopyLeft = strconv.Itoa(stats.networkCopyLeftCount)
		res.LatestApprovalAndLicenseExp = strconv.Itoa(stats.andCount)
		res.LatestApprovalOrLicenseExp = strconv.Itoa(stats.orCount)
		res.LatestApprovalWithLicenseExp = strconv.Itoa(stats.withCount)
		res.LatestApprovalMixedLicenseExp = strconv.Itoa(stats.mixedCount)
		res.LatestApprovalMassiveAndExp = strconv.Itoa(stats.massiveAnd)
		res.LatestApprovalMassiveOrExp = strconv.Itoa(stats.massiveOr)
		res.LatestApprovalKeepSourceCode = strconv.Itoa(stats.keepOfSourceCodeCount)
		res.LatestApprovalGNU_CCSObligation = strconv.Itoa(stats.GNU_CCSObligationCount)
		res.LatestApprovalNoFoss = strconv.Itoa(stats.noFossCount)
		if latest.Internal.Aborted {
			res.LatestApprovalStatus = "aborted"
		} else if latest.Internal.IsDeclined() {
			res.LatestApprovalStatus = "declined"
		} else if latest.Internal.SupplierDone() {
			res.LatestApprovalStatus = "pending"
			res.LatestApprovalStatusDetails = "developer approved"
			if latest.Internal.CustomerDone() {
				res.LatestApprovalStatus = "approved"
				res.LatestApprovalStatusDetails = "customer approved"
			}
		} else {
			res.LatestApprovalStatus = "pending"
		}

		latestApprovalUpdated = latest.Updated
		res.LatestApprovalSourceCodeReference = strconv.Itoa(g.countSourceRefsForApprovals(latest))
		res.LatestApprovalAllowed = strconv.Itoa(latest.Info.CompStats.Allowed)
		res.LatestApprovalTotal = strconv.Itoa(latest.Info.CompStats.Allowed +
			latest.Info.CompStats.Warned + latest.Info.CompStats.Denied +
			latest.Info.CompStats.Questioned + latest.Info.CompStats.NoAssertion)
		res.LatestApprovalDenied = strconv.Itoa(latest.Info.CompStats.Denied)
		res.LatestApprovalWarned = strconv.Itoa(latest.Info.CompStats.Warned)
		res.LatestApprovalQuestioned = strconv.Itoa(latest.Info.CompStats.Questioned)
		res.LatestApprovalUnasserted = strconv.Itoa(latest.Info.CompStats.NoAssertion)
	}
	if approvedFound {
		approvedApprovalUpdated = latestApproved.Updated
		res.ApprovedApprovalAllowed = strconv.Itoa(latestApproved.Info.CompStats.Allowed)
		res.ApprovedApprovalTotal = strconv.Itoa(latestApproved.Info.CompStats.Allowed +
			latestApproved.Info.CompStats.Warned + latestApproved.Info.CompStats.Denied +
			latestApproved.Info.CompStats.Questioned + latestApproved.Info.CompStats.NoAssertion)
		res.ApprovedApprovalDenied = strconv.Itoa(latestApproved.Info.CompStats.Denied)
		res.ApprovedApprovalWarned = strconv.Itoa(latestApproved.Info.CompStats.Warned)
		res.ApprovedApprovalQuestioned = strconv.Itoa(latestApproved.Info.CompStats.Questioned)
		res.ApprovedApprovalUnasserted = strconv.Itoa(latestApproved.Info.CompStats.NoAssertion)
		res.ApprovedApprovalLink = conf.Config.Server.DisukoHost + "/#/dashboard/"
		if pr.IsGroup {
			res.ApprovedApprovalLink += "groups/" + pr.Key
		} else {
			res.ApprovedApprovalLink += "projects/" + pr.Key
		}
		res.ApprovedApprovalLink += "/approvals"
	}
	if externalFound {
		latest := latestExternal
		stats := g.sumSbomStats(latest)
		res.LatestExternalApprovalWeakCopyLeft = strconv.Itoa(stats.weakCopyLeftCount)
		res.LatestExternalApprovalStrongCopyLeft = strconv.Itoa(stats.strongCopyLeftCount)
		res.LatestExternalApprovalNetworkCopyLeft = strconv.Itoa(stats.networkCopyLeftCount)
		res.LatestExternalApprovalAndLicenseExp = strconv.Itoa(stats.andCount)
		res.LatestExternalApprovalOrLicenseExp = strconv.Itoa(stats.orCount)
		res.LatestExternalApprovalWithLicenseExp = strconv.Itoa(stats.withCount)
		res.LatestExternalApprovalMixedLicenseExp = strconv.Itoa(stats.mixedCount)
		res.LatestExternalApprovalMassiveAndExp = strconv.Itoa(stats.massiveAnd)
		res.LatestExternalApprovalMassiveOrExp = strconv.Itoa(stats.massiveOr)
		res.LatestExternalApprovalKeepSourceCode = strconv.Itoa(stats.keepOfSourceCodeCount)
		res.LatestExternalApprovalGNU_CCSObligation = strconv.Itoa(stats.GNU_CCSObligationCount)
		res.LatestExternalApprovalNoFoss = strconv.Itoa(stats.noFossCount)
		res.LatestExternalApprovalSourceCodeReference = strconv.Itoa(g.countSourceRefsForApprovals(latest))
		res.LatestExternalApprovalStatus = strings.ToLower(string(latest.External.State))
		latestExternalApprovalUpdated = latest.Updated
		res.LatestExternalApprovalAllowed = strconv.Itoa(latest.Info.CompStats.Allowed)
		res.LatestExternalApprovalTotal = strconv.Itoa(latest.Info.CompStats.Allowed +
			latest.Info.CompStats.Warned + latest.Info.CompStats.Denied +
			latest.Info.CompStats.Questioned + latest.Info.CompStats.NoAssertion)
		res.LatestExternalApprovalDenied = strconv.Itoa(latest.Info.CompStats.Denied)
		res.LatestExternalApprovalWarned = strconv.Itoa(latest.Info.CompStats.Warned)
		res.LatestExternalApprovalQuestioned = strconv.Itoa(latest.Info.CompStats.Questioned)
		res.LatestExternalApprovalUnasserted = strconv.Itoa(latest.Info.CompStats.NoAssertion)
		res.LatestExternalApprovalLink = conf.Config.Server.DisukoHost + "/#/dashboard/"
		if pr.IsGroup {
			res.LatestExternalApprovalLink += "groups/" + pr.Key
		} else {
			res.LatestExternalApprovalLink += "projects/" + pr.Key
		}
		res.LatestExternalApprovalLink += "/approvals"
	}
	if !latestApprovalUpdated.IsZero() {
		res.LatestApprovalUpdated = latestApprovalUpdated.String()
	}
	if !approvedApprovalUpdated.IsZero() {
		res.ApprovedApprovalUpdated = approvedApprovalUpdated.String()
	}
	if !latestExternalApprovalUpdated.IsZero() {
		res.LatestExternalApprovalUpdated = latestExternalApprovalUpdated.String()
	}
}

func (g *generation) fillCustomIds(pr *project.Project, res *report.Project) {
	res.CustomIDs = make([]string, len(g.customIdNames))
	for i, name := range g.customIdNames {
		var values []string
		for _, cid := range pr.CustomIds {
			if cid.TechnicalId == name {
				values = append(values, cid.Value)
			}
		}
		res.CustomIDs[i] = strings.Join(values, "|")
	}
}

func renderBool(in bool) string {
	if in {
		return "true"
	}
	return "false"
}

func hasDummyLabel(currentProject *project.Project, dummyLabel *label.Label) bool {
	if dummyLabel == nil {
		return false
	}
	return slices.Contains(currentProject.ProjectLabels, dummyLabel.GetKey())
}

func getSbomSubscribersCount(pr *project.Project) int {
	var sbomSubscription int
	for _, user := range pr.UserManagement.Users {
		if user.Subscriptions.Spdx {
			sbomSubscription++
		}
	}
	return sbomSubscription
}
