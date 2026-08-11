// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {Rights} from '@shared/user/models/Rights';
import {useUserStore} from '@shared/user/stores/user.store';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {reactive, ref} from 'vue';
import {flushPromises} from '@vue/test-utils';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import AdminClassifications from '../AdminClassifications.vue';

const {getAllObligationsMock, deleteObligationMock, getCountOfLicencesUsingThisObligationMock, snackbarInfoMock} =
  vi.hoisted(() => ({
    getAllObligationsMock: vi.fn(),
    deleteObligationMock: vi.fn(),
    getCountOfLicencesUsingThisObligationMock: vi.fn(),
    snackbarInfoMock: vi.fn(),
  }));

vi.mock('@disclosure-portal/services/admin', () => ({
  default: {
    getAllObligations: getAllObligationsMock,
    deleteObligation: deleteObligationMock,
  },
}));

vi.mock('@disclosure-portal/services/license', () => ({
  default: {
    getCountOfLicencesUsingThisObligation: getCountOfLicencesUsingThisObligationMock,
  },
}));

vi.mock('@shared/composables/useSnackbar', () => ({
  default: () => ({
    info: snackbarInfoMock,
  }),
}));

// AdminClassifications reads the current language via storeToRefs(useLanguageStore()), which
// (as a real store) needs vuetify's own locale composable/plugin installed. None of that setup
// is relevant to this view's own logic, so stub the store instead. `appLanguage` must be a ref on
// the raw target object - storeToRefs() only promotes properties that are themselves refs.
vi.mock('@shared/stores/language.store', () => ({
  useLanguageStore: () => reactive({appLanguage: ref('en')}),
}));

const obligation = {
  _key: 'obl-1',
  type: 'obligation',
  warnLevel: 'green',
  name: 'Obligation One',
  nameDe: 'Verpflichtung Eins',
  description: 'desc',
  descriptionDe: 'desc de',
  created: '2026-01-01T10:00:00Z',
  updated: '2026-01-02T10:00:00Z',
};

const childStubs = {
  TableLayout: {template: '<div><slot name="buttons"/><slot name="table"/></div>'},
  NewClassificationDialog: true,
  ConfirmationDialog: true,
  AuditDialog: {template: '<div><slot :open="() => {}" /></div>'},
  'v-data-table': {
    props: ['items'],
    template: '<div class="v-data-table">{{ JSON.stringify(items) }}</div>',
  },
};

describe('AdminClassifications', () => {
  beforeEach(() => {
    getAllObligationsMock.mockReset();
    deleteObligationMock.mockReset();
    getCountOfLicencesUsingThisObligationMock.mockReset();
    snackbarInfoMock.mockReset();
    getAllObligationsMock.mockResolvedValue({data: {items: [obligation]}});
  });

  it('sets the expected breadcrumbs on mount', async () => {
    const {pinia} = mountView(AdminClassifications, {childStubs});
    await flushPromises();

    const breadcrumbs = useBreadcrumbsStore(pinia);
    expect(breadcrumbs.setCurrentBreadcrumbs).toHaveBeenCalledWith([
      ...breadcrumbs.dashboardCrumbs,
      {title: 'Classification', href: '/dashboard/admin/obligations'},
    ]);
  });

  it('renders without throwing and loads classifications into the table', async () => {
    const {wrapper} = mountView(AdminClassifications, {childStubs});
    await flushPromises();

    expect(wrapper.exists()).toBe(true);
    expect(getAllObligationsMock).toHaveBeenCalledTimes(1);
    expect(wrapper.find('.v-data-table').text()).toContain('Obligation One');
  });

  it('shows a plain delete confirmation when the classification is unused, and deletes on confirm', async () => {
    getCountOfLicencesUsingThisObligationMock.mockResolvedValue({data: {count: 0}});
    deleteObligationMock.mockResolvedValue({data: {}});

    const {wrapper} = mountView(AdminClassifications, {
      childStubs,
      beforePiniaMount: (p) => {
        useUserStore(p).setSimpleProfileData({
          rights: {allowObligation: {delete: true}} as unknown as Rights,
          profile: {} as never,
          allowed: true,
        });
      },
    });
    await flushPromises();

    const vm = wrapper.vm as unknown as {
      showDeletionConfirmationDialog: (item: typeof obligation) => Promise<void>;
      confirmConfig: {okButtonIsDisabled: boolean; key: string};
      confirmVisible: boolean;
    };
    await vm.showDeletionConfirmationDialog(obligation);
    await flushPromises();

    expect(getCountOfLicencesUsingThisObligationMock).toHaveBeenCalledWith('obl-1');
    expect(vm.confirmVisible).toBe(true);
    expect(vm.confirmConfig.okButtonIsDisabled).toBe(false);

    await (
      wrapper.vm as unknown as {doDelete: (config: {key: string; okButtonIsDisabled: boolean}) => Promise<void>}
    ).doDelete(vm.confirmConfig as never);
    await flushPromises();

    expect(deleteObligationMock).toHaveBeenCalledWith('obl-1');
    expect(snackbarInfoMock).toHaveBeenCalledWith('classification deleted');
  });

  it('disables the confirm button when the classification is still in use by licenses', async () => {
    getCountOfLicencesUsingThisObligationMock.mockResolvedValue({data: {count: 3}});

    const {wrapper} = mountView(AdminClassifications, {childStubs});
    await flushPromises();

    const vm = wrapper.vm as unknown as {
      showDeletionConfirmationDialog: (item: typeof obligation) => Promise<void>;
      confirmConfig: {okButtonIsDisabled: boolean};
    };
    await vm.showDeletionConfirmationDialog(obligation);
    await flushPromises();

    expect(vm.confirmConfig.okButtonIsDisabled).toBe(true);
  });
});
