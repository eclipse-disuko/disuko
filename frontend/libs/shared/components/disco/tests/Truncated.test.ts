// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {mount} from '@vue/test-utils';
import {describe, expect, it} from 'vitest';
import Truncated from '../Truncated.vue';

const TooltipStub = {
  template: '<div class="tooltip-stub"><slot /></div>',
};

describe('Truncated', () => {
  const createWrapper = (slots = {}) => {
    return mount(Truncated, {
      slots,
      global: {
        stubs: {Tooltip: TooltipStub},
      },
    });
  };

  it('renders without throwing', () => {
    const wrapper = createWrapper({default: 'Some long text'});

    expect(wrapper.exists()).toBe(true);
  });

  it('renders the slot content in the truncated span', () => {
    const wrapper = createWrapper({default: 'Some long text'});

    expect(wrapper.find('.truncate').text()).toBe('Some long text');
  });

  it('also forwards the slot content into the Tooltip', () => {
    const wrapper = createWrapper({default: 'Some long text'});

    expect(wrapper.findComponent(TooltipStub).text()).toBe('Some long text');
  });
});
