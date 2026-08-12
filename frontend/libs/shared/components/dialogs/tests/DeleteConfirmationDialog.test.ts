// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {vuetifyStubs} from '@disclosure-portal/test-utils/vuetify-stubs';
import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, describe, expect, it} from 'vitest';
import DeleteConfirmationDialog from '../DeleteConfirmationDialog.vue';

const DCloseButtonStub = {
  template: '<button class="d-close-button" type="button" />',
};

const DCActionButtonStub = {
  props: ['text', 'isDialogButton', 'variant'],
  template: '<button type="button">{{ text }}</button>',
};

describe('DeleteConfirmationDialog', () => {
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

  const createWrapper = (props: Record<string, unknown> = {}) => {
    return mount(DeleteConfirmationDialog, {
      props,
      slots: {
        default: '<template #default="{ showDialog }"><button class="open-btn" @click="showDialog">open</button></template>',
      },
      global: {
        stubs: {...vuetifyStubs, DCloseButton: DCloseButtonStub, DCActionButton: DCActionButtonStub},
      },
    });
  };

  it('does not render the dialog before it is opened', () => {
    const wrapper = createWrapper();

    expect(wrapper.find('.v-dialog').exists()).toBe(false);
  });

  it('opens the dialog when the exposed showDialog is invoked via the default slot', async () => {
    const wrapper = createWrapper();

    await wrapper.find('.open-btn').trigger('click');

    expect(wrapper.find('.v-dialog').exists()).toBe(true);
  });

  it('shows the title prop text when provided', async () => {
    const wrapper = createWrapper({title: 'Project X'});
    await wrapper.find('.open-btn').trigger('click');

    expect(wrapper.text()).toContain('Project X');
  });

  it('shows the message prop text when provided', async () => {
    const wrapper = createWrapper({message: 'This cannot be undone.'});
    await wrapper.find('.open-btn').trigger('click');

    expect(wrapper.text()).toContain('This cannot be undone.');
  });

  it('closes without emitting confirmed when cancel is clicked', async () => {
    const wrapper = createWrapper();
    await wrapper.find('.open-btn').trigger('click');

    await wrapper.find('.d-close-button').trigger('click');

    expect(wrapper.find('.v-dialog').exists()).toBe(false);
    expect(wrapper.emitted('confirmed')).toBeUndefined();
  });

  it('uses the default BTN_DELETE label when buttonText is not provided', async () => {
    const wrapper = createWrapper();
    await wrapper.find('.open-btn').trigger('click');

    expect(wrapper.text()).toContain('Delete');
  });

  it('uses the custom buttonText label when provided', async () => {
    const wrapper = createWrapper({buttonText: 'Remove forever'});
    await wrapper.find('.open-btn').trigger('click');

    expect(wrapper.text()).toContain('Remove forever');
  });

  it('emits confirmed and closes the dialog when the confirm button is clicked', async () => {
    const wrapper = createWrapper();
    await wrapper.find('.open-btn').trigger('click');

    const buttons = wrapper.findAllComponents(DCActionButtonStub);
    await buttons[1].trigger('click');

    expect(wrapper.emitted('confirmed')).toHaveLength(1);
    expect(wrapper.find('.v-dialog').exists()).toBe(false);
  });
});
