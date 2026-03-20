// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package obligation

import (
	"fmt"

	auditHelper "mercedes-benz.ghe.com/foss/disuko/helper/audit"
)

func (entity *Obligation) ToAuditString() string {
	if len(entity.Description) > 0 {
		return fmt.Sprintf("%s (Type:%s, WarnLevel:%s, Description:%s), Key: %s", entity.Name, entity.Type, entity.WarnLevel, entity.Description, entity.Key)
	}
	return fmt.Sprintf("%s (Type:%s, WarnLevel:%s), Key: %s", entity.Name, entity.Type, entity.WarnLevel, entity.Key)
}

type ObligationAudit struct {
	Created       string    `json:"created"`
	Updated       string    `json:"updated"`
	Name          string    `json:"name"`
	NameDe        string    `json:"nameDe"`
	Type          Type      `json:"type"`
	WarnLevel     WarnLevel `json:"warnLevel"`
	Description   string    `json:"description"`
	DescriptionDe string    `json:"descriptionDe"`
	AutoApproved  bool      `json:"autoApproved"`
}

func (entity *Obligation) ToAudit() *ObligationAudit {
	return &ObligationAudit{
		Created:       auditHelper.ConvertDateTime(entity.Created),
		Updated:       auditHelper.ConvertDateTime(entity.Updated),
		Name:          entity.Name,
		NameDe:        entity.NameDe,
		Type:          entity.Type,
		WarnLevel:     entity.WarnLevel,
		Description:   entity.Description,
		DescriptionDe: entity.DescriptionDe,
		AutoApproved:  entity.AutoApproved,
	}
}
