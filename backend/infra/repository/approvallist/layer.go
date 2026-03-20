// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package approvallist

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/approval"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
)

const ApprovalListCollectionName = "approvalList"

type IApprovalListRepository interface {
	base.IBaseRepositoryWithSoftDelete[*approval.ApprovalList]
}
type IApprovalListRepositoryMigration interface {
	base.IBaseRepositoryWithHardDelete[*approval.ApprovalList]
}
