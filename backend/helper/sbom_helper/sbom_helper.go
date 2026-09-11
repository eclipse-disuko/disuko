package sbom_helper

import (
	"github.com/eclipse-disuko/disuko/domain/project"
	projectRepo "github.com/eclipse-disuko/disuko/infra/repository/project"
	"github.com/eclipse-disuko/disuko/infra/repository/sbomlist"
	"github.com/eclipse-disuko/disuko/logy"
)

func EnsureSbomIsInUse(requestSession *logy.RequestSession, sbomlistRepo sbomlist.ISbomListRepository, versionKey string, sbomKey string, retentionReason string) bool {
	sbomList := sbomlistRepo.FindByKey(requestSession, versionKey, false)
	if sbomList == nil {
		return false
	}

	for _, sbom := range sbomList.SpdxFileHistory {
		if sbom.Key != sbomKey {
			continue
		}
		if sbom.EnsureIsInUse(retentionReason) {
			sbomlistRepo.Update(requestSession, sbomList)
		}
		return true
	}
	return false
}

func EnsureProjectHasSbomToRetain(requestSession *logy.RequestSession, projectRepo projectRepo.IProjectRepository, prj *project.Project) {
	if prj.EnsureSbomToRetain() {
		projectRepo.Update(requestSession, prj)
	}
}
