// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package approvable

const (
	APPROVAL_TYPE_INTERNAL = "internal"
	APPROVAL_TYPE_EXTERNAL = "external"
	APPROVAL_TYPE_VEHICLE  = "vehicle"
	APPROVAL_TYPE_PLAUSI   = "plausi"
)

type ApprovableSPDX struct {
	VersionName string
	VersionKey  string
	SpdxKey     string
}

type ApprovableSPDXDto struct {
	VersionName string `json:"versionName"`
	VersionKey  string `json:"versionkey" validate:"lte=36" example:"dummy-c490-4cbe-8373-1396ba6e5b06"`
	SpdxKey     string `json:"spdxkey" validate:"lte=36" example:"dummy-c490-4cbe-8373-1396ba6e5b06"`
}

func (entity *ApprovableSPDX) ToDto() ApprovableSPDXDto {
	return ApprovableSPDXDto{
		VersionName: entity.VersionName,
		VersionKey:  entity.VersionKey,
		SpdxKey:     entity.SpdxKey,
	}
}

func (a *ApprovableSPDXDto) ToEntity() ApprovableSPDX {
	return ApprovableSPDX{
		VersionName: a.VersionName,
		VersionKey:  a.VersionKey,
		SpdxKey:     a.SpdxKey,
	}
}
