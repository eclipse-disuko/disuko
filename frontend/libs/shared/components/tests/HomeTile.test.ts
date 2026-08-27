// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, describe, expect, it} from 'vitest';
import HomeTile from '../HomeTile.vue';

describe('HomeTile', () => {
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

  const stubs = {
    'v-card': {template: '<div class="v-card"><slot /></div>'},
    'v-card-item': {template: '<div class="v-card-item"><slot /></div>'},
    'v-card-text': {template: '<div class="v-card-text"><slot /></div>'},
    'v-card-actions': {template: '<div class="v-card-actions"><slot /></div>'},
    'v-divider': {template: '<hr />'},
    'v-expand-transition': {template: '<div><slot /></div>'},
    'v-icon': {template: '<i><slot /></i>', props: ['icon', 'size', 'color']},
    'v-btn': {
      template: '<button type="button"><slot /></button>',
      props: ['href', 'size', 'variant', 'color'],
    },
    'v-badge': {template: '<span class="v-badge" :data-content="content"></span>', props: ['content', 'color']},
    Stack: {template: '<div class="stack"><slot /></div>'},
  };

  const createWrapper = (props: Record<string, unknown>) => {
    return mount(HomeTile, {
      props,
      global: {stubs},
    });
  };

  it('renders an action tile with title and icon', () => {
    const wrapper = createWrapper({type: 'action', title: 'My Action', icon: 'mdi-star'});

    expect(wrapper.text()).toContain('My Action');
    expect(wrapper.find('.v-card').exists()).toBe(true);
  });

  it('emits click when an action tile is clicked', async () => {
    const wrapper = createWrapper({type: 'action', title: 'My Action', icon: 'mdi-star'});

    await wrapper.find('.v-card').trigger('click');

    expect(wrapper.emitted('click')).toHaveLength(1);
  });

  it('does not show the more/less toggle when there is no description', () => {
    const wrapper = createWrapper({type: 'action', title: 'My Action', icon: 'mdi-star'});

    expect(wrapper.find('.v-card-actions').exists()).toBe(false);
  });

  it('toggles the description visibility when more/less is clicked', async () => {
    const wrapper = createWrapper({
      type: 'action',
      title: 'My Action',
      icon: 'mdi-star',
      description: 'Some details',
    });

    expect(wrapper.text()).not.toContain('Some details');

    await wrapper.find('.v-card-actions button').trigger('click');
    expect(wrapper.text()).toContain('Some details');
    expect(wrapper.text()).toContain('Less');

    await wrapper.find('.v-card-actions button').trigger('click');
    expect(wrapper.text()).not.toContain('Some details');
  });

  it('does not show the more/less toggle when showExpand is false', () => {
    const wrapper = createWrapper({
      type: 'action',
      title: 'My Action',
      icon: 'mdi-star',
      description: 'Some details',
      showExpand: false,
    });

    expect(wrapper.find('.v-card-actions').exists()).toBe(false);
  });

  it('renders a navigation tile with title and icon', () => {
    const wrapper = createWrapper({type: 'navigation', title: 'Go Somewhere', icon: 'mdi-arrow-right'});

    expect(wrapper.find('button').exists()).toBe(true);
    expect(wrapper.text()).toContain('Go Somewhere');
  });

  it('emits click when a navigation tile is clicked', async () => {
    const wrapper = createWrapper({type: 'navigation', title: 'Go Somewhere', icon: 'mdi-arrow-right'});

    await wrapper.find('button').trigger('click');

    expect(wrapper.emitted('click')).toHaveLength(1);
  });

  it('shows the badge when showBadge and badgeContent are set', () => {
    const wrapper = createWrapper({
      type: 'navigation',
      title: 'Go Somewhere',
      icon: 'mdi-arrow-right',
      showBadge: true,
      badgeContent: 5,
    });

    const badge = wrapper.find('.v-badge');
    expect(badge.exists()).toBe(true);
    expect(badge.attributes('data-content')).toBe('5');
  });

  it('does not show the badge when showBadge is false', () => {
    const wrapper = createWrapper({
      type: 'navigation',
      title: 'Go Somewhere',
      icon: 'mdi-arrow-right',
      showBadge: false,
      badgeContent: 5,
    });

    expect(wrapper.find('.v-badge').exists()).toBe(false);
  });

  it('does not show the badge when badgeContent is missing', () => {
    const wrapper = createWrapper({
      type: 'navigation',
      title: 'Go Somewhere',
      icon: 'mdi-arrow-right',
      showBadge: true,
    });

    expect(wrapper.find('.v-badge').exists()).toBe(false);
  });
});
