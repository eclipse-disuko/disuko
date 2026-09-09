// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {vuetifyStubs} from '@disclosure-portal/test-utils/vuetify-stubs';
import {PolicyState} from '@shared/model/PolicyRule';
import {IRuleBtnCallbacks} from '@shared/components/disco/interfaces';
import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, beforeEach, describe, expect, it, vi} from 'vitest';
import DRuleButtons from '../DRuleButtons.vue';

const {openUrlMock} = vi.hoisted(() => ({openUrlMock: vi.fn()}));

// DRuleButtons calls useUrls(), which pulls in useEventKeysStore() via pinia - mock the
// composable module directly rather than wiring up createTestingPinia + router internals.
vi.mock('@shared/composables/useUrls', () => ({
  useUrls: () => ({openUrl: openUrlMock}),
}));

const buildCallbacks = (overrides: Partial<IRuleBtnCallbacks> = {}): IRuleBtnCallbacks => ({
  getUrlToComponents: vi.fn().mockReturnValue('/components?policy=allow'),
  handlePolicySelect: vi.fn(),
  getCountForPolicyFilterBtn: vi.fn().mockReturnValue(3),
  getToolTipKeyForPolicyFilterBtn: vi.fn().mockReturnValue('TOOLTIP_KEY'),
  getActiveClassForPolicyFilterBtn: vi.fn().mockReturnValue('active-allow'),
  ...overrides,
});

describe('DRuleButtons', () => {
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

  const createWrapper = (props: Record<string, unknown> = {}) => {
    return mount(DRuleButtons, {
      props: {
        policies: [PolicyState.ALLOW],
        callbacks: buildCallbacks(),
        selectedPolicies: [],
        ...props,
      },
      global: {
        stubs: vuetifyStubs,
      },
    });
  };

  it('renders nothing when policies is falsy', () => {
    const wrapper = createWrapper({policies: null});

    expect(wrapper.find('[data-testid="ruleButtons"]').exists()).toBe(false);
  });

  it('renders one button per policy with its count', () => {
    const callbacks = buildCallbacks({getCountForPolicyFilterBtn: vi.fn().mockReturnValue(5)});
    const wrapper = createWrapper({policies: [PolicyState.ALLOW, PolicyState.DENY], callbacks});

    const buttons = wrapper.findAllComponents(vuetifyStubs['v-btn']);
    expect(buttons).toHaveLength(2);
    expect(buttons[0].text()).toContain('5');
  });

  it('disables the button and shows an alert icon for a warn policy when the count is zero and not force-clickable', () => {
    const callbacks = buildCallbacks({getCountForPolicyFilterBtn: vi.fn().mockReturnValue(0)});
    const wrapper = createWrapper({callbacks, policies: [PolicyState.WARN], forceClickable: false});

    const button = wrapper.findComponent(vuetifyStubs['v-btn']);
    expect(button.props('disabled')).toBe(true);
    expect(wrapper.findComponent(vuetifyStubs['v-icon']).props('icon')).toBe('mdi-alert');
  });

  it('keeps the button enabled with a zero count when forceClickable is true', () => {
    const callbacks = buildCallbacks({getCountForPolicyFilterBtn: vi.fn().mockReturnValue(0)});
    const wrapper = createWrapper({callbacks, forceClickable: true});

    const button = wrapper.findComponent(vuetifyStubs['v-btn']);
    expect(button.props('disabled')).toBe(false);
  });

  it('applies the highlighted color and active class when the policy is selected', () => {
    const callbacks = buildCallbacks({getActiveClassForPolicyFilterBtn: vi.fn().mockReturnValue('active-allow')});
    const wrapper = createWrapper({callbacks, selectedPolicies: [PolicyState.ALLOW]});

    const button = wrapper.findComponent(vuetifyStubs['v-btn']);
    expect(button.props('variant')).toBe('flat');
    expect(button.classes()).toContain('active');
    expect(button.classes()).toContain('active-allow');
  });

  it('uses the tonal variant and no active class when the policy is not selected', () => {
    const wrapper = createWrapper({selectedPolicies: []});

    const button = wrapper.findComponent(vuetifyStubs['v-btn']);
    expect(button.props('variant')).toBe('tonal');
    expect(button.classes()).not.toContain('active');
  });

  it('navigates and invokes handlePolicySelect on click when clickable', async () => {
    const callbacks = buildCallbacks();
    const wrapper = createWrapper({callbacks, policies: [PolicyState.ALLOW]});

    await wrapper.findComponent(vuetifyStubs['v-btn']).trigger('click');

    expect(openUrlMock).toHaveBeenCalledWith('/components?policy=allow', expect.anything(), expect.any(Function));
    expect(callbacks.handlePolicySelect).toHaveBeenCalledWith(PolicyState.ALLOW, [PolicyState.ALLOW]);
  });

  it('ignores clicks entirely when the button is disabled (count zero, not force-clickable)', async () => {
    const callbacks = buildCallbacks({getCountForPolicyFilterBtn: vi.fn().mockReturnValue(0)});
    const wrapper = createWrapper({callbacks, forceClickable: false});

    await wrapper.findComponent(vuetifyStubs['v-btn']).trigger('click');

    expect(openUrlMock).not.toHaveBeenCalled();
    expect(callbacks.handlePolicySelect).not.toHaveBeenCalled();
  });

  it("falls back to the 'unknown_policy' translation key for a policy outside the known switch cases", () => {
    const wrapper = createWrapper({policies: ['some-unknown-policy' as PolicyState]});

    expect(wrapper.text()).toContain('unknown_policy');
  });
});
