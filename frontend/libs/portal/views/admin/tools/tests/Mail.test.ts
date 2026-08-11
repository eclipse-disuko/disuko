// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {describe, expect, it, vi} from 'vitest';
import Mail from '../Mail.vue';

const {sendEmailMock} = vi.hoisted(() => ({
  sendEmailMock: vi.fn(),
}));

vi.mock('@disclosure-portal/services/admin', () => ({
  default: {
    sendEmail: sendEmailMock,
  },
}));

const actionButtonStub = {
  template: '<button type="button">{{ text }}</button>',
  props: ['text', 'hint', 'large'],
};

describe('Mail', () => {
  it('renders without throwing', () => {
    const {wrapper} = mountView(Mail, {
      childStubs: {DCActionButton: actionButtonStub},
    });

    expect(wrapper.exists()).toBe(true);
  });

  it('sends mail with the entered type when the send button is clicked', async () => {
    sendEmailMock.mockResolvedValue({data: {}});

    const {wrapper} = mountView(Mail, {
      childStubs: {DCActionButton: actionButtonStub},
    });

    await wrapper.find('input[type="text"]').setValue('welcome');
    await wrapper.find('button').trigger('click');
    await wrapper.vm.$nextTick();

    expect(sendEmailMock).toHaveBeenCalledTimes(1);
    expect(sendEmailMock).toHaveBeenCalledWith(expect.objectContaining({type: 'welcome'}));
  });
});
