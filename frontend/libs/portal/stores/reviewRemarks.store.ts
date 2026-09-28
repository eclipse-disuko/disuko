// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {
  CommentReviewRemarkRequest,
  ReviewRemark,
  ReviewRemarkRequest,
  SetReviewRemarkStatusRequest,
} from '@disclosure-portal/model/Quality';
import {BulkSetReviewRemarkStatusRequest} from '@disclosure-portal/model/ReviewRemarkBulkOperations';
import versionService from '@disclosure-portal/services/version';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import {defineStore} from 'pinia';
import {reactive, toRefs} from 'vue';

const scopeKey = (projectKey: string, versionKey: string): string => JSON.stringify([projectKey, versionKey]);

export const useReviewRemarksStore = defineStore('reviewRemarks', () => {
  const sbomStore = useSbomStore();
  const requestIds: Record<string, number> = {};

  const state = reactive({
    remarksByScope: {} as Record<string, ReviewRemark[]>,
    loadedByScope: {} as Record<string, boolean>,
    pendingByScope: {} as Record<string, number>,
  });

  const getRemarks = (projectKey: string, versionKey: string): ReviewRemark[] =>
    state.remarksByScope[scopeKey(projectKey, versionKey)] ?? [];

  const getRemark = (projectKey: string, versionKey: string, remarkKey: string): ReviewRemark | undefined =>
    getRemarks(projectKey, versionKey).find((remark) => remark.key === remarkKey);

  const getRemarksForComponent = (
    projectKey: string,
    versionKey: string,
    sbomKey: string,
    componentId: string,
  ): ReviewRemark[] =>
    getRemarks(projectKey, versionKey).filter(
      (remark) =>
        remark.sbomId === sbomKey && remark.components?.some((component) => component.componentId === componentId),
    );

  const isLoading = (projectKey: string, versionKey: string): boolean =>
    (state.pendingByScope[scopeKey(projectKey, versionKey)] ?? 0) > 0;

  const fetchRemarks = async (projectKey: string, versionKey: string, force = false): Promise<void> => {
    if (!projectKey || !versionKey) return;

    const key = scopeKey(projectKey, versionKey);
    if (state.loadedByScope[key] && !force) return;

    const requestId = (requestIds[key] ?? 0) + 1;
    requestIds[key] = requestId;
    state.pendingByScope[key] = (state.pendingByScope[key] ?? 0) + 1;
    try {
      const response = await versionService.getReviewRemarks(projectKey, versionKey);
      if (requestIds[key] === requestId) {
        state.remarksByScope[key] = response.data;
        state.loadedByScope[key] = true;
      }
    } finally {
      state.pendingByScope[key] = Math.max((state.pendingByScope[key] ?? 1) - 1, 0);
    }
  };

  const refreshAfterMutation = async (projectKey: string, versionKey: string): Promise<void> => {
    await Promise.allSettled([fetchRemarks(projectKey, versionKey, true), sbomStore.fetchAllSBOMsFlat(true)]);
  };

  const createRemark = async (projectKey: string, versionKey: string, request: ReviewRemarkRequest) => {
    const response = await versionService.createReviewRemark(projectKey, versionKey, request);
    if (response.data.success) {
      await refreshAfterMutation(projectKey, versionKey);
    }
    return response;
  };

  const editRemark = async (
    projectKey: string,
    versionKey: string,
    remarkKey: string,
    request: ReviewRemarkRequest,
  ): Promise<void> => {
    await versionService.editReviewRemark(projectKey, versionKey, remarkKey, request);
    await refreshAfterMutation(projectKey, versionKey);
  };

  const setRemarkStatus = async (
    projectKey: string,
    versionKey: string,
    remarkKey: string,
    request: SetReviewRemarkStatusRequest,
  ): Promise<void> => {
    await versionService.setReviewRemarkStatus(projectKey, versionKey, remarkKey, request);
    await refreshAfterMutation(projectKey, versionKey);
  };

  const setBulkRemarkStatus = async (
    projectKey: string,
    versionKey: string,
    request: BulkSetReviewRemarkStatusRequest,
  ): Promise<void> => {
    await versionService.bulkSetReviewRemarkStatus(projectKey, versionKey, request);
    await refreshAfterMutation(projectKey, versionKey);
  };

  const commentOnRemark = async (
    projectKey: string,
    versionKey: string,
    remarkKey: string,
    request: CommentReviewRemarkRequest,
  ): Promise<void> => {
    await versionService.commentReviewRemark(projectKey, versionKey, remarkKey, request);
    await Promise.allSettled([fetchRemarks(projectKey, versionKey, true)]);
  };

  const reset = (): void => {
    for (const key of Object.keys(requestIds)) {
      requestIds[key] += 1;
    }
    state.remarksByScope = {};
    state.loadedByScope = {};
    state.pendingByScope = {};
  };

  return {
    ...toRefs(state),
    getRemarks,
    getRemark,
    getRemarksForComponent,
    isLoading,
    fetchRemarks,
    createRemark,
    editRemark,
    setRemarkStatus,
    setBulkRemarkStatus,
    commentOnRemark,
    reset,
  };
});
