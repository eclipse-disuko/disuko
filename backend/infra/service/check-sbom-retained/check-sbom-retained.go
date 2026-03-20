// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package sbomLockRetained

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/overallreview"
	"mercedes-benz.ghe.com/foss/disuko/domain/project"
	Iproject "mercedes-benz.ghe.com/foss/disuko/infra/repository/project"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/sbomlist"
	"mercedes-benz.ghe.com/foss/disuko/logy"
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

// CheckVersionHasNonDeletableSboms checks if a specific version has retained SBOMs
func (s *Service) CheckVersionHasNonDeletableSboms(requestSession *logy.RequestSession, version *project.ProjectVersion) bool {
	sbomList := s.sbomListRepository.FindByKey(requestSession, version.Key, false)
	if sbomList == nil || len(sbomList.SpdxFileHistory) == 0 {
		return false
	}
	for _, spdxFile := range sbomList.SpdxFileHistory {
		if IsSpdxToRetain(spdxFile, version) {
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
				if s.CheckVersionHasNonDeletableSboms(requestSession, &versions[i]) {
					return true
				}
			}
		}
	} else {
		// Iterate over the project's own versions.
		versions := currentProject.GetVersions()
		for i := 0; i < len(versions); i++ {
			if s.CheckVersionHasNonDeletableSboms(requestSession, &versions[i]) {
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
		return s.CheckVersionHasNonDeletableSboms(requestSession, version)
	}

	// 3. For a project- (or group-) level deletion (version is nil), check each version (or each child project's version) for a retained SBOM.
	if s.HasAnyVersionWithRetainedSbom(requestSession, currentProject) {
		return true
	}
	return false
}

// Backward compatibility functions - keeping the original function signatures
func CheckVersionHasNonDeletableSboms(requestSession *logy.RequestSession, sbomListRepository sbomlist.ISbomListRepository, version *project.ProjectVersion) bool {
	sbomList := sbomListRepository.FindByKey(requestSession, version.Key, false)
	if sbomList == nil || len(sbomList.SpdxFileHistory) == 0 {
		return false
	}
	for _, spdxFile := range sbomList.SpdxFileHistory {
		if IsSpdxToRetain(spdxFile, version) {
			return true
		}
	}
	return false
}

func HasAnyVersionWithRetainedSbom(requestSession *logy.RequestSession, ProjectRepository Iproject.IProjectRepository, sbomListRepository sbomlist.ISbomListRepository, currentProject *project.Project) bool {
	if currentProject.IsGroup {
		// Iterate over each child project.
		for _, childKey := range currentProject.Children {
			childProj := ProjectRepository.FindByKey(requestSession, childKey, false)
			if childProj == nil {
				continue
			}
			versions := childProj.GetVersions()
			for i := 0; i < len(versions); i++ {
				if CheckVersionHasNonDeletableSboms(requestSession, sbomListRepository, &versions[i]) {
					return true
				}
			}
		}
	} else {
		// Iterate over the project's own versions.
		versions := currentProject.GetVersions()
		for i := 0; i < len(versions); i++ {
			if CheckVersionHasNonDeletableSboms(requestSession, sbomListRepository, &versions[i]) {
				return true
			}
		}
	}
	return false
}

func CheckIfRetainedSbom(requestSession *logy.RequestSession, ProjectRepository Iproject.IProjectRepository, sbomListRepository sbomlist.ISbomListRepository, version *project.ProjectVersion, currentProject *project.Project) bool {
	// 2. If a specific version is provided, check only its retained SBOM status.
	if version != nil {
		return CheckVersionHasNonDeletableSboms(requestSession, sbomListRepository, version)
	}

	// 3. For a project- (or group-) level deletion (version is nil), check each version (or each child project's version) for a retained SBOM.
	if HasAnyVersionWithRetainedSbom(requestSession, ProjectRepository, sbomListRepository, currentProject) {
		return true
	}
	return false
}

func IsSpdxToRetain(spdx *project.SpdxFileBase, version *project.ProjectVersion) bool {
	spdxIsInUse := AnyOverallReviewMatches(spdx.Key, version.OverallReviews) ||
		spdx.ApprovalInfo.IsInApproval ||
		spdx.IsLocked ||
		spdx.IsInUse
	return spdxIsInUse
}

func AnyOverallReviewMatches(spdxKey string, overallReviews []overallreview.OverallReview) bool {
	for _, overallReview := range overallReviews {
		if spdxKey == overallReview.SBOMId {
			return true
		}
	}
	return false
}
