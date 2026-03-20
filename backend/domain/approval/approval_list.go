// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package approval

import "mercedes-benz.ghe.com/foss/disuko/domain"

type ApprovalList struct {
	domain.RootEntity `bson:"inline"`
	domain.SoftDelete `bson:"inline"`

	Approvals []Approval
}

func (approvalList *ApprovalList) GetApproval(appId string) *Approval {
	var approval *Approval
	for i := 0; i < len(approvalList.Approvals); i++ {
		if approvalList.Approvals[i].Key == appId {
			approval = &(approvalList.Approvals[i])
		}
	}
	return approval
}
