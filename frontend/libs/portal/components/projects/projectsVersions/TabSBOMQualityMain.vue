// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script lang="ts" setup>
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import Stack from '@shared/layouts/Stack.vue';
import TableLayout from '@shared/layouts/TableLayout.vue';
import {computed, defineAsyncComponent, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute, useRouter} from 'vue-router';

type QualityTab = {
  id: string;
  buttonIcon: string;
  buttonText: string;
  expandText: string;
};

const componentMap: Record<string, ReturnType<typeof defineAsyncComponent>> = {
  scanRemarks: defineAsyncComponent(
    () => import('@disclosure-portal/components/projects/projectsVersions/sbom-quality/TabScanRemarks.vue'),
  ),
  licenseRemarks: defineAsyncComponent(
    () => import('@disclosure-portal/components/projects/projectsVersions/sbom-quality/TabLicenseRemarks.vue'),
  ),
  reviewRemarks: defineAsyncComponent(
    () => import('@disclosure-portal/components/projects/projectsVersions/sbom-quality/TabReviewRemarks.vue'),
  ),
  generalRemarks: defineAsyncComponent(
    () => import('@disclosure-portal/components/projects/projectsVersions/sbom-quality/TabGeneralRemarks.vue'),
  ),
};

const route = useRoute();
const router = useRouter();
const sbomStore = useSbomStore();
const {t} = useI18n();

const tabs = ref<QualityTab[]>([
  {
    id: 'scanRemarks',
    buttonIcon: 'mdi-text-search-variant',
    buttonText: 'TAB_SCAN_REMARKS',
    expandText: 'QT_INTRO_TEXT_SCAN_REMARKS',
  },
  {
    id: 'licenseRemarks',
    buttonIcon: 'mdi-gavel',
    buttonText: 'TAB_LICENSE_REMARKS',
    expandText: 'QT_INTRO_TEXT_SCAN_REMARKS',
  },
  {
    id: 'reviewRemarks',
    buttonIcon: 'mdi-message-draw',
    buttonText: 'TAB_REVIEW_REMARKS',
    expandText: 'QT_INTRO_TEXT_REVIEW_REMARKS',
  },
  {
    id: 'generalRemarks',
    buttonIcon: 'mdi-bank-outline',
    buttonText: 'TAB_GENERAL_REMARKS',
    expandText: '',
  },
]);
const selectedTabId = ref<string>('scanRemarks');

const selectedTab = computed(() => {
  return tabs.value.find((tab) => tab.id === selectedTabId.value);
});

const currentComponent = computed(() => {
  return selectedTab.value ? componentMap[selectedTab.value.id] : null;
});

const reload = () => {
  const routeName = route.name?.toString();
  if (routeName && routeName !== 'SbomQuality') {
    selectedTabId.value = routeName;
  } else {
    selectedTabId.value = 'scanRemarks';
  }

  if (!tabs.value.some((tab) => tab.id === selectedTabId.value)) {
    selectedTabId.value = tabs.value[0].id;
  }
};

const version = sbomStore.getCurrentVersion;

const changeTab = (currentTab: QualityTab, query = '') => {
  const tabName = 'sbomQuality';

  const sbomKey = route.params.currentSbom as string;
  const projectUuid = encodeURIComponent(route.params.uuid as string);
  const versionKey = encodeURIComponent(version._key);

  let url: string;
  if (sbomKey) {
    url = `/dashboard/projects/${projectUuid}/versions/${versionKey}/${tabName}/${encodeURIComponent(sbomKey)}/${currentTab.id}${query}`;
  } else {
    url = `/dashboard/projects/${projectUuid}/versions/${versionKey}/${tabName}/${currentTab.id}${query}`;
  }
  router.push(url);
};

watch(
  () => route.name,
  async () => {
    reload();
  },
  {immediate: true},
);
</script>
<template>
  <TableLayout has-tab has-title gap="0">
    <template #buttons>
      <Stack direction="row" id="top-quality-remarks">
        <v-btn
          v-for="tab in tabs"
          size="small"
          :key="tab.id"
          @click="changeTab(tab)"
          :variant="selectedTabId === tab.id ? 'tonal' : 'text'"
          :class="{active: selectedTabId === tab.id}"
          class="text-none card-border"
          min-width="130px">
          <v-icon color="primary" class="pr-2">{{ tab.buttonIcon }}</v-icon>
          {{ t(tab.buttonText) }}
        </v-btn>
      </Stack>
      <v-spacer></v-spacer>
    </template>
    <template #table>
      <div class="h-full mt-3">
        <component v-if="currentComponent" :is="currentComponent" />
      </div>
    </template>
  </TableLayout>
</template>
