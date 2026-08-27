// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {vuetifyStubs} from '@disclosure-portal/test-utils/vuetify-stubs';
import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, describe, expect, it} from 'vitest';
import DIconButton from '../DIconButton.vue';

const TooltipStub = {
  template: '<div class="tooltip-stub">{{ text }}</div>',
  props: ['text', 'location'],
};

describe('DIconButton', () => {
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

  const createWrapper = (props: Record<string, unknown>) => {
    return mount(DIconButton, {
      props,
      global: {
        stubs: {...vuetifyStubs, Tooltip: TooltipStub},
      },
    });
  };

  it('renders without throwing', () => {
    const wrapper = createWrapper({icon: 'mdi-pencil'});

    expect(wrapper.findComponent(vuetifyStubs['v-btn']).exists()).toBe(true);
  });

  it('forwards the icon, color and disabled props to the v-btn', () => {
    const wrapper = createWrapper({icon: 'mdi-pencil', color: 'red', disabled: true});
    const btn = wrapper.findComponent(vuetifyStubs['v-btn']);

    expect(btn.props('icon')).toBe('mdi-pencil');
    expect(btn.props('color')).toBe('red');
    expect(btn.props('disabled')).toBe(true);
  });

  it('defaults the color to primary when not given', () => {
    const wrapper = createWrapper({icon: 'mdi-pencil'});

    expect(wrapper.findComponent(vuetifyStubs['v-btn']).props('color')).toBe('primary');
  });

  it('emits clicked when the button is clicked', async () => {
    const wrapper = createWrapper({icon: 'mdi-pencil'});

    await wrapper.find('button').trigger('click');

    expect(wrapper.emitted('clicked')).toBeTruthy();
  });

  it('renders the hint tooltip only when hint is given', () => {
    const withHint = createWrapper({icon: 'mdi-pencil', hint: 'Edit item'});
    expect(withHint.findComponent(TooltipStub).text()).toBe('Edit item');

    const withoutHint = createWrapper({icon: 'mdi-pencil'});
    expect(withoutHint.findComponent(TooltipStub).exists()).toBe(false);
  });
});
