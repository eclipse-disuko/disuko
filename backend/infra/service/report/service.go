// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package report

import (
	"encoding/json"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain/job"
	"github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/domain/obligation"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/report"
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

// GenResult holds the outcome of a report generation run.
type GenResult struct {
	ProjectCnt int
	FileName   string
}

type Service struct {
	Repo                projectRepo.IProjectRepository
	RepoUser            userRepo.IUsersRepository
	RepoDept            department.IDepartmentRepository
	RepoLabel           labels.ILabelRepository
	RepoSboms           sbomRepo.ISbomListRepository
	RepoCustomId        customid.ICustomIdRepository
	RepoApprovals       approvallist.IApprovalListRepository
	RepoPolicyRule      policyrules.IPolicyRulesRepository
	RepoLic             licenserules.ILicenseRulesRepository
	RepoObligation      obligationRepo.IObligationRepository
	SpdxService         *spdx.Service
	ProjectLabelService *projectLabelService.ProjectLabelService
	PolicyDecisionsRepo policydecisions.IPolicyDecisionsRepository
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
) *Service {
	return &Service{
		Repo:                repo,
		RepoUser:            repoUser,
		RepoDept:            repoDept,
		RepoLabel:           repoLabel,
		RepoSboms:           repoSboms,
		RepoApprovals:       repoApprovals,
		RepoObligation:      obligationRepository,
		RepoPolicyRule:      policyRuleRepository,
		RepoLic:             licenseRulesRepository,
		SpdxService:         spdxService,
		RepoCustomId:        repoCustomId,
		ProjectLabelService: prjLabelService,
		PolicyDecisionsRepo: policyDecisionsRepository,
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

func (s *Service) Generate(rs *logy.RequestSession, log *job.Log) GenResult {
	tempHelper := temp.TempHelper{RequestSession: rs}
	tempHelper.CreateRandomFolder()
	defer tempHelper.RemoveAll()
	tmpFileName := tempHelper.GetCompleteFileName(GetCurrentName())

	g := newGeneration(rs, s)
	rep := report.Report{
		CustomIDNames: g.customIdNames,
		Projects:      g.buildProjects(),
	}

	data, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorJsonMarshalling, "report json file"), err)
	}
	if err := os.WriteFile(tmpFileName, data, 0o600); err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCreateFile, "report json file", "header"), err)
	}
	s.uploadReport(rs, GetCurrentName(), tmpFileName)

	now := time.Now()
	if now.Day() == 1 {
		log.AddEntry(job.Info, "saving monthly report %s", GetMonthlyName(now))
		monthlyFileName := tempHelper.GetCompleteFileName(GetMonthlyName(now))
		if err := os.WriteFile(monthlyFileName, data, 0o600); err != nil {
			exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCreateFile, "report monthly json file", "header"), err)
		}
		s.uploadReport(rs, GetMonthlyName(now), monthlyFileName)
	}

	log.AddEntry(job.Info, "successfully report created of %d projects", len(rep.Projects))
	s.cleanupOldMonthlyReports(rs, log, now)

	return GenResult{
		ProjectCnt: len(rep.Projects),
		FileName:   tmpFileName,
	}
}

func (s *Service) uploadReport(rs *logy.RequestSession, reportName string, tmpFileName string) {
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

func (s *Service) cleanupOldMonthlyReports(rs *logy.RequestSession, log *job.Log, now time.Time) {
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
