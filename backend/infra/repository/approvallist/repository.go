// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package approvallist

import (
	"github.com/eclipse-disuko/disuko/domain/approval"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
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
