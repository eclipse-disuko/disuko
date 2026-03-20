// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package fossdd

import (
	"html/template"
	textTemplate "text/template"
	"time"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain/approval"
	"github.com/eclipse-disuko/disuko/domain/job"
	"github.com/eclipse-disuko/disuko/domain/project/pdocument"
	"github.com/eclipse-disuko/disuko/helper/s3Helper"
	"github.com/eclipse-disuko/disuko/helper/temp"
	labelsRepo "github.com/eclipse-disuko/disuko/infra/repository/labels"
	licenseRepo "github.com/eclipse-disuko/disuko/infra/repository/license"
	"github.com/eclipse-disuko/disuko/infra/repository/policydecisions"
	policyrulesRepo "github.com/eclipse-disuko/disuko/infra/repository/policyrules"
	projectRepo "github.com/eclipse-disuko/disuko/infra/repository/project"
	"github.com/eclipse-disuko/disuko/infra/repository/reviewremarks"
	"github.com/eclipse-disuko/disuko/infra/repository/sbomlist"
	projectLabelService "github.com/eclipse-disuko/disuko/infra/service/project-label"
	spdxService "github.com/eclipse-disuko/disuko/infra/service/spdx"
	"github.com/eclipse-disuko/disuko/logy"
	"golang.org/x/text/language"
)

type Service struct {
	ProjectRepo         projectRepo.IProjectRepository
	LabelRepo           labelsRepo.ILabelRepository
	SbomListRepo        sbomlist.ISbomListRepository
	SpdxService         *spdxService.Service
	PolicyRuleRepo      policyrulesRepo.IPolicyRulesRepository
	LicenseRepo         licenseRepo.ILicensesRepository
	ReviewRemarksRepo   reviewremarks.IReviewRemarksRepository
	PolicyDecisionsRepo policydecisions.IPolicyDecisionsRepository
	ProjectLabelService *projectLabelService.ProjectLabelService
	metaTmpl            *textTemplate.Template
	tmpls               map[string]tmpl
}

type GenOpts struct {
	Approval        *approval.Approval
	MainProjectID   string
	SubProjectsRefs []SubProjectRefs
	Flags           Flags
	WithZIP         bool
	Template        string
}

type SubProjectRefs struct {
	ProjectID string
	VersionID *string
	SpdxID    *string
}

type Flags struct {
	C1 bool
	C2 bool
	C3 bool
	C4 bool
	C5 bool
	C6 bool
}

type gen struct {
	rs      *logy.RequestSession
	service *Service
	opts    GenOpts
	jobLog  *job.Log

	tempHelper temp.TempHelper
	meta       []fileMeta

	refTime time.Time

	data data

	projectDocs []*pdocument.PDocument

	// TODO: move this to own service and ditch the handler passing
	externCheckCreator ExternCheckCreator
}

var (
	dateFormatDE     = "02.01.2006"
	datetimeFormatDE = "02.01.2006 15:04 UTC"
)

func (s *Service) ReadTemplates(names []string) {
	r := s3Helper.ReadFileFromLocalFileSystem("resources/zipMeta.md")
	tmplBytes := s3Helper.ReadAllAndClose(r)
	s.metaTmpl = textTemplate.Must(textTemplate.New("").Parse(string(tmplBytes)))
	s.tmpls = make(map[string]tmpl)
	for _, name := range names {
		fm := template.FuncMap{
			"nl2br": nl2br,
			"inc":   inc,
		}
		parsed := template.Must(template.New("").Funcs(fm).ParseGlob(conf.Config.Server.TemplateGlob(name)))
		s.tmpls[name] = tmpl{
			base: parsed,
			// TODO: make this fully configurable?
			contentPaths: map[language.Tag]string{
				language.English: conf.Config.Server.TemplateContentName(language.English),
				language.German:  conf.Config.Server.TemplateContentName(language.German),
			},
		}
	}
}

func (s *Service) Generate(rs *logy.RequestSession, opts GenOpts, jobLog *job.Log, externCheckCreator ExternCheckCreator) {
	g := gen{
		rs:                 rs,
		service:            s,
		opts:               opts,
		externCheckCreator: externCheckCreator,
		jobLog:             jobLog,
	}
	g.run()
}

func (g *gen) run() {
	g.refTime = time.Now().UTC()
	g.tempHelper = temp.TempHelper{
		RequestSession: g.rs,
	}
	g.tempHelper.CreateFolder()
	defer g.tempHelper.RemoveAll()

	g.jobLog.AddEntry(job.Info, "collecting data")
	g.collectData()
	g.jobLog.AddEntry(job.Info, "creating pdfs")
	g.createPDFs()
	g.jobLog.AddEntry(job.Info, "creating policy rule snapshot")
	g.createPRSnapshot()
	g.jobLog.AddEntry(job.Info, "creating policy check snapshot")
	g.createPCSnapshot()
	if g.opts.WithZIP {
		g.jobLog.AddEntry(job.Info, "creating zip")
		g.createZIP()
	}
	g.jobLog.AddEntry(job.Info, "finalizing approval files")
	g.finalize()
}
