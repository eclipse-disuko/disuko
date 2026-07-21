<!-- SPDX-FileCopyrightText: 2026 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import SettingsMenu from '@cli/components/SettingsMenu.vue';
</script>

<template>
  <v-app>
    <div class="fixed top-3 right-3 z-[1100]">
      <SettingsMenu />
    </div>

    <v-main>
      <slot></slot>
    </v-main>

    <DSnackbar />
  </v-app>
</template>
