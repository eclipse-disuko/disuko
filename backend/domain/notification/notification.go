// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package notification

import "github.com/eclipse-disuko/disuko/domain"

type Notification struct {
	domain.RootEntity `bson:"inline"`
	Text              string `json:"text"`
	Enabled           bool   `json:"enabled"`
}

type NotificationDto struct {
	Text    string `json:"text"`
	Enabled bool   `json:"enabled"`
	Rev     string `json:"rev"`
}

func (entity *Notification) ToDto() *NotificationDto {
	return &NotificationDto{
		Text:    entity.Text,
		Enabled: entity.Enabled,
		Rev:     entity.Rev,
	}
}
