// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {UserDto} from '@shared/types/Users';
import {createTestingPinia} from '@pinia/testing';
import {config, mount} from '@vue/test-utils';
import {afterAll, afterEach, beforeAll, beforeEach, describe, expect, it, vi} from 'vitest';
import {nextTick} from 'vue';
import DAutocompleteUser from '../DAutocompleteUser.vue';

const {profileSearchMock, projectSearchMock} = vi.hoisted(() => ({
  profileSearchMock: vi.fn(),
  projectSearchMock: vi.fn(),
}));

vi.mock('@shared/user/services/profile.service', () => ({
  default: {getUsersBySearchFragment: profileSearchMock},
}));

vi.mock('@disclosure-portal/services/projects', () => ({
  default: {getUsersBySearchFragment: projectSearchMock},
}));

// DAutocompleteUser imports VForm from 'vuetify/components' purely for its TS type, but that's a
// value import at runtime, and Vue's SFC compiler binds the <v-form> tag directly to that local
// import rather than resolving it dynamically - so it must resolve to a real renderable component,
// not just any stub, and without mocking the module Vitest would load real Vuetify CSS that jsdom
// can't parse (mirrors the fix in widgets/tests/Idle.test.ts).
vi.mock('vuetify/components', () => ({
  VForm: {
    template: '<form><slot /></form>',
    methods: {
      validate: vi.fn(() => Promise.resolve({valid: true})),
      resetValidation: vi.fn(),
    },
  },
}));

const buildUser = (overrides: Partial<UserDto> = {}): UserDto => ({
  user: 'jdoe',
  forename: 'Jane',
  lastname: 'Doe',
  email: 'jdoe@example.com',
  termsOfUse: false,
  termsOfUseDate: '',
  termsOfUseVersion: '',
  created: '',
  updated: '',
  isSelectable: true,
  _key: 'jdoe',
  roles: [],
  metaData: {companyIdentifier: '', department: '', departmentDescription: ''},
  active: true,
  isInternal: true,
  deprovisioned: '',
  ...overrides,
});

// The real v-autocomplete needs Vuetify's full runtime for its search box and item list, so this
// stub reproduces just the bits DAutocompleteUser relies on: v-model, v-model:search and the
// no-data slot, driven by plain DOM elements the tests can interact with.
const vAutocompleteStub = {
  props: ['modelValue', 'search', 'items', 'rules', 'hideNoData', 'itemTitle', 'disabled'],
  emits: ['update:modelValue', 'update:search', 'click:clear'],
  template: `
    <div class="v-autocomplete">
      <input class="search-input" type="text" :value="search" @input="$emit('update:search', $event.target.value)" />
      <button class="clear-btn" type="button" @click="$emit('click:clear')">clear</button>
      <ul class="items">
        <li v-for="item in items" :key="item.user" class="item" @click="$emit('update:modelValue', item)">
          {{ itemTitle(item) }}
        </li>
      </ul>
      <div v-if="!hideNoData" class="no-data"><slot name="no-data" /></div>
    </div>`,
};

const DCActionButtonStub = {
  props: ['icon', 'hint', 'text'],
  template: '<button type="button" class="dc-action-button">{{ text }}</button>',
};

