// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package audit

import "time"

type AuditDto struct {
	Key      string    `json:"_key"`
	Message  string    `json:"message"`
	MetaJSON string    `json:"meta"`
	Created  time.Time `json:"created,omitempty"`
	User     string    `json:"user"`
}

func (entity *Audit) ToDto() AuditDto {
	return AuditDto{
		Key:      entity.Key,
		Message:  entity.Message,
		MetaJSON: entity.MetaJSON,
		Created:  entity.Created,
		User:     entity.User,
	}
}
