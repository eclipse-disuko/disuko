// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {vuetifyStubs} from '@disclosure-portal/test-utils/vuetify-stubs';
import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, describe, expect, it} from 'vitest';
import DExternalLinkIcon from '../DExternalLinkIcon.vue';

const TooltipStub = {
  template: '<div class="tooltip-stub">{{ text }}</div>',
  props: ['text'],
};
const ExternalLinkDialogStub = {
  template: '<div class="external-link-dialog-stub" :data-url="url" :data-visible="isDialogVisible"></div>',
  props: ['url', 'isDialogVisible'],
};

describe('DExternalLinkIcon', () => {
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
    return mount(DExternalLinkIcon, {
      props,
      global: {
        stubs: {...vuetifyStubs, Tooltip: TooltipStub, ExternalLinkDialog: ExternalLinkDialogStub},
      },
    });
  };

  it('renders without throwing', () => {
    const wrapper = createWrapper({url: 'https://example.com', hint: 'Open'});

    expect(wrapper.exists()).toBe(true);
  });

  it('renders the link with the given href when url is set', () => {
    const wrapper = createWrapper({url: 'https://example.com', hint: 'Open'});
    const link = wrapper.find('a');

    expect(link.exists()).toBe(true);
    expect(link.attributes('href')).toBe('https://example.com');
  });

  it('renders no link when url is empty', () => {
    const wrapper = createWrapper({url: '', hint: 'Open'});

    expect(wrapper.find('a').exists()).toBe(false);
  });

  it('renders the hint tooltip when hint is given', () => {
    const wrapper = createWrapper({url: 'https://example.com', hint: 'Open in new tab'});

    expect(wrapper.findComponent(TooltipStub).text()).toBe('Open in new tab');
  });

  it('opens the confirmation dialog instead of navigating when clicked', async () => {
    const wrapper = createWrapper({url: 'https://example.com', hint: 'Open'});

    expect(wrapper.findComponent(ExternalLinkDialogStub).props('isDialogVisible')).toBe(false);

    await wrapper.find('a').trigger('click');

    expect(wrapper.findComponent(ExternalLinkDialogStub).props('isDialogVisible')).toBe(true);
  });
});
