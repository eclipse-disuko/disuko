// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, describe, expect, it} from 'vitest';
import ExpansionPanel from '../ExpansionPanel.vue';

const vExpansionPanelsStub = {
  template: '<div class="v-expansion-panels"><slot /></div>',
};
const vExpansionPanelStub = {
  template: '<div class="v-expansion-panel"><slot /></div>',
};
const vExpansionPanelTitleStub = {
  template: '<div class="v-expansion-panel-title"><slot :expanded="false" /></div>',
};
const vExpansionPanelTextStub = {
  template: '<div class="v-expansion-panel-text"><slot /></div>',
};
const vIconStub = {
  template: '<i class="v-icon" :data-icon="icon"></i>',
  props: ['icon', 'color'],
};

describe('ExpansionPanel', () => {
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

  const createWrapper = (props: Record<string, unknown>, slots = {}) => {
    return mount(ExpansionPanel, {
      props,
      slots,
      global: {
        stubs: {
          'v-expansion-panels': vExpansionPanelsStub,
          'v-expansion-panel': vExpansionPanelStub,
          'v-expansion-panel-title': vExpansionPanelTitleStub,
          'v-expansion-panel-text': vExpansionPanelTextStub,
          'v-icon': vIconStub,
        },
      },
    });
  };

  it('renders without throwing', () => {
    const wrapper = createWrapper({title: 'My Panel'});

    expect(wrapper.findComponent(vExpansionPanelsStub).exists()).toBe(true);
  });

  it('renders the given title', () => {
    const wrapper = createWrapper({title: 'My Panel'});

    expect(wrapper.text()).toContain('My Panel');
  });

  it('renders the optional text and the body slot inside the panel content', () => {
    const wrapper = createWrapper({title: 'My Panel', text: 'Some description'}, {body: '<span class="body">Body content</span>'});

    expect(wrapper.text()).toContain('Some description');
    expect(wrapper.find('.body').exists()).toBe(true);
  });

  it('does not render descriptive text when the text prop is absent', () => {
    const wrapper = createWrapper({title: 'My Panel'});

    expect(wrapper.find('.text-caption.d-block').exists()).toBe(false);
  });
});
