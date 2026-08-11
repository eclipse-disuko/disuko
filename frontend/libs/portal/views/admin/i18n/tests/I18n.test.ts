// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {flushPromises} from '@vue/test-utils';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import I18n from '../I18n.vue';

const {getLocalesMock, pushMock} = vi.hoisted(() => ({
  getLocalesMock: vi.fn(),
  pushMock: vi.fn(),
}));

vi.mock('@disclosure-portal/services/i18n.service', () => ({
  default: {
    getLocales: getLocalesMock,
  },
}));

vi.mock('vue-router', () => ({
  useRoute: vi.fn(() => ({path: '/', params: {}})),
  useRouter: vi.fn(() => ({push: pushMock, replace: vi.fn(), go: vi.fn(), back: vi.fn(), forward: vi.fn()})),
}));

const vDataTableStub = {
  props: ['headers', 'items'],
  emits: ['click:row'],
  template:
    '<div class="v-data-table"><div v-for="row in items" :key="row.code" class="row" @click="$emit(\'click:row\', $event, {item: row})">{{ row.code }}:{{ row.keyCount }}:{{ row.missingCount }}:{{ row.extraCount }}</div></div>',
};

const childStubs = {
  'v-data-table': vDataTableStub,
  TableLayout: {template: '<div><slot name="description"/><slot name="buttons"/><slot name="table"/></div>'},
  DSearchField: true,
};

describe('I18n', () => {
  beforeEach(() => {
    getLocalesMock.mockReset();
    pushMock.mockReset();
  });

  it('fetches locales then sets the expected breadcrumbs on mount', async () => {
    getLocalesMock.mockResolvedValue({data: []});

    const {pinia} = mountView(I18n, {childStubs});
    await flushPromises();

    expect(getLocalesMock).toHaveBeenCalledTimes(1);
    const breadcrumbs = useBreadcrumbsStore(pinia);
    expect(breadcrumbs.setCurrentBreadcrumbs).toHaveBeenCalledWith([
      ...breadcrumbs.dashboardCrumbs,
      {title: 'Internationalization'},
    ]);
  });

  it('computes missing/extra key counts relative to the default locale and passes rows to the table', async () => {
    getLocalesMock.mockResolvedValue({
      data: [
        {
          localeCode: 'en',
          displayName: 'English',
          nativeName: 'English',
          isDefault: true,
          scope: 'app',
          entryCount: 100,
        },
        {
          localeCode: 'de',
          displayName: 'German',
          nativeName: 'Deutsch',
          isDefault: false,
          scope: 'app',
          entryCount: 90,
        },
        {
          localeCode: 'fr',
          displayName: 'French',
          nativeName: 'Français',
          isDefault: false,
          scope: 'app',
          entryCount: 105,
        },
      ],
    });

    const {wrapper} = mountView(I18n, {childStubs});
    await flushPromises();

    const text = wrapper.find('.v-data-table').text();
    expect(text).toContain('en:100:0:0');
    expect(text).toContain('de:90:10:0');
    expect(text).toContain('fr:105:0:5');
  });

  it('navigates to the locale details page when a row is clicked', async () => {
    getLocalesMock.mockResolvedValue({
      data: [
        {
          localeCode: 'en',
          displayName: 'English',
          nativeName: 'English',
          isDefault: true,
          scope: 'app',
          entryCount: 10,
        },
      ],
    });

    const {wrapper} = mountView(I18n, {childStubs});
    await flushPromises();

    await wrapper.find('.row').trigger('click');

    expect(pushMock).toHaveBeenCalledWith({name: 'I18nLocaleDetails', params: {localeCode: 'en'}});
  });
});
