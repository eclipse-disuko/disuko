// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import eventBus from '@shared/utils/eventbus';
import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, describe, expect, it} from 'vitest';
import DSnackbar from '../DSnackbar.vue';

const VSnackbarStub = {
  props: ['modelValue', 'timeout', 'color'],
  emits: ['update:modelValue'],
  template: '<div v-if="modelValue" class="v-snackbar" :data-timeout="timeout" :data-color="color"><slot /></div>',
};

describe('DSnackbar', () => {
  const originalTMock = config.global.mocks?.$t;

  beforeAll(() => {
    if (config.global.mocks && '$t' in config.global.mocks) {
      delete config.global.mocks.$t;
    }
  });

  afterAll(() => {
    if (originalTMock) {
      config.global.mocks = {
        ...config.global.mocks,
        $t: originalTMock,
      };
    }
  });

  const createWrapper = () => {
    return mount(DSnackbar, {
      global: {stubs: {'v-snackbar': VSnackbarStub}},
    });
  };

  it('is hidden until a show-snackbar event is emitted', () => {
    const wrapper = createWrapper();

    expect(wrapper.find('.v-snackbar').exists()).toBe(false);
  });

  it('shows the message with the default timeout when eventBus emits show-snackbar', async () => {
    const wrapper = createWrapper();

    eventBus.emit('show-snackbar', {message: 'Saved successfully', level: 'info'});
    await wrapper.vm.$nextTick();

    const snackbar = wrapper.find('.v-snackbar');
    expect(snackbar.exists()).toBe(true);
    expect(snackbar.text()).toContain('Saved successfully');
    expect(snackbar.attributes('data-timeout')).toBe('3000');
  });

  it('uses a custom timeout when provided', async () => {
    const wrapper = createWrapper();

    eventBus.emit('show-snackbar', {message: 'Custom timeout', timeout: 1000, level: 'info'});
    await wrapper.vm.$nextTick();

    expect(wrapper.find('.v-snackbar').attributes('data-timeout')).toBe('1000');
  });

  it('stops listening for events after being unmounted', async () => {
    const wrapper = createWrapper();
    wrapper.unmount();

    eventBus.emit('show-snackbar', {message: 'After unmount', level: 'info'});
    await wrapper.vm.$nextTick();

    expect(wrapper.find('.v-snackbar').exists()).toBe(false);
  });
});
