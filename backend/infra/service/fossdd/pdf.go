// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package fossdd

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/language"
	"mercedes-benz.ghe.com/foss/disuko/conf"
	"mercedes-benz.ghe.com/foss/disuko/domain/job"
	"mercedes-benz.ghe.com/foss/disuko/domain/license"
	"mercedes-benz.ghe.com/foss/disuko/domain/project"
	"mercedes-benz.ghe.com/foss/disuko/domain/project/components"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/helper/s3Helper"
	"mercedes-benz.ghe.com/foss/disuko/logy"
	"mercedes-benz.ghe.com/foss/disuko/pdf"
)

type tmpl struct {
	base         *template.Template
	contentPaths map[language.Tag]string
}

type tmplData struct {
	Header  *headerTmplData
	Project *projectTmplData
	Footer  *footerTmplData
}

type headerTmplData struct {
	ProjectName     string
	SupplierName    string
	SupplierAddress string
	SupplierNr      string
	DisclosureDate  string
	Flags           Flags
}

type componentEvalTmplData struct {
	Name          string
	Version       string
	Declared      string
	Concluded     string
	Effective     string
	IsRoot        bool
	RuleApplied   bool
	AliasUsed     bool
	DownloadLink  string
	PURL          string
	Homepage      string
	CopyrightText string
}

type spdxTmplData struct {
	Key        string
	Name       string
	Origin     string
	Uploader   string
	Uploaded   string
	Hash       string
	Evaluation []componentEvalTmplData
}

type sourceTmplData struct {
	URL     string
	Comment string
	Origin  string
	Created string
	Updated string
}

type versionTmplData struct {
	Key     string
	Name    string
	Sources []sourceTmplData
}

type projectRuleTmplData struct {
	Name        string
	Description string
	Key         string
	Created     string
	Updated     string
}

type projectLicenceTmplData struct {
	Name string
	ID   string
}

type projectTmplData struct {
	Nr            int
	Name          string
	Key           string
	SchemaLabel   string
	PolicyLabels  string
	ReferenceTime string

	Spdx          *spdxTmplData
	Version       *versionTmplData
	SbomCompsLink string

	PolicyRules []projectRuleTmplData
	Licenses    []projectLicenceTmplData
}

type ruleLicenseTmplData struct {
	ID   string
	Name string
	URL  string
}

type ruleTmplData struct {
	Nr          int
	Name        string
	Description string
	ID          string
	Allowed     []ruleLicenseTmplData
	Warned      []ruleLicenseTmplData
	Denied      []ruleLicenseTmplData
}

type licenseTmplData struct {
	Name string
	ID   string
	Text string
}

type footerTmplData struct {
	Rules    []ruleTmplData
	Licenses []licenseTmplData
}

func (s *Service) TemplateExist(name string) bool {
	_, ok := s.tmpls[name]
	return ok
}

func (g *gen) createPDFs() {
	tmpl, ok := g.service.tmpls[g.opts.Template]
	if !ok {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.FileNotFound), "template "+g.opts.Template+" not found")
	}
	g.copyCSS(g.opts.Template)

	datas := g.tmplDatas()
	for lang, path := range tmpl.contentPaths {
		htmls := g.createHTMLs(tmpl.base, path, datas)
		outfile := g.tempHelper.GetCompleteFileName(fmt.Sprintf("disclosure-%s.pdf", lang.String()))
		pageHeaderPath := conf.Config.Server.TemplatePageHeaderPath(g.opts.Template, lang)
		pageFooterPath := conf.Config.Server.TemplatePageFooterPath(g.opts.Template, lang)
		err := pdf.ConvertAndMerge(g.rs, outfile, readFileIfExists(pageHeaderPath), readFileIfExists(pageFooterPath), htmls)
		exception.HandleErrorServerMessage(err, message.GetI18N(message.Error))
		err = pdf.StampPageNumbers(g.rs, outfile)
		exception.HandleErrorServerMessage(err, message.GetI18N(message.Error))
	}
}

func readFileIfExists(path string) *string {
	b, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			exception.HandleErrorServerMessage(err, message.GetI18N(message.Error))
		}
		return nil
	}
	s := string(b)
	return &s
}

