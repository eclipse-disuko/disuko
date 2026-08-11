// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useUserStore} from '@shared/user/stores/user.store';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {flushPromises} from '@vue/test-utils';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import Licenses from '../Licenses.vue';

const {
  searchMock,
  deleteMock,
  getMock,
  getCountOfPolicyRuleUsingThisLicenceMock,
  getAllObligationsMock,
  getFilterSetsMock,
  getFilterSetMock,
  pushMock,
  replaceMock,
} = vi.hoisted(() => ({
  searchMock: vi.fn(),
  deleteMock: vi.fn(),
  getMock: vi.fn(),
  getCountOfPolicyRuleUsingThisLicenceMock: vi.fn(),
  getAllObligationsMock: vi.fn(),
  getFilterSetsMock: vi.fn(),
  getFilterSetMock: vi.fn(),
  pushMock: vi.fn(),
  replaceMock: vi.fn(),
}));

vi.mock('@disclosure-portal/services/license', () => ({
  default: {
    search: searchMock,
    delete: deleteMock,
    get: getMock,
    getCountOfPolicyRuleUsingThisLicence: getCountOfPolicyRuleUsingThisLicenceMock,
  },
}));

vi.mock('@disclosure-portal/services/admin', () => ({
  default: {
    getAllObligations: getAllObligationsMock,
  },
}));

vi.mock('@disclosure-portal/services/filtersets', () => ({
  default: {
    getFilterSets: getFilterSetsMock,
    getFilterSet: getFilterSetMock,
  },
}));

const {snackbarInfoMock, snackbarErrorMock} = vi.hoisted(() => ({
  snackbarInfoMock: vi.fn(),
  snackbarErrorMock: vi.fn(),
}));

vi.mock('@shared/composables/useSnackbar', () => ({
  default: () => ({
    info: snackbarInfoMock,
    error: snackbarErrorMock,
  }),
}));

// useViewTools() pulls in useLanguageStore(), which calls Vuetify's useLocale() composable -
// that requires a real Vuetify plugin instance to be installed, which these view tests don't set up.
// Stub the composable directly (keeping the other named exports of the module intact) so mounting
// doesn't need a full Vuetify instance.
vi.mock('@disclosure-portal/utils/View', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@disclosure-portal/utils/View')>();
  return {
    ...actual,
    default: () => ({
      getNameForLanguage: (obligation: {name?: string} | null) => obligation?.name ?? '',
      getDescriptionForLanguage: (obligation: {description?: string} | null) => obligation?.description ?? '',
      gridPolicyRulesAssignmentsHeaderClassByLanguage: () => '',
      gridPolicyRulesAssignmentsRowClassByLanguage: () => '',
    }),
  };
});

vi.mock('vue-router', () => ({
  useRoute: vi.fn(() => ({path: '/dashboard/licenses', params: {}})),
  useRouter: vi.fn(() => ({
    push: pushMock,
    replace: replaceMock,
    currentRoute: {value: {path: '/dashboard/licenses'}},
    go: vi.fn(),
    back: vi.fn(),
    forward: vi.fn(),
  })),
}));

const newOrEditLicenseDialogStub = {
  props: ['mode', 'initialData'],
  emits: ['closed:successfully'],
  template: '<div><slot :showDialog="noop" /></div>',
  methods: {
    showDialog: vi.fn(),
    noop() {},
  },
};

const childStubs = {
  TableLayout: {template: '<div><slot name="description"/><slot name="buttons"/><slot name="table"/></div>'},
  NewOrEditLicenseDialog: newOrEditLicenseDialogStub,
  LicenseCompareDialog: true,
  DCActionButton: true,
  DSearchField: true,
  FilterSets: true,
  ConfirmationDialog: true,
  ClassificationsPerLicenseDialog: true,
  ConfigurePoliciesForLicenseDialog: true,
  'v-data-table-server': {
    props: ['headers', 'items', 'itemsLength', 'loading'],
    template: '<div class="v-data-table-server">{{ JSON.stringify(items) }}</div>',
  },
};

const fullRights = {
  allowLicense: {read: true, create: true, update: true, delete: true},
  allowPolicy: {read: true, create: true, update: true, delete: true},
  allowObligation: {read: true, create: false, update: false, delete: false},
  groups: [],
};

