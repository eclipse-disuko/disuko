// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {vuetifyStubs} from '@disclosure-portal/test-utils/vuetify-stubs';
import {mount} from '@vue/test-utils';
import {describe, expect, it} from 'vitest';
import ExtraMenu from '../ExtraMenu.vue';

describe('ExtraMenu', () => {
  const createWrapper = (slots = {}) => {
    return mount(ExtraMenu, {
      slots,
      global: {
        stubs: {...vuetifyStubs},
      },
    });
  };

  it('renders without throwing', () => {
    const wrapper = createWrapper();

    expect(wrapper.findComponent(vuetifyStubs['v-btn']).exists()).toBe(true);
  });

  it('renders the mdi-dots-horizontal icon button', () => {
    const wrapper = createWrapper();
    const btn = wrapper.findComponent(vuetifyStubs['v-btn']);

    expect(btn.props('icon')).toBe('mdi-dots-horizontal');
  });

  it('renders the default slot content inside the menu list', () => {
    const wrapper = createWrapper({default: '<button class="my-action">Do it</button>'});

    expect(wrapper.find('.my-action').exists()).toBe(true);
  });
});
