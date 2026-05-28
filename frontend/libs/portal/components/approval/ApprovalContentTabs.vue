<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import DApprovalComponents from '@disclosure-portal/components/disco/DApprovalComponents.vue';
import GridSPDXList from '@disclosure-portal/components/grids/GridSPDXList.vue';
import {ProjectApprovable} from '@disclosure-portal/model/Approval';
import {ComponentStats, VersionSlim} from '@disclosure-portal/model/VersionDetails';
import {useI18n} from 'vue-i18n';

withDefaults(
  defineProps<{
    stats: ComponentStats;
    showRedWarnDeniedDecisionsMessage: boolean;
    projects: ProjectApprovable[];
    channels: Record<string, VersionSlim> | Map<string, VersionSlim>;
    isGroup: boolean;
    noFOSS: boolean;
    fossVersion: 'default' | 'legacy';
    doFilter?: boolean;
    selectable?: boolean;
  }>(),
  {
    doFilter: false,
    selectable: false,
  },
);

const tab = defineModel<string | number>('tab');
const selectedProjects = defineModel<string[]>('selectedProjects');

const {t} = useI18n();
</script>

<template>
  <v-tabs v-model="tab" slider-color="brand" show-arrows bg-color="tabsHeader">
    <v-tab value="general">{{ t('TAB_TITLE_GENERAL') }}</v-tab>
    <v-tab value="approvable" v-if="isGroup">{{ t('TAB_TITLE_DETAILS') }}</v-tab>
  </v-tabs>
  <v-tabs-window v-model="tab">
    <v-tabs-window-item value="general">
      <DApprovalComponents :stats="stats" :showRedWarnDeniedDecisionsMessage="showRedWarnDeniedDecisionsMessage" />
    </v-tabs-window-item>
    <v-tabs-window-item eager value="approvable" v-if="isGroup">
      <GridSPDXList
        :projects="projects"
        :channels="channels"
        showSbomExtras
        :selectable="selectable"
        :do-filter="doFilter"
        :filter-is-f-o-s-s="noFOSS"
        @update:selectedProjects="selectedProjects = $event" />
    </v-tabs-window-item>
  </v-tabs-window>
</template>
