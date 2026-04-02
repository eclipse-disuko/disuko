<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<template>
  <TableLayout has-tab has-title>
    <template #description>
      <div id="introTextGridOB">{{ t('LIC_LICENSE_CLASSIFICATION') }}</div>
    </template>
    <template #table>
      <v-data-table
        v-if="obligations?.length > 0"
        ref="gridOb"
        :headers="headers"
        fixed-header
        density="compact"
        class="striped-table custom-data-table fill-height"
        item-value="_key"
        :items-per-page="-1"
        :footer-props="{
          'items-per-page-options': [10, 50, 100, -1],
        }"
        :items="obligations"
        :search="search"
        :sort-by="sortItems">
        <template #item.created="{item}">
          <DDateCellWithTooltip :value="item.created" />
        </template>
        <template #item.type="{item}">
          {{ getTextOfType(item.type) }}
        </template>
        <template #item.warnLevel="{item}">
          <span>
            <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" location="bottom">
              <template #activator="{props}">
                <v-icon v-bind="props" :color="getIconColorOfLevel(item.warnLevel)">
                  {{ getIconOfLevel(item.warnLevel) }}
                </v-icon>
              </template>
              <span>{{ getTextOfLevel(item.warnLevel) }}</span>
            </v-tooltip>
          </span>
        </template>
        <template #item.name="{item}">
          {{ viewTools.getNameForLanguage(item) }}
        </template>
        <template #item.description="{item}">
          {{ viewTools.getDescriptionForLanguage(item) }}
        </template>
      </v-data-table>
    </template>
  </TableLayout>
</template>

<script setup lang="ts">
import {useView} from '@disclosure-portal/composables/useView';
import {IObligation} from '@disclosure-portal/model/IObligation';
import License from '@disclosure-portal/model/License';
import {compareLevel} from '@disclosure-portal/model/Quality';
import {useUserStore} from '@disclosure-portal/stores/user';
import useViewTools, {getIconColorOfLevel, getIconOfLevel} from '@disclosure-portal/utils/View';
import DDateCellWithTooltip from '@shared/components/disco/DDateCellWithTooltip.vue';
import TableLayout from '@shared/layouts/TableLayout.vue';
import {DataTableHeader, SortItem} from '@shared/types/table';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import {reactive, ref} from 'vue';
import {useI18n} from 'vue-i18n';
const props = defineProps<{
  license: License;
}>();
const {t} = useI18n();
const sortItems = ref<SortItem[]>([{key: 'warnLevel', order: 'desc'}]);
const rights = ref();
const search = ref('');
const viewTools = useViewTools();
const {getTextOfLevel, getTextOfType} = useView();
const obligations = reactive<IObligation[]>(props.license.meta.obligationsList);
const gridOb = ref<HTMLElement | null>(null);

const userStore = useUserStore();
rights.value = userStore.getRights;

const headers: DataTableHeader[] = [
  {
    title: t('COL_TYPE'),
    align: 'start',
    class: 'tableHeaderCell',
    value: 'type',
    sortable: true,
    width: 100,
  },
  {
    title: t('COL_WARN_LEVEL'),
    align: 'start',
    sortable: true,
    sort: compareLevel,
    class: 'tableHeaderCell',
    value: 'warnLevel',
    width: 80,
  },
  {
    title: t('COL_SHORT_NAME'),
    align: 'start',
    sortable: true,
    width: 180,
    class: 'tableHeaderCell',
    value: 'name',
  },
  {
    title: t('COL_DESCRIPTION'),
    align: 'start',
    class: 'tableHeaderCell',
    value: 'description',
    width: 400,
    sortable: true,
  },
  {
    title: t('COL_CREATED'),
    align: 'center',
    class: 'tableHeaderCell',
    value: 'created',
    width: 120,
  },
];
</script>
