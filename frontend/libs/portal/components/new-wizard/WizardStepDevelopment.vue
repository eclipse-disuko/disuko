<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {WizardCardProps} from '@disclosure-portal/components/new-wizard/WizardCard.vue';
import {Development, developments} from '@disclosure-portal/model/NewWizard';
import {useWizardStore} from '@disclosure-portal/stores/wizard.store';
import {ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const wizardStore = useWizardStore();

const cardList = ref<WizardCardProps[]>([
  {
    key: developments.inhouse,
    image: new URL('@disclosure-portal/assets/wizard/projectDevelopment/in-house.svg', import.meta.url).href,
    title: t('WIZARD_page_developer_target_inhouse'),
    subtitle: t('WIZARD_developer_target_inhouse_description'),
    helptext: t('TT_WIZARD_TD_IN_HOUSE_DEVELOPMENT'),
    isFlipped: false,
    isActive: wizardStore.project?.development === developments.inhouse,
  },
  {
    key: developments.internal,
    image: new URL('@disclosure-portal/assets/wizard/projectDevelopment/internal-supplier.svg', import.meta.url).href,
    title: t('WIZARD_page_developer_target_internal'),
    subtitle: t('WIZARD_developer_target_internal_description'),
    helptext: t('TT_WIZARD_TD_INTERNAL_DEVELOPMENT'),
    isFlipped: false,
    isActive: wizardStore.project?.development === developments.internal,
  },
  {
    key: developments.external,
    image: new URL('@disclosure-portal/assets/wizard/projectDevelopment/external-supplier.svg', import.meta.url).href,
    title: t('WIZARD_page_developer_target_external'),
    subtitle: t('WIZARD_developer_target_external_description'),
    helptext: t('TT_WIZARD_TD_EXTERNAL_DEVELOPMENT'),
    isFlipped: false,
    isActive: wizardStore.project?.development === developments.external,
  },
]);

const onCardSelect = (development: Development) => {
  if (wizardStore.project.development === development) {
    wizardStore.nextStep();
    return;
  }

  wizardStore.project.development = development;
  if (wizardStore.project.projectSettings?.documentMeta?.supplierName) {
    wizardStore.project.projectSettings.documentMeta.supplierName = '';
  }
  if (Object.keys(wizardStore.project.projectSettings?.documentMeta?.supplierDept ?? {}).length !== 0) {
    wizardStore.project.projectSettings.documentMeta.supplierDept = null;
  }
  wizardStore.updateProjectSettingsBasedOnDevelopment();
  wizardStore.setAvailableSteps();
  wizardStore.nextStep();
};
</script>

<template>
  <Stack class="overflow-hidden">
    <h2 class="text-body-1 py-0">{{ t('WIZARD_page_developer_target_hint') }}</h2>
    <div class="mb-1 grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3">
      <WizardCard v-for="card in cardList" :key="card.key" :card="card" @update="onCardSelect($event)" />
    </div>
  </Stack>
</template>