func (g *gen) copyCSS(tmplName string) {
	matches, err := filepath.Glob(conf.Config.Server.TemplateCSSGlob(tmplName))
	exception.HandleErrorServerMessage(err, message.GetI18N(message.Error))
	for _, path := range matches {
		logy.Infof(g.rs, "copying style file %s to %s", path, g.tempHelper.Path)
		g.jobLog.AddEntry(job.Info, "copying style file %s to %s", path, g.tempHelper.Path)
		s3Helper.CopyLocalFilesystem(path, g.tempHelper.GetCompleteFileName(filepath.Base(path)))
	}
}

func (g *gen) createHTMLs(tmpl *template.Template, path string, datas []tmplData) []string {
	var res []string
	for _, data := range datas {
		var outPath string
		if data.Project == nil {
			outPath = g.tempHelper.GetCompleteFileName("group.html")
		} else {
			outPath = g.tempHelper.GetCompleteFileName("pr_" + data.Project.Key + ".html")
		}
		f, err := os.Create(outPath)
		exception.HandleErrorServerMessage(err, message.GetI18N(message.Error))
		err = tmpl.ExecuteTemplate(f, path, data)
		if err != nil {
			f.Close()
			exception.HandleErrorServerMessage(err, message.GetI18N(message.Error))
		}
		logy.Infof(g.rs, "template %s executed -> %s", path, outPath)
		g.jobLog.AddEntry(job.Info, "template %s executed -> %s", path, outPath)
		res = append(res, outPath)
		f.Close()
	}
	return res
}

func (g *gen) tmplDatas() []tmplData {
	header := g.headerTmpl(g.data.pr)
	if !g.data.pr.IsGroup {
		data := tmplData{
			Header:  &header,
			Project: g.projectTmpl(&g.data.subDatas[0], 0),
			Footer:  g.footerTmpl(),
		}
		return []tmplData{data}
	}

	var res []tmplData
	if len(g.opts.SubProjectsRefs) == 0 {
		res = append(res, tmplData{
			Header:  &header,
			Project: nil,
			Footer:  g.footerTmpl(),
		})
	} else {
		for i, sub := range g.data.subDatas {
			g.jobLog.AddEntry(job.Info, "constructing template data for %s", sub.refs.ProjectID)
			data := tmplData{
				Project: g.projectTmpl(&sub, i),
			}
			if i == 0 {
				data.Header = &header
			}
			if i == len(g.opts.SubProjectsRefs)-1 {
				data.Footer = g.footerTmpl()
			}
			res = append(res, data)
		}
	}
	return res
}

func (g *gen) footerTmpl() *footerTmplData {
	var res footerTmplData
	i := 0
	for _, r := range g.data.allRules {
		res.Rules = append(res.Rules, ruleTmplData{
			Nr:          i,
			Name:        r.Name,
			Description: r.Description,
			ID:          r.Key,
			Allowed:     g.ruleLicensesTmpl(r.allowed),
			Warned:      g.ruleLicensesTmpl(r.warned),
			Denied:      g.ruleLicensesTmpl(r.denied),
		})
		i++
	}
	for _, l := range g.data.allLicenses {
		res.Licenses = append(res.Licenses, licenseTmplData{
			Name: l.Name,
			ID:   l.LicenseId,
			Text: l.Text,
		})
	}
	return &res
}

func (g *gen) ruleLicensesTmpl(lics []*license.License) []ruleLicenseTmplData {
	var res []ruleLicenseTmplData
	for _, lic := range lics {
		res = append(res, ruleLicenseTmplData{
			ID:   lic.LicenseId,
			Name: lic.Name,
			URL:  lic.Meta.LicenseUrl,
		})
	}
	return res
}

func (g *gen) headerTmpl(pr *project.Project) headerTmplData {
	res := headerTmplData{
		ProjectName:     pr.Name,
		SupplierName:    pr.DocumentMeta.SupplierName,
		SupplierAddress: pr.DocumentMeta.SupplierAddress,
		SupplierNr:      pr.DocumentMeta.SupplierNr,
		DisclosureDate:  g.refTime.Format(dateFormatDE),
		Flags:           g.opts.Flags,
	}
	return res
}

