// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package domain

import "time"

type BaseDto struct {
	Key     string    `json:"_key" validate:"lte=36"`
	Rev     string    `json:"_rev"`
	Updated time.Time `json:"updated"`
	Created time.Time `json:"created"`
}

func (dto *BaseDto) GetKey() string {
	return dto.Key
}
func (dto *BaseDto) SetKey(key string) {
	dto.Key = key
}

func (dto *BaseDto) GetRef() string {
	return dto.Rev
}
func (dto *BaseDto) SetRef(ref string) {
	dto.Rev = ref
}
func (dto *BaseDto) GetCreated() time.Time {
	return dto.Created
}
func (dto *BaseDto) SetCreated(created time.Time) {
	dto.Created = created
}

func (dto *BaseDto) GetUpdated() time.Time {
	return dto.Updated
}
func (dto *BaseDto) SetUpdated(updated time.Time) {
	dto.Updated = updated
}

type BaseResponseStatus struct {
	Success bool   `default:"true" json:"success"`
	Message string `json:"message"`
}
