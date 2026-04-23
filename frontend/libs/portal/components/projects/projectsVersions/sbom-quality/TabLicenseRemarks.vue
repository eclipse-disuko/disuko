<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->
<script setup lang="ts">
import {useView} from '@disclosure-portal/composables/useView';
import {IObligation, ObligationDTO} from '@disclosure-portal/model/IObligation';
import {compareLevel, LicenseRemarks} from '@disclosure-portal/model/Quality';
import ProjectService, {RemarkTypes} from '@disclosure-portal/services/projects';
import VersionService from '@disclosure-portal/services/version';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import {downloadFile} from '@disclosure-portal/utils/download';
import {formatDateAndTime} from '@disclosure-portal/utils/Table';
import useViewTools, {getIconColorOfLevel, getIconOfLevel, getStrWithMaxLength} from '@disclosure-portal/utils/View';
import {DataTableHeader, DataTableHeaderFilterItems, DataTableItem, SortItem} from '@shared/types/table';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import _ from 'lodash';
import {computed, onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const appStore = useAppStore();
const projectStore = useProjectStore();
const sbomStore = useSbomStore();
const viewTools = useViewTools();
const {getTextOfLevel, getTextOfType} = useView();

const selectedLicenseRemarks = ref<LicenseRemarks>({
  license: '',
  obligations: [],
  warnings: false,
  alarms: false,
  affected: [],
});

const expanded = ref<string[]>([]);
const remarks = ref<LicenseRemarks[]>([]);
const filteredRemarks = ref<LicenseRemarks[]>([]);
const search = ref('');
const sortBy = ref<SortItem[]>([{key: 'warnLevel', order: 'desc'}]);
const headers = ref<DataTableHeader[]>([
  {title: '', value: 'data-table-expand', width: 53},
  {
    title: t('COL_LEVEL'),
    align: 'center',
    width: 130,
    key: 'warnLevel',
    sort: compareLevel,
  },
  {
    title: t('COL_TYPE'),
    align: 'start',
    width: 130,
    key: 'type',
  },
  {
    title: t('COL_QUALITY_REMARK'),
    width: 210,
    align: 'start',
    key: 'name',
  },
  {
    title: t('COL_DESCRIPTION'),
    align: 'start',
    key: 'description',
  },
]);

const innerHeaders = ref<DataTableHeader[]>([
  {
    title: t('COL_NAME'),
    align: 'start',
    key: 'name',
  },
  {
    title: t('COL_VERSION'),
    align: 'start',
    width: 150,
    key: 'version',
  },
]);

const dataAreLoaded = ref(false);
const selectedFilterStatus = ref<string[]>([]);
const selectedFilterQualityRemark = ref<string[]>([]);
const selectedFilterTypes = ref<string[]>([]);
const tableHeight = ref(0);
const searchFieldInput = ref<string>('');

const projectModel = computed(() => projectStore.currentProject!);
const version = computed(() => sbomStore.getCurrentVersion);
const spdx = computed(() => sbomStore.getSelectedSBOM);

const possibleRemarks = computed((): DataTableHeaderFilterItems[] => {
  const remarkSet = new Set();

  selectedLicenseRemarks.value.obligations.forEach((item: ObligationDTO) => {
    remarkSet.add(item.name);
  });

  return [...remarkSet].map((value) => ({value: value as string}));
});

const possibleTypes = computed((): DataTableHeaderFilterItems[] => {
  if (!selectedLicenseRemarks.value.obligations) {
    return [];
  }

  const uniqueTypes = [...new Set(selectedLicenseRemarks.value.obligations.map(({type}) => type))];

  return uniqueTypes.map(
    (type: string) =>
      ({
        text: getTextOfType(type),
        value: type,
      }) as DataTableHeaderFilterItems,
  );
});

const possibleStatuses = computed((): DataTableHeaderFilterItems[] => {
  if (!selectedLicenseRemarks.value.obligations) {
    return [];
  }

  const uniqueLicenseRemarks = [...new Set(selectedLicenseRemarks.value.obligations.map(({warnLevel}) => warnLevel))];

  return uniqueLicenseRemarks.map(
    (warnLevel: string) =>
      ({
        text: getTextOfLevel(warnLevel),
        value: warnLevel,
        iconColor: getIconColorOfLevel(warnLevel),
        icon: getIconOfLevel(warnLevel),
      }) as DataTableHeaderFilterItems,
  );
});

const filteredList = computed(() => {
  if (!selectedLicenseRemarks.value) {
    return [];
  }
  return selectedLicenseRemarks.value.obligations.filter((item: ObligationDTO) => {
    return filterOnStatus(item) && filterOnRemark(item) && filterOnType(item);
  });
});

const classificationsCustomFilterTable = (value: string, searchTerm: string, item: IObligation ) => {
  if (value != null && value) {
    const dateTime = formatDateAndTime(value);
    if (dateTime && dateTime !== 'Invalid date') {
      return dateTime.indexOf(searchTerm) > -1;
    }

    let found = ('' + value).toLowerCase().indexOf(searchTerm.toLowerCase()) > -1;
    if (!found && value === item.type) {
      found = ('' + getTextOfType(value)).toLowerCase().indexOf(searchTerm.toLowerCase()) > -1;
    }
    if (!found && value === item.name && appStore.getAppLanguage === 'de') {
      found = ('' + item.nameDe).toLowerCase().indexOf(searchTerm.toLowerCase()) > -1;
    }
    if (!found && value === item.description && appStore.getAppLanguage === 'de') {
      found = ('' + item.descriptionDe).toLowerCase().indexOf(searchTerm.toLowerCase()) > -1;
    }
    return found;
  }
  return false;
};

const selectedLicenseChanged = () => {
  if (!selectedLicenseRemarks.value) {
    return;
  }

  filteredRemarks.value = remarks.value;
  searchFieldInput.value = '';
  expanded.value = [];
};

const reload = async (): Promise<void> => {
  if (!projectModel.value || !projectModel.value._key || !spdx.value) {
    dataAreLoaded.value = true;
    return;
  }

  dataAreLoaded.value = false;
  remarks.value = await VersionService.getLicenseRemarksForSbom(
    projectModel.value._key,
    version.value._key,
    spdx.value._key,
  );

  dataAreLoaded.value = true;
  if (remarks.value.length === 0) {
    return;
  }
  filteredRemarks.value = remarks.value;
  selectedLicenseRemarks.value = remarks.value[0];
  selectedLicenseChanged();
};

const searchForLicense = async () => {
  if (!searchFieldInput.value) {
    return [];
  }
  return _.chain(remarks.value)
    .filter((r) => r.license.toLowerCase().includes(searchFieldInput.value.toLowerCase()))
    .value();
};

const filterOnRemark = (item: ObligationDTO): boolean => {
  if (!selectedFilterQualityRemark.value.length) {
    return true;
  }
  return selectedFilterQualityRemark.value.includes(item.name);
};

const filterOnStatus = (item: ObligationDTO): boolean => {
  if (!selectedFilterStatus.value.length) {
    return true;
  }
  return selectedFilterStatus.value.includes(item.warnLevel.toUpperCase());
};

const filterOnType = (item: ObligationDTO): boolean => {
  if (!selectedFilterTypes.value.length) {
    return true;
  }
  return selectedFilterTypes.value.includes(item.type);
};

const downloadLicenseRemarksCsv = async () => {
  downloadFile(
    `${projectModel.value.name}_${version.value.name}_license_remarks.csv`,
    ProjectService.downloadScanOrLicenseRemarksForSbomCsv(
      projectModel.value._key,
      version.value._key,
      RemarkTypes.license,
      spdx.value!._key,
    ),
    true,
  );
};

onMounted(async () => {
  await reload();
});

watch(() => spdx.value, reload);
</script>

<template>
  <div class="h-[calc(100%-56px)]">
    <Stack direction="row" class="pb-1">
      <v-autocomplete
        v-model="selectedLicenseRemarks"
        :items="filteredRemarks"
        :search-input.sync="searchFieldInput"
        :label="t('LABEL_LICENSE_CURRENT')"
        @keyup="searchForLicense()"
        variant="outlined"
        density="compact"
        max-width="500"
        open-on-clear
        auto-select-first
        clearable
        hide-details
        item-text="license"
        return-object
        color="inputActiveBorderColor"
        @change="selectedLicenseChanged"
        style="max-height: 40px !important">
        <template #item="{item, props}">
          <v-list-item v-bind="props" :title="undefined">
            <v-icon v-if="item.value.alarms" :color="getIconColorOfLevel('alarm')" dense
            >{{ getIconOfLevel('alarm') }}
            </v-icon>
            <v-icon v-else-if="item.value.warnings" :color="getIconColorOfLevel('warning')" dense
            >{{ getIconOfLevel('warning') }}
            </v-icon>
            <span class="d-text d-secondary-text">{{ item.value.license }} ({{ item.value.affected.length }})</span>
          </v-list-item>
        </template>
        <template #selection="{item}">
          <div class="d-inline">
            <v-icon v-if="item.value.alarms" :color="getIconColorOfLevel('alarm')" dense
            >{{ getIconOfLevel('alarm') }}
            </v-icon>
            <v-icon v-else-if="item.value.warnings" :color="getIconColorOfLevel('warning')" dense
            >{{ getIconOfLevel('warning') }}
            </v-icon>
            <span class="d-text d-secondary-text">{{ item.value.license }} ({{ item.value.affected.length }})</span>
          </div>
        </template>
      </v-autocomplete>
      <DCActionButton
        :text="t('BTN_DOWNLOAD')"
        icon="mdi-download"
        :hint="t('TT_download_license_remarks')"
        @click="downloadLicenseRemarksCsv"
        class="ml-2 pr-4" />
      <div class="grow"></div>
      <DSearchField v-model="search" />
    </Stack>

    <v-data-table
      :loading="!dataAreLoaded"
      density="compact"
      :headers="headers"
      fixed-header
      :height="tableHeight"
      class="striped-table custom-data-table my-0 h-full py-0"
      item-value="_key"
      :sort-by="sortBy"
      sort-desc
      :search="search"
      :items-per-page="-1"
      :items="filteredList"
      :expanded.sync="expanded"
      @click:row.stop="
        (_: Event, tableItem: DataTableItem<any>) =>
          expanded.some((e: string) => e === tableItem.item._key) ? (expanded = []) : (expanded = [tableItem.item._key])
      "
      :footer-props="{
        'items-per-page-options': [10, 50, 100, -1],
      }"
      :custom-filter="classificationsCustomFilterTable">
      <template #[`header.warnLevel`]="{column, toggleSort, getSortIcon}">
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
      <template #[`header.type`]="{column, toggleSort, getSortIcon}">
        <GridFilterHeader :column="column" :getSortIcon="getSortIcon" :toggleSort="toggleSort">
          <template #filter>
            <GridHeaderFilterIcon
              v-model="selectedFilterTypes"
              :column="column"
              :label="t('COL_TYPE')"
              :allItems="possibleTypes">
            </GridHeaderFilterIcon>
          </template>
        </GridFilterHeader>
      </template>
      <template #[`header.name`]="{column, toggleSort, getSortIcon}">
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
      <template #[`item.type`]="{item}">
        {{ getTextOfType(item.type) }}
      </template>
      <template #[`item.warnLevel`]="{item}">
        <span>
          <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" bottom>
            <template #activator="{props, targetRef}">
              <v-icon v-bind="props" v-on="targetRef" :color="getIconColorOfLevel(item.warnLevel)" dense>{{
                  getIconOfLevel(item.warnLevel)
                }}</v-icon>
            </template>
            <span>{{ getTextOfLevel(item.warnLevel) }}</span>
          </v-tooltip>
        </span>
      </template>
      <template #[`item.remark`]="{item}">
        {{ t('' + item.remark) }}
      </template>
      <template #[`item.name`]="{item}">
        {{ viewTools.getNameForLanguage(item) }}
      </template>
      <template #[`item.description`]="{item}">
        <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" bottom>
          <template #activator="{props, targetRef}">
            <span v-bind="props" v-on="targetRef">
              {{ getStrWithMaxLength(180, t(viewTools.getDescriptionForLanguage(item))) }}
            </span>
          </template>
          <span>{{ t(viewTools.getDescriptionForLanguage(item)) }}</span>
        </v-tooltip>
      </template>
      <template #[`item.data-table-expand`]="{item}">
        <v-icon
          color="primary"
          @click.stop="expanded.some((e: string) => e === item._key) ? (expanded = []) : (expanded = [item._key])">
          {{ expanded.some((e: string) => e === item._key) ? 'mdi-chevron-up' : 'mdi-chevron-down' }}
        </v-icon>
      </template>
      <template #expanded-row="{columns, item}">
        <td v-if="selectedLicenseRemarks" :colspan="columns.length" style="height: 10%">
          <v-data-table
            :headers="innerHeaders"
            :item-key="item.spdxid + '-' + item.name"
            :items="selectedLicenseRemarks.affected"
            :hide-default-header="true"
            :hide-default-footer="true"
            disable-pagination
            class="custom-data-table"
            density="compact">
            <template #item="{item}">
              <tr>
                <td style="width: 100px"></td>
                <td style="width: 150px">{{ item.name }}</td>
                <td>{{ item.version }}</td>
              </tr>
            </template>
          </v-data-table>
        </td>
      </template>
    </v-data-table>
  </div>
</template>
