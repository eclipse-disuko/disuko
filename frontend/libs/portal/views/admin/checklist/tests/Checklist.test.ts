// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {describe, expect, it} from 'vitest';
import Checklist from '../Checklist.vue';

describe('Checklist', () => {
  it('sets the expected breadcrumbs on mount', () => {
    const {pinia} = mountView(Checklist, {
      childStubs: {GridChecklists: true},
    });

    const breadcrumbs = useBreadcrumbsStore(pinia);
    expect(breadcrumbs.setCurrentBreadcrumbs).toHaveBeenCalledWith([
      ...breadcrumbs.dashboardCrumbs,
      {title: 'Checklists', href: '/dashboard/admin/checklist'},
    ]);
  });

  it('renders without throwing and shows the checklists grid', () => {
    const {wrapper} = mountView(Checklist, {
      childStubs: {GridChecklists: true},
    });

    expect(wrapper.exists()).toBe(true);
    expect(wrapper.find('grid-checklists-stub').exists()).toBe(true);
  });
});
