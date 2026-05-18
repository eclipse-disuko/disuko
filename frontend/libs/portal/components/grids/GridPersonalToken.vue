<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import PersonalTokenDialog from '@disclosure-portal/components/dialog/PersonalTokenDialog.vue';
import {PersonalToken} from '@disclosure-portal/model/PersonalToken';
import profileService from '@disclosure-portal/services/profile';
import useSnackbar from '@shared/composables/useSnackbar';
import {SortItem} from '@shared/types/table';
import {computed, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const isLoading = ref(false);
const {t} = useI18n();

const dialog = ref<InstanceType<typeof PersonalTokenDialog>>();
const confirmVisible = ref(false);
const confirmConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const sortBy: SortItem[] = [{key: 'created', order: 'desc'}];
const items = ref<PersonalToken[]>([]);
const {info: snack} = useSnackbar();

const reload = async () => {
  isLoading.value = true;
  try {
    const res = await profileService.getTokens();
    items.value = res || [];
  } catch (error) {
    items.value = [];
    console.error('Failed to load personal tokens:', error);
  } finally {
    isLoading.value = false;
  }
};

const getStatus = (item: PersonalToken) => {
  if (item.expired) {
    return {text: t('STATUS_EXPIRED'), color: 'error', icon: 'mdi-clock-alert'};
  }

  return {text: t('STATUS_ACTIVE'), color: 'success', icon: 'mdi-check-circle'};
};

const headers = computed(() => {
  return [
    {
      title: t('COL_ACTIONS'),
      align: 'center' as const,
      width: 80,
      value: 'actions',
    },
    {
      title: t('COL_DESCRIPTION'),
      align: 'start' as const,
      class: 'tableHeaderCell',
      value: 'description',
      sortable: true,
    },
    {
      title: t('COL_STATUS'),
      key: 'status',
      align: 'start' as const,
      width: 160,
    },
    {
      title: t('COL_EXPIRY'),
      key: 'expiry',
      align: 'start' as const,
      width: 160,
    },
    {
      title: t('COL_CREATED'),
      key: 'created',
      align: 'start' as const,
      width: 160,
    },
  ];
});

onMounted(async () => {
  await reload();
});

const handleAddToken = () => {
  dialog.value?.open();
};

const showExpireConfirm = (item: PersonalToken) => {
  confirmConfig.value = {
    type: ConfirmationType.DELETE,
    key: item.key,
    name: item.description,
    description: t('CONFIRM_EXPIRE_PERSONAL_TOKEN_DESCRIPTION'),
    okButton: t('BTN_EXPIRE'),
  };
  confirmVisible.value = true;
};

const doExpire = async (config: IConfirmationDialogConfig) => {
  await profileService.expireToken(config.key);
  snack(t('PERSONAL_TOKEN_EXPIRED'));
  await reload();
};
</script>

<template>
  <TableLayout has-tab has-title>
    <template #buttons>
      <v-spacer></v-spacer>
      <DCActionButton
        large
        icon="mdi-plus"
        :hint="t('PERSONAL_TOKEN_ADD_TITLE')"
        :text="t('BTN_ADD')"
        @click="handleAddToken"
        class="mx-2" />
    </template>
    <template #table>
      <v-data-table
        density="compact"
        class="striped-table fill-height"
        :loading="isLoading"
        item-key="key"
        :items="items"
        :headers="headers"
        :items-per-page="50"
        fixed-header
        :sort-by="sortBy"
        sort-desc>
        <template v-slot:[`item.created`]="{item}">
          <DDateCellWithTooltip :value="item.created" />
        </template>
        <template v-slot:[`item.expiry`]="{item}">
          <DDateCellWithTooltip :value="item.expiry" />
        </template>
        <template v-slot:[`item.status`]="{item}">
          <v-chip :color="getStatus(item).color" size="small" :prepend-icon="getStatus(item).icon">
            {{ getStatus(item)?.text }}
          </v-chip>
        </template>
        <template v-slot:[`item.actions`]="{item}">
          <TableActionButtons
            :buttons="[
              {
                icon: 'mdi-cancel',
                event: 'expire',
                show: !item.expired,
              },
            ]"
            @expire="showExpireConfirm(item)" />
        </template>
      </v-data-table>
    </template>
  </TableLayout>
  <PersonalTokenDialog ref="dialog" @reload="reload" />
  <ConfirmationDialog v-model:showDialog="confirmVisible" :config="confirmConfig" @confirm="doExpire" />
</template>
