<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import {Checklist} from '@disclosure-portal/model/Checklist';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useChecklistsStore} from '@disclosure-portal/stores/checklists.store';
import TableActionButtons, {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import TableLayout from '@shared/layouts/TableLayout.vue';
import {DataTableHeader, DataTableItem, SortItem} from '@shared/types/table';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';

const {t} = useI18n();
const router = useRouter();
const appStore = useAppStore();
const checklistsStore = useChecklistsStore();
const {info: snack} = useSnackbar();

const dialog = ref();
const confirmVisible = ref(false);
const confirmConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const sortBy: SortItem[] = [{key: 'updated', order: 'desc'}];

const headers = computed<DataTableHeader[]>(() => [
  {
    title: t('COL_ACTIONS'),
    align: 'center',
    width: 80,
    value: 'actions',
  },
  {
    title: t('AL_DIALOG_TF_NAME_EN'),
    align: 'start',
    class: 'tableHeaderCell',
    value: 'name',
    width: 140,
    sortable: true,
  },
  {
    title: t('AL_DIALOG_TF_NAME_DE'),
    align: 'start',
    class: 'tableHeaderCell',
    value: 'nameDE',
    width: 150,
    sortable: true,
  },
  {
    title: t('ACTIVE'),
    align: 'start',
    class: 'tableHeaderCell',
    value: 'active',
    width: 80,
    sortable: true,
  },
  {
    title: t('AL_DIALOG_TF_DESCRIPTION_EN'),
    align: 'start',
    class: 'tableHeaderCell',
    value: 'description',
    width: 180,
    sortable: true,
  },
  {
    title: t('AL_DIALOG_TF_DESCRIPTION_DE'),
    align: 'start',
    class: 'tableHeaderCell',
    value: 'descriptionDE',
    width: 180,
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
]);

const openItem = (_event: Event, item: DataTableItem<Checklist>) => {
  const url = `/dashboard/admin/checklist/${encodeURIComponent(item.item._key)}`;
  router.push(url);
};

const doDelete = async (config: IConfirmationDialogConfig) => {
  await checklistsStore.deleteChecklist(config.key);
  snack(t('DIALOG_CHECKLIST_DELETE'));
};

const showConfirm = async (item: Checklist) => {
  confirmConfig.value = {
    type: ConfirmationType.DELETE,
    key: item._key!,
    name: appStore.getAppLanguage == 'en' ? item.name : item.nameDE,
    description: 'DLG_CONFIRMATION_DESCRIPTION',
    okButton: 'Btn_delete',
  };
  confirmVisible.value = true;
};

const getActionButtons = (item: Checklist): TableActionButtonsProps['buttons'] => {
  return [
    {
      icon: 'mdi-pencil',
      event: 'edit',
    },
    {
      icon: 'mdi-delete',
      event: 'delete',
      disabled: item.active,
      hint: item.active ? t('ACTIVE_CHECKLIST_DELETE') : undefined,
    },
  ];
};
</script>

<template>
  <TableLayout>
    <template #buttons>
      <h1 class="d-headline">{{ t('CHECKLISTS') }}</h1>
      <DCActionButton
        large
        icon="mdi-plus"
        :hint="t('TT_add_project')"
        :text="t('BTN_ADD')"
        @click="dialog?.open()"
        class="mx-2" />
    </template>

    <template #table>
      <v-data-table
        density="compact"
        class="striped-table fill-height"
        :loading="checklistsStore.isLoading"
        item-key="_key"
        :items="checklistsStore.checklists"
        :headers="headers"
        :items-per-page="50"
        fixed-header
        :sort-by="sortBy"
        @click:row="openItem"
        sort-desc>
        <template v-slot:[`item.active`]="{item}">
          <v-icon icon="mdi-check" class="mr-2" :color="item.active ? 'primary' : 'tableBorderColor'"></v-icon>
        </template>
        <template v-slot:[`item.description`]="{item}">
          <Truncated>{{ item.description }}</Truncated>
        </template>
        <template v-slot:[`item.descriptionDE`]="{item}">
          <Truncated>{{ item.descriptionDE }}</Truncated>
        </template>
        <template v-slot:[`item.created`]="{item}">
          <DDateCellWithTooltip :value="item.created" />
        </template>
        <template v-slot:[`item.updated`]="{item}">
          <DDateCellWithTooltip :value="item.updated" />
        </template>
        <template v-slot:[`item.actions`]="{item}">
          <TableActionButtons
            :buttons="getActionButtons(item)"
            variant="normal"
            @edit="dialog?.open(item)"
            @delete="showConfirm(item)" />
        </template>
      </v-data-table>
    </template>
  </TableLayout>

  <ChecklistDialog ref="dialog"></ChecklistDialog>
  <ConfirmationDialog
    v-model:showDialog="confirmVisible"
    :config="confirmConfig"
    @confirm="doDelete"></ConfirmationDialog>
</template>
