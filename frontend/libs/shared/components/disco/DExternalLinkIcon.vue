<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ref} from 'vue';

interface Prop {
  url: string;
  hint: string;
}

defineProps<Prop>();

const isDialogVisible = ref(false);

const displayDialog = ($event: MouseEvent) => {
  $event.preventDefault();
  $event.stopPropagation();
  isDialogVisible.value = true;
};
</script>

<template>
  <span v-if="url">
    <a size="small" :href="url" target="_blank" variant="plain" plain @click="displayDialog($event)">
      <v-icon color="primary" class="mt-n1 mr-1" size="x-small">mdi mdi-open-in-new</v-icon>
      <Tooltip v-if="hint" :text="hint"></Tooltip>
    </a>
    <ExternalLinkDialog v-model:isDialogVisible="isDialogVisible" :url="url"></ExternalLinkDialog>
  </span>
</template>
