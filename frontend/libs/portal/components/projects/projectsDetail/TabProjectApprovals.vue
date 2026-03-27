// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import GridApprovals from '@disclosure-portal/components/grids/GridApprovals.vue';
</script>
<template>
  <GridApprovals ref="tableDocuments" v-on="$attrs"></GridApprovals>
</template>
