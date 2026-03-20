// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package approval

func (a *InternalApproval) FirstPendingApproverRole(username string) Approver {
	for i := 0; i < 4; i++ {
		if a.Approver[i] == username && a.ApproveStates[i].State == Pending {
			return Approver(i)
		}
	}
	return None
}

func (a *InternalApproval) IsActive() bool {
	if a.Aborted {
		return false
	}

	acceptedCount := 0
	for i := 0; i < 4; i++ {
		if a.ApproveStates[i].State == Declined {
			return false
		}
		if a.ApproveStates[i].State == Approved {
			acceptedCount++
		}
	}
	return (acceptedCount != 4)
}

func (a *InternalApproval) IsDeclined() bool {
	for i := 0; i < 4; i++ {
		if a.ApproveStates[i].State == Declined {
			return true
		}
	}
	return false
}

func (a *InternalApproval) IsApprover(username string) Approver {
	for i := 0; i < 4; i++ {
		if a.Approver[i] == username {
			return Approver(i)
		}
	}
	return None
}

func (a *InternalApproval) CustomerDone() bool {
	return a.ApproveStates[Customer1].State == Approved && a.ApproveStates[Customer2].State == Approved
}

func (a *InternalApproval) SupplierDone() bool {
	return a.ApproveStates[Supplier1].State == Approved && a.ApproveStates[Supplier2].State == Approved
}

func (a *InternalApproval) GetApproverName(ai Approver) string {
	return a.Approver[ai]
}
