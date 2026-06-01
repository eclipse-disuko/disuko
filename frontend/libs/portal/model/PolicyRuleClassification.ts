// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {BaseDto} from '@disclosure-portal/model/BaseClass';

export type RuleStatus = 'allowed' | 'warned' | 'denied' | 'forbidden';

export const DEFAULT_CLASSIFICATION_NAMES: ReadonlyArray<string> = [
  'Permissive',
  'Copyleft (weak)',
  'Copyleft (strong)',
  'Copyleft (network protective)',
  'Anti-Tivoization',
  'Severe patent retaliation',
  'Doing Business with US',
  'Unclear or Ambiguous',
];

export interface ClassificationInfo {
  key: string;
  name: string;
}

export interface MatrixResponseDto {
  classifications: ClassificationInfo[];
  useCases: PolicyRuleClassificationDto[];
}

export interface PolicyRuleClassificationRequestDto {
  name: string;
  rules: Record<string, RuleStatus>;
}

export class PolicyRuleClassificationDto extends BaseDto {
  public name = '';
  public rules: Record<string, RuleStatus> = {};

  public constructor(dto?: PolicyRuleClassificationDto | null) {
    super(dto);
  }
}
