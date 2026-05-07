<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {useIdleStore} from '@shared/stores/idle.store';
import {DialogLayoutConfig} from '@shared/layouts/DialogLayout.vue';
import {watch} from 'vue';

const emit = defineEmits(['close', 'secondaryAction', 'primaryAction']);

export interface ReactiveDialogLayoutConfig {
  title: string;
  loading: boolean;
  titleTooltip?: string;
  secondaryButton?: string;
  primaryButton?: string;
  icon?: string;
  iconColor?: string; // optional icon color override
}

const props = defineProps<{
  config: ReactiveDialogLayoutConfig;
}>();

const idleStore = useIdleStore();

const dialogLayoutConfig: DialogLayoutConfig = {
  title: props.config.title,
  titleTooltip: props.config.titleTooltip,
  ...(props.config.secondaryButton ? {secondaryButton: {text: props.config.secondaryButton}} : {}),
  ...(props.config.primaryButton ? {primaryButton: {text: props.config.primaryButton}} : {}),
  ...(props.config.icon ? {icon: props.config.icon} : {}),
  ...(props.config.iconColor ? {iconColor: props.config.iconColor} : {}),
};

watch(
  () => props.config.loading,
  (loading: boolean) => {
    if (loading) {
      idleStore.show();
    } else {
      idleStore.hide();
    }
  },
);
</script>

<template>
  <DialogLayout
    :config="dialogLayoutConfig"
    @close="emit('close')"
    @secondaryAction="emit('secondaryAction')"
    @primaryAction="emit('primaryAction')">
    <template v-if="$slots['title-right']">
      <slot name="title-right"></slot>
    </template>
    <slot></slot>
    <template v-if="$slots.left">
      <slot name="left"></slot>
    </template>
  </DialogLayout>
</template>
