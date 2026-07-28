// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {mount} from '@vue/test-utils';
import {describe, expect, it} from 'vitest';
import {defineComponent, h, nextTick} from 'vue';
import Tooltip from '../Tooltip.vue';

const VTooltipStub = defineComponent({
  name: 'VTooltip',
  setup: (_props, {slots}) => () => h('div', {class: 'v-tooltip-stub'}, slots.default?.()),
});

const mountInParent = () =>
  mount(
    defineComponent({
      components: {Tooltip},
      template: '<div class="parent"><Tooltip text="hint" /></div>',
    }),
    {
      attachTo: document.body,
      global: {
        components: {VTooltip: VTooltipStub},
      },
    },
  );

describe('Tooltip', () => {
  it('does not mount the overlay before the parent is hovered', () => {
    const wrapper = mountInParent();

    expect(wrapper.find('.v-tooltip-stub').exists()).toBe(false);

    wrapper.unmount();
  });

  it('mounts the overlay on parent hover', async () => {
    const wrapper = mountInParent();

    await wrapper.find('.parent').trigger('pointerenter');
    await nextTick();

    expect(wrapper.find('.v-tooltip-stub').exists()).toBe(true);

    wrapper.unmount();
  });

  it('mounts the overlay eagerly for an external activator', () => {
    const wrapper = mount(Tooltip, {
      props: {activator: '#target', text: 'hint'},
      global: {
        components: {VTooltip: VTooltipStub},
      },
    });

    expect(wrapper.find('.v-tooltip-stub').exists()).toBe(true);

    wrapper.unmount();
  });
});
