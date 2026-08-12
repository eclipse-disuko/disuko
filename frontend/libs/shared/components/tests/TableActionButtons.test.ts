// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, beforeEach, describe, expect, it, vi} from 'vitest';
import TableActionButtons, {TableActionButtonsProps} from '../TableActionButtons.vue';

const {toggleSlideMock, setupTableActionSliderMock} = vi.hoisted(() => ({
  toggleSlideMock: vi.fn(),
  setupTableActionSliderMock: vi.fn(),
}));

// The composable pulls in useTableActionSliderStore() (pinia) purely to share one timeout/width
// across sliders in the app - irrelevant to this component's branch logic, so mock it directly.
vi.mock('@shared/composables/useTableActionSlider', () => ({
  useTableActionSlider: () => ({
    sliderWidth: {value: 40},
    baseWidth: {value: 40},
    setupTableActionSlider: setupTableActionSliderMock,
    toggleSlide: toggleSlideMock,
  }),
}));

const DIconButtonStub = {
  props: ['icon', 'hint', 'color', 'disabled'],
  template:
    '<button type="button" class="d-icon-button" :data-icon="icon" :disabled="disabled" @click="$emit(\'clicked\')">{{ icon }}</button>',
};

const ExtraMenuStub = {
  template: '<div class="extra-menu"><slot /></div>',
};

describe('TableActionButtons', () => {
  const buttons = [
    {icon: 'mdi-pencil', event: 'edit'},
    {icon: 'mdi-delete', event: 'remove'},
    {icon: 'mdi-eye', event: 'view', show: false},
  ];

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

  beforeEach(() => {
    vi.clearAllMocks();
  });

  const createWrapper = (props: Partial<TableActionButtonsProps> = {}) => {
    return mount(TableActionButtons, {
      props: {buttons, ...props},
      global: {
        stubs: {DIconButton: DIconButtonStub, ExtraMenu: ExtraMenuStub},
      },
    });
  };

  it('renders a button for every shown entry in the normal variant', () => {
    const wrapper = createWrapper();

    const rendered = wrapper.findAllComponents(DIconButtonStub);
    expect(rendered).toHaveLength(2);
    expect(rendered.map((b) => b.props('icon'))).toEqual(['mdi-pencil', 'mdi-delete']);
  });

  it('emits the button-specific event when clicked in the normal variant', async () => {
    const wrapper = createWrapper();

    await wrapper.findAllComponents(DIconButtonStub)[0].trigger('click');

    expect(wrapper.emitted('edit')).toHaveLength(1);
  });

  it('puts all shown buttons inside an ExtraMenu in the minimal variant', () => {
    const wrapper = createWrapper({variant: 'minimal'});

    expect(wrapper.findComponent(ExtraMenuStub).exists()).toBe(true);
    expect(wrapper.findComponent(ExtraMenuStub).findAllComponents(DIconButtonStub)).toHaveLength(2);
  });

  it('shows shown buttons inline without a menu in the compact variant when there are 2 or fewer', () => {
    const wrapper = createWrapper({variant: 'compact', buttons: buttons.slice(0, 2)});

    expect(wrapper.findAllComponents(DIconButtonStub)).toHaveLength(2);
    expect(wrapper.findComponent(ExtraMenuStub).exists()).toBe(false);
  });

  it('shows the first shown button inline and the rest in an ExtraMenu in the compact variant when there are more than 2', () => {
    const manyButtons = [
      {icon: 'mdi-a', event: 'a'},
      {icon: 'mdi-b', event: 'b'},
      {icon: 'mdi-c', event: 'c'},
    ];
    const wrapper = createWrapper({variant: 'compact', buttons: manyButtons});

    const extraMenu = wrapper.findComponent(ExtraMenuStub);
    expect(extraMenu.exists()).toBe(true);
    expect(extraMenu.findAllComponents(DIconButtonStub)).toHaveLength(2);
    expect(extraMenu.findAllComponents(DIconButtonStub).map((b) => b.props('icon'))).toEqual(['mdi-b', 'mdi-c']);
  });

  it('sets up the slider and shows a toggle button when 2 or more buttons are shown in the slider variant', () => {
    const wrapper = createWrapper({variant: 'slider'});

    expect(setupTableActionSliderMock).toHaveBeenCalledWith(expect.any(Function), 2);
    const toggle = wrapper.findComponent(DIconButtonStub);
    expect(toggle.props('icon')).toBe('mdi-dots-horizontal');
  });

  it('toggles the slider when the dots button is clicked in the slider variant', async () => {
    const wrapper = createWrapper({variant: 'slider'});

    await wrapper.findComponent(DIconButtonStub).trigger('click');

    expect(toggleSlideMock).toHaveBeenCalledTimes(1);
  });

  it('shows a single inline button without the toggle when only one button is shown in the slider variant', () => {
    const wrapper = createWrapper({variant: 'slider', buttons: [buttons[0]]});

    const rendered = wrapper.findAllComponents(DIconButtonStub);
    expect(rendered).toHaveLength(1);
    expect(rendered[0].props('icon')).toBe('mdi-pencil');
  });

  it('emits the slideToggle event with the current slider width when toggled', () => {
    const wrapper = createWrapper({variant: 'slider'});

    expect(setupTableActionSliderMock).toHaveBeenCalled();
    const onToggle = setupTableActionSliderMock.mock.calls[0][0] as () => void;
    onToggle();

    expect(wrapper.emitted('slideToggle')?.[0]).toEqual([40]);
  });
});
