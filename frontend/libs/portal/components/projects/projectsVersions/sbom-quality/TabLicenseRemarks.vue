// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {useView} from '@disclosure-portal/composables/useView';
import {IDefaultSelectItem, IObligation, ObligationDTO} from '@disclosure-portal/model/IObligation';
import {compareLevel, LicenseRemarks} from '@disclosure-portal/model/Quality';
import ProjectService, {RemarkTypes} from '@disclosure-portal/services/projects';
import VersionService from '@disclosure-portal/services/version';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import {downloadFile} from '@disclosure-portal/utils/download';
import {formatDateAndTime} from '@disclosure-portal/utils/Table';
import useViewTools, {getIconColorOfLevel, getIconOfLevel, getStrWithMaxLength} from '@disclosure-portal/utils/View';
import {DataTableHeader, DataTableItem, SortItem} from '@shared/types/table';
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

const classificationsCustomFilterTable = (value: string, search: string, item: IObligation) => {
  if (value != null && value) {
    const dateTime = formatDateAndTime(value);
    if (dateTime && dateTime !== 'Invalid date') {
      return dateTime.indexOf(search) > -1;
    }

    let found = ('' + value).toLowerCase().indexOf(search.toLowerCase()) > -1;
    if (!found && value === item.type) {
      found = ('' + getTextOfType(value)).toLowerCase().indexOf(search.toLowerCase()) > -1;
    }
    if (!found && value === item.name && appStore.getAppLanguage === 'de') {
      found = ('' + item.nameDe).toLowerCase().indexOf(search.toLowerCase()) > -1;
    }
    if (!found && value === item.description && appStore.getAppLanguage === 'de') {
      found = ('' + item.descriptionDe).toLowerCase().indexOf(search.toLowerCase()) > -1;
    }
    return found;
  }
  return false;
};

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
  {title: '', class: 'tableHeaderCell', value: 'data-table-expand', width: '53'},
  {
    title: '' + t('COL_LEVEL'),
    align: 'center',
    class: 'tableHeaderCell',
    width: 130,
    key: 'warnLevel',
    sort: compareLevel,
  },
  {
    title: '' + t('COL_TYPE'),
    align: 'start',
    class: 'tableHeaderCell',
    width: 130,
    key: 'type',
  },
  {
    title: '' + t('COL_QUALITY_REMARK'),
    width: 210,
    align: 'start',
    class: 'tableHeaderCell',
    key: 'name',
  },
  {
    title: '' + t('COL_DESCRIPTION'),
    align: 'start',
    class: 'tableHeaderCell',
    key: 'description',
  },
]);

const innerHeaders = ref<DataTableHeader[]>([
  {
    title: '' + t('COL_NAME'),
    align: 'start',
    class: 'tableHeaderCell',
    key: 'name',
  },
  {
    title: '' + t('COL_VERSION'),
    align: 'start',
    class: 'tableHeaderCell',
    width: 150,
    key: 'version',
  },
]);

const dataAreLoaded = ref(false);
const selectedFilterStatus = ref<string[]>([]);
const selectedFilterQualityRemark = ref<string[]>([]);
const selectedFilterTypes = ref<string[]>([]);
const allRemarks = ref<IDefaultSelectItem[]>([]);
const menu = ref(false);
const menu2 = ref(false);
const menu3 = ref(false);
const tableHeight = ref(0);

// Menu states
watch(menu, () => (menu2.value = menu3.value = false));
watch(menu2, () => (menu.value = menu3.value = false));
watch(menu3, () => (menu.value = menu2.value = false));

const searchFieldInput = ref<string>('');

const projectModel = computed(() => projectStore.currentProject!);
const version = computed(() => sbomStore.getCurrentVersion);
const spdx = computed(() => sbomStore.selectedSpdx);

const possibleTypes = computed(() => {
  if (!selectedLicenseRemarks.value.obligations) {
    return [];
  }
  return _.chain(selectedLicenseRemarks.value.obligations)
    .uniqBy((item: ObligationDTO) => {
      return item.type;
    })
    .map((item: ObligationDTO) => {
      return {
        text: getTextOfType(item.type),
        value: item.type,
      } as IDefaultSelectItem;
    })
    .value();
});

