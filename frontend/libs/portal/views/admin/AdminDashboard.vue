<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import ITile from '@disclosure-portal/model/ITile';
import adminService, {GetDashboardCounts} from '@disclosure-portal/services/admin';
import {RightsUtils} from '@disclosure-portal/utils/Rights';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {computed, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const {dashboardCrumbs, ...breadcrumbs} = useBreadcrumbsStore();

const counts = ref<GetDashboardCounts | undefined>();

const initBreadcrumbs = () => {
  breadcrumbs.setCurrentBreadcrumbs(dashboardCrumbs);
};

const tiles = computed<ITile[]>(() => {
  const res: ITile[] = [];
  if (RightsUtils.hasAllProjectsReadonly()) {
    res.push({
      color: 'primary',
      cnt: counts.value?.projectCount || -1,
      visible: !!counts.value,
      title: 'AllProjects',
      url: '/dashboard/admin/projects',
      icon: 'mdi-list-box-outline',
      expandGroup: false,
      expand: false,
    });
  }
  if (RightsUtils.hasLabelAccess()) {
    res.push({
      color: 'primary',
      cnt: counts.value?.labelCount || -1,
      visible: !!counts.value,
      title: 'Labels',
      url: '/dashboard/admin/labels/schema',
      icon: 'mdi-label-outline',
      expandGroup: false,
      expand: false,
    });
  }
  if (RightsUtils.hasReviewTemplatesAccess()) {
    res.push({
      color: 'primary',
      cnt: counts.value?.reviewTemplateCount || -1,
      visible: true,
      title: 'DB_TITLE_REVIEW_TEMPLATES',
      url: '/dashboard/admin/templates/review',
      icon: 'mdi-text-box-outline',
      expandGroup: false,
      expand: false,
    });
  }
  if (RightsUtils.hasSchemaAccess()) {
    res.push({
      color: 'primary',
      cnt: counts.value?.schemaCount || -1,
      visible: true,
      title: 'Schemes',
      url: '/dashboard/admin/schemas',
      icon: 'mdi-file-tree-outline',
      expandGroup: false,
      expand: false,
    });
  }
  if (RightsUtils.hasClassificationsAccess()) {
    res.push({
      color: 'primary',
      cnt: counts.value?.obligationCount || -1,
      visible: true,
      title: 'CLASSIFICATIONS',
      url: '/dashboard/admin/classifications',
      icon: 'mdi-grain',
      expandGroup: false,
      expand: false,
    });
  }
  if (RightsUtils.hasToolsAccess() || RightsUtils.hasSampleDataAccess()) {
    res.push({
      color: 'primary',
      cnt: -1,
      visible: true,
      title: 'Tools',
      url: '/dashboard/admin/tools/analytics',
      icon: 'mdi-hammer-screwdriver',
      expandGroup: false,
      expand: false,
    });
  }
  if (RightsUtils.isApplicationAdmin()) {
    res.push({
      color: 'primary',
      cnt: -1,
      visible: true,
      title: 'ADMIN_JOBS',
      url: '/dashboard/admin/jobs',
      icon: 'mdi-window-shutter-cog',
      expandGroup: false,
      expand: false,
    });
  }
  if (RightsUtils.hasUsersAccess()) {
    res.push({
      color: 'primary',
      cnt: counts.value?.userCount || -1,
      visible: true,
      title: 'USERS',
      url: '/dashboard/admin/users',
      icon: 'mdi-account-multiple-outline',
      expandGroup: false,
      expand: false,
    });
  }
  if (RightsUtils.isProjectAnalyst()) {
    res.push({
      color: 'primary',
      cnt: -1,
      visible: true,
      title: 'Analytics',
      url: '/dashboard/analytics/overview',
      icon: 'mdi-chart-box-outline',
      expandGroup: false,
      expand: false,
    });
  }
  if (RightsUtils.isApplicationAdmin()) {
    // TODO: ??
    res.push({
      color: 'primary',
      cnt: -1,
      visible: true,
      title: 'Custom IDs',
      url: '/dashboard/admin/customids',
      icon: 'mdi-identifier',
      expandGroup: false,
      expand: false,
    });
  }
  if (RightsUtils.isApplicationAdmin()) {
    res.push({
      color: 'primary',
      cnt: -1,
      visible: true,
      title: 'Internal Token',
      url: '/dashboard/admin/internaltoken',
      icon: 'mdi-shield-lock',
      expandGroup: false,
      expand: false,
    });
  }

  if (RightsUtils.isDomainAdmin() || RightsUtils.isFOSSOffice()) {
    res.push({
      color: 'primary',
      cnt: -1,
      visible: true,
      title: 'CHECKLISTS',
      url: '/dashboard/admin/checklist',
      icon: 'mdi-format-list-checks',
      expandGroup: false,
      expand: false,
    });
  }

  if (RightsUtils.isApplicationAdmin() || RightsUtils.isDomainAdmin()) {
    res.push({
      color: 'primary',
      cnt: -1,
      visible: true,
      title: 'NEWSBOX',
      url: '/dashboard/admin/newsbox',
      icon: 'mdi-newspaper-variant',
      expandGroup: false,
      expand: false,
    });
  }

  if (RightsUtils.isApplicationAdmin()) {
    res.push({
      color: 'primary',
      cnt: -1,
      visible: true,
      title: 'FEATURE_FLAGS',
      url: '/dashboard/admin/featureflags',
      icon: 'mdi-flag-variant',
      expandGroup: false,
      expand: false,
    });
    res.push({
      color: 'primary',
      cnt: -1,
      visible: true,
      title: 'Internationalization',
      url: '/dashboard/admin/i18n',
      icon: 'mdi-translate',
      expandGroup: false,
      expand: false,
    });
  }
  if (RightsUtils.isDomainAdmin()) {
    res.push({
      color: 'primary',
      cnt: -1,
      visible: true,
      title: 'TAB_ADMIN_USER_MANAGEMENT',
      url: '/dashboard/admin/userManagement',
      icon: 'mdi-account-cog-outline',
      expandGroup: false,
      expand: false,
    });
  }
  if (RightsUtils.isDomainAdmin()) {
    res.push({
      color: 'primary',
      cnt: -1,
      visible: true,
      title: 'UPCOMING_DELETIONS',
      url: '/dashboard/admin/deletions',
      icon: 'mdi-delete-clock-outline',
      expandGroup: false,
      expand: false,
    });
  }
  return res;
});

onMounted(() => {
  adminService.getDashboardCounts().then((res) => (counts.value = res));
  initBreadcrumbs();
});
</script>

<template>
  <v-row class="pa-4 m-0 pt-8">
    <v-col cols="12" xs="12" sm="6" md="4" lg="3" xl="2" v-for="(module, index) in tiles" :key="index" class="pa-4">
      <v-card :href="'#' + module.url" style="overflow: visible" class="card-border">
        <v-col cols="4" style="position: relative">
          <v-card
            width="80px"
            height="60px"
            class="d-flex align-center mt-n8 card-border justify-center pt-2"
            style="position: absolute"
            :href="'#' + module.url">
            <v-icon :icon="module.icon" size="x-large"></v-icon>
          </v-card>
        </v-col>
        <v-card-title class="font-weight-light text-right">
          {{ t(module.title) }}
        </v-card-title>
        <v-divider :thickness="1" class="border-opacity-100" color="brand" inset></v-divider>
        <v-card-text class="text-body-2 text-right">
          {{ module.cnt !== -1 ? module.cnt : '&#10240;' }}
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>
