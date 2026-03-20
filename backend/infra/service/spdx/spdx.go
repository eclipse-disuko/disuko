// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package spdx

import (
	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/project/components"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/s3Helper"
	"github.com/eclipse-disuko/disuko/infra/repository/license"
	"github.com/eclipse-disuko/disuko/infra/repository/licenserules"
	"github.com/eclipse-disuko/disuko/infra/service/cache"
	"github.com/eclipse-disuko/disuko/infra/service/locks"
	"github.com/eclipse-disuko/disuko/logy"
)

type Service struct {
	LicenseRepo      license.ILicensesRepository
	LicenseRulesRepo licenserules.ILicenseRulesRepository
	LockService      *locks.Service
}

// Change this to invalidate every cached SPDX evaluation
var salt = "06/03/2026"

func (s *Service) GetComponentInfos(rs *logy.RequestSession, prj *project.Project, versionKey string, spdx *project.SpdxFileBase) components.ComponentInfos {
	var compInfoList components.ComponentInfoList
	cacheService := cache.CacheService{
		RequestSession: rs,
		WithLock:       true,
		LockService:    s.LockService,
	}
	found := cacheService.GetCacheEntry(spdx.Key, &compInfoList)
	if !found {
		logy.Infof(rs, "compinfo file %s not present yet", spdx.Key)
		spdxString := s3Helper.ReadTextFile(rs, prj.GetFilePathSbom(spdx.Key, versionKey), spdx.Hash)
		ci := project.FileContent(*spdxString).ExtractComponentInfo(rs)
		ci.EnrichComponentInfos(rs)
		compInfoList = components.ComponentInfoList{
			RootEntity:     domain.NewRootEntityWithKey(spdx.Key),
			ComponentInfos: ci,
		}
	}

	currentRefs := s.LicenseRepo.GetLicenseRefs(rs)
	currentHash := currentRefs.GenHash(rs)
	currentHash += salt

	// license rules hash
	licenseRules := s.LicenseRulesRepo.FindByKey(rs, prj.Key, false)
	currentLicenseRulesHash := licenseRules.GenHash(rs)
	currentLicenseRulesHash += salt

	if currentHash == compInfoList.UsedRefsListHash && currentLicenseRulesHash == compInfoList.UsedLicenseRulesHash {
		return compInfoList.ComponentInfos
	}

	logy.Infof(rs, "The referencelist has changed, reapplying refs")

	// #6642: reset applied license rules and re-apply later again, after new refs was applied
	compInfoList.ComponentInfos.CleanLicenseRules()

	compInfoList.ComponentInfos.EnrichComponentInfos(rs)
	compInfoList.ComponentInfos.ApplyRefs(currentRefs)

	// #6642: apply license rules
	compInfoList.ComponentInfos.ApplyLicenseRules(licenseRules, spdx.Uploaded, spdx.Key)

	compInfoList.UsedRefsListHash = currentHash
	compInfoList.UsedLicenseRulesHash = currentLicenseRulesHash

	exception.TryCatch(func() {
		cacheService.SaveCacheEntry(spdx.Key, compInfoList)
	}, func(e exception.Exception) {
		exception.LogException(rs, e)
	})
	return compInfoList.ComponentInfos
}

func (s *Service) WriteComponentInfos(rs *logy.RequestSession, spdxKey string, info *components.ComponentInfoList) {
	cacheService := cache.CacheService{
		RequestSession: rs,
		LockService:    s.LockService,
	}
	cacheService.SaveCacheEntry(spdxKey, info)
}
