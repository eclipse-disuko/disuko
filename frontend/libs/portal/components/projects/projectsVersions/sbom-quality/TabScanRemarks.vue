// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

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
      <v-text-field
        autocomplete="off"
        variant="outlined"
        v-model="search"
        class="w-full md:max-w-[400px] md:w-auto"
        min-width="50"
        density="compact"
        append-inner-icon="mdi-magnify"
        :label="t('labelSearch')"
        clearable
        single-line
        hide-details></v-text-field>
    </Stack>
    <v-data-table
      :items="filteredList"
      density="compact"
      class="striped-table py-0 my-0 h-full"
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
      <template v-slot:header.status="{column, getSortIcon, toggleSort}">
        <div class="v-data-table-header__content">
          <span>{{ column.title }}</span>
          <v-menu :close-on-content-click="false" v-model="menu">
            <template v-slot:activator="{props}">
              <DIconButton
                :parentProps="props"
                icon="mdi-filter-variant"
                :hint="t('TT_SHOW_FILTER')"
                :color="selectedFilterStatus.length > 0 ? 'primary' : 'default'" />
            </template>
            <div class="bg-background" style="width: 280px">
              <v-row class="d-flex justify-end ma-1 mr-2">
                <DIconButton icon="mdi-close" @clicked="menu = false" color="default" />
              </v-row>
              <v-select
                v-model="selectedFilterStatus"
                :items="possibleStatuses"
                class="mx-2 pa-2 pb-4"
                :label="t('Lbl_filter_status')"
                clearable
                multiple
                item-title="text"
                item-value="value"
                variant="outlined"
                density="compact"
                menu
                transition="scale-transition"
                persistent-clear
                :list-props="{class: 'striped-filter-dd py-0'}">
                <template v-slot:item="{item, props}">
                  <v-list-item v-bind="props">
                    <template v-slot:prepend="{isSelected}">
                      <v-checkbox hide-details :model-value="isSelected" />
                    </template>
                    <template v-slot:title>
                      <v-icon class="mr-1" :color="getIconColorScanRemarkLevel(item.value)" small>mdi-circle</v-icon>
                      <span class="pFilterEntry">{{ item.value }}</span>
                    </template>
                  </v-list-item>
                </template>
                <template v-slot:selection="{item, index}">
                  <div v-if="index === 0" class="d-flex align-center">
                    <v-icon class="mr-1" :color="getIconColorScanRemarkLevel(item.value)" small>mdi-circle</v-icon>
                    <span class="pStatusFilter">{{ item.value }}</span>
                  </div>
                  <span v-if="index === 1" class="additionalFilter">
                    +{{ selectedFilterStatus.length - 1 }} others
                  </span>
                </template>
              </v-select>
            </div>
          </v-menu>
          <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
        </div>
      </template>
      <template v-slot:header.remarkKey="{column, getSortIcon, toggleSort}">
        <div class="v-data-table-header__content">
          <span>{{ column.title }}</span>
          <v-menu :close-on-content-click="false" v-model="menu2">
            <template v-slot:activator="{props}">
              <DIconButton
                :parentProps="props"
                icon="mdi-filter-variant"
                :hint="t('TT_SHOW_FILTER')"
                :color="selectedFilterQualityRemark.length > 0 ? 'primary' : 'default'" />
            </template>
            <div class="bg-background" style="width: 280px">
              <v-row class="d-flex justify-end ma-1 mr-2">
                <DIconButton icon="mdi-close" @clicked="menu2 = false" color="default" />
              </v-row>
              <v-select
                v-model="selectedFilterQualityRemark"
                :items="possibleRemarks"
                class="mx-2 pa-2 pb-4"
                :label="t('LABEL_FILTER_QUALITY_REMARK')"
                clearable
                multiple
                item-title="text"
                item-value="value"
                variant="outlined"
                density="compact"
                menu
                transition="scale-transition"
                persistent-clear
                :list-props="{class: 'striped-filter-dd py-0'}">
                <template v-slot:item="{item, props}">
                  <v-list-item v-bind="props" class="py-0 px-2">
                    <template v-slot:prepend="{isSelected}">
                      <v-checkbox hide-details :model-value="isSelected" />
                    </template>
                    <template v-slot:title>
                      <span class="pFilterEntry">{{ t(item.value) }}</span>
                    </template>
                  </v-list-item>
                </template>
                <template v-slot:selection="{item, index}">
                  <div v-if="index === 0" class="d-flex align-center">
                    <span class="pFilterEntry">{{ t(item.value) }}</span>
                  </div>
                  <span v-if="index === 1" class="pAdditionalFilter">
                    +{{ selectedFilterQualityRemark.length - 1 }} others
                  </span>
                </template>
              </v-select>
            </div>
          </v-menu>
          <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
        </div>
      </template>
      <template v-slot:header.type="{column, getSortIcon, toggleSort}">
        <div class="v-data-table-header__content">
          <span>{{ column.title }}</span>
          <v-menu :close-on-content-click="false" v-model="menu3">
            <template v-slot:activator="{props}">
              <DIconButton
                :parentProps="props"
                icon="mdi-filter-variant"
                :hint="t('TT_SHOW_FILTER')"
                :color="selectedFilterTypes.length > 0 ? 'primary' : 'default'" />
            </template>
            <div class="bg-background" style="width: 280px">
              <v-row class="d-flex justify-end ma-1 mr-2">
                <DIconButton icon="mdi-close" @clicked="menu3 = false" color="default" />
              </v-row>
              <v-select
                v-model="selectedFilterTypes"
                :items="possibleTypes"
                class="mx-2 pa-2 pb-4"
                :label="t('Lbl_filter_on_type')"
                clearable
                multiple
                item-title="text"
                item-value="value"
                variant="outlined"
                density="compact"
                menu
                transition="scale-transition"
                persistent-clear
                :list-props="{class: 'striped-filter-dd py-0'}">
                <template v-slot:item="{item, props}">
                  <v-list-item v-bind="props" class="py-0 px-2">
                    <template v-slot:prepend="{isSelected}">
                      <v-checkbox hide-details :model-value="isSelected" />
                    </template>
                    <template v-slot:title>
                      <span class="pFilterEntry">{{ t(item.value) }}</span>
                    </template>
                  </v-list-item>
                </template>
                <template v-slot:selection="{item, index}">
                  <div v-if="index === 0" class="d-flex align-center">
                    <span class="pFilterEntry">{{ t(item.value) }}</span>
                  </div>
                  <span v-if="index === 1" class="pAdditionalFilter">
                    +{{ selectedFilterTypes.length - 1 }} others
                  </span>
                </template>
              </v-select>
            </div>
          </v-menu>
          <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
        </div>
      </template>

      <template v-slot:item.status="{item}">
        <div>
          <v-icon :color="getIconColorScanRemarkLevel(item.status)" x-small>mdi-circle</v-icon>
        </div>
      </template>
      <template v-slot:item.remarkKey="{item}">
        {{ t('' + item.remarkKey) }}
      </template>
      <template v-slot:item.descriptionKey="{item}">
        <span>
          {{ getStrWithMaxLength(120, t(item.descriptionKey)) }}
          <Tooltip>
            {{ t(item.descriptionKey) }}
          </Tooltip>
        </span>
      </template>
      <template v-slot:item.Actions="{item}">
        <DCopyClipboardButton
          tableButton="true"
          :hint="t('TT_COPY_SCAN_REMARKS')"
          :content="getRemarkTextForClipboard(item)" />
      </template>
    </v-data-table>

    <ComponentDetailsDialog ref="newComponentDetailsDlg" />
  </div>
