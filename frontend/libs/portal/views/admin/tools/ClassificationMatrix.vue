<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import PolicyRuleClassificationDialog from '@disclosure-portal/components/dialog/PolicyRuleClassificationDialog.vue';
import {
  ClassificationInfo,
  DEFAULT_CLASSIFICATION_NAMES,
  PolicyRuleClassificationDto,
  RuleStatus,
} from '@disclosure-portal/model/PolicyRuleClassification';
import policyRuleClassificationService from '@disclosure-portal/services/policyruleclassification';
import {getStatusColor, getStatusIcon, statusConfig} from '@disclosure-portal/utils/classificationStatus';
import Tooltip from '@shared/components/disco/Tooltip.vue';
import {useHeaderSettings} from '@shared/composables/useHeaderSettings';
import useSnackbar from '@shared/composables/useSnackbar';
import {DataTableHeader} from '@shared/types/table';
import {computed, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const {info: snack, error} = useSnackbar();

const dialog = ref<InstanceType<typeof PolicyRuleClassificationDialog>>();
const confirmVisible = ref(false);
const confirmConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const isLoading = ref(false);
const search = ref('');
const classifications = ref<ClassificationInfo[]>([]);
const useCases = ref<PolicyRuleClassificationDto[]>([]);

const headers = computed<DataTableHeader[]>(() => [
  {title: 'COL_ACTIONS', align: 'center', width: 140, value: 'actions', sortable: false},
  {title: 'COL_USE_CASE', align: 'start', width: 150, value: 'name', sortable: true},
  ...classifications.value.map((c) => ({
    key: c.key,
    value: c.key,
    title: c.name,
    align: 'center' as const,
    width: 140,
    sortable: false,
  })),
]);

const tableName = 'ClassificationMatrix';
const initiallyHiddenList = computed(() =>
  classifications.value.filter((c) => !DEFAULT_CLASSIFICATION_NAMES.includes(c.name)).map((c) => c.key),
);
const headerSettings = useHeaderSettings({tableName});
const {filteredHeaders} = headerSettings;

const getStatusLabel = (status: RuleStatus | undefined) => (status ? t(statusConfig[status]?.labelKey) : '—');

const statusPropsMap = computed(() => {
  const map = new Map<string, {icon: string; color: string; label: string}>();
  for (const useCase of useCases.value) {
    for (const [key, status] of Object.entries(useCase.rules ?? {})) {
      const mapKey = `${useCase._key}-${key}`;
      map.set(mapKey, {
        icon: getStatusIcon(status),
        color: getStatusColor(status),
        label: getStatusLabel(status),
      });
    }
  }
  return map;
});

const getStatusPropsComputed = (itemKey: string, classificationKey: string) =>
  statusPropsMap.value.get(`${itemKey}-${classificationKey}`);

const loadMatrix = async () => {
  isLoading.value = true;
  try {
    const res = await policyRuleClassificationService.getMatrix();
    classifications.value = res.data.classifications ?? [];
    useCases.value = res.data.useCases ?? [];
    headerSettings.resetHeaderSettings({
      tableName,
      headers: headers.value,
      initiallyHiddenList: initiallyHiddenList.value,
    });
  } catch {
    error(t('MATRIX_LOAD_ERROR'));
  } finally {
    isLoading.value = false;
  }
};

const headerExpands = () => {
  headerSettings.resetHeaderSettings({tableName, headers: headers.value});
};

const showConfirm = (item: PolicyRuleClassificationDto) => {
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
  await policyRuleClassificationService.delete(config.key);
  snack(t('DIALOG_prc_delete_success'));
  await loadMatrix();
};

onMounted(loadMatrix);
</script>

<template>
  <TableLayout has-tab has-title>
    <template #buttons>
      <h1 class="text-h5">{{ t('CLASSIFICATION_MATRIX') }}</h1>
      <DCActionButton
        large
        icon="mdi-plus"
        :text="t('BTN_ADD')"
        :hint="t('TT_ADD_CLASSIFICATION_MATRIX')"
        @click="dialog?.open()" />
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
          :items="useCases"
          v-model:search="search"
          item-value="_key"
          fixed-header
          :items-per-page="-1"
          :footer-props="{'items-per-page-options': [10, 50, 100, -1]}"
          :sort-by="[{key: 'name', order: 'asc'}]">
          <template #[`header.actions`]="{column}">
            <GridFilterHeader :column="column">
              <template #settings>
                <HeaderSettings :grid-name="tableName" :column="column" />
              </template>
            </GridFilterHeader>
          </template>
          <template #[`item.actions`]="{item}">
            <TableActionButtons
              variant="normal"
              :buttons="[
                {icon: 'mdi-pencil', event: 'edit'},
                {icon: 'mdi-delete', event: 'delete'},
              ]"
              @edit="dialog?.open(item)"
              @delete="showConfirm(item)"
              @slideToggle="headerExpands" />
          </template>
          <template v-for="c in classifications" :key="c.key" #[`item.${c.key}`]="{item}">
            <Tooltip v-if="item.rules?.[c.key]" :text="getStatusPropsComputed(item._key, c.key)?.label" as-parent>
              <v-icon
                :icon="getStatusPropsComputed(item._key, c.key)?.icon"
                :color="getStatusPropsComputed(item._key, c.key)?.color" />
            </Tooltip>
            <span v-else class="text-medium-emphasis">—</span>
          </template>
        </v-data-table>
      </div>
    </template>
  </TableLayout>

  <PolicyRuleClassificationDialog ref="dialog" @reload="loadMatrix" />
  <ConfirmationDialog v-model:showDialog="confirmVisible" :config="confirmConfig" @confirm="doDelete" />
</template>
