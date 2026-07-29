// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package report

import (
	"strings"
	"time"

	"github.com/eclipse-disuko/disuko/domain/approval"
	"github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/project/components"
	"github.com/eclipse-disuko/disuko/domain/project/sbomlist"
	"github.com/eclipse-disuko/disuko/helper/exception"
)

func countSourceRefs(version *project.ProjectVersion) int {
	sourceCodeReference := 0
	for _, source := range version.SourceExternal {
		if source.URL != "" {
			sourceCodeReference++
		}
	}
	return sourceCodeReference
}

func (g *generation) countSourceRefsForApprovals(approvals approval.Approval) int {
	sourceCodeReference := 0
	for _, projectInfo := range approvals.Info.Projects {
		pr := g.findProject(projectInfo.ProjectKey)
		if pr == nil {
			continue
		}
		for _, version := range pr.Versions {
			sourceCodeReference += countSourceRefs(version)
		}
	}
	return sourceCodeReference
}

func (g *generation) sumSbomStats(approval approval.Approval) sbomStats {
	var res sbomStats
	for _, projectInfo := range approval.Info.Projects {
		approvableSPDX := projectInfo.ApprovableSPDX
		if approvableSPDX.VersionKey == "" {
			continue
		}

		sboms := g.service.RepoSboms.FindByKey(g.rs, approvableSPDX.VersionKey, false)
		if sboms == nil {
			continue
		}
		spdx := sboms.SpdxFileHistory.GetByKey(approvableSPDX.SpdxKey)
		if spdx == nil {
			continue
		}

		pr := g.findProject(projectInfo.ProjectKey)
		if pr == nil {
			continue
		}

		var (
			excHappened bool
			compInfo    components.ComponentInfos
		)
		exception.TryCatch(func() {
			compInfo = g.service.SpdxService.GetComponentInfos(g.rs, pr, approvableSPDX.VersionKey, spdx)
		}, func(exception exception.Exception) {
			excHappened = true
		})
		if excHappened {
			continue
		}
		stats, _ := g.processSbom(pr, compInfo, spdx.Uploaded, spdx.Key, false)

		res.andCount += stats.andCount
		res.keepOfSourceCodeCount += stats.keepOfSourceCodeCount
		res.massiveAnd += stats.massiveAnd
		res.massiveOr += stats.massiveOr
		res.mixedCount += stats.mixedCount
		res.networkCopyLeftCount += stats.networkCopyLeftCount
		res.noFossCount += stats.noFossCount
		res.notDeclaredCount += stats.notDeclaredCount
		res.orCount += stats.orCount
		res.permissiveCount += stats.permissiveCount
		res.strongCopyLeftCount += stats.strongCopyLeftCount
		res.totalComponentCount += stats.totalComponentCount
		res.weakCopyLeftCount += stats.weakCopyLeftCount
		res.withCount += stats.withCount
		res.GNU_CCSObligationCount += stats.GNU_CCSObligationCount

	}
	return res
}

func (g *generation) processSbom(pr *project.Project, ci components.ComponentInfos, sbomUpload *time.Time, sbomKey string, withEval bool) (res sbomStats, evalRes *components.EvaluationResult) {
	if withEval {
		policyRules := g.service.RepoPolicyRule.FindPolicyRulesForLabel(g.rs, pr.PolicyLabels)
		policyDecisions := g.service.PolicyDecisionsRepo.FindByKey(g.rs, pr.Key, false)
		isVehicle := g.service.ProjectLabelService.HasVehiclePlatformLabel(g.rs, pr)
		evalRes = ci.EvaluatePolicyRules(policyRules, policyDecisions, isVehicle, sbomUpload, sbomKey)
	}
	for _, comp := range ci {
		res.totalComponentCount++
		worst := comp.WorstFamily()
		switch worst {
		case license.NetworkCopyleft:
			res.networkCopyLeftCount++
		case license.StrongCopyleft:
			res.strongCopyLeftCount++
		case license.WeakCopyleft:
			res.weakCopyLeftCount++
		case license.Permissive:
			res.permissiveCount++
		default:
			res.notDeclaredCount++
		}
		for _, li := range comp.GetLicensesEffective().List {
			var lic *license.License
			if cached, ok := g.licCache[li.ReferencedLicense]; ok {
				lic = cached
			} else {
				lic = g.service.SpdxService.LicenseRepo.FindById(g.rs, li.ReferencedLicense)
				if lic == nil {
					continue
				}
				g.licCache[li.ReferencedLicense] = lic
			}
			for _, oblKey := range lic.Meta.ObligationsKeyList {
				obligation, ok := g.oblCache[oblKey]
				if !ok {
					obligation = g.service.RepoObligation.FindByKey(g.rs, oblKey, false)
					if obligation == nil {
						continue
					}
					g.oblCache[oblKey] = obligation
				}
				switch obligation.Name {
				case "Keep copy of source code available":
					res.keepOfSourceCodeCount++
				case "GNU-type CCS Obligation":
					res.GNU_CCSObligationCount++
				case "Non-FOSS":
					res.noFossCount++
				}
			}
		}
		operator := comp.GetLicensesEffective().Op
		switch operator {
		case components.AND:
			res.andCount++
			if len(comp.GetLicensesEffective().List) >= 5 {
				res.massiveAnd++
			}
		case components.OR:
			res.orCount++
			if len(comp.GetLicensesEffective().List) >= 5 {
				res.massiveOr++
			}
		}

		if strings.Contains(strings.ToLower(comp.LicenseDeclared), strings.ToLower(WITH)) {
			res.withCount++
		}
		if comp.ComplexExpression {
			res.mixedCount++
		}
	}
	return
}

func countLockedSboms(sboms *sbomlist.SbomList) (total, manual int) {
	for _, sbom := range sboms.SpdxFileHistory {
		if sbom.IsInUse || sbom.IsLocked {
			if sbom.IsLocked {
				manual++
			}
			total++
		}
	}
	return
}
