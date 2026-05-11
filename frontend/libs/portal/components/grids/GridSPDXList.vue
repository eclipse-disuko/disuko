<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script lang="ts">
import {useApprovalCheck} from '@disclosure-portal/composables/useApprovalCheck';
import Icons from '@disclosure-portal/constants/icons';
import {ProjectApprovable} from '@disclosure-portal/model/Approval';
import {ApprovableDto} from '@disclosure-portal/model/Project';
import {VersionSlim} from '@disclosure-portal/model/VersionDetails';
import {formatDateAndTime, getCssClassForTableRow, sbomOutdated} from '@disclosure-portal/utils/Table';
import {createSBOMURL, createVersionURL} from '@shared/utils/apiUrls';
import {DataTableHeader, DataTableItem} from '@shared/types/table';
import {PropType, computed, defineComponent, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {openUrlInNewTab} from '@shared/utils/url';

export default defineComponent({
  name: 'GridSPDXList',
  props: {
    projects: {
      type: Array as PropType<ProjectApprovable[]>,
      required: true,
    },

    showSbomExtras: {
      type: Boolean,
      default: false,
    },
    selectable: {
      type: Boolean,
      default: false,
    },
    doFilter: {
      type: Boolean,
      default: false,
    },
    filterIsFOSS: {
      type: Boolean,
      default: false,
    },
    channels: {
      type: [Map, Object] as PropType<Map<string, VersionSlim> | Record<string, VersionSlim>>,
      default: () => new Map(),
    },
  },
  emits: ['update:selectedProjects'],
  setup(props, {emit}) {
    const icons = Icons;
    const {t} = useI18n();
    const {isAudited} = useApprovalCheck();

    const selectedItems = ref<ProjectApprovable[]>(props.selectable ? [...props.projects] : []);

    if (props.selectable && props.projects.length > 0) {
      const selectedKeys = props.projects.map((item) => item.projectKey);
      emit('update:selectedProjects', selectedKeys);
    }

    const headers: DataTableHeader[] = [
      {
        title: t('COL_APPROVABLE_SPDX'),
        value: 'spdxname',
        align: 'start',
        width: 420,
      },
      {
        title: t('COL_STATS'),
        value: 'stats',
        width: 250,
      },
      {title: '', key: 'data-table-group', align: 'start'},
    ];

    const onRowClick = (event: Event, item: DataTableItem<ApprovableDto>) => {
      openSpdx(item.item);
    };

    const openSpdx = (approvable: ApprovableDto) => {
      if (!approvable.approvablespdx || !approvable.approvablespdx.spdxkey) {
        return;
      }
      const targetUrl = createSBOMURL(
        approvable.projectKey,
        approvable.approvablespdx.versionkey,
        approvable.approvablespdx.spdxkey,
      );
      openUrlInNewTab(targetUrl);
    };
    const openProject = (key: string) => {
      openUrlInNewTab(`/dashboard/projects/${encodeURIComponent(key)}`);
    };
    const openVersion = (approvable: ApprovableDto) => {
      if (!approvable.approvablespdx.versionkey) {
        return openProject(approvable.projectKey);
      }
      const targetUrl = createVersionURL(approvable.projectKey, approvable.approvablespdx.versionkey);
      openUrlInNewTab(targetUrl);
    };
    const groupBy = () => {
      const s = [];
      s.push({key: 'projectKey'});
      return s;
    };

    const isApproved = (approvable: ProjectApprovable) => {
      if (!approvable?.approvablespdx?.spdxkey || !approvable?.approvablespdx?.versionkey) {
        return false;
      }

      const channel =
        props.channels instanceof Map
          ? props.channels.get(approvable.approvablespdx.versionkey)
          : props.channels[approvable.approvablespdx.versionkey];

      if (!channel) {
        return false;
      }

      return isAudited(channel, approvable.approvablespdx.spdxkey);
    };

    const handleSelectionChange = (selected: unknown) => {
      if (!Array.isArray(selected)) return;
      const selectedKeys = (selected as ProjectApprovable[]).map((item) => item.projectKey);
      emit('update:selectedProjects', selectedKeys);
    };

    const filteredList = computed(() => {
      if (props.doFilter) {
        return props.projects.filter((p) => {
          if (props.filterIsFOSS) {
            return !p.isNonFoss;
          } else {
            return p.isNonFoss;
          }
        });
      }
      return props.projects;
    });

    return {
      headers,
      icons,
      openSpdx,
      openProject,
      filteredList,
      onRowClick,
      openVersion,
      getCssClassForTableRow,
      formatDateAndTime,
      sbomOutdated,
      t,
      groupBy,
      isApproved,
      handleSelectionChange,
      selectedItems,
    };
  },
});
</script>

<template>
  <v-data-table
    density="compact"
    fixed-header
    sort-desc
    class="striped-table"
    hide-default-footer
    :headers="headers"
    :items="projects"
    :items-per-page="-1"
    :item-class="getCssClassForTableRow"
    :group-by="groupBy()"
    :show-select="selectable"
    :select-strategy="selectable ? 'page' : undefined"
    item-value="projectKey"
    return-object
    v-model="selectedItems"
    @click:row="onRowClick"
    @update:model-value="handleSelectionChange"
    v-if="projects">
    <template v-slot:group-header="{item, isGroupOpen, toggleGroup}">
      <template
        :ref="
          (_el: any) => {
            if (!isGroupOpen(item)) toggleGroup(item);
          }
        "></template>
      <th colspan="3" class="text-caption expand-header p-1 px-3 text-start">
        <span @click="openProject(item.items[0].raw.projectKey)" class="cursor-pointer">
          <span class="font-color-table">{{ t('PROJECT') }}:</span>
          {{ item.items[0].raw.projectName }}
        </span>
        <v-chip
          size="x-small"
          class="ml-2"
          color="warning"
          variant="outlined"
          selected-class="blue"
          label
          v-if="item.items[0].raw.isNonFoss">
          <span class="font-weight-bold text-uppercase">{{ t('BADGE_NO_FOSS') }}</span>
        </v-chip>

        <span
          v-if="item.items[0].raw.approvablespdx.versionName"
          class="ml-4 cursor-pointer"
          @click="openVersion(item.items[0].raw)">
          <span class="font-color-table">&nbsp;{{ t('VERSION') }}:</span>
          {{ item.items[0].raw.approvablespdx.versionName }}
        </span>
      </th>
    </template>
    <template v-slot:[`item.spdxname`]="{item}">
      <span v-if="item.spdxname == ''">{{ t('NO_APPROVABLE_SPDX') }}</span>
      <v-row class="align-center pl-2" v-else>
        <v-col cols="auto" class="pa-0">
          <v-icon v-if="showSbomExtras" color="primary" size="small" class="pb-1">mdi-star</v-icon>
        </v-col>
        <v-col cols="auto" class="pa-0">
          <v-icon v-if="isApproved(item)" color="green" size="small" class="ml-1 pb-1">
            mdi-clipboard-check-outline
          </v-icon>
        </v-col>
        <v-col>
          <span v-if="item.spdxUploaded">{{ formatDateAndTime(item.spdxUploaded) }}&nbsp;-&nbsp;</span>
          <span>{{ item.spdxname }}</span>
          <span v-if="item.spdxtag">&nbsp;({{ item.spdxtag }})</span>
          <span v-if="showSbomExtras">
            <span v-if="item.isSpdxRecent">&nbsp;{{ '[' + t('SBOM_LATEST') + ']' }}</span>
            <span v-else>&nbsp;{{ '[' + t('SBOM_FORMER') + ']' }}</span>
          </span>
          <span v-if="showSbomExtras && sbomOutdated(item.spdxUploaded)" class="d-text d-secondary-text inline-block">
            <v-icon class="-mt-0.5 mr-1" color="red" x-small>mdi-priority-high</v-icon>
            <span>{{ t('SBOM_IS_OUTDATED') }}</span>
          </span>
        </v-col>
      </v-row>
    </template>
    <template v-slot:[`item.stats`]="{item}">
      <div v-if="item.spdxname != ''">
        <div class="flex flex-row justify-between">
          <div class="flex flex-col items-center justify-center">
            <v-icon size="small">mdi-layers</v-icon>
            <div class="text-no-wrap pt-1 text-center">{{ item.stats.Total }}</div>
          </div>
          <div class="flex flex-col items-center justify-center">
            <v-icon size="small" :color="item.stats.Denied > 0 ? 'policyStatusDeniedColor' : ''">
              mdi-minus-circle
            </v-icon>
            <div class="pt-1 text-center">{{ item.stats.Denied }}</div>
          </div>
          <div class="flex flex-col items-center justify-center">
            <v-icon size="small" :color="item.stats.NoAssertion > 0 ? 'policyStatusUnassertedColor' : ''">
              mdi-lightning-bolt-circle
            </v-icon>
            <div class="pt-1 text-center">{{ item.stats.NoAssertion }}</div>
          </div>
          <div class="flex flex-col items-center justify-center">
            <v-icon size="small" :color="item.stats.Warned > 0 ? 'policyStatusWarnedColor' : ''">mdi-alert</v-icon>
            <div class="pt-1 text-center">{{ item.stats.Warned }}</div>
          </div>
          <div class="flex flex-col items-center justify-center">
            <v-icon size="small" :color="item.stats.Questioned > 0 ? 'green' : ''">mdi-help-circle</v-icon>
            <div class="pt-1 text-center">{{ item.stats.Questioned }}</div>
          </div>
          <div class="flex flex-col items-center justify-center">
            <v-icon size="small" :color="item.stats.Allowed > 0 ? 'green' : ''">mdi-check-circle</v-icon>
            <div class="text-no-wrap pt-1 text-center">{{ item.stats.Allowed }}</div>
          </div>
        </div>
      </div>
    </template>
  </v-data-table>
</template>
