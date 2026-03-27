// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import {
  JOB_EXECUTION_MANUAL,
  JOB_EXECUTION_ONETIME,
  JOB_EXECUTION_PERIODIC,
  JOB_STATUS_SUCCESS,
  JobDto,
  jobExecutionToString,
  jobStatusToString,
  jobTypeToString,
} from '@disclosure-portal/model/Job';
import AdminService from '@disclosure-portal/services/admin';
import {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {DataTableHeader, DataTableHeaderFilterItems} from '@shared/types/table';
import {computed, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const {info} = useSnackbar();
const breadcrumbs = useBreadcrumbsStore();

const selectedFilterExecution = ref<string[]>([]);

const jobs = ref<JobDto[]>([]);
const confirmConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const confirmOnetimeConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const configDialog = ref();
const confirmVisible = ref(false);
const confirmOnetimeVisible = ref(false);

const headers = computed<DataTableHeader[]>(() => {
  return [
    {
      title: t('COL_ACTIONS'),
      align: 'center',
      width: 120,
      maxWidth: 130,
      value: 'actions',
      sortable: false,
    },
    {title: '', value: 'data-table-expand', width: 53, maxWidth: 53, align: 'center'},
    {
      title: t('CD_NAME'),
      align: 'start',
      value: 'name',
      sortable: true,
      width: 235,
      minWidth: 235,
    },
    {
      title: t('TYPE'),
      align: 'start',
      value: 'jobType',
      sortable: true,
      width: 235,
      minWidth: 235,
    },
    {
      title: t('JOB_EXECUTION'),
      align: 'start',
      value: 'execution',
      sortable: true,
      width: 150,
      maxWidth: 160,
    },
    {
      title: t('JOB_STATUS'),
      align: 'start',
      value: 'status',
      sortable: true,
      width: 120,
      maxWidth: 130,
    },
    {
      title: t('JOB_NEXT_SCHEDULED_EXECUTION'),
      align: 'start',
      value: 'nextScheduledExecution',
      sortable: true,
      width: 200,
      minWidth: 200,
    },
    {
      title: t('CREATED'),
      align: 'start',
      value: 'created',
      sortable: true,
      width: 120,
      maxWidth: 130,
    },
    {
      title: t('UPDATED'),
      align: 'start',
      value: 'updated',
      sortable: true,
      width: 120,
      maxWidth: 130,
    },
  ];
});

const innerHeaders = computed<DataTableHeader[]>(() => {
  return [
    {
      title: t('CREATED'),
      align: 'start',
      value: 'created',
      width: 140,
      sortable: true,
    },
    {title: t('COL_LEVEL'), align: 'start', value: 'level', width: 90, sortable: false},
    {
      title: t('COL_INSTANCE'),
      align: 'start',
      width: 180,
      value: 'instance',
      sortable: false,
    },
    {title: t('COL_MESSAGE'), align: 'start', value: 'msg', sortable: false},
  ];
});

const search = ref('');
const jobsLoaded = ref(false);
const toStart = ref<JobDto | null>();
const expanded = ref<string[]>([]);

const reload = async () => {
  jobs.value = await AdminService.getJobsAll();
};

const showConfirm = (job: JobDto) => {
  toStart.value = job;
  confirmConfig.value = {
    key: job._key,
    name: job.name,
    okButtonIsDisabled: false,
    okButton: 'Btn_confirm',
    description: 'DLG_CONFIRMATION_START_JOB',
  } as IConfirmationDialogConfig;
  confirmVisible.value = true;
};

const showOnetimeConfirm = (job: JobDto) => {
  toStart.value = job;
  confirmOnetimeConfig.value = {
    key: job._key,
    name: job.name,
    okButtonIsDisabled: false,
    okButton: 'Btn_confirm',
    description: 'DLG_CONFIRMATION_START_ONETIME_JOB',
  } as IConfirmationDialogConfig;

  confirmOnetimeVisible.value = true;
};

const possibleExecutions = computed((): DataTableHeaderFilterItems[] => {
  const uniqueJobExecutions = [...new Set(jobs.value.map(({execution}) => execution))];
  return uniqueJobExecutions.map((execution) => ({
    text: t(jobExecutionToString(execution)),
    value: String(execution),
  }));
});

const filteredItems = computed((): JobDto[] =>
  jobs.value.filter((j) =>
    selectedFilterExecution.value.length >= 1
      ? selectedFilterExecution.value.some((execution) => execution === String(j.execution))
      : true,
  ),
);

const start = async () => {
  if (!toStart.value) {
    return;
  }
  await AdminService.startJob(toStart.value.jobType);
  info(t('JOB_MANUAL_START'));
  toStart.value = null;
};
const startOnetime = async () => {
  if (!toStart.value) {
    return;
  }
  await AdminService.rerunOnetimeJob(toStart.value._key);
  info(t('JOB_ONETIME_START'));
  toStart.value = null;
};

const getActionButtons = (item: JobDto): TableActionButtonsProps['buttons'] => {
  return [
    {
      icon: 'mdi-play-circle-outline',
      hint: t('TT_run_job'),
      event: 'start',
      show: item.execution === JOB_EXECUTION_MANUAL,
    },
    {
      icon: 'mdi-replay',
      hint: t('TT_run_onetime_job'),
      event: 'startOnetime',
      show: item.execution === JOB_EXECUTION_ONETIME && item.status !== JOB_STATUS_SUCCESS,
    },
    {
      icon: 'mdi-cog',
      hint: t('JOB_CONFIG'),
      event: 'config',
      show: !!item.config,
    },
  ];
};

const toggleExpand = (item: JobDto) => {
  const index = expanded.value.indexOf(item._key);
  if (index > -1) {
    expanded.value.splice(index, 1);
  } else {
    expanded.value.push(item._key);
  }
};

const isExpanded = (item: JobDto) => {
  return expanded.value.includes(item._key);
};

const initBreadcrumbs = () => {
  breadcrumbs.setCurrentBreadcrumbs([
    {title: t('BC_Dashboard'), href: '/dashboard/home', disabled: false},
    {title: t('BC_ADMIN'), href: '/dashboard/admin', disabled: false},
    {title: t('BC_JOBS'), href: '/dashboard/admin/jobs', disabled: false},
  ]);
};

onMounted(async () => {
  initBreadcrumbs();
  await reload();
  jobsLoaded.value = true;
});
</script>

<template>
  <TableLayout>
    <template #buttons>
      <h1 class="text-h5">{{ t('ADMIN_JOBS') }}</h1>
      <DCActionButton large icon="mdi-refresh" :hint="t('TT_reload')" :text="t('BTN_RELOAD')" @click="reload" />
      <v-spacer></v-spacer>
      <v-text-field
        autocomplete="off"
        :max-width="500"
        v-model="search"
        :label="t('labelSearch')"
        append-inner-icon="mdi-magnify"
        variant="outlined"
        clearable
        density="compact"
        hide-details />
    </template>
    <template #table>
      <div ref="tableGridJobs" class="fill-height">
        <v-data-table
          density="compact"
          class="striped-table custom-data-table fill-height"
          fixed-header
          :headers="headers"
          :search="search"
          :items="filteredItems"
          v-bind:expanded="expanded"
          :single-expand="true"
          :sort-by="[{key: 'updated', order: 'desc'}]"
          items-per-page="100"
          item-value="_key">
          <template #[`header.execution`]="{column, getSortIcon, toggleSort}">
            <span class="mr-1">{{ column.title }}</span>
            <GridHeaderFilterIcon
              v-model="selectedFilterExecution"
              :column="column"
              :label="t('LICENSE_CHART_STATUS')"
              :initial-selected="[String(JOB_EXECUTION_MANUAL), String(JOB_EXECUTION_PERIODIC)]"
              :allItems="possibleExecutions">
            </GridHeaderFilterIcon>
            <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
          </template>
          <template v-slot:expanded-row="{item}">
            <tr>
              <td colspan="9" class="px-0">
                <v-data-table
                  :headers="innerHeaders"
                  :items="item.log"
                  density="compact"
                  hide-default-footer
                  v-if="item.log"
                  sort-desc
                  items-per-page="-1">
                  <template #[`item.created`]="{item}">
                    <DDateCellWithTooltip show-time :value="item.created" />
                  </template>
                </v-data-table>
                <span v-else>{{ t('NO_DATA_AVAILABLE') }}</span>
              </td>
            </tr>
          </template>
          <template #[`item.jobType`]="{item}">
            <span>{{ t(jobTypeToString(item.jobType)) }}</span>
          </template>
          <template #[`item.execution`]="{item}">
            <span>{{ t(jobExecutionToString(item.execution)) }}</span>
          </template>
          <template #[`item.status`]="{item}">
            <span>{{ t(jobStatusToString(item.status)) }}</span>
          </template>
          <template #[`item.updated`]="{item}">
            <DDateCellWithTooltip :value="item.updated" />
          </template>
          <template #[`item.created`]="{item}">
            <DDateCellWithTooltip :value="item.created" />
          </template>
          <template #[`item.nextScheduledExecution`]="{item}">
            <DDateCellWithTooltip :value="item.nextScheduledExecution" />
          </template>
          <template #[`item.data-table-expand`]="{item}">
            <v-icon color="primary" @click.stop="toggleExpand(item)">
              {{ isExpanded(item) ? 'mdi-chevron-up' : 'mdi-chevron-down' }}
            </v-icon>
          </template>
          <template #[`item.actions`]="{item}">
            <TableActionButtons
              variant="compact"
              :buttons="getActionButtons(item)"
              @start="showConfirm(item)"
              @startOnetime="showOnetimeConfirm(item)"
              @config="configDialog?.open(item.config)" />
          </template>
        </v-data-table>
      </div>
    </template>
  </TableLayout>

  <ConfirmationDialog v-model:showDialog="confirmVisible" :config="confirmConfig" @confirm="start" />
  <ConfirmationDialog
    v-model:showDialog="confirmOnetimeVisible"
    :config="confirmOnetimeConfig"
    @confirm="startOnetime" />
  <JobConfigDialog ref="configDialog" />
</template>
