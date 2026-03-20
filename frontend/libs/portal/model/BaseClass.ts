// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export class BaseDto {
  public _key = '';
  public _rev = '';
  public updated = '';
  public created = '';

  public constructor(dto: BaseDto | null | undefined = null) {
    if (dto) {
      Object.assign(this, dto);
    }
  }
}
