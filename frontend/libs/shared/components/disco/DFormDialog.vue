<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<template>
  <v-dialog v-model="showDialog" :key="renderKey" :max-width="props.maxWidth" :persistent="props.persistent">
    <slot />
  </v-dialog>
</template>

<script setup lang="ts">
import {ref, watch} from 'vue';

const props = withDefaults(
  defineProps<{
    persistent?: boolean;
    maxWidth?: string | number;
  }>(),
  {maxWidth: 500},
);

const showDialog = defineModel<boolean>('showDialog');

// Reactive variable for forcing re-render
const renderKey = ref(0);

// Watch the dialog visibility to increment renderKey when opened
watch(
  () => showDialog.value,
  (newVal) => {
    if (newVal) {
      // Re-render dialog content each time it is opened
      renderKey.value++;
    }
  },
);
</script>
