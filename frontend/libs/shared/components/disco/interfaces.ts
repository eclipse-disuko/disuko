// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {PolicyState} from '@disclosure-portal/model/PolicyRule';

export interface IRuleBtnCallbacks {
  getUrlToComponents(policy: PolicyState): string | null;
  handlePolicySelect(policy: PolicyState, selectedFilterPolicyTypes: PolicyState[]): void;
  getCountForPolicyFilterBtn(policy: PolicyState): number;
  getToolTipKeyForPolicyFilterBtn(policy: PolicyState): string;
  getActiveClassForPolicyFilterBtn(policy: PolicyState): string;
}
