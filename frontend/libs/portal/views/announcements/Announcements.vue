// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {IAnnouncement} from '@disclosure-portal/model/AnnouncementsResponse';
import announcementService from '@disclosure-portal/services/announcements';
import {RightsUtils} from '@disclosure-portal/utils/Rights';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {DataTableHeader, DataTableHeaderFilterItems, SortItem} from '@shared/types/table';
import {computed, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const breadcrumbs = useBreadcrumbsStore();

const announcements = ref<IAnnouncement[]>([]);
const search = ref('');
const selectedAttributeFilters = ref<string[]>([]);
const gridAnnouncement = ref<HTMLElement | null>(null);
const sortItems = ref<SortItem[]>([{key: 'when', order: 'desc'}]);

const headers = computed((): DataTableHeader[] => [
  ...(RightsUtils.rights().isInternal
    ? [
        {
          title: t('COL_ACTIONS'),
          align: 'center',
          width: 100,
          maxWidth: 110,
          value: 'actions',
          sortable: false,
        } as DataTableHeader,
      ]
    : []),
  {
    title: t('ANNOUNCEMENTS_COL_WHEN'),
    align: 'start',
    width: 140,
    maxWidth: 150,
    value: 'when',
    sortable: true,
  },
  {
    title: t('ANNOUNCEMENTS_COL_LICENSE_NAME'),
    align: 'start',
    width: 200,
    minWidth: 200,
    maxWidth: 220,
    value: 'content.licenseName',
    sortable: true,
  },
  {
    title: t('ANNOUNCEMENTS_COL_LICENSE_ID'),
    align: 'start',
    width: 180,
    minWidth: 180,
    value: 'content.licenseId',
    sortable: true,
  },
  {
    title: t('ANNOUNCEMENTS_COL_ATTRIBUTE'),
    align: 'start',
    width: 220,
    minWidth: 220,
    value: 'content.changeType',
    sortable: true,
  },
  {
    title: t('ANNOUNCEMENTS_COL_OLD_VALUE'),
    align: 'start',
    width: 120,
    maxWidth: 140,
    value: 'content.oldVal',
    sortable: true,
  },
  {
    title: t('ANNOUNCEMENTS_COL_NEW_VALUE'),
    align: 'start',
    width: 120,
    maxWidth: 140,
    value: 'content.newVal',
    sortable: true,
  },
]);

const filteredAnnouncements = computed(() =>
  announcements.value.filter(
    (announcement) =>
      selectedAttributeFilters.value.length === 0 ||
      selectedAttributeFilters.value.includes(announcement.content.changeType),
  ),
);

const possibleAttributeFilters = computed((): DataTableHeaderFilterItems[] => {
  const uniqueValues = [...new Set(announcements.value.map((ann) => ann.content.changeType))];
  return uniqueValues.map((val) => ({
    text: t('ANNOUNCEMENTS_TYPE_' + val),
    value: val,
  }));
});

const initBreadcrumbs = () => {
  breadcrumbs.setCurrentBreadcrumbs([
    {title: 'Dashboard', disabled: false, href: '/dashboard/home'},
    {title: 'Announcements', disabled: false, href: '/dashboard/announcements/'},
  ]);
};

const reloadAnnouncements = async () => {
  const announcementsUnparsed = (await announcementService.getAll()).data;
  if (!announcementsUnparsed) return;

  announcements.value = announcementsUnparsed.map((ann: any) => {
    ann.content = JSON.parse(ann.content);
    return ann;
  });
};

onMounted(async () => {
  initBreadcrumbs();
  await reloadAnnouncements();
});
</script>

<template>
  <TableLayout>
    <template #buttons>
      <h2 class="text-h5">{{ t('ANNOUNCEMENTS_HEADLINE') }}</h2>
      <v-spacer></v-spacer>
      <v-text-field
        autocomplete="off"
        :max-width="500"
        append-inner-icon="mdi-magnify"
        variant="outlined"
        density="compact"
        v-model="search"
        :label="t('labelSearch')"
        single-line
        hide-details></v-text-field>
    </template>
    <template #table>
      <div ref="gridAnnouncement" class="fill-height">
        <v-data-table
          density="compact"
          class="striped-table fill-height"
          fixed-header
          items-per-page="100"
          :headers="headers"
          :items="filteredAnnouncements"
          :search="search"
          :sort-by="sortItems">
          <template #[`item.when`]="{item}">
            <DDateCellWithTooltip :value="item.when" />
          </template>
          <template #[`header.content.changeType`]="{column, toggleSort, getSortIcon}">
            <span class="mr-1">{{ column.title }}</span>
            <GridHeaderFilterIcon
              v-model="selectedAttributeFilters"
              :column="column"
              :label="t('ANNOUNCEMENTS_COL_ATTRIBUTE')"
              :allItems="possibleAttributeFilters">
            </GridHeaderFilterIcon>
            <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
          </template>

          <template #[`item.content.changeType`]="{item}">
            {{ t('ANNOUNCEMENTS_TYPE_' + item.content.changeType) }}
          </template>
          <template #[`item.content.oldVal`]="{item}">
            <span v-if="item.content.changeType !== 'custom_license_deleted'">{{
              t(item.content.oldVal || 'not declared')
            }}</span>
          </template>
          <template #[`item.content.newVal`]="{item}">
            <span v-if="item.content.changeType !== 'custom_license_deleted'">{{
              t(item.content.newVal || 'not declared')
            }}</span>
          </template>
          <template #[`item.actions`]="{item}">
            <router-link
              v-if="item.content.changeType !== 'custom_license_deleted'"
              :to="'/dashboard/licenses/' + item.content.licenseId"
              target="_blank">
              <v-icon color="primary" size="large" icon="mdi-open-in-new"></v-icon>
            </router-link>
          </template>
        </v-data-table>
      </div>
    </template>
  </TableLayout>
</template>
