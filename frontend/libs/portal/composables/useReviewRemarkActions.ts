// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import {ReviewRemark, ReviewRemarkStatus, SetReviewRemarkStatusRequest} from '@disclosure-portal/model/Quality';
import {useReviewRemarksStore} from '@disclosure-portal/stores/reviewRemarks.store';
import useSnackbar from '@shared/composables/useSnackbar';
import {ref} from 'vue';
import {useI18n} from 'vue-i18n';

export const useReviewRemarkActions = () => {
  const {t} = useI18n();
  const {info: snack} = useSnackbar();
  const reviewRemarksStore = useReviewRemarksStore();

  const confirmCloseConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
  const confirmCancelConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
  const confirmReopenConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
  const confirmInProgressConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);

  const closeVisible = ref(false);
  const cancelVisible = ref(false);
  const reopenVisible = ref(false);
  const inProgressVisible = ref(false);

  const isOpen = (item: ReviewRemark): boolean => {
    return item.status === ReviewRemarkStatus.OPEN;
  };

  const isInProgress = (item: ReviewRemark): boolean => {
    return item.status === ReviewRemarkStatus.IN_PROGRESS;
  };

  const isClosed = (item: ReviewRemark): boolean => {
    return item.status === ReviewRemarkStatus.CLOSED;
  };

  const isCancelled = (item: ReviewRemark): boolean => {
    return item.status === ReviewRemarkStatus.CANCELLED;
  };

  const openCloseRemarkDialog = (remark: ReviewRemark) => {
    confirmCloseConfig.value = {
      key: remark.key,
      name: '',
      type: ConfirmationType.CONFIRM,
      description: t('DLG_CONFIRMATION_DESCRIPTION_CLOSE_REMARK', {name: remark.title}),
      okButton: 'Btn_confirm',
    };
    closeVisible.value = true;
  };

  const openCancelRemarkDialog = (remark: ReviewRemark) => {
    confirmCancelConfig.value = {
      key: remark.key,
      name: '',
      type: ConfirmationType.CONFIRM,
      description: t('DLG_CONFIRMATION_DESCRIPTION_CANCEL_REMARK', {name: remark.title}),
      okButton: 'Btn_confirm',
    };
    cancelVisible.value = true;
  };

  const openReopenRemarkDialog = (remark: ReviewRemark) => {
    confirmReopenConfig.value = {
      key: remark.key,
      name: '',
      type: ConfirmationType.CONFIRM,
      description: t('DLG_CONFIRMATION_DESCRIPTION_REOPEN_REMARK', {name: remark.title}),
      okButton: 'Btn_confirm',
    };
    reopenVisible.value = true;
  };

  const openInProgressRemarkDialog = (remark: ReviewRemark) => {
    confirmInProgressConfig.value = {
      key: remark.key,
      name: '',
      type: ConfirmationType.CONFIRM,
      description: t('DLG_CONFIRMATION_DESCRIPTION_IN_PROGRESS_REMARK', {name: remark.title}),
      okButton: 'Btn_confirm',
    };
    inProgressVisible.value = true;
  };

  const doCloseRemark = async (config: IConfirmationDialogConfig, projectKey: string, versionKey: string) => {
    const req: SetReviewRemarkStatusRequest = {
      status: ReviewRemarkStatus.CLOSED,
    };
    await reviewRemarksStore.setRemarkStatus(projectKey, versionKey, config.key, req);
    snack(t('DIALOG_remark_closed'));
  };

  const doCancelRemark = async (config: IConfirmationDialogConfig, projectKey: string, versionKey: string) => {
    const req: SetReviewRemarkStatusRequest = {
      status: ReviewRemarkStatus.CANCELLED,
    };
    await reviewRemarksStore.setRemarkStatus(projectKey, versionKey, config.key, req);
    snack(t('DIALOG_remark_cancelled'));
  };

  const doReopenRemark = async (config: IConfirmationDialogConfig, projectKey: string, versionKey: string) => {
    const req: SetReviewRemarkStatusRequest = {
      status: ReviewRemarkStatus.OPEN,
    };
    await reviewRemarksStore.setRemarkStatus(projectKey, versionKey, config.key, req);
    snack(t('DIALOG_remark_reopened'));
  };

  const doMarkInProgress = async (config: IConfirmationDialogConfig, projectKey: string, versionKey: string) => {
    const req: SetReviewRemarkStatusRequest = {
      status: ReviewRemarkStatus.IN_PROGRESS,
    };
    await reviewRemarksStore.setRemarkStatus(projectKey, versionKey, config.key, req);
    snack(t('DIALOG_remark_in_progress'));
  };

  return {
    confirmCloseConfig,
    confirmCancelConfig,
    confirmReopenConfig,
    confirmInProgressConfig,
    closeVisible,
    cancelVisible,
    reopenVisible,
    inProgressVisible,
    isOpen,
    isInProgress,
    isClosed,
    isCancelled,
    openCloseRemarkDialog,
    openCancelRemarkDialog,
    openReopenRemarkDialog,
    openInProgressRemarkDialog,
    doCloseRemark,
    doCancelRemark,
    doReopenRemark,
    doMarkInProgress,
  };
};
