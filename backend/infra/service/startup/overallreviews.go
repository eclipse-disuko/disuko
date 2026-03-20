// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package startup

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/project"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

func (startUpHandler *StartUpHandler) migrateOverallReviews(requestSession *logy.RequestSession) {
	logy.Infof(requestSession, "migrateOverallReviews - START")

	exception.TryCatchAndLog(requestSession, func() {
		projects := startUpHandler.ProjectRepository.FindAllKeys(requestSession)
		logy.Infof(requestSession, "migrateOverallReviews - Found %d projects to process", len(projects))

		for _, projectID := range projects {
			exception.TryCatch(func() {
				project := startUpHandler.ProjectRepository.FindByKey(requestSession, projectID, false)
				if project == nil {
					logy.Warnf(requestSession, "migrateOverallReviews - Project not found: %s", projectID)
					return
				}
				for _, v := range project.Versions {
					startUpHandler.migrateVersionReviews(requestSession, v)
				}
			}, func(e exception.Exception) {
				exception.LogException(requestSession, e)
				logy.Errorf(requestSession, "migrateOverallReviews - Failed to update project: %s", projectID)
			})
		}
	})

	logy.Infof(requestSession, "migrateOverallReviews - END")
}

func (startUpHandler *StartUpHandler) migrateVersionReviews(requestSession *logy.RequestSession, version *project.ProjectVersion) {
	sbomList := startUpHandler.SbomListRepository.FindByKey(requestSession, version.Key, false)
	if sbomList == nil || sbomList.SpdxFileHistory == nil {
		return
	}

	sbomUpdated := false
	for _, r := range version.OverallReviews {
		if r.SBOMId == "" {
			continue
		}
		sbom := sbomList.SpdxFileHistory.GetByKey(r.SBOMId)
		if sbom == nil {
			continue
		}
		sbom.OverallReview = &r
		logy.Infof(requestSession, "migrateOverallReviews - Adding review %s to sbom %s", r.Key, sbom.Key)
		sbomUpdated = true
	}
	if sbomUpdated {
		logy.Infof(requestSession, "migrateOverallReviews - Updated sbom list %s", version.Key)
		startUpHandler.SbomListRepository.Update(requestSession, sbomList)
	}
}
