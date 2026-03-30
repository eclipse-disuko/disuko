// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import ConfirmationDialog from '@disclosure-portal/components/dialog/ConfirmationDialog.vue';
import ApprovalInfoTabs from '@disclosure-portal/components/projects/approval/ApprovalInfoTabs.vue';
import {
  Approval,
  ApprovalStates,
  ApprovalType,
  ApprovalUpdate,
  ApproverRoles,
  PowerOfAttorneyType,
} from '@disclosure-portal/model/Approval';
import {ApprovalResponse} from '@disclosure-portal/model/ApprovalRequest';
import {FillCustomerReq, ProjectModel} from '@disclosure-portal/model/Project';
import {TaskDto, UserDto} from '@disclosure-portal/model/Users';
import ProjectService from '@disclosure-portal/services/projects';
import {escapeHtml} from '@disclosure-portal/utils/Validation';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import DFormDialog from '@shared/components/disco/DFormDialog.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import {ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';

const showDialog = defineModel<boolean>('showDialog');
const props = defineProps<{selectedRow: TaskDto | null}>();

const {t} = useI18n();
const {info} = useSnackbar();

const maxCommentLength = 1000;
const responseComment = ref('');
const responseCommentError = ref('');
const approvalTaskDescription = ref('');
const plausibilityTaskDescription = ref('');
const item = ref<TaskDto>({} as TaskDto);
const approval = ref<Approval>({} as Approval);
const projectModel = ref<ProjectModel>({} as ProjectModel);
const showViewQuestion = ref(false);
const showFill = ref(false);
const showButton = ref(false);
const isReadOnly = ref(true);
const commentWith = ref(6);
const powerOfAttorney = ref<PowerOfAttorneyType | ''>('');
const powerOfAttorneyError = ref(false);
const fillCustomerDialogOpen = ref(false);
const customer1User = ref<UserDto>({} as UserDto);
const customer2User = ref<UserDto>({} as UserDto);
const abortTaskConfirmationVisible = ref(false);
const confirmationDialogConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);

const powerOfAttorneyItems = [
  {
    text: t('POWER_OF_ATTORNEY_IV'),
    value: PowerOfAttorneyType.iV,
  },
  {
    text: t('POWER_OF_ATTORNEY_PPA'),
    value: PowerOfAttorneyType.ppa,
  },
  {
    text: t('POWER_OF_ATTORNEY_OTHER'),
    value: PowerOfAttorneyType.other,
  },
];

const getColorForApproval = (status: ApprovalStates) => {
  switch (status) {
    case ApprovalStates.Approved:
      return 'var(--v-approvalApproved-base)';
    case ApprovalStates.Declined:
      return 'var(--v-approvalDeclined-base)';
    case ApprovalStates.Pending:
      return 'var(--v-approvalPending-base)';
  }
};

const reset = () => {
  item.value = {} as TaskDto;
  powerOfAttorney.value = '';
  powerOfAttorneyError.value = false;
  responseComment.value = '';
  responseCommentError.value = '';
  approvalTaskDescription.value = '';
  plausibilityTaskDescription.value = '';
};

const open = async (task: TaskDto) => {
  reset();
  isReadOnly.value = false;
  if (task.type.toUpperCase().endsWith('_INFO') || task.status.toUpperCase() === 'DONE') {
    isReadOnly.value = true;
  }
  approval.value = (await ProjectService.getApproval(task.approvalGuid, task.projectGuid)).data;
  projectModel.value = await ProjectService.get(task.projectGuid);

  showButton.value = false;
  item.value = task;
  showViewQuestion.value = false;
  showFill.value =
    approval.value.type === ApprovalType.Internal &&
    (approval.value.internal.approver[ApproverRoles.Customer1] === '' ||
      approval.value.internal.approver[ApproverRoles.Customer2] === '');

  let groupOrProject = t('BTN_PROJECT');
  let includesProjects = '';
  if (projectModel.value.isGroup) {
    groupOrProject = t('BTN_GROUP');
    let cntProjects = 0;
    if (approval.value.info.projects) {
      cntProjects = approval.value.info.projects.length;
    }
    includesProjects = t('GROUP_INCLUDES_PROJECTS', {
      cntProjects: cntProjects,
    });
  }

  approvalTaskDescription.value = t(
    approval.value.comment ? 'APPROVAL_TEXT_REQUEST' : 'APPROVAL_TEXT_REQUEST_NO_COMMENT',
    {
      requestor: escapeHtml(`${approval.value.creatorFullName} (${approval.value.creator})`),
      groupOrProject,
      includesProjects,
      projectName: escapeHtml(projectModel.value.name),
      comment: escapeHtml(approval.value.comment),
    },
  );
  if (item.value.type === 'APPROVALINFO') {
    plausibilityTaskDescription.value = t('PLAUSBILITY_TEXT_CREATOR_INFO', {
      requested: escapeHtml(
        `${approval.value.plausibility.approverFullName} (${approval.value.plausibility.approver})`,
      ),
    });
  } else {
    plausibilityTaskDescription.value = t(
      approval.value.comment ? 'PLAUSBILITY_TEXT_REQUEST' : 'PLAUSBILITY_TEXT_REQUEST_NO_COMMENT',
      {
        requestor: escapeHtml(`${approval.value.creatorFullName} (${approval.value.creator})`),
        groupOrProject,
        comment: escapeHtml(approval.value.comment),
      },
    );
  }
  commentWith.value = approval.value.type === ApprovalType.Plausibility ? 12 : 6;
};

