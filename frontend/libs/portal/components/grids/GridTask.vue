// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG // // SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {TaskDto} from '@disclosure-portal/model/Users';
import Profile from '@disclosure-portal/services/profile';
import {downloadFile} from '@disclosure-portal/utils/download';
import {RightsUtils} from '@disclosure-portal/utils/Rights';
import {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import {DataTableHeader, DataTableHeaderFilterItems, DataTableItem, SortItem} from '@shared/types/table';
import dayjs from 'dayjs';
import {indexOf} from 'lodash';
import {computed, onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute, useRouter} from 'vue-router';
import {useHeaderSettings} from '@shared/composables/useHeaderSettings';
import {useUrls} from '@shared/composables/useUrls';

const router = useRouter();
const route = useRoute();
const {t} = useI18n();
const {info, error} = useSnackbar();
const {openUrlInNewTab} = useUrls();

interface Props {
  hideEditAction?: boolean;
  readOnly?: boolean;
  fixedHeight?: boolean;
  inOwnView?: boolean;
  fetchMethod?: () => Promise<TaskDto[]>;
}

const props = withDefaults(defineProps<Props>(), {
  hideEditAction: false,
  readOnly: false,
  fixedHeight: false,
  inOwnView: false,
  fetchMethod: async () => Profile.getTasks(),
});

const loading = ref(false);
const items = ref<TaskDto[]>([]);
const search = ref('');
const taskApprovalVisible = ref(false);
const selectedRow = ref<TaskDto | null>(null);
const selectedFilterStatus = ref<string[]>([]);
const tableGridTasks = ref<HTMLElement | null>(null);
const delegateDialogVisible = ref(false);
const selectedTaskForDelegate = ref<TaskDto | null>(null);

const initialSelectedStatus = ['ACTIVE', 'PENDING'];

const possibleStatus = computed((): DataTableHeaderFilterItems[] => {
  const values: string[] = [...items.value.map((item: TaskDto) => item.status), ...initialSelectedStatus];
  return [...new Set(values)].map((item) => ({value: item}));
});

const showTaskApprovalDialog = (item: TaskDto) => {
  if (props.readOnly) return;
  selectedRow.value = {...item};
  taskApprovalVisible.value = true;
  router.push({name: 'TasksApprovalDialog', params: {id: encodeURIComponent(item.id)}});
};

const sortType = (a: TaskDto, b: TaskDto) => {
  return t(`TASK_TYPE_${b.approvalType}_${b.type}`).localeCompare(t(`TASK_TYPE_${a.approvalType}_${a.type}`));
};

const headers: DataTableHeader[] = [
  {title: 'COL_ACTIONS', align: 'center', value: 'actions', width: 120},
  {title: 'COL_STATUS', sortable: true, align: 'start', width: 120, value: 'status'},
  {title: 'TYPE', sortable: true, sortRaw: sortType, align: 'start', width: 250, value: 'approvalType'},
  {title: 'COL_CREATOR', sortable: true, align: 'start', width: 180, value: 'creator'},
  {title: 'COL_DEPARTMENT', sortable: true, align: 'start', width: 180, value: 'creatorDepartment'},
  {title: 'COL_DELEGATED_TO', sortable: true, align: 'start', width: 180, value: 'delegatedTo'},
  {title: 'COL_REFERENCE', sortable: true, align: 'start', width: 180, value: 'projectName'},
  {title: 'COL_RESULT', sortable: true, align: 'start', width: 160, value: 'resultStatus'},
  {title: 'COL_UPDATED', sortable: true, align: 'start', width: 110, value: 'updated'},
  {title: 'COL_CREATED', sortable: true, align: 'start', width: 110, value: 'created'},
];

const tableName = 'TasksGrid';
const headerSettings = useHeaderSettings({tableName, headers});
const {filteredHeaders} = headerSettings;

watch(
  () => route.params.id,
  async (newId) => {
    if (newId && (!selectedRow.value || selectedRow.value.id !== newId)) {
      const id = Array.isArray(newId) ? newId[0] : newId;
      const task = await Profile.getTask(id);
      if (task) {
        selectedRow.value = task;
        taskApprovalVisible.value = true;
      }
    }
  },
  {immediate: true, deep: true},
);

function closeTaskApprovalDialog(value: boolean) {
  if (!value) {
    selectedRow.value = null;
    router.push({path: '/dashboard/tasks'});
    reload();
  }
}
const filterOnApproval = (item: TaskDto) => {
  return selectedFilterStatus.value.length === 0 || indexOf(selectedFilterStatus.value, item.status) !== -1;
};
const filteredList = computed(() => {
  return items.value.filter(filterOnApproval);
});

const reload = async () => {
  loading.value = true;
  items.value = await props.fetchMethod();
  loading.value = false;
};

const openProject = (item: TaskDto) => {
  const linkToProject = item.isProjectGroup
    ? `/dashboard/groups/${encodeURIComponent(item.projectGuid)}`
    : `/dashboard/projects/${encodeURIComponent(item.projectGuid)}`;
  openUrlInNewTab(linkToProject);
};

const showDelegateDialog = (item: TaskDto) => {
  selectedTaskForDelegate.value = item;
  delegateDialogVisible.value = true;
};

const handleDelegateConfirm = async (delegateUserId: string) => {
  if (selectedTaskForDelegate.value) {
    try {
      await Profile.delegateTask(selectedTaskForDelegate.value.id, delegateUserId);
      info(t('TASK_DELEGATED_SUCCESS'));
      await reload();
    } catch (err) {
      console.error('Failed to update delegate task', err);
      error(t('TASK_DELEGATE_FAILED'));
    }
  }
};

const canDelegate = (item: TaskDto): boolean => {
  return (
    !props.readOnly &&
    item.status === 'ACTIVE' &&
    item.type === 'APPROVAL' &&
    item.approvalType === 'PLAUSIBILITY' &&
    item.projectType === 'Vehicle' &&
    RightsUtils.isFOSSOffice()
  );
};

const downloadCsv = async () => {
  info(t('Downloading csv file...'));
  downloadFile(`tasks_${dayjs(new Date()).format('YYYY-MM-DD_hh_mm_ss')}.csv`, Profile.downloadTasksCsv(), true);
};

const getActionButtons = (item: TaskDto): TableActionButtonsProps['buttons'] => {
  return [
    {
      icon: 'mdi-pencil',
      hint: t('TT_OPEN_TASK'),
      event: 'edit',
      show: !props.hideEditAction,
    },
    {
      icon: 'mdi-open-in-new',
      hint: t('TT_OPEN_REFERENCE'),
      event: 'openReference',
      show: true,
    },
    {
      icon: 'mdi-account-arrow-right',
      hint: t('TT_DELEGATE_TASK'),
      event: 'delegate',
      show: canDelegate(item),
    },
  ];
};

const sortItems: SortItem[] = [{key: 'updated', order: 'desc'}];

const customFilter = (value: any, query: string, item?: any) => {
  if (!query) return true;
  const searchLower = query.toLowerCase();
  const taskItem = item?.raw as TaskDto;
  if (!taskItem) return false;

  const searchableFields = [
    taskItem.status,
    taskItem.creator,
    taskItem.creatorFullName,
    taskItem.creatorDepartment,
    taskItem.creatorDepartmentDescription,
    taskItem.delegatedTo,
    taskItem.delegatedToFullName,
    taskItem.projectName,
    taskItem.approvalType,
    taskItem.type,
    taskItem.resultStatus,
  ];

  return searchableFields.some((field) => field?.toLowerCase().includes(searchLower));
};

onMounted(async () => {
  await reload();
});
</script>

<template>
  <TableLayout :has-tab="!inOwnView" :has-title="!inOwnView">
    <template v-if="inOwnView" #buttons>
      <h2 class="text-h5">{{ t('TASKS') }}</h2>
      <v-spacer></v-spacer>
      <DCActionButton
        v-if="items && items.length > 0"
        :text="t('BTN_DOWNLOAD')"
        icon="mdi-download"
        :hint="t('TT_download_tasks_csv')"
        class="align-content-center mx-2"
        @click="downloadCsv" />
      <DSearchField v-model="search" />
    </template>
    <template #table>
      <div ref="tableGridTasks" class="fill-height">
        <v-data-table
          density="compact"
          class="striped-table custom-data-table fill-height"
          :class="[props.readOnly ? 'force-border tableNoHandCursor' : 'force-border']"
          fixed-header
          :search="search"
          :custom-filter="customFilter"
          :headers="filteredHeaders"
          :items-per-page="50"
          :footer-props="{'items-per-page-options': [10, 50, 100, -1]}"
          :sort-by="sortItems"
          item-key="_key"
          :loading="loading"
          :items="filteredList"
          item-value="_key"
          @click:row="(_: Event, dataItem: DataTableItem<TaskDto>) => showTaskApprovalDialog(dataItem.item)">
          <template v-slot:[`header.actions`]="{column}">
            <GridFilterHeader :column="column">
              <template #settings>
                <HeaderSettings :column="column" :grid-name="tableName" />
              </template>
            </GridFilterHeader>
          </template>
          <template v-slot:[`header.status`]="{column, toggleSort, getSortIcon}">
            <GridFilterHeader :column="column" :getSortIcon="getSortIcon" :toggleSort="toggleSort">
              <template #filter>
                <GridHeaderFilterIcon
                  v-model="selectedFilterStatus"
                  :column="column"
                  :label="t('STATUS')"
                  :allItems="possibleStatus"
                  :initialSelected="initialSelectedStatus">
                </GridHeaderFilterIcon>
              </template>
            </GridFilterHeader>
          </template>
          <template v-slot:[`item.created`]="{item}">
            <DDateCellWithTooltip :value="item.created" />
          </template>
          <template v-slot:[`item.updated`]="{item}">
            <DDateCellWithTooltip :value="item.updated" />
          </template>
          <template v-slot:[`item.approvalType`]="{item}">
            {{ t(`TASK_TYPE_${item.approvalType}_${item.type}`) }}
          </template>
          <template v-slot:[`item.resultStatus`]="{item}">
            {{ t(`COL_APPROVAL_STATUS_${item.approvalType}_${item.resultStatus}`) }}
          </template>
          <template v-slot:[`item.creator`]="{item}">
            <span>{{ item.creatorFullName }} ({{ item.creator }})</span>
          </template>
          <template v-slot:[`item.creatorDepartment`]="{item}">
            <span v-if="item.creatorDepartmentDescription && item.creatorDepartment"
              >{{ item.creatorDepartmentDescription }} ({{ item.creatorDepartment }})</span
            >
            <span v-else>-</span>
          </template>
          <template v-slot:[`item.delegatedTo`]="{item}">
            <span v-if="item.delegatedToFullName">{{ item.delegatedToFullName }} ({{ item.delegatedTo }})</span>
            <span v-else>-</span>
          </template>
          <template v-slot:[`item.projectName`]="{item}">
            {{ `Project: ${item.projectName}` }}
          </template>
          <template v-slot:[`item.actions`]="{item}">
            <TableActionButtons
              variant="compact"
              :buttons="getActionButtons(item)"
              @edit="showTaskApprovalDialog(item)"
              @openReference="openProject(item)"
              @delegate="showDelegateDialog(item)" />
          </template>
        </v-data-table>
      </div>
    </template>
  </TableLayout>
  <TaskApprovalDialog
    v-model:showDialog="taskApprovalVisible"
    :selected-row="selectedRow"
    @update:showDialog="closeTaskApprovalDialog($event!)" />
  <DelegateTaskDialog
    v-model:showDialog="delegateDialogVisible"
    :project-key="selectedTaskForDelegate?.projectGuid || ''"
    @confirm="handleDelegateConfirm" />
</template>
