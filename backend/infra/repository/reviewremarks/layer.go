// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package reviewremarks

import (
	rr "github.com/eclipse-disuko/disuko/domain/reviewremarks"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
)

type IReviewRemarksRepository interface {
	base.IBaseRepositoryWithSoftDelete[*rr.ReviewRemarks]
	FindByKeyFilteredByComponentId(requestSession *logy.RequestSession, key string, spdxId string) *rr.ReviewRemarks
}
