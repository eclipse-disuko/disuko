// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {JOB_EXECUTION_MANUAL, JOB_STATUS_SUCCESS} from '@disclosure-portal/model/Job';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {flushPromises} from '@vue/test-utils';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import Jobs from '../Jobs.vue';

const {getJobsAllMock, startJobMock, rerunOnetimeJobMock, snackbarInfoMock} = vi.hoisted(() => ({
  getJobsAllMock: vi.fn(),
  startJobMock: vi.fn(),
  rerunOnetimeJobMock: vi.fn(),
  snackbarInfoMock: vi.fn(),
}));

vi.mock('@disclosure-portal/services/admin', () => ({
  default: {
    getJobsAll: getJobsAllMock,
    startJob: startJobMock,
    rerunOnetimeJob: rerunOnetimeJobMock,
  },
}));

vi.mock('@shared/composables/useSnackbar', () => ({
  default: () => ({
    info: snackbarInfoMock,
  }),
}));

const job = {
  _key: 'job-1',
  name: 'Job One',
  jobType: 1,
  execution: JOB_EXECUTION_MANUAL,
  status: JOB_STATUS_SUCCESS,
  nextScheduledExecution: '',
  created: '2026-01-01T10:00:00Z',
  updated: '2026-01-02T10:00:00Z',
  config: null,
  log: [],
};

const childStubs = {
  TableLayout: {template: '<div><slot name="buttons"/><slot name="table"/></div>'},
  ConfirmationDialog: true,
  JobConfigDialog: true,
  'v-data-table': {
    props: ['items'],
    template: '<div class="v-data-table">{{ JSON.stringify(items) }}</div>',
  },
};

describe('Jobs', () => {
  beforeEach(() => {
    getJobsAllMock.mockReset();
    startJobMock.mockReset();
    rerunOnetimeJobMock.mockReset();
    snackbarInfoMock.mockReset();
    getJobsAllMock.mockResolvedValue([job]);
  });

  it('sets the expected breadcrumbs on mount', async () => {
    const {pinia} = mountView(Jobs, {childStubs});
    await flushPromises();

    const breadcrumbs = useBreadcrumbsStore(pinia);
    expect(breadcrumbs.setCurrentBreadcrumbs).toHaveBeenCalledWith([
      {title: 'Dashboard', href: '/dashboard/home', disabled: false},
      {title: 'Admin', href: '/dashboard/admin', disabled: false},
      {title: 'Jobs', href: '/dashboard/admin/jobs', disabled: false},
    ]);
  });

  it('renders without throwing and loads jobs into the table', async () => {
    const {wrapper} = mountView(Jobs, {childStubs});
    await flushPromises();

    expect(wrapper.exists()).toBe(true);
    expect(getJobsAllMock).toHaveBeenCalledTimes(1);
    expect(wrapper.find('.v-data-table').text()).toContain('Job One');
  });

  it('shows the start confirmation and starts the manual job on confirm', async () => {
    startJobMock.mockResolvedValue({data: {}});

    const {wrapper} = mountView(Jobs, {childStubs});
    await flushPromises();

    const vm = wrapper.vm as unknown as {
      showConfirm: (job: typeof job) => void;
      confirmVisible: boolean;
      confirmConfig: {key: string};
      start: () => Promise<void>;
    };
    vm.showConfirm(job);
    expect(vm.confirmVisible).toBe(true);
    expect(vm.confirmConfig.key).toBe('job-1');

    await vm.start();
    await flushPromises();

    expect(startJobMock).toHaveBeenCalledWith(job.jobType);
    expect(snackbarInfoMock).toHaveBeenCalledWith('Job started');
  });

  it('shows the one-time confirmation and reruns the one-time job on confirm', async () => {
    rerunOnetimeJobMock.mockResolvedValue({data: {}});

    const {wrapper} = mountView(Jobs, {childStubs});
    await flushPromises();

    const vm = wrapper.vm as unknown as {
      showOnetimeConfirm: (job: typeof job) => void;
      confirmOnetimeVisible: boolean;
      startOnetime: () => Promise<void>;
    };
    vm.showOnetimeConfirm(job);
    expect(vm.confirmOnetimeVisible).toBe(true);

    await vm.startOnetime();
    await flushPromises();

    expect(rerunOnetimeJobMock).toHaveBeenCalledWith('job-1');
    expect(snackbarInfoMock).toHaveBeenCalledWith('Job restarted');
  });
});
