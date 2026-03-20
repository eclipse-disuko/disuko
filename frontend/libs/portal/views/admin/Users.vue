<script setup lang="ts">
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {onBeforeMount} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const {dashboardCrumbs, ...breadcrumbs} = useBreadcrumbsStore();

onBeforeMount(() => {
  breadcrumbs.setCurrentBreadcrumbs([
    ...dashboardCrumbs,
    {
      title: t('TITLE_USERS'),
      disabled: false,
      href: '/dashboard/admin/users',
    },
  ]);
});
</script>

<template>
  <TableUser />
</template>
