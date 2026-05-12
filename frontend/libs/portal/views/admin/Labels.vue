<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import useDimensions from '@disclosure-portal/composables/useDimensions';
import Icons from '@disclosure-portal/constants/icons';
import Label from '@disclosure-portal/model/Label';
import {Rights} from '@shared/types/Rights';
import AdminService from '@disclosure-portal/services/admin';
import {useUserStore} from '@disclosure-portal/stores/user';
import {downloadFile} from '@disclosure-portal/utils/download';
import eventBus from '@shared/utils/eventbus';
import {formatDateAndTime, getCssClassForTableRow} from '@disclosure-portal/utils/Table';
import useSnackbar from '@shared/composables/useSnackbar';
import {useTabsWindows} from '@shared/composables/useTabsWindows';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {DataTableHeader, SortItem} from '@shared/types/table';
import dayjs from 'dayjs';
import {computed, nextTick, onMounted, onUnmounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const {calculateHeight} = useDimensions();
const {info} = useSnackbar();
const {dashboardCrumbs, ...breadcrumbs} = useBreadcrumbsStore();
const {tabUrl, selectedTab} = useTabsWindows('/dashboard/admin/labels', ['schema', 'policy', 'project']);
const userStore = useUserStore();

const dialogNewLabelVisible = ref(false);
const schemaLabels = ref<Label[]>([]);
const schemaLabelSearch = ref('');
const policyLabels = ref<Label[]>([]);
const policyLabelSearch = ref('');
const projectLabels = ref<Label[]>([]);
const projectLabelSearch = ref('');
const rights = ref<Rights>({} as Rights);
const icons = Icons;

const headers = computed<DataTableHeader[]>(() => [
  {
    title: t('COL_ACTIONS'),
    align: 'center',
    width: 120,
    value: 'actions',
  },
  {
    title: t('CD_NAME'),
    align: 'start',
    value: 'name',
    width: 200,
    sortable: true,
  },
  {
    title: t('DESCRIPTION'),
    width: 200,
    align: 'start',
    value: 'description',
  },
  {
    title: t('CREATED'),
    align: 'center',
    width: 120,
    value: 'created',
    sortable: true,
  },
]);

const customFilterTable = (value: string | null, search: string) => {
  if (value) {
    const dateTime = formatDateAndTime(value);
    if (dateTime && dateTime !== 'Invalid date') {
      return dateTime.indexOf(search) > -1;
    }
    return value.toLowerCase().indexOf(search.toLowerCase()) > -1;
  }
  return false;
};

const deleteLabel = async (config: IConfirmationDialogConfig) => {
  await AdminService.deleteLabel(config.key);
  info(t('DIALOG_LABEL_DELETE_SUCCESS'));
  await reloadAll();
};

const confirmVisible = ref(false);
const confirmConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const showConfirm = (item: Label) => {
  confirmConfig.value = {
    type: ConfirmationType.DELETE,
    key: item._key,
    name: item.name,
    okButtonIsDisabled: false,
    okButton: 'BTN_DELETE',
    description: 'DLG_CONFIRMATION_DESCRIPTION',
  } as IConfirmationDialogConfig;
  confirmVisible.value = true;
};

const reloadAll = async () => {
  const labels = (await AdminService.getLabels()).data;
  policyLabels.value = [];
  schemaLabels.value = [];
  projectLabels.value = [];
  for (const label of labels) {
    if (label.type === 'SCHEMA') {
      schemaLabels.value.push(label);
    } else if (label.type === 'PROJECT') {
      projectLabels.value.push(label);
    } else {
      policyLabels.value.push(label);
    }
  }
};

const downloadCsv = async () => {
  downloadFile(
    'policies_and_labels_' + dayjs(new Date()).format('YYYY-MM-DD_hh_mm_ss') + '.csv',
    AdminService.downloadPLcsv(),
    true,
  );
};

const transferData = ref<Label | undefined>();
const formMode = ref<'create' | 'edit'>('create');
const showEditLabelDialog = (editData: Label) => {
  formMode.value = 'edit';
  dialogNewLabelVisible.value = true;
  transferData.value = editData;
};

const labelType = ref<'SCHEMA' | 'POLICY' | 'PROJECT'>('SCHEMA');
const showCreateSchemaLabelDialog = () => {
  formMode.value = 'create';
  dialogNewLabelVisible.value = true;
  labelType.value = 'SCHEMA';
};

const showCreatePolicyLabelDialog = () => {
  formMode.value = 'create';
  dialogNewLabelVisible.value = true;
  labelType.value = 'POLICY';
};

const showCreateProjectLabelDialog = () => {
  formMode.value = 'create';
  dialogNewLabelVisible.value = true;
  labelType.value = 'PROJECT';
};

const sortItems = [{key: 'name', order: 'asc'} as SortItem];

const labelButtons = computed(() => {
  const editHints = {
    schema: t('TT_edit_schema_label'),
    policy: t('TT_edit_policy_label'),
    project: t('TT_edit_project_label'),
  };
  const deleteHints = {
    schema: t('TT_delete_schema_label'),
    policy: t('TT_delete_policy_label'),
    project: t('TT_delete_project_label'),
  };
  const tab = selectedTab.value as 'schema' | 'policy' | 'project';
  return [
    {
      icon: 'mdi-pencil',
      event: 'edit',
      hint: editHints[tab],
      show: !!(rights.value.allowLabel && rights.value.allowLabel.update),
    },
    {
      icon: 'mdi-delete',
      event: 'delete',
      hint: deleteHints[tab],
      show: !!(rights.value.allowLabel && rights.value.allowLabel.delete),
    },
  ];
});

const initBreadcrumbs = () => {
  breadcrumbs.setCurrentBreadcrumbs([
    ...dashboardCrumbs,
    {
      title: t('BC_Labels'),
      href: '/dashboard/admin/labels',
    },
  ]);
};

onMounted(() => {
  rights.value = userStore.getRights;
  initBreadcrumbs();
  reloadAll();

  updateTableHeight();
  eventBus.on('window-resize', updateTableHeight);
});

onUnmounted(() => {
  eventBus.off('window-resize', updateTableHeight);
});

const tableHeight = ref(0);
const dataTableAsElement = ref<HTMLElement | null>(null);

const updateTableHeight = () => {
  nextTick(() => {
    tableHeight.value = calculateHeight(dataTableAsElement.value, false, true);
  });
};
</script>

<template>
  <v-container fluid class="h-full px-6">
    <div class="d-flex align-center pb-3">
      <h1 class="text-h5">{{ t('Labels') }}</h1>
    </div>
    <v-row expand class="align-center">
      <v-col>
        <v-card>
          <v-tabs v-model="selectedTab" slider-color="brand" active-class="active" show-arrows bg-color="tabsHeader">
            <v-tab value="schema" :to="tabUrl.schema">
              {{ t('TAB_SCHEMA_LABELS') }}
            </v-tab>
            <v-tab value="policy" :to="tabUrl.policy">
              {{ t('TAB_POLICY_LABELS') }}
            </v-tab>
            <v-tab value="project" :to="tabUrl.project">
              {{ t('TAB_PROJECT_LABELS') }}
            </v-tab>
          </v-tabs>
          <v-tabs-window v-model="selectedTab">
            <v-tabs-window-item value="schema">
              <TableLayout has-title has-tab>
                <template #buttons>
                  <span class="text-h6">{{ t('TITLE_SCHEMA_LABELS') }}</span>
                  <DCActionButton
                    v-if="rights.allowLabel && rights.allowLabel.create"
                    large
                    :text="t('BTN_ADD')"
                    class="mx-2"
                    icon="mdi-plus"
                    :hint="t('TT_add_schema_label')"
                    @click="showCreateSchemaLabelDialog" />
                  <v-spacer></v-spacer>
                  <DSearchField v-model="schemaLabelSearch" />
                </template>
                <template #table>
                  <v-data-table
                    class="striped-table fill-height"
                    density="compact"
                    fixed-header
                    :headers="headers"
                    :sort-by="sortItems"
                    :items="schemaLabels"
                    :items-per-page="-1"
                    :custom-filter="customFilterTable"
                    :item-class="getCssClassForTableRow"
                    :search="schemaLabelSearch">
                    <template #[`item.name`]="{item}">
                      <DLabel v-if="item.name" :labelName="item.name" :iconName="icons.SCHEMA" class="mt-1"></DLabel>
                    </template>

                    <template #[`item.description`]="{item}">
                      <Truncated>{{ item.description }}</Truncated>
                    </template>

                    <template #[`item.created`]="{item}">
                      <DDateCellWithTooltip :value="item.created"></DDateCellWithTooltip>
                    </template>

                    <template #[`item.actions`]="{item}">
                      <TableActionButtons
                        :buttons="labelButtons"
                        @edit="showEditLabelDialog(item)"
                        @delete="showConfirm(item)" />
                    </template>
                  </v-data-table>
                </template>
              </TableLayout>
            </v-tabs-window-item>
            <v-tabs-window-item value="policy">
              <TableLayout has-title has-tab>
                <template #buttons>
                  <span class="text-h6">{{ t('TITLE_POLICY_LABELS') }}</span>
                  <DCActionButton
                    large
                    :text="t('BTN_ADD')"
                    class="mx-2"
                    icon="mdi-plus"
                    :hint="t('TT_add_policy_label')"
                    @click="showCreatePolicyLabelDialog"
                    v-if="rights.allowLabel && rights.allowLabel.create" />
                  <v-spacer></v-spacer>

                  <DCActionButton
                    large
                    icon="mdi-download"
                    :text="t('BTN_DOWNLOAD')"
                    :hint="t('TT_download_label_csv')"
                    @click="downloadCsv" />
                  <DSearchField v-model="policyLabelSearch" />
                </template>
                <template #table>
                  <v-data-table
                    class="striped-table fill-height"
                    density="compact"
                    fixed-header
                    :headers="headers"
                    :sort-by="sortItems"
                    :custom-filter="customFilterTable"
                    :items-per-page="-1"
                    :items="policyLabels"
                    :search="policyLabelSearch"
                    :item-class="getCssClassForTableRow"
                    :height="tableHeight">
                    <template #[`item.name`]="{item}">
                      <DLabel v-if="item.name" :iconName="icons.POLICY" :labelName="item.name" class="mt-1"></DLabel>
                    </template>

                    <template #[`item.created`]="{item}">
                      <DDateCellWithTooltip :value="item.created"></DDateCellWithTooltip>
                    </template>

                    <template #[`item.actions`]="{item}">
                      <TableActionButtons
                        :buttons="labelButtons"
                        @edit="showEditLabelDialog(item)"
                        @delete="showConfirm(item)" />
                    </template>

                    <template #[`item.description`]="{item}">
                      <Truncated>{{ item.description }}</Truncated>
                    </template>
                  </v-data-table>
                </template>
              </TableLayout>
            </v-tabs-window-item>
            <v-tabs-window-item value="project">
              <TableLayout has-title has-tab>
                <template #buttons>
                  <span class="text-h6">{{ t('TITLE_PROJECT_LABELS') }}</span>
                  <DCActionButton
                    v-if="rights.allowLabel && rights.allowLabel.create"
                    large
                    :text="t('BTN_ADD')"
                    class="mx-2"
                    icon="mdi-plus"
                    :hint="t('TT_add_project_label')"
                    @click="showCreateProjectLabelDialog" />
                  <v-spacer></v-spacer>
                  <DCActionButton
                    :text="t('BTN_DOWNLOAD')"
                    large
                    icon="mdi-download"
                    :hint="t('TT_download_label_csv')"
                    @click="downloadCsv" />
                  <DSearchField v-model="projectLabelSearch" />
                </template>
                <template #table>
                  <v-data-table
                    class="striped-table fill-height"
                    density="compact"
                    fixed-header
                    :headers="headers"
                    :sort-by="sortItems"
                    :custom-filter="customFilterTable"
                    :items-per-page="-1"
                    :items="projectLabels"
                    :search="projectLabelSearch"
                    :item-class="getCssClassForTableRow"
                    :height="tableHeight">
                    <template #[`item.name`]="{item}">
                      <DLabel
                        v-if="item.name"
                        :iconName="icons.PROJECT_LABEL"
                        :labelName="item.name"
                        class="mt-1"></DLabel>
                    </template>

                    <template #[`item.created`]="{item}">
                      <DDateCellWithTooltip :value="item.created"></DDateCellWithTooltip>
                    </template>

                    <template #[`item.actions`]="{item}">
                      <TableActionButtons
                        :buttons="labelButtons"
                        @edit="showEditLabelDialog(item)"
                        @delete="showConfirm(item)" />
                    </template>

                    <template #[`item.description`]="{item}">
                      <Truncated>{{ item.description }}</Truncated>
                    </template>
                  </v-data-table>
                </template>
              </TableLayout>
            </v-tabs-window-item>
          </v-tabs-window>
        </v-card>
        <NewLabelDialog
          :is-open="dialogNewLabelVisible"
          @update:is-open="dialogNewLabelVisible = $event"
          :mode="formMode"
          :initial-data="transferData"
          :type="labelType"
          @reload="reloadAll"></NewLabelDialog>

        <ConfirmationDialog
          v-model:showDialog="confirmVisible"
          :config="confirmConfig"
          @confirm="deleteLabel"></ConfirmationDialog>
      </v-col>
    </v-row>
  </v-container>
</template>