func (g *gen) projectTmpl(subData *subData, nr int) *projectTmplData {
	var policyLabels []string
keyLoop:
	for _, lk := range subData.pr.PolicyLabels {
		for _, l := range g.data.allPolicyLabels {
			if l.Key == lk {
				policyLabels = append(policyLabels, l.Name)
				continue keyLoop
			}
		}
	}

	res := projectTmplData{
		Nr:            nr,
		Name:          subData.pr.Name,
		Key:           subData.pr.Key,
		SchemaLabel:   subData.schemaLabel.Name,
		PolicyLabels:  strings.Join(policyLabels, ", "),
		ReferenceTime: g.refTime.Format(datetimeFormatDE),
		Version:       g.versionTmpl(subData),
		SbomCompsLink: g.sbomCompLink(subData),
		Spdx:          g.spdxTmpl(subData),
		PolicyRules:   g.projectRulesTmpl(subData),
		Licenses:      g.projectLicsTmpl(subData),
	}

	return &res
}

func (g *gen) spdxTmpl(subData *subData) *spdxTmplData {
	if subData.refs.SpdxID == nil {
		return nil
	}
	return &spdxTmplData{
		Key:        *subData.refs.SpdxID,
		Name:       subData.spdx.MetaInfo.Name,
		Origin:     subData.spdx.Origin,
		Uploader:   subData.spdx.Uploader,
		Uploaded:   subData.spdx.Uploaded.UTC().Format(datetimeFormatDE),
		Hash:       subData.spdx.Hash,
		Evaluation: g.evalTmpl(subData),
	}
}

func (g *gen) evalTmpl(subData *subData) []componentEvalTmplData {
	var res []componentEvalTmplData
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
		res = append(res, componentEvalTmplData{
			Name:          compRes.Component.Name,
			Version:       compRes.Component.Version,
			Declared:      compRes.Component.LicenseDeclared,
			Concluded:     compRes.Component.License,
			Effective:     compRes.Component.EffectiveLicensesString(),
			IsRoot:        compRes.Component.Type == components.ROOT,
			RuleApplied:   compRes.Component.LicenseRuleApplied != nil && compRes.Component.LicenseRuleApplied.Active,
			AliasUsed:     compRes.Component.IsAliasUsed(),
			DownloadLink:  compRes.Component.DownloadLocation,
			PURL:          compRes.Component.PURL,
			Homepage:      compRes.Component.HomepageURL,
			CopyrightText: compRes.Component.CopyrightText,
		})
	}
	return res
}

func (g *gen) projectLicsTmpl(subData *subData) []projectLicenceTmplData {
	var res []projectLicenceTmplData
	for _, lic := range subData.lics {
		res = append(res, projectLicenceTmplData{
			Name: lic.Name,
			ID:   lic.LicenseId,
		})
	}
	return res
}

func (g *gen) projectRulesTmpl(subData *subData) []projectRuleTmplData {
	var res []projectRuleTmplData
	for _, r := range subData.rules {
		if r.Auxiliary {
			continue
		}
		res = append(res, projectRuleTmplData{
			Name:        r.Name,
			Description: r.Description,
			Key:         r.Key,
			Created:     r.Created.UTC().Format(dateFormatDE),
			Updated:     r.Updated.UTC().Format(dateFormatDE),
		})
	}
	return res
}

func (g *gen) sourcesTmpl(version *project.ProjectVersion) []sourceTmplData {
	var res []sourceTmplData
	for _, s := range version.SourceExternal {
		res = append(res, sourceTmplData{
			URL:     s.URL,
			Comment: s.Comment,
			Origin:  s.Origin,
			Created: s.Created.UTC().UTC().Format(dateFormatDE),
			Updated: s.Updated.UTC().UTC().Format(dateFormatDE),
		})
	}
	return res
}

func (g *gen) versionTmpl(subData *subData) *versionTmplData {
	if subData.refs.VersionID == nil {
		return nil
	}

	v, ok := subData.pr.Versions[*subData.refs.VersionID]
	if !ok {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbNotFound), *subData.refs.VersionID+" not found in DB")
	}

	res := versionTmplData{
		Key:     *subData.refs.VersionID,
		Name:    v.Name,
		Sources: g.sourcesTmpl(v),
	}
	return &res
}

func (g *gen) sbomCompLink(subData *subData) string {
	if subData.refs.SpdxID == nil || subData.refs.VersionID == nil {
		return ""
	}
	return conf.Config.Server.GetSBomComponentsLink(subData.refs.ProjectID, *subData.refs.VersionID, *subData.refs.SpdxID)
}

func nl2br(str string) template.HTML {
	escaped := template.HTMLEscapeString(str)
	return template.HTML(strings.ReplaceAll(escaped, "\n", "<br>"))
}

func inc(i int) int {
	return i + 1
}
