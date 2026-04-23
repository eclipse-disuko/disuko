<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script lang="ts" setup>
import ComponentDetailsDialog from '@disclosure-portal/components/dialog/ComponentDetailsDialog.vue';
import {ScanRemark, ScanRemarkLevel} from '@disclosure-portal/model/Quality';
import ProjectService, {RemarkTypes} from '@disclosure-portal/services/projects';
import VersionService from '@disclosure-portal/services/version';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import {downloadFile} from '@disclosure-portal/utils/download';
import {
  getIconColorScanRemarkLevel,
  getScanRemarkStatusSortIndex,
  getStrWithMaxLength,
} from '@disclosure-portal/utils/View';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import DCopyClipboardButton from '@shared/components/disco/DCopyClipboardButton.vue';
import {DataTableHeader, DataTableHeaderFilterItems, DataTableItem, SortItem} from '@shared/types/table';
import _ from 'lodash';
import {computed, onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute} from 'vue-router';

const projectStore = useProjectStore();
const sbomStore = useSbomStore();
const {t} = useI18n();

const search = ref('');
const sortBy = ref<SortItem[]>([{key: 'status', order: 'desc'}]);
const dataAreLoaded = ref(false);
const selectedFilterStatus = ref<ScanRemarkLevel[]>([]);
const selectedFilterQualityRemark = ref<string[]>([]);
const selectedFilterTypes = ref<string[]>([]);
const tableItems = ref<ScanRemark[]>([]);
const route = useRoute();
const newComponentDetailsDlg = ref<InstanceType<typeof ComponentDetailsDialog> | null>(null);

const possibleTypes = computed((): DataTableHeaderFilterItems[] => {
  if (!tableItems.value) {
    return [];
  }

  const uniqueTypes = [...new Set(tableItems.value.map(({type}) => type))];

  return uniqueTypes.map((types: string) => ({
    value: types,
  }));
});

const possibleRemarks = computed((): DataTableHeaderFilterItems[] => {
  if (!tableItems.value) {
    return [];
  }

  const uniqueRemarkKeys = [...new Set(tableItems.value.map(({remarkKey}) => remarkKey))];

  return uniqueRemarkKeys.map(
    (remarkKey: string) =>
      ({
        value: remarkKey,
        text: t(remarkKey),
      }) as DataTableHeaderFilterItems,
  );
});

const possibleStatuses = computed((): DataTableHeaderFilterItems[] => {
  if (!tableItems.value) {
    return [];
  }

  const uniqueStatuses = [...new Set(tableItems.value.map(({status}) => status))];

  return uniqueStatuses.map(
    (status: string) =>
      ({
        value: status,
        text: t('SCAN_REMARK_STATUS_' + status),
        icon: 'mdi-circle',
        iconColor: getIconColorScanRemarkLevel(status as ScanRemarkLevel),
      }) as DataTableHeaderFilterItems,
  );
});

const headers: DataTableHeader[] = [
  {
    title: t('COL_ACTIONS'),
    sortable: false,
    align: 'center',
    width: 60,
    value: 'Actions',
  },
  {
    title: t('COL_LEVEL'),
    width: 100,
    align: 'center',
    key: 'status',
    sortable: true,
  },
  {
    title: t('COL_QUALITY_REMARK'),
    width: 210,
    align: 'start',
    key: 'remarkKey',
    sortable: true,
  },
  {
    title: t('COL_COMPONENT_NAME'),
    width: 240,
    align: 'start',
    key: 'name',
    sortable: true,
  },
  {
    title: t('COL_COMPONENT_VERSION'),
    width: 80,
    align: 'start',
    value: 'version',
    sortable: true,
  },
  {
    title: t('COL_COMPONENT_TYPE'),
    align: 'start',
    value: 'type',
    width: 100,
    sortable: true,
  },
  {
    title: t('COL_DESCRIPTION'),
    align: 'start',
    value: 'descriptionKey',
    width: 320,
    sortable: true,
  },
];

const handleFilterQuery = () => {
  const filter = route.query.scanRemarkLevel as string;
  if (filter) {
    selectedFilterStatus.value = [filter as ScanRemarkLevel];
  }
};

const projectModel = computed(() => projectStore.currentProject!);
const version = computed(() => sbomStore.getCurrentVersion);
const spdx = computed(() => sbomStore.getSelectedSBOM);

const filteredList = computed(() => {
  return tableItems.value.filter((item: ScanRemark) => {
    return filterOnStatus(item) && filterOnRemark(item) && filterOnType(item);
  });
});

const reload = async () => {
  if (!spdx.value) {
    dataAreLoaded.value = true;
    return;
  }
  tableItems.value = await VersionService.getScanRemarksForSbom(
    projectModel.value._key,
    version.value._key,
    spdx.value._key,
  );
  dataAreLoaded.value = true;
};

const customFilterTable = (value: string, searchTerm: string) => {
  if (value != null && value) {
    const valueTranslated = t(value);
    return ('' + valueTranslated).toLowerCase().indexOf(searchTerm.toLowerCase()) > -1;
  }
  return false;
};

const showDetails = async (_: Event, row: DataTableItem<ScanRemark>) => {
  if (!row.item.spdxId) {
    return;
  }
  await ProjectService.getComponentDetailsForSbom(
    projectModel.value._key,
    version.value._key,
    spdx.value!._key,
    row.item.spdxId,
  ).then((response) => {
    if (newComponentDetailsDlg.value) {
      newComponentDetailsDlg.value?.open(response.data, '', row.item.policyRuleStatus, row.item.unmatchedLicenses);
    }
  });
};
const downloadScanRemarksCsv = async () => {
  downloadFile(
    projectModel.value.name + '_' + version.value.name + '_scan_remarks.csv',
    ProjectService.downloadScanOrLicenseRemarksForSbomCsv(
      projectModel.value._key,
      version.value._key,
      RemarkTypes.scan,
      spdx.value!._key,
    ),
    true,
  );
};

