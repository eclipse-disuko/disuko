<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ProjectApprovable} from '@disclosure-portal/model/Approval';
import {ComponentStats, VersionSlim} from '@disclosure-portal/model/VersionDetails';
import {useI18n} from 'vue-i18n';

interface Props {
  stats: ComponentStats;
  showRedWarnDeniedDecisionsMessage: boolean;
  projects: ProjectApprovable[];
  channels: Map<string, VersionSlim> | Record<string, VersionSlim>;
  isGroup: boolean;
  noFOSS: boolean;
  fossVersion: 'default' | 'legacy';
  selectedProjects: string[];
  doFilter?: boolean;
  tab: 'general' | 'approvable';
}

const props = withDefaults(defineProps<Props>(), {
  doFilter: false,
});

const emit = defineEmits<{
  'update:selectedProjects': [value: string[]];
  'update:tab': [value: 'general' | 'approvable'];
}>();

const {t} = useI18n();
</script>

<template>
  <v-tabs
    :model-value="props.tab"
    slider-color="mbti"
    show-arrows
    bg-color="tabsHeader"
    @update:model-value="emit('update:tab', $event)">
    <v-tab value="general">{{ t('TAB_TITLE_GENERAL') }}</v-tab>
    <v-tab value="approvable" v-if="props.isGroup">{{ t('TAB_TITLE_DETAILS') }}</v-tab>
  </v-tabs>
  <v-tabs-window :model-value="props.tab">
    <v-tabs-window-item value="general">
      <DApprovalComponents
        :stats="props.stats"
        :showRedWarnDeniedDecisionsMessage="props.showRedWarnDeniedDecisionsMessage" />
    </v-tabs-window-item>
    <v-tabs-window-item value="approvable" eager v-if="props.isGroup">
      <GridSPDXList
        :projects="props.projects"
        :channels="props.channels"
        :do-filter="props.doFilter"
        :filter-is-f-o-s-s="!props.noFOSS"
        :foss-version="props.fossVersion"
        :selected-projects="props.selectedProjects"
        showSbomExtras
        selectable
        @update:selectedProjects="emit('update:selectedProjects', $event)" />
    </v-tabs-window-item>
  </v-tabs-window>
</template>
