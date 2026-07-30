// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export interface PersonalToken {
  key: string;
  description: string;
  expiry: string;
  expired: boolean;
  created: string;
}

export interface CreatePersonalTokenRequest {
  description: string;
  expiry: string;
}

export interface CreatePersonalTokenResponse {
  token: string;
}
