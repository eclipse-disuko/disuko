// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {useAppStore} from '@disclosure-portal/stores/app';
import {UseWindowSize} from '@vueuse/components';
import {useWindowSize} from '@vueuse/core';
import {ref, watch} from 'vue';
import {useRoute} from 'vue-router';
import Stack from './Stack.vue';

interface Props {
  hasTab?: boolean;
  hasTitle?: boolean;
  gap?: string;
}

const props = withDefaults(defineProps<Props>(), {
  hasTab: false,
  hasTitle: false,
  gap: '3',
});

const route = useRoute();
const windowSize = useWindowSize();
const appStore = useAppStore();

const desc = ref<HTMLElement | null>(null);
const buttons = ref<HTMLElement | null>(null);
const table = ref<HTMLElement | null>(null);
const layout = ref<HTMLElement | null>(null);

const calculateTotalContentHeight = (height: number, hasTab = props.hasTab, hasTitle = props.hasTitle) => {
  const notificationHeight = appStore.notificationClosed || !appStore.notificationMessage ? 0 : 32;
  const headerHeight = 56;
  const footerHeight = 40;
  const tabHeight = hasTab ? 48 + 32 : 0;
  const titleHeight = hasTitle ? 64 : 0;
  const gap = Number(props.gap ?? 0) * 4;
  const buttonsHeight = buttons.value ? buttons.value.clientHeight + 4 + gap : 0;
  const descHeight = desc.value ? desc.value.clientHeight + 12 : 0;
  const padding = 32; // p-4
  return (
    height -
    (tabHeight + padding + titleHeight + buttonsHeight + descHeight + headerHeight + footerHeight + notificationHeight)
  );
};

watch([() => route.path, windowSize.width], () => {
  if (table.value) {
    table.value.style.height = `${calculateTotalContentHeight(windowSize.height.value)}px`;
  }
});
</script>

<template>
  <Stack ref="layout" :class="`p-4 px-5 gap-${gap}`">
    <div v-if="$slots.description" ref="desc">
      <slot name="description"></slot>
    </div>
    <div v-if="$slots.buttons" ref="buttons" class="flex flex-row gap-3 items-center">
      <slot name="buttons"></slot>
    </div>
    <UseWindowSize v-slot="{height}">
      <div ref="table" class="overflow-auto" :style="`height:${calculateTotalContentHeight(height)}px`">
        <slot name="table"></slot>
      </div>
    </UseWindowSize>
  </Stack>
</template>
