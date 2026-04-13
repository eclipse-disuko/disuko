<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ProjectSlim} from '@disclosure-portal/model/ProjectsResponse';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import type {SearchOptions} from '@disclosure-portal/utils/Table';
import {openProjectUrlByKey} from '@disclosure-portal/utils/url';
import TableLayout from '@shared/layouts/TableLayout.vue';
import {DataTableHeader, DataTableItem, SortItem} from '@shared/types/table';
import {storeToRefs} from 'pinia';
import {computed, onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';
import {debounce} from 'lodash';
import {useProjectUtils} from '@disclosure-portal/utils/projects';

const {t} = useI18n();
const router = useRouter();
const projectsUtils = useProjectUtils();
const projectStore = useProjectStore();
const {projects, projectsCount, loading, projectPossibleStatuses} = storeToRefs(projectStore);

const abort = ref<AbortController | null>(null);
const search = ref('');
const selectedFilterStatus = ref<string[]>([]);
const sortBy = ref<SortItem[]>([{key: 'updated', order: 'desc'}]);
const itemsPerPage = ref(100);

const options = computed(
  (): SearchOptions => ({
    page: 1,
    itemsPerPage: itemsPerPage.value,
    sortBy: sortBy.value,
    groupBy: [],
    search: search.value,
    filterString: search.value,
    filterBy: {
      status: selectedFilterStatus.value,
    },
  }),
);

const headers = computed<DataTableHeader[]>(() => [
  {title: '', value: 'data-table-expand', width: '38'},
  {title: t('COL_ACTIONS'), align: 'center', width: 80, value: 'actions', sortable: false},
  {title: t('COL_STATUS'), sortable: true, value: 'status', width: '155'},
  {title: t('COL_GROUP'), align: 'center', sortable: true, value: 'isGroup', width: '120'},
  {title: t('COL_NAME'), align: 'start', value: 'name', width: 270, sortable: true},
  {
    title: t('COL_DEVELOPER_COMPANY'),
    align: 'start',
    width: 270,
    value: 'supplier',
    sortable: true,
  },
  {title: t('COL_OWNER_COMPANY'), align: 'start', width: 270, value: 'company', sortable: true},
  {
    title: t('COL_OWNER_DEPARTMENT'),
    align: 'start',
    width: 270,
    value: 'department',
    sortable: true,
  },
  {title: t('COL_APPID'), align: 'start', width: 155, value: 'applicationId', sortable: true},
  {title: t('COL_UPDATED'), align: 'start', width: 103, value: 'updated', sortable: true},
  {title: t('COL_CREATED'), align: 'start', width: 103, value: 'created', sortable: true},
]);

const reload = async () => {
  if (abort.value) {
    abort.value.abort();
  }

  abort.value = new AbortController();

  await projectStore.fetchProjects(options.value, abort.value.signal);

  abort.value = null;
};

const searchChanged = async () => {
  if (search.value && search.value.length > 80) {
    return;
  }

  options.value.page = 1;

  await reload();
};

const debounceReload = debounce(searchChanged, 300);

const onRowClick = (event: Event, item: DataTableItem<ProjectSlim>) => {
  const project: ProjectSlim = item.item;
  openProjectUrlByKey(project._key, router);
};

const expanded = ref<string[]>([]);
const toggleExpand = (item: ProjectSlim) => {
  const index = expanded.value.indexOf(item._key);
  if (index > -1) {
    expanded.value.splice(index, 1);
  } else {
    expanded.value.push(item._key);
  }
};

const isExpanded = (item: ProjectSlim) => {
  return expanded.value.includes(item._key);
};

watch(options, debounceReload, {deep: true});

onMounted(() => {
  reload();
});
</script>

<template>
  <TableLayout data-testid="projects">
    <template #buttons>
      <h1 class="text-h5">{{ t('AllProjects') }}</h1>
      <v-spacer></v-spacer>
      <v-text-field
        v-model="search"
        autocomplete="off"
        :max-width="500"
        append-inner-icon="mdi-magnify"
        variant="outlined"
        density="compact"
        :label="t('labelSearch')"
        single-line
        hide-details
        clearable />
    </template>
    <template #table>
      <div class="table-wrapper fill-height">
        <v-data-table-server
          :loading="loading"
          density="compact"
          class="striped-table custom-data-table fill-height"
          :headers="headers"
          :items="projects"
          fixed-header
          item-value="_key"
          :items-length="projectsCount"
          :row-props="{
            class: {
              'py-8': true,
            },
          }"
          :footer-props="{'items-per-page-options': [10, 50, 100, -1]}"
          :options="options"
          v-model:items-per-page="itemsPerPage"
          v-model:sort-by="sortBy"
          v-model:expanded="expanded"
          @click:row="onRowClick">
          <template #[`header.status`]="{column, getSortIcon, toggleSort}">
            <GridFilterHeader :column="column" :getSortIcon="getSortIcon" :toggleSort="toggleSort">
              <template #filter>
                <GridHeaderFilterIcon
                  v-model="selectedFilterStatus"
                  :column="column"
                  :label="t('COL_PROJECT_STATUS')"
                  :initial-selected="['active', 'ready']"
                  :allItems="projectPossibleStatuses">
                </GridHeaderFilterIcon>
              </template>
            </GridFilterHeader>
          </template>
          <template #[`item.updated`]="{item}">
            <DDateCellWithTooltip :value="item.updated"></DDateCellWithTooltip>
          </template>
          <template #[`item.created`]="{item}">
            <DDateCellWithTooltip :value="item.created"></DDateCellWithTooltip>
          </template>
          <template #[`item.status`]="{item}">
            <span :style="{color: projectsUtils.getTextStatusColor(item.status)}">
              {{ t('STATUS_' + (!item.status ? 'new' : item.status)) }}
            </span>
          </template>
          <template #[`item.isGroup`]="{item}">
            <v-icon icon="mdi-check" class="mr-2" :color="item.isGroup ? 'primary' : 'tableBorderColor'" />
          </template>
          <template #[`item.actions`]="{item}">
            <ProjectsTableAction :item="item" @reload="reload()"></ProjectsTableAction>
          </template>
          <template #[`item.company`]="{item}">
            <span v-if="!item.missing">{{ item.company }}</span>
            <div v-else>
              <v-icon class="pr-2" icon="mdi-alert" color="warning" small></v-icon>
              <span>{{ t('WARNING_MISSING_DEPT') }}</span>
            </div>
          </template>
          <template #[`item.department`]="{item}">
            <span v-if="!item.missing">{{ item.department }}</span>
            <div v-else>
              <v-icon class="pr-2" color="warning" icon="mdi-alert" small></v-icon>
              <span>{{ t('WARNING_MISSING_DEPT') }}</span>
            </div>
          </template>
          <template #[`item.supplier`]="{item}">
            <span v-if="!item.supplierMissing">{{ item.supplier }}</span>
            <div v-else>
              <v-icon class="pr-2" color="warning" icon="mdi-alert" small></v-icon>
              <span>{{ t('WARNING_MISSING_DEPT') }}</span>
            </div>
          </template>
          <template #[`item.data-table-expand`]="{item}">
            <v-icon color="primary" @click.stop="toggleExpand(item)">
              {{ isExpanded(item) ? 'mdi-chevron-up' : 'mdi-chevron-down' }}
            </v-icon>
          </template>

          <template #expanded-row="{item}">
            <td :colspan="headers.length" class="cursor-default h-full overflow-y-clip bg-table-header">
              <GridProjectsExpandContent :item="item" :is-async="true"></GridProjectsExpandContent>
            </td>
          </template>
        </v-data-table-server>
      </div>
    </template>
  </TableLayout>
</template>

<style scoped lang="scss">
.bg-table-header {
  @apply bg-[rgb(var(--v-theme-tableHeaderBackgroundColor))];
}
</style>
