// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {vuetifyStubs} from '@disclosure-portal/test-utils/vuetify-stubs';
import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, describe, expect, it} from 'vitest';
import MenuItem from '../MenuItem.vue';

const StackStub = {
  template: '<div class="stack-stub"><slot /></div>',
};
const TooltipStub = {
  template: '<div class="tooltip-stub"><slot /></div>',
};

describe('MenuItem', () => {
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

  const createWrapper = (props: Record<string, unknown>, slots = {}) => {
    return mount(MenuItem, {
      props,
      slots,
      global: {
        stubs: {...vuetifyStubs, Stack: StackStub, Tooltip: TooltipStub},
      },
    });
  };

  it('renders without throwing', () => {
    const wrapper = createWrapper({text: 'Edit'});

    expect(wrapper.findComponent(vuetifyStubs['v-list-item']).exists()).toBe(true);
  });

  it('renders the given text and icon', () => {
    const wrapper = createWrapper({text: 'Edit', icon: 'mdi-pencil'});

    expect(wrapper.text()).toContain('Edit');
    expect(wrapper.findComponent(vuetifyStubs['v-icon']).exists()).toBe(true);
  });

  it('emits click with the event when clicked and not disabled', async () => {
    const wrapper = createWrapper({text: 'Edit'});

    await wrapper.find('.v-list-item').trigger('click');

    expect(wrapper.emitted('click')).toBeTruthy();
  });

  it('does not emit click when disabled', async () => {
    const wrapper = createWrapper({text: 'Edit', disabled: true});

    await wrapper.find('.v-list-item').trigger('click');

    expect(wrapper.emitted('click')).toBeFalsy();
  });

  it('renders the tooltip when a tooltip prop is given', () => {
    const wrapper = createWrapper({text: 'Edit', tooltip: 'Some hint'});

    expect(wrapper.findComponent(TooltipStub).text()).toContain('Some hint');
  });
});
