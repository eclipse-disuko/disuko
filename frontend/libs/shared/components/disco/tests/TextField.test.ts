// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {vuetifyStubs} from '@disclosure-portal/test-utils/vuetify-stubs';
import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, describe, expect, it} from 'vitest';
import TextField from '../TextField.vue';

const TooltipStub = {
  template: '<div class="tooltip-stub"><slot />{{ text }}</div>',
  props: ['text', 'asParent'],
};

// vuetifyStubs' generic v-text-field stub only renders the default slot; TextField.vue uses the
// named #append-inner slot for the help Tooltip, so render it here too.
const vTextFieldStub = {
  ...vuetifyStubs['v-text-field'],
  template:
    '<input type="text" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" /><span class="append-inner-slot"><slot name="append-inner" /></span>',
};

describe('TextField', () => {
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
    return mount(TextField, {
      props: {modelValue: '', ...props},
      global: {
        stubs: {...vuetifyStubs, 'v-text-field': vTextFieldStub, Tooltip: TooltipStub},
      },
    });
  };

  it('renders without throwing', () => {
    const wrapper = createWrapper();

    expect(wrapper.findComponent(vTextFieldStub).exists()).toBe(true);
  });

  it('emits update:modelValue when the value changes', async () => {
    const wrapper = createWrapper();

    await wrapper.find('input').setValue('hello');

    expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['hello']);
  });

  it('uses the solo variant when readonly and outlined otherwise', () => {
    const readonly = createWrapper({readonly: true});
    expect(readonly.findComponent(vTextFieldStub).props('variant')).toBe('solo');

    const editable = createWrapper({readonly: false});
    expect(editable.findComponent(vTextFieldStub).props('variant')).toBe('outlined');
  });

  it('uses custom rules when given instead of the default required rule', () => {
    const customRules = [(v: string) => !!v || 'Custom message'];
    const wrapper = createWrapper({rules: customRules, required: true});

    expect(wrapper.findComponent(vTextFieldStub).props('rules')).toEqual(customRules);
  });

  it('applies the default required rule when required is true and no custom rules are given', () => {
    const wrapper = createWrapper({required: true});
    const rules = wrapper.findComponent(vTextFieldStub).props('rules') as Array<
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
