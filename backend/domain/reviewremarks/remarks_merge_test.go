// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package reviewremarks

import (
	"encoding/json"
	"testing"
	"time"
)

func TestFindRemarkForMerge(t *testing.T) {
	t1 := mustTime("2026-10-01T10:00:00Z")
	t2 := mustTime("2026-10-02T10:00:00Z")
	t3 := mustTime("2026-10-03T10:00:00Z")
	t4 := mustTime("2026-10-04T10:00:00Z")
	makeRemark := func(sbom string, uploaded *time.Time, template string, status Status) *Remark {
		return &Remark{SBOMId: sbom, SBOMUploaded: uploaded, TemplateKey: template, Status: status}
	}
	older := makeRemark("sbom1", &t1, "template", Closed)
	previous := makeRemark("sbom2", &t2, "template", InProgress)
	current := makeRemark("sbom3", &t3, "template", Open)
	previousCancelled := makeRemark("sbom2", &t2, "template", Cancelled)
	currentCancelled := makeRemark("sbom3", &t3, "template", Cancelled)
	otherTemplate := makeRemark("sbom3", &t3, "other", Closed)
	future := makeRemark("sbom4", &t4, "template", Closed)
	secondActive := makeRemark("sbom2", &t2, "template", Open)
	newRemark := makeRemark("sbom3", &t3, "template", Open)
	tests := []struct {
		name     string
		existing []*Remark
		want     *Remark
	}{
		{"no existing remarks", nil, nil},
		{"latest previous SBOM", []*Remark{older, previous}, previous},
		{"current SBOM preferred", []*Remark{older, previous, current}, current},
		{"other template ignored", []*Remark{otherTemplate, previous}, previous},
		{"cancelled previous only", []*Remark{older, previousCancelled}, nil},
		{"cancelled current prevents older fallback", []*Remark{previous, currentCancelled}, nil},
		{"one active and one cancelled in previous SBOM", []*Remark{previousCancelled, previous}, previous},
		{"one active and one cancelled in current SBOM", []*Remark{previous, currentCancelled, current}, current},
		{"multiple active previous remarks are ambiguous", []*Remark{previous, secondActive}, nil},
		{"future SBOM ignored", []*Remark{older, future}, older},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindRemarkForMerge(tt.existing, newRemark)
			if got != tt.want {
				t.Fatalf("FindRemarkForMerge() = %p (%+v), want %p (%+v)", got, got, tt.want, tt.want)
			}
		})
	}
}

func TestMergeCarriesStateButKeepsNewCombination(t *testing.T) {
	previousUploaded := mustTime("2026-10-02T10:00:00Z")
	currentUploaded := mustTime("2026-10-03T10:00:00Z")
	closedAt := mustTime("2026-10-02T14:00:00Z")
	previous := &Remark{
		SBOMId: "sbom2", SBOMUploaded: &previousUploaded, TemplateKey: "template",
		Status: Closed, Level: Green, Author: "reviewer", Closed: &closedAt,
		Components: comps("A", "B"), Licenses: lics("L1"),
		Events: []Event{{Type: CommentEvent}},
	}
	current := &Remark{
		SBOMId: "sbom3", SBOMUploaded: &currentUploaded, TemplateKey: "template",
		Status: Open, Components: comps("A", "C"), Licenses: lics("L2"),
	}
	source := FindRemarkForMerge([]*Remark{previous}, current)
	if source != previous {
		t.Fatal("expected latest matching template source regardless of component/license changes")
	}
	current.CarryOverStateFrom(source)
	assertComponents(t, current.Components, []string{"A", "C"})
	assertLicenses(t, current.Licenses, []string{"L2"})
	if current.Status != Closed || current.Level != Green || current.Author != "reviewer" || current.Closed == nil || !current.Closed.Equal(closedAt) {
		t.Fatalf("state not carried over: %+v", current)
	}
	if len(current.Events) != 2 || current.Events[0].Type != CommentEvent || current.Events[1].Type != CarriedOverEvent {
		t.Fatalf("expected inherited comment and CARRIED_OVER event, got %+v", current.Events)
	}
	if current.Events[1].Author != "SYSTEM" {
		t.Fatalf("carry-over author = %q, want SYSTEM", current.Events[1].Author)
	}
	var info CarryOverInfo
	if err := json.Unmarshal(current.Events[1].Content, &info); err != nil {
		t.Fatalf("invalid carry-over info: %v", err)
	}
	if len(previous.Events) != 1 || len(previous.Components) != 2 {
		t.Fatal("source remark must not be modified")
	}
}

func TestMergeSameSBOMDoesNotAddCarryOverEvent(t *testing.T) {
	uploaded := mustTime("2026-10-03T10:00:00Z")
	previous := &Remark{SBOMId: "sbom3", SBOMUploaded: &uploaded, TemplateKey: "template", Status: InProgress, Events: []Event{{Type: CommentEvent}}}
	incoming := &Remark{SBOMId: "sbom3", SBOMUploaded: &uploaded, TemplateKey: "template", Status: Open, Components: comps("A", "B")}
	source := FindRemarkForMerge([]*Remark{previous}, incoming)
	if source != previous {
		t.Fatal("expected current-SBOM remark")
	}
	incoming.CarryOverStateFrom(source)
	if incoming.Status != InProgress || len(incoming.Events) != 1 || incoming.Events[0].Type != CommentEvent {
		t.Fatalf("same SBOM should reuse state without a new carry-over event: %+v", incoming)
	}
	assertComponents(t, incoming.Components, []string{"A", "B"})
}
