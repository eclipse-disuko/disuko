<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import PolicyRuleClassificationDialog from '@disclosure-portal/components/dialog/PolicyRuleClassificationDialog.vue';
import {DEFAULT_CLASSIFICATION_NAMES} from '@disclosure-portal/model/IObligation';
import PolicyRule from '@disclosure-portal/model/PolicyRule';
import adminService from '@disclosure-portal/services/admin';
import {policyRulesMatrixStore} from '@disclosure-portal/stores/classificationMatrix.store';
import {RightsUtils} from '@disclosure-portal/utils/Rights';
import Tooltip from '@shared/components/disco/Tooltip.vue';
import {useHeaderSettings} from '@shared/composables/useHeaderSettings';
import useSnackbar from '@shared/composables/useSnackbar';
import {DataTableHeader} from '@shared/types/table';
import {computed, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {storeToRefs} from 'pinia';

const {t} = useI18n();
const {info: snack, error} = useSnackbar();
const store = policyRulesMatrixStore();
const {classifications, isLoading, policyRulesWithStatus} = storeToRefs(store);

const dialog = ref<InstanceType<typeof PolicyRuleClassificationDialog>>();
const confirmVisible = ref(false);
const confirmConfig = ref<IConfirmationDialogConfig>();
const search = ref('');
const isPolicyManager = ref(false);
const showActions = computed(() => isPolicyManager.value);
const actionHeader: DataTableHeader = {
  title: 'COL_ACTIONS',
  align: 'center',
  width: 100,
  value: 'actions',
  sortable: false,
};

const headers = computed<DataTableHeader[]>(() => [
  ...(showActions.value ? [actionHeader] : []),
  {title: 'POLICY_RULES', align: 'start', width: 150, value: 'name', sortable: true},
  ...classifications.value.map((c) => ({
    key: c._key,
    value: c._key,
    title: c.name,
    align: 'center' as const,
    width: 150,
    sortable: false,
  })),
]);

const tableName = 'ClassificationMatrix';
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

const loadMatrix = async () => {
  try {
    await store.loadMatrix();
    syncHeaderSettings();
  } catch {
    error(t('MATRIX_LOAD_ERROR'));
  }
};

const actionButtons = [
  {icon: 'mdi-pencil', event: 'edit'},
  {icon: 'mdi-delete', event: 'delete'},
];

const showConfirm = (item: PolicyRule) => {
  confirmConfig.value = {
    type: ConfirmationType.DELETE,
    key: item._key,
    name: item.name,
    description: 'DLG_CONFIRMATION_DESCRIPTION',
    okButton: 'Btn_delete',
  };
  confirmVisible.value = true;
};

const doDelete = async (config: IConfirmationDialogConfig) => {
  await adminService.deletePolicyRule(config.key);
  snack(t('DIALOG_prc_delete_success'));
  await loadMatrix();
};

onMounted(() => {
  isPolicyManager.value = RightsUtils.isPolicyManager();
  loadMatrix();
});
</script>

<template>
  <TableLayout has-tab has-title>
    <template #buttons>
      <h1 class="text-h5">{{ t('CLASSIFICATION_MATRIX') }}</h1>
      <v-spacer />
      <DSearchField v-model="search" />
    </template>
    <template #table>
      <div class="fill-height">
        <v-data-table
          density="compact"
          class="striped-table fill-height"
          :loading="isLoading"
          :headers="filteredHeaders"
          :items="policyRulesWithStatus"
          v-model:search="search"
          item-value="_key"
          fixed-header
          :items-per-page="-1"
          :footer-props="{'items-per-page-options': [10, 50, 100, -1]}"
          :sort-by="[{key: 'name', order: 'asc'}]">
          <template v-if="showActions" #[`header.actions`]="{column}">
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
          <template v-if="showActions" #[`item.actions`]="{item}">
            <TableActionButtons
              variant="normal"
              :buttons="actionButtons"
              @edit="dialog?.open(item)"
              @delete="showConfirm(item)"
              @slideToggle="syncHeaderSettings" />
          </template>
          <template v-for="c in classifications" :key="c._key" #[`item.${c._key}`]="{item}">
            <Tooltip v-if="item.statusProps?.[c._key]" :text="item.statusProps?.[c._key]?.label" as-parent>
              <v-icon :icon="item.statusProps?.[c._key]?.icon" :color="item.statusProps?.[c._key]?.color" />
            </Tooltip>
          </template>
        </v-data-table>
      </div>
    </template>
  </TableLayout>

  <PolicyRuleClassificationDialog ref="dialog" @reload="loadMatrix" />
  <ConfirmationDialog v-model:showDialog="confirmVisible" :config="confirmConfig" @confirm="doDelete" />
</template>
