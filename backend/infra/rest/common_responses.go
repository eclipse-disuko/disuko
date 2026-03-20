// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

type CountResponse struct {
	Count int `json:"count"`
}

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message" example:"Resource created"`
} //	@name	SuccessResponse

type FoundResponse struct {
	Found bool `json:"found"`
}
