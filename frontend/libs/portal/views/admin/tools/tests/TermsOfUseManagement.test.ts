// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {Group, Rights} from '@shared/user/models/Rights';
import {useUserStore} from '@shared/user/stores/user.store';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {nextTick} from 'vue';
import {afterEach, beforeEach, describe, expect, it, vi} from 'vitest';
import TermsOfUseManagement from '../TermsOfUseManagement.vue';

const {getJobLatestMock, getTermsOfUseCurrentVersionMock, setJobConfigMock, startJobMock} = vi.hoisted(() => ({
  getJobLatestMock: vi.fn(),
  getTermsOfUseCurrentVersionMock: vi.fn(),
  setJobConfigMock: vi.fn(),
  startJobMock: vi.fn(),
}));

vi.mock('@disclosure-portal/services/admin', () => ({
  default: {
    getJobLatest: getJobLatestMock,
    getTermsOfUseCurrentVersion: getTermsOfUseCurrentVersionMock,
    setJobConfig: setJobConfigMock,
    startJob: startJobMock,
  },
}));

const createWrapper = () =>
  mountView(TermsOfUseManagement, {
    childStubs: {
      ConfirmationDialog: {
        template: '<div />',
        methods: {openWithoutDetails: vi.fn()},
      },
    },
    beforePiniaMount: (pinia) => {
      useUserStore(pinia).setSimpleProfileData({
        rights: {groups: [Group.UserApplicationAdmin]} as Rights,
        profile: {} as never,
        allowed: true,
      });
    },
  });

describe('TermsOfUseManagement', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    getJobLatestMock.mockResolvedValue({data: {status: 'success', updated: '2026-01-01', customRes: ''}});
    getTermsOfUseCurrentVersionMock.mockResolvedValue({termsOfUseCurrentVersion: '1.0'});
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('loads the current terms of use version on mount', async () => {
    const {wrapper} = createWrapper();
    await nextTick();
    await vi.advanceTimersByTimeAsync(1000);

    expect(getTermsOfUseCurrentVersionMock).toHaveBeenCalledTimes(1);
    const inputs = wrapper.findAll('input[type="text"]');
    expect(inputs[0].element.value).toBe('1.0');
  });

  it('starts the reset job when confirmed', async () => {
    setJobConfigMock.mockResolvedValue({});
    startJobMock.mockResolvedValue({});

    const {wrapper} = createWrapper();
    await nextTick();
    await vi.advanceTimersByTimeAsync(1000);

    const versionInputs = wrapper.findAll('input[type="text"]');
    await versionInputs[1].setValue('2.0');

    await (
      wrapper.vm as unknown as {doTriggerResetTermsAcceptanceJob: () => Promise<void>}
    ).doTriggerResetTermsAcceptanceJob();
    await nextTick();

    expect(setJobConfigMock).toHaveBeenCalledTimes(1);
    expect(startJobMock).toHaveBeenCalledTimes(1);
  });
});
