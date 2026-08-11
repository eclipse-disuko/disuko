// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {Group, Rights} from '@shared/user/models/Rights';
import {useUserStore} from '@shared/user/stores/user.store';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {flushPromises} from '@vue/test-utils';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import Stats from '../Stats.vue';

const {getStatsMock, downloadReportMock, downloadReportXLSXMock, downloadCombinedReportXLSXMock} = vi.hoisted(() => ({
  getStatsMock: vi.fn(),
  downloadReportMock: vi.fn(),
  downloadReportXLSXMock: vi.fn(),
  downloadCombinedReportXLSXMock: vi.fn(),
}));

vi.mock('@disclosure-portal/services/analytics', () => ({
  default: {
    getStats: getStatsMock,
    downloadReport: downloadReportMock,
    downloadReportXLSX: downloadReportXLSXMock,
    downloadCombinedReportXLSX: downloadCombinedReportXLSXMock,
  },
}));

const actionButtonStub = {
  template: '<button type="button">{{ text }}</button>',
  props: ['text', 'hint', 'large', 'icon', 'disabled'],
};

const childStubs = {
  DCActionButton: actionButtonStub,
  TableLayout: {template: '<div><slot name="buttons"/><slot name="table"/></div>'},
};

const createWrapper = (groups: Group[] = []) =>
  mountView(Stats, {
    childStubs,
    beforePiniaMount: (pinia) => {
      useUserStore(pinia).setSimpleProfileData({
        rights: {groups} as unknown as Rights,
        profile: {} as never,
        allowed: true,
      });
    },
  });

describe('Stats', () => {
  beforeEach(() => {
    getStatsMock.mockReset();
    downloadReportMock.mockReset();
    downloadReportXLSXMock.mockReset();
    downloadCombinedReportXLSXMock.mockReset();
    downloadReportMock.mockResolvedValue({data: 'csv'});
    downloadReportXLSXMock.mockResolvedValue({data: 'xlsx'});
    downloadCombinedReportXLSXMock.mockResolvedValue({data: 'xlsx'});
  });

  // Stats.vue has no breadcrumbs call, so that step of the standard checklist is skipped here.

  it('fetches stats on mount and renders them in the table', async () => {
    getStatsMock.mockResolvedValue({data: {projectCount: 7, licenseCount: 3, userCount: 12}});

    const {wrapper} = createWrapper();
    await flushPromises();

    expect(getStatsMock).toHaveBeenCalledTimes(1);
    expect(wrapper.text()).toContain('7');
    expect(wrapper.text()).toContain('3');
    expect(wrapper.text()).toContain('12');
  });

  it('hides the download buttons for users who are not project analysts', async () => {
    getStatsMock.mockResolvedValue({data: {}});

    const {wrapper} = createWrapper([]);
    await flushPromises();

    expect(wrapper.findAll('button').length).toBe(0);
  });

  it('downloads the CSV and XLSX reports for project analysts', async () => {
    getStatsMock.mockResolvedValue({data: {}});

    const {wrapper} = createWrapper([Group.UserProjectAnalyst]);
    await flushPromises();

    const buttons = wrapper.findAll('button');
    await buttons[0].trigger('click');
    await buttons[1].trigger('click');

    expect(downloadReportMock).toHaveBeenCalledTimes(1);
    expect(downloadReportXLSXMock).toHaveBeenCalledTimes(1);
  });

  it('downloads a combined report for the selected months', async () => {
    getStatsMock.mockResolvedValue({data: {}});

    const {wrapper} = createWrapper([Group.UserProjectAnalyst]);
    await flushPromises();

    (wrapper.vm as unknown as {selectedMonths: unknown[]}).selectedMonths = [
      {month: 1, year: 2026, label: 'January 2026'},
    ];
    await (wrapper.vm as unknown as {downloadCombinedReportXLSX: () => Promise<void>}).downloadCombinedReportXLSX();

    expect(downloadCombinedReportXLSXMock).toHaveBeenCalledWith({
      months: [{month: 1, year: 2026, label: 'January 2026'}],
    });
    expect((wrapper.vm as unknown as {showMonthMenu: boolean}).showMonthMenu).toBe(false);
  });
});
