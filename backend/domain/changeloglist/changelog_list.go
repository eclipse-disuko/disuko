// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package changeloglist

import "mercedes-benz.ghe.com/foss/disuko/domain"

type ChangeLogList struct {
	domain.RootEntity `bson:"inline"`

	ChangeLogs []*ChangeLog
}
