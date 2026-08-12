// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, beforeEach, describe, expect, it, vi} from 'vitest';
import ExternalLinkDialog from '../ExternalLinkDialog.vue';

const {openUrlInNewTabMock} = vi.hoisted(() => ({openUrlInNewTabMock: vi.fn()}));

vi.mock('@shared/utils/url', () => ({
  openUrlInNewTab: openUrlInNewTabMock,
}));

const VDialogStub = {
  props: ['modelValue', 'width'],
  emits: ['update:modelValue'],
  template: '<div v-if="modelValue" class="v-dialog"><slot /></div>',
};

const DialogLayoutStub = {
  props: ['config'],
  emits: ['close', 'secondaryAction', 'primaryAction'],
  template: `
    <div class="dialog-layout">
      <h4>{{ config.title }}</h4>
      <slot />
      <button class="primary-btn" type="button" @click="$emit('primaryAction')">{{ config.primaryButton?.text }}</button>
      <button class="secondary-btn" type="button" @click="$emit('secondaryAction')">{{ config.secondaryButton?.text }}</button>
    </div>`,
};

describe('ExternalLinkDialog', () => {
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

  beforeEach(() => {
    vi.clearAllMocks();
  });

  const createWrapper = (props: Record<string, unknown> = {}) => {
    return mount(ExternalLinkDialog, {
      props: {url: 'https://example.com', isDialogVisible: true, ...props},
      global: {stubs: {'v-dialog': VDialogStub, DialogLayout: DialogLayoutStub}},
    });
  };

  it('renders the dialog with the given url when visible', () => {
    const wrapper = createWrapper();

    expect(wrapper.find('.v-dialog').exists()).toBe(true);
    expect(wrapper.text()).toContain('https://example.com');
  });

  it('does not render the dialog when isDialogVisible is false', () => {
    const wrapper = createWrapper({isDialogVisible: false});

    expect(wrapper.find('.v-dialog').exists()).toBe(false);
  });

  it('opens the url in a new tab, closes the dialog, and emits close when the primary button is clicked', async () => {
    const wrapper = createWrapper();

    await wrapper.find('.primary-btn').trigger('click');

    expect(openUrlInNewTabMock).toHaveBeenCalledWith('https://example.com');
    expect(wrapper.emitted('update:isDialogVisible')?.[0]).toEqual([false]);
    expect(wrapper.emitted('close')).toHaveLength(1);
  });

  it('closes the dialog and emits close without opening the url when secondary is clicked', async () => {
    const wrapper = createWrapper();

    await wrapper.find('.secondary-btn').trigger('click');

    expect(openUrlInNewTabMock).not.toHaveBeenCalled();
    expect(wrapper.emitted('update:isDialogVisible')?.[0]).toEqual([false]);
    expect(wrapper.emitted('close')).toHaveLength(1);
  });
});
