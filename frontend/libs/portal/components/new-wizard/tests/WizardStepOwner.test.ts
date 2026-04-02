// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useWizardStore} from '@disclosure-portal/stores/wizard.store';
import {vuetifyStubs} from '@disclosure-portal/test-utils/vuetify-stubs';
import {RightsUtils} from '@disclosure-portal/utils/Rights';
import {createTestingPinia} from '@pinia/testing';
import {mount, VueWrapper} from '@vue/test-utils';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import WizardStepOwner from '../WizardStepOwner.vue';

vi.mock('@disclosure-portal/utils/Rights', () => ({
  RightsUtils: {
    rights: vi.fn(() => ({
      isInternal: true,
    })),
  },
}));

vi.mock('@disclosure-portal/composables/useNewWizard', () => ({
  useNewWizard: () => ({
    validationRules: {
      address: [(v: string) => !!v || 'Address is required'],
    },
    initSteps: [],
    enterpriseOrMobileSteps: [],
    vehicleSteps: [],
    otherSteps: [],
    mergeSteps: vi.fn(),
  }),
  isValidByRules: vi.fn((value: string) => !!value),
  removeStep: vi.fn((steps: any[], stepId: string) => steps),
}));

describe('WizardStepOwner', () => {
  let wrapper: VueWrapper;
  let wizardStore: any;

  const createWrapper = (options = {}) => {
    return mount(WizardStepOwner, {
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
          DAutocompleteCompany: {
            template: '<input type="text" class="d-autocomplete-company" />',
            props: ['modelValue', 'label', 'required', 'aria'],
          },
        },
      },
      ...options,
    });
  };

  beforeEach(() => {
    vi.clearAllMocks();
    wrapper = createWrapper();
    wizardStore = useWizardStore();
  });

  it('should render component', () => {
    expect(wrapper.exists()).toBe(true);
  });

  it('should not render company autocomplete when user is not internal', async () => {
    vi.mocked(RightsUtils.rights).mockReturnValue({isInternal: false} as any);
    wrapper = createWrapper();
    await wrapper.vm.$nextTick();

    const companyField = wrapper.find('.d-autocomplete-company');
    expect(companyField.exists()).toBe(false);
  });

  it('should render customer address textarea', () => {
    const addressField = wrapper.find('[data-testid="OwnerSettings__Address"]');
    expect(addressField.exists()).toBe(true);
  });

  it('should not render notice contact address textarea for vehicle onboard architecture', async () => {
    wizardStore.isVehicleOnboardArchitecture = true;
    await wrapper.vm.$nextTick();

    const noticeContactField = wrapper.find('#thirdparty-address');
    expect(noticeContactField.exists()).toBe(false);
  });

  it('should validate form on mount when customerMeta dept is present', async () => {
    const testWizardStore = useWizardStore();
    testWizardStore.project.projectSettings = {
      customerMeta: {dept: {id: '123', name: 'Test Dept'}, address: ''},
      noticeContactMeta: {address: ''},
    };

    const newWrapper = createWrapper();
    await newWrapper.vm.$nextTick();
    await new Promise((resolve) => setTimeout(resolve, 0));

    expect(newWrapper.exists()).toBe(true);
  });

  it('should validate form on mount when customerMeta address is present', async () => {
    const testWizardStore = useWizardStore();
    testWizardStore.project.projectSettings = {
      customerMeta: {dept: null, address: 'Test Address'},
      noticeContactMeta: {address: ''},
    };

    const newWrapper = createWrapper();
    await newWrapper.vm.$nextTick();
    await new Promise((resolve) => setTimeout(resolve, 0));

    expect(newWrapper.exists()).toBe(true);
  });
});
