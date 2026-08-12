// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {vuetifyStubs} from '@disclosure-portal/test-utils/vuetify-stubs';
import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, describe, expect, it} from 'vitest';
import DExternalLink from '../DExternalLink.vue';

const TooltipStub = {
  template: '<div class="tooltip-stub">{{ text }}</div>',
  props: ['text'],
};
const ExternalLinkDialogStub = {
  template: '<div class="external-link-dialog-stub" :data-visible="isDialogVisible"></div>',
  props: ['url', 'isDialogVisible'],
};

describe('DExternalLink', () => {
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
    return mount(DExternalLink, {
      props,
      global: {
        stubs: {...vuetifyStubs, Tooltip: TooltipStub, ExternalLinkDialog: ExternalLinkDialogStub},
      },
    });
  };

  it('renders without throwing', () => {
    const wrapper = createWrapper({url: 'https://example.com', text: 'Example'});

    expect(wrapper.exists()).toBe(true);
  });

  it('renders the link with href and text', () => {
    const wrapper = createWrapper({url: 'https://example.com', text: 'Example'});
    const link = wrapper.find('a');

    expect(link.attributes('href')).toBe('https://example.com');
    expect(link.text()).toContain('Example');
  });

  it('renders no link when url is empty', () => {
    const wrapper = createWrapper({url: '', text: 'Example'});

    expect(wrapper.find('a').exists()).toBe(false);
  });

  it('shows the confirmation dialog on click for an http(s) url and does not for a relative url', async () => {
    const external = createWrapper({url: 'https://example.com', text: 'Example'});
    await external.find('a').trigger('click');
    expect(external.findComponent(ExternalLinkDialogStub).props('isDialogVisible')).toBe(true);

    const internal = createWrapper({url: '/internal/path', text: 'Example'});
    await internal.find('a').trigger('click');
    expect(internal.findComponent(ExternalLinkDialogStub).props('isDialogVisible')).toBe(false);
  });

  it('renders the tooltip only when the tooltip prop is true', () => {
    const withTooltip = createWrapper({url: 'https://example.com', text: 'Example', tooltip: true});
    expect(withTooltip.findComponent(TooltipStub).exists()).toBe(true);

    const withoutTooltip = createWrapper({url: 'https://example.com', text: 'Example'});
    expect(withoutTooltip.findComponent(TooltipStub).exists()).toBe(false);
  });
});
