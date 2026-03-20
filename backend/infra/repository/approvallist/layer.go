// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package approvallist

import (
	"github.com/eclipse-disuko/disuko/domain/approval"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
)

const ApprovalListCollectionName = "approvalList"

type IApprovalListRepository interface {
	base.IBaseRepositoryWithSoftDelete[*approval.ApprovalList]
}
type IApprovalListRepositoryMigration interface {
	base.IBaseRepositoryWithHardDelete[*approval.ApprovalList]
}
