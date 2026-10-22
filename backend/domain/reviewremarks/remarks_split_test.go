// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package reviewremarks

import "testing"

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
