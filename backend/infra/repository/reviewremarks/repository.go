// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package reviewremarks

import (
	rr "mercedes-benz.ghe.com/foss/disuko/domain/reviewremarks"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

const ReviewRemarksCollectionName = "reviewremarks"

type reviewRemarksRepositoryStruct struct {
	base.BaseRepositoryWithSoftDelete[*rr.ReviewRemarks]
}

func NewReviewRemarskRepositry(requestSession *logy.RequestSession) IReviewRemarksRepository {
	return &reviewRemarksRepositoryStruct{
		BaseRepositoryWithSoftDelete: base.CreateRepositoryWithSoftDelete[*rr.ReviewRemarks](
			requestSession,
			ReviewRemarksCollectionName,
			func() *rr.ReviewRemarks {
				return &rr.ReviewRemarks{}
			},
			nil,
			nil,
			nil),
	}
}

func (repo *reviewRemarksRepositoryStruct) FindByKeyFilteredByComponentId(requestSession *logy.RequestSession, key string, spdxId string) *rr.ReviewRemarks {
	reviewRemarks := repo.FindByKey(requestSession, key, false)
	if reviewRemarks == nil {
		return nil
	}

	var filteredRemarks []*rr.Remark
	for _, remark := range reviewRemarks.Remarks {
		for _, component := range remark.Components {
			if component.ComponentId == spdxId {
				filteredRemarks = append(filteredRemarks, remark)
				break
			}
		}
	}
	reviewRemarks.Remarks = filteredRemarks
	return reviewRemarks
}
