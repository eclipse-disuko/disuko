// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package approval

type ApprovalAudit struct {
	ProjectGuid   string           `json:"projectGuid"`
	Creator       string           `json:"creator"`
	Comment       string           `json:"comment"`
	Info          Info             `json:"info"`
	DocumentFlags TaskMetaDocument `json:"flags"`

	Type         ApprovalType      `json:"type"`
	Internal     InternalApproval  `json:"internal"`
	Plausibility PlausibilityCheck `json:"plausibility"`
	External     ExternalApproval  `json:"external"`
}

// TODO: refactor approval audit logs
func (a *Approval) ToAudit() ApprovalAudit {
	return ApprovalAudit{
		ProjectGuid:   a.ProjectGuid,
		Creator:       a.Creator,
		Comment:       a.Comment,
		Info:          a.Info,
		DocumentFlags: a.DocumentFlags,

		Type:         a.Type,
		Internal:     a.Internal,
		Plausibility: a.Plausibility,
		External:     a.External,

		//CustomerApprover1: a.CustomerApprover1,
		//CustomerApprover2: a.CustomerApprover2,
		//SupplierApprover1: a.SupplierApprover1,
		//SupplierApprover2: a.SupplierApprover2,
		//ApproveComments:   a.ApproveComments,
	}
}
