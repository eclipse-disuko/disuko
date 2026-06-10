<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
defineProps<{
  icon: string;
  value: number;
  color?: string;
  nowrap?: boolean;
}>();
</script>

<template>
  <div class="flex flex-col items-center justify-center">
    <v-icon size="small" :color="color">{{ icon }}</v-icon>
    <div :class="['pt-1', 'text-center', {'text-gray-500': value === 0, 'text-no-wrap': nowrap}]">
      {{ value }}
    </div>
  </div>
</template>
