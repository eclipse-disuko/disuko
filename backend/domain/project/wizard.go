// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package project

import "errors"

type TargetPlatform string

const (
	PlatformEnterprise TargetPlatform = "Enterprise IT"
	PlatformMobile     TargetPlatform = "Mobile"
	PlatformVehicle    TargetPlatform = "Product"
	PlatformOther      TargetPlatform = "Other"
)

func (t TargetPlatform) Validate() error {
	switch t {
	case PlatformEnterprise, PlatformMobile, PlatformVehicle, PlatformOther:
		return nil
	}
	return errors.New("unknown target platform")
}

type Architecture string

const (
	ArchitectureFrontend        Architecture = "Frontend or Client"
	ArchitectureBackend         Architecture = "Backend"
	ArchitectureVehicleOnboard  Architecture = "Product onboard"
	ArchitectureVehicleOffboard Architecture = "Product offboard"
	ArchitectureNone            Architecture = "None"
	ArchitectureEmpty           Architecture = ""
)

func (a Architecture) Validate() error {
	switch a {
	case ArchitectureBackend, ArchitectureFrontend, ArchitectureNone, ArchitectureVehicleOffboard, ArchitectureVehicleOnboard, ArchitectureEmpty:
		return nil
	}
	return errors.New("unknown architecture")
}

type TargetUsers string

const (
	TargetUsersCompany         TargetUsers = "Company"
	TargetUsersBusinessPartner TargetUsers = "Business Partner"
	TargetUsersCustomer        TargetUsers = "End Customer"
	TargetUsersEmpty           TargetUsers = ""
)

func (t TargetUsers) Validate() error {
	switch t {
	case TargetUsersCompany, TargetUsersBusinessPartner, TargetUsersCustomer, TargetUsersEmpty:
		return nil
	}
	return errors.New("unknown target users")
}

type DistributionTarget string

const (
	DistributionTargetsCompany         DistributionTarget = "Company"
	DistributionTargetsBusinessPartner DistributionTarget = "Business Partner"
	DistributionTargetsEmpty           DistributionTarget = ""
)

func (d DistributionTarget) Validate() error {
	switch d {
	case DistributionTargetsCompany, DistributionTargetsBusinessPartner, DistributionTargetsEmpty:
		return nil
	}
	return errors.New("unknown distribution target")
}

type Development string

const (
	DevelopmentsInhouse  Development = "In-house"
	DevelopmentsInternal Development = "Internal Developer"
	DevelopmentsExternal Development = "External Developer"
)

func (d Development) Validate() bool {
	switch d {
	case DevelopmentsExternal, DevelopmentsInhouse, DevelopmentsInternal:
		return true
	default:
		return false
	}
}

type WizardGroupDto struct {
	Name        string `json:"name" validate:"required,gte=3,lte=80"`
	Description string `json:"description" validate:"lte=10000"`

	Architecture       Architecture       `json:"architecture" validate:"validateFn"`
	DistributionTarget DistributionTarget `json:"distributionTarget" validate:"validateFn"`
	Development        Development        `json:"development" validate:"validateFn"`

	IsDummy   *bool  `json:"isDummy"`
	IsGroup   *bool  `json:"isGroup" validate:"required"`
	ParentKey string `json:"parentKey"`

	Settings ProjectSettingsDto `json:"projectSettings" validate:"required"`

	Labels []string `json:"labels"`
}
type WizardProjectDto struct {
	Name            string             `json:"name" validate:"required,gte=3,lte=80"`
	Description     string             `json:"description" validate:"lte=10000"`
	ApplicationMeta ApplicationMetaDto `json:"applicationMeta"`

	TargetPlatform     TargetPlatform     `json:"targetPlatform" validate:"validateFn"`
	Architecture       Architecture       `json:"architecture" validate:"validateFn"`
	DistributionTarget DistributionTarget `json:"distributionTarget" validate:"validateFn"`
	TargetUsers        TargetUsers        `json:"targetUsers" validate:"validateFn"`
	Development        Development        `json:"development" validate:"validateFn"`

	IsDummy   *bool  `json:"isDummy"`
	IsGroup   *bool  `json:"isGroup" validate:"required"`
	ParentKey string `json:"parentKey"`

	Settings ProjectSettingsDto `json:"projectSettings" validate:"required"`

	Labels []string `json:"labels"`
}

func (d WizardProjectDto) GetIsGroup() bool {
	return d.IsGroup != nil && *d.IsGroup
}

func (d WizardProjectDto) GetIsDummy() bool {
	return d.IsDummy != nil && *d.IsDummy
}

func (d WizardGroupDto) GetIsDummy() bool {
	return d.IsDummy != nil && *d.IsDummy
}

type WizardAttributesDto struct {
	TargetPlatform     TargetPlatform     `json:"targetPlatform" validate:"validateFn"`
	Architecture       Architecture       `json:"architecture" validate:"validateFn"`
	DistributionTarget DistributionTarget `json:"distributionTarget" validate:"validateFn"`
	TargetUsers        TargetUsers        `json:"targetUsers" validate:"validateFn"`
	Development        Development        `json:"development" validate:"validateFn"`
	HasDeniedDecisions bool               `json:"hasDeniedDecisions"`
}
