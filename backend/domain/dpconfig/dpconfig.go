// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package dpconfig

import (
	"github.com/eclipse-disuko/disuko/domain"
)

type DPConfig struct {
	domain.RootEntity `bson:"inline"`
	DatabaseVersion   int `default:"0"`
}

type SampleDataCreationState struct {
	domain.RootEntity `bson:"inline"`
	domain.BaseState  `bson:"inline"`
	HasErrors         bool   `default:"false" json:"hasErrors"`
	LastError         string `json:"lastError"`
	TargetCount       int    `json:"targetCount"`
	CreatedCount      int    `json:"createdCount"`
	WithFileUpload    bool   `json:"withFileUpload"`
}
