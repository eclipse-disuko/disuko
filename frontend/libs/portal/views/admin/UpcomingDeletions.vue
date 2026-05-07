<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script lang="ts" setup>
import GridUpcomingDeletions from '@disclosure-portal/components/grids/GridUpcomingDeletions.vue';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {onMounted} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const {dashboardCrumbs, ...breadcrumbs} = useBreadcrumbsStore();

onMounted(() => {
  breadcrumbs.setCurrentBreadcrumbs([...dashboardCrumbs, {title: t('UPCOMING_DELETIONS'), disabled: true}]);
});
</script>

<template>
  <GridUpcomingDeletions />
</template>
