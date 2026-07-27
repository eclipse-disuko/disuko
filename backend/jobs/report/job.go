// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package report

import (
	"encoding/json"
	"os"
	"path"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain/approval"
	"github.com/eclipse-disuko/disuko/domain/job"
	"github.com/eclipse-disuko/disuko/domain/label"
	"github.com/eclipse-disuko/disuko/domain/license"
	license2 "github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/domain/obligation"
	"github.com/eclipse-disuko/disuko/domain/overallreview"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/project/components"
	"github.com/eclipse-disuko/disuko/domain/project/sbomlist"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/s3Helper"
	"github.com/eclipse-disuko/disuko/helper/temp"
	"github.com/eclipse-disuko/disuko/infra/repository/approvallist"
	"github.com/eclipse-disuko/disuko/infra/repository/customid"
	"github.com/eclipse-disuko/disuko/infra/repository/department"
	"github.com/eclipse-disuko/disuko/infra/repository/labels"
	"github.com/eclipse-disuko/disuko/infra/repository/licenserules"
	obligationRepo "github.com/eclipse-disuko/disuko/infra/repository/obligation"
	"github.com/eclipse-disuko/disuko/infra/repository/policydecisions"
	"github.com/eclipse-disuko/disuko/infra/repository/policyrules"
	projectRepo "github.com/eclipse-disuko/disuko/infra/repository/project"
	sbomRepo "github.com/eclipse-disuko/disuko/infra/repository/sbomlist"
	userRepo "github.com/eclipse-disuko/disuko/infra/repository/user"
	projectLabelService "github.com/eclipse-disuko/disuko/infra/service/project-label"
	"github.com/eclipse-disuko/disuko/infra/service/spdx"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/eclipse-disuko/disuko/scheduler"
)

const WITH = " with "

type sbomStats struct {
	weakCopyLeftCount      int
	strongCopyLeftCount    int
	permissiveCount        int
	networkCopyLeftCount   int
	notDeclaredCount       int
	andCount               int
	orCount                int
	withCount              int
	mixedCount             int
	massiveAnd             int
	massiveOr              int
	keepOfSourceCodeCount  int
	GNU_CCSObligationCount int
	noFossCount            int
	totalComponentCount    int
}

type (
	licCache map[string]*license.License
	prCache  map[string]*project.Project
	oblCache map[string]*obligation.Obligation
)

type Job struct {
	repo                projectRepo.IProjectRepository
	repoUser            userRepo.IUsersRepository
	repoDept            department.IDepartmentRepository
	repoLabel           labels.ILabelRepository
	repoSboms           sbomRepo.ISbomListRepository
	repoCustomId        customid.ICustomIdRepository
	repoApprovals       approvallist.IApprovalListRepository
	repoPolicyRule      policyrules.IPolicyRulesRepository
	repoLic             licenserules.ILicenseRulesRepository
	repoObligation      obligationRepo.IObligationRepository
	spdxService         *spdx.Service
	projectLabelService *projectLabelService.ProjectLabelService
	policyDecisionsRepo policydecisions.IPolicyDecisionsRepository
}

func Init(
	repo projectRepo.IProjectRepository,
	repoUser userRepo.IUsersRepository,
	repoDept department.IDepartmentRepository,
	repoLabel labels.ILabelRepository,
	repoSboms sbomRepo.ISbomListRepository,
	repoApprovals approvallist.IApprovalListRepository,
	obligationRepository obligationRepo.IObligationRepository,
	policyRuleRepository policyrules.IPolicyRulesRepository,
	licenseRulesRepository licenserules.ILicenseRulesRepository,
	spdxService *spdx.Service,
	repoCustomId customid.ICustomIdRepository,
	prjLabelService *projectLabelService.ProjectLabelService,
	policyDecisionsRepository policydecisions.IPolicyDecisionsRepository,
) *Job {
	return &Job{
		repo:                repo,
		repoUser:            repoUser,
		repoDept:            repoDept,
		repoLabel:           repoLabel,
		repoSboms:           repoSboms,
		repoApprovals:       repoApprovals,
		repoObligation:      obligationRepository,
		repoPolicyRule:      policyRuleRepository,
		repoLic:             licenseRulesRepository,
		spdxService:         spdxService,
		repoCustomId:        repoCustomId,
		projectLabelService: prjLabelService,
		policyDecisionsRepo: policyDecisionsRepository,
	}
}

func GetCurrentName() string {
	return "report_all.json"
}

func GetMonthlyName(t time.Time) string {
	return strings.ToLower(t.Month().String()) + "_" + strconv.Itoa(t.Year()) + ".json"
}

