// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {Rights} from '@shared/user/models/Rights';
import {useUserStore} from '@shared/user/stores/user.store';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {flushPromises} from '@vue/test-utils';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import PolicyRules from '../PolicyRules.vue';

const {
  getAllPolicyRulesMock,
  getPolicyLabelsMock,
  deletePolicyRuleMock,
  deprecatePolicyRuleMock,
  copyPolicyRuleMock,
  snackbarInfoMock,
} = vi.hoisted(() => ({
  getAllPolicyRulesMock: vi.fn(),
  getPolicyLabelsMock: vi.fn(),
  deletePolicyRuleMock: vi.fn(),
  deprecatePolicyRuleMock: vi.fn(),
  copyPolicyRuleMock: vi.fn(),
  snackbarInfoMock: vi.fn(),
}));

vi.mock('@disclosure-portal/services/admin', () => ({
  default: {
    getPolicyLabels: getPolicyLabelsMock,
    deletePolicyRule: deletePolicyRuleMock,
    deprecatePolicyRule: deprecatePolicyRuleMock,
    copyPolicyRule: copyPolicyRuleMock,
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
  }),
}));

const policyRule = {
  _key: 'pr-1',
  name: 'Rule One',
  description: 'desc',
  status: 'active',
  deprecated: false,
  componentsAllow: [],
  componentsWarn: [],
  componentsDeny: [],
  created: '2026-01-01T10:00:00Z',
  updated: '2026-01-02T10:00:00Z',
};

const childStubs = {
  TableLayout: {template: '<div><slot name="buttons"/><slot name="table"/></div>'},
  NewPolicyRuleDialog: {template: '<div><slot name="default" :showDialog="() => {}" /></div>'},
  ConfirmationDialog: true,
  ClassificationMatrixDialog: true,
  'v-data-table': {
    props: ['items'],
    template: '<div class="v-data-table">{{ JSON.stringify(items) }}</div>',
  },
};

// PolicyRules calls RightsUtils.isPolicyManager() on mount, which reads `userStore.getRights.groups`
// directly - seed the user store so that lookup doesn't throw on an empty Rights object.
const createWrapper = () =>
  mountView(PolicyRules, {
    childStubs,
    beforePiniaMount: (p) => {
      useUserStore(p).setSimpleProfileData({
        rights: {groups: []} as unknown as Rights,
        profile: {} as never,
        allowed: true,
      });
    },
  });

describe('PolicyRules', () => {
  beforeEach(() => {
    getAllPolicyRulesMock.mockReset();
    getPolicyLabelsMock.mockReset();
    deletePolicyRuleMock.mockReset();
    deprecatePolicyRuleMock.mockReset();
    copyPolicyRuleMock.mockReset();
    snackbarInfoMock.mockReset();
    getAllPolicyRulesMock.mockResolvedValue({data: [policyRule]});
    getPolicyLabelsMock.mockResolvedValue({data: []});
  });

  it('sets the expected breadcrumbs on mount', async () => {
    const {pinia} = createWrapper();
    await flushPromises();

    const breadcrumbs = useBreadcrumbsStore(pinia);
    expect(breadcrumbs.setCurrentBreadcrumbs).toHaveBeenCalledWith([
      {title: 'Dashboard', disabled: false, href: '/dashboard/home'},
      {title: 'Policy Rules', disabled: false, href: '/dashboard/policyrules'},
    ]);
  });

  it('renders without throwing and loads policy rules into the table', async () => {
    const {wrapper} = createWrapper();
    await flushPromises();

    expect(wrapper.exists()).toBe(true);
    expect(getAllPolicyRulesMock).toHaveBeenCalledTimes(1);
    expect(wrapper.find('.v-data-table').text()).toContain('Rule One');
  });

  it('deletes a policy rule on confirm and shows a success snackbar', async () => {
    deletePolicyRuleMock.mockResolvedValue({data: {}});

    const {wrapper} = createWrapper();
    await flushPromises();

    await (wrapper.vm as unknown as {doDeletePolicyRule: (config: {key: string}) => Promise<void>}).doDeletePolicyRule({
      key: 'pr-1',
    } as never);
    await flushPromises();

    expect(deletePolicyRuleMock).toHaveBeenCalledWith('pr-1');
    expect(snackbarInfoMock).toHaveBeenCalledWith('policy rule deleted');
    expect(getAllPolicyRulesMock).toHaveBeenCalledTimes(2);
  });

  it('deprecates a policy rule on confirm and shows a success snackbar', async () => {
    deprecatePolicyRuleMock.mockResolvedValue({data: {}});

    const {wrapper} = createWrapper();
    await flushPromises();

    await (wrapper.vm as unknown as {doDeprecate: (config: {key: string}) => Promise<void>}).doDeprecate({
      key: 'pr-1',
    } as never);
    await flushPromises();

    expect(deprecatePolicyRuleMock).toHaveBeenCalledWith('pr-1');
    expect(snackbarInfoMock).toHaveBeenCalledWith('policy rule deprecated');
  });

  it('copies a policy rule on confirm and shows a success snackbar', async () => {
    copyPolicyRuleMock.mockResolvedValue({data: {}});

    const {wrapper} = createWrapper();
    await flushPromises();

    await (wrapper.vm as unknown as {doCopy: (config: {key: string}) => Promise<void>}).doCopy({
      key: 'pr-1',
    } as never);
    await flushPromises();

    expect(copyPolicyRuleMock).toHaveBeenCalledWith('pr-1');
    expect(snackbarInfoMock).toHaveBeenCalledWith('policy rule copied');
  });
});
