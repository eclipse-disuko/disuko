// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package integrity

import (
	"strconv"

	"golang.org/x/text/language"
	approval2 "mercedes-benz.ghe.com/foss/disuko/domain/approval"
	"mercedes-benz.ghe.com/foss/disuko/domain/integrity"
	"mercedes-benz.ghe.com/foss/disuko/domain/project/pdocument"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/approvallist"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/database"
	project2 "mercedes-benz.ghe.com/foss/disuko/infra/repository/project"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

func checkIfEachApprovalHasDocumentsOnProject(requestSession *logy.RequestSession,
	projectRepository project2.IProjectRepository,
	approvalListRepo approvallist.IApprovalListRepository,
	fixIt bool,
	state *integrity.DbIntegrityResult) {

	approvalsCount := approvalListRepo.CountAllWithDeleted(requestSession)
	logy.Infof(requestSession, "Start analyse approvals, approval count: "+strconv.Itoa(approvalsCount))

	offset := 0
	limit := 100
	for {
		qc := database.New().SetLimit(offset, limit)
		qbRawRes := approvalListRepo.Query(requestSession, qc)
		var qbRes []string
		for _, e := range qbRawRes {
			qbRes = append(qbRes, e.GetKey())
		}
		for _, approvalListKey := range qbRes {
			state.CountApprovals++
			approvalList := approvalListRepo.FindByKeyWithDeleted(requestSession, approvalListKey, false)
			if approvalList == nil || approvalList.Approvals == nil || len(approvalList.Approvals) == 0 {
				continue
			}

			prj := projectRepository.FindByKeyWithDeleted(requestSession, approvalListKey, false)
			documents := prj.GetDocuments()
			for _, approval := range approvalList.Approvals {
				if approval.Type == approval2.TypePlausibility {
					// PLAUSIBILITY has no Docs
					continue
				}

				if approval.Type == approval2.TypeExternal {
					expectedDocuments := createExpectedDocsForExternalApproval(approval)

					missingDocsMeta := make([]*integrity.DocumentMeta, 0)
					if len(documents) == 0 {
						missingDocsMeta = expectedDocuments
					} else {
						for _, expectedDocument := range expectedDocuments {
							if !contains(documents, expectedDocument) {
								missingDocsMeta = append(missingDocsMeta, expectedDocument)
							}
						}
					}
					if len(missingDocsMeta) > 0 {
						missingDocumentRefs := &integrity.MissingDocumentRefsOnProjectForApproval{
							ApprovalId:            approval.Key,
							ApprovalCreated:       approval.Created,
							ApprovalUpdated:       approval.Updated,
							ApprovalType:          approval.Type,
							ApprovalDocVersion:    -1,
							ApprovalListIsDeleted: approvalList.IsDeleted(),
							ProjectUuid:           prj.UUID(),
							ProjectCreated:        prj.Created,
							ProjectName:           prj.Name,
							ProjectIsDeleted:      prj.IsDeleted(),
							MissingDocsCount:      len(missingDocsMeta),
							MissingDocs:           missingDocsMeta,
						}
						state.MissingDocRefsOnProject = append(state.MissingDocRefsOnProject, missingDocumentRefs)
					}
				}
				if approval.Type == approval2.TypeInternal {
					expectedDocuments := createExpectedDocsForInternalApproval(approval)

					missingDocsMeta := make([]*integrity.DocumentMeta, 0)
					if len(documents) == 0 {
						missingDocsMeta = expectedDocuments
					} else {
						for _, expectedDocument := range expectedDocuments {
							if !contains(documents, expectedDocument) {
								missingDocsMeta = append(missingDocsMeta, expectedDocument)
							}
						}
					}
					if len(missingDocsMeta) > 0 {
						missingDocumentRefs := &integrity.MissingDocumentRefsOnProjectForApproval{
							ApprovalId:            approval.Key,
							ApprovalCreated:       approval.Created,
							ApprovalUpdated:       approval.Updated,
							ApprovalType:          approval.Type,
							ApprovalDocVersion:    approval.Internal.DocVersion,
							ApprovalListIsDeleted: approvalList.IsDeleted(),
							ProjectUuid:           prj.UUID(),
							ProjectCreated:        prj.Created,
							ProjectName:           prj.Name,
							ProjectIsDeleted:      prj.IsDeleted(),
							MissingDocsCount:      len(missingDocsMeta),
							MissingDocs:           missingDocsMeta,
						}
						state.MissingDocRefsOnProject = append(state.MissingDocRefsOnProject, missingDocumentRefs)
					}
				}
			}

		}

		logy.Infof(requestSession, "analysed approvals: "+strconv.Itoa(state.CountApprovals)+" / "+strconv.Itoa(approvalsCount))

		if len(qbRes) < limit {
			break
		}
		offset += limit
	}

	if state.CountApprovals != approvalsCount {
		addError2(requestSession, state, "Different analysed approval count: "+strconv.Itoa(state.CountApprovals)+" / "+strconv.Itoa(approvalsCount))
	}

	logy.Infof(requestSession, "End analyse approvals: "+strconv.Itoa(state.CountApprovals)+" / "+strconv.Itoa(approvalsCount))
}

func createExpectedDocsForInternalApproval(approval approval2.Approval) []*integrity.DocumentMeta {
	// INTERNAL Approvals has always policies, policyChecks and count of disclosures depending on approval.Internal.DocVersion,
	expectedDocuments := make([]*integrity.DocumentMeta, 0)
	version := int(pdocument.NONE_VERSION)
	policiesDoc := integrity.DocumentMeta{
		ApprovalId: approval.Key,
		Type:       pdocument.PD_POLICY_RULES,
		Version:    version,
		Language:   "",
	}
	policieChecksDoc := integrity.DocumentMeta{
		ApprovalId: approval.Key,
		Type:       pdocument.PD_POLICY_CHECK,
		Version:    version,
		Language:   "",
	}
	expectedDocuments = append(expectedDocuments, &policiesDoc)
	expectedDocuments = append(expectedDocuments, &policieChecksDoc)

	for i := 0; i <= approval.Internal.DocVersion; i++ {
		currentVersion := i - 1
		if currentVersion < 0 {
			currentVersion = int(pdocument.NONE_VERSION)
		}
		disclosureDe := integrity.DocumentMeta{
			ApprovalId: approval.Key,
			Type:       pdocument.PD_DISCLOSURE_DOC,
			Version:    currentVersion,
			Language:   language.German.String(),
		}
		disclosureEn := integrity.DocumentMeta{
			ApprovalId: approval.Key,
			Type:       pdocument.PD_DISCLOSURE_DOC,
			Version:    currentVersion,
			Language:   language.English.String(),
		}
		expectedDocuments = append(expectedDocuments, &disclosureDe)
		expectedDocuments = append(expectedDocuments, &disclosureEn)
	}
	return expectedDocuments
}

func createExpectedDocsForExternalApproval(approval approval2.Approval) []*integrity.DocumentMeta {
	// EXTERNAL must have 4 Docs: policies, policyCheck, disclosures en/de
	expectedDocuments := make([]*integrity.DocumentMeta, 0)
	version := int(pdocument.NONE_VERSION)
	policiesDoc := integrity.DocumentMeta{
		ApprovalId: approval.Key,
		Type:       pdocument.PD_POLICY_RULES,
		Version:    version,
		Language:   "",
	}
	policieChecksDoc := integrity.DocumentMeta{
		ApprovalId: approval.Key,
		Type:       pdocument.PD_POLICY_CHECK,
		Version:    version,
		Language:   "",
	}
	disclosureDe := integrity.DocumentMeta{
		ApprovalId: approval.Key,
		Type:       pdocument.PD_DISCLOSURE_DOC,
		Version:    version,
		Language:   language.German.String(),
	}
	disclosureEn := integrity.DocumentMeta{
		ApprovalId: approval.Key,
		Type:       pdocument.PD_DISCLOSURE_DOC,
		Version:    version,
		Language:   language.English.String(),
	}
	expectedDocuments = append(expectedDocuments, &policiesDoc)
	expectedDocuments = append(expectedDocuments, &policieChecksDoc)
	expectedDocuments = append(expectedDocuments, &disclosureDe)
	expectedDocuments = append(expectedDocuments, &disclosureEn)
	return expectedDocuments
}

func contains(docs []*pdocument.PDocument, expectedDoc *integrity.DocumentMeta) bool {
	skipLanguageCheck := expectedDoc.Type == pdocument.PD_POLICY_CHECK || expectedDoc.Type == pdocument.PD_POLICY_RULES
	for _, doc := range docs {
		if doc.ApprovalId == expectedDoc.ApprovalId &&
			doc.Type == expectedDoc.Type &&
			int(*doc.VersionIndex) == expectedDoc.Version &&
			(skipLanguageCheck || doc.Lang == expectedDoc.Language) {
			return true
		}
	}
	return false
}
