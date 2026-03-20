// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package department

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain"
)

type Department struct {
	domain.RootEntity  `bson:"inline"`
	ParentDeptId       string
	ValidFrom          *time.Time
	DescriptionEnglish string
	OrgAbbreviation    string
	Skz                string
	CompanyCode        string
	CompanyName        string
	Level              int
}
