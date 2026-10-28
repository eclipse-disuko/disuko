// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package reviewremarks

import (
	"testing"
	"time"
)

func TestSplitAndCarryOverFrom(t *testing.T) {
	tests := []struct {
		name     string
		previous []*Remark
		current  *Remark
		want     []struct {
			components []string
			licenses   []string
			status     Status
		}
	}{
		{
			name: "exact combination is carried over",
			previous: []*Remark{{
				Status:     InProgress,
				Components: comps("A", "B"),
				Licenses:   lics("L1"),
			}},
			current: &Remark{Status: Open, Components: comps("A", "B"), Licenses: lics("L1")},
			want: []struct {
				components []string
				licenses   []string
				status     Status
			}{{[]string{"A", "B"}, []string{"L1"}, InProgress}},
		},
		{
			name: "new component is split into fresh remark",
			previous: []*Remark{{
				Status:     Closed,
				Components: comps("A", "B"),
				Licenses:   lics("L1"),
			}},
			current: &Remark{Status: Open, Components: comps("A", "B", "C"), Licenses: lics("L1")},
			want: []struct {
				components []string
				licenses   []string
				status     Status
			}{
				{[]string{"A", "B"}, []string{"L1"}, Closed},
				{[]string{"C"}, nil, Open},
			},
		},
		{
			name: "smaller current set carries previous state",
			previous: []*Remark{{
				Status:     InProgress,
				Components: comps("A", "B", "C"),
				Licenses:   lics("L1"),
			}},
			current: &Remark{Status: Open, Components: comps("A", "B"), Licenses: lics("L1")},
			want: []struct {
				components []string
				licenses   []string
				status     Status
			}{{[]string{"A", "B"}, []string{"L1"}, InProgress}},
		},
		{
			name: "partial component overlap is split",
			previous: []*Remark{{
				Status:     InProgress,
				Components: comps("A", "B"),
				Licenses:   lics("L1"),
			}},
			current: &Remark{Status: Open, Components: comps("B", "C"), Licenses: lics("L1")},
			want: []struct {
				components []string
				licenses   []string
				status     Status
			}{
				{[]string{"B"}, []string{"L1"}, InProgress},
				{[]string{"C"}, nil, Open},
			},
		},
		{
			name: "components and licenses are split independently",
			previous: []*Remark{{
				Status:     Closed,
				Components: comps("A", "B"),
				Licenses:   lics("L1", "L2"),
			}},
			current: &Remark{Status: Open, Components: comps("B", "C"), Licenses: lics("L2", "L3")},
			want: []struct {
				components []string
				licenses   []string
				status     Status
			}{
				{[]string{"B"}, []string{"L2"}, Closed},
				{[]string{"C"}, []string{"L3"}, Open},
			},
		},
		{
			name: "cancelled previous state is not carried over",
			previous: []*Remark{{
				Status:     Cancelled,
				Components: comps("A", "B"),
				Licenses:   lics("L1"),
			}},
			current: &Remark{Status: Open, Components: comps("A", "B"), Licenses: lics("L1")},
			want: []struct {
				components []string
				licenses   []string
				status     Status
			}{{[]string{"A", "B"}, []string{"L1"}, Open}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.current.SplitAndCarryOverFrom(tt.previous)
			if len(got) != len(tt.want) {
				t.Fatalf("got %d remarks, want %d", len(got), len(tt.want))
			}
			for i := range tt.want {
				assertComponents(t, got[i].Components, tt.want[i].components)
				assertLicenses(t, got[i].Licenses, tt.want[i].licenses)
				if got[i].Status != tt.want[i].status {
					t.Fatalf("remark %d status = %s, want %s", i, got[i].Status, tt.want[i].status)
				}
			}
		})
	}
}

func comps(names ...string) []ComponentMeta {
	res := make([]ComponentMeta, 0, len(names))
	for _, name := range names {
		res = append(res, ComponentMeta{ComponentName: name})
	}
	return res
}

func lics(ids ...string) []LicenseMeta {
	res := make([]LicenseMeta, 0, len(ids))
	for _, id := range ids {
		res = append(res, LicenseMeta{LicenseId: id})
	}
	return res
}

func assertComponents(t *testing.T, got []ComponentMeta, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("components = %v, want %v", got, want)
	}
	for i := range want {
		if got[i].ComponentName != want[i] {
			t.Fatalf("component %d = %q, want %q", i, got[i].ComponentName, want[i])
		}
	}
}

func assertLicenses(t *testing.T, got []LicenseMeta, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("licenses = %v, want %v", got, want)
	}
	for i := range want {
		if got[i].LicenseId != want[i] {
			t.Fatalf("license %d = %q, want %q", i, got[i].LicenseId, want[i])
		}
	}
}

