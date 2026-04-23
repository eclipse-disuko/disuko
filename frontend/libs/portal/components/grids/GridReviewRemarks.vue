<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import {Checklist} from '@disclosure-portal/model/Checklist';
import {useReviewRemarkActions} from '@disclosure-portal/composables/useReviewRemarkActions';
import type {Project} from '@disclosure-portal/model/Project';
import {ReviewRemark, ReviewRemarkLevel, ReviewRemarkStatus, compareRRLevel} from '@disclosure-portal/model/Quality';
import type {BulkSetReviewRemarkStatusRequest} from '@disclosure-portal/model/ReviewRemarkBulkOperations';
import {default as ProjectService} from '@disclosure-portal/services/projects';
import versionService from '@disclosure-portal/services/version';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import {downloadFile} from '@disclosure-portal/utils/download';
import {formatDateAndTime} from '@disclosure-portal/utils/Table';
import {getIconColorReviewRemarkLevel, getIconReviewRemarkLevel} from '@disclosure-portal/utils/View';
import {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import {DataTableHeader, DataTableHeaderFilterItems} from '@shared/types/table';
import {useClipboard} from '@shared/utils/clipboard';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import {chain} from 'lodash';
import {computed, onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute} from 'vue-router';
import {useHeaderSettings} from '@shared/composables/useHeaderSettings';

const sbomStore = useSbomStore();
const projectStore = useProjectStore();
const route = useRoute();
const {t} = useI18n();
const {info: snack} = useSnackbar();
const {copyToClipboard} = useClipboard();

const {
  confirmCloseConfig,
  confirmCancelConfig,
  confirmReopenConfig,
  closeVisible,
  cancelVisible,
  reopenVisible,
  openCloseRemarkDialog,
  openCancelRemarkDialog,
  openReopenRemarkDialog,
  doCloseRemark: doCloseRemarkAction,
  doCancelRemark: doCancelRemarkAction,
  doReopenRemark: doReopenRemarkAction,
} = useReviewRemarkActions();

const items = ref<ReviewRemark[]>([]);
const on = ref(false);
const search = ref('');
const loading = ref(false);
const selectedFilterLevel = ref<string[]>([]);
const selectedFilterStatus = ref<string[]>([]);
const selectedFilterSbom = ref<string[]>([]);
const tableHeight = ref(0);
const reviewRemarkDialog = ref();
const viewRemarkDialog = ref();
const confirmBulkCloseConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const confirmBulkCancelConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const executeDialog = ref();
const bulkCloseVisible = ref(false);
const bulkCancelVisible = ref(false);
const lists = ref<Checklist[]>([]);

const projectModel = computed((): Project => projectStore.currentProject!);
const version = computed(() => sbomStore.getCurrentVersion);
const checklistAvailable = computed(() => lists.value.length > 0);

const possibleLevel = computed((): DataTableHeaderFilterItems[] => {
  if (!items.value) {
    return [];
  }

  const allReviewRemarkLevels = items.value.map((remark: ReviewRemark) => remark.level);

  return [...new Set(allReviewRemarkLevels)].map(
    (remarkLevel: string) =>
      ({
        value: remarkLevel,
        text: t('REMARK_LEVEL_' + remarkLevel),
        icon: getIconReviewRemarkLevel(remarkLevel as ReviewRemarkLevel),
        iconColor: getIconColorReviewRemarkLevel(remarkLevel as ReviewRemarkLevel),
      }) as DataTableHeaderFilterItems,
  );
});

const possibleSbom = computed((): DataTableHeaderFilterItems[] => {
  if (!items.value) {
    return [];
  }

  const allReviewRemarkSboms = items.value.filter(
    (remark: ReviewRemark) => !!(remark.sbomName && remark.sbomName.trim() && remark.sbomId),
  );

  return (
    chain(allReviewRemarkSboms)
      .sortBy('sbomUploaded')
      .map((remark: ReviewRemark) => {
        const uploadedDate = remark.sbomUploaded ? formatDateAndTime(remark.sbomUploaded.toString()) : '';
        return {
          value: remark.sbomId || '',
          text: `${uploadedDate}- ${remark.sbomName}`,
        } as DataTableHeaderFilterItems;
      })
      .uniqBy('value')
      // oxlint-disable-next-line unicorn/no-array-reverse
      .reverse()
      .value()
  );
});

const possibleStatus = computed((): DataTableHeaderFilterItems[] => {
  if (!items.value) {
    return [];
  }

  const allReviewRemarkStatus = items.value.map((remark: ReviewRemark) => remark.status);

  return [...new Set(allReviewRemarkStatus)].map(
    (remarkStatus: string) =>
      ({
        value: remarkStatus,
        text: t('REMARK_STATUS_' + remarkStatus),
      }) as DataTableHeaderFilterItems,
  );
});

const headers: DataTableHeader[] = [
  {
    title: 'COL_ACTIONS',
    align: 'center',
    width: 120,
    value: 'actions',
  },
  {
    title: 'COL_LEVEL',
    align: 'center',
    value: 'level',
    sortable: true,
    width: 150,
    sort: compareRRLevel,
  },
  {
    title: 'COL_STATUS',
    align: 'start',
    value: 'status',
    width: 150,
    sortable: true,
  },
  {
    title: 'COL_REVIEW_REMARK',
    align: 'start',
    value: 'title',
    width: 250,
  },
  {
    title: 'COMPONENTS',
    align: 'start',
    value: 'components',
    width: 125,
    sortable: true,
  },
  {
    title: 'SBOM_REFERENCE',
    align: 'start',
    value: 'sbomName',
    width: 200,
    sortable: true,
  },
  {
    title: 'LICENSES',
    align: 'start',
    value: 'licenses',
    width: 125,
    sortable: true,
  },
  {
    title: 'COL_CREATOR',
    align: 'start',
    value: 'author',
    width: 120,
    sortable: true,
  },
  {
    title: 'COL_ORIGIN',
    align: 'start',
    value: 'origin',
    width: 80,
    sortable: true,
  },
  {
    title: 'COL_CLOSED',
    align: 'start',
    value: 'closed',
    width: 100,
    sortable: true,
  },
  {
    title: 'COL_CREATED',
    align: 'start',
    value: 'created',
    width: 100,
    sortable: true,
  },
  {
    title: 'COL_UPDATED',
    align: 'start',
    value: 'updated',
    width: 100,
    sortable: true,
  },
];

const tableName = 'ReviewRemarksGrid';
const headerSettings = useHeaderSettings({tableName, headers});
const {filteredHeaders} = headerSettings;

const sortItems = ref([{key: 'level', order: 'desc' as const}]);

watch(
  () => route.query,
  async () => {
    filterOnReviewRemarkLevel();
  },
);

const getRemarkTextForClipboard = (item: ReviewRemark): string => {
  return `ID:
${item.key}

Title:
${item.title}

Description:
${item.description}

Level:
${t('REMARK_LEVEL_' + item.level)}

SBOM:
${item.sbomName}

Components:
${item.components ? item.components.map((c) => `${c.componentName} (${c.componentVersion})`).join(';\n') : ''}

Licenses:
${item.licenses ? item.licenses.map((l) => (l.licenseName === '' ? `${l.licenseId} (${t('TT_REVIEW_REMARK_DIALOG_LICENSE_UNKNOWN')})` : l.licenseName)).join(';\n') : ''}
`;
};

const filterOnReviewRemarkLevel = () => {
  const reviewRemarkLevel = route.query.reviewRemarkLevel as string;
  if (reviewRemarkLevel) {
    selectedFilterLevel.value = [reviewRemarkLevel];
    selectedFilterStatus.value = [];
    if (possibleStatus.value.some((status) => status.value === ReviewRemarkStatus.OPEN)) {
      selectedFilterStatus.value.push(ReviewRemarkStatus.OPEN as string);
    }
    if (possibleStatus.value.some((status) => status.value === ReviewRemarkStatus.IN_PROGRESS)) {
      selectedFilterStatus.value.push(ReviewRemarkStatus.IN_PROGRESS as string);
    }
    if (possibleStatus.value.some((status) => status.value === ReviewRemarkStatus.CLOSED)) {
      selectedFilterStatus.value.push(ReviewRemarkStatus.CLOSED as string);
    }
  }
};

const reload = async (): Promise<void> => {
  loading.value = true;
  const uniqueSelectedKeys = new Set(selected.value.map(({key}) => key));
  items.value = (await versionService.getReviewRemarks(projectModel.value._key, version.value._key)).data;
  selected.value = items.value.filter((item) => uniqueSelectedKeys.has(item.key));
  loading.value = false;
};

onMounted(async () => {
  filterOnReviewRemarkLevel();
  if (projectModel.value.accessRights.allowExecuteChecklist) {
    await ProjectService.getApplicableChecklists(projectModel.value._key).then((res) => (lists.value = res.data));
  }
  await reload();
});

const customFilter = (_: unknown, searchTerm: string | null, item: unknown): boolean => {
  if (!searchTerm) return true;

  const itemObj = item as {raw?: ReviewRemark};
  if (!itemObj?.raw) return false;

  const reviewRemark = itemObj.raw;
  const searchTerms = searchTerm.toLowerCase().split(' ');
  const itemText = [
    reviewRemark.title,
    reviewRemark.description,
    reviewRemark.sbomName,
    ...(reviewRemark.components?.map((c) => `${c.componentName} ${c.componentVersion}`) || []),
    ...(reviewRemark.licenses?.flatMap((l) => [l.licenseName, l.licenseId]) || []),
  ]
    .filter(Boolean)
    .join(' ')
    .toLowerCase();

  return searchTerms.every((term) => itemText.includes(term));
};

const filteredList = computed(() => {
  return items.value.filter((item: ReviewRemark) => filterOnLevel(item) && filterOnStatus(item) && filterOnSbom(item));
});

const filterOnLevel = (item: ReviewRemark): boolean => {
  if (!selectedFilterLevel.value.length) {
    return true;
  }
  return selectedFilterLevel.value.includes(item.level);
};

const filterOnStatus = (item: ReviewRemark): boolean => {
  if (!selectedFilterStatus.value.length) {
    return true;
  }
  return selectedFilterStatus.value.includes(item.status);
};

const filterOnSbom = (item: ReviewRemark): boolean => {
  if (!selectedFilterSbom.value.length) {
    return true;
  }
  return selectedFilterSbom.value.includes(item.sbomId || '');
};

const doCloseRemark = async (config: IConfirmationDialogConfig) => {
  await doCloseRemarkAction(config, projectModel.value._key, version.value._key, reload);
};

const doCancelRemark = async (config: IConfirmationDialogConfig) => {
  await doCancelRemarkAction(config, projectModel.value._key, version.value._key, reload);
};

const doReopenRemark = async (config: IConfirmationDialogConfig) => {
  await doReopenRemarkAction(config, projectModel.value._key, version.value._key, reload);
};

const openReviewRemarkDialog = (toEdit?: ReviewRemark) => {
  reviewRemarkDialog.value?.open({presetItem: toEdit});
};

const copyRemarkToClipboard = (item: ReviewRemark) => {
  const content = getRemarkTextForClipboard(item);
  copyToClipboard(content);
};

const downloadReviewRemarksCsv = async () => {
  downloadFile(
    `${projectModel.value.name}_${version.value.name}_review_remarks.csv`,
    ProjectService.downloadReviewRemarksForSbomCsv(projectModel.value._key, version.value._key),
    true,
  );
};

const viewRemark = (row: ReviewRemark) => {
  if (row) {
    viewRemarkDialog.value!.open(row);
  }
};

const selected = ref<ReviewRemark[]>([]);

const openBulkCloseDialog = () => {
  if (!selected.value.length) return;
  const openRemarks = selected.value.filter(
    (remark) => remark.status === ReviewRemarkStatus.OPEN || remark.status === ReviewRemarkStatus.IN_PROGRESS,
  );
  if (!openRemarks.length) return;

  confirmBulkCloseConfig.value = {
    key: 'bulk-close',
    name: '',
    type: ConfirmationType.CONFIRM,
    description: t('DLG_CONFIRMATION_DESCRIPTION_BULK_CLOSE_REMARKS'),
    okButton: 'Btn_confirm',
  };
  bulkCloseVisible.value = true;
};

const openBulkCancelDialog = () => {
  if (!selected.value.length) {
    return;
  }
  const openRemarks = selected.value.filter(
    (remark) => remark.status === ReviewRemarkStatus.OPEN || remark.status === ReviewRemarkStatus.IN_PROGRESS,
  );
  if (!openRemarks.length) {
    return;
  }

  confirmBulkCancelConfig.value = {
    key: 'bulk-cancel',
    name: '',
    type: ConfirmationType.CONFIRM,
    description: t('DLG_CONFIRMATION_DESCRIPTION_BULK_CANCEL_REMARKS'),
    okButton: 'Btn_confirm',
  };
  bulkCancelVisible.value = true;
};

const doBulkCloseRemarks = async () => {
  loading.value = true;
  const openRemarks = selected.value.filter(
    (remark) => remark.status === ReviewRemarkStatus.OPEN || remark.status === ReviewRemarkStatus.IN_PROGRESS,
  );

  try {
    const req: BulkSetReviewRemarkStatusRequest = {
      remarkKeys: openRemarks.map((remark) => remark.key),
      status: ReviewRemarkStatus.CLOSED,
    };

    await versionService.bulkSetReviewRemarkStatus(projectModel.value._key, version.value._key, req).catch((error) => {
      console.error('Failed to process bulk close operation:', error);
    });
  } catch (error) {
    console.error('Error in bulk close operation:', error);
  } finally {
    snack(t('DIALOG_remark_closed'));
    selected.value = [];
    await reload();
    loading.value = false;
  }
};

const doBulkCancelRemarks = async () => {
  loading.value = true;
  const openRemarks = selected.value.filter(
    (remark) => remark.status === ReviewRemarkStatus.OPEN || remark.status === ReviewRemarkStatus.IN_PROGRESS,
  );

  try {
    const req: BulkSetReviewRemarkStatusRequest = {
      remarkKeys: openRemarks.map((remark) => remark.key),
      status: ReviewRemarkStatus.CANCELLED,
    };

    await versionService
      .bulkSetReviewRemarkStatus(projectModel.value._key, version.value._key, req)
      .catch((error: Error) => {
        console.error('Failed to process bulk cancel operation:', error);
      });
  } catch (error) {
    console.error('Error in bulk cancel operation:', error);
  } finally {
    snack(t('DIALOG_remark_cancelled'));
    selected.value = [];
    await reload();
    loading.value = false;
  }
};

const rowClick = (_: MouseEvent, row: {item: ReviewRemark}) => {
  selected.value = row.item ? [row.item] : [];
  viewRemark(row.item);
};

const getActionButtons = (item: ReviewRemark): TableActionButtonsProps['buttons'] => {
  const isOpenOrInProgress = item.status === ReviewRemarkStatus.OPEN || item.status === ReviewRemarkStatus.IN_PROGRESS;
  const isClosedOrCancelled = item.status === ReviewRemarkStatus.CLOSED || item.status === ReviewRemarkStatus.CANCELLED;
  const hasComments = item && item.events && item.events.length > 0;

  return [
    {
      icon: hasComments ? 'mdi-comment-multiple-outline' : 'mdi-comment-outline',
      event: 'view',
      show: true,
    },
    {
      icon: 'mdi-pencil',
      hint: t('TT_edit_review_remark'),
      event: 'edit',
      show: isOpenOrInProgress,
      disabled: projectModel.value.isDeprecated,
    },
    {
      icon: 'mdi-content-copy',
      hint: t('TT_COPY_REVIEW_REMARK'),
      event: 'copy',
      show: true,
    },
    {
      icon: 'mdi-check',
      hint: t('TT_close_review_remark'),
      event: 'close',
      show: isOpenOrInProgress,
      disabled: projectModel.value.isDeprecated,
    },
    {
      icon: 'mdi-cancel',
      hint: t('TT_cancel_review_remark'),
      event: 'cancel',
      show: isOpenOrInProgress,
      disabled: projectModel.value.isDeprecated,
    },
    {
      icon: 'mdi-refresh',
      hint: t('TT_reopen_review_remark'),
      event: 'reopen',
      show: isClosedOrCancelled,
      disabled: projectModel.value.isDeprecated,
    },
  ];
};

// Reset selection on load
onMounted(() => {
  selected.value = [];
});
</script>

<template>
  <div class="h-[calc(100%-56px)]">
    <div class="flex flex-col justify-between gap-5 pb-1 md:flex-row">
      <Stack direction="row">
        <DCActionButton
          v-if="projectModel.accessRights.allowExecuteChecklist"
          :text="t('CHECKLIST')"
          icon="mdi-plus"
          :disabled="!checklistAvailable"
          @click="executeDialog?.open(lists)" />
        <DCActionButton
          v-if="!projectModel.isDeprecated"
          :text="t('BTN_ADD')"
          icon="mdi-plus"
          :hint="t('TT_new_remark')"
          @click="() => openReviewRemarkDialog()" />
        <DCActionButton
          :text="t('BTN_DOWNLOAD')"
          icon="mdi-download"
          :hint="t('TT_DOWNLOAD_REVIEW_REMARKS')"
          @click="downloadReviewRemarksCsv" />
        <DCActionButton
          :text="t('BTN_CLOSE')"
          icon="mdi-check-all"
          @click="openBulkCloseDialog"
          v-if="!projectModel.isDeprecated && selected.length > 0"
          :hint="selected.length > 0 ? t('TT_BULK_CLOSE_REMARK') : t('TT_BULK_CLOSE_SELECT_REMARK')">
        </DCActionButton>
        <DCActionButton
          :text="t('BTN_CANCEL')"
          icon="mdi-cancel"
          @click="openBulkCancelDialog"
          v-if="!projectModel.isDeprecated && selected.length > 0"
          :hint="selected.length > 0 ? t('TT_BULK_CANCEL_REMARK') : t('TT_BULK_CANCEL_SELECT_REMARK')">
        </DCActionButton>
      </Stack>

      <DSearchField v-model="search" />
    </div>

    <v-data-table
      v-model="selected"
      :loading="loading"
      item-key="key"
      :items="filteredList"
      :headers="filteredHeaders"
      :search="search"
      :custom-filter="customFilter"
      :height="tableHeight"
      density="compact"
      fixed-header
      class="striped-table my-0 h-full py-0"
      :sort-by="sortItems"
      sort-desc
      show-select
      single-select
      return-object
      items-per-page="100"
      @click:row="rowClick"
      @update:modelValue="(val) => (selected = val || [])">
      <template v-slot:[`header.actions`]="{column}">
        <HeaderSettings :column="column" :grid-name="tableName" :show-borders="false" />
        <span>{{ column.title }}</span>
      </template>
      <template v-slot:[`header.level`]="{column, getSortIcon, toggleSort}">
        <GridFilterHeader :column="column" :getSortIcon="getSortIcon" :toggleSort="toggleSort">
          <template #filter>
            <GridHeaderFilterIcon
              v-model="selectedFilterLevel"
              :column="column"
              :label="t('Lbl_filter_status')"
              :allItems="possibleLevel">
            </GridHeaderFilterIcon>
          </template>
        </GridFilterHeader>
      </template>
      <template v-slot:[`header.status`]="{column, getSortIcon, toggleSort}">
        <GridFilterHeader :column="column" :getSortIcon="getSortIcon" :toggleSort="toggleSort">
          <template #filter>
            <GridHeaderFilterIcon
              v-model="selectedFilterStatus"
              :column="column"
              :label="t('Lbl_filter_status')"
              :allItems="possibleStatus">
            </GridHeaderFilterIcon>
          </template>
        </GridFilterHeader>
      </template>
      <template v-slot:[`header.sbomName`]="{column, getSortIcon, toggleSort}">
        <GridFilterHeader :column="column" :getSortIcon="getSortIcon" :toggleSort="toggleSort">
          <template #filter>
            <GridHeaderFilterIcon
              v-model="selectedFilterSbom"
              :column="column"
              :label="t('Lbl_filter_sbom')"
              :allItems="possibleSbom">
            </GridHeaderFilterIcon>
          </template>
        </GridFilterHeader>
      </template>
      <template v-slot:[`item.level`]="{item}">
        <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" bottom content-class="dpTooltip">
          <template v-slot:activator="{}">
            <v-icon v-on="on" :color="getIconColorReviewRemarkLevel(item.level)">
              {{ getIconReviewRemarkLevel(item.level) }}
            </v-icon>
          </template>
          <span>{{ t('REMARK_LEVEL_' + item.level) }}</span>
        </v-tooltip>
      </template>
      <template v-slot:[`item.status`]="{item}">
        <span>{{ t('REMARK_STATUS_' + item.status) }}</span>
      </template>
      <template v-slot:[`item.title`]="{item}">
        <span>{{ item.title }}</span>
      </template>
      <template v-slot:[`item.components`]="{item}">
        <Truncated>
          {{
            item.components ? item.components.map((c) => `${c.componentName} (${c.componentVersion})`).join(';\n') : ''
          }}
        </Truncated>
      </template>
      <template v-slot:[`item.sbomName`]="{item}">
        <div class="d-flex flex-column">
          <span>{{ item.sbomName }}</span>
          <DDateCellWithTooltip v-if="item.sbomName && item.sbomUploaded" :value="item.sbomUploaded.toString()">
          </DDateCellWithTooltip>
        </div>
      </template>
      <template v-slot:[`item.licenses`]="{item}">
        <Truncated>
          {{
            item.licenses
              ? item.licenses
                  .map((l) => (l.licenseName ? `${l.licenseName} (${l.licenseId})` : l.licenseId))
                  .join(';\n')
              : ''
          }}
        </Truncated>
      </template>
      <template v-slot:[`item.created`]="{item}">
        <DDateCellWithTooltip :value="item.created"></DDateCellWithTooltip>
      </template>
      <template v-slot:[`item.updated`]="{item}">
        <DDateCellWithTooltip :value="item.updated"></DDateCellWithTooltip>
      </template>
      <template v-slot:[`item.closed`]="{item}">
        <DDateCellWithTooltip :value="item.closed" v-if="item.closed"></DDateCellWithTooltip>
      </template>
      <template v-slot:[`item.actions`]="{item}">
        <Stack direction="row" justify="center" align="center" class="gap-0">
          <TableActionButtons
            variant="compact"
            :buttons="getActionButtons(item)"
            @view="viewRemark(item)"
            @edit="openReviewRemarkDialog(item)"
            @copy="copyRemarkToClipboard(item)"
            @close="openCloseRemarkDialog(item)"
            @cancel="openCancelRemarkDialog(item)"
            @reopen="openReopenRemarkDialog(item)" />
        </Stack>
      </template>
    </v-data-table>
  </div>
  <ReviewRemarkDialog ref="reviewRemarkDialog" @reload="reload"></ReviewRemarkDialog>
  <ConfirmationDialog v-model:showDialog="closeVisible" :config="confirmCloseConfig" @confirm="doCloseRemark">
  </ConfirmationDialog>
  <ConfirmationDialog v-model:showDialog="cancelVisible" :config="confirmCancelConfig" @confirm="doCancelRemark">
  </ConfirmationDialog>
  <ConfirmationDialog v-model:showDialog="reopenVisible" :config="confirmReopenConfig" @confirm="doReopenRemark">
  </ConfirmationDialog>
  <ConfirmationDialog
    v-model:showDialog="bulkCloseVisible"
    :config="confirmBulkCloseConfig"
    @confirm="doBulkCloseRemarks"></ConfirmationDialog>
  <ConfirmationDialog
    v-model:showDialog="bulkCancelVisible"
    :config="confirmBulkCancelConfig"
    @confirm="doBulkCancelRemarks"></ConfirmationDialog>
  <ReviewRemarksDetailsDialog
    ref="viewRemarkDialog"
    :project-uuid="projectModel._key"
    :version-uuid="version._key"
    @reload="reload"
    @close-remark="openBulkCloseDialog"></ReviewRemarksDetailsDialog>
  <ChecklistExecuteDialog ref="executeDialog" @reload="reload"></ChecklistExecuteDialog>
</template>
