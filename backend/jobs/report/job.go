// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package report

import (
	"encoding/csv"
	"os"
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
	"golang.org/x/text/encoding/unicode"
)

type Col int

const WITH = " with "

// this does define the ordering

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

type sbomType int

const (
	latestInternal sbomType = iota
	latestExternal
	latestUploaded
)

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

func GetReportAllName() string {
	return "report_all.csv"
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
	tmpFileName := tempHelper.GetCompleteFileName(GetReportAllName())
	tmpFile, err := os.Create(tmpFileName)
	if err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCsvGeneration, "report csv file", "header"), err)
	}
	defer j.finishCSV(rs, tmpFile, tmpFileName)

	convWriter := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewEncoder().Writer(tmpFile)
	csvWriter := csv.NewWriter(convWriter)
	csvWriter.Comma = '\t'
	defer csvWriter.Flush()

	h := j.writeHeaders(rs, csvWriter)
	pc := j.processProjects(rs, csvWriter, h)

	log.AddEntry(job.Info, "successfully report created of %d projects", pc)
	log.AddEntry(job.Info, "finished")
	customRes.ProjectCnt = pc
	customRes.FileName = tmpFileName

	return scheduler.ExecutionResult{
		Success:   true,
		Log:       log,
		CustomRes: customRes,
	}
}

func (j *Job) finishCSV(rs *logy.RequestSession, tmpFile *os.File, tmpFileName string) {
	if err := tmpFile.Close(); err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCreateFile, "close file error "+tmpFileName, "header"), err)
		return
	}
	fileReader, err := os.Open(tmpFileName)
	if err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCreateFile, "read file error "+tmpFileName, "header"), err)
	}
	s3FileName := GetReportStorageFileNameOf(GetReportAllName())
	metadata := s3Helper.MetadataForApplication(rs, tmpFileName, rs.ReqID)
	if s3Helper.ExistFile(rs, s3FileName) {
		s3Helper.DeleteFile(rs, s3FileName)
	}
	s3Helper.SaveFile(rs, s3FileName, fileReader, metadata)
	logy.Infof(rs, "Report is uploaded to storage into folder %s", s3FileName)
}

func (j *Job) writeHeaders(rs *logy.RequestSession, csv *csv.Writer) []string {
	var customIdHeaders []string
	cids := j.repoCustomId.FindAll(rs, false)
	for _, cid := range cids {
		customIdHeaders = append(customIdHeaders, cid.Key)
	}
	combinedHeaders := append(colHeaders, customIdHeaders...)
	j.writeRow(csv, combinedHeaders)
	return combinedHeaders
}

// writeProjectRow writes a single project row to the CSV and handles errors
func (j *Job) writeRow(csvWriter *csv.Writer, row []string) {
	if err := csvWriter.Write(row); err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCsvGeneration, "report analyzes", "data"), err)
	}
}

// processProjects processes a list of projects and groups
func (j *Job) processProjects(rs *logy.RequestSession, csvWriter *csv.Writer, headers []string) int {
	prCache := make(prCache)
	licCache := make(licCache)
	oblCache := make(oblCache)
	processed := 0

	dummyLabel := j.repoLabel.FindByNameAndType(rs, label.DUMMY, label.PROJECT)
	projects := j.repo.FindAllKeys(rs)
	for _, projectKey := range projects {
		pr, ok := prCache[projectKey]
		if !ok {
			pr = j.repo.FindByKey(rs, projectKey, false)
			prCache[projectKey] = pr
		}
		if pr == nil {
			continue
		}
		j.processProject(rs, pr, csvWriter, headers, prCache, licCache, oblCache, hasDummyLabel(pr, dummyLabel))
		processed++

	}
	return processed
}

func (j *Job) processProject(rs *logy.RequestSession, pr *project.Project, csv *csv.Writer, headers []string, prCache prCache, licCache licCache, oblCache oblCache, isDummy bool) {
	j.writeRow(csv, j.row(rs, pr, headers, prCache, licCache, oblCache, isDummy))
}

