// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<template>
  <v-container fluid class="h-full" data-testid="licenses-details">
    <div class="pb-3 d-flex justify-end">
      <span class="text-h5 inline-block">{{ t('LICENSE') }}</span>
      <span class="text-h5 px-2 inline-block">
        <q> {{ item.name }} </q>
      </span>

      <span class="d-flex align-center">
        <v-icon v-if="item.meta?.approvalState === 'forbidden'" color="licenceForbidden">
          mdi-shield-alert-outline
        </v-icon>
        <v-icon v-if="item.meta?.isLicenseChart" color="licenceChartIcon">mdi-shield-check-outline </v-icon>
        <v-icon v-else color="licenceNotApproved">mdi-shield-off-outline</v-icon>

        <v-tooltip activator="parent" content-class="dpTooltip" location="bottom">
          <span v-if="item.meta?.approvalState === 'forbidden'">{{ t('ICON_LICENSE_FORBIDDEN_TOOLTIP') }}</span>
          <span v-else-if="item.meta?.isLicenseChart">{{ t('ICON_LICENSE_CHART_STATUS_TOOLTIP') }}</span>
          <span v-else>{{ t('TABLE_LICENSE_CHART_STATUS_IS_NOT') }}</span>
        </v-tooltip>
      </span>
      <v-spacer></v-spacer>

      <NewOrEditLicenseDialog @closed:successfully="reload()" :initial-data="item" mode="edit">
        <template v-slot="{showDialog}">
          <DCActionButton
            v-if="hasEditRights"
            large
            :text="t('BTN_EDIT')"
            icon="mdi-pencil"
            :hint="t('TT_edit_license')"
            data-testid="licenses-edit-button"
            @click="showDialog" />
        </template>
      </NewOrEditLicenseDialog>
    </div>

    <v-card>
      <v-tabs
        v-model="selectedTab"
        data-testid="licenseDetailsTabs"
        slider-color="mbti"
        show-arrows
        bg-color="tabsHeader">
        <v-tab v-for="header in filteredTabHeaders" :key="header.key" :value="header.key" :to="header.key">
          <span> {{ t(header.text) }} </span>
          <v-tooltip activator="parent" content-class="dpTooltip" location="bottom">
            {{ t(header.tooltipText) }}
          </v-tooltip>
        </v-tab>
      </v-tabs>
      <v-tabs-window v-model="selectedTab">
        <v-tabs-window-item value="details">
          <LicenseDetails :license="item"></LicenseDetails>
        </v-tabs-window-item>
        <v-tabs-window-item value="license">
          <LicenseText :license="item"></LicenseText>
        </v-tabs-window-item>
        <v-tabs-window-item value="obligations">
          <LicenseObligations v-if="licenseLoaded" :license="item"></LicenseObligations>
        </v-tabs-window-item>
        <v-tabs-window-item value="legalcomments">
          <LicenseLegalComments :license="item" />
        </v-tabs-window-item>
        <v-tabs-window-item value="policyRules">
          <GridPolicyRulesAssignments
            v-if="licenseId"
            ref="policyRulesAssignments"
            :license-id="licenseId"
            :disabled="true"></GridPolicyRulesAssignments>
        </v-tabs-window-item>
        <v-tabs-window-item value="aliases">
          <GridAliases :license="item" mode="readonly"></GridAliases>
        </v-tabs-window-item>
        <v-tabs-window-item value="evaluation">
          <LicenseEvaluation :license="item"></LicenseEvaluation>
        </v-tabs-window-item>
        <v-tabs-window-item value="changelog" class="expanding-container" v-if="isLicenseAdmin">
          <div v-if="item.meta && item.meta.changelog">
            <div class="licenseText text-caption pa-3" v-text="item.meta.changelog"></div>
          </div>
        </v-tabs-window-item>
        <v-tabs-window-item value="auditLog" v-if="isLicenseAdmin">
          <GridAuditLog :fetch-method="() => LicenseService.getAuditTrail(licenseId)" ref="auditLog"></GridAuditLog>
        </v-tabs-window-item>
      </v-tabs-window>
    </v-card>
  </v-container>
</template>