</template>

<script lang="ts" setup>
import ComponentDetailsDialog from '@disclosure-portal/components/dialog/ComponentDetailsDialog.vue';
import {IDefaultSelectItem} from '@disclosure-portal/model/IObligation';
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
import {DataTableHeader, DataTableItem, SortItem} from '@shared/types/table';
import _ from 'lodash';
import {computed, onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute} from 'vue-router';

const projectStore = useProjectStore();
const sbomStore = useSbomStore();
const {t} = useI18n();

const search = ref('');
const headers: DataTableHeader[] = [
  {
    title: t('COL_ACTIONS'),
    sortable: false,
    align: 'center',
    width: 60,
    class: 'tableHeaderCell',
    value: 'Actions',
  },
  {
    title: t('COL_LEVEL'),
    width: 100,
    align: 'center',
    class: 'tableHeaderCell',
    key: 'status',
    sortable: true,
  },
  {
    title: t('COL_QUALITY_REMARK'),
    width: 210,
    align: 'start',
    class: 'tableHeaderCell',
    key: 'remarkKey',
    sortable: true,
  },
  {
    title: t('COL_COMPONENT_NAME'),
    width: 240,
    align: 'start',
    class: 'tableHeaderCell',
    key: 'name',
    sortable: true,
  },
  {
    title: t('COL_COMPONENT_VERSION'),
    width: 80,
    align: 'start',
    class: 'tableHeaderCell',
    value: 'version',
    sortable: true,
  },
  {
    title: t('COL_COMPONENT_TYPE'),
    align: 'start',
    class: 'tableHeaderCell',
    value: 'type',
    width: 100,
    sortable: true,
  },
  {
    title: t('COL_DESCRIPTION'),
    align: 'start',
    class: 'tableHeaderCell',
    value: 'descriptionKey',
    width: 320,
    sortable: true,
  },
];
const sortBy = ref<SortItem[]>([{key: 'status', order: 'desc'}]);
const dataAreLoaded = ref(false);
const selectedFilterStatus = ref<ScanRemarkLevel[]>([]);
const selectedFilterQualityRemark = ref<string[]>([]);
const selectedFilterTypes = ref<string[]>([]);
const menu = ref(false);
const menu2 = ref(false);
const menu3 = ref(false);
const tableItems = ref<ScanRemark[]>([]);
const route = useRoute();
const newComponentDetailsDlg = ref<InstanceType<typeof ComponentDetailsDialog> | null>(null);
const possibleTypes = computed(() => {
  const types = tableItems.value
    .map((remark: ScanRemark) => {
      return remark.type;
    })
    .sort();
  const distinctTypes = [...new Set(types)];
  return distinctTypes;
});
const possibleRemarks = computed(() => {
  if (!tableItems.value) {
    return [];
  }
  return _.chain(tableItems.value)
    .uniqBy('remarkKey')
    .map((remark: ScanRemark) => {
      return {
        value: remark.remarkKey,
        text: t(remark.remarkKey),
      } as IDefaultSelectItem;
    })
    .value();
});
const possibleStatuses = computed(() => {
  if (!tableItems.value) {
    return [];
  }
  return _.chain(tableItems.value)
    .uniqBy('status')
    .map((remark: ScanRemark) => {
      return {
        value: remark.status,
        text: t('SCAN_REMARK_STATUS_' + remark.status),
      } as IDefaultSelectItem;
    })
    .value();
});

