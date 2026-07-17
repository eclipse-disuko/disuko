<script setup lang="ts">
import type {AugmentedScanRemark, ScanRemark} from '@cli/models/Sbom';
import {compareScanRemarkStatus} from '@cli/models/Sbom';
import {useView} from '@disclosure-portal/composables/useView';
import {getIconColorScanRemarkLevel} from '@disclosure-portal/utils/View';
import Tooltip from '@shared/components/disco/Tooltip.vue';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import type {DataTableHeader, DataTableSortItem} from 'vuetify';

const props = defineProps<{
  scanRemarks: ScanRemark[];
}>();

const {t} = useI18n();
const {getTextOfLevel} = useView();

function statusIcon(_raw: string | undefined | null): string {
  return 'mdi-circle';
}

function statusColor(raw: string | undefined | null): string {
  const v = raw || '';
  return getIconColorScanRemarkLevel(v as any) || 'gray';
}

function statusText(raw: string | undefined | null): string {
  return getTextOfLevel(raw || '');
}

const truncateText = (text: string | undefined, maxLength: number = 120): string => {
  if (!text) return '';
  return text.length > maxLength ? text.slice(0, maxLength) + '...' : text;
};

const sortBy = ref<DataTableSortItem[]>([{key: 'status', order: 'asc'}]);

const statusScanRemarks = computed<AugmentedScanRemark[]>(() =>
  props.scanRemarks.map((s) => ({
    ...s,
    _statusIcon: statusIcon(s.status),
    _statusColor: statusColor(s.status),
    _statusText: statusText(s.status),
    _description: truncateText(s.description),
  })),
);

const baseHeaderClass = 'tableHeaderCell';
const scanHeaders = computed<DataTableHeader[]>(() => [
  {
    title: t('STATUS'),
    key: 'status',
    width: 130,
    align: 'center' as const,
    class: baseHeaderClass,
    sortable: true,
    mustSort: true,
    sort: compareScanRemarkStatus,
  },
  {title: t('COL_REMARK'), key: 'remark', width: 220, align: 'start' as const, class: baseHeaderClass},
  {title: t('COL_DESCRIPTION'), key: 'description', align: 'start' as const, class: baseHeaderClass},
]);
</script>

<template>
  <v-data-table
    height="350"
    :headers="scanHeaders"
    :items="statusScanRemarks"
    :sort-by="sortBy"
    item-key="description"
    density="compact"
    :items-per-page="25"
    :footer-props="{'items-per-page-options': [25, 50, 100, -1]}"
    class="striped-table custom-data-table elevation-0 p-2"
    hover>
    <template #item.status="{item}">
      <div class="flex h-full items-center justify-center">
        <v-icon :color="item._statusColor" dense>{{ item._statusIcon }}</v-icon>
        <Tooltip>{{ item._statusText }}</Tooltip>
      </div>
    </template>
    <template #item.remark="{item}">
      <span>{{ item.remark || '-' }}</span>
    </template>
    <template #item.description="{item}">
      <span>{{ item._description }}</span>
      <Tooltip>{{ item.description || '' }}</Tooltip>
    </template>
  </v-data-table>
</template>
