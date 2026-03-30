// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {architectures, targetPlatforms} from '@disclosure-portal/model/NewWizard';
import {useWizardStore} from '@disclosure-portal/stores/wizard.store';
import {createTestingPinia} from '@pinia/testing';
import {mount} from '@vue/test-utils';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import WizardStepPlatform from '../WizardStepPlatform.vue';

describe('WizardStepPlatform', () => {
  let wrapper: any;
  let wizardStore: any;

  beforeEach(() => {
    wrapper = mount(WizardStepPlatform, {
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

  it('should render title text', () => {
    wrapper = mount(WizardStepPlatform, {
      global: {
        plugins: [createTestingPinia({createSpy: vi.fn, stubActions: true})],
        stubs: {
          Stack: {
            template: '<div class="stack"><slot /></div>',
          },
          WizardCard: true,
        },
      },
    });

    expect(wrapper.find('h2').exists()).toBe(true);
  });

  it('should have 4 cards in cardList', () => {
    const script = wrapper.vm;
    expect(script.cardList).toHaveLength(4);
  });

  it('should call setAvailableSteps and nextStep when onCardSelect is called', () => {
    const script = wrapper.vm;
    script.onCardSelect(targetPlatforms.enterprise);

    expect(wizardStore.setAvailableSteps).toHaveBeenCalled();
    expect(wizardStore.nextStep).toHaveBeenCalled();
  });

  it('should reset architecture when switching from vehicle to non-vehicle platform', () => {
    wizardStore.project.architecture = architectures.vehicleOnboard;

    const script = wrapper.vm;
    script.onCardSelect(targetPlatforms.enterprise);

    expect(wizardStore.project.architecture).toBeNull();
  });

  it('should reset architecture when switching from non-vehicle to vehicle platform', () => {
    wizardStore.project.architecture = architectures.frontend;

    const script = wrapper.vm;
    script.onCardSelect(targetPlatforms.vehicle);

    expect(wizardStore.project.architecture).toBeNull();
    expect(wizardStore.project.targetUsers).toBeNull();
    expect(wizardStore.project.distributionTarget).toBeNull();
  });

  it('should reset all dependent fields when selecting "other" platform', () => {
    wizardStore.project.architecture = architectures.backend;
    wizardStore.project.targetUsers = 'internal';
    wizardStore.project.distributionTarget = 'production';

    const script = wrapper.vm;
    script.onCardSelect(targetPlatforms.other);

    expect(wizardStore.project.architecture).toBeNull();
    expect(wizardStore.project.targetUsers).toBeNull();
    expect(wizardStore.project.distributionTarget).toBeNull();
  });

  it('should set targetPlatform when card is selected', () => {
    const script = wrapper.vm;
    script.onCardSelect(targetPlatforms.mobile);

    expect(wizardStore.project.targetPlatform).toBe(targetPlatforms.mobile);
  });

  it('should mark enterprise card as active when targetPlatform is enterprise', () => {
    const testingPinia = createTestingPinia({
      createSpy: vi.fn,
      stubActions: true,
      initialState: {
        wizard: {
          project: {
            targetPlatform: targetPlatforms.enterprise,
          },
        },
      },
    });

    wrapper = mount(WizardStepPlatform, {
      global: {
        plugins: [testingPinia],
        stubs: ['Stack', 'WizardCard'],
      },
    });

    const script = wrapper.vm;
    const enterpriseCard = script.cardList.find((card: any) => card.key === targetPlatforms.enterprise);

    expect(enterpriseCard.isActive).toBe(true);
  });
});
