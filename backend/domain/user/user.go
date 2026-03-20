// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package user

import (
	"time"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/audit"
	"github.com/eclipse-disuko/disuko/logy"
)

type IUserDtoProvider interface {
	FindByUserId(requestSession *logy.RequestSession, name string) *User
}

type MetaData struct {
	CompanyIdentifier     string
	Department            string
	DepartmentDescription string
}

func NewMetaData(companyIdentifier, department, departmentDescription string) *MetaData {
	return &MetaData{
		CompanyIdentifier:     companyIdentifier,
		Department:            department,
		DepartmentDescription: departmentDescription,
	}
}

type User struct {
	domain.RootEntity `bson:",inline"`
	audit.Container   `bson:",inline"`
	User              string
	Forename          string
	Lastname          string
	Email             string
	TermsOfUse        bool
	TermsOfUseDate    *time.Time
	TermsOfUseVersion string
	Tasks             []Task
	Roles             []string
	MetaData          *MetaData
	Active            bool
	IsInternal        bool
	Deprovisioned     time.Time
	NewsboxLastSeenId string
}

func (m *MetaData) Equal(cmp *MetaData) bool {
	if m == nil || cmp == nil {
		return m == cmp
	}
	return m.CompanyIdentifier == cmp.CompanyIdentifier && m.Department == cmp.Department && m.DepartmentDescription == cmp.DepartmentDescription
}

func (entity *User) UpdateData(data UserRequestDto) {
	entity.Email = data.Email
	entity.Forename = data.Forename
	entity.Lastname = data.Lastname
	entity.User = data.User
}

func (entity *User) UpdateNewsboxLastSeenData(data UserLastSeenDto) {
	entity.NewsboxLastSeenId = data.NewsboxLastSeenId
}

func CreateUser(forename string, lastname string, username string, email string, roles []string, metaData *MetaData, isInternal bool) *User {
	return &User{
		RootEntity:        domain.NewRootEntity(),
		User:              username,
		Forename:          forename,
		Lastname:          lastname,
		Email:             email,
		TermsOfUse:        false,
		TermsOfUseDate:    nil,
		TermsOfUseVersion: conf.Config.Server.TermsOfUseCurrentVersion,
		Roles:             roles,
		MetaData:          metaData,
		Active:            conf.IsProdEnv() || conf.Config.Server.Env == "local",
		IsInternal:        isInternal,
	}
}
