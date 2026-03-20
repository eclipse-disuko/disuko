// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package statistic

import (
	"mercedes-benz.ghe.com/foss/disuko/domain"
)

type SystemStatistic struct {
	domain.RootEntity            `bson:",inline"`
	ProjectCount                 int `json:"projectCount"`
	ProjectActiveCount           int `json:"projectActiveCount"`
	ProjectDeletedCount          int `json:"projectDeletedCount"`
	LicenseCount                 int `json:"licenseCount"`
	LicenseChartCount            int `json:"LicenseChartCount"`
	LicenseActiveCount           int `json:"licenseActiveCount"`
	LicenseDeletedCount          int `json:"licenseDeletedCount"`
	PolicyRuleCount              int `json:"policyRuleCount"`
	PolicyRuleActiveCount        int `json:"policyRuleActiveCount"`
	PolicyRuleDeletedCount       int `json:"policyRuleDeletedCount"`
	LabelCount                   int `json:"labelCount"`
	SchemaCount                  int `json:"schemaCount"`
	ObligationCount              int `json:"obligationCount"`
	ObligationActiveCount        int `json:"obligationActiveCount"`
	ObligationDeletedCount       int `json:"obligationDeletedCount"`
	UploadedFilesCnt             int `json:"uploadFileCnt"`
	UploadedFilesCntPDF          int `json:"uploadFileCntPDF"`
	UploadedFilesCntJSON         int `json:"uploadFileCntJSON"`
	UploadedFilesCntSBOM         int `json:"uploadFileCntSBOM"`
	DbBackupFilesCnt             int `json:"dbBackupFileCnt"`
	UserCount                    int `json:"userCount"`
	UserActiveCount              int `json:"userActiveCount"`
	UserDeprovisionedCount       int `json:"userDeprovisionedCount"`
	UserDeactivateCount          int `json:"userDeactivateCount"`
	UserTermsNotAcceptedCount    int `json:"userTermsNotAcceptedCount"`
	MaxVersionsInOneProject      int `json:"maxVersionsInOneProject"`
	ProjectsOverOrAtVersionLimit int `json:"projectsOverOrAtVersionLimit"`
	VersionLimit                 int `json:"versionLimit"`
}

func (systemStatistic *SystemStatistic) ToDto() *SystemStatisticDto {
	return ToDto(systemStatistic)
}

func ToDto(systemStatistic *SystemStatistic) *SystemStatisticDto {
	if systemStatistic == nil {
		return nil
	}
	dto := &SystemStatisticDto{
		ProjectCount:              systemStatistic.ProjectCount,
		ProjectActiveCount:        systemStatistic.ProjectActiveCount,
		ProjectDeletedCount:       systemStatistic.ProjectDeletedCount,
		LicenseCount:              systemStatistic.LicenseCount,
		LicenseChartCount:         systemStatistic.LicenseChartCount,
		LicenseActiveCount:        systemStatistic.LicenseActiveCount,
		LicenseDeletedCount:       systemStatistic.LicenseDeletedCount,
		PolicyRuleCount:           systemStatistic.PolicyRuleCount,
		PolicyRuleActiveCount:     systemStatistic.PolicyRuleActiveCount,
		PolicyRuleDeletedCount:    systemStatistic.PolicyRuleDeletedCount,
		LabelCount:                systemStatistic.LabelCount,
		SchemaCount:               systemStatistic.SchemaCount,
		ObligationCount:           systemStatistic.ObligationCount,
		ObligationActiveCount:     systemStatistic.ObligationActiveCount,
		ObligationDeletedCount:    systemStatistic.ObligationDeletedCount,
		UploadFileCnt:             systemStatistic.UploadedFilesCnt,
		UploadFileCntPDF:          systemStatistic.UploadedFilesCntPDF,
		UploadFileCntJSON:         systemStatistic.UploadedFilesCntJSON,
		UploadFileCntSBOM:         systemStatistic.UploadedFilesCntSBOM,
		DbBackupFileCnt:           systemStatistic.DbBackupFilesCnt,
		UserCount:                 systemStatistic.UserCount,
		UserActiveCount:           systemStatistic.UserActiveCount,
		UserDeactivateCount:       systemStatistic.UserDeactivateCount,
		UserTermsNotAcceptedCount: systemStatistic.UserTermsNotAcceptedCount,
		UserDeprovisionedCount:    systemStatistic.UserDeprovisionedCount,

		MaxVersionsInOneProject:      systemStatistic.MaxVersionsInOneProject,
		ProjectsOverOrAtVersionLimit: systemStatistic.ProjectsOverOrAtVersionLimit,
		VersionLimit:                 systemStatistic.VersionLimit,
	}
	domain.SetBaseValues(systemStatistic, dto)
	return dto
}

func ToSystemStatisticDtoList(source []*SystemStatistic) []*SystemStatisticDto {
	return domain.MapTo(source, ToDto)
}
