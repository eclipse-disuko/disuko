// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';

interface Props {
  url: string;
  text: string;
  tooltip?: boolean;
}
const props = defineProps<Props>();

const {t} = useI18n();

const isDialogVisible = ref(false);

const isExternalUrl = computed(() => {
  return props.url.startsWith('http://') || props.url.startsWith('https://');
});

const handleClick = ($event: MouseEvent) => {
  if (!isExternalUrl.value) {
    return;
  }

  $event.preventDefault();
  $event.stopPropagation();
  isDialogVisible.value = true;
};
</script>

<template>
  <a v-if="url" :href="url" target="_blank" @click="handleClick($event)">
    <v-icon color="primary" class="mt-n1 mr-1" size="x-small">mdi mdi-open-in-new</v-icon>
    <Tooltip v-if="tooltip" :text="`${t('OPEN_URL_EXTERN')} ${url}`"></Tooltip>
    {{ text }}
  </a>
  <ExternalLinkDialog v-model:isDialogVisible="isDialogVisible" :url="url"></ExternalLinkDialog>
</template>
