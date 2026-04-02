<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
interface Props {
  title: string;
  text?: string;
}
defineProps<Props>();
</script>

<template>
  <v-expansion-panels elevation="0">
    <v-expansion-panel>
      <v-expansion-panel-title height="10" max-height="10">
        <template #default="{expanded}">
          <v-icon color="primary" :icon="expanded ? 'mdi-chevron-up' : 'mdi-chevron-down'" class="mr-2" />
          <span class="text-caption">{{ title }}</span>
        </template>
        <template #actions></template>
      </v-expansion-panel-title>
      <v-expansion-panel-text>
        <span v-if="text" class="text-caption d-block mb-2">{{ text }}</span>
        <slot name="body" />
      </v-expansion-panel-text>
    </v-expansion-panel>
  </v-expansion-panels>
</template>
