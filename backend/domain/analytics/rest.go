// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package analytics

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/domain/user/approval"
)

type RequestSearchOptions struct {
	Component      string `json:"component" validate:"gte=0,lte=200"`
	License        string `json:"license" validate:"gte=0,lte=500"`
	ExactComponent bool   `json:"exactComponent"`
	ExactLicense   bool   `json:"exactLicense"`
}

type ResponseAnalyticsSearch struct {
	Success bool                 `json:"success"`
	Items   []SearchResponseItem `json:"result"`
	Count   int                  `json:"count"`
	Stats   Statistic            `json:"stats"`
}

type ResponseComponentsSearch struct {
	Components []string `json:"result"`
}

type ResponseLicensesSearch struct {
	Licenses []string `json:"result"`
}

type Statistic struct {
	CountProjects   int `json:"cntProjects"`
	CountComponents int `json:"cntComponents"`
	CountVersions   int `json:"cntVersions"`
	TotalProjects   int `json:"totalProjects"`
	TotalComponents int `json:"totalComponents"`
	TotalVersions   int `json:"totalVersions"`
}

type SearchResponseItem struct {
	Key                string                  `json:"key"`
	Name               string                  `json:"name"`
	ComponentName      string                  `json:"componentName"`
	ComponentVersion   string                  `json:"componentVersion"`
	ProjectVersionKey  string                  `json:"projectVersionKey"`
	ProjectVersionName string                  `json:"projectVersionName"`
	Type               ItemType                `json:"type"`
	LicenseDeclared    string                  `json:"licenseDeclared"`
	LicenseConcluded   string                  `json:"licenseConcluded"`
	EntryLicense       string                  `json:"entryLicense"`
	SBomName           string                  `json:"sbomName"`
	SBomStatus         approval.ApprovalStatus `json:"sbomStatus"`
	Responsible        string                  `json:"responsible"`
	LastUpdate         time.Time               `json:"lastUpdate"`
	OwnerDeptId        string                  `bson:"-" json:"-"`
	OwnerCompany       string                  `json:"ownerCompany"`
	OwnerDept          string                  `json:"ownerDep"`
	OwnerDeptMissing   bool                    `json:"ownerDeptMissing"`
}

type ItemType string

const (
	PROJECT   = "PROJECT"
	COMPONENT = "COMPONENT"
	LICENSE   = "LICENSE"
)

type OccurrenceDto struct {
	OrigName          string                  `json:"origName"`
	ReferencedLicense string                  `json:"referencedLicense"`
	Count             int                     `json:"count"`
	License           *license.LicenseSlimDto `json:"license,omitempty"`
}

type OccurrencesResDto struct {
	List           []OccurrenceDto              `json:"list"`
	PossibleValues license.PossibleFilterValues `json:"possibleValues"`
}

type StatsDto struct {
	ProjectCount        int `json:"projectCount"`
	ProjectActiveCount  int `json:"projectActiveCount"`
	ProjectDeletedCount int `json:"projectDeletedCount"`

	LicenseCount        int `json:"licenseCount"`
	LicenseChartCount   int `json:"licenseChartCount"`
	LicenseActiveCount  int `json:"licenseActiveCount"`
	LicenseDeletedCount int `json:"licenseDeletedCount"`
	LicenseForbiden     int `json:"licenseForbiddenCount"`
	LicenseUnknown      int `json:"licenseUnknownCount"`

	UploadFileCntSBOM int `json:"uploadFileCntSBOM"`

	UserCount                 int `json:"userCount"`
	UserActiveCount           int `json:"userActiveCount"`
	UserDeactivateCount       int `json:"userDeactivateCount"`
	UserTermsNotAcceptedCount int `json:"userTermsNotAcceptedCount"`

	CompletedTrainings int `json:"completedTrainings"`
}
