// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {vuetifyStubs} from '@disclosure-portal/test-utils/vuetify-stubs';
import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, describe, expect, it} from 'vitest';
import DCActionButton from '../DCActionButton.vue';

const TooltipStub = {
  template: '<div class="tooltip-stub"><slot />{{ text }}</div>',
  props: ['text'],
};

describe('DCActionButton', () => {
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

  const createWrapper = (props: Record<string, unknown> = {}, slots = {}) => {
    return mount(DCActionButton, {
      props,
      slots,
      global: {
        stubs: {...vuetifyStubs, Tooltip: TooltipStub},
      },
    });
  };

  it('renders without throwing', () => {
    const wrapper = createWrapper({text: 'Save'});

    expect(wrapper.findComponent(vuetifyStubs['v-btn']).exists()).toBe(true);
  });

  it('renders the given icon and text', () => {
    const wrapper = createWrapper({text: 'Save', icon: 'mdi-content-save'});

    expect(wrapper.text()).toContain('Save');
    expect(wrapper.findComponent(vuetifyStubs['v-icon']).exists()).toBe(true);
  });

  it('defaults to variant tonal and color primary', () => {
    const wrapper = createWrapper({text: 'Save'});
    const btn = wrapper.findComponent(vuetifyStubs['v-btn']);

    expect(btn.props('variant')).toBe('tonal');
    expect(btn.props('color')).toBe('primary');
  });

  it('forwards disabled and loading props to the v-btn', () => {
    const wrapper = createWrapper({text: 'Save', disabled: true, loading: true});
    const btn = wrapper.findComponent(vuetifyStubs['v-btn']);

    expect(btn.props('disabled')).toBe(true);
    expect(btn.props('loading')).toBe(true);
  });

  it('emits clicked when clicked', async () => {
    const wrapper = createWrapper({text: 'Save'});

    await wrapper.find('button').trigger('click');

    expect(wrapper.emitted('clicked')).toBeTruthy();
  });

  it('renders the hint tooltip when hint is given', () => {
    const wrapper = createWrapper({text: 'Save', hint: 'Save changes'});

    expect(wrapper.findComponent(TooltipStub).text()).toContain('Save changes');
  });
});
