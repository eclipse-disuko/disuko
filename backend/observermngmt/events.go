// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package observermngmt

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/approval"
	"mercedes-benz.ghe.com/foss/disuko/domain/oauth"
	"mercedes-benz.ghe.com/foss/disuko/domain/overallreview"
	"mercedes-benz.ghe.com/foss/disuko/domain/project"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type EventId int

const (
	ProjectUpdated EventId = iota + 1
	ProjectDeleted
	ProjectVersionDeleted
	SpdxAdded
	SpdxDeleted
	SpdxUpdatedNewest
	LicenseAdded
	LicenseDeleted
	LicenseAliasAdded
	LicenseAliasDeleted
	ApprovalTaskCreated
	ApprovalFinalized
	OverallReviewCreated
	DatabaseEntryAddedOrDeleted
)

type ApprovalTaskData struct {
	RequestSession *logy.RequestSession
	TargetUser     string
	Type           approval.ApprovalType
	TaskId         string
	Creator        string
	ProjectId      string
	Comment        string
	Approvables    []approval.ProjectApprovable
}

type SpdxData struct {
	RequestSession *logy.RequestSession
	Project        *project.Project
	Version        *project.ProjectVersion
	SpdxFile       *project.SpdxFileBase
}

type ProjectDeletedData struct {
	RequestSession *logy.RequestSession
	Project        *project.Project
}

type VersionData struct {
	RequestSession *logy.RequestSession
	Project        *project.Project
	Version        *project.ProjectVersion
}

type LicenseData struct {
	RequestSession *logy.RequestSession
	Id             string
}

type AliasData struct {
	RequestSession *logy.RequestSession
	Id             string
	Alias          string
}

type ApprovalData struct {
	RequestSession *logy.RequestSession
	Approval       *approval.Approval
	DelegatedTo    string
}

type ProjectUpdatedData struct {
	RequestSession *logy.RequestSession
	Old            *project.Project
	New            *project.Project
	NewParent      *project.Project
}

type OverallReviewData struct {
	RequestSession *logy.RequestSession
	Project        *project.Project
	Version        *project.ProjectVersion
	Review         *overallreview.OverallReview
}

type DatabaseSizeChange struct {
	RequestSession *logy.RequestSession
	Username       string
	CollectionName string
	Rights         *oauth.AccessAndRolesRights
}
