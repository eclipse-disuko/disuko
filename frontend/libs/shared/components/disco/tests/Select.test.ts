// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {vuetifyStubs} from '@disclosure-portal/test-utils/vuetify-stubs';
import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, describe, expect, it} from 'vitest';
import Select from '../Select.vue';

const TooltipStub = {
  template: '<div class="tooltip-stub">{{ text }}</div>',
  props: ['text', 'asParent'],
};

// vuetifyStubs' generic v-select stub only renders the default slot; Select.vue uses the
// named #append-inner slot for the help Tooltip, so render it here too.
const vSelectStub = {
  ...vuetifyStubs['v-select'],
  template:
    '<select :value="modelValue" @change="$emit(\'update:modelValue\', $event.target.value)"><slot /><slot name="append-inner" /></select>',
};

describe('Select', () => {
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

  const items = [
    {name: 'Option A', value: 'a'},
    {name: 'Option B', value: 'b'},
  ];

  const createWrapper = (props: Record<string, unknown> = {}) => {
    return mount(Select, {
      props: {
        modelValue: '',
        label: 'My Select',
        items,
        ...props,
      },
      global: {
        stubs: {...vuetifyStubs, 'v-select': vSelectStub, Tooltip: TooltipStub, 'v-icon': {template: '<i></i>'}},
      },
    });
  };

  it('renders a v-select with the given items and label', () => {
    const wrapper = createWrapper();
    const select = wrapper.findComponent(vSelectStub);

    expect(select.exists()).toBe(true);
    expect(select.props('items')).toEqual(items);
    expect(select.props('label')).toBe('My Select');
  });

  it('passes item-title/item-value keys for the default (non-simple) list', () => {
    const wrapper = createWrapper();
    const select = wrapper.findComponent(vSelectStub);

    expect(select.props('itemTitle')).toBe('name');
    expect(select.props('itemValue')).toBe('value');
  });

  it('does not pass item-title/item-value for a simpleList', () => {
    const wrapper = createWrapper({simpleList: true, items: ['a', 'b']});
    const select = wrapper.findComponent(vSelectStub);

    expect(select.props('itemTitle')).toBeUndefined();
    expect(select.props('itemValue')).toBeUndefined();
  });

  it('emits update:modelValue when the selection changes', async () => {
    const wrapper = createWrapper();

    await wrapper.findComponent(vSelectStub).vm.$emit('update:modelValue', 'b');

    expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['b']);
  });

  it('applies the required rule when required is true', () => {
    const wrapper = createWrapper({required: true});
    const select = wrapper.findComponent(vSelectStub);
    const rules = select.props('rules') as Array<(v: unknown) => boolean | string>;

    expect(rules).toHaveLength(1);
    expect(rules[0]('')).toBe('Required.');
    expect(rules[0]('value')).toBe(true);
  });

  it('applies no rules when required is false', () => {
    const wrapper = createWrapper({required: false});
    const select = wrapper.findComponent(vSelectStub);

    expect(select.props('rules')).toEqual([]);
  });

  it('uses the solo variant and disables clearable when readonly', () => {
    const wrapper = createWrapper({readonly: true});
    const select = wrapper.findComponent(vSelectStub);

    expect(select.props('variant')).toBe('solo');
    expect(select.props('readonly')).toBe(true);
  });

  it('uses the outlined variant when not readonly', () => {
    const wrapper = createWrapper({readonly: false});
    const select = wrapper.findComponent(vSelectStub);

    expect(select.props('variant')).toBe('outlined');
  });

  it('renders a help tooltip when the help prop is set', () => {
    const wrapper = createWrapper({help: 'Some help text'});

    const tooltip = wrapper.findComponent(TooltipStub);
    expect(tooltip.exists()).toBe(true);
    expect(tooltip.props('text')).toBe('Some help text');
  });

  it('does not render a help tooltip when the help prop is absent', () => {
    const wrapper = createWrapper();

    expect(wrapper.findComponent(TooltipStub).exists()).toBe(false);
  });
});
