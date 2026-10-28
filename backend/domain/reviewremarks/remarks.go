// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package reviewremarks

import (
	"encoding/json"
	"slices"
	"sort"
	"strings"
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
	Gray       Level = "GRAY"
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

const (
	fieldSeparator = "\x1f"
	itemSeparator  = "\x1e"
)

func norm(s string) string { return strings.ToLower(strings.TrimSpace(s)) }

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

	TemplateKey string
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
	r.Level = Gray
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

type RrKey struct {
	sbomId      string
	components  string
	licenses    string
	templateKey string
}

type RrCombinationKey struct {
	components  string
	licenses    string
	templateKey string
}

func (r *Remark) MakeRrKey() RrKey {
	combinationKey := r.makeCombinationKey()
	return RrKey{
		sbomId:      r.SBOMId,
		components:  combinationKey.components,
		licenses:    combinationKey.licenses,
		templateKey: combinationKey.templateKey,
	}
}

func (r *Remark) makeCombinationKey() RrCombinationKey {
	componentKeys := make([]string, 0, len(r.Components))
	for _, cmp := range r.Components {
		componentKeys = append(componentKeys, componentSetKey(cmp))
	}
	sort.Strings(componentKeys)

	licenseKeys := make([]string, 0, len(r.Licenses))
	for _, lic := range r.Licenses {
		licenseKeys = append(licenseKeys, licenseSetKey(lic))
	}
	sort.Strings(licenseKeys)

	return RrCombinationKey{
		components:  strings.Join(componentKeys, itemSeparator),
		licenses:    strings.Join(licenseKeys, itemSeparator),
		templateKey: r.TemplateKey,
	}
}

func (r *Remark) PreventsDuplicate() bool {
	return r.Status != Cancelled
}

func FindLatestMatchingRemark(existingRemarks []*Remark, newRemark *Remark) *Remark {
	if newRemark.SBOMUploaded == nil {
		return nil
	}

	newKey := newRemark.makeCombinationKey()
	var latest *Remark
	for _, existingRemark := range existingRemarks {
		if existingRemark.SBOMId == newRemark.SBOMId {
			continue
		}
		if existingRemark.SBOMUploaded == nil {
			continue
		}
		if existingRemark.SBOMUploaded.After(*newRemark.SBOMUploaded) {
			continue
		}
		if existingRemark.makeCombinationKey() != newKey {
			continue
		}
		if latest == nil || existingRemark.SBOMUploaded.After(*latest.SBOMUploaded) {
			latest = existingRemark
		}
	}
	if latest != nil && latest.Status == Cancelled {
		return nil
	}
	return latest
}

// FindRemarksForSplit returns all remarks for the same template, preferring the current SBOM.
// If none exist, it uses the latest previous SBOM containing remarks for that template.
func FindRemarksForSplit(existingRemarks []*Remark, newRemark *Remark) []*Remark {
	current := findRemarksByTemplateAndSBOM(existingRemarks, newRemark.TemplateKey, newRemark.SBOMId)
	if len(current) > 0 {
		return current
	}
	return findLatestPreviousRemarksByTemplate(existingRemarks, newRemark)
}

// Find all remarks for the given SBOM and template, regardless of
// their component/license combinations or status (including CANCELLED).
func findRemarksByTemplateAndSBOM(existingRemarks []*Remark, templateKey, sbomId string) []*Remark {
	var result []*Remark
	for _, existingRemark := range existingRemarks {
		if existingRemark.SBOMId != sbomId {
			continue
		}
		if existingRemark.TemplateKey != templateKey {
			continue
		}
		result = append(result, existingRemark)
	}
	return result
}

func findLatestPreviousRemarksByTemplate(existingRemarks []*Remark, newRemark *Remark) []*Remark {
	if newRemark.SBOMUploaded == nil || newRemark.TemplateKey == "" {
		return nil
	}

	// Find the most recent previous SBOM containing remarks for the same
	// template, excluding the current SBOM and any later uploads.
	var latestSBOMId string
	var latestUploaded *time.Time

	for _, existingRemark := range existingRemarks {
		if existingRemark.SBOMId == newRemark.SBOMId || existingRemark.SBOMUploaded == nil {
			continue
		}
		if existingRemark.SBOMUploaded.After(*newRemark.SBOMUploaded) {
			continue
		}
		if existingRemark.TemplateKey != newRemark.TemplateKey {
			continue
		}

		if latestUploaded == nil || existingRemark.SBOMUploaded.After(*latestUploaded) {
			uploaded := *existingRemark.SBOMUploaded
			latestUploaded = &uploaded
			latestSBOMId = existingRemark.SBOMId
		}
	}

	if latestUploaded == nil {
		return nil
	}

	// Collect all remarks from the latest matching previous SBOM.
	return findRemarksByTemplateAndSBOM(existingRemarks, newRemark.TemplateKey, latestSBOMId)
}

// SplitAndCarryOverFrom splits the remark based on matching components and
// licenses from previous remarks. Matching parts inherit their previous state,
// while unmatched parts form a new remark. CANCELLED remarks are not considered.
func (r *Remark) SplitAndCarryOverFrom(previousRemarks []*Remark) []*Remark {
	// No splitting possible without previous remarks or components/licenses.
	if len(previousRemarks) == 0 {
		return []*Remark{r}
	}
	if len(r.Components) == 0 && len(r.Licenses) == 0 {
		return []*Remark{r}
	}

	// Initialize the remaining components/licenses to be processed.
	remainingComponents := append([]ComponentMeta(nil), r.Components...)
	remainingLicenses := append([]LicenseMeta(nil), r.Licenses...)
	result := make([]*Remark, 0, len(previousRemarks)+1)

	// Match components/licenses against previous remarks and carry over their state.
	for _, previous := range previousRemarks {
		if previous.Status == Cancelled {
			continue
		}

		matchedComponents, restComponents := splitComponents(remainingComponents, previous.Components)
		matchedLicenses, restLicenses := splitLicenses(remainingLicenses, previous.Licenses)

		if len(matchedComponents) == 0 && len(matchedLicenses) == 0 {
			continue
		}

		carried := r.cloneWithCombination(matchedComponents, matchedLicenses)
		carried.CarryOverStateFrom(previous)
		result = append(result, carried)

		remainingComponents = restComponents
		remainingLicenses = restLicenses
	}

	// Create a new remark for any components/licenses not matched previously.
	if len(remainingComponents) > 0 || len(remainingLicenses) > 0 {
		result = append(result, r.cloneWithCombination(remainingComponents, remainingLicenses))
	}

	return result
}

// cloneWithCombination creates a new remark from a fresh checklist remark
// with the given subset of components and licenses.
func (r *Remark) cloneWithCombination(components []ComponentMeta, licenses []LicenseMeta) *Remark {
	clone := *r
	clone.ChildEntity = domain.NewChildEntity()
	clone.Components = slices.Clone(components)
	clone.Licenses = slices.Clone(licenses)
	return &clone
}

func splitComponents(current, previous []ComponentMeta) ([]ComponentMeta, []ComponentMeta) {
	previousKeys := make(map[string]struct{}, len(previous))
	for _, component := range previous {
		previousKeys[componentSetKey(component)] = struct{}{}
	}

	matched := make([]ComponentMeta, 0, len(current))
	remaining := make([]ComponentMeta, 0, len(current))
	for _, component := range current {
		if _, exists := previousKeys[componentSetKey(component)]; exists {
			matched = append(matched, component)
		} else {
			remaining = append(remaining, component)
		}
	}
	return matched, remaining
}

func splitLicenses(current, previous []LicenseMeta) ([]LicenseMeta, []LicenseMeta) {
	previousKeys := make(map[string]struct{}, len(previous))
	for _, lic := range previous {
		previousKeys[licenseSetKey(lic)] = struct{}{}
	}

	matched := make([]LicenseMeta, 0, len(current))
	remaining := make([]LicenseMeta, 0, len(current))
	for _, lic := range current {
		if _, exists := previousKeys[licenseSetKey(lic)]; exists {
			matched = append(matched, lic)
		} else {
			remaining = append(remaining, lic)
		}
	}
	return matched, remaining
}

func componentSetKey(component ComponentMeta) string {
	return norm(component.ComponentName) + fieldSeparator + norm(component.ComponentVersion)
}

func licenseSetKey(lic LicenseMeta) string {
	return norm(lic.LicenseId) + fieldSeparator + norm(lic.LicenseName)
}

func (r *Remark) CarryOverStateFrom(previous *Remark) {
	// r.Author = previous.Author // todo: Initial Author or a new one???
	r.Level = previous.Level
	r.Status = previous.Status

	if previous.Closed != nil {
		closed := *previous.Closed
		r.Closed = &closed
	} else {
		r.Closed = nil
	}

	r.Events = append([]Event(nil), previous.Events...)
}
