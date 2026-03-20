// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package announcement

import (
	"encoding/json"
	"time"

	"github.com/eclipse-disuko/disuko/domain"
)

const (
	TYPE_LICENSE = "license_change"
)

type Announcement struct {
	domain.RootEntity `bson:"inline"`
	When              time.Time
	Type              string
	Content           string
}

type LicenseAnnouncementContent struct {
	LicenseName string `json:"licenseName"`
	LicenseId   string `json:"licenseId"`
	ChangeType  string `json:"changeType"`
	OldVal      string `json:"oldVal"`
	NewVal      string `json:"newVal"`
}
type LicenseAnnouncement struct {
	Announcement
	LicenseAnnouncementContent
}

type AnnouncementDto struct {
	Key     string    `json:"_key"`
	When    time.Time `json:"when"`
	Type    string    `json:"type"`
	Content string    `json:"content"`
}

type AnnouncementGenerationResultDto struct {
	ReqID              string
	AnnouncementsCount int
	Errors             []string
}

func NewLicenseAnnouncement() *LicenseAnnouncement {
	return &LicenseAnnouncement{
		Announcement: Announcement{
			RootEntity: domain.NewRootEntity(),
			When:       time.Now(),
			Type:       TYPE_LICENSE,
		},
		LicenseAnnouncementContent: LicenseAnnouncementContent{
			LicenseName: "",
			LicenseId:   "",
			OldVal:      "",
			NewVal:      "",
			ChangeType:  "",
		},
	}
}

func (entity *LicenseAnnouncement) ToAnnouncement() (*Announcement, error) {
	content := LicenseAnnouncementContent{
		LicenseName: entity.LicenseName,
		LicenseId:   entity.LicenseId,
		ChangeType:  entity.ChangeType,
		OldVal:      entity.OldVal,
		NewVal:      entity.NewVal,
	}

	contentJson, err := json.Marshal(content)

	return &Announcement{
		RootEntity: domain.NewRootEntity(),
		When:       entity.When,
		Type:       entity.Type,
		Content:    string(contentJson),
	}, err
}

func (entity *Announcement) ToDto() *AnnouncementDto {
	return ToDto(entity)
}

func ToDto(entity *Announcement) *AnnouncementDto {
	return &AnnouncementDto{
		Key:     entity.Key,
		When:    entity.When,
		Type:    entity.Type,
		Content: entity.Content,
	}
}

func ToDtos(source []*Announcement) []*AnnouncementDto {
	return domain.MapTo(source, ToDto)
}
