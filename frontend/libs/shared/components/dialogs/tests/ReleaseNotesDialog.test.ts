// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, describe, expect, it} from 'vitest';
import ReleaseNotesDialog from '../ReleaseNotesDialog.vue';

const VDialogStub = {
  props: ['modelValue', 'width', 'maxWidth', 'scrim', 'scrollable'],
  emits: ['update:modelValue'],
  template: '<div v-if="modelValue" class="v-dialog"><slot /></div>',
};

const DCloseButtonStub = {
  template: '<button class="d-close-button" type="button" />',
};

const VBtnStub = {
  props: ['text', 'size', 'color'],
  template: '<button type="button">{{ text }}</button>',
};

const MarkdownStub = {
  props: ['text', 'id'],
  template: '<div class="markdown" :id="id">{{ text }}</div>',
};

describe('ReleaseNotesDialog', () => {
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

  const createWrapper = (slots: Record<string, string> = {}) => {
    return mount(ReleaseNotesDialog, {
      props: {releaseNotes: '# What changed'},
      slots,
      global: {
        stubs: {'v-dialog': VDialogStub, 'v-btn': VBtnStub, DCloseButton: DCloseButtonStub, Markdown: MarkdownStub},
      },
    });
  };

  it('renders the default trigger button when no slot is provided', () => {
    const wrapper = createWrapper();

    expect(wrapper.find('button').text()).toBe('Replace me');
  });

  it('does not render the dialog before it is opened', () => {
    const wrapper = createWrapper();

    expect(wrapper.find('.v-dialog').exists()).toBe(false);
  });

  it('opens the dialog when the default trigger button is clicked', async () => {
    const wrapper = createWrapper();

    await wrapper.find('button').trigger('click');

    expect(wrapper.find('.v-dialog').exists()).toBe(true);
  });

  it('opens the dialog via a custom slot invoking showDialog', async () => {
    const wrapper = createWrapper({
      default:
        '<template #default="{ showDialog }"><button class="custom-open" @click="showDialog">open</button></template>',
    });

    await wrapper.find('.custom-open').trigger('click');

    expect(wrapper.find('.v-dialog').exists()).toBe(true);
  });

  it('renders the releaseNotes content via Markdown once opened', async () => {
    const wrapper = createWrapper();
    await wrapper.find('button').trigger('click');

    expect(wrapper.findComponent(MarkdownStub).props('text')).toBe('# What changed');
  });

  it('closes the dialog when the close button is clicked', async () => {
    const wrapper = createWrapper();
    await wrapper.find('button').trigger('click');
    expect(wrapper.find('.v-dialog').exists()).toBe(true);

    await wrapper.find('.d-close-button').trigger('click');

    expect(wrapper.find('.v-dialog').exists()).toBe(false);
  });
});
