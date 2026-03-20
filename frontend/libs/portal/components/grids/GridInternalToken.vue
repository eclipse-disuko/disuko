<script setup lang="ts">
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import InternalTokenDialog from '@disclosure-portal/components/dialog/InternalTokenDialog.vue';
import {InternalToken} from '@disclosure-portal/model/InternalToken';
import adminService from '@disclosure-portal/services/admin';
import useSnackbar from '@shared/composables/useSnackbar';
import TableLayout from '@shared/layouts/TableLayout.vue';
import {SortItem} from '@shared/types/table';
import {computed, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';
const isLoading = ref(false);

const {t} = useI18n();

const dialog = ref<InstanceType<typeof InternalTokenDialog>>();
const confirmVisible = ref(false);
const confirmConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const sortBy: SortItem[] = [{key: 'updated', order: 'desc'}];
const items = ref<InternalToken[]>([]);
const {info: snack} = useSnackbar();

const reload = async () => {
  isLoading.value = true;
  try {
    const res = await adminService.getInternalTokens();
    items.value = res.data || [];
  } catch (error) {
    items.value = [];
    console.error('Failed to load internal tokens:', error);
  } finally {
    isLoading.value = false;
  }
};

// Function to determine status based on revoked and expiry
const getStatus = (item: InternalToken) => {
  if (item.revoked) {
    return {text: t('STATUS_REVOKED'), color: 'error', icon: 'mdi-cancel'};
  }

  if (item.expiry) {
    const expiryDate = new Date(item.expiry);
    const now = new Date();

    if (expiryDate <= now) {
      return {text: t('STATUS_EXPIRED'), color: 'error', icon: 'mdi-clock-alert'};
    }
  }

  return {text: t('STATUS_ACTIVE'), color: 'success', icon: 'mdi-check-circle'};
};

const headers = computed(() => {
  return [
    {
      title: t('COL_NAME'),
      align: 'start' as const,
      class: 'tableHeaderCell',
      value: 'name',
      width: 180,
      sortable: true,
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

    {
      title: t('COL_ACTIONS'),
      align: 'center' as const,
      width: 200,
      value: 'actions',
    },
  ];
});

onMounted(async () => {
  await reload();
});

const handleAddToken = () => {
  dialog.value?.open(); // This should work for creating new tokens
};

const doDelete = async (config: IConfirmationDialogConfig) => {
  if (config.okButton === 'BTN_RENEW') {
    // Handle renew action
    const response = await adminService.renewInternalToken(config.key);
    const renewedToken = response.data;

    // Show success message
    snack(t('DIALOG_internal_token_renew_success'));

    dialog.value?.showToken(renewedToken, true);
  } else {
    // Handle revoke action
    await adminService.revokeInternalToken(config.key);
    snack(t('DIALOG_internal_token_revoke_success'));
  }
  await reload();
};

const showRenewConfirm = async (item: InternalToken) => {
  confirmConfig.value = {
    type: ConfirmationType.DELETE,
    key: item._key!,
    name: item.name,
    description: t('CONFIRM_RENEW_TOKEN_DESCRIPTION'),
    okButton: 'BTN_RENEW',
  };
  confirmVisible.value = true;
};

const showConfirm = async (item: InternalToken) => {
  confirmConfig.value = {
    type: ConfirmationType.DELETE,
    key: item._key!,
    name: item.name,
    description: t('CONFIRM_REVOKE_TOKEN_DESCRIPTION'),
    okButton: t('BTN_REVOKE'),
  };
  confirmVisible.value = true;
};
</script>

<template>
  <TableLayout>
    <template #buttons>
      <h1 class="d-headline">{{ t('INTERNAL_TOKEN') }}</h1>
      <DCActionButton
        large
        icon="mdi-plus"
        :hint="t('INTERNAL_TOKEN')"
        :text="t('BTN_ADD')"
        @click="handleAddToken"
        class="mx-2" />
    </template>
    <template #table>
      <v-data-table
        density="compact"
        class="striped-table fill-height"
        :loading="isLoading"
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
        <template v-slot:[`item.expiry`]="{item}">
          <DDateCellWithTooltip :value="item.expiry" />
        </template>
        <template v-slot:[`item.status`]="{item}">
          <v-chip :color="getStatus(item).color" size="small" :prepend-icon="getStatus(item).icon">
            {{ getStatus(item)?.text }}
          </v-chip>
        </template>
        <template v-slot:[`item.actions`]="{item}">
          <DIconButton
            v-if="
              !item.revoked &&
              item.expiry &&
              !isNaN(new Date(item.expiry).getTime()) &&
              new Date(item.expiry) > new Date()
            "
            icon="mdi-refresh"
            @clicked="showRenewConfirm(item)" />
          <DIconButton
            v-if="!item.revoked && (!item.expiry || new Date(item.expiry) > new Date())"
            icon="mdi-delete"
            @clicked="showConfirm(item)" />
        </template>
      </v-data-table>
    </template>
  </TableLayout>
  <InternalTokenDialog ref="dialog" @reload="reload" />
  <ConfirmationDialog v-model:showDialog="confirmVisible" :config="confirmConfig" @confirm="doDelete" />
</template>
