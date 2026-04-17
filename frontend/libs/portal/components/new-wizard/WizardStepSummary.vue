<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {Department} from '@disclosure-portal/model/Department';
import {useLabelStore} from '@disclosure-portal/stores/label.store';
import {useWizardStore} from '@disclosure-portal/stores/wizard.store';
import {getStrWithMaxLength} from '@disclosure-portal/utils/View';
import {onBeforeMount} from 'vue';
import {useI18n} from 'vue-i18n';
import SummaryItem from './SummaryItem.vue';

const {t} = useI18n();
const wizardStore = useWizardStore();
const labelStore = useLabelStore();

const lineOrValue = (value: string): string => {
  return value && value.length > 0 ? value : '-';
};

const getDeptIdAsText = (dept: Department): string => {
  if (dept && dept.deptId) {
    return `[${dept.companyCode}] ${dept.companyName} / [${dept.deptId},${dept.orgAbbreviation}] ${dept.descriptionEnglish}`;
  } else {
    return '- / -';
  }
};

onBeforeMount(() => {
  wizardStore.preview();
});
</script>

<template>
  <Stack class="overflow-hidden">
    <h2 class="text-body-1 py-0">{{ t('WIZARD_page_summary') }}</h2>
    <div class="grid grid-cols-1 gap-4 md:grid-cols-3">
      <SummaryItem :label="t('WIZARD_project_name')" :value="wizardStore.project.name" />
      <div class="md:col-span-2">
        <SummaryItem :label="t('WIZARD_project_description')">
          <div v-if="!wizardStore.project.description || wizardStore.project.description.length === 0">-</div>
          <div v-else>
            {{ getStrWithMaxLength(250, wizardStore.project.description) }}
            <Tooltip :text="wizardStore.project.description" />
          </div>
        </SummaryItem>
      </div>
    </div>

    <div class="mt-4 grid grid-cols-1 gap-4 md:grid-cols-3">
      <SummaryItem :label="t('WIZARD_target_platform')" :value="wizardStore.project.targetPlatform" />
      <SummaryItem :label="t('WIZARD_architecture')" :value="wizardStore.project.architecture" />
      <SummaryItem :label="t('WIZARD_target_users')" :value="wizardStore.project.targetUsers" />
      <SummaryItem :label="t('WIZARD_distribution_target')" :value="wizardStore.project.distributionTarget" />
      <SummaryItem :label="t('WIZARD_application')">
        <div v-if="!wizardStore.project?.applicationMeta?.id">-</div>
        <div v-else>{{ wizardStore.project.applicationMeta.name }}</div>
      </SummaryItem>
    </div>

    <div class="mt-4 grid grid-cols-1 gap-4 md:grid-cols-3">
      <Stack>
        <SummaryItem
          :label="t('WIZARD_Developer_CompanyCode_Department')"
          :value="getDeptIdAsText(wizardStore.project.projectSettings?.documentMeta?.supplierDept)"
          :show-dash="false" />
        <SummaryItem
          :label="t('WIZARD_Developer_address')"
          :value="lineOrValue(wizardStore.project.projectSettings?.documentMeta?.supplierAddress)"
          :show-dash="false" />
        <SummaryItem :label="t('WIZARD_Developer_Id_And_Name')" :show-dash="false">
          ({{ lineOrValue(wizardStore.project.projectSettings?.documentMeta?.supplierNr) }})
          {{ lineOrValue(wizardStore.project.projectSettings?.documentMeta?.supplierName) }}
        </SummaryItem>
      </Stack>

      <Stack>
        <SummaryItem
          :label="t('WIZARD_Owner_CompanyCode_Department')"
          :value="getDeptIdAsText(wizardStore.project.projectSettings?.customerMeta?.dept)"
          :show-dash="false" />
        <SummaryItem
          :label="t('WIZARD_Owner_Address')"
          :value="lineOrValue(wizardStore.project.projectSettings?.customerMeta?.address)"
          :show-dash="false" />
      </Stack>

      <div v-if="!wizardStore.isVehicleOnboardArchitecture">
        <SummaryItem
          :label="t('WIZARD_Contact_Address')"
          :value="lineOrValue(wizardStore.project.projectSettings?.noticeContactMeta?.address)"
          :show-dash="false" />
      </div>
    </div>

    <div class="d-text d-secondary-text text-xs">{{ t('WIZARD_resulting_labelset') }}</div>
    <Stack v-if="wizardStore.previewLoading" direction="row">
      <div v-for="(_, i) in [0, 1, 2]" :key="i" class="h-7 w-[120px] animate-pulse rounded bg-gray-500"></div>
    </Stack>
    <Stack v-else direction="row" class="pb-4" wrap>
      <ProjectLabel v-for="label in wizardStore.project.labels" :key="label" :label="labelStore.getLabelByKey(label)" />
    </Stack>
  </Stack>
</template>
