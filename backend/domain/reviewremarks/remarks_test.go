// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package reviewremarks

import "testing"

func TestRemarkCloseSetsGreenLevel(t *testing.T) {
	remark := &Remark{Level: Blue, Status: Open}
	remark.Close("user", "User")
	if remark.Level != Green || remark.Status != Closed || remark.Closed == nil {
		t.Fatalf("Close() = level %q, status %q, closed %v", remark.Level, remark.Status, remark.Closed)
	}
}

func TestRemarkReopenResetsOnlyGreenLevel(t *testing.T) {
	for _, level := range []Level{Blue, Green, Yellow, Red, Gray} {
		t.Run(string(level), func(t *testing.T) {
			remark := &Remark{Level: level, Status: Closed}
			remark.Reopen("user", "User")
			want := level
			if level == Green {
				want = Blue
			}
			if remark.Level != want || remark.Status != Open || remark.Closed != nil {
				t.Fatalf("Reopen() = level %q, status %q, closed %v; want level %q, OPEN, nil", remark.Level, remark.Status, remark.Closed, want)
			}
		})
	}
}

func TestParseLevelSupportsBlueAndExistingLevels(t *testing.T) {
	for _, level := range []Level{Blue, Green, Yellow, Red, Gray} {
		valid, parsed := ParseLevel(string(level))
		if !valid || parsed != level {
			t.Errorf("ParseLevel(%q) = %t, %q", level, valid, parsed)
		}
	}
	for _, invalid := range []string{"", "blue", "UNKNOWN"} {
		valid, parsed := ParseLevel(invalid)
		if valid || parsed != UnsetLevel {
			t.Errorf("ParseLevel(%q) = %t, %q; want false, unset", invalid, valid, parsed)
		}
	}
}

func TestRemarkCancelSetsGrayLevel(t *testing.T) {
	remark := &Remark{Level: Red, Status: Open}

	remark.Cancel("user", "User")

	if remark.Status != Cancelled {
		t.Errorf("Status = %q, want %q", remark.Status, Cancelled)
	}
	if remark.Level != Gray {
		t.Errorf("Level = %q, want %q", remark.Level, Gray)
	}
}

func TestRemarkHasComments(t *testing.T) {
	tests := []struct {
		name   string
		remark *Remark
		want   bool
	}{
		{name: "nil remark", remark: nil, want: false},
		{name: "no events", remark: &Remark{}, want: false},
		{name: "non-comment event", remark: &Remark{Events: []Event{{Type: InProgressEvent}}}, want: false},
		{name: "comment event", remark: &Remark{Events: []Event{{Type: CommentEvent}}}, want: true},
		{
			name: "comment among other events",
			remark: &Remark{Events: []Event{
				{Type: InProgressEvent},
				{Type: CommentEvent},
			}},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.remark.HasComments(); got != tt.want {
				t.Errorf("HasComments() = %t, want %t", got, tt.want)
			}
		})
	}
}
