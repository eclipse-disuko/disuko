<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<template>
  <span v-if="value && value !== '0001-01-01T00:00:00Z'" :class="$attrs.class">
    <span>{{ !showTime ? formatDate(value) : formatDateAndTime(value) }}</span>
    <Tooltip :text="formatDateAndTime(value)" />
  </span>
</template>

<script setup lang="ts">
import {formatDate, formatDateAndTime} from '@disclosure-portal/utils/Table';

defineOptions({inheritAttrs: false});

defineProps<{
  value: string | Date | undefined;
  showTime?: boolean;
}>();
</script>
