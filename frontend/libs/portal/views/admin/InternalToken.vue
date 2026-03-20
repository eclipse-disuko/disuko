<script lang="ts" setup>
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const {dashboardCrumbs, ...breadcrumbs} = useBreadcrumbsStore();

breadcrumbs.setCurrentBreadcrumbs([
  ...dashboardCrumbs,
  {
    title: t('INTERNAL_TOKEN'),
    href: '/dashboard/internaltoken',
  },
]);
</script>

<template>
  <GridInternalToken />
</template>
