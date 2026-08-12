// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, describe, expect, it} from 'vitest';
import GridHeaderFilterIcon from '../GridHeaderFilterIcon.vue';

const GridHeaderMenuStub = {
  props: ['showReset', 'resetHint', 'cardTitle', 'allItems', 'selectedItems', 'selectLabel'],
  emits: ['update', 'reset'],
  template: `
    <div class="grid-header-menu">
      <slot name="activator" :props="{}" />
      <button class="update-btn" type="button" @click="$emit('update', ['a'])">update</button>
      <button class="reset-btn" type="button" @click="$emit('reset')">reset</button>
    </div>`,
};

const TooltipStub = {template: '<span class="tooltip"><slot /></span>'};
const VIconStub = {
  props: ['color'],
  template: '<i class="v-icon" :color="color"><slot /></i>',
};

describe('GridHeaderFilterIcon', () => {
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

  const allItems = [
    {value: 'a', text: 'A'},
    {value: 'b', text: 'B'},
  ];

  const createWrapper = (props: Record<string, unknown> = {}) => {
    return mount(GridHeaderFilterIcon, {
      props: {column: {value: 'a'}, label: 'My Column', allItems, ...props},
      global: {
        stubs: {GridHeaderMenu: GridHeaderMenuStub, Tooltip: TooltipStub, 'v-icon': VIconStub},
      },
    });
  };

  it('passes allItems and label through to GridHeaderMenu', () => {
    const wrapper = createWrapper();

    const menu = wrapper.findComponent(GridHeaderMenuStub);
    expect(menu.props('allItems')).toEqual(allItems);
    expect(menu.props('selectLabel')).toBe('My Column');
  });

  it('initializes selectedFilters from initialSelected', () => {
    const wrapper = createWrapper({initialSelected: ['b']});

    const menu = wrapper.findComponent(GridHeaderMenuStub);
    expect(menu.props('selectedItems')).toEqual(['b']);
    expect(menu.props('showReset')).toBe(true);
  });

  it('does not show reset when there is no selection', () => {
    const wrapper = createWrapper();

    expect(wrapper.findComponent(GridHeaderMenuStub).props('showReset')).toBe(false);
  });

  it('updates selectedFilters and the v-model when GridHeaderMenu emits update', async () => {
    const wrapper = createWrapper();

    await wrapper.find('.update-btn').trigger('click');

    expect(wrapper.findComponent(GridHeaderMenuStub).props('selectedItems')).toEqual(['a']);
    expect(wrapper.emitted('update:modelValue')?.[0]).toEqual([['a']]);
  });

  it('resets selectedFilters back to initialSelected when GridHeaderMenu emits reset', async () => {
    const wrapper = createWrapper({initialSelected: ['b']});

    await wrapper.find('.update-btn').trigger('click');
    expect(wrapper.findComponent(GridHeaderMenuStub).props('selectedItems')).toEqual(['a']);

    await wrapper.find('.reset-btn').trigger('click');

    expect(wrapper.findComponent(GridHeaderMenuStub).props('selectedItems')).toEqual(['b']);
  });

  it('highlights the filter icon as primary once a filter is selected', async () => {
    const wrapper = createWrapper();

    expect(wrapper.find('.v-icon').attributes('color')).toBe('default');

    await wrapper.find('.update-btn').trigger('click');

    expect(wrapper.find('.v-icon').attributes('color')).toBe('primary');
  });
});
