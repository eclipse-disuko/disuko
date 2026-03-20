// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package reviewremarks

import (
	rt "github.com/eclipse-disuko/disuko/domain/reviewremarks"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
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
