<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {useWizardStore} from '@disclosure-portal/stores/wizard.store';
import {useDisplay} from 'vuetify/framework';

const wizardStore = useWizardStore();
const {smAndDown, height} = useDisplay();
</script>

<template>
  <v-dialog
    v-model="wizardStore.isWizardOpen"
    scrollable
    :width="smAndDown ? '100%' : '960'"
    :height="smAndDown ? '100%' : height <= 820 ? '100%' : '820'"
    persistent
    :fullscreen="smAndDown">
    <NewWizard data-testid="project-wizard-dialog"></NewWizard>
  </v-dialog>
</template>

<style scoped></style>
