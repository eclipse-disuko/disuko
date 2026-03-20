// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {PolicyState} from '@disclosure-portal/model/PolicyRule';
import DRuleButtons from '@shared/components/disco/DRuleButtons.vue';

export interface IRuleBtnCallbacks {
  getUrlToComponents(policy: PolicyState): string | null;
  handlePolicySelect(policy: PolicyState, selectedFilterPolicyTypes: PolicyState[]): void;
  getCountForPolicyFilterBtn(policy: PolicyState): number;
  getInitSelectedPolicy(): PolicyState;
  getToolTipKeyForPolicyFilterBtn(policy: PolicyState): string;
  getActiveClassForPolicyFilterBtn(policy: PolicyState): string;
  setRuleButtons(ruleButtons: InstanceType<typeof DRuleButtons>): void;
}
