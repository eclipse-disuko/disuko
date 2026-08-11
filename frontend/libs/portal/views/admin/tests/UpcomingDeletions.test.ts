// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {describe, expect, it} from 'vitest';
import UpcomingDeletions from '../UpcomingDeletions.vue';

describe('UpcomingDeletions', () => {
  it('sets the expected breadcrumbs on mount', () => {
    const {pinia} = mountView(UpcomingDeletions, {
      childStubs: {GridUpcomingDeletions: true},
    });

    const breadcrumbs = useBreadcrumbsStore(pinia);
    expect(breadcrumbs.setCurrentBreadcrumbs).toHaveBeenCalledWith([
      ...breadcrumbs.dashboardCrumbs,
      {title: 'Upcoming Deletions', disabled: true},
    ]);
  });

  it('renders without throwing and shows the upcoming deletions grid', () => {
    const {wrapper} = mountView(UpcomingDeletions, {
      childStubs: {GridUpcomingDeletions: true},
    });

    expect(wrapper.exists()).toBe(true);
    expect(wrapper.find('grid-upcoming-deletions-stub').exists()).toBe(true);
  });
});
