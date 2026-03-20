// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package label

import "github.com/eclipse-disuko/disuko/domain"

type LabelRequestDto struct {
	Name        string `json:"name" validate:"required,gte=3,lte=80"`
	Description string `json:"description" validate:"lte=1000"`
	Type        string `json:"type" validate:"required,gt=1,lte=50"`
}

type LabelResponseDto struct {
	domain.BaseDto
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}
