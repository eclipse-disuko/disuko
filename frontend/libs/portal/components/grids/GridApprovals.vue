<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import {Approval, ApprovalStates, ApprovalType, ApprovalUpdate, ApproverRoles} from '@disclosure-portal/model/Approval';
import {ApprovalResponse} from '@disclosure-portal/model/ApprovalRequest';
import {FillCustomerReq} from '@disclosure-portal/model/Project';
import {UserDto} from '@disclosure-portal/model/Users';
import projectService from '@disclosure-portal/services/projects';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useUserStore} from '@disclosure-portal/stores/user';
import {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import {DataTableHeader, SortItem} from '@shared/types/table';
import {computed, onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute} from 'vue-router';

const {t} = useI18n();
const projectStore = useProjectStore();
const userStore = useUserStore();
const {info} = useSnackbar();
const route = useRoute();

const items = ref<Approval[]>([]);
const clickedApp = ref<Approval>({} as Approval);
const fillCustomerDialogOpen = ref(false);
const customer1User = ref<UserDto | null>(null);
const customer2User = ref<UserDto | null>(null);
const expanded = ref<string[]>([]);
const search = ref('');
const dataAreLoaded = ref(false);
const abortConfirmationDialogVisible = ref(false);
const confirmationDialogConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const sortBy: SortItem[] = [{key: 'created', order: 'desc'}];
const tableApprovals = ref<HTMLElement | null>(null);

const headers = computed((): DataTableHeader[] => {
  return [
    {
      title: '',
      value: 'data-table-expand',
      width: 25,
    },
    {
      title: t('COL_ACTIONS'),
      align: 'center',
      width: 80,
      value: 'actions',
      sortable: false,
    },
    {
      title: t('COL_TITLE'),
      align: 'start',
      value: 'title',
      sortable: true,
      width: 300,
      sortRaw: sortByTitle,
    },
    {
      title: t('COL_CREATED'),
      width: 150,
      align: 'start',
      value: 'created',
      sortable: true,
    },
    {
      title: t('COL_UPDATED'),
      width: 150,
      align: 'start',
      value: 'updated',
      sortable: true,
    },
    {
      title: t('COL_USER'),
      width: 160,
      sortable: true,
      align: 'start',
      value: 'creator',
    },
  ];
});

const projectModel = computed(() => projectStore.currentProject!);
const user = computed(() => userStore.getProfile.user);

const sortByTitle = (a: Approval, b: Approval): number => {
  return t('COL_APPROVAL_TITLE_TYPE_' + b.type).localeCompare(t('COL_APPROVAL_TITLE_TYPE_' + a.type));
};
const onRowExpand = (newExpanded: string[]) => {
  if (newExpanded.length > 1) {
    // Keep only the last expanded row
    expanded.value = [newExpanded[newExpanded.length - 1]];
  } else {
    expanded.value = newExpanded;
  }
};

const toggleExpand = (item: Approval) => {
  const index = expanded.value.indexOf(item.key);
  if (index > -1) {
    expanded.value.splice(index, 1);
  } else {
    expanded.value.push(item.key);
  }
};

const isExpanded = (item: Approval) => {
  return expanded.value.includes(item.key);
};

const canBeAborted = (item: Approval) => {
  return (
    (item.type === ApprovalType.Internal || item.type === ApprovalType.Plausibility) &&
    (item.status === ApprovalStates.Pending || item.status === ApprovalStates.SupplierApproved) &&
    item.creator === user.value
  );
};

const abortTooltipKey = (type: ApprovalType) => {
  let key = 'TT_TAD_BTN_ABORT';
  if (type === ApprovalType.Plausibility) {
    key = 'TT_BTN_ABORT_REVIEW';
  }
  return key;
};

const canBeFilledIn = (item: Approval) => {
  return (
    item.type === ApprovalType.Internal &&
    (item.internal.approver[ApproverRoles.Customer1] === '' ||
      item.internal.approver[ApproverRoles.Customer2] === '') &&
    (item.status === ApprovalStates.Pending || item.status === ApprovalStates.SupplierApproved)
  );
};

const getColorForApproval = (status: ApprovalStates) => {
  switch (status) {
    case ApprovalStates.Approved:
      return 'rgb(var(--v-theme-approvalApproved))';
    case ApprovalStates.Declined:
      return 'rgb(var(--v-theme-approvalDeclined))';
    case ApprovalStates.Pending:
      return 'rgb(var(--v-theme-approvalPending))';
    case ApprovalStates.CustomerApproved:
      return 'rgb(var(--v-theme-approvalApproved))';
    case ApprovalStates.SupplierApproved:
      return 'rgb(var(--v-theme-approvalApproved))';
    case ApprovalStates.Aborted:
      return 'rgb(var(--v-theme-approvalDeclined))';
    case ApprovalStates.GenerationFailed:
      return 'rgb(var(--v-theme-approvalDeclined))';
  }
};

const getApprovalStatus = (item: Approval) => {
  return t('COL_APPROVAL_STATUS_' + item.type + '_' + item.status);
};

const reloadInternal = async (forceReload: boolean) => {
  if (!forceReload && dataAreLoaded.value) {
    return;
  }
  dataAreLoaded.value = false;
  items.value = await projectService.getAllApprovals(projectModel.value._key);
  dataAreLoaded.value = true;
};

const abort = async (item: Approval) => {
  clickedApp.value = item;
  let textKey = 'DLG_CONFIRMATION_DESCRIPTION_ABORT_TASK_APPROVAL';
  if (item.type === ApprovalType.Plausibility) {
    textKey = 'DLG_CONFIRMATION_DESCRIPTION_ABORT_TASK_REVIEW';
  }
  confirmationDialogConfig.value = {
    type: ConfirmationType.NOT_SET,
    description: textKey,
    okButton: 'TAD_BTN_ABORT',
    key: '',
    name: '',
    okButtonIsDisabled: false,
  };
  abortConfirmationDialogVisible.value = true;
};

const doAbort = async () => {
  const updateItem = {} as ApprovalUpdate;
  updateItem.state = ApprovalStates.Aborted;

  await projectService.updateApproval(updateItem, projectModel.value._key, clickedApp.value.key).then((response) => {
    if (response.data as unknown as ApprovalResponse) {
      info(t('TASK_COMPLETED'));
      reloadInternal(true);
    }
  });
};

const openFillDialog = async (item: Approval) => {
  clickedApp.value = item;
  if (
    clickedApp.value.internal.approver[ApproverRoles.Customer1] &&
    clickedApp.value.internal.approver[ApproverRoles.Customer1] !== ''
  ) {
    customer1User.value = (
      await projectService.getApproverUser(clickedApp.value.key, projectModel.value._key, ApproverRoles.Customer1)
    ).data;
  } else {
    customer1User.value = null;
  }
  if (
    clickedApp.value.internal.approver[ApproverRoles.Customer2] &&
    clickedApp.value.internal.approver[ApproverRoles.Customer2] !== ''
  ) {
    customer2User.value = (
      await projectService.getApproverUser(clickedApp.value.key, projectModel.value._key, ApproverRoles.Customer2)
    ).data;
  } else {
    customer2User.value = null;
  }
  fillCustomerDialogOpen.value = true;
};

const onSaveApprover = async (customer1: string, customer2: string) => {
  await projectService
    .fillCustomer(new FillCustomerReq(customer1, customer2), projectModel.value._key, clickedApp.value.key)
    .then(async (result) => {
      if (result) {
        info(t('DIALOG_owner_fill_success'));
        await reloadInternal(true);
      }
    });
  fillCustomerDialogOpen.value = false;
};
onMounted(async () => {
  await reloadInternal(true);
});

watch(
  () => projectModel.value._key,
  async () => {
    await reloadInternal(true);
  },
);
watch(
  () => route.path,
  async (_newPath) => {
    if (_newPath.includes('approvals')) {
      await reloadInternal(true);
    }
  },
);
const appStore = useAppStore();
watch(
  () => appStore.shouldReloadApprovals,
  async (_new) => {
    if (_new) {
      await reloadInternal(true);
      appStore.setShouldReloadApprovals(false);
    }
  },
);

const getActionButtons = (item: Approval): TableActionButtonsProps['buttons'] => {
  return [
    {
      icon: 'mdi-pencil',
      hint: t('TAD_BTN_FILL'),
      event: 'fill',
      show: canBeFilledIn(item) && !projectModel.value.isDeprecated,
    },
    {
      icon: 'mdi-close',
      hint: t(abortTooltipKey(item.type)),
      event: 'abort',
      show: canBeAborted(item) && !projectModel.value.isDeprecated,
    },
  ];
};
</script>

<template>
  <TableLayout has-title has-tab>
    <template v-if="projectModel.parent" #description>
      <div>
        <p class="d-headline-2">{{ t('APPROVALS_USE_PARENT_TEXT_EXPLANATION') }}</p>
        <span class="d-subtitle-2 pt-2">{{ t('APPROVALS_USE_PARENT_TEXT_PRE') }}</span>
        <DInternalLink
          :text="t('APPROVALS_USE_PARENT_TEXT_LINK')"
          :url="'/#/dashboard/groups/' + projectModel.parent + '/approvals'" />
        <span class="d-subtitle-2 ps-1 pt-2">{{ t('APPROVALS_USE_PARENT_TEXT_POST') }}</span>
      </div>
    </template>
    <template #buttons>
      <v-spacer></v-spacer>
      <DSearchField v-model="search" />
    </template>
    <template #table>
      <div ref="tableApprovals" class="fill-height">
        <v-data-table
          :loading="!dataAreLoaded"
          item-value="key"
          :items="items"
          :headers="headers"
          v-model:search="search"
          fixed-header
          class="striped-table custom-data-table fill-height"
          :sort-by="sortBy"
          density="compact"
          expand-on-click
          :expanded.sync="expanded"
          @update:expanded="onRowExpand"
          :items-per-page="100">
          <template #[`item.data-table-expand`]="{item}">
            <v-icon color="primary" @click.stop="toggleExpand(item)">
              {{ isExpanded(item) ? 'mdi-chevron-up' : 'mdi-chevron-down' }}
            </v-icon>
          </template>
          <template #expanded-row="{columns, item}">
            <td :colspan="columns.length">
              <ApprovalInfoTabs
                v-if="item.type == ApprovalType.Internal"
                :item="item"
                :tabs-list="['history', 'general', 'details', 'documents', 'attributes']"
                task-description=""></ApprovalInfoTabs>
              <ApprovalInfoTabs
                v-if="item.type == ApprovalType.Plausibility"
                :item="item"
                :tabs-list="['generalReview', 'details']"
                task-description=""></ApprovalInfoTabs>
              <ApprovalInfoTabs
                v-if="item.type == ApprovalType.External"
                :item="item"
                :tabs-list="['generalExternal', 'details', 'documents', 'attributes']"
                @reloads-approvals="reloadInternal(true)"
                task-description=""></ApprovalInfoTabs>
            </td>
          </template>
          <template #[`item.title`]="{item}">
            {{ t('COL_APPROVAL_TITLE_TYPE_' + item.type) }} -
            <span
              v-if="item.type == ApprovalType.Internal || item.type == ApprovalType.Plausibility"
              :style="{color: getColorForApproval(item.status)}">
              {{ getApprovalStatus(item) }}
            </span>
            <span v-if="item.type == ApprovalType.External" :style="{color: getColorForApproval(item.external.state)}">
              {{ t('COL_APPROVAL_STATUS_EXTERNAL_' + item.external.state) }}
            </span>
          </template>
          <template #[`item.created`]="{item}">
            <DDateCellWithTooltip :value="item.created" />
          </template>
          <template #[`item.updated`]="{item}">
            <DDateCellWithTooltip :value="item.updated" />
          </template>
          <template #[`item.actions`]="{item}">
            <Stack direction="row" class="gap-0">
              <TableActionButtons
                variant="compact"
                :buttons="getActionButtons(item)"
                @fill="openFillDialog(item)"
                @abort="abort(item)" />
            </Stack>
          </template>
        </v-data-table></div
      ></template>
  </TableLayout>

  <ConfirmationDialog
    v-model="abortConfirmationDialogVisible"
    :config="confirmationDialogConfig"
    @confirm="doAbort"></ConfirmationDialog>
  <DFormDialog v-model:showDialog="fillCustomerDialogOpen">
    <FillCustomerApprover
      :title="t('UM_DIALOG_TITLE_FILL_OWNER')"
      :confirm-text="t('NP_DIALOG_BTN_CREATE')"
      :customer1User="customer1User!"
      :customer2User="customer2User!"
      :projectKey="projectModel._key"
      @confirm="onSaveApprover"
      @close="fillCustomerDialogOpen = false" />
  </DFormDialog>
</template>
