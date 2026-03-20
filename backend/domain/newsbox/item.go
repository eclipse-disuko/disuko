// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package newsbox

import (
	"time"

	"mercedes-benz.ghe.com/foss/disuko/domain"
)

type Item struct {
	domain.RootEntity `bson:"inline"`
	Title             string
	TitleDE           string
	Description       string
	DescriptionDE     string
	Image             *string
	Link              *string
	Expiry            time.Time
}
