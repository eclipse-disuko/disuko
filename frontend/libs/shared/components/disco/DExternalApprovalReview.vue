<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {Approval, ApprovalStates} from '@disclosure-portal/model/Approval';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {escapeHtml} from '@disclosure-portal/utils/Validation';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {DataTableHeader} from '@shared/types/table';

const {t} = useI18n();
const projectStore = useProjectStore();

const props = defineProps<{
  externalApproval: Approval;
}>();

const emit = defineEmits<{
  reloading: [];
}>();

const editApprovalReviewExternalVisible = ref(false);
const editingExternalApproval = ref<Approval | null>(null);

const isOwner = computed(() => {
  return projectStore.currentProject!.isProjectOwner;
});

const headers = computed((): DataTableHeader[] => [
  ...(isOwner
    ? [
        {
          title: t('COL_ACTIONS'),
          key: 'actions',
          align: 'center',
          width: 100,
          sortable: false,
        } as DataTableHeader,
      ]
    : []),
  {
    title: t('COL_APPROVAL_HISTORY_STATE'),
    align: 'start',
    width: 150,
    value: 'external.state',
    sortable: false,
  },
  {
    title: t('COL_APPROVAL_REVIEW_EXTERNAL_CREATOR'),
    align: 'start',
    value: 'creatorFullName',
    sortable: false,
    width: 200,
  },
  {
    title: t('COL_REQUESTER_COMMENT'),
    width: 200,
    align: 'start',
    value: 'comment',
    sortable: false,
  },
  {
    title: t('COL_REVIEWER_COMMENT'),
    width: 200,
    align: 'start',
    value: 'external.comment',
    sortable: false,
  },
]);

const reload = async () => {
  emit('reloading');
};

const getColorForApproval = (status: ApprovalStates) => {
  switch (status) {
    case ApprovalStates.Approved:
      return 'var(--v-approvalApproved-base)';
    case ApprovalStates.Declined:
      return 'var(--v-approvalDeclined-base)';
    case ApprovalStates.Pending:
      return 'var(--v-approvalPending-base)';
    case ApprovalStates.CustomerApproved:
      return 'var(--v-approvalApproved-base)';
    case ApprovalStates.SupplierApproved:
      return 'var(--v-approvalApproved-base)';
    case ApprovalStates.Aborted:
      return 'var(--v-approvalDeclined-base)';
    case ApprovalStates.GenerationFailed:
      return 'var(--v-approvalDeclined-base)';
  }
};

const showEditExternalApprovalDialog = () => {
  editingExternalApproval.value = {...props.externalApproval};
  editApprovalReviewExternalVisible.value = true;
};

const urlify = (text: string) => {
  text = escapeHtml(text);
  const urlRegex = /(https?:\/\/[^\s]+)/g;

  return text.replace(urlRegex, (url) => {
    return '<a target="_blank" href="' + url + '">' + url + '</a>';
  });
};
</script>

<template>
  <v-data-table
    :headers="headers"
    :items="[externalApproval]"
    density="compact"
    fixed-header
    hide-default-footer
    class="w-full pb-8">
    <template #[`item.comment`]="{item}">
      <span class="userState" v-if="item.external.state" :style="'color: ' + getColorForApproval(item.external.state)">
        {{ t('COL_APPROVAL_STATUS_EXTERNAL_' + item.external.state) }}
      </span>
    </template>
    <template #[`item.external.state`]="{item}">
      <span v-html="urlify(item.comment)"></span>
    </template>
    <template #[`item.external.comment`]="{item}">
      <span v-html="urlify(item.external.comment)"></span>
    </template>
    <template #[`item.actions`]="{item}">
      <DIconButton
        icon="mdi-pencil"
        :hint="t('TT_UPDATE_APPROVAL_REVIEW_EXTERNAL')"
        :disabled="item.external.state == 'GENERATING'"
        @clicked="showEditExternalApprovalDialog" />
    </template>
  </v-data-table>

  <EditApprovalReviewExternalDialog
    v-model:showDialog="editApprovalReviewExternalVisible"
    :approval="editingExternalApproval"
    @reload="reload" />
</template>
