// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package approval

import "mercedes-benz.ghe.com/foss/disuko/domain/project/approvable"

type ChildApprovable struct {
	Key             string
	Name            string
	ApprovableSPDX  approvable.ApprovableSPDX
	CustomerDiffers bool
	SupplierDiffers bool
	ApprovableStats ComponentStats
	SpdxName        string
}

type GroupApprovableInfo struct {
	GroupCompStats ComponentStats
	Children       []ChildApprovable
}

type ChildApprovableDto struct {
	Key             string                       `json:"key"`
	Name            string                       `json:"name"`
	ApprovableSPDX  approvable.ApprovableSPDXDto `json:"approvablespdx"`
	CustomerDiffers bool                         `json:"customerdiff"`
	SupplierDiffers bool                         `json:"supplierdiff"`
	ApprovableStats ComponentStats               `json:"stats"`
	SpdxName        string                       `json:"spdxname"`
}

func (c *ChildApprovableDto) ToEntity() ChildApprovable {
	return ChildApprovable{
		Key:             c.Key,
		Name:            c.Name,
		ApprovableSPDX:  c.ApprovableSPDX.ToEntity(),
		CustomerDiffers: c.CustomerDiffers,
		SupplierDiffers: c.SupplierDiffers,
		ApprovableStats: c.ApprovableStats,
		SpdxName:        c.SpdxName,
	}
}

func (c *ChildApprovable) ToDto() ChildApprovableDto {
	return ChildApprovableDto{
		Key:             c.Key,
		Name:            c.Name,
		ApprovableSPDX:  c.ApprovableSPDX.ToDto(),
		CustomerDiffers: c.CustomerDiffers,
		SupplierDiffers: c.SupplierDiffers,
		ApprovableStats: c.ApprovableStats,
		SpdxName:        c.SpdxName,
	}
}

type GroupApprovableDto struct {
	GroupCompStats ComponentStats       `json:"stats"`
	Children       []ChildApprovableDto `json:"children"`
}

func (g *GroupApprovableDto) ToEntity() (res GroupApprovableInfo) {
	res.GroupCompStats = g.GroupCompStats
	for _, c := range g.Children {
		res.Children = append(res.Children, c.ToEntity())
	}
	return
}

func (g *GroupApprovableInfo) ToDto() (res GroupApprovableDto) {
	res.GroupCompStats = g.GroupCompStats
	for _, c := range g.Children {
		res.Children = append(res.Children, c.ToDto())
	}
	return
}

type GroupApprovalRequestDto struct {
	DocMeta   MetaDocumentDto    `json:"metaDoc"`
	SRI       string             `json:"sri" validate:"required,gte=3,lte=50"`
	FRI       string             `json:"fri" validate:"required,gte=3,lte=50"`
	Comment   string             `json:"comment" validate:"lte=255"`
	GroupInfo GroupApprovableDto `json:"groupinfo"`
}