const possibleStatuses = computed(() => {
  if (!selectedLicenseRemarks.value.obligations) {
    return [];
  }
  return _.chain(selectedLicenseRemarks.value.obligations)
    .uniqBy((item: ObligationDTO) => {
      return item.warnLevel;
    })
    .map((item: ObligationDTO) => {
      return {
        text: item.warnLevel,
        value: item.warnLevel,
      } as IDefaultSelectItem;
    })
    .value();
});

const selectedLicenseChanged = () => {
  if (!selectedLicenseRemarks.value) {
    return;
  }

  filteredRemarks.value = remarks.value;
  searchFieldInput.value = '';
  expanded.value = [];
  const remarkSet = new Set<string>();
  allRemarks.value = [];
  selectedLicenseRemarks.value.obligations.forEach((item: ObligationDTO) => {
    remarkSet.add(item.name);
  });
  Array.from(remarkSet).forEach((item: string) => {
    allRemarks.value.push({value: item, text: item} as IDefaultSelectItem);
  });
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

onMounted(async () => {
  await reload();
});

watch(() => spdx.value, reload);

const searchForLicense = async () => {
  if (!searchFieldInput.value) {
    return [];
  }
  return _.chain(remarks.value)
    .filter((r) => r.license.toLowerCase().includes(searchFieldInput.value.toLowerCase()))
    .value();
};

const filteredList = computed(() => {
  if (!selectedLicenseRemarks.value) {
    return [];
  }
  return selectedLicenseRemarks.value.obligations.filter((item: ObligationDTO) => {
    return filterOnStatus(item) && filterOnRemark(item) && filterOnType(item);
  });
});

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
      spdx.value._key,
    ),
    true,
  );
};
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
        <template v-slot:item="{item, props}">
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
        <template v-slot:selection="{item}">
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
        class="pr-4 ml-2" />
      <div class="grow"></div>
      <v-text-field
        class="grow-2 w-full md:max-w-[400px] md:w-auto"
        min-width="50"
        autocomplete="off"
        variant="outlined"
        density="compact"
        v-model="search"
        append-inner-icon="mdi-magnify"
        :label="t('labelSearch')"
        clearable
        hide-details
        single-line></v-text-field>
    </Stack>

    <v-data-table
      :loading="!dataAreLoaded"
      density="compact"
      :headers="headers"
      fixed-header
      :height="tableHeight"
      class="striped-table py-0 my-0 custom-data-table h-full"
      item-value="_key"
      :sort-by="sortBy"
      sort-desc
      :search="search"
      :items-per-page="-1"
      :items="filteredList"
      :expanded.sync="expanded"
      @click:row.stop="
        (event: Event, tableItem: DataTableItem<any>) =>
          expanded.some((e: string) => e === tableItem.item._key) ? (expanded = []) : (expanded = [tableItem.item._key])
      "
      :footer-props="{
        'items-per-page-options': [10, 50, 100, -1],
      }"
      :custom-filter="classificationsCustomFilterTable">
      <template v-slot:header.warnLevel="{column, toggleSort, getSortIcon}">
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
                transition="scale-transition"
                menu
                persistent-clear
                :list-props="{class: 'striped-filter-dd py-0'}">
                <template v-slot:item="{item, props}">
                  <v-list-item v-bind="props" class="py-0 px-2">
                    <template v-slot:prepend="{isSelected}">
                      <v-checkbox hide-details :model-value="isSelected" />
                    </template>
                    <template v-slot:title>
                      <v-icon :color="getIconColorOfLevel(item.value)" dense>{{ getIconOfLevel(item.value) }}</v-icon>
                      <span class="ml-1 pFilterEntry">{{ getTextOfLevel(item.value) }}</span>
                    </template>
                  </v-list-item>
                </template>
                <template v-slot:selection="{item, index}">
                  <div v-if="index === 0" class="d-flex align-center">
                    <v-icon :color="getIconColorOfLevel(item.value)" dense>{{ getIconOfLevel(item.value) }}</v-icon>
                    <span class="ml-1 pFilterEntry">{{ getTextOfLevel(item.value) }}</span>
                  </div>
                  <span v-if="index === 1" class="pAdditionalFilter">
                    +{{ selectedFilterStatus.length - 1 }} others
                  </span>
                </template>
              </v-select>
            </div>
          </v-menu>
          <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
        </div>
      </template>
      <template v-slot:header.type="{column, toggleSort, getSortIcon}">
        <div class="v-data-table-header__content">
          <span>{{ column.title }}</span>
          <v-menu :close-on-content-click="false" v-model="menu2">
            <template v-slot:activator="{props}">
              <DIconButton
                :parentProps="props"
                icon="mdi-filter-variant"
                :hint="t('TT_SHOW_FILTER')"
                :color="selectedFilterTypes.length > 0 ? 'primary' : 'default'" />
            </template>
            <div class="bg-background" style="width: 280px">
              <v-row class="d-flex justify-end ma-1 mr-2">
                <DIconButton icon="mdi-close" @clicked="menu2 = false" color="default" />
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
                transition="scale-transition"
                menu
                persistent-clear
                :list-props="{class: 'striped-filter-dd py-0'}">
                <template v-slot:item="{props, item}">
                  <v-list-item v-bind="props" class="py-0 px-2" title="">
                    <template v-slot:prepend="{isSelected}">
                      <v-checkbox hide-details :model-value="isSelected" />
                    </template>
                    <span class="pFilterEntry">{{ item.title }}</span>
                  </v-list-item>
                </template>
                <template v-slot:selection="{item, index}">
                  <div v-if="index === 0" class="d-flex align-center">
                    <span class="pFilterEntry">{{ item.title }}</span>
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
      <template v-slot:header.name="{column, toggleSort, getSortIcon}">
        <div class="v-data-table-header__content">
          <span>{{ column.title }}</span>
          <v-menu :close-on-content-click="false" v-model="menu3">
            <template v-slot:activator="{props}">
              <DIconButton
                :parentProps="props"
                icon="mdi-filter-variant"
                :hint="t('TT_SHOW_FILTER')"
                :color="selectedFilterQualityRemark.length > 0 ? 'primary' : 'default'" />
            </template>
            <div class="bg-background" style="width: 280px">
              <v-row class="d-flex justify-end ma-1 mr-2">
                <DIconButton icon="mdi-close" @clicked="menu3 = false" color="default" />
              </v-row>
              <v-select
                v-model="selectedFilterQualityRemark"
                :items="allRemarks"
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
                <template v-slot:item="{props, item}">
                  <v-list-item v-bind="props" class="py-0 px-2" title="">
                    <template v-slot:prepend="{isSelected}">
                      <v-checkbox hide-details :model-value="isSelected" />
                    </template>
                    <span class="pFilterEntry">{{ item.title }}</span>
                  </v-list-item>
                </template>
                <template v-slot:selection="{item, index}">
                  <div v-if="index === 0" class="d-flex align-center">
                    <span class="pFilterEntry">{{ item.title }}</span>
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
      <template v-slot:item.type="{item}">
        {{ getTextOfType(item.type) }}
      </template>
      <template v-slot:item.warnLevel="{item}">
        <span>
          <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" bottom>
            <template v-slot:activator="{props, targetRef}">
              <v-icon v-bind="props" v-on="targetRef" :color="getIconColorOfLevel(item.warnLevel)" dense>{{
                getIconOfLevel(item.warnLevel)
              }}</v-icon>
            </template>
            <span>{{ getTextOfLevel(item.warnLevel) }}</span>
          </v-tooltip>
        </span>
      </template>
      <template v-slot:item.remark="{item}">
        {{ t('' + item.remark) }}
      </template>
      <template v-slot:item.name="{item}">
        {{ viewTools.getNameForLanguage(item) }}
      </template>
      <template v-slot:item.description="{item}">
        <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" bottom>
          <template v-slot:activator="{props, targetRef}">
            <span v-bind="props" v-on="targetRef">
              {{ getStrWithMaxLength(180, t(viewTools.getDescriptionForLanguage(item))) }}
            </span>
          </template>
          <span>{{ t(viewTools.getDescriptionForLanguage(item)) }}</span>
        </v-tooltip>
      </template>
      <template v-slot:item.data-table-expand="{item}">
        <v-icon
          color="primary"
          @click.stop="expanded.some((e: string) => e === item._key) ? (expanded = []) : (expanded = [item._key])">
          {{ expanded.some((e: string) => e === item._key) ? 'mdi-chevron-up' : 'mdi-chevron-down' }}
        </v-icon>
      </template>
      <template v-slot:expanded-row="{columns, item}">
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
            <template v-slot:item="{item}">
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

    <ComponentDetailsDialog ref="newComponentDetailsDlg"></ComponentDetailsDialog>
  </div>
</template>
