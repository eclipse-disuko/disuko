// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package sbomretention

import (
	"github.com/eclipse-disuko/disuko/domain/overallreview"
	"github.com/eclipse-disuko/disuko/domain/project"
	Iproject "github.com/eclipse-disuko/disuko/infra/repository/project"
	"github.com/eclipse-disuko/disuko/infra/repository/sbomlist"
	"github.com/eclipse-disuko/disuko/logy"
)

// Service handles SBOM retention checks
type Service struct {
	projectRepository  Iproject.IProjectRepository
	sbomListRepository sbomlist.ISbomListRepository
}

// NewService creates a new check-sbom-retained service
func NewService(
	projectRepository Iproject.IProjectRepository,
	sbomListRepository sbomlist.ISbomListRepository,
) *Service {
	return &Service{
		projectRepository:  projectRepository,
		sbomListRepository: sbomListRepository,
	}
}

// checkVersionHasNonDeletableSboms checks if a specific version has retained SBOMs
func (s *Service) checkVersionHasNonDeletableSboms(requestSession *logy.RequestSession, version *project.ProjectVersion) bool {
	sbomList := s.sbomListRepository.FindByKey(requestSession, version.Key, false)
	if sbomList == nil || len(sbomList.SpdxFileHistory) == 0 {
		return false
	}
	for _, spdxFile := range sbomList.SpdxFileHistory {
		if isSpdxRetainedOrLocked(spdxFile, version) {
			return true
		}
	}
	return false
}

// HasAnyVersionWithRetainedSbom checks if any version in a project has retained SBOMs
func (s *Service) HasAnyVersionWithRetainedSbom(requestSession *logy.RequestSession, currentProject *project.Project) bool {
	if currentProject.IsGroup {
		// Iterate over each child project.
		for _, childKey := range currentProject.Children {
			childProj := s.projectRepository.FindByKey(requestSession, childKey, false)
			if childProj == nil {
				continue
			}
			versions := childProj.GetVersions()
			for i := 0; i < len(versions); i++ {
				if s.checkVersionHasNonDeletableSboms(requestSession, &versions[i]) {
					return true
				}
			}
		}
	} else {
		// Iterate over the project's own versions.
		versions := currentProject.GetVersions()
		for i := 0; i < len(versions); i++ {
			if s.checkVersionHasNonDeletableSboms(requestSession, &versions[i]) {
				return true
			}
		}
	}
	return false
}

// CheckIfRetainedSbom checks if retained SBOMs exist for deletion operations
func (s *Service) CheckIfRetainedSbom(requestSession *logy.RequestSession, version *project.ProjectVersion, currentProject *project.Project) bool {
	// 2. If a specific version is provided, check only its retained SBOM status.
	if version != nil {
		return s.checkVersionHasNonDeletableSboms(requestSession, version)
	}

	// 3. For a project- (or group-) level deletion (version is nil), check each version (or each child project's version) for a retained SBOM.
	if s.HasAnyVersionWithRetainedSbom(requestSession, currentProject) {
		return true
	}
	return false
}

// IsSpdxToRetain checks if an SPDX file should be retained based on system level or business rules.
func IsSpdxToRetain(spdx *project.SpdxFileBase, version *project.ProjectVersion) bool {
	return anyOverallReviewMatches(spdx.Key, version.OverallReviews) ||
		spdx.ApprovalInfo.IsInApproval ||
		spdx.IsInUse
}

// IsSpdxProtectedFromDeletion checks if an SPDX file is protected from deletion.
func IsSpdxProtectedFromDeletion(spdx *project.SpdxFileBase, prj *project.Project, version *project.ProjectVersion) bool {
	return spdx.Key == prj.ApprovableSPDX.SpdxKey ||
		isSpdxRetainedOrLocked(spdx, version)
}

func isSpdxRetainedOrLocked(spdx *project.SpdxFileBase, version *project.ProjectVersion) bool {
	return spdx.IsLocked ||
		IsSpdxToRetain(spdx, version)
}

func anyOverallReviewMatches(spdxKey string, overallReviews []overallreview.OverallReview) bool {
	for _, overallReview := range overallReviews {
		if spdxKey == overallReview.SBOMId {
			return true
		}
	}
	return false
}

func (s *Service) EnsureSbomIsInUse(requestSession *logy.RequestSession, versionKey string, sbomKey string, retentionReason string) bool {
	sbomList := s.sbomListRepository.FindByKey(requestSession, versionKey, false)
	if sbomList == nil {
		return false
	}

	for _, sbom := range sbomList.SpdxFileHistory {
		if sbom.Key != sbomKey {
			continue
		}
		if sbom.EnsureIsInUse(retentionReason) {
			s.sbomListRepository.Update(requestSession, sbomList)
		}
		return true
	}
	return false
}

func (s *Service) EnsureProjectHasSbomToRetain(requestSession *logy.RequestSession, prj *project.Project) {
	if prj.EnsureSbomToRetain() {
		s.projectRepository.Update(requestSession, prj)
	}
}
