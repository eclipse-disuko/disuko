// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package user

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain"
)

type MetaDataDto struct {
	CompanyIdentifier     string `json:"companyIdentifier"`
	Department            string `json:"department"`
	DepartmentDescription string `json:"departmentDescription"`
}

func (metaData *MetaData) ToDto() *MetaDataDto {
	return &MetaDataDto{
		CompanyIdentifier:     metaData.CompanyIdentifier,
		Department:            metaData.Department,
		DepartmentDescription: metaData.DepartmentDescription,
	}
}

type UserDto struct {
	domain.BaseDto
	User              string       `json:"user"`
	Forename          string       `json:"forename"`
	Lastname          string       `json:"lastname"`
	Email             string       `json:"email"`
	TermsOfUse        bool         `json:"termsOfUse"`
	TermsOfUseDate    *time.Time   `json:"termsOfUseDate"`
	TermsOfUseVersion string       `json:"termsOfUseVersion"`
	Roles             []string     `json:"roles"`
	MetaData          *MetaDataDto `json:"metaData"`
	Active            bool         `json:"active"`
	IsInternal        bool         `json:"isInternal"`
	Deprovisioned     *time.Time   `json:"deprovisioned"`
}

func (d UserDto) GetDepartment() string {
	if d.MetaData != nil {
		return d.MetaData.Department
	}
	return ""
}

func (d UserDto) GetDepartmentDescription() string {
	if d.MetaData != nil {
		return d.MetaData.DepartmentDescription
	}
	return ""
}

func (d UserDto) GetCompanyIdentifier() string {
	if d.MetaData != nil {
		return d.MetaData.CompanyIdentifier
	}
	return ""
}

func (d UserDto) GetTermsOfUseDate() int64 {
	if d.TermsOfUseDate != nil {
		return d.TermsOfUseDate.Unix()
	}
	return 0
}

func (d UserDto) GetDeprovisioned() int64 {
	if d.Deprovisioned != nil {
		return d.Deprovisioned.Unix()
	}
	return 0
}

func (entity *User) ToDto() *UserDto {
	dto := &UserDto{
		User:              entity.User,
		Forename:          entity.Forename,
		Lastname:          entity.Lastname,
		Email:             entity.Email,
		TermsOfUse:        entity.TermsOfUse,
		TermsOfUseDate:    entity.TermsOfUseDate,
		TermsOfUseVersion: entity.TermsOfUseVersion,
		Roles:             entity.Roles,
		Active:            entity.Active,
		IsInternal:        entity.IsInternal,
	}
	if entity.MetaData != nil {
		dto.MetaData = entity.MetaData.ToDto()
	}
	if !entity.Deprovisioned.IsZero() {
		dto.Deprovisioned = &entity.Deprovisioned
	}
	domain.SetBaseValues(entity, dto)
	return dto
}

func (entity *User) GetTask(guid string, taskType TaskType, status TaskStatus) *Task {
	for i := range entity.Tasks {
		if entity.Tasks[i].TargetGuid == guid {
			if taskType != "" && entity.Tasks[i].Type != taskType {
				continue
			}
			if status != "" && entity.Tasks[i].Status != status {
				continue
			}
			return &entity.Tasks[i]
		}
	}
	return nil
}

func ToDtos(users []*User) []*UserDto {
	result := make([]*UserDto, 0)
	for _, entity := range users {
		result = append(result, entity.ToDto())
	}
	return result
}

type AllResponse struct {
	Users []*UserDto `json:"items"`
	Count int        `json:"count"`
}

type UserRequestDto struct {
	Key            string     `json:"_key" validate:"lte=36"`
	User           string     `json:"user" validate:"lte=80,RealUser"`
	Forename       string     `json:"forename" validate:"lte=80"`
	Lastname       string     `json:"lastname" validate:"lte=80"`
	Email          string     `json:"email" validate:"lte=80"`
	TermsOfUse     bool       `json:"termsOfUse"`
	TermsOfUseDate *time.Time `json:"TermsOfUseDate"`
	Active         bool       `json:"active"`
}

type UserLastSeenDto struct {
	NewsboxLastSeenId string `json:"newsboxLastSeenId"`
}

type UserRolesRequestDto struct {
	Roles []string `json:"roles"`
}

type UserMailDto struct {
	User     string `json:"user" validate:"lte=80,RealUser"`
	Forename string `json:"forename" validate:"lte=80"`
	Lastname string `json:"lastname" validate:"lte=80"`
	Email    string `json:"email" validate:"lte=80"`
}

type DeletePersonalDataDto struct {
	Username   string `json:"username" validate:"required"`
	DryRun     bool   `json:"dry_run"`
	EntityType string `json:"entity_type"`
}

func (entity *User) ToUserMailDto() *UserMailDto {
	return &UserMailDto{
		User:     entity.User,
		Forename: entity.Forename,
		Lastname: entity.Lastname,
		Email:    entity.Email,
	}
}
