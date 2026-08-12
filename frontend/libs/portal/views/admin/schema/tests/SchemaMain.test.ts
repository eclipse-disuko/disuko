// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {Rights} from '@shared/user/models/Rights';
import {useUserStore} from '@shared/user/stores/user.store';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {flushPromises} from '@vue/test-utils';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import {useRoute} from 'vue-router';
import SchemaMain from '../SchemaMain.vue';

const {getSchemaMock, getSchemaLabelsMock, activateMock, downloadSchemaMock} = vi.hoisted(() => ({
  getSchemaMock: vi.fn(),
  getSchemaLabelsMock: vi.fn(),
  activateMock: vi.fn(),
  downloadSchemaMock: vi.fn(),
}));

vi.mock('@disclosure-portal/services/admin', () => ({
  default: {
    getSchema: getSchemaMock,
    getSchemaLabels: getSchemaLabelsMock,
    activate: activateMock,
    downloadSchema: downloadSchemaMock,
  },
}));

const schema = {
  _key: 'schema-1',
  name: 'My Schema',
  version: '1.0',
  description: 'A schema',
  content: '{"a":1}',
  type: 0,
  label: 'label-1',
  active: false,
  created: '2026-01-01T10:00:00Z',
  updated: '2026-01-02T10:00:00Z',
};

const childStubs = {
  'json-viewer': true,
};

describe('SchemaMain', () => {
  beforeEach(() => {
    getSchemaMock.mockReset();
    getSchemaLabelsMock.mockReset();
    activateMock.mockReset();
    downloadSchemaMock.mockReset();
    getSchemaMock.mockResolvedValue({data: {...schema}});
    getSchemaLabelsMock.mockResolvedValue({data: [{_key: 'label-1', name: 'Label One', description: 'desc'}]});
    // The globally mocked useRoute() has no route params by default; give it the schema id
    // this view reads via `route.params.id` on mount.
    vi.mocked(useRoute).mockReturnValue({
      path: '/dashboard/schemas/schema-1',
      name: 'schema',
      params: {id: 'schema-1'},
      query: {},
      hash: '',
      fullPath: '/dashboard/schemas/schema-1',
      matched: [],
      meta: {},
      redirectedFrom: undefined,
    } as never);
  });

  it('renders without throwing', async () => {
    const {wrapper} = mountView(SchemaMain, {childStubs});
    await flushPromises();

    expect(wrapper.exists()).toBe(true);
  });

  it('loads the schema and labels on mount, and sets the read-only breadcrumbs', async () => {
    const {wrapper, pinia} = mountView(SchemaMain, {
      childStubs,
      beforePiniaMount: (p) => {
        useUserStore(p).setSimpleProfileData({
          rights: {allowSchema: {read: true}} as unknown as Rights,
          profile: {} as never,
          allowed: true,
        });
      },
    });
    await flushPromises();

    expect(getSchemaMock).toHaveBeenCalledTimes(1);
    expect(getSchemaLabelsMock).toHaveBeenCalledTimes(1);
    expect(wrapper.text()).toContain('My Schema');

    // The globally mocked useRoute() resolves to path '/', so the view falls back to its
    // read-only breadcrumb branch (no admin/schemas listing crumb).
    const breadcrumbs = useBreadcrumbsStore(pinia);
    expect(breadcrumbs.setCurrentBreadcrumbs).toHaveBeenCalledWith([
      {title: 'Dashboard', disabled: false, href: '/dashboard/home'},
      {title: 'My Schema', disabled: false, href: '/dashboard/admin/schemas/schema-1'},
    ]);
  });

  it('activates the schema and calls the service with the schema id', async () => {
    activateMock.mockResolvedValue({data: {success: false}});

    const {wrapper} = mountView(SchemaMain, {childStubs});
    await flushPromises();

    await (wrapper.vm as unknown as {activate: () => Promise<void>}).activate();
    await flushPromises();

    expect(activateMock).toHaveBeenCalledWith('schema-1');
  });

  it('downloads the schema content', async () => {
    downloadSchemaMock.mockResolvedValue({data: '{"a":1}'});

    const {wrapper} = mountView(SchemaMain, {childStubs});
    await flushPromises();

    await (wrapper.vm as unknown as {download: () => Promise<void>}).download();
    await flushPromises();

    expect(downloadSchemaMock).toHaveBeenCalledWith('schema-1');
  });
});
