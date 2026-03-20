// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package department

import (
	"time"
)

type DepartmentDto struct {
	DeptId             string     `json:"deptId"`
	ParentDeptId       string     `json:"parentDeptId"`
	ValidFrom          *time.Time `json:"validFrom"`
	DescriptionEnglish string     `json:"descriptionEnglish"`
	OrgAbbreviation    string     `json:"orgAbbreviation"`
	Skz                string     `json:"skz"`
	CompanyCode        string     `json:"companyCode"`
	CompanyName        string     `json:"companyName"`
	Level              int        `json:"level"`
}

func (entity *Department) ToDto() *DepartmentDto {
	return &DepartmentDto{
		DeptId:             entity.Key,
		ParentDeptId:       entity.ParentDeptId,
		ValidFrom:          entity.ValidFrom,
		DescriptionEnglish: entity.DescriptionEnglish,
		OrgAbbreviation:    entity.OrgAbbreviation,
		Skz:                entity.Skz,
		CompanyCode:        entity.CompanyCode,
		CompanyName:        entity.CompanyName,
		Level:              entity.Level,
	}
}
