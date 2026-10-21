// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {FilterSetDto} from '@shared/model/FilterSet';
import {vuetifyStubs} from '@disclosure-portal/test-utils/vuetify-stubs';
import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, beforeEach, describe, expect, it, vi} from 'vitest';
import {nextTick, ref} from 'vue';
import FilterSets from '../FilterSets.vue';

const {
  getFilterSetsMock,
  getFilterSetMock,
  createMock,
  updateMock,
  deleteMock,
  getAllWithOptionsMock,
  getAllObligationsMock,
  useLanguageStoreMock,
} = vi.hoisted(() => ({
  getFilterSetsMock: vi.fn(),
  getFilterSetMock: vi.fn(),
  createMock: vi.fn(),
  updateMock: vi.fn(),
  deleteMock: vi.fn(),
  getAllWithOptionsMock: vi.fn(),
  getAllObligationsMock: vi.fn(),
  useLanguageStoreMock: vi.fn(),
}));

vi.mock('@disclosure-portal/services/filtersets', () => ({
  default: {
    getFilterSets: getFilterSetsMock,
    getFilterSet: getFilterSetMock,
    create: createMock,
    update: updateMock,
    delete: deleteMock,
  },
}));

vi.mock('@disclosure-portal/services/license', () => ({
  default: {getAllWithOptions: getAllWithOptionsMock},
}));

vi.mock('@disclosure-portal/services/admin', () => ({
  default: {getAllObligations: getAllObligationsMock},
}));

// setOptionsAndFilters() calls useViewTools(), which calls useLanguageStore() internally - that
// needs Vuetify's real useLocale() composable, unavailable in this plain component test. Mock the
// store module directly, mirroring disco/tests/PolicyRulesTable.test.ts.
vi.mock('@shared/stores/language.store', () => ({
  useLanguageStore: useLanguageStoreMock,
}));

const licenseMetaResponse = {
  data: {
    meta: {
      possibleCharts: {true: 3, false: 2},
      possibleSources: {SRC_A: 2},
      possibleFamilies: {MIT: 4},
      possibleApproval: {APPROVED: 1},
      possibleType: {PERMISSIVE: 2},
      possibleClassifications: [],
    },
  },
};

const buildFilterSet = (overrides: Partial<FilterSetDto> = {}): FilterSetDto => ({
  _key: 'fs-1',
  name: 'My Filter',
  includedFilters: [],
  excludedFilters: [],
  tableName: 'licenses',
  ...overrides,
});

const ConfirmationDialogStub = {
  props: ['config', 'showDialog'],
  emits: ['update:showDialog', 'confirm'],
  template: `<div v-if="showDialog" class="confirmation-dialog"><button type="button" class="confirm-btn" @click="$emit('confirm', config)">confirm</button></div>`,
};

const DCActionButtonStub = {
  props: ['icon', 'hint', 'text', 'large'],
  template: '<button type="button" class="dc-action-button" :data-text="text">{{ text }}</button>',
};

const DCloseButtonStub = {template: '<button type="button" class="dc-close-button" />'};
const StatisticsStub = {template: '<div class="statistics-stub" />'};

