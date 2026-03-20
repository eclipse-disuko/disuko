// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package integrity

import (
	"time"

	"mercedes-benz.ghe.com/foss/disuko/domain"
	"mercedes-benz.ghe.com/foss/disuko/domain/approval"
	"mercedes-benz.ghe.com/foss/disuko/domain/project/pdocument"
)

type MissingS3File struct {
	ProjectUuid      string
	ProjectName      string
	ProjectIsDeleted bool
	VersionName      string
	VersionUuid      string
	VersionIsDeleted bool
	S3FileName       string
	OrgFileName      string
	Upload           *time.Time
	Message          string
	Fixed            bool
	DbEntityKey      string
	ErrorCode        string
	ErrorMessage     string
	ErrorRaw         string
}

type DbIntegrityResult struct {
	domain.RootEntity    `bson:"inline"`
	domain.BaseState     `bson:"inline"`
	FixIt                bool
	MissingFileOnS3      []*MissingS3File
	MissingFileOnS3Count int
	MissingFileOnDB      []*MissingDBFile
	MissingFileOnDBCount int
	Errors               []string
	ErrorsCount          int

	CountFilesOnDB   int
	CountUploadsOnS3 int
	CountProjects    int
	CountApprovals   int

	MissingDocRefsOnProject      []*MissingDocumentRefsOnProjectForApproval
	MissingDocRefsOnProjectCount int
}
type MissingDBFile struct {
	FilePath     string
	Upload       time.Time
	Message      string
	Fixed        bool
	MetaData     map[string]string
	ErrorCode    string
	ErrorMessage string
	ErrorRaw     string
	FileDeleted  bool
}

type DocumentMeta struct {
	ApprovalId string
	Type       pdocument.PDocumentType
	Version    int
	Language   string
}

type MissingDocumentRefsOnProjectForApproval struct {
	ApprovalId            string
	ApprovalCreated       time.Time
	ApprovalUpdated       time.Time
	ApprovalType          approval.ApprovalType
	ApprovalDocVersion    int
	ApprovalListIsDeleted bool
	ProjectUuid           string
	ProjectCreated        time.Time
	ProjectName           string
	ProjectIsDeleted      bool
	MissingDocsCount      int
	MissingDocs           []*DocumentMeta
}
