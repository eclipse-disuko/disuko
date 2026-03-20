// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package reviewremarks

import (
	rr "mercedes-benz.ghe.com/foss/disuko/domain/reviewremarks"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type IReviewRemarksRepository interface {
	base.IBaseRepositoryWithSoftDelete[*rr.ReviewRemarks]
	FindByKeyFilteredByComponentId(requestSession *logy.RequestSession, key string, spdxId string) *rr.ReviewRemarks
}
