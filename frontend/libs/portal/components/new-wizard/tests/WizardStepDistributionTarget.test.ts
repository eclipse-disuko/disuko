import {distributionTargets} from '@disclosure-portal/model/NewWizard';
import {useWizardStore} from '@disclosure-portal/stores/wizard.store';
import {createTestingPinia} from '@pinia/testing';
import {mount} from '@vue/test-utils';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import WizardStepDistributionTarget from '../WizardStepDistributionTarget.vue';

describe('WizardStepDistributionTarget', () => {
  let wrapper: any;
  let wizardStore: any;

  beforeEach(() => {
    wrapper = mount(WizardStepDistributionTarget, {
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

  it('should have 2 cards in cardList', () => {
    const script = wrapper.vm;
    expect(script.cardList).toHaveLength(2);
  });

  it('should set distributionTarget and call nextStep when card is selected', () => {
    const script = wrapper.vm;
    script.onCardSelect(distributionTargets.company);

    expect(wizardStore.project.distributionTarget).toBe(distributionTargets.company);
    expect(wizardStore.nextStep).toHaveBeenCalled();
  });

  it('should mark company card as active when distributionTarget is company', () => {
    const testingPinia = createTestingPinia({
      createSpy: vi.fn,
      stubActions: true,
      initialState: {
        wizard: {
          project: {
            distributionTarget: distributionTargets.company,
          },
        },
      },
    });

    wrapper = mount(WizardStepDistributionTarget, {
      global: {
        plugins: [testingPinia],
        stubs: ['Stack', 'WizardCard'],
      },
    });

    const script = wrapper.vm;
    const companyCard = script.cardList.find((card: any) => card.key === distributionTargets.company);

    expect(companyCard.isActive).toBe(true);
  });
});
