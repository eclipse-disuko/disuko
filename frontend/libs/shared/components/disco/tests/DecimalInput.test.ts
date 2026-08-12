// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {vuetifyStubs} from '@disclosure-portal/test-utils/vuetify-stubs';
import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, describe, expect, it} from 'vitest';
import DecimalInput from '../DecimalInput.vue';

describe('DecimalInput', () => {
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
    return mount(DecimalInput, {
      props: {modelValue: '0', ...props},
      global: {
        stubs: vuetifyStubs,
      },
    });
  };

  it('renders without throwing', () => {
    const wrapper = createWrapper();

    expect(wrapper.find('input').exists()).toBe(true);
  });

  it('formats the modelValue with thousand separators and two decimals', () => {
    const wrapper = createWrapper({modelValue: '1234.5'});

    expect(wrapper.find('input').element.value).toBe('1.234,50');
  });

  it('parses and emits a normalized value on blur', async () => {
    const wrapper = createWrapper();
    const input = wrapper.find('input');

    await input.setValue('1.234,56');
    await input.trigger('blur');

    expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['1234.56']);
  });

  it('defaults to 0.00 on blur when the field was cleared', async () => {
    const wrapper = createWrapper();
    const input = wrapper.find('input');

    await input.setValue('');
    await input.trigger('blur');

    expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['0.00']);
  });

  it('resets the value and emits 0.00 when reset becomes true', async () => {
    const wrapper = createWrapper({modelValue: '10'});

    await wrapper.setProps({reset: true});

    expect(wrapper.emitted('update:modelValue')?.at(-1)).toEqual(['0.00']);
  });

  it('passes the readonly prop through to the underlying field', () => {
    const wrapper = createWrapper({readonly: true});

    expect(wrapper.findComponent(vuetifyStubs['v-text-field']).props('readonly')).toBe(true);
  });
});
