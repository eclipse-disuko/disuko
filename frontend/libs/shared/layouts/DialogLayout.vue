// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import {useI18n} from 'vue-i18n';

const emit = defineEmits(['close', 'secondaryAction', 'primaryAction']);

export interface DialogLayoutConfig {
  title: string;
  titleTooltip?: string;
  secondaryButton?: {text: string; disabled?: boolean; loading?: boolean};
  primaryButton?: {text: string; disabled?: boolean; loading?: boolean};
  icon?: string;
  iconColor?: string; // optional icon color override
}

defineProps<{
  config: DialogLayoutConfig;
}>();

const {t} = useI18n();
</script>

<template>
  <v-card class="p-12">
    <Stack direction="row" align="center">
      <v-icon v-if="config?.icon" :icon="config.icon" :color="config?.iconColor || 'primary'" />
      <h4 class="text-h5 truncate">
        {{ config.title }}
        <Tooltip>{{ config?.titleTooltip ? config.titleTooltip : config.title }}</Tooltip>
      </h4>
      <template v-if="$slots['title-right']">
        <slot name="title-right"></slot>
      </template>
      <v-spacer></v-spacer>
      <DCloseButton class="-mr-4" @click="emit('close')" />
    </Stack>

    <v-card-text class="p-0 pt-8">
      <slot></slot>
    </v-card-text>

    <Stack direction="row" class="pt-8" align="center">
      <template v-if="$slots.left">
        <slot name="left"></slot>
      </template>

      <v-spacer></v-spacer>

      <template v-if="$slots.right">
        <slot name="right"></slot>
      </template>

      <template v-else>
        <DCActionButton
          v-if="!config?.secondaryButton && !config?.primaryButton"
          is-dialog-button
          size="small"
          variant="flat"
          @click="emit('close')"
          :text="t('BTN_CLOSE')" />

        <DCActionButton
          v-if="config?.secondaryButton"
          is-dialog-button
          size="small"
          variant="text"
          @click="emit('secondaryAction')"
          :disabled="config?.secondaryButton?.disabled"
          :loading="config?.secondaryButton?.loading"
          :text="config?.secondaryButton?.text" />

        <DCActionButton
          v-if="config?.primaryButton"
          is-dialog-button
          size="small"
          variant="flat"
          @click="emit('primaryAction')"
          :disabled="config?.primaryButton?.disabled"
          :loading="config?.primaryButton?.loading"
          :text="config?.primaryButton?.text" />
      </template>
    </Stack>
  </v-card>
</template>
