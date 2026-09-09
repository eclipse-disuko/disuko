// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import INavItem from '@shared/model/INavItem';
import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, describe, expect, it} from 'vitest';
import SubNavItem from '../SubNavItem.vue';

const VTooltipStub = {
  template:
    '<div class="v-tooltip"><slot name="activator" :props="{}" /><div class="tooltip-text"><slot /></div></div>',
};

const VListItemStub = {
  props: ['to', 'active', 'title'],
  template: `<div class="v-list-item" :to="to" :active="active">
    <slot name="prepend" />
    <slot name="title" />
    <slot />
  </div>`,
};

describe('SubNavItem', () => {
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

  const baseItem: INavItem = {
    title: 'Sub Item',
    path: '/dashboard/sub',
    iconName: 'mdi-file',
    condition: true,
    active: false,
    subItems: [],
  };

  const createWrapper = (item: INavItem) => {
    return mount(SubNavItem, {
      props: {item},
      global: {stubs: {'v-tooltip': VTooltipStub, 'v-list-item': VListItemStub}},
    });
  };

  it('renders a v-list-item using item.path', () => {
    const wrapper = createWrapper(baseItem);

    expect(wrapper.find('.v-list-item').exists()).toBe(true);
    expect(wrapper.text()).toContain('Sub Item');
  });

  it('renders the item icon', () => {
    const wrapper = createWrapper(baseItem);

    expect(wrapper.text()).toContain('mdi-file');
  });

  it('renders the translated title as the tooltip text', () => {
    const wrapper = createWrapper(baseItem);

    expect(wrapper.find('.tooltip-text').text()).toBe('Sub Item');
  });

  it('supports mouseover/mouseleave without throwing', async () => {
    const wrapper = createWrapper(baseItem);

    await wrapper.find('.v-list-item').trigger('mouseover');
    await wrapper.find('.v-list-item').trigger('mouseleave');

    expect(wrapper.exists()).toBe(true);
  });
});
