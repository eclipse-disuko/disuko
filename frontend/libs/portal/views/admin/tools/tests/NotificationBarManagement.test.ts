// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useAppStore} from '@disclosure-portal/stores/app';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {nextTick} from 'vue';
import {describe, expect, it, vi} from 'vitest';
import NotificationBarManagement from '../NotificationBarManagement.vue';

const {getNotificationMock, setNotificationMock} = vi.hoisted(() => ({
  getNotificationMock: vi.fn(),
  setNotificationMock: vi.fn(),
}));

vi.mock('@disclosure-portal/services/admin', () => ({
  default: {
    getNotification: getNotificationMock,
    setNotification: setNotificationMock,
  },
}));

const actionButtonStub = {
  template: '<button type="button" :disabled="disabled">{{ text }}</button>',
  props: ['text', 'hint', 'large', 'disabled', 'loading'],
};

describe('NotificationBarManagement', () => {
  it('loads the current notification on mount', async () => {
    getNotificationMock.mockResolvedValue({data: {enabled: true, text: 'Existing notice'}});

    const {wrapper} = mountView(NotificationBarManagement, {
      childStubs: {DCActionButton: actionButtonStub},
    });
    await nextTick();
    await nextTick();

    expect(getNotificationMock).toHaveBeenCalledTimes(1);
    expect(wrapper.find('textarea').element.value).toBe('Existing notice');
  });

  it('toggles the notification and updates the app store when posted', async () => {
    getNotificationMock.mockResolvedValue({data: {enabled: false, text: ''}});
    setNotificationMock.mockResolvedValue({data: {}});

    const {wrapper, pinia} = mountView(NotificationBarManagement, {
      childStubs: {DCActionButton: actionButtonStub},
    });
    await nextTick();
    await nextTick();

    await wrapper.find('textarea').setValue('New notice');

    const buttons = wrapper.findAll('button');
    await buttons[1].trigger('click');
    await nextTick();

    expect(setNotificationMock).toHaveBeenCalledWith({enabled: true, text: 'New notice'});

    const appStore = useAppStore(pinia);
    expect(appStore.notificationMessage).toBe('New notice');
  });
});
