// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {Group, Rights} from '@shared/user/models/Rights';
import {useUserStore} from '@shared/user/stores/user.store';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {flushPromises} from '@vue/test-utils';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import PolicyRuleClassificationMatrix from '../PolicyRuleClassificationMatrix.vue';

const {getAllObligationsMock, deletePolicyRuleMock, getAllPolicyRulesMock, snackbarInfoMock, snackbarErrorMock} =
  vi.hoisted(() => ({
    getAllObligationsMock: vi.fn(),
    deletePolicyRuleMock: vi.fn(),
    getAllPolicyRulesMock: vi.fn(),
    snackbarInfoMock: vi.fn(),
    snackbarErrorMock: vi.fn(),
  }));

vi.mock('@disclosure-portal/services/admin', () => ({
  default: {
    getAllObligations: getAllObligationsMock,
    deletePolicyRule: deletePolicyRuleMock,
  },
}));

vi.mock('@disclosure-portal/services/policyrules', () => ({
  default: {
    getAllPolicyRules: getAllPolicyRulesMock,
  },
}));

vi.mock('@shared/composables/useSnackbar', () => ({
  default: () => ({
    info: snackbarInfoMock,
    error: snackbarErrorMock,
  }),
}));

const vDataTableStub = {
  props: ['headers', 'items', 'loading'],
  template: '<div class="v-data-table">{{ JSON.stringify(items) }}</div>',
};

const childStubs = {
  'v-data-table': vDataTableStub,
  PolicyRuleClassificationDialog: true,
  ConfirmationDialog: true,
  TableLayout: {template: '<div><slot name="buttons"/><slot name="table"/></div>'},
  DSearchField: true,
};

const createWrapper = (groups: Group[] = [Group.UserPolicyManager]) =>
  mountView(PolicyRuleClassificationMatrix, {
    childStubs,
    beforePiniaMount: (pinia) => {
      useUserStore(pinia).setSimpleProfileData({
        rights: {groups} as unknown as Rights,
        profile: {} as never,
        allowed: true,
      });
    },
  });

describe('PolicyRuleClassificationMatrix', () => {
  beforeEach(() => {
    getAllObligationsMock.mockReset();
    deletePolicyRuleMock.mockReset();
    getAllPolicyRulesMock.mockReset();
    snackbarInfoMock.mockReset();
    snackbarErrorMock.mockReset();
  });

  it('renders without throwing', async () => {
    getAllObligationsMock.mockResolvedValue({data: {items: []}});
    getAllPolicyRulesMock.mockResolvedValue({data: []});

    const {wrapper} = createWrapper();
    await flushPromises();

    expect(wrapper.exists()).toBe(true);
  });

  it('loads the matrix on mount and passes only calculated rules to the data table', async () => {
    getAllObligationsMock.mockResolvedValue({data: {items: []}});
    getAllPolicyRulesMock.mockResolvedValue({
      data: [
        {_key: 'rule1', name: 'RuleOne', calculated: true},
        {_key: 'rule2', name: 'RuleTwo', calculated: false},
      ],
    });

    const {wrapper} = createWrapper();
    await flushPromises();

    expect(getAllPolicyRulesMock).toHaveBeenCalledTimes(1);
    expect(getAllObligationsMock).toHaveBeenCalledTimes(1);
    const table = wrapper.find('.v-data-table');
    expect(table.text()).toContain('RuleOne');
    expect(table.text()).not.toContain('RuleTwo');
  });

  it('shows a snackbar error when the matrix fails to load', async () => {
    getAllObligationsMock.mockResolvedValue({data: {items: []}});
    getAllPolicyRulesMock.mockRejectedValue(new Error('boom'));

    const {wrapper} = createWrapper();
    await (wrapper.vm as unknown as {loadMatrix: () => Promise<void>}).loadMatrix();
    await flushPromises();

    expect(snackbarErrorMock).toHaveBeenCalledWith('Failed to load classification matrix');
  });

  it('deletes a policy rule and reloads the matrix', async () => {
    getAllObligationsMock.mockResolvedValue({data: {items: []}});
    getAllPolicyRulesMock.mockResolvedValue({data: [{_key: 'rule1', name: 'RuleOne', calculated: true}]});
    deletePolicyRuleMock.mockResolvedValue({});

    const {wrapper} = createWrapper();
    await flushPromises();
    getAllPolicyRulesMock.mockClear();

    await (wrapper.vm as unknown as {doDelete: (config: {key: string}) => Promise<void>}).doDelete({
      key: 'rule1',
    } as never);
    await flushPromises();

    expect(deletePolicyRuleMock).toHaveBeenCalledWith('rule1');
    expect(snackbarInfoMock).toHaveBeenCalledWith('Rule deleted');
    expect(getAllPolicyRulesMock).toHaveBeenCalledTimes(1);
  });
});
