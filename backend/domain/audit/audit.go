// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package audit

import (
	"time"

	"github.com/google/uuid"
)

type Audit struct {
	Key      string `bson:"_id" json:"_key"`
	Rev      string `json:"_rev"`
	Message  string
	MetaJSON string
	User     string
	Created  time.Time
}

func NewAudit(user string, message string, metaJSON string) *Audit {
	return &Audit{
		Key:      uuid.NewString(),
		Message:  message,
		MetaJSON: metaJSON,
		Created:  time.Now(),
		User:     user,
	}
}
