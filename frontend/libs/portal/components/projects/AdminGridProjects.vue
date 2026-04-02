<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import ProjectsTableAction from '@disclosure-portal/components/projects/ProjectsTableAction.vue';
import {ProjectSlim} from '@disclosure-portal/model/ProjectsResponse';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import type {SearchOptions} from '@disclosure-portal/utils/Table';
import {openProjectUrlByKey} from '@disclosure-portal/utils/url';
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import DDateCellWithTooltip from '@shared/components/disco/DDateCellWithTooltip.vue';
import DIconButton from '@shared/components/disco/DIconButton.vue';
import TableLayout from '@shared/layouts/TableLayout.vue';
import {DataTableHeader, DataTableItem, SortItem} from '@shared/types/table';
import {useDebounceFn} from '@vueuse/core';
import {storeToRefs} from 'pinia';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';

const {t} = useI18n();
const router = useRouter();
const projectStore = useProjectStore();
const {projects, projectsCount, loading, projectPossibleStatuses} = storeToRefs(projectStore);

const searchInput = ref(''); // debounced input so that the table is not reloaded on every keypress
const search = ref('');
const selectedFilterStatus = ref<string[]>(['active', 'ready']);
const menu = ref(false);
const sortBy = ref<SortItem[]>([{key: 'updated', order: 'desc'}]);

const options = ref<SearchOptions>({
  page: 1,
  itemsPerPage: 50,
  sortBy: [{key: 'updated', order: 'desc'}],
  groupBy: [],
  search: '',
  filterString: '',
  filterBy: {
    status: [],
  },
});

const headers = computed<DataTableHeader[]>(() => [
  {title: '', class: 'tableHeaderCell', value: 'data-table-expand', width: '38'},
  {title: t('COL_ACTIONS'), align: 'center', width: 80, value: 'actions', sortable: false},
  {title: t('COL_STATUS'), sortable: true, filterable: true, value: 'status', width: '155'},
  {title: t('COL_GROUP'), align: 'center', sortable: true, filterable: false, value: 'isGroup', width: '120'},
  {title: t('COL_NAME'), align: 'start', value: 'name', width: 270, sortable: true},
  {title: t('COL_DEVELOPER_COMPANY'), align: 'start', width: 270, value: 'supplier', sortable: true},
  {title: t('COL_OWNER_COMPANY'), align: 'start', width: 270, value: 'company', sortable: true},
  {title: t('COL_OWNER_DEPARTMENT'), align: 'start', width: 270, value: 'department', sortable: true},
  {title: t('COL_APPID'), align: 'start', width: 155, value: 'applicationId', sortable: true},
  {title: t('COL_UPDATED'), align: 'start', width: 103, value: 'updated', sortable: true},
  {title: t('COL_CREATED'), align: 'start', width: 103, class: 'tableHeaderCell', value: 'created', sortable: true},
]);

const reload = async () => {
  options.value.filterString = search.value;
  options.value.filterBy = {
    status: selectedFilterStatus.value,
  };
  await projectStore.fetchProjects(options.value);
};

const searchInputChange = useDebounceFn(() => {
  search.value = searchInput.value;
}, 300);

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
</script>

