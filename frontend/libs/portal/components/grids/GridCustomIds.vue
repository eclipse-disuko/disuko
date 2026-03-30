// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import CustomIdDialog from '@disclosure-portal/components/dialog/CustomIdDialog.vue';
import {CustomId} from '@disclosure-portal/model/CustomId';
import adminService from '@disclosure-portal/services/admin';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useCustomIdStore} from '@disclosure-portal/stores/customid.store';
import useSnackbar from '@shared/composables/useSnackbar';
import {SortItem} from '@shared/types/table';
import {computed, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const appStore = useAppStore();
const customIdStore = useCustomIdStore();
const {info: snack} = useSnackbar();

const dialog = ref<InstanceType<typeof CustomIdDialog>>();
const confirmVisible = ref(false);
const confirmConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const sortBy: SortItem[] = [{key: 'updated', order: 'desc'}];
const items = ref<CustomId[]>([]);
const loaded = ref(false);

const reload = async () => {
  loaded.value = false;
  adminService.getCustomIds().then((res) => {
    items.value = res.data;
    loaded.value = true;
  });
};

const headers = computed(() => [
  {
    title: t('COL_ACTIONS'),
    align: 'center',
    width: 120,
    value: 'actions',
  },
  {
    title: t('COL_TID'),
    width: 140,
    align: 'start',
    class: 'tableHeaderCell',
    value: '_key',
    sortable: true,
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
    title: t('LINK_TEMPLATE'),
    align: 'start',
    class: 'tableHeaderCell',
    value: 'linkTemplate',
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

const doDelete = async (config: IConfirmationDialogConfig) => {
  await adminService.deleteCustomId(config.key);
  snack(t('DIALOG_customid_delete'));
  await customIdStore.updateCustomIds();
  await reload();
};

const showConfirm = async (item: CustomId) => {
  const amount = (await adminService.customIdUsage(item._key!)).data.count;
  confirmConfig.value = {
    type: ConfirmationType.DELETE,
    key: item._key!,
    name: appStore.getAppLanguage == 'en' ? item.name : item.nameDE,
    description: 'DLG_CONFIRMATION_DESCRIPTION',
    okButton: 'Btn_delete',
    emphasiseText: t('CUSTOM_ID_USAGE', {count: amount}),
  };
  confirmVisible.value = true;
};

onMounted(async () => {
  await reload();
});
</script>

<template>
  <TableLayout>
    <template #buttons>
      <h1 class="d-headline">{{ t('CUSTOMIDS') }}</h1>
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
        :loading="!loaded"
        item-key="_key"
        :items="items"
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
          <TableActionButtons
            :buttons="[
              {icon: 'mdi-pencil', event: 'edit'},
              {icon: 'mdi-delete', event: 'delete'},
            ]"
            @edit="dialog?.open(item)"
            @delete="showConfirm(item)" />
        </template>
      </v-data-table>
    </template>
  </TableLayout>
  <CustomIdDialog ref="dialog" @reload="reload" />
  <ConfirmationDialog v-model:showDialog="confirmVisible" :config="confirmConfig" @confirm="doDelete" />
</template>
