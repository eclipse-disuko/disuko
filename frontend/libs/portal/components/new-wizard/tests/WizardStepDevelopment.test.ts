// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {developments} from '@disclosure-portal/model/NewWizard';
import {useWizardStore} from '@disclosure-portal/stores/wizard.store';
import {vuetifyStubs} from '@disclosure-portal/test-utils/vuetify-stubs';
import {createTestingPinia} from '@pinia/testing';
import {mount} from '@vue/test-utils';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import WizardStepDevelopment from '../WizardStepDevelopment.vue';

describe('WizardStepDevelopment', () => {
  let wrapper: any;
  let wizardStore: any;

  beforeEach(() => {
    wrapper = mount(WizardStepDevelopment, {
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
            stubActions: true,
          }),
        ],
        stubs: {
          ...vuetifyStubs,
          Stack: {template: '<div class="stack"><slot /></div>'},
          WizardCard: true,
        },
      },
    });
    wizardStore = useWizardStore();
  });

  it('should render component', () => {
    expect(wrapper.exists()).toBe(true);
  });

  it('should have 3 cards in cardList', () => {
    const script = wrapper.vm;
    expect(script.cardList).toHaveLength(3);
  });

  it('should set development and call necessary methods when card is selected', () => {
    const script = wrapper.vm;
    script.onCardSelect(developments.inhouse);

    expect(wizardStore.project.development).toBe(developments.inhouse);
    expect(wizardStore.updateProjectSettingsBasedOnDevelopment).toHaveBeenCalled();
    expect(wizardStore.setAvailableSteps).toHaveBeenCalled();
    expect(wizardStore.nextStep).toHaveBeenCalled();
  });

  it('should clear supplierName when development is selected', () => {
    wizardStore.project.projectSettings = {
      documentMeta: {
        supplierName: 'Test Supplier',
        supplierDept: null,
      },
    };

    const script = wrapper.vm;
    script.onCardSelect(developments.external);
    expect(wizardStore.project.projectSettings.documentMeta.supplierName).toBe('');
  });

  it('should clear supplierDept when development is selected and dept is not empty', () => {
    wizardStore.project.projectSettings = {
      documentMeta: {
        supplierName: '',
        supplierDept: {id: '123', name: 'Test Dept'},
      },
    };

    const script = wrapper.vm;
    script.onCardSelect(developments.internal);

    expect(wizardStore.project.projectSettings.documentMeta.supplierDept).toBeNull();
  });

  it('should mark inhouse card as active when development is inhouse', () => {
    const testingPinia = createTestingPinia({
      createSpy: vi.fn,
      stubActions: true,
      initialState: {
        wizard: {
          project: {
            development: developments.inhouse,
          },
        },
      },
    });

    wrapper = mount(WizardStepDevelopment, {
      global: {
        plugins: [testingPinia],
        stubs: {
          ...vuetifyStubs,
          Stack: {template: '<div class="stack"><slot /></div>'},
          WizardCard: true,
        },
      },
    });

    const script = wrapper.vm;
    const inhouseCard = script.cardList.find((card: any) => card.key === developments.inhouse);

    expect(inhouseCard.isActive).toBe(true);
  });
});
