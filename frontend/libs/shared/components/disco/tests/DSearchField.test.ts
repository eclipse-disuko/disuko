// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {vuetifyStubs} from '@disclosure-portal/test-utils/vuetify-stubs';
import {config, mount} from '@vue/test-utils';
import {afterAll, afterEach, beforeAll, describe, expect, it, vi} from 'vitest';
import {nextTick} from 'vue';
import DSearchField from '../DSearchField.vue';

// DSearchField calls inputRef.value.focus() on the mounted component instance, so the stub
// needs to expose a focus() method that focuses its root <input>, mirroring real Vuetify's API.
const vTextFieldStub = {
  ...vuetifyStubs['v-text-field'],
  methods: {
    focus() {
      (this as unknown as {$el: HTMLElement}).$el.focus();
    },
  },
};

describe('DSearchField', () => {
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
    return mount(DSearchField, {
      props: {modelValue: '', ...props},
      attachTo: document.body,
      global: {
        stubs: {...vuetifyStubs, 'v-text-field': vTextFieldStub},
      },
    });
  };

  afterEach(() => {
    vi.useRealTimers();
  });

  it('renders the collapsed search button when modelValue is empty', () => {
    const wrapper = createWrapper();

    expect(wrapper.find('button').exists()).toBe(true);
    expect(wrapper.find('input').exists()).toBe(false);
    expect(wrapper.text()).toContain('Search');
  });

  it('renders the expanded field when modelValue is pre-filled', () => {
    const wrapper = createWrapper({modelValue: 'existing'});

    expect(wrapper.find('input').exists()).toBe(true);
    expect(wrapper.find('button').exists()).toBe(false);
  });

  it('expands and focuses the input when the button is clicked', async () => {
    vi.useFakeTimers();
    const wrapper = createWrapper();

    await wrapper.find('button').trigger('click');
    await nextTick();
    vi.advanceTimersByTime(150);

    expect(wrapper.find('input').exists()).toBe(true);
    expect(document.activeElement).toBe(wrapper.find('input').element);

    wrapper.unmount();
  });

  it('emits update:modelValue when typing in the field', async () => {
    const wrapper = createWrapper({modelValue: 'a'});

    await wrapper.find('input').setValue('abc');

    expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['abc']);
  });

  it('collapses on Escape when the field is empty', async () => {
    const wrapper = createWrapper({modelValue: ''});
    await wrapper.find('button').trigger('click');
    await nextTick();

    await wrapper.find('input').trigger('keydown', {key: 'Escape'});
    await nextTick();

    expect(wrapper.find('button').exists()).toBe(true);
  });

  it('does not collapse on Escape when the field has a value', async () => {
    const wrapper = createWrapper({modelValue: 'text'});

    await wrapper.find('input').trigger('keydown', {key: 'Escape'});
    await nextTick();

    expect(wrapper.find('input').exists()).toBe(true);
  });

  it('collapses on blur when the field is empty', async () => {
    const wrapper = createWrapper({modelValue: ''});
    await wrapper.find('button').trigger('click');
    await nextTick();

    await wrapper.find('input').trigger('blur');
    await nextTick();

    expect(wrapper.find('button').exists()).toBe(true);
  });

  it('does not collapse on blur when the field has a value', async () => {
    const wrapper = createWrapper({modelValue: 'text'});

    await wrapper.find('input').trigger('blur');
    await nextTick();

    expect(wrapper.find('input').exists()).toBe(true);
  });

  it('forwards the disabled prop to the button and text field', () => {
    const collapsed = createWrapper({disabled: true});
    expect(collapsed.find('button').attributes('disabled')).toBeDefined();

    const expanded = createWrapper({modelValue: 'x', disabled: true});
    const textField = expanded.findComponent(vTextFieldStub);
    expect(textField.props('disabled')).toBe(true);
  });
});
