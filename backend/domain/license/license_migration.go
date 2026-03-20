// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package license

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain/obligation"
)

type MetaDataMigration struct {
	Family          FamilyOfLicense         `json:"family"`
	ApprovalState   ApprovalStatus          `json:"approvalState"`
	ReviewState     ReviewStatus            `json:"reviewState"`
	ReviewDate      time.Time               `json:"reviewDate,omitempty"`
	ObligationsList []obligation.Obligation `json:"obligationsList"`
	LicenseUrl      string                  `json:"licenseUrl"`
	OSIApproved     bool                    `json:"osiApproved"`
	FSFApproved     bool                    `json:"fsfApproved"`
	Changelog       string                  `json:"changelog"`
	LicenseType     string                  `json:"licenseType"`
	Evaluation      string                  `json:"evaluation"`
}
type LicenseMigration struct {
	Key       string            `json:"_key"`
	Rev       string            `json:"_rev"`
	LicenseId string            `json:"licenseId"`
	Name      string            `json:"name"`
	Text      string            `json:"text"`
	Active    bool              `json:"active"`
	Meta      MetaDataMigration `json:"meta"`
	Source    Source            `json:"source"`
	Created   time.Time         `json:"created,omitempty"`
	Updated   time.Time         `json:"updated,omitempty"`

	/*
		This means, that the license was not complete loaded!
		Do not save this license otherwise, we will delete the not loaded data in the database.
	*/
	Optimized bool `bson: "-" json:"-"`
}