func GetReportStorageFileNameOf(fileName string) string {
	return project.RemoveDoubleSlash(strings.Join([]string{conf.Config.Server.GetUploadPath(), "reports", fileName}, "/"))
}

func (j *Job) Execute(rs *logy.RequestSession, info job.Job) scheduler.ExecutionResult {
	var log job.Log
	log.AddEntry(job.Info, "started")

	var customRes struct {
		ProjectCnt int    `json:"projectCnt"`
		FileName   string `json:"fileName"`
	}

	tempHelper := temp.TempHelper{RequestSession: rs}
	tempHelper.CreateRandomFolder()
	defer tempHelper.RemoveAll()
	tmpFileName := tempHelper.GetCompleteFileName(GetCurrentName())

	customIdNames := j.customIdNames(rs)
	report := Report{
		CustomIDNames: customIdNames,
		Projects:      j.buildProjects(rs, customIdNames),
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorJsonMarshalling, "report json file"), err)
	}
	if err := os.WriteFile(tmpFileName, data, 0o600); err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCreateFile, "report json file", "header"), err)
	}
	j.uploadReport(rs, GetCurrentName(), tmpFileName)

	now := time.Now()
	if now.Day() == 1 {
		log.AddEntry(job.Info, "saving monthly report %s", GetMonthlyName(now))
		monthlyFileName := tempHelper.GetCompleteFileName(GetMonthlyName(now))
		if err := os.WriteFile(monthlyFileName, data, 0o600); err != nil {
			exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCreateFile, "report monthly json file", "header"), err)
		}
		j.uploadReport(rs, GetMonthlyName(now), monthlyFileName)
	}

	log.AddEntry(job.Info, "successfully report created of %d projects", len(report.Projects))
	j.cleanupOldMonthlyReports(rs, &log, now)
	log.AddEntry(job.Info, "finished")
	customRes.ProjectCnt = len(report.Projects)
	customRes.FileName = tmpFileName

	return scheduler.ExecutionResult{
		Success:   true,
		Log:       log,
		CustomRes: customRes,
	}
}

func (j *Job) uploadReport(rs *logy.RequestSession, reportName string, tmpFileName string) {
	fileReader, err := os.Open(tmpFileName)
	if err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCreateFile, "read file error "+tmpFileName, "header"), err)
	}
	s3FileName := GetReportStorageFileNameOf(reportName)
	metadata := s3Helper.MetadataForApplication(rs, tmpFileName, rs.ReqID)
	if s3Helper.ExistFile(rs, s3FileName) {
		s3Helper.DeleteFile(rs, s3FileName)
	}
	s3Helper.SaveFile(rs, s3FileName, fileReader, metadata)
	logy.Infof(rs, "Report is uploaded to storage into folder %s", s3FileName)
}

func (j *Job) cleanupOldMonthlyReports(rs *logy.RequestSession, log *job.Log, now time.Time) {
	cutoff := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -12, 0)
	folder := GetReportStorageFileNameOf("")
	for obj := range s3Helper.ListObjects(rs, folder) {
		fileName := path.Base(obj.Key)
		monthTime, ok := parseMonthlyReportFileName(fileName)
		if !ok || !monthTime.Before(cutoff) {
			continue
		}
		s3FileName := GetReportStorageFileNameOf(fileName)
		s3Helper.DeleteFile(rs, s3FileName)
		log.AddEntry(job.Info, "deleted monthly report %s (older than 12 months)", fileName)
		logy.Infof(rs, "Deleted old monthly report %s", s3FileName)
	}
}

func parseMonthlyReportFileName(fileName string) (time.Time, bool) {
	if !strings.HasSuffix(fileName, ".json") {
		return time.Time{}, false
	}
	base := strings.TrimSuffix(fileName, ".json")
	parts := strings.SplitN(base, "_", 2)
	if len(parts) != 2 {
		return time.Time{}, false
	}
	year, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, false
	}
	for m := time.January; m <= time.December; m++ {
		if strings.ToLower(m.String()) == parts[0] {
			return time.Date(year, m, 1, 0, 0, 0, 0, time.UTC), true
		}
	}
	return time.Time{}, false
}

func (j *Job) customIdNames(rs *logy.RequestSession) []string {
	var customIdNames []string
	cids := j.repoCustomId.FindAll(rs, false)
	for _, cid := range cids {
		customIdNames = append(customIdNames, cid.Key)
	}
	return customIdNames
}

