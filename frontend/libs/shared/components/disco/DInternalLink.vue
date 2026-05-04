<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
interface Props {
  url: string;
  text: string;
}

defineProps<Props>();
</script>

<template>
  <span>
    <a color="primary" class="d-subtitle-2" :href="url" target="_blank" v-if="url">
      {{ text }}
    </a>
  </span>
</template>
