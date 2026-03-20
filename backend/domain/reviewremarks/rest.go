// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package reviewremarks

import (
	"encoding/json"
	"time"

	"mercedes-benz.ghe.com/foss/disuko/domain"
)

type EventDto struct {
	Key     string    `json:"key"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`

	Type           string          `json:"type"`
	Author         string          `json:"author"`
	AuthorFullName string          `json:"authorFullName"`
	Content        json.RawMessage `json:"content,omitempty"`
}

type ComponentMetaDto struct {
	ComponentId      string `json:"componentId"`
	ComponentName    string `json:"componentName"`
	ComponentVersion string `json:"componentVersion"`
}

func (c *ComponentMeta) ToDto() ComponentMetaDto {
	return ComponentMetaDto{
		ComponentId:      c.ComponentId,
		ComponentName:    c.ComponentName,
		ComponentVersion: c.ComponentVersion,
	}
}

func (r *Remark) ComponentsDto() (res []ComponentMetaDto) {
	for _, c := range r.Components {
		res = append(res, c.ToDto())
	}
	return
}

type LicenseMetaDto struct {
	LicenseId   string `json:"licenseId"`
	LicenseName string `json:"licenseName"`
}

func (l *LicenseMeta) ToDto() LicenseMetaDto {
	return LicenseMetaDto{
		LicenseId:   l.LicenseId,
		LicenseName: l.LicenseName,
	}
}

func (r *Remark) LicensesDto() (res []LicenseMetaDto) {
	for _, l := range r.Licenses {
		res = append(res, l.ToDto())
	}
	return
}

type ReviewRemarkRequestDto struct {
	Title       string   `json:"title" validate:"required,gte=3,lte=80" example:"CSS missing"`
	Description string   `json:"description" validate:"required,gte=10,lte=700" example:"CSS missing please add"`
	Level       string   `json:"level"`
	SBOMId      string   `json:"sbomId"`
	Components  []string `json:"components"`
	Licenses    []string `json:"Licenses"`
}

type CommentRequestDto struct {
	Content string `json:"content" validate:"required,gte=3,lte=500" example:"Problem solved"`
}

type SetRemarkStatusRequestDto struct {
	Status string `json:"status" validate:"required,oneof=OPEN CLOSED CANCELLED IN_PROGRESS" example:"CLOSED"`
}

type BulkReviewRemarkStatusRequest struct {
	RemarkKeys []string `json:"remarkKeys" validate:"required,min=1,max=100,dive,uuid4"`
	Status     string   `json:"status" validate:"required,oneof=CLOSED CANCELLED"`
}

type RemarkDto struct {
	Key     string    `json:"key"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`

	Author string `json:"author"`
	Origin string `json:"origin"`

	Title       string `json:"title"`
	Level       string `json:"level"`
	Description string `json:"description"`

	Status string     `json:"status"`
	Closed *time.Time `json:"closed"`

	Events []EventDto `json:"events"`

	SBOMId       string     `json:"sbomId"`
	SBOMName     string     `json:"sbomName"`
	SBOMUploaded *time.Time `json:"sbomUploaded,omitempty"`

	Components []ComponentMetaDto `json:"components"`
	Licenses   []LicenseMetaDto   `json:"licenses"`
}

// RemarkDto is a public JSON representation for review remark
type RemarkDtoExternV1 struct {
	Key     string    `json:"key"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`

	Author string `json:"author"`
	Origin string `json:"origin"`

	Title       string `json:"title"`
	Level       string `json:"level"`
	Description string `json:"description"`

	Status string     `json:"status"`
	Closed *time.Time `json:"closed"`

	Events []EventDto `json:"events"`

	SBOMId       string     `json:"sbomId"`
	SBOMName     string     `json:"sbomName"`
	SBOMUploaded *time.Time `json:"sbomUploaded,omitempty"`

	ComponentId      string `json:"componentId"`
	ComponentName    string `json:"componentName"`
	ComponentVersion string `json:"componentVersion"`
	LicenseId        string `json:"licenseId"`
	LicenseName      string `json:"licenseName"`
} //	@name	RemarkDto

