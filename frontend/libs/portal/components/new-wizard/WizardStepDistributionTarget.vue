<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {WizardCardProps} from '@disclosure-portal/components/new-wizard/WizardCard.vue';
import {DistributionTarget, distributionTargets} from '@disclosure-portal/model/NewWizard';
import {useWizardStore} from '@disclosure-portal/stores/wizard.store';
import {ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const wizardStore = useWizardStore();

const cardList = ref<WizardCardProps[]>([
  {
    key: distributionTargets.company as string,
    image: new URL('@disclosure-portal/assets/wizard/projectDistribution/company.svg', import.meta.url).href,
    title: t('NEW_WIZARD_DISTRIBUTION_TARGET_MB'),
    subtitle: t('NEW_WIZARD_DISTRIBUTION_TARGET_MB_DESC'),
    helptext: t('NEW_WIZARD_DISTRIBUTION_TARGET_MB_HELP'),
    isFlipped: false,
    isActive: wizardStore.project?.distributionTarget === distributionTargets.company,
  },
  {
    key: distributionTargets.businessPartner as string,
    image: new URL('@disclosure-portal/assets/wizard/projectDistribution/business-partner.svg', import.meta.url).href,
    title: t('NEW_WIZARD_DISTRIBUTION_TARGET_BUSINESS'),
    subtitle: t('NEW_WIZARD_DISTRIBUTION_TARGET_BUSINESS_DESC'),
    helptext: t('NEW_WIZARD_DISTRIBUTION_TARGET_BUSINESS_HELP'),
    isFlipped: false,
    isActive: wizardStore.project?.distributionTarget === distributionTargets.businessPartner,
  },
]);

const onCardSelect = (distributionTarget: DistributionTarget) => {
  wizardStore.project.distributionTarget = distributionTarget;
  wizardStore.nextStep();
};
</script>

<template>
  <Stack class="overflow-hidden">
    <h2 class="text-body-1 py-0">{{ t('WIZARD_page_distribution_target_hint') }}</h2>
    <div class="mb-1 grid grid-cols-1 gap-4 sm:grid-cols-2">
      <WizardCard v-for="card in cardList" :key="card.key" :card="card" @update="onCardSelect($event)" />
    </div>
  </Stack>
</template>
