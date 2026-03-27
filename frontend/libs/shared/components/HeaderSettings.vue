// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {useHeaderSettingsStore} from '@shared/stores/headerSettings.store';
import {storeToRefs} from 'pinia';
import {computed} from 'vue';
import {useI18n} from 'vue-i18n';
import {InternalDataTableHeader} from 'vuetify/lib/components/VDataTable/types';

interface Props {
  column: InternalDataTableHeader;
  gridName: string;
  showBorders?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  showBorders: true,
});

const {t} = useI18n();
const headerSettingsStore = useHeaderSettingsStore();

headerSettingsStore.setupStore(props.gridName);
headerSettingsStore.setSettingsColumn(props.column);

const {selectableHeaders, selectedHeaders, initialSelectedHeaders} = storeToRefs(headerSettingsStore);

const selectedItems = computed((): string[] =>
  selectedHeaders.value
    .filter((headerIndex) => selectableHeaders.value[headerIndex] ?? false)
    .map((headerIndex) => selectableHeaders.value[headerIndex].value),
);

const updateSelectedHeaders = (newHeaders: string[]) => {
  headerSettingsStore.updateSelectedHeadersFromStringList(newHeaders);
};

const resetSelectedHeaders = () => {
  headerSettingsStore.resetSelectedHeaders();
};
</script>

<template>
  <GridHeaderMenu
    :show-reset="selectedHeaders.length !== initialSelectedHeaders.length"
    :reset-hint="t('Btn_clear_headers')"
    :card-title="t('TT_TABLE_SETTINGS')"
    :all-items="selectableHeaders"
    :selected-items="selectedItems"
    :select-label="t('LABEL_SELECT_COL')"
    @update="updateSelectedHeaders"
    @reset="resetSelectedHeaders">
    <template v-slot:activator="{props}">
      <div
        v-if="showBorders"
        class="absolute top-0 left-0 w-[40px] h-[40px] border-[rgb(var(--v-theme-primary))] border-t-[2px] border-l-[2px]"></div>
      <DIconButton
        :parentProps="props"
        icon="mdi-cog"
        :hint="t('TT_TABLE_SETTINGS')"
        color="primary"
        class="absolute top-0 left-0 w-[40px] h-[40px]" />
    </template>
  </GridHeaderMenu>
</template>
