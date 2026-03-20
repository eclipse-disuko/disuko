// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package fossdd

import (
	"os"

	"github.com/eclipse-disuko/disuko/domain/job"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/notices"
	"github.com/eclipse-disuko/disuko/helper/s3Helper"
)

type fileMeta struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Hash string `json:"sha256"`
}

type metaTmplData struct {
	ProjectName        string
	ProjectUuid        string
	ReferenceTimeStamp string
	Filename           string
	Created            string
	Updated            string
}

func (g *gen) createZIP() {
	g.collectSpdxFiles()
	g.collectNoticeFiles()
	g.collectORFiles()
	g.collectRRFiles()

	g.jobLog.AddEntry(job.Info, "collected all files")

	z, err := initZipfile(g.tempHelper.GetCompleteFileName(g.opts.Approval.Key + "_archive.zip"))
	if err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCreateFile), err)
	}

	z.addDir(g.tempHelper.GetCompleteFileName("sbom/"))
	z.addDir(g.tempHelper.GetCompleteFileName("notices/"))
	z.addDir(g.tempHelper.GetCompleteFileName("overallreviews/"))
	z.addDir(g.tempHelper.GetCompleteFileName("reviewremarks/"))
	g.addPDFStoZIP(z)
	g.addSnapshotsToZIP(z)
	g.addMetaToZIP(z)
	g.addReadmeToZIP(z)
	g.jobLog.AddEntry(job.Info, "added all files to zip")

	z.close()
}

func (g *gen) addReadmeToZIP(z *zipfile) {
	r := s3Helper.ReadFileFromLocalFileSystem("resources/zipReadme.md")
	if _, err := z.copy("readme.md", r); err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCreateFile), err)
	}
	r.Close()
}

func (g *gen) addMetaToZIP(z *zipfile) {
	if _, err := z.writeJson("meta.json", g.meta); err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCreateFile), err)
	}
	tmplData := metaTmplData{
		ProjectName:        g.data.pr.Name,
		ProjectUuid:        g.data.pr.Key,
		ReferenceTimeStamp: g.refTime.Format(datetimeFormatDE),
		Filename:           g.opts.Approval.Key + "_archive.zip",
		Created:            g.opts.Approval.Created.UTC().Format(datetimeFormatDE),
		Updated:            g.opts.Approval.Updated.UTC().Format(datetimeFormatDE),
	}
	if err := z.executeTmpl("meta.md", g.service.metaTmpl, tmplData); err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCreateFile), err)
	}
}

func (g *gen) addSnapshotsToZIP(z *zipfile) {
	f, err := os.Open(g.tempHelper.GetCompleteFileName("pr-snapshot.json"))
	if err != nil {
		exception.HandleErrorServerMessage(err, message.GetI18N(message.FileNotFound))
	}
	hash, err := z.copy("pr-snapshot.json", f)
	if err != nil {
		exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorCreateFile))
	}
	f.Close()
	g.meta = append(g.meta, fileMeta{
		Type: "pr-snapshot",
		Name: "pr-snapshot.json",
		Hash: hash,
	})

	f, err = os.Open(g.tempHelper.GetCompleteFileName("pc-snapshot.json"))
	if err != nil {
		exception.HandleErrorServerMessage(err, message.GetI18N(message.FileNotFound))
	}
	hash, err = z.copy("pc-snapshot.json", f)
	if err != nil {
		exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorCreateFile))
	}
	f.Close()
	g.meta = append(g.meta, fileMeta{
		Type: "pc-snapshot",
		Name: "pc-snapshot.json",
		Hash: hash,
	})
}

func (g *gen) addPDFStoZIP(z *zipfile) {
	tmpl, ok := g.service.tmpls[g.opts.Template]
	if !ok {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.FileNotFound), "template "+g.opts.Template+" not found")
	}
	for lang := range tmpl.contentPaths {
		f, err := os.Open(g.tempHelper.GetCompleteFileName("disclosure-" + lang.String() + ".pdf"))
		if err != nil {
			exception.HandleErrorServerMessage(err, message.GetI18N(message.FileNotFound))
		}
		hash, err := z.copy("disclosure-"+lang.String()+".pdf", f)
		if err != nil {
			exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorCreateFile))
		}
		f.Close()
		g.meta = append(g.meta, fileMeta{
			Type: "doc",
			Name: "disclosure-" + lang.String() + ".pdf",
			Hash: hash,
		})

	}
}