function close() {
  showDialog.value = false;
}

const abort = () => {
  let textKey = 'DLG_CONFIRMATION_DESCRIPTION_ABORT_TASK_APPROVAL';
  if (approval.value.type === ApprovalType.Plausibility) {
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

  abortTaskConfirmationVisible.value = true;
};

const doAbort = async () => {
  const updateItem = {} as ApprovalUpdate;
  updateItem.state = ApprovalStates.Aborted;

  await ProjectService.updateApproval(updateItem, projectModel.value._key, approval.value.key).then((response) => {
    if (response.data as unknown as ApprovalResponse) {
      info(t('TASK_COMPLETED'));
      close();
    }
  });
};

const doAnswer = async (accepted: boolean) => {
  powerOfAttorneyError.value = false;
  responseCommentError.value = '';
  const updateItem = {} as ApprovalUpdate;
  updateItem.state = accepted ? ApprovalStates.Approved : ApprovalStates.Declined;
  updateItem.comment = responseComment.value;

  if (approval.value.type === ApprovalType.Internal) {
    if (accepted) {
      if (powerOfAttorney.value === '') {
        powerOfAttorneyError.value = true;
        return;
      }
      updateItem.powerOfAttorney = powerOfAttorney.value;
    }
  }
  if (approval.value.type === ApprovalType.Plausibility || approval.value.type === ApprovalType.Internal) {
    if (updateItem.comment.length > maxCommentLength) {
      responseCommentError.value = t('TAD_COMMENT_EXCEEDS_MAX_LENGTH');
      return;
    } else if (!accepted && !updateItem.comment) {
      responseCommentError.value = t('TAD_RESPONSE_COMMENT_MANDATORY');
      return;
    }
  }
  await ProjectService.updateApproval(updateItem, projectModel.value._key, approval.value.key).then((response) => {
    if (response.data) {
      info(t('TASK_COMPLETED'));
      close();
    }
  });
};

const openFillDialog = async () => {
  if (
    approval.value.internal.approver[ApproverRoles.Customer1] &&
    approval.value.internal.approver[ApproverRoles.Customer1] !== ''
  ) {
    customer1User.value = (
      await ProjectService.getApproverUser(approval.value.key, projectModel.value._key, ApproverRoles.Customer1)
    ).data;
  }
  if (
    approval.value.internal.approver[ApproverRoles.Customer2] &&
    approval.value.internal.approver[ApproverRoles.Customer2] !== ''
  ) {
    customer2User.value = (
      await ProjectService.getApproverUser(approval.value.key, projectModel.value._key, ApproverRoles.Customer2)
    ).data;
  }
  fillCustomerDialogOpen.value = true;
};

const onSaveApprover = async (customer1: string, customer2: string) => {
  ProjectService.fillCustomer(
    new FillCustomerReq(customer1, customer2),
    projectModel.value._key,
    approval.value.key,
  ).then(async (result) => {
    if (result) {
      info(t('DIALOG_owner_fill_success'));
      approval.value = (await ProjectService.getApproval(item.value.approvalGuid, item.value.projectGuid)).data;
      showFill.value = false;
    }
  });
  fillCustomerDialogOpen.value = false;
};

watch(
  () => props.selectedRow,
  async (newRow) => {
    if (newRow) {
      reset();
      const newItem = newRow ? {...newRow} : ({} as TaskDto);
      await open(newItem);
    }
  },
  {immediate: true},
);
</script>
<template>
  <v-form ref="taskApprovalForm">
    <v-dialog v-model="showDialog" content-class="large" scrollable width="800" persistent>
      <v-card class="pa-8" data-testid="task-approval-dialog">
        <v-card-title>
          <v-row>
            <v-col cols="10" class="align-center">
              <span class="text-h5">
                {{ t('COL_APPROVAL_TITLE_TYPE_' + approval.type) }}
              </span>
              <div :style="'color: ' + getColorForApproval(approval.status)" class="text-button">
                {{ t('COL_APPROVAL_STATUS_' + approval.status) }}
              </div>
            </v-col>
            <v-col cols="2" class="text-right px-0">
              <DCloseButton @click="close" />
            </v-col>
          </v-row>
        </v-card-title>
        <v-card-text class="pt-2">
          <v-row>
            <v-col cols="12" xs="12" class="pa-0">
              <ApprovalInfoTabs
                v-if="approval.type == ApprovalType.Internal"
                :item="approval"
                :projectUuid="projectModel._key"
                :project-name="projectModel.name"
                :taskDescription="approvalTaskDescription"
                :tabs-list="['task', 'details', 'history']"
                :showRedWarnDeniedDecisionsMessage="approval.info?.hasDeniedDecisions" />
              <ApprovalInfoTabs
                v-if="approval.type == ApprovalType.Plausibility"
                :item="approval"
                :projectUuid="projectModel._key"
                :project-name="projectModel.name"
                :taskDescription="plausibilityTaskDescription"
                :tabs-list="['task', 'details']" />
              <ApprovalInfoTabs
                v-if="approval.type == ApprovalType.External"
                :item="approval"
                :projectUuid="projectModel._key"
                :project-name="projectModel.name"
                :taskDescription="plausibilityTaskDescription"
                :tabs-list="['generalExternal', 'task', 'details']"
                :showRedWarnDeniedDecisionsMessage="approval.info?.hasDeniedDecisions" />
            </v-col>
          </v-row>
          <v-row class="shrink" justify="end">
            <v-col cols="6" v-if="!isReadOnly && approval.type == 'INTERNAL'">
              <v-checkbox
                v-if="item.type !== 'APPROVALINFO'"
                v-model="showButton"
                hide-details
                color="primary"
                class="shrink mr-2 mt-0 text-caption">
                <template v-slot:label>
                  <span class="custom-checkbox-label">{{ t('LABEL_CHECKBOX_TASK_APPROVAL') }}</span>
                </template>
              </v-checkbox>
            </v-col>
            <v-col :cols="commentWith" v-if="!isReadOnly && item.type == 'APPROVAL'">
              <v-row>
                <v-col cols="12" xs="12" class="py-0">
                  <v-textarea
                    rows="3"
                    auto-grow
                    variant="outlined"
                    :label="t('TAD_RESPONSE_COMMENT')"
                    v-model="responseComment"
                    :counter="maxCommentLength"
                    :error-messages="responseCommentError"></v-textarea>
                </v-col>
              </v-row>
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions class="justify-end" v-if="isReadOnly">
          <DCActionButton size="small" variant="flat" @click="close" :text="t('BTN_CLOSE')" />
        </v-card-actions>
        <v-card-actions v-else class="justify-end">
          <v-col>
            <DCActionButton
              isDialogButton
              size="small"
              variant="text"
              @click="close"
              class="mr-5"
              :text="t('BTN_CANCEL')" />
          </v-col>
          <v-col v-if="!showViewQuestion && item.type == 'APPROVALINFO'">
            <DCActionButton
              isDialogButton
              v-if="showFill"
              class="px-2"
              :text="t('TAD_BTN_FILL')"
              :hint="t('TT_TAD_BTN_FILL')"
              @click="openFillDialog"
              :disabled="approval.status == 'GENERATING'" />
            &nbsp;
            <DCActionButton
              isDialogButton
              class="px-3"
              :text="t('TAD_BTN_ABORT')"
              icon="mdi-cancel"
              :hint="t('TT_TAD_BTN_ABORT')"
              @click="abort"
              color="error"
              :disabled="approval.status == 'GENERATING'" />
          </v-col>
          <v-col v-if="item.type == 'APPROVAL'" class="text-right">
            <DCActionButton
              isDialogButton
              class="px-3"
              :text="t('TAD_BTN_DECLINE_' + approval.type)"
              icon="mdi-cancel"
              :hint="t('TT_TAD_BTN_DECLINE_' + approval.type)"
              @click="doAnswer(false)"
              color="error" />
          </v-col>
          <v-col v-if="approval.type == 'INTERNAL' && item.type !== 'APPROVALINFO'" class="mt-5 text-right">
            <v-select
              v-model="powerOfAttorney"
              :items="powerOfAttorneyItems"
              :error="powerOfAttorneyError"
              variant="outlined"
              density="compact"
              v-bind:menu-props="{location: 'bottom'}"
              :label="t('POWER_OF_ATTORNEY')"
              item-title="text">
            </v-select>
          </v-col>
          <v-col v-if="item.type == 'APPROVAL'" class="text-right">
            <DCActionButton
              isDialogButton
              :text="t('TAD_BTN_ACCEPT_' + approval.type)"
              icon="mdi-check-circle"
              :hint="t('TT_TAD_BTN_ACCEPT_' + approval.type)"
              @click="doAnswer(true)"
              color="success"
              :disabled="!showButton && approval.type == 'INTERNAL'" />
          </v-col>
        </v-card-actions>
        <ConfirmationDialog
          v-model:showDialog="abortTaskConfirmationVisible"
          :config="confirmationDialogConfig"
          @confirm="doAbort"></ConfirmationDialog>
        <DFormDialog v-model:showDialog="fillCustomerDialogOpen">
          <FillCustomerApprover
            :title="t('UM_DIALOG_TITLE_FILL_OWNER')"
            :confirm-text="t('NP_DIALOG_BTN_CREATE')"
            :customer1User="customer1User"
            :customer2User="customer2User"
            :projectKey="projectModel._key"
            @confirm="onSaveApprover"
            @close="fillCustomerDialogOpen = false" />
        </DFormDialog>
      </v-card>
    </v-dialog>
  </v-form>
</template>

<style>
.custom-checkbox-label {
  font-size: 12px;
}
</style>
