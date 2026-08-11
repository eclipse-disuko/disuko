// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {nextTick} from 'vue';
import {describe, expect, it, vi} from 'vitest';
import S3DataIntegrity from '../S3DataIntegrity.vue';

const {getDbS3CheckStartMock, getDbS3CheckStopMock, getDbS3CheckGetResultMock} = vi.hoisted(() => ({
  getDbS3CheckStartMock: vi.fn(),
  getDbS3CheckStopMock: vi.fn(),
  getDbS3CheckGetResultMock: vi.fn(),
}));

vi.mock('@disclosure-portal/services/admin', () => ({
  default: {
    getDbS3CheckStart: getDbS3CheckStartMock,
    getDbS3CheckStop: getDbS3CheckStopMock,
    getDbS3CheckGetResult: getDbS3CheckGetResultMock,
  },
}));

const actionButtonStub = {
  template: '<button type="button">{{ text }}</button>',
  props: ['text', 'hint', 'large', 'icon'],
};

describe('S3DataIntegrity', () => {
  it('fetches the current check result on mount', async () => {
    getDbS3CheckGetResultMock.mockResolvedValue({data: '{"ok":true}'});

    const {wrapper} = mountView(S3DataIntegrity, {
      childStubs: {DCActionButton: actionButtonStub, JsonViewer: true},
    });
    await nextTick();
    await nextTick();

    expect(getDbS3CheckGetResultMock).toHaveBeenCalledTimes(1);
    expect(wrapper.find('json-viewer-stub').exists()).toBe(true);
  });

  it('starts the check and refreshes the result when the start button is clicked', async () => {
    getDbS3CheckGetResultMock.mockResolvedValue({data: null});
    getDbS3CheckStartMock.mockResolvedValue({});

    const {wrapper} = mountView(S3DataIntegrity, {
      childStubs: {DCActionButton: actionButtonStub, JsonViewer: true},
    });
    await nextTick();
    await nextTick();
    getDbS3CheckGetResultMock.mockClear();

    const buttons = wrapper.findAll('button');
    await buttons[1].trigger('click');
    await nextTick();
    await nextTick();

    expect(getDbS3CheckStartMock).toHaveBeenCalledTimes(1);
    expect(getDbS3CheckGetResultMock).toHaveBeenCalledTimes(1);
  });
});
