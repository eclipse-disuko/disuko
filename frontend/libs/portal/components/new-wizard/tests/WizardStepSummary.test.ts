// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useLabelStore} from '@disclosure-portal/stores/label.store';
import {useWizardStore} from '@disclosure-portal/stores/wizard.store';
import {vuetifyStubs} from '@disclosure-portal/test-utils/vuetify-stubs';
import {createTestingPinia} from '@pinia/testing';
import {mount, VueWrapper} from '@vue/test-utils';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import WizardStepSummary from '../WizardStepSummary.vue';

vi.mock('@disclosure-portal/utils/View', () => ({
  getStrWithMaxLength: vi.fn((max: number, str: string) => str.substring(0, max)),
}));

describe('WizardStepSummary', () => {
  let wrapper: VueWrapper;
  let wizardStore: any;
  let labelStore: any;

  const createWrapper = (options = {}) => {
    return mount(WizardStepSummary, {
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
            stubActions: false,
          }),
        ],
        stubs: {
          ...vuetifyStubs,
          Stack: {template: '<div class="stack"><slot /></div>'},
          SummaryItem: {
            template: '<div class="summary-item">{{ value }}<slot /></div>',
            props: ['label', 'value', 'showDash'],
          },
          Tooltip: {template: '<span class="tooltip"></span>', props: ['text']},
          ProjectLabel: {template: '<div class="project-label"></div>', props: ['label']},
        },
      },
      ...options,
    });
  };

  beforeEach(() => {
    vi.clearAllMocks();
    wrapper = createWrapper();
    wizardStore = useWizardStore();
    labelStore = useLabelStore();

    // Set default values
    wizardStore.project = {
      name: 'Test Project',
      description: 'Test Description',
      targetPlatform: 'Enterprise',
      architecture: 'Backend',
      targetUsers: 'Internal',
      distributionTarget: 'Company',
      labels: [],
      projectSettings: {
        documentMeta: {
          supplierDept: null,
          supplierAddress: '',
          supplierNr: '',
          supplierName: '',
        },
        customerMeta: {
          dept: null,
          address: '',
        },
        noticeContactMeta: {
          address: '',
        },
      },
    };
    wizardStore.previewLoading = false;
    wizardStore.isVehicleOnboardArchitecture = false;
  });

  it('should render component', () => {
    expect(wrapper.exists()).toBe(true);
  });

  it('should call wizard preview endpoint on mount', () => {
    expect(wizardStore.preview).toHaveBeenCalled();
  });

  it('should display project name and description', () => {
    expect(wrapper.text()).toContain('Test Project');
    expect(wrapper.text()).toContain('Test Description');
  });

  it('should display dash when a field is empty', async () => {
    wizardStore.project.description = '';
    await wrapper.vm.$nextTick();

    const summaryItems = wrapper.findAll('.summary-item');
    const descriptionItem = summaryItems.find((item) => item.text().includes('-'));
    expect(descriptionItem).toBeDefined();
  });

  it('should display project platform details', () => {
    expect(wrapper.text()).toContain('Enterprise');
    expect(wrapper.text()).toContain('Backend');
    expect(wrapper.text()).toContain('Internal');
    expect(wrapper.text()).toContain('Company');
  });

  it('should display dash when application meta id is not present', () => {
    wizardStore.project.applicationMeta = undefined;
    expect(wrapper.text()).toContain('-');
  });

  it('should display application name when application meta is present', async () => {
    wizardStore.project.applicationMeta = {id: '123', name: 'Test App'};
    await wrapper.vm.$nextTick();

    expect(wrapper.text()).toContain('Test App');
  });

  it('should format department information correctly', async () => {
    wizardStore.project.projectSettings.documentMeta.supplierDept = {
      deptId: 'D001',
      companyCode: 'CC001',
      companyName: 'Test Company',
      orgAbbreviation: 'TC',
      descriptionEnglish: 'Test Department',
    };
    await wrapper.vm.$nextTick();

    expect(wrapper.text()).toContain('[CC001] Test Company / [D001,TC] Test Department');
  });

  it('should display "- / -" when department is not present', () => {
    wizardStore.project.projectSettings.documentMeta.supplierDept = null;
    expect(wrapper.text()).toContain('- / -');
  });

  it('should display supplier information', async () => {
    wizardStore.project.projectSettings.documentMeta.supplierAddress = '123 Test';
    wizardStore.project.projectSettings.documentMeta.supplierNr = 'SUP001';
    wizardStore.project.projectSettings.documentMeta.supplierName = 'Test Supplier';
    await wrapper.vm.$nextTick();

    expect(wrapper.text()).toContain('123 Test');
    expect(wrapper.text()).toContain('SUP001');
    expect(wrapper.text()).toContain('Test Supplier');
  });

  it('should display customer information', async () => {
    wizardStore.project.projectSettings.customerMeta.address = 'Stuttgart';
    await wrapper.vm.$nextTick();

    expect(wrapper.text()).toContain('Stuttgart');
  });

  it('should not display notice contact address for vehicle onboard architecture', async () => {
    wizardStore.isVehicleOnboardArchitecture = true;
    wizardStore.project.projectSettings.noticeContactMeta.address = 'Contact Address';
    await wrapper.vm.$nextTick();

    const divs = wrapper.findAll('div');
    const hasContactAddress = divs.some((div) => div.text().includes('Contact Address'));
    expect(hasContactAddress).toBe(false);
  });

  it('should show loading skeleton when preview is loading', async () => {
    wizardStore.previewLoading = true;
    await wrapper.vm.$nextTick();

    const loadingSkeleton = wrapper.findAll('div').filter((div) => div.classes().includes('animate-pulse'));
    expect(loadingSkeleton.length).toBeGreaterThan(0);
  });

  it('should display project labels when labels are loaded', async () => {
    wizardStore.previewLoading = false;
    wizardStore.project.labels = ['label1', 'label2', 'label3'];
    labelStore.getLabelByKey = vi.fn((key) => ({key, name: `Label ${key}`}));
    await wrapper.vm.$nextTick();

    const projectLabels = wrapper.findAll('.project-label');
    expect(projectLabels.length).toBe(3);
  });

  it('should display dash for empty supplier information', () => {
    wizardStore.project.projectSettings.documentMeta.supplierAddress = '';
    wizardStore.project.projectSettings.documentMeta.supplierNr = '';
    wizardStore.project.projectSettings.documentMeta.supplierName = '';

    expect(wrapper.text()).toContain('-');
  });
});
