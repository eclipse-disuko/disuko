// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import PolicyRule from '@disclosure-portal/model/PolicyRule';

export type RuleStatus = 'allowed' | 'warned' | 'denied' | 'forbidden';

export const statusConfig: Record<RuleStatus, {color: string; icon: string; labelKey: string}> = {
  allowed: {color: 'green', icon: 'mdi-check-circle', labelKey: 'ALLOWED'},
  warned: {color: 'policyStatusWarnedColor', icon: 'mdi-alert', labelKey: 'WARNED'},
  denied: {color: 'policyStatusDeniedColor', icon: 'mdi-minus-circle', labelKey: 'DENIED'},
  forbidden: {color: 'error', icon: 'mdi-cancel', labelKey: 'MATRIX_STATUS_FORBIDDEN'},
};

export const getStatusColor = (status: RuleStatus | undefined): string =>
  status ? (statusConfig[status]?.color ?? '') : '';

export const getStatusIcon = (status: RuleStatus | undefined): string =>
  status ? (statusConfig[status]?.icon ?? '') : '';

export const getStatusLabel = (status: RuleStatus | undefined, t: (key: string) => string): string =>
  status ? t(statusConfig[status]?.labelKey) : '';

export const toRuleStatusMap = (rule: PolicyRule): Record<string, RuleStatus> => {
  const bucket = rule.calculatedConfig?.bucketDefinition;
  if (!bucket) return {};
  return Object.fromEntries([
    ...(bucket.allowedClassifications ?? []).map((key): [string, RuleStatus] => [key, 'allowed']),
    ...(bucket.warnedClassifications ?? []).map((key): [string, RuleStatus] => [key, 'warned']),
    ...(bucket.deniedClassifications ?? []).map((key): [string, RuleStatus] => [key, 'denied']),
  ]);
};

export const toBucketDefinition = (ruleMap: Record<string, RuleStatus>) => {
  const entries = Object.entries(ruleMap);
  return {
    allowedClassifications: entries.filter(([, s]) => s === 'allowed').map(([k]) => k),
    warnedClassifications: entries.filter(([, s]) => s === 'warned').map(([k]) => k),
    deniedClassifications: entries.filter(([, s]) => s === 'denied').map(([k]) => k),
  };
};
