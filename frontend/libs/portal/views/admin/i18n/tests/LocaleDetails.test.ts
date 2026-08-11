// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {flushPromises} from '@vue/test-utils';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import LocaleDetails from '../LocaleDetails.vue';

const {
  getLocaleMock,
  upsertTranslationMock,
  deleteTranslationMock,
  exportLocaleMock,
  getLocalesMock,
  importLocaleMock,
  replaceMock,
  snackbarInfoMock,
  snackbarErrorMock,
} = vi.hoisted(() => ({
  getLocaleMock: vi.fn(),
  upsertTranslationMock: vi.fn(),
  deleteTranslationMock: vi.fn(),
  exportLocaleMock: vi.fn(),
  getLocalesMock: vi.fn(),
  importLocaleMock: vi.fn(),
  replaceMock: vi.fn(),
  snackbarInfoMock: vi.fn(),
  snackbarErrorMock: vi.fn(),
}));

vi.mock('@disclosure-portal/services/i18n.service', () => ({
  default: {
    getLocale: getLocaleMock,
    upsertTranslation: upsertTranslationMock,
    deleteTranslation: deleteTranslationMock,
    exportLocale: exportLocaleMock,
    getLocales: getLocalesMock,
    importLocale: importLocaleMock,
  },
}));

vi.mock('@shared/composables/useSnackbar', () => ({
  default: () => ({
    info: snackbarInfoMock,
    error: snackbarErrorMock,
  }),
}));

vi.mock('vue-router', () => ({
  useRoute: vi.fn(() => ({path: '/dashboard/admin/i18n/en', params: {localeCode: 'en'}})),
  useRouter: vi.fn(() => ({push: vi.fn(), replace: replaceMock, go: vi.fn(), back: vi.fn(), forward: vi.fn()})),
}));

const childStubs = {
  TableLayout: {template: '<div><slot name="description"/><slot name="buttons"/><slot name="table"/></div>'},
  DSearchField: true,
  DCActionButton: {template: '<button type="button" @click="$emit(\'clicked\')"><slot /></button>'},
  DIconButton: {template: '<button type="button" @click="$emit(\'clicked\')"><slot /></button>'},
  DialogLayout: {
    template:
      '<div><slot /><slot name="right" /><button class="primary" @click="$emit(\'primary-action\')">primary</button><button class="secondary" @click="$emit(\'secondary-action\')">secondary</button></div>',
  },
  Stack: {template: '<div><slot /></div>'},
  'v-data-table': {
    props: ['headers', 'items'],
    template: '<div class="v-data-table">{{ JSON.stringify(items) }}</div>',
  },
};

const localeResponse = {
  localeCode: 'en',
  nativeName: 'English',
  entries: {GREETING: 'Hello', FAREWELL: 'Bye'},
};

const createWrapper = () => mountView(LocaleDetails, {childStubs});

