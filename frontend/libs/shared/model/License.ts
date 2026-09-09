// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import type {IObligation} from '@shared/model/IObligation';

export class ClassificationWithCount {
  public classification: IObligation = {} as IObligation;
  public count = 0;
}

export enum LicenseFamily {
  PERMISSIVE = 'permissive',
  WEAKCOPYLEFT = 'weak copyleft',
  STRONGCOPYLEFT = 'strong copyleft',
  NETWORKCOPYLEFT = 'network copyleft',
  NOTDECLARED = 'not declared',
}

export const familyWeight: Map<string, number> = new Map<string, number>([
  ['', -1],
  ['unknown', 0],
  [LicenseFamily.NOTDECLARED, 1],
  [LicenseFamily.NETWORKCOPYLEFT, 2],
  [LicenseFamily.STRONGCOPYLEFT, 3],
  [LicenseFamily.WEAKCOPYLEFT, 4],
  [LicenseFamily.PERMISSIVE, 5],
]);

export function compareFamily(aRaw: string, bRaw: string): number {
  const a = aRaw.toLowerCase().replace('_', ' ');
  const b = bRaw.toLowerCase().replace('_', ' ');
  return (familyWeight.get(a) ?? 0) - (familyWeight.get(b) ?? 0);
}
