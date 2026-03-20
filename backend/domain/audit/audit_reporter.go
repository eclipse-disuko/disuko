// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package audit

import (
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
)

type DiffReporter struct {
	path  cmp.Path
	diffs []string
}

func (r *DiffReporter) PushStep(ps cmp.PathStep) {
	r.path = append(r.path, ps)
}

func (r *DiffReporter) Report(rs cmp.Result) {
	if !rs.Equal() {
		vx, vy := r.path.Last().Values()
		vxString := fmt.Sprintf("%+v", vx)
		vyString := fmt.Sprintf("%+v", vy)
		if len(vxString) == 0 {
			r.diffs = append(r.diffs, fmt.Sprintf("%#v:\n\t+: %s\n", r.path, vyString))
		} else if len(vyString) == 0 {
			r.diffs = append(r.diffs, fmt.Sprintf("%#v:\n\t-: %s\n", r.path, vxString))
		} else {
			r.diffs = append(r.diffs, fmt.Sprintf("%#v:\n\t-: %s\n\t+: %s\n", r.path, vxString, vyString))
		}
	}
}

func (r *DiffReporter) PopStep() {
	r.path = r.path[:len(r.path)-1]
}

func (r *DiffReporter) String() string {
	return strings.Join(r.diffs, "\n")
}

func DiffWithReporter(before, after interface{}, options ...cmp.Option) string {
	var r DiffReporter
	reporterOption := cmp.Reporter(&r)
	cmp.Equal(before, after, append(options, reporterOption)...)
	return r.String()
}
