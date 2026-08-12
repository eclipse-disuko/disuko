// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {createTestingPinia} from '@pinia/testing';
import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, beforeEach, describe, expect, it, vi} from 'vitest';
import {nextTick} from 'vue';
import DBreadcrumb from '../DBreadcrumb.vue';

const VBreadcrumbsStub = {
  props: ['items', 'color'],
  template: '<div class="v-breadcrumbs"><slot name="divider" /><slot /></div>',
};
const VIconStub = {props: ['icon', 'size', 'color'], template: '<i class="v-icon"></i>'};

describe('DBreadcrumb', () => {
  const originalTMock = config.global.mocks?.$t;

  beforeAll(() => {
    if (config.global.mocks && '$t' in config.global.mocks) {
      delete config.global.mocks.$t;
    }
  });

  afterAll(() => {
    if (originalTMock) {
      config.global.mocks = {
        ...config.global.mocks,
        $t: originalTMock,
      };
    }
  });

  beforeEach(() => {
    vi.clearAllMocks();
  });

  const createWrapper = () => {
    return mount(DBreadcrumb, {
      global: {
        plugins: [createTestingPinia({createSpy: vi.fn, stubActions: false})],
        stubs: {'v-breadcrumbs': VBreadcrumbsStub, 'v-icon': VIconStub},
      },
    });
  };

  it('renders no breadcrumb items when the store is empty', () => {
    const wrapper = createWrapper();

    expect(wrapper.findComponent(VBreadcrumbsStub).props('items')).toEqual([]);
  });

  it('maps store breadcrumbs to items, marking the last one as disabled', async () => {
    const wrapper = createWrapper();
    const store = useBreadcrumbsStore();
    store.currentBreadcrumbs = [
      {title: 'Dashboard', href: '/dashboard'},
      {title: 'Projects', href: '/dashboard/projects'},
    ];
    await nextTick();

    const items = wrapper.findComponent(VBreadcrumbsStub).props('items');
    expect(items).toEqual([
      {title: 'Dashboard', href: '/dashboard', to: '/dashboard'},
      {title: 'Projects', href: '/dashboard/projects', to: '/dashboard/projects', disabled: true},
    ]);
  });

  it('falls back to an empty string for `to` when href is missing', async () => {
    const wrapper = createWrapper();
    const store = useBreadcrumbsStore();
    store.currentBreadcrumbs = [{title: 'No Link'}];
    await nextTick();

    const items = wrapper.findComponent(VBreadcrumbsStub).props('items');
    expect(items).toEqual([{title: 'No Link', to: '', disabled: true}]);
  });
});
