<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script lang="ts" setup>
import {CombinedSearchOptions, IAnalyticsSearchRequest, SearchResponseItem} from '@disclosure-portal/model/Analytics';
import AnalyticsService from '@disclosure-portal/services/analytics';
import {getCssClassForTableRow, SearchOptions} from '@disclosure-portal/utils/Table';
import {createProjectURL, createSBOMURL, createVersionURL} from '@shared/utils/apiUrls';
import {DataTableHeader} from '@shared/types/table';
import {debounce} from 'lodash';
import {computed, onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute} from 'vue-router';

const {t} = useI18n();
const route = useRoute();

const items = ref<SearchResponseItem[]>([]);
const components = ref<string[]>([]);
const licenses = ref<string[]>([]);
const allHeaders = computed<DataTableHeader[]>(() => [
  {title: t('COL_PROJECT_NAME'), align: 'start', value: 'name', sortable: true, width: '240'},
  {title: t('COL_PROJECT_VERSION'), align: 'start', value: 'projectVersionName', sortable: true, width: '240'},
  {title: t('COL_DELIVERY'), align: 'start', value: 'sbomName', sortable: true, width: '240'},
  {title: t('COL_COMPONENT_NAME'), align: 'start', value: 'componentName', sortable: false, width: '180'},
  {title: t('COL_COMPONENT_VERSION'), align: 'start', value: 'componentVersion', sortable: true, width: '180'},
  {title: t('COL_OWNER_COMPANY'), align: 'start', value: 'ownerCompany', sortable: false, width: '240'},
  {title: t('COL_OWNER_DEPARTMENT'), align: 'start', value: 'ownerDep', sortable: false, width: '240'},
  {title: t('COL_USER_ROLE_RESPONSIBLE'), align: 'start', value: 'responsible', sortable: false, width: '240'},
  {title: t('COL_SPDX_LICENSE_DECLARED'), align: 'start', value: 'licenseDeclared', sortable: false, width: '240'},
  {title: t('COL_SPDX_LICENSE_CONCLUDED'), align: 'start', value: 'licenseConcluded', sortable: false, width: '240'},
  {title: t('COL_LICENSE'), align: 'start', value: 'entryLicense', sortable: false, width: '240'},
  {title: t('COL_DELIVERY_STATUS'), align: 'start', value: 'sbomStatus', sortable: false, width: '240'},
  {title: t('COL_DELIVERY_DATE'), align: 'start', value: 'lastUpdate', sortable: false, width: '240'},
]);

const componentSearch = ref<string | null>(null);
const selectedHeaders = ref<number[]>([]);
const licenseSearch = ref<string | null>(null);
const myProjects = ref(false);
const exactMatch = ref(false);
const loadingAnalytics = ref(false);
const loadingComponents = ref(false);
const loadingLicenses = ref(false);
const options = ref<SearchOptions>({} as SearchOptions);
const total = ref(0);
const menuSettings = ref(false);

watch(
  () => componentSearch.value,
  () => {
    reloadAnalytics();
  },
);

watch(
  () => licenseSearch.value,
  () => {
    reloadAnalytics();
  },
);

const reloadAnalytics = async (newInput = true) => {
  loadingAnalytics.value = true;
  items.value = [];

  const analyticsSearchRequest = {
    component: componentSearch.value,
    license: licenseSearch.value,
    exactComponent: exactMatch.value,
    exactLicense: exactMatch.value,
  };
  const combinedSearchOptions = {
    analyticsRequestSearchOptions: analyticsSearchRequest,
    requestSearchOptions: options.value,
  } as CombinedSearchOptions;

  const result = await AnalyticsService.searchAnalytics(combinedSearchOptions, !myProjects.value);
  items.value = result.data.result;
  total.value = result.data.count;
  loadingAnalytics.value = false;
};

const debouncedSearchComponents = debounce(async (query: string) => {
  if (!query || query.length < 3 || query.length > 200) {
    components.value = [];
    return;
  }
  loadingComponents.value = true;

  const request = {} as IAnalyticsSearchRequest;
  request.component = query;
  request.exactComponent = exactMatch.value;
  request.exactLicense = exactMatch.value;

  const result = await AnalyticsService.searchComponents(request);
  components.value = result.data.result;
  loadingComponents.value = false;
}, 400);

const debouncedSearchLicenses = debounce(async (query: string) => {
  if (!query || query.length < 3 || query.length > 500) {
    licenses.value = [];
    return;
  }
  loadingLicenses.value = true;

  const request = {} as IAnalyticsSearchRequest;
  request.component = query;
  request.exactComponent = exactMatch.value;
  request.exactLicense = exactMatch.value;

  const result = await AnalyticsService.searchLicenses({license: query} as IAnalyticsSearchRequest);
  licenses.value = result.data.result;
  loadingLicenses.value = false;
}, 400);

const myProjectsChanged = () => {
  reloadAnalytics(true);
};

const selectableHeaders = computed(() => {
  return [...Array(allHeaders.value.length).keys()];
});

onMounted(async () => {
  selectedHeaders.value = selectableHeaders.value;
  await reloadAnalytics();
});

const filteredHeaders = computed(() => {
  const res: DataTableHeader[] = [];
  res.push({
    title: '',
    value: 'settings',
    width: '5',
    sortable: false,
    class: 'pa-0',
    selectable: false,
  });
  selectedHeaders.value.sort().forEach((i) => {
    res.push(allHeaders.value.at(i)!);
  });
  return res;
});

