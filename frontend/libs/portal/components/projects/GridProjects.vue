<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ProjectSlim} from '@disclosure-portal/model/ProjectsResponse';
import {Rights} from '@disclosure-portal/model/Rights';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useCustomIdStore} from '@disclosure-portal/stores/customid.store';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useUserStore} from '@disclosure-portal/stores/user';
import {useWizardStore} from '@disclosure-portal/stores/wizard.store';
import {openProjectUrlByKey} from '@disclosure-portal/utils/url';
import {DataTableHeader, DataTableItem, SortItem} from '@shared/types/table';
import {chain} from 'lodash';
import {storeToRefs} from 'pinia';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';
import {useProjectUtils} from '@disclosure-portal/utils/projects';

const {t} = useI18n();
const router = useRouter();
const appStore = useAppStore();
const rights = useUserStore().getRights as Rights;
const customIdsStore = useCustomIdStore();
const projectStore = useProjectStore();
const wizardStore = useWizardStore();
const projectsUtils = useProjectUtils();

const {projects, projectsCount, loading, projectPossibleStatuses} = storeToRefs(projectStore);

const search = ref('');
const selectedFilterStatus = ref<string[]>([]);
const filterGroups = ref(false);
const menuIsGroup = ref(false);
const sortBy = ref<SortItem[]>([{key: 'updated', order: 'desc'}]);
const page = ref<string | number>(1);

const labelTools = computed(() => appStore.getLabelsTools);
const filteredList = computed<ProjectSlim[]>(() => {
  if (!Array.isArray(projects.value)) {
    return [];
  }
  let result = projects.value;
  if (selectedFilterStatus.value.length > 0) {
    result = chain(projects.value).filter(filterOnApproval).value();
  }
  if (filterGroups.value) {
    result = result.filter((item) => item.isGroup);
  }
  return result;
});
const headers = computed<DataTableHeader[]>(() => [
  {title: '', value: 'data-table-expand', width: '38'},
  {title: t('COL_ACTIONS'), align: 'center', width: 80, value: 'actions', sortable: false},
  {title: t('COL_STATUS'), sortable: true, value: 'status', width: '120'},
  {title: t('COL_GROUP'), align: 'center', sortable: true, value: 'isGroup', width: '120'},
  {title: t('COL_NAME'), align: 'start', value: 'name', width: 270, sortable: true},
  {title: t('COL_DEVELOPER_COMPANY'), align: 'start', width: 270, value: 'supplier', sortable: true},
  {title: t('COL_OWNER_COMPANY'), align: 'start', width: 270, value: 'company', sortable: true},
  {title: t('COL_OWNER_DEPARTMENT'), align: 'start', width: 270, value: 'department', sortable: true},
  {title: t('COL_APPID'), align: 'start', width: 155, value: 'applicationId', sortable: true},
  {title: t('COL_UPDATED'), align: 'start', width: 103, value: 'updated', sortable: true},
  {title: t('COL_CREATED'), align: 'start', width: 103, value: 'created', sortable: true},
]);

const filterOnApproval = (item: ProjectSlim): boolean => {
  return selectedFilterStatus.value.length === 0 || selectedFilterStatus.value.includes(item.status);
};

const reload = async () => {
  await projectStore.fetchProjects();
};

const onRowClick = (event: Event, item: DataTableItem<ProjectSlim>) => {
  const project: ProjectSlim = item.item;
  openProjectUrlByKey(project._key, router);
};

reload();

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

const customFilterTable = (rawCellValue: unknown, search: string, internalItem: any) => {
  const project = internalItem.raw as ProjectSlim;
  const lowerSearch = search.toLowerCase();

  const foundPolicy = project.policyLabels.some((l) => {
    const labelStr = labelTools.value.policyLabelsMap[l].name;
    return labelStr.toLowerCase().indexOf(lowerSearch) !== -1;
  });

  const foundProject = project.projectLabels.some((l) => {
    const labelStr = labelTools.value.projectLabelsMap[l].name;
    return labelStr.toLowerCase().indexOf(lowerSearch) !== -1;
  });

  const foundFree = (project.freeLabels || []).some((l) => l.toLowerCase().indexOf(lowerSearch) !== -1);

  const foundCell =
    (rawCellValue !== undefined && rawCellValue !== null ? String(rawCellValue) : '')
      .toLowerCase()
      .indexOf(lowerSearch) !== -1;

  let foundCustomIds: boolean;
  if (search.includes(':')) {
    const [customIdSearch, valueSearch] = search.split(':').map((s) => s.trim().toLowerCase());
    foundCustomIds = project.customIds.some((id) => {
      return (
        id.technicalId.toLowerCase().indexOf(customIdSearch) !== -1 &&
        id.value.toLowerCase().indexOf(valueSearch) !== -1
      );
    });
  } else {
    foundCustomIds = project.customIds.some((id) => {
      const cid = customIdsStore.customIds.map[id.technicalId];
      if (!cid) {
        return false;
      }
      return (
        id.value.indexOf(lowerSearch) !== -1 ||
        cid._key?.indexOf(lowerSearch) !== -1 ||
        cid.description.indexOf(lowerSearch) !== -1 ||
        cid.descriptionDE.indexOf(lowerSearch) !== -1 ||
        cid.name.indexOf(lowerSearch) !== -1 ||
        cid.nameDE.indexOf(lowerSearch) !== -1 ||
        cid.linkTemplate.indexOf(lowerSearch) !== -1
      );
    });
  }

  return foundCell || foundProject || foundPolicy || foundFree || foundCustomIds;
};
</script>

