// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {Decision} from '@disclosure-portal/model/Decision';
import {PolicyDecisionSlim} from '@disclosure-portal/model/PolicyDecision';
import {getIconColorForPolicyType, getIconForPolicyType} from '@disclosure-portal/utils/View';

defineProps<{
  decision: PolicyDecisionSlim | Decision;
  previewMode: boolean;
  licenseId: string;
}>();
</script>

<template>
  <span>
    <v-icon size="small" :color="getIconColorForPolicyType(decision.policyEvaluated)">
      {{ getIconForPolicyType(decision.policyEvaluated) }}
    </v-icon>
    {{ previewMode ? ' ⇢ ' : ' → ' }}
    <v-icon size="small" :color="getIconColorForPolicyType(decision.policyDecision)">
      {{ getIconForPolicyType(decision.policyDecision) }}
    </v-icon>
    {{ `: ${licenseId}` }}
  </span>
</template>
