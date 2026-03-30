// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script lang="ts" setup>
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const {dashboardCrumbs, ...breadcrumbs} = useBreadcrumbsStore();

breadcrumbs.setCurrentBreadcrumbs([
  ...dashboardCrumbs,
  {
    title: t('BC_CUSTOMIDS'),
    href: '/dashboard/customids/',
  },
]);
</script>

<template>
  <GridCustomIds></GridCustomIds>
</template>
