// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package reviewremarks

import (
	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/audit"
	"github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/project/components"
	"github.com/eclipse-disuko/disuko/domain/project/sbomlist"
	"github.com/eclipse-disuko/disuko/domain/reviewremarks"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/infra/repository/auditloglist"
	licenseRepo "github.com/eclipse-disuko/disuko/infra/repository/license"
	"github.com/eclipse-disuko/disuko/infra/repository/licenserules"
	reviewRemarksRepo "github.com/eclipse-disuko/disuko/infra/repository/reviewremarks"
	"github.com/eclipse-disuko/disuko/infra/service/spdx"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/google/go-cmp/cmp"
	"github.com/jinzhu/copier"
)

type spdxRetriever interface {
	RetrieveSbomListAndFile(*logy.RequestSession, string, string) (*sbomlist.SbomList, *project.SpdxFileBase)
}

type ReviewRemarksService struct {
	RequestSession *logy.RequestSession

	LicenseRepo             licenseRepo.ILicensesRepository
	AuditLogListRepository  auditloglist.IAuditLogListRepository
	ReviewRemarksRepository reviewRemarksRepo.IReviewRemarksRepository
	Retriever               spdxRetriever
	LicenseRulesRepo        licenserules.ILicenseRulesRepository
	SpdxService             *spdx.Service
}

func (s *ReviewRemarksService) CreateReviewRemark(prj *project.Project, versionId string, data reviewremarks.ReviewRemarkRequestDto, author string) bool {
	valid, r := s.parseAndValidate(prj, versionId, data)
	if !valid {
		return false
	}
	r.Author = author
	r.Origin = project.OriginUi
	remarks := s.ReviewRemarksRepository.FindByKey(s.RequestSession, versionId, false)
	if remarks == nil {
		remarks = &reviewremarks.ReviewRemarks{
			RootEntity: domain.NewRootEntityWithKey(versionId),
			Remarks: []*reviewremarks.Remark{
				&r,
			},
		}
		s.ReviewRemarksRepository.Save(s.RequestSession, remarks)
		return true
	}
	remarks.Remarks = append(remarks.Remarks, &r)
	s.ReviewRemarksRepository.Update(s.RequestSession, remarks)
	s.AuditLogListRepository.CreateAuditEntryByKey(s.RequestSession, versionId, author, message.ReviewRemarkCreated, audit.DiffWithReporter, r, reviewremarks.Remark{})
	return true
}

func (s *ReviewRemarksService) EditReviewRemark(prj *project.Project, versionId, reviewId, author, fullName string, data reviewremarks.ReviewRemarkRequestDto) bool {
	valid, r := s.parseAndValidate(prj, versionId, data)
	if !valid {
		return false
	}
	r.Author = author
	remarks := s.ReviewRemarksRepository.FindByKey(s.RequestSession, versionId, false)
	var remark *reviewremarks.Remark
	for _, r := range remarks.Remarks {
		if r.Key == reviewId {
			remark = r
			break
		}
	}
	if remark == nil || remark.Status == reviewremarks.Closed || remark.Status == reviewremarks.Cancelled {
		exception.ThrowExceptionBadRequestResponse()
	}

	var before reviewremarks.Remark
	copier.Copy(&before, remark)
	remark.Set(&r, author, fullName)
	s.ReviewRemarksRepository.Update(s.RequestSession, remarks)
	s.AuditLogListRepository.CreateAuditEntryByKey(s.RequestSession, versionId, author, message.ReviewRemarkChanged, cmp.Diff, *remark, before)
	return true
}

func (s *ReviewRemarksService) parseAndValidate(prj *project.Project, versionId string, data reviewremarks.ReviewRemarkRequestDto) (bool, reviewremarks.Remark) {
	res := reviewremarks.Remark{
		ChildEntity: domain.NewChildEntity(),
		Title:       data.Title,
		Description: data.Description,
		Status:      reviewremarks.Open,
	}
	valid, level := reviewremarks.ParseLevel(data.Level)
	if !valid {
		return false, res
	}
	res.Level = level

	if data.SBOMId == "" {
		return true, res
	}
	sbomList, file := s.Retriever.RetrieveSbomListAndFile(s.RequestSession, versionId, data.SBOMId)
	if file == nil {
		return false, res
	}
	res.SBOMName = file.MetaInfo.Name
	res.SBOMId = data.SBOMId

	// Set SBOMUploaded from the history
	if sbomList != nil {
		for _, sbom := range sbomList.SpdxFileHistory {
			if sbom.Key == data.SBOMId {
				res.SBOMUploaded = sbom.Uploaded
				break
			}
		}
	}

	if len(data.Components) > 0 {
		comps := s.SpdxService.GetComponentInfos(s.RequestSession, prj, versionId, file)
		compMap := make(map[string]components.ComponentInfo, len(comps))
		for _, c := range comps {
			compMap[c.SpdxId] = c
		}
		seen := make(map[string]bool, len(data.Components))
		res.Components = make([]reviewremarks.ComponentMeta, 0, len(data.Components))
		for _, id := range data.Components {
			if seen[id] {
				continue
			}
			seen[id] = true
			if c, ok := compMap[id]; ok {
				res.Components = append(res.Components, reviewremarks.ComponentMeta{
					ComponentId:      id,
					ComponentName:    c.Name,
					ComponentVersion: c.Version,
				})
			} else {
				logy.Warnf(s.RequestSession, "Component %s not found in SBOM %s (project/version: %s/%s)", id, data.SBOMId, prj.Name, versionId)
			}
		}
	}

	if len(data.Licenses) > 0 {
		licenses := s.LicenseRepo.FindByIds(s.RequestSession, data.Licenses)
		licMap := make(map[string]*license.License, len(licenses))
		for _, lic := range licenses {
			licMap[lic.LicenseId] = lic
		}
		seen := make(map[string]bool, len(data.Licenses))
		res.Licenses = make([]reviewremarks.LicenseMeta, 0, len(data.Licenses))
		for _, id := range data.Licenses {
			if seen[id] {
				continue
			}
			seen[id] = true
			l := reviewremarks.LicenseMeta{LicenseId: id}
			if lic, ok := licMap[id]; ok {
				l.LicenseName = lic.Name
			}
			res.Licenses = append(res.Licenses, l)
		}
	}

	return true, res
}
