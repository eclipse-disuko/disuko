// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {BaseDto} from '@disclosure-portal/model/BaseClass';
import {ReviewRemarkLevel} from '@disclosure-portal/model/Quality';

export class ReviewTemplate extends BaseDto {
  public title = '';
  public description = '';
  public level: ReviewRemarkLevel = ReviewRemarkLevel.NOT_SET;
  public source = '';

  public constructor(dto: BaseDto | null | undefined = null) {
    super(dto);
    if (dto) {
      Object.assign(this, dto);
    }
  }
}
