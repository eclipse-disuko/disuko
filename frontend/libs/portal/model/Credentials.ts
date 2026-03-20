// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export class ActionRights {
  public upload = false;
  public download = true;
  public delete = false;
}

export interface RefreshTokenResponseDto {
  accessToken: string;
  expiry: string;
}

export class RefreshTokenRequestDto {
  public refreshToken = '';
}
