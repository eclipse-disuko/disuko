// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {ApprovalFormHooks, useApprovalFormBase} from '@disclosure-portal/composables/approval/useApprovalFormBase';
import {DocumentMeta, InternalApprovalRequest} from '@disclosure-portal/model/ApprovalRequest';
import {ApprovableSPDXDto} from '@disclosure-portal/model/Project';
import {UserDto} from '@shared/types/Users';
import ErrorDialogConfig from '@shared/types/ErrorDialogConfig';
import projectService from '@disclosure-portal/services/projects';
import eventBus from '@shared/utils/eventbus';
import useSnackbar from '@shared/composables/useSnackbar';
import {useAppStore} from '@disclosure-portal/stores/app';
import {nextTick, ref} from 'vue';

export function useInternalApprovalForm(externalHooks: Partial<ApprovalFormHooks> = {}) {
  const snackbar = useSnackbar();
  const appStore = useAppStore();

  const approverTab = ref(0);
  const ownerApprover1 = ref('');
  const ownerApprover2 = ref('');
  const developerApprover1 = ref('');
  const developerApprover2 = ref('');
  const ownerApproverIn1 = ref();
  const ownerApproverPre1 = ref<UserDto>();
  const ownerApproverIn2 = ref();
  const ownerApproverPre2 = ref<UserDto>();
  const developerApproverIn1 = ref();
  const developerApproverPre1 = ref<UserDto>();
  const developerApproverIn2 = ref();
  const developerApproverPre2 = ref<UserDto>();
  const isVehicle = ref(false);

  const hooks: ApprovalFormHooks = {
    ...externalHooks,
    afterResetFormState: () => {
      approverTab.value = 0;
      ownerApprover1.value = '';
      ownerApprover2.value = '';
      developerApprover1.value = '';
      developerApprover2.value = '';
      ownerApproverPre1.value = undefined;
      ownerApproverPre2.value = undefined;
      developerApproverPre1.value = undefined;
      developerApproverPre2.value = undefined;
      externalHooks.afterResetFormState?.();
    },
  };

  const base = useApprovalFormBase(hooks);

  const open = async (isVehicleProject: boolean) => {
    base.idle.showIdle = true;
    isVehicle.value = isVehicleProject;
    if (base.projectModel.value.customerMeta.userFRI) {
      ownerApproverPre1.value = base.projectModel.value.customerMeta.userFRI;
    }
    if (base.projectModel.value.customerMeta.userSRI) {
      ownerApproverPre2.value = base.projectModel.value.customerMeta.userSRI;
    }
    if (base.projectModel.value.supplierExtraData.userFRI) {
      developerApproverPre1.value = base.projectModel.value.supplierExtraData.userFRI;
    }
    if (base.projectModel.value.supplierExtraData.userSRI) {
      developerApproverPre2.value = base.projectModel.value.supplierExtraData.userSRI;
    }
    base.noFOSS.value = base.projectModel.value.isNoFoss;
    base.approvableInfo.value = await projectService.getApprovableInfo(base.projectModel.value._key);

    await base.autoSelect();
    base.setDefaultFlags();
    developerApproverIn1.value?.resetForm();
    developerApproverIn2.value?.resetForm();
    base.idle.showIdle = false;
    base.isVisible.value = true;
  };

  const doDialogAction = async () => {
    await nextTick(async () => {
      base.form.value?.validate().then(async (info) => {
        if (!info.valid) {
          return;
        }
        const isDev1Valid = await developerApproverIn1.value?.validateOnCreate();
        const isDev2Valid = await developerApproverIn2.value?.validateOnCreate();
        const isOwner1Valid = await ownerApproverIn1.value?.validateOnCreate();
        const isOwner2Valid = await ownerApproverIn2.value?.validateOnCreate();
        if (!isDev1Valid || !isDev2Valid) {
          approverTab.value = 1;
          return;
        }
        if (!isOwner1Valid || !isOwner2Valid) {
          approverTab.value = 0;
          return;
        }
        if (
          (ownerApprover1.value !== '' || ownerApprover2.value !== '') &&
          ownerApprover1.value === ownerApprover2.value
        ) {
          approverTab.value = 0;
          const d = new ErrorDialogConfig();
          d.title = '' + base.t('SBOM_REQUEST_INTERNAL_APPROVAL');
          d.description = '' + base.t('EQUAL_OWNER_APPROVERS_ERROR_MESSAGE');
          eventBus.emit('on-error', {error: d});
          return;
        }
        if (
          (ownerApprover1.value !== '' && ownerApprover2.value === '') ||
          (ownerApprover1.value === '' && ownerApprover2.value !== '')
        ) {
          approverTab.value = 0;
          const d = new ErrorDialogConfig();
          d.title = '' + base.t('SBOM_REQUEST_INTERNAL_APPROVAL');
          d.description = '' + base.t('BOTH_OR_NONE_OWNER_APPROVERS_ALLOWED_ERROR_MESSAGE');
          eventBus.emit('on-error', {error: d});
          return;
        }
        if (developerApprover1.value === developerApprover2.value) {
          approverTab.value = 1;
          const d = new ErrorDialogConfig();
          d.title = '' + base.t('SBOM_REQUEST_INTERNAL_APPROVAL');
          d.description = '' + base.t('EQUAL_DEVELOPER_APPROVERS_ERROR_MESSAGE');
          eventBus.emit('on-error', {error: d});
          return;
        }
        base.idle.showIdle = true;
        base.idle.idleMessage = base.t('SBOM_REQUEST_APPROVAL_PROGRESS');
        if (!base.projectModel.value.isGroup && base.selectedSbom.value) {
          const approvableSpdx = {
            spdxkey: '',
            versionkey: '',
          } as ApprovableSPDXDto;
          approvableSpdx.spdxkey = base.selectedSbom.value?._key ?? '';
          approvableSpdx.versionkey = base.selectedChannel.value?._key ?? '';
          await projectService.updateApprovableSpdx(approvableSpdx, base.projectModel.value._key);
        }

        const metaDoc = new DocumentMeta();
        metaDoc.c1 = base.c1.value;
        metaDoc.c2 = base.c2.value;
        metaDoc.c3 = base.c3.value;
        metaDoc.c4 = base.c4.value;
        metaDoc.c5 = base.c5.value;
        metaDoc.c6 = base.noFOSS.value;

        const req: InternalApprovalRequest = {
          withZip: base.withZip.value,
          comment: base.comment.value,
          guidProject: base.projectModel.value._key,
          metaDoc: metaDoc,
          customerApprover1: ownerApprover1.value,
          customerApprover2: ownerApprover2.value,
          supplierApprover1: developerApprover1.value,
          supplierApprover2: developerApprover2.value,
          fossVersion: 'vanilla',
        };

        projectService.createInternalApproval(req, base.projectModel.value._key).then(async (response) => {
          if (response) {
            await base.jobStore.pollJobStatus(base.projectModel.value._key, response.jobKey);
            base.isVisible.value = false;
            snackbar.info(base.t('DIALOG_request_internal_approval_success'));
            appStore.setShouldReloadApprovals(true);
            if (!base.projectModel.value.isGroup) {
              await base.projectStore.fetchProjectByKey(base.projectModel.value._key);
            }
          } else {
            base.idle.showIdle = false;
            base.idle.idleMessage = '';
          }
        });
      });
    });
  };

  const close = () => {
    base.isVisible.value = false;
  };

  const dialogConfig = {
    title: base.t('SBOM_REQUEST_INTERNAL_APPROVAL'),
    secondaryButton: {text: base.t('BTN_CANCEL')},
    primaryButton: {text: base.t('BTN_REQUEST')},
  };

  return {
    ...base,
    approverTab,
    ownerApprover1,
    ownerApprover2,
    developerApprover1,
    developerApprover2,
    ownerApproverIn1,
    ownerApproverPre1,
    ownerApproverIn2,
    ownerApproverPre2,
    developerApproverIn1,
    developerApproverPre1,
    developerApproverIn2,
    developerApproverPre2,
    isVehicle,
    open,
    doDialogAction,
    close,
    dialogConfig,
  };
}
