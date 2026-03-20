// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package fossdd

import (
	"fmt"
	"os"

	"golang.org/x/text/language"
	"mercedes-benz.ghe.com/foss/disuko/domain/project"
	"mercedes-benz.ghe.com/foss/disuko/domain/project/pdocument"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/helper/s3Helper"
	"mercedes-benz.ghe.com/foss/disuko/helper/temp"
	"mercedes-benz.ghe.com/foss/disuko/logy"
	"mercedes-benz.ghe.com/foss/disuko/pdf"
)

func (s *Service) SignDocuments(rs *logy.RequestSession, pr *project.Project, approvalID, signingTetxt string, version int) {
	tempHelper := temp.TempHelper{RequestSession: rs}
	tempHelper.CreateFolder()
	defer tempHelper.RemoveAll()

	prevVersion := version - 1
	if prevVersion < 0 {
		prevVersion = int(pdocument.NONE_VERSION)
	}

	pr = s.ProjectRepo.FindByKey(rs, pr.Key, false)
	if pr == nil {
		logy.Infof(rs, "project deleted in the meantime")
		exception.ThrowExceptionServerMessage(message.GetI18N(message.MissingProject), "")
	}

	signedEn := s.createSignedDoc(rs, tempHelper, pr, approvalID, signingTetxt, prevVersion, version, language.English)
	s.addDocToProject(rs, pr, approvalID, signedEn, language.English, version)
	signedGer := s.createSignedDoc(rs, tempHelper, pr, approvalID, signingTetxt, prevVersion, version, language.German)
	s.addDocToProject(rs, pr, approvalID, signedGer, language.German, version)

	s.ProjectRepo.Update(rs, pr)
}

func (s *Service) addDocToProject(rs *logy.RequestSession, pr *project.Project, approvalKey, path string, lang language.Tag, version int) {
	destName := pdocument.GetFileNameWithIndex(pdocument.PD_DISCLOSURE_DOC, approvalKey, &lang, version)
	destPath := pr.GetFilePathDocumentForProject(destName)
	meta := s3Helper.Metadata(rs, pr, "", destName, rs.ReqID)
	f, err := os.Open(path)
	if err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.FileNotFound), err)
	}
	hash := s3Helper.SaveFileAndGetHash(rs, destPath, f, meta)
	f.Close()
	pr.AddDocument(projectDocument(pdocument.PD_DISCLOSURE_DOC, approvalKey, message.PDDescription, lang.String(), hash, version))
}

func (s *Service) createSignedDoc(rs *logy.RequestSession, tempHelper temp.TempHelper, pr *project.Project, approvalKey, signingTetxt string, prevVersion, version int, lang language.Tag) string {
	pdfInBytes := readDocument(rs, pr, approvalKey, prevVersion, &lang)
	// TODO: refactor to copy
	tempFileName := fmt.Sprintf("disclosure_%s.pdf", lang.String())
	tempHelper.WriteFile(tempFileName, pdfInBytes)
	tempFileNameTarget := fmt.Sprintf("disclosure_%d_%s.pdf", version, lang)
	lastHash := documentHash(pr, approvalKey, lang, prevVersion)
	offset := fmt.Sprintf("offset: -55 -%d", 390+version*30)
	pdf.AddStampToPdf(tempHelper, tempFileName, tempFileNameTarget, signingTetxt+lastHash, offset)
	return tempHelper.GetCompleteFileName(tempFileNameTarget)
}

func documentHash(pr *project.Project, approvalID string, targetLanguage language.Tag, version int) string {
	name := pdocument.GetFileNameWithIndex(pdocument.PD_DISCLOSURE_DOC, approvalID, &targetLanguage, version)
	doc := pr.GetDocumentByFileNameWithIndex(name, version)
	return doc.Hash
}

func readDocument(rs *logy.RequestSession, pr *project.Project, approvalID string, docIndex int, lang *language.Tag) []byte {
	docFile := pdocument.GetFileNameWithIndex(pdocument.PD_DISCLOSURE_DOC, approvalID, lang, docIndex)
	docPath := pr.GetFilePathDocumentForProject(docFile)
	document := pr.GetDocumentByFileNameWithIndex(docFile, docIndex)
	pdfInBytes := s3Helper.ReadFileFully(rs, docPath, document.Hash)
	return pdfInBytes
}