watch(
  () => route.query,
  (newQuery) => {
    // Update licenseSearch and reloadAnalytics based on query params
    if (!newQuery.license) {
      return;
    }
    licenseSearch.value = Array.isArray(newQuery.license) ? newQuery.license[0] : newQuery.license || '';
    reloadAnalytics();
  },
);
</script>

<template>
  <TableLayout has-tab has-title>
    <template #buttons>
      <div class="grid w-full grid-cols-12 gap-3">
        <div class="flex flex-row sm:col-span-5 md:col-span-4 lg:col-span-3">
          <v-checkbox
            v-model="myProjects"
            hide-details
            color="primary"
            @change="myProjectsChanged"
            :label="t('LBL_MY_PROJECTS')" />
          <v-checkbox v-model="exactMatch" hide-details color="primary" :label="t('LBL_EXACT_MATCH')" />
        </div>
        <v-spacer class="sm:col-span-1 md:col-span-2 lg:col-span-3"></v-spacer>
        <div class="sm:col-span-3 md:col-span-3 lg:col-span-3">
          <v-autocomplete
            :label="t('labelSearchComponent')"
            variant="outlined"
            clearable
            density="compact"
            v-model="componentSearch"
            single-line
            @update:search="debouncedSearchComponents"
            :loading="loadingComponents"
            :items="components"
            return-object
            hide-details="auto"
            :no-filter="true"
            autocomplete="off" />
        </div>
        <div class="sm:col-span-3 md:col-span-3 lg:col-span-3">
          <v-autocomplete
            :label="t('labelSearchLicense')"
            variant="outlined"
            clearable
            v-model="licenseSearch"
            density="compact"
            single-line
            @change="reloadAnalytics"
            @update:search="debouncedSearchLicenses"
            :loading="loadingLicenses"
            :items="licenses"
            return-object
            hide-details="auto"
            :no-filter="true"
            autocomplete="off" />
        </div>
      </div>
    </template>

    <template #table>
      <v-data-table
        :loading="loadingAnalytics"
        :headers="filteredHeaders"
        fixed-header
        density="compact"
        :items-per-page="50"
        :footer-props="{'items-per-page-options': [10, 50, 100, -1]}"
        :items="items"
        class="striped-table fill-height"
        :item-class="getCssClassForTableRow"
        v-model:options="options"
        :server-items-length="total">
        <template #[`header.settings`]="{column}">
          <span>{{ column.title }}</span>
          <v-menu offset-y :close-on-content-click="false" v-model="menuSettings">
            <template v-slot:activator="{props}">
              <DIconButton
                :parentProps="props"
                icon="mdi-cog"
                :hint="t('TT_TABLE_SETTINGS')"
                color="primary"
                class="pl-2" />
            </template>
            <div class="bg-background" style="width: 280px">
              <v-row class="d-flex ma-1 mr-2 justify-end">
                <DCloseButton @click="menuSettings = false" />
              </v-row>
              <v-select
                v-model="selectedHeaders"
                :items="selectableHeaders"
                class="pa-2 dp-select mx-2"
                :label="t('LABEL_SELECT_COL')"
                multiple
                return-object
                v-bind:menu-props="{location: 'bottom'}"
                variant="outlined"
                density="compact">
                <template v-slot:item="{props, item}">
                  <v-list-item v-bind="props" class="px-2 py-0">
                    <template v-slot:prepend="{isSelected}">
                      <v-checkbox hide-details :model-value="isSelected" />
                    </template>
                    <template v-slot:title>
                      <span class="pFilterEntry"> {{ allHeaders[item.value].title }}</span>
                    </template>
                  </v-list-item>
                </template>
                <template v-slot:selection="{item, index}">
                  <div v-if="index === 0" class="d-flex align-center">
                    <span class="pFilterEntry">{{ allHeaders[item.value].title }}</span>
                  </div>
                  <span v-if="index === 1" class="pAdditionalFilter"> +{{ selectedHeaders.length - 1 }} others </span>
                </template>
              </v-select>
            </div>
          </v-menu>
        </template>

        <template #[`item.lastUpdate`]="{item}">
          <DDateCellWithTooltip :value="item.lastUpdate" />
        </template>

        <template #[`item.name`]="{item}">
          <router-link target="_blank" :to="createProjectURL(item.key)">{{ item.name }}</router-link>
        </template>

        <template #[`item.projectVersionName`]="{item}">
          <router-link target="_blank" :to="createVersionURL(item.key, item.projectVersionKey)">{{
            item.projectVersionName
          }}</router-link>
        </template>

        <template #[`item.sbomName`]="{item}">
          <router-link target="_blank" :to="createSBOMURL(item.key, item.projectVersionKey, item.projectVersionKey)">{{
            item.sbomName
          }}</router-link>
        </template>

        <template #[`item.ownerCompanyName`]="{item}">
          <span v-if="!item.ownerCompanyMissing">{{ item.ownerCompanyName }}</span>
          <div v-else>
            <v-icon class="pr-2" color="warning" small>warning</v-icon>
            <span>{{ t('WARNING_MISSING_DEPT') }}</span>
          </div>
        </template>
      </v-data-table>
    </template>
  </TableLayout>
</template>
