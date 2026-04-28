// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"fmt"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/overallreview"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/infra/repository/auditloglist"
	projectRepo "github.com/eclipse-disuko/disuko/infra/repository/project"
	"github.com/eclipse-disuko/disuko/infra/repository/sbomlist"
	"github.com/eclipse-disuko/disuko/infra/repository/user"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/eclipse-disuko/disuko/observermngmt"
	"github.com/google/go-cmp/cmp"
	"github.com/jinzhu/copier"
)

type OverallReviewService struct {
	AuditlogRepo auditloglist.IAuditLogListRepository
	ProjectRepo  projectRepo.IProjectRepository
	SbomListRepo sbomlist.ISbomListRepository
	UserRepo     user.IUsersRepository
}

func (s *OverallReviewService) AddToProjectFromDTO(
	rs *logy.RequestSession,
	pr *project.Project,
	version *project.ProjectVersion,
	username string,
	dto overallreview.OverallReviewRequestDto,
) {
	s.AddToProject(rs, pr, version, username, dto.State, dto.Comment, dto.SBOMId)
}

func (s *OverallReviewService) AddToProject(
	rs *logy.RequestSession,
	pr *project.Project,
	version *project.ProjectVersion,
	username string,
	state overallreview.State,
	comment string,
	sbomID string,
) {
	oldVersion := project.ProjectVersion{}
	copier.Copy(&oldVersion, version)

	version.Status = versionStatus(state)

	sbomList := s.SbomListRepo.FindByKey(rs, version.Key, false)
	if sbomList == nil {
		exception.ThrowExceptionBadRequestResponse()
	}

	spdx := sbomList.SpdxFileHistory.GetByKey(sbomID)
	if sbomList == nil {
		exception.ThrowExceptionBadRequestResponse()
	}

	review := overallreview.OverallReview{
		ChildEntity:     domain.NewChildEntity(),
		Creator:         username,
		CreatorFullName: s.fullname(rs, username),
		Comment:         comment,
		State:           state,
		SBOMId:          spdx.Key,
		SBOMName:        spdx.MetaInfo.Name,
		SBOMUploaded:    spdx.Uploaded.Format("2006-01-02 15:04:05"),
	}

	version.OverallReviews = append(version.OverallReviews, review)

	s.ProjectRepo.Update(rs, pr)

	spdx.OverallReview = &review

	s.SbomListRepo.Update(rs, sbomList)

	observermngmt.FireEvent(observermngmt.OverallReviewCreated, observermngmt.OverallReviewData{
		RequestSession: rs,
		Project:        pr,
		Version:        &oldVersion,
		Review:         &review,
	})

	s.AuditlogRepo.CreateAuditEntryByKey(rs, version.Key, username, message.OverallReviewUpdated, cmp.Diff, version, oldVersion)
}

func (s OverallReviewService) fullname(rs *logy.RequestSession, userid string) string {
	user := s.UserRepo.FindByUserId(rs, userid)
	if user != nil {
		return fmt.Sprintf("%s %s", user.Forename, user.Lastname)
	}
	return userid
}

func versionStatus(state overallreview.State) project.ProjectVersionStatusType {
	switch state {
	case overallreview.Audited:
		return project.PV_Audited
	case overallreview.Acceptable:
		return project.PV_Acceptable
	case overallreview.AcceptableAfterChanges:
		return project.PV_AcceptableAfterChanges
	case overallreview.NotAcceptable:
		return project.PV_NotAcceptable
	default:
		return project.PV_Unreviewed
	}
}
