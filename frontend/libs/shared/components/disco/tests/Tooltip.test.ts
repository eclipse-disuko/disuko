// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {vuetifyStubs} from '@disclosure-portal/test-utils/vuetify-stubs';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, describe, expect, it} from 'vitest';
import Tooltip from '../Tooltip.vue';

const DExternalLinkStub = {
  template: '<a class="d-external-link" :href="url">{{ text }}</a>',
  props: ['url', 'text', 'tooltip'],
};

describe('Tooltip', () => {
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

  const createWrapper = (props = {}, slots = {}) => {
    return mount(Tooltip, {
      props,
      slots,
      global: {
        stubs: {...vuetifyStubs, DExternalLink: DExternalLinkStub},
      },
    });
  };

  it('renders the default slot when no text is given', () => {
    const wrapper = createWrapper({}, {default: '<span class="activator">hover me</span>'});

    expect(wrapper.find('.activator').exists()).toBe(true);
  });

  it('renders plain text with no links as a span', () => {
    const wrapper = createWrapper({text: 'Just a plain hint'});

    expect(wrapper.text()).toContain('Just a plain hint');
    expect(wrapper.findComponent(DExternalLinkStub).exists()).toBe(false);
  });

  it('splits text containing a URL into text and link segments', () => {
    const wrapper = createWrapper({text: 'Check https://example.com for details'});

    const link = wrapper.findComponent(DExternalLinkStub);
    expect(link.exists()).toBe(true);
    expect(link.props('url')).toBe('https://example.com');
    expect(wrapper.text()).toContain('Check');
    expect(wrapper.text()).toContain('for details');
  });

  it('forwards the tooltip open delay constant', () => {
    const wrapper = createWrapper({text: 'hint'});

    expect(wrapper.attributes('open-delay')).toBe(String(TOOLTIP_OPEN_DELAY_IN_MS));
  });

  it('forwards the disabled prop', () => {
    const wrapper = createWrapper({text: 'hint', disabled: true});

    expect(wrapper.attributes('disabled')).toBe('true');
  });

  it('renders the activator slot with the default slot content when asParent is true', () => {
    const wrapper = createWrapper({text: 'wrapped hint', asParent: true}, {default: '<button>trigger</button>'});

    expect(wrapper.find('button').exists()).toBe(true);
    expect(wrapper.text()).toContain('wrapped hint');
  });

  it('splits links correctly when asParent is true', () => {
    const wrapper = createWrapper(
      {text: 'See https://example.com now', asParent: true},
      {default: '<button>trigger</button>'},
    );

    const link = wrapper.findComponent(DExternalLinkStub);
    expect(link.exists()).toBe(true);
    expect(link.props('url')).toBe('https://example.com');
  });
});
