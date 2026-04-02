<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {WizardCardProps} from '@disclosure-portal/components/new-wizard/WizardCard.vue';
import {Architecture, architectures, targetPlatforms} from '@disclosure-portal/model/NewWizard';
import {useWizardStore} from '@disclosure-portal/stores/wizard.store';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const wizardStore = useWizardStore();

const isVehicleArchitecture = computed(() => wizardStore.project?.targetPlatform === targetPlatforms.vehicle);

const architectureCardList = ref<WizardCardProps[]>([
  {
    key: architectures.frontend,
    image: new URL('@disclosure-portal/assets/wizard/projectArchitecture/frontend.svg', import.meta.url).href,
    title: t('WIZARD_architecture_frontend'),
    subtitle: t('WIZARD_architecture_frontend_description'),
    helptext: t('TT_WIZARD_ARCH_FRONTEND_OR_CLIENT'),
    isFlipped: false,
    isActive: wizardStore.project?.architecture === architectures.frontend,
  },
  {
    key: architectures.backend,
    image: new URL('@disclosure-portal/assets/wizard/projectArchitecture/backend.svg', import.meta.url).href,
    title: t('WIZARD_architecture_backend'),
    subtitle: t('WIZARD_architecture_backend_description'),
    helptext: t('TT_WIZARD_ARCH_BACKEND'),
    isFlipped: false,
    isActive: wizardStore.project?.architecture === architectures.backend,
  },
]);

const architectureVehicleCardList = ref<WizardCardProps[]>([
  {
    key: architectures.vehicleOnboard,
    image: new URL('@disclosure-portal/assets/wizard/projectArchitecture/product-onboard.svg', import.meta.url).href,
    title: t('WIZARD_vehicle_architecture_onboard'),
    subtitle: t('WIZARD_vehicle_architecture_onboard_description'),
    helptext: t('TT_WIZARD_ARCH_VEHICLE_ONBOARD'),
    isFlipped: false,
    isActive: wizardStore.project?.architecture === architectures.vehicleOnboard,
  },
  {
    key: architectures.vehicleOffboard,
    image: new URL('@disclosure-portal/assets/wizard/projectArchitecture/product-offboard.svg', import.meta.url).href,
    title: t('WIZARD_vehicle_architecture_offboard'),
    subtitle: t('WIZARD_vehicle_architecture_offboard_description'),
    helptext: t('TT_WIZARD_ARCH_VEHICLE_OFFBOARD'),
    isFlipped: false,
    isActive: wizardStore.project?.architecture === architectures.vehicleOffboard,
  },
]);

const cardList = computed(() =>
  isVehicleArchitecture.value ? architectureVehicleCardList.value : architectureCardList.value,
);

const onCardSelect = (architecture: Architecture) => {
  wizardStore.project.architecture = architecture;
  wizardStore.nextStep();
};
</script>

<template>
  <Stack class="overflow-hidden">
    <h2 class="text-body-1 py-0">{{ t('WIZARD_page_architecture_hint') }}</h2>
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-1">
      <WizardCard v-for="card in cardList" :key="card.key" :card="card" @update="onCardSelect($event)" />
    </div>
  </Stack>
</template>