<script setup lang="ts">
import {usePageTitle} from '@disclosure-portal/composables/usePageTitle';
import LicenseModel from '@disclosure-portal/model/License';
import LicenseService from '@disclosure-portal/services/license';
import {useUserStore} from '@disclosure-portal/stores/user';
import {RightsUtils} from '@disclosure-portal/utils/Rights';
import {useTabsWindows} from '@shared/composables/useTabsWindows';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {computed, onBeforeMount, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute} from 'vue-router';

const {t} = useI18n();
const route = useRoute();
const breadcrumbs = useBreadcrumbsStore();
const userStore = useUserStore();
const {useReactiveTitle} = usePageTitle();

const licenseLoaded = ref(false);
const item = ref<LicenseModel>({} as LicenseModel);

const licenseId = computed(() => (Array.isArray(route.params.id) ? route.params.id[0] : route.params.id));
const hasEditRights = computed(() => userStore.getRights.allowLicense.update);
const isLicenseAdmin = computed(() => RightsUtils.hasLicenseAccess());

const filteredTabHeaders = computed(() => {
  return tabHeaders.filter((header) => {
    if (header.requiresLicenseAdmin) {
      return isLicenseAdmin.value;
    }
    return true;
  });
});

const baseUrl = computed(() => `/dashboard/licenses/${encodeURIComponent(licenseId.value)}`);

const tabHeaders = [
  {
    key: 'details',
    text: 'TAB_OVERVIEW',
    tooltipText: 'TAB_TITLE_OVERVIEW_TOOLTIP',
  },
  {
    key: 'license',
    text: 'TAB_TITLE_LICENSE',
    tooltipText: 'TAB_TITLE_LICENSE_TOOLTIP',
  },
  {
    key: 'obligations',
    text: 'TAB_TITLE_CLASSIFICATIONS',
    tooltipText: 'TAB_TITLE_CLASSIFICATIONS_TOOLTIP',
  },
  {
    key: 'legalcomments',
    text: 'TAB_TITLE_LEGALCOMMENTS',
    tooltipText: 'TAB_TITLE_LEGALCOMMENTS_TOOLTIP',
  },
  {
    key: 'policyRules',
    text: 'TAB_TITLE_POLICY_RULES',
    tooltipText: 'TAB_TITLE_POLICY_RULES_TOOLTIP',
  },
  {
    key: 'aliases',
    text: 'TAB_TITLE_ALIASES',
    tooltipText: 'TAB_TITLE_ALIASES_TOOLTIP',
  },
  {
    key: 'evaluation',
    text: 'TAB_TITLE_EVALUATION',
    tooltipText: 'TAB_TITLE_EVALUATION_TOOLTIP',
  },
  {
    key: 'changelog',
    text: 'TAB_TITLE_CHANGELOG',
    tooltipText: 'TAB_TITLE_CHANGELOG_TOOLTIP',
    requiresLicenseAdmin: true,
  },
  {
    key: 'auditLog',
    text: 'TAB_PROJECT_AUDIT',
    tooltipText: 'TAB_PROJECT_AUDIT',
    requiresLicenseAdmin: true,
  },
];

const {selectedTab} = useTabsWindows(
  baseUrl,
  tabHeaders.map((h) => h.key),
);

const reload = async () => {
  licenseLoaded.value = false;
  const response = await LicenseService.get(licenseId.value);
  if (response) {
    licenseLoaded.value = true;
    item.value = response.data;
  }
};

const initBreadcrumbs = () => {
  breadcrumbs.setCurrentBreadcrumbs([
    {
      title: t('BC_Dashboard'),
      href: '/dashboard/home',
    },
    {
      title: t('BC_License'),
      href: '/dashboard/licenses/',
    },
    {
      title: item.value?.name,
      href: `/dashboard/licenses/${encodeURIComponent(item.value?.licenseId)}`,
    },
  ]);
};

watch(licenseId, async (newId) => {
  if (newId) {
    await reload();
  }
});

onBeforeMount(async () => {
  await reload();
  initBreadcrumbs();
});

// Set up reactive title based on license name
watch(
  () => item.value?.name,
  (item) => {
    if (item) {
      useReactiveTitle(item + ' | License');
    }
  },
  {immediate: true},
);
</script>
