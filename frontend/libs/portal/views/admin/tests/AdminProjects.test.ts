// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {describe, expect, it} from 'vitest';
import AdminProjects from '../AdminProjects.vue';

describe('AdminProjects', () => {
  it('sets the expected breadcrumbs on mount', () => {
    const {pinia} = mountView(AdminProjects, {
      childStubs: {AdminGridProjects: true},
    });

    const breadcrumbs = useBreadcrumbsStore(pinia);
    expect(breadcrumbs.setCurrentBreadcrumbs).toHaveBeenCalledWith([
      ...breadcrumbs.dashboardCrumbs,
      breadcrumbs.projectsCrumb,
    ]);
  });

  it('renders without throwing and shows the projects grid', () => {
    const {wrapper} = mountView(AdminProjects, {
      childStubs: {AdminGridProjects: true},
    });

    expect(wrapper.exists()).toBe(true);
    expect(wrapper.find('admin-grid-projects-stub').exists()).toBe(true);
  });
});
