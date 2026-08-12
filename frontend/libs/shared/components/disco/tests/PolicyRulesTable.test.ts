// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {vuetifyStubs} from '@disclosure-portal/test-utils/vuetify-stubs';
import {PolicyRulesAssignmentsDto} from '@disclosure-portal/model/PolicyRule';
import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, beforeEach, describe, expect, it, vi} from 'vitest';
import {ref} from 'vue';
import PolicyRulesTable from '../PolicyRulesTable.vue';

const {useLanguageStoreMock} = vi.hoisted(() => ({useLanguageStoreMock: vi.fn()}));

// PolicyRulesTable (via its own import and via useViewTools()) calls useLanguageStore(), which
// internally calls Vuetify's useLocale() composable - that requires a real Vuetify plugin instance
// these component tests don't set up. Mock the store module directly instead of using
// createTestingPinia, which would still execute the real store setup (and hit the same issue).
vi.mock('@shared/stores/language.store', () => ({
  useLanguageStore: useLanguageStoreMock,
}));

const TooltipStub = {
  template: '<div class="tooltip-stub">{{ text }}</div>',
  props: ['text'],
};

const DIconButtonStub = {
  props: ['icon', 'hint', 'color', 'parentProps'],
  template: '<button class="d-icon-button" type="button">{{ icon }}</button>',
};

const DCloseButtonStub = {
  template: '<button class="d-close-button" type="button" @click="$emit(\'click\')" />',
};

// The generic v-data-table-virtual stub renders the headers and the item slots this
// component actually uses, so we can assert on language/status-driven rendering per row.
const vDataTableVirtualStub = {
  props: ['loading', 'headers', 'items', 'search', 'sortBy'],
  template: `
    <div class="v-data-table-virtual">
      <div class="headers">
        <template v-for="h in headers" :key="h.value">
          <span class="header-title" :data-value="h.value" :data-width="h.width">{{ h.title }}</span>
          <slot :name="'header.' + h.value" :column="h" :getSortIcon="() => ''" :toggleSort="() => {}" />
        </template>
      </div>
      <div class="rows">
        <div v-for="item in items" :key="item.key" class="row" :data-key="item.key">
          <slot name="item.status" :item="item" />
          <slot name="item.description" :item="item" />
          <slot name="item.type" :item="item" />
        </div>
      </div>
    </div>`,
};

const buildItem = (overrides: Partial<PolicyRulesAssignmentsDto> = {}): PolicyRulesAssignmentsDto => ({
  status: 'active',
  key: 'k1',
  name: 'Rule 1',
  description: 'A short description',
  type: 'allow',
  ...overrides,
});

describe('PolicyRulesTable', () => {
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
  });

  const createWrapper = (props: Record<string, unknown> = {}, modelValue: PolicyRulesAssignmentsDto[] = []) => {
    return mount(PolicyRulesTable, {
      props: {modelValue, ...props},
      global: {
        stubs: {
          ...vuetifyStubs,
          'v-data-table-virtual': vDataTableVirtualStub,
          Tooltip: TooltipStub,
          DIconButton: DIconButtonStub,
          DCloseButton: DCloseButtonStub,
        },
      },
    });
  };

  it('filters out items whose status is not in the active status filter by default', () => {
    const items = [
      buildItem({key: 'active-rule', status: 'active'}),
      buildItem({key: 'deprecated-rule', status: 'deprecated'}),
    ];
    const wrapper = createWrapper({}, items);

    expect(wrapper.find('[data-key="active-rule"]').exists()).toBe(true);
    expect(wrapper.find('[data-key="deprecated-rule"]').exists()).toBe(false);
  });

  it('omits the description column when isDialog is true', () => {
    const withDialog = createWrapper({isDialog: true});
    const withoutDialog = createWrapper({isDialog: false});

    expect(withDialog.find('.header-title[data-value="description"]').exists()).toBe(false);
    expect(withoutDialog.find('.header-title[data-value="description"]').exists()).toBe(true);
  });

  it('renders the translated status label per row', () => {
    const items = [buildItem({key: 'k1', status: 'active'})];
    const wrapper = createWrapper({}, items);

    expect(wrapper.find('[data-key="k1"]').text()).toContain('ACTIVE');
  });

  it('renders a truncated description with a tooltip for the full text', () => {
    const longDescription = 'x'.repeat(150);
    const items = [buildItem({key: 'k1', description: longDescription})];
    const wrapper = createWrapper({}, items);

    const row = wrapper.find('[data-key="k1"]');
    expect(row.text()).toContain('x'.repeat(117) + '...');

    const tooltip = wrapper.findComponent(TooltipStub);
    expect(tooltip.props('text')).toBe(longDescription);
  });

  it('disables the type radio group when not in edit mode', () => {
    const items = [buildItem({key: 'k1', status: 'active'})];
    const wrapper = createWrapper({edit: false}, items);

    const radioGroup = wrapper.findComponent(vuetifyStubs['v-radio-group']);
    expect(radioGroup.attributes('disabled')).toBe('true');
  });

  it('enables the type radio group for active rows in edit mode', () => {
    const items = [buildItem({key: 'k1', status: 'active'})];
    const wrapper = createWrapper({edit: true}, items);

    const radioGroup = wrapper.findComponent(vuetifyStubs['v-radio-group']);
    expect(radioGroup.attributes('disabled')).toBe('false');
  });

  it('uses a narrower type column width for English than German', () => {
    useLanguageStoreMock.mockReturnValue({appLanguage: ref('en')});
    const englishWrapper = createWrapper();
    const englishHeader = englishWrapper.findAll('.header-title')[3];

    useLanguageStoreMock.mockReturnValue({appLanguage: ref('de')});
    const germanWrapper = createWrapper();
    const germanHeader = germanWrapper.findAll('.header-title')[3];

    expect(englishHeader.attributes('data-width')).toBe('450');
    expect(germanHeader.attributes('data-width')).toBe('470');
  });
});
