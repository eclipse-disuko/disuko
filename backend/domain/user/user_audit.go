// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package user

import (
	"fmt"

	auditHelper "github.com/eclipse-disuko/disuko/helper/audit"
)

func (user *User) ToProjectMemberAudit() string {
	return fmt.Sprintf("%s, %s (%s), Key: %s", user.Lastname, user.Forename, user.User, user.Key)
}

type UserAudit struct {
	TermsOfUse        bool        `json:"termsOfUse"`
	TermsOfUseDate    string      `json:"termsOfUseDate"`
	TermsOfUseVersion string      `json:"termsOfUseVersion"`
	Created           string      `json:"created"`
	Updated           string      `json:"updated"`
	Active            bool        `json:"active"`
	Deprovisioned     bool        `json:"deprovisioned"`
	Tasks             []TaskAudit `json:"tasks"`
}

func (user *User) ToUserAudit() *UserAudit {
	termOfUseDate := ""
	if user.TermsOfUseDate != nil {
		termOfUseDate = auditHelper.ConvertDateTime(*user.TermsOfUseDate)
	}
	taskAudits := make([]TaskAudit, 0)
	for _, t := range user.Tasks {
		taskAudits = append(taskAudits, t.ToAudit())
	}
	return &UserAudit{
		TermsOfUse:        user.TermsOfUse,
		TermsOfUseDate:    termOfUseDate,
		TermsOfUseVersion: user.TermsOfUseVersion,
		Created:           auditHelper.ConvertDateTime(user.Created),
		Updated:           auditHelper.ConvertDateTime(user.Updated),
		Active:            user.Active,
		Deprovisioned:     !user.Deprovisioned.IsZero(),
		Tasks:             taskAudits,
	}
}
