// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, describe, expect, it, vi} from 'vitest';
import DCopyClipboardButton from '../DCopyClipboardButton.vue';

const copyToClipboard = vi.hoisted(() => vi.fn());

vi.mock('@shared/utils/clipboard', () => ({
  useClipboard: () => ({copyToClipboard}),
}));

const DIconButtonStub = {
  template: '<button class="d-icon-button" :data-icon="icon" :data-hint="hint" @click="$emit(\'clicked\')"></button>',
  props: ['icon', 'hint'],
  emits: ['clicked'],
};

describe('DCopyClipboardButton', () => {
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
    return mount(DCopyClipboardButton, {
      props,
      global: {
        stubs: {DIconButton: DIconButtonStub},
      },
    });
  };

  it('renders without throwing', () => {
    const wrapper = createWrapper({hint: 'Copy', content: 'some content'});

    expect(wrapper.findComponent(DIconButtonStub).exists()).toBe(true);
  });

  it('forwards the hint and defaults to the mdi-content-copy icon', () => {
    const wrapper = createWrapper({hint: 'Copy me', content: 'value'});
    const btn = wrapper.findComponent(DIconButtonStub);

    expect(btn.props('hint')).toBe('Copy me');
    expect(btn.props('icon')).toBe('mdi-content-copy');
  });

  it('forwards a custom icon when given', () => {
    const wrapper = createWrapper({icon: 'mdi-clipboard', hint: 'Copy', content: 'value'});

    expect(wrapper.findComponent(DIconButtonStub).props('icon')).toBe('mdi-clipboard');
  });

  it('copies the content to the clipboard when the button is clicked', async () => {
    copyToClipboard.mockClear();
    const wrapper = createWrapper({hint: 'Copy', content: 'text to copy'});

    await wrapper.findComponent(DIconButtonStub).trigger('click');

    expect(copyToClipboard).toHaveBeenCalledWith('text to copy');
  });
});