// RemarkDto is a public JSON representation for review remark
type RemarkDtoExternV2 struct {
	Key     string    `json:"key"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`

	Author string `json:"author"`
	Origin string `json:"origin"`

	Title       string `json:"title"`
	Level       string `json:"level"`
	Description string `json:"description"`

	Status string     `json:"status"`
	Closed *time.Time `json:"closed"`

	Events []EventDto `json:"events"`

	SBOMId       string     `json:"sbomId"`
	SBOMName     string     `json:"sbomName"`
	SBOMUploaded *time.Time `json:"sbomUploaded,omitempty"`

	Components []ComponentMetaDto `json:"components"`
	Licenses   []LicenseMetaDto   `json:"licenses"`
} //	@name	RemarkDtoV2

type RRCommentExternDTO struct {
	Content string `json:"content" validate:"required,gte=3,lte=500" example:"Problem solved"`
} //	@name	RRCommentExternDTO

func (e *Event) ToDto() EventDto {
	return EventDto{
		Key:            e.Key,
		Created:        e.Created,
		Updated:        e.Updated,
		Type:           string(e.Type),
		Author:         e.Author,
		AuthorFullName: e.AuthorFullName,
		Content:        e.Content,
	}
}

type ReviewTemplateRequestDto struct {
	Title       string `json:"title"`
	Level       string `json:"level"`
	Description string `json:"description"`
	Source      string `json:"source"`
}

type ReviewTemplateResponseDto struct {
	domain.BaseDto

	Title       string `json:"title"`
	Level       string `json:"level"`
	Description string `json:"description"`
	Source      string `json:"source"`
}

func (r *Remark) EventsDto() (res []EventDto) {
	for _, c := range r.Events {
		res = append(res, c.ToDto())
	}
	return
}

func (r *Remark) ToDto() RemarkDto {
	return RemarkDto{
		Key:          r.Key,
		Created:      r.Created,
		Updated:      r.Updated,
		Author:       r.Author,
		Origin:       r.Origin,
		Title:        r.Title,
		Level:        string(r.Level),
		Description:  r.Description,
		Status:       string(r.Status),
		Closed:       r.Closed,
		Events:       r.EventsDto(),
		SBOMId:       r.SBOMId,
		SBOMName:     r.SBOMName,
		SBOMUploaded: r.SBOMUploaded,
		Components:   r.ComponentsDto(),
		Licenses:     r.LicensesDto(),
	}
}

func (r *Remark) ToExternV2Dto() RemarkDtoExternV2 {
	return RemarkDtoExternV2{
		Key:          r.Key,
		Created:      r.Created,
		Updated:      r.Updated,
		Author:       r.Author,
		Origin:       r.Origin,
		Title:        r.Title,
		Level:        string(r.Level),
		Description:  r.Description,
		Status:       string(r.Status),
		Closed:       r.Closed,
		Events:       r.EventsDto(),
		SBOMId:       r.SBOMId,
		SBOMName:     r.SBOMName,
		SBOMUploaded: r.SBOMUploaded,
		Components:   r.ComponentsDto(),
		Licenses:     r.LicensesDto(),
	}
}

func (r *Remark) ToExternV1Dto() RemarkDtoExternV1 {
	res := RemarkDtoExternV1{
		Key:          r.Key,
		Created:      r.Created,
		Updated:      r.Updated,
		Author:       r.Author,
		Origin:       r.Origin,
		Title:        r.Title,
		Level:        string(r.Level),
		Description:  r.Description,
		Status:       string(r.Status),
		Closed:       r.Closed,
		Events:       r.EventsDto(),
		SBOMId:       r.SBOMId,
		SBOMName:     r.SBOMName,
		SBOMUploaded: r.SBOMUploaded,
	}
	if len(r.Components) > 0 {
		res.ComponentId = r.Components[0].ComponentId
		res.ComponentName = r.Components[0].ComponentName
		res.ComponentVersion = r.Components[0].ComponentVersion
	}
	if len(r.Licenses) > 0 {
		res.LicenseId = r.Licenses[0].LicenseId
		res.LicenseName = r.Licenses[0].LicenseName
	}
	return res
}

func ParseLevel(level string) (valid bool, res Level) {
	switch level {
	case string(Green):
		valid, res = true, Green
	case string(Yellow):
		valid, res = true, Yellow
	case string(Red):
		valid, res = true, Red
	default:
		valid, res = false, UnsetLevel
	}
	return
}

func ParseStatus(status string) (valid bool, res Status) {
	switch status {
	case string(Open):
		valid, res = true, Open
	case string(Closed):
		valid, res = true, Closed
	case string(Cancelled):
		valid, res = true, Cancelled
	case string(InProgress):
		valid, res = true, InProgress
	default:
		valid, res = false, UnsetStatus
	}
	return
}
