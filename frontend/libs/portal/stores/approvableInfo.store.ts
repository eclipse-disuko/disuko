// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {defineStore} from 'pinia';
import {reactive, toRefs} from 'vue';
import {ApprovableInfo} from '@disclosure-portal/model/Approval';
import projectService from '@disclosure-portal/services/projects';
import {useProjectStore} from '@disclosure-portal/stores/project.store';

export const useApprovableInfoStore = defineStore('approvableInfo', () => {
  const projectStore = useProjectStore();

  const state = reactive({
    approvableInfoLatestSbom: new ApprovableInfo(),
    approvableInfo: new ApprovableInfo(),
    loading: false,
  });

  const fetchApprovableInfo = async (latestSbom: boolean = false) => {
    const projectKey = projectStore.currentProject?._key;
    if (!projectKey) {
      state.approvableInfoLatestSbom = new ApprovableInfo();
      state.approvableInfo = new ApprovableInfo();
      return;
    }
    state.loading = true;
    try {
      const data = await projectService.getApprovableInfo(projectKey, latestSbom);
      if (latestSbom) {
        state.approvableInfoLatestSbom = data;
      } else {
        state.approvableInfo = data;
      }
    } catch (error) {
      state.approvableInfoLatestSbom = new ApprovableInfo();
      state.approvableInfo = new ApprovableInfo();
      console.error('Error fetching approvable info:', error);
    } finally {
      state.loading = false;
    }
  };

  return {
    ...toRefs(state),
    fetchApprovableInfo,
  };
});
