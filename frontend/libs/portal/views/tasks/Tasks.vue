// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import GridTask from '@disclosure-portal/components/grids/GridTask.vue';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {onMounted} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const breadcrumbs = useBreadcrumbsStore();

const initBreadcrumbs = () => {
  breadcrumbs.setCurrentBreadcrumbs([
    {
      title: t('BC_Dashboard'),
      disabled: false,
      href: '/dashboard/home',
    },
    {
      title: t('BC_Tasks'),
      disabled: false,
      href: '/dashboard/tasks',
    },
  ]);
};

onMounted(() => {
  initBreadcrumbs();
});
</script>
<template>
  <GridTask :in-own-view="true" :fixed-height="false" />
</template>