describe('DAutocompleteUser', () => {
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

  const createWrapper = (props: Record<string, unknown> = {}) => {
    return mount(DAutocompleteUser, {
      props,
      global: {
        plugins: [createTestingPinia({createSpy: vi.fn, stubActions: false})],
        stubs: {
          'v-autocomplete': vAutocompleteStub,
          DCActionButton: DCActionButtonStub,
        },
      },
    });
  };

  it('does not search while the query is shorter than 3 characters', async () => {
    vi.useFakeTimers();
    const wrapper = createWrapper();

    await wrapper.find('.search-input').setValue('ab');
    vi.advanceTimersByTime(500);
    await nextTick();

    expect(profileSearchMock).not.toHaveBeenCalled();
  });

  it('debounces the search and calls profileService once the query settles', async () => {
    vi.useFakeTimers();
    profileSearchMock.mockResolvedValue({data: [buildUser()]});
    const wrapper = createWrapper();

    await wrapper.find('.search-input').setValue('jd');
    vi.advanceTimersByTime(100);
    await wrapper.find('.search-input').setValue('jdo');
    await vi.advanceTimersByTimeAsync(300);

    expect(profileSearchMock).toHaveBeenCalledTimes(1);
    expect(profileSearchMock).toHaveBeenCalledWith('jdo', true);
  });

  it('routes the search through projectService when projectKey is set', async () => {
    vi.useFakeTimers();
    projectSearchMock.mockResolvedValue({data: [buildUser()]});
    const wrapper = createWrapper({projectKey: 'proj-1'});

    await wrapper.find('.search-input').setValue('jdo');
    await vi.advanceTimersByTimeAsync(300);

    expect(projectSearchMock).toHaveBeenCalledWith('proj-1', 'jdo', true);
    expect(profileSearchMock).not.toHaveBeenCalled();
  });

  it('filters results to internal users only when onlyInternalUsers is set', async () => {
    vi.useFakeTimers();
    profileSearchMock.mockResolvedValue({
      data: [buildUser({user: 'internal', isInternal: true}), buildUser({user: 'external', isInternal: false})],
    });
    const wrapper = createWrapper({onlyInternalUsers: true});

    await wrapper.find('.search-input').setValue('jdo');
    await vi.advanceTimersByTimeAsync(300);

    const autocomplete = wrapper.findComponent(vAutocompleteStub);
    expect((autocomplete.props('items') as UserDto[]).map((u) => u.user)).toEqual(['internal']);
  });

  it('emits userChanged and updates the v-model when a user is selected', async () => {
    vi.useFakeTimers();
    profileSearchMock.mockResolvedValue({data: [buildUser()]});
    const wrapper = createWrapper();

    await wrapper.find('.search-input').setValue('jdo');
    await vi.advanceTimersByTimeAsync(300);

    await wrapper.find('.item').trigger('click');
    await nextTick();

    expect(wrapper.emitted('userChanged')?.[0]).toEqual([buildUser()]);
    expect(wrapper.emitted('update:modelValue')?.at(-1)).toEqual(['jdoe']);
  });

  it('preselects the given user on mount', async () => {
    const preselect = buildUser({user: 'preselected'});
    const wrapper = createWrapper({preselect});
    await nextTick();

    const autocomplete = wrapper.findComponent(vAutocompleteStub);
    expect(autocomplete.props('modelValue')).toEqual(preselect);
  });

  it('shows a "not found" no-data slot with a mailto action once a search returns nothing', async () => {
    vi.useFakeTimers();
    profileSearchMock.mockResolvedValue({data: []});
    const wrapper = createWrapper();

    await wrapper.find('.search-input').setValue('jdo');
    await vi.advanceTimersByTimeAsync(300);

    expect(wrapper.find('.no-data').text()).toContain('User not registered? Send an invite...');
    expect(wrapper.findComponent(DCActionButtonStub).exists()).toBe(true);
  });

  it('flags the required rule as failing when nothing has been selected or typed', async () => {
    const wrapper = createWrapper({required: true});

    await (wrapper.vm as unknown as {validateOnCreate: () => Promise<boolean>}).validateOnCreate();

    const autocomplete = wrapper.findComponent(vAutocompleteStub);
    const rules = autocomplete.props('rules') as Array<() => boolean | string>;
    expect(rules[0]()).toBe('User is required');
  });

  it('flags the alphanumeric rule as failing for special characters', async () => {
    vi.useFakeTimers();
    const wrapper = createWrapper({onlyAlphanumeric: true});

    await wrapper.find('.search-input').setValue('jd@!');
    await vi.advanceTimersByTimeAsync(300);
    await (wrapper.vm as unknown as {validateOnCreate: () => Promise<boolean>}).validateOnCreate();

    const autocomplete = wrapper.findComponent(vAutocompleteStub);
    const rules = autocomplete.props('rules') as Array<() => boolean | string>;
    expect(rules[0]()).toBe('Alpha-numeric characters only are allowed');
  });

  it('resetForm clears the search input and items', async () => {
    vi.useFakeTimers();
    profileSearchMock.mockResolvedValue({data: [buildUser()]});
    const wrapper = createWrapper();

    await wrapper.find('.search-input').setValue('jdo');
    await vi.advanceTimersByTimeAsync(300);
    expect(wrapper.findComponent(vAutocompleteStub).props('items')).toHaveLength(1);

    (wrapper.vm as unknown as {resetForm: () => void}).resetForm();
    await nextTick();

    expect(wrapper.findComponent(vAutocompleteStub).props('items')).toHaveLength(0);
  });

  it('clears the search state when the clear button is clicked', async () => {
    vi.useFakeTimers();
    profileSearchMock.mockResolvedValue({data: [buildUser()]});
    const wrapper = createWrapper();

    await wrapper.find('.search-input').setValue('jdo');
    await vi.advanceTimersByTimeAsync(300);

    await wrapper.find('.clear-btn').trigger('click');
    await nextTick();

    expect(wrapper.findComponent(vAutocompleteStub).props('items')).toHaveLength(0);
  });

  it('prepareForCreate resets the selected user and item list', async () => {
    const preselect = buildUser({user: 'preselected'});
    const wrapper = createWrapper({preselect});

    (wrapper.vm as unknown as {prepareForCreate: () => void}).prepareForCreate();
    await nextTick();

    const autocomplete = wrapper.findComponent(vAutocompleteStub);
    expect(autocomplete.props('modelValue')).toBeNull();
    expect(autocomplete.props('items')).toHaveLength(0);
  });
});
