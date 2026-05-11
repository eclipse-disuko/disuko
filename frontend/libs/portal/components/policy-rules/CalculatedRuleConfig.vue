<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {storeToRefs} from 'pinia';
import {computed} from 'vue';
import {useI18n} from 'vue-i18n';
import {useCalculatedPolicyRuleStore} from '@disclosure-portal/stores/calculatedPolicyRule.store';

type CalculatedBucketName = 'deniedClassifications' | 'warnedClassifications' | 'allowedClassifications';
type ScopeFilterName = 'isLicenseChart' | 'approvalState' | 'family' | 'licenseType' | 'source';

const {t} = useI18n();
const calculatedPolicyRuleStore = useCalculatedPolicyRuleStore();
const {calculatedRuleConfig} = storeToRefs(calculatedPolicyRuleStore);

const config = computed(() => calculatedRuleConfig.value);
const classificationOptions = computed(() => config.value.classificationOptions);

const handleBucketUpdate = (bucketName: CalculatedBucketName, values: unknown) => {
  calculatedPolicyRuleStore.setBucketClassifications(bucketName, values as string[]);
};

const handleScopeUpdate = (filterName: ScopeFilterName, values: unknown) => {
  calculatedPolicyRuleStore.setScopeFilterValues(filterName, values as Array<string | boolean>);
};
</script>

<template>
  <v-card class="mb-2 w-full basis-full" variant="flat">
    <Stack v-if="config.calculated" class="gap-6">
      <v-card variant="flat" class="pa-4">
        <div class="d-subtitle-2 mb-5">{{ t('CALCULATED_BUCKETS_TITLE') }}</div>
        <div class="grid grid-cols-1 gap-x-4 gap-y-6 md:grid-cols-2 xl:grid-cols-3">
          <DMultiSelect
            :label="t('CALCULATED_DENIED_CLASSIFICATIONS')"
            :items="classificationOptions"
            :model-value="config.buckets.deniedClassifications"
            @update:modelValue="handleBucketUpdate('deniedClassifications', $event)" />
          <DMultiSelect
            :label="t('CALCULATED_WARNED_CLASSIFICATIONS')"
            :items="classificationOptions"
            :model-value="config.buckets.warnedClassifications"
            @update:modelValue="handleBucketUpdate('warnedClassifications', $event)" />
          <DMultiSelect
            :label="t('CALCULATED_ALLOWED_CLASSIFICATIONS')"
            :items="classificationOptions"
            :model-value="config.buckets.allowedClassifications"
            @update:modelValue="handleBucketUpdate('allowedClassifications', $event)" />
        </div>
      </v-card>

      <v-card variant="flat" class="pa-4">
        <div class="d-subtitle-2 mb-5">{{ t('CALCULATED_SCOPE_FILTERS_TITLE') }}</div>
        <div class="grid grid-cols-1 gap-x-4 gap-y-6 md:grid-cols-2">
          <DMultiSelect
            :label="t('CALCULATED_SCOPE_LICENSE_CHART_INCLUDE')"
            :items="config.scopeConfig.isLicenseChart.options"
            :model-value="config.scopeConfig.isLicenseChart.values"
            @update:modelValue="handleScopeUpdate('isLicenseChart', $event)" />

          <DMultiSelect
            :label="t('CALCULATED_SCOPE_APPROVAL_INCLUDE')"
            :items="config.scopeConfig.approvalState.options"
            :model-value="config.scopeConfig.approvalState.values"
            @update:modelValue="handleScopeUpdate('approvalState', $event)" />

          <DMultiSelect
            :label="t('CALCULATED_SCOPE_FAMILY_INCLUDE')"
            :items="config.scopeConfig.family.options"
            :model-value="config.scopeConfig.family.values"
            @update:modelValue="handleScopeUpdate('family', $event)" />

          <DMultiSelect
            :label="t('CALCULATED_SCOPE_TYPE_INCLUDE')"
            :items="config.scopeConfig.licenseType.options"
            :model-value="config.scopeConfig.licenseType.values"
            @update:modelValue="handleScopeUpdate('licenseType', $event)" />

          <DMultiSelect
            :label="t('CALCULATED_SCOPE_SOURCE_INCLUDE')"
            :items="config.scopeConfig.source.options"
            :model-value="config.scopeConfig.source.values"
            @update:modelValue="handleScopeUpdate('source', $event)" />
        </div>
      </v-card>
    </Stack>
  </v-card>
</template>
