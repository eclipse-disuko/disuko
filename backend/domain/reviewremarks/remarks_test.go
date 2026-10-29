// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package reviewremarks

import "testing"

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
