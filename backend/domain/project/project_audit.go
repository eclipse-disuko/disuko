// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"fmt"

	"github.com/eclipse-disuko/disuko/domain/label"
	"github.com/eclipse-disuko/disuko/domain/user"
	auditHelper "github.com/eclipse-disuko/disuko/helper/audit"
	"github.com/eclipse-disuko/disuko/logy"
)

type ProjectMemberAudit struct {
	Key           string
	UserProfile   string
	Created       string
	Updated       string
	Comment       string
	IsResponsible bool
}

func (entity *ProjectMemberEntity) ToAudit(requestSession *logy.RequestSession, userProvider user.IUserDtoProvider) *ProjectMemberAudit {
	userProfile := &user.User{}
	if userProvider != nil {
		foundUserProfile := userProvider.FindByUserId(requestSession, entity.UserId)
		if foundUserProfile != nil {
			userProfile = foundUserProfile
		}
	}
	return &ProjectMemberAudit{
		Key:           entity.Key,
		UserProfile:   createUserProfileString(entity.UserId, entity.UserType, userProfile.ToProjectMemberAudit()),
		Created:       auditHelper.ConvertDateTime(entity.Created),
		Updated:       auditHelper.ConvertDateTime(entity.Updated),
		IsResponsible: entity.IsResponsible,
		Comment:       entity.Comment,
	}
}

func createUserProfileString(userId string, userType UserType, userProfile string) string {
	return fmt.Sprintf("%s (%s): %s", userId, userType, userProfile)
}

type PolicyLabelsAudit map[string]string

func PolicyLabelsToAudit(in []string, labelsByKey map[string]*label.Label) PolicyLabelsAudit {
	res := make(PolicyLabelsAudit)
	for _, lk := range in {
		l, ok := labelsByKey[lk]
		if !ok {
			res[lk] = ""
			continue
		}
		res[lk] = l.Name
	}
	return res
}
