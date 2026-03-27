// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
interface SummaryItemProps {
  label: string;
  value?: string | null;
  showDash?: boolean;
}

const props = withDefaults(defineProps<SummaryItemProps>(), {
  showDash: true,
});

const displayValue = (val: string | null | undefined): string => {
  if (!val || val.length === 0) {
    return props.showDash ? '-' : '';
  }
  return val;
};
</script>

<template>
  <div>
    <div class="text-xs d-text d-secondary-text mb-2">{{ label }}</div>
    <div>
      <slot>{{ displayValue(value) }}</slot>
    </div>
  </div>
</template>