describe('FilterSets', () => {
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
    useLanguageStoreMock.mockReturnValue({appLanguage: ref('en')});
    getAllObligationsMock.mockResolvedValue({data: {items: []}});
    getFilterSetsMock.mockResolvedValue([]);
    getAllWithOptionsMock.mockResolvedValue(licenseMetaResponse);
    // Assigning selectedFilterSet (directly, or via the v-select) re-fetches it through
    // onFilterSetChange; echo back a filter set matching the requested key so that flow resolves
    // instead of leaving an unhandled rejection.
    getFilterSetMock.mockImplementation((key: string) =>
      Promise.resolve(buildFilterSet({_key: key, name: 'Existing'})),
    );
  });

  const createWrapper = () => {
    return mount(FilterSets, {
      global: {
        stubs: {
          ...vuetifyStubs,
          ConfirmationDialog: ConfirmationDialogStub,
          Statistics: StatisticsStub,
          DCActionButton: DCActionButtonStub,
          DCloseButton: DCloseButtonStub,
        },
      },
    });
  };

  it('loads classifications on mount', async () => {
    createWrapper();
    await nextTick();
    await nextTick();

    expect(getAllObligationsMock).toHaveBeenCalledTimes(1);
  });

  it('loads filter sets and license filter options when the menu is opened', async () => {
    getFilterSetsMock.mockResolvedValue([buildFilterSet()]);
    const wrapper = createWrapper();
    await nextTick();

    (wrapper.vm as unknown as {filterSetMenu: boolean}).filterSetMenu = true;
    await nextTick();
    await nextTick();
    await nextTick();

    expect(getFilterSetsMock).toHaveBeenCalledWith('licenses');
    expect(getAllWithOptionsMock).toHaveBeenCalledTimes(1);
    expect((wrapper.vm as unknown as {filterSets: FilterSetDto[]}).filterSets).toHaveLength(1);
  });

  it('doNew switches into create mode with a default name and empty filters', async () => {
    const wrapper = createWrapper();
    const vm = wrapper.vm as unknown as {
      doNew: () => void;
      isNew: boolean;
      filterSetName: string;
      selectedFilters: Record<string, string[]>;
    };

    vm.doNew();
    await nextTick();

    expect(vm.isNew).toBe(true);
    expect(vm.filterSetName).toBe('New Filter');
    expect(vm.selectedFilters).toEqual({
      isLicenseChart: [],
      source: [],
      family: [],
      approvalState: [],
      licenseType: [],
    });
    expect(wrapper.find('input.required').exists()).toBe(true);
  });

  it('doDelete does nothing when no filter set is selected', () => {
    const wrapper = createWrapper();
    const vm = wrapper.vm as unknown as {doDelete: () => void; confirmVisible: boolean};

    vm.doDelete();

    expect(vm.confirmVisible).toBe(false);
  });

  it('doDelete opens the confirmation dialog with the selected filter set', async () => {
    const wrapper = createWrapper();
    const vm = wrapper.vm as unknown as {
      doDelete: () => void;
      selectedFilterSet: FilterSetDto | null;
      confirmVisible: boolean;
      confirmConfig: {key: string; name: string};
    };

    // Assigning selectedFilterSet re-fetches it (mocked to resolve name: 'Existing'), so the
    // fetch must settle before doDelete reads it.
    vm.selectedFilterSet = buildFilterSet({_key: 'fs-42', name: 'Team Filter'});
    await nextTick();
    await nextTick();
    vm.doDelete();
    await nextTick();

    expect(vm.confirmVisible).toBe(true);
    expect(vm.confirmConfig.key).toBe('fs-42');
    expect(vm.confirmConfig.name).toBe('Existing');
  });

  it('deletes the filter set and closes the menu once the deletion is confirmed', async () => {
    const wrapper = createWrapper();
    const vm = wrapper.vm as unknown as {
      selectedFilterSet: FilterSetDto | null;
      confirmVisible: boolean;
      filterSetMenu: boolean;
    };
    vm.filterSetMenu = true;
    await nextTick();
    vm.selectedFilterSet = buildFilterSet({_key: 'fs-42'});
    await nextTick();
    vm.confirmVisible = true;
    await nextTick();

    await wrapper.find('.confirm-btn').trigger('click');
    await nextTick();
    await nextTick();

    expect(deleteMock).toHaveBeenCalledWith('fs-42');
    expect(wrapper.emitted('reloadFilter')?.[0]).toEqual(['fs-42']);
    expect(vm.filterSetMenu).toBe(false);
  });

  it('onClickCreate creates a new filter set from the selected filters when valid', async () => {
    createMock.mockResolvedValue({});
    const wrapper = createWrapper();
    const vm = wrapper.vm as unknown as {
      isValid: boolean;
      doNew: () => void;
      filterSetName: string;
      selectedFilters: Record<string, string[]>;
      onClickCreate: () => Promise<void>;
      filterSetMenu: boolean;
    };

    vm.filterSetMenu = true;
    await nextTick();
    vm.doNew();
    vm.filterSetName = 'My New Filter';
    vm.selectedFilters.source = ['SRC_A'];
    vm.isValid = true;
    await nextTick();

    await vm.onClickCreate();
    await nextTick();

    expect(createMock).toHaveBeenCalledTimes(1);
    const payload = createMock.mock.calls[0][0];
    expect(payload.name).toBe('My New Filter');
    expect(payload.tableName).toBe('licenses');
    expect(payload.includedFilters).toContainEqual({name: 'source', values: ['SRC_A']});
    expect(wrapper.emitted('reloadFilter')).toBeTruthy();
    expect(vm.filterSetMenu).toBe(false);
  });

  it('onClickCreate does not create a filter set when the form is invalid', async () => {
    const wrapper = createWrapper();
    const vm = wrapper.vm as unknown as {isValid: boolean; onClickCreate: () => Promise<void>};

    vm.isValid = false;
    await vm.onClickCreate();

    expect(createMock).not.toHaveBeenCalled();
  });

  it('onClickSave updates the currently selected filter set when valid', async () => {
    updateMock.mockResolvedValue({});
    const wrapper = createWrapper();
    const vm = wrapper.vm as unknown as {
      isValid: boolean;
      selectedFilterSet: FilterSetDto | null;
      onClickSave: () => Promise<void>;
      filterSetMenu: boolean;
    };

    vm.filterSetMenu = true;
    await nextTick();
    vm.selectedFilterSet = buildFilterSet({_key: 'fs-7', name: 'Existing'});
    await nextTick();
    vm.isValid = true;

    await vm.onClickSave();
    await nextTick();

    expect(updateMock).toHaveBeenCalledTimes(1);
    const [payload, key] = updateMock.mock.calls[0];
    expect(key).toBe('fs-7');
    expect(payload.name).toBe('Existing');
    expect(vm.filterSetMenu).toBe(false);
  });

  it('onClickSave does nothing when no filter set is selected', async () => {
    const wrapper = createWrapper();
    const vm = wrapper.vm as unknown as {isValid: boolean; onClickSave: () => Promise<void>};

    vm.isValid = true;
    await vm.onClickSave();

    expect(updateMock).not.toHaveBeenCalled();
  });

  it('onClickCancel exits create mode without closing the menu when isNew is true', async () => {
    const wrapper = createWrapper();
    const vm = wrapper.vm as unknown as {
      doNew: () => void;
      onClickCancel: () => void;
      isNew: boolean;
      filterSetMenu: boolean;
    };

    vm.filterSetMenu = true;
    await nextTick();
    vm.doNew();
    await nextTick();

    vm.onClickCancel();
    await nextTick();

    expect(vm.isNew).toBe(false);
    expect(vm.filterSetMenu).toBe(true);
  });

  it('onClickCancel closes the menu when not in create mode', async () => {
    const wrapper = createWrapper();
    const vm = wrapper.vm as unknown as {onClickCancel: () => void; filterSetMenu: boolean};

    vm.filterSetMenu = true;
    await nextTick();

    vm.onClickCancel();
    await nextTick();

    expect(vm.filterSetMenu).toBe(false);
  });

  it('close() collapses the menu', async () => {
    const wrapper = createWrapper();
    const vm = wrapper.vm as unknown as {close: () => void; filterSetMenu: boolean};

    vm.filterSetMenu = true;
    await nextTick();

    vm.close();
    await nextTick();

    expect(vm.filterSetMenu).toBe(false);
  });

  it('classificationCheckboxCss marks a classification checked only when it is selected', () => {
    const wrapper = createWrapper();
    const vm = wrapper.vm as unknown as {
      selectedFilters: Record<string, string[]>;
      classificationCheckboxCss: (state: string) => string;
    };

    vm.selectedFilters.classifications = ['Copyleft'];

    expect(vm.classificationCheckboxCss('Copyleft')).toContain('checkbox-marked');
    expect(vm.classificationCheckboxCss('Permissive')).toContain('checkbox-blank-outline');
  });

  it('getWarnLevel falls back to INFORMATION for an unknown classification', () => {
    const wrapper = createWrapper();
    const vm = wrapper.vm as unknown as {getWarnLevel: (name: string) => string};

    expect(vm.getWarnLevel('unknown-classification')).toBe('INFORMATION');
  });

  it('setFilterData replaces the selected filters', async () => {
    const wrapper = createWrapper();
    const vm = wrapper.vm as unknown as {
      setFilterData: (filters: Record<string, string[]>) => void;
      selectedFilters: Record<string, string[]>;
    };

    vm.setFilterData({source: ['SRC_B']});
    await nextTick();

    expect(vm.selectedFilters).toEqual({source: ['SRC_B']});
  });
});
