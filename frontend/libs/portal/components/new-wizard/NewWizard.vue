<script setup lang="ts">
import {useNewWizard} from '@disclosure-portal/composables/useNewWizard';
import {useWizardStore} from '@disclosure-portal/stores/wizard.store';

const wizardStore = useWizardStore();
const {componentMap} = useNewWizard();
</script>

<template>
  <WizardLayout>
    <component :is="componentMap[wizardStore.currentStep.id]" />
  </WizardLayout>
</template>
