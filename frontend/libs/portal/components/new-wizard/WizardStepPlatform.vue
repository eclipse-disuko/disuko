// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {WizardCardProps} from '@disclosure-portal/components/new-wizard/WizardCard.vue';
import {type TargetPlatform, targetPlatforms} from '@disclosure-portal/model/NewWizard';
import {useWizardStore} from '@disclosure-portal/stores/wizard.store';
import {ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const wizardStore = useWizardStore();
const cardList = ref<WizardCardProps[]>([
  {
    key: targetPlatforms.enterprise,
    image: new URL('@disclosure-portal/assets/wizard/projectPlatform/enterpriseIT.svg', import.meta.url).href,
    title: t('WIZARD_target_enterprise'),
    subtitle: t('WIZARD_target_enterprise_description'),
    helptext: t('TT_WIZARD_TP_ENTERPRISE_IT'),
    isFlipped: false,
    isActive: wizardStore.project?.targetPlatform === targetPlatforms.enterprise,
  },
  {
    key: targetPlatforms.mobile,
    image: new URL('@disclosure-portal/assets/wizard/projectPlatform/mobile.svg', import.meta.url).href,
    title: t('WIZARD_target_mobile'),
    subtitle: t('WIZARD_target_mobile_description'),
    helptext: t('TT_WIZARD_TP_MOBILE'),
    isFlipped: false,
    isActive: wizardStore.project?.targetPlatform === targetPlatforms.mobile,
  },
  {
    key: targetPlatforms.vehicle,
    image: new URL('@disclosure-portal/assets/wizard/projectPlatform/product.svg', import.meta.url).href,
    title: t('WIZARD_target_vehicle'),
    subtitle: t('WIZARD_target_vehicle_description'),
    helptext: t('TT_WIZARD_TP_VEHICLE'),
    isFlipped: false,
    isActive: wizardStore.project?.targetPlatform === targetPlatforms.vehicle,
  },
  {
    key: targetPlatforms.other,
    image: new URL('@disclosure-portal/assets/wizard/projectPlatform/other.svg', import.meta.url).href,
    title: t('WIZARD_target_other'),
    subtitle: t('WIZARD_target_other_description'),
    helptext: t('TT_WIZARD_TP_OTHER'),
    isFlipped: false,
    isActive: wizardStore.project?.targetPlatform === targetPlatforms.other,
  },
]);

const resetIncompatibleFields = () => {
  // Clear application meta because it's not intended to be filled
  if (wizardStore.isVehiclePlatform || wizardStore.isOtherPlatform) {
    wizardStore.project.applicationMeta = null;
  }
  // Clear fields if user had different platform & fields selected
  if (wizardStore.isVehiclePlatform && !wizardStore.isVehicleArchitectures) {
    wizardStore.project.architecture = null;
    wizardStore.project.targetUsers = null;
    wizardStore.project.distributionTarget = null;
  } else if (!wizardStore.isVehiclePlatform && wizardStore.isVehicleArchitectures) {
    wizardStore.project.architecture = null;
  }

  if (wizardStore.isOtherPlatform) {
    wizardStore.project.architecture = null;
    wizardStore.project.targetUsers = null;
    wizardStore.project.distributionTarget = null;
  }
};

const onCardSelect = (targetPlatform: TargetPlatform) => {
  wizardStore.project.targetPlatform = targetPlatform;
  resetIncompatibleFields();
  wizardStore.setAvailableSteps();
  wizardStore.nextStep();
};
</script>

<template>
  <Stack class="overflow-hidden">
    <h2 class="text-body-1 py-0">{{ t('WIZARD_page_select_target') }}</h2>
    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4 mb-1">
      <WizardCard v-for="card in cardList" :key="card.key" :card="card" @update="onCardSelect($event)" />
    </div>
  </Stack>
</template>
