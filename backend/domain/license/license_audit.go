// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package license

import (
	"fmt"

	"github.com/eclipse-disuko/disuko/domain/obligation"
	auditHelper "github.com/eclipse-disuko/disuko/helper/audit"
	"github.com/eclipse-disuko/disuko/helper/exception"
	obligation2 "github.com/eclipse-disuko/disuko/infra/repository/obligation"
	"github.com/eclipse-disuko/disuko/logy"
)

type LicenseAudit struct {
	Key                   string
	IsDeprecatedLicenseId bool          `json:"isDeprecatedLicenseId"`
	LicenseId             string        `json:"licenseId"`
	Name                  string        `json:"name"`
	Text                  string        `json:"text"`
	Aliases               []string      `json:"aliases"`
	Meta                  MetaDataAudit `json:"meta"`
	Source                Source        `json:"source"`
	Created               string        `json:"created"`
	Updated               string        `json:"updated"`
}

type MetaDataAudit struct {
	Family         FamilyOfLicense `json:"family"`
	ApprovalState  ApprovalStatus  `json:"approvalState"`
	ReviewState    ReviewStatus    `json:"reviewState"`
	ReviewDate     string          `json:"reviewDate,omitempty"`
	Obligations    []string        `json:"obligations"`
	LicenseUrl     string          `json:"licenseUrl"`
	SourceUrl      string          `json:"sourceUrl"`
	Changelog      string          `json:"changelog"`
	LicenseType    TypeOfLicenses  `json:"licenseType"`
	Evaluation     string          `json:"evaluation"`
	LegalComments  string          `json:"legalComments"`
	IsLicenseChart bool            `json:"isLicenseChart"`
}

func (entity *MetaData) ToAudit(requestSession *logy.RequestSession, obligationProvider obligation2.IObligationRepository) *MetaDataAudit {
	obligations := make([]string, 0)
	if obligationProvider != nil {
		// with loaded obligations?
		for _, obligationKey := range entity.ObligationsKeyList {
			var o *obligation.Obligation
			exception.TryCatch(func() {
				o = obligationProvider.FindByKey(requestSession, obligationKey, false)
			}, func(exception2 exception.Exception) {
				exception.LogException(requestSession, exception2)
			})
			if o == nil {
				o = &obligation.Obligation{}
			}
			obligations = append(obligations, o.ToAuditString())
		}
	}
	return &MetaDataAudit{
		Family:         entity.Family,
		ApprovalState:  entity.ApprovalState,
		ReviewState:    entity.ReviewState,
		ReviewDate:     entity.ReviewDateStr,
		Obligations:    obligations,
		LicenseUrl:     entity.LicenseUrl,
		SourceUrl:      entity.SourceUrl,
		Changelog:      entity.Changelog,
		LicenseType:    entity.LicenseType,
		Evaluation:     entity.Evaluation,
		LegalComments:  entity.LegalComments,
		IsLicenseChart: entity.IsLicenseChart,
	}
}

func (entity *Alias) ToAudit() string {
	return fmt.Sprintf("%s(%s), Key: %s", entity.LicenseId, entity.Description, entity.Key)
}

func (entity *License) ToAudit(requestSession *logy.RequestSession, obligationProvider obligation2.IObligationRepository) *LicenseAudit {
	res := &LicenseAudit{
		Key:                   entity.Key,
		IsDeprecatedLicenseId: entity.IsDeprecatedLicenseId,
		LicenseId:             entity.LicenseId,
		Name:                  entity.Name,
		Text:                  entity.Text,
		Meta:                  *entity.Meta.ToAudit(requestSession, obligationProvider),
		Source:                entity.Source,
		Created:               auditHelper.ConvertDateTime(entity.Created),
		Updated:               auditHelper.ConvertDateTime(entity.Updated),
	}
	for _, alias := range entity.Aliases {
		res.Aliases = append(res.Aliases, alias.ToAudit())
	}
	return res
}
