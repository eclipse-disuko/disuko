// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useIdleStore} from '@shared/stores/idle.store';
import {createTestingPinia} from '@pinia/testing';
import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, beforeEach, describe, expect, it, vi} from 'vitest';
import {nextTick} from 'vue';
import Idle from '../Idle.vue';

const {VProgressCircularStub, SelfBuildingSquareSpinnerStub} = vi.hoisted(() => ({
  VProgressCircularStub: {
    props: ['modelValue'],
    template: '<div class="v-progress-circular" :data-model-value="modelValue"></div>',
  },
  SelfBuildingSquareSpinnerStub: {template: '<div class="spinner-square" />'},
}));

// Idle.vue statically imports VProgressCircular/SelfBuildingSquareSpinner rather than using the
// auto-imported <v-progress-circular> tag resolution, so VTU's global.stubs (which only intercepts
// name-based component resolution) can't stub them - the real Vuetify component would otherwise
// pull in CSS imports that Vitest's environment can't load. Mock the modules directly instead.
vi.mock('vuetify/components', () => ({VProgressCircular: VProgressCircularStub}));
vi.mock('epic-spinners', () => ({SelfBuildingSquareSpinner: SelfBuildingSquareSpinnerStub}));

describe('Idle', () => {
  const originalTMock = config.global.mocks?.$t;

  beforeAll(() => {
    if (config.global.mocks && '$t' in config.global.mocks) {
      delete config.global.mocks.$t;
    }
  });

  afterAll(() => {
    if (originalTMock) {
      config.global.mocks = {
        ...config.global.mocks,
        $t: originalTMock,
      };
    }
  });

  const createWrapper = () => {
    return mount(Idle, {
      global: {
        plugins: [createTestingPinia({createSpy: vi.fn, stubActions: false})],
        stubs: {
          'v-layout': {template: '<div class="v-layout"><slot /></div>'},
          Stack: {template: '<div class="stack"><slot /></div>'},
        },
      },
    });
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders nothing when the idle store is hidden', () => {
    const wrapper = createWrapper();

    expect(wrapper.find('[data-testid="idle"]').exists()).toBe(false);
  });

  it('renders the overlay when showIdle is true', async () => {
    const wrapper = createWrapper();
    const idle = useIdleStore();
    idle.showIdle = true;
    await nextTick();

    expect(wrapper.find('[data-testid="idle"]').exists()).toBe(true);
  });

  it('shows the indeterminate spinner when progress is -1', async () => {
    const wrapper = createWrapper();
    const idle = useIdleStore();
    idle.showIdle = true;
    idle.progress = -1;
    await nextTick();

    expect(wrapper.find('.spinner-square').exists()).toBe(true);
    expect(wrapper.findComponent(VProgressCircularStub).exists()).toBe(false);
  });

  it('shows the determinate progress circle with the progress value once known', async () => {
    const wrapper = createWrapper();
    const idle = useIdleStore();
    idle.showIdle = true;
    idle.progress = 42;
    await nextTick();

    expect(wrapper.find('.spinner-square').exists()).toBe(false);
    const progress = wrapper.findComponent(VProgressCircularStub);
    expect(progress.exists()).toBe(true);
    expect(progress.props('modelValue')).toBe(42);
  });

  it('renders the idle message and progress unit', async () => {
    const wrapper = createWrapper();
    const idle = useIdleStore();
    idle.showIdle = true;
    idle.idleMessage = 'Loading data...';
    idle.progress = 50;
    idle.progressUnit = '%';
    await nextTick();

    expect(wrapper.text()).toContain('Loading data...');
    expect(wrapper.text()).toContain('50%');
  });

  it('does not render a progress suffix span when progress is unknown', async () => {
    const wrapper = createWrapper();
    const idle = useIdleStore();
    idle.showIdle = true;
    idle.idleMessage = 'Loading data...';
    idle.progress = -1;
    await nextTick();

    expect(wrapper.text()).toContain('Loading data...');
    expect(wrapper.text()).not.toContain('-1');
  });
});