<template>
  <TableLayout data-testid="projects">
    <template #buttons>
      <h1 class="text-h5">{{ t('AllProjects') }}</h1>
      <v-spacer></v-spacer>
      <v-text-field
        v-model="searchInput"
        autocomplete="off"
        :max-width="500"
        append-inner-icon="mdi-magnify"
        variant="outlined"
        density="compact"
        :label="t('labelSearch')"
        single-line
        hide-details
        clearable
        @keyup="searchInputChange"
        @click:clear="searchInputChange"></v-text-field>
    </template>
    <template #table>
      <div class="h-full">
        <v-data-table-server
          v-model:search="search"
          v-model:options="options"
          v-model:expanded="expanded"
          :loading="loading"
          density="compact"
          class="striped-table h-full"
          :headers="headers"
          :items="projects"
          :sort-by="sortBy"
          fixed-header
          items-per-page="50"
          item-value="_key"
          :items-length="projectsCount"
          :row-props="{
            class: {
              'py-8': true,
            },
          }"
          @click:row="onRowClick"
          @update:options="reload">
          <template v-slot:[`header.status`]="{column, getSortIcon, toggleSort}">
            <div class="v-data-table-header__content">
              <span>{{ column.title }}</span>
              <v-menu offset-y :close-on-content-click="false" v-model="menu">
                <template v-slot:activator="{props}">
                  <DIconButton
                    :parentProps="props"
                    icon="mdi-filter-variant"
                    :hint="t('TT_SHOW_FILTER')"
                    :color="
                      selectedFilterStatus && selectedFilterStatus.length > 0 ? 'primary' : 'secondary'
                    "></DIconButton>
                </template>
                <div style="width: 320px" class="bg-background">
                  <v-card>
                    <v-row class="d-flex justify-end ma-1 mr-2">
                      <DCloseButton @click="menu = false"></DCloseButton>
                    </v-row>

                    <v-select
                      v-model="selectedFilterStatus"
                      variant="outlined"
                      density="compact"
                      class="mx-2 pa-2"
                      autofocus
                      clearable
                      :items="projectPossibleStatuses"
                      :label="t('lbl_filter_on_status')"
                      hide-details
                      multiple
                      menu
                      transition="scale-transition"
                      persistent-clear
                      :list-props="{class: 'striped-filter-dd py-0'}"
                      @update:modelValue="reload">
                      <template v-slot:item="{props, item}">
                        <v-list-item v-bind="props" :title="undefined" class="py-0 px-2">
                          <template v-slot:prepend="{isSelected}">
                            <v-checkbox hide-details :model-value="isSelected"></v-checkbox>
                          </template>
                          <span :class="'pStatus' + (!item.value ? 'new' : item.value) + ' pStatusFilter'">
                            {{ !item.value ? 'new' : t('STATUS_' + item.value) }}
                          </span>
                        </v-list-item>
                      </template>
                      <template v-slot:selection="{item, index}">
                        <div v-if="index === 0" class="d-flex align-center">
                          <span :class="'pStatus' + (!item.value ? 'new' : item.value) + ' pStatusFilter'">
                            {{ !item.value ? 'new' : t('STATUS_' + item.value) }}
                          </span>
                        </div>
                        <span v-if="index === 1" class="pAdditionalFilter">
                          +{{ selectedFilterStatus.length - 1 }} others
                        </span>
                      </template>
                    </v-select>
                  </v-card>
                </div>
              </v-menu>
              <v-icon
                class="v-data-table-header__sort-icon"
                :icon="getSortIcon(column)"
                @click="toggleSort(column)"></v-icon>
            </div>
          </template>
          <template v-slot:[`item.updated`]="{item}">
            <DDateCellWithTooltip :value="item.updated"></DDateCellWithTooltip>
          </template>
          <template v-slot:[`item.created`]="{item}">
            <DDateCellWithTooltip :value="item.created"></DDateCellWithTooltip>
          </template>
          <template v-slot:[`item.status`]="{item}">
            <span :class="'pStatus' + (!item.status ? 'new' : item.status)">
              {{ !item.status ? 'new' : t('STATUS_' + item.status) }}
            </span>
          </template>
          <template v-slot:item.isGroup="{item}">
            <v-icon icon="mdi-check" class="mr-2" :color="item.isGroup ? 'primary' : 'tableBorderColor'" />
          </template>
          <template v-slot:[`item.actions`]="{item}">
            <ProjectsTableAction :item="item" @reload="reload()"></ProjectsTableAction>
          </template>
          <template v-slot:[`item.company`]="{item}">
            <span v-if="!item.missing">{{ item.company }}</span>
            <div v-else>
              <v-icon class="pr-2" icon="mdi-alert" color="warning" small></v-icon>
              <span>{{ t('WARNING_MISSING_DEPT') }}</span>
            </div>
          </template>
          <template v-slot:[`item.department`]="{item}">
            <span v-if="!item.missing">{{ item.department }}</span>
            <div v-else>
              <v-icon class="pr-2" color="warning" icon="mdi-alert" small></v-icon>
              <span>{{ t('WARNING_MISSING_DEPT') }}</span>
            </div>
          </template>
          <template v-slot:[`item.supplier`]="{item}">
            <span v-if="!item.supplierMissing">{{ item.supplier }}</span>
            <div v-else>
              <v-icon class="pr-2" color="warning" icon="mdi-alert" small></v-icon>
              <span>{{ t('WARNING_MISSING_DEPT') }}</span>
            </div>
          </template>
          <template v-slot:[`item.data-table-expand`]="{item}">
            <v-icon color="primary" @click.stop="toggleExpand(item)">
              {{ isExpanded(item) ? 'mdi-chevron-up' : 'mdi-chevron-down' }}
            </v-icon>
          </template>

          <template v-slot:expanded-row="{item}">
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

:deep(.v-data-table tbody tr:has(.pStatusdeprecated)) {
  color: rgb(var(--v-theme-projectDeprecated));
}
</style>
