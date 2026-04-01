// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package approval

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain/project/approvable"
	"github.com/eclipse-disuko/disuko/domain/project/components"
	"github.com/eclipse-disuko/disuko/domain/project/pdocument"
)

type DocumentFlagsDto struct {
	C1 bool `json:"c1"`
	C2 bool `json:"c2"`
	C3 bool `json:"c3"`
	C4 bool `json:"c4"`
	C5 bool `json:"c5"`
	C6 bool `json:"c6"`
}

type RequestInternalApprovalDto struct {
	MetaDoc           DocumentFlagsDto `json:"metaDoc"`
	CustomerApprover1 string           `json:"customerApprover1" validate:"RealInternalUser,NeFieldIfSet=CustomerApprover2"`
	CustomerApprover2 string           `json:"customerApprover2" validate:"RealInternalUser"`
	SupplierApprover1 string           `json:"supplierApprover1" validate:"required,RealInternalUser,NeFieldIfSet=SupplierApprover2"`
	SupplierApprover2 string           `json:"supplierApprover2" validate:"required,RealInternalUser"`
	Comment           string           `json:"comment" validate:"lte=1000"`
	GuidProject       string           `json:"guidProject" validate:"required,gte=3,lte=50"`
	WithZip           bool             `json:"withZip"`
	FOSSVersion       string           `json:"fossVersion"`
}

type RequestExternalApprovalDto struct {
	MetaDoc          DocumentFlagsDto `json:"metaDoc"`
	Comment          string           `json:"comment" validate:"lte=1000"`
	GuidProject      string           `json:"guidProject" validate:"required,gte=3,lte=50"`
	WithZip          bool             `json:"withZip"`
	FOSSVersion      string           `json:"fossVersion"`
	SelectedProjects []string         `json:"selectedProjects"`
}

type RequestPlausibilityCheckDto struct {
	MetaDoc     DocumentFlagsDto `json:"metaDoc"`
	Comment     string           `json:"comment" validate:"lte=1000"`
	GuidProject string           `json:"guidProject" validate:"required,gte=3,lte=50"`
	Approver    string           `json:"approver" validate:"required,RealUser"`
}

type UpdateApprovalDto struct {
	Comment string `json:"comment" validate:"lte=1000"`
	State   string `json:"state"`

	PowerOfAttorney PowerOfAttorneyType `json:"powerOfAttorney"`
}

type FillCustomerDto struct {
	CustomerApprover1 string `json:"customerApprover1" validate:"required,RealInternalUser,NeFieldIfSet=CustomerApprover2"`
	CustomerApprover2 string `json:"customerApprover2" validate:"required,RealInternalUser"`
}

func (p *ProjectApprovable) ToDto() ProjectApprovableDto {
	return ProjectApprovableDto{
		ProjectKey:      p.ProjectKey,
		ProjectName:     p.ProjectName,
		ApprovableSPDX:  p.ApprovableSPDX.ToDto(),
		CustomerDiffers: p.CustomerDiffers,
		SupplierDiffers: p.SupplierDiffers,
		ApprovableStats: p.ApprovableStats,
		SpdxName:        p.SpdxName,
		SpdxTag:         p.SpdxTag,
		SpdxUploaded:    p.SpdxUploaded,
		IsSpdxRecent:    p.IsSpdxRecent,
	}
}

type ProjectApprovableDto struct {
	ProjectKey      string                       `json:"projectKey"`
	ProjectName     string                       `json:"projectName"`
	ApprovableSPDX  approvable.ApprovableSPDXDto `json:"approvablespdx"`
	CustomerDiffers bool                         `json:"customerdiff"`
	SupplierDiffers bool                         `json:"supplierdiff"`
	ApprovableStats components.ComponentStats    `json:"stats"`
	SpdxName        string                       `json:"spdxname"`
	SpdxTag         string                       `json:"spdxtag"`
	SpdxUploaded    *time.Time                   `json:"spdxUploaded"`
	IsSpdxRecent    bool                         `json:"isSpdxRecent"`
}

type ApproveStateDto struct {
	State   string     `json:"state"`
	Updated *time.Time `json:"updated"`
}

func (s *ApproveState) ToDto() ApproveStateDto {
	return ApproveStateDto{
		State:   string(s.State),
		Updated: s.Updated,
	}
}

type ExternalApprovalDto struct {
	Vehicel        bool      `json:"vehicle"`
	State          StateInfo `json:"state"`
	ApproveComment string    `json:"comment"`
}

