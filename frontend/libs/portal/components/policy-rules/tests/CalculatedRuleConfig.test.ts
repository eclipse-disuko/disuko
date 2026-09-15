// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import PolicyRule from '@disclosure-portal/model/PolicyRule';
import {useCalculatedPolicyRuleStore} from '@disclosure-portal/stores/calculatedPolicyRule.store';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {describe, expect, it, vi} from 'vitest';
import {nextTick} from 'vue';
import CalculatedRuleConfig from '../CalculatedRuleConfig.vue';

vi.mock('@disclosure-portal/utils/View', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@disclosure-portal/utils/View')>();
  return {
    ...actual,
    default: () => ({
      getNameForLanguage: (obligation: {name?: string} | null) => obligation?.name ?? '',
    }),
  };
});

const DMultiSelectStub = {
  props: ['label', 'items', 'modelValue'],
  emits: ['update:modelValue'],
  template: '<div class="d-multi-select">{{ label }}</div>',
};

const createWrapper = () => {
  const rule = new PolicyRule();
  rule.calculated = true;
  rule.calculatedConfig.bucketDefinition = {
    deniedClassifications: ['denied'],
    warnedClassifications: ['warned'],
    allowedClassifications: ['allowed'],
  };
  rule.calculatedConfig.licenseScope = {
    isLicenseChart: [true],
    approvalState: ['approved'],
    family: ['permissive'],
    licenseType: ['open source'],
    source: ['spdx'],
  };

  const {wrapper, pinia} = mountView(CalculatedRuleConfig, {
    childStubs: {DMultiSelect: DMultiSelectStub},
    piniaOptions: {
      initialState: {
        calculatedPolicyRule: {rule},
      },
    },
  });

  return {wrapper, store: useCalculatedPolicyRuleStore(pinia)};
};

describe('CalculatedRuleConfig', () => {
  it('renders bucket and scope selectors for calculated rules', () => {
    const {wrapper} = createWrapper();

    const selectors = wrapper.findAllComponents(DMultiSelectStub);

    expect(selectors).toHaveLength(8);
    expect(selectors.map((selector) => selector.props('modelValue'))).toEqual([
      ['denied'],
      ['warned'],
      ['allowed'],
      ['true'],
      ['approved'],
      ['permissive'],
      ['open source'],
      ['spdx'],
    ]);
  });

  it('updates the selected bucket and scope through the store', async () => {
    const {wrapper, store} = createWrapper();
    const selectors = wrapper.findAllComponents(DMultiSelectStub);

    await selectors[0].vm.$emit('update:modelValue', ['new-denied']);
    await selectors[3].vm.$emit('update:modelValue', ['false']);

    expect(store.rule.calculatedConfig.bucketDefinition.deniedClassifications).toEqual(['new-denied']);
    expect(store.rule.calculatedConfig.licenseScope.isLicenseChart).toEqual([false]);
  });

  it('does not render selectors for non-calculated rules', async () => {
    const {wrapper, store} = createWrapper();
    store.setCalculated(false);
    await nextTick();

    expect(wrapper.findAllComponents(DMultiSelectStub)).toHaveLength(0);
  });
});
