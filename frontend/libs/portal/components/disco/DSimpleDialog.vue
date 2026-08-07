<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<template>
  <v-dialog v-model="localValue" width="auto" max-width="1000px" content-class="d-simple-dialog" scrollable>
    <DialogLayout :config="dialogConfig" @close="localValue = false">
      <slot />
    </DialogLayout>
  </v-dialog>
</template>

<script setup lang="ts">
import DialogLayout, {type DialogLayoutConfig} from '@shared/layouts/DialogLayout.vue';
import {computed} from 'vue';

interface Props {
  title: string;
}

const props = defineProps<Props>();
const localValue = defineModel<boolean>({required: true});

const dialogConfig = computed<DialogLayoutConfig>(() => ({
  title: props.title,
}));
</script>
