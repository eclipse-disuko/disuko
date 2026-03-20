// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package license

import (
	"github.com/eclipse-disuko/disuko/domain/audit"

	"github.com/eclipse-disuko/disuko/domain"

	"github.com/eclipse-disuko/disuko/helper/hash"
	"github.com/eclipse-disuko/disuko/logy"
)

type Alias struct {
	domain.ChildEntity `bson:"inline"`
	LicenseId          string
	Description        string
}

type LicenseRef struct {
	ID            string
	Family        FamilyOfLicense
	ApprovalState ApprovalStatus
}

type LicenseRefs map[string]LicenseRef

type MetaData struct {
	Family             FamilyOfLicense `json:"family"`
	ApprovalState      ApprovalStatus  `json:"approvalState"`
	ReviewState        ReviewStatus    `json:"reviewState"`
	ReviewDateStr      string          `json:"reviewDate,omitempty"`
	ObligationsKeyList []string        `json:"obligationsKeyList"`
	LicenseUrl         string          `json:"licenseUrl"`
	SourceUrl          string          `json:"sourceUrl"`
	OSIApproved        bool            `json:"osiApproved"`
	FSFApproved        bool            `json:"fsfApproved"`
	Changelog          string          `json:"changelog"`
	LicenseType        TypeOfLicenses  `json:"licenseType"`
	Evaluation         string          `json:"evaluation"`
	LegalComments      string          `json:"legalComments"`
	IsLicenseChart     bool
}

type License struct {
	domain.RootEntity     `bson:"inline"`
	audit.Container       `bson:"inline"`
	domain.SoftDelete     `bson:"inline"`
	IsDeprecatedLicenseId bool     `json:"isDeprecatedLicenseId"`
	LicenseId             string   `json:"licenseId"`
	Name                  string   `json:"name"`
	Text                  string   `json:"text"`
	Meta                  MetaData `json:"meta"`
	Source                Source   `json:"source"`
	Aliases               []Alias
}

type Source string

const (
	PublicLicenseDb = "spdx"
	CUSTOM          = "custom"
)

type ApprovalStatus string

const (
	NotSet     = ""
	Pending    = "pending"
	Check      = "check"
	Assigning  = "assigning"
	Approved   = "approved"
	Forbidden  = "forbidden"
	Deprecated = "deprecated"
)

func (as ApprovalStatus) Value() string {
	if len(as) == 0 {
		return "not set"
	}
	return string(as)
}

type ReviewStatus string

const (
	NotReviewed   = ""
	ReviewRequest = "review request"
	InReview      = "in review"
	Reviewed      = "reviewed"
)

type FamilyOfLicense string

const (
	NotDeclared     FamilyOfLicense = ""
	Permissive      FamilyOfLicense = "permissive"
	WeakCopyleft    FamilyOfLicense = "weak copyleft"
	StrongCopyleft  FamilyOfLicense = "strong copyleft"
	NetworkCopyleft FamilyOfLicense = "network copyleft"
)

func (lf FamilyOfLicense) Value() string {
	if len(lf) == 0 {
		return "not declared"
	}
	return string(lf)
}

type TypeOfLicenses string

const (
	NotDeclaredType = ""
	Freeware        = "freeware"
	Proprietary     = "proprietary"
	PublicDomain    = "public domain"
	OpenSource      = "open source"
	NonFoss         = "non foss"
)

func (lt TypeOfLicenses) Value() string {
	if len(lt) == 0 {
		return "not declared"
	}
	return string(lt)
}

func (l *LicenseRefs) GenHash(requestSession *logy.RequestSession) string {
	return hash.Hash(requestSession, l)
}
