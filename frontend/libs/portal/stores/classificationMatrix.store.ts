// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {IObligation} from '@disclosure-portal/model/IObligation';
import PolicyRule from '@disclosure-portal/model/PolicyRule';
import AdminService from '@disclosure-portal/services/admin';
import policyRuleService from '@disclosure-portal/services/policyrules';
import {
  RuleStatus,
  getStatusColor,
  getStatusIcon,
  getStatusLabel,
  toRuleStatusMap,
} from '@disclosure-portal/utils/PolicyRuleClassification';
import {defineStore} from 'pinia';
import {computed, reactive, toRefs} from 'vue';
import {useI18n} from 'vue-i18n';

export const policyRulesMatrixStore = defineStore('classificationMatrix', () => {
  const {t} = useI18n();

  const state = reactive({
    classifications: [] as IObligation[],
    classificationsLoaded: false,
    policyRules: [] as PolicyRule[],
    isLoading: false,
  });

  const policyRulesWithStatus = computed(() =>
    state.policyRules.map((rule) => {
      const rules = toRuleStatusMap(rule);
      return {
        ...rule,
        rules,
        statusProps: Object.fromEntries(
          Object.entries(rules).map(([classKey, status]) => [
            classKey,
            status
              ? {
                  icon: getStatusIcon(status as RuleStatus),
                  color: getStatusColor(status as RuleStatus),
                  label: getStatusLabel(status as RuleStatus, t),
                }
              : undefined,
          ]),
        ),
      };
    }),
  );

  const loadClassifications = async (): Promise<void> => {
    if (state.classificationsLoaded) return;
    try {
      const response = (await AdminService.getAllObligations()).data;
      state.classifications = response.items ?? response ?? [];
    } catch {
      state.classifications = [];
    } finally {
      state.classificationsLoaded = true;
    }
  };

  const loadMatrix = async (): Promise<void> => {
    state.isLoading = true;
    try {
      const [rulesRes] = await Promise.all([policyRuleService.getAllPolicyRules(), loadClassifications()]);
      state.policyRules = (rulesRes.data ?? []).filter((rule) => rule.calculated);
    } finally {
      state.isLoading = false;
    }
  };

  return {
    ...toRefs(state),
    policyRulesWithStatus,
    loadMatrix,
  };
});