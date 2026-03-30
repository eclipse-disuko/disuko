// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import {ChecklistItem} from '@disclosure-portal/model/Checklist';
import {useChecklistsStore} from '@disclosure-portal/stores/checklists.store';
import {useLabelStore} from '@disclosure-portal/stores/label.store';
import useSnackbar from '@shared/composables/useSnackbar';
import TableLayout from '@shared/layouts/TableLayout.vue';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {DataTableHeader, SortItem} from '@shared/types/table';
import {storeToRefs} from 'pinia';
import {computed, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const {dashboardCrumbs, ...breadcrumbs} = useBreadcrumbsStore();
const checklistsStore = useChecklistsStore();
const {info: snack} = useSnackbar();
const labelStore = useLabelStore();

const {checklist} = storeToRefs(checklistsStore);

const confirmVisible = ref(false);
const confirmConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const sortBy: SortItem[] = [{key: 'updated', order: 'desc'}];
const dialogItem = ref();

const policyLabels = computed(() =>
  checklist.value ? checklist.value?.policyLabels.map((labelKey) => labelStore.getLabelByKey(labelKey)) : [],
);

const initBreadcrumbs = () => {
  console.log(checklist.value);
  breadcrumbs.setCurrentBreadcrumbs([
    ...dashboardCrumbs,
    {
      title: t('BC_CHECKLIST'),
      href: '/dashboard/admin/checklist',
    },
    {
      title: checklist.value!.name,
      disabled: false,
      href: `/dashboard/admin/checklist/${encodeURIComponent(checklist.value!._key)}`,
    },
  ]);
};

const headers = computed<DataTableHeader[]>(() => [
  {
    title: t('CD_NAME'),
    align: 'start',
    class: 'tableHeaderCell',
    value: 'name',
    width: 140,
    sortable: true,
  },
  {
    title: t('LBL_CHECKLIST_TEMPLATE'),
    align: 'start',
    class: 'tableHeaderCell',
    value: 'targetTemplateName',
    width: 140,
    sortable: true,
  },
  {
    title: t('COL_CREATED'),
    key: 'created',
    align: 'start',
    width: 120,
  },
  {
    title: t('COL_UPDATED'),
    key: 'updated',
    align: 'start',
    width: 120,
  },
  {
    title: t('COL_ACTIONS'),
    align: 'center',
    width: 120,
    value: 'actions',
  },
]);

const doDelete = async (config: IConfirmationDialogConfig) => {
  await checklistsStore.deleteItem(config.key);
  snack(t('DIALOG_CHECKLIST_DELETE'));
};

const showConfirm = async (item: ChecklistItem) => {
  confirmConfig.value = {
    type: ConfirmationType.DELETE,
    key: item._key!,
    name: item.name,
    description: 'DLG_CONFIRMATION_DESCRIPTION',
    okButton: 'Btn_delete',
  };
  confirmVisible.value = true;
};

watch(
  checklist,
  (value) => {
    if (value) {
      initBreadcrumbs();
    }
  },
  {immediate: true},
);
</script>

<template>
  <TableLayout>
    <template v-if="$slots.default" #description>
      <slot></slot>
    </template>

    <template #buttons>
      <Stack>
        <div>
          <span class="text-h5 pr-2">{{ t('CHECKLIST') }}</span>
          <q class="text-h5">{{ checklist?.name }}</q>
        </div>
        <div>
          <ProjectLabel v-for="(label, i) in policyLabels" :label="label" :key="i"></ProjectLabel>
        </div>
        <div>
          <DCActionButton
            large
            icon="mdi-plus"
            :hint="t('TT_add_project')"
            :text="t('BTN_ADD')"
            @click="dialogItem?.open()" />
        </div>
      </Stack>
    </template>

    <template #table>
      <v-data-table
        density="compact"
        class="striped-table fill-height"
        item-key="_key"
        :items="checklist?.items"
        :headers="headers"
        :items-per-page="50"
        fixed-header
        :sort-by="sortBy"
        sort-desc>
        <template v-slot:[`item.created`]="{item}">
          <DDateCellWithTooltip :value="item.created" />
        </template>
        <template v-slot:[`item.updated`]="{item}">
          <DDateCellWithTooltip :value="item.updated" />
        </template>
        <template v-slot:[`item.actions`]="{item}">
          <DIconButton icon="mdi-pencil" @clicked="dialogItem?.open(item)" />
          <DIconButton icon="mdi-delete" @clicked="showConfirm(item)" />
        </template>
      </v-data-table>
    </template>
  </TableLayout>

  <ChecklistItemDialog ref="dialogItem"></ChecklistItemDialog>
  <ConfirmationDialog
    v-model:showDialog="confirmVisible"
    :config="confirmConfig"
    @confirm="doDelete"></ConfirmationDialog>
</template>
