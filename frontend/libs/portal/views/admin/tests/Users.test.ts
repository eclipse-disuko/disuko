// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {describe, expect, it} from 'vitest';
import Users from '../Users.vue';

describe('Users', () => {
  it('sets the expected breadcrumbs on mount', () => {
    const {pinia} = mountView(Users, {
      childStubs: {TableUser: true},
    });

    const breadcrumbs = useBreadcrumbsStore(pinia);
    expect(breadcrumbs.setCurrentBreadcrumbs).toHaveBeenCalledWith([
      ...breadcrumbs.dashboardCrumbs,
      {title: 'Users', disabled: false, href: '/dashboard/admin/users'},
    ]);
  });

  it('renders without throwing and shows the user table', () => {
    const {wrapper} = mountView(Users, {
      childStubs: {TableUser: true},
    });

    expect(wrapper.exists()).toBe(true);
    expect(wrapper.find('table-user-stub').exists()).toBe(true);
  });
});