func (ea *ExternalApproval) ToDto() ExternalApprovalDto {
	return ExternalApprovalDto{
		Vehicel:        ea.Vehicle,
		State:          ea.State,
		ApproveComment: ea.Comment,
	}
}

type InternalApprovalDto struct {
	ApproveStates     [4]ApproveStateDto `json:"states"`
	Approver          [4]string          `json:"approver"`
	ApproverFullNames [4]string          `json:"approverFullName"`
	Comments          [4]string          `json:"comments"`
	DocVersion        int                `json:"docVersion"`
	Aborted           bool               `json:"aborted"`
}

func (ia *InternalApproval) ToDto(approverFullNames [4]string) InternalApprovalDto {
	var s [4]ApproveStateDto
	for i := 0; i < 4; i++ {
		s[i] = ia.ApproveStates[i].ToDto()
	}

	return InternalApprovalDto{
		ApproveStates:     s,
		Approver:          ia.Approver,
		ApproverFullNames: approverFullNames,
		Comments:          ia.ApproveComments,
		DocVersion:        ia.DocVersion,
		Aborted:           ia.Aborted,
	}
}

type PlausibilityDto struct {
	State            ApproveStateDto `json:"state"`
	ApproveComment   string          `json:"comment"`
	Approver         string          `json:"approver"`
	ApproverFullName string          `json:"approverFullName"`
}

func (p *PlausibilityCheck) ToDto(approverFullName string) PlausibilityDto {
	return PlausibilityDto{
		State:            p.State.ToDto(),
		ApproveComment:   p.ApproveComment,
		Approver:         p.Approver,
		ApproverFullName: approverFullName,
	}
}

type ApprovableInfoDto struct {
	CompStats          *components.ComponentStats `json:"stats"`
	Projects           []ProjectApprovableDto     `json:"projects"`
	HasDeniedDecisions bool                       `json:"hasDeniedDecisions"`
}

func (i *Info) ToDto() (res ApprovableInfoDto) {
	res.CompStats = i.CompStats
	for _, p := range i.Projects {
		res.Projects = append(res.Projects, p.ToDto())
	}
	return
}

type ApprovalDto struct {
	Key     string    `json:"key"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`

	ProjectGuid     string                   `json:"projectKey"`
	Creator         string                   `json:"creator"`
	CreatorFullName string                   `json:"creatorFullName"`
	Comment         string                   `json:"comment"`
	Info            ApprovableInfoDto        `json:"info"`
	Documents       []pdocument.PDocumentDto `json:"documents"`
	DocumentFlags   DocumentFlagsDto         `json:"flags"`

	Type         ApprovalType        `json:"type"`
	Internal     InternalApprovalDto `json:"internal"`
	Plausibility PlausibilityDto     `json:"plausibility"`
	External     ExternalApprovalDto `json:"external"`

	Status StateInfo `json:"status"`
}

func (a *Approval) ToDto(creatorFullName string, approverFullNames [4]string) (res ApprovalDto) {
	res.Key = a.Key
	res.Created = a.Created
	res.Updated = a.Updated
	res.ProjectGuid = a.ProjectGuid
	res.Creator = a.Creator
	res.CreatorFullName = creatorFullName
	res.Comment = a.Comment
	res.DocumentFlags = a.DocumentFlags.ToDto()
	res.Type = a.Type
	if res.Type == TypePlausibility {
		res.Plausibility = a.Plausibility.ToDto(approverFullNames[0])
	} else if res.Type == TypeInternal {
		res.Internal = a.Internal.ToDto(approverFullNames)
	} else {
		res.External = a.External.ToDto()
	}
	res.Info = a.Info.ToDto()

	res.Status = a.ToApprovalDtoStatus()
	res.Documents = make([]pdocument.PDocumentDto, 0)
	return
}

func (a *Approval) ToApprovalDtoStatus() StateInfo {
	if a.Type == TypePlausibility {
		if a.Plausibility.Aborted {
			return Aborted
		}
		if a.Plausibility.State.State == Unset {
			return Pending
		}
		return a.Plausibility.State.State
	} else if a.Type == TypeInternal {
		if a.Internal.Aborted {
			return Aborted
		}
		if a.Internal.Generating {
			return GeneratingDocs
		}
		if a.Internal.IsDeclined() {
			return Declined
		}
		if a.Internal.CustomerDone() {
			return Approved
		}
		if a.Internal.SupplierDone() {
			return SupplierApproved
		}
		return Pending

	}
	return Pending
}
