// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {RuleStatus} from '@disclosure-portal/model/PolicyRuleClassification';

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
