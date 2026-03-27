// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {architectures, targetPlatforms} from '@disclosure-portal/model/NewWizard';
import {useWizardStore} from '@disclosure-portal/stores/wizard.store';
import {createTestingPinia} from '@pinia/testing';
import {mount} from '@vue/test-utils';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import WizardStepArchitecture from '../WizardStepArchitecture.vue';

describe('WizardStepArchitecture', () => {
  let wrapper: any;
  let wizardStore: any;

  beforeEach(() => {
    wrapper = mount(WizardStepArchitecture, {
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
            stubActions: true,
          }),
        ],
        stubs: ['Stack', 'WizardCard'],
      },
    });
    wizardStore = useWizardStore();
  });

  it('should render component', () => {
    expect(wrapper.exists()).toBe(true);
  });

  it('should display 2 enterprise cards for non-vehicle architecture', () => {
    wizardStore.project = {
      targetPlatform: targetPlatforms.enterprise,
      architecture: null,
    };

    const script = wrapper.vm as any;
    expect(script.cardList).toHaveLength(2);
    expect(script.cardList[0].key).toBe(architectures.frontend);
    expect(script.cardList[1].key).toBe(architectures.backend);
  });

  it('should display 2 vehicle cards for vehicle architecture', () => {
    wizardStore.project = {
      targetPlatform: targetPlatforms.vehicle,
      architecture: null,
    };
    const script = wrapper.vm as any;
    expect(script.cardList).toHaveLength(2);
    expect(script.cardList[0].key).toBe(architectures.vehicleOnboard);
    expect(script.cardList[1].key).toBe(architectures.vehicleOffboard);
  });

  it('should go to next step when onCardSelect is called', () => {
    const script = wrapper.vm as any;
    script.onCardSelect(architectures.frontend);
    expect(wizardStore.nextStep).toHaveBeenCalled();
  });

  it('should set architecture when card is selected', () => {
    const script = wrapper.vm as any;
    script.onCardSelect(architectures.backend);
    expect(wizardStore.project.architecture).toBe(architectures.backend);
  });

  it('should mark frontend card as active when architecture is frontend', () => {
    const testingPinia = createTestingPinia({
      createSpy: vi.fn,
      stubActions: true,
      initialState: {
        wizard: {
          project: {
            targetPlatform: targetPlatforms.enterprise,
            architecture: architectures.frontend,
          },
        },
      },
    });

    wrapper = mount(WizardStepArchitecture, {
      global: {
        plugins: [testingPinia],
        stubs: ['Stack', 'WizardCard'],
      },
    });

    const script = wrapper.vm as any;
    const frontendCard = script.architectureCardList.find((card: any) => card.key === architectures.frontend);

    expect(frontendCard.isActive).toBe(true);
  });

  it('should switch card list based on target platform', () => {
    wizardStore.project = {
      targetPlatform: targetPlatforms.enterprise,
    };

    let script = wrapper.vm as any;
    expect(script.isVehicleArchitecture).toBe(false);
    expect(script.cardList[0].key).toBe(architectures.frontend);

    wizardStore.project.targetPlatform = targetPlatforms.vehicle;

    wrapper = mount(WizardStepArchitecture, {
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
            stubActions: true,
            initialState: {
              wizard: {
                project: {
                  targetPlatform: targetPlatforms.vehicle,
                },
              },
            },
          }),
        ],
        stubs: ['Stack', 'WizardCard'],
      },
    });

    script = wrapper.vm as any;
    expect(script.isVehicleArchitecture).toBe(true);
    expect(script.cardList[0].key).toBe(architectures.vehicleOnboard);
  });
});
