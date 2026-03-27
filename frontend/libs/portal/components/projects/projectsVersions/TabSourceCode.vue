// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script lang="ts" setup>
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import useDimensions from '@disclosure-portal/composables/useDimensions';
import {ExternalSource} from '@disclosure-portal/model/VersionDetails';
import VersionService from '@disclosure-portal/services/version';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import eventBus from '@disclosure-portal/utils/eventbus';
import {formatDateAndTime} from '@disclosure-portal/utils/Table';
import {formatDateTimeShort, getStrWithMaxLength} from '@disclosure-portal/utils/View';
import TableActionButtons, {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import {DataTableHeader, SortItem} from '@shared/types/table';
import {useClipboard} from '@shared/utils/clipboard';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import dayjs from 'dayjs';
import {computed, nextTick, onMounted, onUnmounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const appStore = useAppStore();
const currentProject = computed(() => useProjectStore().currentProject!);
const {t} = useI18n();
const snackbar = useSnackbar();
const {copyToClipboard} = useClipboard();
const {calculateHeight} = useDimensions();

const confirmConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const sourceCodeHistory = ref<ExternalSource[]>([]);
const sourceCode = ref<ExternalSource[]>([]);
const search = ref('');
const loading = ref(false);
const confirmationDialog = ref(null);
const dlgExternalSource = ref(null);
const tableHeight = ref(0);
const scGrid = ref<HTMLElement | null>(null);

const labelTools = computed(() => appStore.getLabelsTools);
const version = computed(() => appStore.currentVersion);
const headersSources = computed<DataTableHeader[]>(() => {
  return [
    {
      width: 80,
      title: t('COL_ACTIONS'),
      sortable: false,
      align: 'center',
      class: 'tableHeaderCell',
      value: 'actions',
    },
    {
      title: t('COL_URL'),
      align: 'start',
      class: 'tableHeaderCell',
      value: 'url',
      width: 350,
      sortable: true,
    },
    {
      title: t('COL_DESCRIPTION'),
      align: 'start',
      class: 'tableHeaderCell',
      value: 'comment',
      width: 300,
      sortable: true,
    },
    {
      title: t('COL_ORIGIN'),
      align: 'start',
      class: 'tableHeaderCell',
      value: 'origin',
      width: 120,
      sortable: true,
    },
    {
      title: t('COL_UPLOADER'),
      align: 'start',
      class: 'tableHeaderCell',
      width: 120,
      value: 'uploader',
      sortable: true,
    },
    {
      title: t('COL_CREATED'),
      width: 120,
      align: 'start',
      class: 'tableHeaderCell',
      value: 'created',
      sortable: true,
    },
  ];
});

const updateTableHeight = async () => {
  await nextTick();
  if (scGrid.value) {
    tableHeight.value = calculateHeight(scGrid.value, false, true, ['divDisclaimerText']);
  }
};
const sortBy: SortItem[] = [{key: 'created', order: 'desc'}];

const reload = async () => {
  if (loading.value) return;
  loading.value = true;
  if (currentProject.value._key) {
    const response = await VersionService.getExternalSources(currentProject.value._key, version.value._key);
    sourceCode.value = response.data.filter((source: ExternalSource) => source.sourceType === 'upload');
    sourceCodeHistory.value = response.data.filter((source: ExternalSource) => source.sourceType === 'external');
    loading.value = false;
  }
};

const getReferenceInfoForClipboard = (item: ExternalSource): string => {
  const schemaLabelName = labelTools.value.schemaLabelsMap[currentProject.value.schemaLabel]
    ? labelTools.value.schemaLabelsMap[currentProject.value.schemaLabel].name
    : 'UNKNOWN_LABEL';
  const policyLabelNames = currentProject.value.policyLabels
    .map((l: string) =>
      labelTools.value.policyLabelsMap[l] ? labelTools.value.policyLabelsMap[l].name : 'UNKNOWN_LABEL',
    )
    .join(', ');
  if (item.sourceType === 'upload') {
    return `Disclosure Portal CCS Reference

Project Name: ${currentProject.value.name}
Project Identifier: ${currentProject.value._key}
Project Schema Label: ${schemaLabelName}
Project Policy Labels: ${policyLabelNames}
Project Version: ${version.value.name}
Version Identifier:  ${version.value._key}
Reference Timestamp: ${formatDateAndTime(dayjs().toISOString())} (UTC)
CCS Name: ${item.url}
Origin: ${item.origin}
Uploader: ${item.uploader}
Upload Date: ${formatDateTimeShort(item.created, true)} (UTC)
CCS SHA-256: ${item.hash}
Deliveries Link: ${window.location.href}`;
  } else {
    return `Disclosure Portal Resource Reference

Project Name: ${currentProject.value.name}
Project Identifier: ${currentProject.value._key}
Project Schema Label: ${schemaLabelName}
Project Policy Labels: ${policyLabelNames}
Project Version: ${version.value.name}
Version Identifier:  ${version.value._key}
Reference Timestamp: ${formatDateAndTime(dayjs().toISOString())} (UTC)
Resource: ${item.url}
Resource Description: ${item.comment}
Origin: ${item.origin}
Uploader: ${item.uploader}
Upload Date: ${formatDateTimeShort(item.created, true)} (UTC)`;
  }
};

const showDeletionConfirmationDialog = (sourceId: string, sourceUrl: string) => {
  confirmConfig.value = {
    type: ConfirmationType.DELETE,
    key: sourceId,
    name: sourceUrl,
    description: 'DLG_CONFIRMATION_DESCRIPTION',
    okButton: 'Btn_delete',
  };
  (confirmationDialog.value as any)?.makeVisible();
};
const showCreateExternalSourceDialog = () => {
  (dlgExternalSource.value as any)?.open(currentProject.value._key, version.value._key);
};

const showEditExternalSourceDialog = (item: ExternalSource) => {
  (dlgExternalSource.value as any)?.edit(currentProject.value._key, version.value._key, item);
};

const doDeleteExternalSource = async (config: IConfirmationDialogConfig) => {
  await VersionService.deleteExternalSource(config.key, currentProject.value._key, version.value._key);
  snackbar.info(t('DIALOG_source_code_delete_success'));
  await reload();
};

const copySourceToClipboard = (item: ExternalSource) => {
  const content = getReferenceInfoForClipboard(item);
  copyToClipboard(content);
};

const getActionButtons = (item: ExternalSource): TableActionButtonsProps['buttons'] => {
  const isDeprecated = currentProject.value.isDeprecated;

  return [
    {
      icon: 'mdi-content-copy',
      hint: t('TT_COPY_REFERENCE_INFO'),
      event: 'copy',
    },
    {
      icon: 'mdi-pencil',
      hint: t('BTN_EDIT'),
      event: 'edit',
      disabled: isDeprecated,
    },
    {
      icon: 'mdi-delete',
      hint: t('Btn_delete'),
      event: 'delete',
      disabled: isDeprecated,
    },
  ];
};

onUnmounted(async () => {
  eventBus.off('window-resize', updateTableHeight);
});

onMounted(async () => {
  eventBus.on('window-resize', updateTableHeight);
  await updateTableHeight();
  await reload();
});
</script>

<template>
  <TableLayout has-title has-tab>
    <template #description>
      <span class="d-text" v-html="t('SOURCE_CODE_DISCLAIMER_TEXT')"></span>
    </template>
    <template #buttons>
      <span class="d-headline-2">{{ t('TITLE_CCS_RESOURCES') }}</span>

      <DCActionButton
        v-if="currentProject.accessRights?.allowCCSAction.upload"
        large
        :text="t('BTN_ADD')"
        class="text-none"
        icon="mdi-plus"
        :hint="t('TT_ADD_SOURCE')"
        @click="showCreateExternalSourceDialog" />
    </template>
    <template #table>
      <div ref="scGrid" class="fill-height">
        <v-data-table
          density="compact"
          class="striped-table fill-height"
          :headers="headersSources"
          fixed-header
          item-key="_key"
          :height="tableHeight"
          :sort-by="sortBy"
          :search="search"
          sort-desc
          :items="sourceCodeHistory"
          :footer-props="{
            'items-per-page-options': [10, 50, 100, -1],
          }">
          <template v-slot:item.url="{item}">
            <v-tooltip v-if="!item.url.startsWith('file://')" :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" location="bottom">
              <template v-slot:activator="{props}">
                <span v-bind="props">
                  <DExternalLink :text="getStrWithMaxLength(45, item.url)" :url="item.url"></DExternalLink>
                </span>
              </template>
              <span>{{ t('OPEN_URL_EXTERN') }} {{ item.url }}</span>
            </v-tooltip>
            <span v-else>{{ item.url }}</span>
          </template>
          <template v-slot:item.created="{item}">
            <DDateCellWithTooltip :value="'' + item.created" />
          </template>
          <template v-slot:item.fileSize="{item}">
            {{ item.fileSize }}
          </template>
          <template v-slot:[`item.actions`]="{item}">
            <TableActionButtons
              variant="compact"
              :buttons="getActionButtons(item)"
              @copy="copySourceToClipboard(item)"
              @edit="showEditExternalSourceDialog(item)"
              @delete="showDeletionConfirmationDialog(item._key, item.url)" />
          </template>
        </v-data-table>
      </div>
    </template>
  </TableLayout>

  <ConfirmationDialog ref="confirmationDialog" @confirm="doDeleteExternalSource" :config="confirmConfig" />
  <NewExternalSourceDialog ref="dlgExternalSource" @reload="reload" />
</template>
