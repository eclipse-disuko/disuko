<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {AuditLog} from '@disclosure-portal/model/VersionDetails';
import {escapeHtml} from '@disclosure-portal/utils/Validation';
import {DataTableHeader, SortItem} from '@shared/types/table';
import {computed, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const props = defineProps<{
  fetchMethod: () => Promise<AuditLog[]>;
}>();

const items = ref<AuditLog[]>([]);
const headers = computed<DataTableHeader[]>(() => [
  {
    title: '',
    value: 'data-table-expand',
    width: 20,
  },
  {
    title: t('COL_TITLE'),
    width: 200,
    align: 'start',
    value: 'message',
    sortable: true,
  },
  {
    title: t('COL_CREATED'),
    width: 160,
    align: 'start',
    value: 'created',
    sortable: true,
  },
  {
    width: 200,
    title: t('COL_USER'),
    sortable: false,
    align: 'start',
    value: 'user',
  },
]);

const expanded = ref<string[]>([]);
const search = ref('');
const dataAreLoaded = ref(false);
const sortBy: SortItem[] = [{key: 'created', order: 'desc'}];
const tableAuditLog = ref<HTMLElement | null>(null);

const reload = async (forceReload = false) => {
  if (!forceReload && dataAreLoaded.value) return;
  dataAreLoaded.value = false;
  items.value = await props.fetchMethod();
  dataAreLoaded.value = true;
};
onMounted(async () => {
  await reload(true);
});

const onRowExpand = (newExpanded: string[]) => {
  if (newExpanded.length > 1) {
    // Keep only the last expanded row
    expanded.value = [newExpanded[newExpanded.length - 1]];
  } else {
    expanded.value = newExpanded;
  }
};

const toggleExpand = (item: AuditLog) => {
  const index = expanded.value.indexOf(item._key);
  if (index > -1) {
    expanded.value.splice(index, 1);
  } else {
    expanded.value.push(item._key);
  }
};

const isExpanded = (item: AuditLog) => {
  return expanded.value.includes(item._key);
};
</script>

<template>
  <TableLayout has-tab has-title>
    <template #buttons>
      <v-spacer></v-spacer>
      <DSearchField v-model="search" />
    </template>
    <template #table>
      <div ref="tableAuditLog" class="fill-height">
        <v-data-table
          :loading="!dataAreLoaded"
          item-value="_key"
          :items="items"
          :headers="headers"
          :search="search"
          fixed-header
          class="striped-table custom-data-table fill-height"
          density="compact"
          :sort-by="sortBy"
          expand-on-click
          :expanded.sync="expanded"
          @update:expanded="onRowExpand">
          <template v-slot:item.data-table-expand="{item}">
            <v-icon color="primary" @click.stop="toggleExpand(item)">
              {{ isExpanded(item) ? 'mdi-chevron-up' : 'mdi-chevron-down' }}
            </v-icon>
          </template>
          <template v-slot:expanded-row="{columns, item}">
            <tr>
              <td :colspan="columns.length">
                <pre v-html="escapeHtml(item.meta)" class="auditPre"></pre>
              </td>
            </tr>
          </template>
          <template v-slot:item.created="{item}">
            <DDateCellWithTooltip :value="item.created.toString()" />
          </template>
        </v-data-table>
      </div>
    </template>
  </TableLayout>
</template>
