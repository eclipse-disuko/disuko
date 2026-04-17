<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ChangeLogResponse, IPolicyRuleChangeLog} from '@disclosure-portal/model/ChangeLog';
import {useUserStore} from '@disclosure-portal/stores/user';
import {getCssClassForTableRow} from '@disclosure-portal/utils/Table';
import {
  getIconColorForPolicyType,
  getIconForChange,
  getIconForPolicyType,
  policyStateToTranslationKey,
} from '@disclosure-portal/utils/View';
import DDateCellWithTooltip from '@shared/components/disco/DDateCellWithTooltip.vue';
import TableLayout from '@shared/layouts/TableLayout.vue';
import {DataTableHeader, SortItem} from '@shared/types/table';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import _ from 'lodash';
import {computed, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const props = defineProps<{
  fetchMethod: () => Promise<ChangeLogResponse[]>;
}>();
const {t} = useI18n();
const items = ref<IPolicyRuleChangeLog[]>([]);
const search = ref('');
const userStore = useUserStore();
const dataAreLoaded = ref(false);
const menuFilterLicenseIds = ref(false);
const menuFilterPolicyStatus = ref(false);
const menuFilterChange = ref(false);
const selectedFilterLicenseIds = ref<string[]>([]);
const selectedFilterPolicyStatus = ref<string[]>([]);
const selectedFilterChange = ref<string[]>([]);
const headers = computed<DataTableHeader[]>(() => {
  const baseHeaders = [
    {
      title: t('CHANGE_LOG_COL_CHANGE'),
      align: 'center',
      width: 130,
      class: 'tableHeaderCell',
      value: 'content.change',
      sortable: true,
    },
    {
      title: t('POLICY_STATE'),
      align: 'center',
      width: 170,
      class: 'tableHeaderCell',
      value: 'content.policyStatus',
      sortable: true,
    },
    {
      title: t('CHANGE_LOG_COL_LICENSE_NAME'),
      align: 'start',
      class: 'tableHeaderCell',
      value: 'content.licenseName',
      sortable: true,
      filterable: true,
    },
    {
      title: t('COL_LICENSE_ID'),
      align: 'start',
      width: 300,
      class: 'tableHeaderCell',
      value: 'content.licenseId',
      sortable: true,
      filterable: true,
    },
    {
      title: t('CHANGE_LOG_COL_WHEN'),
      align: 'start',
      width: 140,
      class: 'tableHeaderCell',
      value: 'when',
      sortable: true,
    },
  ];
  const rights = userStore.getRights;
  if (rights.isInternal) {
    baseHeaders.push({
      title: t('COL_ACTIONS'),
      align: 'center',
      width: 160,
      class: 'tableHeaderCell',
      value: 'actions',
      sortable: false,
    });
  }
  return baseHeaders as DataTableHeader[];
});

const fetchMethod = async (): Promise<ChangeLogResponse[]> => {
  return await props.fetchMethod();
};

const reloadInternal = async (forceReload: boolean): Promise<void> => {
  if (!forceReload && dataAreLoaded.value) {
    return;
  }
  dataAreLoaded.value = false;
  const changeLogUnparsed = await fetchMethod();
  if (!changeLogUnparsed) {
    dataAreLoaded.value = true;
    return;
  }

  const unsortedItems = changeLogUnparsed.map((cl) => {
    cl.content = JSON.parse(cl.content);
    return cl;
  }) as unknown as IPolicyRuleChangeLog[];
  items.value = sortItems(unsortedItems);
  dataAreLoaded.value = true;
};

const sortItems = (unsortedItems: IPolicyRuleChangeLog[]): IPolicyRuleChangeLog[] => {
  return unsortedItems.sort((a, b) => {
    // Prioritize "Added" over "Removed"
    if (a.content.change === 'Added' && b.content.change === 'Removed') return -1;
    if (a.content.change === 'Removed' && b.content.change === 'Added') return 1;
    return 0;
  });
};

const filteredItems = computed(() => {
  return _.chain(items.value).filter(filterOnLicenseId).filter(filterOnPolicyStatus).filter(filterOnChange).value();
});

const filterOnLicenseId = (item: IPolicyRuleChangeLog): boolean => {
  return selectedFilterLicenseIds.value.length === 0 || selectedFilterLicenseIds.value.includes(item.content.licenseId);
};

const filterOnPolicyStatus = (item: IPolicyRuleChangeLog): boolean => {
  return (
    selectedFilterPolicyStatus.value.length === 0 ||
    selectedFilterPolicyStatus.value.includes(item.content.policyStatus)
  );
};

const filterOnChange = (item: IPolicyRuleChangeLog): boolean => {
  return selectedFilterChange.value.length === 0 || selectedFilterChange.value.includes(item.content.change);
};

const possibleFilterLicenseIds = computed(() => {
  return _.chain(items.value)
    .uniqBy('content.licenseId')
    .map((item: IPolicyRuleChangeLog) => item.content.licenseId)
    .map((val: string) => {
      return {
        text: val,
        value: val,
      };
    })
    .value();
});

const possibleFilterPolicyStatus = computed(() => {
  return _.chain(items.value)
    .uniqBy('content.policyStatus')
    .map((item: IPolicyRuleChangeLog) => item.content.policyStatus)
    .map((val: string) => {
      return {
        text: val,
        value: val,
      };
    })
    .value();
});

const possibleFilterChange = computed(() => {
  return _.chain(items.value)
    .uniqBy('content.change')
    .map((item: IPolicyRuleChangeLog) => item.content.change)
    .map((val: string) => {
      return {
        text: val,
        value: val,
      };
    })
    .value();
});

const sortBy = ref<SortItem[]>([]);
sortBy.value = [{key: 'when', order: 'desc'}];

onMounted(async () => {
  await reloadInternal(true);
});
</script>
<template>
  <TableLayout has-tab has-title>
    <template #buttons>
      <v-spacer></v-spacer>
      <v-text-field
        autocomplete="off"
        :max-width="500"
        density="compact"
        variant="outlined"
        v-model="search"
        append-inner-icon="mdi-magnify"
        :label="t('labelSearch')"
        clearable
        single-line
        hide-details></v-text-field>
    </template>
    <template #table>
      <div ref="gridChangeLog" class="fill-height">
        <v-data-table
          item-key="_key"
          :headers="headers"
          :items="filteredItems"
          :search="search"
          fixed-header
          density="compact"
          class="striped-table fill-height"
          :sort-by="sortBy"
          :item-class="getCssClassForTableRow"
          :items-per-page="10"
          :footer-props="{'items-per-page-options': [10, 50, 100, 200, 500, -1]}">
          <template v-slot:header.content.licenseId="{column, getSortIcon, toggleSort}">
            <div class="v-data-table-header__content">
              <span>{{ column.title }}</span>
              <v-menu :close-on-content-click="false" v-model="menuFilterLicenseIds">
                <template v-slot:activator="{props}">
                  <DIconButton
                    :parentProps="props"
                    icon="mdi-filter-variant"
                    :hint="t('TT_SHOW_FILTER')"
                    :color="selectedFilterLicenseIds.length > 0 ? 'primary' : 'default'"
                    location="top" />
                </template>
                <div class="bg-background" style="width: 280px">
                  <v-row class="d-flex ma-1 mr-2 justify-end">
                    <DIconButton icon="mdi-close" @clicked="menuFilterLicenseIds = false" color="default" />
                  </v-row>
                  <v-autocomplete
                    v-model="selectedFilterLicenseIds"
                    :items="possibleFilterLicenseIds"
                    class="pa-2 mx-2 pb-4"
                    :label="t('Lbl_filter_license')"
                    clearable
                    multiple
                    chips
                    item-title="text"
                    item-value="value"
                    variant="outlined"
                    density="compact">
                    <template v-slot:item="{props}">
                      <v-list-item v-bind="props">
                        <template v-slot:prepend="{isSelected}">
                          <v-checkbox hide-details :model-value="isSelected" />
                        </template>
                      </v-list-item>
                    </template>
                    <template v-slot:selection="{item, index}">
                      <div v-if="index === 0" class="d-flex align-center">
                        <span>{{ item.title }}</span>
                      </div>
                      <span v-if="index === 1" class="d-secondary-text d-subtitle-2 pl-1">
                        +{{ selectedFilterLicenseIds.length - 1 }} others
                      </span>
                    </template>
                  </v-autocomplete>
                </div>
              </v-menu>
              <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
            </div>
          </template>
          <template v-slot:header.content.policyStatus="{column, getSortIcon, toggleSort}">
            <div class="v-data-table-header__content">
              <span>{{ column.title }}</span>
              <v-menu :close-on-content-click="false" v-model="menuFilterPolicyStatus">
                <template v-slot:activator="{props}">
                  <DIconButton
                    :parentProps="props"
                    icon="mdi-filter-variant"
                    :hint="t('TT_SHOW_FILTER')"
                    :color="selectedFilterPolicyStatus.length > 0 ? 'primary' : 'default'"
                    location="top" />
                </template>
                <div class="bg-background" style="width: 280px">
                  <v-row class="d-flex ma-1 mr-2 justify-end">
                    <DIconButton icon="mdi-close" @clicked="menuFilterPolicyStatus = false" color="default" />
                  </v-row>
                  <v-select
                    v-model="selectedFilterPolicyStatus"
                    :items="possibleFilterPolicyStatus"
                    class="pa-2 mx-2 pb-4"
                    :label="t('Lbl_filter_status')"
                    clearable
                    multiple
                    item-title="text"
                    item-value="value"
                    variant="outlined"
                    density="compact">
                    <template v-slot:item="{item, props}">
                      <v-list-item v-bind="props" title="">
                        <template v-slot:prepend="{isSelected}">
                          <v-checkbox hide-details :model-value="isSelected" />
                        </template>
                        <v-icon small :color="getIconColorForPolicyType(item.value.toLowerCase())">{{
                          getIconForPolicyType(item.value.toLowerCase())
                        }}</v-icon>
                        <span class="d-subtitle-2 handpointer mx-2">{{
                          t(policyStateToTranslationKey(item.value.toLowerCase()))
                        }}</span>
                      </v-list-item>
                    </template>
                    <template v-slot:selection="{item, index}">
                      <div v-if="index === 0" class="d-flex align-center">
                        <v-icon small :color="getIconColorForPolicyType(item.value.toLowerCase())">{{
                          getIconForPolicyType(item.value.toLowerCase())
                        }}</v-icon>
                        <span class="d-subtitle-2 handpointer mx-2">{{
                          t(policyStateToTranslationKey(item.value.toLowerCase()))
                        }}</span>
                      </div>
                      <span v-if="index === 1" class="d-secondary-text d-subtitle-2 pl-1">
                        +{{ selectedFilterPolicyStatus.length - 1 }} others
                      </span>
                    </template>
                  </v-select>
                </div>
              </v-menu>
              <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
            </div>
          </template>
          <template v-slot:header.content.change="{column, getSortIcon, toggleSort}">
            <div class="v-data-table-header__content">
              <span>{{ column.title }}</span>
              <v-menu :close-on-content-click="false" v-model="menuFilterChange">
                <template v-slot:activator="{props}">
                  <DIconButton
                    :parentProps="props"
                    icon="mdi-filter-variant"
                    :hint="t('TT_SHOW_FILTER')"
                    :color="selectedFilterChange.length > 0 ? 'primary' : 'default'"
                    location="top" />
                </template>
                <div class="bg-background" style="width: 280px">
                  <v-row class="d-flex ma-1 mr-2 justify-end">
                    <DIconButton icon="mdi-close" @clicked="menuFilterChange = false" color="default" />
                  </v-row>
                  <v-select
                    v-model="selectedFilterChange"
                    :items="possibleFilterChange"
                    class="pa-2 mx-2 pb-4"
                    :label="t('lbl_filter_on_change')"
                    clearable
                    multiple
                    item-title="text"
                    item-value="value"
                    variant="outlined"
                    density="compact">
                    <template v-slot:item="{item, props}">
                      <v-list-item v-bind="props" title="">
                        <template v-slot:prepend="{isSelected}">
                          <v-checkbox hide-details :model-value="isSelected" />
                        </template>
                        <v-icon
                          size="small"
                          :color="getIconColorForPolicyType(item.value.toLowerCase())"
                          class="mr-2"
                          :icon="getIconForChange(item.value)">
                        </v-icon>

                        <span class="mr-2">{{ item.props.title }}</span>
                      </v-list-item>
                    </template>
                    <template v-slot:selection="{item, index}">
                      <div v-if="index === 0" class="d-flex align-center">
                        <v-icon
                          size="small"
                          :color="getIconColorForPolicyType(item.value.toLowerCase())"
                          class="mr-2"
                          :icon="getIconForChange(item.value)">
                        </v-icon>

                        <span class="mr-2">{{ item.props.title }}</span>
                      </div>
                      <span v-if="index === 1" class="d-secondary-text d-subtitle-2 pl-1">
                        +{{ selectedFilterChange.length - 1 }} others
                      </span>
                    </template>
                  </v-select>
                </div>
              </v-menu>
              <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
            </div>
          </template>
          <template v-slot:item.when="{item}">
            <DDateCellWithTooltip :value="item.when" />
          </template>
          <template v-slot:item.content.policyStatus="{item}">
            <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" bottom content-class="dpTooltip">
              <template v-slot:activator="{props}">
                <v-icon
                  small
                  :color="getIconColorForPolicyType(item.content.policyStatus.toLowerCase())"
                  v-bind="props">
                  {{ getIconForPolicyType(item.content.policyStatus.toLowerCase()) }}
                </v-icon>
              </template>
              {{ t('POLICY_STATUS_' + item.content.policyStatus.toUpperCase()) }}
            </v-tooltip>
          </template>
          <template v-slot:item.content.change="{item}">
            <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" bottom content-class="dpTooltip">
              <template v-slot:activator="{props}">
                <v-icon small v-bind="props" color="textColor">
                  {{ getIconForChange(item.content.change) }}
                </v-icon>
              </template>
              {{ t('CHANGE_' + item.content.change.toUpperCase()) }}
            </v-tooltip>
          </template>
          <template v-slot:item.actions="{item}">
            <router-link
              v-if="item.content.licenseName !== ''"
              :to="'/dashboard/licenses/' + item.content.licenseId"
              target="_blank">
              <v-icon color="primary" size="large">mdi mdi-chevron-right</v-icon>
            </router-link>
          </template>
        </v-data-table>
      </div>
    </template>
  </TableLayout>
</template>
