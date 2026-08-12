// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {nextTick} from 'vue';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import AdminProfile from '../AdminProfile.vue';

const {getUserMock} = vi.hoisted(() => ({
  getUserMock: vi.fn(),
}));

vi.mock('@disclosure-portal/services/admin', () => ({
  default: {
    getUser: getUserMock,
  },
}));

describe('AdminProfile', () => {
  beforeEach(() => {
    getUserMock.mockReset();
  });

  it('sets the expected breadcrumbs after loading the user profile', async () => {
    getUserMock.mockResolvedValue({data: {_key: 'user-1', user: 'jdoe'}});

    const {pinia} = mountView(AdminProfile, {
      childStubs: {UserMain: true},
    });
    await nextTick();
    await nextTick();

    const breadcrumbs = useBreadcrumbsStore(pinia);
    expect(breadcrumbs.setCurrentBreadcrumbs).toHaveBeenCalledWith([
      ...breadcrumbs.dashboardCrumbs,
      {title: 'Users', disabled: false, href: '/dashboard/admin/users'},
      {title: 'jdoe', disabled: false, href: `/dashboard/admin/users/${encodeURIComponent('user-1')}`},
    ]);
  });

  it('passes the loaded user profile to UserMain', async () => {
    getUserMock.mockResolvedValue({data: {_key: 'user-1', user: 'jdoe'}});

    const {wrapper} = mountView(AdminProfile, {
      childStubs: {UserMain: true},
    });
    await nextTick();
    await nextTick();

    expect(getUserMock).toHaveBeenCalledTimes(1);
    const userMain = wrapper.find('user-main-stub');
    expect(userMain.exists()).toBe(true);
    expect(userMain.attributes('has-users-access')).toBe('true');
  });
});
