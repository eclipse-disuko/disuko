// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package newsbox

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain"
)

type ItemDto struct {
	domain.BaseDto
	Title         string    `json:"title"`
	TitleDE       string    `json:"titleDE"`
	Description   string    `json:"description"`
	DescriptionDE string    `json:"descriptionDE"`
	Image         *string   `json:"image,omitempty"`
	Link          *string   `json:"link,omitempty"`
	Expiry        time.Time `json:"expiry"`
}
type NewsboxResponse struct {
	Items  []ItemDto `json:"items"`
	ToShow int       `json:"toShow"`
}

func (i *Item) ToDto() ItemDto {
	res := ItemDto{
		Title:         i.Title,
		TitleDE:       i.TitleDE,
		Description:   i.Description,
		DescriptionDE: i.DescriptionDE,
		Image:         i.Image,
		Link:          i.Link,
		Expiry:        i.Expiry,
	}
	domain.SetBaseValues(i, &res)
	return res
}

func (d *ItemDto) ToEntity() Item {
	return Item{
		Title:         d.Title,
		TitleDE:       d.TitleDE,
		Description:   d.Description,
		DescriptionDE: d.DescriptionDE,
		Image:         d.Image,
		Link:          d.Link,
		Expiry:        d.Expiry,
	}
}
