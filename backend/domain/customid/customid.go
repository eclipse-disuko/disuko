// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package customid

import (
	"github.com/eclipse-disuko/disuko/domain"
)

type CustomId struct {
	domain.RootEntity `bson:"inline"`
	Name              string `json:"name"`
	NameDE            string `json:"nameDE"`
	Description       string `json:"description"`
	DescriptionDE     string `json:"descriptionDE"`
	LinkTemplate      string `json:"linkTemplate"`
}