// row generates a report row for a given project, filling all relevant columns using helper functions.
func (j *Job) row(rs *logy.RequestSession, pr *project.Project, headers []string, prCache prCache, licCache licCache, oblCache oblCache, isDummy bool) []string {
	res := make([]string, len(headers))
	j.fillBasicProjectInfo(pr, &res, isDummy)
	j.fillParentAndSupplierInfo(rs, pr, prCache, &res)
	j.fillLabelsAndTags(rs, pr, &res)
	j.fillLicenseDecisionRuleStats(rs, pr, prCache, &res)
	j.fillPolicyDecisionRuleStats(rs, pr, prCache, &res)
	j.fillDeniedPolicyDecisionStats(rs, pr, prCache, &res)
	j.fillReviewStats(pr, &res)
	j.fillSourceStats(pr, &res)
	j.fillTokenStats(pr, &res)
	j.fillCustomIds(pr, headers, &res)
	j.fillSbomStats(rs, pr, &res, prCache, licCache, oblCache)
	j.fillApprovalStats(rs, pr, &res, prCache, licCache, oblCache)
	return res
}

// fillBasicProjectInfo fills basic project fields
func (j *Job) fillBasicProjectInfo(pr *project.Project, res *[]string, isDummy bool) {
	(*res)[colName] = pr.Name
	(*res)[colGroup] = renderBool(pr.IsGroup)
	(*res)[colNonFoss] = renderBool(pr.IsNoFoss)
	(*res)[colIsDummy] = renderBool(isDummy)
	(*res)[colStatus] = string(pr.Status)
	(*res)[colGuid] = pr.Key
	(*res)[colCreated] = pr.Created.String()
	(*res)[colUpdated] = pr.Updated.String()
	(*res)[colLink] = conf.Config.Server.DisukoHost + "/#/dashboard/"
	(*res)[colSubscribers] = strconv.Itoa(j.getSbomSubscribersCount(pr))
	if pr.IsGroup {
		(*res)[colLink] += "groups/" + pr.Key
	} else {
		(*res)[colLink] += "projects/" + pr.Key
	}
}

// fillParentAndSupplierInfo fills parent, owner, and supplier info
func (j *Job) fillParentAndSupplierInfo(rs *logy.RequestSession, pr *project.Project, prCache prCache, res *[]string) {
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
			(*res)[colOwnerCompanyName] = "deleted in the meantime"
			(*res)[colOwnerCompanyId] = "deleted in the meantime"
			(*res)[colOwnerDepartmentId] = "deleted in the meantime"
			(*res)[colOwnerDepartmentTitle] = "deleted in the meantime"
			(*res)[colOwnerDepartmentAbbreviation] = "deleted in the meantime"
		} else {
			(*res)[colOwnerCompanyName] = dep.CompanyName
			(*res)[colOwnerCompanyId] = dep.CompanyCode
			(*res)[colOwnerDepartmentId] = dep.Key
			(*res)[colOwnerDepartmentTitle] = dep.DescriptionEnglish
			(*res)[colOwnerDepartmentAbbreviation] = dep.OrgAbbreviation
		}
	}
	if depPrj.DocumentMeta.SupplierDeptId != "" && !depPrj.SupplierExtraData.External {
		dep := j.repoDept.GetByDeptId(rs, depPrj.DocumentMeta.SupplierDeptId)
		if dep == nil {
			(*res)[colSupplierCompanyName] = "deleted in the meantime"
			(*res)[colSupplierCompanyId] = "deleted in the meantime"
			(*res)[colSupplierDepartmentId] = "deleted in the meantime"
			(*res)[colSupplierDepartmentTitle] = "deleted in the meantime"
			(*res)[colSupplierDepartmentAbbreviation] = "deleted in the meantime"
		} else {
			(*res)[colSupplierCompanyName] = dep.CompanyName
			(*res)[colSupplierCompanyId] = dep.CompanyCode
			(*res)[colSupplierDepartmentId] = dep.Key
			(*res)[colSupplierDepartmentTitle] = dep.DescriptionEnglish
			(*res)[colSupplierDepartmentAbbreviation] = dep.OrgAbbreviation
		}
	} else if depPrj.SupplierExtraData.External && depPrj.DocumentMeta.SupplierName != "" {
		(*res)[colSupplierCompanyName] = depPrj.DocumentMeta.SupplierName
	}
	(*res)[colSupplierDepartmentExternal] = renderBool(depPrj.SupplierExtraData.External)
	responsible := pr.ProjectResponsible()
	if responsible != nil {
		(*res)[colProjectResponsibleUserid] = responsible.UserId
		// TODO: maybe cache that too?
		user := j.repoUser.FindByUserId(rs, responsible.UserId)
		if user != nil {
			(*res)[colProjectResponsibleEmail] = user.Email
			(*res)[colProjectResponsibleFullName] = user.Lastname + "," + user.Forename
		}
	}
	(*res)[colApplicationId] = pr.ApplicationMeta.Id
	(*res)[colApplicationSecondaryId] = pr.ApplicationMeta.SecondaryId
	(*res)[colApplicationName] = pr.ApplicationMeta.Name
	(*res)[colGroupId] = pr.Parent
	(*res)[colGroupName] = pr.ParentName
}

