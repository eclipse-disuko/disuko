<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import Markdown from '@shared/components/Markdown.vue';
import {useI18n} from 'vue-i18n';

defineProps<{
  license: {
    meta?: {
      evaluation?: string;
    };
  };
}>();

const {t} = useI18n();
</script>

<template>
  <TableLayout has-title has-tab>
    <template #description>
      <p class="text-caption">{{ t('LICENSE_CHARTS_DISCLAIMER') }}</p>
    </template>
    <template #table>
      <div class="text-justify [&_*]:text-[14px]">
        <Markdown v-if="license && license.meta && license.meta.evaluation" :text="license.meta?.evaluation"></Markdown>
      </div>
    </template>
  </TableLayout>
</template>
