// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {architectures, distributionTargets, targetUsers} from '@disclosure-portal/model/NewWizard';
import {useWizardStore} from '@disclosure-portal/stores/wizard.store';
import {createTestingPinia} from '@pinia/testing';
import {mount} from '@vue/test-utils';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import WizardStepTargetUsers from '../WizardStepTargetUsers.vue';

describe('WizardStepTargetUsers', () => {
  let wrapper: any;
  let wizardStore: any;

  beforeEach(() => {
    wrapper = mount(WizardStepTargetUsers, {
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

  it('should have 3 cards in cardList', () => {
    const script = wrapper.vm;
    expect(script.cardList).toHaveLength(3);
  });

  it('should set targetUsers when card is selected', () => {
    const script = wrapper.vm;
    script.onCardSelect(targetUsers.company);

    expect(wizardStore.project.targetUsers).toBe(targetUsers.company);
  });

  it('should call nextStep when architecture is not frontend', () => {
    wizardStore.project.architecture = architectures.backend;

    const script = wrapper.vm;
    script.onCardSelect(targetUsers.company);

    expect(wizardStore.nextStep).toHaveBeenCalled();
    expect(wizardStore.nextTwoSteps).not.toHaveBeenCalled();
  });

  it('should set distributionTarget to company and call nextTwoSteps when company is selected and architecture is frontend', () => {
    wizardStore.project.architecture = architectures.frontend;

    const script = wrapper.vm;
    script.onCardSelect(targetUsers.company);

    expect(wizardStore.project.distributionTarget).toBe(distributionTargets.company);
    expect(wizardStore.nextTwoSteps).toHaveBeenCalled();
    expect(wizardStore.nextStep).not.toHaveBeenCalled();
  });

  it('should set distributionTarget to businessPartner and call nextTwoSteps when businessPartner is selected and architecture is frontend', () => {
    wizardStore.project.architecture = architectures.frontend;

    const script = wrapper.vm;
    script.onCardSelect(targetUsers.businessPartner);

    expect(wizardStore.project.distributionTarget).toBe(distributionTargets.businessPartner);
    expect(wizardStore.nextTwoSteps).toHaveBeenCalled();
    expect(wizardStore.nextStep).not.toHaveBeenCalled();
  });

  it('should set distributionTarget to businessPartner and call nextTwoSteps when customer is selected and architecture is frontend', () => {
    wizardStore.project.architecture = architectures.frontend;

    const script = wrapper.vm;
    script.onCardSelect(targetUsers.customer);

    expect(wizardStore.project.distributionTarget).toBe(distributionTargets.businessPartner);
    expect(wizardStore.nextTwoSteps).toHaveBeenCalled();
    expect(wizardStore.nextStep).not.toHaveBeenCalled();
  });

  it('should mark customer card as active when targetUsers is customer', () => {
    const testingPinia = createTestingPinia({
      createSpy: vi.fn,
      stubActions: true,
      initialState: {
        wizard: {
          project: {
            targetUsers: targetUsers.customer,
          },
        },
      },
    });

    wrapper = mount(WizardStepTargetUsers, {
      global: {
        plugins: [testingPinia],
        stubs: ['Stack', 'WizardCard'],
      },
    });

    const script = wrapper.vm;
    const customerCard = script.cardList.find((card: any) => card.key === targetUsers.customer);

    expect(customerCard.isActive).toBe(true);
  });
});
