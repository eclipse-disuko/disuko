// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {BaseDto} from '@disclosure-portal/model/BaseClass';

export const labelTypes = {
  SCHEMA: 'SCHEMA',
  POLICY: 'POLICY',
  PROJECT: 'PROJECT',
} as const;

export type LabelType = keyof typeof labelTypes;

export interface Label extends BaseDto {
  type: LabelType;
  name: string;
  description: string;
}

export default Label;
