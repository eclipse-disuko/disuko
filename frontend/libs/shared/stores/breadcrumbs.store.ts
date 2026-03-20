// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {defineStore} from 'pinia';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute} from 'vue-router';
import {InternalBreadcrumbItem} from 'vuetify/lib/components/VBreadcrumbs/VBreadcrumbs';

export const useBreadcrumbsStore = defineStore('breadcrumbStore', () => {
  const {t} = useI18n();
  const route = useRoute();
  const currentBreadcrumbs = ref<InternalBreadcrumbItem[]>([]);

  const isInAdminArea = computed<boolean>(() => route.path.includes('admin'));
  const dashboardCrumbs = computed(() => (isInAdminArea.value ? [dashboard, adminDashboard] : [dashboard]));

  const setCurrentBreadcrumbs = (breadcrumbs: InternalBreadcrumbItem[]) => {
    currentBreadcrumbs.value = breadcrumbs;
  };

  const dashboard = {title: t('BC_Dashboard'), href: '/dashboard/home'};
  const adminDashboard = {title: t('BC_ADMIN'), href: '/dashboard/admin'};
  const projectsCrumb = {title: t('BC_Projects'), href: '/dashboard/projects/'};

  return {
    currentBreadcrumbs,
    isInAdminArea,
    setCurrentBreadcrumbs,
    dashboard,
    adminDashboard,
    projectsCrumb,
    dashboardCrumbs,
  };
});
