// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, beforeEach, describe, expect, it, vi} from 'vitest';
import {ref} from 'vue';
import HeaderSettings from '../HeaderSettings.vue';

const {useHeaderSettingsMock, updateSelectedHeadersFromStringListMock, resetSelectedHeadersMock} = vi.hoisted(() => ({
  useHeaderSettingsMock: vi.fn(),
  updateSelectedHeadersFromStringListMock: vi.fn(),
  resetSelectedHeadersMock: vi.fn(),
}));

vi.mock('@shared/composables/useHeaderSettings', () => ({
  useHeaderSettings: useHeaderSettingsMock,
}));

const GridHeaderMenuStub = {
  props: ['showReset', 'resetHint', 'cardTitle', 'allItems', 'selectedItems', 'selectLabel'],
  emits: ['update', 'reset'],
  template: `
    <div class="grid-header-menu">
      <slot name="activator" :props="{}" />
      <button class="update-btn" type="button" @click="$emit('update', ['col-a'])">update</button>
      <button class="reset-btn" type="button" @click="$emit('reset')">reset</button>
    </div>`,
};

const DIconButtonStub = {
  props: ['icon', 'hint', 'color', 'parentProps'],
  template: '<button class="d-icon-button" type="button">{{ icon }}</button>',
};

describe('HeaderSettings', () => {
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

  const selectableHeaders = [
    {value: 'col-a', text: 'Column A'},
    {value: 'col-b', text: 'Column B'},
    {value: 'col-c', text: 'Column C'},
  ];

  const setupHeaderSettingsMock = (overrides: Record<string, unknown> = {}) => {
    useHeaderSettingsMock.mockReturnValue({
      selectableHeaders: ref(selectableHeaders),
      selectedHeaders: ref([0, 2]),
      initialSelectedHeaders: ref([0, 1, 2]),
      updateSelectedHeadersFromStringList: updateSelectedHeadersFromStringListMock,
      resetSelectedHeaders: resetSelectedHeadersMock,
      ...overrides,
    });
  };

  beforeEach(() => {
    vi.clearAllMocks();
    setupHeaderSettingsMock();
  });

  const createWrapper = (props: Record<string, unknown> = {}) => {
    return mount(HeaderSettings, {
      props: {column: {value: 'col-a'}, gridName: 'myGrid', ...props},
      global: {
        stubs: {GridHeaderMenu: GridHeaderMenuStub, DIconButton: DIconButtonStub},
      },
    });
  };

  it('renders the settings icon button in the activator slot', () => {
    const wrapper = createWrapper();

    const button = wrapper.findComponent(DIconButtonStub);
    expect(button.exists()).toBe(true);
    expect(button.props('icon')).toBe('mdi-cog');
  });

  it('renders the top-left border indicator when showBorders is true', () => {
    const wrapper = createWrapper({showBorders: true});

    expect(wrapper.find('.border-t-\\[2px\\]').exists()).toBe(true);
  });

  it('does not render the border indicator when showBorders is false', () => {
    const wrapper = createWrapper({showBorders: false});

    expect(wrapper.find('.border-t-\\[2px\\]').exists()).toBe(false);
  });

  it('maps selectedHeaders indices to selectable header values, skipping out-of-range indices', () => {
    setupHeaderSettingsMock({selectedHeaders: ref([0, 2, 5])});
    const wrapper = createWrapper();

    const menu = wrapper.findComponent(GridHeaderMenuStub);
    expect(menu.props('selectedItems')).toEqual(['col-a', 'col-c']);
  });

  it('shows the reset option only when selection differs from the initial selection', () => {
    setupHeaderSettingsMock({selectedHeaders: ref([0, 1, 2]), initialSelectedHeaders: ref([0, 1, 2])});
    const unchanged = createWrapper();
    expect(unchanged.findComponent(GridHeaderMenuStub).props('showReset')).toBe(false);

    setupHeaderSettingsMock({selectedHeaders: ref([0]), initialSelectedHeaders: ref([0, 1, 2])});
    const changed = createWrapper();
    expect(changed.findComponent(GridHeaderMenuStub).props('showReset')).toBe(true);
  });

  it('forwards GridHeaderMenu update events to updateSelectedHeadersFromStringList', async () => {
    const wrapper = createWrapper();

    await wrapper.find('.update-btn').trigger('click');

    expect(updateSelectedHeadersFromStringListMock).toHaveBeenCalledWith(['col-a']);
  });

  it('forwards GridHeaderMenu reset events to resetSelectedHeaders', async () => {
    const wrapper = createWrapper();

    await wrapper.find('.reset-btn').trigger('click');

    expect(resetSelectedHeadersMock).toHaveBeenCalledTimes(1);
  });
});
