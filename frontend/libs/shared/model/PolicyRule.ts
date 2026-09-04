// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export interface License {
  key: string;
  identifier: string;
  name: string;
  policyType?: string;
}

export interface PolicyRule {
  id?: string;
  key: string;
  name: string;
  type: string;
  description: string;
  created: string;
  updated: string;
  licenses: License[];
}

export type PolicyCat = 'allow' | 'warn' | 'deny';

export interface PolicyStatusDef {
  key: 'all' | PolicyCat;
  icon: string;
  outlineIcon?: string;
  color: string;
  labelKey: string;
}

export const POLICY_STATUS_DEFS: PolicyStatusDef[] = [
  {key: 'all', icon: 'mdi-format-list-bulleted', color: 'primary', labelKey: 'FILTER_ALL'},
  {
    key: 'allow',
    icon: 'mdi-check-circle',
    outlineIcon: 'mdi-check-circle-outline',
    color: 'success',
    labelKey: 'FILTER_ALLOWED',
  },
  {key: 'warn', icon: 'mdi-alert', outlineIcon: 'mdi-alert-outline', color: 'warning', labelKey: 'FILTER_WARNED'},
  {
    key: 'deny',
    icon: 'mdi-close-circle',
    outlineIcon: 'mdi-close-circle-outline',
    color: 'error',
    labelKey: 'FILTER_DENIED',
  },
];

export enum PolicyState {
  ALLOW = 'allow',
  DENY = 'deny',
  WARN = 'warn',
  NOT_SET = 'NOT_SET',
  NOASSERTION = 'noassertion',
  QUESTIONED = 'questioned',
}

export const PolicyRules: PolicyState[] = [PolicyState.ALLOW, PolicyState.WARN, PolicyState.DENY];
export const PolicyStates: PolicyState[] = [
  PolicyState.NOT_SET,
  PolicyState.DENY,
  PolicyState.NOASSERTION,
  PolicyState.WARN,
  PolicyState.QUESTIONED,
  PolicyState.ALLOW,
];

export class PolicyRulesAssignmentsDto {
  public status = '';
  public key = '';
  public name = '';
  public description = '';
  public type = '';
}

export function comparePolicyStatus(a: string, b: string): number {
  const statusWeight: Record<string, number> = {
    deny: 0,
    noassertion: 1,
    questioned: 2,
    warn: 3,
    notset: 4,
    allow: 5,
  };
  const aWeight = statusWeight[a?.toLowerCase()] ?? statusWeight.notset;
  const bWeight = statusWeight[b?.toLowerCase()] ?? statusWeight.notset;
  return aWeight - bWeight;
}
