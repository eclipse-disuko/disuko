// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package analytics

import (
	"mercedes-benz.ghe.com/foss/disuko/domain"
)

type Component struct {
	domain.RootEntity `bson:"inline"`
	domain.SoftDelete `bson:"inline"`
	Name              string
	Version           string
}
