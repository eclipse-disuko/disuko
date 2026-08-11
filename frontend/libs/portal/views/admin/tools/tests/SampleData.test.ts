// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {nextTick} from 'vue';
import {describe, expect, it, vi, beforeEach} from 'vitest';
import SampleData from '../SampleData.vue';

const {triggerCreateSampleDataMock, stopCreateSampleDataMock, getCreateSampleDataStateMock} = vi.hoisted(() => ({
  triggerCreateSampleDataMock: vi.fn(),
  stopCreateSampleDataMock: vi.fn(),
  getCreateSampleDataStateMock: vi.fn(),
}));

vi.mock('@disclosure-portal/services/admin', () => ({
  default: {
    triggerCreateSampleData: triggerCreateSampleDataMock,
    stopCreateSampleData: stopCreateSampleDataMock,
    getCreateSampleDataState: getCreateSampleDataStateMock,
  },
}));

const actionButtonStub = {
  template: '<button type="button">{{ text }}</button>',
  props: ['text', 'hint', 'large'],
};

describe('SampleData', () => {
  beforeEach(() => {
    triggerCreateSampleDataMock.mockReset();
    stopCreateSampleDataMock.mockReset();
    getCreateSampleDataStateMock.mockReset();
  });

  it('renders without throwing', () => {
    const {wrapper} = mountView(SampleData, {
      childStubs: {DCActionButton: actionButtonStub},
    });

    expect(wrapper.exists()).toBe(true);
  });

  it('triggers sample data creation and then fetches the resulting state', async () => {
    triggerCreateSampleDataMock.mockResolvedValue({data: 'started'});
    getCreateSampleDataStateMock.mockResolvedValue({
      data: {isRunning: true, createdCount: 3, targetCount: 10, startTime: '', endTime: '', reqID: 'r1'},
    });

    const {wrapper} = mountView(SampleData, {
      childStubs: {DCActionButton: actionButtonStub},
    });

    const buttons = wrapper.findAll('button');
    await buttons[0].trigger('click');
    await nextTick();
    await nextTick();

    expect(triggerCreateSampleDataMock).toHaveBeenCalledWith(1, false);
    expect(getCreateSampleDataStateMock).toHaveBeenCalledTimes(1);
    expect(wrapper.text()).toContain('r1');
  });

  it('stops sample data creation when the stop button is clicked', async () => {
    stopCreateSampleDataMock.mockResolvedValue({data: 'stopped'});

    const {wrapper} = mountView(SampleData, {
      childStubs: {DCActionButton: actionButtonStub},
    });

    const buttons = wrapper.findAll('button');
    await buttons[3].trigger('click');
    await nextTick();

    expect(stopCreateSampleDataMock).toHaveBeenCalledTimes(1);
  });
});
