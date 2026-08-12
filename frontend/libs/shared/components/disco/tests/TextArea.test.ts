// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {vuetifyStubs} from '@disclosure-portal/test-utils/vuetify-stubs';
import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, describe, expect, it} from 'vitest';
import TextArea from '../TextArea.vue';

const TooltipStub = {
  template: '<div class="tooltip-stub"><slot />{{ text }}</div>',
  props: ['text', 'asParent'],
};

// vuetifyStubs' generic v-textarea stub only renders no slot at all; TextArea.vue uses the
// named #append-inner slot for the help Tooltip, so render it here too.
const vTextareaStub = {
  ...vuetifyStubs['v-textarea'],
  template:
    '<textarea :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)"></textarea><span class="append-inner-slot"><slot name="append-inner" /></span>',
};

describe('TextArea', () => {
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
    return mount(TextArea, {
      props: {modelValue: '', ...props},
      global: {
        stubs: {...vuetifyStubs, 'v-textarea': vTextareaStub, Tooltip: TooltipStub},
      },
    });
  };

  it('renders without throwing', () => {
    const wrapper = createWrapper();

    expect(wrapper.findComponent(vTextareaStub).exists()).toBe(true);
  });

  it('emits update:modelValue when the value changes', async () => {
    const wrapper = createWrapper();

    await wrapper.find('textarea').setValue('hello');

    expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['hello']);
  });

  it('uses the solo variant and readonly when readonly is true', () => {
    const wrapper = createWrapper({readonly: true});
    const textarea = wrapper.findComponent(vTextareaStub);

    expect(textarea.props('variant')).toBe('solo');
  });

  it('uses the outlined variant when not readonly', () => {
    const wrapper = createWrapper({readonly: false});

    expect(wrapper.findComponent(vTextareaStub).props('variant')).toBe('outlined');
  });

  it('applies the required rule when required is true', () => {
    const wrapper = createWrapper({required: true});
    const rules = wrapper.findComponent(vTextareaStub).props('rules') as Array<
      (v: unknown) => boolean | string
    >;

    expect(rules).toHaveLength(1);
    expect(rules[0]('')).toBe('Required.');
  });

  it('renders a help tooltip when help is set', () => {
    const wrapper = createWrapper({help: 'Some help text'});

    expect(wrapper.findComponent(TooltipStub).props('text')).toBe('Some help text');
  });
});
