<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import Tooltip from '@shared/components/disco/Tooltip.vue';
import Stack from '@shared/layouts/Stack.vue';
import {ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';

const props = defineProps<{
  headerText: string;
  helpText: string;
  navigationPath?: string;
  total?: number;
}>();

const showMenu = ref(false);
const router = useRouter();
const {t} = useI18n();

const navigateTo = () => {
  if (props.navigationPath) {
    router.push(props.navigationPath);
  }
};
</script>

<template>
  <Stack
    direction="row"
    align="center"
    justify="between"
    class="min-h-9 border-b border-[rgba(var(--v-theme-on-surface),0.08)] pb-3">
    <button
      v-if="navigationPath"
      type="button"
      class="group hover:text-primary focus-visible:outline-primary inline-flex min-w-0 items-center gap-2 rounded text-left text-[15px] leading-snug font-semibold tracking-normal transition-colors duration-200 focus-visible:outline-2 focus-visible:outline-offset-4 motion-reduce:transition-none"
      @click="navigateTo()">
      <span>{{ t(headerText) }}</span>
      <v-icon
        icon="mdi-arrow-right"
        size="16"
        class="shrink-0 opacity-50 transition-[transform,opacity] duration-200 group-hover:translate-x-0.5 group-hover:opacity-100 group-focus-visible:opacity-100 motion-reduce:transform-none motion-reduce:transition-none" />
      <Tooltip :text="t('CHARTS_TT_DETAILS')" />
    </button>
    <h3 v-else class="text-[15px] leading-snug font-semibold tracking-normal">{{ t(headerText) }}</h3>

    <span
      v-if="total !== undefined"
      class="ml-auto shrink-0 text-[15px] font-semibold tabular-nums"
      :aria-label="`${t('TOTAL')}: ${total}`">
      {{ total }}
    </span>

    <v-menu v-model="showMenu" :close-on-content-click="false" location="bottom end">
      <template #activator="{props: menuProps}">
        <v-btn
          v-bind="menuProps"
          icon
          variant="text"
          size="small"
          density="comfortable"
          class="focus-visible:outline-primary shrink-0 opacity-60 transition-opacity duration-200 hover:opacity-100 focus-visible:opacity-100 focus-visible:outline-2 motion-reduce:transition-none"
          :aria-label="t('CHARTS_TT_HELP')">
          <v-icon icon="mdi-help-circle-outline" />
          <Tooltip :text="t('CHARTS_TT_HELP')" />
        </v-btn>
      </template>
      <v-card max-width="320">
        <v-card-text class="pa-3 helptext-container">
          {{ t(helpText) }}
        </v-card-text>
      </v-card>
    </v-menu>
  </Stack>
</template>
