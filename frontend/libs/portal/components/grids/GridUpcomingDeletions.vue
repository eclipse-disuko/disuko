<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {UpcomingDeletion} from '@disclosure-portal/model/UpcomingDeletion';
import {BlockingProject} from '@disclosure-portal/model/UpcomingDeletion';
import adminService from '@disclosure-portal/services/admin';
import Icons from '@disclosure-portal/constants/icons';
import {openUrlInNewTab} from '@shared/utils/url';
import {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import {DataTableHeader, SortItem} from '@shared/types/table';
import {computed, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useLabelStore} from '@disclosure-portal/stores/label.store';

const {t} = useI18n();

const labelStore = useLabelStore();
const icons = Icons;

const items = ref<UpcomingDeletion[]>([]);
const loaded = ref(false);
const expanded = ref<string[]>([]);
const sortBy: SortItem[] = [{key: 'overdue', order: 'desc'}];

const blockingProjectHeaders = computed((): DataTableHeader[] => [
  {
    title: '',
    sortable: false,
    align: 'center',
    width: 50,
    value: 'actions',
  },
  {
    title: t('COL_NAME'),
    align: 'start',
    value: 'name',
    sortable: true,
  },
  {
    title: t('COL_APPID'),
    align: 'start',
    value: 'applicationId',
    sortable: true,
  },
  {
    title: t('COL_LABELS'),
    align: 'start',
    value: 'labels',
    sortable: false,
  },
]);

const headers = computed((): DataTableHeader[] => [
  {
    title: '',
    value: 'data-table-expand',
    width: 25,
  },
  {
    title: t('COL_OVERDUE'),
    align: 'center',
    value: 'overdue',
    sortable: true,
    width: 100,
  },
  {
    title: t('COL_USER'),
    align: 'start',
    value: 'user',
    sortable: true,
  },
  {
    title: t('COL_FORENAME'),
    align: 'start',
    value: 'forename',
    sortable: true,
  },
  {
    title: t('COL_LASTNAME'),
    align: 'start',
    value: 'lastname',
    sortable: true,
  },
  {
    title: t('COL_DEPARTMENT'),
    align: 'start',
    value: 'department',
    sortable: true,
  },
  {
    title: t('COL_DEPROVISIONED'),
    align: 'start',
    value: 'deprovisioned',
    sortable: true,
  },
  {
    title: t('COL_DELETION_DATE'),
    align: 'start',
    value: 'deletionDate',
    sortable: true,
  },
]);

const onRowExpand = (newExpanded: string[]) => {
  if (newExpanded.length > 1) {
    expanded.value = [newExpanded[newExpanded.length - 1]];
  } else {
    expanded.value = newExpanded;
  }
};

const toggleExpand = (item: UpcomingDeletion) => {
  const index = expanded.value.indexOf(item.user);
  if (index > -1) {
    expanded.value.splice(index, 1);
  } else {
    expanded.value.push(item.user);
  }
};

const isExpanded = (item: UpcomingDeletion) => {
  return expanded.value.includes(item.user);
};

const blockingProjectActionButtons = computed((): TableActionButtonsProps['buttons'] => [
  {
    icon: 'mdi-open-in-new',
    hint: t('TT_OPEN_PROJECT'),
    event: 'open',
    show: true,
  },
]);

const openProject = (project: BlockingProject) => {
  openUrlInNewTab(`/dashboard/projects/${project.key}`);
};

onMounted(async () => {
  const res = await adminService.getUpcomingDeletions();
  items.value = res.data;
  loaded.value = true;
});
</script>

<template>
  <TableLayout>
    <template #buttons>
      <h1 class="text-h5">{{ t('UPCOMING_DELETIONS') }}</h1>
    </template>
    <template #table>
      <v-data-table
        density="compact"
        class="striped-table fill-height"
        :loading="!loaded"
        item-value="user"
        :items="items"
        :headers="headers"
        :items-per-page="50"
        fixed-header
        :sort-by="sortBy"
        expand-on-click
        :expanded.sync="expanded"
        @update:expanded="onRowExpand">
        <template #[`item.data-table-expand`]="{item}">
          <v-icon v-if="item.blockingProjects?.length" color="primary" @click.stop="toggleExpand(item)">
            {{ isExpanded(item) ? 'mdi-chevron-up' : 'mdi-chevron-down' }}
          </v-icon>
        </template>
        <template #[`item.department`]="{item}">
          <span v-if="item.departmentDescription && item.department"
            >{{ item.departmentDescription }} ({{ item.department }})</span
          >
          <span v-else>-</span>
        </template>
        <template #[`item.deprovisioned`]="{item}">
          <DDateCellWithTooltip :value="item.deprovisioned" />
        </template>
        <template #[`item.deletionDate`]="{item}">
          <DDateCellWithTooltip :value="item.deletionDate" />
        </template>
        <template #[`item.overdue`]="{item}">
          <v-icon v-if="item.overdue" color="error" icon="mdi-exclamation-thick" />
        </template>
        <template #expanded-row="{columns, item}">
          <tr>
            <td :colspan="columns.length" class="pa-4">
              <h3 class="text-subtitle-2 mb-2">{{ t('COL_BLOCKING_PROJECTS') }}</h3>
              <v-data-table
                density="compact"
                :items="item.blockingProjects"
                :headers="blockingProjectHeaders"
                item-key="key"
                hide-default-footer
                class="striped-table">
                <template #[`item.actions`]="{item: project}">
                  <TableActionButtons
                    variant="compact"
                    :buttons="blockingProjectActionButtons"
                    @open="openProject(project)" />
                </template>
                <template #[`item.labels`]="{item: project}">
                  <div class="flex flex-wrap gap-1 py-1">
                    <Tooltip
                      v-for="(l, i) in project.freeLabels"
                      :key="'free' + i"
                      :text="t('TT_free_label')"
                      as-parent>
                      <DLabel :labelName="l" :iconName="icons.TAG" />
                    </Tooltip>
                    <Tooltip
                      v-for="(l, i) in project.policyLabels"
                      :key="'policy' + i"
                      :text="`${t('TT_policy_label_with_description')}${labelStore.getLabelByKey(l).description}`"
                      as-parent>
                      <DLabel
                        :labelName="labelStore.getLabelByKey(l).name ?? 'UNKNOWN_LABEL'"
                        :iconName="icons.POLICY" />
                    </Tooltip>
                    <Tooltip
                      v-for="(l, i) in project.projectLabels"
                      :key="'proj' + i"
                      :text="`${t('TT_project_label_with_description')}${labelStore.getLabelByKey(l).description}`"
                      as-parent>
                      <DLabel
                        :labelName="labelStore.getLabelByKey(l).name ?? 'UNKNOWN_LABEL'"
                        :iconName="icons.PROJECT_LABEL" />
                    </Tooltip>
                  </div>
                </template>
              </v-data-table>
            </td>
          </tr>
        </template>
      </v-data-table>
    </template>
  </TableLayout>
</template>
