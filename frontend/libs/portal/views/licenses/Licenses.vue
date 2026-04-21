<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import {useLicense} from '@disclosure-portal/composables/useLicense';
import {FilterSetDto} from '@disclosure-portal/model/FilterSet';
import {IObligation} from '@disclosure-portal/model/IObligation';
import License, {
  ClassificationWithCount,
  compareFamily,
  getLicenseApprovalTypeKeys,
  LicenseSlim,
  PossibleFilterValues,
} from '@disclosure-portal/model/License';
import {Group, Rights} from '@disclosure-portal/model/Rights';
import AdminService from '@disclosure-portal/services/admin';
import filterSetService from '@disclosure-portal/services/filtersets';
import licenseService from '@disclosure-portal/services/license';
import {useUserStore} from '@disclosure-portal/stores/user';
import {RightsUtils} from '@disclosure-portal/utils/Rights';
import {SearchOptions} from '@disclosure-portal/utils/Table';
import useViewTools, {getIconColorOfLevel, getIconOfLevel, openUrl} from '@disclosure-portal/utils/View';
import {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {useHeaderSettingsStore} from '@shared/stores/headerSettings.store';
import {DataTableHeader, DataTableHeaderFilterItems, DataTableItem} from '@shared/types/table';
import {debounce} from 'lodash';
import {storeToRefs} from 'pinia';
import {computed, nextTick, onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute, useRouter} from 'vue-router';
import {SortItem} from 'vuetify/lib/components/VDataTable/composables/sort';

interface FilterCondition {
  field: string;
  include: string[];
}

const {t} = useI18n();
const {getI18NTextOfPrefixKey} = useLicense();
const router = useRouter();
const breadcrumbs = useBreadcrumbsStore();
const snackbar = useSnackbar();
const route = useRoute();
const userStore = useUserStore();
const viewTools = useViewTools();

const gridName = 'License';
const headerSettingsStore = useHeaderSettingsStore();
const {filteredHeaders} = storeToRefs(headerSettingsStore);

const page = ref(1);
const sortItems = ref<SortItem[]>([{key: 'name', order: 'asc'}]);
const dataGridLicenses = ref<HTMLElement | null>(null);
const confirmVisible = ref(false);
const confirmConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const items = ref<LicenseSlim[]>([]);
const licensesLoading = ref(false);
const search = ref('');
const rights = ref<Rights>({} as Rights);
const selectedFilterIsLicenseChart = ref<string[]>([]);
const selectedFilterSource = ref<string[]>([]);
const selectedFilterFamily = ref<string[]>([]);
const selectedFilterApproval = ref<string[]>([]);
const selectedFilterType = ref<string[]>([]);
const selectedFilterClassification = ref<string[]>([]);
const total = ref(0);
const metaData = ref<PossibleFilterValues | null>(null);
const selectedFilterSet = ref<FilterSetDto | null>(null);
const filterSets = ref<FilterSetDto[]>([]);
const classifications = ref<IObligation[]>([]);
const dlgFilterSets = ref();
const dlgCompareLicense = ref();

const classificationsDialogRef = ref();
const configurePoliciesForLicenseDialogRef = ref();
const licenseDialogRef = ref();
const currentLicenseForAction = ref<License | null>(null);
const licenseDialogMode = ref<'edit' | 'duplicate'>('edit');
const abort = ref<AbortController | null>(null);
const itemsPerPage = ref(100);

const options = computed(
  (): SearchOptions =>
    ({
      page: page.value,
      itemsPerPage: itemsPerPage.value,
      sortBy: sortItems.value,
      groupBy: [],
      search: search.value,
      filterString: search.value,
      filterBy: {
        isLicenseChart: selectedFilterIsLicenseChart.value,
        source: selectedFilterSource.value,
        family: selectedFilterFamily.value,
        approvalState: selectedFilterApproval.value,
        licenseType: selectedFilterType.value,
        classifications: selectedFilterClassification.value,
      },
    }) as SearchOptions,
);

const possibleIsLicenseChart = computed((): DataTableHeaderFilterItems[] =>
  metaData.value
    ? Object.entries(metaData.value.possibleCharts)
        .map(([k, count]) => ({
          text: k === 'true' ? t('TABLE_LICENSE_CHART_STATUS_IS') : t('TABLE_LICENSE_CHART_STATUS_IS_NOT'),
          value: k,
          chip: String(count),
        }))
        .sort()
    : [],
);

const possibleSources = computed((): DataTableHeaderFilterItems[] =>
  metaData.value
    ? Object.entries(metaData.value.possibleSources).map(([k, count]) => ({
        text: k,
        value: k,
        chip: String(count),
      }))
    : [],
);

const possibleFamilies = computed((): DataTableHeaderFilterItems[] =>
  metaData.value
    ? Object.entries(metaData.value.possibleFamilies)
        .sort((a, b) => compareFamily(a[0], b[0]))
        .map(([k, count]) => ({
          text: getI18NTextOfPrefixKey('LIC_FAMILY_', k),
          value: k.length === 0 ? 'not declared' : k,
          chip: String(count),
        }))
    : [],
);

const possibleApproval = computed((): DataTableHeaderFilterItems[] =>
  metaData.value
    ? Object.entries(metaData.value.possibleApproval)
        .sort(
          ([keyA], [keyB]) => getLicenseApprovalTypeKeys().indexOf(keyA) - getLicenseApprovalTypeKeys().indexOf(keyB),
        )
        .map(([k, count]) => ({
          text: getI18NTextOfPrefixKey('LT_APP_', k),
          value: k.length === 0 ? 'not set' : k,
          chip: String(count),
        }))
    : [],
);

const possibleType = computed((): DataTableHeaderFilterItems[] =>
  metaData.value
    ? Object.entries(metaData.value.possibleType).map(([k, count]) => ({
        text: getI18NTextOfPrefixKey('LT_', k),
        value: k.length === 0 ? 'not declared' : k,
        chip: String(count),
      }))
    : [],
);

const possibleClassifications = computed((): DataTableHeaderFilterItems[] =>
  metaData.value
    ? metaData.value.possibleClassifications.map(({classification, count}: ClassificationWithCount) => {
        const value = viewTools.getNameForLanguage(classification) ? classification.name : '';
        return {
          text: viewTools.getNameForLanguage(classification) || t('NO_CLASSIFICATIONS'),
          value: value,
          icon: getIconOfLevel(getWarnLevel(value).toUpperCase()),
          iconColor: getIconColorOfLevel(getWarnLevel(value)),
          chip: String(count),
        };
      })
    : [],
);

const resetPagination = () => {
  page.value = 1;
};

const resetFilter = () => {
  resetPagination();
  items.value = [];
  selectedFilterSet.value = null;
  selectedFilterIsLicenseChart.value = [];
  selectedFilterSource.value = [];
  selectedFilterFamily.value = [];
  selectedFilterApproval.value = ['approved', 'deprecated', 'forbidden'];
  selectedFilterType.value = [];
  selectedFilterClassification.value = [];

  router.replace({path: '/dashboard/licenses'});
};

const redirectNotAllowed = () => {
  if (!rights.value.allowLicense?.read) {
    router.replace({path: '/dashboard/home'});
    return true;
  }
  return false;
};

const initBreadcrumbs = () => {
  breadcrumbs.setCurrentBreadcrumbs([
    {
      title: t('BC_Dashboard'),
      href: '/dashboard/home',
    },
    {
      title: t('BC_License'),
      href: '/dashboard/licenses/',
    },
  ]);
};

const searchForSimilarLicenseText = (licenseText: string) => {
  dlgCompareLicense.value?.search(licenseText);
};

const showDeletionConfirmationDialog = async (license: LicenseSlim) => {
  const r = await licenseService.getCountOfPolicyRuleUsingThisLicence(license.licenseId);
  const count = r.data.count;
  if (count > 0) {
    const userHasRightsToDeleteEvenIsUsed = false;
    if (userHasRightsToDeleteEvenIsUsed) {
      confirmConfig.value = {
        type: ConfirmationType.NOT_SET, // allows to separate in onConfirm Callback see below if not need set to ConfirmationType.NOT_SET
        key: license.licenseId,
        name: license.name,
        description: 'DLG_CONFIRMATION_DESCRIPTION',
        extendedDetails: t('LICENSE_IN_USE_BY_POLICY_RULES', {count: count}),
        okButton: 'Btn_delete',
      };
    } else {
      confirmConfig.value = {
        type: ConfirmationType.NOT_SET, // allows to separate in onConfirm Callback see below if not need set to ConfirmationType.NOT_SET
        key: license.licenseId,
        name: license.name,
        description: 'DLG_CAN_NOT_DELETE_IN_USE',
        extendedDetails: t('LICENSE_IN_USE_BY_POLICY_RULES', {count: count}),
        okButton: 'BTN_CLOSE',
        okButtonIsDisabled: true,
        title: 'DLG_WARNING_TITLE',
      };
    }
  } else {
    //  dlgConfirmDeleteProject.open(license.licenseId, license.name, 'DLG_CONFIRMATION_DESCRIPTION', 'Btn_delete');
    confirmConfig.value = {
      type: ConfirmationType.DELETE, // allows to separate in onConfirm Callback see below if not need set to ConfirmationType.NOT_SET
      key: license.licenseId,
      name: license.name,
      description: 'DLG_CONFIRMATION_DESCRIPTION',
      extendedDetails: '',
      okButton: 'Btn_delete',
      okButtonIsDisabled: false,
    };
  }
  confirmVisible.value = true;
};

const configurePoliciesForLicense = async (license: LicenseSlim) => {
  const fullModel: License = (await licenseService.get(license.licenseId)).data;
  configurePoliciesForLicenseDialogRef.value?.open(fullModel.licenseId, fullModel.name);
};

const doDeleteLicense = async (license: string) => {
  await licenseService.delete(license);
  snackbar.info(t('DIALOG_license_delete_success'));

  await reload();
};

const onClickRow = (event: Event, table: DataTableItem<LicenseSlim>) => {
  if (rights.value.allowLicense && rights.value.allowLicense.read) {
    openUrl('/dashboard/licenses/' + table.item.licenseId, router);
  }
};

const editLicense = async (license: LicenseSlim) => {
  currentLicenseForAction.value = (await licenseService.get(license.licenseId)).data;
  licenseDialogMode.value = 'edit';

  await nextTick();

  licenseDialogRef.value?.showDialog();
};

const duplicateLicense = async (license: LicenseSlim) => {
  currentLicenseForAction.value = (await licenseService.get(license.licenseId)).data;
  licenseDialogMode.value = 'duplicate';

  await nextTick();

  licenseDialogRef.value?.showDialog();
};

const onLicenseDialogClosed = async () => {
  currentLicenseForAction.value = null;

  await reload();
};

const getActionButtons = (item: LicenseSlim): TableActionButtonsProps['buttons'] => {
  const canEdit = rights.value?.allowLicense?.delete;
  const canCreate = rights.value?.allowLicense?.create;
  const canDelete = rights.value?.allowLicense?.delete && item.source !== 'spdx';
  const canConfigurePolicies =
    rights.value?.allowLicense?.read &&
    rights.value?.allowPolicy?.create &&
    rights.value?.allowPolicy?.read &&
    rights.value?.allowPolicy?.update &&
    rights.value?.allowPolicy?.delete;

  return [
    {
      icon: 'mdi-pencil',
      hint: t('TT_edit_license'),
      event: 'edit',
      show: canEdit,
    },
    {
      icon: 'mdi-content-copy',
      hint: t('LICENSE_COPY_BTN_TOOLTIP'),
      event: 'duplicate',
      show: canCreate,
    },
    {
      icon: 'mdi-bank-outline',
      hint: t('CONFIGURE_POLICIES_FOR_LICENSE_TOOLTIP'),
      event: 'configure',
      show: canConfigurePolicies,
    },
    {
      icon: 'mdi-delete',
      hint: t('TT_delete_license'),
      event: 'delete',
      show: canDelete,
    },
  ];
};

const allowActions = computed(() => RightsUtils.hasLicenseAccess() || RightsUtils.hasPolicyAccess());

const headers = computed((): DataTableHeader[] => [
  ...(allowActions.value
    ? [
        {
          title: 'COL_ACTIONS',
          align: 'center',
          width: 120,
          value: 'actions',
        } as DataTableHeader,
      ]
    : []),
  {
    title: 'COL_LICENSE_CHART_STATUS',
    tooltipText: 'TABLE_LICENSE_CHART_STATUS_TOOLTIP',
    align: 'center',
    value: 'meta.isLicenseChart',
    key: 'meta.isLicenseChart',
    width: 150,
    maxWidth: 150,
    sortable: true,
  },
  {
    title: 'CLASSIFICATIONS',
    tooltipText: 'LC_CLASSIFICATION_TT',
    align: 'center',
    width: 180,
    maxWidth: 180,
    value: 'meta.classifications',
    sortable: true,
  },
  {
    title: 'COL_LICENSE_NAME',
    tooltipText: 'COL_LICENSE_NAME_TOOLTIP',
    align: 'start',
    value: 'name',
    minWidth: 360,
    width: 360,
    sortable: true,
  },
  {
    title: 'COL_LICENSE_ID',
    tooltipText: 'COL_LICENSE_ID_TOOLTIP',
    align: 'start',
    value: 'licenseId',
    width: 200,
    sortable: true,
  },
  {
    title: 'COL_LICENSE_ALIASES',
    tooltipText: 'COL_LICENSE_ALIASES_TOOLTIP',
    align: 'start',
    value: 'aliases',
    width: 220,
  },
  {
    title: 'COL_APPROVAL_STATUS',
    tooltipText: 'COL_APPROVAL_STATUS_TOOLTIP',
    width: 190,
    value: 'meta.approvalState',
    sortable: true,
  },
  {
    title: 'COL_LICENSE_FAMILY',
    tooltipText: 'COL_LICENSE_FAMILY_TOOLTIP',
    align: 'start',
    value: 'meta.family',
    width: 190,
    sortable: true,
  },
  {
    title: 'COL_TYPE',
    tooltipText: 'COL_LICENSE_TYPE_TOOLTIP',
    width: 145,
    value: 'meta.licenseType',
    sortable: true,
  },
  {
    title: 'COL_LICENSE_SOURCE',
    tooltipText: 'COL_LICENSE_SOURCE_TOOLTIP',
    align: 'start',
    value: 'source',
    width: 130,
    maxWidth: 140,
    sortable: true,
  },
  {
    title: 'COL_UPDATED',
    align: 'start',
    width: 120,
    maxWidth: 130,
    value: 'updated',
    sortable: true,
  },
  {
    title: 'COL_CREATED',
    align: 'start',
    width: 120,
    maxWidth: 130,
    value: 'created',
    sortable: true,
  },
]);

headerSettingsStore.setupStore(gridName, headers.value);

const filterForCondition = (condition: FilterCondition) => {
  const filterAndMap = (
    possibleItems: DataTableHeaderFilterItems[],
    include: string[] = [],
    exclude: string[] = [],
  ) => {
    return possibleItems
      .filter((item) => {
        const value = item.value;
        const shouldInclude = include.includes(value);
        const shouldExclude = exclude.includes(value) || exclude.length === 0;
        return shouldInclude || !shouldExclude;
      })
      .map((item) => item.value);
  };
  switch (condition.field) {
    case 'isLicenseChart':
      selectedFilterIsLicenseChart.value = filterAndMap(possibleIsLicenseChart.value, condition.include);
      break;
    case 'source':
      selectedFilterSource.value = filterAndMap(possibleSources.value, condition.include);
      break;
    case 'family':
      selectedFilterFamily.value = filterAndMap(possibleFamilies.value, condition.include);
      break;
    case 'approvalState':
      selectedFilterApproval.value = filterAndMap(possibleApproval.value, condition.include);
      break;
    case 'licenseType':
      selectedFilterType.value = filterAndMap(possibleType.value, condition.include);
      break;
    case 'classifications':
      selectedFilterClassification.value = filterAndMap(possibleClassifications.value, condition.include);
      break;
  }
};
const applyFilterSet = async (filter: FilterSetDto) => {
  selectedFilterSet.value = filter;

  Object.keys(options.value.filterBy).forEach((key) => {
    const includedFilter = filter.includedFilters.find((filter) => filter.name === key);
    const filterValues = includedFilter ? includedFilter.values : [];

    filterForCondition({
      field: key,
      include: filterValues || [],
    });
  });

  await updateFilterSets();

  const newRoute = `/dashboard/licenses/filtersets/${encodeURIComponent(selectedFilterSet.value._key)}`;

  if (router.currentRoute.value.path !== newRoute) {
    await router.push(newRoute);
  }
};

const onFilterSetChange = async (oldVal: FilterSetDto) => {
  if (
    selectedFilterSet.value &&
    selectedFilterSet.value._key &&
    (oldVal === null || oldVal._key != selectedFilterSet.value._key)
  ) {
    selectedFilterSet.value = await filterSetService.getFilterSet(selectedFilterSet.value._key);
    await applyFilterSet(selectedFilterSet.value);
  }
};

const reloadFilter = async (filterKey: string | string[]) => {
  filterSets.value = await getSortedFilterSets();
  if (filterKey && filterKey.length > 0) {
    const finds = filterSets.value.find((f) => f._key === filterKey);
    if (finds) {
      selectedFilterSet.value = finds;
      await applyFilterSet(selectedFilterSet.value);
    } else {
      resetFilter();
    }
  }
};

const getSortedFilterSets = async () => {
  const filterSets = await filterSetService.getFilterSets('licenses');
  return filterSets.sort((a, b) => {
    const nameA = a.name.toLowerCase();
    const nameB = b.name.toLowerCase();
    if (nameA < nameB) return -1;
    if (nameA > nameB) return 1;
    return 0;
  });
};

const sendSelectedFilters = () => {
  if (dlgFilterSets.value) {
    dlgFilterSets.value.setFilterData(options.value.filterBy);
  }
};
const onConfirm = async (config: IConfirmationDialogConfig) => {
  if (config.okButtonIsDisabled) return;
  if (config.type === ConfirmationType.NOT_SET) {
    // do nothing
  } else if (config.type === ConfirmationType.DELETE) {
    await doDeleteLicense(config.key);
  }
};

const updateFilterSets = async () => {
  filterSets.value = await getSortedFilterSets();
};

const openClassifications = (classifications: IObligation[], licenseName: string, licenseId: string) => {
  if (classificationsDialogRef.value) {
    classificationsDialogRef.value?.open(classifications, licenseName, licenseId);
  }
};

const getWarnLevel = (name: string) => {
  const classification = classifications.value.find((c) => c.name === name || c.nameDe === name);
  return classification ? classification.warnLevel : 'INFORMATION';
};
const reactiveTotal = computed(() => total.value);

const retrieveClassifications = async () => {
  const response = (await AdminService.getAllObligations()).data;
  classifications.value = response.items;
};

const headerExpands = () => {
  headerSettingsStore.setupStore(gridName, headers.value);
};

const reload = async () => {
  if (abort.value) {
    abort.value.abort();
  }

  licensesLoading.value = true;

  abort.value = new AbortController();

  const {licenses, count, meta} = (await licenseService.search(options.value, abort.value.signal)).data;

  abort.value = null;

  items.value = licenses;
  total.value = count;
  metaData.value = meta;

  licensesLoading.value = false;
};

const searchChanged = async () => {
  if (search.value && search.value.length > 80) {
    return;
  }

  resetPagination();

  await reload();
};

const debouncedSearch = debounce(searchChanged, 300);

watch(selectedFilterSet, (_, oldVal) => onFilterSetChange(oldVal!));

watch(
  [
    selectedFilterIsLicenseChart,
    selectedFilterSource,
    selectedFilterFamily,
    selectedFilterApproval,
    selectedFilterType,
    selectedFilterClassification,
    search,
  ],
  debouncedSearch,
  {deep: true},
);

const onUpdateOptions = async (tableOptions: {page: number; itemsPerPage: number; sortBy: SortItem[]}) => {
  page.value = tableOptions.page;
  await reload();
};

onMounted(async () => {
  rights.value = userStore.getRights;

  if (redirectNotAllowed()) {
    return;
  }

  await retrieveClassifications();

  filterSets.value = await getSortedFilterSets();

  initBreadcrumbs();

  if (router.currentRoute.value.path.includes('compare')) {
    const licenseText = localStorage.getItem('licenseText');
    if (licenseText) {
      searchForSimilarLicenseText(licenseText);
    }
  }

  await reload();

  if (route.path.includes('filtersets') && route.params.id) {
    await reloadFilter(route.params.id);
  } else {
    await reloadFilter('');
  }
});
</script>

<template>
  <TableLayout data-testid="licenses">
    <template #description>
      <div class="d-flex align-center ga-3 flex-row">
        <h1 class="text-h5">{{ t('Licenses') }}</h1>
        <NewOrEditLicenseDialog
          mode="create"
          @closed:successfully="onLicenseDialogClosed"
          v-slot="{showDialog}"
          v-if="rights && rights.allowLicense && rights.allowLicense.create">
          <DCActionButton
            large
            class="mx-2"
            :text="t('BTN_ADD')"
            icon="mdi-plus"
            :hint="t('TT_add_license')"
            @click="showDialog" />
        </NewOrEditLicenseDialog>
        <LicenseCompareDialog v-slot="{showDialog}" ref="dlgCompareLicense">
          <DCActionButton
            :text="t('Btn_compare')"
            large
            icon="mdi-shield-search"
            :hint="t('TT_search_license_text')"
            @click="showDialog" />
        </LicenseCompareDialog>
        <v-spacer></v-spacer>
        <v-select
          v-if="!rights?.groups?.includes(Group.UserLicenseManager)"
          v-model="selectedFilterSet"
          :item-title="['name']"
          :items="filterSets"
          :label="t('FILTER_SET_LABEL')"
          variant="outlined"
          density="compact"
          clearable
          @click:clear="resetFilter"
          hide-details
          return-object></v-select>
        <DSearchField v-model="search" />
      </div>
    </template>
    <template v-if="rights?.groups?.includes(Group.UserLicenseManager)" #buttons>
      <v-spacer></v-spacer>
      <FilterSets
        ref="dlgFilterSets"
        @requestFilterDataFromTable="sendSelectedFilters"
        @reloadFilter="reloadFilter"></FilterSets>
      <v-select
        style="max-width: 500px"
        v-model="selectedFilterSet"
        :item-title="['name']"
        :items="filterSets"
        :label="t('FILTER_SET_LABEL')"
        variant="outlined"
        density="compact"
        clearable
        @click:clear="resetFilter"
        hide-details="auto"
        v-bind:menu-props="{location: 'bottom'}"
        return-object></v-select>
    </template>
    <template #table>
      <div ref="dataGridLicenses" class="table-wrapper fill-height">
        <v-data-table-server
          :headers="filteredHeaders"
          fixed-header
          :items-length="reactiveTotal"
          :loading="licensesLoading"
          :items="items"
          :footer-props="{'items-per-page-options': [10, 50, 100, -1]}"
          density="compact"
          show-footer
          class="striped-table custom-data-table fill-height"
          v-model:page="page"
          v-model:items-per-page="itemsPerPage"
          v-model:sort-by="sortItems"
          @update:options="onUpdateOptions"
          @click:row="onClickRow">
          <!-- Settings column header slot -->
          <template v-if="allowActions" #[`header.actions`]="{column}">
            <GridFilterHeader :column="column">
              <template #settings>
                <HeaderSettings :grid-name="gridName" :column="column" />
              </template>
            </GridFilterHeader>
          </template>
          <template #[`header.meta.isLicenseChart`]="{column, getSortIcon, toggleSort}">
            <GridFilterHeader :column="column" :getSortIcon="getSortIcon" :toggleSort="toggleSort">
              <template #settings>
                <HeaderSettings v-if="!allowActions" :grid-name="gridName" :column="column" />
              </template>
              <template #filter>
                <GridHeaderFilterIcon
                  v-model="selectedFilterIsLicenseChart"
                  :column="column"
                  :label="t('LICENSE_CHART_STATUS')"
                  :allItems="possibleIsLicenseChart">
                </GridHeaderFilterIcon>
              </template>
            </GridFilterHeader>
          </template>
          <template #[`header.source`]="{column, getSortIcon, toggleSort}">
            <GridFilterHeader :column="column" :getSortIcon="getSortIcon" :toggleSort="toggleSort">
              <template #filter>
                <GridHeaderFilterIcon
                  v-model="selectedFilterSource"
                  :column="column"
                  :label="t('SOURCE')"
                  :allItems="possibleSources">
                </GridHeaderFilterIcon>
              </template>
            </GridFilterHeader>
          </template>
          <template #[`header.meta.family`]="{column, getSortIcon, toggleSort}">
            <GridFilterHeader :column="column" :getSortIcon="getSortIcon" :toggleSort="toggleSort">
              <template #filter>
                <GridHeaderFilterIcon
                  v-model="selectedFilterFamily"
                  :column="column"
                  :label="t('LICENSE_FAMILY')"
                  :allItems="possibleFamilies">
                </GridHeaderFilterIcon>
              </template>
            </GridFilterHeader>
          </template>
          <template #[`header.meta.approvalState`]="{column, getSortIcon, toggleSort}">
            <GridFilterHeader :column="column" :getSortIcon="getSortIcon" :toggleSort="toggleSort">
              <template #filter>
                <GridHeaderFilterIcon
                  v-model="selectedFilterApproval"
                  :column="column"
                  :label="t('APPROVAL_STATUS')"
                  :initial-selected="['approved', 'deprecated', 'forbidden']"
                  :allItems="possibleApproval">
                </GridHeaderFilterIcon>
              </template>
            </GridFilterHeader>
          </template>
          <template #[`header.meta.licenseType`]="{column, getSortIcon, toggleSort}">
            <GridFilterHeader :column="column" :getSortIcon="getSortIcon" :toggleSort="toggleSort">
              <template #filter>
                <GridHeaderFilterIcon
                  v-model="selectedFilterType"
                  :column="column"
                  :label="t('LICENSE_TYPE')"
                  :allItems="possibleType">
                </GridHeaderFilterIcon>
              </template>
            </GridFilterHeader>
          </template>
          <template #[`header.meta.classifications`]="{column, getSortIcon, toggleSort}">
            <GridFilterHeader :column="column" :getSortIcon="getSortIcon" :toggleSort="toggleSort">
              <template #filter>
                <GridHeaderFilterIcon
                  v-model="selectedFilterClassification"
                  :column="column"
                  :label="t('CLASSIFICATION')"
                  :allItems="possibleClassifications">
                </GridHeaderFilterIcon>
              </template>
            </GridFilterHeader>
          </template>
          <template #[`item.meta.classifications`]="{item}">
            <span @click.stop="openClassifications(item.meta.classifications, item.name, item.licenseId)">
              <v-icon color="primary" size="small" class="mr-2">mdi-chevron-right</v-icon>
              <v-icon
                :class="item.meta.prevalentClassificationLevel.toUpperCase() === 'WARNING' ? 'mr-1' : 'mr-2'"
                :color="getIconColorOfLevel(item.meta.prevalentClassificationLevel)"
                >{{ getIconOfLevel(item.meta.prevalentClassificationLevel) }}
              </v-icon>
              <Tooltip location="bottom">
                {{ t('TT_OPEN_CLASSIFICATIONS', {license: item.name}) }}
              </Tooltip>
            </span>
          </template>
          <template #[`item.meta.isLicenseChart`]="{item}">
            <DLicenseChartIcon :meta="item.meta" />
          </template>
          <template #[`item.name`]="{item}">
            {{ item.name }}
          </template>
          <template #[`item.licenseId`]="{item}">
            {{ item.licenseId }}
          </template>
          <template #[`item.meta.family`]="{item}">
            {{ getI18NTextOfPrefixKey('LIC_FAMILY_', '' + item.meta.family) }}
          </template>
          <template #[`item.meta.approvalState`]="{item}">
            {{ getI18NTextOfPrefixKey('LT_APP_', item.meta.approvalState) }}
          </template>
          <template #[`item.meta.licenseType`]="{item}">
            {{ getI18NTextOfPrefixKey('LT_', item.meta.licenseType) }}
          </template>
          <template #[`item.created`]="{item}">
            <DDateCellWithTooltip :value="item.created" />
          </template>
          <template #[`item.updated`]="{item}">
            <DDateCellWithTooltip :value="item.updated" />
          </template>
          <template v-if="allowActions" #[`item.actions`]="{item}">
            <TableActionButtons
              variant="compact"
              :buttons="getActionButtons(item)"
              @edit="editLicense(item)"
              @duplicate="duplicateLicense(item)"
              @delete="showDeletionConfirmationDialog(item)"
              @configure="configurePoliciesForLicense(item)"
              @slideToggle="headerExpands" />
          </template>
          <template #[`item.aliases`]="{item}">
            <span v-if="Array.isArray(item.aliases) && item.aliases.length > 0">
              {{ item.aliases.map((a) => a.licenseId).join(', ') }}
            </span>
            <span v-else>-</span>
          </template>
        </v-data-table-server>
      </div>
    </template>
  </TableLayout>
  <ConfirmationDialog
    v-model:showDialog="confirmVisible"
    :config="confirmConfig"
    @confirm="onConfirm"></ConfirmationDialog>
  <ClassificationsPerLicenseDialog ref="classificationsDialogRef"></ClassificationsPerLicenseDialog>
  <ConfigurePoliciesForLicenseDialog ref="configurePoliciesForLicenseDialogRef"></ConfigurePoliciesForLicenseDialog>
  <NewOrEditLicenseDialog
    v-if="currentLicenseForAction"
    ref="licenseDialogRef"
    :initial-data="currentLicenseForAction"
    :mode="licenseDialogMode"
    @closed:successfully="onLicenseDialogClosed">
  </NewOrEditLicenseDialog>
</template>
