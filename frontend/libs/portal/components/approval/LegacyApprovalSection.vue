<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {useI18n} from 'vue-i18n';

interface Props {
  noFOSS: boolean;
  isVehicle: boolean;
}

const props = defineProps<Props>();
const c1 = defineModel<boolean>('c1', {default: false});
const c2 = defineModel<boolean>('c2', {default: false});
const c3 = defineModel<boolean>('c3', {default: false});
const c4 = defineModel<boolean>('c4', {default: false});
const c5 = defineModel<boolean>('c5', {default: false});
const radioGroup = defineModel<number>('radioGroup', {default: 0});

const {t} = useI18n();
</script>

<template>
  <Stack class="gap-1">
    <Stack v-if="props.noFOSS" direction="row" align="center">
      <v-icon size="small" color="warning">mdi-alert</v-icon>
      <span class="d-block">{{ t('NO_FOSS_WARNING') }}</span>
    </Stack>
    <Stack direction="row" align="center">
      <v-icon size="small" color="warning">mdi-alert</v-icon>
      <span class="d-block">{{ t('NO_FOSS_DISABLED_TOOLTIP') }}</span>
    </Stack>
    <v-switch
      class="ml-2"
      :model-value="props.noFOSS"
      color="primary"
      density="compact"
      :label="t('NO_FOSS_MARKER')"
      hide-details
      disabled></v-switch>
  </Stack>

  <ExpansionPanel :title="t('SBOM_APPROVAL_ATTRIBUTES')">
    <template #body>
      <template v-if="props.isVehicle">
        <v-radio-group v-model="radioGroup">
          <v-radio
            :label="t('SBOM_APPROVAL_VEHICLE_CHECK2')"
            :key="2"
            :value="2"
            class="py-2"
            :readonly="props.noFOSS"></v-radio>
          <v-radio
            :label="t('SBOM_APPROVAL_VEHICLE_CHECK3')"
            :key="3"
            :value="3"
            class="py-2"
            :readonly="props.noFOSS"
            disabled></v-radio>
        </v-radio-group>
      </template>
      <template v-else>
        <v-checkbox v-model="c1" :readonly="props.noFOSS" :label="t('SBOM_APPROVAL_CHECK1')" hide-details />
        <v-checkbox v-model="c2" :readonly="props.noFOSS" :label="t('SBOM_APPROVAL_CHECK2')" hide-details />
        <v-checkbox v-model="c3" :readonly="props.noFOSS" :label="t('SBOM_APPROVAL_CHECK3')" hide-details />
        <v-checkbox v-model="c4" :readonly="props.noFOSS" :label="t('SBOM_APPROVAL_CHECK4')" hide-details />
        <v-checkbox v-model="c5" :readonly="props.noFOSS" :label="t('SBOM_APPROVAL_CHECK5')" hide-details />
        <v-checkbox :model-value="props.noFOSS" :label="t('SBOM_APPROVAL_CHECK6')" disabled></v-checkbox>
      </template>
    </template>
  </ExpansionPanel>
</template>
