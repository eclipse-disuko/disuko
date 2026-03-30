// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {DataTableHeaderFilterItems} from '@shared/types/table';
import {useI18n} from 'vue-i18n';
import {InternalDataTableHeader} from 'vuetify/lib/components/VDataTable/types';

interface Props {
  column: InternalDataTableHeader;
  label: string;
  allItems: DataTableHeaderFilterItems[];
  initialSelected?: string[];
}

const props = withDefaults(defineProps<Props>(), {
  initialSelected: () => [],
});

const selectedFilters = defineModel<string[]>({
  default: [],
});

const {t} = useI18n();

if (props.initialSelected.length >= 1) {
  selectedFilters.value = props.initialSelected;
}

const resetSelectedFilters = () => {
  selectedFilters.value = props.initialSelected;
};
</script>

<template>
  <GridHeaderMenu
    :show-reset="selectedFilters.length >= 1"
    :reset-hint="t('Btn_clear_filters')"
    :card-title="t('FILTER_ON')"
    :all-items="allItems"
    :selected-items="selectedFilters"
    :select-label="label"
    @update="(value: string[]) => (selectedFilters = value)"
    @reset="resetSelectedFilters">
    <template v-slot:activator="{props}">
      <span>
        <v-icon class="mr-1" v-bind="props" :color="selectedFilters.length > 0 ? 'primary' : 'default'">
          mdi-filter-variant
        </v-icon>
        <Tooltip>{{ t('TT_SHOW_FILTER') }}</Tooltip>
      </span>
    </template>
  </GridHeaderMenu>
</template>