func (j *Job) buildProjects(rs *logy.RequestSession, customIdNames []string) []Project {
	prCache := make(prCache)
	licCache := make(licCache)
	oblCache := make(oblCache)
	var projects []Project

	dummyLabel := j.repoLabel.FindByNameAndType(rs, label.DUMMY, label.PROJECT)
	projectKeys := j.repo.FindAllKeys(rs)
	for _, projectKey := range projectKeys {
		pr, ok := prCache[projectKey]
		if !ok {
			pr = j.repo.FindByKey(rs, projectKey, false)
			prCache[projectKey] = pr
		}
		if pr == nil {
			continue
		}
		projects = append(projects, j.project(rs, pr, customIdNames, prCache, licCache, oblCache, hasDummyLabel(pr, dummyLabel)))
	}
	return projects
}

func (j *Job) project(rs *logy.RequestSession, pr *project.Project, customIdNames []string, prCache prCache, licCache licCache, oblCache oblCache, isDummy bool) Project {
	res := Project{}
	j.fillBasicProjectInfo(pr, &res, isDummy)
	j.fillParentAndSupplierInfo(rs, pr, prCache, &res)
	j.fillLabelsAndTags(rs, pr, &res)
	j.fillLicenseDecisionRuleStats(rs, pr, prCache, &res)
	j.fillPolicyDecisionRuleStats(rs, pr, prCache, &res)
	j.fillDeniedPolicyDecisionStats(rs, pr, prCache, &res)
	j.fillReviewStats(pr, &res)
	j.fillSourceStats(pr, &res)
	j.fillTokenStats(pr, &res)
	j.fillCustomIds(pr, customIdNames, &res)
	j.fillSbomStats(rs, pr, &res, prCache, licCache, oblCache)
	j.fillApprovalStats(rs, pr, &res, prCache, licCache, oblCache)
	return res
}

func (j *Job) fillBasicProjectInfo(pr *project.Project, res *Project, isDummy bool) {
	res.Name = pr.Name
	res.Group = renderBool(pr.IsGroup)
	res.NonFoss = renderBool(pr.IsNoFoss)
	res.IsDummy = renderBool(isDummy)
	res.Status = string(pr.Status)
	res.Guid = pr.Key
	res.Created = pr.Created.String()
	res.Updated = pr.Updated.String()
	res.Link = conf.Config.Server.DisukoHost + "/#/dashboard/"
	res.Subscribers = strconv.Itoa(j.getSbomSubscribersCount(pr))
	if pr.IsGroup {
		res.Link += "groups/" + pr.Key
	} else {
		res.Link += "projects/" + pr.Key
	}
}

