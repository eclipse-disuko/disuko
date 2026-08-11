// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {ConfirmationType} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {flushPromises} from '@vue/test-utils';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import ReviewTemplates from '../ReviewTemplates.vue';

const {getReviewTemplatesMock, deleteReviewTemplateMock, createReviewTemplateMock, snackbarInfoMock} = vi.hoisted(
  () => ({
    getReviewTemplatesMock: vi.fn(),
    deleteReviewTemplateMock: vi.fn(),
    createReviewTemplateMock: vi.fn(),
    snackbarInfoMock: vi.fn(),
  }),
);

vi.mock('@disclosure-portal/services/admin', () => ({
  default: {
    getReviewTemplates: getReviewTemplatesMock,
    deleteReviewTemplate: deleteReviewTemplateMock,
    createReviewTemplate: createReviewTemplateMock,
  },
}));

vi.mock('@shared/composables/useSnackbar', () => ({
  default: () => ({
    info: snackbarInfoMock,
  }),
}));

const reviewTemplates = [
  {
    _key: 'rt-1',
    title: 'Template One',
    level: 'GREEN',
    description: 'desc',
    source: 'source',
    created: '2026-01-01T10:00:00Z',
    updated: '2026-01-02T10:00:00Z',
  },
];

const childStubs = {
  TableLayout: {template: '<div><slot name="buttons"/><slot name="table"/></div>'},
  ReviewTemplateDialog: true,
  ConfirmationDialog: true,
  'v-data-table': {
    props: ['items'],
    template: '<div class="v-data-table">{{ JSON.stringify(items) }}</div>',
  },
};

describe('ReviewTemplates', () => {
  beforeEach(() => {
    getReviewTemplatesMock.mockReset();
    deleteReviewTemplateMock.mockReset();
    createReviewTemplateMock.mockReset();
    snackbarInfoMock.mockReset();
    getReviewTemplatesMock.mockResolvedValue({data: reviewTemplates});
  });

  it('sets the expected breadcrumbs on mount', async () => {
    const {pinia} = mountView(ReviewTemplates, {childStubs});
    await flushPromises();

    const breadcrumbs = useBreadcrumbsStore(pinia);
    expect(breadcrumbs.setCurrentBreadcrumbs).toHaveBeenCalledWith([
      ...breadcrumbs.dashboardCrumbs,
      {title: 'Review Templates', href: '/dashboard/templates/review'},
    ]);
  });

  it('renders without throwing and loads review templates into the table', async () => {
    const {wrapper} = mountView(ReviewTemplates, {childStubs});
    await flushPromises();

    expect(wrapper.exists()).toBe(true);
    expect(getReviewTemplatesMock).toHaveBeenCalledTimes(1);
    expect(wrapper.find('.v-data-table').text()).toContain('Template One');
  });

  it('deletes a review template on confirm and shows a success snackbar', async () => {
    deleteReviewTemplateMock.mockResolvedValue({data: {}});

    const {wrapper} = mountView(ReviewTemplates, {childStubs});
    await flushPromises();

    await (wrapper.vm as unknown as {onConfirm: (config: {type: ConfirmationType; key: string}) => void}).onConfirm({
      type: ConfirmationType.DELETE,
      key: 'rt-1',
    });
    await flushPromises();

    expect(deleteReviewTemplateMock).toHaveBeenCalledWith('rt-1');
    expect(snackbarInfoMock).toHaveBeenCalledWith('review template deleted');
    expect(getReviewTemplatesMock).toHaveBeenCalledTimes(2);
  });

  it('creates a review template via the form submission composable and shows a success snackbar', async () => {
    createReviewTemplateMock.mockResolvedValue({data: {}});

    const {wrapper} = mountView(ReviewTemplates, {childStubs});
    await flushPromises();

    await (wrapper.vm as unknown as {openDialog: (mode: 'create' | 'edit') => Promise<void>}).openDialog('create');
    await (wrapper.vm as unknown as {handleSave: (formData: {title: string}) => Promise<void>}).handleSave({
      title: 'New Template',
    } as never);
    await flushPromises();

    expect(createReviewTemplateMock).toHaveBeenCalledWith({title: 'New Template'});
    expect(snackbarInfoMock).toHaveBeenCalledWith('review template created');
  });
});
