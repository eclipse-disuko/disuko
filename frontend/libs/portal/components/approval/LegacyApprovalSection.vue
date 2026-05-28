<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import ExpansionPanel from '@shared/components/disco/ExpansionPanel.vue';
import {useI18n} from 'vue-i18n';

withDefaults(
  defineProps<{
    noFOSS: boolean;
    isVehicle?: boolean;
  }>(),
  {
    isVehicle: false,
  },
);

const c1 = defineModel<boolean>('c1');
const c2 = defineModel<boolean>('c2');
const c3 = defineModel<boolean>('c3');
const c4 = defineModel<boolean>('c4');
const c5 = defineModel<boolean>('c5');
const radioGroup = defineModel<number>('radioGroup', {default: 0});

const {t} = useI18n();
</script>

<template>
  <div>
    <Stack direction="row" align="center">
      <v-icon v-if="noFOSS" size="small">mdi-alert</v-icon>
      <span class="d-block" v-if="noFOSS">{{ t('NO_FOSS_WARNING') }}</span>
    </Stack>
    <v-switch :model-value="noFOSS" color="primary" :label="t('NO_FOSS_MARKER')" hide-details disabled></v-switch>

    <ExpansionPanel :title="t('SBOM_APPROVAL_ATTRIBUTES')">
      <template #body>
        <Stack v-if="!isVehicle" class="gap-2">
          <v-checkbox v-model="c1" color="primary" :label="t('SBOM_APPROVAL_CHECK1')" hide-details density="compact" />
          <v-checkbox v-model="c2" color="primary" :label="t('SBOM_APPROVAL_CHECK2')" hide-details density="compact" />
          <v-checkbox v-model="c3" color="primary" :label="t('SBOM_APPROVAL_CHECK3')" hide-details density="compact" />
          <v-checkbox v-model="c4" color="primary" :label="t('SBOM_APPROVAL_CHECK4')" hide-details density="compact" />
          <v-checkbox v-model="c5" color="primary" :label="t('SBOM_APPROVAL_CHECK5')" hide-details density="compact" />
          <v-checkbox
            :model-value="noFOSS"
            color="primary"
            :label="t('SBOM_APPROVAL_CHECK6')"
            hide-details
            density="compact"
            disabled />
        </Stack>
        <!-- [FORK] -->
        <v-radio-group v-else v-model="radioGroup" hide-details>
          <v-radio :label="t('SBOM_APPROVAL_VEHICLE_CHECK2')" :value="2"></v-radio>
          <v-radio :label="t('SBOM_APPROVAL_VEHICLE_CHECK3')" :value="3"></v-radio>
        </v-radio-group>
      </template>
    </ExpansionPanel>
  </div>
</template>
