// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package overallreview

import (
	"errors"

	"mercedes-benz.ghe.com/foss/disuko/domain"
)

type State string

const (
	Unreviewed             State = "UNREVIEWED"
	Acceptable             State = "ACCEPTABLE"
	AcceptableAfterChanges State = "ACCEPTABLE_AFTER_CHANGES"
	Audited                State = "AUDITED"
	NotAcceptable          State = "NOT_ACCEPTABLE"
)

func (s State) Validate() error {
	switch s {
	case Unreviewed, Acceptable, AcceptableAfterChanges, Audited, NotAcceptable:
		return nil
	}
	return errors.New("unknown state")
}

type OverallReview struct {
	domain.ChildEntity `bson:",inline"`

	Creator         string
	CreatorFullName string
	Comment         string
	State           State
	SBOMId          string
	SBOMName        string
	SBOMUploaded    string
}
