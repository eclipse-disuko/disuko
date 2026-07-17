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

// Canonical policy categories applied to licenses
export type PolicyCat = 'allow' | 'warn' | 'deny';

// Definition for policy status filter buttons (UI metadata)
export interface PolicyStatusDef {
  key: 'all' | PolicyCat;
  icon: string;
  outlineIcon?: string;
  color: string; // Vuetify theme color name
  labelKey: string; // i18n translation key
}

// Shared status button definitions for reuse in dialogs/components
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

// CLI-side PolicyState (duplicated to avoid portal model dependency)
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

/**
 * Compare function for policy status to enable proper sorting
 * Order: DENY (0) < NOASSERTION (1) < QUESTIONED (2) < WARN (3) < NOT_SET (4)  < ALLOW (5)
 */
export function comparePolicyStatus(a: string, b: string): number {
  const statusWeight: Record<string, number> = {
    deny: 0,
    noassertion: 1,
    questioned: 2,
    warn: 3,
    notSet: 4,
    allow: 5,
  };
  const aWeight = statusWeight[a?.toLowerCase()] ?? statusWeight['notSet'];
  const bWeight = statusWeight[b?.toLowerCase()] ?? statusWeight['notSet'];
  return aWeight - bWeight;
}
