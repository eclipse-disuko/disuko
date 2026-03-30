// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {Rights} from '@disclosure-portal/model/Rights';
import {useUserStore} from '@disclosure-portal/stores/user';
import {RightsUtils} from '@disclosure-portal/utils/Rights';
import {useTabsWindows} from '@shared/composables/useTabsWindows';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import config from '@shared/utils/config';
import {computed, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const getTABURL = (tabName: string) => {
  return '/dashboard/admin/tools/' + tabName;
};
const rights = ref(new Rights());
const userStore = useUserStore();
const {t} = useI18n();
const {dashboardCrumbs, ...breadcrumbs} = useBreadcrumbsStore();

const {tabUrl, selectedTab} = useTabsWindows('/dashboard/admin/tools', [
  'analytics',
  'accessRights',
  'export_import',
  'storageConsistency',
  'sampleData',
  'termsOfUseManagement',
  'mail',
  'notificationBar',
]);

onMounted(() => {
  rights.value = userStore.getRights;

  RightsUtils.redirectRestrictedAccess((): boolean => {
    return rights.value.hasToolsAccess() || rights.value.hasSampleDataAccess();
  });
  rights.value = RightsUtils.rights();
  initBreadcrumbs();
});

const isProd = computed(() => {
  return config.isProd;
});

const initBreadcrumbs = () => {
  breadcrumbs.setCurrentBreadcrumbs([...dashboardCrumbs, {title: t('BC_Tools'), disabled: false}]);
};

defineExpose({
  isProd,
  rights,
  getTABURL,
  initBreadcrumbs,
  t,
});
</script>

<template>
  <v-container fluid>
    <v-row class="pl-2 pt-2 align-center measure-height">
      <v-col xs="12" md="auto" class="d-flex align-center">
        <h1 class="text-h5 pr-2">{{ t('TOOLS') }}</h1>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-tabs v-model="selectedTab" slider-color="mbti" active-class="active" show-arrows bg-color="tabsHeader">
            <v-tab v-if="rights.hasToolsAccess()" value="analytics" :to="tabUrl.analytics">
              {{ t('ANALYTICS') }}
            </v-tab>
            <v-tab v-if="rights.hasToolsAccess()" value="accessRights" :to="tabUrl.accessRights">
              {{ t('TAB_ADMIN_ACCESS_RIGHTS') }}
            </v-tab>
            <v-tab v-if="rights.hasToolsAccess()" value="export_import" :to="tabUrl.export_import">
              {{ t('TAB_ADMIN_EXPORT_IMPORT') }}
            </v-tab>
            <v-tab v-if="rights.hasToolsAccess()" value="storageConsistency" :to="tabUrl.storageConsistency">
              {{ t('TAB_ADMIN_S3') }}
            </v-tab>
            <v-tab v-if="!isProd && rights.hasSampleDataAccess()" value="sampleData" :to="tabUrl.sampleData">
              {{ t('TAB_ADMIN_SAMPLE_DATA') }}
            </v-tab>
            <v-tab v-if="rights.hasToolsAccess()" value="termsOfUseManagement" :to="tabUrl.termsOfUseManagement">
              {{ t('TAB_ADMIN_TERMS_MANAGEMENT') }}
            </v-tab>
            <v-tab v-if="rights.hasToolsAccess()" value="mail" :to="tabUrl.mail">
              {{ t('TAB_ADMIN_MAIL') }}
            </v-tab>
            <v-tab v-if="rights.hasToolsAccess()" value="notificationBar" :to="tabUrl.notificationBar">
              {{ t('TAB_ADMIN_NOTIFICATION_BAR') }}
            </v-tab>
          </v-tabs>

          <v-tabs-window v-model="selectedTab">
            <v-tabs-window-item v-if="rights.hasToolsAccess()" value="analytics">
              <Statistics ref="analytics"></Statistics>
            </v-tabs-window-item>
            <v-tabs-window-item v-if="rights.hasToolsAccess()" value="accessRights">
              <AccessRights ref="accessRights"></AccessRights>
            </v-tabs-window-item>
            <v-tabs-window-item v-if="rights.hasToolsAccess()" value="export_import">
              <ExportImportTools ref="export_import"></ExportImportTools>
            </v-tabs-window-item>
            <v-tabs-window-item v-if="rights.hasToolsAccess()" value="storageConsistency">
              <S3DataIntegrity ref="storageConsistency"></S3DataIntegrity>
            </v-tabs-window-item>
            <v-tabs-window-item v-if="!isProd && rights.hasSampleDataAccess()" value="sampleData">
              <SampleData ref="sampleData"></SampleData>
            </v-tabs-window-item>
            <v-tabs-window-item v-if="rights.hasToolsAccess()" value="termsOfUseManagement">
              <TermsOfUseManagement ref="termsOfUseManagement"></TermsOfUseManagement>
            </v-tabs-window-item>
            <v-tabs-window-item v-if="rights.hasToolsAccess()" value="mail">
              <Mail ref="mail"></Mail>
            </v-tabs-window-item>
            <v-tabs-window-item v-if="rights.hasToolsAccess()" value="notificationBar">
              <NotificationBarManagement ref="notificationBar"></NotificationBarManagement>
            </v-tabs-window-item>
          </v-tabs-window>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>
