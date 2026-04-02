<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import {useView} from '@disclosure-portal/composables/useView';
import IObligation from '@disclosure-portal/model/IObligation';
import {compareLevel} from '@disclosure-portal/model/Quality';
import {Rights} from '@disclosure-portal/model/Rights';
import AdminService from '@disclosure-portal/services/admin';
import licenseService from '@disclosure-portal/services/license';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useUserStore} from '@disclosure-portal/stores/user';
import {downloadFile} from '@disclosure-portal/utils/download';
import {RightsUtils} from '@disclosure-portal/utils/Rights';
import {getCssClassForTableRow} from '@disclosure-portal/utils/Table';
import useViewTools, {getIconColorOfLevel, getIconOfLevel} from '@disclosure-portal/utils/View';
import {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {DataTableHeader, DataTableHeaderFilterItems, DataTableItem, SortItem} from '@shared/types/table';
import dayjs from 'dayjs';
import _, {indexOf} from 'lodash';
import {computed, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const viewTools = useViewTools();
const {dashboardCrumbs, ...breadcrumbs} = useBreadcrumbsStore();
const userStore = useUserStore();
const snackbar = useSnackbar();
const {getTextOfLevel, getTextOfType} = useView();
const appStore = useAppStore();

const items = ref<IObligation[]>([]);
const search = ref('');
const classificationDlg = ref();
const auditDialogOpen = ref<((key: string, name: string) => void) | null>(null);

const headers = computed<DataTableHeader[]>(() => {
  return [
    {
      title: t('COL_ACTIONS'),
      align: 'center',
      width: 80,
      value: 'actions',
      sortable: false,
    },
    {
      title: t('COL_TYPE'),
      align: 'start',
      sortable: true,
      width: 140,
      value: 'type',
    },
    {
      title: t('COL_WARN_LEVEL'),
      align: 'start',
      width: 150,
      value: 'warnLevel',
      sortable: true,
      sort: compareLevel,
    },
    {
      title: t('COL_SHORT_NAME'),
      width: 280,
      align: 'start',
      value: 'name',
      sortable: true,
    },
    {
      title: t('COL_DESCRIPTION'),
      width: 280,
      align: 'start',
      value: 'description',
      sortable: false,
    },
    {
      title: t('COL_CREATED'),
      value: 'created',
      width: 115,
      sortable: true,
    },
    {
      title: t('COL_UPDATED'),
      width: 115,
      sortable: true,
      value: 'updated',
    },
  ];
});
const rights = ref<Rights>(new Rights());
const selectedFilterTypes = ref<string[]>([]);
const possibleTypes = computed((): DataTableHeaderFilterItems[] => {
  const uniqueTypes = [...new Set(items.value.map((item: IObligation) => item.type))];
  return uniqueTypes.map((item: string) => ({
    value: item,
  }));
});
const selectedFilterLevel = ref<string[]>([]);
const possibleLevel = computed((): DataTableHeaderFilterItems[] => {
  const uniqueLevels = [...new Set(items.value.map((item: IObligation) => item.warnLevel))];
  return uniqueLevels.map((item: string) => ({
    text: getTextOfLevel(item),
    value: item,
    iconColor: getIconColorOfLevel(item),
    icon: getIconOfLevel(item),
  }));
});

const confirmConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const confirmVisible = ref(false);

onMounted(async () => {
  rights.value = userStore.getRights;
  initBreadcrumbs();
  await reload();
});

const sortItem = (): SortItem[] => [{key: 'type', order: 'asc'}];

const initBreadcrumbs = () => {
  breadcrumbs.setCurrentBreadcrumbs([
    ...dashboardCrumbs,
    {
      title: t('BC_OBLIGATION'),
      href: '/dashboard/admin/obligations',
    },
  ]);
};

const doDelete = async (config: IConfirmationDialogConfig) => {
  if (config.okButtonIsDisabled) return;
  await AdminService.deleteObligation(config.key);
  snackbar.info(t('DIALOG_classification_delete_success'));
  await reload();
};

const showDeletionConfirmationDialog = async (item: IObligation) => {
  await licenseService.getCountOfLicencesUsingThisObligation(item._key).then((r) => {
    const count = r.data.count;
    if (count > 0) {
      const userHasRightsToDeleteEvenIsUsed = false;
      if (userHasRightsToDeleteEvenIsUsed) {
        confirmConfig.value = {
          key: item._key,
          name: item.name,
          okButtonIsDisabled: false,
          extendedDetails: '' + t('OBLIGATION_IN_USE_BY_LICENSES', {count: count}),
          okButton: 'Btn_delete',
          description: 'DLG_CONFIRMATION_DESCRIPTION',
        } as IConfirmationDialogConfig;
      } else {
        confirmConfig.value = {
          title: 'DLG_WARNING_TITLE',
          key: item._key,
          name: item.name,
          okButtonIsDisabled: true,
          extendedDetails: '' + t('OBLIGATION_IN_USE_BY_LICENSES', {count: count}),
          okButton: 'Btn_delete',
          description: 'DLG_CAN_NOT_DELETE_IN_USE',
        } as IConfirmationDialogConfig;
      }
    } else {
      confirmConfig.value = {
        key: item._key,
        name: item.name,
        okButtonIsDisabled: false,
        okButton: 'Btn_delete',
        description: 'DLG_CONFIRMATION_DESCRIPTION',
      } as IConfirmationDialogConfig;
    }
    confirmVisible.value = true;
  });
};

const filterOnType = (item: IObligation): boolean => {
  return selectedFilterTypes.value.length === 0 || indexOf(selectedFilterTypes.value, item.type) !== -1;
};

const filterOnLevel = (item: IObligation): boolean => {
  return selectedFilterLevel.value.length === 0 || indexOf(selectedFilterLevel.value, item.warnLevel) !== -1;
};
const openDialog = (item?: IObligation) => {
  classificationDlg.value?.open(item);
};

const onRowClick = (event: Event, item: DataTableItem<IObligation>) => {
  openDialog(item.item);
};
const filteredList = computed(() => {
  return _.chain(items.value).filter(filterOnType).filter(filterOnLevel).value();
});

const reload = async () => {
  const response = (await AdminService.getAllObligations()).data;
  items.value = response.items;
};

const downloadCsv = async () => {
  downloadFile(
    `licenses_and_classifications_${dayjs(new Date()).format('YYYY-MM-DD_hh_mm_ss')}.csv`,
    AdminService.downloadLCcsv(),
    true,
  );
};

const getActionButtons = (_: IObligation): TableActionButtonsProps['buttons'] => {
  return [
    {
      icon: 'mdi-pencil',
      hint: t('TT_edit_classification'),
      event: 'edit',
      show: rights.value?.allowObligation?.update,
    },
    {
      icon: 'mdi-note-text-outline',
      hint: t('TT_AUDIT_LOG'),
      event: 'audit',
      show: RightsUtils.isLicenseManager() || RightsUtils.isDomainAdmin(),
    },
    {
      icon: 'mdi-delete',
      hint: t('TT_delete_classification'),
      event: 'delete',
      show: rights.value?.allowObligation?.delete,
    },
  ];
};
</script>

<template>
  <TableLayout>
    <template #buttons>
      <h1 class="d-headline">{{ t('CLASSIFICATIONS') }}</h1>
      <DCActionButton
        large
        :text="t('BTN_ADD')"
        icon="mdi-plus"
        :hint="t('TT_ADD_CLASSIFICATION')"
        @click="openDialog()"
        v-if="rights.allowObligation?.create" />
      <v-spacer></v-spacer>
      <DCActionButton
        :text="t('BTN_DOWNLOAD')"
        large
        icon="mdi-download"
        :hint="t('TT_download_license_csv')"
        @click="downloadCsv" />
      <v-text-field
        autocomplete="off"
        :max-width="500"
        v-model="search"
        append-inner-icon="mdi-magnify"
        :label="t('labelSearch')"
        variant="outlined"
        clearable
        density="compact"
        hide-details="auto"></v-text-field>
    </template>
    <template #table>
      <div ref="tableGridClassifications" class="fill-height">
        <v-data-table
          density="compact"
          class="striped-table fill-height"
          :headers="headers"
          fixed-header
          @click:row="onRowClick"
          :items-per-page="-1"
          :footer-props="{'items-per-page-options': [10, 50, 100, -1]}"
          :items="filteredList"
          :sort-by="sortItem()"
          :item-class="getCssClassForTableRow"
          :search="search">
          <template #[`header.type`]="{column, getSortIcon, toggleSort}">
            <span class="mr-1">{{ column.title }}</span>
            <GridHeaderFilterIcon
              v-model="selectedFilterTypes"
              :column="column"
              :label="t('COL_TYPE')"
              :allItems="possibleTypes">
            </GridHeaderFilterIcon>
            <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
          </template>

          <template #[`header.warnLevel`]="{column, getSortIcon, toggleSort}">
            <span class="mr-1">{{ column.title }}</span>
            <GridHeaderFilterIcon
              v-model="selectedFilterLevel"
              :column="column"
              :label="t('COL_WARN_LEVEL')"
              :allItems="possibleLevel">
            </GridHeaderFilterIcon>
            <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
          </template>

          <template #[`item.type`]="{item}">
            {{ getTextOfType(item.type) }}
          </template>

          <template #[`item.name`]="{item}">
            {{ appStore.appLanguage === 'en' ? item.name : item.nameDe }}
          </template>

          <template #[`item.description`]="{item}">
            <Truncated>{{ appStore.appLanguage === 'en' ? item.description : item.descriptionDe }}</Truncated>
          </template>

          <template #[`item.created`]="{item}">
            <DDateCellWithTooltip :value="item.created" />
          </template>

          <template #[`item.updated`]="{item}">
            <DDateCellWithTooltip :value="item.updated" />
          </template>

          <template #[`item.warnLevel`]="{item}">
            <v-icon :color="getIconColorOfLevel(item.warnLevel)">{{ getIconOfLevel(item.warnLevel) }}</v-icon>
            <Tooltip location="bottom">{{ getTextOfLevel(item.warnLevel) }}</Tooltip>
          </template>

          <template #[`item.actions`]="{item}">
            <AuditDialog v-slot="{open}">
              <template v-if="!auditDialogOpen">
                {{ ((auditDialogOpen = open), '') }}
              </template>
              <TableActionButtons
                variant="compact"
                :buttons="getActionButtons(item)"
                @edit="openDialog(item)"
                @audit="open(item._key, viewTools.getNameForLanguage(item))"
                @delete="showDeletionConfirmationDialog(item)" />
            </AuditDialog>
          </template>
        </v-data-table>
      </div>
    </template>
  </TableLayout>

  <ConfirmationDialog v-model:showDialog="confirmVisible" :config="confirmConfig" @confirm="doDelete" />
  <NewClassificationDialog ref="classificationDlg" @reload="reload" />
</template>