// countLicenseDecisionRules counts the number of active and inactive license decision rules for a project version.
func (j *Job) fillLicenseDecisionRuleStats(rs *logy.RequestSession, pr *project.Project, prCache prCache, res *[]string) {
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
	(*res)[colActiveLicenseDecisionRules] = strconv.Itoa(active)
	(*res)[colInactiveLicenseDecisionRules] = strconv.Itoa(inactive)
}

// countPolicyDecisionRules counts the number of active and inactive policy decision rules for a project version.
func (j *Job) fillPolicyDecisionRuleStats(rs *logy.RequestSession, pr *project.Project, prCache prCache, res *[]string) {
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
	(*res)[colActivePolicyDecisionRules] = strconv.Itoa(active)
	(*res)[colInactivePolicyDecisionRules] = strconv.Itoa(inactive)
}

func (j *Job) fillDeniedPolicyDecisionStats(rs *logy.RequestSession, pr *project.Project, prCache prCache, res *[]string) {
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
	(*res)[colActiveDeniedPolicyDecision] = strconv.Itoa(active)
	(*res)[colInactiveDeniedPolicyDecision] = strconv.Itoa(inactive)
}

// fillLabelsAndTags fills label and tag columns
func (j *Job) fillLabelsAndTags(rs *logy.RequestSession, pr *project.Project, res *[]string) {
	if pr.SchemaLabel != "" {
		l := j.repoLabel.FindByKey(rs, pr.SchemaLabel, false)
		if l != nil {
			(*res)[colSchemaLabel] = l.Name
		}
	}
	for _, k := range pr.PolicyLabels {
		l := j.repoLabel.FindByKey(rs, k, false)
		if l != nil {
			(*res)[colPolicyLabels] += l.Name + ","
		}
	}
	(*res)[colPolicyLabels] = strings.TrimSuffix((*res)[colPolicyLabels], ",")
	if pr.FreeLabels != nil {
		(*res)[colTags] = strings.Join(pr.FreeLabels, ",")
	}
	for _, k := range pr.ProjectLabels {
		l := j.repoLabel.FindByKey(rs, k, false)
		if l != nil {
			(*res)[colProjectLabels] += l.Name + ","
		}
	}
	(*res)[colProjectLabels] = strings.TrimSuffix((*res)[colProjectLabels], ",")
}

func (j *Job) fillReviewStats(pr *project.Project, res *[]string) {
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
	(*res)[colLatestStatusReviewDate] = latestReviewDate.Format(time.RFC3339)
	(*res)[colLatestStatusReviewStatus] = string(latestReviewState)
	if latestReviewState == overallreview.Audited {
		(*res)[colLatestE2ReviewDate] = latestReviewDate.Format(time.RFC3339)
		(*res)[colLatestE2ReviewStatus] = string(latestReviewState)
		(*res)[colLatestE2ReviewComment] = strings.ReplaceAll(latestReviewComment, "\n", " ")
	}
}

func (j *Job) fillSourceStats(pr *project.Project, res *[]string) {
	codeReferenceCount := 0
	for k := range pr.Versions {
		codeReferenceCount += j.countSourceRefs(pr.Versions[k])
	}
	(*res)[colNumberOfCodeReference] = strconv.Itoa(codeReferenceCount)
}

