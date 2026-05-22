<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {useI18n} from 'vue-i18n';

interface Props {
  modelValue: 'default' | 'legacy';
  disabled: boolean;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  'update:modelValue': [value: 'default' | 'legacy'];
}>();

const {t} = useI18n();
</script>

<template>
  <Stack direction="row" align="center" class="rounded bg-gray-500/20 py-1">
    <v-radio-group
      inline
      hide-details
      :model-value="props.modelValue"
      :disabled="props.disabled"
      @update:model-value="emit('update:modelValue', $event)">
      <v-radio :label="t('FOSSDD_STANDARD')" value="default"></v-radio>
      <v-radio :label="t('FOSSDD_LEGACY')" value="legacy"></v-radio>
    </v-radio-group>
    <v-spacer></v-spacer>
    <DIconButton icon="mdi-information-outline" :hint="t('FOSSDD_VERSION_TOOLTIP')" />
  </Stack>
</template>
