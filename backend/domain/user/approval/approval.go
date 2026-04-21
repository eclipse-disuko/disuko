// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package approval

type ApprovalStatus string

const (
	ApRequested ApprovalStatus = "REQUESTED"
	ApApproved  ApprovalStatus = "APPROVED"
	ApDeclined  ApprovalStatus = "DECLINED"
)

type RequestApprovalDto struct {
	RequestCreateDisclosureDocumentDto
	UserRequested string           `json:"requestUser" validate:"required,gte=3,lte=50"`
	Comment       string           `json:"comment" validate:"lte=255"`
	Stats         ApprovalStatsDto `json:"stats"`
	GuidVersion   string           `json:"guidVersion" validate:"required,gte=3,lte=50"`
	GuidProject   string           `json:"guidProject" validate:"required,gte=3,lte=50"`
}

type RequestCreateDisclosureDocumentDto struct {
	GuidSBOM string          `json:"guidSBOM" validate:"required,gte=3,lte=50"`
	DocMeta  MetaDocumentDto `json:"metaDoc"`
}

type MetaDocumentDto struct {
	// client/dportal/src/components/dialog/project/RequestApproval.vue
	C1 bool `json:"c1"`
	C2 bool `json:"c2"`
	C3 bool `json:"c3"`
	C4 bool `json:"c4"`
	C5 bool `json:"c5"`
	C6 bool `json:"c6"`
}

type ApprovalStatsDto struct {
	Allowed     int `json:"allowed"`
	Denied      int `json:"denied"`
	NoAssertion int `json:"noAssertion"`
	Questioned  int `json:"questioned"`
	Total       int `json:"total"`
	Warned      int `json:"warned"`
}

type UpdateApprovalDto struct {
	Comment        string `json:"comment" validate:"lte=255"`
	FromUserGuid   string `json:"guidFromUser" validate:"lte=36"`
	TargetUserGuid string `json:"guidTargetUser" validate:"lte=36"`
	ProjectGuid    string `json:"guidProject" validate:"lte=36"`
	VersionGuid    string `json:"guidVersion" validate:"lte=36"`
	SBOMGuid       string `json:"guidSBOM" validate:"lte=36"`
	TaskGuid       string `json:"guidTask" validate:"lte=36"`
	IsAccepted     bool   `json:"accepted"`
}

type ResponseApprovalDto struct {
	ApprovalGuid string `json:"approvalGuid"`
	Success      bool   `json:"success"`
	JobKey       string `json:"jobKey"`
}
