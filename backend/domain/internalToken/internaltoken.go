// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package internalToken

import (
	"time"

	"mercedes-benz.ghe.com/foss/disuko/domain"
)

type Capability int

const (
	StatisticsCSV Capability = iota
	CustomLicenses
)

type InternalToken struct {
	domain.RootEntity `bson:"inline"`
	Name              string
	Revoked           bool
	Token             string
	Salt              string
	Description       string
	Expiry            time.Time
	Capabilities      []Capability
}