const filterOnRemark = (item: ScanRemark): boolean => {
  if (selectedFilterQualityRemark.value.length <= 0) {
    return true;
  }
  let found = false;
  selectedFilterQualityRemark.value.forEach((filter: string) => {
    if (!found && '' + item.remarkKey === filter) {
      found = true;
    }
  });
  return found;
};

const filterOnStatus = (item: ScanRemark): boolean => {
  if (selectedFilterStatus.value.length <= 0) {
    return true;
  }
  let found = false;
  selectedFilterStatus.value.forEach((filter: string) => {
    if (!found && '' + item.status === filter) {
      found = true;
    }
  });
  return found;
};

const filterOnType = (item: ScanRemark): boolean => {
  if (selectedFilterTypes.value.length > 0) {
    return selectedFilterTypes.value.some((filterType) => item.type === filterType);
  } else {
    return true;
  }
};

const customKeySort = {
  status: (a: ScanRemarkLevel, b: ScanRemarkLevel) => {
    const status1Index = getScanRemarkStatusSortIndex(a);
    const status2Index = getScanRemarkStatusSortIndex(b);
    return status2Index - status1Index;
  },
};

const getRemarkTextForClipboard = (item: ScanRemark): string => {
  return `${t(item.remarkKey)} in ${item.name} ${item.version}

Component:
${item.name}

Component Version:
${item.version}

Scan Remark:
${t(item.remarkKey)}

Additional information:
${t(item.descriptionKey)}
`;
};

watch(
  () => spdx.value,
  async () => {
    dataAreLoaded.value = false;
    await reload();
  },
);

watch(
  () => route.path,
  async (_newPath, _oldPath) => {
    handleFilterQuery();
  },
);

onMounted(async () => {
  handleFilterQuery();

  await reload();
});
</script>

<template>
  <div class="h-[calc(100%-56px)]">
    <Stack direction="row" class="pb-1">
      <DCActionButton
        :text="t('BTN_DOWNLOAD')"
        large
        icon="mdi-download"
        :hint="t('TT_download_scan_remarks')"
        @click="downloadScanRemarksCsv"
        class="pr-4" />
      <v-spacer></v-spacer>
      <DSearchField v-model="search" />
    </Stack>
    <v-data-table
      :items="filteredList"
      density="compact"
      class="striped-table my-0 h-full py-0"
      fixed-header
      :headers="headers"
      item-key="_key"
      :custom-key-sort="customKeySort"
      :sort-by.sync="sortBy"
      :loading="!dataAreLoaded"
      @click:row="showDetails"
      :search="search"
      :custom-filter="customFilterTable"
      :items-per-page="100"
      :footer-props="{
        'items-per-page-options': [10, 50, 100, 500],
      }">
      <template #[`header.status`]="{column, getSortIcon, toggleSort}">
        <GridFilterHeader :column="column" :getSortIcon="getSortIcon" :toggleSort="toggleSort">
          <template #filter>
            <GridHeaderFilterIcon
              v-model="selectedFilterStatus"
              :column="column"
              :label="t('COL_LEVEL')"
              :allItems="possibleStatuses">
            </GridHeaderFilterIcon>
          </template>
        </GridFilterHeader>
      </template>
      <template #[`header.remarkKey`]="{column, getSortIcon, toggleSort}">
        <GridFilterHeader :column="column" :getSortIcon="getSortIcon" :toggleSort="toggleSort">
          <template #filter>
            <GridHeaderFilterIcon
              v-model="selectedFilterQualityRemark"
              :column="column"
              :label="t('COL_QUALITY_REMARK')"
              :allItems="possibleRemarks">
            </GridHeaderFilterIcon>
          </template>
        </GridFilterHeader>
      </template>
      <template #[`header.type`]="{column, getSortIcon, toggleSort}">
        <GridFilterHeader :column="column" :getSortIcon="getSortIcon" :toggleSort="toggleSort">
          <template #filter>
            <GridHeaderFilterIcon
              v-model="selectedFilterTypes"
              :column="column"
              :label="t('COL_COMPONENT_TYPE')"
              :allItems="possibleTypes">
            </GridHeaderFilterIcon>
          </template>
        </GridFilterHeader>
      </template>

      <template #[`item.status`]="{item}">
        <div>
          <v-icon :color="getIconColorScanRemarkLevel(item.status)" x-small>mdi-circle</v-icon>
        </div>
      </template>
      <template #[`item.remarkKey`]="{item}">
        {{ t('' + item.remarkKey) }}
      </template>
      <template #[`item.descriptionKey`]="{item}">
        <span>
          {{ getStrWithMaxLength(120, t(item.descriptionKey)) }}
          <Tooltip>
            {{ t(item.descriptionKey) }}
          </Tooltip>
        </span>
      </template>
      <template #[`item.Actions`]="{item}">
        <DCopyClipboardButton
          tableButton="true"
          :hint="t('TT_COPY_SCAN_REMARKS')"
          :content="getRemarkTextForClipboard(item)" />
      </template>
    </v-data-table>

    <ComponentDetailsDialog ref="newComponentDetailsDlg" />
  </div>
</template>
