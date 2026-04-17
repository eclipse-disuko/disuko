<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {PolicyDecisionSlim} from '@disclosure-portal/model/PolicyDecision';
import {formatDateAndTime} from '@disclosure-portal/utils/Table';
import {getIconColorForPolicyType, getIconForPolicyType} from '@disclosure-portal/utils/View';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();

defineProps<{
  decisions: PolicyDecisionSlim[];
  title: string;
  icon: string;
  iconColor?: string;
  arrow: string;
}>();
</script>

<template>
  <div>
    <v-divider></v-divider>
    <v-icon size="small" :color="iconColor">{{ icon }}</v-icon>
    {{ title }}
    <br />
    <div v-for="(pd, index) in decisions" :key="index" class="d-text d-secondary-text">
      <v-icon>mdi-blank</v-icon>
      <v-icon size="small" :color="getIconColorForPolicyType(pd.policyEvaluated)">
        {{ getIconForPolicyType(pd.policyEvaluated) }}
      </v-icon>
      {{ ` ${arrow} ` }}
      <v-icon size="small" :color="getIconColorForPolicyType(pd.policyDecision)">
        {{ getIconForPolicyType(pd.policyDecision) }}
      </v-icon>
      {{
        `: ${pd.licenseId} (${t('TT_POLICY_DECISION_BY_AT', {
          creator: pd.creator,
          created: formatDateAndTime(pd.created),
        })})`
      }}
      <br />
    </div>
  </div>
</template>
