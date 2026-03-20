// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package fossdd

import (
	"fmt"
	"os"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/job"
	"github.com/eclipse-disuko/disuko/domain/project/pdocument"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/s3Helper"
	"github.com/eclipse-disuko/disuko/logy"
	"golang.org/x/text/language"
)

func (g *gen) finalize() {
	pr := g.service.ProjectRepo.FindByKey(g.rs, g.opts.MainProjectID, false)
	if pr == nil {
		logy.Infof(g.rs, "project deleted in the meantime")
		g.jobLog.AddEntry(job.Error, "project deleted in the meantime")
		exception.ThrowExceptionServerMessage(message.GetI18N(message.MissingProject), "")
		return
	}
	g.copyPDFsToProject()
	g.jobLog.AddEntry(job.Info, "copied pdfs to project")
	g.copySnapshotsToProject()
	g.jobLog.AddEntry(job.Info, "copied snapshots to project")
	if g.opts.WithZIP {
		g.copyZIPtoProject()
		g.jobLog.AddEntry(job.Info, "copied zip to project")
	}

	pr.AddDocument(g.projectDocs...)
	g.service.ProjectRepo.Update(g.rs, pr)
}

func (g *gen) copyZIPtoProject() {
	f, err := os.Open(g.tempHelper.GetCompleteFileName(g.opts.Approval.Key + "_archive.zip"))
	if err != nil {
		exception.HandleErrorServerMessage(err, message.GetI18N(message.FileNotFound))
	}
	destPath := g.data.pr.GetFilePathDocumentForProject(g.opts.Approval.Key + "_archive.zip")
	metadata := s3Helper.Metadata(g.rs, g.data.pr, "", g.opts.Approval.Key+"_archive.zip", g.rs.ReqID)
	hash := s3Helper.SaveFileAndGetHash(g.rs, destPath, f, metadata)
	f.Close()
	g.projectDocs = append(g.projectDocs, g.projectDocument(pdocument.PD_ARCHIVE, message.Archive, language.English.String(), hash, int(pdocument.NONE_VERSION)))
}

func (g *gen) copySnapshotsToProject() {
	f, err := os.Open(g.tempHelper.GetCompleteFileName("pr-snapshot.json"))
	if err != nil {
		exception.HandleErrorServerMessage(err, message.GetI18N(message.FileNotFound))
	}
	snapName := pdocument.GetFileName(pdocument.PD_POLICY_RULES, g.opts.Approval.Key, &language.English)
	destPath := g.data.pr.GetFilePathDocumentForProject(snapName)
	meta := s3Helper.Metadata(g.rs, g.data.pr, "", snapName, g.rs.ReqID)
	hash := s3Helper.SaveFileAndGetHash(g.rs, destPath, f, meta)
	f.Close()
	g.projectDocs = append(g.projectDocs, g.projectDocument(pdocument.PD_POLICY_RULES, message.FileDescriptionPolicyRules, language.English.String(), hash, int(pdocument.NONE_VERSION)))

	f, err = os.Open(g.tempHelper.GetCompleteFileName("pc-snapshot.json"))
	if err != nil {
		exception.HandleErrorServerMessage(err, message.GetI18N(message.FileNotFound))
	}
	snapName = pdocument.GetFileName(pdocument.PD_POLICY_CHECK, g.opts.Approval.Key, &language.English)
	destPath = g.data.pr.GetFilePathDocumentForProject(snapName)
	meta = s3Helper.Metadata(g.rs, g.data.pr, "", snapName, g.rs.ReqID)
	hash = s3Helper.SaveFileAndGetHash(g.rs, destPath, f, meta)
	f.Close()
	g.projectDocs = append(g.projectDocs, g.projectDocument(pdocument.PD_POLICY_CHECK, message.FileDescriptionPolicyCheck, language.English.String(), hash, int(pdocument.NONE_VERSION)))
}

func (g *gen) copyPDFsToProject() {
	tmpl, ok := g.service.tmpls[g.opts.Template]
	if !ok {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.FileNotFound), "template "+g.opts.Template+" not found")
	}
	for lang := range tmpl.contentPaths {
		f, err := os.Open(g.tempHelper.GetCompleteFileName(fmt.Sprintf("disclosure-%s.pdf", lang.String())))
		if err != nil {
			exception.HandleErrorServerMessage(err, message.GetI18N(message.FileNotFound))
		}

		pdfName := pdocument.GetFileName(pdocument.PD_DISCLOSURE_DOC, g.opts.Approval.Key, &lang)
		destPath := g.data.pr.GetFilePathDocumentForProject(pdfName)
		meta := s3Helper.Metadata(g.rs, g.data.pr, "", pdfName, g.rs.ReqID)
		hash := s3Helper.SaveFileAndGetHash(g.rs, destPath, f, meta)
		f.Close()
		g.projectDocs = append(g.projectDocs, g.projectDocument(pdocument.PD_DISCLOSURE_DOC, message.PDDescription, lang.String(), hash, int(pdocument.NONE_VERSION)))
	}
}

func (g *gen) projectDocument(fileType pdocument.PDocumentType, description string, lang string, hash string, index int) *pdocument.PDocument {
	return projectDocument(fileType, g.opts.Approval.Key, description, lang, hash, index)
}

func projectDocument(fileType pdocument.PDocumentType, approvalID, description, lang, hash string, index int) *pdocument.PDocument {
	docVersion := pdocument.DocumentVersion(index)
	pDocument := &pdocument.PDocument{
		ChildEntity:  domain.NewChildEntity(),
		Description:  description,
		ApprovalId:   approvalID,
		Type:         fileType,
		Lang:         lang,
		VersionIndex: &docVersion,
		Hash:         hash,
	}
	return pDocument
}
