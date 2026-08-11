// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {describe, expect, it} from 'vitest';
import Tasks from '../Tasks.vue';

describe('Tasks', () => {
  it('sets the expected breadcrumbs on mount', () => {
    const {pinia} = mountView(Tasks, {
      childStubs: {GridTask: true},
    });

    const breadcrumbs = useBreadcrumbsStore(pinia);
    expect(breadcrumbs.setCurrentBreadcrumbs).toHaveBeenCalledWith([
      {title: 'Dashboard', disabled: false, href: '/dashboard/home'},
      {title: 'Tasks', disabled: false, href: '/dashboard/tasks'},
    ]);
  });

  it('renders without throwing and shows the task grid', () => {
    const {wrapper} = mountView(Tasks, {
      childStubs: {GridTask: true},
    });

    expect(wrapper.exists()).toBe(true);
    const grid = wrapper.find('grid-task-stub');
    expect(grid.exists()).toBe(true);
    expect(grid.attributes('in-own-view')).toBe('true');
    expect(grid.attributes('fixed-height')).toBe('false');
  });
});
