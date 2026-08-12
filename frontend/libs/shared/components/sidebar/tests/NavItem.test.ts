// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import INavItem from '@disclosure-portal/model/INavItem';
import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, describe, expect, it} from 'vitest';
import NavItem from '../NavItem.vue';

const VTooltipStub = {
  template: '<div class="v-tooltip"><slot name="activator" :props="{}" /><div class="tooltip-text"><slot /></div></div>',
};

const VListItemStub = {
  props: ['to', 'href', 'active', 'title'],
  template: `<div class="v-list-item" :href="href" :to="to" :active="active">
    <slot name="prepend" />
    <slot name="title" />
    <slot />
  </div>`,
};

describe('NavItem', () => {
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
    title: 'Dashboard',
    path: '/dashboard',
    iconName: 'mdi-view-dashboard',
    condition: true,
    active: false,
    subItems: [],
  };

  const createWrapper = (item: INavItem) => {
    return mount(NavItem, {
      props: {item},
      global: {stubs: {'v-tooltip': VTooltipStub, 'v-list-item': VListItemStub}},
    });
  };

  it('renders an internal v-list-item using item.path when there is no externalPath', () => {
    const wrapper = createWrapper(baseItem);

    const listItem = wrapper.find('.v-list-item');
    expect(listItem.exists()).toBe(true);
    expect(listItem.attributes('href')).toBeUndefined();
  });

  it('renders an external v-list-item using item.externalPath when set', () => {
    const wrapper = createWrapper({...baseItem, externalPath: 'https://example.com'});

    const listItem = wrapper.find('.v-list-item');
    expect(listItem.attributes('href')).toBe('https://example.com');
  });

  it('shows the outline icon when not hovered and not active', () => {
    const wrapper = createWrapper(baseItem);

    expect(wrapper.text()).toContain('mdi-view-dashboard-outline');
  });

  it('shows the solid icon when active', () => {
    const wrapper = createWrapper({...baseItem, active: true});

    expect(wrapper.text()).toContain('mdi-view-dashboard');
    expect(wrapper.text()).not.toContain('mdi-view-dashboard-outline');
  });

  it('shows the solid icon while hovered', async () => {
    const wrapper = createWrapper(baseItem);

    await wrapper.find('.v-list-item').trigger('mouseover');

    expect(wrapper.text()).toContain('mdi-view-dashboard');
    expect(wrapper.text()).not.toContain('mdi-view-dashboard-outline');

    await wrapper.find('.v-list-item').trigger('mouseleave');

    expect(wrapper.text()).toContain('mdi-view-dashboard-outline');
  });

  it('renders the translated tooltip text when tooltip is set', () => {
    const wrapper = createWrapper({...baseItem, tooltip: 'TT_DASHBOARD'});

    expect(wrapper.find('.tooltip-text').text()).toBe('TT_DASHBOARD');
  });

  it('renders an empty tooltip when tooltip is not set', () => {
    const wrapper = createWrapper(baseItem);

    expect(wrapper.find('.tooltip-text').text()).toBe('');
  });
});
