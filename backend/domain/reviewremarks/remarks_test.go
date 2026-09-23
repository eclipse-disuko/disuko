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