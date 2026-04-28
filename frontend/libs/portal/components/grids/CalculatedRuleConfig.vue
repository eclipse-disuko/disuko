<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {IDefaultSelectItem} from '@disclosure-portal/model/IObligation';
import DMultiSelect from '@shared/components/disco/DMultiSelect.vue';
import Stack from '@shared/layouts/Stack.vue';
import {useI18n} from 'vue-i18n';

type CalculatedBucketName = 'deniedClassifications' | 'warnedClassifications' | 'allowedClassifications';
type ScopeFilterName = 'isLicenseChart' | 'approvalState' | 'family' | 'licenseType' | 'source';

interface Props {
  calculated: boolean;
  classificationOptions: IDefaultSelectItem[];
  classificationOptionsLoaded: boolean;
  deniedClassifications: string[];
  warnedClassifications: string[];
  allowedClassifications: string[];
  calculatedIsLicenseChartOptions: IDefaultSelectItem[];
  calculatedApprovalOptions: IDefaultSelectItem[];
  calculatedFamilyOptions: IDefaultSelectItem[];
  calculatedTypeOptions: IDefaultSelectItem[];
  calculatedSourceOptions: IDefaultSelectItem[];
  getScopeFilterValues: (filterName: ScopeFilterName) => Array<string | boolean>;
}

defineProps<Props>();

const emit = defineEmits<{
  'update-calculated': [value: boolean];
  'update-bucket': [payload: {bucketName: CalculatedBucketName; values: string[]}];
  'update-scope': [payload: {filterName: ScopeFilterName; values: Array<string | boolean>}];
}>();

const {t} = useI18n();

const emitScope = (filterName: ScopeFilterName, values: unknown) => {
  emit('update-scope', {filterName, values: values as Array<string | boolean>});
};
</script>

<template>
  <v-card class="mb-2 w-full basis-full" variant="flat">
    <Stack v-if="calculated" class="gap-4">
      <v-card variant="flat" class="pa-3">
        <div class="d-subtitle-2 mb-4">{{ t('CALCULATED_BUCKETS_TITLE') }}</div>
        <div class="grid grid-cols-1 gap-3 md:grid-cols-3">
          <DMultiSelect
            class="pa-2 mx-2"
            :label="t('CALCULATED_DENIED_CLASSIFICATIONS')"
            :items="classificationOptionsLoaded ? classificationOptions : []"
            :model-value="classificationOptionsLoaded ? deniedClassifications : []"
            @update:modelValue="
              emit('update-bucket', {bucketName: 'deniedClassifications', values: $event as string[]})
            " />
          <DMultiSelect
            class="pa-2 mx-2"
            :label="t('CALCULATED_WARNED_CLASSIFICATIONS')"
            :items="classificationOptionsLoaded ? classificationOptions : []"
            :model-value="classificationOptionsLoaded ? warnedClassifications : []"
            @update:modelValue="
              emit('update-bucket', {bucketName: 'warnedClassifications', values: $event as string[]})
            " />
          <DMultiSelect
            class="pa-2 mx-2"
            :label="t('CALCULATED_ALLOWED_CLASSIFICATIONS')"
            :items="classificationOptionsLoaded ? classificationOptions : []"
            :model-value="classificationOptionsLoaded ? allowedClassifications : []"
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
            :items="calculatedIsLicenseChartOptions"
            :model-value="getScopeFilterValues('isLicenseChart')"
            @update:modelValue="emitScope('isLicenseChart', $event)" />

          <DMultiSelect
            class="pa-2 mx-2"
            :label="t('CALCULATED_SCOPE_APPROVAL_INCLUDE')"
            :items="calculatedApprovalOptions"
            :model-value="getScopeFilterValues('approvalState')"
            @update:modelValue="emitScope('approvalState', $event)" />

          <DMultiSelect
            class="pa-2 mx-2"
            :label="t('CALCULATED_SCOPE_FAMILY_INCLUDE')"
            :items="calculatedFamilyOptions"
            :model-value="getScopeFilterValues('family')"
            @update:modelValue="emitScope('family', $event)" />

          <DMultiSelect
            class="pa-2 mx-2"
            :label="t('CALCULATED_SCOPE_TYPE_INCLUDE')"
            :items="calculatedTypeOptions"
            :model-value="getScopeFilterValues('licenseType')"
            @update:modelValue="emitScope('licenseType', $event)" />

          <DMultiSelect
            class="pa-2 mx-2"
            :label="t('CALCULATED_SCOPE_SOURCE_INCLUDE')"
            :items="calculatedSourceOptions"
            :model-value="getScopeFilterValues('source')"
            @update:modelValue="emitScope('source', $event)" />
        </div>
      </v-card>
    </Stack>
  </v-card>
</template>
