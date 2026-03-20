// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package obligation

import (
	"time"

	"github.com/google/uuid"
	"mercedes-benz.ghe.com/foss/disuko/domain"
)

type Obligation struct {
	domain.RootEntity `bson:"inline"`
	domain.SoftDelete `bson:"inline"`
	Name              string    `json:"name"`
	NameDe            string    `json:"nameDe"`
	Type              Type      `json:"type"`
	WarnLevel         WarnLevel `json:"warnLevel"`
	Description       string    `json:"description"`
	DescriptionDe     string    `json:"descriptionDe"`
	AutoApproved      bool      `json:"autoApproved"`
}

type Type string

const (
	Rights      = "rights"
	Obligations = "obligations"
	Liability   = "liability"
	Other       = "other"
)

type WarnLevel string

const (
	Remark      = "remark"
	Information = "information"
	Warning     = "warning"
	Alarm       = "alarm"
)

var levelWeight = map[WarnLevel]int{
	Information: 0,
	Warning:     1,
	Alarm:       2,
}

func GetLevelWeight(level WarnLevel) int64 {
	return int64(levelWeight[level])
}

func (entity *Obligation) Update(obligation *ObligationDto) {
	entity.Name = obligation.Name
	entity.NameDe = obligation.NameDe
	entity.WarnLevel = obligation.WarnLevel
	entity.Type = obligation.Type
	entity.Description = obligation.Description
	entity.DescriptionDe = obligation.DescriptionDe
	entity.AutoApproved = obligation.AutoApproved
	entity.Updated = time.Now()
}

func CreateFrom(obligationDto *ObligationDto) *Obligation {
	newItem := obligationDto.ToEntity()
	newItem.Key = uuid.New().String()
	newItem.Created = time.Now()
	newItem.Updated = time.Now()
	return newItem
}
