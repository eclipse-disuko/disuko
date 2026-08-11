// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {Group, Rights} from '@shared/user/models/Rights';
import {useUserStore} from '@shared/user/stores/user.store';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {nextTick} from 'vue';
import {describe, expect, it} from 'vitest';
import UserProfile from '../UserProfile.vue';

describe('UserProfile', () => {
  it('sets the expected breadcrumbs on mount', async () => {
    const {pinia} = mountView(UserProfile, {
      childStubs: {UserMain: true},
      beforePiniaMount: (p) => {
        useUserStore(p).setSimpleProfileData({
          rights: {groups: []} as unknown as Rights,
          profile: {} as never,
          allowed: true,
        });
      },
    });
    await nextTick();

    const breadcrumbs = useBreadcrumbsStore(pinia);
    expect(breadcrumbs.setCurrentBreadcrumbs).toHaveBeenCalledWith([
      ...breadcrumbs.dashboardCrumbs,
      {title: 'Profile', disabled: false, href: '/dashboard/user'},
    ]);
  });

  it('passes the current user profile to UserMain and marks roles switchable for domain admins', async () => {
    const {wrapper, pinia} = mountView(UserProfile, {
      childStubs: {UserMain: true},
      beforePiniaMount: (p) => {
        useUserStore(p).setSimpleProfileData({
          rights: {groups: [Group.UserDomainAdmin]} as unknown as Rights,
          profile: {user: 'jdoe'} as never,
          allowed: true,
        });
      },
    });
    await nextTick();

    const userMain = wrapper.findComponent({name: 'UserMain'});
    expect(userMain.props('userProfile')).toEqual(useUserStore(pinia).getProfile);
    expect(userMain.props('rolesSwitchable')).toBe(true);
    expect(userMain.props('hasUsersAccess')).toBe(false);
  });
});
