<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<template>
  <Stack class="h-full p-6">
    <h1 class="text-h5 pb-3 ga-2">{{ t('ANALYTICS') }}</h1>
    <v-card>
      <v-tabs v-model="selectedTab" slider-color="mbti" active-class="active" show-arrows bg-color="tabsHeader">
        <v-tab value="overview" :to="tabUrl.overview">
          {{ t('TAB_ANALYTICS') }}
        </v-tab>
        <v-tab value="occurrences" :to="tabUrl.occurrences">
          {{ t('TAB_OCCURRENCES') }}
        </v-tab>
        <v-tab value="stats" :to="tabUrl.stats">
          {{ t('TAB_STATS') }}
        </v-tab>
      </v-tabs>
      <v-tabs-window v-model="selectedTab">
        <v-tabs-window-item value="overview">
          <GridAnalytics ref="overview" />
        </v-tabs-window-item>
        <v-tabs-window-item value="occurrences">
          <GridOccurrences ref="occurrences" />
        </v-tabs-window-item>
        <v-tabs-window-item value="stats">
          <Stats ref="stats" />
        </v-tabs-window-item>
      </v-tabs-window>
    </v-card>
  </Stack>
</template>

<script lang="ts" setup>
import {useTabsWindows} from '@shared/composables/useTabsWindows';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {onMounted} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const {tabUrl, selectedTab} = useTabsWindows('/dashboard/analytics', ['overview', 'occurrences', 'stats']);
const breadcrumbs = useBreadcrumbsStore();

onMounted(() => {
  breadcrumbs.setCurrentBreadcrumbs([
    {
      title: t('BC_Dashboard'),
      href: '/dashboard/home',
    },
    {
      title: t('BC_Analytics'),
      href: '/dashboard/analytics/overview',
    },
  ]);
});
</script>
