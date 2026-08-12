// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {vuetifyStubs} from '@disclosure-portal/test-utils/vuetify-stubs';
import {mount} from '@vue/test-utils';
import {describe, expect, it, vi} from 'vitest';
import DCloseButton from '../DCloseButton.vue';

describe('DCloseButton', () => {
  const createWrapper = (attrs: Record<string, unknown> = {}) => {
    return mount(DCloseButton, {
      attrs,
      global: {
        stubs: {...vuetifyStubs},
      },
    });
  };

  it('renders without throwing', () => {
    const wrapper = createWrapper();

    expect(wrapper.findComponent(vuetifyStubs['v-btn']).exists()).toBe(true);
  });

  it('forwards the close icon and aria-label to the underlying v-btn', () => {
    const wrapper = createWrapper();
    const btn = wrapper.findComponent(vuetifyStubs['v-btn']);

    expect(btn.props('icon')).toBe('mdi-close');
    expect(wrapper.attributes('aria-label')).toBe('close');
  });

  it('forwards arbitrary attributes (e.g. disabled) via $attrs to the v-btn', () => {
    const wrapper = createWrapper({disabled: true});
    const btn = wrapper.findComponent(vuetifyStubs['v-btn']);

    expect(btn.props('disabled')).toBe(true);
  });

  it('invokes a click listener forwarded via $attrs when clicked', async () => {
    const onClick = vi.fn();
    const wrapper = createWrapper({onClick});

    await wrapper.find('button').trigger('click');

    expect(onClick).toHaveBeenCalledTimes(1);
  });
});
