// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {describe, expect, it} from 'vitest';
import CustomIds from '../CustomIds.vue';

describe('CustomIds', () => {
  it('sets the expected breadcrumbs on mount', () => {
    const {pinia} = mountView(CustomIds, {
      childStubs: {GridCustomIds: true},
    });

    const breadcrumbs = useBreadcrumbsStore(pinia);
    expect(breadcrumbs.setCurrentBreadcrumbs).toHaveBeenCalledWith([
      ...breadcrumbs.dashboardCrumbs,
      {title: 'Custom IDs', href: '/dashboard/customids/'},
    ]);
  });

  it('renders without throwing and shows the custom ids grid', () => {
    const {wrapper} = mountView(CustomIds, {
      childStubs: {GridCustomIds: true},
    });

    expect(wrapper.exists()).toBe(true);
    expect(wrapper.find('grid-custom-ids-stub').exists()).toBe(true);
  });
});
