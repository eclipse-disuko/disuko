// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package checklist

import (
	"mercedes-benz.ghe.com/foss/disuko/domain"
	"mercedes-benz.ghe.com/foss/disuko/domain/project"
)

type TriggerType string

const (
	Default           TriggerType = "DEFAULT"
	ClassificationAND TriggerType = "CLASS_AND"
	ClassificationOR  TriggerType = "CLASS_OR"
	PolicyStatus      TriggerType = "POLICY_STATUS"
	ScanRemark        TriggerType = "SCAN_REMARK"
	License           TriggerType = "LICENSE"
	ComponentName     TriggerType = "COMPONENT_NAME"
)

type PolicyStatusType string

const (
	Allowed    PolicyStatusType = "ALLOWED"
	Denied     PolicyStatusType = "DENIED"
	Unasserted PolicyStatusType = "UNASSERTED"
	Warned     PolicyStatusType = "WARNED"
	Questioned PolicyStatusType = "QUESTIONED"
)

type Item struct {
	domain.ChildEntity `bson:"inline"`
	Name               string
	TriggerType        TriggerType
	Classifications    []string
	PolicyStatus       []PolicyStatusType
	PolicyLabels       []string
	ScanRemarks        *project.ScanRemarkStatus
	LicenseIds         []string
	ComponentNames     []string
	TargetTemplateName string
	TargetTemplateKey  string
}

type Checklist struct {
	domain.RootEntity `bson:"inline"`
	Name              string
	NameDE            string
	Description       string
	DescriptionDE     string
	PolicyLabels      []string
	Items             []Item
	Active            bool
}
