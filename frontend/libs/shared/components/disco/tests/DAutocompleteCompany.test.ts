// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {Department} from '@shared/model/Department';
import {config, mount} from '@vue/test-utils';
import {afterAll, afterEach, beforeAll, beforeEach, describe, expect, it, vi} from 'vitest';
import {nextTick} from 'vue';
import DAutocompleteCompany from '../DAutocompleteCompany.vue';

const {companyFindMock} = vi.hoisted(() => ({companyFindMock: vi.fn()}));

vi.mock('@disclosure-portal/services/companies', () => ({
  default: {find: companyFindMock},
}));

const buildDept = (overrides: Partial<Department> = {}): Department =>
  new Department({
    deptId: 'D1',
    parentDeptId: '',
    validFrom: '',
    descriptionEnglish: 'Engineering',
    orgAbbreviation: 'ENG',
    skz: '',
    companyCode: 'C1',
    companyName: 'Acme',
    level: 0,
    ...overrides,
  });

// Mirrors just the pieces DAutocompleteCompany drives: v-model, search updates and the
// no-data-text prop, via plain DOM elements the tests can interact with directly.
const vAutocompleteStub = {
  props: ['modelValue', 'items', 'noDataText', 'rules', 'disabled', 'clearable'],
  emits: ['update:modelValue', 'update:search'],
  template: `
    <div class="v-autocomplete">
      <input class="search-input" type="text" @input="$emit('update:search', $event.target.value)" />
      <div class="no-data-text">{{ noDataText }}</div>
      <ul class="items">
        <li v-for="item in items" :key="item.deptId" class="item" @click="$emit('update:modelValue', item)">
          {{ item.companyName }}
        </li>
      </ul>
      <slot name="append-inner" />
    </div>`,
};

const TooltipStub = {props: ['text'], template: '<div class="tooltip-stub">{{ text }}</div>'};

describe('DAutocompleteCompany', () => {
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

  afterEach(() => {
    vi.useRealTimers();
  });

  const createWrapper = (props: Record<string, unknown> = {}, modelValue: Department | null = null) => {
    return mount(DAutocompleteCompany, {
      props: {modelValue, ...props},
      global: {
        stubs: {'v-autocomplete': vAutocompleteStub, Tooltip: TooltipStub},
      },
    });
  };

  it('does not search while the query is shorter than 3 characters', async () => {
    vi.useFakeTimers();
    const wrapper = createWrapper();

    await wrapper.find('.search-input').setValue('ac');
    vi.advanceTimersByTime(500);

    expect(companyFindMock).not.toHaveBeenCalled();
    expect(wrapper.findComponent(vAutocompleteStub).props('noDataText')).toBe('Please type at least 3 letters');
  });

  it('debounces the search and calls companyService with a lowercased, trimmed query', async () => {
    vi.useFakeTimers();
    companyFindMock.mockResolvedValue([buildDept()]);
    const wrapper = createWrapper();

    await wrapper.find('.search-input').setValue('AC');
    vi.advanceTimersByTime(100);
    await wrapper.find('.search-input').setValue(' ACME ');
    vi.advanceTimersByTime(300);
    await nextTick();
    await nextTick();

    expect(companyFindMock).toHaveBeenCalledTimes(1);
    expect(companyFindMock).toHaveBeenCalledWith('acme');
  });

  it('sorts results so level-0 departments come first', async () => {
    vi.useFakeTimers();
    companyFindMock.mockResolvedValue([buildDept({deptId: 'child', level: 2}), buildDept({deptId: 'root', level: 0})]);
    const wrapper = createWrapper();

    await wrapper.find('.search-input').setValue('acme');
    vi.advanceTimersByTime(300);
    await nextTick();
    await nextTick();

    const items = wrapper.findComponent(vAutocompleteStub).props('items') as Department[];
    expect(items.map((i) => i.deptId)).toEqual(['root', 'child']);
  });

  it('shows the no-results message once a search of 3+ chars returns nothing', async () => {
    vi.useFakeTimers();
    companyFindMock.mockResolvedValue([]);
    const wrapper = createWrapper();

    await wrapper.find('.search-input').setValue('acme');
    vi.advanceTimersByTime(300);
    await nextTick();
    await nextTick();

    expect(wrapper.findComponent(vAutocompleteStub).props('noDataText')).toBe('Search returned no results');
  });

  it('emits depChanged and updates the v-model when a department is selected', async () => {
    vi.useFakeTimers();
    const dept = buildDept();
    companyFindMock.mockResolvedValue([dept]);
    const wrapper = createWrapper();

    await wrapper.find('.search-input').setValue('acme');
    vi.advanceTimersByTime(300);
    await nextTick();
    await nextTick();

    await wrapper.find('.item').trigger('click');
    await nextTick();

    expect(wrapper.emitted('depChanged')?.[0]).toEqual([dept]);
    expect(wrapper.emitted('update:modelValue')?.at(-1)).toEqual([dept]);
  });

  it('passes the currently selected department through as v-model when it has a deptId', () => {
    const dept = buildDept();
    const wrapper = createWrapper({}, dept);

    expect(wrapper.findComponent(vAutocompleteStub).props('modelValue')).toEqual(dept);
  });

  it('treats a modelValue without a deptId as unselected', () => {
    const wrapper = createWrapper({}, new Department());

    expect(wrapper.findComponent(vAutocompleteStub).props('modelValue')).toBeNull();
  });

  it('applies the required rule only when the required prop is set', () => {
    const withoutRequired = createWrapper();
    expect(withoutRequired.findComponent(vAutocompleteStub).props('rules')).toEqual([]);

    const withRequired = createWrapper({required: true});
    const rules = withRequired.findComponent(vAutocompleteStub).props('rules') as Array<
      (v: Department | null) => boolean | string
    >;
    expect(rules[0](null)).toBe('Company / Department is required');
    expect(rules[0](buildDept())).toBe(true);
  });

  it('renders the help tooltip when the help prop is provided', () => {
    const wrapper = createWrapper({help: 'Pick your department'});

    expect(wrapper.findComponent(TooltipStub).props('text')).toBe('Pick your department');
  });

  it('does not render the help tooltip when no help text is given', () => {
    const wrapper = createWrapper();

    expect(wrapper.findComponent(TooltipStub).exists()).toBe(false);
  });
});
