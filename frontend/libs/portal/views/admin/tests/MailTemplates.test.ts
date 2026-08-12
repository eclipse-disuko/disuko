// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {describe, expect, it} from 'vitest';
import MailTemplates from '../MailTemplates.vue';

describe('MailTemplates', () => {
  it('sets the expected breadcrumbs on mount', () => {
    const {pinia} = mountView(MailTemplates, {
      childStubs: {GridMailTemplates: true},
    });

    const breadcrumbs = useBreadcrumbsStore(pinia);
    expect(breadcrumbs.setCurrentBreadcrumbs).toHaveBeenCalledWith([
      ...breadcrumbs.dashboardCrumbs,
      {title: 'Mail Templates', disabled: true, href: '/dashboard/admin/mailtemplates'},
    ]);
  });

  it('renders without throwing and shows the mail templates grid', () => {
    const {wrapper} = mountView(MailTemplates, {
      childStubs: {GridMailTemplates: true},
    });

    expect(wrapper.exists()).toBe(true);
    expect(wrapper.find('grid-mail-templates-stub').exists()).toBe(true);
  });
});
