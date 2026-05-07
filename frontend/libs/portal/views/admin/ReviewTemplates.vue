<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import {useFormSubmission} from '@disclosure-portal/components/dialog/reviewTemplate/formSubmission';
import useDimensions from '@disclosure-portal/composables/useDimensions';
import {ReviewRemarkLevel, compareRRLevel} from '@disclosure-portal/model/Quality';
import {ReviewTemplate} from '@disclosure-portal/model/ReviewTemplate';
import adminService from '@disclosure-portal/services/admin';
import {getCssClassForTableRow} from '@disclosure-portal/utils/Table';
import {getIconColorReviewRemarkLevel, getIconReviewRemarkLevel} from '@disclosure-portal/utils/View';
import {downloadFile} from '@disclosure-portal/utils/download';
import eventBus from '@shared/utils/eventbus';
import {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {DataTableHeader, DataTableHeaderFilterItems, DataTableItem, SortItem} from '@shared/types/table';
import {AxiosResponse} from 'axios';
import dayjs from 'dayjs';
import {computed, nextTick, onMounted, onUnmounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const {calculateHeight} = useDimensions();
const {submitForm} = useFormSubmission();
const {info} = useSnackbar();
const {dashboardCrumbs, ...breadcrumbs} = useBreadcrumbsStore();

const headers = computed<DataTableHeader[]>(() => {
  return [
    {title: t('COL_ACTIONS'), align: 'center', width: 80, maxWidth: 100, value: 'actions'},
    {title: t('NPV_DIALOG_TF_TITLE'), align: 'start', value: 'title', width: 200, minWidth: 200, sortable: true},
    {
      title: t('COL_LEVEL'),
      align: 'start',
      value: 'level',
      width: 130,
      minWidth: 130,
      maxWidth: 140,
      sortable: true,
      sort: compareRRLevel,
    },
    {
      title: t('NP_DIALOG_TF_DESCRIPTION'),
      align: 'start',
      width: 130,
      minWidth: 130,
      maxWidth: 180,
      value: 'description',
      sortable: true,
    },
    {
      title: t('NPV_DIALOG_TF_SOURCE'),
      align: 'start',
      width: 130,
      minWidth: 130,
      maxWidth: 180,
      value: 'source',
      sortable: true,
    },
    {title: t('COL_UPDATED'), sortable: true, align: 'start', width: 110, maxWidth: 120, value: 'updated'},
    {title: t('COL_CREATED'), sortable: true, align: 'start', width: 110, maxWidth: 120, value: 'created'},
  ];
});
const dialogVisible = ref(false);
const formMode = ref<'create' | 'edit'>('create');
const errorMessage = ref<string | undefined>();
const reviewTemplates = ref<ReviewTemplate[]>([]);
const search = ref('');
const selectedFilterLevel = ref<ReviewRemarkLevel[]>([]);
const allLevel = ref([ReviewRemarkLevel.GREEN, ReviewRemarkLevel.YELLOW, ReviewRemarkLevel.RED]);
const tableHeight = ref(0);
const dataTableAsElement = ref<HTMLElement | null>(null);
const editingData = ref<ReviewTemplate | null>(null);
const confirmVisible = ref(false);
const confirmConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);

const possibleLevel = computed((): DataTableHeaderFilterItems[] => {
  const uniqueLevels = [...new Set(reviewTemplates.value.map((template) => template.level))];
  return uniqueLevels.map((level) => ({
    value: level,
    text: t('REMARK_LEVEL_' + level),
    icon: getIconReviewRemarkLevel(level),
    iconColor: getIconColorReviewRemarkLevel(level),
  }));
});

const reload = async () => {
  reviewTemplates.value = (await adminService.getReviewTemplates()).data;
};

const doDelete = async (key?: string) => {
  await adminService.deleteReviewTemplate(key!);
  info(t('DIALOG_REVIEW_TEMPLATE_DELETE_SUCCESS'));
  await reload();
};

const sortTable = () => {
  return [{key: 'title', order: 'asc'} as SortItem];
};

const showConfirmDelete = (item: ReviewTemplate) => {
  confirmConfig.value = {
    type: ConfirmationType.DELETE,
    key: item._key,
    name: item.title,
    okButtonIsDisabled: false,
    okButton: 'BTN_DELETE',
    description: 'DLG_CONFIRMATION_DESCRIPTION',
  } as IConfirmationDialogConfig;
  confirmVisible.value = true;
};

const initBreadcrumbs = () => {
  breadcrumbs.setCurrentBreadcrumbs([
    ...dashboardCrumbs,
    {title: t('DB_TITLE_REVIEW_TEMPLATES'), href: '/dashboard/templates/review'},
  ]);
};

const onRowClick = (_: Event, table: DataTableItem<ReviewTemplate>) => {
  openDialog('edit', table.item._key);
};

const openDialog = async (mode: 'create' | 'edit', id = '') => {
  errorMessage.value = undefined;
  formMode.value = mode;

  if (mode === 'edit' && id) {
    const response: AxiosResponse<ReviewTemplate> = await adminService.getReviewTemplate(id);
    editingData.value = response.data;
  } else {
    editingData.value = null;
  }
  dialogVisible.value = true;
};

const filteredReviewTemplates = computed(() => {
  if (selectedFilterLevel.value.length === 0) {
    return reviewTemplates.value;
  }

  return reviewTemplates.value.filter((template) => selectedFilterLevel.value.includes(template.level));
});

const handleSave = async (formData: ReviewTemplate) => {
  errorMessage.value = undefined;
  await submitForm(formData, formMode.value);
  dialogVisible.value = false;
  await reload();
};

const downloadCsv = async () => {
  downloadFile(
    'review_remark_templates' + dayjs(new Date()).format('YYYY-MM-DD_hh_mm_ss') + '.csv',
    adminService.downloadReviewTemplateCSV(),
    true,
  );
};

const updateTableHeight = () => {
  nextTick(() => {
    if (dataTableAsElement.value) {
      tableHeight.value = calculateHeight(dataTableAsElement.value, false);
    }
  });
};
const onConfirm = (config: IConfirmationDialogConfig) => {
  if (config.type === ConfirmationType.DELETE) {
    doDelete(config.key);
  }
};

const getActionButtons = (_: ReviewTemplate): TableActionButtonsProps['buttons'] => {
  return [
    {
      icon: 'mdi-pencil',
      hint: t('TT_EDIT_REVIEW_TEMPLATE'),
      event: 'edit',
    },
    {
      icon: 'mdi-delete',
      hint: t('TT_DELETE_REVIEW_TEMPLATE'),
      event: 'delete',
    },
  ];
};

onMounted(async () => {
  initBreadcrumbs();
  updateTableHeight();
  eventBus.on('window-resize', updateTableHeight);

  await reload();
});

onUnmounted(() => {
  eventBus.off('window-resize', updateTableHeight);
});
</script>

<template>
  <TableLayout>
    <template #buttons>
      <h1 class="text-h5">{{ t('DB_TITLE_REVIEW_TEMPLATES') }}</h1>
      <DCActionButton
        large
        icon="mdi-plus"
        :text="t('BTN_ADD')"
        :hint="t('TT_ADD_REVIEW_TEMPLATE')"
        @click="openDialog('create')" />
      <v-spacer></v-spacer>
      <DCActionButton
        large
        class="align-content-center"
        icon="mdi-download"
        :text="t('BTN_DOWNLOAD')"
        :hint="t('TT_download_label_csv')"
        @click="downloadCsv" />
      <DSearchField v-model="search" />
    </template>
    <template #table>
      <div ref="dataTableAsElement" class="fill-height">
        <v-data-table
          class="striped-table fill-height"
          density="compact"
          fixed-header
          :sort-by="sortTable()"
          @click:row="onRowClick"
          :headers="headers"
          :items-per-page="-1"
          :footer-props="{'items-per-page-options': [10, 50, 100, -1]}"
          :items="filteredReviewTemplates"
          :item-class="getCssClassForTableRow"
          v-model:search="search"
          :height="tableHeight">
          <template #[`header.level`]="{column, getSortIcon, toggleSort}">
            <span class="mr-1">{{ column.title }}</span>
            <GridHeaderFilterIcon
              v-model="selectedFilterLevel"
              :column="column"
              :label="t('LABEL_FILTER_STATUS')"
              :allItems="possibleLevel">
            </GridHeaderFilterIcon>
            <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
          </template>

          <template #[`item.level`]="{item}">
            <v-icon :color="getIconColorReviewRemarkLevel(item.level)">
              {{ getIconReviewRemarkLevel(item.level) }}
            </v-icon>
            <Tooltip location="bottom">
              <span>{{ t('REMARK_LEVEL_' + item.level) }}</span>
            </Tooltip>
          </template>

          <template #[`item.source`]="{item}">
            {{ item.source }}
          </template>
          <template #[`item.created`]="{item}">
            <DDateCellWithTooltip :value="item.created" />
          </template>
          <template #[`item.updated`]="{item}">
            <DDateCellWithTooltip :value="item.updated" />
          </template>
          <template #[`item.actions`]="{item}">
            <TableActionButtons
              variant="normal"
              :buttons="getActionButtons(item)"
              @edit="openDialog('edit', item._key)"
              @delete="showConfirmDelete(item)" />
          </template>
        </v-data-table>
      </div>
    </template>
  </TableLayout>

  <ReviewTemplateDialog
    v-model:dialog="dialogVisible"
    :initial-data="editingData"
    :mode="formMode"
    :errorMessage="errorMessage"
    @save="handleSave"
    :levels="allLevel" />
  <ConfirmationDialog v-model:showDialog="confirmVisible" :config="confirmConfig" @confirm="onConfirm" />
</template>
