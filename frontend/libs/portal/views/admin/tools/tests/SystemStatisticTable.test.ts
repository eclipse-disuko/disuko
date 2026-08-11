// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {vuetifyStubs} from '@disclosure-portal/test-utils/vuetify-stubs';
import {config, mount} from '@vue/test-utils';
import {describe, expect, it} from 'vitest';
import SystemStatisticTable from '../SystemStatisticTable.vue';

if (config.global.mocks && '$t' in config.global.mocks) {
  delete config.global.mocks.$t;
}

const baseStat = {
  _key: 's1',
  created: '2026-01-01T10:00:00Z',
  projectCount: 10,
  projectActiveCount: 8,
  projectDeletedCount: 2,
  missingProjects: false,
  maxVersionsInOneProject: 3,
  projectsOverOrAtVersionLimit: 1,
  versionLimit: 5,
  licenseCount: 20,
  missingLicenses: false,
  licenseActiveCount: 18,
  licenseChartCount: 4,
  licenseDeletedCount: 2,
  policyRuleCount: 6,
  missingPolicyRules: false,
  policyRuleActiveCount: 5,
  policyRuleDeletedCount: 1,
  labelCount: 9,
  schemaCount: 2,
  obligationCount: 15,
  missingObligations: false,
  obligationActiveCount: 14,
  obligationDeletedCount: 1,
  userCount: 30,
  missingUsers: false,
  userActiveCount: 28,
  userDeactivateCount: 1,
  userTermsNotAcceptedCount: 1,
  userDeprovisionedCount: 0,
  uploadFileCnt: 40,
  missingUploadFiles: false,
  uploadFileCntSBOM: 12,
  uploadFileCntPDF: 20,
  uploadFileCntJSON: 8,
  dbBackupFileCnt: 3,
};

describe('SystemStatisticTable', () => {
  it('renders nothing when stats is not set', () => {
    const wrapper = mount(SystemStatisticTable, {
      props: {stats: undefined as unknown as never[]},
      global: {stubs: vuetifyStubs},
    });

    expect(wrapper.find('table').exists()).toBe(false);
  });

  it('renders one column per stats snapshot with the expected values', () => {
    const wrapper = mount(SystemStatisticTable, {
      props: {stats: [baseStat]},
      global: {stubs: vuetifyStubs},
    });

    const text = wrapper.text();
    expect(text).toContain('10'); // projectCount
    expect(text).toContain('20'); // licenseCount
    expect(text).toContain('30'); // userCount
    expect(text).not.toContain('!!!');
  });

  it('flags missing entities with an error marker', () => {
    const wrapper = mount(SystemStatisticTable, {
      props: {stats: [{...baseStat, missingProjects: true, missingUsers: true}]},
      global: {stubs: vuetifyStubs},
    });

    const rows = wrapper.findAll('tr');
    const projectsRow = rows.find((row) => row.text().startsWith('Projects'));
    expect(projectsRow?.text()).toContain('!!!');
    const usersRow = rows.find((row) => row.text().startsWith('Users') && !row.text().startsWith('Users ('));
    expect(usersRow?.text()).toContain('!!!');
  });
});
