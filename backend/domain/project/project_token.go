// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"time"

	"mercedes-benz.ghe.com/foss/disuko/domain"
)

type Token struct {
	domain.ChildEntity `bson:",inline"`
	Company            string `validate:"required,gte=1,lte=100"`
	Description        string `validate:"lte=1000"`
	Expiry             string `validate:"lte=100"`
	TokenSecret        string `validate:"lte=36"`
	Status             TokenStatus
}

type TokenStatus string

const (
	REVOKED = "revoked"
	EXPIRED = "expired"
	ACTIVE  = "active"
)

func (token *Token) IsExpired() bool {
	if token.Expiry == "" {
		return false
	}
	expiry, err := token.GetExpired()
	if err != nil {
		return true
	}
	return token.Status == EXPIRED || expiry.Before(time.Now())
}

func (token *Token) GetExpired() (time.Time, error) {
	return time.Parse(time.RFC3339, token.Expiry)
}
