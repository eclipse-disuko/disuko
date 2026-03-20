// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package reviewremarks

import (
	"encoding/json"
	"time"

	"github.com/eclipse-disuko/disuko/domain"
)

type LevelChange struct {
	Before Level `json:"before"`
	After  Level `json:"after"`
}

type TitleChange struct {
	Before string `json:"before"`
	After  string `json:"after"`
}

type DescriptionChange struct {
	Before string `json:"before"`
	After  string `json:"after"`
}

type Comment string

type Status string

const (
	Open        Status = "OPEN"
	Closed      Status = "CLOSED"
	Cancelled   Status = "CANCELLED"
	InProgress  Status = "IN_PROGRESS"
	UnsetStatus Status = ""
)

type Level string

const (
	Green      Level = "GREEN"
	Yellow     Level = "YELLOW"
	Red        Level = "RED"
	UnsetLevel Level = ""
)

type EventType string

const (
	ClosedEvent       EventType = "CLOSED"
	CancelledEvent    EventType = "CANCELLED"
	InProgressEvent   EventType = "IN_PROGRESS"
	ReopenedEvent     EventType = "REOPENED"
	CommentEvent      EventType = "COMMENT"
	ChangedTitleEvent EventType = "CHANGED_TITLE"
	ChangedDescEvent  EventType = "CHANGED_DESCRIPTION"
	ChangedLevelEvent EventType = "CHANGED_LEVEL"
)

type Event struct {
	domain.ChildEntity `bson:",inline"`

	Type           EventType
	Author         string
	AuthorFullName string
	Content        json.RawMessage
}

type ComponentMeta struct {
	ComponentId      string
	ComponentName    string
	ComponentVersion string
}

type LicenseMeta struct {
	LicenseId   string
	LicenseName string
}

type Remark struct {
	domain.ChildEntity `bson:",inline"`

	Author string
	Origin string

	Title       string
	Level       Level
	Description string

	Status Status
	Closed *time.Time

	Events []Event

	SBOMId       string
	SBOMName     string
	SBOMUploaded *time.Time

	// Deprecated
	ComponentId string
	// Deprecated
	ComponentName string
	// Deprecated
	ComponentVersion string

	// Deprecated
	LicenseId string
	// Deprecated
	LicenseName string

	Components []ComponentMeta
	Licenses   []LicenseMeta
}

type ReviewTemplate struct {
	domain.RootEntity `bson:",inline"`
	domain.SoftDelete `bson:",inline"`

	Title       string
	Level       Level
	Description string
	Source      string
}

func (t ReviewTemplate) ToDto() *ReviewTemplateResponseDto {
	return &ReviewTemplateResponseDto{
		BaseDto: domain.BaseDto{
			Key:     t.Key,
			Rev:     t.Rev,
			Updated: t.Updated,
			Created: t.Created,
		},
		Title:       t.Title,
		Level:       string(t.Level),
		Description: t.Description,
		Source:      t.Source,
	}
}

func (t ReviewTemplate) Update(new ReviewTemplateRequestDto) *ReviewTemplate {
	return &ReviewTemplate{
		RootEntity: domain.RootEntity{
			ChildEntity: domain.ChildEntity{
				Key:       t.Key,
				Created:   t.Created,
				Updated:   time.Now(),
				Optimized: false,
			},
			Rev: t.Rev,
		},
		Title:       new.Title,
		Level:       Level(new.Level),
		Description: new.Description,
		Source:      new.Source,
	}
}

type ReviewRemarks struct {
	domain.RootEntity `bson:",inline"`
	domain.SoftDelete `bson:",inline"`

	Remarks []*Remark
}

func (r *Remark) Set(new *Remark, author, fullName string) {
	if r.Level != new.Level {
		data, _ := json.Marshal(LevelChange{
			Before: r.Level,
			After:  new.Level,
		})
		r.Events = append(r.Events, Event{
			ChildEntity:    domain.NewChildEntity(),
			Type:           ChangedLevelEvent,
			Author:         author,
			AuthorFullName: fullName,
			Content:        data,
		})
	}
	r.Level = new.Level
	if r.Title != new.Title {
		data, _ := json.Marshal(TitleChange{
			Before: r.Title,
			After:  new.Title,
		})
		r.Events = append(r.Events, Event{
			ChildEntity:    domain.NewChildEntity(),
			Type:           ChangedTitleEvent,
			Author:         author,
			AuthorFullName: fullName,
			Content:        data,
		})
	}
	r.Title = new.Title
	if r.Description != new.Description {
		data, _ := json.Marshal(DescriptionChange{
			Before: r.Description,
			After:  new.Description,
		})
		r.Events = append(r.Events, Event{
			ChildEntity:    domain.NewChildEntity(),
			Type:           ChangedDescEvent,
			Author:         author,
			AuthorFullName: fullName,
			Content:        data,
		})
	}
	r.Description = new.Description
	r.SBOMId = new.SBOMId
	r.SBOMName = new.SBOMName

	r.Components = new.Components
	r.Licenses = new.Licenses
}

type GetFullNameFn func(user string) string

func (r *Remark) Close(author, fullName string) {
	r.Events = append(r.Events, Event{
		ChildEntity:    domain.NewChildEntity(),
		Type:           ClosedEvent,
		Author:         author,
		AuthorFullName: fullName,
	})
	r.Level = Green
	r.Status = Closed
	now := time.Now()
	r.Closed = &now
}

func (r *Remark) InProgress(author, fullName string) {
	r.Events = append(r.Events, Event{
		ChildEntity:    domain.NewChildEntity(),
		Type:           InProgressEvent,
		Author:         author,
		AuthorFullName: fullName,
	})
	r.Status = InProgress
}

func (r *Remark) Cancel(author, fullName string) {
	r.Events = append(r.Events, Event{
		ChildEntity:    domain.NewChildEntity(),
		Type:           CancelledEvent,
		Author:         author,
		AuthorFullName: fullName,
	})
	r.Status = Cancelled
}

func (r *Remark) Reopen(author, fullName string) {
	r.Events = append(r.Events, Event{
		ChildEntity:    domain.NewChildEntity(),
		Type:           ReopenedEvent,
		Author:         author,
		AuthorFullName: fullName,
	})
	r.Status = Open
	r.Closed = nil
}

func (r *Remark) Comment(author, fullName, content string) {
	data, _ := json.Marshal(Comment(content))
	r.Events = append(r.Events, Event{
		ChildEntity:    domain.NewChildEntity(),
		Type:           CommentEvent,
		Author:         author,
		AuthorFullName: fullName,
		Content:        data,
	})
}
