<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import Icons from '@disclosure-portal/constants/icons';
import PolicyRule from '@disclosure-portal/model/PolicyRule';
import {UserDto} from '@disclosure-portal/model/Users';
import adminService from '@disclosure-portal/services/admin';
import {getCssClassForTableRow, SearchOptions} from '@disclosure-portal/utils/Table';
import {openUrl} from '@disclosure-portal/utils/url';
import {DataTableHeader, DataTableHeaderFilterItems, DataTableItem, SortItem} from '@shared/types/table';
import _ from 'lodash';
import {computed, nextTick, onBeforeMount, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';

const {t} = useI18n();
const router = useRouter();

const search = ref('');
const icons = Icons;
const options = ref<SearchOptions>({} as SearchOptions);
const total = ref(0);
const loading = ref(false);
const users = ref<UserDto[]>([]);
const sortBy = ref<SortItem[]>([{key: 'user', order: 'asc'}]);
const abort = ref<AbortController | null>(null);
const tableAllUsers = ref<HTMLElement | null>(null);
const searchField = ref<InstanceType<typeof DSearchField> | null>(null);
const selectedFilterStatus = ref<string[]>([]);
const selectedFilterScopes = ref<string[]>([]);

const possibleFilterScopes = computed((): DataTableHeaderFilterItems[] => [
  {text: t('USER_ACCESS_SCOPE_INTERNAL'), value: 'true'},
  {text: t('USER_ACCESS_SCOPE_EXTERNAL'), value: 'false'},
]);
const possibleFilterStatus = computed((): DataTableHeaderFilterItems[] => [
  {text: t('ICON_LABEL_TEXT_ACTIVE'), value: 'true'},
  {text: t('ICON_LABEL_TEXT_INACTIVE'), value: 'false'},
]);
const reactiveTotal = computed(() => total.value);

const headers = computed<DataTableHeader[]>(() => [
  {
    title: t('COL_STATUS'),
    align: 'center',
    value: 'active',
    sortable: true,
    width: 140,
    maxWidth: 150,
  },
  {
    title: t('COL_USER_ID'),
    align: 'start',
    width: 120,
    maxWidth: 160,
    value: 'user',
    sortable: true,
  },
  {
    title: t('COL_USER_ACCESS_SCOPE'),
    align: 'start',
    width: 180,
    value: 'isInternal',
    sortable: true,
  },
  {
    title: t('COL_FORENAME'),
    align: 'start',
    width: 140,
    minWidth: 140,
    value: 'forename',
    sortable: true,
  },
  {
    title: t('COL_LASTNAME'),
    width: 140,
    minWidth: 140,
    align: 'start',
    value: 'lastname',
    sortable: true,
  },
  {
    title: t('COL_EMAIL'),
    align: 'start',
    width: 300,
    minWidth: 300,
    value: 'email',
    sortable: true,
  },
  {
    title: t('DEPARTMENT'),
    align: 'start',
    width: 112,
    value: 'metaData.department',
    sortable: true,
  },
  {
    title: t('DEPARTMENT_DESCRIPTION'),
    align: 'start',
    width: 200,
    value: 'metaData.departmentDescription',
    sortable: true,
  },
  {
    title: t('COMPANY_IDENTIFIER'),
    align: 'start',
    width: 188,
    value: 'metaData.companyIdentifier',
    sortable: true,
  },
  {
    title: t('COL_CREATED'),
    align: 'start',
    width: 110,
    maxWidth: 120,
    value: 'created',
    sortable: true,
  },
  {
    title: t('COL_UPDATED'),
    align: 'start',
    width: 110,
    maxWidth: 120,
    value: 'updated',
    sortable: true,
  },
  {
    title: t('DEPROVISIONED_DATE'),
    align: 'start',
    width: 180,
    value: 'deprovisioned',
    sortable: true,
  },
  {
    title: t('COL_TERMS_DATE'),
    width: 220,
    align: 'start',
    value: 'termsOfUseDate',
    sortable: true,
  },
  {
    title: t('COL_TERMS_VERSION'),
    width: 220,
    align: 'start',
    value: 'termsOfUseVersion',
    sortable: true,
  },
  {
    title: t('COL_TERMS_ACCEPTANCE'),
    width: 220,
    align: 'center',
    value: 'termsOfUse',
    sortable: true,
  },
]);

const loadData = async () => {
  loading.value = true;
  options.value.filterString = search.value;
  options.value.filterBy = {
    isActive: selectedFilterStatus.value,
    isInternal: selectedFilterScopes.value,
  };

  abort.value = new AbortController();

  const {items, count} = (await adminService.getUsersWithOptions(options.value, abort.value.signal)).data;

  abort.value = null;

  users.value = items;
  total.value = count;
  loading.value = false;
};

const reload = async () => {
  if (abort.value) {
    abort.value.abort();
  }

  await loadData();
};

const filterChanged = ref(false);
const optionsChanged = async () => {
  if (filterChanged.value) {
    filterChanged.value = false;
    return;
  }
  await reload();
};

const searchChanged = async () => {
  if (search.value && search.value.length > 80) {
    return;
  }
  setFilterChangedAndResetPagination();
  if (!loading.value) {
    loading.value = true;
    await nextTick(async () => {
      await reload();
    });
  }
};

const setFilterChangedAndResetPagination = () => {
  filterChanged.value = true;
  if (options.value.page > 1) options.value.page = 1;
};

const onClickRow = (_: Event, table: DataTableItem<PolicyRule>) => {
  openUrl('/dashboard/admin/users/' + table.item._key, router);
};

const wait = 300;
const debouncedOptions = _.debounce(() => optionsChanged(), wait);
const debouncedSearch = _.debounce(() => searchChanged(), wait);

onBeforeMount(async () => {
  await reload();
});

watch(search, debouncedSearch);
watch([options, selectedFilterStatus, selectedFilterScopes], debouncedOptions);

watch(loading, async (newValue) => {
  if (!newValue) {
    await nextTick(() => {
      searchField.value?.focus();
    });
  }
});
</script>

<template>
  <TableLayout>
    <template #buttons>
      <h1 class="text-h5">{{ t('TITLE_USERS') }}</h1>
      <v-spacer></v-spacer>
      <DSearchField ref="searchField" v-model="search" :disabled="loading" />
    </template>
    <template #table>
      <div ref="tableAllUsers" class="fill-height">
        <v-data-table-server
          item-key="_key"
          :items="users"
          :headers="headers"
          :sort-by="sortBy"
          @click:row="onClickRow"
          fixed-header
          density="compact"
          class="striped-table fill-height"
          :item-class="getCssClassForTableRow"
          items-per-page="100"
          :footer-props="{'items-per-page-options': [10, 50, 100, -1]}"
          :items-length="reactiveTotal"
          :loading="loading"
          v-model:options="options"
          @update:options="reload">
          <template #[`header.active`]="{column, getSortIcon, toggleSort}">
            <span class="mr-1">{{ column.title }}</span>
            <GridHeaderFilterIcon
              v-model="selectedFilterStatus"
              :column="column"
              :label="t('COL_USER_ACCESS_SCOPE')"
              :allItems="possibleFilterStatus">
            </GridHeaderFilterIcon>
            <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
          </template>
          <template #[`header.isInternal`]="{column, getSortIcon, toggleSort}">
            <span class="mr-1">{{ column.title }}</span>
            <GridHeaderFilterIcon
              v-model="selectedFilterScopes"
              :column="column"
              :label="t('COL_USER_ACCESS_SCOPE')"
              :allItems="possibleFilterScopes">
            </GridHeaderFilterIcon>
            <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
          </template>
          <template #[`item.active`]="{item}">
            <div class="flex justify-center">
              <v-icon size="x-small" :color="item.active ? 'success' : 'warning'">{{ icons.CIRCLE_FILLED }}</v-icon>
            </div>
          </template>
          <template #[`item.isInternal`]="{item}">
            <span>{{ item.isInternal ? t('USER_ACCESS_SCOPE_INTERNAL') : t('USER_ACCESS_SCOPE_EXTERNAL') }}</span>
          </template>

          <template #[`item.updated`]="{item}">
            <DDateCellWithTooltip :value="item.updated" />
          </template>
          <template #[`item.deprovisioned`]="{item}">
            <DDateCellWithTooltip v-if="item.deprovisioned" :value="item.deprovisioned" />
          </template>
          <template #[`item.created`]="{item}">
            <DDateCellWithTooltip :value="item.created" />
          </template>
          <template #[`item.termsOfUseDate`]="{item}">
            <DDateCellWithTooltip :value="item.termsOfUseDate" />
          </template>
          <template #[`item.termsOfUse`]="{item}">
            <div class="flex justify-center">
              <v-icon size="small" :color="item.termsOfUse ? 'primary' : 'greyCheck'">mdi-check</v-icon>
            </div>
          </template>
        </v-data-table-server>
      </div>
    </template>
  </TableLayout>
</template>
