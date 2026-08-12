// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, describe, expect, it} from 'vitest';
import DDateCellWithTooltip from '../DDateCellWithTooltip.vue';

const TooltipStub = {
  template: '<div class="tooltip-stub">{{ text }}</div>',
  props: ['text'],
};

describe('DDateCellWithTooltip', () => {
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
    return mount(DDateCellWithTooltip, {
      props,
      global: {
        stubs: {Tooltip: TooltipStub},
      },
    });
  };

  it('renders without throwing', () => {
    const wrapper = createWrapper({value: '2024-05-01T00:00:00Z'});

    expect(wrapper.exists()).toBe(true);
  });

  it('renders nothing when value is the zero-date sentinel', () => {
    const wrapper = createWrapper({value: '0001-01-01T00:00:00Z'});

    expect(wrapper.find('span').exists()).toBe(false);
  });

  it('renders nothing when value is empty', () => {
    const wrapper = createWrapper({value: ''});

    expect(wrapper.find('span').exists()).toBe(false);
  });

  it('shows the formatted date (not date+time) by default', () => {
    const wrapper = createWrapper({value: '2024-05-01T12:34:00Z'});

    expect(wrapper.find('span > span').text()).not.toContain(':');
  });

  it('shows the formatted date and time when showTime is true', () => {
    const wrapper = createWrapper({value: '2024-05-01T12:34:00Z', showTime: true});

    expect(wrapper.findComponent(TooltipStub).props('text')).toContain('2024');
  });
});
