<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script lang="ts" setup>
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import {ApprovalType} from '@disclosure-portal/model/Approval';
import {UserApproval} from '@disclosure-portal/model/Users';
import AdminService from '@disclosure-portal/services/admin';
import {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import Stack from '@shared/layouts/Stack.vue';
import {DataTableHeader} from '@shared/types/table';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute} from 'vue-router';

const props = defineProps<{
  approvals: UserApproval[];
  loading: boolean;
}>();

const emit = defineEmits<{
  reload: [];
}>();

const {t} = useI18n();
const route = useRoute();

const selected = ref<UserApproval[]>([]);
const abortTarget = ref<UserApproval | null>(null);
const abortAllPending = ref(false);
const abortSelectedPending = ref(false);
const confirmVisible = ref(false);
const confirmConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);

const items = computed(() => props.approvals.filter((a) => a.isActive));

const headers = computed<DataTableHeader[]>(() => [
  {
    title: t('COL_ACTIONS'),
    align: 'center',
    value: 'actions',
    width: 80,
    sortable: false,
  },
  {
    title: t('COL_APPROVAL_TYPE'),
    align: 'center',
    value: 'approvalType',
    width: 140,
    sortable: true,
  },
  {
    title: t('COL_IS_CREATOR'),
    align: 'center',
    value: 'isCreator',
    width: 120,
  },
  {
    title: t('COL_IS_APPROVER'),
    align: 'center',
    value: 'isApprover',
    width: 120,
  },
  {
    title: t('COL_PROJECT_NAME'),
    align: 'start',
    value: 'projectName',
    sortable: true,
  },
]);

const getActionButtons = (_: UserApproval): TableActionButtonsProps['buttons'] => [
  {
    icon: 'mdi-close',
    hint: t('TT_TAD_BTN_ABORT'),
    event: 'abort',
    show: true,
  },
];

const onAbort = (item: UserApproval) => {
  abortTarget.value = item;
  confirmConfig.value = {
    type: ConfirmationType.NOT_SET,
    description:
      item.approvalType === ApprovalType.Plausibility
        ? 'DLG_CONFIRMATION_DESCRIPTION_ABORT_TASK_REVIEW'
        : 'DLG_CONFIRMATION_DESCRIPTION_ABORT_TASK_APPROVAL',
    okButton: 'TAD_BTN_ABORT',
    key: '',
    name: '',
    okButtonIsDisabled: false,
  };
  confirmVisible.value = true;
};

const doAbort = () => {
  if (abortAllPending.value) {
    AdminService.abortAllUserApprovals(route.params.uuid as string).then(() => emit('reload'));
    abortAllPending.value = false;
  } else if (abortSelectedPending.value) {
    const keys = selected.value.map((item) => item.approvalUUID);
    AdminService.abortUserApprovals(route.params.uuid as string, keys).then(() => emit('reload'));
    selected.value = [];
    abortSelectedPending.value = false;
  } else if (abortTarget.value) {
    AdminService.abortUserApprovals(route.params.uuid as string, [abortTarget.value.approvalUUID]).then(() =>
      emit('reload'),
    );
    abortTarget.value = null;
  }
};

const onAbortAll = () => {
  abortAllPending.value = true;
  abortTarget.value = null;
  confirmConfig.value = {
    type: ConfirmationType.NOT_SET,
    description: 'DLG_CONFIRMATION_DESCRIPTION_ABORT_ALL_APPROVALS',
    okButton: 'BTN_ABORT_ALL',
    key: '',
    name: '',
    okButtonIsDisabled: false,
  };
  confirmVisible.value = true;
};

const onAbortSelected = () => {
  if (selected.value.length === 0) return;
  abortSelectedPending.value = true;
  abortTarget.value = null;
  abortAllPending.value = false;
  confirmConfig.value = {
    type: ConfirmationType.NOT_SET,
    description: 'DLG_CONFIRMATION_DESCRIPTION_ABORT_SELECTED_APPROVALS',
    okButton: 'BTN_ABORT',
    key: '',
    name: '',
    okButtonIsDisabled: false,
  };
  confirmVisible.value = true;
};
</script>

<template>
  <div class="h-[calc(100%-64px)]">
    <Stack direction="row" class="pb-3">
      <DCActionButton
        :text="t('BTN_ABORT_ALL')"
        icon="mdi-close-circle-outline"
        :hint="t('TT_ABORT_ALL_APPROVALS')"
        @click="onAbortAll" />
      <DCActionButton
        v-if="selected.length > 0"
        :text="`${t('BTN_ABORT')} (${selected.length})`"
        icon="mdi-close-circle-outline"
        :hint="t('TT_ABORT_SELECTED_APPROVALS')"
        @click="onAbortSelected" />
    </Stack>

    <v-data-table
      v-model="selected"
      :loading="loading"
      :items="items"
      :headers="headers"
      fixed-header
      item-key="approvalUUID"
      class="striped-table custom-data-table h-full"
      show-select
      return-object
      :items-per-page="50">
      <template #item.approvalType="{item}">
        {{ t('COL_APPROVAL_TITLE_TYPE_' + item.approvalType) }}
      </template>
      <template #item.actions="{item}">
        <TableActionButtons variant="compact" :buttons="getActionButtons(item)" @abort="onAbort(item)" />
      </template>
      <template #item.isCreator="{item}">
        <v-icon icon="mdi-check" :color="item.isCreator ? 'primary' : 'tableBorderColor'"></v-icon>
      </template>
      <template #item.isApprover="{item}">
        <v-icon icon="mdi-check" :color="item.isApprover ? 'primary' : 'tableBorderColor'"></v-icon>
      </template>
    </v-data-table>

    <ConfirmationDialog v-model="confirmVisible" :config="confirmConfig" @confirm="doAbort" />
  </div>
</template>
