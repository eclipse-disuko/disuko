<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import EditApprovalReviewExternalDialog from '@disclosure-portal/components/dialog/project/EditApprovalReviewExternalDialog.vue';
import {Approval, ApprovalStates} from '@disclosure-portal/model/Approval';
import {Group} from '@disclosure-portal/model/Rights';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {escapeHtml} from '@disclosure-portal/utils/Validation';
import DIconButton from '@shared/components/disco/DIconButton.vue';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';

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
  return projectStore.currentProject!.accessRights.groups.includes(Group.ProjectOwner);
});

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
  <table width="100%" class="pb-8">
    <thead>
      <tr>
        <td width="150" class="pa-2 font-weight-bold">
          {{ t('COL_APPROVAL_HISTORY_STATE') }}
        </td>
        <td width="200" class="pa-2 font-weight-bold">
          {{ t('COL_APPROVAL_REVIEW_EXTERNAL_CREATOR') }}
        </td>
        <td width="400" class="pa-2 font-weight-bold">
          {{ t('COL_REQUESTER_COMMENT') }}
        </td>
        <td width="400" class="pa-2 font-weight-bold">
          {{ t('COL_REVIEWER_COMMENT') }}
        </td>
        <td v-if="isOwner" width="100" class="pa-2 font-weight-bold">
          {{ t('COL_ACTIONS') }}
        </td>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td class="pa-2">
          <span
            class="userState"
            v-if="externalApproval.external.state"
            :style="'color: ' + getColorForApproval(externalApproval.external.state)"
            >{{ t('COL_APPROVAL_STATUS_EXTERNAL_' + externalApproval.external.state) }}</span
          >
        </td>
        <td class="pa-2">
          <span class="userTitle">{{ externalApproval.creatorFullName }}</span>
        </td>
        <td class="pa-2" style="max-width: 40px; word-wrap: break-word">
          <span v-html="urlify(externalApproval.comment)"></span>
        </td>
        <td class="pa-2" style="max-width: 40px; word-wrap: break-word">
          <span v-html="urlify(externalApproval.external.comment)"></span>
        </td>
        <td v-if="isOwner">
          <DIconButton
            icon="mdi-pencil"
            :hint="t('TT_UPDATE_APPROVAL_REVIEW_EXTERNAL')"
            @clicked="showEditExternalApprovalDialog"
            :disabled="externalApproval.external.state == 'GENERATING'" />
        </td>
      </tr>
    </tbody>
    <EditApprovalReviewExternalDialog
      v-model:showDialog="editApprovalReviewExternalVisible"
      :approval="editingExternalApproval"
      @reload="reload"></EditApprovalReviewExternalDialog>
  </table>
</template>
