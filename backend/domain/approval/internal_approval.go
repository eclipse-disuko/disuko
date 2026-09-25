// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package approval

func (a *InternalApproval) FirstPendingApproverRole(username string) Approver {
	for i := range 4 {
		if a.Approver[i] == username && a.ApproveStates[i].State == Pending {
			return Approver(i)
		}
	}
	return None
}

func (a *InternalApproval) Finalized() bool {
	return a.Aborted || a.Declined() || (a.CustomerDone() && a.SupplierDone())
}

func (a *InternalApproval) ApprovedBySuppliersOnly() bool {
	return a.SupplierDone() && len(a.Approver) == 2
}

func (a *InternalApproval) Pending() bool {
	if a.Generating || a.GenerationFailed {
		return false
	}
	if a.Finalized() {
		return false
	}
	if a.ApprovedBySuppliersOnly() {
		return false
	}
	return true
}

func (a *InternalApproval) Declined() bool {
	for i := range 4 {
		if a.ApproveStates[i].State == Declined {
			return true
		}
	}
	return false
}

func (a *InternalApproval) IsApprover(username string) Approver {
	for i := range 4 {
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
