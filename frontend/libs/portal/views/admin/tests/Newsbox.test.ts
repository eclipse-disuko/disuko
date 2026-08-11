// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useNewsboxStore} from '@disclosure-portal/stores/newsbox.store';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {flushPromises} from '@vue/test-utils';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import Newsbox from '../Newsbox.vue';

// Newsbox.vue imports {VForm} from 'vuetify/components' purely as a TypeScript type
// (`ref<VForm | null>`). That barrel import still executes at runtime and pulls in Vuetify's
// full component bundle including raw `.css` side-effect imports, which Vitest's Node-based
// module loader cannot resolve. Stub the barrel out entirely to avoid loading it.
vi.mock('vuetify/components', () => ({VForm: {}}));

const {getAllNewsboxItemsMock, createNewsboxItemMock, updateNewsboxItemMock, deleteItemsAdminMock, snackbarInfoMock} =
  vi.hoisted(() => ({
    getAllNewsboxItemsMock: vi.fn(),
    createNewsboxItemMock: vi.fn(),
    updateNewsboxItemMock: vi.fn(),
    deleteItemsAdminMock: vi.fn(),
    snackbarInfoMock: vi.fn(),
  }));

vi.mock('@disclosure-portal/services/newsbox.service', () => ({
  default: {
    getAllNewsboxItems: getAllNewsboxItemsMock,
    createNewsboxItem: createNewsboxItemMock,
    updateNewsboxItem: updateNewsboxItemMock,
    deleteItemsAdmin: deleteItemsAdminMock,
  },
}));

vi.mock('@shared/composables/useSnackbar', () => ({
  default: () => ({
    info: snackbarInfoMock,
  }),
}));

const newsItem = {
  _key: 'news-1',
  title: 'News One',
  titleDE: '',
  description: 'desc',
  descriptionDE: '',
  image: null,
  link: null,
  expiry: null,
};

const childStubs = {
  TableLayout: {template: '<div><slot name="buttons"/><slot name="table"/></div>'},
  DialogLayout: {template: '<div><slot /></div>'},
  ConfirmationDialog: true,
  TextField: true,
  TextArea: true,
  DImageUpload: true,
  'v-data-table': {
    props: ['items'],
    template: '<div class="v-data-table">{{ JSON.stringify(items) }}</div>',
  },
};

describe('Newsbox', () => {
  beforeEach(() => {
    getAllNewsboxItemsMock.mockReset();
    createNewsboxItemMock.mockReset();
    updateNewsboxItemMock.mockReset();
    deleteItemsAdminMock.mockReset();
    snackbarInfoMock.mockReset();
    getAllNewsboxItemsMock.mockResolvedValue({data: [newsItem]});
  });

  it('sets the expected breadcrumbs on mount', async () => {
    const {pinia} = mountView(Newsbox, {childStubs});
    await flushPromises();

    const breadcrumbs = useBreadcrumbsStore(pinia);
    expect(breadcrumbs.setCurrentBreadcrumbs).toHaveBeenCalledWith([
      {title: 'Dashboard', href: '/dashboard/home'},
      {title: 'Admin', href: '/dashboard/admin'},
      {title: 'Newsbox', href: '/dashboard/admin/newsbox'},
    ]);
  });

  it('renders without throwing and loads newsbox items into the table', async () => {
    const {wrapper} = mountView(Newsbox, {childStubs});
    await flushPromises();

    expect(wrapper.exists()).toBe(true);
    expect(getAllNewsboxItemsMock).toHaveBeenCalledTimes(1);
    expect(wrapper.find('.v-data-table').text()).toContain('News One');
  });

  it('creates a newsbox item and reloads the list', async () => {
    createNewsboxItemMock.mockResolvedValue({data: 'new-key'});

    const {wrapper, pinia} = mountView(Newsbox, {childStubs});
    await flushPromises();

    const vm = wrapper.vm as unknown as {
      form: {title: string; description: string};
      submit: () => Promise<void>;
    };
    vm.form.title = 'Brand New';
    vm.form.description = 'a description';
    await vm.submit();
    await flushPromises();

    expect(createNewsboxItemMock).toHaveBeenCalledWith(
      expect.objectContaining({title: 'Brand New', description: 'a description'}),
    );
    expect(snackbarInfoMock).toHaveBeenCalledWith('Newsbox item created successfully');
    expect(getAllNewsboxItemsMock).toHaveBeenCalledTimes(2);
    const newsboxStore = useNewsboxStore(pinia);
    expect(newsboxStore.adminNewsItems?.items).toEqual([newsItem]);
  });

  it('deletes a newsbox item on confirm and shows a success snackbar', async () => {
    deleteItemsAdminMock.mockResolvedValue({data: {}});

    const {wrapper} = mountView(Newsbox, {childStubs});
    await flushPromises();

    const vm = wrapper.vm as unknown as {
      showConfirmDelete: (item: typeof newsItem) => void;
      confirmVisible: boolean;
      onConfirm: () => Promise<void>;
    };
    vm.showConfirmDelete(newsItem);
    expect(vm.confirmVisible).toBe(true);

    await vm.onConfirm();
    await flushPromises();

    expect(deleteItemsAdminMock).toHaveBeenCalledWith('news-1');
    expect(snackbarInfoMock).toHaveBeenCalledWith('Newsbox item deleted successfully');
  });
});
