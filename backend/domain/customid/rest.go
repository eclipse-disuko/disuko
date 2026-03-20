// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package customid

import "github.com/eclipse-disuko/disuko/domain"

type CustomIdDto struct {
	domain.BaseDto
	Name          string `json:"name" validate:"required,gte=3,lte=80"`
	NameDE        string `json:"nameDE" validate:"required,gte=3,lte=80"`
	Description   string `json:"description" validate:"lte=1000"`
	DescriptionDE string `json:"descriptionDE" validate:"lte=1000"`
	LinkTemplate  string `json:"linkTemplate" validate:"lte=1000"`
}

type CustomIdUsage struct {
	Count int `json:"count"`
}

func (c *CustomId) ToDto() CustomIdDto {
	res := CustomIdDto{
		Name:          c.Name,
		NameDE:        c.NameDE,
		Description:   c.Description,
		DescriptionDE: c.DescriptionDE,
		LinkTemplate:  c.LinkTemplate,
	}
	domain.SetBaseValues(c, &res)
	return res
}

func (d *CustomIdDto) ToEntity() CustomId {
	res := CustomId{
		Name:          d.Name,
		NameDE:        d.NameDE,
		Description:   d.Description,
		DescriptionDE: d.DescriptionDE,
		LinkTemplate:  d.LinkTemplate,
	}
	domain.SetBaseValues(d, &res)
	return res
}
