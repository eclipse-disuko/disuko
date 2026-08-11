// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {Rights} from '@shared/user/models/Rights';
import {useUserStore} from '@shared/user/stores/user.store';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {flushPromises} from '@vue/test-utils';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import Schemas from '../Schemas.vue';

const {getAllSchemasMock, getSchemaLabelsMock, createSchemaMock, snackbarInfoMock} = vi.hoisted(() => ({
  getAllSchemasMock: vi.fn(),
  getSchemaLabelsMock: vi.fn(),
  createSchemaMock: vi.fn(),
  snackbarInfoMock: vi.fn(),
}));

vi.mock('@disclosure-portal/services/admin', () => ({
  default: {
    getAllSchemas: getAllSchemasMock,
    getSchemaLabels: getSchemaLabelsMock,
    createSchema: createSchemaMock,
  },
}));

vi.mock('@shared/composables/useSnackbar', () => ({
  default: () => ({
    info: snackbarInfoMock,
  }),
}));

const schemas = [
  {
    _key: 'schema-1',
    name: 'My Schema',
    version: '1.0',
    description: 'desc',
    type: 0,
    label: 'label-1',
    active: true,
    created: '2026-01-01T10:00:00Z',
    updated: '2026-01-02T10:00:00Z',
  },
];

const childStubs = {
  TableLayout: {template: '<div><slot name="buttons"/><slot name="table"/></div>'},
  NewSchemaDialog: true,
  'v-data-table': {
    props: ['items'],
    template: '<div class="v-data-table">{{ JSON.stringify(items) }}</div>',
  },
};

describe('Schemas', () => {
  beforeEach(() => {
    getAllSchemasMock.mockReset();
    getSchemaLabelsMock.mockReset();
    createSchemaMock.mockReset();
    snackbarInfoMock.mockReset();
    getAllSchemasMock.mockResolvedValue({data: schemas});
    getSchemaLabelsMock.mockResolvedValue({data: [{_key: 'label-1', name: 'Label One', description: 'desc'}]});
  });

  it('sets the expected breadcrumbs on mount', async () => {
    const {pinia} = mountView(Schemas, {childStubs});
    await flushPromises();

    const breadcrumbs = useBreadcrumbsStore(pinia);
    expect(breadcrumbs.setCurrentBreadcrumbs).toHaveBeenCalledWith([
      {title: 'Dashboard', href: '/dashboard/home'},
      {title: 'Admin', href: '/dashboard/admin'},
      {title: 'SBOM Schemas', href: '/dashboard/admin/schemas/'},
    ]);
  });

  it('renders without throwing and loads schemas into the table', async () => {
    const {wrapper} = mountView(Schemas, {childStubs});
    await flushPromises();

    expect(wrapper.exists()).toBe(true);
    expect(getAllSchemasMock).toHaveBeenCalledTimes(1);
    expect(wrapper.find('.v-data-table').text()).toContain('My Schema');
  });

  it('creates a schema and shows a success snackbar, then reloads the list', async () => {
    createSchemaMock.mockResolvedValue({data: {}});

    const {wrapper} = mountView(Schemas, {
      childStubs,
      beforePiniaMount: (p) => {
        useUserStore(p).setSimpleProfileData({
          rights: {allowSchema: {create: true}} as unknown as Rights,
          profile: {} as never,
          allowed: true,
        });
      },
    });
    await flushPromises();

    const formData = new FormData();
    await (wrapper.vm as unknown as {createSchema: (fd: FormData) => Promise<void>}).createSchema(formData);
    await flushPromises();

    expect(createSchemaMock).toHaveBeenCalledWith(formData);
    expect(snackbarInfoMock).toHaveBeenCalledWith('SBOM schema created');
    expect(getAllSchemasMock).toHaveBeenCalledTimes(2);
  });
});
