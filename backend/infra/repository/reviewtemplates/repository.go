// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package reviewremarks

import (
	rt "mercedes-benz.ghe.com/foss/disuko/domain/reviewremarks"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

const ReviewTemplateCollectionName = "reviewtemplate"

type reviewRemarksRepositoryStruct struct {
	base.BaseRepositoryWithSoftDelete[*rt.ReviewTemplate]
}

func NewReviewTemplateRepositry(requestSession *logy.RequestSession) IReviewTemplateRepository {
	return &reviewRemarksRepositoryStruct{
		BaseRepositoryWithSoftDelete: base.CreateRepositoryWithSoftDelete[*rt.ReviewTemplate](
			requestSession,
			ReviewTemplateCollectionName,
			func() *rt.ReviewTemplate {
				return &rt.ReviewTemplate{}
			},
			nil,
			nil,
			nil),
	}
}
