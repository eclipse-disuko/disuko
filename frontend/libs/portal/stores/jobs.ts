// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {
  JOB_STATUS_FAILURE,
  JOB_STATUS_IDLE,
  JOB_STATUS_IN_PROGRESS,
  JOB_STATUS_SUCCESS,
} from '@disclosure-portal/model/Job';
import projectService from '@disclosure-portal/services/projects';
import {useIdleStore} from '@disclosure-portal/stores/idle.store';
import {defineStore} from 'pinia';
import {reactive, toRefs} from 'vue';
import {useI18n} from 'vue-i18n';

const jobStatuses = JOB_STATUS_IDLE | JOB_STATUS_FAILURE | JOB_STATUS_SUCCESS | JOB_STATUS_IN_PROGRESS | -1;
export type JobStatus = typeof jobStatuses;
export const useJobStore = defineStore('jobs', () => {
  const state = reactive({
    jobStatus: -1 as JobStatus,
  });

  const idle = useIdleStore();
  const {t} = useI18n();

  const getJobByKey = async (projectId: string, jobId: string) => {
    const job = await projectService.getJobByKey(projectId, jobId);
    state.jobStatus = job.data.status;
  };

  const pollJobStatus = async (projectId: string, jobId: string, interval = 1000) => {
    idle.showIdle = true;
    idle.idleMessage = t('TITLE_IDLE_WAIT');

    return new Promise((resolve, reject) => {
      const poll = async () => {
        try {
          const job = await getJobByKey(projectId, jobId);
          if (state.jobStatus === JOB_STATUS_SUCCESS) {
            resolve(job);
            idle.showIdle = false;
            return;
          } else if (state.jobStatus === JOB_STATUS_FAILURE) {
            reject(new Error('Job failed'));
            idle.showIdle = false;
            return;
          }

          setTimeout(poll, interval);
        } catch (error) {
          console.error('Error while polling:', error);
          idle.showIdle = false;
          idle.idleMessage = '';
          reject(error);
        }
      };

      poll();
    });
  };

  return {
    ...toRefs(state),
    getJobByKey,
    pollJobStatus,
  };
});
