// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useChecklistsStore} from '@disclosure-portal/stores/checklists.store';
import {useLabelStore} from '@disclosure-portal/stores/label.store';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {flushPromises} from '@vue/test-utils';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import ChecklistMain from '../ChecklistMain.vue';

const {getChecklistMock, deleteChecklistItemMock, snackbarInfoMock} = vi.hoisted(() => ({
  getChecklistMock: vi.fn(),
  deleteChecklistItemMock: vi.fn(),
  snackbarInfoMock: vi.fn(),
}));

vi.mock('@disclosure-portal/services/admin', () => ({
  default: {
    getChecklist: getChecklistMock,
    deleteChecklistItem: deleteChecklistItemMock,
  },
}));

vi.mock('@shared/composables/useSnackbar', () => ({
  default: () => ({
    info: snackbarInfoMock,
  }),
}));

// The checklists store resolves the active checklist from the current route's `id` param via
// vuetify's own useRoute composable (not vue-router's). Point it at a fixed checklist key.
vi.mock('vuetify/lib/composables/router', () => ({
  useRoute: () => ({value: {params: {id: 'cl-1'}}}),
  useRouter: () => ({push: vi.fn(), replace: vi.fn()}),
}));

const checklist = {
  _key: 'cl-1',
  name: 'My Checklist',
  nameDE: '',
  description: '',
  descriptionDE: '',
  policyLabels: [],
  active: true,
  items: [{_key: 'item-1', name: 'Item One', targetTemplateName: 'Template A', created: '', updated: ''}],
};

const childStubs = {
  ChecklistItemDialog: true,
  ConfirmationDialog: true,
  TableActionButtons: true,
  ProjectLabel: true,
  DCActionButton: true,
  'v-data-table': {
    props: ['headers', 'items'],
    template: '<div class="v-data-table">{{ JSON.stringify(items) }}</div>',
  },
};

const createWrapper = () =>
  mountView(ChecklistMain, {
    childStubs,
    beforePiniaMount: (pinia) => {
      const checklistsStore = useChecklistsStore(pinia);
      checklistsStore.checklists = [checklist as never];
      checklistsStore.loaded = true;
    },
  });

describe('ChecklistMain', () => {
  beforeEach(() => {
    getChecklistMock.mockReset();
    deleteChecklistItemMock.mockReset();
    snackbarInfoMock.mockReset();
  });

  it('sets breadcrumbs including the checklist name once the checklist is available', async () => {
    const {pinia} = createWrapper();
    await flushPromises();

    const breadcrumbs = useBreadcrumbsStore(pinia);
    expect(breadcrumbs.setCurrentBreadcrumbs).toHaveBeenCalledWith([
      ...breadcrumbs.dashboardCrumbs,
      {title: 'Checklists', href: '/dashboard/admin/checklist'},
      {title: 'My Checklist', disabled: false, href: '/dashboard/admin/checklist/cl-1'},
    ]);
  });

  it('renders without throwing and passes checklist items to the data table', async () => {
    const {wrapper} = createWrapper();
    await flushPromises();

    expect(wrapper.exists()).toBe(true);
    expect(wrapper.find('.v-data-table').text()).toContain('Item One');
  });

  it('resolves unknown policy labels via the label store', async () => {
    const {pinia} = createWrapper();
    useChecklistsStore(pinia).checklists = [{...checklist, policyLabels: ['missing-label']} as never];
    await flushPromises();

    const labelStore = useLabelStore(pinia);
    expect(labelStore.getLabelByKey('missing-label')).toEqual({name: 'UNKNOWN_LABEL', description: ''});
  });

  it('deletes a checklist item and shows a success snackbar', async () => {
    deleteChecklistItemMock.mockResolvedValue({data: {...checklist, items: []}});

    const {wrapper} = createWrapper();
    await flushPromises();

    await (wrapper.vm as unknown as {doDelete: (config: {key: string}) => Promise<void>}).doDelete({
      key: 'item-1',
    } as never);
    await flushPromises();

    expect(deleteChecklistItemMock).toHaveBeenCalledWith('cl-1', 'item-1');
    expect(snackbarInfoMock).toHaveBeenCalledWith('Checklist deleted');
  });
});