func TestFindRemarksForSplit_PrefersCurrentSbom(t *testing.T) {
	currentUploaded := mustTime("2026-10-03T10:00:00Z")
	previousUploaded := mustTime("2026-10-02T10:00:00Z")

	previous := &Remark{
		SBOMId:       "sbom-2",
		SBOMUploaded: &previousUploaded,
		TemplateKey:  "weak-copyleft",
		Status:       Closed,
		Components:   comps("A", "B"),
	}
	currentCancelledA := &Remark{
		SBOMId:       "sbom-3",
		SBOMUploaded: &currentUploaded,
		TemplateKey:  "weak-copyleft",
		Status:       Cancelled,
		Components:   comps("A"),
	}
	currentCancelledB := &Remark{
		SBOMId:       "sbom-3",
		SBOMUploaded: &currentUploaded,
		TemplateKey:  "weak-copyleft",
		Status:       Cancelled,
		Components:   comps("B"),
	}
	newRemark := &Remark{
		SBOMId:       "sbom-3",
		SBOMUploaded: &currentUploaded,
		TemplateKey:  "weak-copyleft",
		Status:       Open,
		Components:   comps("A", "B"),
	}

	got := FindRemarksForSplit([]*Remark{previous, currentCancelledA, currentCancelledB}, newRemark)
	if len(got) != 2 {
		t.Fatalf("got %d remarks, want 2 current-SBOM remarks", len(got))
	}
	for _, remark := range got {
		if remark.SBOMId != "sbom-3" {
			t.Fatalf("got remark from %s, want current SBOM sbom-3", remark.SBOMId)
		}
	}
}

func TestSplitAndCarryOverFrom_AllCurrentCancelledCreatesFreshAggregate(t *testing.T) {
	current := &Remark{
		Status:     Open,
		Components: comps("A", "B", "C"),
		Licenses:   lics("L1", "L2"),
	}
	currentExisting := []*Remark{
		{Status: Cancelled, Components: comps("A", "B"), Licenses: lics("L1")},
		{Status: Cancelled, Components: comps("C")},
		{Status: Cancelled, Licenses: lics("L2")},
	}

	got := current.SplitAndCarryOverFrom(currentExisting)
	if len(got) != 1 {
		t.Fatalf("got %d remarks, want one fresh aggregate remark", len(got))
	}
	assertComponents(t, got[0].Components, []string{"A", "B", "C"})
	assertLicenses(t, got[0].Licenses, []string{"L1", "L2"})
	if got[0].Status != Open {
		t.Fatalf("status = %s, want %s", got[0].Status, Open)
	}
}

func TestSplitAndCarryOverFrom_CurrentMixedStatusRecreatesOnlyCancelledPart(t *testing.T) {
	current := &Remark{
		Status:     Open,
		Components: comps("A", "B", "C"),
	}
	currentExisting := []*Remark{
		{Status: Closed, Components: comps("A", "B")},
		{Status: Cancelled, Components: comps("C")},
	}

	got := current.SplitAndCarryOverFrom(currentExisting)
	if len(got) != 2 {
		t.Fatalf("got %d remarks, want carried and fresh remark", len(got))
	}
	assertComponents(t, got[0].Components, []string{"A", "B"})
	if got[0].Status != Closed {
		t.Fatalf("carried status = %s, want %s", got[0].Status, Closed)
	}
	assertComponents(t, got[1].Components, []string{"C"})
	if got[1].Status != Open {
		t.Fatalf("fresh status = %s, want %s", got[1].Status, Open)
	}
}

func TestFindRemarksForSplit_FallsBackToLatestPreviousSbom(t *testing.T) {
	currentUploaded := mustTime("2026-10-03T10:00:00Z")
	previousUploaded := mustTime("2026-10-02T10:00:00Z")
	olderUploaded := mustTime("2026-10-01T10:00:00Z")

	existing := []*Remark{
		{SBOMId: "sbom-1", SBOMUploaded: &olderUploaded, TemplateKey: "weak-copyleft", Status: Closed},
		{SBOMId: "sbom-2", SBOMUploaded: &previousUploaded, TemplateKey: "weak-copyleft", Status: InProgress},
	}
	newRemark := &Remark{
		SBOMId:       "sbom-3",
		SBOMUploaded: &currentUploaded,
		TemplateKey:  "weak-copyleft",
	}

	got := FindRemarksForSplit(existing, newRemark)
	if len(got) != 1 {
		t.Fatalf("got %d remarks, want 1", len(got))
	}
	if got[0].SBOMId != "sbom-2" {
		t.Fatalf("got SBOM %s, want latest previous sbom-2", got[0].SBOMId)
	}
}

func mustTime(value string) time.Time {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		panic(err)
	}
	return t
}
