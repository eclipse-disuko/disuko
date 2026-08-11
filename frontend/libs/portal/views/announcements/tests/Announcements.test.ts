// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {Group, Rights} from '@shared/user/models/Rights';
import {useUserStore} from '@shared/user/stores/user.store';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {flushPromises} from '@vue/test-utils';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import Announcements from '../Announcements.vue';

const {getAllMock} = vi.hoisted(() => ({
  getAllMock: vi.fn(),
}));

vi.mock('@disclosure-portal/services/announcements', () => ({
  default: {
    getAll: getAllMock,
  },
}));

const vDataTableStub = {
  props: ['headers', 'items'],
  template: '<div class="v-data-table">{{ JSON.stringify(items) }}</div>',
};

const childStubs = {
  'v-data-table': vDataTableStub,
  TableLayout: {template: '<div><slot name="buttons"/><slot name="table"/></div>'},
  DSearchField: true,
};

const createWrapper = (groups: Group[] = []) =>
  mountView(Announcements, {
    childStubs,
    beforePiniaMount: (pinia) => {
      useUserStore(pinia).setSimpleProfileData({
        rights: {groups} as unknown as Rights,
        profile: {} as never,
        allowed: true,
      });
    },
  });

const announcement = (overrides: Record<string, unknown> = {}) => ({
  _key: 'a1',
  when: '2026-01-01T00:00:00Z',
  type: 'license_change',
  content: JSON.stringify({
    licenseName: 'MIT',
    licenseId: 'lic-1',
    changeType: 'license_forbidden',
    oldVal: 'old',
    newVal: 'new',
    ...overrides,
  }),
});

describe('Announcements', () => {
  beforeEach(() => {
    getAllMock.mockReset();
  });

  it('sets the expected breadcrumbs on mount', () => {
    getAllMock.mockResolvedValue({data: []});

    const {pinia} = createWrapper();

    const breadcrumbs = useBreadcrumbsStore(pinia);
    expect(breadcrumbs.setCurrentBreadcrumbs).toHaveBeenCalledWith([
      {title: 'Dashboard', disabled: false, href: '/dashboard/home'},
      {title: 'Announcements', disabled: false, href: '/dashboard/announcements/'},
    ]);
  });

  it('loads announcements on mount, parses their content, and passes them to the data table', async () => {
    getAllMock.mockResolvedValue({data: [announcement()]});

    const {wrapper} = createWrapper();
    await flushPromises();

    expect(getAllMock).toHaveBeenCalledTimes(1);
    const table = wrapper.find('.v-data-table');
    expect(table.text()).toContain('MIT');
    expect(table.text()).toContain('license_forbidden');
  });

  it('only shows the actions column for internal users', async () => {
    getAllMock.mockResolvedValue({data: []});

    const {wrapper} = createWrapper();
    await flushPromises();

    const headersNonInternal = (wrapper.vm as unknown as {headers: {value: string}[]}).headers;
    expect(headersNonInternal.some((h) => h.value === 'actions')).toBe(false);
  });

  it('shows the actions column when the user is internal', async () => {
    getAllMock.mockResolvedValue({data: []});

    const {wrapper, pinia} = createWrapper();
    await flushPromises();
    // isInternal is read directly off RightsUtils.rights() as a plain property (not derived from a
    // group), so flip it directly on the store's rights object.
    useUserStore(pinia).simpleProfileData.rights.isInternal = true;
    await wrapper.vm.$nextTick();

    const headers = (wrapper.vm as unknown as {headers: {value: string}[]}).headers;
    expect(headers.some((h) => h.value === 'actions')).toBe(true);
  });
});