// fillVersionStats fills version-related statistics
func (j *Job) fillSbomStats(rs *logy.RequestSession, pr *project.Project, res *[]string, prCache prCache, licCache licCache, oblCache oblCache) {
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
			(*res)[colLastUpload] = latestSbom.Created.String()
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

	(*res)[colManuallyLockedSBOM] = strconv.Itoa(manuallyLockedCount)
	(*res)[colTotalLockedSBOM] = strconv.Itoa(totalLockedSboms)
	(*res)[colProjectSboms] = strconv.Itoa(uploaded)

	if !sbomFound {
		return
	}

	(*res)[colLatestSbomSourceCodeReference] = strconv.Itoa(latestSbomSourceCodeReference)
	(*res)[colLatestSbomAllowed] = strconv.Itoa(compStats.Allowed)
	(*res)[colLatestSbomDenied] = strconv.Itoa(compStats.Denied)
	(*res)[colLatestSbomUnasserted] = strconv.Itoa(compStats.NoAssertion)
	(*res)[colLatestSbomWarned] = strconv.Itoa(compStats.Warned)
	(*res)[colLatestSbomQuestioned] = strconv.Itoa(compStats.Questioned)
	(*res)[colLatestSbomWeakCopyLeft] = strconv.Itoa(stats.weakCopyLeftCount)
	(*res)[colLatestSbomStrongCopyLeft] = strconv.Itoa(stats.strongCopyLeftCount)
	(*res)[colLatestSbomNetworkCopyLeft] = strconv.Itoa(stats.networkCopyLeftCount)
	(*res)[colLatestSbomAndLicenseExp] = strconv.Itoa(stats.andCount)
	(*res)[colLatestSbomOrLicenseExp] = strconv.Itoa(stats.orCount)
	(*res)[colLatestSbomWithLicenseExp] = strconv.Itoa(stats.withCount)
	(*res)[colLatestSbomMixedLicenseExp] = strconv.Itoa(stats.mixedCount)
	(*res)[colLatestSbomMassiveAndExp] = strconv.Itoa(stats.massiveAnd)
	(*res)[colLatestSbomMassiveOrExp] = strconv.Itoa(stats.massiveOr)
	(*res)[colLatestSbomKeepSourceCode] = strconv.Itoa(stats.keepOfSourceCodeCount)
	(*res)[colLatestSbomGNU_CCSObligation] = strconv.Itoa(stats.GNU_CCSObligationCount)
	(*res)[colLatestSbomNoFoss] = strconv.Itoa(stats.noFossCount)
	(*res)[colLatestSbomTotal] = strconv.Itoa(stats.totalComponentCount)
}

// fillTokenStats fills token-related statistics
func (j *Job) fillTokenStats(pr *project.Project, res *[]string) {
	(*res)[colProjectTokens] = strconv.Itoa(len(pr.Token))
	active := 0
	for _, t := range pr.Token {
		if t.IsExpired() || t.Status == project.REVOKED {
			continue
		}
		active++
	}
	(*res)[colActiveTokens] = strconv.Itoa(active)
}

