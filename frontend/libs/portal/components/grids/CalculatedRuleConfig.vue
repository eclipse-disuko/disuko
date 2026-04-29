<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {CalculatedRuleConfigType} from '@disclosure-portal/model/CalculatedPolicyRules';
import DMultiSelect from '@shared/components/disco/DMultiSelect.vue';
import Stack from '@shared/layouts/Stack.vue';
import {computed} from 'vue';
import {useI18n} from 'vue-i18n';

type CalculatedBucketName = 'deniedClassifications' | 'warnedClassifications' | 'allowedClassifications';
type ScopeFilterName = 'isLicenseChart' | 'approvalState' | 'family' | 'licenseType' | 'source';

interface Props {
  config: CalculatedRuleConfigType;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  'update-calculated': [value: boolean];
  'update-bucket': [payload: {bucketName: CalculatedBucketName; values: string[]}];
  'update-scope': [payload: {filterName: ScopeFilterName; values: Array<string | boolean>}];
}>();

const {t} = useI18n();

const classificationOptions = computed(() => props.config.classificationOptions);

const emitScope = (filterName: ScopeFilterName, values: unknown) => {
  emit('update-scope', {filterName, values: values as Array<string | boolean>});
};
</script>

<template>
  <v-card class="mb-2 w-full basis-full" variant="flat">
    <Stack v-if="config.calculated" class="gap-4">
      <v-card variant="flat" class="pa-3">
        <div class="d-subtitle-2 mb-4">{{ t('CALCULATED_BUCKETS_TITLE') }}</div>
        <div class="grid grid-cols-1 gap-3 md:grid-cols-3">
          <DMultiSelect
            class="pa-2 mx-2"
            :label="t('CALCULATED_DENIED_CLASSIFICATIONS')"
            :items="classificationOptions"
            :model-value="config.buckets.deniedClassifications"
            @update:modelValue="
              emit('update-bucket', {bucketName: 'deniedClassifications', values: $event as string[]})
            " />
          <DMultiSelect
            class="pa-2 mx-2"
            :label="t('CALCULATED_WARNED_CLASSIFICATIONS')"
            :items="classificationOptions"
            :model-value="config.buckets.warnedClassifications"
            @update:modelValue="
              emit('update-bucket', {bucketName: 'warnedClassifications', values: $event as string[]})
            " />
          <DMultiSelect
            class="pa-2 mx-2"
            :label="t('CALCULATED_ALLOWED_CLASSIFICATIONS')"
            :items="classificationOptions"
            :model-value="config.buckets.allowedClassifications"
            @update:modelValue="
              emit('update-bucket', {bucketName: 'allowedClassifications', values: $event as string[]})
            " />
        </div>
      </v-card>

      <v-card variant="flat" class="pa-3">
        <Stack direction="row" align="center" class="mb-4">
          <div class="d-subtitle-2">{{ t('CALCULATED_SCOPE_FILTERS_TITLE') }}</div>
        </Stack>
        <div class="grid grid-cols-1 gap-3 md:grid-cols-2 xl:grid-cols-3">
          <DMultiSelect
            class="pa-2 mx-2"
            :label="t('CALCULATED_SCOPE_LICENSE_CHART_INCLUDE')"
            :items="config.scopeConfig.isLicenseChart.options"
            :model-value="config.scopeConfig.isLicenseChart.values"
            @update:modelValue="emitScope('isLicenseChart', $event)" />

          <DMultiSelect
            class="pa-2 mx-2"
            :label="t('CALCULATED_SCOPE_APPROVAL_INCLUDE')"
            :items="config.scopeConfig.approvalState.options"
            :model-value="config.scopeConfig.approvalState.values"
            @update:modelValue="emitScope('approvalState', $event)" />

          <DMultiSelect
            class="pa-2 mx-2"
            :label="t('CALCULATED_SCOPE_FAMILY_INCLUDE')"
            :items="config.scopeConfig.family.options"
            :model-value="config.scopeConfig.family.values"
            @update:modelValue="emitScope('family', $event)" />

          <DMultiSelect
            class="pa-2 mx-2"
            :label="t('CALCULATED_SCOPE_TYPE_INCLUDE')"
            :items="config.scopeConfig.licenseType.options"
            :model-value="config.scopeConfig.licenseType.values"
            @update:modelValue="emitScope('licenseType', $event)" />

          <DMultiSelect
            class="pa-2 mx-2"
            :label="t('CALCULATED_SCOPE_SOURCE_INCLUDE')"
            :items="config.scopeConfig.source.options"
            :model-value="config.scopeConfig.source.values"
            @update:modelValue="emitScope('source', $event)" />
        </div>
      </v-card>
    </Stack>
  </v-card>
</template>
