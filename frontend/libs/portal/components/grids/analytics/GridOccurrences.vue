<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script lang="ts" setup>
import {SearchOccurrenciesItem} from '@disclosure-portal/model/Analytics';
import {ISelectItemWithCount} from '@disclosure-portal/model/ISelectItem';
import {PossibleFilterValues, compareFamily, getLicenseApprovalTypeKeys} from '@disclosure-portal/model/License';
import AnalyticsService from '@disclosure-portal/services/analytics';
import {getCssClassForTableRow} from '@disclosure-portal/utils/Table';
import useSnackbar from '@shared/composables/useSnackbar';
import {DataTableHeader, DataTableItem, SortItem} from '@shared/types/table';
import {computed, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';

const {t} = useI18n();
const router = useRouter();
const {info: snack} = useSnackbar();

// Reactive variables
const unreferencedOnly = ref(false);
const search = ref('');
const isLoading = ref(false);
const items = ref<SearchOccurrenciesItem[]>([]);
const menuIsLicenseChart = ref(false);
const selectedFilterIsLicenseChart = ref<string[]>([]);
const possibleIsLicenseChart = ref<ISelectItemWithCount[]>([]);
const menuSource = ref(false);
const selectedFilterSource = ref<string[]>([]);
const possibleSources = ref<ISelectItemWithCount[]>([]);
const menuFamily = ref(false);
const selectedFilterFamily = ref<string[]>([]);
const possibleFamilies = ref<ISelectItemWithCount[]>([]);
const menuApproval = ref(false);
const selectedFilterApproval = ref<string[]>([]);
const possibleApproval = ref<ISelectItemWithCount[]>([]);
const menuType = ref(false);
const selectedFilterType = ref<string[]>([]);
const possibleType = ref<ISelectItemWithCount[]>([]);
const sortBy = ref<SortItem[]>([]);
const headers = ref<DataTableHeader[]>([
  {title: t('COL_LICENSE_ID'), align: 'start', value: 'origName', width: '240'},
  {title: t('COL_COUNT'), align: 'start', value: 'count', width: '120', sortable: true},
  {title: t('COL_REFERENCE'), align: 'start', value: 'license.licenseId', width: '240'},
  {title: t('COL_LICENSE_NAME'), align: 'start', value: 'license.name', width: '240'},
  {
    title: t('COL_LICENSE_CHART_STATUS'),
    align: 'center',
    value: 'license.meta.isLicenseChart',
    width: 115,
    sortable: false,
  },
  {title: t('COL_APPROVAL_STATUS'), value: 'license.meta.approvalState', width: 165, sortable: false},
  {title: t('COL_LICENSE_FAMILY'), align: 'start', value: 'license.meta.family', width: 165, sortable: false},
  {title: t('COL_TYPE'), value: 'license.meta.licenseType', width: 115, sortable: false},
  {title: t('COL_LICENSE_SOURCE'), value: 'license.source', width: 115, sortable: false},
  {title: t('COL_ACTIONS'), align: 'center', value: 'actions', width: 115, sortable: false},
]);

sortBy.value = [{key: 'count', order: 'desc'}];

const renderChartValue = (value: string): string => {
  switch (value) {
    case 'true':
      return t('TABLE_LICENSE_CHART_STATUS_IS');
    case 'false':
      return t('TABLE_LICENSE_CHART_STATUS_IS_NOT');
    default:
      return t('TABLE_LICENSE_CHART_UNREFERENCED');
  }
};

const toI18n = (title: string) => {
  return title.toUpperCase().replace(/ /g, '_').replace(/-/g, '_');
};

const filteredItems = computed(() => {
  return items.value.filter((f) => {
    if (unreferencedOnly.value && f.license) {
      return false;
    }
    if (!unreferencedOnly.value) {
      if (!isMatchedByChartFilter(f)) {
        return false;
      }
      if (
        selectedFilterApproval.value.length > 0 &&
        !selectedFilterApproval.value.some((fa) => f.license?.meta.approvalState === fa)
      ) {
        return false;
      }
      if (
        selectedFilterFamily.value.length > 0 &&
        !selectedFilterFamily.value.some((ff) => f.license?.meta.family === ff)
      ) {
        return false;
      }
      if (
        selectedFilterType.value.length > 0 &&
        !selectedFilterType.value.some((ft) => f.license?.meta.licenseType === ft)
      ) {
        return false;
      }
      if (selectedFilterSource.value.length > 0 && !selectedFilterSource.value.some((fs) => f.license?.source === fs)) {
        return false;
      }
    }
    if (!search.value) {
      return true;
    }
    return (
      f.origName.includes(search.value) ||
      f.license?.name.includes(search.value) ||
      f.license?.licenseId.includes(search.value)
    );
  });
});

const isMatchedByChartFilter = (item: SearchOccurrenciesItem): boolean => {
  if (selectedFilterIsLicenseChart.value.length === 0) {
    return true;
  }
  if (!item.license && selectedFilterIsLicenseChart.value.indexOf('unreferenced') !== -1) {
    return true;
  }
  if (item.license) {
    if (item.license.meta.isLicenseChart && selectedFilterIsLicenseChart.value.indexOf('true') !== -1) {
      return true;
    }
    if (!item.license.meta.isLicenseChart && selectedFilterIsLicenseChart.value.indexOf('false') !== -1) {
      return true;
    }
  }
  return false;
};

const setPossibleValues = (possibleValues: PossibleFilterValues) => {
  possibleIsLicenseChart.value = Object.entries(possibleValues.possibleCharts)
    .map(
      ([k, count]) =>
        ({
          text: renderChartValue(k),
          value: k,
          count: count,
        }) as ISelectItemWithCount,
    )
    .sort();
  possibleSources.value = Object.entries(possibleValues.possibleSources).map(
    ([k, count]) =>
      ({
        text: k,
        value: k,
        count: count,
      }) as ISelectItemWithCount,
  );
  possibleFamilies.value = Object.entries(possibleValues.possibleFamilies)
    .sort((a, b) => compareFamily(a[0], b[0]))
    .map(
      ([k, count]) =>
        ({
          text: t('LIC_FAMILY_' + (toI18n(k) || 'UNKNOWN')),
          value: k,
          count: count,
        }) as ISelectItemWithCount,
    );
  possibleApproval.value = Object.entries(possibleValues.possibleApproval)
    .sort(([keyA], [keyB]) => getLicenseApprovalTypeKeys().indexOf(keyA) - getLicenseApprovalTypeKeys().indexOf(keyB))
    .map(
      ([k, count]) =>
        ({
          text: t('LT_APP_' + (toI18n(k) || 'UNKNOWN')),
          value: k,
          count: count,
        }) as ISelectItemWithCount,
    );
  possibleType.value = Object.entries(possibleValues.possibleType).map(
    ([k, count]) =>
      ({
        text: t('LT_' + (toI18n(k) || 'UNKNOWN')),
        value: k,
        count: count,
      }) as ISelectItemWithCount,
  );
};

// Methods
const reload = async () => {
  isLoading.value = true;
  items.value = [];
  try {
    const result = await AnalyticsService.searchOccurrencies();
    setPossibleValues(result.data.possibleValues);
    items.value = result.data.list;
  } catch (error: any) {
    snack(t('ERROR_OCCURRENCES_API'));
    console.error(error);
  } finally {
    isLoading.value = false;
  }
};

const onClickRow = (event: Event, table: DataTableItem<SearchOccurrenciesItem>) => {
  router.push('/dashboard/analytics/overview?license=' + encodeURIComponent(table.item.origName));
};

const getReferenceInfoForClipboard = (item: SearchOccurrenciesItem): string => {
  return `Disclosure Portal Project License Occurence Reference

License Id: ${item.origName}
Count: ${item.count}
Referenced License Id: ${item.license ? item.license.licenseId : 'unreferenced'}
Is chart: ${item.license ? item.license.meta.isLicenseChart : 'unreferenced'}
Evaluation State: ${item.license ? item.license.meta.approvalState : 'unreferenced'}
Family: ${item.license ? item.license.meta.family : 'unreferenced'}
Type: ${item.license ? item.license.meta.licenseType : 'unreferenced'}
Source: ${item.license ? item.license.source : 'unreferenced'}`;
};

onMounted(async () => {
  await reload();
});
</script>

<template>
  <TableLayout has-tab has-title>
    <template #buttons>
      <div class="grid w-full grid-cols-12">
        <div class="sm:col-span-5 md:col-span-4 lg:col-span-2">
          <v-checkbox v-model="unreferencedOnly" hide-details color="primary" :label="t('LABEL_UNREF_ONLY')" />
        </div>
        <v-spacer class="sm:col-span-2 md:col-span-4 lg:col-span-7"></v-spacer>
        <div class="sm:col-span-5 md:col-span-4 lg:col-span-3">
          <v-text-field
            autocomplete="off"
            v-model="search"
            :label="t('labelSearch')"
            append-inner-icon="mdi-magnify"
            variant="outlined"
            clearable
            density="compact" />
        </div>
      </div>
    </template>
    <template #table>
      <v-data-table
        :loading="isLoading"
        :headers="headers"
        :sort-by="sortBy"
        sort-desc
        density="compact"
        fixed-header
        @click:row="onClickRow"
        class="striped-table fill-height"
        :items-per-page="50"
        :footer-props="{'items-per-page-options': [10, 50, 100, -1]}"
        :items="filteredItems"
        :item-class="getCssClassForTableRow">
        <template #[`header.license.meta.isLicenseChart`]="{column, getSortIcon, toggleSort}" v-if="!unreferencedOnly">
          <div class="v-data-table-header__content">
            <span>{{ column.title }}</span>
            <v-menu :close-on-content-click="false" v-model="menuIsLicenseChart">
              <template v-slot:activator="{props}">
                <DIconButton
                  :parentProps="props"
                  icon="mdi-filter-variant"
                  :hint="t('TT_SHOW_FILTER')"
                  :color="
                    selectedFilterIsLicenseChart && selectedFilterIsLicenseChart.length > 0 ? 'primary' : 'secondary'
                  " />
              </template>
              <div class="bg-background" style="width: 320px">
                <v-row class="d-flex ma-1 mr-2 justify-end">
                  <DCloseButton @click="menuIsLicenseChart = false" />
                </v-row>
                <v-select
                  v-model="selectedFilterIsLicenseChart"
                  :items="possibleIsLicenseChart"
                  class="pa-2 mx-2 pb-4"
                  :label="t('Lbl_filter_License_Chart_Status')"
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
                  <template v-slot:item="{props}">
                    <v-list-item v-bind="props" class="px-2 py-0">
                      <template v-slot:prepend="{isSelected}">
                        <v-checkbox hide-details :model-value="isSelected" />
                      </template>
                      <template v-slot:title="{title}">
                        <span :class="'pFilterEntry'"> {{ title }}</span>
                      </template>
                    </v-list-item>
                  </template>
                  <template v-slot:selection="{item, index}">
                    <div v-if="index === 0" class="d-flex align-center">
                      <span :class="'pFilterEntry'">{{ item.title }}</span>
                    </div>
                    <span v-if="index === 1" class="pAdditionalFilter">
                      +{{ selectedFilterIsLicenseChart.length - 1 }} others
                    </span>
                  </template>
                </v-select>
              </div>
            </v-menu>
            <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
          </div>
        </template>
        <template #[`header.license.source`]="{column, getSortIcon, toggleSort}" v-if="!unreferencedOnly">
          <div class="v-data-table-header__content">
            <span>{{ column.title }}</span>
            <v-menu :close-on-content-click="false" v-model="menuSource">
              <template v-slot:activator="{props}">
                <DIconButton
                  :parentProps="props"
                  icon="mdi-filter-variant"
                  :hint="t('TT_SHOW_FILTER')"
                  :color="selectedFilterSource && selectedFilterSource.length > 0 ? 'primary' : 'secondary'" />
              </template>
              <div class="bg-background" style="width: 320px">
                <v-row class="d-flex ma-1 mr-2 justify-end">
                  <DCloseButton @click="menuSource = false" />
                </v-row>
                <v-select
                  v-model="selectedFilterSource"
                  :items="possibleSources"
                  class="pa-2 mx-2 pb-4"
                  :label="t('Lbl_filter_License_Chart_Status')"
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
                  <template v-slot:item="{props}">
                    <v-list-item v-bind="props" class="px-2 py-0">
                      <template v-slot:prepend="{isSelected}">
                        <v-checkbox hide-details :model-value="isSelected" />
                      </template>
                      <template v-slot:title="{title}">
                        <span :class="'pFilterEntry'"> {{ title }}</span>
                      </template>
                    </v-list-item>
                  </template>
                  <template v-slot:selection="{item, index}">
                    <div v-if="index === 0" class="d-flex align-center">
                      <span :class="'pFilterEntry'">{{ item.title }}</span>
                    </div>
                    <span v-if="index === 1" class="pAdditionalFilter">
                      +{{ selectedFilterSource.length - 1 }} others
                    </span>
                  </template>
                </v-select>
              </div>
            </v-menu>
            <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
          </div>
        </template>
        <template #[`header.license.meta.family`]="{column, getSortIcon, toggleSort}" v-if="!unreferencedOnly">
          <div class="v-data-table-header__content">
            <span>{{ column.title }}</span>
            <v-menu :close-on-content-click="false" v-model="menuFamily">
              <template v-slot:activator="{props}">
                <DIconButton
                  :parentProps="props"
                  icon="mdi-filter-variant"
                  :hint="t('TT_SHOW_FILTER')"
                  :color="selectedFilterFamily && selectedFilterFamily.length > 0 ? 'primary' : 'secondary'" />
              </template>
              <div class="bg-background" style="width: 320px">
                <v-row class="d-flex ma-1 mr-2 justify-end">
                  <DCloseButton @click="menuFamily = false" />
                </v-row>
                <v-select
                  v-model="selectedFilterFamily"
                  :items="possibleFamilies"
                  class="pa-2 mx-2 pb-4"
                  :label="t('Lbl_filter_License_Chart_Status')"
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
                  <template v-slot:item="{props}">
                    <v-list-item v-bind="props" class="px-2 py-0">
                      <template v-slot:prepend="{isSelected}">
                        <v-checkbox hide-details :model-value="isSelected" />
                      </template>
                      <template v-slot:title="{title}">
                        <span :class="'pFilterEntry'"> {{ title }}</span>
                      </template>
                    </v-list-item>
                  </template>
                  <template v-slot:selection="{item, index}">
                    <div v-if="index === 0" class="d-flex align-center">
                      <span :class="'pFilterEntry'">{{ item.title }}</span>
                    </div>
                    <span v-if="index === 1" class="pAdditionalFilter">
                      +{{ selectedFilterFamily.length - 1 }} others
                    </span>
                  </template>
                </v-select>
              </div>
            </v-menu>
            <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
          </div>
        </template>
        <template #[`header.license.meta.licenseType`]="{column, getSortIcon, toggleSort}" v-if="!unreferencedOnly">
          <div class="v-data-table-header__content">
            <span>{{ column.title }}</span>
            <v-menu :close-on-content-click="false" v-model="menuType">
              <template v-slot:activator="{props}">
                <DIconButton
                  :parentProps="props"
                  icon="mdi-filter-variant"
                  :hint="t('TT_SHOW_FILTER')"
                  :color="selectedFilterType && selectedFilterType.length > 0 ? 'primary' : 'secondary'" />
              </template>
              <div class="bg-background" style="width: 320px">
                <v-row class="d-flex ma-1 mr-2 justify-end">
                  <DCloseButton @click="menuType = false" />
                </v-row>
                <v-select
                  v-model="selectedFilterType"
                  :items="possibleType"
                  class="pa-2 mx-2 pb-4"
                  :label="t('Lbl_filter_License_Chart_Status')"
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
                  <template v-slot:item="{props}">
                    <v-list-item v-bind="props" class="px-2 py-0">
                      <template v-slot:prepend="{isSelected}">
                        <v-checkbox hide-details :model-value="isSelected" />
                      </template>
                      <template v-slot:title="{title}">
                        <span :class="'pFilterEntry'"> {{ title }}</span>
                      </template>
                    </v-list-item>
                  </template>
                  <template v-slot:selection="{item, index}">
                    <div v-if="index === 0" class="d-flex align-center">
                      <span :class="'pFilterEntry'">{{ item.title }}</span>
                    </div>
                    <span v-if="index === 1" class="pAdditionalFilter">
                      +{{ selectedFilterType.length - 1 }} others
                    </span>
                  </template>
                </v-select>
              </div>
            </v-menu>
            <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
          </div>
        </template>
        <template #[`header.license.meta.approvalState`]="{column, getSortIcon, toggleSort}" v-if="!unreferencedOnly">
          <div class="v-data-table-header__content">
            <span>{{ column.title }}</span>
            <v-menu :close-on-content-click="false" v-model="menuApproval">
              <template v-slot:activator="{props}">
                <DIconButton
                  :parentProps="props"
                  icon="mdi-filter-variant"
                  :hint="t('TT_SHOW_FILTER')"
                  :color="selectedFilterApproval && selectedFilterApproval.length > 0 ? 'primary' : 'secondary'" />
              </template>
              <div class="bg-background" style="width: 320px">
                <v-row class="d-flex ma-1 mr-2 justify-end">
                  <DCloseButton @click="menuApproval = false" />
                </v-row>
                <v-select
                  v-model="selectedFilterApproval"
                  :items="possibleApproval"
                  class="pa-2 mx-2 pb-4"
                  :label="t('Lbl_filter_License_Chart_Status')"
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
                  <template v-slot:item="{props}">
                    <v-list-item v-bind="props" class="px-2 py-0">
                      <template v-slot:prepend="{isSelected}">
                        <v-checkbox hide-details :model-value="isSelected" />
                      </template>
                      <template v-slot:title="{title}">
                        <span :class="'pFilterEntry'"> {{ title }}</span>
                      </template>
                    </v-list-item>
                  </template>
                  <template v-slot:selection="{item, index}">
                    <div v-if="index === 0" class="d-flex align-center">
                      <span :class="'pFilterEntry'">{{ item.title }}</span>
                    </div>
                    <span v-if="index === 1" class="pAdditionalFilter">
                      +{{ selectedFilterApproval.length - 1 }} others
                    </span>
                  </template>
                </v-select>
              </div>
            </v-menu>
            <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
          </div>
        </template>
        <template #[`item.license.meta.isLicenseChart`]="{item}">
          <DLicenseChartIcon v-if="item.license" :meta="item.license.meta" />
        </template>
        <template #[`item.actions`]="{item}">
          <DCopyClipboardButton
            :tableButton="true"
            class="ml-2"
            :hint="t('TT_COPY_REFERENCE_INFO')"
            :content="getReferenceInfoForClipboard(item)" />
        </template>
        <template #[`item.license.meta.approvalState`]="{item}">
          {{ item.license ? t('LT_APP_' + (toI18n(item.license.meta.approvalState) || 'UNKNOWN')) : '' }}
        </template>
        <template #[`item.license.meta.family`]="{item}">
          {{ item.license ? t('LIC_FAMILY_' + (toI18n(item.license.meta.family || '') || 'UNKNOWN')) : '' }}
        </template>
        <template #[`item.license.meta.licenseType`]="{item}">
          {{ item.license ? t('LT_' + (toI18n(item.license.meta.licenseType) || 'UNKNOWN')) : '' }}
        </template>
      </v-data-table>
    </template>
  </TableLayout>
</template>
