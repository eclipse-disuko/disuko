// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package reviewremarks

import (
	rt "github.com/eclipse-disuko/disuko/domain/reviewremarks"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
)

type IReviewTemplateRepository interface {
	base.IBaseRepositoryWithSoftDelete[*rt.ReviewTemplate]
}
