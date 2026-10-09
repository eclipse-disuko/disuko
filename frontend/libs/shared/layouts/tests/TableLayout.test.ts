// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {createTestingPinia} from '@pinia/testing';
import {config, mount} from '@vue/test-utils';
import {nextTick} from 'vue';
import {afterAll, beforeAll, describe, expect, it} from 'vitest';
import {useNotificationStore} from '@shared/stores/notification.store';
import TableLayout from '../TableLayout.vue';

describe('TableLayout', () => {
  const originalTMock = config.global.mocks?.$t;

  beforeAll(() => {
    if (config.global.mocks && '$t' in config.global.mocks) {
      delete config.global.mocks.$t;
    }
  });

  afterAll(() => {
    if (originalTMock) {
      config.global.mocks = {...config.global.mocks, $t: originalTMock};
    }
  });

  it('reactively reserves space for a visible notification bar', async () => {
    const pinia = createTestingPinia({stubActions: false});
    const notificationStore = useNotificationStore(pinia);
    notificationStore.message = 'System notice';

    const wrapper = mount(TableLayout, {
      global: {
        plugins: [pinia],
        stubs: {
          Stack: {template: '<div><slot /></div>'},
          UseWindowSize: {template: '<slot :height="1000" />'},
        },
      },
      slots: {
        table: '<div data-testid="content" />',
      },
    });

    const table = wrapper.find('[data-testid="content"]').element.parentElement;
    expect(table?.getAttribute('style')).toContain('height: 840px');

    notificationStore.closed = true;
    await nextTick();

    expect(table?.getAttribute('style')).toContain('height: 872px');
  });
});