const readOnlyRights = {
  allowLicense: {read: false, create: false, update: false, delete: false},
  allowPolicy: {read: false, create: false, update: false, delete: false},
  allowObligation: {read: false, create: false, update: false, delete: false},
  groups: [],
};

const emptyMeta = {
  possibleCharts: {},
  possibleSources: {},
  possibleFamilies: {},
  possibleApproval: {},
  possibleType: {},
  possibleClassifications: [],
};

const licenseSlim = {licenseId: 'MIT', name: 'MIT License', source: 'spdx', meta: {}, aliases: []};

const createWrapper = (rights: unknown = fullRights) =>
  mountView(Licenses, {
    childStubs,
    beforePiniaMount: (pinia) => {
      const userStore = useUserStore(pinia);
      userStore.setSimpleProfileData({rights, profile: {}, allowed: true} as never);
    },
  });

describe('Licenses', () => {
  beforeEach(() => {
    searchMock.mockReset();
    deleteMock.mockReset();
    getMock.mockReset();
    getCountOfPolicyRuleUsingThisLicenceMock.mockReset();
    getAllObligationsMock.mockReset();
    getFilterSetsMock.mockReset();
    getFilterSetMock.mockReset();
    pushMock.mockReset();
    replaceMock.mockReset();
    snackbarInfoMock.mockReset();
    snackbarErrorMock.mockReset();

    getAllObligationsMock.mockResolvedValue({data: {items: []}});
    getFilterSetsMock.mockResolvedValue([]);
    searchMock.mockResolvedValue({data: {licenses: [licenseSlim], count: 1, meta: emptyMeta}});
  });

  it('redirects to home and skips fetching when the user has no license read rights', async () => {
    createWrapper(readOnlyRights);
    await flushPromises();

    expect(replaceMock).toHaveBeenCalledWith({path: '/dashboard/home'});
    expect(searchMock).not.toHaveBeenCalled();
  });

  it('fetches classifications, filter sets, and licenses on mount, sets breadcrumbs and populates the table', async () => {
    const {wrapper, pinia} = createWrapper();
    await flushPromises();

    expect(getAllObligationsMock).toHaveBeenCalledTimes(1);
    expect(getFilterSetsMock).toHaveBeenCalledWith('licenses');
    expect(searchMock).toHaveBeenCalledTimes(1);
    expect(wrapper.find('.v-data-table-server').text()).toContain('MIT License');

    const breadcrumbs = useBreadcrumbsStore(pinia);
    expect(breadcrumbs.setCurrentBreadcrumbs).toHaveBeenCalledWith([
      {title: 'Dashboard', href: '/dashboard/home'},
      {title: 'Licenses', href: '/dashboard/licenses/'},
    ]);
  });

  it('shows action buttons according to the current user rights, both ways', async () => {
    const {wrapper} = createWrapper(fullRights);
    await flushPromises();

    const vmWithRights = wrapper.vm as unknown as {
      getActionButtons: (item: typeof licenseSlim) => {event: string; show: boolean}[];
    };
    const buttonsWithFullRights = vmWithRights.getActionButtons(licenseSlim);
    expect(buttonsWithFullRights.find((b) => b.event === 'edit')?.show).toBe(true);
    expect(buttonsWithFullRights.find((b) => b.event === 'duplicate')?.show).toBe(true);
    // spdx-sourced licenses can never be deleted, regardless of rights
    expect(buttonsWithFullRights.find((b) => b.event === 'delete')?.show).toBe(false);

    const {wrapper: readOnlyWrapper} = createWrapper(readOnlyRights);
    await flushPromises();
    const vmReadOnly = readOnlyWrapper.vm as unknown as {
      getActionButtons: (item: typeof licenseSlim) => {event: string; show: boolean}[];
    };
    const buttonsReadOnly = vmReadOnly.getActionButtons({...licenseSlim, source: 'manual'});
    expect(buttonsReadOnly.every((b) => !b.show)).toBe(true);
  });

  it('deletes a license after confirmation and shows a success snackbar', async () => {
    getCountOfPolicyRuleUsingThisLicenceMock.mockResolvedValue({data: {count: 0}});
    deleteMock.mockResolvedValue({});

    const {wrapper} = createWrapper();
    await flushPromises();
    searchMock.mockClear();

    const vm = wrapper.vm as unknown as {
      showDeletionConfirmationDialog: (license: typeof licenseSlim) => Promise<void>;
      confirmConfig: {type: string; key: string; okButtonIsDisabled?: boolean};
      onConfirm: (config: {type: string; key: string; okButtonIsDisabled?: boolean}) => Promise<void>;
    };

    await vm.showDeletionConfirmationDialog({...licenseSlim, source: 'manual'});
    expect(vm.confirmConfig.type).toBe('DELETE');
    expect(vm.confirmConfig.key).toBe('MIT');

    await vm.onConfirm(vm.confirmConfig);
    await flushPromises();

    expect(deleteMock).toHaveBeenCalledWith('MIT');
    expect(snackbarInfoMock).toHaveBeenCalledWith('license deleted');
    expect(searchMock).toHaveBeenCalledTimes(1);
  });

  it('blocks deletion and disables the confirm button when the license is still in use', async () => {
    getCountOfPolicyRuleUsingThisLicenceMock.mockResolvedValue({data: {count: 3}});

    const {wrapper} = createWrapper();
    await flushPromises();

    const vm = wrapper.vm as unknown as {
      showDeletionConfirmationDialog: (license: typeof licenseSlim) => Promise<void>;
      confirmConfig: {okButtonIsDisabled?: boolean};
      onConfirm: (config: {okButtonIsDisabled?: boolean}) => Promise<void>;
    };

    await vm.showDeletionConfirmationDialog({...licenseSlim, source: 'manual'});
    expect(vm.confirmConfig.okButtonIsDisabled).toBe(true);

    await vm.onConfirm(vm.confirmConfig);
    await flushPromises();

    expect(deleteMock).not.toHaveBeenCalled();
  });

  it('loads the full license and opens the dialog in edit mode', async () => {
    const fullLicense = {...licenseSlim, description: 'A permissive license'};
    getMock.mockResolvedValue({data: fullLicense});

    const {wrapper} = createWrapper();
    await flushPromises();

    const vm = wrapper.vm as unknown as {
      editLicense: (license: typeof licenseSlim) => Promise<void>;
      currentLicenseForAction: typeof fullLicense | null;
      licenseDialogMode: string;
    };
    await vm.editLicense(licenseSlim);
    await flushPromises();

    expect(getMock).toHaveBeenCalledWith('MIT');
    expect(vm.currentLicenseForAction).toEqual(fullLicense);
    expect(vm.licenseDialogMode).toBe('edit');
  });

  it('loads the full license and opens the dialog in duplicate mode', async () => {
    const fullLicense = {...licenseSlim, description: 'A permissive license'};
    getMock.mockResolvedValue({data: fullLicense});

    const {wrapper} = createWrapper();
    await flushPromises();

    const vm = wrapper.vm as unknown as {
      duplicateLicense: (license: typeof licenseSlim) => Promise<void>;
      licenseDialogMode: string;
    };
    await vm.duplicateLicense(licenseSlim);
    await flushPromises();

    expect(vm.licenseDialogMode).toBe('duplicate');
  });

  it('applies a filter set and updates the route to the filter set url', async () => {
    getFilterSetsMock.mockResolvedValue([]);
    // Applying a filter set intersects the requested values with the facets known from the last
    // search response, so the license source facet must include 'spdx' for it to end up selected.
    searchMock.mockResolvedValue({
      data: {licenses: [licenseSlim], count: 1, meta: {...emptyMeta, possibleSources: {spdx: 1, mit: 2}}},
    });

    const {wrapper} = createWrapper();
    await flushPromises();

    const filterSet = {
      _key: 'fs-1',
      name: 'My Filter',
      table: 'licenses',
      includedFilters: [{name: 'source', values: ['spdx']}],
    };
    getFilterSetMock.mockResolvedValue(filterSet);

    const vm = wrapper.vm as unknown as {
      applyFilterSet: (filter: typeof filterSet) => Promise<void>;
      selectedFilterSource: string[];
    };
    await vm.applyFilterSet(filterSet as never);
    await flushPromises();

    expect(vm.selectedFilterSource).toEqual(['spdx']);
    expect(pushMock).toHaveBeenCalledWith('/dashboard/licenses/filtersets/fs-1');
  });
});
