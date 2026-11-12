// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import SharedIcons from '@shared/constants/icons';

enum Icons {
  POLICY = 'mdi-shield-search',
  LABEL = 'mdi-label',
  PROJECT = 'mdi-link',
  PROJECT_LABEL = 'mdi-link',
  SCHEMA = 'mdi-file-tree-outline',
  TAG = 'mdi-tag',
  BACKUP = 'mdi-card-multiple-outline',
  CHILD = 'mdi-account-supervisor-outline',
  OBLIGATION_REMARK = 'priority_high',
  OBLIGATION_WARNING = SharedIcons.OBLIGATION_WARNING,
  OBLIGATION_INFORMATION = SharedIcons.OBLIGATION_INFORMATION,
  OBLIGATION_ALARM = SharedIcons.OBLIGATION_ALARM,
  REMARK_WARNING = 'warning',
  REMARK_PROBLEM = 'error',
  REMARK_INFO = 'info_outline',
  INVALID = 'do_not_disturb_on',
  WARNING = SharedIcons.WARNING,
  ALLOW = SharedIcons.ALLOW,
  LINK = 'link',
  CIRCLE_FILLED = 'mdi-checkbox-blank-circle',
  GROUP = 'mdi-check',
  CANCEL = 'cancel',
  DENY = SharedIcons.DENY,
  QUESTIONED = SharedIcons.QUESTIONED,
  SUPPORT = 'support',
  NO_ASSERTION = SharedIcons.NO_ASSERTION,
  COMPONENTS = SharedIcons.COMPONENTS,
  ADD = SharedIcons.ADD,
  REMOVED = SharedIcons.REMOVED,
  CHANGED = 'mdi-compare-horizontal',
  SECURITY = 'mdi-shield-lock-outline',
  NON_FOSS = 'mdi-circle-off-outline',
}

export default Icons;