const handleFilterQuery = () => {
  const filter = route.query.scanRemarkLevel as string;
  if (filter) {
    selectedFilterStatus.value = [filter as ScanRemarkLevel];
  }
};

onMounted(async () => {
  handleFilterQuery();

  await reload();
});

const projectModel = computed(() => projectStore.currentProject!);
const version = computed(() => sbomStore.getCurrentVersion);
const spdx = computed(() => sbomStore.selectedSpdx);

const filteredList = computed(() => {
  return tableItems.value.filter((item: ScanRemark) => {
    return filterOnStatus(item) && filterOnRemark(item) && filterOnType(item);
  });
});

watch(
  () => spdx.value,
  async () => {
    dataAreLoaded.value = false;
    await reload();
  },
);
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
onMounted(reload);

const customFilterTable = (value: string, search: string) => {
  if (value != null && value) {
    const valueTranslated = t(value);
    return ('' + valueTranslated).toLowerCase().indexOf(search.toLowerCase()) > -1;
  }
  return false;
};

const showDetails = async (event: Event, row: DataTableItem<ScanRemark>) => {
  if (!row.item.spdxId) {
    return;
  }
  await ProjectService.getComponentDetailsForSbom(
    projectModel.value._key,
    version.value._key,
    spdx.value._key,
    row.item.spdxId,
  ).then((response) => {
    if (newComponentDetailsDlg.value) {
      newComponentDetailsDlg.value?.open(
        projectModel.value,
        version.value._key,
        spdx.value._key,
        response.data,
        row.item.policyRuleStatus,
        row.item.unmatchedLicenses,
      );
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
      spdx.value._key,
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
  () => route.path,
  async (_newPath, _oldPath) => {
    handleFilterQuery();
  },
);
</script>
