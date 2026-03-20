// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package statistic

import "github.com/eclipse-disuko/disuko/domain"

type DashboardCounts struct {
	ProjectCount        int  `json:"projectCount"`
	LicenseCount        int  `json:"licenseCount"`
	PolicyRuleCount     int  `json:"policyRuleCount"`
	LabelCount          int  `json:"labelCount"`
	SchemaCount         int  `json:"schemaCount"`
	ObligationCount     int  `json:"obligationCount"`
	UserCount           int  `json:"userCount"`
	DisclosureCount     int  `json:"disclosureCount"`
	ReviewTemplateCount int  `json:"reviewTemplateCount"`
	ActiveJobCount      int  `json:"activeJobCount"`
	HasNewNewsboxItem   bool `json:"hasNewNewsboxItem"`
}

type SystemStatsResponseDto struct {
	DayStats    []*SystemStatisticDto `json:"dayStats"`
	MonthsStats []*SystemStatisticDto `json:"monthsStats"`
}

type SystemStatisticDto struct {
	domain.BaseDto
	ProjectCount        int  `json:"projectCount"`
	ProjectActiveCount  int  `json:"projectActiveCount"`
	ProjectDeletedCount int  `json:"projectDeletedCount"`
	MissingProjects     bool `json:"missingProjects"`

	LicenseCount        int  `json:"licenseCount"`
	LicenseChartCount   int  `json:"licenseChartCount"`
	LicenseActiveCount  int  `json:"licenseActiveCount"`
	LicenseDeletedCount int  `json:"licenseDeletedCount"`
	MissingLicenses     bool `json:"missingLicenses"`

	PolicyRuleCount        int  `json:"policyRuleCount"`
	PolicyRuleActiveCount  int  `json:"policyRuleActiveCount"`
	PolicyRuleDeletedCount int  `json:"policyRuleDeletedCount"`
	MissingPolicyRules     bool `json:"missingPolicyRules"`

	LabelCount  int `json:"labelCount"`
	SchemaCount int `json:"schemaCount"`

	ObligationCount        int  `json:"obligationCount"`
	ObligationActiveCount  int  `json:"obligationActiveCount"`
	ObligationDeletedCount int  `json:"obligationDeletedCount"`
	MissingObligations     bool `json:"missingObligations"`

	UploadFileCnt      int  `json:"uploadFileCnt"`
	UploadFileCntPDF   int  `json:"uploadFileCntPDF"`
	UploadFileCntJSON  int  `json:"uploadFileCntJSON"`
	UploadFileCntSBOM  int  `json:"uploadFileCntSBOM"`
	MissingUploadFiles bool `json:"missingUploadFiles"`

	DbBackupFileCnt int `json:"dbBackupFileCnt"`

	UserCount                 int  `json:"userCount"`
	UserActiveCount           int  `json:"userActiveCount"`
	UserDeactivateCount       int  `json:"userDeactivateCount"`
	UserTermsNotAcceptedCount int  `json:"userTermsNotAcceptedCount"`
	UserDeprovisionedCount    int  `json:"userDeprovisionedCount"`
	MissingUsers              bool `json:"missingUsers"`

	MaxVersionsInOneProject      int `json:"maxVersionsInOneProject"`
	ProjectsOverOrAtVersionLimit int `json:"projectsOverOrAtVersionLimit"`
	VersionLimit                 int `json:"versionLimit"`
}
