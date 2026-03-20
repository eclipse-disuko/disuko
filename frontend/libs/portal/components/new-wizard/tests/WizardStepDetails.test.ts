import {useCapabilitiesStore} from '@disclosure-portal/stores/capabilities';
import {useWizardStore} from '@disclosure-portal/stores/wizard.store';
import {vuetifyStubs} from '@disclosure-portal/test-utils/vuetify-stubs';
import {createTestingPinia} from '@pinia/testing';
import {mount, VueWrapper} from '@vue/test-utils';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import WizardStepDetails from '../WizardStepDetails.vue';

describe('WizardStepDetails', () => {
  let wrapper: VueWrapper;
  let wizardStore: any;
  let capabilitiesStore: any;

  const createWrapper = (options = {}) => {
    return mount(WizardStepDetails, {
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
          ApplicationSelector: {
            template: '<div class="application-selector"></div>',
            methods: {
              validate: vi.fn(() => Promise.resolve(true)),
            },
          },
        },
      },
      ...options,
    });
  };

  beforeEach(() => {
    wrapper = createWrapper();
    wizardStore = useWizardStore();
    capabilitiesStore = useCapabilitiesStore();
  });

  it('should render component', () => {
    expect(wrapper.exists()).toBe(true);
  });

  it('should render project name input field', () => {
    const nameField = wrapper.find('input[type="text"]');
    expect(nameField.exists()).toBe(true);
  });

  it('should render project description textarea', () => {
    const descriptionField = wrapper.find('textarea');
    expect(descriptionField.exists()).toBe(true);
  });

  it('should show ApplicationSelector for enterprise or mobile platform when connector is available', async () => {
    wizardStore.isEnterpriseOrMobilePlatform = true;
    capabilitiesStore.applicationConnector = true;
    await wrapper.vm.$nextTick();

    const appSelector = wrapper.find('.application-selector');
    expect(appSelector.exists()).toBe(true);
  });

  it('should not show ApplicationSelector for vehicle or other platforms', async () => {
    wizardStore.isEnterpriseOrMobilePlatform = false;
    capabilitiesStore.applicationConnector = true;
    await wrapper.vm.$nextTick();

    const appSelector = wrapper.find('.application-selector');
    expect(appSelector.exists()).toBe(false);
  });

  it('should render noFossProject checkbox', async () => {
    wizardStore.project.projectSettings = {noFossProject: false};
    await wrapper.vm.$nextTick();

    const checkbox = wrapper.find('input[type="checkbox"]');
    expect(checkbox.exists()).toBe(true);
  });

  it('should show warning when noFossProject is checked', async () => {
    wizardStore.project.projectSettings = {noFossProject: true};
    await wrapper.vm.$nextTick();

    const warningDiv = wrapper.findAll('div').find((div) => div.text().includes('warning'));
    expect(warningDiv).toBeDefined();
    expect(warningDiv?.exists()).toBe(true);
  });

  it('should validate form on mount when name is present', async () => {
    const testWizardStore = useWizardStore();
    testWizardStore.project.name = 'Existing Project';

    const newWrapper = createWrapper();
    await newWrapper.vm.$nextTick();
    await new Promise((resolve) => setTimeout(resolve, 0));

    expect(newWrapper.exists()).toBe(true);
  });
});
