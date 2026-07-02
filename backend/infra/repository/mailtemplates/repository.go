// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package mailtemplates

import (
	"github.com/eclipse-disuko/disuko/domain/mailtemplate"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
)

type mailTemplatesRepositoryStruct struct {
	base.BaseRepositoryWithHardDelete[*mailtemplate.MailTemplate]
}

func NewMailTemplatesRepository(requestSession *logy.RequestSession) IMailTemplatesRepository {
	return &mailTemplatesRepositoryStruct{
		BaseRepositoryWithHardDelete: base.CreateRepositoryWithHardDelete[*mailtemplate.MailTemplate](
			requestSession,
			MailTemplatesCollectionName,
			func() *mailtemplate.MailTemplate {
				return &mailtemplate.MailTemplate{}
			},
			nil,
			"",
			nil,
			nil),
	}
}
