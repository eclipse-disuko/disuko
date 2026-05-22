<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {useI18n} from 'vue-i18n';

interface Props {
  noFOSS: boolean;
  c1: boolean;
  c2: boolean;
  c3: boolean;
  c4: boolean;
  c5: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  'update:c1': [value: boolean];
  'update:c2': [value: boolean];
  'update:c3': [value: boolean];
  'update:c4': [value: boolean];
  'update:c5': [value: boolean];
}>();

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
      <v-checkbox
        :model-value="props.c1"
        :readonly="props.noFOSS"
        :label="t('SBOM_APPROVAL_CHECK1')"
        hide-details
        @update:model-value="emit('update:c1', $event as boolean)" />
      <v-checkbox
        :model-value="props.c2"
        :readonly="props.noFOSS"
        :label="t('SBOM_APPROVAL_CHECK2')"
        hide-details
        @update:model-value="emit('update:c2', $event as boolean)" />
      <v-checkbox
        :model-value="props.c3"
        :readonly="props.noFOSS"
        :label="t('SBOM_APPROVAL_CHECK3')"
        hide-details
        @update:model-value="emit('update:c3', $event as boolean)" />
      <v-checkbox
        :model-value="props.c4"
        :readonly="props.noFOSS"
        :label="t('SBOM_APPROVAL_CHECK4')"
        hide-details
        @update:model-value="emit('update:c4', $event as boolean)" />
      <v-checkbox
        :model-value="props.c5"
        :readonly="props.noFOSS"
        :label="t('SBOM_APPROVAL_CHECK5')"
        hide-details
        @update:model-value="emit('update:c5', $event as boolean)" />
      <v-checkbox :model-value="props.noFOSS" :label="t('SBOM_APPROVAL_CHECK6')" disabled></v-checkbox>
    </template>
  </ExpansionPanel>
</template>
