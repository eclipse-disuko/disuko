// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package publicauth

import "github.com/golang-jwt/jwt/v4"

type AccessClaims struct {
	jwt.RegisteredClaims
	ProjectUUID string `json:"projectUUID"`
	TokenKey    string `json:"tokenKey"`
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	ProjectUUID string `json:"projectUUID"`
	TokenKey    string `json:"tokenKey"`
}

type LoginReq struct {
	ProjectUUID string `json:"projectUUID" validate:"required"`
	Token       string `json:"token" validate:"required"`
}

type InfoDto struct {
	ProjectUUID string `json:"projectUUID"`
	IsGroup     bool   `json:"isGroup"`
}