// fillApprovalStats fills approval-related statistics
func (j *Job) fillApprovalStats(rs *logy.RequestSession, pr *project.Project, res *[]string, prCache prCache, licCache licCache, oblCache oblCache) {
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
		(*res)[colLatestApprovalWeakCopyLeft] = strconv.Itoa(stats.weakCopyLeftCount)
		(*res)[colLatestApprovalStrongCopyLeft] = strconv.Itoa(stats.strongCopyLeftCount)
		(*res)[colLatestApprovalNetworkCopyLeft] = strconv.Itoa(stats.networkCopyLeftCount)
		(*res)[colLatestApprovalAndLicenseExp] = strconv.Itoa(stats.andCount)
		(*res)[colLatestApprovalOrLicenseExp] = strconv.Itoa(stats.orCount)
		(*res)[colLatestApprovalWithLicenseExp] = strconv.Itoa(stats.withCount)
		(*res)[colLatestApprovalMixedLicenseExp] = strconv.Itoa(stats.mixedCount)
		(*res)[colLatestApprovalMassiveAndExp] = strconv.Itoa(stats.massiveAnd)
		(*res)[colLatestApprovalMassiveOrExp] = strconv.Itoa(stats.massiveOr)
		(*res)[colLatestApprovalKeepSourceCode] = strconv.Itoa(stats.keepOfSourceCodeCount)
		(*res)[colLatestApprovalGNU_CCSObligation] = strconv.Itoa(stats.GNU_CCSObligationCount)
		(*res)[colLatestApprovalNoFoss] = strconv.Itoa(stats.noFossCount)
		if latest.Internal.Aborted {
			(*res)[colLatestApprovalStatus] = "aborted"
		} else if latest.Internal.IsDeclined() {
			(*res)[colLatestApprovalStatus] = "declined"
		} else if latest.Internal.SupplierDone() {
			(*res)[colLatestApprovalStatus] = "pending"
			(*res)[colLatestApprovalStatusDetails] = "developer approved"
			if latest.Internal.CustomerDone() {
				(*res)[colLatestApprovalStatus] = "approved"
				(*res)[colLatestApprovalStatusDetails] = "customer approved"
			}
		} else {
			(*res)[colLatestApprovalStatus] = "pending"
		}

		latestApprovalUpdated = latest.Updated
		(*res)[colLatestApprovalSourceCodeReference] = strconv.Itoa(j.countSourceRefsForApprovals(rs, latest, prCache))
		(*res)[colLatestApprovalAllowed] = strconv.Itoa(latest.Info.CompStats.Allowed)
		(*res)[colLatestApprovalTotal] = strconv.Itoa(latest.Info.CompStats.Allowed +
			latest.Info.CompStats.Warned + latest.Info.CompStats.Denied +
			latest.Info.CompStats.Questioned + latest.Info.CompStats.NoAssertion)
		(*res)[colLatestApprovalDenied] = strconv.Itoa(latest.Info.CompStats.Denied)
		(*res)[colLatestApprovalWarned] = strconv.Itoa(latest.Info.CompStats.Warned)
		(*res)[colLatestApprovalQuestioned] = strconv.Itoa(latest.Info.CompStats.Questioned)
		(*res)[colLatestApprovalUnasserted] = strconv.Itoa(latest.Info.CompStats.NoAssertion)
	}
	if approvedFound {
		approvedApprovalUpdated = latestApproved.Updated
		(*res)[colApprovedApprovalAllowed] = strconv.Itoa(latestApproved.Info.CompStats.Allowed)
		(*res)[colApprovedApprovalTotal] = strconv.Itoa(latestApproved.Info.CompStats.Allowed +
			latestApproved.Info.CompStats.Warned + latestApproved.Info.CompStats.Denied +
			latestApproved.Info.CompStats.Questioned + latestApproved.Info.CompStats.NoAssertion)
		(*res)[colApprovedApprovalDenied] = strconv.Itoa(latestApproved.Info.CompStats.Denied)
		(*res)[colApprovedApprovalWarned] = strconv.Itoa(latestApproved.Info.CompStats.Warned)
		(*res)[colApprovedApprovalQuestioned] = strconv.Itoa(latestApproved.Info.CompStats.Questioned)
		(*res)[colApprovedApprovalUnasserted] = strconv.Itoa(latestApproved.Info.CompStats.NoAssertion)
		(*res)[colApprovedApprovalLink] = conf.Config.Server.DisukoHost + "/#/dashboard/"
		if pr.IsGroup {
			(*res)[colApprovedApprovalLink] += "groups/" + pr.Key
		} else {
			(*res)[colApprovedApprovalLink] += "projects/" + pr.Key
		}
		(*res)[colApprovedApprovalLink] += "/approvals"
	}
	if externalFound {
		latest := latestExternal
		stats := j.sumSbomStats(rs, latest, prCache, licCache, oblCache)
		(*res)[colLatestExternalApprovalWeakCopyLeft] = strconv.Itoa(stats.weakCopyLeftCount)
		(*res)[colLatestExternalApprovalStrongCopyLeft] = strconv.Itoa(stats.strongCopyLeftCount)
		(*res)[colLatestExternalApprovalNetworkCopyLeft] = strconv.Itoa(stats.networkCopyLeftCount)
		(*res)[colLatestExternalApprovalAndLicenseExp] = strconv.Itoa(stats.andCount)
		(*res)[colLatestExternalApprovalOrLicenseExp] = strconv.Itoa(stats.orCount)
		(*res)[colLatestExternalApprovalWithLicenseExp] = strconv.Itoa(stats.withCount)
		(*res)[colLatestExternalApprovalMixedLicenseExp] = strconv.Itoa(stats.mixedCount)
		(*res)[colLatestExternalApprovalMassiveAndExp] = strconv.Itoa(stats.massiveAnd)
		(*res)[colLatestExternalApprovalMassiveOrExp] = strconv.Itoa(stats.massiveOr)
		(*res)[colLatestExternalApprovalKeepSourceCode] = strconv.Itoa(stats.keepOfSourceCodeCount)
		(*res)[colLatestExternalApprovalGNU_CCSObligation] = strconv.Itoa(stats.GNU_CCSObligationCount)
		(*res)[colLatestExternalApprovalNoFoss] = strconv.Itoa(stats.noFossCount)
		(*res)[colLatestExternalApprovalSourceCodeReference] = strconv.Itoa(j.countSourceRefsForApprovals(rs, latest, prCache))
		(*res)[colLatestExternalApprovalStatus] = strings.ToLower(string(latest.External.State))
		latestExternalApprovalUpdated = latest.Updated
		(*res)[colLatestExternalApprovalAllowed] = strconv.Itoa(latest.Info.CompStats.Allowed)
		(*res)[colLatestExternalApprovalTotal] = strconv.Itoa(latest.Info.CompStats.Allowed +
			latest.Info.CompStats.Warned + latest.Info.CompStats.Denied +
			latest.Info.CompStats.Questioned + latest.Info.CompStats.NoAssertion)
		(*res)[colLatestExternalApprovalDenied] = strconv.Itoa(latest.Info.CompStats.Denied)
		(*res)[colLatestExternalApprovalWarned] = strconv.Itoa(latest.Info.CompStats.Warned)
		(*res)[colLatestExternalApprovalQuestioned] = strconv.Itoa(latest.Info.CompStats.Questioned)
		(*res)[colLatestExternalApprovalUnasserted] = strconv.Itoa(latest.Info.CompStats.NoAssertion)
		(*res)[colLatestExternalApprovalLink] = conf.Config.Server.DisukoHost + "/#/dashboard/"
		if pr.IsGroup {
			(*res)[colLatestExternalApprovalLink] += "groups/" + pr.Key
		} else {
			(*res)[colLatestExternalApprovalLink] += "projects/" + pr.Key
		}
		(*res)[colLatestExternalApprovalLink] += "/approvals"
	}
	if !latestApprovalUpdated.IsZero() {
		(*res)[colLatestApprovalUpdated] = latestApprovalUpdated.String()
	}
	if !approvedApprovalUpdated.IsZero() {
		(*res)[colApprovedApprovalUpdated] = approvedApprovalUpdated.String()
	}
	if !latestExternalApprovalUpdated.IsZero() {
		(*res)[colLatestExternalApprovalUpdated] = latestExternalApprovalUpdated.String()
	}
}

// fillCustomIds fills custom ID columns
func (j *Job) fillCustomIds(pr *project.Project, headers []string, res *[]string) {
	for i := len(colHeaders); i < len(headers); i++ {
		for _, cid := range pr.CustomIds {
			if cid.TechnicalId == headers[i] {
				(*res)[i] += cid.Value + "|"
			}
		}
		(*res)[i] = strings.TrimSuffix((*res)[i], "|")
	}
}

// countSourceRefs counts the number of non-empty source code references in a project version.
func (j *Job) countSourceRefs(version *project.ProjectVersion) int {
	sourceCodeReference := 0
	for _, source := range version.SourceExternal {
		if source.URL != "" {
			sourceCodeReference++
		}
	}
	return sourceCodeReference
}

// countSourceRefsForApprovals counts the total number of source code references for all projects in the given approvals.
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
				case "no-FOSS":
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

// countLockedSboms counts the number of manually locked SBOMs in the given SBOM list.
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

// renderBool returns "true" if the input is true, otherwise "false".
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