describe('LocaleDetails', () => {
  beforeEach(() => {
    getLocaleMock.mockReset();
    upsertTranslationMock.mockReset();
    deleteTranslationMock.mockReset();
    exportLocaleMock.mockReset();
    getLocalesMock.mockReset();
    importLocaleMock.mockReset();
    replaceMock.mockReset();
    snackbarInfoMock.mockReset();
    snackbarErrorMock.mockReset();
  });

  it('fetches the locale on mount, populates the table, and sets breadcrumbs with the locale label', async () => {
    getLocaleMock.mockResolvedValue({data: localeResponse});

    const {wrapper, pinia} = createWrapper();
    await flushPromises();

    expect(getLocaleMock).toHaveBeenCalledWith('en');
    expect(wrapper.find('.v-data-table').text()).toContain('GREETING');
    expect(wrapper.find('.v-data-table').text()).toContain('Hello');

    const breadcrumbs = useBreadcrumbsStore(pinia);
    expect(breadcrumbs.setCurrentBreadcrumbs).toHaveBeenCalledWith([
      ...breadcrumbs.dashboardCrumbs,
      {title: 'Internationalization', disabled: false, href: '/dashboard/admin/i18n'},
      {title: 'English'},
    ]);
  });

  it('shows an error and clears entries when fetching the locale fails', async () => {
    getLocaleMock.mockRejectedValue(new Error('boom'));

    const {wrapper} = createWrapper();
    await flushPromises();

    expect((wrapper.vm as unknown as {actionError: string}).actionError).toBe('ERROR_500_TITLE');
    expect(wrapper.find('.v-data-table').text()).toBe('[]');
  });

  it('adds a new entry and shows a success snackbar', async () => {
    getLocaleMock.mockResolvedValue({data: localeResponse});
    upsertTranslationMock.mockResolvedValue({});

    const {wrapper} = createWrapper();
    await flushPromises();

    const vm = wrapper.vm as unknown as {
      newEntryKey: string;
      newEntryTranslation: string;
      addEntry: () => Promise<void>;
      entries: {key: string; translation: string}[];
    };
    vm.newEntryKey = 'NEW_KEY';
    vm.newEntryTranslation = 'New Value';
    await vm.addEntry();
    await flushPromises();

    expect(upsertTranslationMock).toHaveBeenCalledWith('en', 'NEW_KEY', 'New Value');
    expect(snackbarInfoMock).toHaveBeenCalledWith('Entry saved successfully.');
    expect(vm.entries.some((entry) => entry.key === 'NEW_KEY')).toBe(true);
  });

  it('shows an error snackbar when adding an entry fails', async () => {
    getLocaleMock.mockResolvedValue({data: localeResponse});
    upsertTranslationMock.mockRejectedValue(new Error('boom'));

    const {wrapper} = createWrapper();
    await flushPromises();

    const vm = wrapper.vm as unknown as {
      newEntryKey: string;
      newEntryTranslation: string;
      addEntry: () => Promise<void>;
    };
    vm.newEntryKey = 'NEW_KEY';
    vm.newEntryTranslation = 'New Value';
    await vm.addEntry();
    await flushPromises();

    expect(snackbarErrorMock).toHaveBeenCalledWith('ERROR_500_TITLE');
  });

  it('saves an edited translation and shows a success snackbar', async () => {
    getLocaleMock.mockResolvedValue({data: localeResponse});
    upsertTranslationMock.mockResolvedValue({});

    const {wrapper} = createWrapper();
    await flushPromises();

    const vm = wrapper.vm as unknown as {
      editRowKey: string | null;
      draftTranslation: string;
      saveEdit: () => Promise<void>;
      entries: {key: string; translation: string}[];
    };
    vm.editRowKey = 'GREETING';
    vm.draftTranslation = 'Hi there';
    await vm.saveEdit();
    await flushPromises();

    expect(upsertTranslationMock).toHaveBeenCalledWith('en', 'GREETING', 'Hi there');
    expect(snackbarInfoMock).toHaveBeenCalledWith('Entry saved successfully.');
    expect(vm.entries.find((entry) => entry.key === 'GREETING')?.translation).toBe('Hi there');
  });

  it('shows an error snackbar when saving an edited translation fails', async () => {
    getLocaleMock.mockResolvedValue({data: localeResponse});
    upsertTranslationMock.mockRejectedValue(new Error('boom'));

    const {wrapper} = createWrapper();
    await flushPromises();

    const vm = wrapper.vm as unknown as {
      editRowKey: string | null;
      draftTranslation: string;
      saveEdit: () => Promise<void>;
    };
    vm.editRowKey = 'GREETING';
    vm.draftTranslation = 'Hi there';
    await vm.saveEdit();
    await flushPromises();

    expect(snackbarErrorMock).toHaveBeenCalledWith('ERROR_500_TITLE');
  });

  it('deletes an entry and shows a success snackbar', async () => {
    getLocaleMock.mockResolvedValue({data: localeResponse});
    deleteTranslationMock.mockResolvedValue({});

    const {wrapper} = createWrapper();
    await flushPromises();

    const vm = wrapper.vm as unknown as {
      deleteEntryKey: string | null;
      deleteGlobally: boolean;
      onDeleteConfirm: () => Promise<void>;
      entries: {key: string; translation: string}[];
    };
    vm.deleteEntryKey = 'GREETING';
    vm.deleteGlobally = false;
    await vm.onDeleteConfirm();
    await flushPromises();

    expect(deleteTranslationMock).toHaveBeenCalledWith('en', 'GREETING');
    expect(snackbarInfoMock).toHaveBeenCalledWith('Entry deleted successfully.');
    expect(vm.entries.some((entry) => entry.key === 'GREETING')).toBe(false);
  });

  it('sets an error message when deleting the current locale entry fails', async () => {
    getLocaleMock.mockResolvedValue({data: localeResponse});
    deleteTranslationMock.mockRejectedValue({response: {data: {message: 'Delete failed'}}});

    const {wrapper} = createWrapper();
    await flushPromises();

    const vm = wrapper.vm as unknown as {
      deleteEntryKey: string | null;
      deleteGlobally: boolean;
      onDeleteConfirm: () => Promise<void>;
      actionError: string;
    };
    vm.deleteEntryKey = 'GREETING';
    vm.deleteGlobally = false;
    await vm.onDeleteConfirm();
    await flushPromises();

    expect(vm.actionError).toBe('Delete failed');
    expect(snackbarInfoMock).not.toHaveBeenCalled();
  });

  it('exports the locale as JSON and shows a success snackbar', async () => {
    getLocaleMock.mockResolvedValue({data: localeResponse});
    exportLocaleMock.mockResolvedValue({data: new Blob(['{}'])});

    const createObjectURLSpy = vi.fn(() => 'blob:mock-url');
    const revokeObjectURLSpy = vi.fn();
    URL.createObjectURL = createObjectURLSpy;
    URL.revokeObjectURL = revokeObjectURLSpy;

    const {wrapper} = createWrapper();
    await flushPromises();

    const vm = wrapper.vm as unknown as {exportAsJson: () => Promise<void>};
    await vm.exportAsJson();
    await flushPromises();

    expect(exportLocaleMock).toHaveBeenCalledWith('en');
    expect(createObjectURLSpy).toHaveBeenCalled();
    expect(revokeObjectURLSpy).toHaveBeenCalled();
    expect(snackbarInfoMock).toHaveBeenCalledWith('Export completed successfully.');
  });

  it('shows an error snackbar when exporting the locale fails', async () => {
    getLocaleMock.mockResolvedValue({data: localeResponse});
    exportLocaleMock.mockRejectedValue(new Error('boom'));

    const {wrapper} = createWrapper();
    await flushPromises();

    const vm = wrapper.vm as unknown as {exportAsJson: () => Promise<void>};
    await vm.exportAsJson();
    await flushPromises();

    expect(snackbarErrorMock).toHaveBeenCalledWith('ERROR_500_TITLE');
  });

  it('redirects to the locale list instead of fetching when there is no locale code in the route', async () => {
    const routerModule = await import('vue-router');
    vi.mocked(routerModule.useRoute).mockReturnValueOnce({path: '/dashboard/admin/i18n', params: {}} as never);

    createWrapper();
    await flushPromises();

    expect(replaceMock).toHaveBeenCalledWith({name: 'I18nAdmin'});
    expect(getLocaleMock).not.toHaveBeenCalled();
  });
});