func (g *gen) collectRRFiles() {
	g.tempHelper.CreateSubFolder("reviewremarks")
	for _, subData := range g.data.subDatas {
		if subData.refs.VersionID == nil {
			continue
		}
		v := subData.pr.GetVersion(*subData.refs.VersionID)
		remarksList := g.service.ReviewRemarksRepo.FindByKey(g.rs, v.Key, false)
		if remarksList == nil {
			continue
		}
		for _, r := range remarksList.Remarks {
			creationTime := r.Created.Format("2006-01-02-15-04")
			hash := s3Helper.SaveObjectToLocalFileAndGetHash(g.rs, g.tempHelper.GetCompleteFileName("reviewremarks/"+creationTime+"_"+r.Key+".json"), r.ToDto(), nil)
			g.jobLog.AddEntry(job.Info, "collected review remark %s", r.Key)
			g.meta = append(g.meta, fileMeta{
				Type: "reviewremark",
				Name: creationTime + "_" + r.Key + ".json",
				Hash: hash,
			})
		}
	}
}

func (g *gen) collectORFiles() {
	g.tempHelper.CreateSubFolder("overallreviews")
	for _, subData := range g.data.subDatas {
		if subData.refs.VersionID == nil || subData.refs.SpdxID == nil {
			continue
		}
		v := subData.pr.GetVersion(*subData.refs.VersionID)
		for _, or := range v.OverallReviews {
			if or.SBOMId != *subData.refs.SpdxID {
				continue
			}
			creationTime := or.Created.Format("2006-01-02-15-04")
			hash := s3Helper.SaveObjectToLocalFileAndGetHash(g.rs, g.tempHelper.GetCompleteFileName("overallreviews/"+creationTime+"_"+or.Key+".json"), or.ToDto(), nil)

			g.jobLog.AddEntry(job.Info, "collected over all file for spdx %s", subData.spdx.Key)
			g.meta = append(g.meta, fileMeta{
				Type: "overallreview",
				Name: creationTime + "_" + or.Key + ".json",
				Hash: hash,
			})

		}
	}
}

func (g *gen) collectNoticeFiles() {
	g.tempHelper.CreateSubFolder("notices")
	for _, subData := range g.data.subDatas {
		if subData.compInfos == nil {
			continue
		}
		notice := notices.GenerateJSONNotices(g.rs, *subData.pr, g.service.LicenseRepo, subData.compInfos, subData.pr.NoticeContactMeta, g.service.LabelRepo)
		hash := s3Helper.SaveObjectToLocalFileAndGetHash(g.rs, g.tempHelper.GetCompleteFileName("notices/"+subData.spdx.Key+".json"), notice, nil)
		g.jobLog.AddEntry(job.Info, "collected notice file for spdx %s", subData.spdx.Key)
		g.meta = append(g.meta, fileMeta{
			Type: "notice",
			Name: subData.spdx.Key + ".json",
			Hash: hash,
		})
	}
}

func (g *gen) collectSpdxFiles() {
	g.tempHelper.CreateSubFolder("sbom")
	for _, subData := range g.data.subDatas {
		if subData.spdx == nil {
			continue
		}
		spdxString := s3Helper.ReadTextFile(g.rs, subData.pr.GetFilePathSbom(subData.spdx.Key, *subData.refs.VersionID), subData.spdx.Hash)
		if spdxString == nil {
			continue
		}
		g.jobLog.AddEntry(job.Info, "collected spdx %s", subData.spdx.Key)
		g.tempHelper.WriteFile("sbom/"+subData.spdx.Key+".json", []byte(*spdxString))
		g.meta = append(g.meta, fileMeta{
			Type: "spdx",
			Name: subData.spdx.Key + ".json",
			Hash: subData.spdx.Hash,
		})

	}
}
