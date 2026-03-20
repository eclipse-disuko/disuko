// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package approval

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/audit"
	"github.com/eclipse-disuko/disuko/domain/project/approvable"
	"github.com/eclipse-disuko/disuko/domain/project/components"
)

type ApproveState struct {
	State   StateInfo
	Updated *time.Time
}

type StateInfo string

const (
	Unset            StateInfo = ""
	Pending          StateInfo = "PENDING"
	Declined         StateInfo = "DECLINED"
	Approved         StateInfo = "APPROVED"
	SupplierApproved StateInfo = "SUPPLIER_APPROVED"
	CustomerApproved StateInfo = "CUSTOMER_APPROVED"
	Aborted          StateInfo = "ABORTED"
	GeneratingDocs   StateInfo = "GENERATING"
)

func ParseStateInfo(state string) (valid bool, result StateInfo) {
	switch state {
	case string(Pending):
		valid, result = true, Pending
	case string(Declined):
		valid, result = true, Declined
	case string(Approved):
		valid, result = true, Approved
	case string(SupplierApproved):
		valid, result = true, SupplierApproved
	case string(CustomerApproved):
		valid, result = true, CustomerApproved
	case string(Aborted):
		valid, result = true, Aborted
	case string(GeneratingDocs):
		valid, result = true, GeneratingDocs
	default:
		valid, result = false, Unset
	}
	return
}

type Approver int

const (
	Supplier1 Approver = 0
	Supplier2 Approver = 1
	Customer1 Approver = 2
	Customer2 Approver = 3
	None      Approver = 4
)

type ApprovalType string

const (
	TypeInternal     ApprovalType = "INTERNAL"
	TypePlausibility ApprovalType = "PLAUSIBILITY"
	TypeExternal     ApprovalType = "EXTERNAL"
)

type TaskMetaDocument struct {
	C1 bool
	C2 bool
	C3 bool
	C4 bool
	C5 bool
	C6 bool
}

func (d TaskMetaDocument) ToDto() DocumentFlagsDto {
	return DocumentFlagsDto{
		C1: d.C1,
		C2: d.C2,
		C3: d.C3,
		C4: d.C4,
		C5: d.C5,
		C6: d.C6,
	}
}

type ProjectApprovable struct {
	ProjectKey     string
	ProjectName    string
	ApprovableSPDX approvable.ApprovableSPDX
	// TODO: We don't need them anymore?
	CustomerDiffers bool
	SupplierDiffers bool
	ApprovableStats components.ComponentStats
	SpdxName        string
	SpdxTag         string
	SpdxUploaded    *time.Time
	IsSpdxRecent    bool
}

type Info struct {
	CompStats *components.ComponentStats
	Projects  []ProjectApprovable
}

type InternalApproval struct {
	ApproveStates   [4]ApproveState
	Approver        [4]string
	ApproveComments [4]string
	DocVersion      int
	Aborted         bool
	Generating      bool
}

type PlausibilityCheck struct {
	State          ApproveState
	ApproveComment string
	Approver       string
	Aborted        bool
}

type ExternalApproval struct {
	Vehicle bool
	Comment string
	State   StateInfo
}

type Approval struct {
	domain.ChildEntity `bson:"inline"`
	audit.Container    `bson:"inline"`

	ProjectGuid   string
	Creator       string
	Comment       string
	Info          Info
	DocumentFlags TaskMetaDocument

	Type         ApprovalType
	Internal     InternalApproval
	Plausibility PlausibilityCheck
	External     ExternalApproval
}

func (p *PlausibilityCheck) GetPendingApprover() string {
	if p.State.State != Pending {
		return ""
	}
	return p.Approver
}

func (p *PlausibilityCheck) IsActive() bool {
	return p.State.State == Pending && !p.Aborted
}

type PowerOfAttorneyType string

const (
	PowerOfAttorneyIV    PowerOfAttorneyType = "iV"
	PowerOfAttorneyPPA   PowerOfAttorneyType = "ppa"
	PowerOfAttorneyOther PowerOfAttorneyType = "other"
)
