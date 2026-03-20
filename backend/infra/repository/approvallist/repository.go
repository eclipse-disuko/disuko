// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package approvallist

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/approval"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type approvalListRepositoryStruct struct {
	base.BaseRepositoryWithSoftDelete[*approval.ApprovalList]
}

func NewApprovalListRepository(requestSession *logy.RequestSession) IApprovalListRepository {
	return &approvalListRepositoryStruct{
		BaseRepositoryWithSoftDelete: base.CreateRepositoryWithSoftDelete[*approval.ApprovalList](
			requestSession,
			ApprovalListCollectionName,
			func() *approval.ApprovalList {
				return &approval.ApprovalList{}
			},
			nil,
			nil,
			nil),
	}
}
