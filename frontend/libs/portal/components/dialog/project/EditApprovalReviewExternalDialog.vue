// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {Approval, ApprovalStates, ApprovalUpdate} from '@disclosure-portal/model/Approval';
import {IDefaultSelectItem} from '@disclosure-portal/model/IObligation';
import projectService from '@disclosure-portal/services/projects';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import _ from 'lodash';
import {ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const {info: infoSnackbar} = useSnackbar();

const showDialog = defineModel<boolean>('showDialog');

const props = defineProps<{
  approval: Approval | null;
}>();

const emit = defineEmits<{
  reload: [];
}>();

const editingApproval = ref<Approval | null>(null);
const selectedState = ref(ApprovalStates.Pending);
const comment = ref('');

function close() {
  showDialog.value = false;
}

const possibleStates: IDefaultSelectItem[] = [
  {text: t('COL_APPROVAL_STATUS_EXTERNAL_SUPPLIER_APPROVED'), value: ApprovalStates.SupplierApproved},
  {text: t('COL_APPROVAL_STATUS_EXTERNAL_CUSTOMER_APPROVED'), value: ApprovalStates.CustomerApproved},
  {text: t('COL_APPROVAL_STATUS_EXTERNAL_PENDING'), value: ApprovalStates.Pending},
  {text: t('COL_APPROVAL_STATUS_EXTERNAL_ABORTED'), value: ApprovalStates.Aborted},
  {text: t('COL_APPROVAL_STATUS_EXTERNAL_DECLINED'), value: ApprovalStates.Declined},
];

watch(
  () => props.approval,
  (newApproval) => {
    if (newApproval) {
      editingApproval.value = _.cloneDeep(newApproval);
      selectedState.value = newApproval.external.state;
      comment.value = newApproval.external.comment;
    }
  },
  {immediate: true},
);

const save = async () => {
  if (editingApproval.value) {
    editingApproval.value.external.state = selectedState.value;
    editingApproval.value.external.comment = comment.value;
    const approvalUpdate = new ApprovalUpdate();
    approvalUpdate.comment = comment.value;
    approvalUpdate.state = selectedState.value;
    await projectService.updateApproval(approvalUpdate, editingApproval.value.projectKey, editingApproval.value.key);
    emit('reload');
    close();
    infoSnackbar(t('DIALOG_UPDATE_APPROVAL_REVIEW_EXTERNAL_SUCCESS'));
  }
};
</script>

<template>
  <v-form ref="editExternalApprovalReviewForm">
    <v-dialog v-model="showDialog" width="650">
      <v-card class="pa-8">
        <v-card-title>
          <v-row>
            <v-col cols="10">
              <h1 class="d-headline">{{ t('HEADER_UPDATE_APPROVAL_REVIEW_EXTERNAL') }}</h1>
            </v-col>
            <v-col cols="2" align="right">
              <DCloseButton @click="close" />
            </v-col>
          </v-row>
        </v-card-title>
        <v-card-text>
          <v-row>
            <v-col>
              <v-select
                variant="outlined"
                density="compact"
                :items="possibleStates"
                v-bind:menu-props="{location: 'bottom'}"
                v-model="selectedState"
                :label="t('SELECT_APPROVAL_STATE')"
                item-title="text">
              </v-select>
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-textarea
                auto-grow
                variant="outlined"
                :label="t('COL_REVIEWER_COMMENT')"
                v-model="comment"
                :counter="1000"></v-textarea>
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <DCActionButton
            isDialogButton
            size="small"
            variant="text"
            @click="close"
            class="mr-4"
            :text="t('BTN_CANCEL')" />
          <DCActionButton
            size="small"
            variant="flat"
            color="primary"
            isDialogButton
            @click="save"
            :text="t('Btn_save')" />
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-form>
</template>
