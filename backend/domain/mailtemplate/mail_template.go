// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package mailtemplate

import (
	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/audit"
)

type MailTemplateKey string

const (
	MailTemplateKeyApprovalFinalized MailTemplateKey = "approvalFinalized"
	MailTemplateKeyApprovalInactive  MailTemplateKey = "approvalInactiveMail"
	MailTemplateKeyDummyDeletion     MailTemplateKey = "dummyDeletion"
	MailTemplateKeyReviewCreated     MailTemplateKey = "reviewCreated"
	MailTemplateKeyReviewFinalized   MailTemplateKey = "reviewFinalized"
	MailTemplateKeySpdxUploaded      MailTemplateKey = "spdxUploaded"
	MailTemplateKeyTaskApproval      MailTemplateKey = "taskApproval"
)

type MailTemplate struct {
	domain.RootEntity `bson:"inline"`
	audit.Container   `bson:"inline"`
	Subject           string            `json:"subject"`
	Message           string            `json:"message"`
	Bcc               string            `json:"bcc"`
	Cc                string            `json:"cc"`
	Values            map[string]string `json:"values"`
}

type MailTemplateDto struct {
	Key     string            `json:"_key"`
	Subject string            `json:"subject"`
	Message string            `json:"message"`
	Bcc     string            `json:"bcc"`
	Cc      string            `json:"cc"`
	Values  map[string]string `json:"values"`
}

type UpdateMailTemplateDto struct {
	Subject string `json:"subject"`
	Message string `json:"message"`
	Bcc     string `json:"bcc"`
	Cc      string `json:"cc"`
}

func NewMailTemplate() *MailTemplate {
	return &MailTemplate{
		RootEntity: domain.NewRootEntity(),
		Values:     make(map[string]string),
	}
}

func (entity *MailTemplate) ToDto() *MailTemplateDto {
	return &MailTemplateDto{
		Key:     entity.Key,
		Subject: entity.Subject,
		Message: entity.Message,
		Bcc:     entity.Bcc,
		Cc:      entity.Cc,
		Values:  entity.Values,
	}
}

func ToDto(entity *MailTemplate) *MailTemplateDto {
	return entity.ToDto()
}

func ToDtos(source []*MailTemplate) []*MailTemplateDto {
	return domain.MapTo(source, ToDto)
}
