<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script lang="ts" setup>
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import DialogLayout, {type DialogLayoutConfig} from '@shared/layouts/DialogLayout.vue';
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
const dialog = ref<'disclaimer' | 'info' | null>(null);
const selectedTabId = ref<string>('scanRemarks');

const selectedTab = computed(() => {
  return tabs.value.find((tab) => tab.id === selectedTabId.value);
});

const currentComponent = computed(() => {
  return selectedTab.value ? componentMap[selectedTab.value.id] : null;
});

const showExpansionPanel = computed(() => {
  return selectedTab.value && selectedTab.value.expandText && selectedTab.value.id !== 'generalRemarks';
});

const expansionPanelTitle = computed(() => {
  if (!selectedTab.value) return '';
  return `${t('HEADER_OBLIGATIONS_AND_CONFIDENTIALITY')} ${t(selectedTab.value.buttonText)}`;
});

const expansionPanelText = computed(() => {
  if (!selectedTab.value || !selectedTab.value.expandText) return '';
  return t(selectedTab.value.expandText);
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

const getDialogConfig = (): DialogLayoutConfig => ({
  title: dialog.value === 'info' ? expansionPanelTitle.value : t('BTN_DISCLAIMER'),
});

const dialogContent = computed(() => {
  return dialog.value === 'info' ? expansionPanelText.value : t('QT_SCAN_REMARKS_DISCLAIMER');
});

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
      <v-spacer></v-spacer>
      <v-btn variant="text" size="small" class="text-none" v-if="showExpansionPanel" @click="dialog = 'info'">
        <v-icon color="primary" icon="mdi-chevron-right" class="mr-2" />
        <span class="text-caption">{{ expansionPanelTitle }}</span>
      </v-btn>
    </template>
    <template #table>
      <div class="relative h-[calc(100%-12px)] pt-3">
        <component v-if="currentComponent" :is="currentComponent" />
        <div id="bottom-quality-remarks" class="absolute bottom-5 left-4">
          <v-btn variant="text" size="small" class="text-none" @click="dialog = 'disclaimer'">
            <v-icon color="primary">mdi mdi-chevron-right</v-icon>
            {{ t('BTN_DISCLAIMER') }}
          </v-btn>
        </div>
      </div>
    </template>
  </TableLayout>

  <v-dialog persistent :model-value="Boolean(dialog)" width="700" scrollable>
    <DialogLayout :config="getDialogConfig()" @close="dialog = null">
      <span class="text-caption" v-html="dialogContent"></span>
    </DialogLayout>
  </v-dialog>
</template>
