<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {DEFAULT_CLASSIFICATION_NAMES} from '@disclosure-portal/model/IObligation';
import {policyRulesMatrixStore} from '@disclosure-portal/stores/classificationMatrix.store';
import Tooltip from '@shared/components/disco/Tooltip.vue';
import {useHeaderSettings} from '@shared/composables/useHeaderSettings';
import useSnackbar from '@shared/composables/useSnackbar';
import DialogLayout, {type DialogLayoutConfig} from '@shared/layouts/DialogLayout.vue';
import {DataTableHeader} from '@shared/types/table';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {storeToRefs} from 'pinia';

const {t} = useI18n();
const {error} = useSnackbar();
const store = policyRulesMatrixStore();
const {classifications, isLoading, policyRulesWithStatus} = storeToRefs(store);

const isVisible = ref(false);
const search = ref('');

const tableName = 'ClassificationMatrixReadonlyDialog';

const headers = computed<DataTableHeader[]>(() => [
  {title: t('POLICY_RULES'), align: 'start', width: 150, value: 'name', sortable: true},
  ...classifications.value.map((c) => ({
    key: c._key,
    value: c._key,
    title: c.name,
    align: 'center' as const,
    width: 150,
    sortable: false,
  })),
]);

const initiallyHiddenList = computed(() =>
  classifications.value.filter((c) => !DEFAULT_CLASSIFICATION_NAMES.has(c.name)).map((c) => c._key),
);
const headerSettings = useHeaderSettings({tableName});
const {filteredHeaders} = headerSettings;

const syncHeaderSettings = () => {
  headerSettings.resetHeaderSettings({
    tableName,
    headers: headers.value,
    initiallyHiddenList: initiallyHiddenList.value,
  });
};

const dialogConfig = computed(
  (): DialogLayoutConfig => ({
    title: t('CLASSIFICATION_MATRIX'),
  }),
);

const open = async () => {
  isVisible.value = true;
  try {
    await store.loadMatrix();
    syncHeaderSettings();
  } catch {
    error(t('MATRIX_LOAD_ERROR'));
  }
};

defineExpose({open});
</script>

<template>
  <v-dialog v-model="isVisible" scrollable width="95vw" max-width="95vw">
    <DialogLayout :config="dialogConfig" @close="isVisible = false">
      <div class="mb-2 flex justify-end">
        <DSearchField v-model="search" />
      </div>
      <v-data-table
        density="compact"
        class="striped-table"
        :loading="isLoading"
        :headers="filteredHeaders"
        :items="policyRulesWithStatus"
        v-model:search="search"
        item-value="_key"
        fixed-header
        height="65vh"
        :items-per-page="-1"
        :sort-by="[{key: 'name', order: 'asc'}]">
        <template #[`header.name`]="{column}">
          <GridFilterHeader :column="column">
            <template #settings>
              <HeaderSettings :grid-name="tableName" :column="column" />
            </template>
          </GridFilterHeader>
        </template>
        <template v-for="c in classifications" :key="`header-${c._key}`" #[`header.${c._key}`]="{column}">
          <Tooltip :text="c.name" as-parent>
            <span class="inline-block max-w-[150px] truncate">{{ column.title }}</span>
          </Tooltip>
        </template>
        <template v-for="c in classifications" :key="c._key" #[`item.${c._key}`]="{item}">
          <Tooltip v-if="item.statusProps?.[c._key]" :text="item.statusProps?.[c._key]?.label" as-parent>
            <v-icon :icon="item.statusProps?.[c._key]?.icon" :color="item.statusProps?.[c._key]?.color" />
          </Tooltip>
        </template>
      </v-data-table>
    </DialogLayout>
  </v-dialog>
</template>
