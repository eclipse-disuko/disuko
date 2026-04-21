<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ProjectRoleDto} from '@disclosure-portal/model/Users';
import {openUrl} from '@disclosure-portal/utils/View';
import TableLayout from '@shared/layouts/TableLayout.vue';
import {SortItem} from '@shared/types/table';
import {computed, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';
import {DataTableItem} from 'vuetify/lib/components/VDataTable/types';

interface Props {
  fetchMethod: () => Promise<ProjectRoleDto[]>;
}

const props = defineProps<Props>();

const {t} = useI18n();
const router = useRouter();

const sortBy = ref<SortItem[]>([]);
const items = ref<ProjectRoleDto[]>([]);
const search = ref('');
const dataAreLoaded = ref(false);

const reloadInternal = async (forceReload: boolean) => {
  if (!forceReload && dataAreLoaded.value) return;
  dataAreLoaded.value = false;
  items.value = await props.fetchMethod();
  dataAreLoaded.value = true;
};

const openProject = (item: ProjectRoleDto) => {
  const url = `/dashboard/projects/${encodeURIComponent(item.projectKey)}`;
  openUrl(url, router);
};

const headers = computed(() => {
  return [
    {
      title: t('COL_PROJECT_NAME'),
      align: 'start',
      class: 'tableHeaderCell',
      value: 'projectName',
      width: 180,
      sortable: true,
    },
    {
      title: t('COL_USER_TYPE'),
      width: 180,
      align: 'start',
      class: 'tableHeaderCell',
      value: 'userType',
      sortable: true,
    },
    {
      width: 200,
      title: t('ROLES'),
      align: 'start',
      class: 'tableHeaderCell',
      value: 'responsible',
      sortable: true,
    },
    {
      title: '',
      class: 'tableHeaderCell',
      value: '',
    },
  ];
});

onMounted(() => {
  sortBy.value = [{key: 'updated', order: 'desc'}];
  reloadInternal(true);
});
</script>

<template>
  <TableLayout has-tab has-title>
    <template #buttons>
      <v-spacer></v-spacer>
      <DSearchField v-model="search" />
    </template>
    <template #table>
      <v-data-table
        density="compact"
        class="striped-table fill-height"
        :loading="!dataAreLoaded"
        item-key="_key"
        :items="items"
        :headers="headers"
        :search="search"
        @click:row="(event: Event, dataItem: DataTableItem<ProjectRoleDto>) => openProject(dataItem.item)"
        :items-per-page="50"
        fixed-header
        :sort-by="sortBy"
        sort-desc>
        <template v-slot:[`item.responsible`]="{item}">
          <span v-if="item.responsible">{{ t('COL_USER_ROLE_RESPONSIBLE') }}</span>
        </template>
        <template v-slot:[`item.userType`]="{item}">
          <span>{{ t(`USER_ROLE_${item.userType}`) }}</span>
        </template>
      </v-data-table>
    </template>
  </TableLayout>
</template>
