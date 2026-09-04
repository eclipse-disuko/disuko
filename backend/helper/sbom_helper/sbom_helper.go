package sbom_helper

import (
	"github.com/eclipse-disuko/disuko/infra/repository/sbomlist"
	"github.com/eclipse-disuko/disuko/logy"
)

func EnsureSbomIsInUse(requestSession *logy.RequestSession, repository sbomlist.ISbomListRepository, versionKey string, sbomKey string, retentionReason string) bool {
	sbomList := repository.FindByKey(requestSession, versionKey, false)
	if sbomList == nil {
		return false
	}

	for _, sbom := range sbomList.SpdxFileHistory {
		if sbom.Key != sbomKey {
			continue
		}
		if sbom.EnsureIsInUse(retentionReason) {
			repository.Update(requestSession, sbomList)
		}
		return true
	}
	return false
}