func (j *Job) fillParentAndSupplierInfo(rs *logy.RequestSession, pr *project.Project, prCache prCache, res *Project) {
	depPrj := pr
	if pr.Parent != "" {
		var ok bool
		depPrj, ok = prCache[pr.Parent]
		if !ok {
			depPrj = j.repo.FindByKey(rs, pr.Parent, false)
			prCache[pr.Parent] = depPrj
		}
		if depPrj == nil {
			logy.Warnf(rs, "parent project %s for %s not found anymore", pr.Parent, pr.Key)
			return
		}
	}
	if depPrj.CustomerMeta.DeptId != "" {
		dep := j.repoDept.GetByDeptId(rs, depPrj.CustomerMeta.DeptId)
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
		dep := j.repoDept.GetByDeptId(rs, depPrj.DocumentMeta.SupplierDeptId)
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
		user := j.repoUser.FindByUserId(rs, responsible.UserId)
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

func (j *Job) fillLicenseDecisionRuleStats(rs *logy.RequestSession, pr *project.Project, prCache prCache, res *Project) {
	prs := []*project.Project{pr}
	if pr.IsGroup {
		prs = make([]*project.Project, 0)
		for _, ck := range pr.Children {
			child, ok := prCache[ck]
			if !ok {
				child = j.repo.FindByKey(rs, ck, false)
				prCache[ck] = child
			}
			if child == nil || child.Deleted {
				continue
			}
			prs = append(prs, child)
		}
	}
	var (
		active   int
		inactive int
	)
	for _, iP := range prs {
		licenseRules := j.repoLic.FindByKey(rs, iP.Key, false)
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

func (j *Job) fillPolicyDecisionRuleStats(rs *logy.RequestSession, pr *project.Project, prCache prCache, res *Project) {
	prs := []*project.Project{pr}
	if pr.IsGroup {
		prs = make([]*project.Project, 0)
		for _, ck := range pr.Children {
			child, ok := prCache[ck]
			if !ok {
				child = j.repo.FindByKey(rs, ck, false)
				prCache[ck] = child
			}
			if child == nil || child.Deleted {
				continue
			}
			prs = append(prs, child)
		}
	}
	var (
		active   int
		inactive int
	)
	for _, iP := range prs {
		policyRules := j.policyDecisionsRepo.FindByKey(rs, iP.Key, false)
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

func (j *Job) fillDeniedPolicyDecisionStats(rs *logy.RequestSession, pr *project.Project, prCache prCache, res *Project) {
	prs := []*project.Project{pr}
	if pr.IsGroup {
		prs = make([]*project.Project, 0)
		for _, ck := range pr.Children {
			child, ok := prCache[ck]
			if !ok {
				child = j.repo.FindByKey(rs, ck, false)
				prCache[ck] = child
			}
			if child == nil || child.Deleted {
				continue
			}
			prs = append(prs, child)
		}
	}
	var (
		active   int
		inactive int
	)
	for _, iP := range prs {
		policyRules := j.policyDecisionsRepo.FindByKey(rs, iP.Key, false)
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

func (j *Job) fillLabelsAndTags(rs *logy.RequestSession, pr *project.Project, res *Project) {
	if pr.SchemaLabel != "" {
		l := j.repoLabel.FindByKey(rs, pr.SchemaLabel, false)
		if l != nil {
			res.SchemaLabel = l.Name
		}
	}
	for _, k := range pr.PolicyLabels {
		l := j.repoLabel.FindByKey(rs, k, false)
		if l != nil {
			res.PolicyLabels += l.Name + ","
		}
	}
	res.PolicyLabels = strings.TrimSuffix(res.PolicyLabels, ",")
	if pr.FreeLabels != nil {
		res.Tags = strings.Join(pr.FreeLabels, ",")
	}
	for _, k := range pr.ProjectLabels {
		l := j.repoLabel.FindByKey(rs, k, false)
		if l != nil {
			res.ProjectLabels += l.Name + ","
		}
	}
	res.ProjectLabels = strings.TrimSuffix(res.ProjectLabels, ",")
}

func (j *Job) fillReviewStats(pr *project.Project, res *Project) {
	var (
		latestReviewState   overallreview.State
		latestReviewDate    time.Time
		latestReviewComment string
	)
	for k := range pr.Versions {
		for _, review := range pr.Versions[k].OverallReviews {
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
		res.LatestE2ReviewStatus = string(latestReviewState)
		res.LatestE2ReviewComment = strings.ReplaceAll(latestReviewComment, "\n", " ")
	}
}

func (j *Job) fillSourceStats(pr *project.Project, res *Project) {
	codeReferenceCount := 0
	for k := range pr.Versions {
		codeReferenceCount += j.countSourceRefs(pr.Versions[k])
	}
	res.NumberOfCodeReference = strconv.Itoa(codeReferenceCount)
}

func (j *Job) fillSbomStats(rs *logy.RequestSession, pr *project.Project, res *Project, prCache prCache, licCache licCache, oblCache oblCache) {
	prs := []*project.Project{pr}
	if pr.IsGroup {
		prs = make([]*project.Project, 0)
		for _, ck := range pr.Children {
			child, ok := prCache[ck]
			if !ok {
				child = j.repo.FindByKey(rs, ck, false)
				prCache[ck] = child
			}
			if child == nil || child.Deleted {
				continue
			}
			prs = append(prs, child)
		}
	}

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
			sboms := j.repoSboms.FindByKey(rs, k, false)
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
			t, m := j.countLockedSboms(sboms)
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

		latestSbomSourceCodeReference += j.countSourceRefs(latestSbomVersion)
		var (
			excHappened bool
			compInfo    components.ComponentInfos
		)
		exception.TryCatch(func() {
			compInfo = j.spdxService.GetComponentInfos(rs, pr, latestSbomVersion.Key, latestSbom)
		}, func(exception exception.Exception) {
			excHappened = true
		})
		if excHappened {
			continue
		}

		sbomStats, evalRes := j.processSbom(rs, iP, compInfo, licCache, oblCache, nil, "", true)
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

func (j *Job) fillTokenStats(pr *project.Project, res *Project) {
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

func (j *Job) fillApprovalStats(rs *logy.RequestSession, pr *project.Project, res *Project, prCache prCache, licCache licCache, oblCache oblCache) {
	approvals := j.repoApprovals.FindByKey(rs, pr.Key, false)
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
		stats := j.sumSbomStats(rs, latest, prCache, licCache, oblCache)
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
		res.LatestApprovalSourceCodeReference = strconv.Itoa(j.countSourceRefsForApprovals(rs, latest, prCache))
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
		stats := j.sumSbomStats(rs, latest, prCache, licCache, oblCache)
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
		res.LatestExternalApprovalSourceCodeReference = strconv.Itoa(j.countSourceRefsForApprovals(rs, latest, prCache))
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

func (j *Job) fillCustomIds(pr *project.Project, customIdNames []string, res *Project) {
	res.CustomIDs = make([]string, len(customIdNames))
	for i, name := range customIdNames {
		var values []string
		for _, cid := range pr.CustomIds {
			if cid.TechnicalId == name {
				values = append(values, cid.Value)
			}
		}
		res.CustomIDs[i] = strings.Join(values, "|")
	}
}

func (j *Job) countSourceRefs(version *project.ProjectVersion) int {
	sourceCodeReference := 0
	for _, source := range version.SourceExternal {
		if source.URL != "" {
			sourceCodeReference++
		}
	}
	return sourceCodeReference
}

func (j *Job) countSourceRefsForApprovals(rs *logy.RequestSession, approvals approval.Approval, prCache map[string]*project.Project) int {
	sourceCodeReference := 0
	for _, projectInfo := range approvals.Info.Projects {
		pr, ok := prCache[projectInfo.ProjectKey]
		if !ok || pr == nil {
			pr = j.repo.FindByKey(rs, projectInfo.ProjectKey, false)
			prCache[projectInfo.ProjectKey] = pr
		}
		if pr == nil {
			continue
		}
		for _, version := range pr.Versions {
			sourceCodeReference += j.countSourceRefs(version)
		}
	}
	return sourceCodeReference
}

func (j *Job) sumSbomStats(rs *logy.RequestSession, approval approval.Approval, prCache prCache, licCache licCache, oblCache oblCache) sbomStats {
	var res sbomStats
	for _, projectInfo := range approval.Info.Projects {
		approvableSPDX := projectInfo.ApprovableSPDX
		if approvableSPDX.VersionKey == "" {
			continue
		}

		sboms := j.repoSboms.FindByKey(rs, approvableSPDX.VersionKey, false)
		if sboms == nil {
			continue
		}
		spdx := sboms.SpdxFileHistory.GetByKey(approvableSPDX.SpdxKey)
		if spdx == nil {
			continue
		}

		pr, ok := prCache[projectInfo.ProjectKey]
		if !ok || pr == nil {
			pr = j.repo.FindByKey(rs, projectInfo.ProjectKey, false)
			prCache[projectInfo.ProjectKey] = pr
		}
		if pr == nil {
			continue
		}

		var (
			excHappened bool
			compInfo    components.ComponentInfos
		)
		exception.TryCatch(func() {
			compInfo = j.spdxService.GetComponentInfos(rs, pr, approvableSPDX.VersionKey, spdx)
		}, func(exception exception.Exception) {
			excHappened = true
		})
		if excHappened {
			continue
		}
		stats, _ := j.processSbom(rs, pr, compInfo, licCache, oblCache, spdx.Uploaded, spdx.Key, false)

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

func (j *Job) processSbom(rs *logy.RequestSession, pr *project.Project, ci components.ComponentInfos, licCache licCache, oblCache oblCache, sbomUpload *time.Time, sbomKey string, withEval bool) (res sbomStats, evalRes *components.EvaluationResult) {
	if withEval {
		policyRules := j.repoPolicyRule.FindPolicyRulesForLabel(rs, pr.PolicyLabels)
		policyDecisions := j.policyDecisionsRepo.FindByKey(rs, pr.Key, false)
		isVehicle := j.projectLabelService.HasVehiclePlatformLabel(rs, pr)
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
			if cached, ok := licCache[li.ReferencedLicense]; ok {
				lic = cached
			} else {
				lic = j.spdxService.LicenseRepo.FindById(rs, li.ReferencedLicense)
				if lic == nil {
					continue
				}
				licCache[li.ReferencedLicense] = lic
			}
			for _, oblKey := range lic.Meta.ObligationsKeyList {
				obligation, ok := oblCache[oblKey]
				if !ok {
					obligation = j.repoObligation.FindByKey(rs, oblKey, false)
					if obligation == nil {
						continue
					}
					oblCache[oblKey] = obligation
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

func (j *Job) countLockedSboms(sboms *sbomlist.SbomList) (total, manual int) {
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

func (j *Job) getSbomSubscribersCount(pr *project.Project) int {
	var sbomSubscription int
	for _, user := range pr.UserManagement.Users {
		if user.Subscriptions.Spdx {
			sbomSubscription++
		}
	}
	return sbomSubscription
}
