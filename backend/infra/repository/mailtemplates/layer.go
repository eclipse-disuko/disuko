// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package mailtemplates

import (
	"github.com/eclipse-disuko/disuko/domain/mailtemplate"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
)

const MailTemplatesCollectionName = "mailtemplates"

type IMailTemplatesRepository interface {
	base.IBaseRepositoryWithHardDelete[*mailtemplate.MailTemplate]
}
