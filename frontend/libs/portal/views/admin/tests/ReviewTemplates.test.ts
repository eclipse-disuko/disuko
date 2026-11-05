// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {ConfirmationType} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import {ReviewRemarkLevel, compareRRLevel} from '@disclosure-portal/model/Quality';
import {getIconColorReviewRemarkLevel, getIconReviewRemarkLevel} from '@disclosure-portal/utils/View';
import en from '@disclosure-portal/i18n/locales/en.json';
import de from '@disclosure-portal/i18n/locales/de.json';
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
    level: 'BLUE',
    description: 'desc',
    source: 'source',
    created: '2026-01-01T10:00:00Z',
    updated: '2026-01-02T10:00:00Z',
  },
];

const childStubs = {
  TableLayout: {template: '<div><slot name="buttons"/><slot name="table"/></div>'},
  ReviewTemplateDialog: {props: ['levels'], template: '<div class="template-levels">{{ levels.join(",") }}</div>'},
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

  it('offers BLUE Information instead of GREEN Closed in the template selector', async () => {
    const {wrapper} = mountView(ReviewTemplates, {childStubs});
    await flushPromises();

    expect(wrapper.find('.template-levels').text()).toBe('BLUE,YELLOW,RED');
    expect(wrapper.find('.v-data-table').text()).toContain('BLUE');
  });

  it('distinguishes Information and Closed in labels, icons, colors and sorting', () => {
    expect([en.REMARK_LEVEL_BLUE, de.REMARK_LEVEL_BLUE]).toEqual(['Information', 'Information']);
    expect([en.REMARK_LEVEL_GREEN, de.REMARK_LEVEL_GREEN]).toEqual(['Closed', 'Geschlossen']);
    expect([en.RR_CHART_INFORMATION, en.RR_CHART_ACCEPTABLE]).toEqual(['Information', 'Closed']);
    expect(getIconColorReviewRemarkLevel(ReviewRemarkLevel.BLUE)).toBe('blue');
    expect(getIconColorReviewRemarkLevel(ReviewRemarkLevel.GREEN)).toBe('green');
    expect(getIconReviewRemarkLevel(ReviewRemarkLevel.BLUE)).toBe('mdi-comment-text');
    expect(getIconReviewRemarkLevel(ReviewRemarkLevel.GREEN)).toBe('mdi-comment-check');
    expect(
      [
        ReviewRemarkLevel.RED,
        ReviewRemarkLevel.BLUE,
        ReviewRemarkLevel.GREEN,
        ReviewRemarkLevel.YELLOW,
        ReviewRemarkLevel.GRAY,
      ].sort(compareRRLevel),
    ).toEqual([
      ReviewRemarkLevel.GRAY,
      ReviewRemarkLevel.GREEN,
      ReviewRemarkLevel.BLUE,
      ReviewRemarkLevel.YELLOW,
      ReviewRemarkLevel.RED,
    ]);
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