<template>
  <TableLayout data-testid="projects">
    <template #buttons>
      <h1 class="text-h5">{{ t('Projects') }}</h1>
      <DCActionButton
        v-if="rights.allowProject.create"
        large
        icon="mdi-plus"
        :hint="t('TT_add_project')"
        :text="t('BTN_ADD')"
        @click="wizardStore.openWizard()" />
      <DCActionButton
        v-if="rights.allowProject.create"
        large
        icon="mdi-plus"
        :hint="t('TT_add_group')"
        :text="t('BTN_GROUP')"
        @click="wizardStore.openWizard({isGroup: true})"></DCActionButton>
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
        clearable></v-text-field>
    </template>
    <template #table>
      <div class="fill-height">
        <v-data-table
          density="comfortable"
          class="striped-table fill-height"
          fixed-header
          :headers="headers"
          :items="filteredList"
          :items-length="projectsCount"
          :page="page"
          :sort-by="sortBy"
          :custom-filter="customFilterTable"
          items-per-page="25"
          item-value="_key"
          :loading="loading"
          :cell-props="{class: 'py-3'}"
          v-model:search="search"
          v-model:expanded="expanded"
          @click:row="onRowClick">
          <template v-slot:[`header.status`]="{column, getSortIcon, toggleSort}">
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
          <template v-slot:[`header.isGroup`]="{column, getSortIcon, toggleSort}">
            <GridFilterHeader :column="column" :getSortIcon="getSortIcon" :toggleSort="toggleSort">
              <template #filter>
                <v-menu offset-y :close-on-content-click="false" v-model="menuIsGroup">
                  <template v-slot:activator="{props}">
                    <span>
                      <v-icon class="mr-1" v-bind="props" :color="filterGroups ? 'primary' : 'default'">
                        mdi-filter-variant
                      </v-icon>
                      <Tooltip>{{ t('TT_SHOW_FILTER') }}</Tooltip>
                    </span>
                  </template>
                  <div class="bg-background w-[320px]">
                    <v-card class="d-flex justify-space-between align-center">
                      <v-checkbox hide-details v-model="filterGroups" :label="t('lbl_filter_on_group')" />
                      <DCloseButton @click="menuIsGroup = false" />
                    </v-card>
                  </div>
                </v-menu>
              </template>
            </GridFilterHeader>
          </template>
          <template v-slot:[`item.updated`]="{item}">
            <DDateCellWithTooltip :value="item.updated"></DDateCellWithTooltip>
          </template>
          <template v-slot:[`item.created`]="{item}">
            <DDateCellWithTooltip :value="item.created"></DDateCellWithTooltip>
          </template>
          <template v-slot:[`item.status`]="{item}">
            <span :style="{color: projectsUtils.getTextStatusColor(item.status)}">
              {{ t('STATUS_' + (!item.status ? 'new' : item.status)) }}
            </span>
          </template>
          <template v-slot:[`item.isGroup`]="{item}">
            <v-icon icon="mdi-check" class="mr-2" :color="item.isGroup ? 'primary' : 'tableBorderColor'"></v-icon>
          </template>
          <template v-slot:[`item.actions`]="{item}">
            <ProjectsTableAction :item="item" @reload="reload"></ProjectsTableAction>
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
            <td :colspan="headers.length" class="bg-table-header h-full cursor-default overflow-y-clip">
              <GridProjectsExpandContent :item="item" :is-async="false" />
            </td>
          </template>
        </v-data-table>
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
